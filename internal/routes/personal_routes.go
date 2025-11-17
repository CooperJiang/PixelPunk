package routes

import (
	activityController "pixelpunk/internal/controllers/activity"

	"github.com/gin-gonic/gin"
)

func RegisterPersonalRoutes(r *gin.RouterGroup) {
	r.GET("/activities", activityController.GetUserActivities)
}
