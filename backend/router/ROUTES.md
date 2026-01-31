# API 路由文档

本文档描述了抽奖系统后端的所有 API 路由。

---

## 基础信息

- **路由模块**: `router/router.go`
- **路由配置**: 分为用户端 API 和管理后台 API
- **认证方式**: JWT Token

---

## 🏥 健康检查端点

### `GET /api/health` 和 `HEAD /api/health`

**描述**: 服务健康检查

**认证**: 无需认证

**响应**: 200 OK

---

## 👤 用户端 API (`/api`)

### 🔓 公开接口（无需认证）

#### `POST /api/register`

**描述**: 用户注册或登录

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Query 参数**:
- `company_code` (必填): 公司代码

#### `POST /api/unified-login`

**描述**: 统一登录接口（支持用户和管理员）

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Query 参数**:
- `company_code` (必填): 公司代码

#### `POST /api/login`

**描述**: 统一登录接口（alias）

**同 `/api/unified-login`**

---

### 🔒 需要用户认证的接口

#### 公司信息

##### `GET /api/company`

**描述**: 获取公司信息

**Query 参数**:
- `company_code` (必填): 公司代码

#### 用户信息

##### `GET /api/user`

**描述**: 获取用户信息

**Query 参数**:
- `phone` (必填): 用户手机号
- `company_code` (必填): 公司代码

##### `POST /api/user/change-password`

**描述**: 修改用户密码

**请求体**:
```json
{
  "old_password": "string",
  "new_password": "string"
}
```

#### 奖品相关

##### `GET /api/prize-levels`

**描述**: 获取启用的奖项等级

**Query 参数**:
- `company_code` (必填): 公司代码

#### 抽奖相关

##### `POST /api/draw`

**描述**: 执行抽奖

**Query 参数**:
- `company_code` (必填): 公司代码

**请求体**:
```json
{
  "level_id": 0,
  "count": 1,
  "user_phone": "string"
}
```

##### `GET /api/my-prize`

**描述**: 获取我的奖品

**Query 参数**:
- `phone` (必填): 用户手机号
- `company_code` (必填): 公司代码

##### `GET /api/user-stats`

**描述**: 获取用户统计

**Query 参数**:
- `company_code` (必填): 公司代码

##### `GET /api/draw-records`

**描述**: 获取抽奖记录（公开）

**Query 参数**:
- `company_code` (必填): 公司代码

##### `GET /api/available-users`

**描述**: 获取未抽奖的用户列表

**Query 参数**:
- `company_code` (必填): 公司代码

---

## 🔐 管理后台 API (`/admin`)

### 🔓 登录接口（无需认证）

#### `POST /admin/login`

**描述**: 管理员登录

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

---

### 🔒 需要管理员认证的接口

#### 管理员信息

##### `GET /admin/info`

**描述**: 获取当前管理员信息

##### `POST /admin/change-password`

**描述**: 修改管理员密码

**请求体**:
```json
{
  "old_password": "string",
  "new_password": "string"
}
```

#### 管理员管理（仅超级管理员）

##### `GET /admin/admins`

**描述**: 获取管理员列表

##### `POST /admin/admins`

**描述**: 创建管理员

**请求体**:
```json
{
  "username": "string",
  "password": "string",
  "name": "string",
  "company_id": 0
}
```

##### `PUT /admin/admins/:id`

**描述**: 更新管理员

**路径参数**:
- `id`: 管理员 ID

##### `DELETE /admin/admins/:id`

**描述**: 删除管理员

**路径参数**:
- `id`: 管理员 ID

#### 用户管理

##### `GET /admin/users`

**描述**: 获取用户列表

**Query 参数**:
- `page`: 页码
- `page_size`: 每页数量
- `company_id`: 公司 ID
- `search`: 搜索关键词

##### `POST /admin/users`

**描述**: 创建用户

**请求体**:
```json
{
  "username": "string",
  "password": "string",
  "name": "string",
  "phone": "string",
  "company_id": 0
}
```

##### `POST /admin/users/batch`

**描述**: 批量创建用户

**请求体**:
```json
{
  "users": [
    {
      "username": "string",
      "password": "string",
      "name": "string",
      "phone": "string"
    }
  ],
  "company_id": 0
}
```

##### `PUT /admin/users/:id`

**描述**: 更新用户

**路径参数**:
- `id`: 用户 ID

##### `DELETE /admin/users/:id`

**描述**: 删除用户

**路径参数**:
- `id`: 用户 ID

#### 公司管理（超级管理员）

##### `GET /admin/companies`

**描述**: 获取公司列表

##### `POST /admin/companies`

**描述**: 创建公司

**请求体**:
```json
{
  "name": "string",
  "code": "string",
  "theme_color": "string"
}
```

##### `PUT /admin/companies/:id`

**描述**: 更新公司

**路径参数**:
- `id`: 公司 ID

##### `DELETE /admin/companies/:id`

**描述**: 删除公司

**路径参数**:
- `id`: 公司 ID

##### `GET /admin/company-stats`

**描述**: 获取公司统计信息

**Query 参数**:
- `company_id`: 公司 ID

#### 奖项等级管理

##### `POST /admin/prize-levels`

**描述**: 创建奖项等级

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "probability": 0.0,
  "total_stock": 100,
  "sort_order": 1,
  "company_id": 0
}
```

##### `GET /admin/prize-levels`

**描述**: 获取奖项等级列表

**Query 参数**:
- `company_id`: 公司 ID

##### `PUT /admin/prize-levels/:id`

**描述**: 更新奖项等级

**路径参数**:
- `id`: 奖项等级 ID

##### `DELETE /admin/prize-levels/:id`

**描述**: 删除奖项等级

**路径参数**:
- `id`: 奖项等级 ID

#### 奖品管理

##### `POST /admin/prizes`

**描述**: 创建奖品

**请求体**:
```json
{
  "name": "string",
  "level_id": 0,
  "image": "string"
}
```

##### `GET /admin/prizes/:levelId`

**描述**: 获取指定等级的奖品列表

**路径参数**:
- `levelId`: 奖项等级 ID

##### `PUT /admin/prizes/:id`

**描述**: 更新奖品

**路径参数**:
- `id`: 奖品 ID

##### `DELETE /admin/prizes/:id`

**描述**: 删除奖品

**路径参数**:
- `id`: 奖品 ID

#### 抽奖记录和统计

##### `GET /admin/draw-records`

**描述**: 获取抽奖记录

**Query 参数**:
- `page`: 页码
- `page_size`: 每页数量
- `company_id`: 公司 ID
- `level_id`: 奖项等级 ID

##### `GET /admin/stats`

**描述**: 获取统计数据

**Query 参数**:
- `company_id`: 公司 ID

#### 操作日志（仅超级管理员）

##### `GET /admin/operation-logs`

**描述**: 获取操作日志

**Query 参数**:
- `page`: 页码
- `page_size`: 每页数量
- `admin_id`: 管理员 ID
- `action`: 操作类型
- `start_date`: 开始日期
- `end_date`: 结束日期

##### `GET /admin/operation-stats`

**描述**: 获取操作统计

**Query 参数**:
- `start_date`: 开始日期
- `end_date`: 结束日期

---

## 认证说明

### 用户认证

**Header**:
```
Authorization: Bearer <token>
```

**Token 获取**: 通过 `/api/register` 或 `/api/unified-login` 获取

### 管理员认证

**Header**:
```
Authorization: Bearer <token>
```

**Token 获取**: 通过 `/admin/login` 获取

---

## 错误响应

所有错误响应格式：

```json
{
  "error": "错误信息描述"
}
```

**常见 HTTP 状态码**:
- `200 OK`: 请求成功
- `201 Created`: 创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 禁止访问
- `404 Not Found`: 资源不存在
- `429 Too Many Requests`: 请求过于频繁
- `500 Internal Server Error`: 服务器内部错误

---

## 中间件

### CORS 中间件
- 自动处理跨域请求
- 支持的源在配置文件中设置

### 限流中间件
- **内存限流**: 默认每秒 10 个请求
- **Redis 限流**: 支持分布式限流
- **配置**: 可在配置文件中调整

### 认证中间件
- **用户认证**: `UserAuthMiddleware()`
- **管理员认证**: `AuthMiddleware()`

---

*本文档由路由提取工具自动生成*
