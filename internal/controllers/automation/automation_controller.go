package automation

import (
	"net/http"
	"pixelpunk/internal/controllers/automation/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/automation"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)


// @Summary 获取用户自动任务总览
// @Tags 用户自动任务
// @Accept json
// @Produce json
// @Success 200 {object} dto.AutomationOverviewResponse
// @Router /user/automation/overview [get]
func GetUserAutomationOverview(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	overview, err := automation.GetUserAutomationOverview(userID)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, overview, "获取成功")
}

// @Summary 获取用户打标任务列表
// @Tags 用户自动任务
// @Accept json
// @Produce json
// @Param status query string false "状态过滤"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} dto.TaskListResponse
// @Router /user/automation/tagging/tasks [get]
func GetUserTaggingTasks(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var query dto.TaskListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	result, err := automation.GetUserTaggingTasks(userID, query)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, result, "获取成功")
}

// @Summary 获取用户向量任务列表
// @Tags 用户自动任务
// @Accept json
// @Produce json
// @Param status query string false "状态过滤"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} dto.TaskListResponse
// @Router /user/automation/vector/tasks [get]
func GetUserVectorTasks(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var query dto.TaskListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	result, err := automation.GetUserVectorTasks(userID, query)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, result, "获取成功")
}
