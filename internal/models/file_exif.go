package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/* FileEXIF 文件 EXIF 元数据模型 */
type FileEXIF struct {
	ID        string    `gorm:"primarykey;size:32" json:"id"`
	FileID    string    `gorm:"uniqueIndex;size:32;not null" json:"file_id"` // 关联 file 表
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Make             string `gorm:"size:100;index" json:"make,omitempty"`         // 制造商 (Apple, Canon, Nikon)
	Model            string `gorm:"size:100;index" json:"model,omitempty"`        // 型号 (iPhone 15 Pro, EOS R5)
	LensModel        string `gorm:"size:255" json:"lens_model,omitempty"`         // 镜头型号
	LensMake         string `gorm:"size:100" json:"lens_make,omitempty"`          // 镜头制造商
	SerialNumber     string `gorm:"size:100" json:"serial_number,omitempty"`      // 机身序列号
	LensSerialNumber string `gorm:"size:100" json:"lens_serial_number,omitempty"` // 镜头序列号

	FNumber           *float64 `gorm:"index" json:"f_number,omitempty"`        // 光圈值 (1.8, 2.8, 5.6)
	ExposureTime      string   `gorm:"size:20" json:"exposure_time,omitempty"` // 曝光时间 (1/1000, 1/250)
	ISO               *int     `gorm:"index" json:"iso,omitempty"`             // ISO感光度 (100, 400, 3200)
	FocalLength       *float64 `gorm:"index" json:"focal_length,omitempty"`    // 焦距 (24mm, 50mm, 200mm)
	FocalLengthIn35mm *int     `json:"focal_length_in_35mm,omitempty"`         // 35mm等效焦距

	GPSLatitude     *float64 `gorm:"index" json:"gps_latitude,omitempty"`       // 纬度
	GPSLongitude    *float64 `gorm:"index" json:"gps_longitude,omitempty"`      // 经度
	GPSAltitude     *float64 `json:"gps_altitude,omitempty"`                    // 海拔 (米)
	GPSLatitudeRef  string   `gorm:"size:1" json:"gps_latitude_ref,omitempty"`  // 纬度参考 (N/S)
	GPSLongitudeRef string   `gorm:"size:1" json:"gps_longitude_ref,omitempty"` // 经度参考 (E/W)

	DateTime          *time.Time `json:"date_time,omitempty"`                       // 文件修改时间
	DateTimeOriginal  *time.Time `gorm:"index" json:"date_time_original,omitempty"` // 原始拍摄时间
	DateTimeDigitized *time.Time `json:"date_time_digitized,omitempty"`             // 数字化时间

	Orientation    *int     `json:"orientation,omitempty"`     // 旋转方向 (1-8)
	ColorSpace     *int     `json:"color_space,omitempty"`     // 色彩空间 (1=sRGB)
	XResolution    *float64 `json:"x_resolution,omitempty"`    // X轴分辨率 (DPI)
	YResolution    *float64 `json:"y_resolution,omitempty"`    // Y轴分辨率 (DPI)
	ResolutionUnit *int     `json:"resolution_unit,omitempty"` // 分辨率单位 (2=英寸, 3=厘米)

	ExposureProgram  *int     `json:"exposure_program,omitempty"`   // 曝光程序 (0=未定义, 1=手动, 2=程序自动, 3=光圈优先, 4=快门优先)
	ExposureBias     *float64 `json:"exposure_bias,omitempty"`      // 曝光补偿 (-2, -1, 0, +1, +2 EV)
	MeteringMode     *int     `json:"metering_mode,omitempty"`      // 测光模式 (1=平均, 2=中央重点, 3=点测光, 5=矩阵)
	Flash            *int     `json:"flash,omitempty"`              // 闪光灯状态
	WhiteBalance     *int     `json:"white_balance,omitempty"`      // 白平衡 (0=自动, 1=手动)
	SceneCaptureType *int     `json:"scene_capture_type,omitempty"` // 场景类型 (0=标准, 1=风景, 2=人像, 3=夜景)
	DigitalZoomRatio *float64 `json:"digital_zoom_ratio,omitempty"` // 数字变焦比
	Contrast         *int     `json:"contrast,omitempty"`           // 对比度 (0=普通, 1=柔和, 2=强烈)
	Saturation       *int     `json:"saturation,omitempty"`         // 饱和度
	Sharpness        *int     `json:"sharpness,omitempty"`          // 锐度

	SubjectDistance  *float64 `json:"subject_distance,omitempty"`   // 对焦距离 (米)
	MaxApertureValue *float64 `json:"max_aperture_value,omitempty"` // 最大光圈值

	Software         string `gorm:"size:100" json:"software,omitempty"`           // 编辑软件 (Photoshop, Lightroom)
	Artist           string `gorm:"size:100" json:"artist,omitempty"`             // 作者/摄影师
	Copyright        string `gorm:"size:255" json:"copyright,omitempty"`          // 版权信息
	ImageDescription string `gorm:"type:text" json:"image_description,omitempty"` // 图片描述

	File *File `gorm:"foreignKey:FileID;references:ID" json:"file,omitempty"`
}

func (FileEXIF) TableName() string {
	return "file_exif"
}

func (fe *FileEXIF) BeforeCreate(tx *gorm.DB) error {
	if fe.ID == "" {
		fe.ID = uuid.New().String()[:32]
	}
	return nil
}

func (fe *FileEXIF) HasGPS() bool {
	return fe.GPSLatitude != nil && fe.GPSLongitude != nil
}

func (fe *FileEXIF) GetLocation() string {
	if !fe.HasGPS() {
		return ""
	}
	return formatGPSCoordinate(*fe.GPSLatitude, *fe.GPSLongitude)
}

func (fe *FileEXIF) GetCameraInfo() string {
	if fe.Make == "" && fe.Model == "" {
		return ""
	}
	if fe.Make == "" {
		return fe.Model
	}
	if fe.Model == "" {
		return fe.Make
	}
	return fe.Make + " " + fe.Model
}

func (fe *FileEXIF) GetShootingParams() string {
	params := ""

	if fe.FocalLength != nil {
		params += formatFloat(*fe.FocalLength, 0) + "mm "
	}

	if fe.FNumber != nil {
		params += "f/" + formatFloat(*fe.FNumber, 1) + " "
	}

	if fe.ExposureTime != "" {
		params += fe.ExposureTime + "s "
	}

	if fe.ISO != nil {
		params += "ISO" + formatInt(*fe.ISO) + " "
	}

	return params
}

// 辅助格式化函数
func formatGPSCoordinate(lat, lon float64) string {
	return formatFloat(lat, 6) + ", " + formatFloat(lon, 6)
}

func formatFloat(val float64, precision int) string {
	if precision == 0 {
		return fmt.Sprintf("%.0f", val)
	}
	return fmt.Sprintf("%."+fmt.Sprintf("%d", precision)+"f", val)
}

func formatInt(val int) string {
	return fmt.Sprintf("%d", val)
}
