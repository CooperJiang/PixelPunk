import { del, get, post, put } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type {
  CreateStorageChannelRequest,
  StorageChannel,
  StorageChannelExport,
  StorageChannelsBatchExport,
  StorageConfigItem,
  UpdateStorageChannelRequest,
  UpdateStorageConfigRequest,
} from '../types'

const baseUrl = '/storage'

export function getChannelList(): Promise<ApiResult<StorageChannel[]>> {
  return get<StorageChannel[]>(`${baseUrl}/list`)
}

export function getChannelDetail(id: string): Promise<ApiResult<StorageChannel>> {
  return get<StorageChannel>(`${baseUrl}/${id}`)
}

export function createChannel(data: CreateStorageChannelRequest): Promise<ApiResult<StorageChannel>> {
  return post<StorageChannel>(`${baseUrl}/create`, data)
}

export function updateChannel(id: string, data: UpdateStorageChannelRequest): Promise<ApiResult<StorageChannel>> {
  return put<StorageChannel>(`${baseUrl}/${id}`, data)
}

export function deleteChannel(id: string): Promise<ApiResult<null>> {
  return del<null>(`${baseUrl}/${id}`)
}

export function getChannelConfigs(id: string): Promise<ApiResult<StorageConfigItem[]>> {
  return get<StorageConfigItem[]>(`${baseUrl}/${id}/configs`)
}

export function updateChannelConfigs(id: string, data: UpdateStorageConfigRequest): Promise<ApiResult<null>> {
  return put<null>(`${baseUrl}/${id}/configs`, data)
}

export function testChannel(id: string): Promise<ApiResult<null>> {
  return post<null>(`${baseUrl}/${id}/test`)
}

export function setDefaultChannel(id: string): Promise<ApiResult<null>> {
  return post<null>(`${baseUrl}/${id}/default`)
}

export function enableChannel(id: string): Promise<ApiResult<null>> {
  return post<null>(`${baseUrl}/${id}/enable`)
}

export function disableChannel(id: string): Promise<ApiResult<null>> {
  return post<null>(`${baseUrl}/${id}/disable`)
}

export function exportChannel(id: string): Promise<ApiResult<StorageChannelExport>> {
  return get<StorageChannelExport>(`${baseUrl}/${id}/export`)
}

export function exportAllChannels(): Promise<ApiResult<StorageChannelsBatchExport>> {
  return get<StorageChannelsBatchExport>(`${baseUrl}/export/all`)
}

export function importChannel(file: File): Promise<ApiResult<null>> {
  const formData = new FormData()
  formData.append('file', file)
  return post<null>(`${baseUrl}/import`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

export function getSupportedTypes(): Promise<ApiResult<string[]>> {
  return get<string[]>(`${baseUrl}/supported-types`)
}

export interface StorageConfigTemplate {
  name: string
  key_name: string
  type: string // 'string' | 'int' | 'bool' | 'password'
  is_secret?: boolean
  required?: boolean
  description?: string
  options?: string[]
  default?: string
}

export function getConfigTemplates(type: string): Promise<ApiResult<StorageConfigTemplate[]>> {
  return get<StorageConfigTemplate[]>(`${baseUrl}/config-templates/${type}`)
}

export default {
  getChannelList,
  getChannelDetail,
  createChannel,
  updateChannel,
  deleteChannel,
  getChannelConfigs,
  updateChannelConfigs,
  testChannel,
  setDefaultChannel,
  enableChannel,
  disableChannel,
  exportChannel,
  exportAllChannels,
  importChannel,
  getSupportedTypes,
  getConfigTemplates,
}
