# NATS Subjects / JetStream Streams 设计与使用速查

目标：把“谁 publish 到哪里、谁 subscribe 什么、durable 名称是什么、是否 manual ack”这件事一次性梳清楚。

相关源码：

- NATS/JetStream 初始化： [nats.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/infrastructure/nats.go)
- 行情接入 publish： [worker.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/worker.go)
- K线处理 subscribe/publish： [kline.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/processor/kline.go)
- 落库订阅： [worker.go:startPersistenceService](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/worker.go#L92-L119)
- 策略信号： [strategy_runner.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/engine/strategy_runner.go)
- 撮合 WAL/快照/租约： [wal_jetstream.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/wal_jetstream.go)

## 1) 主题（Subjects）清单（按域）

### Market 域

- 原始成交（producer：connectors / consumer：processor/storage 等）
  - `market.raw.{exchange}.{symbol}`
  - 例：`market.raw.binance.BTCUSDT`
- K线（producer：KlineProcessor / consumer：storage/alert/paper/strategy/push）
  - `market.kline.{period}.{symbol}`
  - 例：`market.kline.1m.BTCUSDT`

### Strategy 域

- 策略信号（producer：StrategyRunner / consumer：push/UI）
  - `strategy.signal.{strategyName}.{symbol}`

### Notification 域

- 用户通知（producer：AlertService / consumer：push/UI）
  - `notification.user.{userID}`

### Matching 域（撮合）

- 撮合事件 WAL（producer：matching.Engine / consumer：恢复/审计/下游服务）
  - `matching.event.{symbol}`

## 2) Streams / KV Buckets（JetStream）

### MARKET stream

由 [InitNATS](file:///Users/yiying/GoProjects/quant-trader/backend/internal/infrastructure/nats.go) 确保存在：

- subjects 覆盖：`market.raw.*.*` 与 `market.kline.*.*`
- 作用：让 raw/kline 具备持久化与 durable consumer 语义

### MATCHING stream

由 [ensureMatchingEventStream](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/wal_jetstream.go#L21-L50) 确保存在：

- subjects 覆盖：`matching.event.*`
- 作用：作为撮合 WAL（Write-Ahead Log），支持恢复回放

### MATCHING_SNAPSHOT（KV bucket）

由 [NewJetStreamSnapshotStore](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/wal_jetstream.go#L81-L127) 创建/获取：

- key：symbol
- value：`EngineSnapshot` JSON
- 作用：缩短恢复时间（不用从 0 回放全部 WAL）

### MATCHING_LEASE（KV bucket）

由 [NewJetStreamLease](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/wal_jetstream.go#L129-L168) 创建/获取：

- key：symbol
- value：owner（实例标识）
- TTL：租约过期控制
- 作用：多实例下保证同一 symbol 仅一个 leader 运行撮合并写入 WAL/快照

## 3) Durable / ManualAck 使用规则（为什么这么做）

### Durable 的意义

以 `KlineProcessor` 为例（[kline.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/processor/kline.go#L41-L57)）：

```go
p.js.Subscribe("market.raw.*.*", handler,
  nats.Durable("kline-processor"),
  nats.ManualAck(),
)
```

- durable 名称 = “消费位点”，进程重启后会从上次 ack 的位置继续
- manual ack = 只有处理成功才 ack（否则允许重投）

### 至少一次投递的代价：下游要幂等

因此落库侧通常需要：

- upsert / conflict do nothing
- 或引入消息 id 做去重表

（本项目的 `storage/*` 以 batch + upsert/冲突处理为主。）

## 4) 前端如何对齐 subjects（WS 订阅）

前端通过 WS 发送：

- `market.kline.{period}.{symbol}`
- `market.raw.*.{symbol}`
- `strategy.signal.*.{symbol}`

对应实现见：

- [useWebSocket.ts](file:///Users/yiying/GoProjects/quant-trader/frontend/src/hooks/useWebSocket.ts)
- [PushGateway](file:///Users/yiying/GoProjects/quant-trader/backend/internal/push/gateway.go)

