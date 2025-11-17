package vector

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"sync"
	"time"

	"gorm.io/gorm"
)

/* VectorVerificationWorker 向量验证工作器 */
type VectorVerificationWorker struct {
	db          *gorm.DB
	taskService *VectorVerificationTaskService
	isRunning   bool
	currentTask *models.VectorVerificationTask
	mutex       sync.RWMutex
	stopChan    chan struct{}
}

/* NewVectorVerificationWorker 创建向量验证工作器实例 */
func NewVectorVerificationWorker() *VectorVerificationWorker {
	return &VectorVerificationWorker{
		db:          database.GetDB(),
		taskService: NewVectorVerificationTaskService(),
		isRunning:   false,
		stopChan:    make(chan struct{}),
	}
}

/* ExecuteTask 异步执行验证任务 */
func (w *VectorVerificationWorker) ExecuteTask(taskID string) error {
	w.mutex.Lock()
	if w.isRunning {
		w.mutex.Unlock()
		return fmt.Errorf("工作器正忙，已有任务在执行中")
	}
	w.isRunning = true
	w.mutex.Unlock()

	go func() {
		defer func() {
			w.mutex.Lock()
			w.isRunning = false
			w.currentTask = nil
			w.mutex.Unlock()
		}()

		if err := w.processTask(taskID); err != nil {
			logger.Error("执行验证任务失败 [%s]: %v", taskID, err)
		}
	}()

	return nil
}

func (w *VectorVerificationWorker) processTask(taskID string) error {
	task := &models.VectorVerificationTask{}
	err := w.db.Where("task_id = ?", taskID).First(task).Error
	if err != nil {
		return fmt.Errorf("获取任务失败: %v", err)
	}

	w.currentTask = task

	task.MarkAsStarted()
	if err := w.db.Save(task).Error; err != nil {
		logger.Error("更新任务状态失败: %v", err)
		return err
	}

	err = w.executeVerification(task)

	if err != nil {
		task.MarkAsFailed(fmt.Sprintf("验证执行失败: %v", err))
		logger.Error("验证任务执行失败 [%s]: %v", taskID, err)
	} else {
		task.MarkAsCompleted()
	}

	if updateErr := w.db.Save(task).Error; updateErr != nil {
		logger.Error("更新任务完成状态失败: %v", updateErr)
	}

	return err
}

func (w *VectorVerificationWorker) executeVerification(task *models.VectorVerificationTask) error {
	filters, err := task.GetFilterConditions()
	if err != nil {
		return fmt.Errorf("获取筛选条件失败: %v", err)
	}

	offset := 0
	batchSize := task.BatchSize

	for {
		select {
		case <-w.stopChan:
			return fmt.Errorf("任务被手动停止")
		default:
		}

		vectors, err := w.getVectorsBatch(offset, batchSize, filters)
		if err != nil {
			return fmt.Errorf("获取向量批次失败: %v", err)
		}

		if len(vectors) == 0 {
			break // 没有更多数据
		}

		batchResult, err := w.processBatch(vectors)
		if err != nil {
			logger.Error("处理批次失败 (offset: %d): %v", offset, err)
		}

		task.ProcessedCount += len(vectors)
		if batchResult != nil {
			task.VerifiedCount += batchResult.VerifiedCount
			task.MissingCount += batchResult.MissingCount
			task.ErrorCount += batchResult.ErrorCount
		}

		if err := w.db.Model(&models.VectorVerificationTask{}).
			Where("task_id = ?", task.TaskID).
			Updates(map[string]interface{}{
				"processed_count": task.ProcessedCount,
				"verified_count":  task.VerifiedCount,
				"missing_count":   task.MissingCount,
				"error_count":     task.ErrorCount,
			}).Error; err != nil {
			logger.Error("更新任务进度失败: %v", err)
		}

		offset += batchSize

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

/* BatchResult 批次处理结果 */
type BatchResult struct {
	ProcessedCount int
	VerifiedCount  int
	MissingCount   int
	ErrorCount     int
}

func (w *VectorVerificationWorker) processBatch(vectors []models.FileVector) (*BatchResult, error) {
	result := &BatchResult{
		ProcessedCount: len(vectors),
	}

	for _, vector := range vectors {
		exists, err := w.taskService.VerifyVectorExists(vector.FileID)

		var status string
		var errorMsg string

		if err != nil {
			status = models.FileVectorActualStatusMissing
			errorMsg = fmt.Sprintf("验证失败，视为缺失: %v", err)
			result.MissingCount++
			logger.Error("验证向量失败，标记为缺失 [%s]: %v", vector.FileID, err)
		} else if exists {
			status = models.FileVectorActualStatusVerified
			result.VerifiedCount++
		} else {
			status = models.FileVectorActualStatusMissing
			result.MissingCount++
		}

		vector.SetVerificationStatus(status, errorMsg)

		if err := w.db.Save(&vector).Error; err != nil {
			logger.Error("保存向量验证状态失败 [%s]: %v", vector.FileID, err)
			result.ErrorCount++
		}
	}

	return result, nil
}

func (w *VectorVerificationWorker) getVectorsBatch(offset, limit int, filters map[string]interface{}) ([]models.FileVector, error) {
	var vectors []models.FileVector

	query := w.db.Model(&models.FileVector{})

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

	err := query.Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&vectors).Error

	return vectors, err
}

/* IsRunning 检查工作器是否正在运行 */
func (w *VectorVerificationWorker) IsRunning() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.isRunning
}

/* GetCurrentTask 获取当前正在执行的任务 */
func (w *VectorVerificationWorker) GetCurrentTask() *models.VectorVerificationTask {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.currentTask
}

/* Stop 停止当前任务 */
func (w *VectorVerificationWorker) Stop() error {
	w.mutex.RLock()
	if !w.isRunning {
		w.mutex.RUnlock()
		return fmt.Errorf("没有正在运行的任务")
	}
	w.mutex.RUnlock()

	select {
	case w.stopChan <- struct{}{}:
		return nil
	default:
		return fmt.Errorf("停止信号发送失败")
	}
}

/* GetProgress 获取当前任务进度 */
func (w *VectorVerificationWorker) GetProgress() (float64, *models.TaskSummary, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	if !w.isRunning || w.currentTask == nil {
		return 0.0, nil, fmt.Errorf("没有正在运行的任务")
	}

	progress := w.currentTask.GetProgress()
	summary := w.currentTask.ToSummary()

	return progress, summary, nil
}

var (
	globalWorker *VectorVerificationWorker
	workerOnce   sync.Once
)

/* GetGlobalWorker 获取全局工作器实例 */
func GetGlobalWorker() *VectorVerificationWorker {
	workerOnce.Do(func() {
		globalWorker = NewVectorVerificationWorker()
	})
	return globalWorker
}

/* StartGlobalVerificationTask 启动全局验证任务 */
func StartGlobalVerificationTask(taskID string) error {
	worker := GetGlobalWorker()
	return worker.ExecuteTask(taskID)
}
