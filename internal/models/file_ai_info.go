package models

import (
	"encoding/json"
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* FileAIInfo 文件AI识别信息 */
type FileAIInfo struct {
	ID               uint            `gorm:"primarykey" json:"id"`
	FileID           string          `gorm:"size:32;not null;uniqueIndex" json:"file_id"`
	Description      string          `gorm:"type:text" json:"description"`       // AI生成的文件描述（用户友好）
	SearchContent    string          `gorm:"type:text" json:"search_content"`    // 专为语义搜索优化的详细描述
	SemanticKeywords json.RawMessage `gorm:"type:json" json:"semantic_keywords"` // 语义关键词JSON数组
	Tags             json.RawMessage `gorm:"type:json" json:"tags"`              // 标签JSON数组，由AI生成
	CreatedAt        common.JSONTime `json:"created_at"`                         // 创建时间
	UpdatedAt        common.JSONTime `json:"updated_at"`                         // 更新时间

	Width         int     `json:"width"`                           // 文件/视频宽度
	Height        int     `json:"height"`                          // 文件/视频高度
	AspectRatio   float64 `json:"aspect_ratio"`                    // 宽高比
	Resolution    string  `gorm:"size:20;index" json:"resolution"` // 分辨率等级
	FileType      string  `gorm:"size:10" json:"file_type"`        // 文件格式
	EstimatedSize string  `gorm:"size:20" json:"estimated_size"`   // 估计的文件大小

	DominantColor string          `gorm:"size:10;index" json:"dominant_color"` // 主色调HEX值
	ColorPalette  json.RawMessage `gorm:"type:json" json:"color_palette"`      // 颜色调色板(JSON数组)
	ObjectsCount  int             `json:"objects_count"`                       // 识别到的物体数量
	Composition   string          `gorm:"size:20" json:"composition"`          // 主体构图

	DocumentText    string `gorm:"type:longtext" json:"document_text"` // 文档提取的文本内容
	DocumentSummary string `gorm:"type:text" json:"document_summary"`  // 文档摘要
	PageCount       int    `json:"page_count"`                         // 文档页数
	Language        string `gorm:"size:10" json:"language"`            // 文档语言
	DocumentType    string `gorm:"size:20" json:"document_type"`       // 文档类型细分

	Duration   float64 `json:"duration"`                   // 音视频时长（秒）
	Bitrate    int     `json:"bitrate"`                    // 比特率
	SampleRate int     `json:"sample_rate"`                // 采样率（音频）
	Channels   int     `json:"channels"`                   // 声道数（音频）
	VideoCodec string  `gorm:"size:20" json:"video_codec"` // 视频编码
	AudioCodec string  `gorm:"size:20" json:"audio_codec"` // 音频编码

	IsNSFW         bool            `gorm:"default:false;index" json:"is_nsfw"` // 是否包含不适内容
	NSFWScore      float64         `gorm:"index" json:"nsfw_score"`            // 不适内容综合评分
	NSFWCategories json.RawMessage `gorm:"type:json" json:"nsfw_categories"`   // NSFW分类评分
	NSFWEvaluation string          `gorm:"size:20" json:"nsfw_evaluation"`     // NSFW评估结果
	NSFWReason     string          `gorm:"type:text" json:"nsfw_reason"`       // NSFW原因详细说明
}

func (FileAIInfo) TableName() string {
	return "file_ai_info"
}

func (info *FileAIInfo) BeforeCreate(tx *gorm.DB) error {
	if info.NSFWEvaluation == "" {
		info.NSFWEvaluation = "安全"
	}

	return nil
}
