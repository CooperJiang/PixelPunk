package stats

import (
	"strconv"

	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func DashboardUserStats(c *gin.Context) {
	data, err := stats.GetDashboardUserStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取用户统计成功")
}

func DashboardFileStats(c *gin.Context) {
	data, err := stats.GetDashboardFileStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取文件统计成功")
}

func DashboardStorageStats(c *gin.Context) {
	data, err := stats.GetDashboardStorageStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取存储统计成功")
}

func DashboardSystemResources(c *gin.Context) {
	data, err := stats.GetDashboardSystemResources()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取系统资源成功")
}

func DashboardUploadTrends(c *gin.Context) {
	days := 7 // 默认7天
	daysParam := c.DefaultQuery("days", "7")
	if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
		if d > 90 {
			d = 90 // 最多90天
		}
		days = d
	}

	data, err := stats.GetDashboardUploadTrends(days)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取上传趋势成功")
}

func DashboardAIServices(c *gin.Context) {
	data, err := stats.GetDashboardAIServices()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取AI服务状态成功")
}

func DashboardShareStats(c *gin.Context) {
	data, err := stats.GetDashboardShareStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取分享统计成功")
}

func DashboardTagStats(c *gin.Context) {
	data, err := stats.GetDashboardTagStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取标签统计成功")
}

func DashboardSystemInfo(c *gin.Context) {
	data, err := stats.GetDashboardSystemInfo()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取系统信息成功")
}
