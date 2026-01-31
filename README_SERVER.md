# 🚀 服务器部署指南

## 快速部署（3步完成）

### 1️⃣ 安装 Docker 和 Docker Compose

```bash
# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2️⃣ 克隆项目并配置

```bash
# 克隆项目
git clone <your-repo-url> /opt/lottery-system
cd /opt/lottery-system

# 创建配置文件
cp .env.production.template docker-compose-production.env

# 修改密码（重要！）
vim docker-compose-production.env
```

**必须修改**：
- `MYSQL_ROOT_PASSWORD`
- `MYSQL_PASSWORD`
- `REDIS_PASSWORD`
- `JWT_SECRET`

### 3️⃣ 一键部署

```bash
chmod +x docker-compose-deploy.sh
./docker-compose-deploy.sh
```

---

## 📍 访问地址

- **前端**: http://makerroot.com
- **HTTPS**: https://makerroot.com
- **管理后台**: https://makerroot.com/#/admin/
- **默认账号**: makerroot / 123456

---

## 💡 常用命令

```bash
# 部署/更新
git pull && ./docker-compose-deploy.sh

# 查看状态
docker-compose --env-file docker-compose-production.env ps

# 查看日志
docker-compose --env-file docker-compose-production.env logs -f

# 重启服务
docker-compose --env-file docker-compose-production.env restart

# 停止服务
docker-compose --env-file docker-compose-production.env down
```

---

## 📚 详细文档

查看 `DOCKER_DEPLOY.md` 获取完整部署指南。
