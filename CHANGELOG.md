# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- WebSocket support for real-time market data streaming
- Multiple exchange connectors (Binance, OKX, Bybit, Coinbase, Kraken)
- Paper trading engine with portfolio management
- Strategy marketplace with Stripe integration
- Alert system for price and indicator notifications
- WASM-based strategy execution sandbox
- Backtesting engine with performance analytics

### Improved
- Enhanced matching engine with lease-based state management
- Optimized K-line aggregation with microsecond precision
- NATS JetStream integration for reliable message distribution
- TimescaleDB persistence with automatic partitioning

---

## [0.1.0] - 2024-01-15

### Added

#### Core Engine
- High-performance matching engine with in-memory order book
- Multi-exchange WebSocket connectors with auto-reconnect
- Real-time K-line aggregation (1m, 5m, 15m, 1h, 4h, 1d)
- NATS JetStream event bus integration

#### Trading Features
- Paper trading engine with virtual balances
- Risk management module (position limits, loss prevention)
- Alert system with flexible rules
- Portfolio management

#### API & Frontend
- RESTful API with Gin framework
- JWT authentication
- Rate limiting middleware (Free/Pro/Enterprise tiers)
- React-based dashboard with real-time updates
- ECharts integration for market visualization

#### Infrastructure
- Docker Compose setup for all services
- Grafana dashboards for monitoring
- Prometheus metrics collection
- Database migrations

---

## [0.0.1] - 2023-12-01

### Added
- Initial project structure
- Basic project setup
- Docker configuration
- README documentation

---

## Roadmap

### Planned Features

#### v0.2.0 (Q2 2024)
- [ ] Advanced order types (OCO, STOP_LIMIT, TRAILING_STOP)
- [ ] Portfolio rebalancing automation
- [ ] Advanced technical indicators library
- [ ] Strategy parameter optimization

#### v0.3.0 (Q3 2024)
- [ ] Multi-account management
- [ ] Institutional-grade reporting
- [ ] Real trading integration (Binance, OKX API)
- [ ] Advanced backtesting features

#### v0.4.0 (Q4 2024)
- [ ] AI-powered strategy suggestions
- [ ] Social trading features
- [ ] Mobile app support
- [ ] Advanced risk analytics

---

## Versioning

We use [SemVer](http://semver.org/) for versioning. Available versions are listed on the [releases page](https://github.com/your-repo/quant-trader/releases).

---

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details on how to contribute.

---

## Credits

Thank you to all our contributors!

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
