#!/bin/bash

# ============================================
# 快速回滚脚本 - 回滚到上一个版本
# ============================================

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⚠️  开始回滚到上一个版本..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 确认回滚
read -p "⚠️  警告：即将回滚到上一个版本，是否继续？(yes/no): " -r
echo
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ 回滚已取消"
    exit 0
fi

echo ""
echo "📦 步骤 1/4: 回滚代码..."
git reset --hard HEAD@{1}
echo "✅ 代码已回滚"
echo ""

echo "🔧 步骤 2/4: 重新构建后端..."
docker compose --env-file docker-compose-production.env build backend
echo "✅ 后端已构建"
echo ""

echo "🚀 步骤 3/4: 重启服务..."
docker compose --env-file docker-compose-production.env up -d --no-deps backend
docker compose --env-file docker-compose-production.env restart caddy
echo "✅ 服务已重启"
echo ""

echo "⏳ 步骤 4/4: 等待服务就绪..."
sleep 10
docker compose --env-file docker-compose-production.env ps
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 回滚完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 如果回滚后仍有问题，请检查日志:"
echo "   docker compose --env-file docker-compose-production.env logs -f"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
