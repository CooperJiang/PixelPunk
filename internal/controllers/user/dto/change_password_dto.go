package dto

type ChangePasswordDTO struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20"`
}

func (d *ChangePasswordDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"OldPassword.required": "原密码不能为空",
		"NewPassword.required": "新密码不能为空",
		"NewPassword.min":      "新密码长度不能小于6位",
		"NewPassword.max":      "新密码长度不能大于20位",
	}
}
