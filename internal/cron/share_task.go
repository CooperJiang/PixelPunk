package cron

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/internal/services/share"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/logger"
	"time"
)

func registerShareTask() {
	// 清理过期的分享访问令牌 - 每小时执行一次
	_, err := cronManager.AddFunc("0 0 * * * *", func() {
		cleanedCount, err := share.CleanExpiredTokens()
		if err != nil {
			logger.Error("清理过期的分享访问令牌失败: %v", err)
		} else {
			if cleanedCount > 0 {
				activity.LogSystemCleanup(int(cleanedCount), "过期访问令牌")
			}
		}
	})
	if err != nil {
		logger.Error("注册清理过期分享访问令牌任务失败: %v", err)
	}

	// 检查即将过期的分享并发送通知 - 每天早上9点执行
	_, err = cronManager.AddFunc("0 0 9 * * *", func() {
		checkAndNotifyExpiringShares()
	})
	if err != nil {
		logger.Error("注册分享过期提醒任务失败: %v", err)
	}
}

func checkAndNotifyExpiringShares() {
	// 查找24小时内即将过期的分享
	var shares []models.Share
	now := time.Now()
	warningTime := now.Add(24 * time.Hour)

	err := db.Where("status = ? AND expired_at IS NOT NULL AND expired_at > ? AND expired_at <= ?",
		common.ShareStatusNormal,
		now,
		warningTime,
	).Find(&shares).Error

	if err != nil {
		logger.Error("查询即将过期的分享失败: %v", err)
		return
	}

	if len(shares) == 0 {
		return
	}

	msgService := messageService.GetMessageService()
	notifiedCount := 0

	for _, share := range shares {
		remaining := time.Until(time.Time(*share.ExpiredAt))
		var expiresIn string
		if remaining < time.Hour {
			expiresIn = "1小时内"
		} else {
			hours := int(remaining.Hours())
			expiresIn = fmt.Sprintf("%d小时", hours)
		}

		variables := map[string]interface{}{
			"share_id":     share.ID,
			"share_name":   share.Name,
			"expires_in":   expiresIn,
			"related_type": "share",
			"related_id":   share.ID,
		}

		if err := msgService.SendTemplateMessage(share.UserID, common.MessageTypeShareExpiryWarning, variables); err != nil {
			logger.Warn("发送分享过期提醒失败: userID=%d, shareID=%s, error=%v", share.UserID, share.ID, err)
		} else {
			notifiedCount++
		}
	}

	if notifiedCount > 0 {
		logger.Info("分享过期提醒: 发送了 %d 条通知", notifiedCount)
	}
}
