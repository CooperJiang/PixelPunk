package dto

type WorkspaceStatsDTO struct {
	TotalImages           int    `json:"total_files"`             // 文件总数
	UsedStorage           int64  `json:"used_storage"`            // 已用存储(字节)
	UsedStorageFormatted  string `json:"used_storage_formatted"`  // 已用存储格式化
	TotalStorage          int64  `json:"total_storage"`           // 总存储(字节)
	TotalStorageFormatted string `json:"total_storage_formatted"` // 总存储格式化
	TotalViews            int64  `json:"total_views"`             // 总访问量
	TotalShares           int    `json:"total_shares"`            // 总分享数

	UsedBandwidth           int64  `json:"used_bandwidth"`            // 已用带宽(字节)
	UsedBandwidthFormatted  string `json:"used_bandwidth_formatted"`  // 已用带宽格式化
	TotalBandwidth          int64  `json:"total_bandwidth"`           // 总带宽(字节)
	TotalBandwidthFormatted string `json:"total_bandwidth_formatted"` // 总带宽格式化
}

func (req *WorkspaceStatsDTO) GetValidationMessages() map[string]string {
	return map[string]string{}
}
