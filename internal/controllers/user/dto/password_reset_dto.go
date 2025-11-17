package dto

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// VerifyResetTokenRequest 验证重置token请求
type VerifyResetTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// ForgotPasswordResponse 忘记密码响应
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

// ResetPasswordResponse 重置密码响应
type ResetPasswordResponse struct {
	Message string `json:"message"`
}

// VerifyResetTokenResponse 验证重置token响应
type VerifyResetTokenResponse struct {
	Valid bool   `json:"valid"`
	Email string `json:"email,omitempty"`
}
