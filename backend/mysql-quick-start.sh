#!/bin/bash
# MySQL 数据库快速设置脚本

set -e

echo "🚀 MySQL 数据库快速设置"
echo "======================="
echo ""

# 检查 MySQL 是否安装
if ! command -v mysql &> /dev/null; then
    echo "❌ 错误: 未检测到 MySQL"
    echo ""
    echo "请先安装 MySQL:"
    echo "  macOS:   brew install mysql"
    echo "  Ubuntu:  sudo apt install mysql-server"
    echo "  Windows: https://dev.mysql.com/downloads/installer/"
    exit 1
fi

echo "✅ 检测到 MySQL 已安装"
echo ""

# 提示输入 MySQL root 密码
echo "请输入 MySQL root 密码（用于创建数据库）:"
read -s MYSQL_ROOT_PASSWORD

echo ""
echo "📝 创建数据库和用户..."

# 创建数据库和用户
mysql -u root -p"${MYSQL_ROOT_PASSWORD}" << EOF
-- 创建数据库
CREATE DATABASE IF NOT EXISTS lottery_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建专用用户（如果不存在）
CREATE USER IF NOT EXISTS 'lottery_user'@'localhost' IDENTIFIED BY 'lottery_password_2024';

-- 授权
GRANT ALL PRIVILEGES ON lottery_db.* TO 'lottery_user'@'localhost';
FLUSH PRIVILEGES;

-- 显示结果
SELECT 'Database setup completed!' AS Status;
SELECT DATABASE() AS CurrentDB;
SELECT COUNT(*) AS TableCount FROM information_schema.tables WHERE table_schema = 'lottery_db';
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 数据库创建成功！"
    echo ""
    echo "📋 数据库信息:"
    echo "   数据库名: lottery_db"
    echo "   用户名: lottery_user"
    echo "   密码: lottery_password_2024"
    echo ""
    echo "📝 请更新 .env 文件:"
    echo "   DATABASE_URL=lottery_user:lottery_password_2024@tcp(localhost:3306)/lottery_db?charset=utf8mb4&parseTime=True&loc=Local"
    echo ""
    echo "🎯 下一步:"
    echo "   1. 更新 backend/.env 文件中的 DATABASE_URL"
    echo "   2. 运行应用: cd backend && go run main.go"
else
    echo ""
    echo "❌ 数据库创建失败，请检查密码是否正确"
    exit 1
fi
