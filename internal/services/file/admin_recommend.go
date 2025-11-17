package file

// Admin recommend/random helpers (no behavior change).

import (
	"encoding/json"
	"math/rand"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/storage"
	"pixelpunk/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// ToggleImageRecommendStatus 切换文件推荐状态
func ToggleFileRecommendStatus(fileID string) (*AdminFileDetailResponse, error) {
	var file models.File
	if err := database.DB.Where("id = ?", fileID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeFileNotFound, "文件不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	file.IsRecommended = !file.IsRecommended
	if err := database.DB.Save(&file).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新文件推荐状态失败")
	}
	var userName string
	if file.UserID > 0 {
		var user models.User
		if err := database.DB.Select("username").Where("id = ?", file.UserID).First(&user).Error; err == nil {
			userName = user.Username
		}
	}
	aiInfo, _ := GetFileAIInfo(file.ID)
	resp := BuildAdminFileDetailResponse(file, 0, userName, aiInfo)
	return &resp, nil
}

// getRandomFileGlobal 全局随机获取推荐文件（回退方法）
func getRandomFileGlobal() (*AdminFileDetailResponse, error) {
	var file models.File
	var totalCount int64
	if err := database.DB.Model(&models.File{}).Where("is_recommended = ? AND access_level = ?", true, AccessPublic).Where("status <> ?", StatusPendingDeletion).Count(&totalCount).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询推荐文件总数失败")
	}
	if totalCount == 0 {
		return nil, errors.New(errors.CodeNotFound, "暂无推荐文件")
	}
	offset := rand.Int63n(totalCount)
	if err := database.DB.Where("is_recommended = ? AND access_level = ?", true, AccessPublic).Where("status <> ?", StatusPendingDeletion).Offset(int(offset)).Limit(1).First(&file).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询推荐文件失败")
	}
	return buildFileResponse(file)
}

// getUserDetailInfo 获取用户详细信息
func getUserDetailInfo(userID uint, excludeFileID string) (*UserDetailInfo, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	daysJoined := int(time.Since(time.Time(user.CreatedAt)).Hours() / 24)

	var totalImages int64
	database.DB.Model(&models.File{}).Where("user_id = ? AND access_level = ? AND status <> ?",
		userID, AccessPublic, StatusPendingDeletion).Count(&totalImages)

	// 统计所有文件总浏览量（从file_stats表获取）
	var totalViews int64
	database.DB.Table("file i").
		Select("COALESCE(SUM(s.views), 0)").
		Joins("LEFT JOIN file_stats s ON i.id = s.file_id").
		Where("i.user_id = ? AND i.status <> ?", userID, StatusPendingDeletion).
		Row().Scan(&totalViews)

	// 先检查用户的总文件数量
	var userTotalCount int64
	database.DB.Model(&models.File{}).Where("user_id = ? AND status <> ?", userID, StatusPendingDeletion).Count(&userTotalCount)

	// 获取其他作品（最多6张公开文件，排除当前文件，按浏览量排序）
	var otherImagesData []models.File
	err := database.DB.Table("file i").
		Select("i.*").
		Joins("LEFT JOIN file_stats s ON i.id = s.file_id").
		Where("i.user_id = ? AND i.access_level = ? AND i.status <> ? AND i.id <> ?",
			userID, AccessPublic, StatusPendingDeletion, excludeFileID).
		Order("COALESCE(s.views, 0) DESC, i.created_at DESC").
		Limit(6).Find(&otherImagesData).Error

	// 转换为FileThumbnail格式
	var otherImages []FileThumbnail
	for _, img := range otherImagesData {
		fullURL, fullThumbURL, _ := storage.GetFullURLs(img)

		otherImages = append(otherImages, FileThumbnail{
			ID:           img.ID,
			DisplayName:  img.DisplayName,
			FullURL:      fullURL,
			FullThumbURL: fullThumbURL,
			Width:        img.Width,
			Height:       img.Height,
			IsNSFW:       img.NSFW,
		})
	}

	var topTags []string
	var aiInfos []models.FileAIInfo

	// 先获取用户所有文件的AI信息
	err = database.DB.Table("file_ai_info iai").
		Select("iai.tags").
		Joins("JOIN file i ON iai.file_id = i.id").
		Where("i.user_id = ? AND i.status <> ?", userID, StatusPendingDeletion).
		Scan(&aiInfos).Error

	if err == nil {
		tagCount := make(map[string]int)
		for _, aiInfo := range aiInfos {
			if len(aiInfo.Tags) > 0 {
				var tags []string
				if err := json.Unmarshal(aiInfo.Tags, &tags); err == nil {
					for _, tag := range tags {
						if tag != "" {
							tagCount[tag]++
						}
					}
				}
			}
		}

		type tagFreq struct {
			tag   string
			count int
		}
		var tagFreqs []tagFreq
		for tag, count := range tagCount {
			tagFreqs = append(tagFreqs, tagFreq{tag, count})
		}

		for i := 0; i < len(tagFreqs); i++ {
			for j := i + 1; j < len(tagFreqs); j++ {
				if tagFreqs[j].count > tagFreqs[i].count {
					tagFreqs[i], tagFreqs[j] = tagFreqs[j], tagFreqs[i]
				}
			}
		}

		for i := 0; i < len(tagFreqs) && i < 5; i++ {
			topTags = append(topTags, tagFreqs[i].tag)
		}
	}

	avatarFullURL := ""
	if user.Avatar != "" {
		avatarFullURL = utils.GetSystemFileURL(user.Avatar)
	}

	return &UserDetailInfo{
		ID:          user.ID,
		Username:    user.Username,
		Avatar:      avatarFullURL,
		CreatedAt:   user.CreatedAt,
		TotalImages: totalImages,
		TotalViews:  totalViews,
		DaysJoined:  daysJoined,
		OtherImages: otherImages,
		TopTags:     topTags,
	}, nil
}

func buildFileResponse(file models.File) (*AdminFileDetailResponse, error) {
	fullURL, fullThumbURL, shortURL := storage.GetFullURLs(file)
	var userName string
	var userInfo *UserDetailInfo

	if file.UserID > 0 {
		var user models.User
		if err := database.DB.Select("username").Where("id = ?", file.UserID).First(&user).Error; err == nil {
			userName = user.Username
		}

		if userDetail, err := getUserDetailInfo(file.UserID, file.ID); err == nil {
			userInfo = userDetail
		}
	}

	aiInfo, _ := GetFileAIInfo(file.ID)
	return &AdminFileDetailResponse{
		ID: file.ID, URL: file.URL, ThumbnailURL: file.ThumbURL, FullURL: fullURL, FullThumbURL: fullThumbURL, ShortURL: shortURL,
		OriginalName: file.OriginalName, DisplayName: file.DisplayName, Size: file.Size, Width: file.Width, Height: file.Height, Format: file.Format,
		AccessLevel: file.AccessLevel, FolderID: file.FolderID, CreatedAt: file.CreatedAt, UpdatedAt: file.UpdatedAt, IsRecommended: file.IsRecommended,
		StorageProviderID: file.StorageProviderID, IsDuplicate: file.IsDuplicate, MD5Hash: file.MD5Hash,
		UserName: userName, UserInfo: userInfo, AIInfo: aiInfo,
	}, nil
}

func GetRandomRecommendedFile() (*AdminFileDetailResponse, error) {
	return getRandomFileGlobal()
}
