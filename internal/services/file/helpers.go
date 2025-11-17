package file

import "pixelpunk/pkg/imagex/formats"

/* GetCorrectFileExtension 根据图像格式返回正确的文件扩展名 */
func GetCorrectFileExtension(format string) string {
	normalized := formats.NormalizeFormat(format)
	if normalized == "jpeg" {
		return "jpg"
	}
	if normalized == "tif" {
		return "tiff"
	}
	if formats.IsSupported(normalized) {
		return normalized
	}
	return normalized
}

/* GetContentTypeByFormat 根据文件格式返回正确的Content-Type */
func GetContentTypeByFormat(format string) string { return formats.GetContentType(format) }
