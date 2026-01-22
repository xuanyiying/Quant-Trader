# 启动链路注释式解读（main → App.Init → App.Run）

目标：用“能画出调用链”的粒度理解后端如何把 DB/NATS/各服务装配起来，并掌握哪些 goroutine 长期运行、如何退出。

相关源码：

- 入口： [main.go](file:///Users/yiying/GoProjects/quant-trader/backend/cmd/main.go)
- 应用生命周期： [app.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/app.go)
- NATS/JetStream 初始化： [nats.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/infrastructure/nats.go)

## 1) main.go：最薄入口

```go
// main.go (精简片段)
application, err := app.NewApp() // 读取配置 + 初始化logger
if err != nil { log.Fatalf(...) }

ctx := context.Background()
if err := application.Init(ctx); err != nil { log.Fatalf(...) } // DB + NATS + service new
if err := application.Run(ctx); err != nil { log.Fatalf(...) }  // 启动后台组件 + HTTP
```

理解要点：

- `NewApp()` 只做“轻量、纯创建”的事（配置/日志），不做网络连接
- `Init()` 才会做外部依赖连接（DB/NATS）
- `Run()` 启动各后台循环（订阅/flush/worker）并启动 HTTP server

## 2) App.Init：把外部依赖接入进来

源码： [App.Init](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/app.go#L58-L85)

```go
// App.Init (精简片段)
dbPool, err := pgxpool.New(ctx, a.Config.DB_DSN) // 建连接池
a.DB = dbPool

_ = a.initDatabase(ctx) // 读取 scripts/init.sql 并执行（创建表/超表等）

nc, js, err := infrastructure.InitNATS(a.Config.NatsURL, a.Logger) // nats connect + js
a.NC = nc
a.JS = js

// 依赖注入：把 db/js/logger 传给各 service
a.PushGateway = push.NewPushGateway(js, a.Logger)
a.AlertService = alert.NewAlertService(a.DB, js, a.Logger)
a.PaperEngine = paper.NewPaperEngine(a.DB, js, a.Logger)
```

理解要点：

- `InitNATS()` 会确保 `MARKET` stream 存在（subjects：raw/kline），详见 [nats.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/infrastructure/nats.go)
- 此时只是“new service”，并没有开始订阅/处理；真正开始是在 `Run()`

## 3) App.Run：启动后台组件与对外接口

源码： [App.Run](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/app.go#L88-L130)

核心顺序（按依赖关系）：

1. 启动持久化订阅：订阅 raw/kline 并批量落库（trade_saver/kline_saver durable）
2. 启动 KlineProcessor：订阅 raw → 聚合 → publish kline
3. 启动 AlertService：订阅 1m kline → 规则评估 → publish notification
4. 启动 PaperEngine：订阅 1m kline → 撮合 → DB 事务更新
5. 启动 Ingestion Worker：connectors → publish raw
6. 启动 StrategyRunner：订阅 kline → publish strategy.signal
7. 启动 HTTP Server + `/ws`（WS 由 PushGateway 处理）

其中 matching 的启动逻辑在：

- [matching_boot.go](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/matching_boot.go)

它会在 `Run()` 开头初始化 JetStream WAL/KV，并为每个 symbol 运行 leader 循环（租约→恢复→快照）。

## 4) 退出与优雅停机

源码： [waitForShutdown](file:///Users/yiying/GoProjects/quant-trader/backend/internal/app/app.go#L132-L151)

- 监听 SIGINT/SIGTERM
- cancel runCtx（让后台 goroutine 通过 ctx.Done() 退出）
- 关闭 HTTP server、NATS、DB

建议你阅读时重点标出：哪些 goroutine 有 `ctx.Done()` 分支，哪些没有（没有的需要后续补齐，避免卡住退出）。

