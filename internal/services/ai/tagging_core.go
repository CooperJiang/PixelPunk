package ai

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	qqueue "pixelpunk/internal/queue"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/storage"

	"gorm.io/gorm"
)

type TaggingService struct {
	db                *gorm.DB
	storage           *storage.Storage
	concurrentNum     int
	mutex             sync.Mutex
	stopping          bool
	paused            bool
	activeWorkers     int // Pipeline模式下不使用
	ctx               context.Context
	cancel            context.CancelFunc
	resizeCh          chan struct{}
	failureTimestamps []time.Time
	failureThreshold  int

	taskQueue  qqueue.Queue
	reaperStop chan struct{}

	pipeline *PipelineProcessor

	pushMu        sync.Mutex
	lastPush      time.Time
	pushScheduled bool
}

var errMissingFile = errors.New("ai:missing_file")

func getMissingFilePolicy() string {
	s := strings.ToLower(strings.TrimSpace(setting.GetString("ai", "missing_file_policy", "drop")))
	switch s {
	case "drop", "ignore", "retry":
		return s
	default:
		return "drop"
	}
}

func getAIConcurrency() int {
	v := setting.GetInt("ai", "ai_concurrency", 5)
	if v <= 0 {
		v = 5
	}
	return v
}

func NewTaggingService(db *gorm.DB, concurrentNum int) *TaggingService {
	if concurrentNum <= 0 {
		concurrentNum = getAIConcurrency()
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &TaggingService{
		db:                db,
		storage:           storage.NewGlobalStorage(),
		concurrentNum:     concurrentNum,
		ctx:               ctx,
		cancel:            cancel,
		resizeCh:          make(chan struct{}),
		paused:            false,
		failureTimestamps: make([]time.Time, 0, 64),
		failureThreshold:  20,
		reaperStop:        make(chan struct{}),
	}

	if cache.IsRedisEnabled() {
		if rq := qqueue.NewRedisQueue(); rq != nil {
			service.taskQueue = rq
			rq.StartReaper(1*time.Second, service.reaperStop)
		} else {
			service.taskQueue = qqueue.NewDBQueue()
		}
	} else {
		service.taskQueue = qqueue.NewDBQueue()
	}

	service.pipeline = NewPipelineProcessor(service, concurrentNum)
	service.pipeline.Start()

	return service
}

func NewTaggingServiceWithConfig(db *gorm.DB) *TaggingService {
	return NewTaggingService(db, 0)
}
