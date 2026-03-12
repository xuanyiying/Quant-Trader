#!/bin/bash

# ============================================================
# 量化交易平台自动化部署脚本
# 功能: 代码拉取、依赖安装、构建打包、环境配置、健康检查
# 作者: Quant Trader Team
# 版本: 1.0.0
# ============================================================

set -e  # 遇到错误立即退出

# ==================== 配置区域 ====================

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_NAME="quant-trader"
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="${PROJECT_DIR}/backend"
FRONTEND_DIR="${PROJECT_DIR}/frontend"
DEPLOY_DIR="${PROJECT_DIR}/deploy"
LOG_DIR="${PROJECT_DIR}/logs"
BACKUP_DIR="${PROJECT_DIR}/backups"

# 版本控制
VERSION_FILE="${PROJECT_DIR}/VERSION"
CURRENT_VERSION=$(cat "${VERSION_FILE}" 2>/dev/null || echo "0.0.0")
DEPLOY_TIME=$(date +"%Y%m%d_%H%M%S")
DEPLOY_ID="${DEPLOY_TIME}_$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

# 日志文件
LOG_FILE="${LOG_DIR}/deploy_${DEPLOY_ID}.log"

# 环境配置
ENVIRONMENT="${1:-development}"  # 默认开发环境
DOCKER_COMPOSE_FILE="docker-compose.${ENVIRONMENT}.yml"

# 健康检查配置
HEALTH_CHECK_RETRIES=30
HEALTH_CHECK_INTERVAL=5
BACKEND_HEALTH_URL="http://localhost:8080/health"
FRONTEND_HEALTH_URL="http://localhost:5173"

# ==================== 工具函数 ====================

# 日志函数
log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${timestamp} [${level}] ${message}" | tee -a "${LOG_FILE}"
}

log_info() {
    log "INFO" "${GREEN}$@${NC}"
}

log_warn() {
    log "WARN" "${YELLOW}$@${NC}"
}

log_error() {
    log "ERROR" "${RED}$@${NC}"
}

log_step() {
    log "STEP" "${BLUE}========== $@ ==========${NC}"
}

# 错误处理
error_exit() {
    log_error "$1"
    # 发送通知（如果配置了）
    send_notification "部署失败" "$1"
    exit 1
}

# 发送通知
send_notification() {
    local title="$1"
    local message="$2"
    
    # 钉钉通知（如果配置了）
    if [ -n "${DINGTALK_WEBHOOK}" ]; then
        curl -s -X POST "${DINGTALK_WEBHOOK}" \
            -H 'Content-Type: application/json' \
            -d "{\"msgtype\":\"text\",\"text\":{\"content\":\"[${PROJECT_NAME}] ${title}: ${message}\"}}" > /dev/null 2>&1 || true
    fi
    
    # 企业微信通知（如果配置了）
    if [ -n "${WECHAT_WEBHOOK}" ]; then
        curl -s -X POST "${WECHAT_WEBHOOK}" \
            -H 'Content-Type: application/json' \
            -d "{\"msgtype\":\"text\",\"text\":{\"content\":\"[${PROJECT_NAME}] ${title}: ${message}\"}}" > /dev/null 2>&1 || true
    fi
}

# 检查命令是否存在
check_command() {
    if ! command -v "$1" &> /dev/null; then
        error_exit "命令 $1 未安装，请先安装"
    fi
}

# 创建目录
ensure_dirs() {
    mkdir -p "${LOG_DIR}"
    mkdir -p "${BACKUP_DIR}"
    mkdir -p "${DEPLOY_DIR}"
}

# ==================== 部署前检查 ====================

pre_deploy_check() {
    log_step "部署前检查"
    
    # 检查必要命令
    log_info "检查必要命令..."
    check_command git
    check_command docker
    check_command docker-compose || check_command docker
    
    # 检查 Docker 服务
    log_info "检查 Docker 服务..."
    if ! docker info &> /dev/null; then
        error_exit "Docker 服务未运行，请启动 Docker"
    fi
    
    # 检查磁盘空间
    log_info "检查磁盘空间..."
    local available_space=$(df -BG "${PROJECT_DIR}" | awk 'NR==2 {print $4}' | sed 's/G//')
    if [ "${available_space}" -lt 5 ]; then
        error_exit "磁盘空间不足 (当前: ${available_space}GB, 需要: 5GB)"
    fi
    
    # 检查端口占用
    log_info "检查端口占用..."
    local ports=("8080" "5173" "5432" "6379" "4222" "9090")
    for port in "${ports[@]}"; do
        if lsof -i:"${port}" &> /dev/null; then
            log_warn "端口 ${port} 已被占用"
        fi
    done
    
    # 检查环境变量文件
    log_info "检查环境变量..."
    if [ ! -f "${BACKEND_DIR}/.env" ]; then
        log_warn "未找到 .env 文件，从模板创建..."
        cp "${BACKEND_DIR}/.env.example" "${BACKEND_DIR}/.env"
    fi
    
    log_info "部署前检查完成 ✓"
}

# ==================== 版本控制 ====================

backup_current_version() {
    log_step "备份当前版本"
    
    if [ -d "${DEPLOY_DIR}/current" ]; then
        local backup_path="${BACKUP_DIR}/backup_${DEPLOY_ID}"
        log_info "备份到 ${backup_path}..."
        cp -r "${DEPLOY_DIR}/current" "${backup_path}"
        log_info "备份完成 ✓"
    else
        log_info "无当前版本，跳过备份"
    fi
}

create_version_tag() {
    log_step "创建版本标签"
    
    # 更新版本文件
    echo "${DEPLOY_ID}" > "${VERSION_FILE}"
    
    # 创建 Git 标签（如果是 Git 仓库）
    if git rev-parse --git-dir > /dev/null 2>&1; then
        git tag -a "deploy_${DEPLOY_ID}" -m "Deploy at ${DEPLOY_TIME}" 2>/dev/null || true
    fi
    
    log_info "版本标签: ${DEPLOY_ID} ✓"
}

rollback() {
    log_step "执行回滚"
    
    local last_backup=$(ls -t "${BACKUP_DIR}" | head -1)
    if [ -z "${last_backup}" ]; then
        error_exit "没有可用的备份版本"
    fi
    
    log_info "回滚到版本: ${last_backup}"
    
    # 停止当前服务
    stop_services
    
    # 恢复备份
    rm -rf "${DEPLOY_DIR}/current"
    cp -r "${BACKUP_DIR}/${last_backup}" "${DEPLOY_DIR}/current"
    
    # 启动服务
    start_services
    
    log_info "回滚完成 ✓"
}

# ==================== 代码拉取 ====================

pull_code() {
    log_step "拉取最新代码"
    
    cd "${PROJECT_DIR}"
    
    # 检查是否有未提交的更改
    if git diff-index --quiet HEAD -- 2>/dev/null; then
        log_info "拉取代码..."
        git pull origin main || git pull origin master
        log_info "代码拉取完成 ✓"
    else
        log_warn "有未提交的更改，跳过拉取"
        git status --short
    fi
}

# ==================== 依赖安装 ====================

install_dependencies() {
    log_step "安装依赖"
    
    # 后端依赖
    log_info "安装后端依赖..."
    cd "${BACKEND_DIR}"
    go mod download
    go mod tidy
    log_info "后端依赖安装完成 ✓"
    
    # 前端依赖
    log_info "安装前端依赖..."
    cd "${FRONTEND_DIR}"
    if [ -f "package-lock.json" ]; then
        npm ci
    else
        npm install
    fi
    log_info "前端依赖安装完成 ✓"
}

# ==================== 构建打包 ====================

build_backend() {
    log_step "构建后端"
    
    cd "${BACKEND_DIR}"
    
    # 设置 Go 构建参数
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64
    
    # 构建
    log_info "编译后端服务..."
    go build -ldflags="-s -w -X main.Version=${DEPLOY_ID}" -o "${DEPLOY_DIR}/backend/quant-trader" ./cmd/main.go
    
    log_info "后端构建完成 ✓"
}

build_frontend() {
    log_step "构建前端"
    
    cd "${FRONTEND_DIR}"
    
    # 构建生产版本
    log_info "编译前端应用..."
    npm run build
    
    # 复制构建产物
    mkdir -p "${DEPLOY_DIR}/frontend"
    cp -r dist/* "${DEPLOY_DIR}/frontend/"
    
    log_info "前端构建完成 ✓"
}

build_docker_images() {
    log_step "构建 Docker 镜像"
    
    cd "${PROJECT_DIR}"
    
    # 构建后端镜像
    log_info "构建后端镜像..."
    docker build -t "${PROJECT_NAME}-backend:${DEPLOY_ID}" -t "${PROJECT_NAME}-backend:latest" \
        -f "${DEPLOY_DIR}/docker/Dockerfile.backend" .
    
    # 构建前端镜像
    log_info "构建前端镜像..."
    docker build -t "${PROJECT_NAME}-frontend:${DEPLOY_ID}" -t "${PROJECT_NAME}-frontend:latest" \
        -f "${DEPLOY_DIR}/docker/Dockerfile.frontend" .
    
    log_info "Docker 镜像构建完成 ✓"
}

# ==================== 环境配置 ====================

setup_environment() {
    log_step "配置环境"
    
    # 复制配置文件
    log_info "复制配置文件..."
    mkdir -p "${DEPLOY_DIR}/config"
    
    if [ "${ENVIRONMENT}" = "production" ]; then
        cp "${DEPLOY_DIR}/docker/.env.production" "${DEPLOY_DIR}/.env"
    else
        cp "${DEPLOY_DIR}/docker/.env.development" "${DEPLOY_DIR}/.env"
    fi
    
    # 设置文件权限
    chmod 600 "${DEPLOY_DIR}/.env"
    
    log_info "环境配置完成 ✓"
}

# ==================== 服务管理 ====================

stop_services() {
    log_step "停止服务"
    
    cd "${PROJECT_DIR}"
    
    if [ -f "${DEPLOY_DIR}/docker/${DOCKER_COMPOSE_FILE}" ]; then
        docker-compose -f "${DEPLOY_DIR}/docker/${DOCKER_COMPOSE_FILE}" down
    fi
    
    log_info "服务已停止 ✓"
}

start_services() {
    log_step "启动服务"
    
    cd "${PROJECT_DIR}"
    
    docker-compose -f "${DEPLOY_DIR}/docker/${DOCKER_COMPOSE_FILE}" up -d
    
    log_info "服务启动中..."
}

restart_services() {
    log_step "重启服务"
    
    stop_services
    start_services
}

# ==================== 健康检查 ====================

health_check() {
    log_step "健康检查"
    
    local retries=0
    local backend_healthy=false
    local frontend_healthy=false
    
    log_info "等待服务启动..."
    
    while [ $retries -lt $HEALTH_CHECK_RETRIES ]; do
        # 检查后端
        if ! $backend_healthy; then
            if curl -sf "${BACKEND_HEALTH_URL}" > /dev/null 2>&1; then
                backend_healthy=true
                log_info "后端服务健康 ✓"
            fi
        fi
        
        # 检查前端（仅开发环境）
        if [ "${ENVIRONMENT}" = "development" ] && ! $frontend_healthy; then
            if curl -sf "${FRONTEND_HEALTH_URL}" > /dev/null 2>&1; then
                frontend_healthy=true
                log_info "前端服务健康 ✓"
            fi
        fi
        
        # 检查是否全部健康
        if $backend_healthy; then
            if [ "${ENVIRONMENT}" = "production" ] || $frontend_healthy; then
                log_info "所有服务健康检查通过 ✓"
                return 0
            fi
        fi
        
        retries=$((retries + 1))
        sleep $HEALTH_CHECK_INTERVAL
    done
    
    error_exit "健康检查失败，请检查服务日志"
}

# ==================== 数据库迁移 ====================

run_migrations() {
    log_step "执行数据库迁移"
    
    cd "${BACKEND_DIR}"
    
    # 等待数据库启动
    sleep 5
    
    # 执行迁移脚本
    for migration in scripts/migrations/*.sql; do
        if [ -f "$migration" ]; then
            log_info "执行迁移: $(basename $migration)"
            # 使用 Docker 执行迁移
            docker-compose -f "${DEPLOY_DIR}/docker/${DOCKER_COMPOSE_FILE}" exec -T timescaledb \
                psql -U postgres -d quant_trader -f "/docker-entrypoint-initdb.d/$(basename $migration)" || true
        fi
    done
    
    log_info "数据库迁移完成 ✓"
}

# ==================== 主流程 ====================

main() {
    # 初始化
    ensure_dirs
    
    log_info "============================================"
    log_info "开始部署: ${PROJECT_NAME}"
    log_info "环境: ${ENVIRONMENT}"
    log_info "版本: ${DEPLOY_ID}"
    log_info "============================================"
    
    # 记录开始时间
    local start_time=$(date +%s)
    
    # 执行部署流程
    pre_deploy_check
    backup_current_version
    pull_code
    install_dependencies
    build_backend
    build_frontend
    build_docker_images
    setup_environment
    stop_services
    start_services
    run_migrations
    health_check
    create_version_tag
    
    # 计算耗时
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    log_info "============================================"
    log_info "部署成功! 耗时: ${duration}秒"
    log_info "============================================"
    
    # 发送成功通知
    send_notification "部署成功" "版本: ${DEPLOY_ID}, 耗时: ${duration}秒"
}

# ==================== 命令行参数处理 ====================

case "${1}" in
    deploy)
        main
        ;;
    rollback)
        rollback
        ;;
    stop)
        stop_services
        ;;
    start)
        start_services
        ;;
    restart)
        restart_services
        ;;
    health)
        health_check
        ;;
    build)
        build_backend
        build_frontend
        ;;
    *)
        echo "用法: $0 {deploy|rollback|stop|start|restart|health|build} [environment]"
        echo ""
        echo "命令:"
        echo "  deploy    - 完整部署流程"
        echo "  rollback  - 回滚到上一版本"
        echo "  stop      - 停止服务"
        echo "  start     - 启动服务"
        echo "  restart   - 重启服务"
        echo "  health    - 健康检查"
        echo "  build     - 构建应用"
        echo ""
        echo "环境:"
        echo "  development  - 开发环境 (默认)"
        echo "  production   - 生产环境"
        exit 1
        ;;
esac
