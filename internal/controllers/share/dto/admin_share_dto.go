package dto

import "time"

type AdminShareListQueryDTO struct {
	Page      int       `form:"page" binding:"omitempty,min=1"`
	Size      int       `form:"size" binding:"omitempty,min=1,max=100"`
	Status    int       `form:"status" binding:"omitempty,oneof=1 2 3 4"`
	Keyword   string    `form:"keyword" binding:"omitempty,max=100"`
	UserID    uint      `form:"user_id" binding:"omitempty,min=1"`
	StartDate time.Time `form:"start_date" binding:"omitempty"`
	EndDate   time.Time `form:"end_date" binding:"omitempty"`
	OrderBy   string    `form:"order_by" binding:"omitempty"`
}

func (d *AdminShareListQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":     "页码必须大于等于1",
		"Size.min":     "每页数量必须大于等于1",
		"Size.max":     "每页数量必须小于等于100",
		"Status.oneof": "状态值必须是1、2、3或4",
		"Keyword.max":  "关键字不能超过100个字符",
		"UserID.min":   "用户ID必须大于0",
	}
}

type AdminShareStatusUpdateDTO struct {
	Status int    `json:"status" binding:"required,oneof=1 2 3 4"`
	Reason string `json:"reason" binding:"omitempty,max=200"`
}

func (d *AdminShareStatusUpdateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Status.required": "状态不能为空",
		"Status.oneof":    "状态值必须是1、2、3或4",
		"Reason.max":      "操作原因不能超过200个字符",
	}
}

type AdminShareStatsResponseDTO struct {
	TotalShares   int64             `json:"total_shares"`   // 总分享数
	ActiveShares  int64             `json:"active_shares"`  // 活跃分享数
	ViewsToday    int64             `json:"views_today"`    // 今日浏览量
	PopularShares []PopularShareDTO `json:"popular_shares"` // 热门分享
	StatsByDate   []ShareStatByDate `json:"stats_by_date"`  // 按日期统计
}

type PopularShareDTO struct {
	ID          string `json:"id"`
	ShareKey    string `json:"share_key"`
	Name        string `json:"name"`
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	TotalViews  int    `json:"total_views"`
	RecentViews int    `json:"recent_views"` // 最近7天的浏览量
}

type ShareStatByDate struct {
	Date       string `json:"date"`        // 日期，格式：YYYY-MM-DD
	NewShares  int    `json:"new_shares"`  // 新增分享数
	TotalViews int    `json:"total_views"` // 浏览总数
}

type AdminShareConfigDTO struct {
	DefaultExpiredDays int      `json:"default_expired_days"`
	MaxSharesPerUser   int      `json:"max_shares_per_user"`
	MaxViewsAllowed    int      `json:"max_views_allowed"`
	ForbiddenKeywords  []string `json:"forbidden_keywords"`
}

func (d *AdminShareConfigDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"DefaultExpiredDays.min": "默认过期天数不能为负数",
		"MaxSharesPerUser.min":   "每用户最大分享数不能为负数",
		"MaxViewsAllowed.min":    "最大访问次数不能为负数",
	}
}

type AdminGenerateAccessTokenResponse struct {
	AccessToken string `json:"access_token"` // 临时访问令牌
	ShareKey    string `json:"share_key"`    // 分享密钥
	ShareURL    string `json:"share_url"`    // 分享URL
	FullURL     string `json:"full_url"`     // 带令牌的完整URL
}
