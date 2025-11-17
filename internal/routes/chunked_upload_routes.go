package routes

import (
	fileController "pixelpunk/internal/controllers/file"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterChunkedUploadRoutes(r *gin.RouterGroup) {
	chunked := r.Group("/chunked")
	chunked.Use(middleware.RequireAuth()) // 需要认证
	{
		chunked.POST("/init", fileController.InitChunkedUpload)

		chunked.POST("/upload", fileController.UploadChunk)

		chunked.POST("/complete", fileController.CompleteChunkedUpload)

		chunked.GET("/status", fileController.GetChunkedUploadStatus)

		chunked.DELETE("/cancel", fileController.CancelChunkedUpload)
	}

	adminGroup := chunked.Group("/admin")
	adminGroup.Use(middleware.RequireAdmin())
	{
		adminGroup.POST("/cleanup", fileController.ManualCleanupChunkedUploads)

		adminGroup.GET("/stats", fileController.GetChunkedUploadStats)
	}
}
