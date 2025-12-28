# Quant-Trader

[English](./README.md) | [简体中文](./README-zh.md)

专为高并发、低延迟和机构级安全设计的专业量化交易基础设施。

---

## 📂 项目结构

- [backend/](./backend/): 高性能 Go 交易引擎、API 服务及核心逻辑。
  - [技术架构文档](./backend/ARCHITECTURE-zh.md): 深入了解引擎设计与系统拓扑。
- [frontend/](./frontend/): 基于 React 的现代仪表盘，用于实时监控和策略管理。

---

## 🚀 核心特性

- **多交易所集成**: 原生支持 Binance, OKX, Bybit 等交易所的 WebSocket 接入。
- **WASM 沙箱隔离**: 安全、隔离的策略执行环境。
- **TimescaleDB 存储**: 针对高频时序数据优化的存储方案。
- **模拟交易 (Paper Trading)**: 具备风险管理的低延迟模拟撮合系统。
- **Stripe 支付**: 集成订阅管理和策略市场。

---

## 🛠 技术栈

- **后端**: Go 1.24+, Gin, NATS JetStream, TimescaleDB, Wazero.
- **前端**: React, Vite, TailwindCSS, ECharts.

---

## 🏁 快速开始

详细的安装和运行步骤请参考各子目录下的文档：

- [后端配置](./backend/README.md)
- [前端配置](./frontend/README.md)
