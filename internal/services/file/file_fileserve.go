package file

/* File serving helpers (no behavior change). */

import (
	"io"
	"os"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage"
	pathutil "pixelpunk/pkg/storage/path"
	"strings"
)

/* GetFileLocalPath 从上下文中获取文件信息并返回本地路径 */
func GetFileLocalPath(file models.File, isThumb bool) (string, error) {
	var localPath string
	if isThumb {
		localPath = file.LocalThumbPath
		if localPath == "" {
			return "", errors.New(errors.CodeFileNotFound, "缩略图文件不存在")
		}
	} else {
		localPath = file.LocalFilePath
		if localPath == "" {
			return "", errors.New(errors.CodeFileNotFound, "文件不存在")
		}
	}
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		if isThumb {
			return "", errors.New(errors.CodeFileNotFound, "缩略图文件不存在")
		}
		return "", errors.New(errors.CodeFileNotFound, "文件不存在")
	}
	return localPath, nil
}

/* ServeFile 根据文件存储类型获取访问信息 */
func ServeFile(file models.File, isThumb bool) (interface{}, bool, bool, error) {
	provider, err := storage.GetStorageProviderByChannelID(file.StorageProviderID)
	if err != nil {
		return nil, false, false, err
	}
	if provider.IsDirectAccess() {
		localPath := file.LocalFilePath
		if isThumb {
			localPath = file.LocalThumbPath
		}
		return localPath, true, false, nil
	}

	globalSettings, err := setting.GetSettingsByGroupAsMap("global")
	var globalHideRemoteURL bool
	if err == nil {
		if val, exists := globalSettings.Settings["hide_remote_url"]; exists {
			if boolVal, ok := val.(bool); ok {
				globalHideRemoteURL = boolVal
			}
		}
	}
	useProxy := false
	if channelConfigMap, err := storage.GetChannelConfigMapFromService(file.StorageProviderID); err == nil {
		isPrivateAccess := false
		if val, exists := channelConfigMap["access_control"]; exists {
			if v, ok := val.(string); ok {
				isPrivateAccess = (v == "private")
			}
		}
		if isPrivateAccess {
			useProxy = true
		} else {
			var channelHideRemoteURL bool
			var channelHasHideRemoteURLSetting bool
			if val, exists := channelConfigMap["hide_remote_url"]; exists {
				channelHasHideRemoteURLSetting = true
				switch v := val.(type) {
				case bool:
					channelHideRemoteURL = v
				case string:
					channelHideRemoteURL = (v == "true")
				}
			}
			if globalHideRemoteURL {
				useProxy = true
			} else if channelHasHideRemoteURLSetting {
				useProxy = channelHideRemoteURL
			} else {
				useProxy = false
			}
		}
	} else {
		useProxy = globalHideRemoteURL
	}

	var candidate string
	if isThumb {
		if file.RemoteThumbURL != "" && !pathutil.IsHTTPURL(file.RemoteThumbURL) {
			candidate = file.RemoteThumbURL
		} else {
			candidate = file.ThumbURL
		}
	} else {
		if file.RemoteURL != "" && !pathutil.IsHTTPURL(file.RemoteURL) {
			candidate = file.RemoteURL
		} else {
			candidate = file.URL
		}
	}
	candidate = strings.TrimPrefix(candidate, "/")
	remoteUrl := candidate

	if useProxy {
		content, contentType, err := provider.GetRemoteContent(remoteUrl, isThumb, file.UserID)
		if err != nil {
			logger.Error("代理模式获取内容失败: %v", err)
			return nil, false, false, err
		}
		return &ProxyResponse{Content: content, ContentType: contentType, ContentLength: 0}, false, true, nil
	}
	fileURL, err := provider.GetFileURL(remoteUrl, isThumb)
	if err != nil {
		return nil, false, false, err
	}
	return fileURL, false, false, nil
}

/* ProxyResponse 代理响应 */
type ProxyResponse struct {
	Content       io.ReadCloser
	ContentType   string
	ContentLength int64
}
