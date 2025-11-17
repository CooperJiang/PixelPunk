package watermark

import (
	"image"
)

type WatermarkType string

const (
	TypeImage WatermarkType = "image"
	TypeFile  WatermarkType = "image" // 向后兼容别名
)

type WatermarkPosition string

const (
	PositionTopLeft      WatermarkPosition = "top-left"
	PositionTopCenter    WatermarkPosition = "top-center"
	PositionTopRight     WatermarkPosition = "top-right"
	PositionMiddleLeft   WatermarkPosition = "middle-left"
	PositionMiddleCenter WatermarkPosition = "middle-center"
	PositionMiddleRight  WatermarkPosition = "middle-right"
	PositionBottomLeft   WatermarkPosition = "bottom-left"
	PositionBottomCenter WatermarkPosition = "bottom-center"
	PositionBottomRight  WatermarkPosition = "bottom-right"
	PositionCustom       WatermarkPosition = "custom"
)

type WatermarkConfig struct {
	Enabled bool          `json:"enabled"`
	Type    WatermarkType `json:"type"`

	FileURL    string `json:"fileUrl"`
	FileBase64 string `json:"fileBase64"` // 文件base64数据（前端生成的水印也用这个字段）
	FileData   []byte `json:"-"`          // 文件字节数据，不序列化

	// 位置配置 - 锚点+边距模式（语义：距离参考边缘的距离）
	Position   WatermarkPosition `json:"position"`
	OffsetX    float64           `json:"offsetX"`              // 距离参考边缘的距离
	OffsetY    float64           `json:"offsetY"`              // 距离参考边缘的距离
	OffsetUnit string            `json:"offsetUnit,omitempty"` // "px" | "percent"
	Opacity    float64           `json:"opacity"`              // 0-1
	Scale      float64           `json:"scale"`                // 0.1-2
	Rotation   int               `json:"rotation"`             // -180 到 180 度

	Shadow        bool   `json:"shadow"`
	ShadowColor   string `json:"shadowColor"`
	ShadowBlur    int    `json:"shadowBlur"`
	ShadowOffsetX int    `json:"shadowOffsetX"`
	ShadowOffsetY int    `json:"shadowOffsetY"`

	// Canvas生成相关（兼容字段，不推荐使用）
	GeneratedImage string `json:"generatedImage,omitempty"` // 别名字段，等同于 FileBase64
	GeneratedFile  string `json:"generatedFile,omitempty"`  // 别名字段，等同于 FileBase64
}

type ProcessorInterface interface {
	ProcessImage(img image.Image, config *WatermarkConfig) (image.Image, error)
	ProcessImageBytes(imgData []byte, config *WatermarkConfig) ([]byte, error)
	SetWatermarkImagePath(path string)
	ValidateConfig(config *WatermarkConfig) error
	GetSupportedFormats() []string
	SetEnableTextWatermark(enabled bool) // 预留的文字水印开关（当前无效）
}

type Position struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}

type WatermarkResult struct {
	Success        bool        `json:"success"`
	ProcessedImage image.Image `json:"-"`
	ProcessedData  []byte      `json:"-"`
	ErrorMessage   string      `json:"error_message,omitempty"`
	ProcessTime    int64       `json:"process_time"` // 毫秒
}

var DefaultWatermarkConfig = &WatermarkConfig{
	Enabled: false,
	Type:    TypeImage,

	FileURL: "/pixelpunk.png",

	// 位置默认值 - 锚点+边距模式
	Position:   PositionBottomRight,
	OffsetX:    20, // 距离右边缘20px
	OffsetY:    20, // 距离底边缘20px
	OffsetUnit: "px",
	Opacity:    0.7,
	Scale:      1.0,
	Rotation:   0,

	Shadow:        true,
	ShadowColor:   "#000000",
	ShadowBlur:    4,
	ShadowOffsetX: 1,
	ShadowOffsetY: 1,
}
