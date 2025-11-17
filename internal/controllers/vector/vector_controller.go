package vector

import (
	"encoding/json"
	"fmt"
	"net/http"
	setdto "pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/controllers/vector/dto"
	settingSvc "pixelpunk/internal/services/setting"
	vectorService "pixelpunk/internal/services/vector"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetVectorList(c *gin.Context) {
	req, err := common.ValidateRequest[dto.VectorListRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	results, total, err := vectorService.GetVectorWithFileInfo(
		req.Page, req.PageSize, req.Status, req.Model, req.Keyword)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, "获取向量列表失败: "+err.Error()))
		return
	}

	var vectorItems []dto.VectorItemResponse
	for _, result := range results {
		item := dto.VectorItemResponse{
			ID:                 result.ID,
			FileID:             result.FileID,
			Status:             result.Status,
			Model:              result.Model,
			Description:        result.Description,
			Dimension:          result.Dimension,
			RetryCount:         result.RetryCount,
			ErrorMessage:       result.ErrorMessage,
			ProcessingDuration: result.ProcessingDuration,
			CreatedAt:          result.CreatedAt,
			UpdatedAt:          result.UpdatedAt,
			LastRetryAt:        result.LastRetryAt,
			FileName:           result.FileName,
			FileURL:            result.FileURL,
			ThumbnailURL:       result.ThumbnailURL,
		}

		vectorItems = append(vectorItems, item)
	}

	totalPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	response := dto.VectorListResponse{
		Data: vectorItems,
		Pagination: dto.PaginationResponse{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	}

	errors.ResponseSuccess(c, response, "获取向量列表成功")
}

func BatchVectorAction(c *gin.Context) {
	req, err := common.ValidateRequest[dto.VectorBatchActionRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = vectorService.BatchVectorAction(req.FileIDs, req.Action)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBUpdateFailed, "批量操作失败: "+err.Error()))
		return
	}

	var message string
	switch req.Action {
	case "reset":
		message = "批量重置向量成功"
	case "retry":
		message = "批量重试向量成功"
	case "delete":
		message = "批量删除向量成功"
	default:
		message = "批量操作完成"
	}

	errors.ResponseSuccess(c, nil, message)
}

func RetryVector(c *gin.Context) {
	req, err := common.ValidateRequest[dto.VectorRetryRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = vectorService.RetryVector(req.FileID)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBUpdateFailed, "重试向量失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "重试向量成功")
}

func GetVectorStats(c *gin.Context) {
	if svc := vectorService.GetGlobalVectorQueueService(); svc != nil {
		stats := svc.GetQueueStats()
		errors.ResponseSuccess(c, stats, "获取统计信息成功")
		return
	}

	// 向量服务未初始化，返回空统计数据
	response := dto.VectorStatsResponse{
		TotalCount:     0,
		PendingCount:   0,
		RunningCount:   0,
		CompletedCount: 0,
		FailedCount:    0,
		ResetCount:     0,
		ActiveWorkers:  0,
		MaxWorkers:     0,
	}

	errors.ResponseSuccess(c, response, "向量服务未启用")
}

func SetAutoProcessing(c *gin.Context) {
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误"))
		return
	}
	if svc := vectorService.GetGlobalVectorQueueService(); svc != nil {
		svc.SetPaused(!body.Enabled)
	}
	_, _ = settingSvc.BatchUpsertSettings(&setdto.BatchUpsertSettingDTO{Settings: []setdto.SettingCreateDTO{{
		Key: "vector_auto_processing_enabled", Value: body.Enabled, Type: "boolean", Group: "vector", Description: "向量队列自动处理开关", IsSystem: true,
	}}})
	errors.ResponseSuccess(c, gin.H{"enabled": body.Enabled}, "已更新向量队列开关")
}

func GetAutoProcessing(c *gin.Context) {
	if svc := vectorService.GetGlobalVectorQueueService(); svc != nil {
		rt := svc.GetQueueStats()["runtime"].(map[string]interface{})
		paused, _ := rt["paused"].(bool)
		errors.ResponseSuccess(c, gin.H{"enabled": !paused}, "获取成功")
		return
	}
	errors.ResponseSuccess(c, gin.H{"enabled": true}, "获取成功")
}

func SetConcurrency(c *gin.Context) {
	var body struct {
		Concurrency int `json:"concurrency"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Concurrency <= 0 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误"))
		return
	}
	if svc := vectorService.GetGlobalVectorQueueService(); svc != nil {
		svc.UpdateConcurrency(body.Concurrency)
	}
	_, _ = settingSvc.BatchUpsertSettings(&setdto.BatchUpsertSettingDTO{Settings: []setdto.SettingCreateDTO{{
		Key: "vector_concurrency", Value: body.Concurrency, Type: "number", Group: "vector", Description: "向量生成并发数量", IsSystem: true,
	}}})
	errors.ResponseSuccess(c, gin.H{"configured": body.Concurrency}, "并发已更新")
}

func GetVectorLogs(c *gin.Context) {
	req, err := common.ValidateRequest[dto.VectorLogListRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	logs, total, err := vectorService.GetVectorLogs(
		req.Page, req.PageSize, req.FileID, req.Action, req.Model)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, "获取日志列表失败: "+err.Error()))
		return
	}

	var logItems []dto.VectorLogItemResponse
	for _, log := range logs {
		var data map[string]interface{}
		if log.Data != "" {
			if err := json.Unmarshal([]byte(log.Data), &data); err != nil {
				data = make(map[string]interface{})
			}
		} else {
			data = make(map[string]interface{})
		}

		logItems = append(logItems, dto.VectorLogItemResponse{
			ID:        log.ID,
			FileID:    log.FileID,
			Action:    log.Action,
			Type:      log.Type,
			Data:      data,
			Message:   log.Message,
			Model:     log.Model,
			Duration:  log.Duration,
			ErrorCode: log.ErrorCode,
			TaskID:    log.TaskID,
			CreatedAt: log.CreatedAt,
		})
	}

	totalPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	response := dto.VectorLogListResponse{
		Data: logItems,
		Pagination: dto.PaginationResponse{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	}

	errors.ResponseSuccess(c, response, "获取日志列表成功")
}

func RegenerateAllVectors(c *gin.Context) {
	req, err := common.ValidateRequest[dto.VectorRegenerateAllRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = vectorService.RegenerateAllVectors(req.Force)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBUpdateFailed, "重新生成向量失败: "+err.Error()))
		return
	}

	message := "重新生成向量任务已启动"
	if req.Force {
		message = "强制重新生成所有向量任务已启动"
	}

	errors.ResponseSuccess(c, nil, message)
}

func RetryAllFailedVectors(c *gin.Context) {
	err := vectorService.RetryAllFailedVectors()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBUpdateFailed, "重试失败向量任务失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "重试失败向量任务已启动")
}

func GetAvailableModels(c *gin.Context) {
	models := vectorService.GetAvailableModels()
	errors.ResponseSuccess(c, gin.H{"models": models}, "获取模型列表成功")
}

func ReconcileMissing(c *gin.Context) {
	var body struct {
		Limit  int  `json:"limit"`
		DryRun bool `json:"dry_run"`
	}
	_ = c.ShouldBindJSON(&body)
	if svc := vectorService.GetGlobalVectorQueueService(); svc != nil {
		total, enq, err := svc.ReconcileMissing(body.Limit, body.DryRun)
		if err != nil {
			errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, err.Error()))
			return
		}
		errors.ResponseSuccess(c, gin.H{"found": total, "enqueued": enq, "dry_run": body.DryRun}, "补齐缺失完成")
		return
	}
	errors.HandleError(c, errors.New(errors.CodeInternal, "服务不可用"))
}

func CleanOrphans(c *gin.Context) {
	var body struct {
		Limit  int  `json:"limit"`
		DryRun bool `json:"dry_run"`
	}
	_ = c.ShouldBindJSON(&body)
	if svc := vectorService.GetGlobalVectorQueueService(); svc != nil {
		total, removed, err := svc.CleanOrphans(body.Limit, body.DryRun)
		if err != nil {
			errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, err.Error()))
			return
		}
		errors.ResponseSuccess(c, gin.H{"found": total, "removed": removed, "dry_run": body.DryRun}, "清理孤儿完成")
		return
	}
	errors.HandleError(c, errors.New(errors.CodeInternal, "服务不可用"))
}

func RebuildStale(c *gin.Context) {
	var body struct {
		Limit int `json:"limit"`
	}
	_ = c.ShouldBindJSON(&body)
	if svc := vectorService.GetGlobalVectorQueueService(); svc != nil {
		enq, err := svc.RebuildStale(body.Limit)
		if err != nil {
			errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, err.Error()))
			return
		}
		errors.ResponseSuccess(c, gin.H{"enqueued": enq}, "重建过期已入队")
		return
	}
	errors.HandleError(c, errors.New(errors.CodeInternal, "服务不可用"))
}

func GetVectorDetail(c *gin.Context) {
	fileID := c.Param("fileId")
	if fileID == "" {
		errors.HandleError(c, errors.NewValidationError("fileId", "文件ID不能为空"))
		return
	}

	results, _, err := vectorService.GetVectorWithFileInfo(1, 1, "", "", fileID)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, "获取向量详情失败: "+err.Error()))
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "向量记录不存在",
		})
		return
	}

	result := results[0]
	item := dto.VectorItemResponse{
		ID:                 result.ID,
		FileID:             result.FileID,
		Status:             result.Status,
		Model:              result.Model,
		Description:        result.Description,
		Dimension:          result.Dimension,
		RetryCount:         result.RetryCount,
		ErrorMessage:       result.ErrorMessage,
		ProcessingDuration: result.ProcessingDuration,
		CreatedAt:          result.CreatedAt,
		UpdatedAt:          result.UpdatedAt,
		LastRetryAt:        result.LastRetryAt,
		FileName:           result.FileName,
		FileURL:            result.FileURL,
		ThumbnailURL:       result.ThumbnailURL,
	}

	errors.ResponseSuccess(c, item, "获取向量详情成功")
}

func RecoverStuckTasks(c *gin.Context) {
	affected, err := vectorService.RecoverStuckRunningTasks()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBUpdateFailed, "恢复卡住任务失败: "+err.Error()))
		return
	}

	if affected == 0 {
		errors.ResponseSuccess(c, gin.H{"affected_count": 0}, "没有发现卡住的任务")
		return
	}

	errors.ResponseSuccess(c, gin.H{"affected_count": affected},
		fmt.Sprintf("成功恢复 %d 个卡住的向量处理任务", affected))
}

func GetQdrantRealStats(c *gin.Context) {
	realStats, err := vectorService.GetQdrantRealStats()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, "获取向量统计信息失败: "+err.Error()))
		return
	}

	// 根据集合状态返回不同的消息
	message := "获取 Qdrant 实际统计信息成功"
	switch realStats.CollectionStatus {
	case "engine_not_initialized":
		message = "向量引擎未初始化，仅显示 MySQL 数据"
	case "qdrant_connection_failed":
		message = "Qdrant 连接失败，仅显示 MySQL 数据"
	case "empty_but_should_have_data":
		message = "Qdrant 数据缺失，请检查向量同步状态"
	}

	errors.ResponseSuccess(c, realStats, message)
}
