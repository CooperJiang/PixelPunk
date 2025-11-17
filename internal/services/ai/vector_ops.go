package ai

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	vector2 "pixelpunk/internal/services/vector"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/vector"
	"strings"

	"gorm.io/gorm"
)

// buildEnrichedVectorContent 构建专为语义搜索优化的向量内容
func buildEnrichedVectorContent(fileID, description string) (string, error) {
	db := GetDBFromContext()
	if db == nil {
		return description, fmt.Errorf("无法获取数据库连接")
	}
	var aiInfo models.FileAIInfo
	if err := db.Where("file_id = ?", fileID).First(&aiInfo).Error; err != nil {
		return description, err
	}
	var file models.File
	if err := db.Where("id = ?", fileID).First(&file).Error; err != nil {
		logger.Warn("获取文件基础信息失败: %v", err)
	}
	var contentParts []string
	var tags []string
	if len(aiInfo.Tags) > 0 {
		if err := json.Unmarshal(aiInfo.Tags, &tags); err == nil && len(tags) > 0 {
			maxTags := min(5, len(tags))
			topTags := tags[:maxTags]
			contentParts = append(contentParts, strings.Join(topTags, " "))
		}
	}
	if aiInfo.SearchContent != "" {
		contentParts = append(contentParts, aiInfo.SearchContent)
	}
	var semanticKeywords []string
	if len(aiInfo.SemanticKeywords) > 0 {
		if err := json.Unmarshal(aiInfo.SemanticKeywords, &semanticKeywords); err == nil && len(semanticKeywords) > 0 {
			tagSet := make(map[string]bool)
			for _, tag := range tags {
				tagSet[strings.ToLower(tag)] = true
			}
			var uniqueKeywords []string
			maxKeywords := min(3, len(semanticKeywords))
			for i := 0; i < maxKeywords && len(uniqueKeywords) < 3; i++ {
				keyword := semanticKeywords[i]
				if !tagSet[strings.ToLower(keyword)] {
					uniqueKeywords = append(uniqueKeywords, keyword)
				}
			}
			if len(uniqueKeywords) > 0 {
				contentParts = append(contentParts, strings.Join(uniqueKeywords, " "))
			}
		}
	}
	if description != "" && description != aiInfo.SearchContent && len(description) <= 100 {
		contentParts = append(contentParts, description)
	}
	if file.OriginalName != "" && isFileNameMeaningful(file.OriginalName) {
		filename := processFileName(file.OriginalName)
		if filename != "" {
			contentParts = append(contentParts, filename)
		}
	}
	if file.Width >= 1920 && file.Height >= 1080 {
		if file.Width >= 7680 || file.Height >= 4320 {
			contentParts = append(contentParts, "8K超清")
		} else if file.Width >= 3840 || file.Height >= 2160 {
			contentParts = append(contentParts, "4K超清")
		} else {
			contentParts = append(contentParts, "高清")
		}
	}
	result := strings.Join(contentParts, " ")
	if strings.TrimSpace(result) == "" {
		return description, nil
	}
	return result, nil
}

// 文件名提取关键词相关
func isFileNameMeaningful(filename string) bool {
	if filename == "" {
		return false
	}
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	if len(name) < 3 {
		return false
	}
	if strings.Trim(name, "0123456789") == "" {
		return false
	}
	randomChars := 0
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			randomChars++
		}
	}
	if float64(randomChars)/float64(len(name)) > 0.8 && len(name) > 10 {
		return false
	}
	return true
}
func processFileName(filename string) string {
	if filename == "" {
		return ""
	}
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, ".", " ")
	name = strings.ReplaceAll(name, "+", " ")
	for strings.Contains(name, "  ") {
		name = strings.ReplaceAll(name, "  ", " ")
	}
	name = strings.TrimSpace(name)
	if len(name) > 50 {
		name = name[:50]
		if last := strings.LastIndex(name, " "); last > 30 {
			name = name[:last]
		}
	}
	return name
}

// RegenerateFileVector 重新生成指定文件的向量
func RegenerateFileVector(fileID, description string) error {
	engine := vector.GetGlobalVectorEngine()
	if engine == nil {
		return fmt.Errorf("向量引擎未初始化")
	}
	if !engine.IsEnabled() {
		return fmt.Errorf("向量引擎未启用")
	}
	enrichedContent, err := buildEnrichedVectorContent(fileID, description)
	if err != nil {
		logger.Warn("构建丰富向量内容失败，使用原始描述: %v", err)
		enrichedContent = description
	}
	if err := engine.ProcessFile(fileID, enrichedContent); err != nil {
		logger.Error("重新生成向量失败: %v", err)
		return fmt.Errorf("重新生成向量失败: %v", err)
	}
	return nil
}

// createPendingVectorRecord 创建pending向量记录
func createPendingVectorRecord(fileID, description string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}
	currentModel := getCurrentVectorModel()

	var existingVector models.FileVector
	err := db.Where("file_id = ?", fileID).First(&existingVector).Error
	if err == gorm.ErrRecordNotFound {
		newVector := &models.FileVector{FileID: fileID, Description: description, Model: currentModel, Status: common.VectorStatusPending}
		if err := db.Create(newVector).Error; err != nil {
			logger.Error("创建新向量记录失败 [%s]: %v", fileID, err)
			return fmt.Errorf("创建新向量记录失败: %v", err)
		}
		logVectorProcessing(fileID, models.VectorLogActionStart, "vector.start",
			map[string]interface{}{"model": currentModel}, currentModel, 0, "", "")
		if vqs := vector2.GetGlobalVectorQueueService(); vqs != nil {
			_ = vqs.EnqueueVector(fileID)
		}
	} else if err != nil {
		logger.Error("查询向量记录失败 [%s]: %v", fileID, err)
		return fmt.Errorf("查询向量记录失败: %v", err)
	} else {
		updateData := map[string]interface{}{
			"status":        common.VectorStatusPending,
			"description":   description,
			"model":         currentModel,
			"error_message": "",
		}
		if err := db.Model(&existingVector).Updates(updateData).Error; err != nil {
			logger.Error("更新向量记录为pending失败 [%s]: %v", fileID, err)
			return fmt.Errorf("更新向量记录为pending失败: %v", err)
		}
		logVectorProcessing(fileID, models.VectorLogActionReset, "vector.reset",
			map[string]interface{}{"model": currentModel}, currentModel, 0, "", "")
		if vqs := vector2.GetGlobalVectorQueueService(); vqs != nil {
			_ = vqs.EnqueueVector(fileID)
		}
	}
	return nil
}

// getCurrentVectorModel 获取当前配置的向量模型
func getCurrentVectorModel() string {
	model := setting.GetString("vector", "vector_model", "text-embedding-3-small")
	if model == "" {
		model = "text-embedding-3-small"
	}
	return model
}

// logVectorProcessing 记录向量处理日志
func logVectorProcessing(fileID, action, logType string, data map[string]interface{}, model string, duration int, errorCode, taskID string) {
	db := database.GetDB()
	if db == nil {
		logger.Error("数据库连接不可用，无法记录向量处理日志")
		return
	}

	dataJSON, _ := json.Marshal(data)
	log := &models.VectorProcessingLog{
		FileID:    fileID,
		Action:    action,
		Type:      logType,
		Data:      string(dataJSON),
		Model:     model,
		Duration:  duration,
		ErrorCode: errorCode,
		TaskID:    taskID,
	}
	if err := db.Create(log).Error; err != nil {
		logger.Error("创建向量处理日志失败: %v", err)
	}
}
