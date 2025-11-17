package setting

import (
	"encoding/json"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/pkg/cache"
)

func getSettingCacheKey(key string) string {
	return SettingCachePrefix + key
}

func getSettingGroupCacheKey(group string) string {
	return SettingGroupPrefix + group
}

func getSettingFromCache(key string) (*dto.SettingResponseDTO, bool) {
	cacheKey := getSettingCacheKey(key)
	cachedData, err := cache.Get(cacheKey)
	if err != nil {
		return nil, false
	}

	var setting dto.SettingResponseDTO
	if err := json.Unmarshal([]byte(cachedData), &setting); err != nil {
		return nil, false
	}

	return &setting, true
}

func setSettingToCache(key string, setting *dto.SettingResponseDTO) {
	cacheKey := getSettingCacheKey(key)
	data, err := json.Marshal(setting)
	if err != nil {
		return
	}
	_ = cache.Set(cacheKey, string(data), SettingCacheExpire)
}

func getSettingGroupFromCache(group string) (*dto.SettingMapResponseDTO, bool) {
	cacheKey := getSettingGroupCacheKey(group)
	cachedData, err := cache.Get(cacheKey)
	if err != nil {
		return nil, false
	}

	var settingGroup dto.SettingMapResponseDTO
	if err := json.Unmarshal([]byte(cachedData), &settingGroup); err != nil {
		return nil, false
	}

	return &settingGroup, true
}

func setSettingGroupToCache(group string, settingGroup *dto.SettingMapResponseDTO) {
	cacheKey := getSettingGroupCacheKey(group)
	data, err := json.Marshal(settingGroup)
	if err != nil {
		return
	}
	_ = cache.Set(cacheKey, string(data), SettingGroupCacheExpire)
}

func deleteSettingFromCache(key string) {
	cacheKey := getSettingCacheKey(key)
	_ = cache.Del(cacheKey)
}

func deleteSettingGroupFromCache(group string) {
	cacheKey := getSettingGroupCacheKey(group)
	_ = cache.Del(cacheKey)
}

func invalidateSettingCaches(group, key string) {
	if key != "" {
		deleteSettingFromCache(key)
	}

	if group != "" {
		deleteSettingGroupFromCache(group)
	}
}
