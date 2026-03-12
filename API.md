# Quant-Trader API Documentation

[English](#english) | [简体中文](#中文文档)

---

# English

Complete API reference for the Quant-Trader backend service.

---

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Endpoints Overview](#endpoints-overview)
- [Authentication Endpoints](#authentication-endpoints)
- [Market Data Endpoints](#market-data-endpoints)
- [Paper Trading Endpoints](#paper-trading-endpoints)
- [Alert Endpoints](#alert-endpoints)
- [Portfolio Endpoints](#portfolio-endpoints)
- [Marketplace Endpoints](#marketplace-endpoints)
- [Backtest & Analytics Endpoints](#backtest--analytics-endpoints)
- [Subscription Endpoints](#subscription-endpoints)
- [Error Handling](#error-handling)
- [Rate Limits](#rate-limits)

---

## Base URL

```
http://localhost:8080
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

---

## Endpoints Overview

| Category | Endpoints |
|----------|-----------|
| [Authentication](#authentication-endpoints) | `/api/v1/register`, `/api/v1/login` |
| [Market Data](#market-data-endpoints) | `/api/v1/klines/:symbol`, `/ws` |
| [Paper Trading](#paper-trading-endpoints) | `/api/v1/paper/*` |
| [Alerts](#alert-endpoints) | `/api/v1/alerts` |
| [Portfolios](#portfolio-endpoints) | `/api/v1/portfolios` |
| [Marketplace](#marketplace-endpoints) | `/api/v1/market/strategies` |
| [Backtest & Analytics](#backtest--analytics-endpoints) | `/api/v1/backtest`, `/api/v1/analytics/*` |
| [Subscription](#subscription-endpoints) | `/api/v1/subscription` |

---

## Authentication Endpoints

### Register User

Create a new user account.

**Endpoint:** `POST /api/v1/register`

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User email address |
| `password` | string | Yes | User password (min 8 characters) |

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**Response:**

```json
{
  "id": "user_123",
  "email": "user@example.com",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Login

Authenticate user and receive JWT token.

**Endpoint:** `POST /api/v1/login`

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User email address |
| `password` | string | Yes | User password |

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user_123",
    "email": "user@example.com"
  },
  "expires_at": "2024-01-16T10:30:00Z"
}
```

---

## Market Data Endpoints

### Get Historical K-Lines

Retrieve historical candlestick data for a symbol.

**Endpoint:** `GET /api/v1/klines/:symbol`

**Authentication:** Not required (rate limited)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `symbol` | path | Yes | - | Trading pair (e.g., BTCUSDT) |
| `interval` | query | No | 1m | Time interval (1m, 5m, 15m, 1h, 4h, 1d) |
| `limit` | query | No | 100 | Number of records (max 1000) |
| `start_time` | query | No | - | Start timestamp (Unix ms) |
| `end_time` | query | No | - | End timestamp (Unix ms) |

**Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/klines/BTCUSDT?interval=1h&limit=100"
```

**Response:**

```json
{
  "symbol": "BTCUSDT",
  "interval": "1h",
  "data": [
    {
      "timestamp": 1704067200000,
      "open": "42000.00",
      "high": "42500.00",
      "low": "41800.00",
      "close": "42300.00",
      "volume": "1234.56",
      "trades": 5678
    }
  ]
}
```

### WebSocket Market Stream

Real-time market data WebSocket connection.

**Endpoint:** `GET /ws`

**Authentication:** Optional

**Message Format:**

```javascript
// Subscribe
{
  "action": "subscribe",
  "symbols": ["BTCUSDT", "ETHUSDT"]
}

// Unsubscribe
{
  "action": "unsubscribe",
  "symbols": ["BTCUSDT"]
}

// Market Update
{
  "type": "trade",
  "symbol": "BTCUSDT",
  "price": "42350.00",
  "quantity": "0.1234",
  "timestamp": 1704067200000
}

// Kline Update
{
  "type": "kline",
  "symbol": "BTCUSDT",
  "interval": "1m",
  "open": "42300.00",
  "high": "42400.00",
  "low": "42250.00",
  "close": "42350.00",
  "volume": "123.45",
  "timestamp": 1704067260000
}
```

**JavaScript Example:**

```javascript
const ws = new WebSocket('http://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to WebSocket');
  ws.send(JSON.stringify({
    action: 'subscribe',
    symbols: ['BTCUSDT', 'ETHUSDT']
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected');
};
```

---

## Paper Trading Endpoints

### Get Paper Account

Get current paper trading account balance.

**Endpoint:** `GET /api/v1/paper/account`

**Authentication:** Required

**Example:**

```bash
curl -X GET http://localhost:8080/api/v1/paper/account \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "user_id": "user_123",
  "balances": {
    "USDT": "10000.00",
    "BTC": "0.50",
    "ETH": "5.00"
  },
  "total_equity": "15000.00",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Reset Paper Account

Reset paper trading account to initial balance.

**Endpoint:** `POST /api/v1/paper/account/reset`

**Authentication:** Required

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/paper/account/reset \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "success": true,
  "message": "Account reset successfully"
}
```

### Create Paper Order

Submit a paper trading order.

**Endpoint:** `POST /api/v1/paper/orders`

**Authentication:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `symbol` | string | Yes | Trading pair (e.g., BTCUSDT) |
| `side` | string | Yes | Order side: `buy` or `sell` |
| `type` | string | Yes | Order type: `limit` or `market` |
| `price` | string | No* | Limit price (*required for limit orders) |
| `quantity` | string | Yes | Order quantity |
| `time_in_force` | string | No | Time in force: `GTC`, `IOC`, `FOK` (default: GTC) |

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/paper/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "side": "buy",
    "type": "limit",
    "price": "42000.00",
    "quantity": "0.1"
  }'
```

**Response:**

```json
{
  "order_id": "order_abc123",
  "symbol": "BTCUSDT",
  "side": "buy",
  "type": "limit",
  "price": "42000.00",
  "quantity": "0.1",
  "filled": "0.0",
  "status": "pending",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Get Paper Positions

Get all open positions in paper trading.

**Endpoint:** `GET /api/v1/paper/positions`

**Authentication:** Required

**Example:**

```bash
curl -X GET http://localhost:8080/api/v1/paper/positions \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "positions": [
    {
      "symbol": "BTCUSDT",
      "side": "long",
      "quantity": "0.5",
      "entry_price": "41500.00",
      "current_price": "42300.00",
      "unrealized_pnl": "400.00",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

## Alert Endpoints

### List Alerts

Get all alerts for the authenticated user.

**Endpoint:** `GET /api/v1/alerts`

**Authentication:** Required

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | string | Filter by status: `active`, `triggered`, `disabled` |

**Example:**

```bash
curl -X GET http://localhost:8080/api/v1/alerts \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "alerts": [
    {
      "id": "alert_123",
      "symbol": "BTCUSDT",
      "condition": "price_above",
      "value": "50000.00",
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Create Alert

Create a new price or indicator alert.

**Endpoint:** `POST /api/v1/alerts`

**Authentication:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `symbol` | string | Yes | Trading pair |
| `condition` | string | Yes | Alert condition |
| `value` | number | Yes | Target value |
| `type` | string | Yes | Alert type: `price`, `indicator` |
| `indicator` | string | No* | Indicator name (*required if type is indicator) |
| `enabled` | boolean | No | Enable alert (default: true) |

**Supported Conditions:**
- `price_above` - Price greater than value
- `price_below` - Price less than value
- `price_crosses` - Price crosses value
- `indicator_above` - Indicator greater than value
- `indicator_below` - Indicator less than value

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "condition": "price_above",
    "value": 50000,
    "type": "price",
    "enabled": true
  }'
```

**Response:**

```json
{
  "id": "alert_456",
  "symbol": "BTCUSDT",
  "condition": "price_above",
  "value": "50000.00",
  "type": "price",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Delete Alert

Delete an alert.

**Endpoint:** `DELETE /api/v1/alerts/:id`

**Authentication:** Required

**Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/alerts/alert_456 \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "success": true,
  "message": "Alert deleted successfully"
}
```

---

## Portfolio Endpoints

### List Portfolios

Get all portfolios for the authenticated user.

**Endpoint:** `GET /api/v1/portfolios`

**Authentication:** Required

**Example:**

```bash
curl -X GET http://localhost:8080/api/v1/portfolios \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "portfolios": [
    {
      "id": "portfolio_123",
      "name": "My Trading Portfolio",
      "balance": "15000.00",
      "positions": 3,
      "created_at": "2024-01-10T10:00:00Z"
    }
  ]
}
```

### Create Portfolio

Create a new portfolio.

**Endpoint:** `POST /api/v1/portfolios`

**Authentication:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Portfolio name |
| `initial_balance` | string | No | Initial balance (default: 10000) |

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/portfolios \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Conservative Portfolio",
    "initial_balance": "50000.00"
  }'
```

**Response:**

```json
{
  "id": "portfolio_456",
  "name": "Conservative Portfolio",
  "balance": "50000.00",
  "positions": 0,
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

## Marketplace Endpoints

### List Strategies

Browse available trading strategies in the marketplace.

**Endpoint:** `GET /api/v1/market/strategies`

**Authentication:** Not required

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `category` | string | Filter by category |
| `sort_by` | string | Sort by: `popular`, `newest`, `price` |
| `page` | integer | Page number |
| `limit` | integer | Items per page |

**Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/market/strategies?sort_by=popular&limit=10"
```

**Response:**

```json
{
  "strategies": [
    {
      "id": "strategy_123",
      "name": "MA Crossover Bot",
      "description": "Moving average crossover strategy",
      "category": "trend",
      "price": "29.99",
      "subscribers": 150,
      "rating": 4.5,
      "author": "trader_pro"
    }
  ],
  "total": 50,
  "page": 1,
  "total_pages": 5
}
```

### Purchase Strategy

Subscribe to a trading strategy.

**Endpoint:** `POST /api/v1/market/strategies/:id/purchase`

**Authentication:** Required

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/market/strategies/strategy_123/purchase \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "success": true,
  "subscription_id": "sub_abc123",
  "strategy_id": "strategy_123",
  "expires_at": "2024-02-15T10:30:00Z"
}
```

---

## Backtest & Analytics Endpoints

### Run Backtest

Execute a backtest for a strategy.

**Endpoint:** `POST /api/v1/backtest`

**Authentication:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `strategy_id` | string | Yes | Strategy to backtest |
| `symbol` | string | Yes | Trading pair |
| `start_time` | string | Yes | Start timestamp |
| `end_time` | string | Yes | End timestamp |
| `initial_capital` | string | No | Initial capital |
| `parameters` | object | No | Strategy parameters |

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/backtest \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "strategy_id": "strategy_123",
    "symbol": "BTCUSDT",
    "start_time": "1704067200000",
    "end_time": "1704067200000",
    "initial_capital": "10000.00",
    "parameters": {
      "fast_period": 10,
      "slow_period": 20
    }
  }'
```

**Response:**

```json
{
  "backtest_id": "bt_abc123",
  "status": "running",
  "progress": 0
}
```

### Get Portfolio Analytics

Get detailed portfolio performance analytics.

**Endpoint:** `GET /api/v1/analytics/portfolio`

**Authentication:** Required

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `portfolio_id` | string | Portfolio ID |
| `start_time` | string | Start timestamp |
| `end_time` | string | End timestamp |

**Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/analytics/portfolio?portfolio_id=portfolio_123" \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "portfolio_id": "portfolio_123",
  "total_return": "15.5%",
  "sharpe_ratio": 1.8,
  "max_drawdown": "8.2%",
  "win_rate": "65%",
  "profit_factor": 2.1,
  "total_trades": 150,
  "average_trade": "35.00",
  "daily_returns": [...]
}
```

---

## Subscription Endpoints

### Get Subscription

Get current subscription status.

**Endpoint:** `GET /api/v1/subscription`

**Authentication:** Required

**Example:**

```bash
curl -X GET http://localhost:8080/api/v1/subscription \
  -H "Authorization: Bearer <token>"
```

**Response:**

```json
{
  "tier": "pro",
  "status": "active",
  "expires_at": "2025-01-15T10:30:00Z",
  "features": {
    "api_rate_limit": 1000,
    "concurrent_strategies": 5,
    "backtest_history": 100,
    "alerts": 50
  }
}
```

---

## Error Handling

All endpoints may return error responses in the following format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Invalid or missing authentication |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Invalid request parameters |
| `RATE_LIMITED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Rate Limits

| Tier | Requests/minute |
|------|-----------------|
| Free | 60 |
| Pro | 600 |
| Enterprise | 6000 |

---

# 中文文档

Quant-Trader 后端服务完整 API 参考文档。

---

## 目录

- [基础 URL](#基础-url)
- [认证方式](#认证方式)
- [接口概览](#接口概览)
- [认证接口](#认证接口-1)
- [市场数据接口](#市场数据接口)
- [模拟交易接口](#模拟交易接口)
- [告警接口](#告警接口)
- [投资组合接口](#投资组合接口)
- [策略市场接口](#策略市场接口)
- [回测与分析接口](#回测与分析接口)
- [订阅接口](#订阅接口)
- [错误处理](#错误处理-1)
- [限流规则](#限流规则)

---

## 基础 URL

```
http://localhost:8080
```

## 认证方式

大部分接口需要 JWT 认证。在请求头中包含令牌：

```
Authorization: Bearer <your-jwt-token>
```

---

## 接口概览

| 类别 | 接口 |
|------|------|
| [认证](#认证接口-1) | `/api/v1/register`, `/api/v1/login` |
| [市场数据](#市场数据接口) | `/api/v1/klines/:symbol`, `/ws` |
| [模拟交易](#模拟交易接口) | `/api/v1/paper/*` |
| [告警](#告警接口) | `/api/v1/alerts` |
| [投资组合](#投资组合接口) | `/api/v1/portfolios` |
| [策略市场](#策略市场接口) | `/api/v1/market/strategies` |
| [回测与分析](#回测与分析接口) | `/api/v1/backtest`, `/api/v1/analytics/*` |
| [订阅](#订阅接口) | `/api/v1/subscription` |

---

## 认证接口

### 用户注册

创建新用户账户。

**接口:** `POST /api/v1/register`

**请求体:**

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `email` | string | 是 | 用户邮箱地址 |
| `password` | string | 是 | 用户密码（至少8位） |

**示例:**

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**响应:**

```json
{
  "id": "user_123",
  "email": "user@example.com",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 用户登录

认证用户并获取 JWT 令牌。

**接口:** `POST /api/v1/login`

**请求体:**

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `email` | string | 是 | 用户邮箱地址 |
| `password` | string | 是 | 用户密码 |

**示例:**

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**响应:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user_123",
    "email": "user@example.com"
  },
  "expires_at": "2024-01-16T10:30:00Z"
}
```

---

## 市场数据接口

### 获取历史 K线数据

获取交易对的历史 K线（蜡烛图）数据。

**接口:** `GET /api/v1/klines/:symbol`

**认证:** 不需要（有限流）

**查询参数:**

| 参数 | 类型 | 必需 | 默认值 | 描述 |
|------|------|------|--------|------|
| `symbol` | path | 是 | - | 交易对（如 BTCUSDT） |
| `interval` | query | 否 | 1m | 时间周期（1m, 5m, 15m, 1h, 4h, 1d） |
| `limit` | query | 否 | 100 | 记录数量（最大1000） |
| `start_time` | query | 否 | - | 开始时间戳（Unix 毫秒） |
| `end_time` | query | 否 | - | 结束时间戳（Unix 毫秒） |

**示例:**

```bash
curl -X GET "http://localhost:8080/api/v1/klines/BTCUSDT?interval=1h&limit=100"
```

**响应:**

```json
{
  "symbol": "BTCUSDT",
  "interval": "1h",
  "data": [
    {
      "timestamp": 1704067200000,
      "open": "42000.00",
      "high": "42500.00",
      "low": "41800.00",
      "close": "42300.00",
      "volume": "1234.56",
      "trades": 5678
    }
  ]
}
```

### WebSocket 市场数据流

实时市场数据 WebSocket 连接。

**接口:** `GET /ws`

**认证:** 可选

**消息格式:**

```javascript
// 订阅
{
  "action": "subscribe",
  "symbols": ["BTCUSDT", "ETHUSDT"]
}

// 取消订阅
{
  "action": "unsubscribe",
  "symbols": ["BTCUSDT"]
}

// 成交更新
{
  "type": "trade",
  "symbol": "BTCUSDT",
  "price": "42350.00",
  "quantity": "0.1234",
  "timestamp": 1704067200000
}

// K线更新
{
  "type": "kline",
  "symbol": "BTCUSDT",
  "interval": "1m",
  "open": "42300.00",
  "high": "42400.00",
  "low": "42250.00",
  "close": "42350.00",
  "volume": "123.45",
  "timestamp": 1704067260000
}
```

**JavaScript 示例:**

```javascript
const ws = new WebSocket('http://localhost:8080/ws');

ws.onopen = () => {
  console.log('WebSocket 已连接');
  ws.send(JSON.stringify({
    action: 'subscribe',
    symbols: ['BTCUSDT', 'ETHUSDT']
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('收到数据:', data);
};

ws.onerror = (error) => {
  console.error('WebSocket 错误:', error);
};

ws.onclose = () => {
  console.log('连接已断开');
};
```

---

## 模拟交易接口

### 获取模拟账户

获取当前模拟交易账户余额。

**接口:** `GET /api/v1/paper/account`

**认证:** 必需

**示例:**

```bash
curl -X GET http://localhost:8080/api/v1/paper/account \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "user_id": "user_123",
  "balances": {
    "USDT": "10000.00",
    "BTC": "0.50",
    "ETH": "5.00"
  },
  "total_equity": "15000.00",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### 重置模拟账户

将模拟交易账户重置为初始余额。

**接口:** `POST /api/v1/paper/account/reset`

**认证:** 必需

**示例:**

```bash
curl -X POST http://localhost:8080/api/v1/paper/account/reset \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "success": true,
  "message": "账户重置成功"
}
```

### 创建模拟订单

提交模拟交易订单。

**接口:** `POST /api/v1/paper/orders`

**认证:** 必需

**请求体:**

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `symbol` | string | 是 | 交易对（如 BTCUSDT） |
| `side` | string | 是 | 订单方向：`buy` 或 `sell` |
| `type` | string | 是 | 订单类型：`limit` 或 `market` |
| `price` | string | 否* | 限价（*限价单必需） |
| `quantity` | string | 是 | 订单数量 |
| `time_in_force` | string | 否 | 有效时间：`GTC`, `IOC`, `FOK`（默认：GTC） |

**示例:**

```bash
curl -X POST http://localhost:8080/api/v1/paper/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "side": "buy",
    "type": "limit",
    "price": "42000.00",
    "quantity": "0.1"
  }'
```

**响应:**

```json
{
  "order_id": "order_abc123",
  "symbol": "BTCUSDT",
  "side": "buy",
  "type": "limit",
  "price": "42000.00",
  "quantity": "0.1",
  "filled": "0.0",
  "status": "pending",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 获取模拟持仓

获取模拟交易中的所有持仓。

**接口:** `GET /api/v1/paper/positions`

**认证:** 必需

**示例:**

```bash
curl -X GET http://localhost:8080/api/v1/paper/positions \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "positions": [
    {
      "symbol": "BTCUSDT",
      "side": "long",
      "quantity": "0.5",
      "entry_price": "41500.00",
      "current_price": "42300.00",
      "unrealized_pnl": "400.00",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

## 告警接口

### 获取告警列表

获取认证用户的所有告警。

**接口:** `GET /api/v1/alerts`

**认证:** 必需

**查询参数:**

| 参数 | 类型 | 描述 |
|------|------|------|
| `status` | string | 按状态筛选：`active`, `triggered`, `disabled` |

**示例:**

```bash
curl -X GET http://localhost:8080/api/v1/alerts \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "alerts": [
    {
      "id": "alert_123",
      "symbol": "BTCUSDT",
      "condition": "price_above",
      "value": "50000.00",
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 创建告警

创建新的价格或指标告警。

**接口:** `POST /api/v1/alerts`

**认证:** 必需

**请求体:**

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `symbol` | string | 是 | 交易对 |
| `condition` | string | 是 | 告警条件 |
| `value` | number | 是 | 目标值 |
| `type` | string | 是 | 告警类型：`price`, `indicator` |
| `indicator` | string | 否* | 指标名称（*类型为 indicator 时必需） |
| `enabled` | boolean | 否 | 启用告警（默认：true） |

**支持的条件:**
- `price_above` - 价格大于目标值
- `price_below` - 价格小于目标值
- `price_crosses` - 价格穿越目标值
- `indicator_above` - 指标大于目标值
- `indicator_below` - 指标小于目标值

**示例:**

```bash
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "condition": "price_above",
    "value": 50000,
    "type": "price",
    "enabled": true
  }'
```

**响应:**

```json
{
  "id": "alert_456",
  "symbol": "BTCUSDT",
  "condition": "price_above",
  "value": "50000.00",
  "type": "price",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 删除告警

删除告警。

**接口:** `DELETE /api/v1/alerts/:id`

**认证:** 必需

**示例:**

```bash
curl -X DELETE http://localhost:8080/api/v1/alerts/alert_456 \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "success": true,
  "message": "告警删除成功"
}
```

---

## 投资组合接口

### 获取投资组合列表

获取认证用户的所有投资组合。

**接口:** `GET /api/v1/portfolios`

**认证:** 必需

**示例:**

```bash
curl -X GET http://localhost:8080/api/v1/portfolios \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "portfolios": [
    {
      "id": "portfolio_123",
      "name": "My Trading Portfolio",
      "balance": "15000.00",
      "positions": 3,
      "created_at": "2024-01-10T10:00:00Z"
    }
  ]
}
```

### 创建投资组合

创建新的投资组合。

**接口:** `POST /api/v1/portfolios`

**认证:** 必需

**请求体:**

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `name` | string | 是 | 投资组合名称 |
| `initial_balance` | string | 否 | 初始余额（默认：10000） |

**示例:**

```bash
curl -X POST http://localhost:8080/api/v1/portfolios \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Conservative Portfolio",
    "initial_balance": "50000.00"
  }'
```

**响应:**

```json
{
  "id": "portfolio_456",
  "name": "Conservative Portfolio",
  "balance": "50000.00",
  "positions": 0,
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

## 策略市场接口

### 获取策略列表

浏览市场中可用的交易策略。

**接口:** `GET /api/v1/market/strategies`

**认证:** 不需要

**查询参数:**

| 参数 | 类型 | 描述 |
|------|------|------|
| `category` | string | 按类别筛选 |
| `sort_by` | string | 排序方式：`popular`, `newest`, `price` |
| `page` | integer | 页码 |
| `limit` | integer | 每页数量 |

**示例:**

```bash
curl -X GET "http://localhost:8080/api/v1/market/strategies?sort_by=popular&limit=10"
```

**响应:**

```json
{
  "strategies": [
    {
      "id": "strategy_123",
      "name": "MA Crossover Bot",
      "description": "Moving average crossover strategy",
      "category": "trend",
      "price": "29.99",
      "subscribers": 150,
      "rating": 4.5,
      "author": "trader_pro"
    }
  ],
  "total": 50,
  "page": 1,
  "total_pages": 5
}
```

### 购买策略

订阅交易策略。

**接口:** `POST /api/v1/market/strategies/:id/purchase`

**认证:** 必需

**示例:**

```bash
curl -X POST http://localhost:8080/api/v1/market/strategies/strategy_123/purchase \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "success": true,
  "subscription_id": "sub_abc123",
  "strategy_id": "strategy_123",
  "expires_at": "2024-02-15T10:30:00Z"
}
```

---

## 回测与分析接口

### 运行回测

执行策略回测。

**接口:** `POST /api/v1/backtest`

**认证:** 必需

**请求体:**

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `strategy_id` | string | 是 | 要回测的策略 |
| `symbol` | string | 是 | 交易对 |
| `start_time` | string | 是 | 开始时间戳 |
| `end_time` | string | 是 | 结束时间戳 |
| `initial_capital` | string | 否 | 初始资金 |
| `parameters` | object | 否 | 策略参数 |

**示例:**

```bash
curl -X POST http://localhost:8080/api/v1/backtest \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "strategy_id": "strategy_123",
    "symbol": "BTCUSDT",
    "start_time": "1704067200000",
    "end_time": "1704067200000",
    "initial_capital": "10000.00",
    "parameters": {
      "fast_period": 10,
      "slow_period": 20
    }
  }'
```

**响应:**

```json
{
  "backtest_id": "bt_abc123",
  "status": "running",
  "progress": 0
}
```

### 获取投资组合分析

获取详细的投资组合绩效分析。

**接口:** `GET /api/v1/analytics/portfolio`

**认证:** 必需

**查询参数:**

| 参数 | 类型 | 描述 |
|------|------|------|
| `portfolio_id` | string | 投资组合 ID |
| `start_time` | string | 开始时间戳 |
| `end_time` | string | 结束时间戳 |

**示例:**

```bash
curl -X GET "http://localhost:8080/api/v1/analytics/portfolio?portfolio_id=portfolio_123" \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "portfolio_id": "portfolio_123",
  "total_return": "15.5%",
  "sharpe_ratio": 1.8,
  "max_drawdown": "8.2%",
  "win_rate": "65%",
  "profit_factor": 2.1,
  "total_trades": 150,
  "average_trade": "35.00",
  "daily_returns": [...]
}
```

---

## 订阅接口

### 获取订阅状态

获取当前订阅状态。

**接口:** `GET /api/v1/subscription`

**认证:** 必需

**示例:**

```bash
curl -X GET http://localhost:8080/api/v1/subscription \
  -H "Authorization: Bearer <token>"
```

**响应:**

```json
{
  "tier": "pro",
  "status": "active",
  "expires_at": "2025-01-15T10:30:00Z",
  "features": {
    "api_rate_limit": 1000,
    "concurrent_strategies": 5,
    "backtest_history": 100,
    "alerts": 50
  }
}
```

---

## 错误处理

所有接口可能返回以下格式的错误响应：

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

### 常见错误码

| 错误码 | HTTP 状态 | 描述 |
|--------|-----------|------|
| `UNAUTHORIZED` | 401 | 认证无效或缺失 |
| `FORBIDDEN` | 403 | 权限不足 |
| `NOT_FOUND` | 404 | 资源未找到 |
| `VALIDATION_ERROR` | 400 | 请求参数无效 |
| `RATE_LIMITED` | 429 | 请求过于频繁 |
| `INTERNAL_ERROR` | 500 | 服务器错误 |

---

## 限流规则

| 等级 | 每分钟请求数 |
|------|--------------|
| 免费版 | 60 |
| 专业版 | 600 |
| 企业版 | 6000 |

---

## 相关文档

- [技术架构](./backend/ARCHITECTURE.md)
- [贡献指南](./CONTRIBUTING.md)
- [快速开始](./README.md#quick-start)
