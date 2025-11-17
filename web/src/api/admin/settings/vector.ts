/**
 * 向量搜索设置相关 API

 */
import type { TranslationFunction } from '@/composables/useTexts'
import { post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type {
  VectorTestParams,
  VectorTestResult,
  RegenerateVectorsParams,
  RegenerateVectorsResult,
  QdrantTestParams,
  QdrantTestResult,
  Setting,
} from './types'

/* 测试向量配置 */
export function testVectorConfig(params: VectorTestParams): Promise<ApiResult<VectorTestResult>> {
  return post<VectorTestResult>('/settings/vector/test', params)
}

export function testQdrantConnection(params: QdrantTestParams): Promise<ApiResult<QdrantTestResult>> {
  return post<QdrantTestResult>('/settings/vector/test-qdrant', params)
}

export function regenerateVectors(params?: RegenerateVectorsParams): Promise<ApiResult<RegenerateVectorsResult>> {
  return post<RegenerateVectorsResult>('/search/vector/regenerate', null, { params })
}

/* 获取向量设置（带翻译） */
export function getVectorDefaults($t: TranslationFunction): Setting[] {
  return [
    {
      key: 'vector_enabled',
      value: false,
      type: 'boolean',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_enabled'),
      is_system: true,
    },
    {
      key: 'vector_provider',
      value: 'openai',
      type: 'string',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_provider'),
      is_system: true,
    },
    {
      key: 'vector_model',
      value: 'text-embedding-3-small',
      type: 'string',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_model'),
      is_system: true,
    },
    {
      key: 'vector_api_key',
      value: '',
      type: 'string',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_api_key'),
      is_system: true,
    },
    {
      key: 'vector_base_url',
      value: 'https://api.openai.com/v1',
      type: 'string',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_base_url'),
      is_system: true,
    },
    {
      key: 'vector_timeout',
      value: 30,
      type: 'number',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_timeout'),
      is_system: true,
    },
    {
      key: 'vector_similarity_threshold',
      value: 0.3,
      type: 'number',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_similarity_threshold'),
      is_system: true,
    },
    {
      key: 'vector_max_results',
      value: 50,
      type: 'number',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_max_results'),
      is_system: true,
    },
    {
      key: 'vector_concurrency',
      value: 3,
      type: 'number',
      group: 'vector',
      description: $t('api.settingsDefaults.vector.vector_concurrency'),
      is_system: true,
    },
  ]
}

/* 向后兼容 - 默认中文 */
export const vectorDefaults: Setting[] = [
  {
    key: 'vector_enabled',
    value: false,
    type: 'boolean',
    group: 'vector',
    description: '启用向量搜索功能',
    is_system: true,
  },
  {
    key: 'vector_provider',
    value: 'openai',
    type: 'string',
    group: 'vector',
    description: '向量化提供者',
    is_system: true,
  },
  {
    key: 'vector_model',
    value: 'text-embedding-3-small',
    type: 'string',
    group: 'vector',
    description: '向量化模型',
    is_system: true,
  },
  {
    key: 'vector_api_key',
    value: '',
    type: 'string',
    group: 'vector',
    description: 'OpenAI API密钥',
    is_system: true,
  },
  {
    key: 'vector_base_url',
    value: 'https://api.openai.com/v1',
    type: 'string',
    group: 'vector',
    description: 'OpenAI API代理地址',
    is_system: true,
  },
  {
    key: 'vector_timeout',
    value: 30,
    type: 'number',
    group: 'vector',
    description: 'API调用超时时间(秒)',
    is_system: true,
  },
  {
    key: 'vector_similarity_threshold',
    value: 0.3,
    type: 'number',
    group: 'vector',
    description: '相似度阈值(0-1)',
    is_system: true,
  },
  {
    key: 'vector_max_results',
    value: 50,
    type: 'number',
    group: 'vector',
    description: '最大搜索结果数',
    is_system: true,
  },
  {
    key: 'vector_concurrency',
    value: 3,
    type: 'number',
    group: 'vector',
    description: '向量生成并发数量',
    is_system: true,
  },
]
