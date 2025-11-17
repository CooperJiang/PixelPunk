package stats

import (
	"pixelpunk/internal/models"
	"time"

	"gorm.io/gorm"
)

type GlobalStatsService struct {
	DB *gorm.DB
}

func NewGlobalStatsService(db *gorm.DB) *GlobalStatsService {
	return &GlobalStatsService{
		DB: db,
	}
}

func (s *GlobalStatsService) InitTodayStats() error {
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	var count int64
	if err := s.DB.Model(&models.GlobalStats{}).Where("date = ?", todayStart).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		yesterdayStart := todayStart.AddDate(0, 0, -1)
		var yesterdayStats models.GlobalStats

		err := s.DB.Where("date = ?", yesterdayStart).First(&yesterdayStats).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		todayStats := models.GlobalStats{
			Date:           todayStart,
			TotalImages:    yesterdayStats.TotalImages,
			TotalStorage:   yesterdayStats.TotalStorage,
			TotalBandwidth: yesterdayStats.TotalBandwidth,
			TotalFolders:   yesterdayStats.TotalFolders,
			TotalUsers:     yesterdayStats.TotalUsers,
			TotalViews:     yesterdayStats.TotalViews,
		}

		return s.DB.Create(&todayStats).Error
	}

	return nil
}

func (s *GlobalStatsService) IncrementFileStats(size int64) error {
	if err := s.InitTodayStats(); err != nil {
		return err
	}

	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	return s.DB.Model(&models.GlobalStats{}).
		Where("date = ?", todayStart).
		Updates(map[string]interface{}{
			"total_images":  gorm.Expr("total_images + ?", 1),
			"total_storage": gorm.Expr("total_storage + ?", size),
			"new_images":    gorm.Expr("new_images + ?", 1),
			"new_storage":   gorm.Expr("new_storage + ?", size),
		}).Error
}

func (s *GlobalStatsService) IncrementViewStats(bandwidth int64) error {
	if err := s.InitTodayStats(); err != nil {
		return err
	}

	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	return s.DB.Model(&models.GlobalStats{}).
		Where("date = ?", todayStart).
		Updates(map[string]interface{}{
			"total_views":     gorm.Expr("total_views + ?", 1),
			"total_bandwidth": gorm.Expr("total_bandwidth + ?", bandwidth),
			"new_views":       gorm.Expr("new_views + ?", 1),
			"new_bandwidth":   gorm.Expr("new_bandwidth + ?", bandwidth),
		}).Error
}

func (s *GlobalStatsService) IncrementUserStats() error {
	if err := s.InitTodayStats(); err != nil {
		return err
	}

	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	return s.DB.Model(&models.GlobalStats{}).
		Where("date = ?", todayStart).
		Updates(map[string]interface{}{
			"total_users": gorm.Expr("total_users + ?", 1),
			"new_users":   gorm.Expr("new_users + ?", 1),
		}).Error
}

func (s *GlobalStatsService) IncrementFolderStats() error {
	if err := s.InitTodayStats(); err != nil {
		return err
	}

	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	return s.DB.Model(&models.GlobalStats{}).
		Where("date = ?", todayStart).
		Updates(map[string]interface{}{
			"total_folders": gorm.Expr("total_folders + ?", 1),
			"new_folders":   gorm.Expr("new_folders + ?", 1),
		}).Error
}

func (s *GlobalStatsService) ReconcileAllStats() error {
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	var totalImages, totalFolders, totalUsers int64
	var totalStorage, totalBandwidth, totalViews int64

	if err := s.DB.Model(&models.File{}).Count(&totalImages).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.File{}).Select("COALESCE(SUM(size), 0)").Scan(&totalStorage).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.FileStats{}).Select("COALESCE(SUM(bandwidth), 0)").Scan(&totalBandwidth).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.FileStats{}).Select("COALESCE(SUM(views), 0)").Scan(&totalViews).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.Folder{}).Count(&totalFolders).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return err
	}

	var newImages, newFolders, newUsers int64
	var newStorage, newBandwidth, newViews int64

	if err := s.DB.Model(&models.File{}).Where("created_at >= ?", todayStart).Count(&newImages).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.File{}).Where("created_at >= ?", todayStart).
		Select("COALESCE(SUM(size), 0)").Scan(&newStorage).Error; err != nil {
		return err
	}

	var yesterdayBandwidth, yesterdayViews int64

	if err := s.DB.Model(&models.FileStats{}).
		Where("updated_at < ?", todayStart).
		Select("COALESCE(SUM(bandwidth), 0)").Scan(&yesterdayBandwidth).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.FileStats{}).
		Where("updated_at < ?", todayStart).
		Select("COALESCE(SUM(views), 0)").Scan(&yesterdayViews).Error; err != nil {
		return err
	}

	var currentBandwidth, currentViews int64

	if err := s.DB.Model(&models.FileStats{}).
		Select("COALESCE(SUM(bandwidth), 0)").Scan(&currentBandwidth).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.FileStats{}).
		Select("COALESCE(SUM(views), 0)").Scan(&currentViews).Error; err != nil {
		return err
	}

	newBandwidth = currentBandwidth - yesterdayBandwidth
	newViews = currentViews - yesterdayViews

	if newBandwidth < 0 {
		newBandwidth = 0
	}
	if newViews < 0 {
		newViews = 0
	}

	if err := s.DB.Model(&models.Folder{}).Where("created_at >= ?", todayStart).Count(&newFolders).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.User{}).Where("created_at >= ?", todayStart).Count(&newUsers).Error; err != nil {
		return err
	}

	todayStats := models.GlobalStats{
		Date:           todayStart,
		TotalImages:    totalImages,
		TotalStorage:   totalStorage,
		TotalBandwidth: totalBandwidth,
		TotalFolders:   totalFolders,
		TotalUsers:     totalUsers,
		TotalViews:     totalViews,
		NewImages:      newImages,
		NewStorage:     newStorage,
		NewBandwidth:   newBandwidth,
		NewFolders:     newFolders,
		NewUsers:       newUsers,
		NewViews:       newViews,
	}

	result := s.DB.Model(&models.GlobalStats{}).Where("date = ?", todayStart).Updates(map[string]interface{}{
		"total_images":    totalImages,
		"total_storage":   totalStorage,
		"total_bandwidth": totalBandwidth,
		"total_folders":   totalFolders,
		"total_users":     totalUsers,
		"total_views":     totalViews,
		"new_images":      newImages,
		"new_storage":     newStorage,
		"new_bandwidth":   newBandwidth,
		"new_folders":     newFolders,
		"new_users":       newUsers,
		"new_views":       newViews,
	})

	if result.RowsAffected == 0 {
		if err := s.DB.Create(&todayStats).Error; err != nil {
			return err
		}
	} else if result.Error != nil {
		return result.Error
	}

	return nil
}

/* ReconcileTodayStats 仅校准今日统计数据（更轻量级） */
func (s *GlobalStatsService) ReconcileTodayStats() error {
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	var newImages, newFolders, newUsers int64
	var newStorage, newBandwidth, newViews int64

	if err := s.DB.Model(&models.File{}).Where("created_at >= ?", todayStart).Count(&newImages).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.File{}).Where("created_at >= ?", todayStart).
		Select("COALESCE(SUM(size), 0)").Scan(&newStorage).Error; err != nil {
		return err
	}

	var yesterdayBandwidth, yesterdayViews int64

	if err := s.DB.Model(&models.FileStats{}).
		Where("updated_at < ?", todayStart).
		Select("COALESCE(SUM(bandwidth), 0)").Scan(&yesterdayBandwidth).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.FileStats{}).
		Where("updated_at < ?", todayStart).
		Select("COALESCE(SUM(views), 0)").Scan(&yesterdayViews).Error; err != nil {
		return err
	}

	var currentBandwidth, currentViews int64

	if err := s.DB.Model(&models.FileStats{}).
		Select("COALESCE(SUM(bandwidth), 0)").Scan(&currentBandwidth).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.FileStats{}).
		Select("COALESCE(SUM(views), 0)").Scan(&currentViews).Error; err != nil {
		return err
	}

	newBandwidth = currentBandwidth - yesterdayBandwidth
	newViews = currentViews - yesterdayViews

	if newBandwidth < 0 {
		newBandwidth = 0
	}
	if newViews < 0 {
		newViews = 0
	}

	if err := s.DB.Model(&models.Folder{}).Where("created_at >= ?", todayStart).Count(&newFolders).Error; err != nil {
		return err
	}

	if err := s.DB.Model(&models.User{}).Where("created_at >= ?", todayStart).Count(&newUsers).Error; err != nil {
		return err
	}

	var yesterdayStats models.GlobalStats
	if err := s.DB.Where("date < ?", todayStart).Order("date DESC").Limit(1).First(&yesterdayStats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return s.ReconcileAllStats()
		}
		return err
	}

	totalImages := yesterdayStats.TotalImages + newImages
	totalStorage := yesterdayStats.TotalStorage + newStorage
	totalBandwidth := yesterdayStats.TotalBandwidth + newBandwidth
	totalFolders := yesterdayStats.TotalFolders + newFolders
	totalUsers := yesterdayStats.TotalUsers + newUsers
	totalViews := yesterdayStats.TotalViews + newViews

	todayStats := models.GlobalStats{
		Date:           todayStart,
		TotalImages:    totalImages,
		TotalStorage:   totalStorage,
		TotalBandwidth: totalBandwidth,
		TotalFolders:   totalFolders,
		TotalUsers:     totalUsers,
		TotalViews:     totalViews,
		NewImages:      newImages,
		NewStorage:     newStorage,
		NewBandwidth:   newBandwidth,
		NewFolders:     newFolders,
		NewUsers:       newUsers,
		NewViews:       newViews,
	}

	result := s.DB.Model(&models.GlobalStats{}).Where("date = ?", todayStart).Updates(map[string]interface{}{
		"total_images":    totalImages,
		"total_storage":   totalStorage,
		"total_bandwidth": totalBandwidth,
		"total_folders":   totalFolders,
		"total_users":     totalUsers,
		"total_views":     totalViews,
		"new_images":      newImages,
		"new_storage":     newStorage,
		"new_bandwidth":   newBandwidth,
		"new_folders":     newFolders,
		"new_users":       newUsers,
		"new_views":       newViews,
	})

	if result.RowsAffected == 0 {
		if err := s.DB.Create(&todayStats).Error; err != nil {
			return err
		}
	} else if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *GlobalStatsService) GetLatestStats() (*models.GlobalStats, error) {
	var stats models.GlobalStats

	err := s.DB.Order("date DESC").First(&stats).Error
	return &stats, err
}

func (s *GlobalStatsService) GetStatsRange(startDate, endDate time.Time) ([]models.GlobalStats, error) {
	var stats []models.GlobalStats

	err := s.DB.Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date").Find(&stats).Error
	return stats, err
}
