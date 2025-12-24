# Market Ingestor - 高性能行情摄取与处理系统

`quant-trader` 是量化交易系统的核心基础设施，专为**高并发、低延迟、高可用**场景设计。它集成了多交易所数据摄取、流式 K 线聚合、异步批量持久化及实时行情推送四大核心功能。

---

## 1. 深度系统架构

本系统基于 **Golang 并发模型** 与 **NATS JetStream 事件驱动架构** 构建，通过解耦生产者与消费者，实现了极高的吞吐能力。

### 1.1 数据流转链路

1. **摄取层 (Ingestion)**: 维护与交易所（Binance, OKX, Bybit, Coinbase, Kraken）的 WebSocket 长连接，将原始消息转换为归一化 `model.Trade`。
2. **消息总线 (Messaging)**: 原始成交通过 NATS JetStream 发布到 `market.raw.{exchange}.{symbol}` 主题。
3. **聚合层 (Processing)**: `KlineProcessor` 订阅原始成交流，在内存中维护分钟级窗口状态，生成 K 线数据。
4. **持久化层 (Storage)**: `BatchSaver` 订阅 `market.raw` 和 `market.kline` 流，利用批量插入技术高效写入 TimescaleDB。
5. **推送层 (Push)**: `PushGateway` 订阅聚合后的 K 线流，通过 WebSocket 广播给终端用户。

### 1.2 并发拓扑与组件通信

系统利用 Goroutine 和 Channel 构建了高度并行的执行环境：

- 每个交易所连接、每个处理逻辑、每个客户端发送都运行在独立的 Goroutine 中。
- **通信模式**：`Connector -> tradeChan (Go Channel) -> NATS JetStream -> Subscriptions`。

---

## 2. 核心算法与设计细节

### 2.1 高可靠摄取算法 (Connector Implementation)

**设计目标**：在不稳定的网络环境下确保数据不丢失。

- **双重超时控制**：
  - `ReadDeadline`: 设定为 60s。若 60s 内未收到任何数据或消息，主动断开重连。
  - `WriteDeadline`: 确保心跳包发送不阻塞主逻辑。
- **指数退避重连逻辑 (Exponential Backoff)**：
    当检测到网络中断或握手失败时，系统会自动重试，重试间隔随失败次数指数增长（1s, 2s, 4s... 最高 60s），防止对交易所 API 造成过大压力。
- **无锁解析**：利用 `json.RawMessage` 延迟解析非关键字段，提升 CPU 利用率。

### 2.2 数据归一化 (Symbol Normalization)

各交易所对交易对的命名格式各异（如 `BTC-USDT`, `btcusdt`, `XBT/USD`）。系统在进入 NATS 消息总线前，会通过 `NormalizeSymbol` 函数统一转换为系统标准格式（如 `BTCUSDT`），确保下游所有聚合和存储逻辑的全局一致性。

### 2.3 流式分钟 K 线聚合 (Kline Processor)

**设计目标**：毫秒级生成 K 线，支持 O(1) 查询效率。

- **内存窗口状态机**：
  - 使用 `map[string]*model.KLine` 存储当前活跃窗口。
  - **聚合逻辑**：
    - `High = max(High, newPrice)`
    - `Low = min(Low, newPrice)`
    - `Close = newPrice`
    - `Volume += newAmount`
- **滑动式状态刷新 (Sliding Flush)**：
    聚合器不依赖“下一分钟”的数据到达来触发刷新。`flushLoop` 每 5 秒扫描一次内存 Map。如果 `candle.Timestamp.Unix() < currentMinute.Unix()`，则认为该分钟已绝对封闭，触发 NATS 发布并从内存安全移除。

### 2.4 异步批量持久化 (Batch Persistence)

**设计目标**：解决每秒万级写入导致的磁盘 I/O 瓶颈。

- **双触发刷新策略**：
  - 缓冲区 `buffer []model.Trade` 达到 1000 条或 1 秒计时器到期。
- **TimescaleDB 超表 (Hypertable) 优化**：
  - 数据按 `time` 字段自动分区（Chunking），提升过期数据清理和新数据写入效率。
  - **索引策略**：建立 `(symbol, time DESC)` 复合索引，优化最近 N 条数据的查询速度。
- **SQL 性能**：使用 `pgx.Batch` 减少网络往返（RTT），吞吐量较单条 `INSERT` 提升 10-20 倍。

### 2.5 非阻塞广播与背压控制 (Push Gateway)

**设计目标**：支持 10k+ 客户端并发订阅，防止“慢客户端”拖慢全系统。

- **客户端隔离 (Isolation)**：
    每个客户端拥有一个独立的 `send chan []byte` (容量 256)。
- **丢弃策略 (Drop Policy)**：
    当 `send` channel 满时，执行 `default: break`（即丢弃最新行情）。
    **理由**：行情数据具有时效性，丢弃旧数据比阻塞整个推送网关更优。

---

## 3. 数据模型与精度保证

### 3.1 金额精度处理

系统**全线严禁使用 `float64`** 处理任何价格和数量。

- **Go 层**：使用 `shopspring/decimal` 库进行任意精度运算。
- **DB 层**：使用 PostgreSQL 的 `NUMERIC(20, 8)` 类型。
- **理由**：避免浮点数舍入误差带来的财务对账问题。

### 3.2 数据库 Schema 核心字段

- `market_trades`: 记录每一笔成交，包含 `trade_id`（交易所原始 ID，用于幂等去重）。
- `market_klines`: 存储聚合后的 K 线，包含 `period` 字段（支持 1m, 5m, 1h 等多周期扩展）。

---

## 4. 系统稳定性与运维

### 4.1 优雅停机 (Graceful Shutdown)

系统监听 `SIGTERM` 和 `SIGINT` 信号。触发后：

1. 停止所有 Ingestor Connector，不再接收新行情。
2. 触发所有 `BatchSaver` 的强制 `Flush`，确保内存中的数据全部落库。
3. 等待 NATS 消费者完成当前消息的 Ack。
4. 最后关闭数据库连接池。

### 4.2 监控指标 (Prometheus Metrics)

系统通过 `/metrics` 接口暴露以下关键指标：

- `ingest_latency_seconds`: 链路全过程延迟。
- `db_insert_total`: 成功写入数据库的记录总数。
- `ws_connections_total`: 当前活跃的 WebSocket 客户端连接数。
- `trade_process_total`: 系统每秒处理的成交笔数。

---

## 5. 性能调优参数

| 参数名 | 默认值 | 描述 | 调优建议 |
| :--- | :--- | :--- | :--- |
| `BATCH_SIZE` | 1000 | 持久化批量大小 | 写入压力大时可调至 5000 |
| `FLUSH_INTERVAL` | 1s | 持久化强制刷新间隔 | 对实时性要求极高时可调至 500ms |
| `WS_SEND_BUFFER` | 256 | WebSocket 发送缓冲区 | 客户端连接多且网络差时需谨慎调大 |

---

## 6. 开发者指南

### 如何添加新交易所？

1. 在 `internal/connector/` 下创建 `new_exchange.go`。
2. 实现 `Run(ctx context.Context, tradeChan chan<- model.Trade)` 方法。
3. 在 `internal/app/worker.go` 中注册。

### 如何运行测试？

```bash
# 运行所有测试
go test ./internal/...

# 运行压力测试
go test -bench=. ./internal/processor/...
```
