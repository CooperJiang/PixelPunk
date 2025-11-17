<script setup lang="ts">
  import { computed, nextTick, reactive, ref, watch } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import {
    type RegenerateVectorsParams,
    type Setting,
    type VectorTestParams,
    regenerateVectors as apiRegenerateVectors,
    testVectorConfig as apiTestVectorConfig,
  } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  /* Events */
  const emit = defineEmits<{
    update: [settings: Setting[]]
  }>()

  const toast = useToast()

  /* 测试状态 */
  const testing = ref(false)
  const regenerating = ref(false)
  const testResult = ref<{
    success: boolean
    message: string
    details?: string
  } | null>(null)

  /* 本地设置对象 */
  const localSettings = reactive({
    vector_enabled: false,
    vector_provider: 'openai',
    vector_model: 'text-embedding-3-small',
    vector_api_key: '',
    vector_base_url: 'https://api.openai.com/v1',
    vector_timeout: 30,
    vector_similarity_threshold: 0.3,
    vector_max_results: 50,
    vector_concurrency: 3,
  })

  /* 提供者选项 */
  const providerOptions = [{ label: 'OpenAI', value: 'openai' }]

  /* 模型选项 */
  const modelOptions = computed(() => [
    { label: $t('admin.settings.vector.models.small'), value: 'text-embedding-3-small' },
    { label: $t('admin.settings.vector.models.large'), value: 'text-embedding-3-large' },
    { label: $t('admin.settings.vector.models.ada'), value: 'text-embedding-ada-002' },
  ])

  /* 配置有效性检查 */
  const isConfigValid = computed(
    () => localSettings.vector_enabled && localSettings.vector_api_key.trim() !== '' && localSettings.vector_model.trim() !== ''
  )

  /* 标记是否为初始化阶段 */
  const isInitializing = ref(true)

  /* 标记是否正在恢复模型值（避免递归触发watch） */
  const isRestoringModel = ref(false)

  /* 监听props变化，更新本地设置 */
  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        newSettings.forEach((setting) => {
          if (setting.key in localSettings) {
            const settingsRecord = localSettings as Record<string, unknown>
            settingsRecord[setting.key] = setting.value
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
    async (newModel, oldModel) => {
      if (isInitializing.value || isRestoringModel.value || !oldModel || newModel === oldModel) {
        return
      }

      const confirmed = await showModelChangeWarning(oldModel, newModel)
      if (!confirmed) {
        isRestoringModel.value = true
        localSettings.vector_model = oldModel
        nextTick(() => {
          isRestoringModel.value = false
        })
        return
      }

      toast.success($t('admin.settings.vector.toast.modelChanged').replace('{model}', newModel))
    }
  )

  watch(
    localSettings,
    () => {
      const settingsRecord = localSettings as Record<string, unknown>
      const updatedSettings: Setting[] = Object.keys(localSettings).map((key) => {
        const originalSetting = props.settings.find((s) => s.key === key)
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

  function getSettingType(value: any): 'string' | 'number' | 'boolean' {
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
      vector_enabled: $t('admin.settings.vector.desc.enabled'),
      vector_provider: $t('admin.settings.vector.desc.provider'),
      vector_model: $t('admin.settings.vector.desc.model'),
      vector_api_key: $t('admin.settings.vector.desc.apiKey'),
      vector_base_url: $t('admin.settings.vector.desc.baseUrl'),
      vector_timeout: $t('admin.settings.vector.desc.timeout'),
      vector_similarity_threshold: $t('admin.settings.vector.desc.threshold'),
      vector_max_results: $t('admin.settings.vector.desc.maxResults'),
      vector_concurrency: $t('admin.settings.vector.desc.concurrency'),
    }
    return descriptions[key] || ''
  }

  const testVectorConfig = async () => {
    if (testing.value) {
      return
    }

    if (!isConfigValid.value) {
      toast.error($t('admin.settings.vector.toast.configIncomplete'))
      return
    }

    testing.value = true
    testResult.value = null

    try {
      const params: VectorTestParams = {
        provider: localSettings.vector_provider,
        model: localSettings.vector_model,
        api_key: localSettings.vector_api_key,
        base_url: localSettings.vector_base_url,
        timeout: localSettings.vector_timeout,
      }

      const result = await apiTestVectorConfig(params)

      if (result.success && result.data) {
        testResult.value = {
          success: result.data.success,
          message: result.data.message,
          details: result.data.details,
        }

        if (result.data.success) {
          toast.success($t('admin.settings.vector.toast.testSuccess'))
        } else {
          toast.error($t('admin.settings.vector.toast.testFailed'))
        }
      } else {
        throw new Error(result.message || $t('admin.settings.vector.toast.testRequestFailed'))
      }
    } catch (error) {
      testResult.value = {
        success: false,
        message: $t('admin.settings.vector.toast.testFailed'),
        details: (error as Error).message || $t('admin.settings.vector.toast.networkRequestFailed'),
      }
      toast.error($t('admin.settings.vector.toast.testFailed'))
    } finally {
      testing.value = false
    }
  }

  const regenerateVectors = async () => {
    if (regenerating.value) {
      return
    }

    if (!isConfigValid.value) {
      toast.error($t('admin.settings.vector.toast.configIncomplete'))
      return
    }

    const confirmed = await showDangerousConfirmDialog()
    if (!confirmed) {
      return
    }

    regenerating.value = true

    try {
      const params: RegenerateVectorsParams = {}

      const result = await apiRegenerateVectors(params)

      if (result.success && result.data) {
        toast.success(
          $t('admin.settings.vector.toast.regenerateSuccess').replace('{count}', result.data.updated_count.toString())
        )

        if (result.data.failed > 0) {
          toast.warning($t('admin.settings.vector.toast.regenerateWarning').replace('{count}', result.data.failed.toString()))
        }
      } else {
        throw new Error(result.message || $t('admin.settings.vector.toast.regenerateRequestFailed'))
      }
    } catch (error) {
      toast.error(
        $t('admin.settings.vector.toast.regenerateFailed').replace(
          '{message}',
          (error as Error).message || $t('admin.settings.vector.toast.unknownError')
        )
      )
    } finally {
      regenerating.value = false
    }
  }

  const showDangerousConfirmDialog = (): Promise<boolean> =>
    new Promise((resolve) => {
      const modalContainer = document.createElement('div')
      modalContainer.className = 'fixed inset-0 bg-background-900 bg-opacity-50 flex items-center justify-center z-50'

      const modalContent = document.createElement('div')
      modalContent.className = 'bg-background-700 border border-red-500/30 rounded-xl p-6 max-w-md mx-4 shadow-2xl'

      modalContent.innerHTML = `
      <div class="flex items-center mb-4">
        <i class="fas fa-exclamation-triangle text-red-400 text-2xl mr-3"></i>
        <h3 class="text-lg font-bold text-content-heading">${$t('admin.settings.vector.confirmDialog.title')}</h3>
      </div>
      <div class="text-content text-sm mb-4 space-y-2">
        <p><strong class="text-red-400">${$t('admin.settings.vector.confirmDialog.warning')}</strong></p>
        <p>${$t('admin.settings.vector.confirmDialog.point1')}</p>
        <p>${$t('admin.settings.vector.confirmDialog.point2')}</p>
        <p>${$t('admin.settings.vector.confirmDialog.point3')}</p>
        <p>${$t('admin.settings.vector.confirmDialog.point4')}</p>
      </div>
      <div class="text-xs text-content-content-muted mb-4 p-3 bg-red-900/20 border border-red-500/30 rounded">
        <p>${$t('admin.settings.vector.confirmDialog.riskNote')}</p>
      </div>
      <input type="text" id="confirmInput" class="w-full p-2 mb-4 bg-background-700 border border-default rounded text-content-heading text-sm" placeholder="${$t('admin.settings.vector.confirmDialog.placeholder')}">
      <div id="buttonsContainer" class="flex gap-3">
      </div>
    `

      modalContainer.appendChild(modalContent)
      document.body.appendChild(modalContainer)

      const buttonsContainer = modalContent.querySelector('#buttonsContainer')
      const confirmInput = modalContent.querySelector('#confirmInput') as HTMLInputElement

      const cancelBtn = document.createElement('button')
      cancelBtn.className =
        'flex-1 px-4 py-2 bg-background-600 hover:bg-hover-bg text-content-heading rounded transition-colors text-sm border border-default'
      cancelBtn.innerHTML = `<i class="fas fa-times mr-2"></i>${$t('admin.settings.vector.confirmDialog.cancel')}`

      const confirmBtn = document.createElement('button')
      confirmBtn.className =
        'flex-1 px-4 py-2 bg-background-600 cursor-not-allowed text-content-content-muted rounded transition-colors text-sm border border-default'
      confirmBtn.innerHTML = `<i class="fas fa-exclamation-triangle mr-2"></i>${$t('admin.settings.vector.confirmDialog.confirm')}`
      confirmBtn.disabled = true

      buttonsContainer?.appendChild(cancelBtn)
      buttonsContainer?.appendChild(confirmBtn)

      const riskKeyword = $t('admin.settings.vector.confirmDialog.riskKeyword')

      confirmInput.addEventListener('input', () => {
        const isValid = confirmInput.value.trim() === riskKeyword
        confirmBtn.disabled = !isValid
        if (isValid) {
          confirmBtn.className =
            'flex-1 px-4 py-2 bg-red-600 hover:bg-red-500 text-content-heading rounded transition-colors text-sm border border-red-500/30'
        } else {
          confirmBtn.className =
            'flex-1 px-4 py-2 bg-background-600 cursor-not-allowed text-content-content-muted rounded transition-colors text-sm border border-default'
        }
      })

      cancelBtn.addEventListener('click', () => {
        document.body.removeChild(modalContainer)
        resolve(false)
      })

      confirmBtn.addEventListener('click', () => {
        if (confirmInput.value.trim() === riskKeyword) {
          document.body.removeChild(modalContainer)
          resolve(true)
        }
      })

      modalContainer.addEventListener('click', (e) => {
        if (e.target === modalContainer) {
          document.body.removeChild(modalContainer)
          resolve(false)
        }
      })

      setTimeout(() => confirmInput.focus(), 100)
    })

  const showModelChangeWarning = (oldModel: string, newModel: string): Promise<boolean> =>
    new Promise((resolve) => {
      const modalContainer = document.createElement('div')
      modalContainer.className = 'fixed inset-0 bg-background-900 bg-opacity-50 flex items-center justify-center z-50'

      const modalContent = document.createElement('div')
      modalContent.className = 'bg-background-700 border border-amber-500/30 rounded-xl p-6 max-w-md mx-4 shadow-2xl'

      modalContent.innerHTML = `
      <div class="flex items-center mb-4">
        <i class="fas fa-exclamation-triangle text-amber-400 text-2xl mr-3"></i>
        <h3 class="text-lg font-bold text-content-heading">${$t('admin.settings.vector.modelChangeDialog.title')}</h3>
      </div>
      <div class="text-content text-sm mb-4 space-y-2">
        <p><strong class="text-amber-400">${$t('admin.settings.vector.modelChangeDialog.warning')}</strong></p>
        <p>• ${$t('admin.settings.vector.modelChangeDialog.changeLine', { old: '<code class="bg-background-700 px-1 rounded text-cyan-300">' + oldModel + '</code>', new: '<code class="bg-background-700 px-1 rounded text-green-300">' + newModel + '</code>' })}</p>
        <p>${$t('admin.settings.vector.modelChangeDialog.point1')}</p>
        <p>${$t('admin.settings.vector.modelChangeDialog.point2')}</p>
        <p>${$t('admin.settings.vector.modelChangeDialog.point3')}</p>
      </div>
      <div class="text-xs text-content-content-muted mb-4 p-3 bg-amber-900/20 border border-amber-500/30 rounded">
        <p>${$t('admin.settings.vector.modelChangeDialog.notice')}</p>
      </div>
      <div id="buttonsContainer" class="flex gap-3">
      </div>
    `

      modalContainer.appendChild(modalContent)
      document.body.appendChild(modalContainer)

      const buttonsContainer = modalContent.querySelector('#buttonsContainer')

      const cleanup = () => {
        try {
          if (modalContainer && modalContainer.parentNode) {
            modalContainer.parentNode.removeChild(modalContainer)
          }
        } catch {}
      }

      const cancelBtn = document.createElement('button')
      cancelBtn.className =
        'flex-1 px-4 py-2 bg-background-600 hover:bg-hover-bg text-content-heading rounded transition-colors text-sm border border-default'
      cancelBtn.innerHTML = `<i class="fas fa-times mr-2"></i>${$t('admin.settings.vector.modelChangeDialog.cancel')}`

      const confirmBtn = document.createElement('button')
      confirmBtn.className =
        'flex-1 px-4 py-2 bg-amber-600 hover:bg-amber-500 text-content-heading rounded transition-colors text-sm border border-amber-500/30'
      confirmBtn.innerHTML = `<i class="fas fa-check mr-2"></i>${$t('admin.settings.vector.modelChangeDialog.confirm')}`

      buttonsContainer?.appendChild(cancelBtn)
      buttonsContainer?.appendChild(confirmBtn)

      cancelBtn.addEventListener('click', () => {
        cleanup()
        resolve(false)
      })

      confirmBtn.addEventListener('click', () => {
        cleanup()
        resolve(true)
      })

      modalContainer.addEventListener('click', (e) => {
        if (e.target === modalContainer) {
          cleanup()
          resolve(false)
        }
      })

      const handleEscape = (e: KeyboardEvent) => {
        if (e.key === 'Escape') {
          document.removeEventListener('keydown', handleEscape)
          cleanup()
          resolve(false)
        }
      }
      document.addEventListener('keydown', handleEscape)
    })
</script>

<template>
  <div class="space-y-5">
    <div class="space-y-5">
      <div class="space-y-5">
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-toggle-on mr-2 text-brand-600" />{{ $t('admin.settings.vector.enable.label') }}
          </label>
          <div class="flex-1">
            <CyberSwitch v-model="localSettings.vector_enabled" />
            <p class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.settings.vector.enable.hint') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-cloud mr-2 text-brand-600" />{{ $t('admin.settings.vector.provider.label') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.vector_provider"
              :options="providerOptions"
              :placeholder="$t('admin.settings.vector.provider.placeholder')"
              width="500px"
            />
            <p class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.settings.vector.provider.hint') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-network-wired mr-2 text-brand-600" />{{ $t('admin.settings.vector.baseUrl.label') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.vector_base_url"
              :placeholder="$t('admin.settings.vector.baseUrl.placeholder')"
              width="500px"
            />
            <p class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.settings.vector.baseUrl.hint') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-key mr-2 text-brand-600" />{{ $t('admin.settings.vector.apiKey.label') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.vector_api_key"
              type="password"
              :placeholder="$t('admin.settings.vector.apiKey.placeholder')"
              width="500px"
            />
            <p class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.settings.vector.apiKey.hint') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-robot mr-2 text-brand-600" />{{ $t('admin.settings.vector.model.label') }}
          </label>
          <div class="flex-1">
            <CyberDropdown
              v-model="localSettings.vector_model"
              :options="modelOptions"
              :placeholder="$t('admin.settings.vector.model.placeholder')"
              width="500px"
            />
            <p class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.settings.vector.model.hint') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-clock mr-2 text-brand-600" />{{ $t('admin.settings.vector.timeout.label') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model.number="localSettings.vector_timeout"
              type="number"
              :placeholder="$t('admin.settings.vector.timeout.placeholder')"
              width="500px"
            >
              <template #unit>{{ $t('admin.settings.vector.timeout.unit') }}</template>
            </CyberInput>
            <p class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.settings.vector.timeout.hint') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-percentage mr-2 text-brand-600" />{{ $t('admin.settings.vector.threshold.label') }}
          </label>
          <div class="flex-1">
            <div class="slider-container flex items-center">
              <CyberSlider
                v-model="localSettings.vector_similarity_threshold"
                :min="0"
                :max="1"
                :step="0.01"
                :description="$t('admin.settings.vector.threshold.description')"
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
            <i class="fas fa-list-ol mr-2 text-brand-600" />{{ $t('admin.settings.vector.maxResults.label') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model.number="localSettings.vector_max_results"
              type="number"
              :placeholder="$t('admin.settings.vector.maxResults.placeholder')"
              width="500px"
            >
              <template #unit>{{ $t('admin.settings.vector.maxResults.unit') }}</template>
            </CyberInput>
            <p class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.settings.vector.maxResults.hint') }}</p>
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
            <i class="fas fa-tasks mr-2 text-brand-600" />{{ $t('admin.settings.vector.concurrency.label') }}
          </label>
          <div class="flex-1">
            <CyberInput
              v-model.number="localSettings.vector_concurrency"
              type="number"
              :placeholder="$t('admin.settings.vector.concurrency.placeholder')"
              width="500px"
              :min="1"
              :max="10"
            >
              <template #unit>{{ $t('admin.settings.vector.concurrency.unit') }}</template>
            </CyberInput>
            <p class="text-content-content-disabled mt-1 text-xs">{{ $t('admin.settings.vector.concurrency.hint') }}</p>
          </div>
        </div>

        <!-- 测试向量配置 -->
        <div class="mt-6 border-t border-subtle pt-5">
          <div class="flex flex-col md:flex-row md:items-center md:justify-between">
            <div class="flex items-center">
              <i class="fas fa-flask mr-2 text-cyan-400" />
              <span class="text-sm text-content">{{ $t('admin.settings.vector.test.title') }}</span>
              <span v-if="!isConfigValid" class="ml-2 text-xs text-amber-400/70">
                ({{ $t('admin.settings.vector.test.needConfig') }})
              </span>
            </div>
            <div class="mt-3 flex gap-3 md:mt-0">
              <CyberButton
                type="outline"
                :loading="testing"
                :show-spinner="false"
                :disabled="testing || !isConfigValid"
                @click="testVectorConfig"
              >
                <template v-if="testing">
                  <i class="fas fa-spinner fa-spin mr-2" />
                  {{ $t('admin.settings.vector.test.testing') }}
                </template>
                <template v-else>
                  <i class="fas fa-vial mr-2" />
                  {{ $t('admin.settings.vector.test.button') }}
                </template>
              </CyberButton>
              <CyberButton
                type="outline"
                :loading="regenerating"
                :show-spinner="false"
                :disabled="regenerating || !isConfigValid"
                @click="regenerateVectors"
              >
                <template v-if="regenerating">
                  <i class="fas fa-spinner fa-spin mr-2" />
                  {{ $t('admin.settings.vector.test.regenerating') }}
                </template>
                <template v-else>
                  <i class="fas fa-sync-alt mr-2" />
                  {{ $t('admin.settings.vector.test.regenerate') }}
                </template>
              </CyberButton>
            </div>
          </div>
        </div>

        <!-- 测试结果 -->
        <div
          v-if="testing || testResult"
          class="mt-4 rounded-lg p-4"
          :class="
            testing
              ? 'border border-blue-500/30 bg-blue-900/20'
              : testResult?.success
                ? 'border border-green-500/30 bg-green-900/20'
                : 'border border-red-500/30 bg-red-900/20'
          "
        >
          <!-- Loading状态 -->
          <div v-if="testing" class="mb-2 flex items-center">
            <i class="fas fa-spinner fa-spin mr-2 text-blue-400" />
            <span class="text-sm font-medium text-content-heading">{{ $t('admin.settings.vector.test.testingStatus') }}</span>
          </div>

          <!-- 测试结果 -->
          <div v-else-if="testResult" class="mb-2 flex items-center">
            <i
              :class="testResult.success ? 'fas fa-check-circle text-green-400' : 'fas fa-times-circle text-red-400'"
              class="mr-2"
            />
            <span class="text-sm font-medium text-content-heading">{{
              testResult.success ? $t('admin.settings.vector.test.valid') : $t('admin.settings.vector.test.invalid')
            }}</span>
          </div>

          <p v-if="testing" class="text-xs text-content-muted">{{ $t('admin.settings.vector.test.verifying') }}</p>
          <p v-else-if="testResult" class="text-xs text-content-muted">
            {{ testResult.message }}
          </p>

          <div v-if="testResult?.details && !testing" class="text-content-content-muted mt-2 text-xs">
            <pre class="whitespace-pre-wrap">{{ testResult.details }}</pre>
          </div>
        </div>

        <!-- 功能说明 -->
        <div
          class="mt-8 overflow-hidden rounded-xl border border-indigo-500/20 bg-gradient-to-br from-indigo-900/10 to-violet-900/10 p-6 backdrop-blur-sm"
        >
          <h4 class="mb-4 flex items-center text-lg font-semibold text-content-heading">
            <div class="mr-3 flex h-8 w-8 items-center justify-center rounded-lg bg-indigo-500/20">
              <i class="fas fa-search text-indigo-400" />
            </div>
            {{ $t('admin.settings.vector.docs.title') }}
          </h4>

          <div class="mb-4 rounded-lg border border-indigo-500/10 bg-indigo-900/10 p-4">
            <div class="mb-2 flex items-center">
              <i class="fas fa-project-diagram mr-2 text-indigo-400" />
              <span class="font-medium text-indigo-300">{{ $t('admin.settings.vector.docs.overview.title') }}</span>
            </div>
            <p class="text-sm leading-relaxed text-content">
              {{ $t('admin.settings.vector.docs.overview.desc1') }}
              {{ $t('admin.settings.vector.docs.overview.desc2') }}
            </p>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <div class="rounded-lg border border-blue-500/20 bg-blue-900/10 p-4">
              <div class="mb-3 flex items-center">
                <i class="fas fa-language mr-2 text-blue-400" />
                <span class="font-medium text-blue-300">{{ $t('admin.settings.vector.docs.naturalSearch.title') }}</span>
              </div>
              <div class="space-y-2 text-sm text-content-muted">
                <p>
                  • <strong class="text-content-heading">{{ $t('admin.settings.vector.docs.naturalSearch.semantic') }}</strong> -
                  {{ $t('admin.settings.vector.docs.naturalSearch.semanticDesc') }}
                </p>
                <p>
                  • <strong class="text-content-heading">{{ $t('admin.settings.vector.docs.naturalSearch.fuzzy') }}</strong> -
                  {{ $t('admin.settings.vector.docs.naturalSearch.fuzzyDesc') }}
                </p>
                <p>
                  • <strong class="text-content-heading">{{ $t('admin.settings.vector.docs.naturalSearch.multilang') }}</strong> -
                  {{ $t('admin.settings.vector.docs.naturalSearch.multilangDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{ $t('admin.settings.vector.docs.naturalSearch.autocomplete') }}</strong>
                  - {{ $t('admin.settings.vector.docs.naturalSearch.autocompleteDesc') }}
                </p>
              </div>
            </div>

            <div class="rounded-lg border border-purple-500/20 bg-purple-900/10 p-4">
              <div class="mb-3 flex items-center">
                <i class="fas fa-images mr-2 text-purple-400" />
                <span class="font-medium text-purple-300">{{ $t('admin.settings.vector.docs.similarSearch.title') }}</span>
              </div>
              <div class="space-y-2 text-sm text-content-muted">
                <p>
                  • <strong class="text-content-heading">{{ $t('admin.settings.vector.docs.similarSearch.byImage') }}</strong> -
                  {{ $t('admin.settings.vector.docs.similarSearch.byImageDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{ $t('admin.settings.vector.docs.similarSearch.styleMatch') }}</strong> -
                  {{ $t('admin.settings.vector.docs.similarSearch.styleMatchDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{
                    $t('admin.settings.vector.docs.similarSearch.contentRelation')
                  }}</strong>
                  - {{ $t('admin.settings.vector.docs.similarSearch.contentRelationDesc') }}
                </p>
                <p>
                  •
                  <strong class="text-content-heading">{{
                    $t('admin.settings.vector.docs.similarSearch.recommendation')
                  }}</strong>
                  - {{ $t('admin.settings.vector.docs.similarSearch.recommendationDesc') }}
                </p>
              </div>
            </div>
          </div>

          <div class="mt-4 grid gap-4 md:grid-cols-2">
            <div class="rounded-lg border border-green-500/20 bg-green-900/10 p-4">
              <div class="mb-3 flex items-center">
                <i class="fas fa-chart-line mr-2 text-green-400" />
                <span class="font-medium text-green-300">{{ $t('admin.settings.vector.docs.performance.title') }}</span>
              </div>
              <div class="space-y-2 text-sm text-content-muted">
                <p>
                  •
                  <strong class="text-green-300">{{ $t('admin.settings.vector.docs.performance.similarityThreshold') }}</strong> -
                  {{ $t('admin.settings.vector.docs.performance.similarityThresholdDesc') }}
                </p>
                <p>
                  • <strong class="text-green-300">{{ $t('admin.settings.vector.docs.performance.concurrency') }}</strong> -
                  {{ $t('admin.settings.vector.docs.performance.concurrencyDesc') }}
                </p>
                <p>
                  • <strong class="text-green-300">{{ $t('admin.settings.vector.docs.performance.timeout') }}</strong> -
                  {{ $t('admin.settings.vector.docs.performance.timeoutDesc') }}
                </p>
              </div>
            </div>

            <div class="rounded-lg border border-amber-500/20 bg-amber-900/10 p-4">
              <div class="mb-3 flex items-center">
                <i class="fas fa-database mr-2 text-amber-400" />
                <span class="font-medium text-amber-300">{{ $t('admin.settings.vector.docs.models.title') }}</span>
              </div>
              <div class="space-y-2 text-sm text-content-muted">
                <p>
                  • <strong class="text-amber-300">text-embedding-3-small</strong> -
                  {{ $t('admin.settings.vector.docs.models.smallDesc') }}
                </p>
                <p>
                  • <strong class="text-amber-300">text-embedding-3-large</strong> -
                  {{ $t('admin.settings.vector.docs.models.largeDesc') }}
                </p>
                <p>
                  • <strong class="text-amber-300">text-embedding-ada-002</strong> -
                  {{ $t('admin.settings.vector.docs.models.adaDesc') }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .slider-container {
    margin-bottom: 6px;
  }
</style>
