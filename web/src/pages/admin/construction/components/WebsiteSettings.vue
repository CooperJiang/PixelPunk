<script setup lang="ts">
  import { onMounted, reactive, watch } from 'vue'
  import type { Setting, SettingType } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  /* 本地设置对象（扁平化）- 初始化为空，完全依赖接口数据 */
  const localSettings = reactive({
    site_name: '',
    site_description: '',
    site_keywords: '',
    icp_number: '',
    show_file_count: false,
    show_storage_usage: false,
    site_logo_url: '',
    favicon_url: '',
    copyright_text: '',
    contact_email: '',
    footer_custom_text: '',
    site_hero_title: '',
    site_features_text: '',
    site_slogan: '',
  })

  /* 将扁平化的设置转换为Setting数组 */
  const getSettingsArray = (): Setting[] => {
    const settingDefinitions = [
      { key: 'site_name', type: 'string', description: $t('admin.construction.websiteSettings.desc.siteName') },
      { key: 'site_description', type: 'string', description: $t('admin.construction.websiteSettings.desc.siteDescription') },
      { key: 'site_hero_title', type: 'string', description: $t('admin.construction.websiteSettings.desc.siteHeroTitle') },
      { key: 'site_features_text', type: 'string', description: $t('admin.construction.websiteSettings.desc.siteFeaturesText') },
      { key: 'site_logo_url', type: 'string', description: $t('admin.construction.websiteSettings.desc.siteLogoUrl') },
      { key: 'favicon_url', type: 'string', description: $t('admin.construction.websiteSettings.desc.faviconUrl') },
      { key: 'show_file_count', type: 'boolean', description: $t('admin.construction.websiteSettings.desc.showFileCount') },
      { key: 'show_storage_usage', type: 'boolean', description: $t('admin.construction.websiteSettings.desc.showStorageUsage') },
      { key: 'site_keywords', type: 'string', description: $t('admin.construction.websiteSettings.desc.siteKeywords') },
      { key: 'icp_number', type: 'string', description: $t('admin.construction.websiteSettings.desc.icpNumber') },
      { key: 'contact_email', type: 'string', description: $t('admin.construction.websiteSettings.desc.contactEmail') },
      { key: 'copyright_text', type: 'string', description: $t('admin.construction.websiteSettings.desc.copyrightText') },
      { key: 'footer_custom_text', type: 'string', description: $t('admin.construction.websiteSettings.desc.footerCustomText') },
    ]

    return settingDefinitions.map((def) => {
      const settingKey = def.key as keyof typeof localSettings
      return {
        key: def.key,
        value: localSettings[settingKey],
        type: def.type as SettingType,
        group: 'website_info',
        description: def.description,
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
  <div class="space-y-5">
    <div class="space-y-5">
      <div class="space-y-5">
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-tag mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.siteName') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.site_name"
              :placeholder="$t('admin.construction.websiteInfo.siteNamePlaceholder')"
              :maxlength="50"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.siteNameHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-align-left mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.siteDescription') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.site_description"
              :placeholder="$t('admin.construction.websiteInfo.siteDescriptionPlaceholder')"
              :maxlength="200"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.siteDescriptionHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-bullhorn mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.heroTitle') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.site_hero_title"
              :placeholder="$t('admin.construction.websiteInfo.heroTitlePlaceholder')"
              :maxlength="100"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.heroTitleHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-list-ul mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.featuresText') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.site_features_text"
              :placeholder="$t('admin.construction.websiteInfo.featuresTextPlaceholder')"
              :maxlength="200"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.featuresTextHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-image mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.siteLogo') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.site_logo_url"
              :placeholder="$t('admin.construction.websiteInfo.siteLogoPlaceholder')"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.siteLogoHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-star mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.favicon') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.favicon_url"
              :placeholder="$t('admin.construction.websiteInfo.faviconPlaceholder')"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.faviconHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-chart-bar mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.showFileCount') }}
          </label>
          <div class="flex-1">
            <CyberSwitch v-model="localSettings.show_file_count" />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.showFileCountHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-hdd mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.showStorageUsage') }}
          </label>
          <div class="flex-1">
            <CyberSwitch v-model="localSettings.show_storage_usage" />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.showStorageUsageHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-key mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.keywords') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.site_keywords"
              :placeholder="$t('admin.construction.websiteInfo.keywordsPlaceholder')"
              :maxlength="200"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.keywordsHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-certificate mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.icpNumber') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.icp_number"
              :placeholder="$t('admin.construction.websiteInfo.icpNumberPlaceholder')"
              :maxlength="50"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.icpNumberHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-envelope mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.contactEmail') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.contact_email"
              :placeholder="$t('admin.construction.websiteInfo.contactEmailPlaceholder')"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.contactEmailHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-copyright mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.copyright') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.copyright_text"
              :placeholder="$t('admin.construction.websiteInfo.copyrightPlaceholder')"
              :maxlength="100"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.copyrightHint') }}</div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-info-circle mr-2 text-brand-600" />{{ $t('admin.construction.websiteInfo.footerCustom') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.footer_custom_text"
              :placeholder="$t('admin.construction.websiteInfo.footerCustomPlaceholder')"
              :maxlength="200"
              width="500px"
            />
            <div class="mt-1 text-xs text-content-muted">{{ $t('admin.construction.websiteInfo.footerCustomHint') }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .hover:shadow-glow:hover {
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }
</style>
