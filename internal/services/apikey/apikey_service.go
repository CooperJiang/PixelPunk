package apikey

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generateAPIKeyID() string {
	for {
		keyID := "ak_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:10]

		var count int64
		if err := database.DB.Model(&models.APIKey{}).Where("id = ?", keyID).Count(&count).Error; err != nil {
			return keyID
		}

		if count == 0 {
			return keyID
		}
	}
}

func generateAPIKeyValue() (string, error) {
	for {
		bytes := make([]byte, 32)
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}

		keyValue := base64.URLEncoding.EncodeToString(bytes)

		if len(keyValue) > 32 {
			keyValue = keyValue[:32]
		}

		var count int64
		if err := database.DB.Model(&models.APIKey{}).Where("key = ?", keyValue).Count(&count).Error; err != nil {
			return keyValue, nil
		}

		if count == 0 {
			return keyValue, nil
		}
	}
}

func formatAllowedTypes(types []string) string {
	if len(types) == 0 {
		return ""
	}
	return strings.Join(types, ",")
}

/* ParseAllowedTypes 解析允许的文件类型为切片 */
func ParseAllowedTypes(typesStr string) []string {
	if typesStr == "" {
		return []string{}
	}
	return strings.Split(typesStr, ",")
}

func calculateExpiresAt(expiresInDays int) *common.JSONTime {
	if expiresInDays <= 0 {
		return nil // 0或负数表示永不过期
	}

	expiresAt := time.Now().AddDate(0, 0, expiresInDays)
	expiresAt = time.Date(expiresAt.Year(), expiresAt.Month(), expiresAt.Day(), 23, 59, 59, 0, expiresAt.Location())

	jsonTime := common.JSONTime(expiresAt)
	return &jsonTime
}

/* CreateAPIKey 创建新的API密钥 */
func CreateAPIKey(userID uint, name string, storageLimit, singleFileLimit int64, uploadCountLimit int, allowedTypes []string, folderID string, expiresInDays int) (*models.APIKey, string, error) {
	db := database.DB

	keyID := generateAPIKeyID()
	keyValue, err := generateAPIKeyValue()
	if err != nil {
		return nil, "", errors.Wrap(err, errors.CodeInternal, "生成API密钥失败")
	}

	if folderID != "" {
		var count int64
		if err := db.Model(&models.Folder{}).Where("id = ? AND user_id = ?", folderID, userID).Count(&count).Error; err != nil {
			return nil, "", errors.Wrap(err, errors.CodeDBQueryFailed, "检查文件夹失败")
		}

		if count == 0 {
			return nil, "", errors.New(errors.CodeFolderNotFound, "指定的文件夹不存在")
		}
	}

	expiresAt := calculateExpiresAt(expiresInDays)

	apiKey := models.APIKey{
		ID:               keyID,
		UserID:           userID,
		Name:             name,
		KeyValue:         keyValue, // 实际值，不会返回给用户
		Status:           models.APIKeyStatusActive,
		StorageLimit:     storageLimit,
		StorageUsed:      0,
		UploadCountLimit: uploadCountLimit,
		UploadCountUsed:  0,
		SingleFileLimit:  singleFileLimit,
		AllowedTypes:     formatAllowedTypes(allowedTypes),
		FolderID:         folderID,
		ExpiresAt:        expiresAt,
		CreatedAt:        common.JSONTimeNow(),
		UpdatedAt:        common.JSONTimeNow(),
	}

	if err := db.Create(&apiKey).Error; err != nil {
		return nil, "", errors.Wrap(err, errors.CodeDBCreateFailed, "创建API密钥失败")
	}

	return &apiKey, keyValue, nil
}

/* GetAPIKey 根据密钥ID获取API密钥详情 */
func GetAPIKey(userID uint, keyID string) (*models.APIKey, error) {
	db := database.DB

	var apiKey models.APIKey
	if err := db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "API密钥不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥失败")
	}

	return &apiKey, nil
}

/* GetAPIKeyList 获取用户的API密钥列表 */
func GetAPIKeyList(userID uint, page, size int, status int, search string) ([]models.APIKey, int64, error) {
	db := database.DB

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = common.DefaultPageSize
	}
	if size > common.MaxPageSize {
		size = common.MaxPageSize
	}

	query := db.Model(&models.APIKey{}).Where("user_id = ?", userID)

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥总数失败")
	}

	var apiKeys []models.APIKey
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&apiKeys).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥列表失败")
	}

	return apiKeys, total, nil
}

/* UpdateAPIKey 更新API密钥 */
func UpdateAPIKey(userID uint, keyID string, updates map[string]interface{}) (*models.APIKey, error) {
	db := database.DB

	var apiKey models.APIKey
	if err := db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "API密钥不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥失败")
	}

	if expiresInDays, ok := updates["expires_in_days"].(int); ok {
		updates["expires_at"] = calculateExpiresAt(expiresInDays)
		delete(updates, "expires_in_days")
	}

	if allowedTypes, ok := updates["allowed_types"].([]string); ok {
		updates["allowed_types"] = formatAllowedTypes(allowedTypes)
	}

	if folderID, ok := updates["folder_id"].(string); ok && folderID != "" {
		var count int64
		if err := db.Model(&models.Folder{}).Where("id = ? AND user_id = ?", folderID, userID).Count(&count).Error; err != nil {
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "检查文件夹失败")
		}

		if count == 0 {
			return nil, errors.New(errors.CodeFolderNotFound, "指定的文件夹不存在")
		}
	}

	updates["updated_at"] = common.JSONTimeNow()

	if err := db.Model(&apiKey).Updates(updates).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新API密钥失败")
	}

	if err := db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询更新后的API密钥失败")
	}

	return &apiKey, nil
}

/* DeleteAPIKey 删除API密钥 */
func DeleteAPIKey(userID uint, keyID string) error {
	db := database.DB

	var apiKey models.APIKey
	if err := db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "API密钥不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥失败")
	}

	var count int64
	if err := db.Model(&models.File{}).Where("api_key_id = ?", keyID).Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "检查API密钥使用情况失败")
	}

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	err := db.Transaction(func(tx *gorm.DB) error {
		if count > 0 {
			if err := tx.Model(&models.File{}).Where("api_key_id = ?", keyID).Update("api_key_id", "").Error; err != nil {
				return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新文件关联失败")
			}
		}

		if err := tx.Delete(&apiKey).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除API密钥失败")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

/* ToggleAPIKeyStatus 切换API密钥状态 */
func ToggleAPIKeyStatus(userID uint, keyID string) (*models.APIKey, error) {
	db := database.DB

	var apiKey models.APIKey
	if err := db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "API密钥不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥失败")
	}

	var newStatus int
	if apiKey.Status == models.APIKeyStatusActive {
		newStatus = models.APIKeyStatusDisabled
	} else {
		newStatus = models.APIKeyStatusActive
	}

	if err := db.Model(&apiKey).Updates(map[string]interface{}{
		"status":     newStatus,
		"updated_at": common.JSONTimeNow(),
	}).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新API密钥状态失败")
	}

	if err := db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询更新后的API密钥失败")
	}

	return &apiKey, nil
}

/* GetAPIKeyByValue 根据密钥值获取API密钥 */
func GetAPIKeyByValue(keyValue string) (*models.APIKey, error) {
	db := database.DB

	var apiKey models.APIKey
	if err := db.Where("key_value = ?", keyValue).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeInvalidParameter, "无效的API密钥")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥失败")
	}

	return &apiKey, nil
}

/* UpdateAPIKeyUsage 更新API密钥使用情况 */
func UpdateAPIKeyUsage(keyID string, fileSize int64) error {
	db := database.DB

	var apiKey models.APIKey
	if err := db.Where("id = ?", keyID).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "API密钥不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥失败")
	}

	now := common.JSONTimeNow()
	updates := map[string]interface{}{
		"storage_used":      apiKey.StorageUsed + fileSize,
		"upload_count_used": apiKey.UploadCountUsed + 1,
		"last_used_at":      now,
		"updated_at":        now,
	}

	if err := db.Model(&apiKey).Updates(updates).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新API密钥使用情况失败")
	}

	return nil
}

/* GetFolderFullPath 获取文件夹的完整路径 */
func GetFolderFullPath(userID uint, folderID string) string {
	if folderID == "" {
		return "根目录"
	}

	db := database.DB
	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		return "根目录"
	}

	var pathParts []string
	current := folder

	for {
		pathParts = append([]string{current.Name}, pathParts...)

		if current.ParentID == "" {
			break
		}

		var parent models.Folder
		if err := db.Where("id = ? AND user_id = ?", current.ParentID, userID).First(&parent).Error; err != nil {
			break
		}
		current = parent
	}

	if len(pathParts) == 0 {
		return "根目录"
	}

	return "根目录/" + strings.Join(pathParts, "/")
}

/* GetAPIKeyStats 获取API密钥统计数据 */
func GetAPIKeyStats(userID uint, keyID string) (map[string]interface{}, error) {
	db := database.DB

	var apiKey models.APIKey
	if err := db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "API密钥不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥失败")
	}

	var imageCount int64
	if err := db.Model(&models.File{}).Where("api_key_id = ?", keyID).Count(&imageCount).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件总数失败")
	}

	var recentImages []models.File
	if err := db.Where("api_key_id = ?", keyID).
		Order("created_at DESC").
		Limit(5).
		Find(&recentImages).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询最近上传文件失败")
	}

	var recentUploads []map[string]interface{}
	for _, file := range recentImages {
		recentUploads = append(recentUploads, map[string]interface{}{
			"id":         file.ID,
			"name":       file.OriginalName,
			"size":       file.Size,
			"created_at": file.CreatedAt,
			"url":        file.URL,
		})
	}

	var storagePercentage float64 = -1 // -1表示不限量
	if apiKey.StorageLimit > 0 {
		storagePercentage = float64(apiKey.StorageUsed) / float64(apiKey.StorageLimit) * 100
	}

	var uploadPercentage float64 = -1 // -1表示不限量
	if apiKey.UploadCountLimit > 0 {
		uploadPercentage = float64(apiKey.UploadCountUsed) / float64(apiKey.UploadCountLimit) * 100
	}

	folderPath := GetFolderFullPath(userID, apiKey.FolderID)

	stats := map[string]interface{}{
		"storage_used":       apiKey.StorageUsed,
		"storage_limit":      apiKey.StorageLimit,
		"storage_percentage": storagePercentage,
		"upload_count_used":  apiKey.UploadCountUsed,
		"upload_count_limit": apiKey.UploadCountLimit,
		"upload_percentage":  uploadPercentage,
		"file_count":         imageCount,
		"last_used":          apiKey.LastUsedAt,
		"recent_uploads":     recentUploads,
		"allowed_types":      ParseAllowedTypes(apiKey.AllowedTypes),
		"single_file_limit":  apiKey.SingleFileLimit,
		"is_expired":         apiKey.IsExpired(),
		"is_active":          apiKey.IsActive(),
		"folder_id":          apiKey.FolderID,
		"folder_path":        folderPath,
	}

	return stats, nil
}

/* RegenerateAPIKey 重新生成API密钥值 */
func RegenerateAPIKey(userID uint, keyID string) (string, error) {
	db := database.DB

	var apiKey models.APIKey
	if err := db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New(errors.CodeNotFound, "API密钥不存在")
		}
		return "", errors.Wrap(err, errors.CodeDBQueryFailed, "查询API密钥失败")
	}

	keyValue, err := generateAPIKeyValue()
	if err != nil {
		return "", errors.Wrap(err, errors.CodeInternal, "生成API密钥失败")
	}

	if err := db.Model(&apiKey).Updates(map[string]interface{}{
		"key_value":  keyValue,
		"updated_at": common.JSONTimeNow(),
	}).Error; err != nil {
		return "", errors.Wrap(err, errors.CodeDBUpdateFailed, "更新API密钥失败")
	}

	return keyValue, nil
}

/* ValidateAPIKey 验证API密钥 */
func ValidateAPIKey(keyValue string) (*models.APIKey, error) {
	apiKey, err := GetAPIKeyByValue(keyValue)
	if err != nil {
		return nil, err
	}

	if !apiKey.IsActive() {
		return nil, errors.New(errors.CodeForbidden, "API密钥已禁用")
	}

	if apiKey.IsExpired() {
		return nil, errors.New(errors.CodeForbidden, "API密钥已过期")
	}

	return apiKey, nil
}

/* CheckAPIKeyLimits 检查API密钥是否超过限制 */
func CheckAPIKeyLimits(apiKey *models.APIKey, fileSize int64) error {
	if !apiKey.CheckUploadCountLimit() {
		return errors.New(errors.CodeForbidden, "已达到上传次数限制")
	}

	if !apiKey.CheckStorageLimit(fileSize) {
		return errors.New(errors.CodeForbidden, "已达到存储容量限制")
	}

	if !apiKey.CheckSingleFileLimit(fileSize) {
		return errors.New(errors.CodeForbidden, fmt.Sprintf("文件大小超过限制 (%d字节)", apiKey.SingleFileLimit))
	}

	return nil
}
