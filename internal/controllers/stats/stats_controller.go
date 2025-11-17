package stats

import (
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserStats(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	data, err := stats.GetUserStats(userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取成功")
}

func LatestFiles(c *gin.Context) {
	limit := 20

	limitParam := c.DefaultQuery("limit", "20")
	if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
		limit = l
	}

	data, err := stats.GetLatestFiles(limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, data, "获取最新文件成功")
}
