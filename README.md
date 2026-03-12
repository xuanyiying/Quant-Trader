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

## 📋 Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [System Architecture](#system-architecture)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage Examples](#usage-examples)
- [API Documentation](#api-documentation)
- [Performance Benchmarks](#performance-benchmarks)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

Quant-Trader is a production-grade algorithmic trading platform that combines high-performance market data processing, sophisticated trading simulation, and enterprise-level risk management. Built with Go and React, it provides a complete infrastructure for quantitative trading strategies from development to deployment.

### Why Quant-Trader?

- **🚀 Enterprise-Grade Performance**: Sub-millisecond matching engine with in-memory order book
- **🔄 Multi-Exchange Support**: Native WebSocket integration with Binance, OKX, Bybit, Coinbase, Kraken
- **🛡️ Security-First Design**: WASM sandbox isolation for untrusted strategy execution
- **📊 Real-Time Analytics**: TimescaleDB-powered time-series data with built-in indicators
- **💰 Complete Trading Pipeline**: From market data ingestion to paper trading and strategy marketplace
- **☁️ Cloud-Native**: Docker-compose ready with NATS JetStream event bus

---

## 📸 Screenshots

### Trading Dashboard
<p align="center">
  <img src="https://raw.githubusercontent.com/xuanyiying/quant-trader/main/assets/screenshot.png" alt="Quant-Trader Dashboard" width="900"/>
</p>

### K-Line Chart
<p align="center">
  <img src="https://raw.githubusercontent.com/xuanyiying/quant-trader/main/assets/chart.png" alt="Quant-Trader K-Line Chart" width="900"/>
</p>

---

## Key Features

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

## System Architecture

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

### Project Structure

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

## Quick Start

### Prerequisites

- **Go 1.24+**
- **Docker & Docker Compose**
- **Node.js 20+** (for frontend)
- **PostgreSQL 15+** (optional, using TimescaleDB docker)

### 1. Clone the Repository

```bash
git clone https://github.com/your-repo/quant-trader.git
cd quant-trader
```

### 2. Start Infrastructure

```bash
make dev
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

## Installation

### Development Environment

#### Backend Setup

```bash
# Clone and navigate to backend
cd quant-trader/backend

# Install Go dependencies
go mod download

# Run database migrations
make db-migrate

# Start the server
go run cmd/main.go
```

#### Frontend Setup

```bash
cd quant-trader/frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

### Production Deployment

See [DEPLOYMENT.md](./backend/DEPLOYMENT.md) for detailed production deployment instructions.

```bash
# Deploy using the deployment script
./deploy/scripts/deploy.sh deploy production
```

---

## Configuration

### Environment Variables

Create a `.env` file in the `backend` directory:

```env
# Basic Configuration
PORT=8080
LOG_LEVEL=info

# Database
DB_DSN=postgres://postgres:password@localhost:5432/quant_trader

# Message Queue
NATS_URL=nats://localhost:4222

# Matching Engine
MATCHING_ENABLED=true
MATCHING_LEASE_TTL=20s
MATCHING_LEASE_REFRESH=5s
MATCHING_SNAPSHOT_EVERY=2s

# Security
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Payments (Stripe)
STRIPE_API_KEY=sk_test_your_stripe_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here

# Monitoring
ENABLE_METRICS=true
METRICS_PORT=9090

# Notifications (Optional)
DINGTALK_WEBHOOK=https://oapi.dingtalk.com/robot/send?access_token=xxx
WECHAT_WEBHOOK=https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx

# AI Integration (Optional)
OPENAI_API_KEY=sk-your-openai-api-key
ANTHROPIC_API_KEY=sk-ant-your-claude-api-key
```

See [.env.example](./backend/.env.example) for the complete configuration template.

### Configuration Guidelines

1. **Security**: Always use strong, unique secrets for `JWT_SECRET` in production
2. **Database**: Use SSL/TLS connections for production database access
3. **Payments**: Use test keys for development, live keys only for production
4. **Monitoring**: Enable metrics and set up proper alerting

---

## Usage Examples

### Paper Trading

#### Submit a Paper Order

```go
// Create order
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

#### Using cURL

```bash
curl -X POST http://localhost:8080/api/v1/paper/orders \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "side": "buy",
    "type": "limit",
    "price": "50000.00",
    "quantity": "0.1"
  }'
```

### Subscribe to Market Data

#### WebSocket Connection

```javascript
// WebSocket subscription
const ws = new WebSocket('ws://localhost:8080/ws/market');

ws.onopen = () => {
  console.log('Connected to market data stream');
  ws.send(JSON.stringify({
    action: 'subscribe',
    symbols: ['BTCUSDT', 'ETHUSDT']
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Market update:', data);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected from market data stream');
};
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

### Create an Alert

```bash
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "condition": "price_above",
    "value": 50000,
    "type": "price",
    "enabled": true
  }'
```

---

## API Documentation

Complete API documentation is available in [API.md](./API.md).

### Authentication

Most endpoints require JWT authentication:

```
Authorization: Bearer <your-jwt-token>
```

### Base URL

```
http://localhost:8080
```

### Available Endpoints

| Category | Endpoints |
|----------|-----------|
| Authentication | `/api/v1/register`, `/api/v1/login` |
| Market Data | `/api/v1/klines/:symbol`, `/ws` |
| Paper Trading | `/api/v1/paper/*` |
| Alerts | `/api/v1/alerts` |
| Portfolios | `/api/v1/portfolios` |
| Marketplace | `/api/v1/market/strategies` |
| Backtest | `/api/v1/backtest`, `/api/v1/analytics/*` |
| Subscription | `/api/v1/subscription` |

---

## Performance Benchmarks

| Layer | Latency (P99) | Throughput |
|-------|---------------|------------|
| **Market Ingestion** | < 2ms | 50,000 msg/s |
| **K-Line Aggregation** | < 5ms | 100 symbols/1m |
| **Order Matching** | < 10ms | 1,000 orders/s |
| **Persistence** | < 20ms | 10,000 records/batch |
| **REST API** | < 50ms | 5,000 req/s |

---

## Troubleshooting

### Common Issues

#### 1. Port Already in Use

```bash
# Check which process is using the port
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### 2. Database Connection Failed

```bash
# Check if TimescaleDB is running
docker-compose ps

# View logs
docker-compose logs timescaledb

# Restart the database
docker-compose restart timescaledb
```

#### 3. NATS Connection Issues

```bash
# Check NATS status
docker-compose logs nats

# Restart NATS
docker-compose restart nats
```

#### 4. Frontend Build Errors

```bash
cd frontend

# Clear npm cache
npm cache clean --force

# Delete node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Reset Development Environment

```bash
# Stop all services and remove volumes
make docker-down
docker volume prune

# Restart fresh
make dev
```

### Getting Help

- **GitHub Issues**: [Report bugs or request features](https://github.com/your-repo/quant-trader/issues)
- **Documentation**: Check [ARCHITECTURE.md](./backend/ARCHITECTURE.md) for technical details
- **Examples**: See [examples/](./examples/) directory for usage examples

---

## Tech Stack

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

## Contributing

We welcome contributions! Please read our [Contributing Guide](./CONTRIBUTING.md) for details.

```bash
# Fork the repo
# Create your feature branch
git checkout -b feature/amazing-feature

# Commit your changes
git commit -m 'feat: add some amazing feature'

# Push to the branch
git push origin feature/amazing-feature

# Open a Pull Request
```

### Development Workflow

1. **Setup**: Follow the [Installation](#installation) guide
2. **Code**: Follow our coding standards (see [CONTRIBUTING.md](./CONTRIBUTING.md))
3. **Test**: Run tests with `make test`
4. **Document**: Update documentation for any changes
5. **Submit**: Create a pull request with a clear description

---

## Related Documentation

- [Technical Architecture](./backend/ARCHITECTURE.md) - Deep dive into engine design
- [API Documentation](./API.md) - Complete API reference
- [Deployment Guide](./backend/DEPLOYMENT.md) - Production deployment instructions
- [Contributing Guide](./CONTRIBUTING.md) - How to contribute
- [Changelog](./CHANGELOG.md) - Version history

---

## License

Distributed under the MIT License. See [LICENSE](./LICENSE) for more information.

---

## Acknowledgments

- [Gin](https://gin-gonic.com/) - Web framework
- [NATS](https://nats.io/) - Messaging system
- [TimescaleDB](https://www.timescale.com/) - Time-series database
- [wazero](https://github.com/tetratelabs/wazero) - Go WASM runtime
- [React](https://react.dev/) - Frontend library
- [ECharts](https://echarts.apache.org/) - Charting library

---

<p align="center">
  <strong>Quant-Trader</strong> — Empowering your strategy with professional infrastructure.
  <br/>
  <a href="https://github.com/your-repo/quant-trader">Star</a> • 
  <a href="https://github.com/your-repo/quant-trader/fork">Fork</a> • 
  <a href="https://github.com/your-repo/quant-trader/issues">Issues</a>
</p>
