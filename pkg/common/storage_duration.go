package common

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// 存储时长常量
const (
	StorageDuration3Days     = "3d"
	StorageDuration7Days     = "7d"
	StorageDuration30Days    = "30d"
	StorageDurationPermanent = "permanent"
)

// 支持的存储时长选项
var SupportedStorageDurations = map[string]bool{
	StorageDuration3Days:     true,
	StorageDuration7Days:     true,
	StorageDuration30Days:    true,
	StorageDurationPermanent: true,
}

// GuestStorageDurations 游客可用的存储时长选项（最长30天）
var GuestStorageDurations = map[string]bool{
	StorageDuration3Days: true,
	StorageDuration7Days: true,
}

// ParseStorageDuration 解析存储时长字符串为时间间隔
func ParseStorageDuration(duration string) time.Duration {
	switch duration {
	case StorageDuration3Days:
		return 3 * 24 * time.Hour
	case StorageDuration7Days:
		return 7 * 24 * time.Hour
	case StorageDuration30Days:
		return 30 * 24 * time.Hour
	case StorageDurationPermanent:
		return 0
	default:
		return parseCustomDuration(duration)
	}
}

// parseCustomDuration 解析自定义时长格式
func parseCustomDuration(duration string) time.Duration {
	// 使用正则表达式匹配数字和单位
	re := regexp.MustCompile(`^(\d+)([smhdSMHD])$`)
	matches := re.FindStringSubmatch(duration)

	if len(matches) != 3 {
		return 0
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil || value <= 0 {
		return 0
	}

	const (
		maxSeconds = 86400 * 365 * 10 // 最多10年
		maxMinutes = 1440 * 365 * 10
		maxHours   = 24 * 365 * 10
		maxDays    = 365 * 10
	)

	unit := matches[2]
	switch unit {
	case "s", "S":
		if value > maxSeconds {
			return 0
		}
		return time.Duration(value) * time.Second
	case "m", "M":
		if value > maxMinutes {
			return 0
		}
		return time.Duration(value) * time.Minute
	case "h", "H":
		if value > maxHours {
			return 0
		}
		return time.Duration(value) * time.Hour
	case "d", "D":
		if value > maxDays {
			return 0
		}
		return time.Duration(value) * 24 * time.Hour
	default:
		return 0
	}
}

// FormatDuration 格式化时间间隔为可读字符串
func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "永久"
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	switch {
	case days > 0:
		return fmt.Sprintf("%d天", days)
	case hours > 0:
		return fmt.Sprintf("%d小时", hours)
	case minutes > 0:
		return fmt.Sprintf("%d分钟", minutes)
	default:
		return "刚刚"
	}
}

// IsValidStorageDuration 检查存储时长是否有效
func IsValidStorageDuration(duration string) bool {
	_, exists := SupportedStorageDurations[duration]
	return exists
}

// IsValidGuestStorageDuration 检查游客存储时长是否有效
func IsValidGuestStorageDuration(duration string) bool {
	_, exists := GuestStorageDurations[duration]
	return exists
}

func GetMaxGuestStorageDuration() string {
	return StorageDuration7Days
}

// CalculateExpiryTime 计算过期时间
func CalculateExpiryTime(duration string) time.Time {
	if duration == StorageDurationPermanent {
		return time.Time{} // 返回零时间表示永久
	}

	d := ParseStorageDuration(duration)
	if d == 0 {
		return time.Time{}
	}

	return time.Now().Add(d)
}

// IsExpired 检查是否已过期
func IsExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}
	return time.Now().After(*expiresAt)
}

// IsExpiringSoon 检查是否即将过期（默认24小时内）
func IsExpiringSoon(expiresAt *time.Time, threshold ...time.Duration) bool {
	if expiresAt == nil {
		return false
	}

	checkThreshold := 24 * time.Hour
	if len(threshold) > 0 {
		checkThreshold = threshold[0]
	}

	remaining := time.Until(*expiresAt)
	return remaining < checkThreshold && remaining > 0
}

func GetRemainingTime(expiresAt *time.Time) time.Duration {
	if expiresAt == nil {
		return -1 // 返回-1表示永久存储
	}

	remaining := time.Until(*expiresAt)
	if remaining < 0 {
		return 0 // 已过期
	}

	return remaining
}

func GetRemainingDays(expiresAt *time.Time) int {
	remaining := GetRemainingTime(expiresAt)
	if remaining < 0 {
		return -1 // 永久存储
	}

	return int(remaining.Hours() / 24)
}

// FormatRemainingTime 格式化剩余时间为易读格式
func FormatRemainingTime(expiresAt *time.Time) string {
	remaining := GetRemainingTime(expiresAt)

	if remaining < 0 {
		return "永久"
	}

	if remaining == 0 {
		return "已过期"
	}

	if remaining < time.Minute {
		return "不到1分钟"
	}

	if remaining < time.Hour {
		return fmt.Sprintf("%d分钟", int(remaining.Minutes()))
	}

	if remaining < 24*time.Hour {
		return fmt.Sprintf("%d小时", int(remaining.Hours()))
	}

	days := int(remaining.Hours() / 24)
	return fmt.Sprintf("%d天", days)
}

func GetStorageOptions(isGuest bool) []StorageOption {
	options := []StorageOption{
		{Value: StorageDuration3Days, Label: "3天", Days: 3},
		{Value: StorageDuration7Days, Label: "7天", Days: 7},
		{Value: StorageDuration30Days, Label: "30天", Days: 30},
	}

	if !isGuest {
		options = append(options, StorageOption{
			Value: StorageDurationPermanent,
			Label: "永久",
			Days:  -1,
		})
	}

	return options
}

// StorageOption 存储时长选项
type StorageOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Days  int    `json:"days"`
}

// ValidateStorageDuration 验证存储时长
func ValidateStorageDuration(duration string, isGuest bool) error {
	if !IsValidStorageDuration(duration) {
		return fmt.Errorf("无效的存储时长: %s", duration)
	}

	if isGuest && !IsValidGuestStorageDuration(duration) {
		return fmt.Errorf("游客仅支持最多30天的存储时长")
	}

	return nil
}

// ValidateStorageDurationWithAllowed 使用允许的时长列表验证存储时长
func ValidateStorageDurationWithAllowed(duration string, allowedDurations []string, isGuest bool) error {
	if duration == "" {
		return fmt.Errorf("存储时长不能为空")
	}

	for _, allowed := range allowedDurations {
		if duration == allowed {
			// 游客模式下不允许永久存储
			if isGuest && duration == StorageDurationPermanent {
				return fmt.Errorf("游客上传不支持永久存储")
			}
			return nil
		}
	}

	var allowedStr string
	for _, allowed := range allowedDurations {
		// 游客模式下过滤掉permanent选项
		if isGuest && allowed == StorageDurationPermanent {
			continue
		}

		if allowedStr != "" {
			allowedStr += "、"
		}
		allowedStr += allowed
	}

	if isGuest {
		return fmt.Errorf("存储时长必须是 %s（游客仅支持限时存储）", allowedStr)
	} else {
		return fmt.Errorf("存储时长必须是 %s", allowedStr)
	}
}
