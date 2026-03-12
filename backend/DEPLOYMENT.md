# 量化交易平台部署文档

## 目录

- [环境要求](#环境要求)
- [快速开始](#快速开始)
- [开发环境](#开发环境)
- [生产环境部署](#生产环境部署)
- [配置说明](#配置说明)
- [故障排查](#故障排查)

---

## 环境要求

### 硬件要求

| 环境 | CPU | 内存 | 磁盘 |
|------|-----|------|------|
| 开发 | 2核 | 4GB | 20GB |
| 生产 | 4核+ | 8GB+ | 100GB+ SSD |

### 软件要求

| 软件 | 版本 | 说明 |
|------|------|------|
| Docker | 24.0+ | 容器运行时 |
| Docker Compose | 2.20+ | 容器编排 |
| Go | 1.24+ | 后端开发 |
| Node.js | 20+ | 前端开发 |

---

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/your-org/quant-trader.git
cd quant-trader
```

### 2. 启动基础设施

```bash
# 启动数据库、缓存、消息队列
make dev
```

### 3. 配置环境变量

```bash
cp backend/.env.example backend/.env
vim backend/.env
```

### 4. 运行后端服务

```bash
make run
```

---

## 开发环境

### 架构说明

开发环境采用**混合模式**：
- **基础设施服务**（Docker）: TimescaleDB、Redis、NATS
- **应用服务**（本地）: 后端 Go 服务、前端 React 服务

### 启动服务

```bash
# 启动基础设施
make docker-up

# 运行后端服务
make run

# 运行前端服务（另一个终端）
cd frontend && npm run dev
```

### 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| TimescaleDB | 5432 | 时序数据库 |
| Redis | 6379 | 缓存服务 |
| NATS | 4222 | 消息队列 |
| Backend | 8080 | API 服务（本地运行） |
| Frontend | 5173 | 前端服务（本地运行） |

### 常用命令

```bash
make help          # 查看所有命令
make docker-up     # 启动基础设施
make docker-down   # 停止基础设施
make docker-logs   # 查看日志
make run           # 运行后端服务
make test          # 运行测试
make db-connect    # 连接数据库
```

---

## 生产环境部署

### 使用部署脚本

```bash
# 执行部署
./deploy/scripts/deploy.sh deploy production
```

### 生产环境服务

生产环境使用 Docker Compose 部署完整服务栈：

| 服务 | 说明 |
|------|------|
| Frontend | Nginx + React 应用 |
| Backend | Go API 服务 |
| TimescaleDB | 时序数据库 |
| Redis | 缓存服务 |
| NATS | 消息队列 |
| Prometheus | 指标采集 |
| Grafana | 监控面板 |
| Alertmanager | 告警管理 |

### 生产环境配置

1. **复制配置文件**
```bash
cp deploy/docker/.env.production.example deploy/docker/.env.production
```

2. **修改关键配置**
```bash
# 必须修改的配置
POSTGRES_PASSWORD=your-secure-password
JWT_SECRET=your-jwt-secret
GRAFANA_ADMIN_PASSWORD=your-grafana-password
```

3. **启动服务**
```bash
docker-compose -f deploy/docker/docker-compose.production.yml up -d
```

---

## 配置说明

### 环境变量

```bash
# 基础配置
ENVIRONMENT=development
PORT=8080
LOG_LEVEL=debug

# 数据库
DB_DSN=postgres://postgres:password@localhost:5432/quant_trader?sslmode=disable

# 消息队列
NATS_URL=nats://localhost:4222

# 缓存
REDIS_URL=redis://localhost:6379

# 安全
JWT_SECRET=your-secret-key
```

---

## 故障排查

### 常见问题

#### 1. 端口占用

```bash
# 检查端口
lsof -i :5432
lsof -i :8080

# 停止占用进程
kill -9 <PID>
```

#### 2. 数据库连接失败

```bash
# 检查数据库状态
make docker-ps

# 查看日志
make docker-logs

# 重启数据库
docker-compose restart timescaledb
```

#### 3. 重置环境

```bash
# 停止并删除所有数据
make docker-down
docker volume prune

# 重新启动
make docker-up
```

### 日志查看

```bash
# 查看所有日志
make docker-logs

# 查看特定服务
docker-compose logs -f timescaledb
docker-compose logs -f redis
docker-compose logs -f nats
```

---

## 联系支持

如有问题，请提交 Issue 或联系开发团队。
