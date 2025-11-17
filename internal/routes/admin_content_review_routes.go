package routes

import (
	adminController "pixelpunk/internal/controllers/admin"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminContentReviewRoutes(r *gin.RouterGroup) {
	reviewGroup := r.Group("/content-review")
	reviewGroup.Use(middleware.RequireAuth())
	reviewGroup.Use(middleware.RequireAdmin())
	{
		reviewGroup.GET("/queue", adminController.GetReviewQueue)

		reviewGroup.GET("/logs", adminController.GetReviewLogs)

		reviewGroup.GET("/stats", adminController.GetReviewStats)

		reviewGroup.GET("/files/:fileId", adminController.GetFileDetail)

		reviewGroup.POST("/review", adminController.ReviewFile)

		reviewGroup.POST("/batch-review", adminController.BatchReview)

		reviewGroup.DELETE("/files/:fileId/hard-delete", adminController.HardDeleteReviewedFile)
	}
}
