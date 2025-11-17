package path

import (
	"strings"
)

// ExtractObjectPathFromURL 从完整URL中提取对象路径
// 兼容 OSS/COS/S3 兼容存储（如 RainyUN）
// 返回对象键（例如：files/.. 或 thumbnails/.. 或其他以 bucket 后路径为主）
func ExtractObjectPathFromURL(fullURL string) string {
	if !strings.Contains(fullURL, "/") {
		return ""
	}

	parts := strings.Split(fullURL, "/")

	// OSS: https://bucket.oss-region.aliyuncs.com/path/to/file
	if len(parts) >= 4 && strings.Contains(parts[2], ".oss-") && strings.Contains(parts[2], ".aliyuncs.com") {
		return strings.Join(parts[3:], "/")
	}

	// COS: https://bucket-appid.cos.region.myqcloud.com/path/to/file
	if len(parts) >= 4 && strings.Contains(parts[2], ".cos.") && strings.Contains(parts[2], ".myqcloud.com") {
		return strings.Join(parts[3:], "/")
	}

	// S3 兼容: https://endpoint/bucket/path/to/file
	if len(parts) >= 5 {
		// 从 bucket 名之后作为对象路径
		return strings.Join(parts[4:], "/")
	}

	return ""
}

// IsHTTPURL 判断是否是 http/https URL
func IsHTTPURL(s string) bool {
	ls := strings.ToLower(s)
	return strings.HasPrefix(ls, "http://") || strings.HasPrefix(ls, "https://")
}
