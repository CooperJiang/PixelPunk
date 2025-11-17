package middleware

import (
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func InstallCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		installManager := common.GetInstallManager()

		path := c.Request.URL.Path
		if shouldSkipInstallCheck(path) {
			c.Next()
			return
		}

		if installManager.IsInstallMode() {
			c.JSON(200, gin.H{
				"code":    int(errors.CodeSystemNotInstalled), // 使用统一错误码：7001
				"message": "系统未安装，请先完成安装配置",
				"data": gin.H{
					"install_required": true,
					"setup_url":        "/setup",
					"status":           installManager.GetStatus(),
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func shouldSkipInstallCheck(path string) bool {
	skipPaths := []string{
		"/api/v1/setup",  // 安装相关API
		"/static",        // 静态资源
		"/resources",     // 资源文件
		"/file",          // 文件服务
		"/api/v1/health", // 健康检查
		"/favicon.ico",   // 网站图标
		"/robots.txt",    // 搜索引擎文件
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	if path == "/" || strings.HasPrefix(path, "/setup") {
		return true
	}

	return false
}
