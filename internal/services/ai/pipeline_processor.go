package ai

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"pixelpunk/internal/models"
	qqueue "pixelpunk/internal/queue"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/logger"
	ai "pixelpunk/pkg/ai"

	"gorm.io/gorm"
)

type DynamicSemaphore struct {
	current int
	max     int
	mu      sync.Mutex
	cond    *sync.Cond
}

func NewDynamicSemaphore(max int) *DynamicSemaphore {
	s := &DynamicSemaphore{
		current: 0,
		max:     max,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *DynamicSemaphore) Acquire() {
	s.mu.Lock()
	for s.current >= s.max {
		s.cond.Wait()
	}
	s.current++
	s.mu.Unlock()
}

func (s *DynamicSemaphore) Release() {
	s.mu.Lock()
	s.current--
	s.cond.Signal()
	s.mu.Unlock()
}

func (s *DynamicSemaphore) SetMax(max int) {
	s.mu.Lock()
	oldMax := s.max
	s.max = max
	s.mu.Unlock()

	if max > oldMax {
		s.cond.Broadcast()
	}

	logger.Info("信号量容量调整: %d -> %d", oldMax, max)
}

func (s *DynamicSemaphore) GetCurrent() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

func (s *DynamicSemaphore) GetMax() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.max
}

type PipelineProcessor struct {
	service *TaggingService

	taskChan   chan *QueueTask
	fileChan   chan *FileTask
	resultChan chan *ProcessResult

	aiSemaphore *DynamicSemaphore
	stopping    bool
	stopMu      sync.Mutex
	wg          sync.WaitGroup

	fileLoaderWg  sync.WaitGroup
	aiProcessorWg sync.WaitGroup

	fileLoaderPoolSize  int
	aiProcessorPoolSize int
	dbSaverPoolSize     int
	batchSize           int
	batchTimeout        time.Duration
}

type QueueTask struct {
	FileID string
	Ack    qqueue.AckFunc
	Nack   qqueue.NackFunc
}

type FileTask struct {
	File        models.File
	Base64Data  string
	ImageFormat string
	Ack         qqueue.AckFunc
	Nack        qqueue.NackFunc
}

type ProcessResult struct {
	FileID         string
	Success        bool
	CategoryID     *uint
	AIResponse     *AIFileResponse
	Error          error
	Ack            qqueue.AckFunc
	Nack           qqueue.NackFunc
	ProcessedAt    time.Time
	CategoryResult *CategoryResult
	HttpDuration   int64 // HTTP调用耗时（毫秒）
}

type CategoryResult struct {
	CategoryID          uint
	CategoryName        string
	CategoryDescription string
}

func NewPipelineProcessor(service *TaggingService, aiConcurrency int) *PipelineProcessor {
	if aiConcurrency <= 0 {
		aiConcurrency = 50
	}

	return &PipelineProcessor{
		service:             service,
		taskChan:            make(chan *QueueTask, 200),
		fileChan:            make(chan *FileTask, 200),
		resultChan:          make(chan *ProcessResult, 200),
		aiSemaphore:         NewDynamicSemaphore(aiConcurrency),
		fileLoaderPoolSize:  10,
		aiProcessorPoolSize: aiConcurrency,
		dbSaverPoolSize:     5,
		batchSize:           10,
		batchTimeout:        100 * time.Millisecond,
	}
}

func (pp *PipelineProcessor) Start() {

	pp.wg.Add(1)
	go pp.taskFetcher()

	for i := 0; i < pp.fileLoaderPoolSize; i++ {
		pp.wg.Add(1)
		pp.fileLoaderWg.Add(1)
		go pp.fileLoader(i)
	}

	for i := 0; i < pp.aiProcessorPoolSize; i++ {
		pp.wg.Add(1)
		pp.aiProcessorWg.Add(1)
		go pp.aiProcessor(i)
	}

	for i := 0; i < pp.dbSaverPoolSize; i++ {
		pp.wg.Add(1)
		go pp.dbSaver(i)
	}

	go pp.channelCloser()
}

func (pp *PipelineProcessor) Stop() {
	pp.stopMu.Lock()
	if pp.stopping {
		pp.stopMu.Unlock()
		return
	}
	pp.stopping = true
	pp.stopMu.Unlock()

	close(pp.taskChan)
	pp.wg.Wait()

	logger.Info("Pipeline已停止")
}

func (pp *PipelineProcessor) IsStopping() bool {
	pp.stopMu.Lock()
	defer pp.stopMu.Unlock()
	return pp.stopping
}

func (pp *PipelineProcessor) AdjustConcurrency(newConcurrency int) {
	if newConcurrency <= 0 {
		return
	}

	oldConcurrency := pp.aiSemaphore.GetMax()
	if oldConcurrency == newConcurrency {
		return
	}

	pp.aiSemaphore.SetMax(newConcurrency)
	pp.aiProcessorPoolSize = newConcurrency

	logger.Info("AI并发数调整: %d -> %d", oldConcurrency, newConcurrency)
}

func (pp *PipelineProcessor) channelCloser() {
	pp.fileLoaderWg.Wait()
	close(pp.fileChan)

	pp.aiProcessorWg.Wait()
	close(pp.resultChan)
}

func (pp *PipelineProcessor) taskFetcher() {
	defer pp.wg.Done()

	for !pp.IsStopping() {
		if pp.service.IsPaused() {
			time.Sleep(1 * time.Second)
			continue
		}

		task, ack, nack, err := pp.service.taskQueue.Fetch(30 * time.Second)

		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		queueTask := &QueueTask{
			FileID: task.FileID,
			Ack:    ack,
			Nack:   nack,
		}

		// 在发送前检查是否正在停止，避免向已关闭的 channel 发送数据
		if pp.IsStopping() {
			nack(1*time.Second, false, "pipeline is stopping")
			continue
		}

		select {
		case pp.taskChan <- queueTask:
		case <-time.After(5 * time.Second):
			nack(1*time.Second, false, "pipeline task channel full")
		}
	}
}

func (pp *PipelineProcessor) fileLoader(workerID int) {
	defer pp.wg.Done()
	defer pp.fileLoaderWg.Done()

	for task := range pp.taskChan {
		var file models.File
		if err := pp.service.db.Where("id = ?", task.FileID).Take(&file).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				task.Ack()
			} else {
				task.Nack(1*time.Second, false, fmt.Sprintf("load file failed: %v", err))
			}
			continue
		}

		base64Data, imageFormat, err := pp.service.readImageAsBase64(file)
		if err != nil {
			if errors.Is(err, errMissingFile) {
				policy := getMissingFilePolicy()
				switch policy {
				case "drop":
					_ = pp.service.db.Model(&models.File{}).Where("id = ?", file.ID).
						Update("ai_tagging_status", common.AITaggingStatusIgnored).Error
					task.Ack()
				case "ignore":
					_ = pp.service.db.Model(&models.File{}).Where("id = ?", file.ID).
						Updates(map[string]interface{}{
							"ai_tagging_status": common.AITaggingStatusNone,
							"ai_tagging_tries":  0,
						}).Error
					task.Ack()
				default:
					task.Nack(3*time.Second, false, fmt.Sprintf("missing file: %v", err))
				}
			} else {
				task.Nack(1*time.Second, false, fmt.Sprintf("read image failed: %v", err))
			}
			continue
		}

		fileTask := &FileTask{
			File:        file,
			Base64Data:  base64Data,
			ImageFormat: imageFormat,
			Ack:         task.Ack,
			Nack:        task.Nack,
		}

		pp.fileChan <- fileTask
	}
}

func (pp *PipelineProcessor) aiProcessor(workerID int) {
	defer pp.wg.Done()
	defer pp.aiProcessorWg.Done()

	for fileTask := range pp.fileChan {
		pp.aiSemaphore.Acquire()

		_ = pp.service.db.Model(&models.File{}).
			Where("id = ?", fileTask.File.ID).
			Update("ai_last_heartbeat_at", time.Now()).Error

		result := &ProcessResult{
			FileID:      fileTask.File.ID,
			Ack:         fileTask.Ack,
			Nack:        fileTask.Nack,
			ProcessedAt: time.Now(),
		}

		err := pp.processImageWithAI(fileTask, result)

		pp.aiSemaphore.Release()

		if err != nil {
			result.Success = false
			result.Error = err
		} else {
			result.Success = true
		}

		pp.resultChan <- result
	}
}

func (pp *PipelineProcessor) processImageWithAI(fileTask *FileTask, result *ProcessResult) error {
	categoryResult, err := performAIImageCategorizationOutsideTx(
		fileTask.File,
		fileTask.Base64Data,
		fileTask.ImageFormat,
	)
	if err != nil {
		return fmt.Errorf("AI分类失败: %v", err)
	}

	categoryName, categoryDescription, categoryID := buildTaggingContext(categoryResult)

	if categoryResult != nil && categoryResult.Success {
		result.CategoryResult = &CategoryResult{
			CategoryID:          categoryResult.CategoryID,
			CategoryName:        categoryResult.CategoryName,
			CategoryDescription: categoryResult.CategoryDescription,
		}
		result.CategoryID = categoryID
	}

	aiResponse, err := performAITagging(
		fileTask.File,
		fileTask.Base64Data,
		fileTask.ImageFormat,
		categoryName,
		categoryDescription,
		categoryID,
	)
	if err != nil {
		return fmt.Errorf("AI标签识别失败: %v", err)
	}

	result.AIResponse = aiResponse

	// 从AI响应中提取HTTP耗时
	if aiResponse != nil {
		result.HttpDuration = aiResponse.HttpDuration
	}

	return nil
}

func (pp *PipelineProcessor) dbSaver(workerID int) {
	defer pp.wg.Done()

	ticker := time.NewTicker(pp.batchTimeout)
	defer ticker.Stop()

	var batch []*ProcessResult

	for {
		select {
		case result, ok := <-pp.resultChan:
			if !ok {
				if len(batch) > 0 {
					pp.saveBatch(batch)
				}
				return
			}

			batch = append(batch, result)

			if len(batch) >= pp.batchSize {
				pp.saveBatch(batch)
				batch = nil
			}

		case <-ticker.C:
			if len(batch) > 0 {
				pp.saveBatch(batch)
				batch = nil
			}
		}
	}
}

func (pp *PipelineProcessor) saveBatch(batch []*ProcessResult) {
	for _, result := range batch {
		if result.Success {
			err := pp.saveResultToDB(result)
			if err != nil {
				result.Nack(1*time.Second, false, fmt.Sprintf("save failed: %v", err))
			} else {
				result.Ack()
				pp.service.notifyQueueStatsChange()
			}
		} else {
			var currentTries int
			_ = pp.service.db.Model(&models.File{}).
				Where("id = ?", result.FileID).
				Select("ai_tagging_tries").
				Scan(&currentTries).Error

			currentTries++
			maxRetries := 3

			if currentTries >= maxRetries {
				_ = pp.service.db.Model(&models.File{}).
					Where("id = ?", result.FileID).
					Updates(map[string]interface{}{
						"ai_tagging_status": common.AITaggingStatusFailed,
						"ai_tagging_tries":  maxRetries,
					}).Error
				result.Ack()
			} else {
				_ = pp.service.db.Model(&models.File{}).
					Where("id = ?", result.FileID).
					Update("ai_tagging_tries", currentTries).Error

				retryDelay := time.Duration(currentTries) * time.Second
				result.Nack(retryDelay, false, result.Error.Error())
			}
		}
	}
}

func (pp *PipelineProcessor) saveResultToDB(result *ProcessResult) error {
	db := pp.service.db

	var fileCheck models.File
	if err := db.Where("id = ?", result.FileID).Select("id").Take(&fileCheck).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errFileDeleted
		}
		return err
	}

	if result.CategoryResult != nil {
		err := saveCategoryResultIfNeeded(db, result.FileID, &ai.FileCategorizationResponse{
			Success:             true,
			CategoryID:          result.CategoryResult.CategoryID,
			CategoryName:        result.CategoryResult.CategoryName,
			CategoryDescription: result.CategoryResult.CategoryDescription,
		})
		if err != nil {
			logger.Warn("保存分类结果失败: %v", err)
		}
	}

	contentDetectionEnabled := getContentDetectionEnabled()
	sensitiveContentHandling := getSensitiveContentHandling()

	err := processAIResponse(
		db,
		fileCheck,
		result.AIResponse,
		contentDetectionEnabled,
		sensitiveContentHandling,
		"",
		"",
	)
	if err != nil {
		return err
	}

	return db.Model(&models.File{}).
		Where("id = ?", result.FileID).
		Updates(map[string]interface{}{
			"ai_tagging_status": common.AITaggingStatusDone,
			"ai_tagging_tries":  0,
			"ai_http_duration":  result.HttpDuration,
		}).Error
}

func getContentDetectionEnabled() bool {
	return true
}

func getSensitiveContentHandling() string {
	return "mark_only"
}
