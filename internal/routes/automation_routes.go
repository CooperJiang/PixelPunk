package routes

import (
	"pixelpunk/internal/controllers/automation"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAutomationRoutes(r *gin.RouterGroup) {
	userAutomation := r.Group("/user/automation")
	userAutomation.Use(middleware.RequireAuth())
	{
		userAutomation.GET("/overview", automation.GetUserAutomationOverview)

		userAutomation.GET("/tagging/tasks", automation.GetUserTaggingTasks)

		userAutomation.GET("/vector/tasks", automation.GetUserVectorTasks)
	}
}
