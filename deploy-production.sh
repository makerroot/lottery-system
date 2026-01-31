#!/bin/bash

# ============================================
# 服务器部署脚本 - 使用本地构建的前端
# ============================================

set -e

echo "🚀 开始部署..."
echo ""

# 1. 拉取最新代码
echo "📦 拉取最新代码..."
git pull origin main

# 2. 验证 dist 存在
echo ""
echo "📋 检查 frontend/dist..."
if [ ! -d "frontend/dist" ]; then
    echo "❌ 错误: frontend/dist 不存在"
    echo "请先在本地构建前端: npm run build"
    exit 1
fi
echo "✅ frontend/dist 存在"
du -sh frontend/dist/

# 3. 停止旧服务
echo ""
echo "🛑 停止旧服务..."
docker compose --env-file docker-compose-production.env down

# 4. 清理旧 volume（可选）
echo ""
read -p "是否删除旧的 frontend-dist volume? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🧹 清理旧 volume..."
    docker volume rm lottery-frontend-dist 2>/dev/null || echo "Volume 不存在"
fi

# 5. 启动服务
echo ""
echo "🚀 启动服务..."
docker compose --env-file docker-compose-production.env up -d

# 6. 等待服务启动
echo ""
echo "⏳ 等待服务启动..."
sleep 10

# 7. 显示服务状态
echo ""
echo "📊 服务状态:"
docker compose --env-file docker-compose-production.env ps

# 8. 验证部署
echo ""
echo "🧪 验证部署..."
if curl -s http://localhost/api/health > /dev/null; then
    echo "✅ 后端 API 正常"
else
    echo "❌ 后端 API 异常"
fi

if curl -s http://localhost/ > /dev/null; then
    echo "✅ 前端访问正常"
else
    echo "❌ 前端访问异常"
fi

# 9. 显示访问信息
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 部署完成！"
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
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
