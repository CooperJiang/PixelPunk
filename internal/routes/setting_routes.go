package routes

import (
	settingController "pixelpunk/internal/controllers/setting"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSettingRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RequireAuth())
	r.Use(middleware.RequireAdmin())
	{
		r.GET("", settingController.GetSettings)

		r.GET("/:key", settingController.GetSetting)

		r.GET("/group/:group/map", settingController.GetSettingsByGroupAsMap)

		r.POST("", settingController.CreateSetting)

		r.POST("/batch", settingController.BatchCreateSettings)

		r.POST("/upsert", settingController.BatchUpsertSettings)

		r.PUT("", settingController.UpdateSetting)

		r.PUT("/batch", settingController.BatchUpdateSettings)

		r.DELETE("/:key", settingController.DeleteSetting)

		r.POST("/mail/test", settingController.TestEmailSettings)
		r.POST("/mail/refresh", settingController.RefreshEmailSettings)

		r.POST("/test-mail", settingController.TestMail)

		r.POST("/vector/test", settingController.TestVectorSettings)
		r.POST("/vector/test-qdrant", settingController.TestQdrantConnection)

		r.POST("/test-proxy", settingController.TestProxy)
	}
}
