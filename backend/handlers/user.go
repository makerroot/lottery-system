package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"lottery-system/config"
	"lottery-system/models"
	"lottery-system/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	CompanyID int    `json:"company_id" binding:"required"`
	Username  string `json:"username"` // 可选：如果不提供，系统不创建可登录账号
	Password  string `json:"password"` // 可选：如果不提供，系统不创建可登录账号
	Name      string `json:"name"`     // 必填：姓名
	Phone     string `json:"phone"`    // 可选：手机号
}

// BatchCreateUserRequest 批量创建用户请求
type BatchCreateUserRequest struct {
	CompanyID int      `json:"company_id" binding:"required"`
	Users     []string `json:"users" binding:"required"` // 格式: ["用户名,密码,姓名", ...]
}

// CreateUser 创建单个用户（权限检查）
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 检查权限 - 普通管理员只能为本公司创建用户
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，强制使用自己的公司ID
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}
		req.CompanyID = int(*companyID.(*int))
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

	// 检查公司是否存在
	var company models.Company
	if err := config.DB.First(&company, req.CompanyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "公司不存在"})
		return
	}

	var user models.User

	// 情况1：提供了 username 和 password -> 创建可登录的用户（扫码注册）
	if req.Username != "" && req.Password != "" {
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

		// 检查用户名是否已存在
		finalUsername := req.Username
		var existingUsers []models.User
		query := config.DB.Where("company_id = ?", req.CompanyID).Where("username = ?", req.Username)

		if req.Phone != "" {
			query = query.Where("phone = ?", req.Phone)
		}

		if err := query.Find(&existingUsers).Error; err == nil && len(existingUsers) > 0 {
			// 有手机号且找到用户：认为已存在
			if req.Phone != "" {
				c.JSON(http.StatusConflict, gin.H{
					"error":          "该用户名和手机号的用户已存在",
					"existing_users": existingUsers,
				})
				return
			}

			// 没有手机号但有重名用户：自动添加序号
			var count int64
			config.DB.Model(&models.User{}).Where("username = ? AND company_id = ?", req.Username, req.CompanyID).Count(&count)
			finalUsername = fmt.Sprintf("%s_%d", req.Username, count+1)

			utils.WithFields(map[string]interface{}{
				"original_username": req.Username,
				"final_username":    finalUsername,
				"count":             count,
			}).Info("检测到重名用户，自动添加序号后缀")
		}

		// 哈希密码
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}

		// 创建可登录的用户
		user = models.User{
			CompanyID: req.CompanyID,
			Username:  finalUsername,
			Password:  hashedPassword,
			Role:      models.RoleUser,
			Name:      req.Name,
			Phone:     req.Phone,
			HasDrawn:  false,
		}
	} else {
		// 情况2：只提供了 name 和 phone -> 创建不可登录的用户（管理员添加）
		// 自动生成 username（使用手机号或时间戳+随机数）
		var username string
		if req.Phone != "" {
			username = req.Phone
		} else {
			username = fmt.Sprintf("u_%d_%d", time.Now().Unix(), utils.RandomInt(10000))
		}

		// 生成随机密码（用户无法登录，但密码字段不能为空）
		randomPassword := utils.GenerateRandomPassword(8)
		hashedPassword, err := utils.HashPassword(randomPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}

		// 检查用户名是否已存在
		var existingUserCount int64
		config.DB.Model(&models.User{}).Where("company_id = ? AND username = ?", req.CompanyID, username).Count(&existingUserCount)
		if existingUserCount > 0 {
			// 用户名重复，添加序号
			username = fmt.Sprintf("%s_%d", username, existingUserCount+1)
		}

		user = models.User{
			CompanyID: req.CompanyID,
			Username:  username,
			Password:  hashedPassword,
			Role:      models.RoleUser,
			Name:      req.Name,
			Phone:     req.Phone,
			HasDrawn:  false,
		}
	}

	// 创建用户
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败，请稍后重试"})
		return
	}

	// 记录操作日志
	resourceID := uint(user.ID)
	details := fmt.Sprintf("创建用户: %s", user.Name)
	if user.Username != "" {
		details += fmt.Sprintf(" (@%s)", user.Username)
	}
	LogOperation(c, "create", "user", &resourceID, details)

	// 如果用户名被修改了，返回提示
	response := map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"name":      user.Name,
		"phone":     user.Phone,
		"has_drawn": false,
		"can_login": user.Username != "", // 是否可以登录
	}

	// 如果用户名为空，说明是管理员添加的抽奖用户
	if user.Username == "" {
		response["message"] = "用户已添加到抽奖池（无法登录）"
	}

	c.JSON(http.StatusOK, response)
}

// BatchCreateUsers 批量创建用户（权限检查）
func BatchCreateUsers(c *gin.Context) {
	var req BatchCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限 - 普通管理员只能为本公司批量创建用户
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，强制使用自己的公司ID
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}
		req.CompanyID = int(*companyID.(*int))
	}

	// 检查公司是否存在
	var company models.Company
	if err := config.DB.First(&company, req.CompanyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company not found"})
		return
	}

	var createdUsers []models.User
	var failedUsers []string

	// 用于批量导入时的用户名计数器
	batchIndex := 0
	baseTimestamp := time.Now().Unix()

	for _, userStr := range req.Users {
		// 解析格式: "姓名,手机号（可选）"
		var name, phone string
		if len(userStr) > 0 {
			parts := strings.Split(userStr, ",")
			name = strings.TrimSpace(parts[0])
			if len(parts) >= 2 {
				phone = strings.TrimSpace(parts[1])
			}
		}

		// 验证姓名
		if name == "" {
			failedUsers = append(failedUsers, userStr+" (姓名为空)")
			continue
		}

		// 验证手机号（如果提供）
		if phone != "" {
			if err := utils.ValidatePhone(phone); err != nil {
				failedUsers = append(failedUsers, name+" ("+err.Error()+")")
				continue
			}
		}

		// 检查是否已存在（根据姓名和手机号）
		var existingUser models.User
		query := config.DB.Where("company_id = ? AND name = ?", req.CompanyID, name)
		if phone != "" {
			query = query.Where("phone = ?", phone)
		}
		if err := query.First(&existingUser).Error; err == nil {
			failedUsers = append(failedUsers, name+" (已存在)")
			continue
		}

		// 创建用户（自动生成 username 和 password）
		// 生成 username：优先使用手机号，否则使用时间戳+批量索引
		var username string
		if phone != "" {
			username = phone
		} else {
			// 使用批量索引确保唯一性
			batchIndex++
			username = fmt.Sprintf("u_%d_%d", baseTimestamp, batchIndex)
		}

		// 生成随机密码（用户无法登录，但密码字段不能为空）
		randomPassword := utils.GenerateRandomPassword(8)
		hashedPassword, err := utils.HashPassword(randomPassword)
		if err != nil {
			failedUsers = append(failedUsers, name+" (密码加密失败)")
			continue
		}

		// 检查用户名是否已存在
		var existingUserCount int64
		config.DB.Model(&models.User{}).Where("company_id = ? AND username = ?", req.CompanyID, username).Count(&existingUserCount)
		if existingUserCount > 0 {
			// 用户名重复，添加序号
			username = fmt.Sprintf("%s_%d", username, existingUserCount+1)
		}

		user := models.User{
			CompanyID: req.CompanyID,
			Username:  username,
			Password:  hashedPassword,
			Role:      models.RoleUser,
			Name:      name,
			Phone:     phone,
			HasDrawn:  false,
		}

		if err := config.DB.Create(&user).Error; err != nil {
			failedUsers = append(failedUsers, name)
			continue
		}

		createdUsers = append(createdUsers, user)
	}

	// 记录操作日志（批量创建）
	if len(createdUsers) > 0 {
		details := fmt.Sprintf("批量创建用户: 成功%d个", len(createdUsers))
		if len(failedUsers) > 0 {
			details += fmt.Sprintf(", 失败%d个", len(failedUsers))
		}
		// 使用第一个用户的ID作为资源ID，或者记录0表示批量操作
		var resourceID *uint
		if len(createdUsers) > 0 {
			rid := uint(createdUsers[0].ID)
			resourceID = &rid
		}
		LogOperation(c, "create", "user", resourceID, details)
	}

	c.JSON(http.StatusOK, gin.H{
		"created": len(createdUsers),
		"failed":  len(failedUsers),
		"users":   createdUsers,
		"errors":  failedUsers,
	})
}

// GetUsers 获取用户列表（权限隔离）
func GetUsers(c *gin.Context) {
	companyIDParam := c.Query("company_id")
	hasDrawn := c.Query("has_drawn")

	query := config.DB.Model(&models.User{})

	// 检查是否是超级管理员
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能查看自己公司的用户
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}
		query = query.Where("company_id = ?", companyID)
	} else {
		// 超级管理员，可以按公司过滤
		if companyIDParam != "" {
			query = query.Where("company_id = ?", companyIDParam)
		}
	}

	if hasDrawn != "" {
		query = query.Where("has_drawn = ?", hasDrawn)
	}

	var users []models.User
	query.Order("id ASC").Find(&users)

	c.JSON(http.StatusOK, users)
}

// DeleteUser 删除用户（权限检查）
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能删除自己公司的用户
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != user.CompanyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// 记录操作日志
	resourceID := uint(user.ID)
	LogOperation(c, "delete", "user", &resourceID, fmt.Sprintf("删除用户: %s (@%s)", user.Name, user.Username))

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	HasDrawn *bool  `json:"has_drawn"`
}

// UpdateUser 更新用户（权限检查）
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能更新自己公司的用户
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != user.CompanyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 更新字段
	updates := map[string]interface{}{}

	if req.Name != "" {
		if err := utils.ValidateName(req.Name); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "姓名格式错误: " + err.Error()})
			return
		}
		updates["name"] = req.Name
	}

	if req.Phone != "" {
		if err := utils.ValidatePhone(req.Phone); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updates["phone"] = req.Phone
	}

	if req.HasDrawn != nil {
		updates["has_drawn"] = *req.HasDrawn
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有要更新的字段"})
		return
	}

	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	// 记录操作日志
	resourceID := uint(user.ID)
	details := fmt.Sprintf("更新用户: %s (@%s)", user.Name, user.Username)
	if req.Name != "" {
		details += fmt.Sprintf(" → 姓名: %s", req.Name)
	}
	LogOperation(c, "update", "user", &resourceID, details)

	c.JSON(http.StatusOK, user)
}

// ChangeUserPasswordRequest 修改用户密码请求
type ChangeUserPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangeUserPassword 用户修改自己的密码
func ChangeUserPassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req ChangeUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "当前密码错误"})
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 更新密码
	user.Password = hashedPassword
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}

	utils.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("用户修改密码成功")

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// ScanAddUserRequest 扫码添加用户请求
type ScanAddUserRequest struct {
	CompanyCode string `json:"company_code" binding:"required"`
	QRCodeData  string `json:"qr_code_data" binding:"required"` // 二维码内容
}

// ScanAddUser 扫码添加用户（管理员权限）
func ScanAddUser(c *gin.Context) {
	// 🔒 权限检查：只允许管理员扫码添加用户
	isAdmin, _ := c.Get("is_admin")
	isSuperAdmin, _ := c.Get("is_super_admin")

	if isAdmin == false && isSuperAdmin == false {
		c.JSON(http.StatusForbidden, gin.H{
			"error":      "只有管理员才能扫码添加用户",
			"error_code": "PERMISSION_DENIED",
		})
		return
	}

	var req ScanAddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 获取管理员的公司ID
	var companyID int
	if isSuperAdmin == true {
		// 超级管理员，使用请求中的公司代码
		var company models.Company
		if err := config.DB.Where("code = ?", req.CompanyCode).First(&company).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "公司不存在"})
			return
		}
		companyID = int(company.ID)
	} else {
		// 普通管理员，使用自己的公司ID
		cid, exists := c.Get("company_id")
		if !exists || cid == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}
		companyID = int(*cid.(*int))
	}

	// 解析二维码数据
	// 支持两种格式：
	// 1. JSON格式: {"username":"zhangsan","name":"张三","phone":"13800138000"}
	// 2. 简单格式: username:zhangsan,name:张三,phone:13800138000

	var username, name, phone string

	// 尝试解析为JSON
	if strings.HasPrefix(req.QRCodeData, "{") {
		var qrData map[string]string
		if err := json.Unmarshal([]byte(req.QRCodeData), &qrData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "二维码格式错误，无法解析JSON"})
			return
		}

		username = qrData["username"]
		name = qrData["name"]
		phone = qrData["phone"]
	} else {
		// 解析简单格式: key:value,key:value
		pairs := strings.Split(req.QRCodeData, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				switch key {
				case "username":
					username = value
				case "name":
					name = value
				case "phone":
					phone = value
				}
			}
		}
	}

	// 验证必填字段
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "二维码中缺少username字段"})
		return
	}

	if name == "" {
		name = username // 如果没有姓名，使用用户名
	}

	// 检查用户是否已存在
	// 策略：如果有手机号，用 (username, phone) 判断；如果没有手机号，允许重名
	var existingUsers []models.User
	query := config.DB.Where("company_id = ?", companyID).Where("username = ?", username)

	if phone != "" {
		// 有手机号：检查 (username, phone) 组合
		query = query.Where("phone = ?", phone)
	}

	if err := query.Find(&existingUsers).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	// 如果找到用户
	if len(existingUsers) > 0 {
		// 检查是否已抽奖
		for _, existingUser := range existingUsers {
			if existingUser.HasDrawn {
				c.JSON(http.StatusConflict, gin.H{
					"error": "该用户已经抽过奖",
					"user": gin.H{
						"id":        existingUser.ID,
						"username":  existingUser.Username,
						"name":      existingUser.Name,
						"phone":     existingUser.Phone,
						"has_drawn": true,
					},
				})
				return
			}
		}

		// 有用户存在但未抽奖，返回第一个
		existingUser := existingUsers[0]
		c.JSON(http.StatusOK, gin.H{
			"message": "用户已在抽奖池中",
			"user": gin.H{
				"id":        existingUser.ID,
				"username":  existingUser.Username,
				"name":      existingUser.Name,
				"phone":     existingUser.Phone,
				"has_drawn": false,
			},
		})
		return
	}

	// 如果没有手机号且username重复，自动添加序号
	if phone == "" {
		var count int64
		config.DB.Model(&models.User{}).Where("username = ? AND company_id = ?", username, companyID).Count(&count)
		if count > 0 {
			// 添加序号后缀
			finalUsername := fmt.Sprintf("%s_%d", username, count+1)

			utils.WithFields(map[string]interface{}{
				"original_username": username,
				"final_username":    finalUsername,
				"count":             count,
			}).Info("检测到重名用户，自动添加序号后缀")

			username = finalUsername
		}
	}

	// 创建新用户
	// 生成随机密码（用户无法登录，密码不重要）
	randomPassword := utils.GenerateRandomPassword(8)
	hashedPassword, err := utils.HashPassword(randomPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		CompanyID: companyID,
		Username:  username,
		Password:  hashedPassword,
		Role:      models.RoleUser,
		Name:      name,
		Phone:     phone,
		HasDrawn:  false,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.WithFields(map[string]interface{}{
			"error":      err,
			"username":   username,
			"company_id": companyID,
		}).Error("创建用户失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	utils.WithFields(map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"name":       user.Name,
		"company_id": companyID,
	}).Info("扫码添加用户成功")

	c.JSON(http.StatusOK, gin.H{
		"message": "添加用户成功",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"name":      user.Name,
			"phone":     user.Phone,
			"has_drawn": false,
		},
	})
}
