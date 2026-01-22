# PaperEngine 批量撮合与事务写入解读

目标：理解 PaperEngine 如何用 “行情驱动撮合 + 批量 flush + DB 事务” 保证一致性并提升吞吐。

源码入口： [paper/engine.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/paper/engine.go)

## 1) 启动时加载未完成订单 + 订阅行情

```go
func (e *PaperEngine) Start(ctx context.Context) error {
  // 1) 重启恢复：把 DB 里 status=open 的订单重新装入内存
  if err := e.loadOpenOrders(ctx); err != nil { return err }

  // 2) 订阅 1m K线，作为撮合触发器（不是订阅 raw）
  _, err := e.js.Subscribe("market.kline.1m.*", func(msg *nats.Msg) {
    var candle model.KLine
    if err := json.Unmarshal(msg.Data, &candle); err != nil { return }
    e.processPriceUpdate(candle)
  })
  if err != nil { return err }

  // 3) flush loop：把已撮合的订单批量写入 DB
  go e.batchFlushLoop(ctx)
  return nil
}
```

关键点：

- 只订阅 1m K 线：用“聚合后的行情”降低撮合触发频率（吞吐更友好）
- 重启恢复：从 DB 把 open orders 重新装回内存，避免状态丢失

## 2) 行情驱动撮合：processPriceUpdate

```go
// 对每个 symbol 的 open orders 扫一遍，用 candle high/low 判断是否成交
switch o.Type {
case "market":
  filled = true
  o.FilledPrice = candle.Close
case "limit":
  if buy && candle.Low <= limitPrice  => filled at limitPrice
  if sell && candle.High >= limitPrice => filled at limitPrice
}
```

关键点：

- 这是“模拟撮合”，成交模型依赖 K 线 high/low，因此无法提供 tick 级精确成交，但足够用于策略评估/回测演示
- `fillChan` 把“已成交订单”从撮合阶段投递到 flush 阶段，实现生产者/消费者拆分

## 3) 批量 flush：batchFlushLoop

```go
// 两个触发条件：满 50 单立即 flush；或 500ms 到点 flush
case o := <-e.fillChan:
  batch = append(batch, o)
  if len(batch) >= 50 { flushBatch(batch); batch=nil }
case <-ticker.C:
  if len(batch) > 0 { flushBatch(batch); batch=nil }
```

关键点：

- 通过 batch 提升 DB 写入吞吐并降低事务开销
- 这里属于“吞吐优先”，如果需要更低延迟，可以降低 ticker 或 batch size

## 4) flushBatch：用事务保证一致性

```go
tx, _ := e.db.Begin(ctx)
defer tx.Rollback(ctx)

for _, o := range batch {
  // 1) 更新订单状态（filled + filled_price）
  tx.Exec("UPDATE paper_orders ...", o.FilledPrice, o.ID)

  // 2) 更新账户余额（买扣钱、卖加钱）
  tx.Exec("UPDATE paper_accounts SET balance = balance +/- $1", amount)

  // 3) 更新持仓（买 upsert 平均成本；卖则减仓）
  tx.Exec("INSERT ... ON CONFLICT DO UPDATE ...")
}

tx.Commit(ctx)
```

关键点：

- 三类更新在同一事务里，避免“订单已成交但余额/持仓未同步”的不一致
- 异常处理当前偏宽松（部分 Exec 失败 continue），如果你要更强一致，应在任一失败时回滚并重试

## 5) 可观察性建议（你可以加的指标/日志）

- 撮合命中率：每条 kline 触发的 toFill 数量
- flush 事务耗时直方图（p99）
- 回放恢复耗时（loadOpenOrders 行数与耗时）

