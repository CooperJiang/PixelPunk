package file

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"pixelpunk/internal/models"
	"pixelpunk/internal/services/storage"
	newstorage "pixelpunk/pkg/storage"

	"github.com/google/uuid"
)

/* StorageChannelRepository 存储渠道仓库实现 */
type StorageChannelRepository struct{}

/* GetChannel 获取存储渠道 */
func (r *StorageChannelRepository) GetChannel(channelID string) (*models.StorageChannel, error) {
	return storage.GetChannelByID(channelID)
}

/* GetChannelConfig 获取渠道配置 */
func (r *StorageChannelRepository) GetChannelConfig(channelID string) (map[string]interface{}, error) {
	return storage.GetChannelConfigMap(channelID)
}

/* GetDefaultChannel 获取默认渠道 */
func (r *StorageChannelRepository) GetDefaultChannel() (*models.StorageChannel, error) {
	return storage.GetDefaultChannel()
}

/* GetActiveChannels 获取活跃渠道列表 */
func (r *StorageChannelRepository) GetActiveChannels() ([]*models.StorageChannel, error) {
	defaultChannel, err := r.GetDefaultChannel()
	if err != nil {
		return nil, err
	}

	return []*models.StorageChannel{defaultChannel}, nil
}

var storageService *newstorage.Storage

func ensureStorageServiceInitialized() (*newstorage.Storage, error) {
	if storageService != nil {
		return storageService, nil
	}

	channelRepo := &StorageChannelRepository{}

	storageService = newstorage.New(channelRepo)
	return storageService, nil
}

/* GetStorageServiceInstance 获取存储服务实例 */
func GetStorageServiceInstance() (*newstorage.Storage, error) {
	return ensureStorageServiceInitialized()
}

func convertToNewStorageRequest(ctx *UploadContext) *newstorage.UploadRequest {
	req := &newstorage.UploadRequest{
		File:          ctx.File,
		UserID:        ctx.UserID,
		FolderPath:    ctx.FolderPath,
		FileName:      ctx.File.Filename,
		ContentType:   "", // 将自动检测
		Quality:       90, // 原图使用高质量
		GenerateThumb: true,
		Compress:      false, // 原图默认不压缩尺寸
		WebPEnabled:   false, // 原图不启用WebP转换，保持原始格式
	}

	if ctx.WatermarkWrapper != nil {
		if processedData, ok := ctx.WatermarkWrapper.([]byte); ok {
			req.ProcessedData = processedData
		}
	}

	if ctx.StorageChannel != nil {
		req.ChannelID = ctx.StorageChannel.ID
	}

	if ctx.CompressOptions != nil {
		req.Quality = ctx.CompressOptions.Quality
		req.ThumbWidth = ctx.CompressOptions.MaxWidth
		req.ThumbHeight = ctx.CompressOptions.MaxHeight
		req.ThumbQuality = ctx.CompressOptions.Quality

		if ctx.CompressOptions.MaxWidth > 2000 || ctx.CompressOptions.MaxHeight > 2000 {
			req.MaxWidth = ctx.CompressOptions.MaxWidth
			req.MaxHeight = ctx.CompressOptions.MaxHeight
			req.Compress = true
		}
	}

	req.FileName = generateUniqueFileName(ctx.File.Filename)

	return req
}

func convertFromNewStorageResult(result *newstorage.UploadResult) *UploadResult {
	return &UploadResult{
		URL:                       result.URL,
		LocalUrlPath:              result.OriginalPath,
		ThumbUrl:                  result.ThumbnailURL,
		LocalThumbPath:            result.ThumbnailPath,
		RemoteUrl:                 result.RemoteURL,
		RemoteThumbUrl:            result.RemoteThumbURL,
		Width:                     result.Width,
		Height:                    result.Height,
		ThumbnailGenerationFailed: result.ThumbnailGenerationFailed,
		ThumbnailFailureReason:    result.ThumbnailFailureReason,
	}
}

func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)

	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	timestamp := fmt.Sprintf("%04d", time.Now().Unix()%10000) // 取时间戳后4位
	uniqueName := uuidStr + timestamp

	return uniqueName + ext
}
