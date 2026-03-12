# API 端点文档

本文档描述了 quant-trader 后端所有 API 端点的请求和响应格式。

## 通用说明

### 认证方式

需要认证的端点需要在请求头中携带 JWT Token：

```
Authorization: Bearer <token>
```

### 通用错误响应格式

```json
{
  "error": "错误描述信息"
}
```

---

## Authentication 认证模块

### POST /api/v1/auth/register

用户注册

**请求体**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 用户邮箱，需符合邮箱格式 |
| password | string | 是 | 密码，最少6个字符 |

**请求示例**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**成功响应** (201 Created)

```json
{
  "message": "user created",
  "id": 1
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 Bad Request | 请求参数格式错误 |
| 409 Conflict | 邮箱已被注册 |

---

### POST /api/v1/auth/login

用户登录

**请求体**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 用户邮箱 |
| password | string | 是 | 密码 |

**请求示例**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**成功响应** (200 OK)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 Bad Request | 请求参数格式错误 |
| 401 Unauthorized | 邮箱或密码错误 |
| 500 Internal Server Error | 服务器内部错误 |

---

## Paper Trading 模拟交易模块

### GET /api/v1/paper/account

获取模拟交易账户信息

**认证** 需要

**成功响应** (200 OK)

```json
{
  "balance": "100000.00000000"
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 500 Internal Server Error | 获取账户失败 |

---

### POST /api/v1/paper/account/reset

重置模拟交易账户

**认证** 需要

**成功响应** (200 OK)

```json
{
  "balance": "100000",
  "message": "account reset successful"
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 500 Internal Server Error | 重置账户失败 |

---

### POST /api/v1/paper/orders

创建模拟交易订单

**认证** 需要

**请求体**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| symbol | string | 是 | 交易对符号，如 "BTCUSDT" |
| side | string | 是 | 订单方向："buy" 或 "sell" |
| type | string | 是 | 订单类型："market" 或 "limit" |
| price | decimal | 否 | 限价单价格（限价单时建议提供） |
| qty | decimal | 是 | 交易数量 |

**请求示例**

```json
{
  "symbol": "BTCUSDT",
  "side": "buy",
  "type": "market",
  "qty": "0.01"
}
```

**成功响应** (201 Created)

```json
{
  "id": 1
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 Bad Request | 请求参数格式错误、无效的订单方向或类型 |
| 403 Forbidden | 风控检查未通过（余额不足等） |

---

### GET /api/v1/paper/positions

获取当前持仓列表

**认证** 需要

**成功响应** (200 OK)

```json
[
  {
    "symbol": "BTCUSDT",
    "qty": "0.01",
    "avg_price": "50000.00000000"
  }
]
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 500 Internal Server Error | 获取持仓失败 |

---

## Portfolio 投资组合模块

### GET /api/v1/portfolio/report

获取投资组合报告

**认证** 需要

**成功响应** (200 OK)

```json
{
  "total_return": "10.50",
  "max_drawdown": "5.20",
  "sharpe_ratio": 1.25,
  "win_rate": 0.65,
  "total_trades": 100,
  "profitable_trades": 65
}
```

**响应字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| total_return | decimal | 总收益率（百分比） |
| max_drawdown | decimal | 最大回撤（百分比） |
| sharpe_ratio | float | 夏普比率 |
| win_rate | float | 胜率（0-1之间） |
| total_trades | int64 | 总交易次数 |
| profitable_trades | int64 | 盈利交易次数 |

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 500 Internal Server Error | 生成报告失败 |

---

## Alerts 告警模块

### GET /api/v1/alerts

获取用户告警列表

**认证** 需要

**成功响应** (200 OK)

```json
[
  {
    "id": 1,
    "symbol": "BTCUSDT",
    "condition_type": "price_above",
    "target_value": "60000.00000000",
    "is_triggered": false,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

**响应字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 告警ID |
| symbol | string | 交易对符号 |
| condition_type | string | 条件类型 |
| target_value | decimal | 目标价格值 |
| is_triggered | bool | 是否已触发 |
| created_at | time | 创建时间 |

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 500 Internal Server Error | 获取告警列表失败 |

---

### POST /api/v1/alerts

创建价格告警

**认证** 需要

**请求体**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| symbol | string | 是 | 交易对符号 |
| condition_type | string | 是 | 条件类型（如 "price_above", "price_below"） |
| target_value | decimal | 是 | 目标价格值 |

**请求示例**

```json
{
  "symbol": "BTCUSDT",
  "condition_type": "price_above",
  "target_value": "60000"
}
```

**成功响应** (201 Created)

```json
{
  "id": 1
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 Bad Request | 请求参数格式错误 |
| 500 Internal Server Error | 创建告警失败 |

---

### DELETE /api/v1/alerts/:id

删除告警

**认证** 需要

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int64 | 告警ID |

**成功响应** (200 OK)

```json
{
  "message": "alert deleted"
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 Bad Request | 无效的告警ID格式 |
| 404 Not Found | 告警不存在 |
| 500 Internal Server Error | 删除告警失败 |

---

## Subscription 订阅模块

### GET /api/v1/subscription

获取用户订阅信息

**认证** 需要

**成功响应** (200 OK)

```json
{
  "tier_name": "Pro",
  "max_symbols": 10,
  "status": "active",
  "expires_at": "2024-12-31"
}
```

**响应字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| tier_name | string | 订阅层级名称（"Free", "Pro", "Enterprise"） |
| max_symbols | int | 最大可监控交易对数量 |
| status | string | 订阅状态（"active", "canceled", "expired", "past_due", "trialing", "unpaid"） |
| expires_at | string | 到期日期 |

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 404 Not Found | 订阅不存在 |
| 500 Internal Server Error | 获取订阅信息失败 |

---

### POST /api/v1/subscription/checkout

创建 Stripe 结账会话

**认证** 需要

**请求体**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| price_id | string | 是 | Stripe 价格ID |

**请求示例**

```json
{
  "price_id": "price_xxxxx"
}
```

**成功响应** (200 OK)

```json
{
  "url": "https://checkout.stripe.com/xxx"
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 Bad Request | 请求参数格式错误 |
| 500 Internal Server Error | 创建结账会话失败 |

---

## Klines K线数据模块

### GET /api/v1/klines/latest

获取最新K线数据

**认证** 需要

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| symbol | string | 是 | 交易对符号（URL路径参数） |

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| period | string | 否 | "1m" | K线周期 |

**支持的K线周期**

| 周期 | 说明 |
|------|------|
| 1m | 1分钟 |
| 5m | 5分钟 |
| 15m | 15分钟 |
| 1h | 1小时 |
| 4h | 4小时 |
| 1d | 1天 |

**请求示例**

```
GET /api/v1/klines/latest/BTCUSDT?period=1h
```

**成功响应** (200 OK)

```json
[
  {
    "id": 1,
    "symbol": "BTCUSDT",
    "exchange": "binance",
    "period": "1h",
    "o": "50000.00000000",
    "h": "51000.00000000",
    "l": "49500.00000000",
    "c": "50500.00000000",
    "v": "1000.00000000",
    "t": "2024-01-01T00:00:00Z"
  }
]
```

**响应字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | K线ID |
| symbol | string | 交易对符号 |
| exchange | string | 交易所名称 |
| period | string | K线周期 |
| o | decimal | 开盘价 |
| h | decimal | 最高价 |
| l | decimal | 最低价 |
| c | decimal | 收盘价 |
| v | decimal | 成交量 |
| t | time | 时间戳 |

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 500 Internal Server Error | 查询K线数据失败 |

---

### POST /api/v1/klines/backfill

触发K线数据回填任务

**认证** 需要

**请求体**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| exchange | string | 是 | 交易所名称（"binance", "bybit", "kraken" 等） |
| symbol | string | 是 | 交易对符号 |
| start_time | time | 是 | 开始时间（RFC3339格式） |
| end_time | time | 是 | 结束时间（RFC3339格式） |

**请求示例**

```json
{
  "exchange": "binance",
  "symbol": "BTCUSDT",
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-01-31T23:59:59Z"
}
```

**成功响应** (202 Accepted)

```json
{
  "message": "backfill task started"
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 Bad Request | 请求参数格式错误 |

---

## Marketplace 策略市场模块

### GET /api/v1/marketplace

获取公开策略列表

**认证** 需要

**成功响应** (200 OK)

```json
[
  {
    "id": 1,
    "price": "99.00",
    "description": "一个高胜率的趋势跟踪策略",
    "metrics": {
      "total_return": 150.5,
      "sharpe_ratio": 2.1,
      "max_drawdown": 8.5,
      "win_rate": 0.72,
      "total_trades": 500,
      "backtest_period": "2023-01-01 to 2024-01-01"
    },
    "name": "趋势跟踪策略",
    "author": "author@example.com"
  }
]
```

**响应字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 策略ID |
| price | decimal | 策略价格 |
| description | string | 策略描述 |
| metrics | object | 策略绩效指标 |
| name | string | 策略名称 |
| author | string | 作者邮箱 |

**metrics 对象字段**

| 字段 | 类型 | 说明 |
|------|------|------|
| total_return | float | 总收益率（百分比） |
| sharpe_ratio | float | 夏普比率 |
| max_drawdown | float | 最大回撤（百分比） |
| win_rate | float | 胜率 |
| total_trades | int | 总交易次数 |
| backtest_period | string | 回测时间段 |

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 500 Internal Server Error | 获取策略列表失败 |

---

### POST /api/v1/marketplace/:id/purchase

购买策略

**认证** 需要

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int64 | 策略ID |

**成功响应** (200 OK)

```json
{
  "message": "purchase successful"
}
```

**错误响应**

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 Bad Request | 无效的策略ID格式 |
| 404 Not Found | 策略不存在 |
| 500 Internal Server Error | 购买失败 |

---

## 附录

### HTTP 状态码说明

| 状态码 | 说明 |
|--------|------|
| 200 OK | 请求成功 |
| 201 Created | 资源创建成功 |
| 202 Accepted | 请求已接受，正在异步处理 |
| 400 Bad Request | 请求参数错误 |
| 401 Unauthorized | 未认证或认证失败 |
| 403 Forbidden | 无权限访问 |
| 404 Not Found | 资源不存在 |
| 409 Conflict | 资源冲突 |
| 500 Internal Server Error | 服务器内部错误 |

### 订阅层级说明

| 层级 | 最大交易对数 | 说明 |
|------|-------------|------|
| Free | 1 | 免费版 |
| Pro | 10 | 专业版 |
| Enterprise | 无限制 | 企业版 |

### 订单类型说明

| 类型 | 说明 |
|------|------|
| market | 市价单 |
| limit | 限价单 |

### 订单方向说明

| 方向 | 说明 |
|------|------|
| buy | 买入 |
| sell | 卖出 |

### 订单状态说明

| 状态 | 说明 |
|------|------|
| new | 新建 |
| open | 待成交 |
| partially_filled | 部分成交 |
| filled | 完全成交 |
| cancelled | 已取消 |
| rejected | 已拒绝 |
| expired | 已过期 |
