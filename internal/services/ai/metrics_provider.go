package ai

import "pixelpunk/internal/metrics"

// Register AI stats provider to avoid import cycles with metrics.
func init() {
	metrics.SetAIStatsProvider(GetAIQueueStats)
}
