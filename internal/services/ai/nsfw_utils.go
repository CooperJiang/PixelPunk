package ai

import (
	"pixelpunk/internal/services/setting"
)

// getNSFWThresholdFromSettings 获取当前NSFW阈值配置（直接从数据库读取，绕过缓存）
func getNSFWThresholdFromSettings() float64 {
	return setting.GetFloatDirectFromDB("ai", "nsfw_threshold", 0.6)
}

// applyNSFWThreshold 根据阈值判断是否为NSFW内容
func applyNSFWThreshold(nsfwScore float64, aiDecision bool) bool {
	threshold := getNSFWThresholdFromSettings()
	if nsfwScore >= 0 && nsfwScore <= 1 {
		thresholdBased := nsfwScore >= threshold
		return thresholdBased
	}
	return aiDecision
}
