## 现状对照（基于当前代码与仓库基础设施）
- 撮合内核已具备：单 symbol 单线程 worker、事件流、Replay、基础风控钩子（`RiskManager`）、集合竞价切换、冰山/止损最小实现。
  - 关键入口：[engine.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/engine.go)，订单簿：[orderbook.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/orderbook.go)
- 仓库里可复用的“持久化事件流/WAL”基座是 NATS JetStream（已用于行情 raw/kline，Durable + ManualAck 可重放）。
  - JetStream 初始化：[nats.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/infrastructure/nats.go)
- 当前缺口集中在三类：
  1) 热路径仍含 `decimal` 与 slice 队列（性能/延迟不可证明）；
  2) 事件未落盘（WAL/快照/恢复/主备切换缺失）；
  3) 反操纵/更完整风控规则尚未实现（只有接口/少量规则）。

---

## 目标
- **低延迟/高吞吐**：撮合热路径避免 `decimal`、避免 O(n) 删除、减少 GC；提供可复现基准与 pprof 证据。
- **强一致 + 高可用**：事件持久化（WAL）+ 快照 + 快速恢复；多实例下同一 symbol 单写者（防双撮合）。
- **反操纵风控**：基于事件流的规则引擎与滑窗统计，支持告警与拦截（配置化）。

---

## 计划（分三条主线并行推进，按优先级落地）

### 1) 低延迟与高吞吐（基准 + 热点重构）
**1.1 基准与剖析先行（不盲改）**
- 新增 `go test -bench` 基准：
  - 单 symbol 连续竞价（限价/市价/撤单）
  - 盘口深度 1e4 价位、1e6 订单的极限路径
- 引入 pprof：对 benchmark 输出 CPU/heap profile，定位热点（decimal 运算、内存分配、队列删除、map access）。

**1.2 热路径数据类型整数化（tick/lot int64）**
- 设计 `PriceTicks int64`、`QtyLots int64` 的内部表示：
  - 外部 API 仍可接受 decimal/string，但进入撮合前量化为整数；
  - 事件落盘也记录 ticks/lots（便于回放一致性）。
- 关键改造点：订单结构、撮合计算、VWAP/fee 计算使用整数或定点（避免 decimal 热路径）。

**1.3 队列结构替换为 O(1) 撤单**
- 将价位队列从 slice FIFO 改为链表/环形队列 + `orderID -> node` 索引：
  - 撤单 O(1)
  - 避免 slice copy 导致的 O(n) 与 GC 压力。

**1.4 内存与分配优化**
- 为订单、事件、trade 使用 `sync.Pool` 或对象复用（基准验证收益再保留）。
- 预分配 channel buffer、批量发布事件（减少系统调用/锁）。

交付物：
- 基准与 pprof 报告（可复现）
- 关键路径整数化与队列 O(1) 撤单实现

### 2) 强一致 + 高可用（WAL/快照/恢复/单写者）
**2.1 JetStream 作为撮合 WAL（append-only log）**
- 新增 `MATCHING` Stream：subjects 形如 `matching.event.*`（按 symbol 分流），配置保留策略与最大容量。
- 在 `emit(Event)` 后增加“事件 sink”：将事件编码（建议 protobuf/定长二进制，或 JSON 先落地）发布到 JetStream。
- Durable Consumer 用于恢复：从上次 `seq`（或 JetStream sequence）重放到内存。

**2.2 快照（snapshot + replay tail）加速恢复**
- 参考现有批量落库模式（`BatchSaver`/`KlineSaver`）：
  - [saver.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/storage/saver.go)
- 设计快照内容：至少包含 open orders（含 iceberg/trigger）、market state、lastPrice、symbol seq。
- 恢复流程：加载最近快照 → 从快照 seq 开始拉 JetStream tail → apply/replay。

**2.3 多实例下“同 symbol 单写者”（避免双撮合）**
- 选型：
  - 方案 A：Postgres advisory lock（若 DB 强依赖且稳定）；
  - 方案 B：JetStream KV/lock（更贴近消息总线）；
  - 方案 C：NATS 服务发现 + lease。
- 启动时抢占 symbol 所属权，失败则只做备份/只读恢复；主挂掉后 lease 过期快速接管（目标 <1s 取决于 lease TTL）。

交付物：
- MATCHING stream + 事件发布/消费恢复
- 快照落库与恢复
- 单写者归属机制（leader/lease）

### 3) 反操纵检测与更完整风控（RiskManager 实现 + 事件消费侧）
**3.1 基于事件流的风控引擎**
- 风控不侵入撮合热路径：使用 `Engine.Events()`（以及后续 JetStream WAL）作为输入，维护用户/席位/品种维度的滑窗统计。

**3.2 第一批规则（可配置阈值）**
- 频率/撤单比：高频挂撤、cancel ratio 过高
- 大额冲击：notional 超阈值 + 连续扫单
- 自成交/对倒：同账户/关联账户的对手盘成交
- Spoofing（虚假挂单）：大单挂出后快速撤销且反复出现
- Layering：多档位堆叠撤单行为

**3.3 输出与处置**
- 输出：`EventAlert`（已有）+ 细分 reason/labels；
- 处置：
  - 软处置：告警 + 限频
  - 硬处置：拒单/冻结（通过 RiskManager.ValidateOrder 返回 error）。

交付物：
- `RiskManager` 的默认实现（规则 + 滑窗统计）
- 指标与告警事件字段标准化

---

## 验收标准
- 性能：提供基准数据（TPS、P50/P99 延迟）、pprof 证据；撮合热路径不再依赖 decimal。
- 一致性：从 WAL + 快照恢复后，订单簿与成交序列可重放一致；断电/重启可在目标时间内恢复。
- HA：同 symbol 保证单写者；主挂后备在 lease 期限内接管。
- 风控：至少 3 类反操纵规则可用，能输出告警并可配置是否拦截。

如果你确认该计划，我将按上述顺序开始实现：先做基准/pprof与整数化方案（避免无效重构），同时并行接入 JetStream 作为撮合 WAL，最后补齐反操纵规则引擎。