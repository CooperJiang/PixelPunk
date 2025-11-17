package file

import (
	"github.com/gin-gonic/gin"

	"mime/multipart"
	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/watermark"
)

func processFileAndUploadWithWatermark(ctx *UploadContext) error {
	if err := processFile(ctx); err != nil {
		logger.Error("常规文件处理失败: %v", err)
		return err
	}

	if ctx.WatermarkEnabled && ctx.WatermarkConfig != "" {
		if err := applyWatermarkToFile(ctx); err != nil {
			logger.Warn("水印处理失败，使用原图上传: %v", err)
			// 记录失败原因，不中断上传流程
			ctx.WatermarkApplied = false
			ctx.WatermarkFailureReason = err.Error()
		} else {
			ctx.WatermarkApplied = true
		}
	}

	if err := executeUpload(ctx); err != nil {
		logger.Error("文件上传失败: %v", err)
		return err
	}
	return nil
}

func applyWatermarkToFile(ctx *UploadContext) error {
	if ctx.OriginalFileData == nil {
		return errors.New(errors.CodeFileUploadFailed, "原始文件数据不可用")
	}

	result, err := watermark.ProcessBytesWithConfigJSON(ctx.OriginalFileData, ctx.WatermarkConfig)
	if err != nil {
		logger.Error("水印合成失败: %v", err)
		return errors.Wrap(err, errors.CodeFileUploadFailed, "水印合成失败")
	}
	if !result.Success {
		logger.Error("水印合成返回失败: %s", result.ErrorMessage)
		return errors.New(errors.CodeFileUploadFailed, result.ErrorMessage)
	}

	if len(result.ProcessedData) > 0 {
		ctx.WatermarkWrapper = result.ProcessedData
	} else {
		logger.Warn("水印合成未返回数据，使用原图")
	}

	return nil
}

/* UploadFileWithWatermark 上传单张文件（支持水印） */
func UploadFileWithWatermark(c *gin.Context, userID uint, file *multipart.FileHeader, folderID, accessLevel string, optimize bool, storageDuration string, watermarkEnabled bool, watermarkConfig string) (*FileUploadResponse, error) {
	available, err := stats.CheckUserStorageAvailable(userID, file.Size)
	if err != nil {
		logger.Error("检查用户存储空间失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeInternal, "检查用户存储空间失败")
	}
	if !available {
		return nil, errors.New(errors.CodeStorageLimitExceeded, "存储空间不足，无法上传文件")
	}

	if exceeded, err := checkDailyUploadLimit(userID, 1); err != nil {
		logger.Warn("检查每日上传限制失败: %v", err)
	} else if exceeded {
		return nil, errors.New(errors.CodeUploadLimitExceeded, "已达到每日上传限制")
	}

	ctx := CreateUploadContextWithDuration(c, userID, file, folderID, accessLevel, optimize, storageDuration)

	if watermarkEnabled && watermarkConfig != "" {
		ctx.WatermarkEnabled = watermarkEnabled
		ctx.WatermarkConfig = watermarkConfig
	}

	if err := validateUploadRequest(ctx); err != nil {
		return nil, err
	}

	if err := processFileAndUploadWithWatermark(ctx); err != nil {
		return nil, err
	}

	if err := saveFileRecordAndStats(ctx); err != nil {
		logger.Error("保存文件记录失败: %v", err)
		return nil, err
	}

	return buildUploadResponse(ctx), nil
}

/* UploadFileBatchWithWatermark 批量上传文件（支持水印） */
func UploadFileBatchWithWatermark(c *gin.Context, userID uint, files []*multipart.FileHeader, folderID, accessLevel string, optimize bool, storageDuration string, watermarkEnabled bool, watermarkConfig string) ([]*FileUploadResponse, error) {
	results := make([]*FileUploadResponse, 0, len(files))

	for _, file := range files {
		result, err := UploadFileWithWatermark(c, userID, file, folderID, accessLevel, optimize, storageDuration, watermarkEnabled, watermarkConfig)
		if err != nil {
			logger.Error("批量上传中单个文件失败 %s: %v", file.Filename, err)
			continue
		}
		results = append(results, result)
	}

	return results, nil
}
