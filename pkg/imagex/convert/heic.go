package convert

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"runtime"
	"strings"
)

// HEICToJPEGOptions HEIC 转 JPEG 的选项
type HEICToJPEGOptions struct {
	Quality int // JPEG 质量 (1-100)，默认 95
}

// HEICToJPEGResult HEIC 转 JPEG 的结果
type HEICToJPEGResult struct {
	Reader    io.Reader // 转换后的数据流
	Converted bool      // 是否成功转换
	Width     int       // 图片宽度
	Height    int       // 图片高度
}

// ToJPEGFromHEIC 将 HEIC/HEIF 格式转换为 JPEG
// input: HEIC/HEIF 图片数据
// opts: 转换选项
// 返回: 转换结果和错误信息
func ToJPEGFromHEIC(input []byte, opts HEICToJPEGOptions) (*HEICToJPEGResult, error) {
	// Windows 平台不支持 HEIC 转换（adrium/goheif 依赖 libheif，Windows 编译困难）
	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("HEIC/HEIF conversion is not supported on Windows platform. Please convert the image to JPEG/PNG format before uploading")
	}

	quality := opts.Quality
	if quality <= 0 || quality > 100 {
		quality = 95
	}

	// 尝试解码 HEIC 图片
	img, format, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 检查是否是 HEIC/HEIF 格式
	format = strings.ToLower(format)
	if format != "heic" && format != "heif" {
		// 不是 HEIC 格式，不需要转换
		return &HEICToJPEGResult{
			Reader:    bytes.NewReader(input),
			Converted: false,
		}, nil
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, fmt.Errorf("failed to encode to JPEG: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	return &HEICToJPEGResult{
		Reader:    bytes.NewReader(buf.Bytes()),
		Converted: true,
		Width:     width,
		Height:    height,
	}, nil
}

// IsHEICFormat 检测是否是 HEIC/HEIF 格式
func IsHEICFormat(data []byte) bool {
	if len(data) < 12 {
		return false
	}

	// HEIC/HEIF 文件头特征
	// ftyp box at offset 4, 包含 "heic" 或 "mif1" 等标识
	if bytes.Contains(data[:20], []byte("ftyp")) {
		header := string(data[:32])
		if strings.Contains(header, "heic") ||
			strings.Contains(header, "heif") ||
			strings.Contains(header, "mif1") ||
			strings.Contains(header, "msf1") {
			return true
		}
	}

	return false
}
