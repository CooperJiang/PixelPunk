package dto

// DownloadFileDTO 下载文件DTO
type DownloadFileDTO struct {
	// 暂时不需要任何参数，文件ID从URL路径获取
	// 质量和格式参数暂时不支持，直接下载原图
}

func (d *DownloadFileDTO) GetValidationMessages() map[string]string {
	return map[string]string{}
}
