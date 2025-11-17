package share

import (
	stderrors "errors"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/storage"
	"pixelpunk/pkg/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GetShareForView(shareKey string, folderID string) (map[string]interface{}, error) {
	share, err := GetShareByKey(shareKey)
	if err != nil {
		return nil, err
	}

	shareItems, err := GetShareItems(share.ID)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := database.DB.Select("id, username, avatar").Where("id = ?", share.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	var currentFolder *models.Folder
	var parentFolderID string

	if folderID != "" && folderID != "0" {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", folderID, share.UserID).First(&folder).Error; err != nil {
			return nil, errors.New(errors.CodeFolderNotFound, "指定的文件夹不存在或无权访问")
		}

		isFolderInShare := false

		for _, item := range shareItems {
			if item.ItemType == common.ShareItemTypeFolder && item.ItemID == folderID {
				isFolderInShare = true
				break
			}
		}

		if !isFolderInShare {
			sharedFolderIDs := []string{}
			for _, item := range shareItems {
				if item.ItemType == common.ShareItemTypeFolder {
					sharedFolderIDs = append(sharedFolderIDs, item.ItemID)
				}
			}

			if len(sharedFolderIDs) > 0 {
				var pathFolder models.Folder = folder
				for {
					if pathFolder.ParentID == "" || pathFolder.ParentID == "0" {
						break
					}

					for _, sharedID := range sharedFolderIDs {
						if pathFolder.ParentID == sharedID {
							isFolderInShare = true
							break
						}
					}

					if isFolderInShare {
						break
					}

					var parentFolder models.Folder
					if err := database.DB.Where("id = ? AND user_id = ?", pathFolder.ParentID, share.UserID).First(&parentFolder).Error; err != nil {
						break
					}
					pathFolder = parentFolder
				}
			}
		}

		if !isFolderInShare {
			return nil, errors.New(errors.CodeValidationFailed, "该文件夹不包含在分享内容中")
		}

		currentFolder = &folder
		parentFolderID = folder.ParentID
	}

	folders := []models.Folder{}
	files := []map[string]interface{}{}

	if currentFolder != nil {
		if err := database.DB.Where("parent_id = ? AND user_id = ?", folderID, share.UserID).Find(&folders).Error; err != nil {
			return nil, err
		}

		var folderImages []models.File
		if err := database.DB.Preload("AIInfo").Where("folder_id = ? AND user_id = ?", folderID, share.UserID).
			Where("status <> ?", "pending_deletion").
			Find(&folderImages).Error; err != nil {
			return nil, err
		}

		for _, file := range folderImages {
			fullURL, fullThumbURL, _ := storage.GetFullURLs(file)

			if fullURL != "" {
				if strings.Contains(fullURL, "?") {
					fullURL = fullURL + "&share=" + shareKey
				} else {
					fullURL = fullURL + "?share=" + shareKey
				}
			}

			if fullThumbURL != "" {
				if strings.Contains(fullThumbURL, "?") {
					fullThumbURL = fullThumbURL + "&share=" + shareKey
				} else {
					fullThumbURL = fullThumbURL + "?share=" + shareKey
				}
			}

			fileMap := map[string]interface{}{
				"id":             file.ID,
				"display_name":   file.DisplayName,
				"description":    file.Description,
				"url":            file.URL,
				"thumb_url":      file.ThumbURL,
				"size":           file.Size,
				"size_formatted": file.SizeFormatted,
				"width":          file.Width,
				"height":         file.Height,
				"format":         file.Format,
				"mime":           file.Mime,
				"created_at":     file.CreatedAt,
				"updated_at":     file.UpdatedAt,
				"full_url":       fullURL,            // 添加完整URL
				"full_thumb_url": fullThumbURL,       // 添加完整缩略图URL
				"resolution":     file.Resolution,    // 添加分辨率信息
				"is_recommended": file.IsRecommended, // 添加推荐标记
				"ai_info":        file.AIInfo,        // 添加AI信息
			}

			var tags []map[string]interface{}
			var globalTags []models.GlobalTag
			if err := database.DB.Model(&models.GlobalTag{}).
				Joins("JOIN file_global_tag_relation ON file_global_tag_relation.tag_id = global_tag.id").
				Where("file_global_tag_relation.file_id = ?", file.ID).
				Find(&globalTags).Error; err == nil {
				for _, globalTag := range globalTags {
					tags = append(tags, map[string]interface{}{
						"id":         globalTag.ID,
						"name":       globalTag.Name,
						"created_at": globalTag.CreatedAt,
					})
				}
			}
			fileMap["tags"] = tags

			files = append(files, fileMap)
		}
	} else {
		for _, item := range shareItems {
			if item.ItemType == common.ShareItemTypeFolder {
				var folder models.Folder
				if err := database.DB.Where("id = ? AND user_id = ?", item.ItemID, share.UserID).First(&folder).Error; err == nil {
					folders = append(folders, folder)
				}
			} else if item.ItemType == common.ShareItemTypeFile {
				var file models.File
				if err := database.DB.Preload("AIInfo").Where("id = ? AND user_id = ?", item.ItemID, share.UserID).
					Where("status <> ?", "pending_deletion").
					First(&file).Error; err == nil {
					fullURL, fullThumbURL, _ := storage.GetFullURLs(file)

					if fullURL != "" {
						if strings.Contains(fullURL, "?") {
							fullURL = fullURL + "&share=" + shareKey
						} else {
							fullURL = fullURL + "?share=" + shareKey
						}
					}

					if fullThumbURL != "" {
						if strings.Contains(fullThumbURL, "?") {
							fullThumbURL = fullThumbURL + "&share=" + shareKey
						} else {
							fullThumbURL = fullThumbURL + "?share=" + shareKey
						}
					}

					fileMap := map[string]interface{}{
						"id":             file.ID,
						"display_name":   file.DisplayName,
						"description":    file.Description,
						"url":            file.URL,
						"thumb_url":      file.ThumbURL,
						"size":           file.Size,
						"size_formatted": file.SizeFormatted,
						"width":          file.Width,
						"height":         file.Height,
						"format":         file.Format,
						"mime":           file.Mime,
						"created_at":     file.CreatedAt,
						"updated_at":     file.UpdatedAt,
						"full_url":       fullURL,            // 添加完整URL
						"full_thumb_url": fullThumbURL,       // 添加完整缩略图URL
						"resolution":     file.Resolution,    // 添加分辨率信息
						"is_recommended": file.IsRecommended, // 添加推荐标记
						"ai_info":        file.AIInfo,        // 添加AI信息
					}

					var tags []map[string]interface{}
					var globalTags []models.GlobalTag
					if err := database.DB.Model(&models.GlobalTag{}).
						Joins("JOIN file_global_tag_relation ON file_global_tag_relation.tag_id = global_tag.id").
						Where("file_global_tag_relation.file_id = ?", file.ID).
						Find(&globalTags).Error; err == nil {
						for _, globalTag := range globalTags {
							tags = append(tags, map[string]interface{}{
								"id":         globalTag.ID,
								"name":       globalTag.Name,
								"created_at": globalTag.CreatedAt,
							})
						}
					}
					fileMap["tags"] = tags

					files = append(files, fileMap)
				}
			}
		}
	}

	result := map[string]interface{}{
		"share": map[string]interface{}{
			"id":                     share.ID,
			"name":                   share.Name,
			"description":            share.Description,
			"created_at":             share.CreatedAt,
			"expired_at":             share.ExpiredAt,
			"expired_days":           share.ExpiredDays,
			"current_views":          share.CurrentViews,
			"max_views":              share.MaxViews,
			"has_password":           share.Password != "",
			"collect_visitor_info":   share.CollectVisitorInfo,
			"notification_on_access": share.NotificationOnAccess,
		},
		"user": map[string]interface{}{
			"username": user.Username,
			"avatar":   user.Avatar,
		},
		"folders":        folders,
		"files":          files,
		"current_folder": currentFolder,
		"parent_id":      parentFolderID,
	}

	return result, nil
}

/* GenerateAccessToken 生成临时访问令牌 */
func GenerateAccessToken(shareKey string, password string, clientIP, userAgent string) (string, error) {
	share, err := GetShareByKey(shareKey)
	if err != nil {
		return "", err
	}

	if share.Password != password {
		return "", errors.New(errors.CodeWrongPassword, "密码错误")
	}

	tokenStr := utils.GenerateRandomString(32)

	expiredAt := time.Now().Add(30 * time.Minute)
	tokenExpiredAt := common.JSONTime(expiredAt)

	token := models.ShareAccessToken{
		ID:        generateID(),
		ShareID:   share.ID,
		ShareKey:  shareKey,
		Token:     tokenStr,
		ExpiredAt: tokenExpiredAt,
		ClientIP:  clientIP,
		UserAgent: userAgent,
	}

	if err := database.DB.Create(&token).Error; err != nil {
		return "", err
	}

	return tokenStr, nil
}

/* ValidateAccessToken 验证临时访问令牌 */
func ValidateAccessToken(shareKey string, accessToken string) (bool, error) {
	var share models.Share
	if err := database.DB.Where("share_key = ? AND status = ?", shareKey, common.ShareStatusNormal).First(&share).Error; err != nil {
		return false, errors.New(errors.CodeNotFound, "分享不存在或已失效")
	}

	var token models.ShareAccessToken
	err := database.DB.Where("token = ? AND share_key = ? AND expired_at > ?",
		accessToken, shareKey, time.Now()).
		First(&token).Error

	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // 令牌不存在或已过期
		}
		return false, err
	}

	if token.IsAdmin {
		return true, nil
	}

	if !share.IsAccessible() {
		return false, errors.New(errors.CodeValidationFailed, "分享已过期或已达最大访问次数")
	}

	return true, nil
}

/* CleanExpiredTokens 清理过期的访问令牌 */
func CleanExpiredTokens() (int64, error) {
	result := database.DB.Where("expired_at < ?", time.Now()).Delete(&models.ShareAccessToken{})
	return result.RowsAffected, result.Error
}

/* GenerateAdminAccessToken 管理员生成访问令牌 */
func GenerateAdminAccessToken(shareID string, adminID uint) (string, error) {
	var share models.Share
	if err := database.DB.Where("id = ?", shareID).First(&share).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New(errors.CodeNotFound, "分享不存在")
		}
		return "", err
	}

	tokenStr := utils.GenerateRandomString(32)

	expiredAt := time.Now().Add(24 * time.Hour)
	tokenExpiredAt := common.JSONTime(expiredAt)

	token := models.ShareAccessToken{
		ID:        generateID(),
		ShareID:   shareID,
		ShareKey:  share.ShareKey,
		Token:     tokenStr,
		ExpiredAt: tokenExpiredAt,
		ClientIP:  "admin",
		UserAgent: fmt.Sprintf("Admin:%d", adminID),
		IsAdmin:   true, // 标记为管理员令牌
	}

	if err := database.DB.Create(&token).Error; err != nil {
		return "", err
	}

	return tokenStr, nil
}

/* ValidateSharedImageAccess 验证文件是否属于分享内容 */
func ValidateSharedFileAccess(shareID, fileID string) (bool, error) {
	var share models.Share
	if err := database.DB.Where("id = ?", shareID).First(&share).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New(errors.CodeNotFound, "分享不存在")
		}
		return false, err
	}

	if share.IsExpired() {
		return false, errors.New(errors.CodeShareExpired, "分享已过期")
	}

	if !share.IsAccessible() {
		return false, errors.New(errors.CodeShareAccessExceeded, "分享访问次数已超限")
	}

	var shareItem models.ShareItem
	err := database.DB.Where("share_id = ? AND item_type = ? AND item_id = ?",
		shareID, "file", fileID).First(&shareItem).Error

	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return validateFileInSharedFolder(shareID, fileID)
		}
		return false, err
	}

	return true, nil
}

func validateFileInSharedFolder(shareID, fileID string) (bool, error) {
	var file models.File
	if err := database.DB.Where("id = ?", fileID).First(&file).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New(errors.CodeFileNotFound, "文件不存在")
		}
		return false, err
	}

	var folderItems []models.ShareItem
	if err := database.DB.Where("share_id = ? AND item_type = ?", shareID, "folder").Find(&folderItems).Error; err != nil {
		return false, err
	}

	for _, folderItem := range folderItems {
		if isFileInFolder(file.FolderID, folderItem.ItemID) {
			return true, nil
		}
	}

	return false, nil
}

func isFileInFolder(imageFolderID string, targetFolderID string) bool {
	if imageFolderID == "" {
		return false
	}

	if imageFolderID == targetFolderID {
		return true
	}

	var folder models.Folder
	if err := database.DB.Where("id = ?", imageFolderID).First(&folder).Error; err != nil {
		return false
	}

	if folder.ParentID != "" {
		return isFileInFolder(folder.ParentID, targetFolderID)
	}

	return false
}
