package routes

import (
	activityController "pixelpunk/internal/controllers/activity"
	userController "pixelpunk/internal/controllers/user"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

/* RegisterPublicUserRoutes 注册公开的用户路由（不需要认证，兼容旧API） */
func RegisterPublicUserRoutes(r *gin.RouterGroup) {
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.POST("/send-registration-code", userController.SendRegistrationCode)
	r.POST("/send-reset-password-code", userController.SendResetPasswordCode)
	r.POST("/reset-password", userController.ResetPassword)
}

/* RegisterUserRoutes 注册用户相关路由（需要认证） */
func RegisterUserRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/personal")
	userGroup.Use(middleware.RequireAuth())
	{
		userGroup.POST("/update-password", userController.UpdatePassword)

		userGroup.GET("/profile", userController.GetProfile)

		userGroup.POST("/profile", userController.UpdateProfile)

		userGroup.GET("/folders", userController.GetUserFolders)

		userGroup.POST("/send-change-email-code", userController.SendChangeEmailCode)

		userGroup.POST("/change-email", userController.ChangeEmail)

		userGroup.GET("/access-control", userController.GetUserAccessControl)
		userGroup.POST("/access-control/createOrUpdate", userController.CreateOrUpdateUserAccessControl)
		userGroup.POST("/access-control/reset", userController.ResetUserAccessControl)

		userGroup.GET("/workspace/stats", userController.GetWorkspaceStats)

		userGroup.GET("/activities", activityController.GetUserActivities)
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.RequireAdmin())
	{
	}

	superGroup := r.Group("/super")
	superGroup.Use(middleware.RequireSuperAdmin())
	{
	}
}
