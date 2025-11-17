package vector

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/vector"
	"time"

	"gorm.io/gorm"
)

/* VectorVerificationTaskService 向量验证任务服务 */
type VectorVerificationTaskService struct {
	db           *gorm.DB
	vectorEngine *vector.VectorEngine
}

/* NewVectorVerificationTaskService 创建向量验证任务服务实例 */
func NewVectorVerificationTaskService() *VectorVerificationTaskService {
	return &VectorVerificationTaskService{
		db:           database.GetDB(),
		vectorEngine: vector.GetGlobalVectorEngine(),
	}
}

/* CreateTask 创建验证任务 */
func (s *VectorVerificationTaskService) CreateTask(taskType string, creatorID *uint, filters map[string]interface{}) (*models.VectorVerificationTask, error) {
	if s.IsTaskRunning() {
		runningTask, _ := s.GetRunningTask()
		return nil, fmt.Errorf("已有任务正在运行中 (任务ID: %s)，请等待完成后再试", runningTask.TaskID)
	}

	totalCount, err := s.countVectorsToVerify(filters)
	if err != nil {
		return nil, fmt.Errorf("统计向量数量失败: %v", err)
	}

	if totalCount == 0 {
		return nil, fmt.Errorf("没有找到需要验证的向量数据")
	}

	task := &models.VectorVerificationTask{
		TaskType:   taskType,
		CreatorID:  creatorID,
		TotalCount: totalCount,
		BatchSize:  100,
	}

	if err := task.SetFilterConditions(filters); err != nil {
		return nil, fmt.Errorf("设置筛选条件失败: %v", err)
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, fmt.Errorf("创建任务失败: %v", err)
	}

	return task, nil
}

/* GetTaskStatus 获取任务状态 */
func (s *VectorVerificationTaskService) GetTaskStatus(taskID string) (*models.TaskSummary, error) {
	var task models.VectorVerificationTask
	err := s.db.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, fmt.Errorf("任务不存在: %v", err)
	}

	return task.ToSummary(), nil
}

/* GetTaskList 获取任务列表 */
func (s *VectorVerificationTaskService) GetTaskList(page, size int, taskType string) ([]models.VectorVerificationTask, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	var tasks []models.VectorVerificationTask
	var total int64

	query := s.db.Model(&models.VectorVerificationTask{}).Preload("Creator")

	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&tasks).Error

	return tasks, total, err
}

/* IsTaskRunning 检查是否有任务正在运行 */
func (s *VectorVerificationTaskService) IsTaskRunning() bool {
	var task models.VectorVerificationTask
	err := s.db.Where("status = ?", models.TaskStatusRunning).First(&task).Error
	return err == nil
}

/* GetRunningTask 获取正在运行的任务 */
func (s *VectorVerificationTaskService) GetRunningTask() (*models.VectorVerificationTask, error) {
	var task models.VectorVerificationTask
	err := s.db.Where("status = ?", models.TaskStatusRunning).First(&task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &task, err
}

/* VerifyVectorExists 验证单个向量是否存在 */
func (s *VectorVerificationTaskService) VerifyVectorExists(fileID string) (bool, error) {
	if s.vectorEngine == nil {
		return false, fmt.Errorf("向量引擎未初始化")
	}

	exists, err := s.vectorEngine.VectorExists(fileID)
	if err != nil {
		logger.Error("验证向量存在性失败 [%s]: %v", fileID, err)
		return false, err
	}

	return exists, nil
}

/* GetVerificationStatistics 获取验证统计信息 */
func (s *VectorVerificationTaskService) GetVerificationStatistics() (*VerificationStatistics, error) {
	// 数据库未初始化时返回空统计
	if s.db == nil {
		return &VerificationStatistics{
			Total:            0,
			Verified:         0,
			Missing:          0,
			Unknown:          0,
			LastVerification: nil,
			RunningTask:      nil,
		}, nil
	}

	stats := &VerificationStatistics{}

	err := s.db.Model(&models.FileVector{}).
		Select(`
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN actual_status = 'verified' THEN 1 ELSE 0 END), 0) as verified,
			COALESCE(SUM(CASE WHEN actual_status = 'missing' THEN 1 ELSE 0 END), 0) as missing,
			COALESCE(SUM(CASE WHEN actual_status = 'unknown' THEN 1 ELSE 0 END), 0) as unknown
		`).
		Where("status = ?", common.VectorStatusCompleted).
		Row().
		Scan(&stats.Total, &stats.Verified, &stats.Missing, &stats.Unknown)

	if err != nil {
		return nil, err
	}

	var lastTask models.VectorVerificationTask
	err = s.db.Where("status IN ?", []string{models.TaskStatusCompleted}).
		Order("completed_at DESC").
		First(&lastTask).Error

	if err == nil {
		stats.LastVerification = lastTask.CompletedAt
	}

	return stats, nil
}

func (s *VectorVerificationTaskService) countVectorsToVerify(filters map[string]interface{}) (int, error) {
	var count int64

	query := s.db.Model(&models.FileVector{})

	if filters != nil {
		if imageIDs, ok := filters["file_ids"].([]string); ok && len(imageIDs) > 0 {
			query = query.Where("file_id IN ?", imageIDs)
		} else {
			if forceFullCheck, ok := filters["force_full_check"].(bool); ok && forceFullCheck {
			} else if needsVerification, ok := filters["needs_verification"].(bool); ok && needsVerification {
				query = query.Where(`
					status = ? OR 
					actual_status IN (?, ?) OR
					last_verified IS NULL OR 
					last_verified < ?
				`, common.VectorStatusCompleted,
					models.FileVectorActualStatusUnknown,
					models.FileVectorActualStatusMissing,
					time.Now().Add(-24*time.Hour))
			} else {
				query = query.Where("status = ?", common.VectorStatusCompleted)
			}
		}

		if actualStatus, ok := filters["actual_status"]; ok {
			query = query.Where("actual_status = ?", actualStatus)
		}
	} else {
		query = query.Where("status = ?", common.VectorStatusCompleted)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计向量数量失败: %v", err)
	}

	return int(count), err
}

/* VerificationStatistics 验证统计信息 */
type VerificationStatistics struct {
	Total            int         `json:"total_vectors"`
	Verified         int         `json:"verified_count"`
	Missing          int         `json:"missing_count"`
	Unknown          int         `json:"unknown_count"`
	LastVerification *time.Time  `json:"last_verification"`
	RunningTask      interface{} `json:"running_task,omitempty"`
}

/* GetVerificationProgress 获取验证进度（用于实时显示） */
func (s *VectorVerificationTaskService) GetVerificationProgress() (float64, *models.VectorVerificationTask, error) {
	task, err := s.GetRunningTask()
	if err != nil || task == nil {
		return 0.0, nil, nil
	}

	progress := task.GetProgress()
	return progress, task, nil
}

/* CancelTask 取消任务（如果支持的话） */
func (s *VectorVerificationTaskService) CancelTask(taskID string) error {
	var task models.VectorVerificationTask
	err := s.db.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		return fmt.Errorf("任务不存在: %v", err)
	}

	if !task.IsRunning() {
		return fmt.Errorf("只能取消正在运行的任务")
	}

	task.MarkAsFailed("用户手动取消")

	return s.db.Save(&task).Error
}

/* RepairMissingVector 修复缺失的向量（重新向量化） */
func (s *VectorVerificationTaskService) RepairMissingVector(fileID string) error {
	var imageVector models.FileVector
	err := s.db.Where("file_id = ?", fileID).First(&imageVector).Error
	if err != nil {
		return fmt.Errorf("向量记录不存在: %v", err)
	}

	if !imageVector.IsActualMissing() {
		return fmt.Errorf("向量状态不是缺失状态，无需修复")
	}

	imageVector.Status = common.VectorStatusPending
	imageVector.ActualStatus = models.FileVectorActualStatusUnknown
	imageVector.VerificationError = ""

	err = s.db.Save(&imageVector).Error
	if err != nil {
		return fmt.Errorf("重置向量状态失败: %v", err)
	}

	// 主动触发入队（修复：自动将任务加入队列）
	if queueService := GetGlobalVectorQueueService(); queueService != nil {
		if err := queueService.EnqueueVector(fileID); err != nil {
			logger.Warn("向量 %s 入队失败: %v", fileID, err)
		} else {
		}
	}

	return nil
}

/* BatchRepairMissingVectors 批量修复缺失的向量 */
func (s *VectorVerificationTaskService) BatchRepairMissingVectors(fileIDs []string) (int, error) {
	if len(fileIDs) == 0 {
		return 0, fmt.Errorf("没有指定要修复的向量ID")
	}

	result := s.db.Model(&models.FileVector{}).
		Where("file_id IN ? AND actual_status = ?", fileIDs, models.FileVectorActualStatusMissing).
		Updates(map[string]interface{}{
			"status":             common.VectorStatusPending,
			"actual_status":      models.FileVectorActualStatusUnknown,
			"verification_error": "",
		})

	if result.Error != nil {
		return 0, fmt.Errorf("批量修复失败: %v", result.Error)
	}

	affected := int(result.RowsAffected)

	return affected, nil
}

/* BatchRepairAllProblemVectors 批量修复所有问题向量（缺失和未知状态） */
func (s *VectorVerificationTaskService) BatchRepairAllProblemVectors() (int, error) {
	// 先查询需要修复的文件ID
	var fileIDs []string
	err := s.db.Model(&models.FileVector{}).
		Where("status = ? AND actual_status IN ?",
			common.VectorStatusCompleted,
			[]string{models.FileVectorActualStatusMissing, models.FileVectorActualStatusUnknown}).
		Pluck("file_id", &fileIDs).Error

	if err != nil {
		return 0, fmt.Errorf("查询问题向量失败: %v", err)
	}

	if len(fileIDs) == 0 {
		return 0, nil
	}

	result := s.db.Model(&models.FileVector{}).
		Where("file_id IN ?", fileIDs).
		Updates(map[string]interface{}{
			"status":             common.VectorStatusPending,
			"actual_status":      models.FileVectorActualStatusUnknown,
			"verification_error": "",
			"last_verified":      nil, // 清除上次验证时间
		})

	if result.Error != nil {
		return 0, fmt.Errorf("批量修复失败: %v", result.Error)
	}

	affected := int(result.RowsAffected)

	// 主动触发批量入队（修复：自动将任务加入队列）
	if queueService := GetGlobalVectorQueueService(); queueService != nil {
		enqueuedCount := 0
		for _, fileID := range fileIDs {
			if err := queueService.EnqueueVector(fileID); err == nil {
				enqueuedCount++
			}
		}
	} else {
		logger.Warn("向量队列服务未初始化，已标记为待处理但未入队，请启用向量服务自动处理")
	}

	return affected, nil
}
