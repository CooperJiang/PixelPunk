package utils

import (
	"fmt"
	"time"
)

// FormatBytes 格式化字节数为人类可读格式
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration 将秒数格式化为人类可读的持续时间
func FormatDuration(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟%d秒", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分钟%d秒", minutes, secs)
	}
	return fmt.Sprintf("%d秒", secs)
}

// FormatTimeAgo 将时间格式化为"多久之前"的形式
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	seconds := int(duration.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	months := days / 30
	years := days / 365

	switch {
	case years > 0:
		return fmt.Sprintf("%d年前", years)
	case months > 0:
		return fmt.Sprintf("%d个月前", months)
	case days > 0:
		return fmt.Sprintf("%d天前", days)
	case hours > 0:
		return fmt.Sprintf("%d小时前", hours)
	case minutes > 0:
		return fmt.Sprintf("%d分钟前", minutes)
	default:
		return "刚刚"
	}
}
