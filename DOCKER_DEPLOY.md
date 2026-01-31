# 🚀 Docker Compose 部署完整指南

## 📋 部署流程

```bash
┌─────────────────────────────────────┐
│ 1. 准备服务器环境                   │
│    - 安装 Docker                     │
│    - 安装 Docker Compose              │
└─────────────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│ 2. 准备配置文件                     │
│    - docker-compose-production.env   │
│    - 修改密码                        │
└─────────────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│ 3. 克隆项目                         │
│    - git clone <repo>                │
│    - cd lottery-system               │
└─────────────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│ 4. 一键部署                          │
│    - ./docker-compose-deploy.sh      │
└─────────────────────────────────────┘
```

---

## 🔧 第一步：准备服务器环境

### 安装 Docker

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# CentOS/RHEL
sudo yum install -y docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker

# 验证安装
docker --version
```

### 安装 Docker Compose

```bash
# Linux
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose --version
```

---

## 📝 第二步：准备配置文件

### 创建环境配置

```bash
# 进入项目目录
cd /opt/lottery-system

# 从模板创建配置文件
cp .env.production.template docker-compose-production.env

# 编辑配置（必须修改密码）
vim docker-compose-production.env
```

### 修改以下配置

```env
# ⚠️ 必须修改这些密码！

# MySQL
MYSQL_ROOT_PASSWORD=your_strong_password_here
MYSQL_PASSWORD=your_strong_password_here

# Redis
REDIS_PASSWORD=your_strong_password_here

# JWT
JWT_SECRET=your-super-secret-jwt-key-at-least-32-characters-long
```

---

## 🎯 第三步：一键部署

### 方式1：使用自动化脚本（推荐）⭐

```bash
# 给脚本执行权限
chmod +x docker-compose-deploy.sh

# 运行部署脚本
./docker-compose-deploy.sh
```

**脚本会自动完成**：
- ✅ 检查 Docker 环境
- ✅ 检查配置文件
- ✅ 检查证书
- ✅ 验证 frontend/dist
- ✅ 拉取最新代码
- ✅ 启动所有服务
- ✅ 检查服务健康
- ✅ 测试访问

### 方式2：手动部署

```bash
# 1. 拉取最新代码
git pull origin main

# 2. 启动服务
docker-compose --env-file docker-compose-production.env up -d

# 3. 查看状态
docker-compose --env-file docker-compose-production.env ps

# 4. 查看日志
docker-compose --env-file docker-compose-production.env logs -f
```

---

## 📊 部署后的服务

### 服务列表

| 服务 | 容器名 | 端口 | 说明 |
|------|--------|------|------|
| **MySQL** | lottery-mysql | 3306 | 数据库 |
| **Redis** | lottery-redis | 6379 | 缓存 |
| **Backend** | lottery-backend | 8080 | API服务 |
| **Caddy** | lottery-caddy | 80, 443 | 反向代理 + 前端 |

### 网络架构

```
用户 (80/443)
    ↓
Caddy (反向代理)
    ├─→ /api/* → Backend (8080)
    └─→ /admin/* → Backend (8080)
    └─→ /* → frontend/dist (静态文件)
```

---

## ✅ 验证部署

### 1. 检查服务状态

```bash
docker-compose --env-file docker-compose-production.env ps
```

**预期输出**：
```
NAME              STATUS
lottery-mysql     Up (healthy)
lottery-redis     Up (healthy)
lottery-backend   Up (healthy)
lottery-caddy     Up (running)
```

### 2. 查看日志

```bash
# 查看所有服务日志
docker-compose --env-file docker-compose-production.env logs -f

# 查看特定服务
docker-compose --env-file docker-compose-production.env logs -f caddy
docker-compose --env-file docker-compose-production.env logs -f backend
```

### 3. 测试访问

```bash
# HTTP 访问
curl -I http://localhost/

# HTTPS 访问（如果配置了证书）
curl -I https://localhost/

# API 测试
curl http://localhost/api/health

# 浏览器访问
open http://makerroot.com
```

---

## 🔄 日常管理

### 查看服务状态

```bash
docker-compose --env-file docker-compose-production.env ps
```

### 查看日志

```bash
# 实时日志
docker-compose --env-file docker-compose-production.env logs -f

# 最近100行
docker-compose --env-file docker-compose-production.env logs --tail=100

# 特定服务
docker-compose --env-file docker-compose-production.env logs -f caddy
```

### 重启服务

```bash
# 重启所有服务
docker-compose --env-file docker-compose-production.env restart

# 重启特定服务
docker-compose --env-file docker-compose-production.env restart caddy
docker-compose --env-file docker-compose-production.env restart backend
```

### 停止服务

```bash
# 停止所有服务
docker-compose --env-file docker-compose-production.env down

# 停止并删除数据（危险！）
docker-compose --env-file docker-compose-production.env down -v
```

### 更新服务

```bash
# 更新代码
git pull origin main

# 重新构建并启动（如果修改了代码）
docker-compose --env-file docker-compose-production.env up -d --build backend

# 仅重启 Caddy（如果更新了前端）
docker-compose --env-file docker-compose-production.env restart caddy
```

---

## 🐛 故障排查

### 问题1：服务启动失败

```bash
# 查看详细日志
docker-compose --env-file docker-compose-production.env logs

# 检查配置文件
cat docker-compose-production.env

# 检查端口占用
sudo lsof -i :80
sudo lsof -i :443
```

### 问题2：数据库连接失败

```bash
# 查看 MySQL 日志
docker-compose --env-file docker-compose-production.env logs mysql

# 检查 MySQL 容器
docker exec -it lottery-mysql mysql -u root -p

# 检查网络
docker network ls
docker network inspect lottery-network
```

### 问题3：前端无法访问

```bash
# 检查 frontend/dist
ls -la frontend/dist/

# 检查 Caddy 挂载
docker exec lottery-caddy ls -la /usr/share/caddy/frontend/

# 查看 Caddy 日志
docker logs lottery-caddy 2>&1 | tail -50
```

### 问题4：HTTPS 证书问题

```bash
# 检查证书文件
ls -la /etc/letsencrypt/live/makerroot.com/

# 查看 Caddy TLS 配置
docker-compose --env-file docker-compose-production.env logs caddy 2>&1 | grep -i "certificate\|tls"

# 重启 Caddy
docker-compose --env-file docker-compose-production.env restart caddy
```

---

## 📝 配置文件说明

### docker-compose.yml

- 定义了所有服务
- 使用本地 frontend/dist
- 通过 Caddy 提供服务

### docker-compose-production.env

- 环境变量配置
- 数据库密码
- Redis 密码
- JWT 密钥
- 服务器配置

### Caddyfile

- Caddy 反向代理配置
- HTTPS/SSL 配置
- 静态文件服务

---

## 🎯 快速命令参考

```bash
# 部署
./docker-compose-deploy.sh

# 状态
docker-compose --env-file docker-compose-production.env ps

# 日志
docker-compose --env-file docker-compose-production.env logs -f

# 重启
docker-compose --env-file docker-compose-production.env restart

# 停止
docker-compose --env-file docker-compose-production.env down

# 更新
git pull && docker-compose --env-file docker-compose-production.env up -d
```

---

## 🔐 安全建议

### 首次部署后必须做

1. ✅ 修改所有默认密码
2. ✅ 检查防火墙配置
3. ✅ 验证 HTTPS 正常工作
4. ✅ 登录管理后台修改密码
5. ✅ 定期备份数据

### 数据备份

```bash
# 备份 MySQL
docker exec lottery-mysql mysqldump -u root -p lottery_system > backup.sql

# 备份 Redis
docker exec lottery-redis redis-cli --rdb /data/backup.rdb

# 备份数据卷
docker run --rm -v lottery-system_mysql-data:/data -v $(pwd):/backup alpine tar czf backup.tar.gz /data
```

---

## 🎉 总结

**一键部署命令**：
```bash
./docker-compose-deploy.sh
```

**3步完成部署**：
1. 安装 Docker 和 Docker Compose
2. 修改 `docker-compose-production.env` 中的密码
3. 运行 `./docker-compose-deploy.sh`

**访问地址**：
- 前端: http://makerroot.com
- HTTPS: https://makerroot.com
- 管理后台: https://makerroot.com/#/admin/

**默认账号**：
- 用户名: makerroot
- 密码: 123456
