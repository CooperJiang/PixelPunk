import { ref, computed } from 'vue'
import type { AISettings, VectorSettings } from './types'

/**
 * AI 和向量设置子模块
 * 管理 AI 和向量搜索相关的配置
 */
export function useAISettingsModule() {
  const aiSettings = ref<AISettings>({})
  const vectorSettings = ref<VectorSettings>({})

  const isAIEnabled = computed(() => aiSettings.value.ai_enabled ?? false)

  const isVectorEnabled = computed(() => vectorSettings.value.vector_enabled ?? false)

  const updateAISettings = (newSettings: AISettings) => {
    aiSettings.value = { ...aiSettings.value, ...newSettings }
  }

  const updateVectorSettings = (newSettings: VectorSettings) => {
    vectorSettings.value = { ...vectorSettings.value, ...newSettings }
  }

  const reset = () => {
    aiSettings.value = {}
    vectorSettings.value = {}
  }

  return {
    aiSettings,
    vectorSettings,

    isAIEnabled,
    isVectorEnabled,

    updateAISettings,
    updateVectorSettings,
    resetAISettings: reset,
  }
}
