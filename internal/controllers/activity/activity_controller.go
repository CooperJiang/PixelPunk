package activity

import (
	"pixelpunk/internal/controllers/activity/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/activity"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetUserActivities(c *gin.Context) {
	req, err := common.ValidateRequest[dto.GetUserActivitiesDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	req.SetDefaults()

	userID := middleware.GetCurrentUserID(c)

	activityService := activity.GetService()
	activities, total, todayStats, err := activityService.GetUserActivities(userID, req.Page, req.Size)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	data := gin.H{
		"list":        activities,
		"total":       total,
		"page":        req.Page,
		"size":        req.Size,
		"today_stats": todayStats,
	}

	errors.ResponseSuccess(c, data, "获取成功")
}
