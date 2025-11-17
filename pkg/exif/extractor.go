package exif

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	exif "github.com/dsoprea/go-exif/v3"
)

// FileEXIFData 统一的 EXIF 数据结构
type FileEXIFData struct {
	Make             string `json:"make,omitempty"`
	Model            string `json:"model,omitempty"`
	LensModel        string `json:"lens_model,omitempty"`
	LensMake         string `json:"lens_make,omitempty"`
	SerialNumber     string `json:"serial_number,omitempty"`
	LensSerialNumber string `json:"lens_serial_number,omitempty"`

	FNumber           *float64 `json:"f_number,omitempty"`
	ExposureTime      string   `json:"exposure_time,omitempty"`
	ISO               *int     `json:"iso,omitempty"`
	FocalLength       *float64 `json:"focal_length,omitempty"`
	FocalLengthIn35mm *int     `json:"focal_length_in_35mm,omitempty"`

	GPSLatitude     *float64 `json:"gps_latitude,omitempty"`
	GPSLongitude    *float64 `json:"gps_longitude,omitempty"`
	GPSAltitude     *float64 `json:"gps_altitude,omitempty"`
	GPSLatitudeRef  string   `json:"gps_latitude_ref,omitempty"`
	GPSLongitudeRef string   `json:"gps_longitude_ref,omitempty"`

	DateTime          *time.Time `json:"date_time,omitempty"`
	DateTimeOriginal  *time.Time `json:"date_time_original,omitempty"`
	DateTimeDigitized *time.Time `json:"date_time_digitized,omitempty"`

	Orientation    *int     `json:"orientation,omitempty"`
	ColorSpace     *int     `json:"color_space,omitempty"`
	XResolution    *float64 `json:"x_resolution,omitempty"`
	YResolution    *float64 `json:"y_resolution,omitempty"`
	ResolutionUnit *int     `json:"resolution_unit,omitempty"`

	ExposureProgram  *int     `json:"exposure_program,omitempty"`
	ExposureBias     *float64 `json:"exposure_bias,omitempty"`
	MeteringMode     *int     `json:"metering_mode,omitempty"`
	Flash            *int     `json:"flash,omitempty"`
	WhiteBalance     *int     `json:"white_balance,omitempty"`
	SceneCaptureType *int     `json:"scene_capture_type,omitempty"`
	DigitalZoomRatio *float64 `json:"digital_zoom_ratio,omitempty"`
	Contrast         *int     `json:"contrast,omitempty"`
	Saturation       *int     `json:"saturation,omitempty"`
	Sharpness        *int     `json:"sharpness,omitempty"`

	SubjectDistance  *float64 `json:"subject_distance,omitempty"`
	MaxApertureValue *float64 `json:"max_aperture_value,omitempty"`

	Software         string `json:"software,omitempty"`
	Artist           string `json:"artist,omitempty"`
	Copyright        string `json:"copyright,omitempty"`
	ImageDescription string `json:"image_description,omitempty"`
}

// ExtractEXIFFromBytes 从字节数据提取 EXIF（支持 JPEG, TIFF, PNG, WebP, HEIC 等）
func ExtractEXIFFromBytes(data []byte) (*FileEXIFData, error) {
	rawExif, err := exif.SearchAndExtractExif(data)
	if err != nil {
		return nil, nil
	}

	entries, _, err := exif.GetFlatExifData(rawExif, nil)
	if err != nil {
		return nil, nil
	}

	result := &FileEXIFData{}
	extractFromEntries(entries, result)

	if isEmptyEXIF(result) {
		return nil, nil
	}

	return result, nil
}

// ExtractEXIFFromReader 从 io.Reader 提取 EXIF
func ExtractEXIFFromReader(reader io.Reader) (*FileEXIFData, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %w", err)
	}
	return ExtractEXIFFromBytes(data)
}

// extractFromEntries 从 EXIF 条目列表提取数据
func extractFromEntries(entries []exif.ExifTag, result *FileEXIFData) {
	for i := range entries {
		entry := &entries[i]
		tagName := entry.TagName
		tagValue := entry.FormattedFirst

		switch tagName {
		case "Make":
			result.Make = strings.TrimSpace(tagValue)
		case "Model":
			result.Model = strings.TrimSpace(tagValue)
		case "LensModel":
			result.LensModel = strings.TrimSpace(tagValue)
		case "LensMake":
			result.LensMake = strings.TrimSpace(tagValue)
		case "BodySerialNumber":
			result.SerialNumber = strings.TrimSpace(tagValue)
		case "LensSerialNumber":
			result.LensSerialNumber = strings.TrimSpace(tagValue)

		case "FNumber":
			if val := parseFloat(tagValue); val != nil {
				result.FNumber = val
			}

		case "ExposureTime":
			if entry.ChildIfdPath == "" && tagValue != "" {
				result.ExposureTime = tagValue
			}

		case "ISOSpeedRatings", "PhotographicSensitivity":
			if val := parseIntFromString(tagValue); val != nil {
				result.ISO = val
			}

		case "FocalLength":
			if val := parseFloat(tagValue); val != nil {
				result.FocalLength = val
			}

		case "FocalLengthIn35mmFilm":
			if val := parseIntFromString(tagValue); val != nil {
				result.FocalLengthIn35mm = val
			}

		case "DateTime":
			if t, err := parseDateTime(tagValue); err == nil {
				result.DateTime = &t
			}
		case "DateTimeOriginal":
			if t, err := parseDateTime(tagValue); err == nil {
				result.DateTimeOriginal = &t
			}
		case "DateTimeDigitized":
			if t, err := parseDateTime(tagValue); err == nil {
				result.DateTimeDigitized = &t
			}

		case "Orientation":
			if val := parseIntFromString(tagValue); val != nil {
				result.Orientation = val
			}
		case "ColorSpace":
			if val := parseIntFromString(tagValue); val != nil {
				result.ColorSpace = val
			}
		case "XResolution":
			if val := parseFloat(tagValue); val != nil {
				result.XResolution = val
			}
		case "YResolution":
			if val := parseFloat(tagValue); val != nil {
				result.YResolution = val
			}
		case "ResolutionUnit":
			if val := parseIntFromString(tagValue); val != nil {
				result.ResolutionUnit = val
			}

		case "ExposureProgram":
			if val := parseIntFromString(tagValue); val != nil {
				result.ExposureProgram = val
			}
		case "ExposureBiasValue":
			if val := parseFloat(tagValue); val != nil {
				result.ExposureBias = val
			}
		case "MeteringMode":
			if val := parseIntFromString(tagValue); val != nil {
				result.MeteringMode = val
			}
		case "Flash":
			if val := parseIntFromString(tagValue); val != nil {
				result.Flash = val
			}
		case "WhiteBalance":
			if val := parseIntFromString(tagValue); val != nil {
				result.WhiteBalance = val
			}
		case "SceneCaptureType":
			if val := parseIntFromString(tagValue); val != nil {
				result.SceneCaptureType = val
			}
		case "DigitalZoomRatio":
			if val := parseFloat(tagValue); val != nil {
				result.DigitalZoomRatio = val
			}
		case "Contrast":
			if val := parseIntFromString(tagValue); val != nil {
				result.Contrast = val
			}
		case "Saturation":
			if val := parseIntFromString(tagValue); val != nil {
				result.Saturation = val
			}
		case "Sharpness":
			if val := parseIntFromString(tagValue); val != nil {
				result.Sharpness = val
			}

		case "SubjectDistance":
			if val := parseFloat(tagValue); val != nil {
				result.SubjectDistance = val
			}
		case "MaxApertureValue":
			if val := parseFloat(tagValue); val != nil {
				result.MaxApertureValue = val
			}

		case "Software":
			result.Software = strings.TrimSpace(tagValue)
		case "Artist":
			result.Artist = strings.TrimSpace(tagValue)
		case "Copyright":
			result.Copyright = strings.TrimSpace(tagValue)
		case "ImageDescription":
			result.ImageDescription = strings.TrimSpace(tagValue)

		case "GPSLatitude":
			if lat := parseGPSFromFormatted(tagValue); lat != nil {
				result.GPSLatitude = lat
			}
		case "GPSLatitudeRef":
			result.GPSLatitudeRef = strings.TrimSpace(tagValue)
		case "GPSLongitude":
			if lon := parseGPSFromFormatted(tagValue); lon != nil {
				result.GPSLongitude = lon
			}
		case "GPSLongitudeRef":
			result.GPSLongitudeRef = strings.TrimSpace(tagValue)
		case "GPSAltitude":
			if alt := parseFloat(tagValue); alt != nil {
				result.GPSAltitude = alt
			}
		}
	}

	// 应用 GPS 参考方向（南纬为负，西经为负）
	if result.GPSLatitude != nil && result.GPSLatitudeRef == "S" {
		*result.GPSLatitude = -*result.GPSLatitude
	}
	if result.GPSLongitude != nil && result.GPSLongitudeRef == "W" {
		*result.GPSLongitude = -*result.GPSLongitude
	}
}

// parseFloat 从字符串解析浮点数
func parseFloat(s string) *float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	// 处理分数格式 "1/100"
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) == 2 {
			num, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			den, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil && den != 0 {
				val := num / den
				return &val
			}
		}
	}

	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return &val
	}

	return nil
}

// parseIntFromString 从字符串解析整数
func parseIntFromString(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	if val, err := strconv.Atoi(s); err == nil {
		return &val
	}

	// 尝试解析浮点数然后转整数
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		intVal := int(math.Round(val))
		return &intVal
	}

	return nil
}

// parseDateTime 解析日期时间
func parseDateTime(timeStr string) (time.Time, error) {
	// EXIF 时间格式: "2006:01:02 15:04:05"
	return time.Parse("2006:01:02 15:04:05", timeStr)
}

// parseGPSFromFormatted 从格式化的 GPS 字符串解析坐标
func parseGPSFromFormatted(s string) *float64 {
	// 格式可能是: "39.9042" 或 "39deg 54' 15.12\"" 等
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return &val
	}

	// 尝试解析度分秒格式
	// 例如: "39deg 54' 15.12\"" 或 "[39 54 15.12]"
	s = strings.ReplaceAll(s, "deg", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")

	parts := strings.Fields(s)
	if len(parts) >= 3 {
		deg, err1 := strconv.ParseFloat(parts[0], 64)
		min, err2 := strconv.ParseFloat(parts[1], 64)
		sec, err3 := strconv.ParseFloat(parts[2], 64)
		if err1 == nil && err2 == nil && err3 == nil {
			val := deg + (min / 60.0) + (sec / 3600.0)
			return &val
		}
	}

	return nil
}

// isEmptyEXIF 检查 EXIF 数据是否为空
func isEmptyEXIF(data *FileEXIFData) bool {
	return data.Make == "" &&
		data.Model == "" &&
		data.LensModel == "" &&
		data.FNumber == nil &&
		data.ExposureTime == "" &&
		data.ISO == nil &&
		data.FocalLength == nil &&
		data.GPSLatitude == nil &&
		data.GPSLongitude == nil &&
		data.DateTimeOriginal == nil
}
