import { del, get, post, put } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type {
  ApiKeyInfo,
  ApiKeyListParams,
  ApiKeyListResponse,
  ApiKeyStatsResponse,
  CreateApiKeyRequest,
  RegenerateApiKeyResponse,
  ToggleApiKeyResponse,
  UpdateApiKeyRequest,
} from '../types'

export function createApiKey(data: CreateApiKeyRequest): Promise<ApiResult<ApiKeyInfo>> {
  return post<ApiKeyInfo>('/apikey/create', data)
}

export function getApiKeyList(params?: ApiKeyListParams): Promise<ApiResult<ApiKeyListResponse>> {
  return get<ApiKeyListResponse>('/apikey/list', params)
}

export function getApiKeyDetail(keyId: string): Promise<ApiResult<ApiKeyInfo>> {
  return get<ApiKeyInfo>(`/apikey/${keyId}`)
}

export function updateApiKey(keyId: string, data: UpdateApiKeyRequest): Promise<ApiResult<ApiKeyInfo>> {
  return put<ApiKeyInfo>(`/apikey/${keyId}`, data)
}

export function deleteApiKey(keyId: string): Promise<ApiResult<{ id: string }>> {
  return del<{ id: string }>(`/apikey/${keyId}`)
}

export function toggleApiKey(keyId: string): Promise<ApiResult<ToggleApiKeyResponse>> {
  return post<ToggleApiKeyResponse>(`/apikey/${keyId}/toggle`)
}

export function getApiKeyStats(keyId: string): Promise<ApiResult<ApiKeyStatsResponse>> {
  return get<ApiKeyStatsResponse>(`/apikey/${keyId}/stats`)
}

export function regenerateApiKey(keyId: string): Promise<ApiResult<RegenerateApiKeyResponse>> {
  return post<RegenerateApiKeyResponse>(`/apikey/${keyId}/regenerate`)
}

export default {
  createApiKey,
  getApiKeyList,
  getApiKeyDetail,
  updateApiKey,
  deleteApiKey,
  toggleApiKey,
  getApiKeyStats,
  regenerateApiKey,
}
