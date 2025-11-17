package metrics

import (
	metrics "pixelpunk/internal/metrics"

	"github.com/gin-gonic/gin"
)

func PrometheusHandler(c *gin.Context) {
	c.Header("Content-Type", "text/plain; version=0.0.4")
	metrics.WritePrometheus(c.Writer)
}
