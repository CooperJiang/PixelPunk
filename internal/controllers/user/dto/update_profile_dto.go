package dto

type UpdateProfileDTO struct {
	Username string `json:"username" binding:"omitempty,min=2,max=20"`
	Avatar   string `json:"avatar" binding:"omitempty"`
	Bio      string `json:"bio" binding:"omitempty,max=500"`
	Website  string `json:"website" binding:"omitempty,max=255"`
}

func (d *UpdateProfileDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Username.min": "用户名长度不能小于2个字符",
		"Username.max": "用户名长度不能超过20个字符",
		"Bio.max":      "个人简介不能超过500个字符",
		"Website.max":  "个人网站地址不能超过255个字符",
	}
}
