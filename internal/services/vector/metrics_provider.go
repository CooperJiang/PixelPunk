package vector

import "pixelpunk/internal/metrics"

func vectorStats() map[string]interface{} {
	if s := GetGlobalVectorQueueService(); s != nil {
		return s.GetQueueStats()
	}
	return map[string]interface{}{}
}

func init() {
	metrics.SetVectorStatsProvider(vectorStats)
}
