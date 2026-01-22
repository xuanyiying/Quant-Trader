# quant-trader 学习检查清单（checklist）

用途：把“读懂到能改到能上线”的学习目标拆成可勾选项，适合用于新人入职/自学/团队知识传递。

建议节奏：

- Day 1（2–4h）：跑起来 + 看懂主数据流（raw→kline→db）
- Day 2（4–8h）：补齐 WS/预警/模拟交易的端到端理解
- Day 3+（8–16h）：深入撮合/风控/HA 与性能/可观测性

## A. 能跑起来（Run）

- [ ] 本地依赖启动：TimescaleDB + NATS(-js) + Grafana（见 [backend/docker-compose.yml](file:///Users/yiying/GoProjects/quant-trader/backend/docker-compose.yml)）
- [ ] 后端启动成功，`/health` 返回 ok（见 [app.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/app.go#L174-L177)）
- [ ] `/metrics` 可抓取 Prometheus 指标（见 [metrics.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/infrastructure/metrics.go)）
- [ ] 前端启动成功，能打开页面并看到图表组件加载（见 [frontend/src/App.tsx](file:///Users/yiying/GoProjects/quant-trader/frontend/src/App.tsx)）
- [ ] WS 连通：浏览器 Network 里看到 `/ws` 已连接，且订阅后有消息回推（见 [useWebSocket.ts](file:///Users/yiying/GoProjects/quant-trader/frontend/src/hooks/useWebSocket.ts)）

## B. 能看懂主链路（Understand）

### 启动与装配

- [ ] 能说清 `main.go → App.Init → App.Run` 的职责边界（见 [examples/app_startup_trace.md](file:///Users/yiying/GoProjects/quant-trader/examples/app_startup_trace.md)）
- [ ] 理解哪些 goroutine 常驻运行、哪些依赖 `ctx.Done()` 退出

### 事件总线（NATS）

- [ ] 能列出核心 subjects：`market.raw.*.*`、`market.kline.*.*`、`strategy.signal.*.*`、`notification.user.*`（见 [examples/nats_subjects_and_streams.md](file:///Users/yiying/GoProjects/quant-trader/examples/nats_subjects_and_streams.md)）
- [ ] 理解 durable + manual ack 的语义与代价（至少一次投递 → 下游幂等）

### 行情处理与落库

- [ ] 能画出 `connectors → market.raw → KlineProcessor → market.kline` 的数据流（见 [diagrams/dataflow_market_bus.mmd](file:///Users/yiying/GoProjects/quant-trader/diagrams/dataflow_market_bus.mmd)）
- [ ] 能解释 KlineProcessor 的 key 设计（exchange/symbol/period/window）与 flush 条件（见 [examples/kline_aggregation_walkthrough.md](file:///Users/yiying/GoProjects/quant-trader/examples/kline_aggregation_walkthrough.md)）
- [ ] 能指出 trades/klines 的 DB 写入点与 durable 名称（trade_saver/kline_saver）

### WS 推送

- [ ] 能从前端订阅 topic 追到后端 PushGateway 的 Subscribe/Unsubscribe（见 [diagrams/sequence_ws_subscribe.mmd](file:///Users/yiying/GoProjects/quant-trader/diagrams/sequence_ws_subscribe.mmd)）
- [ ] 能解释 WS 背压风险（慢客户端）以及可能的改进方向（队列/丢弃/断开）

### 预警与模拟交易

- [ ] 能追踪“创建 alert → 触发 → 通知”的时序（见 [diagrams/sequence_alert_trigger.mmd](file:///Users/yiying/GoProjects/quant-trader/diagrams/sequence_alert_trigger.mmd)）
- [ ] 能解释 PaperEngine 为什么订阅 1m K 线而不是 raw（吞吐与复杂度）（见 [examples/paper_engine_txn.md](file:///Users/yiying/GoProjects/quant-trader/examples/paper_engine_txn.md)）

## C. 能改动（Change）

### 新增交易所接入

- [ ] 能在 `internal/connector` 新增一个 Connector，并在 `startIngestionWorker` 中接入
- [ ] 能保证输出 trade 的字段完整且 symbol 经过 Normalize
- [ ] 能为新 connector 增加最小单测（参考 [connector_test.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/connector/connector_test.go)）

### 新增 K 线周期或指标

- [ ] 能在 `model.SupportedPeriods` 里新增周期，并确保 `PeriodToDuration` 支持（见 [backend/internal/model](file:///Users/yiying/GoProjects/quant-trader/backend/internal/model/)）
- [ ] 能在 `internal/indicators` 新增指标并写单测（见 [indicators_test.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/indicators/indicators_test.go)）

### 新增 WS 推送数据类型

- [ ] 能定义新的 subject 并在前端订阅
- [ ] 能保证后端对 topic 做必要的授权控制（如果该 topic 属于用户隐私域）

## D. 能上线（Ship）

### 可观测性

- [ ] 能列出关键指标：ingest_latency、trade_process_total、db_insert_total、ws_connections_total、goroutine_count（见 [metrics.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/infrastructure/metrics.go)）
- [ ] 能指出关键日志点：connectors 重连、processor 队列满、DB flush 失败、WS 写失败
- [ ] 能导入 Grafana dashboard 并解释主要面板含义（见 [quant-trader.json](file:///Users/yiying/GoProjects/quant-trader/backend/monitoring/grafana/dashboards/quant-trader.json)）

### CI 与质量门槛

- [ ] 能解释 CI 做了什么（build + test + coverage 统计），覆盖率阈值为何未强制 fail（见 [ci.yml](file:///Users/yiying/GoProjects/quant-trader/.github/workflows/ci.yml)）
- [ ] 能补齐一个关键路径的测试并把覆盖率拉高（建议：push、storage、api handler）

### 撮合子系统（进阶）

- [ ] 能解释 “单写者 worker” 的一致性来源（`symbolWorker` 串行处理）
- [ ] 能解释 “WAL + 快照 + 租约” 如何实现恢复与接管（见 [examples/matching_recovery.md](file:///Users/yiying/GoProjects/quant-trader/examples/matching_recovery.md)）
- [ ] 能做一次恢复演练：清空内存 → 从快照+WAL 恢复 → 状态一致

