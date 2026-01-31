# Caddy 反向代理配置说明

## 📋 配置文件说明

本Caddy配置文件提供了两种部署模式：
- **HTTP模式** (默认) - 适用于开发/测试环境
- **HTTPS模式** (可选) - 适用于生产环境

---

## 🔧 当前配置：HTTP模式（默认）

### 访问地址
- 前端: http://localhost
- 后端API: http://localhost/api/*
- 管理后台: http://localhost/admin/*

### 特点
- ✅ 无需SSL证书
- ✅ 开箱即用
- ✅ 适合内网或开发环境
- ✅ 自动禁用浏览器缓存（便于调试）

### 启动方式
```bash
# 使用默认HTTP配置即可
./deploy.sh start
```

---

## 🔒 生产环境：HTTPS模式

### 启用HTTPS的两种方式

#### 方式1: 使用 Let's Encrypt 自动证书（推荐）

1. **确保域名已解析到服务器**
   ```bash
   # 检查域名解析
   nslookup your-domain.com
   ```

2. **编辑Caddyfile，取消HTTPS部分的注释**
   ```bash
   vim docker/caddy/Caddyfile
   ```

3. **修改配置**：
   ```caddy
   # 1. 注释掉或删除 :80 部分
   # :80 { ... }

   # 2. 取消 :443 部分的注释
   :443 {
       # 3. 修改为你的域名
       your-domain.com {
           # 4. 修改为你的邮箱
           tls your-email@example.com

           # ... 其他配置保持不变
       }

       # 5. 取消HTTP到HTTPS重定向的注释
       your-domain.com:80 {
           redir https://your-domain.com{uri} permanent
       }
   }
   ```

4. **重启Caddy**
   ```bash
   docker-compose restart caddy
   ```

5. **访问测试**
   ```bash
   curl -I https://your-domain.com
   ```

#### 方式2: 使用已有SSL证书

1. **准备证书文件**
   ```bash
   # 将证书文件放到 docker/caddy/ssl/ 目录
   mkdir -p docker/caddy/ssl
   cp /path/to/fullchain.pem docker/caddy/ssl/
   cp /path/to/privkey.pem docker/caddy/ssl/
   ```

2. **更新docker-compose.yml，挂载证书目录**
   ```yaml
   caddy:
     volumes:
       - ./docker/caddy/Caddyfile:/etc/caddy/Caddyfile:ro
       - ./docker/caddy/ssl:/etc/caddy/ssl:ro  # 添加这行
       - ./docker/caddy/data:/data
       - ./docker/caddy/logs:/var/log/caddy
   ```

3. **编辑Caddyfile**
   ```caddy
   :443 {
       your-domain.com {
           # 使用已有证书
           tls /etc/caddy/ssl/fullchain.pem /etc/caddy/ssl/privkey.pem

           # ... 其他配置
       }

       # HTTP重定向到HTTPS
       your-domain.com:80 {
           redir https://your-domain.com{uri} permanent
       }
   }
   ```

4. **重启服务**
   ```bash
   docker-compose restart caddy
   ```

---

## 📝 配置详解

### 全局配置

```caddy
{
    admin off                          # 关闭管理API
    auto_https disable_redirects       # 禁用自动HTTPS（开发环境）
    log {
        output file /var/log/caddy/access.log {
            roll_size 50MiB            # 日志文件大小限制
            roll_keep 14               # 保留14个日志文件
        }
        format json                    # JSON格式日志
        level INFO                     # 日志级别
    }
}
```

### 路由处理

```caddy
# API接口 - 代理到后端
handle /api/* /admin/* {
    reverse_proxy backend:8080 {
        health_uri /api/health         # 健康检查
        header_up X-Forwarded-Proto {scheme}
    }
}

# 健康检查
handle /health {
    reverse_proxy backend:8080
}

# 前端SPA - 静态资源服务
handle {
    root * /usr/share/nginx/html

    # 静态资源：文件不存在返回404
    @static_files {
        path *.js *.css *.png *.jpg *.jpeg *.gif *.svg *.woff *.woff2 *.ico
    }
    handle @static_files {
        file_server
    }

    # SPA路由：其他请求返回index.html
    handle {
        try_files {path} /index.html
        file_server
    }
}
```

### 安全响应头

```caddy
header {
    # 开发环境：禁用缓存
    Cache-Control "no-cache, no-store, must-revalidate"
    Pragma "no-cache"
    Expires "0"

    # 安全头
    X-Frame-Options "SAMEORIGIN"              # 防止点击劫持
    X-Content-Type-Options "nosniff"          # 防止MIME类型嗅探
    Referrer-Policy "no-referrer-when-downgrade"
    X-XSS-Protection "1; mode=block"          # XSS保护

    # HTTPS专用（生产环境启用）
    # Strict-Transport-Security "max-age=31536000; includeSubDomains"
}
```

---

## 🚀 快速切换HTTP/HTTPS

### HTTP → HTTPS（生产环境）

```bash
# 1. 备份当前配置
cp docker/caddy/Caddyfile docker/caddy/Caddyfile.http.bak

# 2. 编辑Caddyfile
vim docker/caddy/Caddyfile

# 3. 注释掉 :80 部分，取消 :443 部分的注释
# 4. 修改域名和邮箱
# 5. 重启Caddy
docker-compose restart caddy

# 6. 验证HTTPS
curl -I https://your-domain.com
```

### HTTPS → HTTP（回退）

```bash
# 1. 恢复HTTP配置
cp docker/caddy/Caddyfile.http.bak docker/caddy/Caddyfile

# 2. 重启Caddy
docker-compose restart caddy

# 3. 验证HTTP
curl -I http://localhost
```

---

## 🔍 测试和验证

### 测试HTTP配置

```bash
# 测试前端访问
curl -I http://localhost/

# 测试API代理
curl -I http://localhost/api/health

# 测试管理后台
curl -I http://localhost/admin/companies

# 查看Caddy日志
docker-compose logs -f caddy
```

### 测试HTTPS配置

```bash
# 测试HTTPS访问
curl -I https://your-domain.com

# 测试SSL证书
openssl s_client -connect your-domain.com:443 -servername your-domain.com

# 检查证书有效期
echo | openssl s_client -connect your-domain.com:443 2>/dev/null | openssl x509 -noout -dates

# SSL Labs测试（浏览器访问）
# https://www.ssllabs.com/ssltest/analyze.html?d=your-domain.com
```

---

## 🛠️ 常见问题

### Q1: Let's Encrypt证书申请失败？

**原因**:
- 域名未正确解析到服务器
- 80端口被占用或防火墙阻止
- DNS传播未完成

**解决**:
```bash
# 1. 检查域名解析
nslookup your-domain.com

# 2. 检查80端口
lsof -i :80

# 3. 检查防火墙
sudo ufw status
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 4. 查看Caddy日志
docker-compose logs caddy | grep "tls"
```

### Q2: SPA路由刷新404？

**原因**: 静态资源被错误地重定向到index.html

**解决**: 确保Caddyfile中正确配置了静态资源处理
```caddy
@static_files {
    path *.js *.css *.png *.jpg *.jpeg *.gif *.svg *.woff *.woff2 *.ico
}
handle @static_files {
    file_server  # 不使用try_files
}
```

### Q3: API跨域问题？

**解决**: Caddy已自动添加CORS头，如果仍有问题，检查后端配置
```bash
# 检查后端CORS配置
docker-compose exec backend env | grep ALLOWED_ORIGINS
```

### Q4: 如何查看访问日志？

```bash
# 实时查看日志
docker-compose exec -T caddy tail -f /var/log/caddy/access.log

# 查看错误日志
docker-compose logs caddy | grep error

# 日志文件位置
ls -lh docker/caddy/logs/
```

---

## 📊 性能优化

### 启用HTTP/2

HTTPS模式自动启用HTTP/2，无需额外配置。

### 静态资源缓存（生产环境）

修改响应头配置：
```caddy
header {
    # 生产环境：启用缓存
    Cache-Control "public, max-age=3600"

    # 静态资源长缓存
    @static {
        path *.js *.css *.png *.jpg *.jpeg *.gif *.svg *.woff *.woff2
    }
    header @static Cache-Control "public, max-age=86400"
}
```

### Gzip压缩

Caddy默认启用了gzip压缩，无需额外配置。

---

## 🔐 安全建议

### 1. 启用HTTPS（生产环境必须）

```caddy
# 强制HTTPS
:80 {
    redir https://your-domain.com{uri} permanent
}
```

### 2. 启用HSTS

```caddy
header {
    Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
}
```

### 3. 限制请求大小

```caddy
reverse_proxy backend:8080 {
    # 限制请求体大小（防止大文件攻击）
    header_up X-Forwarded-Proto {scheme}
}
```

### 4. 速率限制（使用Caddy插件）

需要编译Caddy时包含ratelimit插件。

---

## 📚 参考资源

- [Caddy官方文档](https://caddyserver.com/docs/)
- [Caddyfile概念](https://caddyserver.com/docs/caddyfile/concepts)
- [Let's Encrypt文档](https://letsencrypt.org/docs/)
- [SSL Labs测试](https://www.ssllabs.com/ssltest/)

---

**配置版本**: 1.0
**最后更新**: 2026-01-24
**维护者**: Lottery System Team
