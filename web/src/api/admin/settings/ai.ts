/**
 * AI 设置相关 API

 */
import type { TranslationFunction } from '@/composables/useTexts'
import { post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { AITestResponse, Setting } from './types'

/* 测试AI配置 */
export function testAIConfig(config?: Record<string, unknown>): Promise<ApiResult<AITestResponse>> {
  if (config) {
    return post<AITestResponse>('/admin/ai/test-config', config)
  }
  return post<AITestResponse>('/admin/ai/test-config')
}

export function checkAIConfiguration(settings: Setting[]): {
  isConfigured: boolean
  missingFields: string[]
  hasApiKey: boolean
  hasModel: boolean
} {
  if (!settings || settings.length === 0) {
    return {
      isConfigured: false,
      missingFields: ['ai_api_key', 'ai_model'],
      hasApiKey: false,
      hasModel: false,
    }
  }

  const settingsMap = settings.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  const _requiredFields = ['ai_api_key', 'ai_model']
  const missingFields: string[] = []

  const hasApiKey = Boolean(settingsMap.ai_api_key && settingsMap.ai_api_key.toString().trim())
  if (!hasApiKey) {
    missingFields.push('ai_api_key')
  }

  const hasModel = Boolean(settingsMap.ai_model && settingsMap.ai_model.toString().trim())
  if (!hasModel) {
    missingFields.push('ai_model')
  }

  const isConfigured = missingFields.length === 0

  return {
    isConfigured,
    missingFields,
    hasApiKey,
    hasModel,
  }
}

/* 获取AI设置（带翻译） */
export function getAiDefaults($t: TranslationFunction): Setting[] {
  return [
    {
      key: 'ai_enabled',
      value: false,
      type: 'boolean',
      group: 'ai',
      description: $t('api.settingsDefaults.ai.ai_enabled'),
      is_system: true,
    },
    {
      key: 'ai_proxy',
      value: '',
      type: 'string',
      group: 'ai',
      description: $t('api.settingsDefaults.ai.ai_proxy'),
      is_system: true,
    },
    {
      key: 'ai_model',
      value: 'gpt-3.5-turbo',
      type: 'string',
      group: 'ai',
      description: $t('api.settingsDefaults.ai.ai_model'),
      is_system: true,
    },
    {
      key: 'ai_api_key',
      value: '',
      type: 'string',
      group: 'ai',
      description: $t('api.settingsDefaults.ai.ai_api_key'),
      is_system: true,
    },
    {
      key: 'ai_temperature',
      value: 0.7,
      type: 'number',
      group: 'ai',
      description: $t('api.settingsDefaults.ai.ai_temperature'),
      is_system: true,
    },
    {
      key: 'ai_max_tokens',
      value: 2000,
      type: 'number',
      group: 'ai',
      description: $t('api.settingsDefaults.ai.ai_max_tokens'),
      is_system: true,
    },
    {
      key: 'ai_concurrency',
      value: 5,
      type: 'number',
      group: 'ai',
      description: $t('api.settingsDefaults.ai.ai_concurrency'),
      is_system: true,
    },
  ]
}

/* 向后兼容 - 默认中文 */
export const aiDefaults: Setting[] = [
  {
    key: 'ai_enabled',
    value: false,
    type: 'boolean',
    group: 'ai',
    description: '启用AI功能',
    is_system: true,
  },
  {
    key: 'ai_proxy',
    value: '',
    type: 'string',
    group: 'ai',
    description: '代理地址',
    is_system: true,
  },
  {
    key: 'ai_model',
    value: 'gpt-3.5-turbo',
    type: 'string',
    group: 'ai',
    description: '模型选择',
    is_system: true,
  },
  {
    key: 'ai_api_key',
    value: '',
    type: 'string',
    group: 'ai',
    description: 'API密钥',
    is_system: true,
  },
  {
    key: 'ai_temperature',
    value: 0.7,
    type: 'number',
    group: 'ai',
    description: '温度值',
    is_system: true,
  },
  {
    key: 'ai_max_tokens',
    value: 2000,
    type: 'number',
    group: 'ai',
    description: '最大生成令牌数',
    is_system: true,
  },
  {
    key: 'ai_concurrency',
    value: 5,
    type: 'number',
    group: 'ai',
    description: 'AI文件识别并发处理数量',
    is_system: true,
  },
]
