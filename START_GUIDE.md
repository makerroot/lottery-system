# 🚀 一键启动部署指南

## 快速开始

### 一键启动（推荐）

```bash
# 赋予执行权限（首次）
chmod +x start.sh

# 运行启动脚本
./start.sh
```

然后选择启动模式：
- **选项1**: Docker Compose模式（推荐，一键启动所有服务）
- **选项2**: 本地开发模式（需要安装Go和Node.js）

---

## 📖 启动模式说明

### 1. 🐳 Docker Compose模式（推荐）

**特点**:
- ✅ 一键启动所有服务
- ✅ 包含MySQL、Redis、后端、前端、Caddy
- ✅ 使用HTTPS（自签名证书）
- ✅ 适合生产环境和完整测试

**启动方式**:
```bash
# 交互式
./start.sh
# 选择: 1

# 或直接启动
./start.sh docker
```

**访问地址**:
- 前端: https://localhost
- 后端API: https://localhost/api/*
- 管理后台: https://localhost/admin/*

**默认账号**:
- 管理员: `makerroot` / `123456`

**首次访问提示**:
- 浏览器会显示"您的连接不是私密连接"（自签名证书警告）
- 点击"高级"→"继续访问"即可

### 2. 💻 本地开发模式

**特点**:
- ✅ 分别启动后端和前端
- ✅ 需要本地安装Go和Node.js
- ✅ 适合开发调试
- ✅ 热重载支持

**启动方式**:
```bash
# 交互式
./start.sh
# 选择: 2

# 或直接启动
./start.sh local
```

**访问地址**:
- 前端: http://localhost:5173
- 后端: http://localhost:8080
- 管理后台: http://localhost:5173/admin

### 3. 🔄 仅启动后端

```bash
./start.sh backend
# 或
./start.sh 3
```

### 4. 🎨 仅启动前端

```bash
./start.sh frontend
# 或
./start.sh 4
```

---

## 🛠️ 常用命令

### 服务管理

```bash
# 查看服务状态
./start.sh status

# 停止所有服务
./start.sh stop

# 重新启动
./start.sh docker  # 或 local
```

### Docker Compose命令

```bash
# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f mysql

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 查看服务状态
docker-compose ps
```

### 本地开发命令

```bash
# 查看后端日志
tail -f backend.log

# 查看前端日志
tail -f frontend.log

# 查看进程
ps aux | grep "go run main"
ps aux | grep "npm run dev"
```

---

## 📋 首次使用前准备

### Docker Compose模式

1. **安装Docker和Docker Compose**
   - Docker: https://docs.docker.com/get-docker/
   - Docker Compose: https://docs.docker.com/compose/install/

2. **修改环境变量**（必须）
   ```bash
   vim docker-compose-production.env
   ```

   修改以下配置：
   ```env
   MYSQL_ROOT_PASSWORD=your_strong_password
   MYSQL_PASSWORD=your_strong_password
   REDIS_PASSWORD=your_strong_password
   JWT_SECRET=your-super-secret-jwt-key-at-least-32-characters
   ```

3. **启动系统**
   ```bash
   ./start.sh docker
   ```

### 本地开发模式

1. **安装Go**
   - 下载: https://golang.org/dl/
   - 版本: Go 1.18+

2. **安装Node.js**
   - 下载: https://nodejs.org/
   - 版本: Node.js 16+

3. **配置后端**
   ```bash
   cd backend
   cp .env.example .env  # 或创建.env文件
   ```

4. **配置前端**
   ```bash
   cd frontend
   npm install
   ```

5. **启动系统**
   ```bash
   ./start.sh local
   ```

---

## 🎯 快速验证

### Docker Compose模式

```bash
# 1. 检查服务状态
docker-compose ps

# 2. 查看服务日志
docker-compose logs

# 3. 测试API
curl -k https://localhost/api/health

# 4. 浏览器访问
# 打开 https://localhost
# 点击"高级"→"继续访问"
```

### 本地开发模式

```bash
# 1. 检查端口
lsof -i :8080  # 后端
lsof -i :5173  # 前端

# 2. 测试API
curl http://localhost:8080/api/health

# 3. 浏览器访问
# 打开 http://localhost:5173
```

---

## ⚠️ 常见问题

### Q1: Docker服务未启动

**错误**: `Cannot connect to the Docker daemon`

**解决**:
```bash
# 启动Docker
# macOS: 打开Docker Desktop
# Linux: sudo systemctl start docker
```

### Q2: 端口已被占用

**错误**: `Bind for 0.0.0.0:80 failed: port is already allocated`

**解决**:
```bash
# 查看占用进程
lsof -i :80
lsof -i :443

# 停止占用进程
./start.sh stop

# 或修改docker-compose.yml中的端口映射
```

### Q3: HTTPS访问显示安全警告

**原因**: 使用自签名证书

**解决**: 这是正常的！点击"高级"→"继续访问"即可

### Q4: 后端无法连接数据库

**解决**:
```bash
# 检查MySQL容器
docker-compose ps mysql

# 查看MySQL日志
docker-compose logs mysql

# 重启MySQL
docker-compose restart mysql
```

### Q5: 前端无法连接后端

**解决**:
```bash
# 检查CORS配置
cat docker-compose-production.env | grep ALLOWED_ORIGINS

# 确保包含:
# ALLOWED_ORIGINS=https://localhost,http://localhost:5173
```

---

## 📚 相关文档

- **HTTPS配置指南**: [HTTPS_MODE_GUIDE.md](HTTPS_MODE_GUIDE.md)
- **Docker部署指南**: [DOCKER_DEPLOYMENT_GUIDE.md](DOCKER_DEPLOYMENT_GUIDE.md)
- **快速参考**: [DOCKER_README.md](DOCKER_README.md)
- **后端架构**: [backend/ARCHITECTURE.md](backend/ARCHITECTURE.md)
- **API文档**: [backend/API.md](backend/API.md)

---

## 🎊 开始使用

```bash
# 1. 修改配置（首次）
vim docker-compose-production.env

# 2. 启动系统
./start.sh docker

# 3. 访问系统
# 浏览器打开: https://localhost
# 点击"高级"→"继续访问"
# 登录: makerroot / 123456
```

**就这么简单！🚀**
