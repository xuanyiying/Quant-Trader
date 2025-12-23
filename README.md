# Market Ingestor - High-Performance Market Data Ingestion & Processing System

`market-ingestor` is the core infrastructure of the Quant Trader system, designed for **high concurrency, low latency, and high availability**. It integrates multi-exchange data ingestion, streaming K-line aggregation, asynchronous batch persistence, and real-time market data push.

---

## 1. System Architecture

The system is built on the **Golang Concurrency Model** and **NATS JetStream Event-Driven Architecture**, achieving high throughput by decoupling producers and consumers.

### 1.1 Data Flow
1.  **Ingestion Layer**: Maintains WebSocket connections with exchanges (Binance, OKX, Bybit, Coinbase, Kraken), converting raw messages into normalized `model.Trade`.
2.  **Messaging Layer**: Raw trades are published to NATS JetStream topics `market.raw.{exchange}.{symbol}`.
3.  **Processing Layer**: `KlineProcessor` subscribes to raw trade streams, maintains minute-level window states in memory, and generates K-line data.
4.  **Storage Layer**: `BatchSaver` subscribes to `market.raw` and `market.kline` streams, efficiently writing to TimescaleDB using batch insertion.
5.  **Push Layer**: `PushGateway` subscribes to aggregated K-line streams and broadcasts them to end-users via WebSocket.

### 1.2 Concurrency Topology
The system utilizes Goroutines and Channels for highly parallel execution:
- Each exchange connection, processing logic, and client transmission runs in an independent Goroutine.
- **Communication Pattern**: `Connector -> tradeChan (Go Channel) -> NATS JetStream -> Subscriptions`.

---

## 2. Core Algorithms & Design Details

### 2.1 High-Reliability Ingestion (Connector)
**Goal**: Ensure no data loss in unstable network environments.
-   **Dual Timeout Control**:
    -   `ReadDeadline`: Set to 60s. Active reconnection if no data is received.
    -   `WriteDeadline`: Ensures heartbeat packets don't block the main logic.
-   **Exponential Backoff Reconnection**:
    Automatic retries with increasing intervals (1s, 2s, 4s... up to 60s) to avoid overwhelming exchange APIs.
-   **Lock-Free Parsing**: Uses `json.RawMessage` to defer parsing non-critical fields, improving CPU utilization.

### 2.2 Symbol Normalization
Exchanges use different naming formats (e.g., `BTC-USDT`, `btcusdt`, `XBT/USD`). The system normalizes these to a standard format (e.g., `BTCUSDT`) before publishing to NATS.

### 2.3 Streaming K-Line Aggregation (Kline Processor)
**Goal**: Millisecond-level K-line generation with O(1) efficiency.
-   **In-Memory Window State Machine**:
    - Uses `map[string]*model.KLine` to store active windows.
    - **Aggregation Logic**:
        - `High = max(High, newPrice)`
        - `Low = min(Low, newPrice)`
        - `Close = newPrice`
        - `Volume += newAmount`
-   **Sliding Flush**:
    Flushes every 5 seconds. If `candle.Timestamp.Unix() < currentMinute.Unix()`, the minute is considered closed, published to NATS, and removed from memory.

### 2.4 Asynchronous Batch Persistence (Batch Persistence)
**Goal**: Solve I/O bottlenecks for high-frequency writes.
-   **Dual-Trigger Flush**:
    - Buffer `buffer []model.Trade` reaches 1000 items OR 1-second timer expires.
-   **TimescaleDB Hypertable Optimization**:
    - Automatic partitioning by `time` field (Chunking).
    - **Indexing Strategy**: Composite index on `(symbol, time DESC)` for optimized queries.
-   **SQL Performance**: Uses `pgx.Batch` to reduce network round-trips (RTT), increasing throughput by 10-20x.

### 2.5 Non-Blocking Broadcast (Push Gateway)
**Goal**: Support 10k+ concurrent client subscriptions without "slow client" issues.
-   **Client Isolation**: Each client has an independent `send chan []byte` (capacity 256).
-   **Drop Policy**: If the `send` channel is full, the latest data is dropped (`default: break`). Data timeliness is prioritized over blocking the system.

---

## 3. Data Model & Precision

### 3.1 Financial Precision
**No `float64` is used** for prices or amounts.
-   **Go Layer**: `shopspring/decimal` library for arbitrary-precision arithmetic.
-   **DB Layer**: PostgreSQL `NUMERIC(20, 8)` type.

### 3.2 Database Schema
-   `market_trades`: Records every trade with `trade_id` for idempotency.
-   `market_klines`: Stores aggregated K-lines with `period` (1m, 5m, 1h, etc.).

---

## 4. Stability & Operations

### 4.1 Graceful Shutdown
Listens for `SIGTERM` and `SIGINT`:
1. Stop all Ingestor Connectors.
2. Force `Flush` all `BatchSaver` buffers.
3. Wait for NATS consumer Acks.
4. Close database connection pools.

### 4.2 Prometheus Metrics
Exposed via `/metrics`:
-   `ingest_latency_seconds`: End-to-end latency.
-   `db_insert_total`: Successfully written records.
-   `ws_connections_total`: Active WebSocket clients.
-   `trade_process_total`: Trades processed per second.

---

## 5. Performance Tuning

| Parameter | Default | Description | Tuning Suggestion |
| :--- | :--- | :--- | :--- |
| `BATCH_SIZE` | 1000 | Persistence batch size | Increase to 5000 for high load |
| `FLUSH_INTERVAL` | 1s | Persistence flush interval | Decrease to 500ms for high real-time needs |
| `WS_SEND_BUFFER` | 256 | WebSocket send buffer | Adjust carefully based on client count |

---

## 6. Developer Guide

### How to add a new exchange?
1. Create `new_exchange.go` in `internal/connector/`.
2. Implement `Run(ctx context.Context, tradeChan chan<- model.Trade)`.
3. Register it in `internal/app/worker.go`.

### How to run tests?
```bash
go test ./internal/...
go test -bench=. ./internal/processor/...
```
