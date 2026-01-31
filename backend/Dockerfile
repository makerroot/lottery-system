# ============================================
# 后端 Dockerfile - 多阶段构建（生产环境）
# ============================================

# 阶段1: 构建阶段
FROM golang:1.18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 禁用 Go 模块缓存和构建缓存
ENV GOCACHE=/tmp/go-cache
ENV GOMODCACHE=/tmp/go-mod

# 安装必要的构建工具和MySQL客户端库
RUN apk add --no-cache git gcc musl-dev mysql-client

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖（不使用缓存）
RUN go mod download -x

# 复制源代码
COPY . .

# 构建应用（启用CGO以支持MySQL，不使用构建缓存）
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# 阶段2: 运行阶段
FROM alpine:latest

# 安装ca证书、MySQL客户端和wget（用于健康检查）
RUN apk --no-cache add ca-certificates mysql-client wget tzdata

# 设置时区为上海
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 创建非root用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 创建日志目录
RUN mkdir -p /app/logs && chown -R appuser:appuser /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# 启动应用
CMD ["./main"]
