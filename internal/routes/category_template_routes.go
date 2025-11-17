package routes

import (
	categoryController "pixelpunk/internal/controllers/category"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryTemplateRoutes(r *gin.RouterGroup) {
	templateController := categoryController.NewTemplateController()

	adminGroup := r.Group("/admin/category-templates")
	adminGroup.Use(middleware.RequireAuth())
	adminGroup.Use(middleware.RequireAdmin())
	{
		adminGroup.POST("/create", templateController.CreateTemplate)

		adminGroup.GET("/detail/:id", templateController.GetTemplate)

		adminGroup.POST("/update", templateController.UpdateTemplate)

		adminGroup.POST("/delete", templateController.DeleteTemplate)

		adminGroup.GET("/list", templateController.ListTemplates)

		adminGroup.POST("/batch-sort", templateController.BatchUpdateSortOrder)
	}

	publicGroup := r.Group("/category-templates")
	{
		publicGroup.GET("/popular", templateController.GetPopularTemplates)

		publicGroup.GET("/all", templateController.GetAllTemplates)
	}
}
