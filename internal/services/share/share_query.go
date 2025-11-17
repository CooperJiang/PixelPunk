package share

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"pixelpunk/internal/controllers/share/dto"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/email"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"time"

	"gorm.io/gorm"
)

func GetShareByKey(shareKey string) (models.Share, error) {
	var share models.Share
	if err := database.DB.Where("share_key = ?", shareKey).First(&share).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return models.Share{}, errors.New(errors.CodeNotFound, "分享不存在")
		}
		return models.Share{}, err
	}

	if share.Status != common.ShareStatusNormal {
		return models.Share{}, errors.New(errors.CodeValidationFailed, "分享已失效")
	}

	if share.IsExpired() {
		database.DB.Model(&share).Update("status", common.ShareStatusExpired)
		return models.Share{}, errors.New(errors.CodeValidationFailed, "分享已过期")
	}

	if share.MaxViews > 0 && share.CurrentViews >= share.MaxViews {
		return models.Share{}, errors.New(errors.CodeValidationFailed, "分享已达到最大访问次数")
	}

	return share, nil
}

/* VerifySharePassword 验证分享密码 */
func VerifySharePassword(shareKey string, password string) (bool, error) {
	share, err := GetShareByKey(shareKey)
	if err != nil {
		return false, err
	}

	return share.CanAccessWithPassword(password), nil
}

/* GetShareItems 获取分享的项目列表 */
func GetShareItems(shareID string) ([]models.ShareItem, error) {
	var items []models.ShareItem
	if err := database.DB.Where("share_id = ?", shareID).Order("sort_order asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

/* IncreaseShareViews 增加分享浏览次数 */
func IncreaseShareViews(shareID string) error {
	var share models.Share
	err := database.DB.Where("id = ?", shareID).First(&share).Error
	if err != nil {
		return err
	}

	previousViews := share.CurrentViews

	err = database.DB.Model(&models.Share{}).Where("id = ?", shareID).Update("current_views", gorm.Expr("current_views + ?", 1)).Error
	if err != nil {
		return err
	}

	newViews := previousViews + 1

	milestones := []int{50, 100, 200, 500, 1000}
	for _, milestone := range milestones {
		if previousViews < milestone && newViews >= milestone {
			go activity.LogShareMilestone(share.UserID, shareID, milestone)
			break // 只记录刚达到的第一个里程碑
		}
	}

	if share.NotificationOnAccess && share.CurrentViews == (share.NotificationThreshold-1) {
		go sendShareViewCountNotification(&share)
	}

	return nil
}

func sendShareViewCountNotification(share *models.Share) {
	var user models.User
	err := database.DB.Select("username, email").Where("id = ?", share.UserID).First(&user).Error
	if err != nil {
		return
	}

	if user.Email == "" {
		logger.Warn("用户邮箱为空，无法发送访问量通知")
		return
	}

	if !email.IsMailEnabled() {
		logger.Warn("邮件服务未启用，无法发送访问量通知")
		return
	}

	subject := fmt.Sprintf("您的分享达到%d次访问提醒", share.NotificationThreshold)
	body := fmt.Sprintf(`
		<h2>分享访问量提醒</h2>
		<p>%s，您好!</p>
		<p>您的分享内容已被访问%d次</p>
		<p>分享信息：</p>
		<ul>
			<li>名称：%s</li>
			<li>创建时间：%s</li>
			<li>当前访问量：%d</li>
			<li>分享ID：%s</li>
		</ul>
		<p>您可以登录系统查看更多详情。</p>
	`, user.Username, share.NotificationThreshold, share.Name, time.Time(share.CreatedAt).Format("2006-01-02 15:04:05"), share.NotificationThreshold, share.ID)

	err = email.SendMail(user.Email, subject, body)
	if err != nil {
		logger.Error("发送分享访问量提醒邮件失败: %v", err)
	} else {
	}
}

/* LogShareAccess 记录分享访问 */
func LogShareAccess(shareID string, viewedItems []map[string]string, c *dto.VisitorInfoDTO, ip, userAgent, referer string) error {
	viewedItemsJSON, err := json.Marshal(viewedItems)
	if err != nil {
		return err
	}

	log := models.ShareAccessLog{
		ID:          generateID(),
		ShareID:     shareID,
		AccessedAt:  common.JSONTime(time.Now()),
		IPAddress:   ip,
		UserAgent:   userAgent,
		Referer:     referer,
		ViewedItems: viewedItemsJSON,
	}

	if c != nil {
		log.VisitorName = c.Name
		log.VisitorEmail = c.Email
	}

	return database.DB.Create(&log).Error
}

/* GetUserShares 获取用户创建的分享列表 */
func GetUserShares(userID uint, query *dto.ShareListQueryDTO) ([]map[string]interface{}, int64, error) {
	db := database.DB.Model(&models.Share{}).Where("user_id = ?", userID)

	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("(name LIKE ? OR description LIKE ?)", keyword, keyword)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := query.Page
	if page <= 0 {
		page = 1
	}

	size := query.Size
	if size <= 0 {
		size = common.DefaultPageSize
	}
	if size > common.MaxPageSize {
		size = common.MaxPageSize
	}

	var shares []models.Share
	if err := db.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&shares).Error; err != nil {
		return nil, 0, err
	}

	if len(shares) == 0 {
		return []map[string]interface{}{}, total, nil
	}

	var shareIDs []string
	for _, share := range shares {
		shareIDs = append(shareIDs, share.ID)
	}

	type ShareItemCount struct {
		ShareID  string `gorm:"column:share_id"`
		ItemType string `gorm:"column:item_type"`
		Count    int64  `gorm:"column:count"`
	}

	var counts []ShareItemCount
	err := database.DB.Table("share_item").
		Select("share_id, item_type, COUNT(*) as count").
		Where("share_id IN ?", shareIDs).
		Group("share_id, item_type").
		Find(&counts).Error

	if err != nil {
	}

	countMap := make(map[string]map[string]int64)
	for _, c := range counts {
		if countMap[c.ShareID] == nil {
			countMap[c.ShareID] = make(map[string]int64)
		}
		countMap[c.ShareID][c.ItemType] = c.Count
	}

	result := make([]map[string]interface{}, len(shares))
	for i, share := range shares {
		folderCount := int64(0)
		fileCount := int64(0)

		if countMap[share.ID] != nil {
			folderCount = countMap[share.ID]["folder"]
			fileCount = countMap[share.ID]["file"]
		}

		shareMap := map[string]interface{}{
			"id":                     share.ID,
			"share_key":              share.ShareKey,
			"name":                   share.Name,
			"description":            share.Description,
			"expired_days":           share.ExpiredDays,
			"expired_at":             share.ExpiredAt,
			"max_views":              share.MaxViews,
			"current_views":          share.CurrentViews,
			"status":                 share.Status,
			"created_at":             share.CreatedAt,
			"updated_at":             share.UpdatedAt,
			"folder_count":           folderCount,
			"file_count":             fileCount,
			"collect_visitor_info":   share.CollectVisitorInfo,
			"notification_on_access": share.NotificationOnAccess,
		}

		result[i] = shareMap
	}

	return result, total, nil
}

/* DeleteShare 删除分享 */
func DeleteShare(shareID string, userID uint, force bool) error {
	var share models.Share
	if err := database.DB.Where("id = ? AND user_id = ?", shareID, userID).First(&share).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.CodeForbidden, "分享不存在或无权操作")
		}
		return err
	}

	if force {
		// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
		if err := database.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("share_id = ?", shareID).Delete(&models.ShareItem{}).Error; err != nil {
				return fmt.Errorf("删除分享项目失败: %w", err)
			}

			if err := tx.Where("share_id = ?", shareID).Delete(&models.ShareAccessLog{}).Error; err != nil {
				return fmt.Errorf("删除分享访问日志失败: %w", err)
			}

			if err := tx.Where("share_id = ?", shareID).Delete(&models.ShareVisitorInfo{}).Error; err != nil {
				return fmt.Errorf("删除分享访客信息失败: %w", err)
			}

			if err := tx.Delete(&share).Error; err != nil {
				return fmt.Errorf("删除分享失败: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}
	} else {
		if err := database.DB.Model(&models.Share{}).Where("id = ?", shareID).Update("status", common.ShareStatusDeleted).Error; err != nil {
			return fmt.Errorf("删除分享失败: %w", err)
		}
	}

	return nil
}

/* GetShareForView 获取分享详情页面所需数据 */
