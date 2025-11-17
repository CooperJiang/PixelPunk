<script setup lang="ts">
  import { computed, onMounted, reactive, watch } from 'vue'
  import type { Setting, SettingType } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const { $t } = useTexts()

  /* 本地设置对象（扁平化）*/
  const localSettings = reactive({
    show_official_site: true,
    official_site_url: 'https://pixelpunk.cc/',
    show_github_link: true,
    github_url: 'https://github.com/CooperJiang/PixelPunk',
    show_wechat_group: true,
    wechat_qr_image_url: '',
    wechat_contact_account: '',
    show_qq_group: true,
    qq_qr_image_url: '',
    qq_group_number: '',
    enable_multi_layout: true,
    default_layout: 'top',

    enable_multi_language: false,
    default_language: 'zh-CN',
  })

  /* 布局选项 */
  const layoutOptions = computed(() => [
    { label: $t('admin.construction.appearance.layoutOptions.top'), value: 'top' },
    { label: $t('admin.construction.appearance.layoutOptions.sidebar'), value: 'sidebar' },
  ])

  /* 语言选项 */
  const languageOptions = computed(() => [
    { label: $t('admin.construction.appearance.languageOptions.auto'), value: 'auto' },
    { label: $t('admin.construction.appearance.languageOptions.zhCN'), value: 'zh-CN' },
    { label: $t('admin.construction.appearance.languageOptions.en'), value: 'en-US' },
    { label: $t('admin.construction.appearance.languageOptions.ja'), value: 'ja-JP' },
  ])

  /* 将扁平化的设置转换为Setting数组 */
  const getSettingsArray = (): Setting[] => {
    const settingDefinitions = [
      { key: 'show_official_site', type: 'boolean', description: $t('admin.construction.appearance.desc.showOfficialSite') },
      { key: 'official_site_url', type: 'string', description: $t('admin.construction.appearance.desc.officialSiteUrl') },
      { key: 'show_github_link', type: 'boolean', description: $t('admin.construction.appearance.desc.showGithubLink') },
      { key: 'github_url', type: 'string', description: $t('admin.construction.appearance.desc.githubUrl') },
      { key: 'show_wechat_group', type: 'boolean', description: $t('admin.construction.appearance.desc.showWechatGroup') },
      { key: 'wechat_qr_image_url', type: 'string', description: $t('admin.construction.appearance.desc.wechatQrImageUrl') },
      {
        key: 'wechat_contact_account',
        type: 'string',
        description: $t('admin.construction.appearance.desc.wechatContactAccount'),
      },
      { key: 'show_qq_group', type: 'boolean', description: $t('admin.construction.appearance.desc.showQqGroup') },
      { key: 'qq_qr_image_url', type: 'string', description: $t('admin.construction.appearance.desc.qqQrImageUrl') },
      { key: 'qq_group_number', type: 'string', description: $t('admin.construction.appearance.desc.qqGroupNumber') },
      { key: 'enable_multi_layout', type: 'boolean', description: $t('admin.construction.appearance.desc.enableMultiLayout') },
      { key: 'default_layout', type: 'string', description: $t('admin.construction.appearance.desc.defaultLayout') },

      {
        key: 'enable_multi_language',
        type: 'boolean',
        description: $t('admin.construction.appearance.desc.enableMultiLanguage'),
      },
      { key: 'default_language', type: 'string', description: $t('admin.construction.appearance.desc.defaultLanguage') },
    ]

    return settingDefinitions.map((def) => {
      const settingKey = def.key as keyof typeof localSettings
      return {
        key: def.key,
        value: localSettings[settingKey],
        type: def.type as SettingType,
        group: 'appearance',
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
  <div class="space-y-6">
    <div class="rounded-lg border border-subtle bg-background-800 p-6">
      <h3 class="mb-4 flex items-center text-lg font-semibold text-content">
        <i class="fas fa-share-alt text-brand-primary mr-3" />
        {{ $t('admin.construction.appearance.sections.social') }}
      </h3>

      <div class="space-y-5">
        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-globe text-brand-primary mr-2" />{{ $t('admin.construction.appearance.officialSite') }}
          </label>
          <div class="flex-1 space-y-3">
            <div class="flex flex-col gap-2">
              <CyberSwitch v-model="localSettings.show_official_site" />
              <span class="text-xs text-content-muted">{{ $t('admin.construction.appearance.officialSiteHint') }}</span>
            </div>
            <CyberInput
              v-model="localSettings.official_site_url"
              :placeholder="$t('admin.construction.appearance.officialSitePlaceholder')"
              :disabled="!localSettings.show_official_site"
              width="500px"
            />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fab fa-github text-brand-primary mr-2" />{{ $t('admin.construction.appearance.githubLink') }}
          </label>
          <div class="flex-1 space-y-3">
            <div class="flex flex-col gap-2">
              <CyberSwitch v-model="localSettings.show_github_link" />
              <span class="text-xs text-content-muted">{{ $t('admin.construction.appearance.githubHint') }}</span>
            </div>
            <CyberInput
              v-model="localSettings.github_url"
              :placeholder="$t('admin.construction.appearance.githubPlaceholder')"
              :disabled="!localSettings.show_github_link"
              width="500px"
            />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fab fa-weixin text-brand-primary mr-2" />{{ $t('admin.construction.appearance.wechatGroup') }}
          </label>
          <div class="flex-1 space-y-3">
            <div class="flex flex-col gap-2">
              <CyberSwitch v-model="localSettings.show_wechat_group" />
              <span class="text-xs text-content-muted">{{ $t('admin.construction.appearance.wechatGroupHint') }}</span>
            </div>
            <CyberInput
              v-model="localSettings.wechat_qr_image_url"
              :placeholder="$t('admin.construction.appearance.wechatQr')"
              :disabled="!localSettings.show_wechat_group"
              width="500px"
            />
            <CyberInput
              v-model="localSettings.wechat_contact_account"
              :placeholder="$t('admin.construction.appearance.wechatContact')"
              :disabled="!localSettings.show_wechat_group"
              width="500px"
            />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fab fa-qq text-brand-primary mr-2" />{{ $t('admin.construction.appearance.qqGroup') }}
          </label>
          <div class="flex-1 space-y-3">
            <div class="flex flex-col gap-2">
              <CyberSwitch v-model="localSettings.show_qq_group" />
              <span class="text-xs text-content-muted">{{ $t('admin.construction.appearance.qqGroupHint') }}</span>
            </div>
            <CyberInput
              v-model="localSettings.qq_qr_image_url"
              :placeholder="$t('admin.construction.appearance.qqQr')"
              :disabled="!localSettings.show_qq_group"
              width="500px"
            />
            <CyberInput
              v-model="localSettings.qq_group_number"
              :placeholder="$t('admin.construction.appearance.qqNumber')"
              :disabled="!localSettings.show_qq_group"
              width="500px"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="rounded-lg border border-subtle bg-background-800 p-6">
      <h3 class="mb-4 flex items-center text-lg font-semibold text-content">
        <i class="fas fa-layout text-brand-primary mr-3" />
        {{ $t('admin.construction.appearance.sections.layout') }}
      </h3>

      <div class="space-y-5">
        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-toggle-on text-brand-primary mr-2" />{{ $t('admin.construction.appearance.multiLayout') }}
          </label>
          <div class="flex-1">
            <div class="flex flex-col gap-2">
              <CyberSwitch v-model="localSettings.enable_multi_layout" />
              <span class="text-xs text-content-muted">{{ $t('admin.construction.appearance.multiLayoutHint') }}</span>
            </div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-home text-brand-primary mr-2" />{{ $t('admin.construction.appearance.defaultLayout') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.default_layout"
              :options="layoutOptions"
              :disabled="!localSettings.enable_multi_layout"
              :placeholder="$t('admin.construction.appearance.defaultLayoutPlaceholder')"
              width="500px"
            />
            <div class="mt-2 text-xs text-content-muted">{{ $t('admin.construction.appearance.defaultLayoutHint') }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="rounded-lg border border-subtle bg-background-800 p-6">
      <h3 class="mb-4 flex items-center text-lg font-semibold text-content">
        <i class="fas fa-language text-brand-primary mr-3" />
        {{ $t('admin.construction.appearance.sections.language') }}
      </h3>

      <div class="space-y-5">
        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-globe text-brand-primary mr-2" />{{ $t('admin.construction.appearance.multiLanguage') }}
          </label>
          <div class="flex-1">
            <div class="flex flex-col gap-2">
              <CyberSwitch v-model="localSettings.enable_multi_language" />
              <span class="text-xs text-content-muted">{{ $t('admin.construction.appearance.multiLanguageHint') }}</span>
            </div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-start">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-flag text-brand-primary mr-2" />{{ $t('admin.construction.appearance.defaultLanguage') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.default_language"
              :options="languageOptions"
              :disabled="!localSettings.enable_multi_language"
              :placeholder="$t('admin.construction.appearance.defaultLanguagePlaceholder')"
              width="500px"
            />
            <div class="mt-2 text-xs text-content-muted">{{ $t('admin.construction.appearance.defaultLanguageHint') }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
