package vector

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"

	"gorm.io/gorm"
)

/* VectorVerificationService 向量验证服务 */
type VectorVerificationService struct{}

/* NewVectorVerificationService 创建新的向量验证服务 */
func NewVectorVerificationService() *VectorVerificationService {
	return &VectorVerificationService{}
}

/* CreateTask 创建验证任务 */
func (s *VectorVerificationService) CreateTask(task *models.VectorVerificationTask) error {
	return database.GetDB().Create(task).Error
}

/* GetTaskByID 根据ID获取任务 */
func (s *VectorVerificationService) GetTaskByID(taskID string) (*models.VectorVerificationTask, error) {
	var task models.VectorVerificationTask
	err := database.GetDB().Where("task_id = ?", taskID).First(&task).Error
	return &task, err
}

/* UpdateTask 更新任务 */
func (s *VectorVerificationService) UpdateTask(task *models.VectorVerificationTask) error {
	return database.GetDB().Save(task).Error
}

/* GetRunningTask 获取正在运行的任务 */
func (s *VectorVerificationService) GetRunningTask() (*models.VectorVerificationTask, error) {
	var task models.VectorVerificationTask
	err := database.GetDB().Where("status = ?", models.TaskStatusRunning).First(&task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &task, err
}

/* GetTaskList 获取任务列表 */
func (s *VectorVerificationService) GetTaskList(page, size int, taskType string) ([]models.VectorVerificationTask, int64, error) {
	var tasks []models.VectorVerificationTask
	var total int64

	query := database.GetDB().Model(&models.VectorVerificationTask{}).Preload("Creator")

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

/* UpdateTaskProgress 更新任务进度 */
func (s *VectorVerificationService) UpdateTaskProgress(taskID string, processed, verified, missing, errors int) error {
	return database.GetDB().Model(&models.VectorVerificationTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"processed_count": processed,
			"verified_count":  verified,
			"missing_count":   missing,
			"error_count":     errors,
		}).Error
}
