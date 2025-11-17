<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { useTheme, type VisualTheme } from '@/composables/useTheme'
  import { useLayoutStore } from '@/store/layout'
  import { useAppearanceSettings } from '@/composables/useAppearanceSettings'
  import type { SupportedLocale } from '@/locales/zh-CN'
  import { StorageUtil } from '@/utils/storage'
  import { useToast } from '@/components/Toast/useToast'
  import SettingCard from '@/components/SettingCard/index.vue'

  const { currentTheme, currentLocale, setTheme: setTextTheme, setLocale, $t } = useTexts()
  const { selectedTheme, setTheme: setVisualTheme, themeOptions, allThemes } = useTheme()
  const layoutStore = useLayoutStore()
  const { layoutSettings } = useAppearanceSettings()
  const toast = useToast()

  /* ==================== 视觉主题选项 - 使用统一的主题配置 ==================== */
  const visualThemeOptions = themeOptions

  /* ==================== 文案风格选项 ==================== */
  const textStyleOptionsWithHints = computed(() => [
    { label: $t('settings.preferences.options.normalHint'), value: 'normal' },
    { label: $t('settings.preferences.options.cyberHint'), value: 'cyber' },
  ])

  /* ==================== 语言选项 ==================== */
  const localeOptions = computed(() => [
    { label: $t('settings.preferences.languages.zhCN'), value: 'zh-CN' },
    { label: $t('settings.preferences.languages.enUS'), value: 'en-US' },
    { label: $t('settings.preferences.languages.jaJP'), value: 'ja-JP' },
  ])

  /* ==================== 布局模式选项 ==================== */
  const layoutModeOptions = computed(() => [
    { label: $t('settings.preferences.layoutModes.top'), value: 'top' },
    { label: $t('settings.preferences.layoutModes.left'), value: 'left' },
  ])

  /* ==================== 状态 ==================== */
  const selectedVisualTheme = ref<VisualTheme>('dark')
  const selectedTextStyle = ref<string>('normal')
  const selectedLocale = ref<SupportedLocale>('zh-CN')
  const selectedLayoutMode = ref<string>('top')

  /* ==================== 加载偏好设置 ==================== */
  const loadPreferences = () => {
    const savedVisualTheme = StorageUtil.get<string>('visual-theme')
    if (savedVisualTheme && allThemes.includes(savedVisualTheme as VisualTheme)) {
      selectedVisualTheme.value = savedVisualTheme as VisualTheme
    } else {
      selectedVisualTheme.value = selectedTheme.value
    }

    const savedTextTheme = StorageUtil.get<string>('text-theme')
    if (savedTextTheme && ['normal', 'cyber'].includes(savedTextTheme)) {
      selectedTextStyle.value = savedTextTheme
      setTextTheme(savedTextTheme as 'normal' | 'cyber')
    } else {
      selectedTextStyle.value = currentTheme.value
    }

    selectedLocale.value = currentLocale.value
    selectedLayoutMode.value = layoutStore.mode
  }

  /* ==================== 监听语言变化 ==================== */
  watch(currentLocale, (newLocale) => {
    selectedLocale.value = newLocale
  })

  const saveVisualTheme = (theme: VisualTheme) => {
    try {
      selectedVisualTheme.value = theme
      setVisualTheme(theme)
      toast.success($t('settings.preferences.notifications.themeSuccess'))
    } catch (_error) {
      toast.error($t('settings.preferences.notifications.themeFailure'))
    }
  }

  const saveTextStyle = (style: string) => {
    try {
      selectedTextStyle.value = style
      setTextTheme(style as 'normal' | 'cyber')
      StorageUtil.set('text-theme', style)
      toast.success($t('settings.preferences.notifications.success'))
    } catch (_error) {
      toast.error($t('settings.preferences.notifications.failure'))
    }
  }

  const saveLocale = (locale: SupportedLocale) => {
    try {
      selectedLocale.value = locale
      setLocale(locale)
      toast.success($t('settings.preferences.notifications.localeSuccess'))
    } catch (_error) {
      toast.error($t('settings.preferences.notifications.localeFailure'))
    }
  }

  const saveLayoutMode = (mode: string) => {
    try {
      selectedLayoutMode.value = mode
      layoutStore.setLayoutMode(mode as 'top' | 'left')
      toast.success($t('settings.preferences.notifications.layoutSuccess'))
    } catch (_error) {
      toast.error($t('settings.preferences.notifications.layoutFailure'))
    }
  }

  onMounted(() => {
    loadPreferences()
  })

  defineExpose({
    loadPreferences,
  })
</script>

<template>
  <div class="preferences-section">
    <div class="section-header">
      <div class="header-content">
        <div class="title-wrapper">
          <div class="title-icon">
            <i class="fas fa-cog" />
          </div>
          <div class="title-text">
            <h2 class="title">{{ $t('settings.preferences.title') }}</h2>
            <p class="subtitle">{{ $t('settings.preferences.description') }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 语言偏好 -->
    <SettingCard
      icon="fas fa-language"
      :title="$t('settings.preferences.section.languageTitle')"
      :hint="$t('settings.preferences.section.languageHint')"
    >
      <CyberRadioGroup
        v-model="selectedLocale"
        :options="localeOptions"
        layout="horizontal"
        @update:model-value="saveLocale"
      />
    </SettingCard>

    <!-- 视觉主题 -->
    <SettingCard
      icon="fas fa-palette"
      :title="$t('settings.preferences.section.visualThemeTitle')"
      :hint="$t('settings.preferences.section.visualThemeHint')"
    >
      <div class="theme-grid">
        <button
          v-for="option in visualThemeOptions"
          :key="option.value"
          class="theme-option"
          :class="{ active: selectedVisualTheme === option.value }"
          @click="saveVisualTheme(option.value)"
        >
          <div class="theme-icon">
            <i :class="['fas', `fa-${option.icon}`]" />
          </div>
          <div class="theme-info">
            <div class="theme-label">{{ option.label }}</div>
            <div class="theme-description">{{ option.description }}</div>
          </div>
          <div v-if="selectedVisualTheme === option.value" class="theme-check">
            <i class="fas fa-check-circle" />
          </div>
        </button>
      </div>
    </SettingCard>

    <!-- 布局模式 -->
    <SettingCard
      icon="fas fa-grip-horizontal"
      :title="$t('settings.preferences.section.layoutTitle')"
      :hint="$t('settings.preferences.section.layoutHint')"
      :show="layoutSettings.multiLayoutEnabled"
    >
      <CyberRadioGroup
        v-model="selectedLayoutMode"
        :options="layoutModeOptions"
        layout="horizontal"
        @update:model-value="saveLayoutMode"
      />
    </SettingCard>

    <!-- 界面风格 -->
    <SettingCard
      icon="fas fa-font"
      :title="$t('settings.preferences.section.interfaceTitle')"
      :hint="$t('settings.preferences.section.interfaceHint')"
    >
      <CyberRadioGroup
        v-model="selectedTextStyle"
        :options="textStyleOptionsWithHints"
        layout="horizontal"
        @update:model-value="saveTextStyle"
      />
    </SettingCard>
  </div>
</template>

<style scoped>
  .preferences-section {
    @apply space-y-6;
  }

  .section-header {
    @apply relative overflow-hidden;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-800-rgb), 0.9), rgba(var(--color-background-900-rgb), 0.95));
    border: 1px solid var(--color-border-subtle);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
  }

  .header-content {
    @apply relative p-6;
    z-index: 1;
  }

  .title-wrapper {
    @apply flex items-start gap-4;
  }

  .title-icon {
    @apply flex h-12 w-12 items-center justify-center;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.18), rgba(var(--color-brand-500-rgb), 0.08));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.24);
    color: var(--color-brand-400);
    font-size: 18px;
    box-shadow:
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.18),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
  }

  .title-text {
    @apply flex flex-col gap-1;
  }

  .title {
    @apply text-xl font-semibold;
    color: var(--color-content-heading);
    margin: 0;
  }

  .subtitle {
    @apply text-sm leading-relaxed;
    color: var(--color-content-muted);
    margin: 0;
    max-width: 600px;
  }

  .theme-grid {
    @apply grid grid-cols-1 gap-2.5 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4;
  }

  .theme-option {
    @apply relative flex items-start gap-3 border p-3 text-left transition-all duration-200;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-background-800-rgb), 0.5);
    border-color: var(--color-border-subtle);
  }

  .theme-option:hover {
    border-color: var(--color-border-default);
    background: rgba(var(--color-background-700-rgb), 0.6);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .theme-option.active {
    border-color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    box-shadow:
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.2),
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .theme-icon {
    @apply flex h-10 w-10 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.05));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-400);
    font-size: 16px;
  }

  .theme-option.active .theme-icon {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25), rgba(var(--color-brand-500-rgb), 0.15));
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .theme-info {
    @apply flex flex-1 flex-col gap-1;
  }

  .theme-label {
    @apply text-sm font-medium;
    color: var(--color-content);
  }

  .theme-description {
    @apply text-xs;
    color: var(--color-content-muted);
  }

  .theme-check {
    @apply flex-shrink-0 text-base;
    color: var(--color-brand-500);
  }

  @media (max-width: 768px) {
    .title-wrapper {
      @apply flex-col gap-2;
    }

    .title-icon {
      @apply self-start;
    }

    .theme-grid {
      @apply grid-cols-1;
    }
  }
</style>
