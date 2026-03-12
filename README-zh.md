# Quant-Trader

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React Version](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://react.dev/)
[![Docker Supported](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![TimescaleDB](https://img.shields.io/badge/TimescaleDB-PostgreSQL-FF5000?style=flat)](https://www.timescale.com/)
[![NATS JetStream](https://img.shields.io/badge/NATS-JetStream-27B6D4?style=flat)](https://nats.io/)
[![WebAssembly](https://img.shields.io/badge/WASM-Security-654FF0?style=flat)](https://webassembly.org/)

[English](./README.md) | [简体中文](./README-zh.md)

> 专为高并发、低延迟和机构级安全设计的专业量化交易基础设施。

<p align="center">
  <img src="https://raw.githubusercontent.com/your-repo/quant-trader/main/diagrams/architecture_overview.svg" alt="Quant-Trader 架构图" width="800"/>
</p>

---

## 📋 目录

- [项目概述](#项目概述)
- [核心特性](#核心特性)
- [系统架构](#系统架构)
- [快速开始](#快速开始)
- [安装指南](#安装指南)
- [配置说明](#配置说明)
- [使用示例](#使用示例)
- [API 文档](#api-文档)
- [性能基准](#性能基准)
- [故障排查](#故障排查)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

---

## 项目概述

Quant-Trader 是一个生产级的量化交易平台，结合了高性能的市场数据处理、复杂的交易模拟和企业级的风险管理。基于 Go 和 React 构建，为量化交易策略从开发到部署提供完整的基础设施。

### 为什么选择 Quant-Trader？

- **🚀 企业级性能**: 亚毫秒级撮合引擎，采用内存订单簿
- **🔄 多交易所支持**: 原生 WebSocket 接入 Binance、OKX、Bybit、Coinbase、Kraken
- **🛡️ 安全优先设计**: WASM 沙箱隔离，确保不可信策略的安全执行
- **📊 实时分析**: TimescaleDB 驱动的时序数据，内置技术指标库
- **💰 完整交易流程**: 从市场数据摄取到模拟交易和策略市场
- **☁️ 云原生架构**: Docker-compose 部署就绪，NATS JetStream 事件总线

---

## 📸 运行截图

### 交易仪表盘
<p align="center">
  <img src="https://raw.githubusercontent.com/xuanyiying/quant-trader/main/assets/screenshot.png" alt="Quant-Trader 交易仪表盘" width="900"/>
</p>

### K线图表
<p align="center">
  <img src="https://raw.githubusercontent.com/xuanyiying/quant-trader/main/assets/chart.png" alt="Quant-Trader K线图表" width="900"/>
</p>

---

## 核心特性

### 1. 高性能市场数据管道

| 特性 | 描述 |
|------|------|
| **多交易所连接器** | 原生 WebSocket，自动故障转移，心跳检测，指数退避重连 |
| **毫秒级聚合** | 使用内存滑动窗口实现毫秒级 K线生成 |
| **事件驱动架构** | NATS JetStream 实现可靠的异步数据分发 |
| **TimescaleDB 持久化** | 优化的时序数据存储，自动分区和压缩 |

### 2. 专业交易套件

| 特性 | 描述 |
|------|------|
| **模拟交易引擎** | 低延迟模拟撮合，支持多资产余额追踪 |
| **风险管理** | 交易前验证（仓位限制、单日亏损防护） |
| **策略市场** | 基于订阅的策略分发，集成 Stripe 支付 |
| **分级限流** | 多租户 API 限流（免费版/专业版/企业版） |

### 3. 高级策略执行

| 特性 | 描述 |
|------|------|
| **WASM 沙箱** | 通过 wazero 实现安全隔离执行 |
| **内置指标库** | RSI、MACD、布林带、SMA、EMA、ATR |
| **告警系统** | 灵活基于规则的通知引擎 |
| **回测支持** | 历史数据回放与性能指标分析 |

---

## 系统架构

```mermaid
graph TD
    subgraph "外部交易所"
        Binance[币安 WS]
        OKX[OKX WS]
        Bybit[Bybit WS]
        Coinbase[Coinbase WS]
    end

    subgraph "接入层"
        Connectors[交易所连接器]
        Normalizer[交易对归一化]
        Raw[NATS: market.raw]
    end

    subgraph "处理核心"
        KP[K线处理器]
        AE[告警引擎]
        PE[模拟交易引擎]
        AS[分析服务]
    end

    subgraph "执行层"
        RM[风险管理器]
        WR[WASM 运行器]
    end

    subgraph "持久化与API"
        Batch[批量存储]
        TSDB[(TimescaleDB)]
        HTTP[Gin REST API]
        WS[WebSocket 推送]
    end

    subgraph "商业层"
        Stripe[Stripe 支付]
        Marketplace[策略市场]
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

### 项目结构

```
quant-trader/
├── backend/                    # 高性能 Go 交易引擎
│   ├── api/                   # REST API 处理器 (Gin)
│   ├── cmd/                   # 应用程序入口
│   ├── internal/
│   │   ├── connector/        # 交易所 WebSocket 连接器
│   │   ├── engine/            # 策略执行引擎
│   │   ├── matching/          # 高性能撮合引擎
│   │   ├── paper/             # 模拟交易引擎
│   │   ├── processor/        # K线聚合处理器
│   │   ├── risk/              # 风险管理
│   │   ├── strategy/          # 交易策略
│   │   └── storage/           # TimescaleDB 持久化
│   ├── monitoring/            # Grafana 监控面板
│   └── scripts/               # 数据库迁移脚本
├── frontend/                  # React 仪表盘
│   └── src/
│       ├── components/        # UI 组件
│       ├── hooks/             # 自定义 React Hooks
│       └── store/             # 状态管理
├── diagrams/                  # 架构图
└── examples/                  # 使用示例
```

---

## 快速开始

### 前置要求

- **Go 1.24+**
- **Docker & Docker Compose**
- **Node.js 20+**（前端开发）
- **PostgreSQL 15+**（可选，使用 TimescaleDB Docker）

### 1. 克隆仓库

```bash
git clone https://github.com/your-repo/quant-trader.git
cd quant-trader
```

### 2. 启动基础设施

```bash
make dev
```

这将启动以下服务：
- **TimescaleDB** 在 `localhost:5432`
- **NATS JetStream** 在 `localhost:4222`
- **Redis** 在 `localhost:6379`
- **Grafana** 在 `localhost:3000`

### 3. 配置环境变量

```bash
cd backend
cp .env.example .env
# 编辑 .env 文件进行配置
```

### 4. 启动后端服务

```bash
go run cmd/main.go
```

### 5. 启动前端服务

```bash
cd frontend
npm install
npm run dev
```

访问仪表盘：`http://localhost:5173`

---

## 安装指南

### 开发环境

#### 后端设置

```bash
# 克隆并进入后端目录
cd quant-trader/backend

# 安装 Go 依赖
go mod download

# 运行数据库迁移
make db-migrate

# 启动服务
go run cmd/main.go
```

#### 前端设置

```bash
cd quant-trader/frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 生产部署

详见 [DEPLOYMENT.md](./backend/DEPLOYMENT.md) 了解详细的生产部署说明。

```bash
# 使用部署脚本部署
./deploy/scripts/deploy.sh deploy production
```

---

## 配置说明

### 环境变量

在 `backend` 目录创建 `.env` 文件：

```env
# 基础配置
PORT=8080
LOG_LEVEL=info

# 数据库
DB_DSN=postgres://postgres:password@localhost:5432/quant_trader

# 消息队列
NATS_URL=nats://localhost:4222

# 撮合引擎
MATCHING_ENABLED=true
MATCHING_LEASE_TTL=20s
MATCHING_LEASE_REFRESH=5s
MATCHING_SNAPSHOT_EVERY=2s

# 安全
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# 支付（Stripe）
STRIPE_API_KEY=sk_test_your_stripe_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here

# 监控
ENABLE_METRICS=true
METRICS_PORT=9090

# 通知（可选）
DINGTALK_WEBHOOK=https://oapi.dingtalk.com/robot/send?access_token=xxx
WECHAT_WEBHOOK=https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx

# AI 集成（可选）
OPENAI_API_KEY=sk-your-openai-api-key
ANTHROPIC_API_KEY=sk-ant-your-claude-api-key
```

详见 [.env.example](./backend/.env.example) 获取完整配置模板。

### 配置指南

1. **安全性**: 生产环境务必为 `JWT_SECRET` 使用强随机字符串
2. **数据库**: 生产环境数据库访问使用 SSL/TLS 连接
3. **支付**: 开发环境使用测试密钥，生产环境才使用正式密钥
4. **监控**: 启用指标并设置适当的告警

---

## 使用示例

### 模拟交易

#### 提交模拟订单

```go
// 创建订单
order := &model.Order{
    Symbol:    "BTCUSDT",
    Side:      model.OrderSideBuy,
    Type:      model.OrderTypeLimit,
    Price:     decimal.NewFromFloat(50000.00),
    Quantity:  decimal.NewFromFloat(0.1),
}

// 通过 API 提交
POST /api/v1/paper/orders
{
    "symbol": "BTCUSDT",
    "side": "buy",
    "type": "limit",
    "price": "50000.00",
    "quantity": "0.1"
}
```

#### 使用 cURL

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

### 订阅市场数据

#### WebSocket 连接

```javascript
// WebSocket 订阅
const ws = new WebSocket('ws://localhost:8080/ws/market');

ws.onopen = () => {
  console.log('已连接到市场数据流');
  ws.send(JSON.stringify({
    action: 'subscribe',
    symbols: ['BTCUSDT', 'ETHUSDT']
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('市场更新:', data);
};

ws.onerror = (error) => {
  console.error('WebSocket 错误:', error);
};

ws.onclose = () => {
  console.log('已断开市场数据流连接');
};
```

### 运行策略

```go
// 初始化策略
strategy := &ma.CrossStrategy{
    FastPeriod: 10,
    SlowPeriod: 20,
    Symbol:     "BTCUSDT",
}

// 通过引擎执行
engine.Execute(strategy)
```

### 创建告警

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

## API 文档

完整的 API 文档请参见 [API.md](./API.md)。

### 认证方式

大部分接口需要 JWT 认证：

```
Authorization: Bearer <your-jwt-token>
```

### 基础 URL

```
http://localhost:8080
```

### 可用接口

| 类别 | 接口 |
|------|------|
| 认证 | `/api/v1/register`, `/api/v1/login` |
| 市场数据 | `/api/v1/klines/:symbol`, `/ws` |
| 模拟交易 | `/api/v1/paper/*` |
| 告警 | `/api/v1/alerts` |
| 投资组合 | `/api/v1/portfolios` |
| 策略市场 | `/api/v1/market/strategies` |
| 回测 | `/api/v1/backtest`, `/api/v1/analytics/*` |
| 订阅 | `/api/v1/subscription` |

---

## 性能基准

| 层级 | 延迟 (P99) | 吞吐量 |
|------|------------|--------|
| **市场数据摄取** | < 2ms | 50,000 消息/秒 |
| **K线聚合** | < 5ms | 100 交易对/分钟 |
| **订单撮合** | < 10ms | 1,000 订单/秒 |
| **数据持久化** | < 20ms | 10,000 记录/批 |
| **REST API** | < 50ms | 5,000 请求/秒 |

---

## 故障排查

### 常见问题

#### 1. 端口被占用

```bash
# 检查占用端口的进程
lsof -i :8080

# 终止进程
kill -9 <PID>
```

#### 2. 数据库连接失败

```bash
# 检查 TimescaleDB 是否运行
docker-compose ps

# 查看日志
docker-compose logs timescaledb

# 重启数据库
docker-compose restart timescaledb
```

#### 3. NATS 连接问题

```bash
# 检查 NATS 状态
docker-compose logs nats

# 重启 NATS
docker-compose restart nats
```

#### 4. 前端构建错误

```bash
cd frontend

# 清除 npm 缓存
npm cache clean --force

# 删除 node_modules 并重新安装
rm -rf node_modules package-lock.json
npm install
```

### 重置开发环境

```bash
# 停止所有服务并删除卷
make docker-down
docker volume prune

# 重新启动
make dev
```

### 获取帮助

- **GitHub Issues**: [报告问题或请求功能](https://github.com/your-repo/quant-trader/issues)
- **文档**: 查看 [ARCHITECTURE.md](./backend/ARCHITECTURE-zh.md) 了解技术细节
- **示例**: 参见 [examples/](./examples/) 目录获取使用示例

---

## 技术栈

| 层级 | 技术 |
|------|------|
| **后端** | Go 1.24+, Gin, NATS JetStream, TimescaleDB, Wazero |
| **数据库** | PostgreSQL 15+, TimescaleDB |
| **消息队列** | NATS Server + JetStream |
| **安全** | WebAssembly (WASM), JWT 认证 |
| **支付** | Stripe API |
| **前端** | React 18, Vite, TypeScript, TailwindCSS, ECharts |
| **监控** | Prometheus, Grafana |
| **DevOps** | Docker, Docker Compose |

---

## 贡献指南

我们欢迎贡献！请阅读我们的 [贡献指南](./CONTRIBUTING.md) 了解详情。

```bash
# Fork 本仓库
# 创建功能分支
git checkout -b feature/amazing-feature

# 提交更改
git commit -m 'feat: 添加某个神奇的功能'

# 推送分支
git push origin feature/amazing-feature

# 打开 Pull Request
```

### 开发工作流

1. **设置**: 按照 [安装指南](#安装指南) 进行配置
2. **编码**: 遵循我们的编码规范（参见 [CONTRIBUTING.md](./CONTRIBUTING.md)）
3. **测试**: 使用 `make test` 运行测试
4. **文档**: 为任何更改更新文档
5. **提交**: 创建带有清晰描述的 Pull Request

---

## 相关文档

- [技术架构文档](./backend/ARCHITECTURE-zh.md) - 深入了解引擎设计
- [API 文档](./API.md) - 完整的 API 参考
- [部署指南](./backend/DEPLOYMENT.md) - 生产部署说明
- [贡献指南](./CONTRIBUTING.md) - 如何贡献代码
- [更新日志](./CHANGELOG.md) - 版本历史

---

## 许可证

基于 MIT 许可证分发。详见 [LICENSE](./LICENSE)。

---

## 致谢

- [Gin](https://gin-gonic.com/) - Web 框架
- [NATS](https://nats.io/) - 消息系统
- [TimescaleDB](https://www.timescale.com/) - 时序数据库
- [wazero](https://github.com/tetratelabs/wazero) - Go WASM 运行时
- [React](https://react.dev/) - 前端库
- [ECharts](https://echarts.apache.org/) - 图表库

---

<p align="center">
  <strong>Quant-Trader</strong> — 用专业基础设施赋能您的交易策略。
  <br/>
  <a href="https://github.com/your-repo/quant-trader">Star</a> • 
  <a href="https://github.com/your-repo/quant-trader/fork">Fork</a> • 
  <a href="https://github.com/your-repo/quant-trader/issues">Issues</a>
</p>
