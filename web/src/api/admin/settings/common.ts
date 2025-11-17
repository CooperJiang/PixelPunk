/**
 * 设置通用 API
 * 包含增删改查等基础操作

 */
import { del, get, post, put } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { Setting, SettingsListResponse, SettingsQuery, GlobalSettingsResponse, LegalDocumentsResponse } from './types'

export function getSettings(params?: SettingsQuery): Promise<ApiResult<SettingsListResponse>> {
  return get<SettingsListResponse>('/settings', params)
}

export function getSetting(key: string): Promise<ApiResult<Setting>> {
  return get<Setting>(`/settings/${key}`)
}

export function createSetting(setting: Setting): Promise<ApiResult<Setting>> {
  return post<Setting>('/settings', setting)
}

export function updateSetting(setting: Setting): Promise<ApiResult<Setting>> {
  return put<Setting>('/settings', setting)
}

export function batchUpsertSettings(settings: Setting[]): Promise<ApiResult<{ success: boolean }>> {
  return post<{ success: boolean }>('/settings/upsert', { settings })
}

export function deleteSetting(key: string): Promise<ApiResult<null>> {
  return del<null>(`/settings/${key}`)
}

export function getGlobalSettings(): Promise<ApiResult<GlobalSettingsResponse>> {
  return get<GlobalSettingsResponse>('/common/settings/global')
}

export function getGlobalSettingsByGroup(
  group: string = 'construction'
): Promise<ApiResult<{ settings: Record<string, unknown>; updated_at: string }>> {
  return get<{ settings: Record<string, unknown>; updated_at: string }>(`/common/settings/global?group=${group}`)
}

export function getLegalDocuments(): Promise<ApiResult<LegalDocumentsResponse>> {
  return get<LegalDocumentsResponse>('/common/settings/legal')
}
