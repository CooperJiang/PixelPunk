package random_api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	mathrand "math/rand"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/storage"
	"time"

	"gorm.io/gorm"
)

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

func GenerateRandomAPIKey() (string, error) {
	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		bytes := make([]byte, 9)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}

		randomStr := base64.URLEncoding.EncodeToString(bytes)
		randomStr = randomStr[:12]
		apiKey := "rnd_" + randomStr

		var count int64
		if err := database.DB.Model(&models.RandomImageAPI{}).Where("api_key = ?", apiKey).Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return apiKey, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique API key after %d attempts", maxRetries)
}

func CreateRandomAPI(userID uint, name string, folderID *string, returnType string) (*models.RandomImageAPI, error) {
	if folderID != nil && *folderID != "" {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *folderID, userID).First(&folder).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New(errors.CodeFolderNotFound, "文件夹不存在")
			}
			return nil, err
		}
	}

	if returnType != models.RandomImageAPIReturnTypeRedirect && returnType != models.RandomImageAPIReturnTypeDirect {
		returnType = models.RandomImageAPIReturnTypeRedirect
	}

	apiKey, err := GenerateRandomAPIKey()
	if err != nil {
		return nil, errors.New(errors.CodeInternal, "生成API密钥失败")
	}

	randomAPI := &models.RandomImageAPI{
		UserID:     userID,
		Name:       name,
		APIKey:     apiKey,
		FolderID:   folderID,
		ReturnType: returnType,
		Status:     models.RandomImageAPIStatusActive,
	}

	if err := database.DB.Create(randomAPI).Error; err != nil {
		return nil, err
	}

	return randomAPI, nil
}

func GetRandomAPIList(userID uint, page, size int, status int, search string) ([]*models.RandomImageAPI, int64, error) {
	var apis []*models.RandomImageAPI
	var total int64

	query := database.DB.Model(&models.RandomImageAPI{}).Where("user_id = ?", userID)

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Preload("Folder").Order("created_at DESC").Offset(offset).Limit(size).Find(&apis).Error; err != nil {
		return nil, 0, err
	}

	return apis, total, nil
}

func GetRandomAPIByID(id, userID uint) (*models.RandomImageAPI, error) {
	var api models.RandomImageAPI
	if err := database.DB.Preload("Folder").Where("id = ? AND user_id = ?", id, userID).First(&api).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "随机API不存在")
		}
		return nil, err
	}
	return &api, nil
}

func UpdateRandomAPIStatus(id, userID uint, status int) error {
	result := database.DB.Model(&models.RandomImageAPI{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(errors.CodeNotFound, "随机API不存在")
	}

	return nil
}

func UpdateRandomAPIConfig(id, userID uint, folderID *string, returnType string) error {
	if folderID != nil && *folderID != "" {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *folderID, userID).First(&folder).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New(errors.CodeFolderNotFound, "文件夹不存在")
			}
			return err
		}
	}

	if returnType != models.RandomImageAPIReturnTypeRedirect && returnType != models.RandomImageAPIReturnTypeDirect {
		return errors.New(errors.CodeInvalidParameter, "无效的返回类型")
	}

	result := database.DB.Model(&models.RandomImageAPI{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"folder_id":   folderID,
			"return_type": returnType,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(errors.CodeNotFound, "随机API不存在")
	}

	return nil
}

func DeleteRandomAPI(id, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.RandomImageAPI{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(errors.CodeNotFound, "随机API不存在")
	}

	return nil
}

func GetRandomImageByAPIKeyWithConfig(apiKey string) (*models.File, *models.RandomImageAPI, error) {
	var api models.RandomImageAPI
	if err := database.DB.Where("api_key = ?", apiKey).First(&api).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, errors.New(errors.CodeUnauthorized, "无效的API密钥")
		}
		return nil, nil, err
	}

	if !api.IsActive() {
		return nil, nil, errors.New(errors.CodeForbidden, "API密钥已被禁用")
	}

	query := database.DB.Model(&models.File{}).
		Where("user_id = ?", api.UserID).
		Where("access_level = ?", "public")

	if api.FolderID != nil {
		query = query.Where("folder_id = ?", *api.FolderID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, nil, err
	}

	if count == 0 {
		return nil, nil, errors.New(errors.CodeNotFound, "未找到符合条件的图片")
	}

	randomOffset := mathrand.Intn(int(count))

	var file models.File
	if err := query.Offset(randomOffset).Limit(1).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, errors.New(errors.CodeNotFound, "未找到符合条件的图片")
		}
		return nil, nil, err
	}

	go func() {
		now := common.JSONTime(time.Now())
		database.DB.Model(&models.RandomImageAPI{}).
			Where("id = ?", api.ID).
			Updates(map[string]interface{}{
				"call_count":     gorm.Expr("call_count + 1"),
				"last_called_at": &now,
			})
	}()

	return &file, &api, nil
}

func GetFolderNameByID(folderID *string) string {
	if folderID == nil || *folderID == "" {
		return "全部图片"
	}

	var folder models.Folder
	if err := database.DB.Where("id = ?", *folderID).First(&folder).Error; err != nil {
		return fmt.Sprintf("文件夹ID:%s", *folderID)
	}

	return folder.Name
}

func GetFileFullURL(file models.File) (string, string, string) {
	fullURL, fullThumbURL, shortURL := storage.GetFullURLs(file)
	return fullURL, fullThumbURL, shortURL
}
