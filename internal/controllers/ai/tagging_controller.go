package ai

import (
	"fmt"
	"pixelpunk/internal/controllers/ai/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/ai"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetTaggingStatus(c *gin.Context) {
	req, err := common.ValidateRequest[dto.TaggingQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	if req.OrderBy == "" {
		req.OrderBy = "updated_at"
	}
	if req.Order == "" {
		req.Order = "desc"
	}
	result, err := ai.GetTaggingStatus(req.Status, req.Page, req.Limit, req.OrderBy, req.Order)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	errors.ResponseSuccess(c, result, "获取成功")
}

func GetTaggingStats(c *gin.Context) {
	result, err := ai.GetTaggingStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	errors.ResponseSuccess(c, result, "获取成功")
}

func ResetStuckTasks(c *gin.Context) {
	req, err := common.ValidateRequest[dto.ResetPendingDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	userID := middleware.GetCurrentUserID(c)
	count, err := ai.ResetStuckTaggingTasks(req.TimeThresholdMinutes, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	errors.ResponseSuccess(c, gin.H{"reset_count": count}, "成功重置卡住的任务")
}

func IgnoreTagging(c *gin.Context) {
	req, err := common.ValidateRequest[dto.IgnoreTaggingDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	operatorID := middleware.GetCurrentUserID(c)
	count, ignoredIDs, svcErr := ai.IgnoreTagging(req.FileIDs, req.Reason, operatorID)
	if svcErr != nil {
		errors.HandleError(c, svcErr)
		return
	}
	errors.ResponseSuccess(c, map[string]interface{}{
		"count":       count,
		"ignored_ids": ignoredIDs,
	}, "已忽略指定文件")
}

func UnignoreTagging(c *gin.Context) {
	req, err := common.ValidateRequest[dto.UnignoreTaggingDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	operatorID := middleware.GetCurrentUserID(c)
	count, updatedIDs, enqueued, skipped, svcErr := ai.UnignoreTagging(req.FileIDs, operatorID)
	if svcErr != nil {
		errors.HandleError(c, svcErr)
		return
	}
	errors.ResponseSuccess(c, map[string]interface{}{
		"count":       count,
		"updated_ids": updatedIDs,
		"enqueued":    enqueued,
		"skipped":     skipped,
	}, "已取消忽略,可重新打标")
}

func GetQueueStats(c *gin.Context) {
	stats, err := ai.GetAIQueueStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	errors.ResponseSuccess(c, stats, "获取队列状态成功")
}

func SetConcurrency(c *gin.Context) {
	type reqBody struct {
		Concurrency int `json:"concurrency"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil || body.Concurrency <= 0 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误"))
		return
	}
	svc := ai.GetGlobalTaggingService()
	if svc == nil {
		if initErr := ai.InitGlobalTaggingQueue(); initErr != nil {
			errors.HandleError(c, errors.New(errors.CodeInternal, "队列服务初始化失败: "+initErr.Error()))
			return
		}
		svc = ai.GetGlobalTaggingService()
		if svc == nil {
			errors.HandleError(c, errors.New(errors.CodeInternal, "AI功能未启用，请先在AI设置中启用AI功能"))
			return
		}
	}
	if err := svc.UpdateConcurrency(body.Concurrency); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "刷新并发失败"))
		return
	}
	stats, _ := ai.GetAIQueueStats()
	errors.ResponseSuccess(c, gin.H{"configured": body.Concurrency, "stats": stats}, "并发已更新")
}

func RetryFailedAll(c *gin.Context) {
	req, err := common.ValidateRequest[dto.RetryFailedAllDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	operatorID := middleware.GetCurrentUserID(c)
	selected, enqueued, skipped, svcErr := ai.RetryFailedAll(req.Limit, operatorID)
	if svcErr != nil {
		errors.HandleError(c, svcErr)
		return
	}
	errors.ResponseSuccess(c, gin.H{
		"selected": selected,
		"enqueued": enqueued,
		"skipped":  skipped,
	}, "已提交批量重试任务")
}

func ToggleAutoProcessing(c *gin.Context) {
	type reqBody struct {
		Enabled bool `json:"enabled"`
	}
	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误"))
		return
	}
	svc := ai.GetGlobalTaggingService()
	if svc == nil {
		if initErr := ai.InitGlobalTaggingQueue(); initErr != nil {
			errors.HandleError(c, errors.New(errors.CodeInternal, "队列服务初始化失败: "+initErr.Error()))
			return
		}
		svc = ai.GetGlobalTaggingService()
		if svc == nil {
			errors.HandleError(c, errors.New(errors.CodeInternal, "AI功能未启用，请先在AI设置中启用AI功能"))
			return
		}
	}
	if body.Enabled {
		svc.Resume()
	} else {
		svc.Pause()
	}
	errors.ResponseSuccess(c, gin.H{"enabled": body.Enabled}, "已更新自动处理状态")
}

func GetAutoProcessingStatus(c *gin.Context) {
	svc := ai.GetGlobalTaggingService()
	enabled := true
	recent := 0
	if svc != nil {
		enabled = !svc.IsPaused()
		recent = svc.RecentFailures()
	}
	errors.ResponseSuccess(c, gin.H{"enabled": enabled, "recent_failures": recent}, "获取成功")
}

func RetryTagging(c *gin.Context) {
	req, err := common.ValidateRequest[dto.RetryTaggingDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	operatorID := middleware.GetCurrentUserID(c)
	requested, enqueued, skipped, svcErr := ai.RetryTagging(req.FileIDs, operatorID)
	if svcErr != nil {
		errors.HandleError(c, svcErr)
		return
	}
	errors.ResponseSuccess(c, gin.H{
		"requested": requested,
		"enqueued":  enqueued,
		"skipped":   skipped,
	}, "已提交重试任务")
}

func TriggerTagging(c *gin.Context) {
	req, err := common.ValidateRequest[dto.TriggerTaggingDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	count, svcErr := ai.TriggerFullTaggingProcess(req.MaxFiles)
	if svcErr != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, svcErr.Error()))
		return
	}
	status := "ok"
	if count == 0 {
		status = "empty"
	}
	errors.ResponseSuccess(c, gin.H{
		"submitted_count": count,
		"status":          status,
	}, "触发打标任务成功")
}

func GetTaggingLogs(c *gin.Context) {
	type QueryParams struct {
		FileID string `form:"file_id"` // 可选,过滤特定文件
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
	}
	var params QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误"))
		return
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 20
	}
	result, err := ai.GetTaggingLogs(params.FileID, params.Page, params.Limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 转换日志数据，将 JSON 字符串解析为对象
	if logs, ok := result["logs"].([]models.FileTaggingLog); ok {
		logResponses := make([]gin.H, len(logs))
		for i, log := range logs {
			logResponses[i] = dto.TaggingLogToResponse(&log)
		}
		result["logs"] = logResponses
	}

	errors.ResponseSuccess(c, result, "获取成功")
}

func GetTaggingDiagnosis(c *gin.Context) {
	result := gin.H{
		"queue_initialized": false,
		"service_running":   false,
		"paused":            false,
		"config_issues":     []string{},
		"recommendations":   []string{},
		"worker_count":      0,
		"active_workers":    0,
		"files_in_queue":    0,
		"pending_files":     0,
	}

	globalService := ai.GetGlobalTaggingService()
	if globalService != nil {
		result["queue_initialized"] = true
		result["service_running"] = true
		result["paused"] = globalService.IsPaused()

		queueStats := globalService.GetQueueStats()
		if v, ok := queueStats["active_workers"].(int); ok {
			result["active_workers"] = v
		}
		if v, ok := queueStats["max_workers"].(int); ok {
			result["worker_count"] = v
		}

		if ext, ok := queueStats["queue_stats_ext"].(map[string]interface{}); ok {
			if queued, ok := ext["queued"].(int); ok {
				result["files_in_queue"] = queued
			}
		}
	}

	aiEnabled := setting.GetBool("ai", "ai_enabled", false)
	apiKey := setting.GetString("ai", "ai_api_key", "")

	configIssues := []string{}
	recommendations := []string{}

	if !aiEnabled {
		configIssues = append(configIssues, "AI功能未启用")
		recommendations = append(recommendations, "请在AI设置中启用AI功能（ai_enabled = true）")
	}

	if aiEnabled && (apiKey == "" || apiKey == "sk-xxxxxxxxxxxxxxx" || apiKey == "your-api-key-here") {
		configIssues = append(configIssues, "API Key未配置或为占位符")
		recommendations = append(recommendations, "请在AI设置中配置有效的API Key")
	}

	if aiEnabled && globalService == nil {
		configIssues = append(configIssues, "队列服务未初始化")
		recommendations = append(recommendations, "AI功能已启用但队列未运行，请重启服务或检查后端日志")
	}

	if globalService != nil && globalService.IsPaused() {
		configIssues = append(configIssues, "自动处理已关闭（队列已暂停）")
		recommendations = append(recommendations, "点击「自动处理」开关以启动队列自动处理")
	}

	if globalService != nil {
		filesInQueue := 0
		activeWorkers := 0
		if v, ok := result["files_in_queue"].(int); ok {
			filesInQueue = v
		}
		if v, ok := result["active_workers"].(int); ok {
			activeWorkers = v
		}

		// 查询数据库中需要处理的文件数量（none + failed且tries<3）
		var noneCount, failedCount int64
		db := ai.GetDBFromContext()
		if db != nil {
			db.Model(&models.File{}).Where("ai_tagging_status = ?", "none").Count(&noneCount)
			db.Model(&models.File{}).Where("ai_tagging_status = ? AND ai_tagging_tries < 3", "failed").Count(&failedCount)
			result["pending_files"] = int(noneCount + failedCount)

			// 只在队列空闲且有待处理文件时提示（这是正常情况，需要手动触发）
			if filesInQueue == 0 && activeWorkers == 0 && (noneCount > 0 || failedCount > 0) {
				if !globalService.IsPaused() {
					configIssues = append(configIssues, fmt.Sprintf("发现%d个待处理文件（%d个未处理，%d个失败可重试）", noneCount+failedCount, noneCount, failedCount))
					recommendations = append(recommendations, "点击「修复队列」按钮将待处理文件加入队列")
				}
			}

			// 检查是否有大量失败任务（tries>=3）
			var stuckFailedCount int64
			db.Model(&models.File{}).Where("ai_tagging_status = ? AND ai_tagging_tries >= 3", "failed").Count(&stuckFailedCount)
			if stuckFailedCount > 0 {
				configIssues = append(configIssues, fmt.Sprintf("发现%d个超过重试次数的失败任务", stuckFailedCount))
				recommendations = append(recommendations, "这些任务已达到重试上限，点击「修复队列」可重置并重新处理")
			}

			var pendingCount int64
			db.Model(&models.File{}).Where("ai_tagging_status = ?", "pending").Count(&pendingCount)
			if pendingCount > 0 && activeWorkers == 0 && filesInQueue == 0 {
				configIssues = append(configIssues, fmt.Sprintf("发现%d个处于pending状态但未在处理的任务", pendingCount))
				recommendations = append(recommendations, "点击「重置卡住」按钮可重置这些任务状态")
			}
		}
	}

	result["config_issues"] = configIssues
	result["recommendations"] = recommendations

	errors.ResponseSuccess(c, result, "诊断完成")
}
