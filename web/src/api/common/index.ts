import { get } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { FileInfo } from '@/api/types/file'

/* 文件列表查询参数接口 */
export interface FileListParams {
  page?: number
  size?: number
  sort?: 'newest' | 'oldest' | 'name' | 'size' | 'width' | 'height' | 'quality' | 'nsfw_score'
  keyword?: string
  tags?: string
  categoryId?: string // 分类ID，支持逗号分隔的多个ID
  dominant_color?: string
  dominant_colors?: string[]
  resolution?: '720p' | '1080p' | '1440p' | '2160p' | '4k' | '8k' | string
  nsfw_min_score?: number
  is_nsfw?: boolean
  storage_type?: string
  min_width?: number
  max_width?: number
  min_height?: number
  max_height?: number
  user_id?: number
  is_recommended?: boolean
}

export function getCommonTagList(page?: number, size?: number): Promise<ApiResult<{ tags: string[]; total: number }>> {
  return get<{ tags: string[]; total: number }>('/common/files/tags', { page, size })
}

export function getCommonColorList(page?: number, size?: number): Promise<ApiResult<{ colors: string[]; total: number }>> {
  return get<{ colors: string[]; total: number }>('/common/files/colors', { page, size })
}

export function getUserTagList(page?: number, size?: number): Promise<ApiResult<{ tags: string[]; total: number }>> {
  return get<{ tags: string[]; total: number }>('/common/user/tags', { page, size })
}

export function getUserColorList(page?: number, size?: number): Promise<ApiResult<{ colors: string[]; total: number }>> {
  return get<{ colors: string[]; total: number }>('/common/user/colors', { page, size })
}

export function getUserCategoryList(): Promise<
  ApiResult<{ items: Array<{ id: number; name: string; file_count: number }>; total: number }>
> {
  return get<{ items: Array<{ id: number; name: string; file_count: number }>; total: number }>('/common/user/categories')
}

export function getGuestTagList(page?: number, size?: number): Promise<ApiResult<{ tags: string[]; total: number }>> {
  return get<{ tags: string[]; total: number }>('/common/guest/tags', { page, size })
}

export function getGuestColorList(page?: number, size?: number): Promise<ApiResult<{ colors: string[]; total: number }>> {
  return get<{ colors: string[]; total: number }>('/common/guest/colors', { page, size })
}

export function getGuestFileList(
  params?: FileListParams
): Promise<ApiResult<{ items: FileInfo[]; pagination: { total: number; page: number; size: number } }>> {
  return get<{ items: FileInfo[]; pagination: { total: number; page: number; size: number } }>('/files/guest/list', params)
}

export interface PublicStatsResponse {
  total: number
  timestamp?: number
  storage?: {
    total_storage: number
    formatted_storage: string
  }
}

export function getPbData(params?: ImageListParams): Promise<ApiResult<PublicStatsResponse>> {
  return get<PublicStatsResponse>('/common/pb/data', params)
}

export function getHealthCheck(): Promise<ApiResult<{ status: string; timestamp: number }>> {
  return get<{ status: string; timestamp: number }>('/common/health-check')
}

export type { FileInfo } from '@/api/types/file'

export default {
  getCommonTagList,
  getCommonColorList,
  getGuestFileList,
}
