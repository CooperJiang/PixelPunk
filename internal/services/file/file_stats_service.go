package file

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"gorm.io/gorm"
)

/* InitFileStats 初始化文件统计记录 */
func InitFileStats(tx *gorm.DB, fileID string) error {
	stats := models.FileStats{
		FileID:    fileID,
		Views:     0,
		Bandwidth: 0,
	}

	if err := tx.Create(&stats).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBCreateFailed, "创建文件统计记录失败")
	}

	return nil
}

/* FetchFileStats 获取文件统计 */
func FetchFileStats(fileID string) (*models.FileStats, error) {
	var stats models.FileStats

	err := database.DB.Where("file_id = ?", fileID).First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			stats = models.FileStats{
				FileID: fileID,
			}
			if err := database.DB.Create(&stats).Error; err != nil {
				return nil, errors.Wrap(err, errors.CodeDBCreateFailed, "创建文件统计记录失败")
			}
		} else {
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取文件统计失败")
		}
	}

	return &stats, nil
}

/* UpdateViews 更新文件浏览次数 */
func UpdateViews(fileID string) error {
	result := database.DB.Model(&models.FileStats{}).
		Where("file_id = ?", fileID).
		UpdateColumn("views", gorm.Expr("views + 1"))

	if result.Error != nil {
		return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "更新文件浏览次数失败")
	}

	if result.RowsAffected == 0 {
		stats := models.FileStats{
			FileID: fileID,
			Views:  1,
		}
		if err := database.DB.Create(&stats).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBCreateFailed, "创建文件统计记录失败")
		}
	}

	return nil
}

/* UpdateBandwidth 更新文件带宽使用量 */
func UpdateBandwidth(fileID string, bandwidth int64) error {
	result := database.DB.Model(&models.FileStats{}).
		Where("file_id = ?", fileID).
		UpdateColumn("bandwidth", gorm.Expr("bandwidth + ?", bandwidth))

	if result.Error != nil {
		return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "更新文件带宽使用量失败")
	}

	if result.RowsAffected == 0 {
		stats := models.FileStats{
			FileID:    fileID,
			Bandwidth: bandwidth,
		}
		if err := database.DB.Create(&stats).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBCreateFailed, "创建文件统计记录失败")
		}
	}

	return nil
}
