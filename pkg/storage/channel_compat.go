package storage

// Minimal compatibility ChannelRepository used by NewGlobalStorage.

import (
	"fmt"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
)

// CompatChannelRepository bridges manager.ChannelRepository to DB-layer helpers without importing internal services (to avoid cycles).
type CompatChannelRepository struct{}

func (r *CompatChannelRepository) GetChannel(channelID string) (*models.StorageChannel, error) {
	var channel models.StorageChannel
	if err := database.GetDB().Where("id = ?", channelID).First(&channel).Error; err != nil {
		return nil, fmt.Errorf("获取渠道失败: %w", err)
	}
	return &channel, nil
}

func (r *CompatChannelRepository) GetChannelConfig(channelID string) (map[string]interface{}, error) {
	var items []models.StorageConfigItem
	if err := database.GetDB().Where("channel_id = ?", channelID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("查询渠道配置失败: %w", err)
	}
	m := make(map[string]interface{}, len(items))
	for _, it := range items {
		m[it.KeyName] = it.Value
	}
	return m, nil
}

func (r *CompatChannelRepository) GetDefaultChannel() (*models.StorageChannel, error) {
	var ch models.StorageChannel
	if err := database.GetDB().Where("is_default = ? AND status = ?", true, 1).First(&ch).Error; err != nil {
		return nil, fmt.Errorf("未找到默认存储渠道: %w", err)
	}
	return &ch, nil
}

func (r *CompatChannelRepository) GetActiveChannels() ([]*models.StorageChannel, error) {
	var channels []models.StorageChannel
	if err := database.GetDB().Where("status = ?", 1).Find(&channels).Error; err != nil {
		return nil, fmt.Errorf("查询活跃渠道失败: %w", err)
	}
	res := make([]*models.StorageChannel, 0, len(channels))
	for i := range channels {
		res = append(res, &channels[i])
	}
	return res, nil
}
