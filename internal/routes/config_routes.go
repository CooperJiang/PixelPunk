package routes

import (
	configController "pixelpunk/internal/controllers/config"

	"github.com/gin-gonic/gin"
)

func RegisterConfigRoutes(r *gin.RouterGroup) {
	configGroup := r.Group("/config")
	{
		configGroup.GET("/upload", configController.GetUploadConfig)
		configGroup.GET("/upload/capabilities", configController.GetUploadCapabilities)
	}
}
