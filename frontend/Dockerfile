# ============================================
# 前端 Dockerfile - 仅构建静态文件
# ============================================
# 说明：只构建 Vue 静态文件，不包含 Web 服务器
# 静态文件由 Caddy 直接提供服务
# ============================================

# 构建阶段
FROM node:18

# 设置工作目录
WORKDIR /app

# 复制package文件
COPY package*.json ./

# 配置npm镜像源（国内镜像，更快更稳定）
RUN npm config set registry https://registry.npmmirror.com

# 安装依赖
RUN npm install

# 复制源代码
COPY . .

# 构建生产版本
RUN npm run build

# 验证构建产物
RUN ls -la dist/ && test -f dist/index.html

# 构建完成，dist 目录即为最终产物
# 此镜像仅用于构建，不作为运行时容器
