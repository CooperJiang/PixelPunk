package share

import (
	"pixelpunk/internal/controllers/share/dto"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

/* GetAllShares 管理员获取所有分享列表 */
func GetAllShares(query *dto.AdminShareListQueryDTO) ([]models.Share, int64, error) {
	db := database.DB.Model(&models.Share{})

	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("(name LIKE ? OR description LIKE ? OR share_key LIKE ?)", keyword, keyword, keyword)
	}

	if query.UserID > 0 {
		db = db.Where("user_id = ?", query.UserID)
	}

	if !query.StartDate.IsZero() {
		db = db.Where("created_at >= ?", query.StartDate)
	}

	if !query.EndDate.IsZero() {
		db = db.Where("created_at <= ?", query.EndDate)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分享总数失败")
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

	orderBy := "created_at DESC"
	if query != nil && strings.TrimSpace(query.OrderBy) != "" {
		orderBy = query.OrderBy
	}

	var shares []models.Share
	if err := db.Order(orderBy).Offset((page - 1) * size).Limit(size).Find(&shares).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分享列表失败")
	}

	return shares, total, nil
}

/* AdminGetShareDetail 管理员获取分享详情 */
func AdminGetShareDetail(shareID string) (*models.Share, []models.ShareItem, error) {
	var share models.Share
	if err := database.DB.Where("id = ?", shareID).First(&share).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, errors.New(errors.CodeNotFound, "分享不存在")
		}
		return nil, nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分享详情失败")
	}

	var items []models.ShareItem
	if err := database.DB.Where("share_id = ?", shareID).Order("sort_order asc").Find(&items).Error; err != nil {
		return nil, nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分享项目失败")
	}

	return &share, items, nil
}

/* AdminUpdateShareStatus 管理员更新分享状态 */
func AdminUpdateShareStatus(shareID string, statusDTO *dto.AdminShareStatusUpdateDTO, adminID uint) error {
	var share models.Share
	if err := database.DB.Where("id = ?", shareID).First(&share).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "分享不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "获取分享详情失败")
	}

	if err := database.DB.Model(&models.Share{}).Where("id = ?", shareID).Update("status", statusDTO.Status).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新分享状态失败")
	}

	return nil
}

/* AdminDeleteShare 管理员删除分享 */
func AdminDeleteShare(shareID string, adminID uint, force bool) error {
	var share models.Share
	if err := database.DB.Where("id = ?", shareID).First(&share).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "分享不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "获取分享详情失败")
	}

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if force {
			if err := tx.Where("share_id = ?", shareID).Delete(&models.ShareItem{}).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除分享项目失败")
			}

			if err := tx.Where("share_id = ?", shareID).Delete(&models.ShareAccessLog{}).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除分享访问日志失败")
			}

			if err := tx.Where("share_id = ?", shareID).Delete(&models.ShareVisitorInfo{}).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除分享访客信息失败")
			}

			if err := tx.Delete(&share).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除分享失败")
			}
		} else {
			if err := tx.Model(&models.Share{}).Where("id = ?", shareID).Update("status", common.ShareStatusDeleted).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBUpdateFailed, "删除分享失败")
			}
		}

		return nil
	})
}

/* GetShareSystemStats 获取分享系统统计数据 */
func GetShareSystemStats() (*dto.AdminShareStatsResponseDTO, error) {
	db := database.DB
	stats := &dto.AdminShareStatsResponseDTO{}

	if err := db.Model(&models.Share{}).Count(&stats.TotalShares).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取总分享数失败")
	}

	if err := db.Model(&models.Share{}).Where("status = ?", common.ShareStatusNormal).Count(&stats.ActiveShares).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取活跃分享数失败")
	}

	todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
	if err := db.Model(&models.ShareAccessLog{}).Where("accessed_at >= ?", todayStart).Count(&stats.ViewsToday).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取今日浏览量失败")
	}

	var popularShares []struct {
		ID          string `json:"id"`
		ShareKey    string `json:"share_key"`
		Name        string `json:"name"`
		UserID      uint   `json:"user_id"`
		TotalViews  int    `json:"total_views"`
		RecentViews int    `json:"recent_views"`
	}

	// 兼容SQLite和MySQL: 使用 GORM ORM 方法替代 Raw SQL
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// 使用子查询和 GORM 方法，确保 SQLite 兼容性
	if err := db.Model(&models.Share{}).
		Select("share.id, share.share_key, share.name, share.user_id, share.current_views as total_views").
		Where("share.status = ?", common.ShareStatusNormal).
		Order("share.current_views DESC").
		Limit(10).
		Find(&popularShares).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取热门分享失败")
	}

	for i := range popularShares {
		var recentViews int64
		db.Model(&models.ShareAccessLog{}).
			Where("share_id = ? AND accessed_at >= ?", popularShares[i].ID, sevenDaysAgo).
			Count(&recentViews)
		popularShares[i].RecentViews = int(recentViews)
	}

	// 按最近访问量重新排序
	sort.Slice(popularShares, func(i, j int) bool {
		if popularShares[i].RecentViews != popularShares[j].RecentViews {
			return popularShares[i].RecentViews > popularShares[j].RecentViews
		}
		return popularShares[i].TotalViews > popularShares[j].TotalViews
	})

	if len(popularShares) > 10 {
		popularShares = popularShares[:10]
	}

	for _, ps := range popularShares {
		var username string
		db.Model(&models.User{}).Where("id = ?", ps.UserID).Select("username").Scan(&username)

		stats.PopularShares = append(stats.PopularShares, dto.PopularShareDTO{
			ID:          ps.ID,
			ShareKey:    ps.ShareKey,
			Name:        ps.Name,
			UserID:      ps.UserID,
			Username:    username,
			TotalViews:  ps.TotalViews,
			RecentViews: ps.RecentViews,
		})
	}

	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		dateStart := dateStr + " 00:00:00"
		dateEnd := dateStr + " 23:59:59"

		var newShares int64
		var views int64

		db.Model(&models.Share{}).Where("created_at BETWEEN ? AND ?", dateStart, dateEnd).Count(&newShares)

		db.Model(&models.ShareAccessLog{}).Where("accessed_at BETWEEN ? AND ?", dateStart, dateEnd).Count(&views)

		stats.StatsByDate = append(stats.StatsByDate, dto.ShareStatByDate{
			Date:       dateStr,
			NewShares:  int(newShares),
			TotalViews: int(views),
		})
	}

	return stats, nil
}

/* AdminGetAllVisitors 管理员获取所有访客信息 */
func AdminGetAllVisitors(query *dto.VisitorQueryDTO) ([]models.ShareVisitorInfo, int64, error) {
	db := database.DB.Model(&models.ShareVisitorInfo{})

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

/* AdminDeleteVisitorInfo 管理员删除访客信息 */
func AdminDeleteVisitorInfo(visitorID string) error {
	var visitor models.ShareVisitorInfo
	if err := database.DB.Where("id = ?", visitorID).First(&visitor).Error; err != nil {
		return errors.New(errors.CodeNotFound, "访客信息不存在")
	}

	if err := database.DB.Delete(&visitor).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除访客信息失败")
	}

	return nil
}
