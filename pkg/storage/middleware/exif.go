package middleware

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"
)

// EXIFData EXIF数据结构
type EXIFData struct {
	Make         string    `json:"make,omitempty"`          // 相机制造商
	Model        string    `json:"model,omitempty"`         // 相机型号
	DateTime     time.Time `json:"date_time,omitempty"`     // 拍摄时间
	GPS          *GPSInfo  `json:"gps,omitempty"`           // GPS信息
	Orientation  int       `json:"orientation,omitempty"`   // 方向
	Flash        bool      `json:"flash,omitempty"`         // 是否使用闪光灯
	FocalLength  float64   `json:"focal_length,omitempty"`  // 焦距
	ExposureTime string    `json:"exposure_time,omitempty"` // 曝光时间
	FNumber      float64   `json:"f_number,omitempty"`      // 光圈值
	ISO          int       `json:"iso,omitempty"`           // ISO值
	Width        int       `json:"width,omitempty"`         // 图像宽度
	Height       int       `json:"height,omitempty"`        // 图像高度
}

// GPSInfo GPS信息
type GPSInfo struct {
	Latitude  float64 `json:"latitude"`  // 纬度
	Longitude float64 `json:"longitude"` // 经度
	Altitude  float64 `json:"altitude"`  // 海拔
}

// ExtractEXIFData 提取EXIF数据
func ExtractEXIFData(reader io.Reader) (*EXIFData, error) {
	data := make([]byte, 65536) // 读取64KB足够提取EXIF
	n, err := reader.Read(data)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("读取文件数据失败: %w", err)
	}
	data = data[:n]

	if !isJPEG(data) {
		return nil, fmt.Errorf("不是JPEG文件，无法提取EXIF数据")
	}

	exifData, err := findEXIFSegment(data)
	if err != nil {
		return nil, err
	}

	return parseEXIFData(exifData)
}

// isJPEG 检查是否为JPEG文件
func isJPEG(data []byte) bool {
	if len(data) < 3 {
		return false
	}
	return data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF
}

// findEXIFSegment 查找EXIF段
func findEXIFSegment(data []byte) ([]byte, error) {
	pos := 2 // 跳过JPEG文件头 0xFF 0xD8

	for pos < len(data)-1 {
		if data[pos] != 0xFF {
			return nil, fmt.Errorf("JPEG格式错误")
		}

		marker := data[pos+1]
		pos += 2

		if marker == 0xE1 { // APP1段，通常包含EXIF
			if pos+2 > len(data) {
				return nil, fmt.Errorf("JPEG段长度不足")
			}

			segmentLength := int(binary.BigEndian.Uint16(data[pos : pos+2]))
			pos += 2

			if pos+segmentLength-2 > len(data) {
				return nil, fmt.Errorf("EXIF段长度超出文件范围")
			}

			segment := data[pos : pos+segmentLength-2]

			if len(segment) >= 6 && string(segment[:4]) == "Exif" {
				return segment[6:], nil // 跳过"Exif\x00\x00"
			}
		} else if marker == 0xDA { // 图像数据开始
			break
		} else {
			if pos+2 > len(data) {
				break
			}
			segmentLength := int(binary.BigEndian.Uint16(data[pos : pos+2]))
			pos += segmentLength
		}
	}

	return nil, fmt.Errorf("未找到EXIF数据")
}

// parseEXIFData 解析EXIF数据（简化版本）
func parseEXIFData(data []byte) (*EXIFData, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("EXIF数据长度不足")
	}

	exif := &EXIFData{}

	var order binary.ByteOrder
	if string(data[:2]) == "II" {
		order = binary.LittleEndian
	} else if string(data[:2]) == "MM" {
		order = binary.BigEndian
	} else {
		return nil, fmt.Errorf("无效的EXIF字节序标记")
	}

	if order.Uint16(data[2:4]) != 0x002A {
		return nil, fmt.Errorf("无效的TIFF标记")
	}

	ifdOffset := order.Uint32(data[4:8])
	if int(ifdOffset) >= len(data) {
		return nil, fmt.Errorf("IFD偏移超出数据范围")
	}

	err := parseIFD(data, int(ifdOffset), order, exif)
	if err != nil {
		return nil, fmt.Errorf("解析IFD失败: %w", err)
	}

	return exif, nil
}

// parseIFD 解析IFD（Image File Directory）
func parseIFD(data []byte, offset int, order binary.ByteOrder, exif *EXIFData) error {
	if offset+2 > len(data) {
		return fmt.Errorf("IFD偏移超出范围")
	}

	entryCount := order.Uint16(data[offset : offset+2])
	offset += 2

	for i := 0; i < int(entryCount); i++ {
		if offset+12 > len(data) {
			break
		}

		tag := order.Uint16(data[offset : offset+2])
		dataType := order.Uint16(data[offset+2 : offset+4])
		count := order.Uint32(data[offset+4 : offset+8])
		valueOffset := order.Uint32(data[offset+8 : offset+12])

		// 根据标签解析对应的值
		switch tag {
		case 0x010F: // Make (制造商)
			if value := extractString(data, int(valueOffset), int(count), dataType == 2); value != "" {
				exif.Make = value
			}
		case 0x0110: // Model (型号)
			if value := extractString(data, int(valueOffset), int(count), dataType == 2); value != "" {
				exif.Model = value
			}
		case 0x0132: // DateTime (拍摄时间)
			if value := extractString(data, int(valueOffset), int(count), dataType == 2); value != "" {
				if t, err := time.Parse("2006:01:02 15:04:05", value); err == nil {
					exif.DateTime = t
				}
			}
		case 0x0112: // Orientation (方向)
			if dataType == 3 && count == 1 {
				exif.Orientation = int(order.Uint16(data[offset+8 : offset+10]))
			}
		}

		offset += 12
	}

	return nil
}

// extractString 提取字符串值
func extractString(data []byte, offset, length int, isString bool) string {
	if !isString {
		return ""
	}

	if offset >= len(data) {
		return ""
	}

	end := offset + length
	if end > len(data) {
		end = len(data)
	}

	value := string(data[offset:end])
	if idx := strings.Index(value, "\x00"); idx >= 0 {
		value = value[:idx]
	}

	return strings.TrimSpace(value)
}

// HasEXIFData 检查文件是否包含EXIF数据
func HasEXIFData(reader io.Reader) bool {
	// 读取前几KB检查EXIF标记
	buffer := make([]byte, 4096)
	n, _ := reader.Read(buffer)
	if n < 10 {
		return false
	}

	// 检查JPEG + EXIF标记
	return isJPEG(buffer) && bytes.Contains(buffer[:n], []byte("Exif"))
}

// StripEXIFData 移除EXIF数据
func StripEXIFData(reader io.Reader) (io.Reader, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取文件数据失败: %w", err)
	}

	if !isJPEG(data) {
		return bytes.NewReader(data), nil
	}

	result := make([]byte, 0, len(data))
	result = append(result, 0xFF, 0xD8) // JPEG文件头

	pos := 2
	for pos < len(data)-1 {
		if data[pos] != 0xFF {
			break
		}

		marker := data[pos+1]
		segmentStart := pos
		pos += 2

		if marker == 0xE1 { // APP1段，可能包含EXIF
			if pos+2 > len(data) {
				break
			}

			segmentLength := int(binary.BigEndian.Uint16(data[pos : pos+2]))
			segment := data[pos : pos+segmentLength]

			// 如果不是EXIF段，保留
			if len(segment) < 6 || string(segment[2:6]) != "Exif" {
				result = append(result, data[segmentStart:pos+segmentLength]...)
			}

			pos += segmentLength
		} else if marker == 0xDA { // 图像数据开始
			result = append(result, data[segmentStart:]...)
			break
		} else {
			if pos+2 > len(data) {
				break
			}
			segmentLength := int(binary.BigEndian.Uint16(data[pos : pos+2]))
			result = append(result, data[segmentStart:pos+segmentLength]...)
			pos += segmentLength
		}
	}

	return bytes.NewReader(result), nil
}
