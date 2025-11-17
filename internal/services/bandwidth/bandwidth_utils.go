package bandwidth

import "pixelpunk/internal/models"

/* GetUserBandwidthStats 获取用户带宽统计（包级别兼容函数） */
func GetUserBandwidthStats(userID uint) (models.UserBandwidthUsage, error) {
	return Service.GetUserBandwidthStats(userID)
}
