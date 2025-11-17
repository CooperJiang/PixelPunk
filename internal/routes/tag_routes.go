package routes

import (
	tagController "pixelpunk/internal/controllers/tag"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTagRoutes(r *gin.RouterGroup) {
	{
		r.GET("/list", tagController.ListTags)

		r.GET("/search", tagController.SearchTags)

		r.GET("/stats", tagController.GetTagStats)
	}

	authRoute := r.Group("")
	authRoute.Use(middleware.RequireAuth())
	{
		authRoute.POST("/files/:file_id/tags", tagController.AddTagsToFile)
		authRoute.DELETE("/files/:file_id/tags", tagController.RemoveTagsFromFile)
		authRoute.GET("/files/:file_id/tags", tagController.GetFileTags)
	}

	adminRoute := r.Group("/admin")
	adminRoute.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	{
		adminRoute.POST("/create", tagController.CreateTag)
		adminRoute.POST("/update", tagController.UpdateTag)
		adminRoute.DELETE("/:tag_id", tagController.DeleteTag)

		adminRoute.POST("/batch", tagController.BatchOperateTags)

		adminRoute.GET("/stats/detailed", tagController.GetDetailedTagStats)
		adminRoute.GET("/analytics", tagController.GetTagAnalytics)
	}
}
