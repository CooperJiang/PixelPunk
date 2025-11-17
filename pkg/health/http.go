package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BasicHealthHandler 基础健康检查处理器
func BasicHealthHandler(c *gin.Context) {
	result := Check(CheckTypeBasic)

	statusCode := http.StatusOK
	if result.Status == StatusDown {
		statusCode = http.StatusServiceUnavailable
	} else if result.Status == StatusDegraded {
		statusCode = http.StatusOK // 降级状态仍返回200，但在响应中标记
	}

	c.JSON(statusCode, result)
}

// CompleteHealthHandler 完整健康检查处理器
func CompleteHealthHandler(c *gin.Context) {
	result := Check(CheckTypeComplete)

	statusCode := http.StatusOK
	if result.Status == StatusDown {
		statusCode = http.StatusServiceUnavailable
	} else if result.Status == StatusDegraded {
		statusCode = http.StatusOK // 降级状态仍返回200，但在响应中标记
	}

	c.JSON(statusCode, result)
}

// SimpleHealthHandler 简单健康检查处理器（兼容旧的健康检查）
func SimpleHealthHandler(c *gin.Context) {
	uptime := time.Since(startTime).Truncate(time.Second).String()

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"uptime":  uptime,
		"version": appVersion,
	})
}
