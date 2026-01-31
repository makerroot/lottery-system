package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查端点（支持 GET 和 HEAD）
func HealthCheck(c *gin.Context) {
	// HEAD 请求不返回 body（健康检查使用）
	if c.Request.Method == "HEAD" {
		c.Status(http.StatusOK)
		return
	}

	// GET 请求返回 JSON 信息
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "系统运行正常",
		"service": "lottery-backend",
		"version": "1.0.0",
	})
}
