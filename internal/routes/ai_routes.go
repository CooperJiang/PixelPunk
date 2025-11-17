package routes

import (
	aiController "pixelpunk/internal/controllers/ai"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAIRoutes(r *gin.RouterGroup) {
	adminGroup := r.Group("")
	adminGroup.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	{
		taggingGroup := adminGroup.Group("/tagging")
		{
			taggingGroup.GET("/status", aiController.GetTaggingStatus)
			taggingGroup.GET("/stats", aiController.GetTaggingStats)
			taggingGroup.GET("/logs", aiController.GetTaggingLogs)           // 获取打标日志(分页)
			taggingGroup.GET("/queue-stats", aiController.GetQueueStats)     // 新增队列状态API
			taggingGroup.GET("/diagnosis", aiController.GetTaggingDiagnosis) // 诊断接口
			taggingGroup.GET("/auto-processing", aiController.GetAutoProcessingStatus)
			taggingGroup.POST("/concurrency", aiController.SetConcurrency)
			taggingGroup.POST("/retry-failed-all", aiController.RetryFailedAll)
			taggingGroup.POST("/reset-stuck", aiController.ResetStuckTasks)
			taggingGroup.POST("/auto-processing", aiController.ToggleAutoProcessing)
			taggingGroup.POST("/ignore", aiController.IgnoreTagging)     // 忽略打标
			taggingGroup.POST("/unignore", aiController.UnignoreTagging) // 取消忽略
			taggingGroup.POST("/retry", aiController.RetryTagging)
			taggingGroup.POST("/trigger", aiController.TriggerTagging)
		}
	}
}
