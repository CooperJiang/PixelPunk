package routes

import (
	"pixelpunk/internal/controllers/setup"

	"github.com/gin-gonic/gin"
)

func RegisterSetupRoutes(r *gin.RouterGroup) {
	setupController := &setup.SetupController{}

	setupGroup := r.Group("/setup")
	{
		setupGroup.GET("/status", setupController.GetStatus)                // 获取安装状态
		setupGroup.POST("/test-connection", setupController.TestConnection) // 测试数据库连接
		setupGroup.POST("/install", setupController.Install)                // 执行安装
	}
}
