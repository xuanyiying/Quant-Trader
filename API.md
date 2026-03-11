# Quant-Trader API Documentation

Complete API reference for the Quant-Trader backend service.

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

## 📋 Endpoints Overview

| Category | Endpoints |
|----------|-----------|
| [Authentication](#authentication) | `/api/v1/register`, `/api/v1/login` |
| [Market Data](#market-data) | `/api/v1/klines/:symbol`, `/ws` |
| [Paper Trading](#paper-trading) | `/api/v1/paper/*` |
| [Alerts](#alerts) | `/api/v1/alerts` |
| [Portfolios](#portfolios) | `/api/v1/portfolios` |
| [Marketplace](#marketplace) | `/api/v1/market/strategies` |
| [Backtest & Analytics](#backtest--analytics) | `/api/v1/backtest`, `/api/v1/analytics/*` |
| [Subscription](#subscription) | `/api/v1/subscription` |

---

## 🔐 Authentication

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

---

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

## 📊 Market Data

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

---

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

## 💰 Paper Trading

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

---

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

---

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

---

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

## 🔔 Alerts

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

---

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

---

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

## 📊 Portfolios

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

---

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

## 🏪 Marketplace

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

---

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

## 📈 Backtest & Analytics

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

---

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

## 💳 Subscription

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

## ⚠️ Error Responses

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

## 📝 Rate Limits

| Tier | Requests/minute |
|------|-----------------|
| Free | 60 |
| Pro | 600 |
| Enterprise | 6000 |

---

## 🔗 Related Documentation

- [Technical Architecture](./backend/ARCHITECTURE.md)
- [Contributing Guide](./CONTRIBUTING.md)
- [Quick Start](./README.md#quick-start)
