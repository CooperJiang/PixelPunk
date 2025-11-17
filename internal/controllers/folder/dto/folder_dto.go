package dto

type CreateFolderDTO struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	ParentID    string `json:"parent_id"`
	Permission  string `json:"permission" binding:"omitempty,oneof=private public"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

func (d *CreateFolderDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.required":    "文件夹名称不能为空",
		"Name.min":         "文件夹名称不能为空",
		"Name.max":         "文件夹名称不能超过100个字符",
		"Permission.oneof": "权限必须是 private 或 public",
		"Description.max":  "描述不能超过500个字符",
	}
}

type UpdateFolderDTO struct {
	FolderID    string `json:"folderId" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
	ParentID    string `json:"parent_id"`
	Permission  string `json:"permission" binding:"omitempty,oneof=private public"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

func (d *UpdateFolderDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FolderID.required": "文件夹ID不能为空",
		"Name.required":     "文件夹名称不能为空",
		"Name.min":          "文件夹名称不能为空",
		"Name.max":          "文件夹名称不能超过100个字符",
		"Permission.oneof":  "权限必须是 private 或 public",
		"Description.max":   "描述不能超过500个字符",
	}
}

type ListFoldersQueryDTO struct {
	ParentID string `form:"parent_id"`
}

func (d *ListFoldersQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{}
}

type MoveFoldersRequest struct {
	FolderIDs   []string `json:"folder_ids" binding:"required,min=1,max=100"` // 要移动的文件夹ID列表
	NewParentID string   `json:"new_parent_id"`                               // 新的父文件夹ID，空字符串表示移动到根目录
}

func (d *MoveFoldersRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"FolderIDs.required": "文件夹ID列表不能为空",
		"FolderIDs.min":      "至少需要选择一个文件夹",
		"FolderIDs.max":      "单次最多只能移动100个文件夹",
	}
}

type ReorderFoldersRequest struct {
	ParentID  string   `json:"parent_id"`                           // 父文件夹ID，空字符串表示根目录
	FolderIDs []string `json:"folder_ids" binding:"required,min=1"` // 按新顺序排列的文件夹ID数组
}

func (d *ReorderFoldersRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"FolderIDs.required": "文件夹ID列表不能为空",
		"FolderIDs.min":      "至少需要一个文件夹",
	}
}

type FolderContentsDTO struct {
	ParentID    string `form:"parent_id"`                                                     // 文件夹ID
	Keyword     string `form:"keyword"`                                                       // 搜索关键词
	AccessLevel string `form:"access_level"`                                                  // 访问级别筛选
	SortBy      string `form:"sort_by" binding:"omitempty,oneof=name created_at size custom"` // 排序字段
	SortOrder   string `form:"sort_order" binding:"omitempty,oneof=asc desc"`                 // 排序方向
}

func (d *FolderContentsDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":        "页码必须大于等于1",
		"Size.min":        "每页数量必须大于等于1",
		"Size.max":        "每页数量不能超过100",
		"SortBy.oneof":    "排序字段必须是name、created_at、size或custom",
		"SortOrder.oneof": "排序方向必须是asc或desc",
	}
}

type ReorderFoldersDTO struct {
	ParentID  string   `json:"parent_id"`                           // 父文件夹ID，空字符串表示根目录
	FolderIDs []string `json:"folder_ids" binding:"required,min=1"` // 按新顺序排列的文件夹ID数组
}

func (d *ReorderFoldersDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FolderIDs.required": "文件夹ID列表不能为空",
		"FolderIDs.min":      "至少需要一个文件夹ID",
	}
}

type FolderPathChainDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id,omitempty"`
	Level    int    `json:"level"`
}

type FolderPathChainResponseDTO struct {
	FolderID    string               `json:"folder_id"`
	FullPath    string               `json:"full_path"`
	PathChain   []FolderPathChainDTO `json:"path_chain"`
	TotalLevels int                  `json:"total_levels"`
}

type BatchFolderPathChainsRequestDTO struct {
	FolderIDs []string `json:"folder_ids" binding:"required"`
}

func (d *BatchFolderPathChainsRequestDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FolderIDs.required": "文件夹ID列表不能为空",
	}
}
