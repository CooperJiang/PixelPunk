package cron

import (
	"pixelpunk/internal/services/vector"
	"pixelpunk/pkg/logger"
	"time"
)

/* VectorVerificationJob 向量验证定时任务 */
type VectorVerificationJob struct {
	taskService *vector.VectorVerificationTaskService
	schedule    string
}

/* NewVectorVerificationJob 创建向量验证定时任务实例 */
func NewVectorVerificationJob() *VectorVerificationJob {
	return &VectorVerificationJob{
		taskService: vector.NewVectorVerificationTaskService(),
		schedule:    "0 0 0 * * *", // 每天午夜12点执行
	}
}

/* GetSchedule 获取调度配置 */
func (j *VectorVerificationJob) GetSchedule() string {
	return j.schedule
}

/* Execute 执行向量验证任务 */
func (j *VectorVerificationJob) Execute() error {
	logger.Info("开始执行每日向量验证定时任务")

	worker := vector.GetGlobalWorker()
	if worker.IsRunning() {
		logger.Warn("向量验证任务已在运行中，跳过本次调度")
		return nil
	}

	filters := map[string]interface{}{
		"needs_verification": true, // 只验证需要验证的向量（未知状态或超过24小时未验证）
	}

	task, err := j.taskService.CreateTask("scheduled", nil, filters)
	if err != nil {
		logger.Error("创建定时向量验证任务失败: %v", err)
		return err
	}

	if err := vector.StartGlobalVerificationTask(task.TaskID); err != nil {
		logger.Error("启动定时向量验证任务失败: %v", err)
		return err
	}

	logger.Info("定时向量验证任务已启动, 任务ID: %s, 总数: %d", task.TaskID, task.TotalCount)
	return nil
}

/* GetName 获取任务名称 */
func (j *VectorVerificationJob) GetName() string {
	return "Vector Verification Task"
}

/* GetDescription 获取任务描述 */
func (j *VectorVerificationJob) GetDescription() string {
	return "每日定时验证向量数据库与MySQL元数据的一致性"
}

/* GetLastRunTime 获取上次运行时间（可选实现） */
func (j *VectorVerificationJob) GetLastRunTime() *time.Time {
	return nil
}
