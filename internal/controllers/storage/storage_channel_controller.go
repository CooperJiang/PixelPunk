package storage

import (
	"io"
	"strings"

	"pixelpunk/internal/controllers/storage/dto"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/storage"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	storagemod "pixelpunk/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListChannels(ctx *gin.Context) {
	channels, err := storage.GetAllChannels()
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeDBQueryFailed, "获取存储渠道列表失败"))
		return
	}

	errors.ResponseSuccess(ctx, channels, "获取存储渠道列表成功")
}

func GetChannel(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	channel, err := storage.GetChannelByID(channelID)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeNotFound, "存储渠道不存在"))
		return
	}

	errors.ResponseSuccess(ctx, channel, "获取存储渠道成功")
}

func CreateChannel(ctx *gin.Context) {
	req, err := common.ValidateRequest[dto.CreateChannelDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	channel := models.StorageChannel{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Type:      req.Type,
		IsDefault: req.IsDefault,
		Remark:    req.Remark,
	}

	if req.Status != nil {
		channel.Status = *req.Status
	} else {
		channel.Status = 1 // 默认启用
	}

	if err := storage.CreateChannel(&channel, req.Configs); err != nil {
		if _, ok := err.(*errors.Error); ok {
			errors.HandleError(ctx, err)
		} else {
			errors.HandleError(ctx, errors.Wrap(err, errors.CodeDBCreateFailed, "创建存储渠道失败"))
		}
		return
	}

	errors.ResponseSuccess(ctx, channel, "创建存储渠道成功")
}

func UpdateChannel(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	req, err := common.ValidateRequest[dto.UpdateChannelDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	channel, err := storage.GetChannelByID(channelID)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeNotFound, "存储渠道不存在"))
		return
	}

	if req.Name != "" {
		channel.Name = req.Name
	}
	if req.Remark != "" {
		channel.Remark = req.Remark
	}
	if req.IsDefault != nil {
		channel.IsDefault = *req.IsDefault
	}
	if req.Status != nil {
		channel.Status = *req.Status
	}

	if err := storage.UpdateChannel(channel); err != nil {
		if _, ok := err.(*errors.Error); ok {
			errors.HandleError(ctx, err)
		} else {
			errors.HandleError(ctx, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新存储渠道失败"))
		}
		return
	}

	_ = storage.RefreshChannelCache(channelID)

	if len(req.Configs) > 0 {
		if err := storage.UpdateChannelConfigs(channelID, req.Configs); err != nil {
			if _, ok := err.(*errors.Error); ok {
				errors.HandleError(ctx, err)
			} else {
				errors.HandleError(ctx, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新渠道配置失败"))
			}
			return
		}
	}

	if err := storage.RefreshChannelCache(channelID); err != nil {
	}

	errors.ResponseSuccess(ctx, channel, "更新存储渠道成功")
}

func DeleteChannel(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	if err := storage.DeleteChannel(channelID); err != nil {
		if _, ok := err.(*errors.Error); ok {
			errors.HandleError(ctx, err)
		} else {
			errors.HandleError(ctx, errors.Wrap(err, errors.CodeDBDeleteFailed, "删除存储渠道失败"))
		}
		return
	}

	errors.ResponseSuccess(ctx, nil, "删除存储渠道成功")
}

func ListSupportedTypes(ctx *gin.Context) {
	types := storagemod.GetSupportedTypes()
	allowed := map[string]bool{
		"local":     true,
		"s3":        true,
		"oss":       true,
		"cos":       true,
		"qiniu":     true,
		"upyun":     true,
		"r2":        true,
		"rainyun":   true,
		"azureblob": true,
		"webdav":    true,
		"sftp":      true,
		"ftp":       true,
	}
	labelMap := map[string]string{
		"local":     "本地存储",
		"s3":        "通用 S3",
		"oss":       "阿里云 OSS",
		"cos":       "腾讯云COS",
		"qiniu":     "七牛云 Kodo",
		"upyun":     "又拍云 Upyun",
		"r2":        "Cloudflare R2",
		"rainyun":   "雨云 RainYun",
		"azureblob": "Azure Blob Storage",
		"webdav":    "WebDAV",
		"sftp":      "SFTP (基于 SSH)",
		"ftp":       "FTP",
	}
	var out []map[string]string
	for _, t := range types {
		if !allowed[t] {
			continue
		}
		out = append(out, map[string]string{
			"value": t,
			"label": func() string {
				if v, ok := labelMap[t]; ok {
					return v
				}
				return t
			}(),
		})
	}
	errors.ResponseSuccess(ctx, out, "获取支持的存储类型成功")
}

func GetConfigTemplates(ctx *gin.Context) {
	t := ctx.Param("type")
	if strings.TrimSpace(t) == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道类型参数"))
		return
	}
	if templates, ok := models.StorageConfigTemplates[t]; ok {
		errors.ResponseSuccess(ctx, templates, "获取配置模板成功")
		return
	}
	errors.HandleError(ctx, errors.New(errors.CodeNotFound, "不支持的存储类型"))
}

func GetChannelConfigs(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	configs, err := storage.GetChannelConfigs(channelID)
	if err != nil {
		if _, ok := err.(*errors.Error); ok {
			errors.HandleError(ctx, err)
		} else {
			errors.HandleError(ctx, errors.Wrap(err, errors.CodeDBQueryFailed, "获取渠道配置失败"))
		}
		return
	}

	errors.ResponseSuccess(ctx, configs, "获取渠道配置成功")
}

func UpdateChannelConfigs(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	req, err := common.ValidateRequest[dto.ChannelConfigDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	if err := storage.UpdateChannelConfigs(channelID, *req); err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeDBUpdateFailed, "更新渠道配置失败"))
		return
	}

	if err := storage.RefreshChannelCache(channelID); err != nil {
	}

	errors.ResponseSuccess(ctx, nil, "更新渠道配置成功")
}

func TestConnection(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	if err := storage.TestConnection(channelID); err != nil {
		errors.HandleError(ctx, errors.Wrap(err, errors.CodeInternal, "连接测试失败: "+err.Error()))
		return
	}

	errors.ResponseSuccess(ctx, nil, "连接测试成功")
}

func SetDefaultChannel(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	if err := storage.SetDefaultChannel(channelID); err != nil {
		if _, ok := err.(*errors.Error); ok {
			errors.HandleError(ctx, err)
		} else {
			errors.HandleError(ctx, errors.Wrap(err, errors.CodeDBUpdateFailed, "设置默认渠道失败"))
		}
		return
	}

	errors.ResponseSuccess(ctx, nil, "设置默认渠道成功")
}

func EnableChannel(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	if err := storage.EnableChannel(channelID); err != nil {
		if _, ok := err.(*errors.Error); ok {
			errors.HandleError(ctx, err)
		} else {
			errors.HandleError(ctx, errors.Wrap(err, errors.CodeDBUpdateFailed, "启用渠道失败"))
		}
		return
	}

	errors.ResponseSuccess(ctx, nil, "启用渠道成功")
}

func DisableChannel(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	if err := storage.DisableChannel(channelID); err != nil {
		if _, ok := err.(*errors.Error); ok {
			errors.HandleError(ctx, err)
		} else {
			errors.HandleError(ctx, errors.Wrap(err, errors.CodeDBUpdateFailed, "禁用渠道失败"))
		}
		return
	}

	errors.ResponseSuccess(ctx, nil, "禁用渠道成功")
}

func ExportChannelConfig(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	channel, err := storage.GetChannelByID(channelID)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeNotFound, "存储渠道不存在"))
		return
	}

	if channel.IsLocal || channel.Type == "local" {
		errors.HandleError(ctx, errors.New(errors.CodeValidationFailed, "本地存储渠道不支持导出"))
		return
	}

	data, err := storage.ExportChannelConfig(channelID)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "导出渠道配置失败"))
		return
	}

	errors.ResponseSuccess(ctx, data, "导出渠道配置成功")
}

func ImportChannelConfig(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "无效的文件"))
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/json" && !strings.HasSuffix(strings.ToLower(header.Filename), ".json") {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "文件类型错误，请上传JSON文件"))
		return
	}

	if header.Size > 5*1024*1024 {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "文件大小超过限制（5MB）"))
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInternal, "读取文件失败"))
		return
	}

	if err := storage.ImportChannelConfig(data); err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, nil, "导入渠道配置成功")
}

func ExportAllChannelConfigs(ctx *gin.Context) {
	data, err := storage.ExportAllChannelConfigs()
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, data, "导出所有渠道配置成功")
}

func RefreshChannelCache(ctx *gin.Context) {
	channelID := ctx.Param("id")
	if channelID == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "缺少渠道ID参数"))
		return
	}

	_, err := storage.GetChannelByID(channelID)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeNotFound, "存储渠道不存在"))
		return
	}

	if err := storage.RefreshChannelCache(channelID); err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, nil, "渠道缓存刷新成功")
}

func ClearAllChannelCache(ctx *gin.Context) {
	if err := storage.ClearAllChannelCache(); err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, nil, "所有渠道缓存清空成功")
}
