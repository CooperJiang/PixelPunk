package dto

type CreateChannelDTO struct {
	Name      string                 `json:"name" binding:"required,max=50"`
	Type      string                 `json:"type" binding:"required,max=20"`
	IsDefault bool                   `json:"is_default"`
	Status    *int8                  `json:"status"`
	Remark    string                 `json:"remark" binding:"max=255"`
	Configs   map[string]interface{} `json:"configs"` // 配置项，一次性提交渠道信息和配置值
}

func (d *CreateChannelDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.required": "渠道名称不能为空",
		"Name.max":      "渠道名称不能超过50个字符",
		"Type.required": "存储类型不能为空",
		"Type.max":      "存储类型不能超过20个字符",
		"Remark.max":    "备注不能超过255个字符",
	}
}

type UpdateChannelDTO struct {
	Name      string                 `json:"name" binding:"omitempty,max=50"`
	IsDefault *bool                  `json:"is_default"`
	Status    *int8                  `json:"status"`
	Remark    string                 `json:"remark" binding:"max=255"`
	Configs   map[string]interface{} `json:"configs"` // 配置项，可一并更新配置值
}

func (d *UpdateChannelDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.max":   "渠道名称不能超过50个字符",
		"Remark.max": "备注不能超过255个字符",
	}
}

type ChannelConfigDTO map[string]interface{}

func (d *ChannelConfigDTO) GetValidationMessages() map[string]string {
	return map[string]string{}
}
