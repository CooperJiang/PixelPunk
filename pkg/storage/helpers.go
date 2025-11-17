package storage

import (
	"io"
	"net/url"
	"os"
	"strings"

	"pixelpunk/internal/models"
	setting "pixelpunk/internal/services/setting"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/storage/adapter"
	middleware "pixelpunk/pkg/storage/middleware"
	pathutil "pixelpunk/pkg/storage/path"
	urlstrategy "pixelpunk/pkg/storage/url"
	su "pixelpunk/pkg/storage/utils"
	"pixelpunk/pkg/utils"
)

// GetFullURLs returns full URLs (original, thumbnail, shortLink) for an file
func GetFullURLs(file models.File) (string, string, string) {
	st := NewGlobalStorage()

	var globalHideRemote bool
	if sec, err := setting.GetSettingsByGroupAsMap("security"); err == nil && sec != nil {
		if v, ok := sec.Settings["hide_remote_url"].(bool); ok {
			globalHideRemote = v
		}
	} else if glb, err2 := setting.GetSettingsByGroupAsMap("global"); err2 == nil && glb != nil {
		if v, ok := glb.Settings["hide_remote_url"].(bool); ok {
			globalHideRemote = v
		}
	}

	channelConfig, _ := GetChannelConfigMapFromService(file.StorageProviderID)
	getBool := func(key string, def bool) bool {
		v, ok := channelConfig[key]
		if !ok {
			return def
		}
		switch t := v.(type) {
		case bool:
			return t
		case string:
			return strings.TrimSpace(strings.ToLower(t)) == "true"
		default:
			return def
		}
	}

	// 判断是否是本地存储（需要先判断，后面要用）
	isLocal := false
	if file.StorageProviderID != "" && database.GetDB() != nil {
		if mgr := st.GetManager(); mgr != nil {
			if chType, err := mgr.GetChannelType(file.StorageProviderID); err == nil {
				isLocal = (chType == "local")
			}
		}
	}
	if !isLocal && strings.TrimSpace(strings.ToLower(file.StorageType)) == "local" {
		isLocal = true
	}

	rawCustomDomain := func() string {
		if v, ok := channelConfig["custom_domain"].(string); ok {
			return strings.TrimSpace(v)
		}
		return ""
	}()

	sanitizeDomain := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.TrimSuffix(s, "/")
		s = strings.TrimPrefix(s, "http://")
		s = strings.TrimPrefix(s, "https://")
		return s
	}

	extractProtocol := func(s string) string {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(strings.ToLower(s), "https://") {
			return "https"
		}
		if strings.HasPrefix(strings.ToLower(s), "http://") {
			return "http"
		}
		return "https" // 默认 HTTPS
	}

	var useHTTPS bool
	var customDomain string

	// 如果配置了自定义域名，优先使用自定义域名的协议
	if rawCustomDomain != "" {
		customDomain = sanitizeDomain(rawCustomDomain)
		useHTTPS = extractProtocol(rawCustomDomain) == "https"
	} else if isLocal {
		// 本地渠道且没有自定义域名，从 site_base_url 获取协议和域名
		if website, err := setting.GetSettingsByGroupAsMap("website"); err == nil && website != nil {
			if siteURL, ok := website.Settings["site_base_url"].(string); ok && siteURL != "" {
				customDomain = sanitizeDomain(siteURL)
				useHTTPS = extractProtocol(siteURL) == "https"
			}
		}
	} else {
		// 第三方渠道，使用 use_https 配置（默认 true）
		useHTTPS = getBool("use_https", true)
		// 如果没有自定义域名，尝试从 site_base_url 获取
		if website, err := setting.GetSettingsByGroupAsMap("website"); err == nil && website != nil {
			if vv, ok2 := website.Settings["site_base_url"].(string); ok2 && vv != "" {
				customDomain = sanitizeDomain(vv)
			}
		}
	}

	hideRemoteURL := false
	if isLocal {
		hideRemoteURL = true
	} else {
		isPrivateAccess := false
		if val, exists := channelConfig["access_control"]; exists {
			if v, ok := val.(string); ok {
				isPrivateAccess = (v == "private")
			}
		}

		if isPrivateAccess {
			hideRemoteURL = true
		} else {
			var channelHideRemoteURL bool
			var channelHasHideRemoteURLSetting bool
			if val, exists := channelConfig["hide_remote_url"]; exists {
				channelHasHideRemoteURLSetting = true
				switch v := val.(type) {
				case bool:
					channelHideRemoteURL = v
				case string:
					channelHideRemoteURL = (v == "true")
				}
			}

			if channelHasHideRemoteURLSetting {
				hideRemoteURL = channelHideRemoteURL
			} else if globalHideRemote {
				hideRemoteURL = true
			} else {
				hideRemoteURL = false
			}
		}
	}

	var fullURL, fullThumbURL string

	if hideRemoteURL {
		displayName := getDisplayNameWithExtension(file)

		var baseURL string
		if customDomain != "" {
			if useHTTPS {
				baseURL = "https://" + customDomain
			} else {
				baseURL = "http://" + customDomain
			}
		} else {
			if website, err := setting.GetSettingsByGroupAsMap("website"); err == nil && website != nil {
				if siteURL, ok := website.Settings["site_base_url"].(string); ok && siteURL != "" {
					baseURL = strings.TrimSuffix(siteURL, "/")
				}
			}
		}

		if baseURL != "" {
			fullURL = baseURL + "/f/" + file.ID + "/" + displayName
			fullThumbURL = baseURL + "/t/" + file.ID + "/" + displayName
		} else {
			fullURL = "/f/" + file.ID + "/" + displayName
			fullThumbURL = "/t/" + file.ID + "/" + displayName
		}
	} else {
		if file.RemoteURL == "" {
			displayName := getDisplayNameWithExtension(file)
			fullURL = utils.GenerateFullURL("/f/"+file.ID+"/"+displayName, "local")
		} else if pathutil.IsHTTPURL(file.RemoteURL) {
			if customDomain != "" {
				urlstrategy := urlstrategy.NewURLStrategy(urlstrategy.StrategyConfig{
					ForceHTTPS:   useHTTPS,
					CustomDomain: customDomain,
				})
				fullURL = urlstrategy.BuildDirectURL(file.RemoteURL)
			} else {
				fullURL = file.RemoteURL
				if useHTTPS && strings.HasPrefix(fullURL, "http://") {
					fullURL = strings.Replace(fullURL, "http://", "https://", 1)
				} else if !useHTTPS && strings.HasPrefix(fullURL, "https://") {
					fullURL = strings.Replace(fullURL, "https://", "http://", 1)
				}
			}
		} else {
			if mgr := st.GetManager(); mgr != nil {
				adp, err := mgr.GetAdapter(file.StorageProviderID)
				if err == nil && adp != nil {
					opts := &adapter.URLOptions{CustomDomain: customDomain}
					if generatedURL, err := adp.GetURL(file.RemoteURL, opts); err == nil {
						fullURL = generatedURL
					}
				}
			}
			if fullURL == "" {
				displayName := getDisplayNameWithExtension(file)
				fullURL = utils.GenerateFullURL("/f/"+file.ID+"/"+displayName, "local")
			}
		}

		if file.ThumbURL != "" {
			if file.RemoteThumbURL == "" {
				displayName := getDisplayNameWithExtension(file)
				fullThumbURL = utils.GenerateFullURL("/t/"+file.ID+"/"+displayName, "local")
			} else if pathutil.IsHTTPURL(file.RemoteThumbURL) {
				if customDomain != "" {
					urlstrategy := urlstrategy.NewURLStrategy(urlstrategy.StrategyConfig{
						ForceHTTPS:   useHTTPS,
						CustomDomain: customDomain,
					})
					fullThumbURL = urlstrategy.BuildDirectURL(file.RemoteThumbURL)
				} else {
					fullThumbURL = file.RemoteThumbURL
					if useHTTPS && strings.HasPrefix(fullThumbURL, "http://") {
						fullThumbURL = strings.Replace(fullThumbURL, "http://", "https://", 1)
					} else if !useHTTPS && strings.HasPrefix(fullThumbURL, "https://") {
						fullThumbURL = strings.Replace(fullThumbURL, "https://", "http://", 1)
					}
				}
			} else {
				if mgr := st.GetManager(); mgr != nil {
					adp, err := mgr.GetAdapter(file.StorageProviderID)
					if err == nil && adp != nil {
						opts := &adapter.URLOptions{CustomDomain: customDomain, IsThumbnail: true}
						if generatedURL, err := adp.GetURL(file.RemoteThumbURL, opts); err == nil {
							fullThumbURL = generatedURL
						}
					}
				}
				if fullThumbURL == "" {
					displayName := getDisplayNameWithExtension(file)
					fullThumbURL = utils.GenerateFullURL("/t/"+file.ID+"/"+displayName, "local")
				}
			}
		}
	}

	var shortLinkURL string
	if file.ShortURL != "" {
		shortLinkURL = utils.GenerateFullURL("/s/"+file.ShortURL, "local")
	} else {
		shortURL := file.GenerateShortURL()
		shortLinkURL = utils.GenerateFullURL("/s/"+shortURL, "local")
	}

	return fullURL, fullThumbURL, shortLinkURL
}

func GetSafeFolderName(folderID string, _ interface{}, _ uint) string {
	return SanitizeFileName(folderID)
}


func EnsureDirExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

func FileExists(path string) bool { return su.FileExists(path) }

func ExtractEXIFData(r io.Reader) (*middleware.EXIFData, error) { return middleware.ExtractEXIFData(r) }

func GetChannelConfigMapFromService(channelID string) (map[string]interface{}, error) {
	db := database.GetDB()
	if db == nil {
		return map[string]interface{}{}, nil
	}
	var items []models.StorageConfigItem
	if err := db.Where("channel_id = ?", channelID).Find(&items).Error; err != nil {
		return nil, err
	}
	m := make(map[string]interface{}, len(items))
	for _, it := range items {
		m[it.KeyName] = it.Value
	}
	return m, nil
}

func getDisplayNameWithExtension(file models.File) string {
	var displayName string
	if file.OriginalName != "" {
		displayName = file.OriginalName
	} else if file.DisplayName != "" {
		displayName = file.DisplayName
		if !strings.Contains(displayName, ".") && file.Format != "" {
			displayName += "." + strings.ToLower(file.Format)
		}
	} else {
		if file.Format != "" {
			displayName = "file." + strings.ToLower(file.Format)
		} else {
			displayName = "file"
		}
	}

	displayName = strings.TrimPrefix(displayName, "thumb-")
	return url.PathEscape(displayName)
}
