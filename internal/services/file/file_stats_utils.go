package file

/* File stats helpers (no behavior change). */

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"gorm.io/gorm"
)

/* GetFileStats 获取文件统计数据 */
func GetFileStats(userID uint, fileID, period string) (map[string]interface{}, error) {
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeFileNotFound, "文件不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	var stats models.FileStats
	if err := database.DB.Where("file_id = ?", fileID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			stats = models.FileStats{FileID: fileID}
		} else {
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件统计失败")
		}
	}
	referrers := []map[string]interface{}{{"domain": "example.com", "views": 532}, {"domain": "anothersite.org", "views": 217}}
	countries := []map[string]interface{}{{"code": "CN", "views": 436}, {"code": "US", "views": 285}}
	return map[string]interface{}{"views": stats.Views, "bandwidth": stats.Bandwidth, "referrers": referrers, "countries": countries}, nil
}
