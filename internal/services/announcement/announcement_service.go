package announcement

import (
	"fmt"
	"pixelpunk/internal/controllers/announcement/dto"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"

	"gorm.io/gorm"
)

/* CreateAnnouncement 创建公告 */
func CreateAnnouncement(userID uint, createDTO *dto.AnnouncementCreateDTO) (*dto.AnnouncementResponseDTO, error) {
	db := database.GetDB()

	announcement := &models.Announcement{
		Title:     createDTO.Title,
		Content:   createDTO.Content,
		Summary:   createDTO.Summary,
		IsPinned:  createDTO.IsPinned,
		Status:    createDTO.Status,
		CreatedBy: userID,
	}

	if err := db.Create(announcement).Error; err != nil {
		return nil, fmt.Errorf("创建公告失败: %v", err)
	}

	return modelToResponseDTO(announcement), nil
}

/* UpdateAnnouncement 更新公告 */
func UpdateAnnouncement(id uint, updateDTO *dto.AnnouncementUpdateDTO) (*dto.AnnouncementResponseDTO, error) {
	db := database.GetDB()

	var announcement models.Announcement
	if err := db.First(&announcement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("公告不存在")
		}
		return nil, fmt.Errorf("查询公告失败: %v", err)
	}

	updates := make(map[string]interface{})

	if updateDTO.Title != nil {
		updates["title"] = *updateDTO.Title
	}
	if updateDTO.Content != nil {
		updates["content"] = *updateDTO.Content
	}
	if updateDTO.Summary != nil {
		updates["summary"] = *updateDTO.Summary
	}
	if updateDTO.IsPinned != nil {
		updates["is_pinned"] = *updateDTO.IsPinned
	}
	if updateDTO.Status != nil {
		updates["status"] = *updateDTO.Status
	}

	if err := db.Model(&announcement).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新公告失败: %v", err)
	}

	if err := db.First(&announcement, id).Error; err != nil {
		return nil, fmt.Errorf("查询更新后的公告失败: %v", err)
	}

	return modelToResponseDTO(&announcement), nil
}

/* DeleteAnnouncement 删除公告（软删除） */
func DeleteAnnouncement(id uint) error {
	db := database.GetDB()

	if err := db.Delete(&models.Announcement{}, id).Error; err != nil {
		return fmt.Errorf("删除公告失败: %v", err)
	}

	return nil
}

/* TogglePinAnnouncement 切换公告置顶状态（单一置顶逻辑） */
func TogglePinAnnouncement(id uint) (*dto.AnnouncementResponseDTO, error) {
	db := database.GetDB()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var announcement models.Announcement
	if err := tx.First(&announcement, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("公告不存在")
		}
		return nil, fmt.Errorf("查询公告失败: %v", err)
	}

	newPinnedStatus := !announcement.IsPinned

	// 如果要置顶此公告，先取消所有其他公告的置顶
	if newPinnedStatus {
		if err := tx.Model(&models.Announcement{}).
			Where("id != ?", id).
			Update("is_pinned", false).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("取消其他公告置顶失败: %v", err)
		}
	}

	if err := tx.Model(&announcement).Update("is_pinned", newPinnedStatus).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新置顶状态失败: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	if err := db.First(&announcement, id).Error; err != nil {
		return nil, fmt.Errorf("查询更新后的公告失败: %v", err)
	}

	return modelToResponseDTO(&announcement), nil
}

/* GetAnnouncementByID 根据ID获取公告 */
func GetAnnouncementByID(id uint) (*dto.AnnouncementResponseDTO, error) {
	db := database.GetDB()

	var announcement models.Announcement
	if err := db.First(&announcement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("公告不存在")
		}
		return nil, fmt.Errorf("查询公告失败: %v", err)
	}

	return modelToResponseDTO(&announcement), nil
}

/* GetAnnouncementList 获取公告列表（管理端） */
func GetAnnouncementList(query *dto.AnnouncementQueryDTO) (*dto.AnnouncementListResponseDTO, error) {
	db := database.GetDB()

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 20
	}

	queryBuilder := db.Model(&models.Announcement{})

	if query.Status != "" {
		queryBuilder = queryBuilder.Where("status = ?", query.Status)
	}

	if query.IsPinned != nil {
		queryBuilder = queryBuilder.Where("is_pinned = ?", *query.IsPinned)
	}

	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		queryBuilder = queryBuilder.Where("title LIKE ? OR summary LIKE ?", keyword, keyword)
	}

	var total int64
	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询公告总数失败: %v", err)
	}

	var announcements []models.Announcement
	offset := (query.Page - 1) * query.PageSize
	if err := queryBuilder.Order("is_pinned DESC, created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&announcements).Error; err != nil {
		return nil, fmt.Errorf("查询公告列表失败: %v", err)
	}

	announcementDTOs := make([]dto.AnnouncementResponseDTO, len(announcements))
	for i, announcement := range announcements {
		announcementDTOs[i] = *modelToResponseDTO(&announcement)
	}

	totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))

	return &dto.AnnouncementListResponseDTO{
		Data: announcementDTOs,
		Pagination: dto.PaginationResponseDTO{
			Page:      query.Page,
			PageSize:  query.PageSize,
			Total:     total,
			TotalPage: totalPages,
		},
	}, nil
}

/* GetPublicAnnouncementList 获取公开的公告列表（用户端） */
func GetPublicAnnouncementList() (*dto.PublicAnnouncementListDTO, error) {
	db := database.GetDB()

	config, err := GetAnnouncementSettings()
	if err != nil {
		return nil, fmt.Errorf("获取公告配置失败: %v", err)
	}

	// 从配置中读取显示数量，默认10条
	limit := 10
	if displayLimit, ok := config["announcement_display_limit"].(float64); ok {
		limit = int(displayLimit)
	} else if displayLimit, ok := config["announcement_display_limit"].(int); ok {
		limit = displayLimit
	}

	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var announcements []models.Announcement

	if err := db.Where("status = ?", "published").
		Order("is_pinned DESC, created_at DESC").
		Limit(limit).
		Find(&announcements).Error; err != nil {
		return nil, fmt.Errorf("查询公告列表失败: %v", err)
	}

	simpleDTOs := make([]dto.AnnouncementSimpleDTO, len(announcements))
	for i, announcement := range announcements {
		simpleDTOs[i] = *modelToSimpleDTO(&announcement)
	}

	return &dto.PublicAnnouncementListDTO{
		Announcements: simpleDTOs,
		Total:         len(simpleDTOs),
		Config:        config, // 返回配置信息
	}, nil
}

/* GetPublicAnnouncementDetail 获取公告详情（用户端） */
func GetPublicAnnouncementDetail(id uint) (*dto.AnnouncementDetailDTO, error) {
	db := database.GetDB()

	var announcement models.Announcement
	if err := db.First(&announcement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("公告不存在")
		}
		return nil, fmt.Errorf("查询公告失败: %v", err)
	}

	// 只返回已发布的公告
	if announcement.Status != "published" {
		return nil, fmt.Errorf("公告未发布")
	}

	if err := db.Model(&announcement).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		// 浏览次数更新失败不影响返回结果，只记录错误
	}

	return modelToDetailDTO(&announcement), nil
}

/* modelToResponseDTO 将模型转换为响应DTO */
func modelToResponseDTO(announcement *models.Announcement) *dto.AnnouncementResponseDTO {
	return &dto.AnnouncementResponseDTO{
		ID:        announcement.ID,
		Title:     announcement.Title,
		Content:   announcement.Content,
		Summary:   announcement.Summary,
		IsPinned:  announcement.IsPinned,
		Status:    announcement.Status,
		ViewCount: announcement.ViewCount,
		CreatedBy: announcement.CreatedBy,
		CreatedAt: announcement.CreatedAt,
		UpdatedAt: announcement.UpdatedAt,
	}
}

/* modelToSimpleDTO 将模型转换为简化DTO */
func modelToSimpleDTO(announcement *models.Announcement) *dto.AnnouncementSimpleDTO {
	return &dto.AnnouncementSimpleDTO{
		ID:        announcement.ID,
		Title:     announcement.Title,
		Summary:   announcement.Summary,
		IsPinned:  announcement.IsPinned,
		ViewCount: announcement.ViewCount,
		CreatedAt: announcement.CreatedAt,
	}
}

/* modelToDetailDTO 将模型转换为详情DTO */
func modelToDetailDTO(announcement *models.Announcement) *dto.AnnouncementDetailDTO {
	return &dto.AnnouncementDetailDTO{
		ID:        announcement.ID,
		Title:     announcement.Title,
		Content:   announcement.Content,
		Summary:   announcement.Summary,
		IsPinned:  announcement.IsPinned,
		ViewCount: announcement.ViewCount,
		CreatedAt: announcement.CreatedAt,
	}
}
