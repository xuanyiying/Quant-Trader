# Quant-Trader 系统维护与操作手册

本手册旨在为技术团队提供系统的日常维护、故障排查及功能扩展指南。

## 1. 系统架构概览

系统采用前后端分离架构：

- **前端**: React + Vite + TailwindCSS + ECharts。
- **后端**: Go (Gin) + TimescaleDB + NATS JetStream。
- **核心组件**:
  - `PaperEngine`: 模拟撮合引擎，支持高并发成交模拟。
  - `RiskManager`: 事前风控逻辑，防止异常下单。
  - `AnalyticsService`: 基于时序数据的量化指标分析。
  - `WasmRunner`: 安全隔离的策略执行沙箱。

## 2. 环境配置与部署

### 关键环境变量

- `DATABASE_URL`: TimescaleDB 连接字符串。
- `NATS_URL`: NATS 服务器地址。
- `STRIPE_API_KEY`: Stripe 支付网关密钥。
- `STRIPE_WEBHOOK_SECRET`: Stripe 回调验证密钥。

### 部署建议

- **数据库**: 建议启用 TimescaleDB 的压缩功能以节省历史 K 线存储空间。
- **NATS**: 建议配置 JetStream 持久化，确保行情与订单数据不丢失。

## 3. 日常维护操作

### 数据库分区管理

系统已配置自动分区（Hypertables）。定期检查分区状态：

```sql
SELECT * FROM timescaledb_information.hypertable WHERE table_name = 'klines';
```

### 策略沙箱 (WASM) 维护

上传的策略必须通过 `WasmRunner` 加载。确保服务器已安装 `wazero` 运行时的相关依赖（Go 工具链会自动处理）。

## 4. 商业功能管理

### 订阅系统

- **升级流程**: 用户通过前端 `handleUpgrade` 跳转至 Stripe Checkout。
- **回调处理**: `subscription_handler.go` 接收 Webhook 并通过 `StripeService` 更新用户 Tier。

### 模拟交易重置

如需重置用户模拟账户余额：

```sql
UPDATE paper_accounts SET balance = 10000 WHERE user_id = ?;
DELETE FROM paper_positions WHERE user_id = ?;
```

## 5. 常见问题排查 (Troubleshooting)

| 问题 | 可能原因 | 解决方法 |
| --- | --- | --- |
| K 线延迟 | Goroutine 池 Worker 过少 | 调整 `processor/kline.go` 中的 Worker 数量 |
| 下单失败 | 触发 RiskManager 规则 | 查看 `RiskManager` 日志，确认资金是否充足 |
| Webhook 400 | Secret 不匹配 | 检查环境变量 `STRIPE_WEBHOOK_SECRET` |
| 前端 build 失败 | 缓存污染 | 运行 `rm -rf node_modules && npm install` |

---
*Quant-Trader 商业版 v1.0.0*
