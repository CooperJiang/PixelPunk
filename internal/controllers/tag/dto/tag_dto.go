package dto

// ListTagsQueryDTO 标签列表查询DTO（版本）
type ListTagsQueryDTO struct {
	Keyword            string `form:"keyword"`                                                            // 关键字搜索
	UserID             *uint  `form:"user_id"`                                                            // 用户ID过滤
	AccessLevel        string `form:"access_level" binding:"omitempty,oneof=public private protected"`    // 访问级别过滤
	IsSystem           *bool  `form:"is_system"`                                                          // 是否系统标签
	IncludeUserStats   bool   `form:"include_user_stats"`                                                 // 包含用户统计
	IncludeAccessStats bool   `form:"include_access_stats"`                                               // 包含访问级别统计
	Page               int    `form:"page" binding:"omitempty,min=1"`                                     // 页码
	Size               int    `form:"size" binding:"omitempty,min=1"`                                     // 每页数量
	SortBy             string `form:"sort_by" binding:"omitempty,oneof=name created_at updated_at count"` // 排序字段
	SortOrder          string `form:"sort_order" binding:"omitempty,oneof=asc desc"`                      // 排序方向
}

func (d *ListTagsQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"AccessLevel.oneof": "访问级别必须是public、private或protected",
		"Page.min":          "页码必须大于等于1",
		"Size.min":          "每页数量必须大于等于1",
		"SortBy.oneof":      "排序字段必须是name、created_at或updated_at",
		"SortOrder.oneof":   "排序方向必须是asc或desc",
	}
}

// AddTagsToFileDTO 为文件添加标签DTO
type AddTagsToFileDTO struct {
	FileID      string   `json:"file_id" binding:"required,len=32"`                               // 文件ID
	TagNames    []string `json:"tag_names" binding:"required,min=1,dive,required,min=1,max=50"`   // 标签名称列表
	AccessLevel string   `json:"access_level" binding:"omitempty,oneof=public private protected"` // 访问级别
	Source      string   `json:"source" binding:"omitempty,oneof=manual ai import"`               // 来源
}

func (d *AddTagsToFileDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileID.required":   "文件ID不能为空",
		"FileID.len":        "文件ID长度必须为32位",
		"TagNames.required": "标签名称列表不能为空",
		"TagNames.min":      "至少需要一个标签",
		"TagNames.dive":     "标签名称格式错误",
		"AccessLevel.oneof": "访问级别必须是public、private或protected",
		"Source.oneof":      "来源必须是manual、ai或import",
	}
}

// RemoveTagsFromFileDTO 从文件移除标签DTO
type RemoveTagsFromFileDTO struct {
	FileID string `json:"file_id" binding:"required,len=32"`          // 文件ID
	TagIDs []uint `json:"tag_ids" binding:"required,min=1,dive,gt=0"` // 标签ID列表
}

func (d *RemoveTagsFromFileDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileID.required": "文件ID不能为空",
		"FileID.len":      "文件ID长度必须为32位",
		"TagIDs.required": "标签ID列表不能为空",
		"TagIDs.min":      "至少需要一个标签ID",
		"TagIDs.dive":     "标签ID格式错误",
	}
}

type GetFileTagsQueryDTO struct {
	IncludeImageCount  bool   `form:"include_file_count"`                                // 包含文件计数
	IncludeUserStats   bool   `form:"include_user_stats"`                                // 包含用户统计
	IncludeAccessStats bool   `form:"include_access_stats"`                              // 包含访问级别统计
	SortBy             string `form:"sort_by" binding:"omitempty,oneof=name created_at"` // 排序字段
}

// AddTagsToFileRequest 为文件添加标签请求
type AddTagsToFileRequest struct {
	TagNames []string `json:"tag_names" binding:"required,min=1,dive,required,min=1,max=50"`
}

// RemoveTagsFromFileRequest 从文件移除标签请求
type RemoveTagsFromFileRequest struct {
	TagIDs []uint `json:"tag_ids" binding:"required,min=1,dive,gt=0"`
}

func (d *GetFileTagsQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SortBy.oneof":    "排序字段必须是name或created_at",
		"SortOrder.oneof": "排序方向必须是asc或desc",
	}
}

type GetTagStatsQueryDTO struct {
	TagIDs      []uint `form:"tag_ids" binding:"omitempty,dive,gt=0"`                           // 标签ID列表
	UserID      *uint  `form:"user_id"`                                                         // 用户ID过滤
	AccessLevel string `form:"access_level" binding:"omitempty,oneof=public private protected"` // 访问级别过滤
	TimeRange   string `form:"time_range" binding:"omitempty,oneof=day week month year"`        // 时间范围
}

func (d *GetTagStatsQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"TagIDs.dive":       "标签ID格式错误",
		"AccessLevel.oneof": "访问级别必须是public、private或protected",
		"TimeRange.oneof":   "时间范围必须是day、week、month或year",
	}
}

// SearchTagsQueryDTO 搜索标签查询DTO
type SearchTagsQueryDTO struct {
	Keyword           string `form:"keyword" binding:"required,min=1"`                                // 搜索关键字
	UserID            *uint  `form:"user_id"`                                                         // 用户ID过滤
	AccessLevel       string `form:"access_level" binding:"omitempty,oneof=public private protected"` // 访问级别过滤
	IncludeImageCount bool   `form:"include_file_count"`                                              // 包含文件计数
	Size              int    `form:"size" binding:"omitempty,min=1,max=100"`                          // 返回数量
}

func (d *SearchTagsQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Keyword.required":  "搜索关键字不能为空",
		"Keyword.min":       "搜索关键字不能为空",
		"AccessLevel.oneof": "访问级别必须是public、private或protected",
		"Size.min":          "返回数量必须大于等于1",
		"Size.max":          "返回数量不能超过100",
	}
}

// CreateTagDTO 创建标签DTO（版本）
type CreateTagDTO struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`    // 标签名称
	Description string `json:"description" binding:"omitempty,max=500"` // 标签描述
	SortOrder   *int   `json:"sort_order"`                              // 排序顺序
}

func (d *CreateTagDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.required":   "标签名称不能为空",
		"Name.min":        "标签名称不能为空",
		"Name.max":        "标签名称不能超过50个字符",
		"Description.max": "标签描述不能超过500个字符",
	}
}

type UpdateTagDTO struct {
	ID          uint   `json:"id" binding:"required"`                   // 标签ID
	Name        string `json:"name" binding:"omitempty,min=1,max=50"`   // 标签名称
	Description string `json:"description" binding:"omitempty,max=500"` // 标签描述
	SortOrder   *int   `json:"sort_order"`                              // 排序顺序
}

func (d *UpdateTagDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required":     "标签ID不能为空",
		"Name.min":        "标签名称不能为空",
		"Name.max":        "标签名称不能超过50个字符",
		"Description.max": "标签描述不能超过500个字符",
	}
}

type DeleteTagDTO struct {
	ID           uint `json:"id" binding:"required"`                   // 标签ID
	ForceDelete  bool `json:"force_delete"`                            // 强制删除（即使有关联文件）
	TransferToID uint `json:"transfer_to_id" binding:"omitempty,gt=0"` // 转移到的标签ID
}

func (d *DeleteTagDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required":     "标签ID不能为空",
		"TransferToID.gt": "转移目标标签ID必须大于0",
	}
}

// BatchOperateTagsDTO 批量操作标签DTO（版本）
type BatchOperateTagsDTO struct {
	Action   string `json:"action" binding:"required,oneof=delete merge transfer"` // 操作类型
	TagIDs   []uint `json:"tag_ids" binding:"required,min=1,dive,gt=0"`            // 标签ID列表
	TargetID uint   `json:"target_id" binding:"omitempty,gt=0"`                    // 目标标签ID（合并/转移时使用）
	Force    bool   `json:"force"`                                                 // 强制执行
}

func (d *BatchOperateTagsDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Action.required": "操作类型不能为空",
		"Action.oneof":    "操作类型必须是delete、merge或transfer",
		"TagIDs.required": "标签ID列表不能为空",
		"TagIDs.min":      "至少需要选择一个标签",
		"TagIDs.dive":     "标签ID格式错误",
		"TargetID.gt":     "目标标签ID必须大于0",
	}
}

// TagStatsResponseDTO 标签统计响应DTO
type TagStatsResponseDTO struct {
	TagID       uint             `json:"tag_id"`
	TagName     string           `json:"tag_name"`
	TotalCount  int64            `json:"total_count"`
	UserStats   map[uint]int64   `json:"user_stats,omitempty"`   // 用户维度统计
	AccessStats map[string]int64 `json:"access_stats,omitempty"` // 访问级别维度统计
	TimeStats   map[string]int64 `json:"time_stats,omitempty"`   // 时间维度统计
	LastUpdated string           `json:"last_updated"`
}

// TagWithStatsResponseDTO 带统计信息的标签响应DTO
type TagWithStatsResponseDTO struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	IsSystem    bool             `json:"is_system"`
	CreatorID   *uint            `json:"creator_id"`
	SortOrder   int              `json:"sort_order"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
	FileCount   int64            `json:"file_count"`
	UserStats   map[uint]int64   `json:"user_stats,omitempty"`
	AccessStats map[string]int64 `json:"access_stats,omitempty"`
}

// TagListResponseDTO 标签列表响应DTO
type TagListResponseDTO struct {
	Tags       []TagWithStatsResponseDTO `json:"tags"`
	Total      int64                     `json:"total"`
	Pagination map[string]interface{}    `json:"pagination,omitempty"`
}
