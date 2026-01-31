package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"lottery-system/config"
	"lottery-system/models"
	"lottery-system/utils"

	"github.com/gin-gonic/gin"
)

// CreateAdminRequest 创建管理员请求
type CreateAdminRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	CompanyID    *int   `json:"company_id"`
}

// UpdateAdminRequest 更新管理员请求
type UpdateAdminRequest struct {
	Username     *string `json:"username"`
	Password     string  `json:"password"`
	IsSuperAdmin *bool   `json:"is_super_admin"`
	CompanyID    *int    `json:"company_id"`
}

// GetAdmins 获取管理员列表
func GetAdmins(c *gin.Context) {
	var admins []models.Admin

	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能看到自己
		userID, _ := c.Get("user_id")
		config.DB.Where("id = ?", userID).Preload("Company").Find(&admins)
	} else {
		// 超级管理员，可以看到所有管理员
		config.DB.Preload("Company").Find(&admins)
	}

	c.JSON(http.StatusOK, admins)
}

// CreateAdmin 创建管理员（仅超级管理员）
func CreateAdmin(c *gin.Context) {
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingAdmin models.Admin
	if err := config.DB.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// 如果指定了公司，检查公司是否存在
	if req.CompanyID != nil {
		var company models.Company
		if err := config.DB.First(&company, *req.CompanyID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Company not found"})
			return
		}
	}

	// 验证：如果设置为普通管理员，必须指定公司
	if !req.IsSuperAdmin && req.CompanyID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "普通管理员必须指定所属公司"})
		return
	}

	// 验证：如果设置为超级管理员，不应该指定公司
	if req.IsSuperAdmin && req.CompanyID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "超级管理员不应该指定所属公司"})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	admin := models.Admin{
		Username:     req.Username,
		Password:     hashedPassword,
		CompanyID:    req.CompanyID,
		IsSuperAdmin: req.IsSuperAdmin,
	}

	// 根据是否是超级管理员设置角色
	if req.IsSuperAdmin {
		admin.Role = models.RoleSuperAdmin
	} else {
		admin.Role = models.RoleAdmin
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		utils.WithFields(map[string]interface{}{
			"error": err,
			"username": req.Username,
			"company_id": req.CompanyID,
		}).Error("创建管理员失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	utils.WithFields(map[string]interface{}{
		"admin_id": admin.ID,
		"username": admin.Username,
		"company_id": admin.CompanyID,
	}).Info("管理员创建成功")

	// 记录操作日志（在创建成功后）
	details := fmt.Sprintf("创建管理员: %s", admin.Username)
	if admin.CompanyID != nil {
		details += fmt.Sprintf(" (公司ID: %d)", *admin.CompanyID)
	}
	resourceID := uint(admin.ID)
	LogOperation(c, "create", "admin", &resourceID, details)

	// 重新查询以获取关联的公司信息
	var result models.Admin
	if err := config.DB.Preload("Company").First(&result, admin.ID).Error; err != nil {
		// 如果查询失败，至少返回已创建的管理员信息
		c.JSON(http.StatusCreated, admin)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// UpdateAdmin 更新管理员信息
func UpdateAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	var admin models.Admin
	if err := config.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	userID, _ := c.Get("user_id")

	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员只能修改自己
		if int(userID.(int)) != id {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
		// 普通管理员不能修改自己的公司关联
		if admin.CompanyID != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify company association"})
			return
		}
	}

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新用户名（仅超级管理员）
	if isSuperAdmin.(bool) && req.Username != nil {
		// 检查新用户名是否已被其他管理员使用
		var existingAdmin models.Admin
		if err := config.DB.Where("username = ? AND id != ?", *req.Username, id).First(&existingAdmin).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		admin.Username = *req.Username
	}

	// 更新密码
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		admin.Password = hashedPassword
	}

	// 更新管理员身份和公司关联（仅超级管理员）
	if isSuperAdmin.(bool) {
		if req.IsSuperAdmin != nil {
			// 如果要设置为普通管理员
			if !*req.IsSuperAdmin {
				// 必须指定公司
				if req.CompanyID == nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "普通管理员必须指定所属公司"})
					return
				}
				// 检查公司是否存在
				var company models.Company
				if err := config.DB.First(&company, *req.CompanyID).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Company not found"})
					return
				}
				admin.CompanyID = req.CompanyID
				admin.IsSuperAdmin = false
				admin.Role = models.RoleAdmin
			} else {
				// 设置为超级管理员
				admin.CompanyID = nil
				admin.IsSuperAdmin = true
				admin.Role = models.RoleSuperAdmin
			}
		} else if req.CompanyID != nil {
			// 只更新公司，保持原有的超级管理员状态
			var company models.Company
			if err := config.DB.First(&company, *req.CompanyID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Company not found"})
				return
			}
			admin.CompanyID = req.CompanyID
			admin.IsSuperAdmin = false
			admin.Role = models.RoleAdmin
		}
	}

	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin"})
		return
	}

	// 记录操作日志
	details := fmt.Sprintf("更新管理员: %s", admin.Username)
	if admin.CompanyID != nil {
		details += fmt.Sprintf(" (公司ID: %d)", *admin.CompanyID)
	}
	resourceID := uint(admin.ID)
	LogOperation(c, "update", "admin", &resourceID, details)

	config.DB.Preload("Company").First(&admin, admin.ID)
	c.JSON(http.StatusOK, admin)
}

// DeleteAdmin 删除管理员（仅超级管理员）
func DeleteAdmin(c *gin.Context) {
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	// 不能删除自己
	userID, _ := c.Get("user_id")
	if int(userID.(int)) == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete yourself"})
		return
	}

	if err := config.DB.Delete(&models.Admin{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete admin"})
		return
	}

	// 记录操作日志（异步，不阻塞响应）
	go func() {
		resourceID := uint(id)
		LogOperation(c, "delete", "admin", &resourceID, fmt.Sprintf("删除管理员: ID=%d", id))
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Admin deleted successfully"})
}

// GetAdminInfo 获取当前管理员信息（带公司信息）
func GetAdminInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var admin models.Admin
	if err := config.DB.Preload("Company").First(&admin, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	c.JSON(http.StatusOK, admin)
}

// ChangeAdminPasswordRequest 修改管理员密码请求
type ChangeAdminPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangeAdminPassword 管理员修改自己的密码
func ChangeAdminPassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req ChangeAdminPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	if err := config.DB.First(&admin, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "管理员不存在"})
		return
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, admin.Password) {
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
	admin.Password = hashedPassword
	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}

	utils.WithFields(map[string]interface{}{
		"admin_id": admin.ID,
		"username": admin.Username,
	}).Info("管理员修改密码成功")

	// 记录操作日志
	resourceID := uint(admin.ID)
	LogOperation(c, "update", "admin", &resourceID, "修改密码")

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
