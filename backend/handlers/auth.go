package handlers

import (
	"fmt"
	"net/http"

	"lottery-system/config"
	"lottery-system/models"
	"lottery-system/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.Admin `json:"user"`
}

type UnifiedLoginResponse struct {
	Token      string      `json:"token"`
	User       interface{} `json:"user"`       // 可以是 User 或 Admin
	UserType   string      `json:"user_type"`   // "user", "admin", "super_admin"
	IsAdmin    bool        `json:"is_admin"`
	IsSuperAdmin bool        `json:"is_super_admin"`
}

// AdminLogin 管理员登录（保留用于管理后台）
func AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	if err := config.DB.Where("username = ?", req.Username).Preload("Company").First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !utils.CheckPassword(req.Password, admin.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Username, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 记录登录日志
	adminID := uint(admin.ID)
	var companyIDPtr *uint
	if admin.CompanyID != nil {
		cid := uint(*admin.CompanyID)
		companyIDPtr = &cid
	}

	log := models.OperationLog{
		AdminID:    adminID,
		AdminName:  admin.Username,
		CompanyID:  companyIDPtr,
		Action:     "login",
		Resource:   "admin",
		ResourceID: &adminID,
		Details:    "管理员登录",
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	if err := config.DB.Create(&log).Error; err != nil {
		fmt.Printf("[AdminLogin Error] Failed to create login log: %v\n", err)
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  admin,
	})
}

// UnifiedLogin 统一登录接口（支持用户和管理员）
func UnifiedLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取公司代码（必须提供）
	companyCode := c.Query("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code parameter is required"})
		return
	}

	// 先尝试作为用户登录
	var user models.User
	userErr := config.DB.Where("username = ? AND company_id = (SELECT id FROM companies WHERE code = ?)",
		req.Username, companyCode).
		Preload("Company").
		First(&user).Error

	if userErr == nil {
		// 验证密码
		if !utils.CheckPassword(req.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 检查用户角色
		if user.Role != models.RoleUser {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户角色不匹配"})
			return
		}

		// 用户登录成功
		token, err := utils.GenerateUserToken(user.ID, user.Username, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// 记录登录日志
		adminID := uint(user.ID)
		var companyIDPtr *uint
		if user.CompanyID > 0 {
			companyID := uint(user.CompanyID)
			companyIDPtr = &companyID
		}
		log := models.OperationLog{
			AdminID:   adminID,
			AdminName:  user.Username,
			CompanyID: companyIDPtr,
			Action:    "login",
			Resource:  "user",
			ResourceID: &adminID,
			Details:   "用户登录",
			IPAddress:  c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
		}

		if err := config.DB.Create(&log).Error; err != nil {
			fmt.Printf("[UserLogin Error] Failed to create login log: %v\n", err)
		}

		c.JSON(http.StatusOK, UnifiedLoginResponse{
			Token:      token,
			User:       user,
			UserType:   models.RoleUser,
			IsAdmin:    false,
			IsSuperAdmin: false,
		})
		return
	}

	// 用户登录失败，尝试作为管理员登录
	var admin models.Admin
	adminErr := config.DB.Where("username = ?", req.Username).Preload("Company").First(&admin).Error

	if adminErr == nil {
		// 验证密码
		if !utils.CheckPassword(req.Password, admin.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 管理员登录成功
		token, err := utils.GenerateToken(admin.ID, admin.Username, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// 记录登录日志
		adminID := uint(admin.ID)
		var companyIDPtr *uint
		if admin.CompanyID != nil {
			cid := uint(*admin.CompanyID)
			companyIDPtr = &cid
		}

		log := models.OperationLog{
			AdminID:    adminID,
			AdminName:  admin.Username,
			CompanyID:  companyIDPtr,
			Action:     "login",
			Resource:   "admin",
			ResourceID:  &adminID,
			Details:    "管理员登录",
			IPAddress:  c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
		}

		if err := config.DB.Create(&log).Error; err != nil {
			fmt.Printf("[AdminLogin Error] Failed to create login log: %v\n", err)
		}

		// 判断管理员类型
		userType := "admin"
		if admin.IsSuperAdmin {
			userType = "super_admin"
		}

		c.JSON(http.StatusOK, UnifiedLoginResponse{
			Token:      token,
			User:       admin,
			UserType:   userType,
			IsAdmin:    true,
			IsSuperAdmin: admin.IsSuperAdmin,
		})
		return
	}

	// 两种登录都失败
	c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
}
