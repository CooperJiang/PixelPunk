package storage

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ExportChannelConfig(channelID string) (interface{}, error) {
	channel, err := GetChannelByID(channelID)
	if err != nil {
		return nil, err
	}

	configs, err := GetChannelConfigsForExport(channelID)
	if err != nil {
		return nil, err
	}

	exportData := struct {
		Channel    models.StorageChannel      `json:"channel"`
		Configs    []models.StorageConfigItem `json:"configs"`
		Version    string                     `json:"version"`
		ExportTime string                     `json:"export_time"`
	}{
		Channel:    *channel,
		Configs:    configs,
		Version:    "1.0",
		ExportTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	return exportData, nil
}

func ExportAllChannelConfigs() (interface{}, error) {
	channels, err := GetAllChannels()
	if err != nil {
		return nil, err
	}

	var exportableChannels []models.StorageChannel
	for _, channel := range channels {
		if !channel.IsLocal && channel.Type != "local" {
			exportableChannels = append(exportableChannels, channel)
		}
	}

	if len(exportableChannels) == 0 {
		return nil, errors.New(errors.CodeValidationFailed, "没有可导出的渠道")
	}

	var exportData []interface{}
	for _, channel := range exportableChannels {
		configs, err := GetChannelConfigsForExport(channel.ID)
		if err != nil {
			continue // 跳过获取配置失败的渠道
		}

		channelExportData := struct {
			Channel    models.StorageChannel      `json:"channel"`
			Configs    []models.StorageConfigItem `json:"configs"`
			Version    string                     `json:"version"`
			ExportTime string                     `json:"export_time"`
		}{
			Channel:    channel,
			Configs:    configs,
			Version:    "1.0",
			ExportTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		exportData = append(exportData, channelExportData)
	}

	finalExportData := struct {
		Channels   []interface{} `json:"channels"`
		TotalCount int           `json:"total_count"`
		ExportType string        `json:"export_type"`
		Version    string        `json:"version"`
		ExportTime string        `json:"export_time"`
	}{
		Channels:   exportData,
		TotalCount: len(exportData),
		ExportType: "batch",
		Version:    "1.0",
		ExportTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	return finalExportData, nil
}

func ImportChannelConfig(data []byte) error {
	db := database.GetDB()

	if len(data) == 0 {
		return errors.New(errors.CodeInvalidParameter, "导入文件为空")
	}

	var batchImportData struct {
		Channels   []json.RawMessage `json:"channels"`
		TotalCount int               `json:"total_count"`
		ExportType string            `json:"export_type"`
		Version    string            `json:"version"`
	}

	if err := json.Unmarshal(data, &batchImportData); err == nil && batchImportData.ExportType == "batch" && len(batchImportData.Channels) > 0 {
		return importMultipleChannels(db, batchImportData.Channels)
	}

	return importSingleChannel(db, data)
}

func importSingleChannel(db *gorm.DB, data []byte) error {
	var importData struct {
		Channel struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Type      string `json:"type"`
			Status    int8   `json:"status"`
			IsDefault bool   `json:"is_default"`
			IsLocal   bool   `json:"is_local"`
			Remark    string `json:"remark"`
		} `json:"channel"`
		Configs []struct {
			Name        string `json:"name"`
			KeyName     string `json:"key_name"`
			Value       string `json:"value"`
			Type        string `json:"type"`
			IsSecret    bool   `json:"is_secret"`
			Required    bool   `json:"required"`
			Description string `json:"description"`
		} `json:"configs"`
		Version string `json:"version"`
	}

	if err := json.Unmarshal(data, &importData); err != nil {
		return errors.New(errors.CodeInvalidParameter, "JSON格式错误："+err.Error())
	}

	return validateAndImportChannel(db, importData.Channel, importData.Configs)
}

func importMultipleChannels(db *gorm.DB, channelsData []json.RawMessage) error {
	if len(channelsData) == 0 {
		return errors.New(errors.CodeValidationFailed, "没有找到可导入的渠道")
	}

	var successCount, errorCount int
	var errorMessages []string

	for i, channelData := range channelsData {
		var singleChannelData struct {
			Channel struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Type      string `json:"type"`
				Status    int8   `json:"status"`
				IsDefault bool   `json:"is_default"`
				IsLocal   bool   `json:"is_local"`
				Remark    string `json:"remark"`
			} `json:"channel"`
			Configs []struct {
				Name        string `json:"name"`
				KeyName     string `json:"key_name"`
				Value       string `json:"value"`
				Type        string `json:"type"`
				IsSecret    bool   `json:"is_secret"`
				Required    bool   `json:"required"`
				Description string `json:"description"`
			} `json:"configs"`
		}

		if err := json.Unmarshal(channelData, &singleChannelData); err != nil {
			errorCount++
			errorMessages = append(errorMessages, fmt.Sprintf("渠道%d解析失败: %v", i+1, err))
			continue
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			return validateAndImportChannel(tx, singleChannelData.Channel, singleChannelData.Configs)
		})

		if err != nil {
			errorCount++
			errorMessages = append(errorMessages, fmt.Sprintf("渠道'%s'导入失败: %v", singleChannelData.Channel.Name, err))
		} else {
			successCount++
		}
	}

	if successCount == 0 {
		return errors.New(errors.CodeValidationFailed, "所有渠道导入失败: "+strings.Join(errorMessages, "; "))
	}

	if errorCount > 0 {
		return errors.New(errors.CodeValidationFailed, fmt.Sprintf("部分导入成功：成功%d个，失败%d个。失败原因: %s", successCount, errorCount, strings.Join(errorMessages, "; ")))
	}

	return nil // 全部成功
}

func validateAndImportChannel(db *gorm.DB, channelData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Status    int8   `json:"status"`
	IsDefault bool   `json:"is_default"`
	IsLocal   bool   `json:"is_local"`
	Remark    string `json:"remark"`
}, configsData []struct {
	Name        string `json:"name"`
	KeyName     string `json:"key_name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	IsSecret    bool   `json:"is_secret"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}) error {

	if channelData.Name == "" {
		return errors.New(errors.CodeValidationFailed, "渠道名称不能为空")
	}

	if channelData.Type == "" {
		return errors.New(errors.CodeValidationFailed, "渠道类型不能为空")
	}

	if _, ok := models.StorageConfigTemplates[channelData.Type]; !ok {
		return errors.New(errors.CodeValidationFailed, "不支持的渠道类型："+channelData.Type)
	}

	var nameCount int64
	if err := db.Model(&models.StorageChannel{}).Where("name = ?", channelData.Name).Count(&nameCount).Error; err != nil {
		return fmt.Errorf("检查渠道名称失败: %v", err)
	}
	if nameCount > 0 {
		return errors.New(errors.CodeValidationFailed, "渠道名称已存在，请修改后重试")
	}

	if channelData.Type == "local" {
		var localCount int64
		if err := db.Model(&models.StorageChannel{}).Where("type = ?", "local").Count(&localCount).Error; err != nil {
			return fmt.Errorf("检查本地存储渠道失败: %v", err)
		}
		if localCount > 0 {
			return errors.New(errors.CodeValidationFailed, "本地存储渠道已存在，不能重复创建")
		}
	}

	if len(configsData) == 0 {
		return errors.New(errors.CodeValidationFailed, "配置项不能为空")
	}

	templates, ok := models.StorageConfigTemplates[channelData.Type]
	if !ok {
		return errors.New(errors.CodeValidationFailed, "不支持的存储类型: "+channelData.Type)
	}

	configMap := make(map[string]string)
	for _, config := range configsData {
		configMap[config.KeyName] = config.Value
	}

	for _, tmpl := range templates {
		if tmpl.Required {
			value, exists := configMap[tmpl.KeyName]
			if !exists || value == "" {
				return errors.New(errors.CodeValidationFailed, "缺少必需的配置项: "+tmpl.Name)
			}
		}
	}

	newChannel := models.StorageChannel{
		Name:      channelData.Name,
		Type:      channelData.Type,
		Status:    channelData.Status,
		IsDefault: false, // 导入的渠道默认不是默认渠道
		Remark:    channelData.Remark,
	}

	if channelData.Type == "local" {
		newChannel.IsLocal = true
	}

	uuidStr := uuid.New().String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	newChannel.ID = uuidStr
	newChannel.CreatedAt = common.JSONTimeNow()
	newChannel.UpdatedAt = common.JSONTimeNow()

	if err := db.Create(&newChannel).Error; err != nil {
		return fmt.Errorf("创建渠道失败: %v", err)
	}

	for _, config := range configsData {
		configItemUUID := uuid.New().String()
		configItemUUID = strings.ReplaceAll(configItemUUID, "-", "")

		configItem := models.StorageConfigItem{
			ID:          configItemUUID,
			ChannelID:   newChannel.ID,
			Name:        config.Name,
			KeyName:     config.KeyName,
			Value:       config.Value,
			Type:        config.Type,
			IsSecret:    config.IsSecret,
			Required:    config.Required,
			Description: config.Description,
			CreatedAt:   common.JSONTimeNow(),
			UpdatedAt:   common.JSONTimeNow(),
		}

		if err := db.Create(&configItem).Error; err != nil {
			return fmt.Errorf("创建配置项失败: %v", err)
		}
	}

	return nil
}
