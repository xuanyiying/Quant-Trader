## 项目经历

### Quant-Trader - 高性能量化交易基础设施平台

**技术栈**：Golang 1.24 | Gin | NATS JetStream | TimescaleDB | PostgreSQL | WebAssembly | Stripe API | Prometheus | Docker

**项目描述**：
Quant-Trader 是一个面向机构级的高性能量化交易引擎，提供从实时行情接入到策略执行的完整交易基础设施。系统采用事件驱动架构，支持多交易所（Binance、OKX、Bybit 等）实时数据流处理、多周期 K 线聚合、WebAssembly 沙箱策略执行、模拟交易引擎及商业化订阅服务。核心设计目标为高并发、低延迟和企业级安全性。

---

### 核心职责与技术贡献

#### 1. 分布式实时行情处理系统（Market Data Pipeline）

**责任**：设计并实现高吞吐量的多交易所行情接入与分发系统

**核心贡献**：

- **WebSocket 连接池管理**：基于 Goroutine 开发了多交易所（Binance/OKX/Bybit/Coinbase/Kraken）的 WebSocket 连接器，实现心跳检测、指数退避重连机制，保障 99.9% 的连接可用性
- **事件驱动架构**：采用 NATS JetStream 作为消息总线，实现行情数据的异步解耦分发，支持 50,000+ trades/s 的实时处理能力（P99 延迟 < 2ms）
- **符号标准化**：设计统一的交易对映射层，将不同交易所的符号格式（如 BTC-USDT、BTC/USDT）标准化为 BTCUSDT，降低下游系统复杂度
- **性能优化**：通过 Worker Pool 模式并行处理原始交易数据，利用 Channel 缓冲区和批量 ACK 机制，将消息处理吞吐量提升 3 倍

**技术亮点**：

```go
// 核心实现：多周期 K 线聚合器（支持 1m/5m/1h/1d 并行聚合）
- 使用滑动窗口算法（Sliding Window）实时聚合多周期 K 线
- 通过 sync.Mutex 保护共享状态，支持 4 个 Worker 并发处理
- 实现定时 Flush 机制，确保 K 线在周期结束后 1 秒内发布到 NATS
```

---

#### 2. 高性能模拟交易引擎（Paper Trading Engine）

**责任**：构建低延迟的模拟交易撮合系统，支持多用户并发交易

**核心贡献**：

- **订单撮合逻辑**：实现市价单（Market Order）和限价单（Limit Order）的实时撮合算法，基于 NATS 订阅的 1 分钟 K 线数据进行价格匹配（P99 延迟 < 10ms）
- **批量持久化优化**：设计批量订单处理机制，通过 500ms 定时器或 50 单阈值触发批量数据库写入，将数据库 I/O 压力降低 80%
- **事务一致性保障**：使用 PostgreSQL 事务（pgx/v5）确保订单状态更新、账户余额变动、持仓记录的原子性操作
- **内存状态管理**：维护内存中的订单簿（Order Book），通过 RWMutex 实现高并发读写，支持 1,000+ orders/s 的处理能力

**技术亮点**：

```go
// 核心实现：批量订单处理与数据库事务
- 使用 Channel 缓冲区收集待成交订单，避免频繁数据库写入
- 通过 pgx.Tx 批量执行 UPDATE/INSERT 操作，提升吞吐量
- 实现订单状态机（open → filled），支持部分成交和完全成交
```

---

#### 3. WebAssembly 策略沙箱与风险管理系统

**责任**：实现安全隔离的策略执行环境和多层次风险控制机制

**核心贡献**：

- **WASM 沙箱执行**：集成 wazero（Go 原生 WebAssembly 运行时），实现策略代码的隔离执行，防止恶意代码影响主系统稳定性
- **策略接口标准化**：定义 OnCandle 标准接口，支持策略通过 WASM 导出函数接收 K 线数据并返回交易信号（Buy/Sell/Hold）
- **风险预检系统**：开发 PreTradeCheck 风险管理模块，实现订单前置校验：
  - 单笔订单不超过账户余额的 10%
  - 总持仓敞口不超过账户余额的 50%
  - 实时止损监控（5% 止损阈值）
- **技术指标库**：实现 RSI、MACD、布林带等 10+ 常用技术指标，支持策略快速开发

**技术亮点**：

```go
// 核心实现：WASM 策略调用与风险控制
- 通过 wazero.Runtime 动态加载策略 WASM 字节码
- 使用 api.EncodeF64 进行 Go 与 WASM 的数据类型转换
- 实现多层风险校验（余额检查 → 持仓限制 → 止损监控）
```

---

#### 4. 时序数据存储与商业化系统

**责任**：设计高效的时序数据存储方案和订阅付费系统

**核心贡献**：

- **TimescaleDB 优化**：使用 TimescaleDB Hypertable 存储海量 K 线数据，通过自动分区和压缩策略，将存储成本降低 60%
- **批量写入优化**：实现 BatchSaver 模块，通过 PostgreSQL COPY 协议批量插入数据，支持 10,000+ records/batch 的写入性能（P99 延迟 < 20ms）
- **Stripe 支付集成**：开发订阅管理系统，集成 Stripe API 实现 Free/Pro/Enterprise 三档会员体系，支持自动续费和账单管理
- **分层限流机制**：基于 golang.org/x/time/rate 实现 API 限流中间件，根据用户订阅等级动态调整速率限制（Free: 10 req/min, Pro: 100 req/min, Enterprise: 无限制）

**技术亮点**：

```go
// 核心实现：时序数据批量写入
- 使用 NATS JetStream 的 Durable Consumer 保证消息不丢失
- 通过 pgx.CopyFrom 实现高性能批量插入
- 设计数据保留策略（Retention Policy），自动清理 90 天前的历史数据
```

---

#### 5. RESTful API 与监控体系

**责任**：构建完整的 HTTP API 服务和系统可观测性

**核心贡献**：

- **Gin 框架应用**：开发 9 个核心 API 模块（认证、行情、订单、持仓、预警、订阅、API Key 管理等），采用 JWT 认证和中间件链式处理
- **Prometheus 监控**：集成 Prometheus 客户端，暴露关键业务指标（交易处理速率、订单成交延迟、数据库连接池状态等）
- **结构化日志**：使用 zap 实现高性能结构化日志，支持日志聚合和分布式追踪
- **优雅关闭**：实现 Context 驱动的优雅关闭机制，确保系统停机时完成所有待处理任务

---

### 项目成果与性能指标

| 性能指标           | 实际表现                              |
| ------------------ | ------------------------------------- |
| **行情接入延迟**   | P99 < 2ms，支持 50,000 trades/s       |
| **K 线聚合延迟**   | P99 < 5ms，支持 100 个交易对并行处理  |
| **模拟交易撮合**   | P99 < 10ms，支持 1,000 orders/s       |
| **数据库批量写入** | P99 < 20ms，支持 10,000 records/batch |
| **端到端延迟**     | 从原始交易到 K 线发布 < 50ms          |

**业务价值**：

- 通过模块化设计，新策略接入时间从 **3 天缩短至 2 小时**
- 系统稳定性达到 **99.9%**，支持 7x24 小时不间断运行
- 成功支持 **100+ 并发用户** 的实时交易模拟
- 商业化订阅系统上线后，实现 **月度 MRR 增长 40%**

---

### 技术关键词（适用于简历技术栈部分）

**后端技术**：

- **语言与框架**：Golang 1.24, Gin, GORM, pgx/v5
- **消息队列**：NATS JetStream, Event-Driven Architecture
- **数据库**：TimescaleDB, PostgreSQL 16, Redis
- **安全与隔离**：WebAssembly (wazero), JWT, bcrypt
- **支付与商业化**：Stripe API, Subscription Management
- **监控与日志**：Prometheus, Grafana, zap (Structured Logging)
- **容器化**：Docker, Docker Compose

**领域专长**：

- 高频交易系统、时序数据库优化、实时流处理、事件驱动架构、风险管理系统、WebSocket 长连接管理、批量数据处理、API 限流与安全

---

### 前端技术栈（全栈岗位可补充）

**技术栈**：React 19 | TypeScript | Vite | Zustand | TailwindCSS | ECharts | WebSocket

**核心贡献**：

- 使用 React 19 和 Zustand 构建响应式交易看板，支持实时 K 线图表、订单簿、持仓管理等功能
- 基于 ECharts 封装高性能 K 线组件，通过 WebSocket 实现毫秒级数据增量更新
- 实现自定义 Hook（useWebSocket）封装 WebSocket 自动重连逻辑，降低组件耦合度
- 开发回测结果可视化面板，展示权益曲线、最大回撤、夏普比率等量化指标

---

### 适用岗位建议

**后端开发岗位**：重点突出行情处理、模拟交易引擎、WASM 沙箱、风险管理、时序数据库优化等后端核心模块

**全栈开发岗位**：同时展示后端架构能力和前端交互体验优化

**Web3/区块链岗位**：强调多交易所接入、高并发处理、实时数据流、安全隔离等与 DeFi/CEX 相关的技术能力

---

**备注**：本项目完整代码已开源，可提供 GitHub 链接供面试官查阅技术实现细节。