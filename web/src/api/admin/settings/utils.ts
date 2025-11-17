/**
 * 设置工具函数

 */
import type { TranslationFunction } from '@/composables/useTexts'
import { logger } from '@/utils/system/logger'
import type { Setting, SettingGroup } from './types'
import { getSettings, batchUpsertSettings } from './common'
import { defaultSettings, getDefaultSettings } from './defaults'

export async function initializeSettings(group: SettingGroup, $t?: TranslationFunction): Promise<Setting[]> {
  try {
    const { data } = await getSettings({ group })

    // 选择带翻译的默认项（若提供 $t），否则使用静态回退
    const defaultsMap: Record<SettingGroup, Setting[]> = $t ? getDefaultSettings($t) : defaultSettings

    if (!data.settings || data.settings.length < (defaultsMap[group]?.length || 0)) {
      const defaults = defaultsMap[group] || []
      const existingKeys = data.settings ? data.settings.map((s) => s.key) : []

      const settingsToUpsert = defaults.filter((setting) => !existingKeys.includes(setting.key))
      if (settingsToUpsert.length > 0) {
        await batchUpsertSettings(settingsToUpsert)
      }

      const { data: refreshedData } = await getSettings({ group })
      return refreshedData.settings
    }

    return data.settings
  } catch (error) {
    if (process.env.NODE_ENV !== 'production') {
      const errorMsg = $t ? $t('api.settings.errors.loadFailed', { group }) : `Failed to load ${group} settings:`
      logger.error(errorMsg, error)
    }
    return []
  }
}
