<script setup lang="ts">
  import { computed, reactive, ref, watch } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { type Setting } from '@/api/admin/settings'
  import { testQdrantConnection } from '@/api/admin/settings/vector'
  import { useSettingsStore } from '@/store/settings'

  const { $t } = useTexts()
  const settingsStore = useSettingsStore()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    update: [settings: Setting[]]
  }>()

  const toast = useToast()
  const isTestingQdrant = ref(false)

  // 判断是否为 Docker/Compose 模式
  const isDockerMode = computed(() => {
    const deployMode = settingsStore.rawSettings?.deploy_mode
    return deployMode === 'docker' || deployMode === 'compose'
  })

  const localSettings = reactive({
    vector_enabled: false,
    vector_provider: 'openai',
    vector_model: 'text-embedding-3-small',
    vector_api_key: '',
    vector_base_url: 'https://api.openai.com/v1',
    vector_timeout: 30,
    vector_similarity_threshold: 0.7,
    vector_search_threshold: 0.3,
    vector_max_results: 50,
    vector_concurrency: 3,
    qdrant_url: 'http://localhost:6333',
    qdrant_timeout: 30,
  })

  const providerOptions = [{ label: 'OpenAI', value: 'openai' }]

  const modelOptions = computed(() => {
    const presetModels = [
      { label: $t('admin.ai.vectorSettings.models.small'), value: 'text-embedding-3-small' },
      { label: $t('admin.ai.vectorSettings.models.large'), value: 'text-embedding-3-large' },
      { label: $t('admin.ai.vectorSettings.models.ada'), value: 'text-embedding-ada-002' },
    ]

    // 检查当前选中的模型是否在预设列表中
    const currentModel = localSettings.vector_model as string
    const presetValues = presetModels.map((m) => m.value)

    if (currentModel && !presetValues.includes(currentModel)) {
      // 如果不在预设列表中，说明是自定义模型，添加到列表开头
      return [{ label: currentModel, value: currentModel }, ...presetModels]
    }

    return presetModels
  })

  const isInitializing = ref(true)

  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        newSettings.forEach((setting) => {
          if (setting.key in localSettings) {
            ;(localSettings as Record<string, unknown>)[setting.key] = setting.value
          }
        })
        if (isInitializing.value) {
          setTimeout(() => {
            isInitializing.value = false
          }, 100)
        }
      }
    },
    { immediate: true }
  )

  watch(
    () => localSettings.vector_model,
    (newModel, oldModel) => {
      if (isInitializing.value || !oldModel || newModel === oldModel) {
        return
      }

      toast.success($t('admin.ai.messages.vectorModelChanged').replace('{model}', newModel))
    }
  )

  watch(
    localSettings,
    () => {
      const updatedSettings: Setting[] = Object.keys(localSettings).map((key) => {
        const originalSetting = props.settings.find((s) => s.key === key)
        const settingsRecord = localSettings as Record<string, unknown>
        return {
          ...originalSetting,
          key,
          value: settingsRecord[key],
          type: getSettingType(settingsRecord[key]),
          group: 'vector',
          description: getSettingDescription(key),
          is_system: true,
        } as Setting
      })

      emit('update', updatedSettings)
    },
    { deep: true }
  )

  function getSettingType(value: unknown): 'string' | 'number' | 'boolean' {
    if (typeof value === 'boolean') {
      return 'boolean'
    }
    if (typeof value === 'number') {
      return 'number'
    }
    return 'string'
  }

  function getSettingDescription(key: string): string {
    const descriptions: Record<string, string> = {
      vector_enabled: $t('admin.ai.vectorSettings.vectorDesc'),
      vector_provider: $t('admin.ai.vectorSettings.providerDesc'),
      vector_model: $t('admin.ai.vectorSettings.modelDesc'),
      vector_api_key: $t('admin.ai.vectorSettings.apiKeyDesc'),
      vector_base_url: $t('admin.ai.vectorSettings.apiBaseUrlDesc'),
      vector_timeout: $t('admin.ai.vectorSettings.apiTimeoutDesc'),
      vector_similarity_threshold: $t('admin.ai.vectorSettings.similarityThresholdDesc'),
      vector_search_threshold: $t('admin.ai.vectorSettings.searchThresholdDesc'),
      vector_max_results: $t('admin.ai.vectorSettings.maxResultsDesc'),
      vector_concurrency: $t('admin.ai.vectorSettings.concurrencyDesc'),
      qdrant_url: $t('admin.ai.vectorSettings.qdrantUrlDesc'),
      qdrant_timeout: $t('admin.ai.vectorSettings.qdrantTimeoutDesc'),
    }
    return descriptions[key] || ''
  }

  async function handleTestQdrant() {
    if (!localSettings.qdrant_url) {
      toast.warning($t('admin.ai.messages.needQdrantUrl'))
      return
    }

    isTestingQdrant.value = true
    try {
      const result = await testQdrantConnection({
        qdrant_url: localSettings.qdrant_url,
        qdrant_timeout: localSettings.qdrant_timeout || 30,
      })

      if (result.data?.success) {
        const details = result.data.details
        let message = $t('admin.ai.messages.qdrantConnectSuccess')
        if (details) {
          if (details.version) {
            message += `\n${$t('admin.ai.testResults.version')}: ${details.version}`
          }
          if (details.total_vectors !== undefined) {
            message += `\n${$t('admin.ai.testResults.totalVectors')}: ${details.total_vectors}`
          }
          if (details.response_time !== undefined) {
            message += `\n${$t('admin.ai.testResults.responseTime')}: ${details.response_time}ms`
          }
        }
        toast.success(message)
      } else {
        toast.error(result.data?.message || $t('admin.ai.messages.qdrantConnectFailed'))
      }
    } catch (error: unknown) {
      const err = error as Error
      toast.error(
        $t('admin.ai.messages.qdrantTestFailed').replace('{error}', err.message || $t('admin.ai.messages.unknownError'))
      )
    } finally {
      isTestingQdrant.value = false
    }
  }
</script>

<template>
  <div class="vector-settings-container">
    <div class="settings-list">
      <div class="settings-group">
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-toggle-on mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.enableVector') }}
          </label>
          <div class="flex-1">
            <CyberSwitch v-model="localSettings.vector_enabled" />
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.vectorDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-search mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.searchThreshold') }}
          </label>
          <div class="flex-1">
            <div class="slider-container flex items-center">
              <CyberSlider
                v-model="localSettings.vector_search_threshold"
                :min="0"
                :max="1"
                :step="0.01"
                :description="$t('admin.ai.vectorSettings.searchThresholdDesc')"
                width="460px"
              />
              <span class="text-content-content-muted ml-3 flex h-full items-center text-sm">{{
                localSettings.vector_search_threshold.toFixed(2)
              }}</span>
            </div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-cloud mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.provider') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.vector_provider"
              :options="providerOptions"
              :placeholder="$t('admin.ai.vectorSettings.providerPlaceholder')"
              width="500px"
            />
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.providerDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-network-wired mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.apiBaseUrl') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.vector_base_url"
              :placeholder="$t('admin.ai.vectorSettings.apiBaseUrlPlaceholder')"
              width="500px"
            />
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.apiBaseUrlDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-key mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.apiKey') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.vector_api_key"
              type="password"
              :placeholder="$t('admin.ai.vectorSettings.apiKeyPlaceholder')"
              width="500px"
            />
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.apiKeyDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-robot mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.model') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.vector_model"
              :options="modelOptions"
              :placeholder="$t('admin.ai.vectorSettings.modelPlaceholder')"
              width="500px"
              searchable
              allow-create
              :create-label="$t('admin.ai.vectorSettings.createModel')"
            />
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.modelDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-clock mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.apiTimeout') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model.number="localSettings.vector_timeout"
              type="number"
              :placeholder="$t('admin.ai.vectorSettings.apiTimeoutPlaceholder')"
              width="500px"
            >
              <template #unit>{{ $t('admin.ai.vectorSettings.unitSeconds') }}</template>
            </CyberInput>
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.apiTimeoutDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-percentage mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.similarityThreshold') }}
          </label>
          <div class="flex-1">
            <div class="slider-container flex items-center">
              <CyberSlider
                v-model="localSettings.vector_similarity_threshold"
                :min="0"
                :max="1"
                :step="0.01"
                :description="$t('admin.ai.vectorSettings.similarityThresholdDesc')"
                width="460px"
              />
              <span class="text-content-content-muted ml-3 flex h-full items-center text-sm">{{
                localSettings.vector_similarity_threshold.toFixed(2)
              }}</span>
            </div>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-list-ol mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.maxResults') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model.number="localSettings.vector_max_results"
              type="number"
              :placeholder="$t('admin.ai.vectorSettings.maxResultsPlaceholder')"
              width="500px"
            >
              <template #unit>{{ $t('admin.ai.vectorSettings.unitResults') }}</template>
            </CyberInput>
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.maxResultsDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-tasks mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.concurrency') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model.number="localSettings.vector_concurrency"
              type="number"
              :placeholder="$t('admin.ai.vectorSettings.concurrencyPlaceholder')"
              width="500px"
              :min="1"
              :max="10"
            >
              <template #unit>{{ $t('admin.ai.vectorSettings.unitCount') }}</template>
            </CyberInput>
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.concurrencyDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-database mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.qdrantUrl') }}
          </label>
          <div class="flex-1">
            <div class="flex items-center gap-2">
              <CyberInput
                v-model="localSettings.qdrant_url"
                :placeholder="$t('admin.ai.vectorSettings.qdrantUrlPlaceholder')"
                :disabled="isDockerMode"
                width="400px"
              />
              <CyberButton
                :loading="isTestingQdrant"
                :disabled="!localSettings.qdrant_url || isDockerMode"
                size="small"
                type="primary"
                @click="handleTestQdrant"
              >
                <i class="fas fa-vial mr-1" />
                {{ isTestingQdrant ? $t('admin.ai.actions.testing') : $t('admin.ai.actions.testConnection') }}
              </CyberButton>
            </div>
            <p v-if="isDockerMode" class="mt-1 text-xs text-yellow-600">
              <i class="fas fa-info-circle mr-1" />{{ $t('admin.ai.vectorSettings.qdrantUrlDockerWarning') }}
            </p>
            <p v-else class="setting-description">{{ $t('admin.ai.vectorSettings.qdrantUrlDesc') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-stopwatch mr-2 text-brand-700" />{{ $t('admin.ai.vectorSettings.qdrantTimeout') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model.number="localSettings.qdrant_timeout"
              type="number"
              :placeholder="$t('admin.ai.vectorSettings.qdrantTimeoutPlaceholder')"
              width="500px"
            >
              <template #unit>{{ $t('admin.ai.vectorSettings.unitSeconds') }}</template>
            </CyberInput>
            <p class="setting-description">{{ $t('admin.ai.vectorSettings.qdrantTimeoutDesc') }}</p>
          </div>
        </div>

        <ul class="feature-list">
          <li>• {{ $t('admin.ai.vectorSettings.features.naturalLanguage') }}</li>
          <li>• {{ $t('admin.ai.vectorSettings.features.imageSimilarity') }}</li>
          <li>• {{ $t('admin.ai.vectorSettings.features.smartRecommend') }}</li>
          <li>• {{ $t('admin.ai.vectorSettings.features.semanticSearch') }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .vector-settings-container {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
  }

  .settings-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
  }

  .settings-group {
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

  .slider-container {
    margin-bottom: var(--space-sm);
  }
</style>
