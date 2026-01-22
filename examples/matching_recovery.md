# 撮合引擎：WAL + 快照 + 租约（HA）恢复链路解读

目标：理解撮合子系统如何做到“可恢复、可接管、单写者”，以及 app 启动时如何接线。

相关源码：

- 启动接线（App 侧）： [matching_boot.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/matching_boot.go)
- 恢复与泵： [recovery.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/recovery.go)
- WAL / Snapshot / Lease（JetStream）： [wal_jetstream.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/wal_jetstream.go)
- 引擎 restore 注入点： [engine.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/engine.go)

## 1) 目标约束：多实例下同一 symbol 只能有一个 writer

不论是订单簿状态还是 WAL 追加写入，都必须避免“双写”。项目采用：

- JetStream KV 租约（key=symbol，TTL）保证 leader 唯一
- leader 才允许：恢复 → 处理订单 → 写 WAL → 周期快照

对应图示： [matching_ha_wal_snapshot_lease.mmd](file:///Users/yiying/GoProjects/quant-trader/diagrams/matching_ha_wal_snapshot_lease.mmd)

## 2) 恢复（RecoverSymbol）：快照优先 + 回放 WAL tail

入口： [RecoverSymbol](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/recovery.go#L9-L88)

关键步骤（伪代码）：

```go
baseSnap := snapshotStore.Get(symbol)     // KV 取快照（可能为空）
afterSeq := baseSnap.Seq

ob := NewOrderBook(symbol)
ob.SetConfig(baseSnap.Config)

// 1) 先把快照里的 open orders 装回订单簿（缩短恢复时间）
for each orderRecord in baseSnap.OpenOrders:
  ob.AddRestingOrder(orderRecord.ToOrder())

// 2) 只回放 afterSeq 之后的事件（WAL tail）
events := wal.Load(symbol, afterSeq)
for each ev in events:
  switch ev.Type:
    case Submit/Trigger: ob.AddOrder(ev.Order.ToOrder())
    case Cancel: ob.CancelOrder(ev.OrderID)

// 3) 汇总恢复后的 open orders，生成新的 EngineSnapshot 返回
return EngineSnapshot{Seq:lastSeq, OpenOrders:ob.Orders ...}
```

理解要点：

- 快照是“状态基线”，WAL tail 是“增量”，组合后得到当前状态
- WAL 采用 `matching.event.{symbol}` 的 append-only 事件流，恢复可以做到确定性（同一输入→同一状态）

## 3) Engine.Restore：把恢复结果注入到正在运行的 worker

入口： [Engine.Restore](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/engine.go#L162-L187)

核心思想：

- `Engine.AddOrderBook(symbol)` 启动一个 `symbolWorker`（单写者）
- `Restore` 通过向 worker 发送 `opRestore`，在该 goroutine 内重建 `OrderBook`/`triggers`/`seq`/`lastPrice`
- 从而保证“恢复过程”和“后续处理订单”在同一串行上下文中进行

## 4) 持久化泵：StartWALPump / StartSnapshotPump

入口：

- [StartWALPump](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/recovery.go#L107-L118)
- [StartSnapshotPump](file:///Users/yiying/GoProjects/quant-trader/backend/internal/matching/recovery.go#L120-L146)

行为：

- WALPump：持续消费 `engine.Events()` 并 `wal.Append(ev)` 写入 JetStream
- SnapshotPump：按固定间隔 `CaptureSnapshot(symbol)` 并 `store.Put(symbol,snap)` 写入 KV

注意：这是“生产级模式”常见的拆分：\n
状态推进（撮合）与状态落盘（WAL/快照）解耦，但通过单写者与事件序列号约束保持一致性。

## 5) App 侧 leader 循环：Acquire → Restore → Snapshot

入口： [startSymbolLeader](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/matching_boot.go#L58-L106)

策略：

- 抢不到租约：sleep 后重试
- 抢到租约：恢复并 restore；启动 snapshot pump；直到 ctx.Done 或进程退出

你可以把它理解成“每个 symbol 一个 lightweight 的 leader election + watchdog”。

