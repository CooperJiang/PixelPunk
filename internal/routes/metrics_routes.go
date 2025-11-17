package routes

import (
	metricsController "pixelpunk/internal/controllers/metrics"

	"github.com/gin-gonic/gin"
)

func RegisterMetricsRoutes(r *gin.RouterGroup) {
	r.GET("/metrics", metricsController.PrometheusHandler)
}
