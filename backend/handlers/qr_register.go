package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"lottery-system/config"
	"lottery-system/models"
	"lottery-system/utils"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// GetRegisterQRCode 生成用户注册二维码
// 返回一个包含注册URL的二维码图片
func GetRegisterQRCode(c *gin.Context) {
	companyCode := c.Query("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code parameter is required"})
		return
	}

	company, err := getCompanyByCode(companyCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company code"})
		return
	}

	// 生成注册URL
	// 格式: https://yourdomain.com/register?company_code=XXX
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// 获取Host头，如果没有则使用配置
	host := c.Request.Host
	if host == "" {
		host = "localhost:5173" // 默认前端开发地址
	}

	registerURL := fmt.Sprintf("%s://%s/register?company_code=%s", scheme, host, company.Code)

	// 生成二维码
	qrCode, err := qrcode.Encode(registerURL, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// 转换为base64
	qrBase64 := base64.StdEncoding.EncodeToString(qrCode)

	c.JSON(http.StatusOK, gin.H{
		"company_code": company.Code,
		"company_name": company.Name,
		"register_url": registerURL,
		"qr_code":      fmt.Sprintf("data:image/png;base64,%s", qrBase64),
	})
}

// UserSelfRegister 用户自助注册（通过扫码）
func UserSelfRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 验证姓名（必填）
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "姓名不能为空"})
		return
	}

	if err := utils.ValidateName(req.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "姓名格式错误: " + err.Error()})
		return
	}

	// 验证手机号（如果提供）
	if req.Phone != "" {
		if err := utils.ValidatePhone(req.Phone); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// 获取公司代码
	companyCode := c.Query("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code parameter is required"})
		return
	}

	company, err := getCompanyByCode(companyCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company code"})
		return
	}

	// 检查用户是否已存在（根据姓名和手机号）
	var existingUser models.User
	query := config.DB.Where("company_id = ? AND name = ?", company.ID, req.Name)
	if req.Phone != "" {
		query = query.Where("phone = ?", req.Phone)
	}

	if err := query.First(&existingUser).Error; err == nil {
		// 用户已存在
		if existingUser.HasDrawn {
			c.JSON(http.StatusConflict, gin.H{
				"error": "该用户已经抽过奖",
				"user":  existingUser,
			})
			return
		} else {
			c.JSON(http.StatusConflict, gin.H{
				"error": "用户已存在于抽奖池中",
				"user":  existingUser,
			})
			return
		}
	}

	// 创建新用户（无用户名和密码，无法登录）
	user := models.User{
		Username:  "", // 空用户名，无法登录
		Password:  "", // 空密码，无法登录
		Name:      req.Name,
		Phone:     req.Phone,
		CompanyID: company.ID,
		HasDrawn:  false,
		Role:      models.RoleUser,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 记录操作日志
	userID := uint(user.ID)
	details := fmt.Sprintf("用户扫码参与: %s", user.Name)
	LogOperation(c, "self_register", "user", &userID, details)

	c.JSON(http.StatusCreated, gin.H{
		"message": "参与成功！您已加入抽奖池",
		"user":    user,
	})
}

// GetCompanyInfo 获取公司信息（用于注册页面）
func GetCompanyInfo(c *gin.Context) {
	companyCode := c.Query("company_code")
	if companyCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_code parameter is required"})
		return
	}

	company, err := getCompanyByCode(companyCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company code"})
		return
	}

	// 统计信息
	var totalUsers int64
	config.DB.Model(&models.User{}).Where("company_id = ?", company.ID).Count(&totalUsers)

	var undrawnUsers int64
	config.DB.Model(&models.User{}).Where("company_id = ? AND has_drawn = ?", company.ID, false).Count(&undrawnUsers)

	c.JSON(http.StatusOK, gin.H{
		"company": company,
		"stats": gin.H{
			"total_users":   totalUsers,
			"undrawn_users": undrawnUsers,
		},
	})
}
