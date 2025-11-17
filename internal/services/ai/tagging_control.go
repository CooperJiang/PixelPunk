package ai

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/logger"
)

func (s *TaggingService) ProcessSingleFile(file models.File) {
	if s.stopping {
		logger.Warn("服务正在停止，跳过入队: %s", file.ID)
		return
	}
	if s.IsPaused() {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("ProcessSingleFile panic: %v, 文件ID: %s", r, file.ID)
			}
		}()
		if s.taskQueue == nil {
			logger.Error("队列后端未初始化，无法入队: %s", file.ID)
			return
		}
		if err := s.taskQueue.EnqueueUnique(file.ID, 0); err != nil {
			logger.Warn("入队失败(%s): %v", file.ID, err)
			return
		}
		s.notifyQueueStatsChange()
	}()
}

func (s *TaggingService) BatchProcessFiles(files []models.File) { s.BatchProcessFilesWithResult(files) }

func (s *TaggingService) BatchProcessFilesWithResult(files []models.File) (int, int) {
	if len(files) == 0 {
		return 0, 0
	}
	if s.stopping {
		logger.Warn("服务正在停止，跳过批量处理")
		return 0, len(files)
	}
	if s.IsPaused() {
		return 0, len(files)
	}
	if s.taskQueue == nil {
		logger.Error("队列后端未初始化")
		return 0, len(files)
	}

	enqueued, skipped := 0, 0
	for _, file := range files {
		if err := s.taskQueue.EnqueueUnique(file.ID, 0); err == nil {
			enqueued++
		} else {
			skipped++
		}
	}
	s.notifyQueueStatsChange()
	return enqueued, skipped
}

func (s *TaggingService) IsProcessing() bool {
	s.mutex.Lock()
	active := s.activeWorkers
	s.mutex.Unlock()
	return active > 0
}

func (s *TaggingService) Stop() {
	s.mutex.Lock()
	s.stopping = true
	s.mutex.Unlock()

	if s.pipeline != nil {
		s.pipeline.Stop()
	}

	if s.reaperStop != nil {
		close(s.reaperStop)
	}
}

func (s *TaggingService) GetQueueStats() map[string]interface{} {
	s.mutex.Lock()
	activeWorkers := s.activeWorkers
	concurrentNum := s.concurrentNum
	queueLength := 0
	paused := s.paused
	s.mutex.Unlock()

	return map[string]interface{}{
		"active_workers":    activeWorkers,
		"max_workers":       concurrentNum,
		"queue_length":      queueLength,
		"available_workers": concurrentNum - activeWorkers,
		"paused":            paused,
		"recent_failures":   s.RecentFailures(),
	}
}

func (s *TaggingService) UpdateConcurrency(newConcurrency int) error {
	if newConcurrency <= 0 {
		return fmt.Errorf("并发数必须大于0")
	}

	s.mutex.Lock()
	oldConcurrency := s.concurrentNum
	if oldConcurrency == newConcurrency {
		s.mutex.Unlock()
		return nil
	}
	s.concurrentNum = newConcurrency
	s.mutex.Unlock()

	if s.pipeline != nil {
		s.pipeline.AdjustConcurrency(newConcurrency)
	}

	s.notifyQueueStatsChange()

	return nil
}

func (s *TaggingService) RefreshConcurrencyFromConfig() error {
	newConcurrency := getAIConcurrency()
	return s.UpdateConcurrency(newConcurrency)
}

func (s *TaggingService) NotifyQueueStatsChange() {
	s.notifyQueueStatsChange()
}
