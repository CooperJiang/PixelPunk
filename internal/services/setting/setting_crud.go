package setting

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/hooks"
	"strconv"
	"time"
)

/* GetSettings 获取设置列表 */
func GetSettings(query *dto.SettingQueryDTO) (*dto.SettingListResponseDTO, error) {
	db := database.GetDB()
	var settings []models.Setting

	dbQuery := db.Model(&models.Setting{})

	if query != nil && query.Group != "" {
		dbQuery = dbQuery.Where("`group` = ?", query.Group)
	}

	if query != nil && query.Key != "" {
		dbQuery = dbQuery.Where("`key` = ?", query.Key)
	}

	if err := dbQuery.Order("`group` ASC, `key` ASC").Find(&settings).Error; err != nil {
		return &dto.SettingListResponseDTO{Settings: []dto.SettingResponseDTO{}}, nil
	}

	result := &dto.SettingListResponseDTO{Settings: make([]dto.SettingResponseDTO, 0, len(settings))}

	for _, s := range settings {
		value := parseSettingValue(s)
		result.Settings = append(result.Settings, dto.SettingResponseDTO{
			ID:          s.ID,
			Key:         s.Key,
			Value:       value,
			Type:        s.Type,
			Group:       s.Group,
			Description: s.Description,
			IsSystem:    s.IsSystem,
			CreatedAt:   time.Time(s.CreatedAt).Format(time.RFC3339),
			UpdatedAt:   time.Time(s.UpdatedAt).Format(time.RFC3339),
		})
	}

	return result, nil
}

/* GetSetting 获取单个设置 */
func GetSetting(key string) (*dto.SettingResponseDTO, error) {
	if cachedSetting, found := getSettingFromCache(key); found {
		return cachedSetting, nil
	}

	db := database.GetDB()
	var setting models.Setting

	if err := db.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBNoRecord, fmt.Sprintf("设置 %s 不存在", key))
	}

	value := parseSettingValue(setting)

	result := &dto.SettingResponseDTO{
		ID:          setting.ID,
		Key:         setting.Key,
		Value:       value,
		Type:        setting.Type,
		Group:       setting.Group,
		Description: setting.Description,
		IsSystem:    setting.IsSystem,
		CreatedAt:   time.Time(setting.CreatedAt).Format(time.RFC3339),
		UpdatedAt:   time.Time(setting.UpdatedAt).Format(time.RFC3339),
	}

	setSettingToCache(key, result)

	return result, nil
}

/* CreateSetting 创建设置 */
func CreateSetting(createDTO *dto.SettingCreateDTO) (*dto.SettingResponseDTO, error) {
	db := database.GetDB()

	var count int64
	if err := db.Model(&models.Setting{}).Where("`key` = ?", createDTO.Key).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "检查设置键名是否存在失败")
	}
	if count > 0 {
		return nil, errors.New(errors.CodeDBDuplicate, "设置键名已存在")
	}

	var valueToStore string
	var valueJSON []byte
	var jsonErr error

	if createDTO.Type == models.SettingTypeString {
		if strValue, ok := createDTO.Value.(string); ok {
			valueJSON, jsonErr = json.Marshal(strValue)
		} else {
			valueJSON, jsonErr = json.Marshal(fmt.Sprintf("%v", createDTO.Value))
		}
	} else {
		valueJSON, jsonErr = json.Marshal(createDTO.Value)
	}
	if jsonErr != nil {
		return nil, errors.Wrap(jsonErr, errors.CodeInvalidParameter, "设置值格式错误")
	}
	valueToStore = string(valueJSON)

	setting := models.Setting{
		Key:         createDTO.Key,
		Value:       valueToStore,
		Type:        createDTO.Type,
		Group:       createDTO.Group,
		Description: createDTO.Description,
		IsSystem:    createDTO.IsSystem,
	}

	if err := db.Create(&setting).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBCreateFailed, "创建设置失败")
	}

	result := &dto.SettingResponseDTO{
		ID:          setting.ID,
		Key:         setting.Key,
		Value:       createDTO.Value,
		Type:        setting.Type,
		Group:       setting.Group,
		Description: setting.Description,
		IsSystem:    setting.IsSystem,
		CreatedAt:   time.Time(setting.CreatedAt).Format(time.RFC3339),
		UpdatedAt:   time.Time(setting.UpdatedAt).Format(time.RFC3339),
	}
	setSettingToCache(setting.Key, result)
	deleteSettingGroupFromCache(setting.Group)
	return result, nil
}

/* UpdateSetting 更新设置 */
func UpdateSetting(updateDTO *dto.SettingUpdateDTO) (*dto.SettingResponseDTO, error) {
	db := database.GetDB()

	var setting models.Setting
	if err := db.Where("`key` = ?", updateDTO.Key).First(&setting).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBNoRecord, fmt.Sprintf("设置 %s 不存在", updateDTO.Key))
	}

	if setting.IsSystem && setting.Group != updateDTO.Group {
		return nil, errors.New(errors.CodeForbidden, "系统设置不能修改分组")
	}

	var valueToStore string
	var valueJSON []byte
	var jsonErr error
	if updateDTO.Type == models.SettingTypeString {
		if strValue, ok := updateDTO.Value.(string); ok {
			valueJSON, jsonErr = json.Marshal(strValue)
		} else {
			valueJSON, jsonErr = json.Marshal(fmt.Sprintf("%v", updateDTO.Value))
		}
	} else {
		valueJSON, jsonErr = json.Marshal(updateDTO.Value)
	}
	if jsonErr != nil {
		return nil, errors.Wrap(jsonErr, errors.CodeInvalidParameter, "设置值格式错误")
	}
	valueToStore = string(valueJSON)

	setting.Value = valueToStore
	setting.Type = updateDTO.Type
	setting.Group = updateDTO.Group
	setting.Description = updateDTO.Description

	if err := db.Save(&setting).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新设置失败")
	}

	result := &dto.SettingResponseDTO{
		ID:          setting.ID,
		Key:         setting.Key,
		Value:       updateDTO.Value,
		Type:        setting.Type,
		Group:       setting.Group,
		Description: setting.Description,
		IsSystem:    setting.IsSystem,
		CreatedAt:   time.Time(setting.CreatedAt).Format(time.RFC3339),
		UpdatedAt:   time.Time(setting.UpdatedAt).Format(time.RFC3339),
	}

	invalidateSettingCaches(setting.Group, setting.Key)

	var valueStr string
	switch v := updateDTO.Value.(type) {
	case string:
		valueStr = v
	default:
		if data, err := json.Marshal(updateDTO.Value); err == nil {
			valueStr = string(data)
		}
	}
	notifySettingChanged(setting.Group, setting.Key, valueStr)

	hooks.TriggerSettingUpdate(setting.Group)

	return result, nil
}

/* DeleteSetting 删除设置 */
func DeleteSetting(key string) error {
	db := database.GetDB()
	var setting models.Setting
	if err := db.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBNoRecord, fmt.Sprintf("设置 %s 不存在", key))
	}
	if setting.IsSystem {
		return errors.New(errors.CodeForbidden, "系统设置不能删除")
	}
	if err := db.Delete(&setting).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除设置失败")
	}
	invalidateSettingCaches(setting.Group, setting.Key)
	return nil
}

/* GetSettingValue 获取设置值（直接返回对应类型的值） */
func GetSettingValue(key string, defaultValue interface{}) (interface{}, error) {
	db := database.GetDB()
	var setting models.Setting
	if err := db.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return defaultValue, nil
	}
	return parseSettingValue(setting), nil
}

/* GetStringValue 获取字符串设置值 */
func GetStringValue(key string, defaultValue string) (string, error) {
	val, err := GetSettingValue(key, defaultValue)
	if err != nil {
		return defaultValue, err
	}
	if strVal, ok := val.(string); ok {
		return strVal, nil
	}
	return defaultValue, nil
}

/* GetNumberValue 获取数字设置值 */
func GetNumberValue(key string, defaultValue float64) (float64, error) {
	val, err := GetSettingValue(key, defaultValue)
	if err != nil {
		return defaultValue, err
	}
	if numVal, ok := val.(float64); ok {
		return numVal, nil
	}
	return defaultValue, nil
}

/* GetBoolValue 获取布尔设置值 */
func GetBoolValue(key string, defaultValue bool) (bool, error) {
	val, err := GetSettingValue(key, defaultValue)
	if err != nil {
		return defaultValue, err
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal, nil
	}
	return defaultValue, nil
}

func parseSettingValue(setting models.Setting) interface{} {
	var valueInterface interface{}
	if err := json.Unmarshal([]byte(setting.Value), &valueInterface); err != nil {
		return setting.Value
	}

	switch setting.Type {
	case models.SettingTypeString, "text": // 支持text类型(用于HTML内容)
		if strValue, ok := valueInterface.(string); ok {
			return strValue
		}
		return fmt.Sprintf("%v", valueInterface)
	case models.SettingTypeNumber:
		switch v := valueInterface.(type) {
		case float64:
			return v
		case json.Number:
			if f, err := v.Float64(); err == nil {
				return f
			}
		}
		if strValue, ok := valueInterface.(string); ok {
			if f, err := strconv.ParseFloat(strValue, 64); err == nil {
				return f
			}
		}
		return 0
	case models.SettingTypeBoolean:
		if boolValue, ok := valueInterface.(bool); ok {
			return boolValue
		}
		if strValue, ok := valueInterface.(string); ok {
			if b, err := strconv.ParseBool(strValue); err == nil {
				return b
			}
		}
		return false
	case models.SettingTypeJSON:
		return valueInterface
	case models.SettingTypeArray:
		if arrayValue, ok := valueInterface.([]interface{}); ok {
			return arrayValue
		}
		return []interface{}{}
	default:
		return valueInterface
	}
}
