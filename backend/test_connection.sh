#!/bin/bash
cd backend
echo "🔍 测试数据库连接..."
echo ""

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "❌ 错误: .env 文件不存在"
    exit 1
fi

# 读取 DATABASE_URL
DATABASE_URL=$(grep DATABASE_URL .env | cut -d '=' -f2)

if [ -z "$DATABASE_URL" ]; then
    echo "❌ 错误: .env 中未设置 DATABASE_URL"
    exit 1
fi

echo "📝 当前数据库配置:"
if [[ "$DATABASE_URL" == *"@tcp("* ]]; then
    echo "   类型: MySQL"
    echo "   连接: $(echo $DATABASE_URL | sed 's/:[^@]*@/:****@/')"
else
    echo "   类型: SQLite"
    echo "   文件: $DATABASE_URL"
fi
echo ""

# 检查是否是 MySQL
if [[ "$DATABASE_URL" == *"@tcp("* ]]; then
    echo "✅ 检测到 MySQL 配置"
    echo ""
    echo "📋 下一步:"
    echo "   1. 确保 MySQL 服务已启动"
    echo "   2. 确保 MySQL 中已创建 lottery_db 数据库"
    echo "   3. 运行: go run main.go"
    echo ""
    echo "💡 创建数据库命令:"
    echo "   mysql -u root -p -e 'CREATE DATABASE lottery_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'"
else
    echo "✅ 检测到 SQLite 配置"
    echo ""
    echo "📋 下一步:"
    echo "   运行: go run main.go"
fi
