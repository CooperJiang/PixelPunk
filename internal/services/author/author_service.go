package author

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/storage"
	"pixelpunk/pkg/utils"
	"time"

	"gorm.io/gorm"
)

/* AuthorHomepage 作者主页信息结构 */
type AuthorHomepage struct {
	Author     AuthorInfo     `json:"author"`
	Folders    []FolderInfo   `json:"folders"`
	Files      []FileInfo     `json:"files"`
	Shares     []ShareInfo    `json:"shares"`
	Stats      AuthorStats    `json:"stats"`
	Pagination PaginationInfo `json:"pagination"`
}

/* AuthorInfo 作者基础信息结构 */
type AuthorInfo struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username"`
	Avatar         string    `json:"avatar"`
	AvatarFullPath string    `json:"avatarFullPath"`
	Bio            string    `json:"bio"`
	Website        string    `json:"website"`
	CreatedAt      time.Time `json:"createdAt"`
	DaysJoined     int       `json:"daysJoined"`
}

/* FolderInfo 文件夹信息结构 */
type FolderInfo struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	FileCount             int       `json:"file_count"`
	Permission            string    `json:"permission"`
	CoverFile             string    `json:"coverFile"`
	CoverFileFullPath     string    `json:"coverFileFullPath"`
	CoverFileThumbURL     string    `json:"coverFileThumbURL"`     // 缩略图URL
	CoverFileFullThumbURL string    `json:"coverFileFullThumbURL"` // 完整缩略图URL
	CreatedAt             time.Time `json:"created_at"`            // 改为created_at
	UpdatedAt             time.Time `json:"updated_at"`            // 添加updated_at字段
	TotalSize             int64     `json:"totalSize"`             // 文件夹总大小（字节）
	TotalSizeFormatted    string    `json:"totalSizeFormatted"`    // 格式化的文件夹大小
}

/* ShareInfo 分享信息结构 */
type ShareInfo struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	ShareKey              string    `json:"shareKey"`
	Views                 int       `json:"views"`
	CreatedAt             time.Time `json:"createdAt"`
	CoverFile             string    `json:"coverFile"`
	CoverFileFullPath     string    `json:"coverFileFullPath"`
	CoverFileThumbURL     string    `json:"coverFileThumbURL"`     // 缩略图URL
	CoverFileFullThumbURL string    `json:"coverFileFullThumbURL"` // 完整缩略图URL
}

/* AuthorStats 作者统计信息结构 */
type AuthorStats struct {
	TotalFiles   int `json:"totalFiles"`
	TotalViews   int `json:"totalViews"`
	TotalFolders int `json:"totalFolders"`
	TotalShares  int `json:"totalShares"`
}

/* FolderContent 文件夹内容结构 */
type FolderContent struct {
	Folder     FolderInfo     `json:"folder"`
	Files      []FileInfo     `json:"files"`
	SubFolders []FolderInfo   `json:"subFolders"`
	Pagination PaginationInfo `json:"pagination"`
}

/* FileAIInfo AI信息结构 */
type FileAIInfo struct {
	Description    string   `json:"description"`
	Tags           []string `json:"tags"`
	DominantColor  string   `json:"dominant_color"`
	Resolution     string   `json:"resolution"`
	IsNSFW         bool     `json:"is_nsfw"`
	NSFWScore      float64  `json:"nsfw_score"`
	NSFWEvaluation string   `json:"nsfw_evaluation"`
}

/* FileInfo 文件信息结构 */
type FileInfo struct {
	ID           string      `json:"id"`
	FileName     string      `json:"file_name"`     // 改为file_name
	OriginalName string      `json:"original_name"` // 改为original_name
	DisplayName  string      `json:"display_name"`  // 添加display_name
	URL          string      `json:"url"`
	FullURL      string      `json:"full_url"`       // 改为full_url
	ThumbURL     string      `json:"thumb_url"`      // 改为thumb_url
	FullThumbURL string      `json:"full_thumb_url"` // 改为full_thumb_url
	Size         int64       `json:"size"`
	Format       string      `json:"format"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	Views        int         `json:"views"`                 // 浏览次数
	AccessLevel  string      `json:"access_level"`          // 添加access_level字段
	IsDuplicate  bool        `json:"is_duplicate"`          // 添加is_duplicate字段
	Description  string      `json:"description,omitempty"` // 顶层描述字段（从AI信息中提取）
	AIInfo       *FileAIInfo `json:"ai_info,omitempty"`     // AI信息
	CreatedAt    time.Time   `json:"created_at"`            // 改为created_at
	UpdatedAt    time.Time   `json:"updated_at"`            // 添加updated_at字段
}

/* PaginationInfo 分页信息结构 */
type PaginationInfo struct {
	CurrentPage int   `json:"currentPage"`
	PerPage     int   `json:"perPage"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"lastPage"`
}

func formatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

/* GetAuthorHomepage 获取作者主页信息 */
func GetAuthorHomepage(authorID uint) (*AuthorHomepage, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, authorID).Error; err != nil {
		return nil, errors.New(errors.CodeUserNotFound, "作者不存在")
	}

	authorInfo := AuthorInfo{
		ID:         user.ID,
		Username:   user.Username,
		Avatar:     user.Avatar,
		Bio:        user.Bio,
		Website:    user.Website,
		CreatedAt:  time.Time(user.CreatedAt),
		DaysJoined: int(time.Since(time.Time(user.CreatedAt)).Hours() / 24),
	}

	if user.Avatar != "" {
		authorInfo.AvatarFullPath = utils.GetSystemFileURL(user.Avatar)
	}

	var folders []models.Folder
	if err := db.Where("user_id = ? AND permission = 'public' AND parent_id = ''", authorID).Find(&folders).Error; err != nil {
		return nil, errors.New(errors.CodeInternal, "获取文件夹信息失败")
	}

	folderInfos := make([]FolderInfo, 0, len(folders))
	for _, folder := range folders {
		var imageCount int64
		if err := db.Model(&models.File{}).
			Where("folder_id = ? AND access_level = 'public'", folder.ID).
			Count(&imageCount).Error; err != nil {
			imageCount = 0
		}

		var totalSize int64
		if err := db.Model(&models.File{}).
			Where("folder_id = ? AND access_level = 'public'", folder.ID).
			Select("COALESCE(SUM(size), 0)").Scan(&totalSize).Error; err != nil {
			totalSize = 0
		}

		var coverImage models.File
		coverImagePath := ""
		coverImageFullPath := ""
		coverImageThumbURL := ""
		coverImageFullThumbURL := ""
		if db.Where("folder_id = ? AND access_level = 'public'", folder.ID).
			Order("created_at DESC").First(&coverImage).Error == nil {
			coverImageFullPath, coverImageFullThumbURL, _ = storage.GetFullURLs(coverImage)
			coverImagePath = coverImage.URL
			coverImageThumbURL = coverImage.ThumbURL
		}

		folderInfos = append(folderInfos, FolderInfo{
			ID:                    folder.ID,
			Name:                  folder.Name,
			FileCount:             int(imageCount),
			Permission:            folder.Permission,
			CoverFile:             coverImagePath,
			CoverFileFullPath:     coverImageFullPath,
			CoverFileThumbURL:     coverImageThumbURL,
			CoverFileFullThumbURL: coverImageFullThumbURL,
			CreatedAt:             time.Time(folder.CreatedAt),
			UpdatedAt:             time.Time(folder.UpdatedAt),
			TotalSize:             totalSize,
			TotalSizeFormatted:    formatFileSize(totalSize),
		})
	}

	var shares []models.Share
	if err := db.Where("user_id = ? AND status = 1 AND (password = '' OR password IS NULL)", authorID).
		Order("created_at DESC").Limit(10).Find(&shares).Error; err != nil {
		return nil, errors.New(errors.CodeInternal, "获取分享信息失败")
	}

	shareInfos := make([]ShareInfo, 0, len(shares))
	for _, share := range shares {
		var coverImage models.File
		coverImagePath := ""
		coverImageFullPath := ""
		coverImageThumbURL := ""
		coverImageFullThumbURL := ""
		if db.Table("share_item si").
			Joins("JOIN file i ON si.item_id = i.id").
			Where("si.share_id = ? AND si.item_type = 'file' AND i.access_level = 'public' AND i.status <> ?", share.ID, "pending_deletion").
			Order("i.created_at DESC").Select("i.*").First(&coverImage).Error == nil {
			coverImageFullPath, coverImageFullThumbURL, _ = storage.GetFullURLs(coverImage)
			coverImagePath = coverImage.URL
			coverImageThumbURL = coverImage.ThumbURL
		}

		shareInfos = append(shareInfos, ShareInfo{
			ID:                    share.ID,
			Name:                  share.Name,
			Description:           share.Description,
			ShareKey:              share.ShareKey,
			Views:                 share.CurrentViews,
			CreatedAt:             time.Time(share.CreatedAt),
			CoverFile:             coverImagePath,
			CoverFileFullPath:     coverImageFullPath,
			CoverFileThumbURL:     coverImageThumbURL,
			CoverFileFullThumbURL: coverImageFullThumbURL,
		})
	}

	var userStats models.UserUsageStats
	var totalViews int64 = 0
	if err := db.Where("user_id = ?", authorID).First(&userStats).Error; err == nil {
		totalViews = userStats.TotalViews
	}

	var rootFiles []models.File
	var totalRootFiles int64

	if err := db.Model(&models.File{}).
		Where("user_id = ? AND (folder_id = '' OR folder_id IS NULL) AND access_level = 'public' AND status <> ?",
			authorID, "pending_deletion").
		Count(&totalRootFiles).Error; err != nil {
		totalRootFiles = 0
	}

	page := 1
	size := 50
	offset := (page - 1) * size

	if err := db.Where("user_id = ? AND (folder_id = '' OR folder_id IS NULL) AND access_level = 'public' AND status <> ?",
		authorID, "pending_deletion").
		Order("created_at DESC").Offset(offset).Limit(size).
		Find(&rootFiles).Error; err != nil {
		rootFiles = []models.File{}
	}

	fileInfos := make([]FileInfo, 0, len(rootFiles))
	for _, file := range rootFiles {
		fullPath, fullThumbURL, _ := storage.GetFullURLs(file)

		var stats models.FileStats
		views := 0
		if err := db.Where("file_id = ?", file.ID).First(&stats).Error; err == nil {
			views = int(stats.Views)
		}

		aiInfo, _ := getFileAIInfo(file.ID)

		// 从AI信息中提取描述到顶层
		description := ""
		if aiInfo != nil && aiInfo.Description != "" {
			description = aiInfo.Description
		}

		fileInfos = append(fileInfos, FileInfo{
			ID:           file.ID,
			FileName:     file.FileName,
			OriginalName: file.OriginalName,
			DisplayName:  file.DisplayName,
			URL:          file.URL,
			FullURL:      fullPath,
			ThumbURL:     file.ThumbURL,
			FullThumbURL: fullThumbURL,
			Size:         file.Size,
			Format:       file.Format,
			Width:        file.Width,
			Height:       file.Height,
			Views:        views,
			AccessLevel:  file.AccessLevel,
			IsDuplicate:  file.IsDuplicate,
			Description:  description,
			AIInfo:       aiInfo,
			CreatedAt:    time.Time(file.CreatedAt),
			UpdatedAt:    time.Time(file.UpdatedAt),
		})
	}

	pagination := PaginationInfo{
		CurrentPage: page,
		PerPage:     size,
		Total:       totalRootFiles,
		LastPage:    int((totalRootFiles + int64(size) - 1) / int64(size)),
	}

	// 统计用户所有公开的文件数量
	var totalFiles int64
	if err := db.Model(&models.File{}).
		Where("user_id = ? AND access_level = 'public' AND status <> ?", authorID, "pending_deletion").
		Count(&totalFiles).Error; err != nil {
		totalFiles = 0
	}

	// 统计用户所有公开的文件夹数量
	var totalFolders int64
	if err := db.Model(&models.Folder{}).
		Where("user_id = ? AND permission = 'public'", authorID).
		Count(&totalFolders).Error; err != nil {
		totalFolders = 0
	}

	stats := AuthorStats{
		TotalFiles:   int(totalFiles),
		TotalViews:   int(totalViews),
		TotalFolders: int(totalFolders),
		TotalShares:  len(shareInfos),
	}

	return &AuthorHomepage{
		Author:     authorInfo,
		Folders:    folderInfos,
		Files:      fileInfos,
		Shares:     shareInfos,
		Stats:      stats,
		Pagination: pagination,
	}, nil
}

/* GetAuthorFolder 获取作者特定文件夹内容 */
func GetAuthorFolder(authorID uint, folderID string, page, size int) (*FolderContent, error) {
	db := database.GetDB()

	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ? AND permission = 'public'", folderID, authorID).
		First(&folder).Error; err != nil {
		return nil, errors.New(errors.CodeNotFound, "文件夹不存在或无权访问")
	}

	var imageCount int64
	if err := db.Model(&models.File{}).
		Where("folder_id = ? AND access_level = 'public' AND status <> ?", folderID, "pending_deletion").
		Count(&imageCount).Error; err != nil {
		imageCount = 0 // 如果查询失败，设为0
	}

	var totalSize int64
	if err := db.Model(&models.File{}).
		Where("folder_id = ? AND access_level = 'public' AND status <> ?", folderID, "pending_deletion").
		Select("COALESCE(SUM(size), 0)").Scan(&totalSize).Error; err != nil {
		totalSize = 0
	}

	var coverImage models.File
	coverImagePath := ""
	coverImageFullPath := ""
	coverImageThumbURL := ""
	coverImageFullThumbURL := ""
	if db.Where("folder_id = ? AND access_level = 'public' AND status <> ?", folderID, "pending_deletion").
		Order("created_at DESC").First(&coverImage).Error == nil {
		coverImageFullPath, coverImageFullThumbURL, _ = storage.GetFullURLs(coverImage)
		coverImagePath = coverImage.URL
		coverImageThumbURL = coverImage.ThumbURL
	}

	folderInfo := FolderInfo{
		ID:                    folder.ID,
		Name:                  folder.Name,
		FileCount:             int(imageCount),
		Permission:            folder.Permission,
		CoverFile:             coverImagePath,
		CoverFileFullPath:     coverImageFullPath,
		CoverFileThumbURL:     coverImageThumbURL,
		CoverFileFullThumbURL: coverImageFullThumbURL,
		CreatedAt:             time.Time(folder.CreatedAt),
		UpdatedAt:             time.Time(folder.UpdatedAt),
		TotalSize:             totalSize,
		TotalSizeFormatted:    formatFileSize(totalSize),
	}

	var subFolders []models.Folder
	if err := db.Where("parent_id = ? AND permission = 'public'", folderID).
		Find(&subFolders).Error; err != nil {
		return nil, errors.New(errors.CodeInternal, "获取子文件夹失败")
	}

	subFolderInfos := make([]FolderInfo, 0, len(subFolders))
	for _, subFolder := range subFolders {
		var subImageCount int64
		if err := db.Model(&models.File{}).
			Where("folder_id = ? AND access_level = 'public' AND status <> ?", subFolder.ID, "pending_deletion").
			Count(&subImageCount).Error; err != nil {
			subImageCount = 0 // 如果查询失败，设为0
		}

		var subTotalSize int64
		if err := db.Model(&models.File{}).
			Where("folder_id = ? AND access_level = 'public' AND status <> ?", subFolder.ID, "pending_deletion").
			Select("COALESCE(SUM(size), 0)").Scan(&subTotalSize).Error; err != nil {
			subTotalSize = 0
		}

		var subCoverImage models.File
		subCoverImagePath := ""
		subCoverImageFullPath := ""
		subCoverImageThumbURL := ""
		subCoverImageFullThumbURL := ""
		if db.Where("folder_id = ? AND access_level = 'public' AND status <> ?", subFolder.ID, "pending_deletion").
			Order("created_at DESC").First(&subCoverImage).Error == nil {
			subCoverImageFullPath, subCoverImageFullThumbURL, _ = storage.GetFullURLs(subCoverImage)
			subCoverImagePath = subCoverImage.URL
			subCoverImageThumbURL = subCoverImage.ThumbURL
		}

		subFolderInfos = append(subFolderInfos, FolderInfo{
			ID:                    subFolder.ID,
			Name:                  subFolder.Name,
			FileCount:             int(subImageCount),
			Permission:            subFolder.Permission,
			CoverFile:             subCoverImagePath,
			CoverFileFullPath:     subCoverImageFullPath,
			CoverFileThumbURL:     subCoverImageThumbURL,
			CoverFileFullThumbURL: subCoverImageFullThumbURL,
			CreatedAt:             time.Time(subFolder.CreatedAt),
			UpdatedAt:             time.Time(subFolder.UpdatedAt),
			TotalSize:             subTotalSize,
			TotalSizeFormatted:    formatFileSize(subTotalSize),
		})
	}

	var images []models.File
	var total int64

	if err := db.Model(&models.File{}).
		Where("folder_id = ? AND access_level = 'public' AND status <> ?", folderID, "pending_deletion").
		Count(&total).Error; err != nil {
		total = 0
	}

	offset := (page - 1) * size
	if err := db.Where("folder_id = ? AND access_level = 'public' AND status <> ?", folderID, "pending_deletion").
		Order("created_at DESC").Offset(offset).Limit(size).
		Find(&images).Error; err != nil {
		return nil, errors.New(errors.CodeInternal, "获取文件列表失败")
	}

	imageInfos := make([]FileInfo, 0, len(images))
	for _, file := range images {
		fullPath, fullThumbURL, _ := storage.GetFullURLs(file)

		var stats models.FileStats
		views := 0
		if err := db.Where("file_id = ?", file.ID).First(&stats).Error; err == nil {
			views = int(stats.Views)
		}

		aiInfo, _ := getFileAIInfo(file.ID)

		// 从AI信息中提取描述到顶层
		description := ""
		if aiInfo != nil && aiInfo.Description != "" {
			description = aiInfo.Description
		}

		imageInfos = append(imageInfos, FileInfo{
			ID:           file.ID,
			FileName:     file.FileName,
			OriginalName: file.OriginalName,
			DisplayName:  file.DisplayName,
			URL:          file.URL,
			FullURL:      fullPath,
			ThumbURL:     file.ThumbURL,
			FullThumbURL: fullThumbURL,
			Size:         file.Size,
			Format:       file.Format,
			Width:        file.Width,
			Height:       file.Height,
			Views:        views,
			AccessLevel:  file.AccessLevel,
			IsDuplicate:  file.IsDuplicate,
			Description:  description,
			AIInfo:       aiInfo,
			CreatedAt:    time.Time(file.CreatedAt),
			UpdatedAt:    time.Time(file.UpdatedAt),
		})
	}

	pagination := PaginationInfo{
		CurrentPage: page,
		PerPage:     size,
		Total:       total,
		LastPage:    int((total + int64(size) - 1) / int64(size)),
	}

	return &FolderContent{
		Folder:     folderInfo,
		Files:      imageInfos,
		SubFolders: subFolderInfos,
		Pagination: pagination,
	}, nil
}

/* getFileAIInfo 获取文件AI信息 */
func getFileAIInfo(fileID string) (*FileAIInfo, error) {
	db := database.GetDB()
	var aiInfo models.FileAIInfo
	if err := db.Where("file_id = ?", fileID).First(&aiInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// 解析 tags JSON
	var tagNames []string
	if len(aiInfo.Tags) > 0 {
		_ = json.Unmarshal(aiInfo.Tags, &tagNames)
	}

	if len(tagNames) == 0 {
		var tags []models.GlobalTag
		db.Model(&models.GlobalTag{}).
			Joins("JOIN file_global_tag_relation ON file_global_tag_relation.tag_id = global_tag.id").
			Where("file_global_tag_relation.file_id = ?", fileID).
			Find(&tags)
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	return &FileAIInfo{
		Description:    aiInfo.Description,
		Tags:           tagNames,
		DominantColor:  aiInfo.DominantColor,
		Resolution:     aiInfo.Resolution,
		IsNSFW:         aiInfo.IsNSFW,
		NSFWScore:      aiInfo.NSFWScore,
		NSFWEvaluation: aiInfo.NSFWEvaluation,
	}, nil
}
