package routes

import (
	userTagController "pixelpunk/internal/controllers/user_tag"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserTagRoutes(r *gin.RouterGroup) {
	// 用户标签管理（需要认证）
	tagGroup := r.Group("/user-tags")
	tagGroup.Use(middleware.RequireAuth())
	{
		tagGroup.GET("/list", userTagController.GetUserTagList)

		// 获取所有标签（不分页）
		tagGroup.GET("/all", userTagController.GetAllUserTags)

		tagGroup.POST("/create", userTagController.CreateUserTag)

		tagGroup.POST("/update", userTagController.UpdateUserTag)

		tagGroup.POST("/delete", userTagController.DeleteUserTag)

		tagGroup.POST("/batch-delete", userTagController.BatchDeleteUserTags)

		tagGroup.POST("/merge", userTagController.MergeUserTags)

		tagGroup.GET("/stats", userTagController.GetUserTagStats)
	}
}
