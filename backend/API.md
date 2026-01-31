# API Documentation

## 基础信息

**Base URL**: `http://localhost:8080`

**Content-Type**: `application/json`

**字符编码**: `UTF-8`

---

## 标准响应格式

### 成功响应

```json
{
  "success": true,
  "data": {
    // 响应数据
  },
  "error": null
}
```

### 错误响应

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "错误描述"
  }
}
```

---

## 认证相关API

### 1. 用户登录/注册

**端点**: `POST /api/register` 或 `POST /api/login`

**请求体**:
```json
{
  "username": "zhangsan",
  "password": "123456"
}
```

**Query参数**:
- `company_code` (可选): 公司代码，默认为 "default"

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "zhangsan",
      "name": "张三",
      "company_id": 1,
      "has_drawn": false,
      "role": "user"
    },
    "user_type": "user"
  }
}
```

**错误响应**:
- 401: 用户名或密码错误
- 404: 用户不存在
- 400: 参数格式错误

---

### 2. 管理员登录

**端点**: `POST /admin/login`

**请求体**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "is_super_admin": true,
      "company_id": null,
      "role": "super_admin"
    }
  }
}
```

**错误响应**:
- 401: 用户名或密码错误
- 404: 管理员不存在

---

## 用户管理API (需要认证)

### 3. 创建用户

**端点**: `POST /admin/users`

**认证**: 需要管理员Token

**请求体**:
```json
{
  "company_id": 1,
  "username": "lisi",
  "password": "123456",
  "name": "李四",
  "phone": "13800138002"
}
```

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "id": 2,
    "username": "lisi",
    "name": "李四",
    "phone": "13800138002",
    "company_id": 1,
    "has_drawn": false,
    "role": "user"
  }
}
```

**错误响应**:
- 400: 用户名已存在
- 400: 手机号已存在
- 403: 权限不足
- 404: 公司不存在

---

### 4. 批量创建用户

**端点**: `POST /admin/users/batch`

**认证**: 需要管理员Token

**请求体**:
```json
{
  "company_id": 1,
  "users": [
    "zhangsan,123456,张三,13800138001",
    "lisi,123456,李四,13800138002",
    "wangwu,123456,王五"
  ]
}
```

**格式说明**: `用户名,密码,姓名,手机号(可选)`

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "created": 2,
    "failed": 1,
    "users": [
      {
        "id": 3,
        "username": "zhangsan",
        "name": "张三"
      },
      {
        "id": 4,
        "username": "lisi",
        "name": "李四"
      }
    ],
    "errors": [
      "wangwu (手机号格式错误)"
    ]
  }
}
```

---

### 5. 获取用户列表

**端点**: `GET /admin/users`

**认证**: 需要管理员Token

**Query参数**:
- `company_id` (可选): 公司ID，超级管理员可筛选
- `has_drawn` (可选): 抽奖状态，true/false

**成功响应** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "username": "zhangsan",
      "name": "张三",
      "phone": "13800138001",
      "company_id": 1,
      "has_drawn": false,
      "company": {
        "id": 1,
        "name": "默认公司",
        "code": "default"
      }
    }
  ]
}
```

---

### 6. 更新用户

**端点**: `PUT /admin/users/:id`

**认证**: 需要管理员Token

**请求体**:
```json
{
  "name": "张三丰",
  "phone": "13900139001",
  "has_drawn": false
}
```

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "zhangsan",
    "name": "张三丰",
    "phone": "13900139001",
    "has_drawn": false
  }
}
```

---

### 7. 删除用户

**端点**: `DELETE /admin/users/:id`

**认证**: 需要管理员Token

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "message": "User deleted successfully"
  }
}
```

---

## 抽奖API (需要认证)

### 8. 获取奖项等级

**端点**: `GET /api/prize-levels`

**认证**: 需要用户Token

**Query参数**:
- `company_code` (可选): 公司代码

**成功响应** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "一等奖",
      "description": "iPhone 15 Pro",
      "probability": 0.05,
      "total_stock": 5,
      "used_stock": 2,
      "sort_order": 1,
      "is_active": true
    }
  ]
}
```

---

### 9. 执行抽奖

**端点**: `POST /api/draw`

**认证**: 需要用户Token

**请求体**:
```json
{
  "level_id": 0,
  "count": 1
}
```

**参数说明**:
- `level_id`: 奖项等级ID，0表示随机抽奖
- `count`: 抽奖数量

**成功响应** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 10,
      "level_id": 2,
      "prize_id": 5,
      "ip": "127.0.0.1",
      "created_at": "2026-01-24T12:00:00Z",
      "user": {
        "id": 10,
        "username": "zhangsan",
        "name": "张三"
      },
      "level": {
        "id": 2,
        "name": "二等奖"
      },
      "prize": {
        "id": 5,
        "name": "AirPods Pro"
      }
    }
  ]
}
```

**错误响应**:
- 400: 用户已经抽过奖
- 400: 奖品已抽完
- 400: 没有可抽奖用户

---

### 10. 获取我的奖品

**端点**: `GET /api/my-prize`

**认证**: 需要用户Token

**Query参数**:
- `company_code` (可选): 公司代码

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 10,
    "level": {
      "name": "二等奖"
    },
    "prize": {
      "name": "AirPods Pro"
    }
  }
}
```

**错误响应**:
- 400: 您还没有参与抽奖

---

## 管理员管理API (需要超级管理员认证)

### 11. 创建管理员

**端点**: `POST /admin/admins`

**认证**: 需要超级管理员Token

**请求体**:
```json
{
  "username": "admin2",
  "password": "admin123",
  "is_super_admin": false,
  "company_id": 1
}
```

**成功响应** (201):
```json
{
  "success": true,
  "data": {
    "id": 2,
    "username": "admin2",
    "is_super_admin": false,
    "company_id": 1,
    "role": "admin",
    "company": {
      "id": 1,
      "name": "默认公司"
    }
  }
}
```

**错误响应**:
- 400: 管理员已存在
- 400: 普通管理员必须指定公司
- 403: 权限不足（非超级管理员）

---

### 12. 更新管理员

**端点**: `PUT /admin/admins/:id`

**认证**: 需要超级管理员Token或本人

**请求体**:
```json
{
  "password": "newpass123"
}
```

**成功响应** (200): 返回更新后的管理员信息

---

### 13. 删除管理员

**端点**: `DELETE /admin/admins/:id`

**认证**: 需要超级管理员Token

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "message": "Admin deleted successfully"
  }
}
```

**错误响应**:
- 400: 不能删除自己
- 403: 权限不足

---

## 公司管理API (需要超级管理员认证)

### 14. 创建公司

**端点**: `POST /admin/companies`

**认证**: 需要超级管理员Token

**请求体**:
```json
{
  "name": "技术部",
  "code": "tech",
  "description": "技术部门",
  "is_active": true
}
```

**成功响应** (201):
```json
{
  "success": true,
  "data": {
    "id": 2,
    "name": "技术部",
    "code": "tech",
    "is_active": true
  }
}
```

---

### 15. 获取公司列表

**端点**: `GET /admin/companies`

**认证**: 需要超级管理员Token

**成功响应** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "默认公司",
      "code": "default",
      "is_active": true
    }
  ]
}
```

---

## 统计API

### 16. 获取统计数据

**端点**: `GET /admin/stats`

**认证**: 需要管理员Token

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "total_users": 100,
    "drawn_users": 50,
    "available_users": 50,
    "total_companies": 3,
    "active_companies": 3
  }
}
```

---

### 17. 获取操作日志

**端点**: `GET /admin/operation-logs`

**认证**: 需要超级管理员Token

**Query参数**:
- `page` (可选): 页码
- `page_size` (可选): 每页数量
- `action` (可选): 操作类型
- `resource` (可选): 资源类型

**成功响应** (200):
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 1,
        "admin_id": 1,
        "admin_name": "admin",
        "action": "create",
        "resource": "user",
        "resource_id": 10,
        "details": "创建用户: 张三 (zhangsan)",
        "ip_address": "127.0.0.1",
        "created_at": "2026-01-24T12:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 错误代码参考

| 错误代码 | HTTP状态 | 说明 |
|---------|---------|------|
| `BAD_REQUEST` | 400 | 请求参数错误 |
| `UNAUTHORIZED` | 401 | 未认证或认证失败 |
| `FORBIDDEN` | 403 | 权限不足 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `CONFLICT` | 409 | 资源冲突（如用户名已存在） |
| `INTERNAL_ERROR` | 500 | 服务器内部错误 |
| `INVALID_CREDENTIALS` | 401 | 用户名或密码错误 |
| `INVALID_TOKEN` | 401 | Token无效或已过期 |
| `USER_NOT_FOUND` | 404 | 用户不存在 |
| `PERMISSION_DENIED` | 403 | 权限不足 |
| `PRIZE_OUT_OF_STOCK` | 400 | 奖品已抽完 |
| `INVALID_INPUT` | 400 | 输入数据无效 |

---

## 限流规则

| 操作类型 | 限制 |
|---------|------|
| 登录 | 每分钟5次，每小时20次 |
| 密码修改 | 每分钟2次，每小时5次 |
| 管理员创建 | 每分钟1次，每小时10次 |

超过限制返回 `429 Too Many Requests`

---

## 请求头

### 认证请求

```
Authorization: Bearer <token>
```

### 通用请求头

```
Content-Type: application/json
X-Request-ID: <request-id>
```

---

## 分页

大多数列表API支持分页：

**Query参数**:
- `page`: 页码（从1开始）
- `page_size`: 每页数量（默认20，最大100）

**响应格式**:
```json
{
  "success": true,
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

---

**最后更新**: 2026-01-24
**版本**: 2.0
