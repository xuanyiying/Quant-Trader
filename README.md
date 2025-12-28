# Quant-Trader

[English](./README.md) | [简体中文](./README-zh.md)

Professional Algorithmic Trading Infrastructure designed for high concurrency, low latency, and institution-grade security.

---

## 📂 Project Structure

- [backend/](./backend/): High-performance Go trading engine, API service, and core logic.
  - [Technical Architecture](./backend/ARCHITECTURE.md): Deep dive into engine design and system topology.
- [frontend/](./frontend/): Modern React-based dashboard for real-time monitoring and strategy management.

---

## 🚀 Key Features

- **Multi-Exchange Integration**: Native WebSocket support for Binance, OKX, Bybit, etc.
- **WASM Sandboxing**: Secure and isolated strategy execution environment.
- **TimescaleDB Persistence**: Optimized storage for high-frequency time-series data.
- **Paper Trading**: Low-latency simulation matching with risk management.
- **Stripe Billing**: Integrated subscription management and strategy marketplace.

---

## 🛠 Tech Stack

- **Backend**: Go 1.24+, Gin, NATS JetStream, TimescaleDB, Wazero.
- **Frontend**: React, Vite, TailwindCSS, ECharts.

---

## 🏁 Quick Start

Please refer to the documentation in each subdirectory for detailed setup instructions.

- [Backend Setup](./backend/README.md)
- [Frontend Setup](./frontend/README.md)
