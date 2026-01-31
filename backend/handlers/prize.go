package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"lottery-system/config"
	"lottery-system/models"

	"github.com/gin-gonic/gin"
)

// CreatePrizeLevel 创建奖项等级（权限检查）
func CreatePrizeLevel(c *gin.Context) {
	var level models.PrizeLevel
	if err := c.ShouldBindJSON(&level); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 检查权限 - 普通管理员只能为本公司创建奖项
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，强制使用自己的公司ID
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}
		level.CompanyID = int(*companyID.(*int))
	}

	// 验证公司是否存在
	var company models.Company
	if err := config.DB.First(&company, level.CompanyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company not found"})
		return
	}

	// 库存由奖品管理，奖项等级的库存字段设置为0
	level.TotalStock = 0
	level.UsedStock = 0

	if err := config.DB.Create(&level).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败，请稍后重试"})
		return
	}

	// 记录操作日志
	resourceID := uint(level.ID)
	LogOperation(c, "create", "prize_level", &resourceID, fmt.Sprintf("创建奖项等级: %s", level.Name))

	c.JSON(http.StatusCreated, level)
}

// GetPrizeLevels 获取所有奖项等级（权限隔离）
func GetPrizeLevels(c *gin.Context) {
	var levels []models.PrizeLevel

	// 检查是否是超级管理员
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只返回自己公司的奖项
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}
		config.DB.Preload("Company").Where("company_id = ?", companyID).Order("sort_order ASC").Find(&levels)
	} else {
		// 超级管理员，返回所有奖项
		config.DB.Preload("Company").Order("sort_order ASC").Find(&levels)
	}

	c.JSON(http.StatusOK, levels)
}

// UpdatePrizeLevel 更新奖项等级（权限检查）
func UpdatePrizeLevel(c *gin.Context) {
	id := c.Param("id")

	var level models.PrizeLevel
	if err := config.DB.First(&level, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "奖项不存在"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能更新自己公司的奖项
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != level.CompanyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	var req models.PrizeLevel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 普通管理员不能修改公司关联
	if !isSuperAdmin.(bool) && req.CompanyID != level.CompanyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify company association"})
		return
	}

	// 库存由奖品管理，不允许通过此接口修改
	// 只允许更新名称、描述、概率、排序、状态
	updateData := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"probability": req.Probability,
		"sort_order":  req.SortOrder,
		"is_active":   req.IsActive,
	}

	if err := config.DB.Model(&level).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败，请稍后重试"})
		return
	}

	// 记录操作日志
	resourceID := uint(level.ID)
	LogOperation(c, "update", "prize_level", &resourceID, fmt.Sprintf("更新奖项等级: %s", level.Name))

	c.JSON(http.StatusOK, level)
}

// DeletePrizeLevel 删除奖项等级（权限检查）
func DeletePrizeLevel(c *gin.Context) {
	id := c.Param("id")

	var level models.PrizeLevel
	if err := config.DB.First(&level, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize level not found"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能删除自己公司的奖项
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != level.CompanyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	if err := config.DB.Delete(&level).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete prize level"})
		return
	}

	// 记录操作日志
	resourceID := uint(level.ID)
	LogOperation(c, "delete", "prize_level", &resourceID, fmt.Sprintf("删除奖项等级: %s", level.Name))

	c.JSON(http.StatusOK, gin.H{"message": "Prize level deleted successfully"})
}

// CreatePrize 创建具体奖品（权限检查）
func CreatePrize(c *gin.Context) {
	var prize models.Prize
	if err := c.ShouldBindJSON(&prize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证奖品名称不为空
	if prize.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "奖品名称不能为空"})
		return
	}

	// 验证库存：总库存必须 >= 已发放
	if prize.TotalStock < prize.UsedStock {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("总库存 (%d) 不能小于已发放 (%d)", prize.TotalStock, prize.UsedStock)})
		return
	}

	// 检查奖项等级是否存在
	var level models.PrizeLevel
	if err := config.DB.First(&level, prize.LevelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize level not found"})
		return
	}

	// 检查权限 - 普通管理员只能为自己公司的奖项创建奖品
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能为本公司的奖项创建奖品
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != level.CompanyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	if err := config.DB.Create(&prize).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prize"})
		return
	}

	// 记录操作日志
	resourceID := uint(prize.ID)
	LogOperation(c, "create", "prize", &resourceID, fmt.Sprintf("创建奖品: %s", prize.Name))

	c.JSON(http.StatusCreated, prize)
}

// GetPrizesByLevel 获取指定等级的奖品列表（权限检查）
func GetPrizesByLevel(c *gin.Context) {
	levelID := c.Param("levelId")

	// 首先检查奖项等级的权限
	var level models.PrizeLevel
	if err := config.DB.First(&level, levelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize level not found"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能查看自己公司的奖项的奖品
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != level.CompanyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	var prizes []models.Prize
	config.DB.Where("level_id = ?", levelID).Find(&prizes)

	c.JSON(http.StatusOK, prizes)
}

// UpdatePrize 更新奖品（权限检查）
func UpdatePrize(c *gin.Context) {
	id := c.Param("id")

	var prize models.Prize
	if err := config.DB.First(&prize, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize not found"})
		return
	}

	// 检查奖项等级的权限
	var level models.PrizeLevel
	if err := config.DB.First(&level, prize.LevelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize level not found"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能更新自己公司奖项的奖品
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != level.CompanyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	var req models.Prize
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证库存：总库存必须 >= 已发放
	if req.TotalStock < req.UsedStock {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("总库存 (%d) 不能小于已发放 (%d)", req.TotalStock, req.UsedStock)})
		return
	}

	// 如果修改了奖项等级，需要检查新奖项的权限
	if req.LevelID != prize.LevelID {
		var newLevel models.PrizeLevel
		if err := config.DB.First(&newLevel, req.LevelID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "New prize level not found"})
			return
		}

		if !isSuperAdmin.(bool) {
			companyID, _ := c.Get("company_id")
			if int(*companyID.(*int)) != newLevel.CompanyID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Cannot move to different company's prize level"})
				return
			}
		}
	}

	if err := config.DB.Model(&prize).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prize"})
		return
	}

	// 记录操作日志
	resourceID := uint(prize.ID)
	LogOperation(c, "update", "prize", &resourceID, fmt.Sprintf("更新奖品: %s", prize.Name))

	c.JSON(http.StatusOK, prize)
}

// DeletePrize 删除奖品（权限检查）
func DeletePrize(c *gin.Context) {
	id := c.Param("id")

	var prize models.Prize
	if err := config.DB.First(&prize, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize not found"})
		return
	}

	// 检查奖项等级的权限
	var level models.PrizeLevel
	if err := config.DB.First(&level, prize.LevelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize level not found"})
		return
	}

	// 检查权限
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能删除自己公司奖项的奖品
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil || int(*companyID.(*int)) != level.CompanyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}
	}

	if err := config.DB.Delete(&prize).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete prize"})
		return
	}

	// 记录操作日志
	resourceID := uint(prize.ID)
	LogOperation(c, "delete", "prize", &resourceID, fmt.Sprintf("删除奖品: %s", prize.Name))

	c.JSON(http.StatusOK, gin.H{"message": "Prize deleted successfully"})
}

// GetAllPrizes 获取所有奖品（权限隔离）
func GetAllPrizes(c *gin.Context) {
	// 检查是否是超级管理员
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只返回自己公司的奖品
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}

		var prizes []models.Prize
		err := config.DB.Raw(`
			SELECT p.* FROM prizes p
			INNER JOIN prize_levels l ON p.level_id = l.id
			WHERE l.company_id = ?
			ORDER BY l.sort_order ASC, p.id ASC
		`, companyID).Find(&prizes).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prizes"})
			return
		}

		c.JSON(http.StatusOK, prizes)
	} else {
		// 超级管理员，返回所有奖品
		var prizes []models.Prize
		if err := config.DB.Raw(`
			SELECT p.* FROM prizes p
			INNER JOIN prize_levels l ON p.level_id = l.id
			ORDER BY l.company_id ASC, l.sort_order ASC, p.id ASC
		`).Find(&prizes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prizes"})
			return
		}

		c.JSON(http.StatusOK, prizes)
	}
}

// GetDrawRecords 获取抽奖记录（权限隔离）
func GetDrawRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100")) // 默认100条
	search := c.Query("search")
	companyIDParam := c.Query("company_id")

	offset := (page - 1) * pageSize

	query := config.DB.Model(&models.DrawRecord{})

	// 检查是否是超级管理员
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只能查看自己公司的记录
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

	if search != "" {
		query = query.Joins("JOIN users ON draw_records.user_id = users.id").
			Where("users.phone LIKE ? OR users.name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var records []models.DrawRecord
	query.Preload("User").Preload("Level").Preload("Prize").Preload("Company").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"data":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetStats 获取统计数据（权限隔离）
func GetStats(c *gin.Context) {
	var totalUsers int64
	var drawnUsers int64
	var totalRecords int64
	var levels []models.PrizeLevel

	// 检查是否是超级管理员
	isSuperAdmin, exists := c.Get("is_super_admin")
	if !exists || !isSuperAdmin.(bool) {
		// 普通管理员，只统计自己公司的数据
		companyID, exists := c.Get("company_id")
		if !exists || companyID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company assigned"})
			return
		}

		config.DB.Model(&models.User{}).Where("company_id = ?", companyID).Count(&totalUsers)
		config.DB.Model(&models.User{}).Where("company_id = ? AND has_drawn = ?", companyID, true).Count(&drawnUsers)
		config.DB.Model(&models.DrawRecord{}).Where("company_id = ?", companyID).Count(&totalRecords)
		config.DB.Where("company_id = ?", companyID).Find(&levels)
	} else {
		// 超级管理员，统计所有数据
		config.DB.Model(&models.User{}).Count(&totalUsers)
		config.DB.Model(&models.User{}).Where("has_drawn = ?", true).Count(&drawnUsers)
		config.DB.Model(&models.DrawRecord{}).Count(&totalRecords)
		config.DB.Find(&levels)
	}

	// 为每个奖项等级计算奖品的库存信息
	type LevelWithStock struct {
		models.PrizeLevel
		TotalStock int `json:"total_stock"`
		UsedStock  int `json:"used_stock"`
	}

	result := make([]LevelWithStock, len(levels))
	for i, level := range levels {
		var stockData struct {
			TotalStock int `json:"total_stock"`
			UsedStock  int `json:"used_stock"`
		}
		config.DB.Model(&models.Prize{}).
			Where("level_id = ?", level.ID).
			Select("COALESCE(SUM(total_stock), 0) as total_stock, COALESCE(SUM(used_stock), 0) as used_stock").
			Scan(&stockData)

		result[i] = LevelWithStock{
			PrizeLevel: level,
			TotalStock: stockData.TotalStock,
			UsedStock:  stockData.UsedStock,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_users":   totalUsers,
		"drawn_users":   drawnUsers,
		"total_records": totalRecords,
		"levels":        result,
	})
}
