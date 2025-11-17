package routes

import (
	oauthController "pixelpunk/internal/controllers/oauth"
	userController "pixelpunk/internal/controllers/user"

	"github.com/gin-gonic/gin"
)

/* RegisterAuthRoutes 注册公开的认证相关路由（不需要JWT认证） */
func RegisterAuthRoutes(r *gin.RouterGroup) {
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	r.POST("/send-registration-code", userController.SendRegistrationCode)
	r.POST("/send-reset-password-code", userController.SendResetPasswordCode)

	// 基于验证码的密码重置（旧方式）
	r.POST("/reset-password", userController.ResetPassword)

	// 基于token的密码重置（新方式）
	r.POST("/forgot-password", userController.ForgotPassword)
	r.POST("/verify-reset-token", userController.VerifyResetToken)
	r.POST("/reset-password-token", userController.ResetPasswordWithToken)

	oauthRoutes := r.Group("/oauth")
	{
		oauthRoutes.POST("/github/login", oauthController.GithubLogin)
		oauthRoutes.POST("/google/login", oauthController.GoogleLogin)
		oauthRoutes.POST("/linuxdo/login", oauthController.LinuxdoLogin)
	}
}
