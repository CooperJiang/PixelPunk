package routes

import (
	randomAPIController "pixelpunk/internal/controllers/random_api"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRandomAPIRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RequireAuth())

	r.POST("/create", randomAPIController.CreateRandomAPI)
	r.GET("/list", randomAPIController.GetRandomAPIList)
	r.PUT("/:id/status", randomAPIController.UpdateRandomAPIStatus)
	r.PUT("/:id/config", randomAPIController.UpdateRandomAPIConfig)
	r.DELETE("/:id", randomAPIController.DeleteRandomAPI)
}
