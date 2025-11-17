package routes

import (
	searchController "pixelpunk/internal/controllers/search"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSearchRoutes(r *gin.RouterGroup) {
	searchGroup := r.Group("/search")
	{
		vectorGroup := searchGroup.Group("/vector")
		{
			vectorGroup.POST("/search", searchController.VectorSearch)

			vectorGroup.GET("/stats", searchController.GetVectorStats)

			vectorGroup.GET("/health", searchController.GetVectorHealth)

			vectorGroup.GET("/similar/:fileId", searchController.SearchSimilarFiles)

			vectorGroup.POST("/regenerate", searchController.RegenerateVectors)

			vectorGroup.GET("/queue/status", searchController.GetVectorQueueStatus)

			vectorGroup.GET("/task/:task_id", searchController.GetVectorTaskStatus)
		}

		userGroup := searchGroup.Group("/user")
		userGroup.Use(middleware.RequireAuth())
		{
			userGroup.POST("/vector/search", searchController.UserVectorSearch)
		}

		galleryGroup := searchGroup.Group("/gallery")
		{
			galleryGroup.POST("/vector/search", searchController.GalleryVectorSearch)

			galleryGroup.GET("/similar/:fileId", searchController.GallerySimilarFiles)
		}

		userSimilarGroup := searchGroup.Group("/user")
		userSimilarGroup.Use(middleware.RequireAuth())
		{
			userSimilarGroup.GET("/similar/:fileId", searchController.UserSimilarFiles)
		}

		adminSimilarGroup := searchGroup.Group("/admin")
		adminSimilarGroup.Use(middleware.RequireAuth())
		adminSimilarGroup.Use(middleware.RequireAdmin())
		{
			adminSimilarGroup.GET("/similar/:fileId", searchController.AdminSimilarFiles)
		}
	}
}
