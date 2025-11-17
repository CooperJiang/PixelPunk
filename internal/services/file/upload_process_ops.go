package file

/* Processing helpers split from upload_service.go (no behavior change). */

import (
	"fmt"
	"io"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/exif"
	pkgStorage "pixelpunk/pkg/storage"
	storageutils "pixelpunk/pkg/storage/utils"
	"strings"
	"time"
)

func processFile(ctx *UploadContext) error {
	src, err := ctx.File.Open()
	if err != nil {
		return errors.Wrap(err, errors.CodeFileUploadFailed, "打开上传文件失败")
	}
	defer src.Close()
	if ctx.OriginalFileData == nil {
		fileData, err := io.ReadAll(src)
		if err != nil {
			return errors.Wrap(err, errors.CodeFileUploadFailed, "读取文件数据失败")
		}
		ctx.OriginalFileData = fileData
	}
	fileHashStr := storageutils.CalculateDataMD5(ctx.OriginalFileData)
	ctx.FileHash = fileHashStr
	if err := checkDuplicateFile(ctx, fileHashStr); err != nil {
		return err
	}

	if exifData, err := exif.ExtractEXIFFromBytes(ctx.OriginalFileData); err == nil && exifData != nil {
		ctx.EXIFData = convertToFileEXIF(exifData)
	}

	src.Seek(0, 0)
	return processFileName(ctx)
}

// convertToFileEXIF 将 EXIF 数据转换为数据库模型
func convertToFileEXIF(data *exif.FileEXIFData) *models.FileEXIF {
	return &models.FileEXIF{
		Make:              data.Make,
		Model:             data.Model,
		LensModel:         data.LensModel,
		LensMake:          data.LensMake,
		SerialNumber:      data.SerialNumber,
		LensSerialNumber:  data.LensSerialNumber,
		FNumber:           data.FNumber,
		ExposureTime:      data.ExposureTime,
		ISO:               data.ISO,
		FocalLength:       data.FocalLength,
		FocalLengthIn35mm: data.FocalLengthIn35mm,
		GPSLatitude:       data.GPSLatitude,
		GPSLongitude:      data.GPSLongitude,
		GPSAltitude:       data.GPSAltitude,
		GPSLatitudeRef:    data.GPSLatitudeRef,
		GPSLongitudeRef:   data.GPSLongitudeRef,
		DateTime:          data.DateTime,
		DateTimeOriginal:  data.DateTimeOriginal,
		DateTimeDigitized: data.DateTimeDigitized,
		Orientation:       data.Orientation,
		ColorSpace:        data.ColorSpace,
		XResolution:       data.XResolution,
		YResolution:       data.YResolution,
		ResolutionUnit:    data.ResolutionUnit,
		ExposureProgram:   data.ExposureProgram,
		ExposureBias:      data.ExposureBias,
		MeteringMode:      data.MeteringMode,
		Flash:             data.Flash,
		WhiteBalance:      data.WhiteBalance,
		SceneCaptureType:  data.SceneCaptureType,
		DigitalZoomRatio:  data.DigitalZoomRatio,
		Contrast:          data.Contrast,
		Saturation:        data.Saturation,
		Sharpness:         data.Sharpness,
		SubjectDistance:   data.SubjectDistance,
		MaxApertureValue:  data.MaxApertureValue,
		Software:          data.Software,
		Artist:            data.Artist,
		Copyright:         data.Copyright,
		ImageDescription:  data.ImageDescription,
	}
}

func checkDuplicateFile(ctx *UploadContext, fileHash string) error {
	var existingImage models.File
	if err := database.DB.Where("user_id = ? AND md5_hash = ?", ctx.UserID, fileHash).
		Where("status <> ?", "pending_deletion").
		First(&existingImage).Error; err == nil {
		ctx.IsDuplicate = true
		ctx.OriginalFileID = existingImage.ID
		ctx.ReuseExistingFile = true
		ctx.ExistingFile = &existingImage
	}
	return nil
}

func processFileName(ctx *UploadContext) error {
	originalNameWithoutExt := strings.TrimSuffix(ctx.File.Filename, ctx.FileExt)
	ctx.OriginalName = originalNameWithoutExt
	ctx.SafeOriginalName = pkgStorage.SanitizeFileName(originalNameWithoutExt)
	if ctx.IsDuplicate {
		timeStr := time.Now().Format("20060102-150405")
		safeName := ctx.SafeOriginalName
		if len(safeName) > 40 {
			safeName = safeName[:40]
		}
		ctx.DisplayName = fmt.Sprintf("%s-%s", safeName, timeStr)
	} else {
		ctx.DisplayName = ctx.SafeOriginalName
	}
	return processFolderPath(ctx)
}

func processFolderPath(ctx *UploadContext) error {
	ctx.FolderPath = ""
	return nil
}

func handleAccessLevel(ctx *UploadContext) error {
	if ctx.AccessLevel == "" {
		ctx.AccessLevel = "public"
	}
	return nil
}
