package admin

import (
	"fmt"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/vector"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type VectorVerificationController struct {
	taskService *vector.VectorVerificationTaskService
	worker      *vector.VectorVerificationWorker
}

func NewVectorVerificationController() *VectorVerificationController {
	return &VectorVerificationController{
		taskService: vector.NewVectorVerificationTaskService(),
		worker:      vector.GetGlobalWorker(),
	}
}

type StartVerificationTaskDTO struct {
	TaskType          string   `json:"task_type" binding:"required,oneof=manual scheduled partial"`
	ActualStatus      string   `json:"actual_status" binding:"omitempty,oneof=unknown verified missing"`
	NeedsVerification bool     `json:"needs_verification"`
	ForceFullCheck    bool     `json:"force_full_check"`   // 强制全量检查，忽略验证状态
	FileIDs           []string `json:"file_ids,omitempty"` // 部分验证时指定的文件ID
}

func (dto *StartVerificationTaskDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"TaskType.required":  "任务类型不能为空",
		"TaskType.oneof":     "任务类型必须是 manual、scheduled 或 partial",
		"ActualStatus.oneof": "状态必须是 unknown、verified 或 missing",
	}
}

func (c *VectorVerificationController) StartVerificationTask(ctx *gin.Context) {
	req, err := common.ValidateRequest[StartVerificationTaskDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	userID := middleware.GetCurrentUserID(ctx)

	filters := make(map[string]interface{})
	if req.ActualStatus != "" {
		filters["actual_status"] = req.ActualStatus
	}
	if req.NeedsVerification {
		filters["needs_verification"] = true
	}
	if req.ForceFullCheck {
		filters["force_full_check"] = true
	}
	if len(req.FileIDs) > 0 {
		filters["file_ids"] = req.FileIDs
	}

	task, err := c.taskService.CreateTask(req.TaskType, &userID, filters)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	if err := vector.StartGlobalVerificationTask(task.TaskID); err != nil {
		logger.Error("启动验证任务失败: %v", err)
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "启动验证任务失败"))
		return
	}

	errors.ResponseSuccess(ctx, gin.H{
		"task_id":     task.TaskID,
		"total_count": task.TotalCount,
		"message":     "验证任务已启动，正在后台处理中",
	}, "任务创建成功")
}

func (c *VectorVerificationController) GetTaskStatus(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	if taskID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "任务ID不能为空"))
		return
	}

	status, err := c.taskService.GetTaskStatus(taskID)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "获取任务状态失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(ctx, status, "获取成功")
}

func (c *VectorVerificationController) GetTaskList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	taskType := ctx.Query("task_type")

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	tasks, total, err := c.taskService.GetTaskList(page, size, taskType)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "获取任务列表失败: "+err.Error()))
		return
	}

	data := gin.H{
		"items": tasks,
		"pagination": gin.H{
			"total":        total,
			"size":         size,
			"current_page": page,
			"last_page":    (total + int64(size) - 1) / int64(size),
		},
	}

	errors.ResponseSuccess(ctx, data, "获取成功")
}

func (c *VectorVerificationController) GetVerificationStatistics(ctx *gin.Context) {
	stats, err := c.taskService.GetVerificationStatistics()
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "获取统计信息失败: "+err.Error()))
		return
	}

	if c.worker.IsRunning() {
		progress, currentTask, _ := c.worker.GetProgress()
		stats.RunningTask = &struct {
			TaskID   string  `json:"task_id"`
			Progress float64 `json:"progress"`
		}{
			TaskID:   currentTask.TaskID,
			Progress: progress,
		}
	}

	errors.ResponseSuccess(ctx, stats, "获取成功")
}

type RepairMissingVectorDTO struct {
	FileID string `json:"file_id" binding:"required"`
}

func (dto *RepairMissingVectorDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileID.required": "文件ID不能为空",
	}
}

func (c *VectorVerificationController) RepairMissingVector(ctx *gin.Context) {
	req, err := common.ValidateRequest[RepairMissingVectorDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	err = c.taskService.RepairMissingVector(req.FileID)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "修复向量失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(ctx, nil, "向量修复请求已提交，系统将重新处理该文件")
}


func (c *VectorVerificationController) BatchRepairMissingVectors(ctx *gin.Context) {
	affected, err := c.taskService.BatchRepairAllProblemVectors()
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "批量修复向量失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(ctx, gin.H{
		"affected_count": affected,
	}, fmt.Sprintf("成功提交 %d 个向量修复请求，系统将重新处理这些文件的向量化", affected))
}

func (c *VectorVerificationController) CancelTask(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	if taskID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "任务ID不能为空"))
		return
	}

	err := c.taskService.CancelTask(taskID)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "取消任务失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(ctx, nil, "任务已取消")
}

func (c *VectorVerificationController) GetRunningTaskProgress(ctx *gin.Context) {
	if !c.worker.IsRunning() {
		errors.ResponseSuccess(ctx, gin.H{
			"is_running": false,
			"progress":   0.0,
		}, "当前没有正在运行的任务")
		return
	}

	progress, summary, err := c.worker.GetProgress()
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "获取任务进度失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(ctx, gin.H{
		"is_running": true,
		"progress":   progress,
		"task":       summary,
	}, "获取成功")
}

func (c *VectorVerificationController) VerifySpecificVectors(ctx *gin.Context) {
	fileIDsStr := ctx.Param("fileIds")
	if fileIDsStr == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "文件ID不能为空"))
		return
	}

	fileIDs := strings.Split(fileIDsStr, ",")
	if len(fileIDs) == 0 {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "文件ID列表不能为空"))
		return
	}

	if len(fileIDs) > 50 {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "一次最多只能验证50个文件"))
		return
	}

	userID := middleware.GetCurrentUserID(ctx)

	filters := map[string]interface{}{
		"file_ids": fileIDs,
	}

	task, err := c.taskService.CreateTask("partial", &userID, filters)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	if err := vector.StartGlobalVerificationTask(task.TaskID); err != nil {
		logger.Error("启动指定向量验证任务失败: %v", err)
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "启动验证任务失败"))
		return
	}

	errors.ResponseSuccess(ctx, gin.H{
		"task_id":     task.TaskID,
		"total_count": task.TotalCount,
		"file_ids":    fileIDs,
	}, "指定向量验证任务已启动")
}

type SmartFixRequest struct {
	IncludeFailedTasks    bool `json:"include_failed_tasks"`    // 包含失败任务
	IncludeStuckTasks     bool `json:"include_stuck_tasks"`     // 包含卡住任务
	IncludeMissingVectors bool `json:"include_missing_vectors"` // 包含缺失向量
	IncludeUnknownVectors bool `json:"include_unknown_vectors"` // 包含未知状态向量
}

func (dto *SmartFixRequest) GetValidationMessages() map[string]string {
	return map[string]string{}
}

func (c *VectorVerificationController) SmartFix(ctx *gin.Context) {
	req, err := common.ValidateRequest[SmartFixRequest](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	fixResults := gin.H{
		"total_fixed":     0,
		"failed_tasks":    0,
		"stuck_tasks":     0,
		"missing_vectors": 0,
		"unknown_vectors": 0,
		"errors":          []string{},
	}

	if req.IncludeFailedTasks {
		if fixedCount, err := c.fixFailedTasks(); err != nil {
			logger.Error("修复失败任务出错: %v", err)
			fixResults["errors"] = append(fixResults["errors"].([]string), "修复失败任务失败: "+err.Error())
		} else {
			fixResults["failed_tasks"] = fixedCount
			fixResults["total_fixed"] = fixResults["total_fixed"].(int) + fixedCount
		}
	}

	if req.IncludeStuckTasks {
		if fixedCount, err := c.fixStuckTasks(); err != nil {
			logger.Error("修复卡住任务出错: %v", err)
			fixResults["errors"] = append(fixResults["errors"].([]string), "修复卡住任务失败: "+err.Error())
		} else {
			fixResults["stuck_tasks"] = fixedCount
			fixResults["total_fixed"] = fixResults["total_fixed"].(int) + fixedCount
		}
	}

	if req.IncludeMissingVectors || req.IncludeUnknownVectors {
		if fixedCount, err := c.fixMissingAndUnknownVectors(req.IncludeMissingVectors, req.IncludeUnknownVectors); err != nil {
			logger.Error("修复缺失/未知向量出错: %v", err)
			fixResults["errors"] = append(fixResults["errors"].([]string), "修复缺失/未知向量失败: "+err.Error())
		} else {
			if req.IncludeMissingVectors {
				fixResults["missing_vectors"] = fixedCount
			}
			if req.IncludeUnknownVectors {
				fixResults["unknown_vectors"] = fixedCount
			}
			fixResults["total_fixed"] = fixResults["total_fixed"].(int) + fixedCount
		}
	}

	totalFixed := fixResults["total_fixed"].(int)
	var message string
	if totalFixed > 0 {
		message = fmt.Sprintf("智能修复已完成，共修复 %d 个问题", totalFixed)
	} else {
		message = "智能修复已完成，未发现需要修复的问题"
	}

	errorList := fixResults["errors"].([]string)
	if len(errorList) > 0 {
		message += fmt.Sprintf("，但部分操作失败（%d个错误）", len(errorList))
	}

	errors.ResponseSuccess(ctx, fixResults, message)
}

func (c *VectorVerificationController) fixFailedTasks() (int, error) {
	err := vector.RetryAllFailedVectors()
	if err != nil {
		return 0, fmt.Errorf("重试失败任务失败: %w", err)
	}

	stats, err := vector.GetVectorStats()
	if err != nil {
		return 0, fmt.Errorf("获取统计信息失败: %w", err)
	}

	failedCount := 0
	if count, ok := stats["failed_count"].(int64); ok {
		failedCount = int(count)
	}

	return failedCount, nil
}

func (c *VectorVerificationController) fixStuckTasks() (int, error) {
	affectedCount, err := vector.RecoverStuckRunningTasks()
	if err != nil {
		return 0, fmt.Errorf("恢复卡住任务失败: %w", err)
	}

	return int(affectedCount), nil
}

func (c *VectorVerificationController) fixMissingAndUnknownVectors(includeMissing, includeUnknown bool) (int, error) {
	db := database.GetDB()
	if db == nil {
		return 0, fmt.Errorf("数据库连接不可用")
	}

	var conditions []string
	var args []interface{}

	if includeMissing {
		conditions = append(conditions, "actual_status = ?")
		args = append(args, "missing")
	}

	if includeUnknown {
		conditions = append(conditions, "actual_status = ?")
		args = append(args, "unknown")
	}

	if len(conditions) == 0 {
		return 0, nil // 没有要修复的
	}

	query := fmt.Sprintf("status = 'completed' AND (%s)", strings.Join(conditions, " OR "))

	var vectors []models.FileVector
	if err := db.Where(query, args...).Find(&vectors).Error; err != nil {
		return 0, fmt.Errorf("查询需要修复的向量失败: %w", err)
	}

	if len(vectors) == 0 {
		return 0, nil // 没有需要修复的向量
	}

	var fileIDs []string
	for _, v := range vectors {
		fileIDs = append(fileIDs, v.FileID)
	}

	updateData := map[string]interface{}{
		"status":             common.VectorStatusPending,
		"actual_status":      "unknown",
		"verification_error": "",
		"last_verified":      nil,
	}

	err := db.Model(&models.FileVector{}).
		Where("file_id IN ?", fileIDs).
		Updates(updateData).Error

	if err != nil {
		return 0, fmt.Errorf("重置向量状态失败: %w", err)
	}

	return len(vectors), nil
}
