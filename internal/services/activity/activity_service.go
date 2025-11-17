package activity

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"sync"
	"time"
)

/* ActivityService 活动日志服务 */
type ActivityService struct {
	uploadBuffer map[string]*UploadBuffer // 用户上传缓冲，key格式：userID:folderID
	bufferMutex  sync.RWMutex             // 缓冲区锁
}

/* UploadBuffer 上传缓冲结构 */
type UploadBuffer struct {
	UserID     uint
	FolderID   string
	FileCount  int
	FolderName string
	TotalSize  int64
	Timer      *time.Timer
}

/* NewActivityService 创建活动日志服务实例 */
func NewActivityService() *ActivityService {
	return &ActivityService{
		uploadBuffer: make(map[string]*UploadBuffer),
	}
}

/* LogActivityParams 记录活动日志的参数 */
type LogActivityParams struct {
	UserID     *uint          `json:"user_id"`
	Type       string         `json:"type"`
	Data       map[string]any `json:"data"`
	Module     string         `json:"module"`
	EntityType string         `json:"entity_type"`
	EntityID   string         `json:"entity_id"`
	IsVisible  bool           `json:"is_visible"`
	Tags       string         `json:"tags"`
}

/* LogActivityAsync 异步记录活动日志 */
func (s *ActivityService) LogActivityAsync(params LogActivityParams) {
	go func() {
		var dataJSON json.RawMessage
		if params.Data != nil {
			data, err := json.Marshal(params.Data)
			if err != nil {
				logger.Error("序列化活动日志数据失败: %v", err)
				dataJSON = json.RawMessage("{}")
			} else {
				dataJSON = data
			}
		} else {
			dataJSON = json.RawMessage("{}")
		}

		activity := &models.ActivityLog{
			UserID:     params.UserID,
			Type:       params.Type,
			Data:       dataJSON,
			Module:     params.Module,
			EntityType: params.EntityType,
			EntityID:   params.EntityID,
			IsVisible:  params.IsVisible,
			Tags:       params.Tags,
		}

		db := database.GetDB()
		if db == nil {
			return // 数据库未初始化，跳过记录
		}

		if err := db.Create(activity).Error; err != nil {
			logger.Error("保存活动日志失败: %v", err)
		}
	}()
}

/* LogImageUploadDebounced 防抖记录文件上传（按文件夹分组，15秒内的上传合并为一条记录） */
func (s *ActivityService) LogImageUploadDebounced(userID uint, fileName string, fileSize int64, folderID string, folderName string) {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	bufferKey := fmt.Sprintf("%d:%s", userID, folderID)

	for key, existingBuffer := range s.uploadBuffer {
		if existingBuffer.UserID == userID && key != bufferKey {
			if existingBuffer.Timer != nil {
				existingBuffer.Timer.Stop()
			}
			s.doFlushUploadBuffer(key, existingBuffer)
		}
	}

	buffer, exists := s.uploadBuffer[bufferKey]
	if !exists {
		buffer = &UploadBuffer{
			UserID:     userID,
			FolderID:   folderID,
			FileCount:  0,
			FolderName: folderName,
			TotalSize:  0,
		}
		s.uploadBuffer[bufferKey] = buffer
	}

	buffer.FileCount++
	buffer.TotalSize += fileSize

	if buffer.Timer != nil {
		buffer.Timer.Stop()
	}

	buffer.Timer = time.AfterFunc(15*time.Second, func() {
		s.flushUploadBufferByKey(bufferKey)
	})
}

func (s *ActivityService) flushUploadBufferByKey(bufferKey string) {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	buffer, exists := s.uploadBuffer[bufferKey]
	if !exists {
		return
	}

	s.doFlushUploadBuffer(bufferKey, buffer)
}

func (s *ActivityService) doFlushUploadBuffer(bufferKey string, buffer *UploadBuffer) {

	params := LogActivityParams{
		UserID:     &buffer.UserID,
		Type:       "batch_image_upload",
		Module:     "file",
		EntityType: "folder",
		EntityID:   buffer.FolderID,
		IsVisible:  true,
		Tags:       "upload,batch",
		Data: map[string]any{
			"file_count":  buffer.FileCount,
			"folder_name": buffer.FolderName,
			"folder_id":   buffer.FolderID,
			"total_size":  buffer.TotalSize,
		},
	}

	s.LogActivityAsync(params)

	delete(s.uploadBuffer, bufferKey)
}

/* GetUserActivities 获取用户活动列表 */
func (s *ActivityService) GetUserActivities(userID uint, page, pageSize int) ([]models.ActivityLog, int64, map[string]interface{}, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, nil, fmt.Errorf("数据库未初始化")
	}

	query := db.Model(&models.ActivityLog{}).Where("user_id = ? AND is_visible = ?", userID, true)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, nil, err
	}

	var activities []models.ActivityLog
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&activities).Error; err != nil {
		return nil, 0, nil, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var todayUploads int64
	if err := db.Model(&models.File{}).Where("user_id = ? AND created_at >= ?", userID, today).Count(&todayUploads).Error; err != nil {
		todayUploads = 0
	}

	var totalViews int64
	var imageIDs []string
	if err := db.Model(&models.File{}).Where("user_id = ?", userID).Pluck("id", &imageIDs).Error; err == nil && len(imageIDs) > 0 {
		if err := db.Model(&models.FileStats{}).
			Where("file_id IN ?", imageIDs).
			Select("COALESCE(SUM(views), 0)").Scan(&totalViews).Error; err != nil {
			totalViews = 0
		}
	}

	todayStats := map[string]interface{}{
		"today_uploads": todayUploads,
		"total_views":   totalViews,
	}

	return activities, total, todayStats, nil
}

var globalService *ActivityService

func init() {
	globalService = NewActivityService()
}

/* GetService 获取全局service实例 */
func GetService() *ActivityService {
	return globalService
}

/* LogUserLogin 记录用户登录 */
func LogUserLogin(userID uint, username, clientIP string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "user_login",
		Module:     "auth",
		EntityType: "user",
		EntityID:   fmt.Sprintf("%d", userID),
		IsVisible:  true,
		Tags:       "login",
		Data: map[string]any{
			"username":   username,
			"ip_address": clientIP,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogShareCreate 记录分享创建 */
func LogShareCreate(userID uint, shareID string, shareType string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "share_create",
		Module:     "share",
		EntityType: "share",
		EntityID:   shareID,
		IsVisible:  true,
		Tags:       "share,create",
		Data: map[string]any{
			"share_id":   shareID,
			"share_type": shareType,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogShareMilestone 记录分享访问里程碑 */
func LogShareMilestone(userID uint, shareID string, visitCount int) {

	params := LogActivityParams{
		UserID:     &userID,
		Type:       "share_milestone",
		Module:     "share",
		EntityType: "share",
		EntityID:   shareID,
		IsVisible:  true,
		Tags:       "share,milestone",
		Data: map[string]any{
			"share_id":    shareID,
			"visit_count": visitCount,
			"milestone":   visitCount,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogBatchDelete 记录用户批量删除 */
func LogBatchDelete(userID uint, imageCount int, folderName string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "batch_delete",
		Module:     "file",
		EntityType: "file",
		EntityID:   folderName,
		IsVisible:  true,
		Tags:       "delete,batch",
		Data: map[string]any{
			"file_count":  imageCount,
			"folder_name": folderName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogAdminDelete 记录管理员删除用户文件 */
func LogAdminDelete(userID uint, fileID string, adminID uint) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "admin_delete",
		Module:     "admin",
		EntityType: "file",
		EntityID:   fileID,
		IsVisible:  true,
		Tags:       "delete,admin",
		Data: map[string]any{
			"file_id":  fileID,
			"admin_id": adminID,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogSystemCleanup 记录系统清理 */
func LogSystemCleanup(cleanedCount int, cleanupType string) {

	params := LogActivityParams{
		UserID:     nil, // 系统操作，无用户ID
		Type:       "system_cleanup",
		Module:     "system",
		EntityType: "cleanup",
		EntityID:   cleanupType,
		IsVisible:  false, // 系统操作对用户不可见
		Tags:       "system,cleanup",
		Data: map[string]any{
			"cleaned_count": cleanedCount,
			"cleanup_type":  cleanupType,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogFolderCreate 记录文件夹创建 */
func LogFolderCreate(userID uint, folderName string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "folder_create",
		Module:     "folder",
		EntityType: "folder",
		EntityID:   folderName,
		IsVisible:  true,
		Tags:       "folder,create",
		Data: map[string]any{
			"folder_name": folderName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogFolderRename 记录文件夹重命名 */
func LogFolderRename(userID uint, oldName, newName string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "folder_rename",
		Module:     "folder",
		EntityType: "folder",
		EntityID:   newName,
		IsVisible:  true,
		Tags:       "folder,rename",
		Data: map[string]any{
			"old_name": oldName,
			"new_name": newName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogFolderDelete 记录文件夹删除 */
func LogFolderDelete(userID uint, folderName string, imageCount int) {

	params := LogActivityParams{
		UserID:     &userID,
		Type:       "folder_delete",
		Module:     "folder",
		EntityType: "folder",
		EntityID:   folderName,
		IsVisible:  true,
		Tags:       "folder,delete",
		Data: map[string]any{
			"folder_name": folderName,
			"file_count":  imageCount,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogFolderAccessLevelChange 记录文件夹权限切换 */
func LogFolderAccessLevelChange(userID uint, folderName string, oldLevel string, newLevel string) {
	var levelText = map[string]string{
		"private": "私密",
		"public":  "公开",
	}

	oldText := levelText[oldLevel]
	newText := levelText[newLevel]
	if oldText == "" {
		oldText = oldLevel
	}
	if newText == "" {
		newText = newLevel
	}

	params := LogActivityParams{
		UserID:     &userID,
		Type:       "folder_access_level_change",
		Module:     "folder",
		EntityType: "folder",
		EntityID:   folderName,
		IsVisible:  true,
		Tags:       "folder,access_level,permission",
		Data: map[string]any{
			"folder_name":      folderName,
			"old_access_level": oldLevel,
			"new_access_level": newLevel,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogImageAccessLevelChange 记录文件权限切换（兼容旧命名） */
func LogImageAccessLevelChange(userID uint, imageName string, oldLevel string, newLevel string) {
	var levelText = map[string]string{
		"private":  "私密",
		"public":   "公开",
		"unlisted": "不公开列表",
	}

	oldText := levelText[oldLevel]
	newText := levelText[newLevel]
	if oldText == "" {
		oldText = oldLevel
	}
	if newText == "" {
		newText = newLevel
	}

	params := LogActivityParams{
		UserID:     &userID,
		Type:       "file_access_level_change",
		Module:     "file",
		EntityType: "file",
		EntityID:   imageName,
		IsVisible:  true,
		Tags:       "file,access_level,permission",
		Data: map[string]any{
			"file_name":        imageName,
			"old_access_level": oldLevel,
			"new_access_level": newLevel,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogFileAccessLevelChange 记录文件权限切换（新语义命名） */
func LogFileAccessLevelChange(userID uint, fileName string, oldLevel string, newLevel string) {
	LogImageAccessLevelChange(userID, fileName, oldLevel, newLevel)
}

/* LogProfileUpdate 记录个人资料更新 */
func LogProfileUpdate(userID uint, fieldType string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "profile_update",
		Module:     "user",
		EntityType: "user",
		EntityID:   fmt.Sprintf("%d", userID),
		IsVisible:  true,
		Tags:       "profile,update",
		Data: map[string]any{
			"field_type": fieldType,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogPasswordChange 记录密码修改 */
func LogPasswordChange(userID uint) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "password_change",
		Module:     "user",
		EntityType: "user",
		EntityID:   fmt.Sprintf("%d", userID),
		IsVisible:  true,
		Tags:       "security,password",
		Data: map[string]any{
			"action": "password_change",
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogEmailChange 记录邮箱更改 */
func LogEmailChange(userID uint, newEmail string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "email_change",
		Module:     "user",
		EntityType: "user",
		EntityID:   fmt.Sprintf("%d", userID),
		IsVisible:  true,
		Tags:       "security,email",
		Data: map[string]any{
			"new_email": newEmail,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogHotlinkProtectionChange 记录防盗链设置修改 */
func LogHotlinkProtectionChange(userID uint, settingsData map[string]any) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "hotlink_protection_change",
		Module:     "security",
		EntityType: "settings",
		EntityID:   fmt.Sprintf("user_%d_hotlink", userID),
		IsVisible:  true,
		Tags:       "security,hotlink,settings",
		Data:       settingsData, // 存储完整的设置数据作为快照
	}

	globalService.LogActivityAsync(params)
}

/* LogImageUploadByID 通过imageID和folderID记录文件上传日志（推荐使用） */
func LogImageUploadByID(fileID string, folderID string) {
	go func() {
		db := database.GetDB()
		if db == nil {
			logger.Error("数据库未初始化，跳过活动日志记录")
			return
		}

		var file models.File
		if err := db.Where("id = ?", fileID).First(&file).Error; err != nil {
			logger.Error("查询文件信息失败，fileID: %s, error: %v", fileID, err)
			return
		}

		folderName := "根目录"
		if folderID != "" {
			var folder models.Folder
			if err := db.Where("id = ?", folderID).First(&folder).Error; err != nil {
				logger.Warn("查询文件夹信息失败，folderID: %s, 使用默认名称, error: %v", folderID, err)
			} else {
				folderName = folder.Name
			}
		}

		globalService.LogImageUploadDebounced(
			file.UserID,
			file.OriginalName,
			file.Size,
			folderID,
			folderName,
		)
	}()
}

/* LogImageExpired 记录用户文件过期删除 */
func LogImageExpired(userID uint, expiredImages []models.File) {
	if len(expiredImages) == 0 {
		return
	}

	imageData := make([]map[string]any, len(expiredImages))
	for i, file := range expiredImages {
		imageData[i] = map[string]any{
			"id":         file.ID,
			"name":       file.OriginalName,
			"expires_at": file.ExpiresAt.Format("2006-01-02 15:04:05"),
			"size":       file.Size,
			"folder_id":  file.FolderID,
		}
	}

	params := LogActivityParams{
		UserID:     &userID,
		Type:       "file_expired",
		Module:     "system",
		EntityType: "file",
		EntityID:   fmt.Sprintf("expired_batch_%d", len(expiredImages)),
		IsVisible:  true, // 对用户可见
		Tags:       "expired,delete,automatic",
		Data: map[string]any{
			"expired_count": len(expiredImages),
			"files":         imageData,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogGuestImageExpired 记录游客文件过期删除（系统日志） */
func LogGuestImageExpired(expiredImages []models.File) {
	if len(expiredImages) == 0 {
		return
	}

	imageData := make([]map[string]any, len(expiredImages))
	for i, file := range expiredImages {
		imageData[i] = map[string]any{
			"id":         file.ID,
			"name":       file.OriginalName,
			"expires_at": file.ExpiresAt.Format("2006-01-02 15:04:05"),
			"size":       file.Size,
			"ip":         getImageIP(file), // 如果有记录IP的话
		}
	}

	params := LogActivityParams{
		UserID:     nil, // 游客无用户ID，系统操作
		Type:       "guest_file_expired",
		Module:     "system",
		EntityType: "guest_file",
		EntityID:   fmt.Sprintf("guest_expired_batch_%d", len(expiredImages)),
		IsVisible:  false, // 系统操作，对普通用户不可见，但管理员可以看到
		Tags:       "guest,expired,delete,automatic",
		Data: map[string]any{
			"expired_count": len(expiredImages),
			"files":         imageData,
		},
	}

	globalService.LogActivityAsync(params)
}

func getImageIP(file models.File) string {
	return ""
}

/* LogUserRegister 记录用户注册 */
func LogUserRegister(userID uint, username string, email string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "user_register",
		Module:     "auth",
		EntityType: "user",
		EntityID:   fmt.Sprintf("%d", userID),
		IsVisible:  true,
		Tags:       "register,welcome",
		Data: map[string]any{
			"username": username,
			"email":    email,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogFileDelete 记录单个文件删除 */
func LogFileDelete(userID uint, fileName string, fileID string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "file_delete",
		Module:     "file",
		EntityType: "file",
		EntityID:   fileID,
		IsVisible:  true,
		Tags:       "delete,file",
		Data: map[string]any{
			"file_id":   fileID,
			"file_name": fileName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogFileRename 记录文件重命名 */
func LogFileRename(userID uint, oldName string, newName string, fileID string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "file_rename",
		Module:     "file",
		EntityType: "file",
		EntityID:   fileID,
		IsVisible:  true,
		Tags:       "rename,file",
		Data: map[string]any{
			"file_id":  fileID,
			"old_name": oldName,
			"new_name": newName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogFileMove 记录文件移动 */
func LogFileMove(userID uint, fileName string, fileID string, oldFolderName string, newFolderName string) {
	if oldFolderName == "" {
		oldFolderName = "根目录"
	}
	if newFolderName == "" {
		newFolderName = "根目录"
	}

	params := LogActivityParams{
		UserID:     &userID,
		Type:       "file_move",
		Module:     "file",
		EntityType: "file",
		EntityID:   fileID,
		IsVisible:  true,
		Tags:       "move,file",
		Data: map[string]any{
			"file_id":         fileID,
			"file_name":       fileName,
			"old_folder_name": oldFolderName,
			"new_folder_name": newFolderName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogAPIKeyCreate 记录API密钥创建 */
func LogAPIKeyCreate(userID uint, keyName string, keyID string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "apikey_create",
		Module:     "apikey",
		EntityType: "apikey",
		EntityID:   keyID,
		IsVisible:  true,
		Tags:       "apikey,create",
		Data: map[string]any{
			"key_id":   keyID,
			"key_name": keyName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogAPIKeyDelete 记录API密钥删除 */
func LogAPIKeyDelete(userID uint, keyName string, keyID string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "apikey_delete",
		Module:     "apikey",
		EntityType: "apikey",
		EntityID:   keyID,
		IsVisible:  true,
		Tags:       "apikey,delete",
		Data: map[string]any{
			"key_id":   keyID,
			"key_name": keyName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogShareDelete 记录分享删除 */
func LogShareDelete(userID uint, shareID string, shareType string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "share_delete",
		Module:     "share",
		EntityType: "share",
		EntityID:   shareID,
		IsVisible:  true,
		Tags:       "share,delete",
		Data: map[string]any{
			"share_id":   shareID,
			"share_type": shareType,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogAPIKeyToggleStatus 记录API密钥状态切换 */
func LogAPIKeyToggleStatus(userID uint, keyName string, keyID string, newStatus int) {
	var actionType string

	if newStatus == 1 { // APIKeyStatusActive
		actionType = "enable"
	} else {
		actionType = "disable"
	}

	params := LogActivityParams{
		UserID:     &userID,
		Type:       "apikey_toggle_status",
		Module:     "apikey",
		EntityType: "apikey",
		EntityID:   keyID,
		IsVisible:  true,
		Tags:       fmt.Sprintf("apikey,%s", actionType),
		Data: map[string]any{
			"key_id":     keyID,
			"key_name":   keyName,
			"new_status": newStatus,
			"action":     actionType,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogAPIKeyRegenerate 记录API密钥重新生成 */
func LogAPIKeyRegenerate(userID uint, keyName string, keyID string) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "apikey_regenerate",
		Module:     "apikey",
		EntityType: "apikey",
		EntityID:   keyID,
		IsVisible:  true,
		Tags:       "apikey,regenerate,security",
		Data: map[string]any{
			"key_id":   keyID,
			"key_name": keyName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogRandomAPICreate 记录随机图片API创建 */
func LogRandomAPICreate(userID uint, apiName string, apiID uint) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "random_api_create",
		Module:     "random_api",
		EntityType: "random_api",
		EntityID:   fmt.Sprintf("%d", apiID),
		IsVisible:  true,
		Tags:       "random_api,create",
		Data: map[string]any{
			"api_id":   apiID,
			"api_name": apiName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogRandomAPIDelete 记录随机图片API删除 */
func LogRandomAPIDelete(userID uint, apiName string, apiID uint) {
	params := LogActivityParams{
		UserID:     &userID,
		Type:       "random_api_delete",
		Module:     "random_api",
		EntityType: "random_api",
		EntityID:   fmt.Sprintf("%d", apiID),
		IsVisible:  true,
		Tags:       "random_api,delete",
		Data: map[string]any{
			"api_id":   apiID,
			"api_name": apiName,
		},
	}

	globalService.LogActivityAsync(params)
}

/* LogRandomAPIToggleStatus 记录随机图片API状态切换 */
func LogRandomAPIToggleStatus(userID uint, apiName string, apiID uint, newStatus int) {
	var actionType string

	if newStatus == 1 { // RandomImageAPIStatusActive
		actionType = "enable"
	} else {
		actionType = "disable"
	}

	params := LogActivityParams{
		UserID:     &userID,
		Type:       "random_api_toggle_status",
		Module:     "random_api",
		EntityType: "random_api",
		EntityID:   fmt.Sprintf("%d", apiID),
		IsVisible:  true,
		Tags:       fmt.Sprintf("random_api,%s", actionType),
		Data: map[string]any{
			"api_id":     apiID,
			"api_name":   apiName,
			"new_status": newStatus,
			"action":     actionType,
		},
	}

	globalService.LogActivityAsync(params)
}
