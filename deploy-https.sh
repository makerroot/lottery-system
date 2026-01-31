#!/bin/bash

# ============================================
# HTTPS 部署脚本 - 使用已有证书
# ============================================

set -e

echo "🔒 HTTPS 部署脚本"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 1. 检查证书文件
echo "📋 检查证书文件..."
if [ ! -f "/etc/letsencrypt/live/makerroot.com/fullchain.pem" ]; then
    echo "❌ 错误: 证书文件不存在"
    echo "   路径: /etc/letsencrypt/live/makerroot.com/fullchain.pem"
    exit 1
fi

if [ ! -f "/etc/letsencrypt/live/makerroot.com/privkey.pem" ]; then
    echo "❌ 错误: 私钥文件不存在"
    echo "   路径: /etc/letsencrypt/live/makerroot.com/privkey.pem"
    exit 1
fi

echo "✅ 证书文件存在"
ls -lh /etc/letsencrypt/live/makerroot.com/*.pem

# 2. 验证证书有效期
echo ""
echo "🔍 验证证书..."
EXPIRY=$(openssl x509 -in /etc/letsencrypt/live/makerroot.com/fullchain.pem -noout -date | grep 'notAfter' | sed 's/notAfter=//')
echo "   过期时间: $EXPIRY"

# 3. 拉取最新代码
echo ""
echo "📦 拉取最新代码..."
git pull origin main

# 4. 验证 frontend/dist
echo ""
echo "📋 检查 frontend/dist..."
if [ ! -d "frontend/dist" ]; then
    echo "❌ 错误: frontend/dist 不存在"
    echo "   请先在本地构建: npm run build"
    exit 1
fi
echo "✅ frontend/dist 存在"
du -sh frontend/dist/

# 5. 停止旧服务
echo ""
echo "🛑 停止旧服务..."
docker compose --env-file docker-compose-production.env down

# 6. 启动服务
echo ""
echo "🚀 启动服务..."
docker compose --env-file docker-compose-production.env up -d

# 7. 等待服务启动
echo ""
echo "⏳ 等待服务启动..."
sleep 10

# 8. 显示服务状态
echo ""
echo "📊 服务状态:"
docker compose --env-file docker-compose-production.env ps

# 9. 测试 HTTPS
echo ""
echo "🧪 测试 HTTPS 访问..."
sleep 3

if curl -skI https://localhost/ 2>&1 | grep -q "HTTP"; then
    echo "✅ HTTPS 访问正常"
else
    echo "❌ HTTPS 访问异常"
    echo "   查看日志: docker logs lottery-caddy"
fi

# 10. 测试健康检查
echo ""
if curl -s http://localhost/api/health | grep -q "正常"; then
    echo "✅ 后端 API 正常"
else
    echo "❌ 后端 API 异常"
fi

# 11. 显示访问信息
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ HTTPS 部署完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📍 访问地址:"
echo "   HTTP:  http://makerroot.com"
echo "   HTTPS: https://makerroot.com"
echo "   管理后台: https://makerroot.com/#/admin/"
echo ""
echo "🔒 证书信息:"
echo "   路径: /etc/letsencrypt/live/makerroot.com/"
echo "   过期: $EXPIRY"
echo ""
echo "🔑 默认账号:"
echo "   用户名: makerroot"
echo "   密码: 123456"
echo ""
echo "💡 查看日志:"
echo "   docker compose --env-file docker-compose-production.env logs -f caddy"
echo ""
echo "💡 证书续期（自动）:"
echo "   Let's Encrypt 证书会自动续期"
echo "   如需手动续期: docker compose restart caddy"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
