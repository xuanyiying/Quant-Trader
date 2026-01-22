# K线聚合（KlineProcessor）逐行解读

目标：理解 `market.raw` 如何聚合为多周期 `market.kline`，以及并发/锁/批量 flush 的取舍。

源码入口： [kline.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/processor/kline.go)

## 1) 数据结构与并发模型

```go
type KlineProcessor struct {
  js      nats.JetStreamContext
  logger  *zap.Logger

  candles map[string]*model.KLine // 内存桶：key = exchange:symbol:period:window
  mu      sync.Mutex              // 保护 candles（多 worker 并发更新）

  jobs       chan model.Trade     // 有界队列：避免回调里做重活
  numWorkers int                  // worker 数量（默认 4）
}
```

关键点：

- `Subscribe` 回调只做轻量解析 + 投递 `jobs`，避免 JetStream 回调阻塞导致 ack 延迟
- `candles` 用一把 `mu` 保护：实现简单，但高吞吐时可能成为热点锁（优化方向：按 symbol 分片锁或单写者 per-symbol）

## 2) 订阅 raw 并投递到 worker 池

```go
_, err := p.js.Subscribe("market.raw.*.*", func(msg *nats.Msg) {
  var trade model.Trade
  if err := json.Unmarshal(msg.Data, &trade); err != nil {
    p.logger.Error("failed to unmarshal trade in processor", zap.Error(err))
    return
  }

  // 核心思想：把重计算从回调移走
  select {
  case p.jobs <- trade:
  default:
    // 队列满时丢弃，属于“保护系统稳定性”的策略（需要业务确认是否允许丢）
    p.logger.Warn("processor job queue full, trade dropped", zap.String("symbol", trade.Symbol))
  }

  msg.Ack() // manual ack：处理成功后确认
}, nats.Durable("kline-processor"), nats.ManualAck())
```

关键点：

- durable + manual ack：确保处理可恢复，但要求后续系统能处理重复消息
- 这里选择 “队列满则丢” 是一种背压策略：更偏向系统稳定性而非“完全不丢”

## 3) processTrade：同一条 trade 更新多个周期桶

```go
for _, period := range model.SupportedPeriods {
  duration := model.PeriodToDuration(period)
  window := trade.Timestamp.Truncate(duration) // 对齐到窗口开始时间
  key := fmt.Sprintf("%s:%s:%s:%s",
    trade.Exchange, trade.Symbol, period, window.Format(time.RFC3339))

  candle, ok := p.candles[key]
  if !ok {
    // 新窗口：Open/High/Low/Close 都是第一笔成交价
    candle = &model.KLine{...}
    p.candles[key] = candle
  } else {
    // 追加：更新高低收与成交量
    if trade.Price.GreaterThan(candle.High) { candle.High = trade.Price }
    if trade.Price.LessThan(candle.Low) { candle.Low = trade.Price }
    candle.Close = trade.Price
    candle.Volume = candle.Volume.Add(trade.Amount)
  }
}
```

关键点：

- 同一笔 trade 会同时更新 1m/5m/15m/… 多个周期，因此 flush 频率与 CPU 消耗成正比
- key 带 exchange：可允许不同交易所的同 symbol 分开聚合（或后续再做合成）

## 4) flushLoop：定时把“已结束窗口”发布出去

```go
ticker := time.NewTicker(1 * time.Second)
for {
  select {
  case <-ctx.Done():
    return
  case <-ticker.C:
    p.flush()
  }
}
```

flush 的核心：遍历 `candles`，判断 `now.After(candle.Timestamp.Add(duration))` 才 flush，确保窗口结束后才发布。

## 5) 发布到 NATS：market.kline.period.symbol

```go
subject := fmt.Sprintf("market.kline.%s.%s", candle.Period, candle.Symbol)
data, _ := json.Marshal(candle)
_, err := p.js.Publish(subject, data)
```

关键点：

- 发布失败只记录日志；如果需要更强一致，可考虑重试队列或落本地 WAL
- 下游（storage/alert/paper/strategy/push）都只需要订阅 kline，即可复用这条“标准化行情管道”

