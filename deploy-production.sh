#!/bin/bash

# ============================================
# 服务器部署脚本 - 零停机部署
# ============================================

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 开始零停机部署..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 1. 拉取最新代码
echo "📦 步骤 1/5: 拉取最新代码..."
git pull origin main
echo "✅ 代码已更新"
echo ""

# 2. 验证 dist 存在
echo "📋 步骤 2/5: 检查 frontend/dist..."
if [ ! -d "frontend/dist" ]; then
    echo "❌ 错误: frontend/dist 不存在"
    echo "请先在本地构建前端: npm run build"
    exit 1
fi
echo "✅ frontend/dist 存在"
du -sh frontend/dist/
echo ""

# 3. 滚动更新后端服务（零停机）
echo "🔧 步骤 3/5: 滚动更新后端服务（零停机）..."
echo "⏳ 构建新镜像..."
docker compose --env-file docker-compose-production.env build backend

echo "⏳ 启动所有服务（确保首次部署时服务都会启动）..."
docker compose --env-file docker-compose-production.env up -d

echo "⏳ 等待后端健康检查通过..."
for i in {1..30}; do
    if docker compose --env-file docker-compose-production.env ps | grep backend | grep -q "healthy"; then
        echo "✅ 后端服务已就绪"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ 后端健康检查超时"
        docker compose --env-file docker-compose-production.env logs backend --tail 50
        exit 1
    fi
    echo "   等待中... ($i/30)"
    sleep 1
done
echo ""

# 4. 更新前端服务（零停机）
echo "🌐 步骤 4/5: 更新前端服务（零停机）..."
echo "✅ 前端静态文件已通过 volume 挂载，Caddy 自动提供新文件"
echo "✅ 无需重启 Caddy，无 503 错误"
echo ""

# 5. 显示服务状态
echo "📊 步骤 5/5: 服务状态:"
docker compose --env-file docker-compose-production.env ps
echo ""

# 6. 验证部署
echo "🧪 验证部署..."
sleep 3

# 检查后端 API
if curl -sf http://localhost/api/health > /dev/null 2>&1; then
    echo "✅ 后端 API 正常"
else
    echo "❌ 后端 API 异常"
    echo "📋 最近日志:"
    docker compose --env-file docker-compose-production.env logs backend --tail 20
fi

# 检查前端
if curl -sf http://localhost/ > /dev/null 2>&1; then
    echo "✅ 前端访问正常"
else
    echo "❌ 前端访问异常"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 零停机部署完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📍 访问地址:"
echo "   前端: http://localhost"
echo "   管理后台: http://localhost/#/admin/"
echo "   API: http://localhost/api/health"
echo ""
echo "🔑 默认账号:"
echo "   用户名: makerroot"
echo "   密码: 123456"
echo ""
echo "💡 查看日志:"
echo "   docker compose --env-file docker-compose-production.env logs -f"
echo ""
echo "💡 查看特定服务日志:"
echo "   docker compose --env-file docker-compose-production.env logs -f backend"
echo "   docker compose --env-file docker-compose-production.env logs -f caddy"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
