package routes

import (
	userCategoryController "pixelpunk/internal/controllers/user_category"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserCategoryRoutes(r *gin.RouterGroup) {
	// 用户分类管理（需要认证）
	categoryGroup := r.Group("/user-categories")
	categoryGroup.Use(middleware.RequireAuth())
	{
		categoryGroup.POST("/create", userCategoryController.CreateCategory)

		categoryGroup.GET("/list", userCategoryController.GetCategoryList)

		// 获取所有分类（不分页）
		categoryGroup.GET("/all", userCategoryController.GetAllCategories)

		categoryGroup.GET("/:id", userCategoryController.GetCategoryDetail)

		categoryGroup.POST("/update", userCategoryController.UpdateCategory)

		categoryGroup.POST("/delete", userCategoryController.DeleteCategory)

		categoryGroup.POST("/batch-delete", userCategoryController.BatchDelete)

		categoryGroup.POST("/update-status", userCategoryController.UpdateStatus)

		categoryGroup.POST("/batch-sort", userCategoryController.BatchUpdateSortOrder)
	}
}
