package routes

import (
	messageController "pixelpunk/internal/controllers/message"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(r *gin.RouterGroup) {
	userMessageGroup := r.Group("/messages")
	userMessageGroup.Use(middleware.RequireAuth())
	{
		userMessageGroup.GET("", messageController.GetUserMessages)

		userMessageGroup.POST("/:id/read", messageController.MarkMessageRead)

		userMessageGroup.POST("/read-all", messageController.MarkAllMessagesRead)

		userMessageGroup.DELETE("/:id", messageController.DeleteMessage)

		userMessageGroup.GET("/unread-count", messageController.GetUnreadCount)
	}

	adminMessageGroup := r.Group("/admin/messages")
	adminMessageGroup.Use(middleware.RequireAuth())
	adminMessageGroup.Use(middleware.RequireAdmin())
	{
		adminMessageGroup.GET("/templates", messageController.GetAllTemplates)

		adminMessageGroup.GET("/templates/:id", messageController.GetTemplateByID)

		adminMessageGroup.POST("/templates", messageController.CreateTemplate)

		adminMessageGroup.PUT("/templates/:id", messageController.UpdateTemplate)

		adminMessageGroup.DELETE("/templates/:id", messageController.DeleteTemplate)

		adminMessageGroup.POST("/templates/:id/toggle", messageController.ToggleTemplate)

		adminMessageGroup.POST("/send", messageController.SendMessage)

		adminMessageGroup.GET("/stats", messageController.GetMessageStats)
	}
}
