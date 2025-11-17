package dto

type RegisterDTO struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Code     string `json:"code" binding:"required"`
}

func (r *RegisterDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Username.required": "用户名不能为空",
		"Username.min":      "用户名长度不能小于2个字符",
		"Username.max":      "用户名长度不能超过20个字符",
		"Email.required":    "邮箱不能为空",
		"Email.email":       "请输入有效的邮箱地址",
		"Password.required": "密码不能为空",
		"Password.min":      "密码长度不能小于6个字符",
		"Password.max":      "密码长度不能超过20个字符",
		"Code.required":     "验证码不能为空",
	}
}
