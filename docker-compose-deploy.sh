#!/bin/bash

# ============================================
# Docker Compose 一键部署脚本
# ============================================

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🚀 Docker Compose 部署${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 1. 检查 Docker
echo -e "${YELLOW}📋 检查环境...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker 未安装${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose 未安装${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker 环境正常${NC}"
docker --version
docker-compose --version

# 2. 检查配置文件
echo ""
echo -e "${YELLOW}📋 检查配置文件...${NC}"
if [ ! -f "docker-compose-production.env" ]; then
    echo -e "${RED}❌ 配置文件不存在: docker-compose-production.env${NC}"
    echo -e "${YELLOW}从模板创建:${NC}"
    cp .env.production.template docker-compose-production.env
    echo -e "${RED}⚠️  请编辑 docker-compose-production.env 并修改密码！${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 配置文件存在${NC}"

# 3. 检查证书
echo ""
echo -e "${YELLOW}📋 检查证书...${NC}"
if [ ! -f "/etc/letsencrypt/live/makerroot.com/fullchain.pem" ]; then
    echo -e "${YELLOW}⚠️  Let's Encrypt 证书不存在${NC}"
    echo -e "${YELLOW}   将使用 HTTP 模式${NC}"
    HAS_CERT="false"
else
    echo -e "${GREEN}✅ Let's Encrypt 证书存在${NC}"
    HAS_CERT="true"
fi

# 4. 检查前端构建产物
echo ""
echo -e "${YELLOW}📋 检查前端构建产物...${NC}"
if [ ! -d "frontend/dist" ]; then
    echo -e "${RED}❌ frontend/dist 不存在${NC}"
    echo -e "${YELLOW}   正在从 Git 拉取...${NC}"
    git pull origin main
    
    if [ ! -d "frontend/dist" ]; then
        echo -e "${RED}❌ frontend/dist 仍然不存在${NC}"
        echo -e "${YELLOW}   请先在本地构建: npm run build${NC}"
        exit 1
    fi
fi
echo -e "${GREEN}✅ frontend/dist 存在${NC}"
du -sh frontend/dist/

# 5. 停止旧服务
echo ""
echo -e "${YELLOW}🛑 停止旧服务...${NC}"
docker-compose --env-file docker-compose-production.env down 2>/dev/null || true

# 6. 清理旧 volumes（可选）
echo ""
read -p "是否清理旧数据? (仅首次部署或需要重置数据时选择 y) [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🧹 清理旧数据...${NC}"
    docker volume rm lottery-frontend-dist 2>/dev/null || true
    docker volume rm lottery-mysql-data 2>/dev/null || true
    docker volume rm lottery-redis-data 2>/dev/null || true
    echo -e "${GREEN}✅ 清理完成${NC}"
fi

# 7. 拉取最新代码
echo ""
echo -e "${YELLOW}📦 拉取最新代码...${NC}"
git pull origin main

# 8. 启动服务
echo ""
echo -e "${YELLOW}🚀 启动服务...${NC}"
docker-compose --env-file docker-compose-production.env up -d

# 9. 等待服务启动
echo ""
echo -e "${YELLOW}⏳ 等待服务启动...${NC}"
sleep 15

# 10. 显示服务状态
echo ""
echo -e "${YELLOW}📊 服务状态:${NC}"
docker-compose --env-file docker-compose-production.env ps

# 11. 检查服务健康
echo ""
echo -e "${YELLOW}🔍 检查服务健康...${NC}"

# 检查 MySQL
if docker-compose --env-file docker-compose-production.env ps | grep -q "lottery-mysql.*Up (healthy)"; then
    echo -e "${GREEN}✅ MySQL - 健康${NC}"
else
    echo -e "${RED}❌ MySQL - 异常${NC}"
fi

# 检查 Redis
if docker-compose --env-file docker-compose-production.env ps | grep -q "lottery-redis.*Up (healthy)"; then
    echo -e "${GREEN}✅ Redis - 健康${NC}"
else
    echo -e "${RED}❌ Redis - 异常${NC}"
fi

# 检查 Backend
if docker-compose --env-file docker-compose-production.env ps | grep -q "lottery-backend.*Up (healthy)"; then
    echo -e "${GREEN}✅ Backend - 健康${NC}"
else
    echo -e "${RED}❌ Backend - 异常${NC}"
fi

# 检查 Caddy
if docker-compose --env-file docker-compose-production.env ps | grep -q "lottery-caddy.*Up"; then
    echo -e "${GREEN}✅ Caddy - 运行中${NC}"
else
    echo -e "${RED}❌ Caddy - 异常${NC}"
fi

# 12. 测试访问
echo ""
echo -e "${YELLOW}🧪 测试访问...${NC}"
sleep 3

if [ "$HAS_CERT" = "true" ]; then
    # 测试 HTTPS
    if curl -skI https://localhost/ 2>&1 | grep -q "HTTP"; then
        echo -e "${GREEN}✅ HTTPS 访问正常${NC}"
    else
        echo -e "${YELLOW}⚠️  HTTPS 访问测试失败（可能需要等待证书生效）${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  未配置证书，使用 HTTP${NC}"
fi

# 测试健康检查
if curl -s http://localhost/api/health | grep -q "正常"; then
    echo -e "${GREEN}✅ 后端 API 正常${NC}"
else
    echo -e "${RED}❌ 后端 API 异常${NC}"
fi

# 13. 显示访问信息
echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ 部署完成！${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if [ "$HAS_CERT" = "true" ]; then
    echo "📍 访问地址:"
    echo "   HTTP:  http://makerroot.com"
    echo "   HTTPS: https://makerroot.com"
    echo "   管理后台: https://makerroot.com/#/admin/"
else
    echo "📍 访问地址:"
    echo "   HTTP:  http://makerroot.com"
    echo "   管理后台: http://makerroot.com/#/admin/"
    echo ""
    echo -e "${YELLOW}⚠️  HTTPS 未配置${NC}"
    echo "   要启用 HTTPS，请确保证书文件存在"
fi

echo ""
echo "🔑 默认账号:"
echo "   用户名: makerroot"
echo "   密码: 123456"
echo ""

echo "💡 常用命令:"
echo "   查看日志: docker-compose --env-file docker-compose-production.env logs -f"
echo "   查看状态: docker-compose --env-file docker-compose-production.env ps"
echo "   重启服务: docker-compose --env-file docker-compose-production.env restart"
echo "   停止服务: docker-compose --env-file docker-compose-production.env down"
echo ""

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
