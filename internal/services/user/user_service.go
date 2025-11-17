package user

import (
	"fmt"
	"math/rand"
	"pixelpunk/internal/controllers/user/dto"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/auth"
	folderService "pixelpunk/internal/services/folder"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/internal/services/setting"
	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/email"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/utils"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

var userService *UserService

type UserService struct {
}

func InitUserService() {
	userService = &UserService{}
}

func GetUserService() *UserService {
	return userService
}

func Login(account, password string) (map[string]interface{}, string, error) {
	db := database.GetDB()
	var user models.User
	result := db.Where("username = ? OR email = ?", account, account).First(&user)
	if result.Error != nil {
		return nil, "", errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	securitySettings, err := setting.GetSettingsByGroupAsMap("security")

	maxLoginAttempts := 5
	accountLockoutMinutes := 30

	if err != nil {
		return nil, "", errors.New(errors.CodeInternal, "安全配置读取失败：security 组缺失")
	}

	jwtSecret := ""
	if val, ok := securitySettings.Settings["jwt_secret"]; ok {
		if secretStr, ok := val.(string); ok {
			jwtSecret = secretStr
		}
	}
	if strings.TrimSpace(jwtSecret) == "" {
		return nil, "", errors.New(errors.CodeInternal, "安全配置缺失：jwt_secret 未设置")
	}

	expiresHours := 0
	if val, ok := securitySettings.Settings["login_expire_hours"]; ok {
		if hours, ok := val.(float64); ok && hours > 0 {
			expiresHours = int(hours)
		}
	}
	if expiresHours <= 0 {
		return nil, "", errors.New(errors.CodeInternal, "安全配置缺失：login_expire_hours 未设置或非法")
	}

	if val, ok := securitySettings.Settings["max_login_attempts"]; ok {
		if attempts, ok := val.(float64); ok && attempts > 0 {
			maxLoginAttempts = int(attempts)
		}
	}
	if val, ok := securitySettings.Settings["account_lockout_minutes"]; ok {
		if minutes, ok := val.(float64); ok && minutes > 0 {
			accountLockoutMinutes = int(minutes)
		}
	}

	lockKey := fmt.Sprintf("user:login:lock:%d", user.ID)
	if cache.GetCache().Exists(lockKey) {
		ttl, err := cache.GetCache().TTL(lockKey)
		if err == nil && ttl > 0 {
			minutes := int(ttl.Minutes())
			seconds := int(ttl.Seconds()) % 60

			var timeMsg string
			if minutes > 0 {
				timeMsg = fmt.Sprintf("%d分钟", minutes)
				if seconds > 0 {
					timeMsg += fmt.Sprintf("%d秒", seconds)
				}
			} else {
				timeMsg = fmt.Sprintf("%d秒", seconds)
			}

			return nil, "", errors.New(errors.CodeForbidden, fmt.Sprintf("账户已被锁定，请%s后再试", timeMsg))
		}

		return nil, "", errors.New(errors.CodeForbidden, "账户已被锁定，请稍后再试")
	}

	attemptKey := fmt.Sprintf("user:login:attempts:%d", user.ID)
	attemptCount := 0
	if val, err := cache.GetCache().Get(attemptKey); err == nil && val != "" {
		if count, err := strconv.Atoi(val); err == nil {
			attemptCount = count
		}
	}

	if !utils.ComparePasswords(user.Password, password) {
		attemptCount++
		_ = cache.GetCache().Set(attemptKey, fmt.Sprintf("%d", attemptCount), time.Minute)

		if attemptCount >= maxLoginAttempts {
			_ = cache.GetCache().Set(lockKey, "1", time.Duration(accountLockoutMinutes)*time.Minute)
			_ = cache.GetCache().Del(attemptKey)

			return nil, "", errors.New(errors.CodeForbidden,
				fmt.Sprintf("密码错误次数过多，账户已被锁定%d分钟", accountLockoutMinutes))
		}

		return nil, "", errors.New(errors.CodeWrongPassword,
			fmt.Sprintf("密码错误，还有%d次尝试机会", maxLoginAttempts-attemptCount))
	}

	_ = cache.GetCache().Del(attemptKey)

	if !user.IsNormal() {
		return nil, "", errors.New(errors.CodeUserDisabled, "账号已被禁用")
	}

	token, err := auth.GenerateToken(user.ID, user.Username, int(user.Role), jwtSecret, expiresHours)
	if err != nil {
		return nil, "", errors.New(errors.CodeInternal, "生成token失败")
	}

	avatarFullPath := ""
	if user.Avatar != "" {
		avatarFullPath = utils.GetSystemFileURL(user.Avatar)
	} else {
		avatarFullPath = ""
	}

	userInfo := map[string]interface{}{
		"id":             user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"avatar":         user.Avatar,
		"avatarFullPath": avatarFullPath,
		"bio":            user.Bio,
		"website":        user.Website,
		"role":           user.Role,
		"status":         user.Status,
	}

	return userInfo, token, nil
}

func FindUsers() ([]models.User, error) {
	db := database.GetDB()
	var users []models.User
	result := db.Find(&users)
	return users, result.Error
}

func FindUserByID(id string) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func FindUserByEmail(email string) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func generateVerificationCode(email string, codeType string) string {
	rand.Seed(time.Now().UnixNano())

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	key := fmt.Sprintf("%s:%s:code", email, codeType)
	err := cache.GetCache().Set(key, code, 5*time.Minute)
	if err != nil {
		logger.Warn("存储验证码到缓存失败: %v", err)
		return ""
	}

	return code
}

func SendRegistrationCode(email string) error {
	registrationSettings, err := setting.GetSettingsByGroupAsMap("registration")
	if err != nil {
		return errors.New(errors.CodeInternal, "获取系统设置失败")
	}

	enableRegistration, ok := registrationSettings.Settings["enable_registration"].(bool)
	if !ok || !enableRegistration {
		return errors.New(errors.CodeForbidden, "管理员已关闭注册功能")
	}

	emailVerification, ok := registrationSettings.Settings["email_verification"].(bool)
	if !ok || !emailVerification {
		key := fmt.Sprintf("%s:%s:code", email, common.CodeTypeRegister)
		_ = cache.GetCache().Set(key, "999999", 5*time.Minute)
		return errors.New(errors.CodeValidationFailed, "管理员已关闭邮箱验证，请填写默认验证码为【999999】")
	}

	rateKey := fmt.Sprintf("%s:rate_limit:register", email)
	if cache.GetCache().Exists(rateKey) {
		return errors.New(errors.CodeRateLimited, "请求过于频繁，请一分钟后再试")
	}

	_, err = FindUserByEmail(email)
	if err == nil {
		return errors.New(errors.CodeEmailExists, "该邮箱已被注册")
	}

	code := generateVerificationCode(email, common.CodeTypeRegister)
	if code == "" {
		return errors.New(errors.CodeValidationFailed, "生成验证码失败")
	}

	err = sendVerificationEmail(email, code, common.CodeTypeRegister)
	if err != nil {
		return err
	}

	_ = cache.GetCache().Set(rateKey, "1", time.Minute)

	return nil
}

func SendResetPasswordCode(email string) error {
	rateKey := fmt.Sprintf("%s:rate_limit:reset_password", email)
	if cache.GetCache().Exists(rateKey) {
		return errors.New(errors.CodeRateLimited, "请求过于频繁，请一分钟后再试")
	}

	_, err := FindUserByEmail(email)
	if err != nil {
		return errors.New(errors.CodeUserNotFound, "该邮箱尚未注册")
	}

	code := generateVerificationCode(email, common.CodeTypeResetPassword)
	if code == "" {
		return errors.New(errors.CodeValidationFailed, "生成验证码失败")
	}

	err = sendVerificationEmail(email, code, common.CodeTypeResetPassword)
	if err != nil {
		return err
	}

	_ = cache.GetCache().Set(rateKey, "1", time.Minute)

	return nil
}

func ValidateCode(email, code, codeType string) bool {
	key := fmt.Sprintf("%s:%s:code", email, codeType)
	cachedCode, err := cache.GetCache().Get(key)
	if err != nil {
		logger.Warn("获取验证码失败: %v", err)
		return false
	}

	if code == cachedCode {
		_ = cache.GetCache().Del(key)
		return true
	}

	return false
}

func RegisterUser(username, email, password, code string) (uint, error) {
	db := database.GetDB()

	registrationSettings, err := setting.GetSettingsByGroupAsMap("registration")
	if err != nil {
		return 0, errors.New(errors.CodeInternal, "获取系统设置失败")
	}

	enableRegistration, ok := registrationSettings.Settings["enable_registration"].(bool)
	if !ok || !enableRegistration {
		return 0, errors.New(errors.CodeForbidden, "管理员已关闭注册功能")
	}

	emailVerification, ok := registrationSettings.Settings["email_verification"].(bool)
	if ok && emailVerification {
		if !ValidateCode(email, code, common.CodeTypeRegister) {
			return 0, errors.New(errors.CodeInvalidVerifyCode, "验证码无效或已过期")
		}
	}

	var count int64
	db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return 0, errors.New(errors.CodeUserExists, "用户名已存在")
	}

	db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return 0, errors.New(errors.CodeEmailExists, "邮箱已被注册")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, errors.New(errors.CodeInternal, "密码加密失败")
	}

	randomBio := common.GetRandomBio()

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	var user models.User
	err = db.Transaction(func(tx *gorm.DB) error {
		user = models.User{
			Username: username,
			Email:    email,
			Password: hashedPassword,
			Status:   common.UserStatusNormal,
			Role:     common.UserRoleUser,
			Bio:      randomBio, // 设置随机签名
		}

		if err := tx.Create(&user).Error; err != nil {
			return errors.New(errors.CodeDBCreateFailed, "创建用户失败")
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	if _, aliasErr := tenant.ResolveAlias(user.ID); aliasErr != nil {
		logger.Warn("生成用户路径别名失败: userID=%d, err=%v", user.ID, aliasErr)
	}

	if err := InitUserSettings(user.ID); err != nil {
		return 0, err
	}

	initialStorage := int64(50) * 1024 * 1024     // 默认50MB，转换为字节
	initialBandwidth := int64(1024) * 1024 * 1024 // 默认1GB，转换为字节

	if storageValue, exists := registrationSettings.Settings["user_initial_storage"]; exists {
		if storageInt, ok := storageValue.(float64); ok {
			initialStorage = int64(storageInt) * 1024 * 1024 // 转换为字节
		}
	}

	if bandwidthValue, exists := registrationSettings.Settings["user_initial_bandwidth"]; exists {
		if bandwidthInt, ok := bandwidthValue.(float64); ok {
			initialBandwidth = int64(bandwidthInt) * 1024 * 1024 // 转换为字节
		}
	}

	_, err = UpdateUserSettings(user.ID, initialStorage, initialBandwidth, "", false)
	if err != nil {
		return 0, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户存储空间和带宽设置失败")
	}

	if err := InitUserUsageStats(user.ID); err != nil {
		return 0, err
	}

	stats.GetStatsAdapter().RecordUserCreated()

	go sendRegistrationWelcomeMessage(user.ID, username, initialStorage, initialBandwidth)

	return user.ID, nil
}

func GetUserFolders(userID uint, query *dto.FolderQueryDTO) (interface{}, error) {
	if query.UserID != 0 && query.UserID != userID && query.PublicOnly {
		var targetUser models.User
		if err := database.GetDB().First(&targetUser, query.UserID).Error; err != nil {
			return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
		}

		folders, err := folderService.ListFolders(query.UserID, query.ParentID)
		if err != nil {
			return nil, err
		}

		var folderDTOs []dto.FolderStatDTO
		for _, folder := range folders {
			if folder.Permission == "public" {
				folderDTOs = append(folderDTOs, dto.FolderStatDTO{
					ID:          folder.ID,
					Name:        folder.Name,
					FileCount:   int(folder.FileCount),
					TotalSize:   0, // 公开文件夹不显示大小信息
					Description: folder.Description,
					Permission:  folder.Permission,
					CreatedAt:   time.Time(folder.CreatedAt).Format("2006-01-02 15:04:05"),
				})
			}
		}

		return dto.FolderListResponseDTO{Folders: folderDTOs}, nil
	}

	targetUserID := userID
	if query.UserID != 0 && query.UserID == userID {
		targetUserID = query.UserID
	}

	folders, err := folderService.ListFolders(targetUserID, query.ParentID)
	if err != nil {
		return nil, err
	}

	var folderDTOs []dto.FolderStatDTO

	if query.WithStats {
		for _, folder := range folders {
			var totalSize int64
			var totalImages int64

			database.GetDB().Model(&models.File{}).
				Where("folder_id = ?", folder.ID).
				Select("COALESCE(SUM(size), 0) as total_size, COUNT(*) as total_files").
				Row().Scan(&totalSize, &totalImages)

			folderDTOs = append(folderDTOs, dto.FolderStatDTO{
				ID:                 folder.ID,
				Name:               folder.Name,
				FileCount:          int(totalImages),
				TotalSize:          totalSize,
				TotalSizeFormatted: formatFileSize(totalSize),
				Description:        folder.Description,
				Permission:         folder.Permission,
				CreatedAt:          time.Time(folder.CreatedAt).Format("2006-01-02 15:04:05"),
			})
		}
	} else {
		for _, folder := range folders {
			folderDTOs = append(folderDTOs, dto.FolderStatDTO{
				ID:          folder.ID,
				Name:        folder.Name,
				FileCount:   int(folder.FileCount),
				TotalSize:   0,
				Description: folder.Description,
				Permission:  folder.Permission,
				CreatedAt:   time.Time(folder.CreatedAt).Format("2006-01-02 15:04:05"),
			})
		}
	}

	return dto.FolderListResponseDTO{Folders: folderDTOs}, nil
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

func sendVerificationEmail(emailAddr string, code string, codeType string) error {
	if !email.IsMailEnabled() {
		return errors.New(errors.CodeEmailServiceError, "邮件服务不可用，请联系管理员")
	}

	var subject string
	if codeType == common.CodeTypeRegister {
		subject = "注册验证码"
	} else {
		subject = "重置密码验证码"
	}

	err := email.SendMail(emailAddr, subject, fmt.Sprintf("您的验证码是: %s，5分钟内有效。", code))
	if err != nil {
		return errors.Wrap(err, errors.CodeEmailSendFailed, "发送邮件失败")
	}

	return nil
}

func ResetPassword(email, code, newPassword string) error {
	db := database.GetDB()

	if !ValidateCode(email, code, common.CodeTypeResetPassword) {
		return errors.New(errors.CodeInvalidVerifyCode, "验证码无效或已过期")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New(errors.CodeInternal, "密码加密失败")
	}

	result := db.Model(&models.User{}).Where("email = ?", email).Update("password", hashedPassword)
	if result.Error != nil {
		return errors.New(errors.CodeDBUpdateFailed, "更新密码失败")
	}

	if result.RowsAffected == 0 {
		return errors.New(errors.CodeUserNotFound, "未找到用户")
	}

	return nil
}

/* UpdatePassword 修改密码（需验证旧密码） */
func UpdatePassword(userID uint, oldPassword, newPassword string) error {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	if !utils.ComparePasswords(user.Password, oldPassword) {
		return errors.New(errors.CodeWrongPassword, "原密码错误")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New(errors.CodeInternal, "密码加密失败")
	}

	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return errors.New(errors.CodeDBUpdateFailed, "更新密码失败")
	}

	return nil
}

func SendChangeEmailCode(userID uint, newEmail string) error {
	currentUser, err := FindUserByID(fmt.Sprintf("%d", userID))
	if err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	rateKey := fmt.Sprintf("%d:rate_limit:change_email", userID)
	if cache.GetCache().Exists(rateKey) {
		return errors.New(errors.CodeRateLimited, "请求过于频繁，请一分钟后再试")
	}

	if currentUser.Email == newEmail {
		return errors.New(errors.CodeInvalidRequest, "新邮箱不能与当前邮箱相同")
	}

	existingUser, err := FindUserByEmail(newEmail)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return errors.New(errors.CodeEmailExists, "该邮箱已被注册")
	}

	code := generateVerificationCode(newEmail, common.CodeTypeChangeEmail)
	if code == "" {
		return errors.New(errors.CodeValidationFailed, "生成验证码失败")
	}

	err = sendVerificationEmail(newEmail, code, common.CodeTypeChangeEmail)
	if err != nil {
		return err
	}

	_ = cache.GetCache().Set(rateKey, "1", time.Minute)

	return nil
}

func ChangeEmail(userID uint, newEmail, code string) error {
	db := database.GetDB()

	if !ValidateCode(newEmail, code, common.CodeTypeChangeEmail) {
		return errors.New(errors.CodeInvalidVerifyCode, "验证码无效或已过期")
	}

	currentUser, err := FindUserByID(fmt.Sprintf("%d", userID))
	if err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	if currentUser.Email == newEmail {
		return errors.New(errors.CodeInvalidRequest, "新邮箱不能与当前邮箱相同")
	}

	existingUser, err := FindUserByEmail(newEmail)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return errors.New(errors.CodeEmailExists, "该邮箱已被注册")
	}

	result := db.Model(&models.User{}).Where("id = ?", userID).Update("email", newEmail)
	if result.Error != nil {
		return errors.New(errors.CodeDBUpdateFailed, "更新邮箱失败")
	}

	if result.RowsAffected == 0 {
		return errors.New(errors.CodeUserNotFound, "未找到用户")
	}

	return nil
}

func sendRegistrationWelcomeMessage(userID uint, username string, initialStorage, initialBandwidth int64) {
	siteName := "PixelPunk" // 默认值
	constructionSettings, err := setting.GetSettingsByGroupAsMap("construction")
	if err == nil {
		if name, ok := constructionSettings.Settings["site_name"].(string); ok && name != "" {
			siteName = name
		}
	}

	variables := map[string]interface{}{
		"username":     username,
		"site_name":    siteName,
		"storage_mb":   initialStorage / (1024 * 1024),   // 转换为MB
		"bandwidth_mb": initialBandwidth / (1024 * 1024), // 转换为MB
	}

	msgService := messageService.GetMessageService()
	if err := msgService.SendTemplateMessage(userID, common.MessageTypeAccountRegister, variables); err != nil {
		logger.Warn("发送注册欢迎消息失败: userID=%d, error=%v", userID, err)
	} else {
	}
}
