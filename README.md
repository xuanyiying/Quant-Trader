# Quant-Trader

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React Version](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://react.dev/)
[![Docker Supported](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![TimescaleDB](https://img.shields.io/badge/TimescaleDB-PostgreSQL-FF5000?style=flat)](https://www.timescale.com/)
[![NATS JetStream](https://img.shields.io/badge/NATS-JetStream-27B6D4?style=flat)](https://nats.io/)
[![WebAssembly](https://img.shields.io/badge/WASM-Security-654FF0?style=flat)](https://webassembly.org/)

[English](./README.md) | [简体中文](./README-zh.md)

> Professional Algorithmic Trading Infrastructure designed for high concurrency, low latency, and institution-grade security.

<p align="center">
  <img src="https://raw.githubusercontent.com/your-repo/quant-trader/main/diagrams/architecture_overview.svg" alt="Quant-Trader Architecture" width="800"/>
</p>

---

## ⭐ Why Quant-Trader?

- **🚀 Enterprise-Grade Performance**: Sub-millisecond matching engine with in-memory order book
- **🔄 Multi-Exchange Support**: Native WebSocket integration with Binance, OKX, Bybit, Coinbase, Kraken
- **🛡️ Security-First Design**: WASM sandbox isolation for untrusted strategy execution
- **📊 Real-Time Analytics**: TimescaleDB-powered time-series data with built-in indicators
- **💰 Complete Trading Pipeline**: From market data ingestion to paper trading and strategy marketplace
- **☁️ Cloud-Native**: Docker-compose ready with NATS JetStream event bus

---

## 📂 Project Structure

```
quant-trader/
├── backend/                    # High-performance Go trading engine
│   ├── api/                   # REST API handlers (Gin)
│   ├── cmd/                   # Application entry point
│   ├── internal/
│   │   ├── connector/        # Exchange WebSocket connectors
│   │   ├── engine/            # Strategy execution engine
│   │   ├── matching/          # High-performance matching engine
│   │   ├── paper/             # Paper trading simulation
│   │   ├── processor/        # K-line aggregation
│   │   ├── risk/              # Risk management
│   │   ├── strategy/          # Trading strategies
│   │   └── storage/           # TimescaleDB persistence
│   ├── monitoring/            # Grafana dashboards
│   └── scripts/               # Database migrations
├── frontend/                  # React dashboard
│   └── src/
│       ├── components/        # UI components
│       ├── hooks/             # Custom React hooks
│       └── store/             # State management
├── diagrams/                  # Architecture diagrams
└── examples/                  # Usage examples
```

---

## 🚀 Key Features

### 1. High-Performance Market Data Pipeline

| Feature | Description |
|---------|-------------|
| **Multi-Exchange Connectors** | Native WebSocket with automatic failover, heartbeat checks, exponential backoff |
| **Micro-aggregation** | Millisecond K-line generation using in-memory sliding windows |
| **Event-Driven Architecture** | NATS JetStream for reliable async data distribution |
| **TimescaleDB Persistence** | Optimized time-series storage with automated partitioning |

### 2. Professional Trading Suite

| Feature | Description |
|---------|-------------|
| **Paper Trading Engine** | Low-latency simulation matching with multi-asset balance tracking |
| **Risk Management** | Pre-trade validation (position limits, daily loss prevention) |
| **Strategy Marketplace** | Subscription-based strategy distribution with Stripe billing |
| **Tier-based Rate Limiting** | Multi-tenant API limiting (Free/Pro/Enterprise) |

### 3. Advanced Strategy Execution

| Feature | Description |
|---------|-------------|
| **WASM Sandboxing** | Secure isolated execution via wazero |
| **Built-in Indicators** | RSI, MACD, Bollinger Bands, SMA, EMA, ATR |
| **Alert System** | Flexible rule-based notifications |
| **Backtesting** | Historical data replay with performance metrics |

---

## 🛠 Tech Stack

| Layer | Technology |
|-------|------------|
| **Backend** | Go 1.24+, Gin, NATS JetStream, TimescaleDB, Wazero |
| **Database** | PostgreSQL 15+, TimescaleDB |
| **Messaging** | NATS Server with JetStream |
| **Security** | WebAssembly (WASM), JWT Auth |
| **Payments** | Stripe API |
| **Frontend** | React 18, Vite, TypeScript, TailwindCSS, ECharts |
| **Monitoring** | Prometheus, Grafana |
| **DevOps** | Docker, Docker Compose |

---

## 🏗 System Architecture

```mermaid
graph TD
    subgraph "External Exchanges"
        Binance[Binance WS]
        OKX[OKX WS]
        Bybit[Bybit WS]
        Coinbase[Coinbase WS]
    end

    subgraph "Ingestion Layer"
        Connectors[Exchange Connectors]
        Normalizer[Symbol Normalizer]
        Raw[NATS: market.raw]
    end

    subgraph "Processing Core"
        KP[Kline Processor]
        AE[Alert Engine]
        PE[Paper Engine]
        AS[Analytics Service]
    end

    subgraph "Execution Layer"
        RM[Risk Manager]
        WR[WASM Runner]
    end

    subgraph "Persistence & API"
        Batch[Batch Saver]
        TSDB[(TimescaleDB)]
        HTTP[Gin REST API]
        WS[WebSocket Push]
    end

    subgraph "Commercial Layer"
        Stripe[Stripe Service]
        Marketplace[Strategy Marketplace]
    end

    Binance & OKX & Bybit & Coinbase --> Connectors
    Connectors --> Normalizer
    Normalizer --> Raw
    Raw --> KP & AE
    KP --> PE & Batch
    Batch --> TSDB
    PE --> RM
    RM --> WR
    HTTP --> TSDB & PE & Stripe & Marketplace
    TSDB --> AS
```

---

## 🏁 Quick Start

### Prerequisites

- **Go 1.24+**
- **Docker & Docker Compose**
- **PostgreSQL 15+** (optional, using TimescaleDB docker)

### 1. Clone the Repository

```bash
git clone https://github.com/your-repo/quant-trader.git
cd quant-trader
```

### 2. Start Infrastructure

```bash
cd backend
docker-compose up -d
```

This will start:
- **TimescaleDB** on `localhost:5432`
- **NATS JetStream** on `localhost:4222`
- **Redis** on `localhost:6379`
- **Grafana** on `localhost:3000`

### 3. Configure Environment

```bash
cd backend
cp .env.example .env
# Edit .env with your configuration
```

Example `.env`:
```env
DB_DSN=postgres://postgres:password@localhost:5432/quant_trader
NATS_URL=nats://localhost:4222
PORT=8080
MATCHING_ENABLED=true
JWT_SECRET=your-secret-key
STRIPE_SECRET_KEY=sk_test_...
```

### 4. Run the Backend

```bash
go run cmd/main.go
```

### 5. Run the Frontend

```bash
cd frontend
npm install
npm run dev
```

Access the dashboard at `http://localhost:5173`

---

## 📖 Usage Examples

### Paper Trading

```go
// Submit a paper order
order := &model.Order{
    Symbol:    "BTCUSDT",
    Side:      model.OrderSideBuy,
    Type:      model.OrderTypeLimit,
    Price:     decimal.NewFromFloat(50000.00),
    Quantity:  decimal.NewFromFloat(0.1),
}

// Submit via API
POST /api/v1/paper/orders
{
    "symbol": "BTCUSDT",
    "side": "buy",
    "type": "limit",
    "price": "50000.00",
    "quantity": "0.1"
}
```

### Subscribe to Market Data

```javascript
// WebSocket subscription
const ws = new WebSocket('ws://localhost:8080/ws/market');

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Market update:', data);
};

ws.send(JSON.stringify({
    action: 'subscribe',
    symbols: ['BTCUSDT', 'ETHUSDT']
}));
```

### Run a Strategy

```go
// Initialize strategy
strategy := &ma.CrossStrategy{
    FastPeriod: 10,
    SlowPeriod: 20,
    Symbol:     "BTCUSDT",
}

// Execute via engine
engine.Execute(strategy)
```

---

## 📊 Performance Benchmarks

| Layer | Latency (P99) | Throughput |
|-------|---------------|------------|
| **Market Ingestion** | < 2ms | 50,000 msg/s |
| **K-Line Aggregation** | < 5ms | 100 symbols/1m |
| **Order Matching** | < 10ms | 1,000 orders/s |
| **Persistence** | < 20ms | 10,000 records/batch |
| **REST API** | < 50ms | 5,000 req/s |

---

## 🔗 Related Documentation

- [Technical Architecture](./backend/ARCHITECTURE.md) - Deep dive into engine design
- [API Documentation](./API.md) - Complete API reference
- [Contributing Guide](./CONTRIBUTING.md) - How to contribute
- [Changelog](./CHANGELOG.md) - Version history

---

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](./CONTRIBUTING.md) for details.

```bash
# Fork the repo
# Create your feature branch
git checkout -b feature/amazing-feature

# Commit your changes
git commit -m 'Add some amazing feature'

# Push to the branch
git push origin feature/amazing-feature

# Open a Pull Request
```

---

## 📄 License

Distributed under the MIT License. See [LICENSE](./LICENSE) for more information.

---

## 🙏 Acknowledgments

- [Gin](https://gin-gonic.com/) - Web framework
- [NATS](https://nats.io/) - Messaging system
- [TimescaleDB](https://www.timescale.com/) - Time-series database
- [wazero](https://github.com/tetratelabs/wazero) - Go WASM runtime

---

<p align="center">
  <strong>Quant-Trader</strong> — Empowering your strategy with professional infrastructure.
  <br/>
  <a href="https://github.com/your-repo/quant-trader">Star</a> • 
  <a href="https://github.com/your-repo/quant-trader/fork">Fork</a> • 
  <a href="https://github.com/your-repo/quant-trader/issues">Issues</a>
</p>
