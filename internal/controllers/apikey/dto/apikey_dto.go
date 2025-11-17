package dto

type CreateAPIKeyDTO struct {
	Name             string   `json:"name" binding:"required,min=1,max=100"`
	StorageLimit     int64    `json:"storage_limit" binding:"omitempty,min=0"`
	SingleFileLimit  int64    `json:"single_file_limit" binding:"omitempty,min=0"`
	UploadCountLimit int      `json:"upload_count_limit" binding:"omitempty,min=0"`
	AllowedTypes     []string `json:"allowed_types" binding:"omitempty"`
	FolderID         string   `json:"folder_id" binding:"omitempty"`
	ExpiresInDays    int      `json:"expires_in_days" binding:"omitempty,min=0"`
}

func (d *CreateAPIKeyDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.required":        "密钥名称不能为空",
		"Name.min":             "密钥名称长度不能小于1个字符",
		"Name.max":             "密钥名称长度不能超过100个字符",
		"StorageLimit.min":     "存储容量限制不能为负数",
		"SingleFileLimit.min":  "单文件大小限制不能为负数",
		"UploadCountLimit.min": "上传次数限制不能为负数",
		"ExpiresInDays.min":    "有效天数不能为负数",
	}
}

type UpdateAPIKeyDTO struct {
	Name             string   `json:"name" binding:"omitempty,min=1,max=100"`
	StorageLimit     int64    `json:"storage_limit" binding:"omitempty,min=0"`
	SingleFileLimit  int64    `json:"single_file_limit" binding:"omitempty,min=0"`
	UploadCountLimit int      `json:"upload_count_limit" binding:"omitempty,min=0"`
	AllowedTypes     []string `json:"allowed_types" binding:"omitempty"`
	FolderID         string   `json:"folder_id" binding:"omitempty"`
	ExpiresInDays    int      `json:"expires_in_days" binding:"omitempty,min=0"`
	Status           int      `json:"status" binding:"omitempty,oneof=1 2"`
}

func (d *UpdateAPIKeyDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.min":             "密钥名称长度不能小于1个字符",
		"Name.max":             "密钥名称长度不能超过100个字符",
		"StorageLimit.min":     "存储容量限制不能为负数",
		"SingleFileLimit.min":  "单文件大小限制不能为负数",
		"UploadCountLimit.min": "上传次数限制不能为负数",
		"ExpiresInDays.min":    "有效天数不能为负数",
		"Status.oneof":         "状态值无效，应为1(启用)或2(禁用)",
	}
}

type APIKeyQueryDTO struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=100"`
	Status int    `form:"status" binding:"omitempty,oneof=1 2"`
	Search string `form:"search" binding:"omitempty,max=100"`
}

func (d *APIKeyQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":     "页码必须大于等于1",
		"Size.min":     "每页数量必须大于等于1",
		"Size.max":     "每页数量必须小于等于100",
		"Status.oneof": "状态值无效，应为1(启用)或2(禁用)",
		"Search.max":   "搜索关键字不能超过100个字符",
	}
}
