package file

import (
	"net/http"
	"pixelpunk/internal/cron"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ManualCleanupChunkedUploads 手动清理分片上传
func ManualCleanupChunkedUploads(c *gin.Context) {
	cleanupJob := cron.NewChunkedUploadCleanupJob()

	if err := cleanupJob.Execute(); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "分片上传清理任务执行完成")
}

func GetChunkedUploadStats(c *gin.Context) {
	// 这里可以添加获取分片上传统计信息的逻辑
	// 比如当前活跃会话数、过期会话数等
	stats := gin.H{
		"message": "分片上传统计功能即将推出",
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "获取统计信息成功",
	})
}
