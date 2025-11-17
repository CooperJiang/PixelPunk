package routes

import (
	vectorController "pixelpunk/internal/controllers/vector"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterVectorRoutes(r *gin.RouterGroup) {
	vectorGroup := r.Group("/vector")
	vectorGroup.Use(middleware.RequireAuth())
	vectorGroup.Use(middleware.RequireAdmin()) // 管理员权限
	{
		vectorGroup.GET("/list", vectorController.GetVectorList)              // 获取向量列表
		vectorGroup.GET("/stats", vectorController.GetVectorStats)            // 获取向量统计
		vectorGroup.GET("/qdrant-stats", vectorController.GetQdrantRealStats) // 获取 Qdrant 实际统计
		vectorGroup.GET("/models", vectorController.GetAvailableModels)       // 获取可用模型

		vectorGroup.GET("/detail/:fileId", vectorController.GetVectorDetail)

		vectorGroup.POST("/batch", vectorController.BatchVectorAction)
		vectorGroup.POST("/batch-action", vectorController.BatchVectorAction)
		vectorGroup.POST("/retry", vectorController.RetryVector)                      // 重试单个向量
		vectorGroup.POST("/regenerate-all", vectorController.RegenerateAllVectors)    // 重新生成所有向量
		vectorGroup.POST("/retry-all-failed", vectorController.RetryAllFailedVectors) // 重试所有失败的向量任务
		vectorGroup.POST("/recover-stuck", vectorController.RecoverStuckTasks)        // 恢复卡住的任务
		vectorGroup.GET("/auto-processing", vectorController.GetAutoProcessing)
		vectorGroup.POST("/auto-processing", vectorController.SetAutoProcessing)
		vectorGroup.POST("/concurrency", vectorController.SetConcurrency)

		vectorGroup.POST("/reconcile/missing", vectorController.ReconcileMissing)
		vectorGroup.POST("/reconcile/orphans", vectorController.CleanOrphans)
		vectorGroup.POST("/rebuild/stale", vectorController.RebuildStale)

		vectorGroup.GET("/logs", vectorController.GetVectorLogs) // 获取处理日志
	}
}
