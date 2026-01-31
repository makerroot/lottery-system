package handlers

import (
	"fmt"
	"net/http"

	"lottery-system/config"
	"lottery-system/models"
	"lottery-system/utils"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username"` // 可选，保留以兼容旧接口
	Password string `json:"password"` // 可选，保留以兼容旧接口
	Name     string `json:"name"`     // 必填：姓名
	Phone    string `json:"phone"`    // 可选：手机号
}

type DrawRequest struct {
	LevelID   int    `json:"level_id"`   // 指定抽取的奖项等级ID，0表示不指定
	Count     int    `json:"count"`      // 抽取人数
	UserPhone string `json:"user_phone"` // 指定中奖用户的手机号（用于前端选择中奖者）
}

// getCompanyByCode 根据代码获取公司（必须提供参数）
func getCompanyByCode(code string) (*models.Company, error) {
	if code == "" {
		return nil, fmt.Errorf("company_code parameter is required")
	}

	var company models.Company
	if err := config.DB.Where("code = ? AND is_active = ?", code, true).First(&company).Error; err != nil {
		return nil, fmt.Errorf("company not found")
	}

	return &company, nil
}

// RegisterOrLogin 用户或管理员登录（通过用户名密码）
// 支持普通用户和管理员在抽奖页面登录
func RegisterOrLogin(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 验证用户名
	if err := utils.ValidateName(req.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名格式错误: " + err.Error()})
		return
	}

	// 验证密码长度
	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度至少6位"})
		return
	}

	// 获取公司代码（必须提供）
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

	// 先尝试查找普通用户
	var user models.User
	userErr := config.DB.Where("username = ? AND company_id = ?", req.Username, company.ID).
		Preload("Company").
		First(&user).Error

	if userErr == nil {
		// 找到用户，验证用户密码
		if !utils.CheckPassword(req.Password, user.Password) {
			utils.WithFields(map[string]interface{}{
				"username":   req.Username,
				"company_id": company.ID,
			}).Warn("用户登录失败：密码错误")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "用户名或密码错误",
				"error_code": "INVALID_PASSWORD",
			})
			return
		}

		utils.WithFields(map[string]interface{}{
			"user_id":    user.ID,
			"username":   user.Username,
			"company_id": company.ID,
			"user_type":  "user",
		}).Info("用户登录成功")

		// 记录操作日志
		userID := uint(user.ID)
		details := fmt.Sprintf("用户登录: %s (%s)", user.Name, user.Username)
		LogOperation(c, "login", "user", &userID, details)

		// 生成用户token
		token, err := utils.GenerateUserToken(user.ID, user.Username, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiration)
		if err != nil {
			utils.WithFields(map[string]interface{}{
				"error": err,
			}).Error("生成用户token失败")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":     token,
			"user":      user,
			"user_type": "user",
		})
		return
	}

	// 用户不存在，尝试查找管理员
	var admin models.Admin
	adminErr := config.DB.Where("username = ?", req.Username).Preload("Company").First(&admin).Error

	if adminErr != nil {
		utils.WithFields(map[string]interface{}{
			"username":   req.Username,
			"company_id": company.ID,
		}).Warn("登录失败：管理员不存在")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "用户名或密码错误",
			"error_code": "INVALID_CREDENTIALS",
		})
		return
	}

	// 找到管理员，验证管理员密码
	if !utils.CheckPassword(req.Password, admin.Password) {
		utils.WithFields(map[string]interface{}{
			"username": req.Username,
		}).Warn("管理员登录失败：密码错误")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "用户名或密码错误",
			"error_code": "INVALID_PASSWORD",
		})
		return
	}

	// 检查管理员权限
	// 超级管理员可以登录任何公司
	// 普通管理员只能登录到所属的公司
	if !admin.IsSuperAdmin && admin.CompanyID != nil {
		if *admin.CompanyID != company.ID {
			utils.WithFields(map[string]interface{}{
				"username":      req.Username,
				"admin_company": *admin.CompanyID,
				"login_company": company.ID,
			}).Warn("管理员登录失败：公司不匹配")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "您不是该公司的管理员",
				"error_code": "COMPANY_MISMATCH",
			})
			return
		}
	}

	utils.WithFields(map[string]interface{}{
		"admin_id":   admin.ID,
		"username":   admin.Username,
		"company_id": company.ID,
		"user_type":  "admin",
		"is_super":   admin.IsSuperAdmin,
	}).Info("管理员登录成功")

	// 记录操作日志
	adminID := uint(admin.ID)
	details := fmt.Sprintf("管理员登录抽奖页面: %s", admin.Username)
	if admin.CompanyID != nil {
		details += fmt.Sprintf(" (公司ID: %d)", *admin.CompanyID)
	}
	LogOperation(c, "login", "admin", &adminID, details)

	// 生成管理员token（使用GenerateToken而不是GenerateUserToken）
	token, err := utils.GenerateToken(admin.ID, admin.Username, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiration)
	if err != nil {
		utils.WithFields(map[string]interface{}{
			"error": err,
		}).Error("生成管理员token失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 构造响应数据（格式与user一致）
	responseData := gin.H{
		"token":          token,
		"user":           admin,
		"user_type":      "admin",
		"is_super_admin": admin.IsSuperAdmin,
	}

	c.JSON(http.StatusOK, responseData)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	phone := c.Query("phone")
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

	var user models.User
	if err := config.DB.Where("phone = ? AND company_id = ?", phone, company.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetActivePrizeLevels 获取启用的奖项等级（用户端，包含奖品库存信息）
func GetActivePrizeLevels(c *gin.Context) {
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

	var levels []models.PrizeLevel
	config.DB.Where("company_id = ? AND is_active = ?", company.ID, true).
		Order("sort_order ASC").
		Find(&levels)

	// 为每个奖项等级计算奖品的库存信息
	type PrizeLevelWithStock struct {
		models.PrizeLevel
		TotalStock int `json:"total_stock"`
		UsedStock  int `json:"used_stock"`
	}

	result := make([]PrizeLevelWithStock, len(levels))
	for i, level := range levels {
		// 查询该奖项下所有奖品的库存总和
		var stockData struct {
			TotalStock int `json:"total_stock"`
			UsedStock  int `json:"used_stock"`
		}
		config.DB.Model(&models.Prize{}).
			Where("level_id = ?", level.ID).
			Select("COALESCE(SUM(total_stock), 0) as total_stock, COALESCE(SUM(used_stock), 0) as used_stock").
			Scan(&stockData)

		result[i] = PrizeLevelWithStock{
			PrizeLevel: level,
			TotalStock: stockData.TotalStock,
			UsedStock:  stockData.UsedStock,
		}
	}

	c.JSON(http.StatusOK, result)
}

// Draw 执行抽奖 - 仅限管理员和超级管理员
func Draw(c *gin.Context) {
	// 🔒 权限检查：只允许管理员和超级管理员抽奖
	isAdmin, _ := c.Get("is_admin")
	isSuperAdmin, _ := c.Get("is_super_admin")

	if isAdmin == false && isSuperAdmin == false {
		c.JSON(http.StatusForbidden, gin.H{
			"error":      "只有管理员才能进行抽奖操作",
			"error_code": "PERMISSION_DENIED",
		})
		return
	}

	// 获取公司代码（必须提供）
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

	// 解析请求参数
	var req DrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有指定level_id，使用用户传递的count
	levelID := req.LevelID
	drawCount := req.Count

	if levelID == 0 {
		// 未指定奖项，从所有未抽奖用户中抽取
		if drawCount <= 0 {
			drawCount = 1
		}

		// 查找该公司所有未抽奖的用户
		var users []models.User
		if err := config.DB.Where("company_id = ? AND has_drawn = ?", company.ID, false).
			Order("id ASC").
			Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		// 检查是否还有未抽奖的用户
		if len(users) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "没有可抽奖的用户"})
			return
		}

		// 如果剩余用户少于抽奖人数，只抽剩余的
		if drawCount > len(users) {
			drawCount = len(users)
		}

		// 随机选择用户
		selectedIndices := utils.RandomIndices(len(users), drawCount)
		var records []models.DrawRecord
		ip := c.ClientIP()

		// 为每个选中的用户执行抽奖
		for _, index := range selectedIndices {
			selectedUser := users[index]
			record, err := utils.DrawLottery(config.DB, &selectedUser, ip, company.ID)
			if err != nil {
				continue // 跳过失败的抽奖
			}
			records = append(records, *record)
		}

		// 如果所有抽奖都失败了
		if len(records) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "抽奖失败"})
			return
		}

		// 加载关联数据
		for i := range records {
			config.DB.Preload("User").Preload("Level").Preload("Prize").First(&records[i], records[i].ID)
		}

		c.JSON(http.StatusOK, records)
		return
	}

	// 指定了奖项等级，只从该奖项中抽取
	// 检查奖项是否存在且有库存
	var level models.PrizeLevel
	if err := config.DB.Where("id = ? AND company_id = ? AND is_active = ?", levelID, company.ID, true).
		First(&level).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "奖项不存在或已禁用"})
		return
	}

	// 计算该奖项等级下所有奖品的实际库存（从 Prize 表聚合）
	type StockInfo struct {
		TotalStock int `json:"total_stock"`
		UsedStock  int `json:"used_stock"`
	}
	var stockInfo StockInfo
	config.DB.Model(&models.Prize{}).
		Where("level_id = ?", levelID).
		Select("COALESCE(SUM(total_stock), 0) as total_stock, COALESCE(SUM(used_stock), 0) as used_stock").
		Scan(&stockInfo)

	// 检查实际库存
	if stockInfo.UsedStock >= stockInfo.TotalStock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该奖项已抽完"})
		return
	}

	// 限制抽取数量
	if drawCount <= 0 {
		drawCount = 1
	}
	available := stockInfo.TotalStock - stockInfo.UsedStock
	if drawCount > available {
		drawCount = available
	}

	// 查找该公司所有未抽奖的用户
	var users []models.User
	if err := config.DB.Where("company_id = ? AND has_drawn = ?", company.ID, false).
		Order("id ASC").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// 检查用户数量
	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有可抽奖的用户"})
		return
	}

	// 如果指定了用户手机号，使用该用户
	var records []models.DrawRecord
	ip := c.ClientIP()

	if req.UserPhone != "" {
		// 查找指定的用户
		var specifiedUser models.User
		if err := config.DB.Where("phone = ? AND company_id = ? AND has_drawn = ?", req.UserPhone, company.ID, false).
			First(&specifiedUser).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "指定的用户不存在或已抽过奖"})
			return
		}

		// 使用指定的用户执行抽奖（作为第1个中奖者）
		record, err := utils.DrawLotteryWithLevel(config.DB, &specifiedUser, ip, company.ID, levelID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "抽奖失败: " + err.Error()})
			return
		}
		records = append(records, *record)

		// 如果需要更多中奖者（count > 1），从其他用户中随机选择
		if drawCount > 1 {
			remainingCount := drawCount - 1

			// 获取其他未抽奖的用户（排除已指定的用户）
			var otherUsers []models.User
			if err := config.DB.Where("company_id = ? AND has_drawn = ? AND id != ?", company.ID, false, specifiedUser.ID).
				Order("id ASC").
				Find(&otherUsers).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
				return
			}

			// 检查是否有足够的其他用户
			if len(otherUsers) == 0 {
				// 没有其他用户，只返回指定的1个
				c.JSON(http.StatusOK, records)
				return
			}

			// 限制抽取数量不超过可用用户数
			if remainingCount > len(otherUsers) {
				remainingCount = len(otherUsers)
			}

			// 随机选择剩余的用户
			selectedIndices := utils.RandomIndices(len(otherUsers), remainingCount)

			// 为每个选中的用户执行抽奖
			for _, index := range selectedIndices {
				selectedUser := otherUsers[index]
				record, err := utils.DrawLotteryWithLevel(config.DB, &selectedUser, ip, company.ID, levelID)
				if err != nil {
					continue // 跳过失败的抽奖
				}
				records = append(records, *record)
			}
		}
	} else {
		// 未指定用户，随机选择
		// 限制抽取人数不超过用户数
		if drawCount > len(users) {
			drawCount = len(users)
		}

		// 随机选择用户
		selectedIndices := utils.RandomIndices(len(users), drawCount)

		// 为每个选中的用户执行抽奖，强制使用指定的奖项等级
		for _, index := range selectedIndices {
			selectedUser := users[index]
			record, err := utils.DrawLotteryWithLevel(config.DB, &selectedUser, ip, company.ID, levelID)
			if err != nil {
				continue // 跳过失败的抽奖
			}
			records = append(records, *record)
		}
	}

	// 如果所有抽奖都失败了
	if len(records) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "抽奖失败"})
		return
	}

	// 加载关联数据
	for i := range records {
		config.DB.Preload("User").Preload("Level").Preload("Prize").First(&records[i], records[i].ID)
	}

	c.JSON(http.StatusOK, records)
}

// GetMyPrize 获取我的奖品
func GetMyPrize(c *gin.Context) {
	phone := c.Query("phone")
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

	var user models.User
	if err := config.DB.Where("phone = ? AND company_id = ?", phone, company.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !user.HasDrawn {
		c.JSON(http.StatusOK, gin.H{"message": "You haven't drawn yet"})
		return
	}

	var record models.DrawRecord
	config.DB.Where("user_id = ? AND company_id = ?", user.ID, company.ID).
		Preload("Level").
		Preload("Prize").
		First(&record)

	c.JSON(http.StatusOK, record)
}

// GetUserStats 获取用户统计（公开API）
func GetUserStats(c *gin.Context) {
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

	// 统计未抽奖的用户数量
	var undrawnCount int64
	config.DB.Model(&models.User{}).
		Where("company_id = ? AND has_drawn = ?", company.ID, false).
		Count(&undrawnCount)

	// 统计总用户数
	var totalCount int64
	config.DB.Model(&models.User{}).
		Where("company_id = ?", company.ID).
		Count(&totalCount)

	// 统计已抽奖的用户数
	var drawnCount int64
	config.DB.Model(&models.User{}).
		Where("company_id = ? AND has_drawn = ?", company.ID, true).
		Count(&drawnCount)

	c.JSON(http.StatusOK, gin.H{
		"total_users":   totalCount,
		"undrawn_users": undrawnCount,
		"drawn_users":   drawnCount,
	})
}

// GetDrawRecordsPublic 获取抽奖记录（公开API）
func GetDrawRecordsPublic(c *gin.Context) {
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

	var records []models.DrawRecord
	config.DB.Where("company_id = ?", company.ID).
		Preload("User").
		Preload("Level").
		Preload("Prize").
		Order("created_at DESC").
		Limit(100).
		Find(&records)

	c.JSON(http.StatusOK, records)
}

// GetAvailableUsersPublic 获取未抽奖的用户列表（公开API）
func GetAvailableUsersPublic(c *gin.Context) {
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

	var users []models.User
	config.DB.Where("company_id = ? AND has_drawn = ?", company.ID, false).
		Order("id ASC").
		Find(&users)

	c.JSON(http.StatusOK, users)
}
