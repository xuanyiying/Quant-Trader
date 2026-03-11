# Quant-Trader Backend

High-performance algorithmic trading engine built with Go, featuring real-time market data ingestion, paper trading simulation, and strategy execution.

---

## 🚀 Key Features

### 1. High-Performance Market Data Pipeline

- **Exchange Connectors**: Native WebSocket integration with Binance, OKX, Bybit, Coinbase, Kraken with automatic failover
- **Micro-aggregation**: Millisecond K-line generation using in-memory window management
- **Event-Driven Architecture**: Powered by NATS JetStream for reliable async data distribution
- **TimescaleDB Persistence**: Optimized time-series storage with automated partitioning and compression

### 2. Commercial Trading Suite

- **Paper Trading Engine**: Low-latency simulation matching with multi-asset balance tracking
- **Risk Management**: Pre-trade risk validation engine (position limits, excessive loss prevention)
- **Strategy Marketplace**: Subscription-based strategy distribution with integrated Stripe billing
- **Tier-based Rate Limiting**: Multi-tenant API rate limiting (Free/Pro/Enterprise tiers)

### 3. Advanced Strategy Implementation

- **WASM Sandboxing**: Secure, isolated strategy execution using `wazero` (WebAssembly for Go)
- **Universal Indicators**: Built-in library for RSI, MACD, Bollinger Bands, and more
- **Alert System**: Flexible rule-based notification engine for price and technical triggers
- **Backtesting Engine**: Historical data replay with performance metrics

---

## 🛠 Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.24+ |
| Web Framework | Gin |
| Database | PostgreSQL 15+, TimescaleDB |
| Messaging | NATS JetStream |
| Security | WebAssembly (WASM), JWT |
| Payments | Stripe API |
| Monitoring | Prometheus, Grafana |

---

## 🏗 System Architecture

```mermaid
graph TD
    subgraph "External Exchanges"
        Binance[Binance WS]
        OKX[OKX WS]
    end

    subgraph "Ingestor (Go)"
        Connectors[Exchange Connectors]
        Raw[NATS: market.raw]
    end

    subgraph "Core Engine"
        KP[Kline Processor]
        AE[Alert Engine]
        PE[Paper Engine]
        WR[Wasm Runner]
    end

    subgraph "Storage & API"
        TSDB[(TimescaleDB)]
        HTTP[Gin REST API]
        WS[Push Gateway]
    end

    Binance & OKX --> Connectors
    Connectors --> Raw
    Raw --> KP & AE
    KP --> PE & TSDB & WS
    PE --> WR
    HTTP --> TSDB & PE
```

---

## 🏁 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- NATS Server

### 1. Installation

```bash
git clone https://github.com/your-repo/quant-trader.git
cd quant-trader/backend
go mod download
```

### 2. Configuration

Copy `config.yaml.example` to `config.yaml` and configure your credentials:

```yaml
database:
  url: "postgres://user:pass@localhost:5432/quant_trader"
nats:
  url: "nats://localhost:4222"
stripe:
  key: "sk_test_..."
app:
  port: "8080"
```

Or use environment variables:

```bash
export DB_DSN="postgres://postgres:password@localhost:5432/quant_trader"
export NATS_URL="nats://localhost:4222"
export PORT="8080"
export JWT_SECRET="your-secret-key"
```

### 3. Running the System

```bash
# Start infrastructure
docker-compose up -d

# Start backend
go run cmd/main.go
```

The API server will start on `http://localhost:8080`

---

## 📡 API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/register` | Register new user |
| POST | `/api/v1/login` | User login |

### Market Data

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/klines/:symbol` | Get historical K-lines |
| GET | `/ws` | WebSocket for real-time updates |

### Paper Trading

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/paper/account` | Get paper trading account |
| POST | `/api/v1/paper/account/reset` | Reset paper account |
| POST | `/api/v1/paper/orders` | Create paper order |
| GET | `/api/v1/paper/positions` | Get open positions |

### Alerts

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/alerts` | List user alerts |
| POST | `/api/v1/alerts` | Create new alert |
| DELETE | `/api/v1/alerts/:id` | Delete alert |

### Portfolios

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/portfolios` | List portfolios |
| POST | `/api/v1/portfolios` | Create portfolio |

### Marketplace

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/market/strategies` | List available strategies |
| POST | `/api/v1/market/strategies/:id/purchase` | Purchase strategy |

### Backtest & Analytics

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/backtest` | Run backtest |
| POST | `/api/v1/backfill` | Trigger historical data backfill |
| GET | `/api/v1/analytics/portfolio` | Get portfolio analytics |

### Subscription

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/subscription` | Get subscription info |

---

## 📖 API Examples

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123",
    "email": "user@example.com"
  }
}
```

### Get K-Lines

```bash
curl -X GET "http://localhost:8080/api/v1/klines/BTCUSDT?interval=1m&limit=100"
```

### Create Paper Order

```bash
curl -X POST http://localhost:8080/api/v1/paper/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "side": "buy",
    "type": "limit",
    "price": "50000.00",
    "quantity": "0.1"
  }'
```

### WebSocket Subscription

```javascript
const ws = new WebSocket('http://localhost:8080/ws');

ws.onopen = () => {
  ws.send(JSON.stringify({
    action: 'subscribe',
    symbols: ['BTCUSDT', 'ETHUSDT']
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Market update:', data);
};
```

---

## 🔧 Configuration Options

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_DSN` | - | PostgreSQL connection string |
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `JWT_SECRET` | - | JWT signing key |
| `STRIPE_SECRET_KEY` | - | Stripe API key |
| `MATCHING_ENABLED` | `true` | Enable matching engine |
| `MATCHING_LEASE_TTL` | `20s` | Order lease TTL |
| `MATCHING_SNAPSHOT_EVERY` | `2s` | Snapshot interval |

---

## 📊 Performance Benchmark

| Layer | Latency (P99) | Throughput |
| :--- | :--- | :--- |
| **Ingestion** | < 2ms | 50,000 trades/s |
| **Aggregation** | < 5ms | 100 symbols (1m period) |
| **Simulation** | < 10ms | 1,000 orders/s |
| **Persistence** | < 20ms | 10,000 records/batch |

---

## 📁 Project Structure

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── api/
│   ├── middleware/          # Auth, rate limiting, subscription
│   ├── auth_handler.go      # Authentication handlers
│   ├── kline_handler.go     # K-line handlers
│   ├── market_handler.go    # Market data handlers
│   ├── paper_handler.go    # Paper trading handlers
│   └── ...
├── internal/
│   ├── connector/           # Exchange WebSocket connectors
│   │   ├── binance.go
│   │   ├── okx.go
│   │   ├── bybit.go
│   │   └── ...
│   ├── engine/              # Strategy execution engine
│   ├── matching/            # High-performance matching engine
│   ├── paper/               # Paper trading engine
│   ├── processor/           # K-line aggregation
│   ├── risk/                # Risk management
│   ├── strategy/            # Trading strategies
│   ├── storage/             # TimescaleDB persistence
│   └── config/              # Configuration
├── monitoring/
│   └── grafana/             # Grafana dashboards
└── scripts/
    └── migrations/           # Database migrations
```

---

## 🔨 Development

### Run Tests

```bash
go test ./...
```

### Run with Custom Config

```bash
go run cmd/main.go -config /path/to/config.yaml
```

### Database Migrations

```bash
cd scripts
./apply_all.sh
```

---

## ⚖️ License

Distributed under the MIT License. See `LICENSE` for more information.

---

*Quant-Trader - Empowering your strategy with professional infrastructure.*
