package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// 按照RFC 5987标准，同时兼容旧版浏览器
func SetContentDispositionFilename(filename string) string {
	if isASCII(filename) {
		// 纯ASCII文件名，使用简单格式
		return fmt.Sprintf(`attachment; filename="%s"`, escapeQuotes(filename))
	}

	// 包含中文或其他Unicode字符，使用RFC 5987格式
	// 同时提供filename和filename*参数以兼容不同浏览器
	asciiFilename := toASCIIFallback(filename)
	encodedFilename := url.QueryEscape(filename)

	return fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`,
		escapeQuotes(asciiFilename), encodedFilename)
}

// isASCII 检查字符串是否只包含ASCII字符
func isASCII(s string) bool {
	for _, char := range s {
		if char > 127 {
			return false
		}
	}
	return true
}

// escapeQuotes 转义文件名中的引号，防止HTTP头注入
func escapeQuotes(filename string) string {
	return strings.ReplaceAll(filename, `"`, `\"`)
}

// toASCIIFallback 将非ASCII字符转换为ASCII兼容版本
// 用于不支持RFC 5987的旧版浏览器
func toASCIIFallback(filename string) string {
	// 这是一个简单的实现，可以根据需要扩展

	// 使用正则表达式移除非ASCII字符，保留字母数字和常见符号
	reg := regexp.MustCompile(`[^\x20-\x7E]`)
	ascii := reg.ReplaceAllString(filename, "")

	// 清理多余的空格和特殊字符
	ascii = strings.TrimSpace(ascii)
	ascii = regexp.MustCompile(`\s+`).ReplaceAllString(ascii, "_")

	if ascii == "" || ascii == "." {
		// 尝试保留文件扩展名
		if idx := strings.LastIndex(filename, "."); idx > 0 {
			ext := filename[idx:]
			if isASCII(ext) {
				return "download" + ext
			}
		}
		return "download"
	}

	return ascii
}

func GetSafeFilename(filename string) string {
	dangerous := []string{"/", "\\", "..", ":", "*", "?", "\"", "<", ">", "|"}
	safe := filename

	for _, char := range dangerous {
		safe = strings.ReplaceAll(safe, char, "_")
	}

	safe = strings.Trim(safe, " .")

	return safe
}
