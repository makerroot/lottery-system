#!/bin/bash

# ============================================
# 抽奖系统 - 一键启动部署脚本
# ============================================

set -e

# 确保常见路径在PATH中（解决某些环境PATH不完整的问题）
export PATH="/usr/bin:/usr/local/bin:/bin:/usr/sbin:/usr/local/sbin:$PATH"

# Docker Compose命令（全局变量，将在check_docker中设置）
DOCKER_COMPOSE_CMD=""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
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

# 显示Logo
show_logo() {
    echo -e "${CYAN}"
    cat << "EOF"
╔════════════════════════════════════════╗
║     🎉 抽奖系统 - 一键启动部署 🎉      ║
║    Lottery System - Quick Start        ║
╚════════════════════════════════════════╝
EOF
    echo -e "${NC}"
}

# 显示菜单
show_menu() {
    echo ""
    echo -e "${CYAN}请选择启动模式:${NC}"
    echo ""
    echo "  1) 🐳 Docker Compose模式（推荐）"
    echo "     - 自动启动所有服务（MySQL、Redis、后端、Caddy）"
    echo "     - Caddy提供前端静态文件服务和反向代理"
    echo "     - 使用HTTPS（Let's Encrypt证书）"
    echo "     - 适合生产环境和完整测试"
    echo ""
    echo "  2) 💻 本地开发模式"
    echo "     - 分别启动后端和前端"
    echo "     - 需要本地安装Go和Node.js"
    echo "     - 适合开发调试"
    echo ""
    echo "  3) 🔄 仅启动后端（本地）"
    echo "     - 启动Go后端服务"
    echo ""
    echo "  4) 🎨 仅启动前端（本地）"
    echo "     - 启动Vue前端开发服务器"
    echo ""
    echo "  5) 🛑 停止所有服务"
    echo ""
    echo "  6) 📊 查看服务状态"
    echo ""
    echo "  0) 🚪 退出"
    echo ""
    echo -n "请输入选项 [0-6]: "
}

# 检查Docker
check_docker() {
    log_info "检查Docker环境..."

    # 方法1：使用command -v检查
    if ! command -v docker &> /dev/null; then
        # 方法2：直接检查常见安装路径
        if [ ! -f "/usr/bin/docker" ] && [ ! -f "/usr/local/bin/docker" ]; then
            log_error "Docker未安装，请先安装Docker"
            echo "安装指南: https://docs.docker.com/get-docker/"
            echo ""
            echo "提示：如果Docker已安装，请尝试："
            echo "  1. 添加到PATH: export PATH=\$PATH:/usr/bin:/usr/local/bin"
            echo "  2. 或创建符号链接: sudo ln -sf /usr/bin/docker /usr/local/bin/docker"
            exit 1
        fi

        # 找到了docker，添加到PATH
        if [ -f "/usr/bin/docker" ]; then
            export PATH="/usr/bin:$PATH"
        elif [ -f "/usr/local/bin/docker" ]; then
            export PATH="/usr/local/bin:$PATH"
        fi
        log_info "已将Docker添加到PATH"
    fi

    # 检查Docker Compose（优先检测插件版本）
    DOCKER_COMPOSE_CMD=""

    # 优先检测 docker compose（插件版本，V2）
    if docker compose version &> /dev/null 2>&1; then
        DOCKER_COMPOSE_CMD="docker compose"
        log_info "使用 docker compose 命令（插件版本 V2）"
    # 然后检测 docker-compose（独立版本，V1）
    elif command -v docker-compose &> /dev/null; then
        DOCKER_COMPOSE_CMD="docker-compose"
        log_info "使用 docker-compose 命令（独立版本 V1）"
    # 最后检查常见路径
    else
        if [ -f "/usr/bin/docker-compose" ] || [ -f "/usr/local/bin/docker-compose" ]; then
            DOCKER_COMPOSE_CMD="docker-compose"
            log_info "使用 docker-compose 命令（直接路径）"
        # 尝试直接调用docker compose（即使command -v找不到）
        elif [ -f "/usr/bin/docker" ]; then
            if /usr/bin/docker compose version &> /dev/null 2>&1; then
                DOCKER_COMPOSE_CMD="docker compose"
                log_info "使用 docker compose 命令（完整路径）"
            else
                log_error "Docker Compose未安装，请先安装Docker Compose"
                echo "安装指南: https://docs.docker.com/compose/install/"
                echo ""
                echo "您的Docker版本似乎包含compose插件，请尝试："
                echo "  /usr/bin/docker compose version"
                exit 1
            fi
        else
            log_error "Docker Compose未安装，请先安装Docker Compose"
            echo "安装指南: https://docs.docker.com/compose/install/"
            exit 1
        fi
    fi

    # 导出到全局变量，确保其他函数可以使用
    export DOCKER_COMPOSE_CMD

    # 检查Docker服务是否运行
    # 使用更可靠的方法：检查dockerd进程或尝试docker ps
    DOCKER_RUNNING=false

    # 方法1：检查docker进程
    if pgrep -x dockerd > /dev/null 2>&1; then
        DOCKER_RUNNING=true
    fi

    # 方法2：尝试docker命令（使用完整路径避免PATH问题）
    if [ -f "/usr/bin/docker" ]; then
        if /usr/bin/docker ps > /dev/null 2>&1; then
            DOCKER_RUNNING=true
        fi
    elif [ -f "/usr/local/bin/docker" ]; then
        if /usr/local/bin/docker ps > /dev/null 2>&1; then
            DOCKER_RUNNING=true
        fi
    fi

    # 方法3：使用docker info（作为最后手段）
    if [ "$DOCKER_RUNNING" = false ]; then
        if docker info > /dev/null 2>&1; then
            DOCKER_RUNNING=true
        fi
    fi

    if [ "$DOCKER_RUNNING" = false ]; then
        log_error "Docker服务未运行，请先启动Docker"
        echo ""
        echo "启动命令："
        echo "  sudo systemctl start docker    # Systemd系统"
        echo "  sudo service docker start      # SysV系统"
        echo "  dockerd &                      # 直接启动"
        echo ""
        echo "调试信息："
        echo "  检查进程: ps aux | grep dockerd"
        echo "  检查服务: systemctl status docker"
        exit 1
    fi

    log_success "Docker环境检查通过"
}

# 检查环境配置
check_env() {
    log_info "检查环境配置..."

    if [ ! -f "docker-compose-production.env" ]; then
        log_error "docker-compose-production.env 文件不存在"
        exit 1
    fi

    log_info "验证环境配置文件..."

    # 检查文件格式（基本检查）
    if grep -E '^[A-Za-z_][A-Za-z0-9_]*=' docker-compose-production.env | head -1 >/dev/null 2>&1; then
        log_info "配置文件格式正确"
    else
        log_warning "配置文件可能为空或格式异常"
    fi

    # 检查关键配置项
    if ! grep -q "MYSQL_ROOT_PASSWORD=" docker-compose-production.env; then
        log_warning "⚠️  MYSQL_ROOT_PASSWORD 未配置"
    fi

    if ! grep -q "JWT_SECRET=" docker-compose-production.env; then
        log_warning "⚠️  JWT_SECRET 未配置"
    fi

    log_success "环境配置检查完成"
}

# 检查HTTPS证书
check_https_cert() {
    log_info "检查HTTPS证书..."

    CERT_DIR="/etc/letsencrypt/live/makerroot.com"
    FULLCHAIN="$CERT_DIR/fullchain.pem"
    PRIVKEY="$CERT_DIR/privkey.pem"

    if [ ! -f "$FULLCHAIN" ]; then
        log_warning "⚠️  SSL证书不存在: $FULLCHAIN"
        echo ""
        echo "HTTP模式将启动，但HTTPS不可用"
        echo ""
        echo "如需启用HTTPS，请先获取证书："
        echo "  sudo certbot certonly --standalone -d makerroot.com"
        echo ""
        read -p "是否继续启动HTTP模式？(y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "取消启动"
            exit 1
        fi
        return 1
    fi

    if [ ! -f "$PRIVKEY" ]; then
        log_warning "⚠️  SSL私钥不存在: $PRIVKEY"
        return 1
    fi

    # 检查并设置证书权限
    log_info "设置证书权限..."
    sudo chmod 644 "$FULLCHAIN" 2>/dev/null || log_warning "无法设置证书权限（需要sudo）"
    sudo chmod 600 "$PRIVKEY" 2>/dev/null || log_warning "无法设置私钥权限（需要sudo）"
    sudo chmod 755 "$CERT_DIR" 2>/dev/null || true

    log_success "SSL证书检查通过"
    return 0
}

# Docker Compose模式
start_docker_compose() {
    echo ""
    log_info "🐳 启动Docker Compose模式..."

    check_docker
    check_env

    # 检查HTTPS证书
    HTTPS_ENABLED=false
    if check_https_cert; then
        HTTPS_ENABLED=true
    fi

    # 创建必要目录
    log_info "创建数据目录..."
    mkdir -p docker/mysql/data
    mkdir -p docker/mysql/conf.d
    mkdir -p docker/mysql/init
    mkdir -p docker/redis/data
    mkdir -p docker/caddy/data
    mkdir -p docker/caddy/logs
    mkdir -p docker/backend/logs

    # 检查是否已有运行的服务
    if $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env ps 2>/dev/null | grep -q "Up"; then
        log_warning "检测到已有运行的服务"
        read -p "是否重启服务？(y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            log_info "重启服务..."
            $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env restart
        else
            log_info "保持现有服务运行"
        fi
    else
        log_info "构建并启动服务..."
        $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env up -d --build

        log_info "等待服务启动（约30秒）..."
        sleep 5
    fi

    echo ""
    log_success "🎉 Docker Compose模式启动成功！"
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    # 根据HTTPS状态显示不同的访问地址
    if [ "$HTTPS_ENABLED" = true ]; then
        echo -e "${CYAN}访问地址 (HTTPS模式):${NC}"
        echo -e "  🌐 前端: ${GREEN}https://makerroot.com${NC}"
        echo -e "  🔌 后端API: ${GREEN}https://makerroot.com/api/*${NC}"
        echo -e "  👨‍💼 管理后台: ${GREEN}https://makerroot.com/admin${NC}"
        echo ""
        echo -e "${CYAN}本地测试 (如需):${NC}"
        echo -e "  HTTP: ${YELLOW}http://localhost${NC} (会重定向到HTTPS)"
        echo -e "  HTTPS: ${GREEN}https://localhost${NC}"
    else
        echo -e "${YELLOW}访问地址 (HTTP模式 - 证书未配置):${NC}"
        echo -e "  🌐 前端: ${YELLOW}http://localhost${NC}"
        echo -e "  🔌 后端API: ${YELLOW}http://localhost/api/*${NC}"
        echo -e "  👨‍💼 管理后台: ${YELLOW}http://localhost/admin${NC}"
        echo ""
        echo -e "${YELLOW}⚠️  HTTPS未启用，配置证书后重启即可启用HTTPS${NC}"
        echo ""
        echo "获取证书命令:"
        echo "  sudo certbot certonly --standalone -d makerroot.com"
    fi

    echo ""
    echo -e "${CYAN}默认账号:${NC}"
    echo -e "  管理员: ${GREEN}makerroot / 123456${NC}"
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${CYAN}常用命令:${NC}"
    echo "  查看日志: $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env logs -f"
    echo "  查看Caddy日志: $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env logs -f caddy"
    echo "  停止服务: $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env down"
    echo "  服务状态: $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env ps"
    echo ""
    echo -e "${CYAN}证书管理:${NC}"
    echo "  检查证书: ./check-cert.sh"
    echo "  续期证书: certbot renew && $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env restart caddy"
}

# 本地开发模式
start_local_dev() {
    echo ""
    log_info "💻 启动本地开发模式..."

    # 检查Go
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go: https://golang.org/dl/"
        exit 1
    fi

    # 检查Node.js
    if ! command -v node &> /dev/null; then
        log_error "Node.js未安装，请先安装Node.js: https://nodejs.org/"
        exit 1
    fi

    log_success "依赖检查通过"

    # 启动后端
    echo ""
    log_info "启动后端服务..."
    cd backend

    # 检查.env
    if [ ! -f ".env" ]; then
        log_warning "backend/.env 不存在，使用默认配置"
    fi

    # 启动后端（后台运行）
    if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
        log_warning "后端服务已在运行 (端口8080)"
    else
        log_info "编译并启动后端..."
        nohup go run main.go > ../backend.log 2>&1 &
        BACKEND_PID=$!
        echo $BACKEND_PID > ../backend.pid
        log_success "后端服务已启动 (PID: $BACKEND_PID)"

        # 等待后端启动
        sleep 3
    fi

    cd ..

    # 启动前端
    echo ""
    log_info "启动前端服务..."
    cd frontend

    # 检查node_modules
    if [ ! -d "node_modules" ]; then
        log_info "安装前端依赖..."
        npm install
    fi

    # 启动前端（后台运行）
    if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null 2>&1; then
        log_warning "前端服务已在运行 (端口5173)"
    else
        log_info "启动前端开发服务器..."
        nohup npm run dev > ../frontend.log 2>&1 &
        FRONTEND_PID=$!
        echo $FRONTEND_PID > ../frontend.pid
        log_success "前端服务已启动 (PID: $FRONTEND_PID)"
    fi

    cd ..

    echo ""
    log_success "🎉 本地开发模式启动成功！"
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}访问地址:${NC}"
    echo -e "  前端: ${GREEN}http://localhost:5173${NC}"
    echo -e "  后端: ${GREEN}http://localhost:8080${NC}"
    echo -e "  管理后台: ${GREEN}http://localhost:5173/admin${NC}"
    echo ""
    echo -e "${CYAN}默认账号:${NC}"
    echo -e "  管理员: ${GREEN}makerroot / 123456${NC}"
    echo ""
    echo -e "${CYAN}日志文件:${NC}"
    echo "  后端日志: tail -f backend.log"
    echo "  前端日志: tail -f frontend.log"
    echo ""
    echo -e "${CYAN}停止服务:${NC}"
    echo "  ./start.sh (选择5)"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# 仅启动后端
start_backend_only() {
    echo ""
    log_info "🔄 启动后端服务..."

    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go: https://golang.org/dl/"
        exit 1
    fi

    cd backend

    if [ ! -f ".env" ]; then
        log_warning "backend/.env 不存在，使用默认配置"
    fi

    if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
        log_warning "⚠️  后端服务已在运行 (端口8080)"
        echo ""
        echo "后端API: http://localhost:8080"
    else
        log_info "编译并启动后端..."
        go run main.go
    fi
}

# 仅启动前端
start_frontend_only() {
    echo ""
    log_info "🎨 启动前端服务..."

    if ! command -v npm &> /dev/null; then
        log_error "npm未安装，请先安装Node.js: https://nodejs.org/"
        exit 1
    fi

    cd frontend

    if [ ! -d "node_modules" ]; then
        log_info "安装前端依赖..."
        npm install
    fi

    log_info "启动前端开发服务器..."
    npm run dev
}

# 停止所有服务
stop_all_services() {
    echo ""
    log_info "🛑 停止所有服务..."

    # 停止Docker Compose服务
    if $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env ps 2>/dev/null | grep -q "Up"; then
        log_info "停止Docker Compose服务..."
        $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env down
        log_success "Docker Compose服务已停止"
    fi

    # 停止本地后端
    if [ -f "backend.pid" ]; then
        BACKEND_PID=$(cat backend.pid)
        if ps -p $BACKEND_PID > /dev/null 2>&1; then
            log_info "停止后端服务 (PID: $BACKEND_PID)..."
            kill $BACKEND_PID 2>/dev/null || true
            rm backend.pid
            log_success "后端服务已停止"
        fi
    fi

    # 停止本地前端
    if [ -f "frontend.pid" ]; then
        FRONTEND_PID=$(cat frontend.pid)
        if ps -p $FRONTEND_PID > /dev/null 2>&1; then
            log_info "停止前端服务 (PID: $FRONTEND_PID)..."
            kill $FRONTEND_PID 2>/dev/null || true
            rm frontend.pid
            log_success "前端服务已停止"
        fi
    fi

    # 尝试停止端口占用
    if lsof -ti :8080 >/dev/null 2>&1; then
        log_info "停止端口8080的进程..."
        kill -9 $(lsof -ti :8080) 2>/dev/null || true
    fi

    if lsof -ti :5173 >/dev/null 2>&1; then
        log_info "停止端口5173的进程..."
        kill -9 $(lsof -ti :5173) 2>/dev/null || true
    fi

    echo ""
    log_success "✅ 所有服务已停止"
}

# 查看服务状态
show_status() {
    echo ""
    log_info "📊 服务状态..."
    echo ""

    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}Docker Compose 服务:${NC}"
    echo ""
    if $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env ps &>/dev/null; then
        $DOCKER_COMPOSE_CMD --env-file docker-compose-production.env ps
    else
        echo "  Docker Compose未运行"
    fi
    echo ""

    echo -e "${CYAN}本地服务端口:${NC}"
    echo ""

    # 检查端口占用
    if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "  ✅ 后端服务 (${GREEN}http://localhost:8080${NC}) - 运行中"
        lsof -i :8080 | grep LISTEN | awk '{printf "     PID: %s, 进程: %s\n", $2, $1}'
    else
        echo -e "  ❌ 后端服务 (端口8080) - 未运行"
    fi

    if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "  ✅ 前端服务 (${GREEN}http://localhost:5173${NC}) - 运行中"
        lsof -i :5173 | grep LISTEN | awk '{printf "     PID: %s, 进程: %s\n", $2, $1}'
    else
        echo -e "  ❌ 前端服务 (端口5173) - 未运行"
    fi

    echo ""

    # 检查HTTP/HTTPS端口
    HTTP_RUNNING=false
    HTTPS_RUNNING=false

    if lsof -Pi :80 -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "  ✅ HTTP服务 (${GREEN}http://localhost${NC}) - 运行中"
        HTTP_RUNNING=true
    fi

    if lsof -Pi :443 -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "  ✅ HTTPS服务 (${GREEN}https://makerroot.com${NC}) - 运行中"
        lsof -i :443 | grep LISTEN | awk '{printf "     PID: %s, 进程: %s\n", $2, $1}'
        HTTPS_RUNNING=true
    fi

    if [ "$HTTP_RUNNING" = false ] && [ "$HTTPS_RUNNING" = false ]; then
        echo -e "  ❌ HTTP/HTTPS服务 (端口80/443) - 未运行"
    fi

    echo ""

    # 检查证书状态
    if [ -f "/etc/letsencrypt/live/makerroot.com/fullchain.pem" ]; then
        echo -e "${CYAN}SSL证书状态:${NC}"
        EXPIRY_DATE=$(echo | openssl s_client -connect makerroot.com:443 2>/dev/null | openssl x509 -noout -enddate | cut -d= -f2 2>/dev/null)
        if [ -n "$EXPIRY_DATE" ]; then
            echo -e "  ✅ 证书到期时间: ${GREEN}$EXPIRY_DATE${NC}"
        else
            echo -e "  ⚠️  证书存在但无法读取到期时间"
        fi
        echo ""
    fi

    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

# 主函数
main() {
    show_logo

    # 如果有参数，直接执行
    if [ $# -gt 0 ]; then
        case $1 in
            docker|d)
                start_docker_compose
                ;;
            local|l|dev)
                start_local_dev
                ;;
            backend|b|go)
                start_backend_only
                ;;
            frontend|f|vue)
                start_frontend_only
                ;;
            stop|s)
                stop_all_services
                ;;
            status|st)
                show_status
                ;;
            *)
                echo "用法: $0 [docker|local|backend|frontend|stop|status]"
                echo ""
                echo "  docker   - Docker Compose模式"
                echo "  local    - 本地开发模式"
                echo "  backend  - 仅启动后端"
                echo "  frontend - 仅启动前端"
                echo "  stop     - 停止所有服务"
                echo "  status   - 查看服务状态"
                exit 1
                ;;
        esac
        exit 0
    fi

    # 交互式菜单
    while true; do
        show_menu
        read -r choice

        case $choice in
            1)
                start_docker_compose
                break
                ;;
            2)
                start_local_dev
                break
                ;;
            3)
                start_backend_only
                break
                ;;
            4)
                start_frontend_only
                break
                ;;
            5)
                stop_all_services
                ;;
            6)
                show_status
                ;;
            0)
                echo ""
                log_info "退出..."
                exit 0
                ;;
            *)
                echo ""
                log_error "无效选项，请重新选择"
                sleep 1
                ;;
        esac
    done
}

# 执行主函数
main "$@"
