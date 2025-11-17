package routes

import (
	storageController "pixelpunk/internal/controllers/storage"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterStorageRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RequireAdmin())

	r.GET("/list", storageController.ListChannels)

	r.GET("/supported-types", storageController.ListSupportedTypes)

	r.GET("/config-templates", storageController.GetConfigTemplates)
	r.GET("/config-templates/", storageController.GetConfigTemplates)
	r.GET("/config-templates/:type", storageController.GetConfigTemplates)

	r.GET("/:id", storageController.GetChannel)

	r.POST("/", storageController.CreateChannel)
	r.POST("/create", storageController.CreateChannel)

	r.PUT("/:id", storageController.UpdateChannel)

	r.DELETE("/:id", storageController.DeleteChannel)

	r.GET("/:id/configs", storageController.GetChannelConfigs)

	r.PUT("/:id/configs", storageController.UpdateChannelConfigs)

	r.POST("/:id/test", storageController.TestConnection)

	r.POST("/:id/default", storageController.SetDefaultChannel)

	r.POST("/:id/enable", storageController.EnableChannel)

	r.POST("/:id/disable", storageController.DisableChannel)

	r.GET("/:id/export", storageController.ExportChannelConfig)

	r.GET("/export/all", storageController.ExportAllChannelConfigs)

	r.POST("/import", storageController.ImportChannelConfig)

	r.POST("/:id/refresh-cache", storageController.RefreshChannelCache)

	r.POST("/clear-cache", storageController.ClearAllChannelCache)
}
