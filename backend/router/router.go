package router

import (
	"lottery-system/handlers"
	"lottery-system/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回 Gin 路由引擎
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 应用全局中间件
	r.Use(middleware.CORSMiddleware())

	// 根据配置选择限流中间件（需要在 main.go 中设置 Redis）
	// 这里我们使用简单的判断，实际配置在 main.go 中处理
	if middleware.IsRedisEnabled() {
		r.Use(middleware.RedisRateLimitMiddleware())
	} else {
		r.Use(middleware.RateLimitMiddleware())
	}

	// 注册健康检查端点（无需认证）
	setupHealthCheckRoutes(r)

	// 注册用户端 API 路由
	setupUserAPIRoutes(r)

	// 注册管理后台 API 路由
	setupAdminAPIRoutes(r)

	return r
}

// setupHealthCheckRoutes 配置健康检查端点
func setupHealthCheckRoutes(r *gin.Engine) {
	// 健康检查端点（无需认证，支持 GET 和 HEAD）
	r.GET("/api/health", handlers.HealthCheck)
	r.HEAD("/api/health", handlers.HealthCheck)
}

// setupUserAPIRoutes 配置用户端 API 路由
func setupUserAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 公开接口（无需认证）：唯一登录入口
		api.POST("/login", handlers.RegisterOrLogin) // 统一登录接口（仅管理员）

		// 用户自助注册相关接口（公开）
		api.GET("/qr-register", handlers.GetRegisterQRCode)       // 获取注册二维码
		api.GET("/company-info", handlers.GetCompanyInfo)         // 获取公司信息
		api.POST("/self-register", handlers.UserSelfRegister)     // 用户自助注册

		// 需要用户认证的接口
		userAuth := api.Group("")
		userAuth.Use(middleware.UserAuthMiddleware())
		{
			// 公司信息
			userAuth.GET("/company", handlers.GetCompanyByCode)

			// 用户信息
			userAuth.GET("/user", handlers.GetUserInfo)
			userAuth.POST("/user/change-password", handlers.ChangeUserPassword)

			// 奖品相关
			userAuth.GET("/prize-levels", handlers.GetActivePrizeLevels)

			// 抽奖相关
			userAuth.POST("/draw", handlers.Draw)
			userAuth.GET("/my-prize", handlers.GetMyPrize)
			userAuth.GET("/user-stats", handlers.GetUserStats)
			userAuth.GET("/draw-records", handlers.GetDrawRecordsPublic)
			userAuth.GET("/available-users", handlers.GetAvailableUsersPublic)
		}
	}
}

// setupAdminAPIRoutes 配置管理后台 API 路由
func setupAdminAPIRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		// 登录（无需认证）
		admin.POST("/login", handlers.AdminLogin)

		// 需要认证的路由
		auth := admin.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 管理员信息
			auth.GET("/info", handlers.GetAdminInfo)
			auth.POST("/change-password", handlers.ChangeAdminPassword)

			// 管理员管理（仅超级管理员）
			auth.GET("/admins", handlers.GetAdmins)
			auth.POST("/admins", handlers.CreateAdmin)
			auth.PUT("/admins/:id", handlers.UpdateAdmin)
			auth.DELETE("/admins/:id", handlers.DeleteAdmin)

			// 用户管理
			auth.GET("/users", handlers.GetUsers)
			auth.POST("/users", handlers.CreateUser)
			auth.POST("/users/batch", handlers.BatchCreateUsers)
			auth.POST("/users/scan-add", handlers.ScanAddUser) // 扫码添加用户
			auth.PUT("/users/:id", handlers.UpdateUser)
			auth.DELETE("/users/:id", handlers.DeleteUser)

			// 公司管理（超级管理员）
			auth.GET("/companies", handlers.GetCompanies)
			auth.POST("/companies", handlers.CreateCompany)
			auth.PUT("/companies/:id", handlers.UpdateCompany)
			auth.DELETE("/companies/:id", handlers.DeleteCompany)
			auth.GET("/company-stats", handlers.GetCompanyStats)

			// 奖项等级管理
			auth.POST("/prize-levels", handlers.CreatePrizeLevel)
			auth.GET("/prize-levels", handlers.GetPrizeLevels)
			auth.PUT("/prize-levels/:id", handlers.UpdatePrizeLevel)
			auth.DELETE("/prize-levels/:id", handlers.DeletePrizeLevel)

			// 奖品管理
			auth.GET("/prizes/all", handlers.GetAllPrizes)
			auth.POST("/prizes", handlers.CreatePrize)
			auth.GET("/prizes/:levelId", handlers.GetPrizesByLevel)
			auth.PUT("/prizes/:id", handlers.UpdatePrize)
			auth.DELETE("/prizes/:id", handlers.DeletePrize)

			// 抽奖记录和统计
			auth.GET("/draw-records", handlers.GetDrawRecords)
			auth.GET("/stats", handlers.GetStats)

			// 操作日志（仅超级管理员）
			auth.GET("/operation-logs", handlers.GetOperationLogs)
			auth.GET("/operation-stats", handlers.GetOperationStats)
		}
	}
}
