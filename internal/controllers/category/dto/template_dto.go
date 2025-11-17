package dto

type CreateTemplateDTO struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	IsPopular   bool   `json:"is_popular"`
	SortOrder   int    `json:"sort_order"`
}

func (d *CreateTemplateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.required":   "模板名称不能为空",
		"Name.min":        "模板名称不能为空",
		"Name.max":        "模板名称不能超过50个字符",
		"Description.max": "模板描述不能超过500个字符",
		"Icon.max":        "图标不能超过50个字符",
	}
}

type UpdateTemplateDTO struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"omitempty,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	IsPopular   *bool  `json:"is_popular"`
	SortOrder   *int   `json:"sort_order"`
}

func (d *UpdateTemplateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required":     "模板ID不能为空",
		"Name.min":        "模板名称不能为空",
		"Name.max":        "模板名称不能超过50个字符",
		"Description.max": "模板描述不能超过500个字符",
		"Icon.max":        "图标不能超过50个字符",
	}
}

type DeleteTemplateDTO struct {
	ID uint `json:"id" binding:"required"`
}

func (d *DeleteTemplateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required": "模板ID不能为空",
	}
}

type TemplateListQueryDTO struct {
	Keyword   string `form:"keyword"`
	IsPopular *bool  `form:"is_popular"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Size      int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name sort_order usage_count created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

func (d *TemplateListQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":        "页码必须大于等于1",
		"Size.min":        "每页数量必须大于等于1",
		"Size.max":        "每页数量不能超过100",
		"SortBy.oneof":    "排序字段必须是name、sort_order、usage_count或created_at",
		"SortOrder.oneof": "排序方向必须是asc或desc",
	}
}

type BatchSortOrderDTO struct {
	SortOrders []SortOrderItem `json:"sort_orders" binding:"required,min=1"`
}

type SortOrderItem struct {
	ID        uint `json:"id" binding:"required"`
	SortOrder int  `json:"sort_order"`
}

func (d *BatchSortOrderDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SortOrders.required": "排序列表不能为空",
		"SortOrders.min":      "排序列表不能为空",
	}
}
