package dto

type CreateCategoryDTO struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=500"`
	SortOrder   int    `json:"sort_order"`
}

func (d *CreateCategoryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.required":   "分类名称不能为空",
		"Name.min":        "分类名称至少1个字符",
		"Name.max":        "分类名称不能超过50个字符",
		"Description.max": "分类描述不能超过500个字符",
	}
}

type UpdateCategoryDTO struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=500"`
	SortOrder   int    `json:"sort_order"`
}

func (d *UpdateCategoryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required":     "分类ID不能为空",
		"Name.required":   "分类名称不能为空",
		"Name.min":        "分类名称至少1个字符",
		"Name.max":        "分类名称不能超过50个字符",
		"Description.max": "分类描述不能超过500个字符",
	}
}

type DeleteCategoryDTO struct {
	ID uint `json:"id" binding:"required"`
}

func (d *DeleteCategoryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required": "分类ID不能为空",
	}
}

type BatchDeleteDTO struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

func (d *BatchDeleteDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"IDs.required": "分类ID列表不能为空",
		"IDs.min":      "至少选择一个分类",
	}
}

type UpdateStatusDTO struct {
	ID     uint   `json:"id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=active archived"`
}

func (d *UpdateStatusDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required":     "分类ID不能为空",
		"Status.required": "分类状态不能为空",
		"Status.oneof":    "分类状态必须是active或archived",
	}
}

type BatchSortOrderDTO struct {
	SortOrders []struct {
		ID        uint `json:"id" binding:"required"`
		SortOrder int  `json:"sort_order"`
	} `json:"sort_orders" binding:"required,min=1,dive"`
}

func (d *BatchSortOrderDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SortOrders.required": "排序列表不能为空",
		"SortOrders.min":      "至少包含一个排序项",
	}
}
