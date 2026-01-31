#!/bin/bash

# ============================================
# Docker Compose 部署脚本
# 生产环境 - 抽奖系统
# ============================================

set -e  # 遇到错误立即退出
set -u  # 使用未定义变量时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 未安装，请先安装"
        exit 1
    fi
}

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        log_warning "端口 $port 已被占用"
        read -p "是否继续？(y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."

    mkdir -p docker/mysql/data
    mkdir -p docker/mysql/conf.d
    mkdir -p docker/mysql/init
    mkdir -p docker/redis/data
    mkdir -p docker/caddy/data
    mkdir -p docker/caddy/logs
    mkdir -p docker/backend/logs

    # 设置权限
    chmod -R 755 docker

    log_success "目录创建完成"
}

# 检查环境配置文件
check_env_file() {
    log_info "检查环境配置文件..."

    if [ ! -f "docker-compose-production.env" ]; then
        log_error "docker-compose-production.env 文件不存在"
        exit 1
    fi

    # 检查关键配置
    source docker-compose-production.env

    if [ -z "${MYSQL_ROOT_PASSWORD:-}" ] || [ "$MYSQL_ROOT_PASSWORD" == "change_in_production" ]; then
        log_error "MYSQL_ROOT_PASSWORD 未设置或使用默认值，请在 docker-compose-production.env 中修改"
        exit 1
    fi

    if [ -z "${JWT_SECRET:-}" ] || [ "$JWT_SECRET" == "change-in-production-min-32-chars" ]; then
        log_error "JWT_SECRET 未设置或使用默认值，请在 docker-compose-production.env 中修改"
        exit 1
    fi

    log_success "环境配置检查通过"
}

# 检查Docker和Docker Compose
check_docker() {
    log_info "检查Docker环境..."

    check_command docker
    check_command docker

    # 检查Docker服务是否运行
    if ! docker info &> /dev/null; then
        log_error "Docker服务未运行，请先启动Docker"
        exit 1
    fi

    log_success "Docker环境检查通过"
}

# 构建镜像
build_images() {
    log_info "构建Docker镜像..."

    docker compose build --no-cache

    log_success "镜像构建完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."

    docker compose up -d

    log_success "服务启动完成"
}

# 查看服务状态
check_services() {
    log_info "查看服务状态..."

    docker compose ps

    echo ""
    log_info "服务健康状态："
    docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
}

# 查看日志
view_logs() {
    local service=${1:-}
    if [ -z "$service" ]; then
        log_info "查看所有服务日志（按Ctrl+C退出）..."
        docker compose logs -f
    else
        log_info "查看 $service 服务日志（按Ctrl+C退出）..."
        docker compose logs -f $service
    fi
}

# 停止服务
stop_services() {
    log_info "停止服务..."

    docker compose down

    log_success "服务已停止"
}

# 重启服务
restart_services() {
    log_info "重启服务..."

    docker compose restart

    log_success "服务已重启"
}

# 清理数据（危险操作）
clean_data() {
    log_warning "警告：此操作将删除所有数据！"
    read -p "确定要清理所有数据吗？(yes/no) " -r
    echo
    if [[ $REPLY == "yes" ]]; then
        log_info "停止服务..."
        docker compose down -v

        log_info "删除数据目录..."
        rm -rf docker/mysql/data/*
        rm -rf docker/redis/data/*
        rm -rf docker/caddy/data/*
        rm -rf docker/backend/logs/*

        log_success "数据清理完成"
    else
        log_info "取消数据清理"
    fi
}

# 备份数据
backup_data() {
    local backup_dir="backups/$(date +%Y%m%d_%H%M%S)"
    log_info "备份数据到 $backup_dir ..."

    mkdir -p $backup_dir

    # 备份MySQL数据
    if [ -d "docker/mysql/data" ]; then
        log_info "备份MySQL数据..."
        cp -r docker/mysql/data $backup_dir/
    fi

    # 备份Redis数据
    if [ -d "docker/redis/data" ]; then
        log_info "备份Redis数据..."
        cp -r docker/redis/data $backup_dir/
    fi

    # 备份配置文件
    log_info "备份配置文件..."
    cp docker-compose-production.env $backup_dir/
    cp -r docker/mysql/conf.d $backup_dir/
    cp -r docker/caddy $backup_dir/

    log_success "数据备份完成：$backup_dir"
}

# 健康检查
health_check() {
    log_info "执行健康检查..."

    echo ""
    echo "=== 服务状态 ==="
    docker compose ps

    echo ""
    echo "=== MySQL健康检查 ==="
    docker compose exec -T mysql mysqladmin ping -h localhost -u root -p${MYSQL_ROOT_PASSWORD} 2>/dev/null && \
        log_success "MySQL: 健康" || log_error "MySQL: 不健康"

    echo ""
    echo "=== Redis健康检查 ==="
    docker compose exec -T redis redis-cli -a ${REDIS_PASSWORD} ping 2>/dev/null && \
        log_success "Redis: 健康" || log_error "Redis: 不健康"

    echo ""
    echo "=== 后端健康检查 ==="
    curl -s http://localhost:8080/api/health > /dev/null 2>&1 && \
        log_success "后端: 健康" || log_error "后端: 不健康"

    echo ""
    echo "=== 前端健康检查 ==="
    curl -s http://localhost/ > /dev/null 2>&1 && \
        log_success "前端: 健康" || log_error "前端: 不健康"
}

# 显示帮助信息
show_help() {
    echo "抽奖系统 - Docker Compose 部署脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  deploy      部署（检查环境、构建镜像、启动服务）"
    echo "  start       启动服务"
    echo "  stop        停止服务"
    echo "  restart     重启服务"
    echo "  status      查看服务状态"
    echo "  logs        查看所有服务日志"
    echo "  logs [srv]  查看指定服务日志（mysql|redis|backend|frontend|caddy）"
    echo "  health      健康检查"
    echo "  backup      备份数据"
    echo "  clean       清理数据（危险操作）"
    echo "  help        显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 deploy           # 完整部署"
    echo "  $0 logs backend     # 查看后端日志"
    echo "  $0 health           # 健康检查"
}

# 主函数
main() {
    local command=${1:-help}

    case $command in
        deploy)
            check_docker
            check_env_file
            create_directories
            check_port 80
            check_port 443
            build_images
            start_services
            sleep 5
            check_services
            health_check
            log_success "部署完成！"
            log_info "前端访问地址: http://localhost"
            log_info "后端API地址: http://localhost:8080"
            ;;
        start)
            check_docker
            create_directories
            start_services
            check_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        status)
            check_services
            ;;
        logs)
            view_logs $2
            ;;
        health)
            health_check
            ;;
        backup)
            backup_data
            ;;
        clean)
            clean_data
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
