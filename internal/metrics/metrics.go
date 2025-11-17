package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// simple counters (uint64) to avoid external deps
var (
	taggingAckTotal  uint64
	taggingNackTotal uint64
	vectorAckTotal   uint64
	vectorNackTotal  uint64
)

func IncTaggingAck()  { atomic.AddUint64(&taggingAckTotal, 1) }
func IncTaggingNack() { atomic.AddUint64(&taggingNackTotal, 1) }
func IncVectorAck()   { atomic.AddUint64(&vectorAckTotal, 1) }
func IncVectorNack()  { atomic.AddUint64(&vectorNackTotal, 1) }

// Providers are set by respective services to avoid import cycles.
var (
	aiStatsProvider     func() (map[string]interface{}, error)
	vectorStatsProvider func() map[string]interface{}
)

// SetAIStatsProvider registers a callback to obtain AI queue stats.
func SetAIStatsProvider(fn func() (map[string]interface{}, error)) { aiStatsProvider = fn }

// SetVectorStatsProvider registers a callback to obtain vector queue stats.
func SetVectorStatsProvider(fn func() map[string]interface{}) { vectorStatsProvider = fn }

// WritePrometheus writes metrics in Prometheus text exposition format without external deps
func WritePrometheus(w io.Writer) {
	now := time.Now().Unix()

	fmt.Fprintf(w, "# HELP tagging_ack_total Total number of successful AI tagging acknowledgements.\n")
	fmt.Fprintf(w, "# TYPE tagging_ack_total counter\n")
	fmt.Fprintf(w, "tagging_ack_total %d\n", atomic.LoadUint64(&taggingAckTotal))

	fmt.Fprintf(w, "# HELP tagging_nack_total Total number of AI tagging negative acknowledgements.\n")
	fmt.Fprintf(w, "# TYPE tagging_nack_total counter\n")
	fmt.Fprintf(w, "tagging_nack_total %d\n", atomic.LoadUint64(&taggingNackTotal))

	fmt.Fprintf(w, "# HELP vector_ack_total Total number of successful vector acknowledgements.\n")
	fmt.Fprintf(w, "# TYPE vector_ack_total counter\n")
	fmt.Fprintf(w, "vector_ack_total %d\n", atomic.LoadUint64(&vectorAckTotal))

	fmt.Fprintf(w, "# HELP vector_nack_total Total number of vector negative acknowledgements.\n")
	fmt.Fprintf(w, "# TYPE vector_nack_total counter\n")
	fmt.Fprintf(w, "vector_nack_total %d\n", atomic.LoadUint64(&vectorNackTotal))

	// Dynamic gauges from services (safe if nil)
	if aiStatsProvider != nil {
		stats, err := aiStatsProvider()
		if err == nil && stats != nil {
			if qs, ok := stats["queue_stats"].(map[string]int64); ok {
				fmt.Fprintf(w, "# HELP ai_queue_total Total images tracked by AI queue (by DB states).\n")
				fmt.Fprintf(w, "# TYPE ai_queue_total gauge\n")
				fmt.Fprintf(w, "ai_queue_total %d\n", qs["total_count"])

				fmt.Fprintf(w, "# HELP ai_queue_none_count Images in none state.\n")
				fmt.Fprintf(w, "# TYPE ai_queue_none_count gauge\n")
				fmt.Fprintf(w, "ai_queue_none_count %d\n", qs["none_count"])

				fmt.Fprintf(w, "# HELP ai_queue_done_count Images in done state.\n")
				fmt.Fprintf(w, "# TYPE ai_queue_done_count gauge\n")
				fmt.Fprintf(w, "ai_queue_done_count %d\n", qs["done_count"])

				fmt.Fprintf(w, "# HELP ai_queue_failed_count Images in failed state.\n")
				fmt.Fprintf(w, "# TYPE ai_queue_failed_count gauge\n")
				fmt.Fprintf(w, "ai_queue_failed_count %d\n", qs["failed_count"])
			}
			if ext, ok := stats["queue_stats_ext"].(map[string]int); ok {
				fmt.Fprintf(w, "# HELP ai_queue_queued Items queued in backend queue.\n")
				fmt.Fprintf(w, "# TYPE ai_queue_queued gauge\n")
				fmt.Fprintf(w, "ai_queue_queued %d\n", ext["queued"])

				fmt.Fprintf(w, "# HELP ai_queue_inflight Items currently being processed (inflight).\n")
				fmt.Fprintf(w, "# TYPE ai_queue_inflight gauge\n")
				fmt.Fprintf(w, "ai_queue_inflight %d\n", ext["processing"])
			}
			if cfg, ok := stats["config"].(map[string]interface{}); ok {
				// active_workers and concurrency
				if v, ok := cfg["active_workers"].(int); ok {
					fmt.Fprintf(w, "# HELP ai_active_workers Number of active AI workers.\n")
					fmt.Fprintf(w, "# TYPE ai_active_workers gauge\n")
					fmt.Fprintf(w, "ai_active_workers %d\n", v)
				}
				if v, ok := cfg["configured_concurrency"].(int); ok {
					fmt.Fprintf(w, "# HELP ai_configured_concurrency Configured AI worker concurrency.\n")
					fmt.Fprintf(w, "# TYPE ai_configured_concurrency gauge\n")
					fmt.Fprintf(w, "ai_configured_concurrency %d\n", v)
				}
				if p, ok := cfg["paused"].(bool); ok {
					fmt.Fprintf(w, "# HELP ai_paused Whether AI queue is paused (1=true,0=false).\n")
					fmt.Fprintf(w, "# TYPE ai_paused gauge\n")
					if p {
						fmt.Fprintf(w, "ai_paused 1\n")
					} else {
						fmt.Fprintf(w, "ai_paused 0\n")
					}
				}
			}
		}

	}

	if vectorStatsProvider != nil {
		st := vectorStatsProvider()
		if ext, ok := st["queue_stats_ext"].(map[string]int); ok {
			fmt.Fprintf(w, "# HELP vector_queue_queued Items queued in vector backend queue.\n")
			fmt.Fprintf(w, "# TYPE vector_queue_queued gauge\n")
			fmt.Fprintf(w, "vector_queue_queued %d\n", ext["queued"])

			fmt.Fprintf(w, "# HELP vector_queue_inflight Items currently being processed (inflight).\n")
			fmt.Fprintf(w, "# TYPE vector_queue_inflight gauge\n")
			fmt.Fprintf(w, "vector_queue_inflight %d\n", ext["processing"])
		}
		if rt, ok := st["runtime"].(map[string]interface{}); ok {
			if v, ok := rt["active_workers"].(int); ok {
				fmt.Fprintf(w, "# HELP vector_active_workers Number of active vector workers.\n")
				fmt.Fprintf(w, "# TYPE vector_active_workers gauge\n")
				fmt.Fprintf(w, "vector_active_workers %d\n", v)
			}
			if v, ok := rt["configured_concurrency"].(int); ok {
				fmt.Fprintf(w, "# HELP vector_configured_concurrency Configured vector worker concurrency.\n")
				fmt.Fprintf(w, "# TYPE vector_configured_concurrency gauge\n")
				fmt.Fprintf(w, "vector_configured_concurrency %d\n", v)
			}
			if p, ok := rt["paused"].(bool); ok {
				fmt.Fprintf(w, "# HELP vector_paused Whether vector queue is paused (1=true,0=false).\n")
				fmt.Fprintf(w, "# TYPE vector_paused gauge\n")
				if p {
					fmt.Fprintf(w, "vector_paused 1\n")
				} else {
					fmt.Fprintf(w, "vector_paused 0\n")
				}
			}
		}
	}

	// A lightweight timestamp
	fmt.Fprintf(w, "# HELP metrics_timestamp_seconds Unix timestamp of this metrics snapshot.\n")
	fmt.Fprintf(w, "# TYPE metrics_timestamp_seconds gauge\n")
	fmt.Fprintf(w, "metrics_timestamp_seconds %d\n", now)
}
