package dto

type CreateRandomAPIDTO struct {
	Name       string  `json:"name" binding:"required,max=100"`                      // API名称
	FolderID   *string `json:"folder_id"`                                            // 文件夹ID，null表示全部图片
	ReturnType string  `json:"return_type" binding:"required,oneof=redirect direct"` // 返回类型：redirect 或 direct
}

type UpdateRandomAPIStatusDTO struct {
	Status int `json:"status" binding:"required,oneof=1 2"` // 1:激活 2:禁用
}

type UpdateRandomAPIConfigDTO struct {
	FolderID   *string `json:"folder_id"`                                            // 文件夹ID，null表示全部图片
	ReturnType string  `json:"return_type" binding:"required,oneof=redirect direct"` // 返回类型：redirect 或 direct
}

type RandomAPIQueryDTO struct {
	Page   int    `form:"page"`   // 页码
	Size   int    `form:"size"`   // 每页数量
	Status int    `form:"status"` // 状态过滤
	Search string `form:"search"` // 名称搜索
}
