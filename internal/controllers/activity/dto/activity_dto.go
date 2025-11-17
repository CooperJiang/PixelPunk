package dto

type GetUserActivitiesDTO struct {
	Page int `form:"page" binding:"omitempty,min=1" json:"page"`         // 页码，从1开始，可选
	Size int `form:"size" binding:"omitempty,min=1,max=100" json:"size"` // 每页数量，最大100，可选
}

func (d *GetUserActivitiesDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min": "页码必须大于0",
		"Size.min": "每页数量必须大于0",
		"Size.max": "每页数量不能超过100",
	}
}

func (d *GetUserActivitiesDTO) SetDefaults() {
	if d.Page <= 0 {
		d.Page = 1
	}
	if d.Size <= 0 {
		d.Size = 20
	}
	if d.Size > 100 {
		d.Size = 100
	}
}
