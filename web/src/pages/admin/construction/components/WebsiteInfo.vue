<script setup lang="ts">
  import { onMounted, reactive, watch } from 'vue'
  import { defaultSettings, type Setting } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  /* 从defaultSettings获取construction设置的默认值 */
  const constructionDefaults = defaultSettings.construction.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  /* 本地设置对象 */
  const localSettings = reactive({
    site_name: constructionDefaults.site_name || 'PixelPunk',
  })

  /* 将扁平化的设置转换为Setting数组 */
  const getSettingsArray = (): Setting[] => {
    const relevantKeys = ['site_name', 'site_mode']
    return defaultSettings.construction
      .filter((defaultSetting) => relevantKeys.includes(defaultSetting.key))
      .map((defaultSetting) => {
        const settingKey = defaultSetting.key as keyof typeof localSettings
        return {
          key: defaultSetting.key,
          value: localSettings[settingKey],
          type: defaultSetting.type,
          group: 'construction',
          description: defaultSetting.description,
        }
      })
  }

  const applySettings = (settings: Setting[]) => {
    settings.forEach((setting) => {
      const key = setting.key as keyof typeof localSettings
      if (key in localSettings) {
        const settingsRecord = localSettings as Record<string, unknown>
        settingsRecord[key] = setting.value
      }
    })
  }

  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        applySettings(newSettings)
      }
    },
    { deep: true, immediate: true }
  )

  watch(
    localSettings,
    () => {
      emit('update', getSettingsArray())
    },
    { deep: true }
  )

  onMounted(() => {
    if (props.settings && props.settings.length > 0) {
      applySettings(props.settings)
    }
  })
</script>

<template>
  <div class="website-info">
    <div class="rounded-xl border border-subtle bg-background-600 p-6">
      <div class="space-y-6">
        <div class="form-group">
          <label class="form-label">{{ $t('admin.construction.websiteInfo.siteName') }}</label>
          <input
            v-model="localSettings.site_name"
            type="text"
            class="form-input"
            :placeholder="$t('admin.construction.websiteInfo.siteNamePlaceholder')"
            maxlength="50"
          />
          <div class="form-help">{{ $t('admin.construction.websiteInfo.siteNameHint') }}</div>
        </div>

        <div v-if="localSettings.site_name" class="form-group">
          <label class="form-label">{{ $t('admin.construction.websiteInfo.preview') }}</label>
          <div class="site-name-preview">
            <div class="preview-title">
              <i class="fas fa-globe mr-2" />
              <span class="site-title">{{ localSettings.site_name }}</span>
            </div>
            <div class="preview-subtitle">{{ $t('admin.construction.websiteInfo.previewSubtitle') }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .website-info {
    margin-bottom: 2rem;
  }

  .site-name-preview {
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 1.5rem;
    background: rgba(var(--color-background-800-rgb), 0.3);
    margin-top: 0.5rem;
    text-align: center;
  }

  .preview-title {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 0.5rem;
    color: var(--color-brand-500);
    font-weight: 600;
  }

  .site-title {
    font-size: 1.5rem;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .preview-subtitle {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.9rem;
    font-style: italic;
  }

  .mode-preview {
    margin-top: 0.5rem;
  }

  .mode-card {
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 1.5rem;
    background: rgba(var(--color-background-800-rgb), 0.4);
    transition: all 0.3s ease;
  }

  .mode-card.active {
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    background: rgba(var(--color-background-800-rgb), 0.6);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .mode-header {
    display: flex;
    align-items: center;
    margin-bottom: 0.75rem;
  }

  .mode-title {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-brand-500);
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .mode-description {
    color: rgba(var(--color-white-rgb), 0.8);
    font-size: 0.9rem;
    line-height: 1.5;
    margin-bottom: 1rem;
  }

  .mode-features {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .feature-tag {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.8rem;
    font-weight: 500;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }
</style>
