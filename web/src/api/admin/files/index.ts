import { get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/* 文件列表查询参数接口 */
export interface FileListParams {
  page?: number
  size?: number
  sort?: 'newest' | 'oldest' | 'name' | 'size' | 'width' | 'height' | 'quality' | 'nsfw_score'
  keyword?: string
  tags?: string
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

/* AI信息接口 */
export interface FileAIInfo {
  description: string
  tags: string[]
  dominant_color: string
  resolution: string
  is_nsfw: boolean
  nsfw_score: number
  nsfw_evaluation: string
  color_palette?: string[]
  aspect_ratio?: number
  composition?: string
  objects_count?: number
  nsfw_categories?: { [key: string]: number }
}

/* 文件信息接口 */
export interface FileInfo {
  id: string
  url: string
  thumb_url: string
  full_url: string
  full_thumb_url: string
  original_name: string
  display_name: string
  size: number
  width: number
  height: number
  format: string
  access_level: 'public' | 'private' | 'protected'
  folder_id: string
  created_at: string
  updated_at: string
  is_duplicate: boolean
  md5_hash: string
  is_recommended: boolean
  storage_provider_id: string
  user_id: number
  user_name: string
  ai_info: FileAIInfo
}

/* 分页信息接口 */
export interface Pagination {
  total: number
  size: number
  current_page: number
  last_page: number
}

/* 文件列表响应接口 */
export interface FileListResponse {
  items: FileInfo[]
  pagination: Pagination
}

/* 标签信息接口 */
export interface TagInfo {
  id: number
  name: string
  count: number
  created_at: string
  updated_at: string
}

/* 标签列表响应接口 */
export interface TagListResponse {
  items: TagInfo[]
  total: number
  pagination?: Pagination
}

/* 颜色列表响应接口 */
export interface ColorListResponse {
  items: string[]
  total: number
  pagination?: Pagination
}

/* 批量操作结果接口 */
export interface BatchOperationResult {
  success_count: number
  fail_count: number
  success_ids: string[]
  fail_ids: string[]
  errors?: { [key: string]: string }
}

/* 批量推荐操作响应接口 */
export interface BatchRecommendationResponse extends BatchOperationResult {
  updated_files: FileInfo[]
}

/* 批量删除操作响应接口 */
export interface BatchDeleteResponse extends BatchOperationResult {
  deleted_ids: string[]
}

export function getAdminFileList(params?: FileListParams): Promise<ApiResult<FileListResponse>> {
  return get<FileListResponse>('/admin/files/list', params)
}

export function getAdminTagList(page?: number, size?: number): Promise<ApiResult<TagListResponse>> {
  return get<TagListResponse>('/admin/files/tags', { page, size })
}

export function getAdminColorList(page?: number, size?: number): Promise<ApiResult<ColorListResponse>> {
  return get<ColorListResponse>('/admin/files/colors', { page, size })
}

export function setFileRecommendation(file_id: string): Promise<ApiResult<FileInfo>> {
  return post<FileInfo>('/admin/files/recommend', { file_id })
}

export function batchSetFileRecommendation(
  file_ids: string[],
  is_recommended: boolean
): Promise<ApiResult<BatchRecommendationResponse>> {
  return post<BatchRecommendationResponse>('/admin/files/batch-recommend', {
    file_ids,
    is_recommended,
  })
}

export function deleteAdminFile(file_id: string): Promise<ApiResult<{ id: string }>> {
  return post<{ id: string }>('/admin/files/delete', { file_id })
}

export function batchDeleteAdminFiles(file_ids: string[]): Promise<ApiResult<BatchDeleteResponse>> {
  return post<BatchDeleteResponse>('/admin/files/batch-delete', { file_ids })
}

export default {
  getAdminFileList,
  getAdminTagList,
  getAdminColorList,
  setFileRecommendation,
  batchSetFileRecommendation,
  deleteAdminFile,
  batchDeleteAdminFiles,
}
