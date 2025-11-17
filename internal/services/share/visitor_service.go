package share

import (
	"pixelpunk/internal/controllers/share/dto"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"strings"
	"time"
)

/* SaveVisitorInfo 保存访客信息 */
func SaveVisitorInfo(shareKey string, visitorInfo *dto.VisitorInfoDTO, ip, userAgent, referer string) error {
	var share models.Share
	if err := database.DB.Where("share_key = ? AND status = ?", shareKey, common.ShareStatusNormal).First(&share).Error; err != nil {
		return errors.New(errors.CodeNotFound, "分享不存在或已失效")
	}

	var existingVisitor models.ShareVisitorInfo
	err := database.DB.Where("share_id = ? AND visitor_email = ?", share.ID, visitorInfo.Email).First(&existingVisitor).Error

	if err == nil {
		updates := map[string]interface{}{
			"visitor_name":  visitorInfo.Name, // 更新姓名，以防用户修改
			"visit_count":   existingVisitor.VisitCount + 1,
			"last_visit_at": common.JSONTime(time.Now()),
			"ip_address":    ip,
			"user_agent":    userAgent,
			"referer":       referer,
		}

		if err := database.DB.Model(&existingVisitor).Updates(updates).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新访客信息失败")
		}

		return nil
	}

	newVisitor := models.ShareVisitorInfo{
		ID:           generateID(),
		ShareID:      share.ID,
		ShareKey:     shareKey,
		VisitorName:  visitorInfo.Name,
		VisitorEmail: visitorInfo.Email,
		IPAddress:    ip,
		UserAgent:    userAgent,
		Referer:      referer,
		VisitCount:   1,
		LastVisitAt:  common.JSONTime(time.Now()),
	}

	if err := database.DB.Create(&newVisitor).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBCreateFailed, "保存访客信息失败")
	}

	return nil
}

/* GetShareVisitors 获取分享的访客信息列表 */
func GetShareVisitors(shareID string, query *dto.VisitorQueryDTO) ([]models.ShareVisitorInfo, int64, error) {
	db := database.DB.Model(&models.ShareVisitorInfo{}).Where("share_id = ?", shareID)

	if query != nil && query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("(visitor_name LIKE ? OR visitor_email LIKE ? OR ip_address LIKE ?)",
			keyword, keyword, keyword)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取访客信息总数失败")
	}

	page := 1
	size := common.DefaultPageSize

	if query != nil {
		if query.Page > 0 {
			page = query.Page
		}
		if query.Size > 0 && query.Size <= common.MaxPageSize {
			size = query.Size
		}
	}

	orderBy := "created_at DESC"
	if query != nil && strings.TrimSpace(query.OrderBy) != "" {
		orderBy = query.OrderBy
	}

	var visitors []models.ShareVisitorInfo
	if err := db.Order(orderBy).Offset((page - 1) * size).Limit(size).Find(&visitors).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询访客信息列表失败")
	}

	return visitors, total, nil
}

/* DeleteVisitorInfo 删除访客信息 */
func DeleteVisitorInfo(shareID string, visitorID string, userID uint) error {
	var share models.Share
	if err := database.DB.Where("id = ? AND user_id = ?", shareID, userID).First(&share).Error; err != nil {
		return errors.New(errors.CodeNotFound, "分享不存在或您无权操作")
	}

	var visitor models.ShareVisitorInfo
	if err := database.DB.Where("id = ? AND share_id = ?", visitorID, shareID).First(&visitor).Error; err != nil {
		return errors.New(errors.CodeNotFound, "访客信息不存在或不属于该分享")
	}

	if err := database.DB.Delete(&visitor).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除访客信息失败")
	}

	return nil
}
