package dto

type GetUserTagListQueryDTO struct {
	Keyword   string `form:"keyword"`                                                      // 关键字搜索
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name file_count created_at"` // 排序字段
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`                // 排序方向
	Page      int    `form:"page" binding:"omitempty,min=1"`                               // 页码
	Size      int    `form:"size" binding:"omitempty,min=1,max=10000"`                     // 每页数量
}

func (d *GetUserTagListQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SortBy.oneof":    "排序字段必须是name、file_count或created_at",
		"SortOrder.oneof": "排序方向必须是asc或desc",
		"Page.min":        "页码必须大于等于1",
		"Size.min":        "每页数量必须大于等于1",
		"Size.max":        "每页数量不能超过10000",
	}
}

type CreateUserTagDTO struct {
	Name string `json:"name" binding:"required,min=1,max=50"` // 标签名称
}

func (d *CreateUserTagDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.required": "标签名称不能为空",
		"Name.min":      "标签名称不能为空",
		"Name.max":      "标签名称不能超过50个字符",
	}
}

type UpdateUserTagDTO struct {
	ID   uint   `json:"id" binding:"required"`                // 标签ID
	Name string `json:"name" binding:"required,min=1,max=50"` // 标签名称
}

func (d *UpdateUserTagDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required":   "标签ID不能为空",
		"Name.required": "标签名称不能为空",
		"Name.min":      "标签名称不能为空",
		"Name.max":      "标签名称不能超过50个字符",
	}
}

type DeleteUserTagDTO struct {
	ID uint `json:"id" binding:"required"` // 标签ID
}

func (d *DeleteUserTagDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required": "标签ID不能为空",
	}
}

type BatchDeleteUserTagsDTO struct {
	IDs []uint `json:"ids" binding:"required,min=1,dive,gt=0"` // 标签ID列表
}

func (d *BatchDeleteUserTagsDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"IDs.required": "标签ID列表不能为空",
		"IDs.min":      "至少需要选择一个标签",
		"IDs.dive":     "标签ID格式错误",
	}
}

type MergeUserTagsDTO struct {
	SourceIDs []uint `json:"source_ids" binding:"required,min=1,dive,gt=0"` // 源标签ID列表
	TargetID  uint   `json:"target_id" binding:"required,gt=0"`             // 目标标签ID
}

func (d *MergeUserTagsDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SourceIDs.required": "源标签ID列表不能为空",
		"SourceIDs.min":      "至少需要选择一个源标签",
		"SourceIDs.dive":     "源标签ID格式错误",
		"TargetID.required":  "目标标签ID不能为空",
		"TargetID.gt":        "目标标签ID必须大于0",
	}
}
