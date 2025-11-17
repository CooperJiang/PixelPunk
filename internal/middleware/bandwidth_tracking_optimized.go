package middleware

import (
	"pixelpunk/internal/services/bandwidth"
	"pixelpunk/pkg/logger"

	"github.com/gin-gonic/gin"
)

type BandwidthTrackingWriter struct {
	gin.ResponseWriter
	bytesWritten int64
	userID       uint
}

func (w *BandwidthTrackingWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.bytesWritten += int64(n)
	return n, err
}

func (w *BandwidthTrackingWriter) GetBytesWritten() int64 {
	return w.bytesWritten
}

func BandwidthTrackingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetCurrentUserID(c)

		// 如果没有用户ID，跳过统计（匿名访问）
		if userID == 0 {
			c.Next()
			return
		}

		trackingWriter := &BandwidthTrackingWriter{
			ResponseWriter: c.Writer,
			userID:         userID,
		}
		c.Writer = trackingWriter

		c.Next()

		// 异步记录带宽使用
		go recordBandwidthUsage(trackingWriter)
	}
}

func recordBandwidthUsage(writer *BandwidthTrackingWriter) {
	if writer.bytesWritten == 0 {
		return
	}

	if err := bandwidth.Service.RecordBandwidthTransfer(writer.userID, writer.bytesWritten); err != nil {
		logger.Error("[BANDWIDTH_TRACKING] 记录带宽传输失败: %v", err)
	}
}
