# ============================================================
# 量化交易平台 Makefile
# ============================================================

.PHONY: help dev build test clean docker-up docker-down docker-logs

# 默认目标
help:
	@echo "量化交易平台 - 可用命令:"
	@echo ""
	@echo "  开发命令:"
	@echo "    make dev          启动基础设施服务"
	@echo "    make run          本地运行后端服务"
	@echo "    make build        构建应用"
	@echo "    make test         运行测试"
	@echo "    make clean        清理构建产物"
	@echo ""
	@echo "  Docker 命令:"
	@echo "    make docker-up    启动基础设施"
	@echo "    make docker-down  停止所有服务"
	@echo "    make docker-logs  查看服务日志"
	@echo "    make docker-ps    查看服务状态"
	@echo ""
	@echo "  部署命令:"
	@echo "    make deploy       部署到生产环境"
	@echo "    make rollback     回滚到上一版本"
	@echo ""
	@echo "  数据库命令:"
	@echo "    make db-migrate   运行数据库迁移"
	@echo "    make db-reset     重置数据库"
	@echo "    make db-backup    备份数据库"
	@echo "    make db-connect   连接数据库"

# ==================== 开发命令 ====================

# 启动基础设施服务（数据库、缓存、消息队列）
dev: docker-up

# 本地运行后端服务
run:
	cd backend && go run cmd/main.go

# 构建应用
build:
	cd backend && go build -o bin/server ./cmd/main.go
	cd frontend && npm run build

# 运行测试
test:
	cd backend && go test -v ./...
	cd frontend && npm test

# 清理构建产物
clean:
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf backend/tmp
	docker system prune -f

# ==================== Docker 命令 ====================

# 启动基础设施
docker-up:
	cd backend && docker-compose up -d
	@echo ""
	@echo "基础设施已启动:"
	@echo "  TimescaleDB: localhost:5432"
	@echo "  Redis: localhost:6379"
	@echo "  NATS: localhost:4222"
	@echo ""
	@echo "运行后端服务: make run"

# 停止所有服务
docker-down:
	cd backend && docker-compose down

# 查看服务日志
docker-logs:
	cd backend && docker-compose logs -f

# 查看服务状态
docker-ps:
	cd backend && docker-compose ps

# 重启基础设施
docker-restart:
	cd backend && docker-compose restart

# ==================== 部署命令 ====================

# 部署到生产环境
deploy:
	./deploy/scripts/deploy.sh deploy production

# 回滚到上一版本
rollback:
	./deploy/scripts/deploy.sh rollback

# 健康检查
health:
	./deploy/scripts/deploy.sh health

# ==================== 数据库命令 ====================

# 运行数据库迁移
db-migrate:
	cd backend && docker-compose exec timescaledb bash -c 'for f in /docker-entrypoint-initdb.d/migrations/*.sql; do echo "Applying $$f..."; psql -U postgres -d quant_trader -f "$$f"; done'

# 重置数据库
db-reset:
	cd backend && docker-compose down -v
	cd backend && docker-compose up -d timescaledb

# 备份数据库
db-backup:
	mkdir -p backups
	cd backend && docker-compose exec timescaledb pg_dump -U postgres quant_trader > ../backups/backup_$$(date +%Y%m%d_%H%M%S).sql

# 连接数据库
db-connect:
	cd backend && docker-compose exec timescaledb psql -U postgres -d quant_trader

# ==================== 代码质量 ====================

# 代码格式化
fmt:
	cd backend && go fmt ./...
	cd frontend && npm run lint --fix

# 静态检查
lint:
	cd backend && go vet ./...
	cd frontend && npm run lint

# ==================== 工具命令 ====================

# 安装开发工具
install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/go-delve/delve/cmd/dlv@latest

# 查看依赖
deps:
	cd backend && go mod graph
	cd frontend && npm list --depth=0
