package handlers

import (
	"fmt"
	"net/http"

	"lottery-system/config"
	"lottery-system/models"

	"github.com/gin-gonic/gin"
)

// GetCompanyByCode 根据公司代码获取公司信息
func GetCompanyByCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company code is required"})
		return
	}

	var company models.Company
	if err := config.DB.Where("code = ? AND is_active = ?", code, true).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// CreateCompany 创建公司（超级管理员）
func CreateCompany(c *gin.Context) {
	// 检查权限：只有超级管理员可以创建公司
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有超级管理员可以创建公司"})
		return
	}

	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	// 记录操作日志
	resourceID := uint(company.ID)
	LogOperation(c, "create", "company", &resourceID, fmt.Sprintf("创建公司: %s (代码: %s)", company.Name, company.Code))

	c.JSON(http.StatusCreated, company)
}

// UpdateCompany 更新公司信息
func UpdateCompany(c *gin.Context) {
	id := c.Param("id")

	var company models.Company
	if err := config.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能更新自己的公司
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != company.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	var req models.Company
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&company).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	// 记录操作日志
	resourceID := uint(company.ID)
	LogOperation(c, "update", "company", &resourceID, fmt.Sprintf("更新公司: %s", company.Name))

	c.JSON(http.StatusOK, company)
}

// CompanyWithStats 带统计数据的公司信息
type CompanyWithStats struct {
	models.Company
	TotalUsers  int64 `json:"total_users"`
	DrawnCount  int64 `json:"drawn_count"`
}

// GetCompanies 获取公司列表（超级管理员可以看到所有，普通管理员只能看到自己的）
func GetCompanies(c *gin.Context) {
	var companies []models.Company

	// 检查是否是超级管理员
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只返回自己所属的公司
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}
		config.DB.Where("id = ?", companyID).Find(&companies)
	} else {
		// 超级管理员，返回所有公司
		config.DB.Find(&companies)
	}

	// 为每个公司附加统计数据
	var companiesWithStats []CompanyWithStats
	for _, company := range companies {
		// 统计总用户数
		var totalUsers int64
		config.DB.Model(&models.User{}).Where("company_id = ?", company.ID).Count(&totalUsers)

		// 统计中奖次数（抽奖记录数）
		var drawnCount int64
		config.DB.Model(&models.DrawRecord{}).Where("company_id = ?", company.ID).Count(&drawnCount)

		companiesWithStats = append(companiesWithStats, CompanyWithStats{
			Company:     company,
			TotalUsers:  totalUsers,
			DrawnCount:  drawnCount,
		})
	}

	c.JSON(http.StatusOK, companiesWithStats)
}

// GetCompanyStats 获取公司统计数据
func GetCompanyStats(c *gin.Context) {
	companyID := c.Query("company_id")

	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能查看自己所属公司的统计
		adminCompanyID, exists := c.Get("company_id")
		if !exists || adminCompanyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
		// 将查询的company_id转为字符串进行比较
		adminCompanyIDStr := string(rune(*adminCompanyID.(*int)))
		if adminCompanyIDStr != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	var totalUsers int64
	config.DB.Model(&models.User{}).Where("company_id = ?", companyID).Count(&totalUsers)

	var drawnUsers int64
	config.DB.Model(&models.User{}).Where("company_id = ? AND has_drawn = ?", companyID, true).Count(&drawnUsers)

	var totalRecords int64
	config.DB.Model(&models.DrawRecord{}).Where("company_id = ?", companyID).Count(&totalRecords)

	c.JSON(http.StatusOK, gin.H{
		"total_users":   totalUsers,
		"drawn_users":   drawnUsers,
		"total_records": totalRecords,
	})
}

// DeleteCompany 删除公司（超级管理员）
func DeleteCompany(c *gin.Context) {
	// 检查权限：只有超级管理员可以删除公司
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有超级管理员可以删除公司"})
		return
	}

	id := c.Param("id")

	var company models.Company
	if err := config.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	// 检查公司下是否还有用户
	var userCount int64
	config.DB.Model(&models.User{}).Where("company_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("该公司下还有 %d 个用户，无法删除", userCount),
		})
		return
	}

	// 检查公司下是否还有抽奖记录
	var recordCount int64
	config.DB.Model(&models.DrawRecord{}).Where("company_id = ?", id).Count(&recordCount)
	if recordCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("该公司下还有 %d 条抽奖记录，无法删除", recordCount),
		})
		return
	}

	if err := config.DB.Delete(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	// 记录操作日志
	resourceID := uint(company.ID)
	LogOperation(c, "delete", "company", &resourceID, fmt.Sprintf("删除公司: %s (代码: %s)", company.Name, company.Code))

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
