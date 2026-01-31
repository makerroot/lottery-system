package middleware

import (
	"net/http"
	"strings"

	"lottery-system/config"
	"lottery-system/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateToken(tokenString, config.AppConfig.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 获取管理员信息
		var admin struct {
			ID           int
			CompanyID    *int
			IsSuperAdmin bool
		}
		if err := config.DB.Table("admins").
			Select("id, company_id, is_super_admin").
			Where("id = ?", claims.UserID).
			First(&admin).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin not found"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("company_id", admin.CompanyID)
		c.Set("is_super_admin", admin.IsSuperAdmin)
		c.Next()
	}
}

// UserAuthMiddleware 用户JWT认证中间件（同时支持用户和管理员token）
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateToken(tokenString, config.AppConfig.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 先尝试作为用户验证
		var user struct {
			ID        int
			CompanyID int
			Phone     string
			Name      string
			HasDrawn  bool
		}
		userErr := config.DB.Table("users").
			Select("id, company_id, phone, name, has_drawn").
			Where("id = ?", claims.UserID).
			First(&user).Error

		if userErr == nil {
			// 是用户
			c.Set("user_id", user.ID)
			c.Set("user_phone", user.Phone)
			c.Set("user_name", user.Name)
			c.Set("company_id", user.CompanyID)
			c.Set("has_drawn", user.HasDrawn)
			c.Set("is_admin", false)
			c.Next()
			return
		}

		// 如果不是用户，尝试作为管理员验证
		var admin struct {
			ID           int
			CompanyID    *int
			IsSuperAdmin bool
		}
		adminErr := config.DB.Table("admins").
			Select("id, company_id, is_super_admin").
			Where("id = ?", claims.UserID).
			First(&admin).Error

		if adminErr == nil {
			// 是管理员
			c.Set("user_id", admin.ID)
			c.Set("company_id", admin.CompanyID)
			c.Set("is_admin", true)
			c.Set("is_super_admin", admin.IsSuperAdmin)
			c.Next()
			return
		}

		// 既不是用户也不是管理员
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查来源是否在允许列表中
		allowedOrigins := strings.Split(config.AppConfig.AllowedOrigins, ",")
		isAllowed := false

		// 清理空格并检查
		for _, allowed := range allowedOrigins {
			if strings.TrimSpace(allowed) == "*" || strings.TrimSpace(allowed) == origin {
				isAllowed = true
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
