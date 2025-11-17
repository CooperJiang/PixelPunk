package user

import (
	"pixelpunk/internal/controllers/user/dto"
	userService "pixelpunk/internal/services/user"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ForgotPassword 发送密码重置邮件
// @Summary 发送密码重置邮件
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body dto.ForgotPasswordRequest true "请求参数"
// @Success 200 {object} dto.ForgotPasswordResponse
// @Router /api/auth/forgot-password [post]
func ForgotPassword(c *gin.Context) {
	req, err := common.ValidateRequest[dto.ForgotPasswordRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	ipAddress := c.ClientIP()

	if err := userService.SendPasswordResetEmail(req.Email, ipAddress); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "密码重置邮件已发送，请查收邮箱")
}

// VerifyResetToken 验证重置token
// @Summary 验证密码重置token
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body dto.VerifyResetTokenRequest true "请求参数"
// @Success 200 {object} dto.VerifyResetTokenResponse
// @Router /api/auth/verify-reset-token [post]
func VerifyResetToken(c *gin.Context) {
	req, err := common.ValidateRequest[dto.VerifyResetTokenRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	user, err := userService.VerifyResetToken(req.Token)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	data := gin.H{
		"valid": true,
		"email": user.Email,
	}

	errors.ResponseSuccess(c, data, "令牌验证成功")
}

// ResetPasswordWithToken 使用token重置密码
// @Summary 使用token重置密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body dto.ResetPasswordRequest true "请求参数"
// @Success 200 {object} dto.ResetPasswordResponse
// @Router /api/auth/reset-password-token [post]
func ResetPasswordWithToken(c *gin.Context) {
	req, err := common.ValidateRequest[dto.ResetPasswordRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := userService.ResetPasswordWithToken(req.Token, req.NewPassword); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "密码重置成功，请使用新密码登录")
}
