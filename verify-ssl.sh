#!/bin/bash

# ============================================
# SSL证书验证脚本
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✅]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[⚠️]${NC} $1"
}

log_error() {
    echo -e "${RED}[❌]${NC} $1"
}

DOMAIN="makerroot.com"
CERT_PATH="/etc/letsencrypt/live/${DOMAIN}"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}    SSL证书验证 - ${DOMAIN}${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 1. 检查证书文件存在性
log_info "1. 检查证书文件..."
if [ -f "${CERT_PATH}/fullchain.pem" ]; then
    log_success "fullchain.pem 存在"
else
    log_error "fullchain.pem 不存在"
    log_error "请先获取证书: sudo certbot certonly --standalone -d ${DOMAIN}"
    exit 1
fi

if [ -f "${CERT_PATH}/privkey.pem" ]; then
    log_success "privkey.pem 存在"
else
    log_error "privkey.pem 不存在"
    log_error "请先获取证书: sudo certbot certonly --standalone -d ${DOMAIN}"
    exit 1
fi

# 2. 检查证书权限
log_info "2. 检查证书权限..."
FULL_PERM=$(stat -c %a "${CERT_PATH}/fullchain.pem")
PRIV_PERM=$(stat -c %a "${CERT_PATH}/privkey.pem")

echo "  fullchain.pem: ${FULL_PERM}"
echo "  privkey.pem: ${PRIV_PERM}"

if [[ "$FULL_PERM" == "644" || "$FULL_PERM" == "444" ]]; then
    log_success "证书文件权限正常"
elif [[ "$FULL_PERM" == "600" || "$PRIV_PERM" == "600" ]]; then
    log_success "私钥权限正确"
else
    log_warning "权限不是标准值，可能需要修改"
    echo ""
    echo "  建议执行:"
    echo "    sudo chmod 644 ${CERT_PATH}/fullchain.pem"
    echo "    sudo chmod 600 ${CERT_PATH}/privkey.pem"
fi

# 3. 检查证书有效期
log_info "3. 检查证书有效期..."
EXPIRY_DATE=$(openssl x509 -in "${CERT_PATH}/fullchain.pem" -noout -enddate | cut -d= -f2)
if [ -n "$EXPIRY_DATE" ]; then
    log_success "证书有效期至: ${EXPIRY_DATE}"
else
    log_error "无法读取证书有效期"
fi

# 4. 检查Caddy配置
log_info "4. 检查Caddy配置..."
if grep -q "tls /etc/letsencrypt/live/${DOMAIN}" docker/caddy/Caddyfile; then
    log_success "Caddyfile已配置证书路径"
else
    log_error "Caddyfile未配置证书路径"
    exit 1
fi

# 5. 检查Docker挂载
log_info "5. 检查Docker挂载配置..."
if grep -q "/etc/letsencrypt:/etc/letsencrypt:ro" docker-compose.yml; then
    log_success "Docker Compose已配置证书挂载"
else
    log_error "Docker Compose未配置证书挂载"
    exit 1
fi

# 6. 测试证书链
log_info "6. 测试证书链..."
if openssl s_client -connect ${DOMAIN}:443 -servername ${DOMAIN} </dev/null 2>&1 | grep -q "Verify return code: 0"; then
    log_success "证书链有效"
else
    log_warning "证书链测试失败（可能服务未启动）"
fi

# 7. 检查Caddy服务
log_info "7. 检查Caddy服务..."
if docker-compose ps caddy 2>/dev/null | grep -q "Up"; then
    log_success "Caddy服务正在运行"
else
    log_warning "Caddy服务未运行"
    echo ""
    echo "  启动服务:"
    echo "    docker-compose up -d"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}     ✅ 验证完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}下一步操作:${NC}"
echo "  1. 重启Caddy: docker-compose restart caddy"
echo "  2. 访问测试: https://${DOMAIN}"
echo "  3. 查看日志: docker-compose logs -f caddy"
echo ""
echo -e "${BLUE}证书续期:${NC}"
echo "  sudo certbot renew"
echo ""

# 8. SSL Labs测试提示
echo -e "${BLUE}SSL Labs测试:${NC}"
echo "  https://www.ssllabs.com/ssltest/analyze.html?d=${DOMAIN}"
echo ""
