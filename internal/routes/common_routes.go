package routes

import (
	commonController "pixelpunk/internal/controllers/common"
	fileController "pixelpunk/internal/controllers/file"
	settingController "pixelpunk/internal/controllers/setting"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCommonRoutes(r *gin.RouterGroup) {
	r.GET("/health-check", commonController.HealthCheck)

	r.GET("/settings/global", settingController.GetGlobalSettings)
	r.GET("/settings/legal", settingController.GetLegalDocuments)
	r.GET("/settings/oauth", settingController.GetOAuthConfig)

	imageRoutes := r.Group("/files")
	{
		imageRoutes.GET("/tags", fileController.AdminGetTagList)
		imageRoutes.GET("/colors", fileController.AdminGetColorList)
	}

	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.RequireAuth())
	{
		userRoutes.GET("/tags", fileController.UserGetTagList)
		userRoutes.GET("/colors", fileController.UserGetColorList)
		userRoutes.GET("/categories", fileController.UserGetCategoryList)
	}

	guestRoutes := r.Group("/guest")
	{
		guestRoutes.GET("/tags", fileController.GuestGetTagList)
		guestRoutes.GET("/colors", fileController.GuestGetColorList)
	}

	publicRoutes := r.Group("/public")
	{
		publicRoutes.GET("/tags", fileController.GetPublicTagList)
	}

	pbRoutes := r.Group("/pb")
	{
		pbRoutes.GET("/data", fileController.GetPublicFileCount)
	}
}
