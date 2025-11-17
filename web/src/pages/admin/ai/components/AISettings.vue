<script setup lang="ts">
  import { computed, onMounted, reactive, watch } from 'vue'
  import { type Setting, defaultSettings } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const aiDefaults = defaultSettings.ai.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  const localSettings = reactive({
    ai_enabled: aiDefaults.ai_enabled ?? false,
    ai_proxy: aiDefaults.ai_proxy || '',
    ai_model: aiDefaults.ai_model || 'gpt-4o-mini',
    ai_api_key: aiDefaults.ai_api_key || '',
    ai_temperature: aiDefaults.ai_temperature ?? 0.1,
    ai_max_tokens: aiDefaults.ai_max_tokens || 8096,
    ai_concurrency: aiDefaults.ai_concurrency || 5,
  })

  const sanitizeInput = (value: string): string => {
    if (!value) return ''
    return value
      .trim() // 去除首尾空格
      .replace(/[\s\n\r\t]/g, '') // 去除所有空白字符（空格、换行、制表符等）
      .replace(/["'`]/g, '') // 去除所有引号（单引号、双引号、反引号）
  }

  watch(
    () => localSettings.ai_proxy,
    (newValue) => {
      const sanitized = sanitizeInput(newValue as string)
      if (sanitized !== newValue) {
        localSettings.ai_proxy = sanitized
      }
    }
  )

  watch(
    () => localSettings.ai_api_key,
    (newValue) => {
      const sanitized = sanitizeInput(newValue as string)
      if (sanitized !== newValue) {
        localSettings.ai_api_key = sanitized
      }
    }
  )

  const modelOptions = computed(() => {
    const presetModels = [
      'gpt-5-mini',
      'gpt-5-nano',
      'gpt-4.1-nano',
      'gpt-4o-mini',
      'gpt-4.1-mini',
      'o3-mini',
      'o4-mini',
      'o4-mini-high',
      'gpt-4o',
      'gpt-4-turbo-with-vision',
      'gpt-4.1-mini-2025-04-14',
      'gpt-4.1-nano-2025-04-14',
    ]

    // 检查当前选中的模型是否在预设列表中
    const currentModel = localSettings.ai_model as string
    if (currentModel && !presetModels.includes(currentModel)) {
      // 如果不在预设列表中，说明是自定义模型，添加到列表开头
      return [currentModel, ...presetModels].map((model) => ({
        label: model,
        value: model,
      }))
    }

    return presetModels.map((model) => ({
      label: model,
      value: model,
    }))
  })

  const getSettingsArray = () =>
    Object.entries(localSettings).map(([key, value]) => ({
      key,
      value,
      type: typeof value === 'number' ? 'number' : typeof value === 'boolean' ? 'boolean' : 'string',
      group: 'ai',
      description: '',
      is_system: true,
    }))

  const initializeSettings = () => {
    props.settings.forEach((setting) => {
      if (setting.key in localSettings) {
        ;(localSettings as Record<string, unknown>)[setting.key] = setting.value
      }
    })
  }

  watch(
    localSettings,
    () => {
      emit('update', getSettingsArray())
    },
    { deep: true }
  )

  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        initializeSettings()
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    initializeSettings()
  })
</script>

<template>
  <div class="ai-settings-container">
    <div class="settings-list">
      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-power-off mr-2 text-brand-700" />{{ $t('admin.ai.aiSettings.enableAI') }}
        </label>
        <div class="flex-1">
          <div class="flex items-center">
            <CyberSwitch v-model="localSettings.ai_enabled" />
            <span class="ml-3 text-sm text-content-muted">{{
              localSettings.ai_enabled ? $t('admin.ai.aiSettings.aiEnabled') : $t('admin.ai.aiSettings.aiDisabled')
            }}</span>
          </div>
          <p class="setting-description">{{ $t('admin.ai.aiSettings.enableDesc') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-network-wired mr-2 text-brand-700" />{{ $t('admin.ai.aiSettings.proxyUrl') }}
        </label>
        <div class="flex-1">
          <CyberInput v-model="localSettings.ai_proxy" :placeholder="$t('admin.ai.aiSettings.proxyPlaceholder')" width="500px" />
          <p class="setting-description">{{ $t('admin.ai.aiSettings.proxyDesc') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-key mr-2 text-brand-700" />{{ $t('admin.ai.aiSettings.apiKey') }}
        </label>
        <div class="flex-1">
          <CyberInput
            v-model="localSettings.ai_api_key"
            type="password"
            :placeholder="$t('admin.ai.aiSettings.apiKeyPlaceholder')"
            width="500px"
          />
          <p class="setting-description">{{ $t('admin.ai.aiSettings.apiKeyDesc') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-robot mr-2 text-brand-700" />{{ $t('admin.ai.aiSettings.modelSelect') }}
        </label>
        <div class="flex-1">
          <CyberDropdown
            v-model="localSettings.ai_model"
            :options="modelOptions"
            :placeholder="$t('admin.ai.aiSettings.modelPlaceholder')"
            width="500px"
            searchable
            allow-create
            :create-label="$t('admin.ai.aiSettings.createModel')"
          />
          <p class="setting-description">{{ $t('admin.ai.aiSettings.modelDesc') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-thermometer-half mr-2 text-brand-700" />{{ $t('admin.ai.aiSettings.temperature') }}
        </label>
        <div class="flex-1">
          <div class="slider-container flex items-center">
            <CyberSlider
              v-model="localSettings.ai_temperature"
              :min="0"
              :max="1"
              :step="0.01"
              :description="$t('admin.ai.aiSettings.temperatureDesc')"
              width="460px"
            />
            <span class="ml-3 flex h-full items-center text-sm text-content-muted">{{
              localSettings.ai_temperature.toFixed(2)
            }}</span>
          </div>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-text-width mr-2 text-brand-700" />{{ $t('admin.ai.aiSettings.maxTokens') }}
        </label>
        <div class="flex-1">
          <CyberInput
            v-model.number="localSettings.ai_max_tokens"
            type="number"
            :placeholder="$t('admin.ai.aiSettings.maxTokensPlaceholder')"
            width="500px"
          />
          <p class="setting-description">{{ $t('admin.ai.aiSettings.maxTokensDesc') }}</p>
        </div>
      </div>

      <div class="flex flex-col md:flex-row md:items-center">
        <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
          <i class="fas fa-tasks mr-2 text-brand-700" />{{ $t('admin.ai.aiSettings.concurrency') }}
        </label>
        <div class="flex-1">
          <CyberInput
            v-model.number="localSettings.ai_concurrency"
            type="number"
            :placeholder="$t('admin.ai.aiSettings.concurrencyPlaceholder')"
            width="500px"
            :min="1"
            :max="20"
          >
            <template #unit>{{ $t('admin.ai.aiSettings.unitCount') }}</template>
          </CyberInput>
          <p class="setting-description">{{ $t('admin.ai.aiSettings.concurrencyDesc') }}</p>
        </div>
      </div>

      <ul class="feature-list">
        <li>• {{ $t('admin.ai.aiSettings.features.autoTag') }}</li>
        <li>• {{ $t('admin.ai.aiSettings.features.smartCategory') }}</li>
        <li>• {{ $t('admin.ai.aiSettings.features.batchProcess') }}</li>
        <li>• {{ $t('admin.ai.aiSettings.features.detailedDesc') }}</li>
      </ul>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .ai-settings-container {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
  }

  .settings-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
  }

  .setting-description {
    margin-top: var(--space-xs);
    font-size: var(--text-xs);
    color: var(--color-content-subtle);
  }

  .feature-list {
    color: var(--color-content-muted);
    margin-top: var(--space-xl);
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
    font-size: var(--text-sm);
  }

  .hover:shadow-glow:hover {
    box-shadow: var(--shadow-glow-md);
  }

  .slider-container {
    margin-bottom: var(--space-sm);
  }
</style>
