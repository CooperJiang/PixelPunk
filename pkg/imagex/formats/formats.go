package formats

import (
	"strings"
)

// 默认支持的扩展名（带点）
var defaultDotExtensions = []string{
	".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg", ".ico", ".apng", ".jp2", ".tiff", ".tif", ".tga", ".heic", ".heif",
}

// 扩展名到MIME映射（不带点，小写）
var extToMIME = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"bmp":  "image/bmp",
	"svg":  "image/svg+xml",
	"ico":  "image/x-icon",
	"apng": "image/apng",
	"jp2":  "image/jp2",
	"tiff": "image/tiff",
	"tif":  "image/tiff",
	"tga":  "image/x-tga",
	"heic": "image/heic",
	"heif": "image/heif",
}

// NormalizeFormat 规格化格式/扩展名（去点、转小写）
func NormalizeFormat(formatOrExt string) string {
	f := strings.TrimSpace(strings.ToLower(formatOrExt))
	f = strings.TrimPrefix(f, ".")
	return f
}

// GetContentType 根据格式或扩展名获取标准MIME
func GetContentType(formatOrExt string) string {
	f := NormalizeFormat(formatOrExt)
	if mime, ok := extToMIME[f]; ok {
		return mime
	}
	return "application/octet-stream"
}

// GetMimeTypeFromExtension 兼容命名，等价于 GetContentType
func GetMimeTypeFromExtension(ext string) string {
	return GetContentType(ext)
}

// SupportedExtensionsWithDot 返回默认支持的扩展名列表（带点）
func SupportedExtensionsWithDot() []string {
	return append([]string{}, defaultDotExtensions...)
}

// SupportedExtensionsWithoutDot 返回默认支持的扩展名列表（不带点）
func SupportedExtensionsWithoutDot() []string {
	result := make([]string, 0, len(defaultDotExtensions))
	for _, e := range defaultDotExtensions {
		result = append(result, strings.TrimPrefix(strings.ToLower(e), "."))
	}
	return result
}

// IsSupported 检查给定格式/扩展名是否在默认支持列表内
func IsSupported(formatOrExt string) bool {
	f := NormalizeFormat(formatOrExt)
	for _, e := range defaultDotExtensions {
		if strings.TrimPrefix(e, ".") == f {
			return true
		}
	}
	return false
}
