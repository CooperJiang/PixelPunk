package routes

import (
	apikeyController "pixelpunk/internal/controllers/apikey"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAPIKeyRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RequireAuth())

	r.POST("/create", apikeyController.CreateAPIKey)

	r.GET("/list", apikeyController.GetAPIKeyList)

	r.GET("/:key_id", apikeyController.GetAPIKeyDetail)

	r.PUT("/:key_id", apikeyController.UpdateAPIKey)

	r.DELETE("/:key_id", apikeyController.DeleteAPIKey)

	r.POST("/:key_id/toggle", apikeyController.ToggleAPIKeyStatus)

	r.GET("/:key_id/stats", apikeyController.GetAPIKeyStats)

	r.POST("/:key_id/regenerate", apikeyController.RegenerateAPIKey)
}
