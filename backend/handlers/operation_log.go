package handlers

import (
	"fmt"
	"net/http"
	"time"

	"lottery-system/config"
	"lottery-system/models"

	"github.com/gin-gonic/gin"
)

// LogOperation 记录操作日志
func LogOperation(c *gin.Context, action, resource string, resourceID *uint, details string) {
	adminIDValue, exists := c.Get("user_id")
	if !exists {
		return
	}

	adminNameValue, _ := c.Get("username")
	companyIDValue, _ := c.Get("company_id")

	// 安全地转换类型
	var adminID uint
	switch v := adminIDValue.(type) {
	case int:
		adminID = uint(v)
	case uint:
		adminID = v
	case int64:
		adminID = uint(v)
	default:
		return
	}

	var adminName string
	if adminNameValue != nil {
		if s, ok := adminNameValue.(string); ok {
			adminName = s
		}
	}

	log := models.OperationLog{
		AdminID:    adminID,
		AdminName:  adminName,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	if companyIDValue != nil {
		var cid uint
		switch v := companyIDValue.(type) {
		case int:
			cid = uint(v)
		case uint:
			cid = v
		case *int:
			if v != nil {
				cid = uint(*v)
			}
		}
		log.CompanyID = &cid
	}

	if err := config.DB.Create(&log).Error; err != nil {
		// 记录错误但不中断业务流程
		// 使用fmt输出错误信息
		fmt.Printf("[LogOperation Error] Failed to create log: %v, Action: %s, Resource: %s, AdminID: %d\n",
			err, action, resource, adminID)
	}
}

// GetOperationLogs 获取操作日志（分页）
func GetOperationLogs(c *gin.Context) {
	// 分页参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	limit := 20
	offset := 0

	if p, err := toInt(page); err == nil && p > 0 {
		if ps, err := toInt(pageSize); err == nil && ps > 0 && ps <= 100 {
			limit = ps
			offset = (p - 1) * limit
		}
	}

	// 获取当前管理员信息
	isSuperAdmin, _ := c.Get("is_super_admin")
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 获取当前管理员的公司ID
	companyID, _ := c.Get("company_id")

	// 查询日志
	var logs []models.OperationLog
	var total int64

	query := config.DB.Model(&models.OperationLog{})

	// 权限控制：超级管理员查看所有，普通管理员只查看自己公司的日志
	if !isSuperAdmin.(bool) {
		// 普通管理员：只查看自己公司的操作日志
		if companyID != nil {
			query = query.Where("company_id = ?", companyID)
		} else {
			// 如果没有公司ID，只能查看自己的日志
			query = query.Where("admin_id = ?", adminID)
		}
	}

	// 筛选条件
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if resource := c.Query("resource"); resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if adminID := c.Query("admin_id"); adminID != "" {
		// 只有超级管理员才能按admin_id筛选
		if isSuperAdmin.(bool) {
			query = query.Where("admin_id = ?", adminID)
		}
	}

	// 获取总数
	query.Count(&total)

	// 获取分页数据
	query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"total": total,
		"page": page,
		"page_size": limit,
	})
}

// GetOperationStats 获取操作统计
func GetOperationStats(c *gin.Context) {
	// 获取当前管理员信息
	isSuperAdmin, _ := c.Get("is_super_admin")
	companyID, _ := c.Get("company_id")
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 统计最近7天的操作数
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	var stats []struct {
		Action string
		Count  int64
	}

	query := config.DB.Model(&models.OperationLog{}).
		Select("action, count(*) as count").
		Where("created_at >= ?", sevenDaysAgo)

	// 权限控制：超级管理员查看所有统计，普通管理员只查看自己公司的统计
	if !isSuperAdmin.(bool) {
		if companyID != nil {
			query = query.Where("company_id = ?", companyID)
		} else {
			// 如果没有公司ID，只统计自己的操作
			query = query.Where("admin_id = ?", adminID)
		}
	}

	query.Group("action").Scan(&stats)

	c.JSON(http.StatusOK, stats)
}

func toInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
