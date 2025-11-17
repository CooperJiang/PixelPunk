package ai

import (
	"time"

	"pixelpunk/internal/controllers/websocket"
	ws "pixelpunk/internal/websocket"
	"pixelpunk/pkg/logger"
)

func (s *TaggingService) GetActiveWorkers() int {
	return s.activeWorkers
}

func (s *TaggingService) notifyQueueStatsChange() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("WebSocket推送goroutine panic: %v", r)
			}
		}()
		s.pushMu.Lock()
		minGap := 500 * time.Millisecond
		since := time.Since(s.lastPush)
		if since >= minGap && !s.pushScheduled {
			s.lastPush = time.Now()
			s.pushMu.Unlock()
			s.performWebSocketPush()
			return
		}
		if s.pushScheduled {
			s.pushMu.Unlock()
			return
		}
		s.pushScheduled = true
		delay := minGap - since
		if delay < 0 {
			delay = minGap
		}
		s.pushMu.Unlock()
		time.AfterFunc(delay, func() {
			s.pushMu.Lock()
			s.lastPush = time.Now()
			s.pushScheduled = false
			s.pushMu.Unlock()
			s.performWebSocketPush()
		})
	}()
}

func (s *TaggingService) performWebSocketPush() {
	stats, err := GetAIQueueStats()
	if err != nil {
		return
	}

	websocket.BroadcastToAdmins(ws.MessageTypeQueueStats, stats)
}

func (s *TaggingService) Pause() {
	s.mutex.Lock()
	s.paused = true
	s.mutex.Unlock()
	s.notifyQueueStatsChange()
	websocket.BroadcastToAdmins(ws.MessageTypeAnnouncement, map[string]interface{}{
		"title":   "AI 队列已暂停",
		"content": "自动处理已关闭，新的上传仍会入队但不会执行。",
		"ts":      time.Now().Unix(),
	})
}

func (s *TaggingService) Resume() {
	s.mutex.Lock()
	s.paused = false
	s.failureTimestamps = s.failureTimestamps[:0]
	s.mutex.Unlock()
	s.notifyQueueStatsChange()
	websocket.BroadcastToAdmins(ws.MessageTypeAnnouncement, map[string]interface{}{
		"title":   "AI 队列已恢复",
		"content": "自动处理已开启，队列将继续执行。",
		"ts":      time.Now().Unix(),
	})
}

func (s *TaggingService) IsPaused() bool {
	s.mutex.Lock()
	p := s.paused
	s.mutex.Unlock()
	return p
}

func (s *TaggingService) RecentFailures() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	cutoff := time.Now().Add(-1 * time.Minute)
	count := 0
	for _, ts := range s.failureTimestamps {
		if ts.After(cutoff) {
			count++
		}
	}
	return count
}

// processImageWithBase64 使用Base64方式处理文件AI分析
