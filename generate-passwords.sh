#!/bin/bash

# ============================================
# 自动生成并更新密码配置
# ============================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ENV_FILE="docker-compose-production.env"

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}自动生成强密码${NC}"
echo -e "${BLUE}============================================${NC}"
echo

# 检查文件是否存在
if [ ! -f "$ENV_FILE" ]; then
    if [ -f ".env.production.template" ]; then
        echo -e "${YELLOW}从模板创建 $ENV_FILE${NC}"
        cp .env.production.template $ENV_FILE
    else
        echo -e "${RED}错误：找不到配置文件${NC}"
        exit 1
    fi
fi

# 生成强密码
JWT_SECRET=$(openssl rand -base64 32)
MYSQL_ROOT_PASSWORD=$(openssl rand -base64 16)
MYSQL_PASSWORD=$(openssl rand -base64 16)
REDIS_PASSWORD=$(openssl rand -base64 16)

echo -e "${GREEN}生成的密码：${NC}"
echo
echo "JWT_SECRET=$JWT_SECRET"
echo "MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD"
echo "MYSQL_PASSWORD=$MYSQL_PASSWORD"
echo "REDIS_PASSWORD=$REDIS_PASSWORD"
echo

# 询问是否更新
read -p "是否更新配置文件？ (y/N): " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo -e "${YELLOW}已取消${NC}"
    exit 0
fi

# 备份原文件
BACKUP_FILE="${ENV_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
cp $ENV_FILE $BACKUP_FILE
echo -e "${GREEN}已备份到: $BACKUP_FILE${NC}"

# 更新配置文件
sed -i.bak "s/^JWT_SECRET=.*/JWT_SECRET=$JWT_SECRET/" $ENV_FILE
sed -i.bak "s/^MYSQL_ROOT_PASSWORD=.*/MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD/" $ENV_FILE
sed -i.bak "s/^MYSQL_PASSWORD=.*/MYSQL_PASSWORD=$MYSQL_PASSWORD/" $ENV_FILE
sed -i.bak "s/^REDIS_PASSWORD=.*/REDIS_PASSWORD=$REDIS_PASSWORD/" $ENV_FILE

# 删除临时文件
rm -f ${ENV_FILE}.bak

echo
echo -e "${GREEN}============================================${NC}"
echo -e "${GREEN}密码已成功更新！${NC}"
echo -e "${GREEN}============================================${NC}"
echo
echo -e "${YELLOW}重要提示：${NC}"
echo "1. 请妥善保管生成的密码"
echo "2. 建议将密码保存到密码管理器"
echo "3. 配置文件: $ENV_FILE"
echo

# 显示更新的配置
echo -e "${BLUE}已更新的配置：${NC}"
grep "PASSWORD\|SECRET" $ENV_FILE | grep -v "^#"
