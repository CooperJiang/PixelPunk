package file

import (
	"pixelpunk/internal/models"

	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/* GetFileForAccess 封装文件查询（按ID或URL/ThumbURL），并对重复文件跳转到原图 */
func GetFileForAccess(filePath string, isThumb bool) (models.File, error) {
	var file models.File
	if err := database.DB.Where("id = ?", filePath).First(&file).Error; err == nil {
		if file.OriginalFileID != "" {
			var origin models.File
			if e := database.DB.Where("id = ?", file.OriginalFileID).First(&origin).Error; e == nil {
				return origin, nil
			}
		}
		return file, nil
	}
	query := "url = ?"
	if isThumb {
		query = "thumb_url = ?"
		candidates := []string{filePath}
		if !strings.HasPrefix(filePath, "thumb/") {
			candidates = append(candidates, "thumb/"+filePath)
		}
		if !strings.HasPrefix(filePath, "thumbnails/") {
			candidates = append(candidates, "thumbnails/"+filePath)
		}
		var lastErr error
		for _, p := range candidates {
			if err := database.DB.Where(query, p).First(&file).Error; err == nil {
				if file.OriginalFileID != "" {
					var origin models.File
					if e := database.DB.Where("id = ?", file.OriginalFileID).First(&origin).Error; e == nil {
						return origin, nil
					}
				}
				return file, nil
			} else {
				lastErr = err
			}
		}
		if lastErr != nil {
			return models.File{}, lastErr
		}
		return file, nil
	}
	if err := database.DB.Where(query, filePath).First(&file).Error; err != nil {
		return models.File{}, err
	}
	if file.OriginalFileID != "" {
		var origin models.File
		if e := database.DB.Where("id = ?", file.OriginalFileID).First(&origin).Error; e == nil {
			return origin, nil
		}
	}
	return file, nil
}

type FileConfig struct {
	UploadDir    string // 上传目录
	ThumbnailDir string // 缩略图目录
}

var DefaultImageConfig = FileConfig{
	UploadDir:    "uploads/files",
	ThumbnailDir: "uploads/thumbnails",
}

var CurrentConfig = DefaultImageConfig

func SetFileConfig(config FileConfig) {
	CurrentConfig = config
}

func generateFileID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:16]
}

/* GetFileAIInfo 获取文件的AI分析信息 */
func GetFileAIInfo(fileID string) (*AIInfoResponse, error) {
	var aiInfo models.FileAIInfo

	if err := database.DB.Where("file_id = ?", fileID).First(&aiInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件AI信息失败")
	}

	var tagNames []string
	if len(aiInfo.Tags) > 0 {
		if err := json.Unmarshal(aiInfo.Tags, &tagNames); err != nil {
			logger.Error("解析标签JSON失败: %v", err)
		}
	}

	if len(tagNames) == 0 {
		var tags []models.GlobalTag
		database.DB.Model(&models.GlobalTag{}).
			Joins("JOIN file_global_tag_relation ON file_global_tag_relation.tag_id = global_tag.id").
			Where("file_global_tag_relation.file_id = ?", fileID).
			Find(&tags)

		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	var colorPalette []string
	if len(aiInfo.ColorPalette) > 0 {
		if err := json.Unmarshal(aiInfo.ColorPalette, &colorPalette); err != nil {
			logger.Error("解析颜色调色板JSON失败: %v", err)
		}
	}

	var nsfwCategories map[string]float64
	if len(aiInfo.NSFWCategories) > 0 {
		if err := json.Unmarshal(aiInfo.NSFWCategories, &nsfwCategories); err != nil {
			logger.Error("解析NSFW分类JSON失败: %v", err)
		}
	}

	response := &AIInfoResponse{
		Description:    aiInfo.Description,
		Tags:           tagNames,
		DominantColor:  aiInfo.DominantColor,
		Resolution:     aiInfo.Resolution,
		IsNSFW:         aiInfo.IsNSFW,
		NSFWScore:      aiInfo.NSFWScore,
		NSFWEvaluation: aiInfo.NSFWEvaluation,

		ColorPalette:   colorPalette,
		AspectRatio:    aiInfo.AspectRatio,
		Composition:    aiInfo.Composition,
		ObjectsCount:   aiInfo.ObjectsCount,
		NSFWCategories: nsfwCategories,
	}

	return response, nil
}

/* GetImageList 获取文件列表 */

/* GetFileDetail 获取单个文件详情 */

/* UpdateFile 更新文件信息 */

/* DeleteFile 删除文件（统一删除入口） */

func deleteFileWithCascade(file *models.File, userID uint) error {
	fileID := file.ID

	var duplicateCount int64
	if file.MD5Hash != "" {
		database.DB.Model(&models.File{}).
			Where("md5_hash = ? AND id != ?", file.MD5Hash, fileID).
			Count(&duplicateCount)
	}

	var referenceCount int64
	database.DB.Model(&models.File{}).
		Where("original_file_id = ?", fileID).
		Count(&referenceCount)

	totalReferences := duplicateCount + referenceCount

	if err := database.DB.Where("file_id = ?", fileID).Delete(&models.FileGlobalTagRelation{}).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除文件标签关联失败")
	}

	if err := database.DB.Where("file_id = ?", fileID).Delete(&models.FileAIInfo{}).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除文件AI分析信息失败")
	}

	if err := database.DB.Unscoped().Where("file_id = ?", fileID).Delete(&models.FileStats{}).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除文件统计失败")
	}

	if err := database.DB.Unscoped().Where("file_id = ?", fileID).Delete(&models.FileVector{}).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除文件向量数据失败")
	}

	var userStats models.UserUsageStats
	if err := database.DB.Where("user_id = ?", userID).First(&userStats).Error; err == nil {
		updates := make(map[string]interface{})
		if userStats.TotalImages > 0 {
			updates["total_images"] = userStats.TotalImages - 1
		}
		if userStats.TotalSize >= file.Size {
			updates["total_size"] = userStats.TotalSize - file.Size
		} else {
			updates["total_size"] = 0
		}
		if len(updates) > 0 {
			if err := database.DB.Model(&models.UserUsageStats{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户统计失败")
			}
		}
	}

	if file.CategoryID != nil && *file.CategoryID > 0 {
		database.DB.Model(&models.FileCategory{}).
			Where("id = ?", *file.CategoryID).
			UpdateColumn("file_count", gorm.Expr("file_count - 1"))
	}

	if err := database.DB.Unscoped().Delete(file).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除文件记录失败")
	}

	go cleanupFileResources(fileID, *file, totalReferences)

	return nil
}

const (
	AvatarUploadDir = "uploads/system/avatar" // 头像上传目录
	FileUploadDir   = "uploads/system/files"  // 文件上传目录
)

/* UploadAvatar 上传用户头像 */

/* ToggleFileAccessLevel 切换文件访问级别 */

/* GetFileLocalPath 从上下文中获取文件信息并返回本地路径 */

/* ServeFile 根据文件存储类型获取访问信息 */

/* ProxyResponse 代理响应 */

/* CleanupPendingDeletionFiles 查找并删除所有被标记为待删除的文件 */

/* GetTotalImageCount 获取系统中的文件总数 */

/* BatchDeleteUserFiles 用户批量删除自己的文件 */

/* ReorderFiles 重新排序文件 */

/* GetExpiredFiles 获取所有过期的文件 */

/* CleanupExpiredFiles 批量清理过期文件 */

/* UploadAdminFile 管理员文件上传 */
