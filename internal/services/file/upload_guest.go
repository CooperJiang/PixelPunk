package file

import (
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"

	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
)

/* GuestUpload 游客上传文件 */
func GuestUpload(c *gin.Context, file *multipart.FileHeader, folderID, accessLevel string, optimize bool, storageDuration, fingerprint string) (*FileDetailResponse, int, error) {

	storageConfig, err := setting.CreateStorageConfig()
	if err != nil {
		logger.Error("GuestUpload服务: 创建存储配置失败, err=%v", err)
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "创建存储配置失败")
	}

	if accessLevel == "" {
		accessLevel = storageConfig.GetGuestDefaultAccessLevel()
	}

	if storageDuration == "" {
		storageDuration = storageConfig.GetDefaultDuration(true) // true表示游客模式
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	service := GetGuestUploadLimitService()
	if service == nil {
		logger.Error("GuestUpload服务: GetGuestUploadLimitService 返回 nil")
		return nil, 0, errors.New(errors.CodeInternal, "游客上传限制服务初始化失败")
	}

	ipAllowed, ipErr := service.CheckIPUploadLimit(ip)
	if ipErr != nil {
		logger.Error("GuestUpload服务: CheckIPUploadLimit 失败, err=%v", ipErr)
		return nil, 0, ipErr
	}
	if !ipAllowed {
		return nil, 0, errors.New(errors.CodeUploadLimitExceeded, "IP每日上传次数已达上限")
	}

	var limit *models.GuestUploadLimit
	allowed, remainingCount, err := service.CheckUploadLimit(fingerprint)
	if allowed && err == nil {
		limit, _ = getGuestUploadLimitByFingerprint(fingerprint)
	}
	if err != nil {
		logService := GetGuestUploadLogService()
		guestLimit, _ := getGuestUploadLimitByFingerprint(fingerprint)
		guestLimitID := uint(0)
		if guestLimit != nil {
			guestLimitID = guestLimit.ID
		}
		log := &models.GuestUploadLog{
			GuestUploadLimitID: guestLimitID,
			Fingerprint:        fingerprint,
			IP:                 ip,
			UserAgent:          userAgent,
			StorageDuration:    storageDuration,
			Status:             "failed",
			Reason:             err.Error(),
			FileSize:           file.Size,
			FileName:           file.Filename,
			OriginalName:       file.Filename,
		}
		if err := logService.RecordGuestUpload(log); err != nil {
			logger.Error("记录游客上传日志失败: %v", err)
		}
		return nil, 0, err
	}
	if !allowed {
		logService := GetGuestUploadLogService()
		guestLimit, _ := getGuestUploadLimitByFingerprint(fingerprint)
		guestLimitID := uint(0)
		if guestLimit != nil {
			guestLimitID = guestLimit.ID
		}
		log := &models.GuestUploadLog{
			GuestUploadLimitID: guestLimitID,
			Fingerprint:        fingerprint,
			IP:                 ip,
			UserAgent:          userAgent,
			StorageDuration:    storageDuration,
			Status:             "blocked",
			Reason:             "上传次数已达限制",
			FileSize:           file.Size,
			FileName:           file.Filename,
			OriginalName:       file.Filename,
		}
		if err := logService.RecordGuestUpload(log); err != nil {
			logger.Error("记录游客上传日志失败: %v", err)
		}
		return nil, 0, errors.New(errors.CodeUploadLimitExceeded, "上传次数已达限制")
	}

	ctx := CreateUploadContextWithDuration(c, 0, file, folderID, accessLevel, optimize, storageDuration)
	ctx.GuestFingerprint = fingerprint
	ctx.GuestIP = ip
	ctx.GuestUserAgent = userAgent

	imgInfo, err := UploadFileWithDuration(c, 0, file, folderID, accessLevel, optimize, storageDuration)
	if err != nil {
		logService := GetGuestUploadLogService()
		guestLimit, _ := getGuestUploadLimitByFingerprint(fingerprint)
		guestLimitID := uint(0)
		if guestLimit != nil {
			guestLimitID = guestLimit.ID
		}
		log := &models.GuestUploadLog{
			GuestUploadLimitID: guestLimitID,
			Fingerprint:        fingerprint,
			IP:                 ip,
			UserAgent:          userAgent,
			StorageDuration:    storageDuration,
			Status:             "failed",
			Reason:             err.Error(),
			FileSize:           file.Size,
			FileName:           file.Filename,
			OriginalName:       file.Filename,
		}
		if err := logService.RecordGuestUpload(log); err != nil {
			logger.Error("记录游客上传日志失败: %v", err)
		}
		return nil, 0, err
	}

	if err := service.RecordUpload(fingerprint, file.Size); err != nil {
		logger.Error("记录游客上传失败: %v", err)
	}

	logService := GetGuestUploadLogService()
	guestLimitID := uint(0)
	if limit != nil {
		guestLimitID = limit.ID
	}
	var expiresAt *time.Time
	expireTime := common.CalculateExpiryTime(storageDuration)
	if !expireTime.IsZero() {
		expiresAt = &expireTime
	}
	log := &models.GuestUploadLog{
		GuestUploadLimitID: guestLimitID,
		FileID:             imgInfo.ID,
		Fingerprint:        fingerprint,
		IP:                 ip,
		UserAgent:          userAgent,
		StorageDuration:    storageDuration,
		ExpiresAt:          expiresAt,
		Status:             "success",
		FileSize:           file.Size,
		FileName:           imgInfo.OriginalName,
		OriginalName:       file.Filename,
	}
	if err := logService.RecordGuestUpload(log); err != nil {
		logger.Error("记录游客上传日志失败: %v", err)
	}

	return imgInfo, remainingCount, nil
}

/* GuestUploadWithWatermark 游客上传文件（支持水印） */
func GuestUploadWithWatermark(c *gin.Context, file *multipart.FileHeader, folderID, accessLevel string, optimize bool, storageDuration string, fingerprint string, watermarkEnabled bool, watermarkConfig string) (*FileUploadResponse, int, error) {

	remainingCount := 10 // 默认剩余次数

	ctx := CreateUploadContextWithDuration(c, 0, file, folderID, accessLevel, optimize, storageDuration)
	ctx.IsGuestUpload = true
	ctx.GuestFingerprint = fingerprint

	if watermarkEnabled && watermarkConfig != "" {
		ctx.WatermarkEnabled = watermarkEnabled
		ctx.WatermarkConfig = watermarkConfig
	}

	if err := validateUploadRequest(ctx); err != nil {
		return nil, remainingCount, err
	}

	if err := processFileAndUploadWithWatermark(ctx); err != nil {
		return nil, remainingCount, err
	}

	if err := saveFileRecordAndStats(ctx); err != nil {
		logger.Error("保存文件记录失败: %v", err)
		return nil, remainingCount, err
	}

	imgInfo := buildUploadResponse(ctx)
	return imgInfo, remainingCount, nil
}
