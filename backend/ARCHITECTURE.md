# Backend Architecture Documentation

## 📐 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                        │
│  (Frontend: Vue 3 + Vite)                                     │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP/REST API
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                      Middleware Layer                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐  │
│  │   CORS   │→ │   Auth   │→ │Rate Limit│→ │RequestID│  │
│  └──────────┘  └──────────┘  └──────────┘  └────────┘  │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                      Handler Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Lottery      │  │ User         │  │ Admin        │      │
│  │ Handlers     │  │ Handlers     │  │ Handlers     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐  │
│  │   Auth   │  │   User   │  │  Admin   │  │ Lottery │  │
│  │ Service  │  │ Service  │  │ Service  │  │ Service │  │
│  └──────────┘  └──────────┘  └──────────┘  └────────┘  │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    Repository Layer                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐  │
│  │   User   │  │   Admin   │  │ Company  │  │  Prize  │  │
│  │   Repo   │  │   Repo    │  │   Repo   │  │   Repo  │  │
│  └──────────┘  └──────────┘  └──────────┘  └────────┘  │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                      Database Layer                          │
│              MySQL Database (GORM)                          │
└─────────────────────────────────────────────────────────────┘
```

## 📦 层级职责

### 1. Handler Layer (处理器层)

**职责**：
- 处理HTTP请求和响应
- 调用Service层方法
- 执行输入验证（基础）
- 返回标准化响应
- 设置HTTP状态码

**特点**：
- 不包含业务逻辑
- 不直接访问数据库
- 使用 `response` 包返回统一格式

**示例**：
```go
func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "invalid request")
        return
    }

    userService := services.NewUserService()
    user, err := userService.CreateUser(&req)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    response.Success(c, user)
}
```

### 2. Service Layer (服务层)

**职责**：
- 实现核心业务逻辑
- 协调多个Repository
- 执行复杂验证
- 处理业务规则
- 管理事务

**特点**：
- 可重用的业务逻辑
- 独立于HTTP层
- 可独立测试
- 使用Repository访问数据

**示例**：
```go
func (s *UserService) CreateUser(req *CreateUserRequest) (*models.User, error) {
    // 验证
    if err := validators.ValidateUsername(req.Username); err != nil {
        return nil, err
    }

    // 业务规则检查
    exists, _ := s.userRepo.ExistsByUsername(req.Username, req.CompanyID)
    if exists {
        return nil, utils.NewBusinessLogicError("用户名已存在")
    }

    // 密码加密
    hashedPassword, err := s.authService.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    // 创建用户
    user := &models.User{
        Username: req.Username,
        Password: hashedPassword,
        // ...
    }

    return user, s.userRepo.Create(user)
}
```

### 3. Repository Layer (数据访问层)

**职责**：
- 数据库CRUD操作
- 封装数据库查询
- 提供查询接口
- 处理数据映射

**特点**：
- 不包含业务逻辑
- 只负责数据访问
- 使用GORM ORM
- 可模拟（mock）

**示例**：
```go
func (r *UserRepository) FindByID(id int) (*models.User, error) {
    var user models.User
    err := config.DB.Preload("Company").First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

## 🔄 请求流程

### 典型请求流程：用户登录

```
1. Client Request
   POST /api/login
   {
     "username": "zhangsan",
     "password": "123456"
   }

2. Middleware Layer
   ├─ CORS Middleware          # 检查跨域
   ├─ Rate Limit Middleware     # 检查限流
   ├─ Request ID Middleware    # 添加请求ID
   └─ Auth Middleware (skip)    # 跳过认证

3. Handler Layer
   ├─ Parse request
   ├─ Basic validation
   └─ Call service

4. Service Layer
   ├─ Validate credentials
   ├─ Hash password check
   ├─ Generate JWT token
   └─ Log successful login

5. Repository Layer
   ├─ Query user by username
   └─ Return user data

6. Response
   {
     "success": true,
     "data": {
       "token": "...",
       "user": {...}
     },
     "error": null
   }
```

## 🎯 设计模式

### 1. Repository Pattern

**目的**: 抽象数据访问

**实现**:
```go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
    // ...
}
```

**好处**:
- 数据访问逻辑集中
- 易于测试（可mock）
- 数据库无关

### 2. Service Layer Pattern

**目的**: 封装业务逻辑

**实现**:
```go
type UserService struct {
    userRepo  *repositories.UserRepository
    authService *services.AuthService
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*models.User, error) {
    // 业务逻辑
}
```

**好处**:
- 业务逻辑可重用
- 独立于HTTP层
- 易于测试

### 3. Dependency Injection

**目的**: 降低耦合

**实现**:
```go
func NewUserService() *UserService {
    return &UserService{
        userRepo:  repositories.NewUserRepository(),
        authService: NewAuthService(),
    }
}
```

**好处**:
- 依赖关系清晰
- 易于替换实现
- 便于单元测试

### 4. Strategy Pattern

**目的**: 算法族

**示例**: 密码验证
```go
// 验证用户密码
ValidatePasswordForUser(password string)

// 验证管理员密码
ValidatePasswordForAdmin(password string)
```

**好处**:
- 算法可替换
- 符合开闭原则

### 5. Error Wrapper Pattern

**目的**: 错误上下文

**实现**:
```go
if err := s.userRepo.Create(user); err != nil {
    return nil, fmt.Errorf("创建用户失败: %w", err)
}
```

**好处**:
- 保留错误链
- 添加上下文信息
- 便于调试

## 🔒 安全架构

### 认证流程

```
1. 用户提交用户名/密码
        ↓
2. Handler调用Service
        ↓
3. AuthService验证
        ↓
4. UserRepository查询
        ↓
5. bcrypt验证密码
        ↓
6. 生成JWT Token
        ↓
7. 返回Token
```

### 权限检查

```
Request → Middleware
    ↓
从Token解析用户信息
    ↓
存入Context (user_id, role, company_id)
    ↓
Handler从Context获取用户信息
    ↓
Service层检查权限
    ↓
Repository层执行操作
```

### 输入验证流程

```
用户输入
    ↓
Handler基础验证
    ↓
Validator详细验证
    ↓
Service业务规则验证
    ↓
Repository约束验证
    ↓
数据库
```

## 📊 数据流

### 创建用户数据流

```
HTTP Request
    ↓
Handler::CreateUser
    ↓
UserService::CreateUser
    ├→ ValidateUsername
    ├→ ValidatePassword
    ├→ CheckExistsByUsername
    ├→ HashPassword
    ├→ Create user object
    └→ UserRepository::Create
        ↓
    Database Transaction
        ↓
    HTTP Response
```

### 抽奖数据流

```
HTTP Request (user_id, level_id)
    ↓
Handler::Draw
    ↓
LotteryService::DrawPrize
    ├→ CheckUserCanDraw
    ├→ GetPrizeLevel (with stock check)
    ├→ Begin Transaction
    ├→ CreateDrawRecord
    ├→ UpdatePrizeStock
    ├→ UpdateUserStatus
    ├→ Commit Transaction
    └→ LoadAssociations
        ↓
    HTTP Response (draw record with prizes)
```

## 🛡️ 安全机制

### 1. 密码安全

```
用户输入明文密码
    ↓
验证长度和复杂度
    ↓
bcrypt哈希 (cost=14)
    ↓
存储哈希值
    ↓
验证时: bcrypt.CompareHashAndPassword()
```

### 2. JWT认证

```
登录成功
    ↓
生成JWT (包含user_id, username, exp)
    ↓
签名: HMACSHA256(JWT_SECRET)
    ↓
返回Token
    ↓
客户端请求携带Token
    ↓
中间件验证Token
    ↓─有效→ 解析用户信息
  └─无效→ 返回401
```

### 3. 限流保护

```
Request
    ↓
Rate Limit Middleware
    ↓
检查IP/用户请求频率
    ├─ 正常 → 放行
    └─ 超限 → 返回429
```

### 4. XSS防护

```
用户输入
    ↓
SanitizeInput (清理HTML标签)
    ↓
ValidateXSS (检测XSS模式)
    ├─ 安全 → 继续处理
    └─ 危险 → 记录日志并拒绝
```

## 🧪 测试策略

### 单元测试

```
tests/
├── services/
│   ├── auth_service_test.go
│   ├── user_service_test.go
│   └── lottery_service_test.go
├── repositories/
│   ├── user_repository_test.go
│   └── admin_repository_test.go
└── validators/
    ├── validator_test.go
    └── password_test.go
```

### 集成测试

```
tests/
└── integration/
    ├── auth_flow_test.go
    ├── lottery_flow_test.go
    └── user_management_test.go
```

### 测试覆盖目标

- 整体覆盖率：80%+
- Service层：90%+
- Repository层：80%+
- Handler层：70%+

## 📈 性能优化

### 1. 数据库优化

- ✅ 添加必要的索引
- ✅ 使用Preload减少N+1查询
- ✅ 批量操作使用事务
- ✅ 连接池配置

### 2. 缓存策略

```
可缓存的资源：
- 用户信息（短期）
- 公司信息（中期）
- 奖项等级（短期）
- 统计数据（短期）

Redis作为可选缓存层
```

### 3. 限流策略

```
多层限流：
1. 全局限流（所有请求）
2. 端点限流（特定API）
3. 敏感操作限流（登录、密码修改）
```

## 🔧 运维友好

### 健康检查

```go
GET /health
Response: {"status": "ok"}
```

### 优雅关闭

```go
监听系统信号:
- SIGTERM: 优雅关闭
- SIGINT: 立即关闭
```

### 日志级别

```
Development: Debug
Staging: Info
Production: Warn/Error
```

---

## 📚 相关文档

- [README.md](README.md) - 项目概述
- [BACKEND_OPTIMIZATION_PLAN.md](BACKEND_OPTIMIZATION_PLAN.md) - 优化计划
- [BACKEND_OPTIMIZATION_FINAL_REPORT.md](BACKEND_OPTIMIZATION_FINAL_REPORT.md) - 优化报告

---

**最后更新**: 2026-01-24
**版本**: 2.0
