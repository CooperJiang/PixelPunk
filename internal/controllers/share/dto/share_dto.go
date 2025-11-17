package dto

type ShareItemDTO struct {
	ItemType string `json:"item_type" binding:"required,oneof=folder file"`
	ItemID   string `json:"item_id" binding:"required"`
}

type CreateShareDTO struct {
	Name        string         `json:"name" binding:"omitempty,max=100"`
	Description string         `json:"description" binding:"omitempty"`
	Password    string         `json:"password" binding:"omitempty,max=100"`
	ExpiredDays int            `json:"expired_days" binding:"min=0"`
	MaxViews    int            `json:"max_views" binding:"min=0"`
	Items       []ShareItemDTO `json:"items" binding:"required,min=1,dive"`

	CollectVisitorInfo    bool `json:"collect_visitor_info"`
	NotificationOnAccess  bool `json:"notification_on_access"`
	NotificationThreshold int  `json:"notification_threshold" binding:"omitempty,min=1"`
}

func (d *CreateShareDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.max":                  "分享名称不能超过100个字符",
		"Password.max":              "密码不能超过100个字符",
		"ExpiredDays.min":           "过期天数不能为负数",
		"MaxViews.min":              "最大访问次数不能为负数",
		"Items.required":            "分享项目不能为空",
		"Items.min":                 "至少需要分享一个项目",
		"ItemType.required":         "项目类型不能为空",
		"ItemType.oneof":            "项目类型必须是folder或file",
		"ItemID.required":           "项目ID不能为空",
		"NotificationThreshold.min": "通知阈值必须大于0",
	}
}

type ShareListQueryDTO struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Size    int    `form:"size" binding:"omitempty,min=1,max=100"`
	Status  int    `form:"status" binding:"omitempty,oneof=1 2 3 4"`
	Keyword string `form:"keyword" binding:"omitempty,max=100"`
}

func (d *ShareListQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":     "页码必须大于等于1",
		"Size.min":     "每页数量必须大于等于1",
		"Size.max":     "每页数量必须小于等于100",
		"Status.oneof": "状态值必须是1、2、3或4",
		"Keyword.max":  "关键字不能超过100个字符",
	}
}

type VerifySharePasswordDTO struct {
	Password string `json:"password" binding:"required"`
}

func (d *VerifySharePasswordDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Password.required": "密码不能为空",
	}
}

type VerifySharePasswordResponse struct {
	AccessToken string `json:"access_token"` // 临时访问令牌
}

type VisitorInfoDTO struct {
	Name  string `json:"name" binding:"required,max=100"`
	Email string `json:"email" binding:"omitempty,email,max=100"`
}

func (d *VisitorInfoDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Name.required": "访客姓名不能为空",
		"Name.max":      "访客姓名不能超过100个字符",
		"Email.email":   "邮箱格式不正确",
		"Email.max":     "邮箱不能超过100个字符",
	}
}

type VisitorQueryDTO struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Size    int    `form:"size" binding:"omitempty,min=1,max=100"`
	Keyword string `form:"keyword" binding:"omitempty,max=100"`
	OrderBy string `form:"order_by" binding:"omitempty"`
}

func (d *VisitorQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":    "页码必须大于等于1",
		"Size.min":    "每页数量必须大于等于1",
		"Size.max":    "每页数量必须小于等于100",
		"Keyword.max": "关键字不能超过100个字符",
	}
}
