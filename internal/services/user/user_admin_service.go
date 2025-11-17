package user

import (
	"pixelpunk/internal/controllers/user/dto"
	"pixelpunk/internal/models"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUserList(query *dto.AdminUserListQueryDTO) (*dto.AdminUserListResponseDTO, error) {
	db := database.GetDB()

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = common.DefaultPageSize
	}
	if query.Size > common.MaxPageSize {
		query.Size = common.MaxPageSize
	}

	offset := (query.Page - 1) * query.Size

	dbQuery := db.Model(&models.User{})

	if query.Role > 0 {
		dbQuery = dbQuery.Where("role = ?", query.Role)
	}

	if query.Status > 0 {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	if query.Keyword != "" {
		dbQuery = dbQuery.Where("username LIKE ? OR email LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询用户总数失败")
	}

	var users []models.User
	if err := dbQuery.Limit(query.Size).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询用户列表失败")
	}

	result := &dto.AdminUserListResponseDTO{
		Total: total,
		List:  make([]dto.AdminUserResponseDTO, 0, len(users)),
	}

	for _, user := range users {
		avatarFullPath := ""
		if user.Avatar != "" {
			avatarFullPath = utils.GetSystemFileURL(user.Avatar)
		}

		var userSettings models.UserSettings
		if err := db.Where("user_id = ?", user.ID).First(&userSettings).Error; err != nil {
			userSettings = models.UserSettings{
				UserID:         user.ID,
				StorageLimit:   models.DefaultStorageLimit,
				BandwidthLimit: models.DefaultBandwidthLimit,
			}
		}

		var userStats models.UserUsageStats
		if err := db.Where("user_id = ?", user.ID).First(&userStats).Error; err != nil {
			userStats = models.UserUsageStats{
				UserID:         user.ID,
				TotalImages:    0,
				TotalSize:      0,
				TotalBandwidth: 0,
				TotalViews:     0,
			}
		}

		result.List = append(result.List, dto.AdminUserResponseDTO{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			Avatar:         user.Avatar,
			AvatarFullPath: avatarFullPath,
			Status:         user.Status,
			Role:           user.Role,
			StorageLimit:   userSettings.StorageLimit,
			BandwidthLimit: userSettings.BandwidthLimit,
			UsedStorage:    userStats.TotalSize,
			UsedBandwidth:  userStats.TotalBandwidth,
			TotalFiles:     userStats.TotalImages,
			TotalViews:     userStats.TotalViews,
			LastActivityAt: user.LastActivityAt,
			LastActivityIP: user.LastActivityIP,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		})
	}

	return result, nil
}

func AdminUpdateUser(updateDTO *dto.AdminUpdateUserDTO, currentUserID uint) error {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, updateDTO.ID).Error; err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	const SystemRootAdminID uint = 1 // 系统默认超级管理员ID

	if user.IsSuperAdmin() {
		if currentUserID != SystemRootAdminID {
			return errors.New(errors.CodeForbidden, "只有系统默认超级管理员可以修改超级管理员信息")
		}
		if user.ID == SystemRootAdminID {
			return errors.New(errors.CodeForbidden, "系统默认超级管理员信息不能被修改")
		}
	}

	if !user.IsSuperAdmin() && updateDTO.Role == common.UserRoleSuperAdmin {
		if currentUserID != SystemRootAdminID {
			return errors.New(errors.CodeForbidden, "只有系统默认超级管理员可以将用户提升为超级管理员")
		}
	}

	if updateDTO.Username != user.Username {
		var count int64
		if err := db.Model(&models.User{}).
			Where("username = ? AND id != ?", updateDTO.Username, updateDTO.ID).
			Count(&count).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBQueryFailed, "检查用户名是否存在失败")
		}

		if count > 0 {
			return errors.New(errors.CodeUserExists, "用户名已被使用")
		}
	}

	oldStatus := user.Status // 保存旧状态

	updates := map[string]interface{}{
		"username": updateDTO.Username,
		"status":   updateDTO.Status,
		"role":     updateDTO.Role,
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户信息失败")
	}

	// 如果状态发生变化，同步到Redis黑名单
	if oldStatus != updateDTO.Status {
		syncUserStatusToRedis(updateDTO.ID, updateDTO.Status)
	}

	return nil
}

func AdminGetUserByID(id uint) (*dto.AdminUserResponseDTO, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	avatarFullPath := ""
	if user.Avatar != "" {
		avatarFullPath = utils.GetSystemFileURL(user.Avatar)
	}

	var userSettings models.UserSettings
	if err := db.Where("user_id = ?", user.ID).First(&userSettings).Error; err != nil {
		userSettings = models.UserSettings{
			UserID:         user.ID,
			StorageLimit:   models.DefaultStorageLimit,
			BandwidthLimit: models.DefaultBandwidthLimit,
		}
	}

	var userStats models.UserUsageStats
	if err := db.Where("user_id = ?", user.ID).First(&userStats).Error; err != nil {
		userStats = models.UserUsageStats{
			UserID:         user.ID,
			TotalImages:    0,
			TotalSize:      0,
			TotalBandwidth: 0,
			TotalViews:     0,
		}
	}

	return &dto.AdminUserResponseDTO{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Avatar:         user.Avatar,
		AvatarFullPath: avatarFullPath,
		Status:         user.Status,
		Role:           user.Role,
		StorageLimit:   userSettings.StorageLimit,
		BandwidthLimit: userSettings.BandwidthLimit,
		UsedStorage:    userStats.TotalSize,
		UsedBandwidth:  userStats.TotalBandwidth,
		TotalFiles:     userStats.TotalImages,
		TotalViews:     userStats.TotalViews,
		LastActivityAt: user.LastActivityAt,
		LastActivityIP: user.LastActivityIP,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

func AdminDeleteUser(id uint) error {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	if user.IsSuperAdmin() {
		return errors.New(errors.CodeForbidden, "不能删除超级管理员")
	}

	if err := db.Model(&user).Update("status", common.UserStatusDeleted).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "删除用户失败")
	}

	return nil
}

func AdminCreateUser(createDTO *dto.AdminCreateUserDTO) (*dto.AdminUserResponseDTO, error) {
	db := database.GetDB()

	var existingUser models.User
	if err := db.Where("username = ?", createDTO.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New(errors.CodeUserExists, "用户名已存在")
	}

	if err := db.Where("email = ?", createDTO.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New(errors.CodeUserExists, "邮箱已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "密码加密失败")
	}

	var user *models.User
	var userSettings models.UserSettings
	var userStats models.UserUsageStats
	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	err = db.Transaction(func(tx *gorm.DB) error {
		user = &models.User{
			Username:  createDTO.Username,
			Email:     createDTO.Email,
			Password:  string(hashedPassword),
			Role:      createDTO.Role,
			Status:    common.UserStatusNormal,        // 默认状态为正常
			PathAlias: utils.GenerateRandomString(16), // 生成临时路径别名，避免唯一索引冲突
		}

		if err := tx.Create(user).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户失败")
		}

		storageLimit := createDTO.StorageLimit
		if storageLimit <= 0 {
			storageLimit = models.DefaultStorageLimit
		}
		bandwidthLimit := createDTO.BandwidthLimit
		if bandwidthLimit <= 0 {
			bandwidthLimit = models.DefaultBandwidthLimit
		}

		userSettings = models.UserSettings{
			UserID:         user.ID,
			StorageLimit:   storageLimit,
			BandwidthLimit: bandwidthLimit,
		}

		if err := tx.Create(&userSettings).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户设置失败")
		}

		userStats = models.UserUsageStats{
			UserID:         user.ID,
			TotalImages:    0,
			TotalSize:      0,
			TotalBandwidth: 0,
			TotalViews:     0,
		}

		if err := tx.Create(&userStats).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户统计失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if _, aliasErr := tenant.ResolveAlias(user.ID); aliasErr != nil {
		logger.Warn("生成用户路径别名失败(管理员创建): userID=%d, err=%v", user.ID, aliasErr)
	}

	avatarFullPath := ""
	if user.Avatar != "" {
		avatarFullPath = utils.GetSystemFileURL(user.Avatar)
	}

	return &dto.AdminUserResponseDTO{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Avatar:         user.Avatar,
		AvatarFullPath: avatarFullPath,
		Status:         user.Status,
		Role:           user.Role,
		StorageLimit:   userSettings.StorageLimit,
		BandwidthLimit: userSettings.BandwidthLimit,
		UsedStorage:    userStats.TotalSize,
		UsedBandwidth:  userStats.TotalBandwidth,
		TotalFiles:     userStats.TotalImages,
		TotalViews:     userStats.TotalViews,
		LastActivityAt: user.LastActivityAt,
		LastActivityIP: user.LastActivityIP,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

func AdminUpdateUserStorage(updateDTO *dto.AdminUpdateUserStorageDTO) error {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, updateDTO.UserID).Error; err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	if user.IsSuperAdmin() && updateDTO.CurrentUserID != updateDTO.UserID {
		return errors.New(errors.CodeForbidden, "只有超级管理员本人可以修改自己的设置")
	}

	var userSettings models.UserSettings
	var oldStorageLimit int64 = 0
	err := db.Where("user_id = ?", updateDTO.UserID).First(&userSettings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userSettings = models.UserSettings{
				UserID:         updateDTO.UserID,
				StorageLimit:   updateDTO.StorageLimit,
				BandwidthLimit: updateDTO.BandwidthLimit,
			}
			if err := db.Create(&userSettings).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户设置失败")
			}
		} else {
			return errors.Wrap(err, errors.CodeDBQueryFailed, "查询用户设置失败")
		}
	} else {
		oldStorageLimit = userSettings.StorageLimit

		userSettings.StorageLimit = updateDTO.StorageLimit
		userSettings.BandwidthLimit = updateDTO.BandwidthLimit
		if err := db.Save(&userSettings).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户设置失败")
		}
	}

	if oldStorageLimit != updateDTO.StorageLimit {
		go sendStorageChangeNotification(updateDTO.UserID, oldStorageLimit, updateDTO.StorageLimit)
	}

	return nil
}

func AdminResetUserPassword(id uint) (*dto.AdminResetUserPasswordResponseDTO, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	if user.IsSuperAdmin() {
		return nil, errors.New(errors.CodeForbidden, "不能重置超级管理员密码")
	}

	newPassword := utils.GenerateRandomString(8)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "密码加密失败")
	}

	if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "重置密码失败")
	}

	return &dto.AdminResetUserPasswordResponseDTO{
		NewPassword: newPassword,
	}, nil
}

func AdminSendUserEmail(emailDTO *dto.AdminSendUserEmailDTO) error {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, emailDTO.UserID).Error; err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	return nil
}

func AdminToggleUserStatus(statusDTO *dto.AdminToggleUserStatusDTO) error {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, statusDTO.UserID).Error; err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	if user.IsSuperAdmin() {
		return errors.New(errors.CodeForbidden, "不能修改超级管理员状态")
	}

	if err := db.Model(&user).Update("status", statusDTO.Status).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户状态失败")
	}

	// 同步到Redis黑名单（实现即时踢出在线用户）
	syncUserStatusToRedis(statusDTO.UserID, statusDTO.Status)

	return nil
}

func AdminBatchOperateUsers(batchDTO *dto.AdminBatchOperateUsersDTO) error {
	db := database.GetDB()

	var users []models.User
	if err := db.Where("id IN (?)", batchDTO.UserIDs).Find(&users).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询用户失败")
	}

	if len(users) == 0 {
		return errors.New(errors.CodeUserNotFound, "未找到指定的用户")
	}

	for _, user := range users {
		if user.IsSuperAdmin() {
			return errors.New(errors.CodeForbidden, "不能对超级管理员进行批量操作")
		}
	}

	var affectedStatus int
	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	err := db.Transaction(func(tx *gorm.DB) error {
		switch batchDTO.Operation {
		case "enable":
			if err := tx.Model(&models.User{}).Where("id IN (?)", batchDTO.UserIDs).Update("status", common.UserStatusNormal).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBUpdateFailed, "批量启用用户失败")
			}
			affectedStatus = common.UserStatusNormal
		case "disable":
			if err := tx.Model(&models.User{}).Where("id IN (?)", batchDTO.UserIDs).Update("status", common.UserStatusDisabled).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBUpdateFailed, "批量禁用用户失败")
			}
			affectedStatus = common.UserStatusDisabled
		case "delete":
			if err := tx.Model(&models.User{}).Where("id IN (?)", batchDTO.UserIDs).Update("status", common.UserStatusDeleted).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBUpdateFailed, "批量删除用户失败")
			}
			affectedStatus = common.UserStatusDeleted
		case "set_role":
			if batchDTO.Role == 0 {
				return errors.New(errors.CodeInvalidParameter, "设置角色时必须指定角色值")
			}
			if err := tx.Model(&models.User{}).Where("id IN (?)", batchDTO.UserIDs).Update("role", batchDTO.Role).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBUpdateFailed, "批量设置角色失败")
			}
			affectedStatus = 0 // 角色变更不涉及状态同步
		default:
			return errors.New(errors.CodeInvalidParameter, "不支持的操作类型")
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 批量同步到Redis黑名单（仅状态变更操作）
	if affectedStatus != 0 {
		for _, userID := range batchDTO.UserIDs {
			syncUserStatusToRedis(userID, affectedStatus)
		}
	}

	return nil
}

func sendStorageChangeNotification(userID uint, oldStorageLimit, newStorageLimit int64) {
	oldSize := formatBytes(oldStorageLimit)
	newSize := formatBytes(newStorageLimit)

	var messageType string
	variables := map[string]interface{}{
		"old_size":     oldSize,
		"new_size":     newSize,
		"related_type": "storage",
		"related_id":   "storage_quota_change",
	}

	if newStorageLimit > oldStorageLimit {
		messageType = common.MessageTypeStorageQuotaIncreased
		variables["increased_size"] = formatBytes(newStorageLimit - oldStorageLimit)
	} else {
		messageType = common.MessageTypeStorageQuotaDecreased
		variables["decreased_size"] = formatBytes(oldStorageLimit - newStorageLimit)
		variables["reason"] = "管理员调整"
	}

	msgService := messageService.GetMessageService()
	if err := msgService.SendTemplateMessage(userID, messageType, variables); err != nil {
		logger.Warn("发送存储空间变更消息失败: userID=%d, oldSize=%s, newSize=%s, error=%v", userID, oldSize, newSize, err)
	} else {
	}
}

func syncUserStatusToRedis(userID uint, status int) {
	redisClient := cache.GetRedisClient()
	if redisClient == nil {
		logger.Warn("syncUserStatusToRedis: Redis不可用，跳过黑名单同步: userID=%d", userID)
		return
	}

	ctx := cache.GetRedisContext()
	ns := cache.GetNamespace()
	key := ns + ":disabled_users"

	switch status {
	case common.UserStatusDisabled, common.UserStatusDeleted:
		// 禁用或删除：加入黑名单
		if err := redisClient.SAdd(ctx, key, userID).Err(); err != nil {
			logger.Error("syncUserStatusToRedis: 用户加入Redis黑名单失败: userID=%d, error=%v", userID, err)
		}
	case common.UserStatusNormal:
		// 启用：移出黑名单
		if err := redisClient.SRem(ctx, key, userID).Err(); err != nil {
			logger.Error("syncUserStatusToRedis: 用户移出Redis黑名单失败: userID=%d, error=%v", userID, err)
		}
	}
}
