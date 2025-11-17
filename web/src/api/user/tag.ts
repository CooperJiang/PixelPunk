/**
 * 用户标签管理 API

 */
import { post, get } from '@/utils/network/http'
import type { ApiResponse } from '@/api/types/common'

/**
 * 用户标签信息
 */
export interface UserTag {
  id: number
  name: string
  slug: string
  description: string
  is_system: boolean
  creator_id: number
  usage_count: number
  sort_order: number
  created_at: string
  updated_at: string
  source: 'manual' | 'ai' | 'system'
  reference_created_at: string
  file_count: number
}

/**
 * 创建标签请求
 */
export interface CreateTagRequest {
  name: string
}

/**
 * 更新标签请求
 */
export interface UpdateTagRequest {
  id: number
  name: string
}

/**
 * 删除标签请求
 */
export interface DeleteTagRequest {
  id: number
}

/**
 * 批量删除请求
 */
export interface BatchDeleteTagRequest {
  ids: number[]
}

/**
 * 合并标签请求
 */
export interface MergeTagsRequest {
  source_ids: number[]
  target_id: number
}

/**
 * 标签列表查询参数
 */
export interface TagListQuery {
  keyword?: string
  page?: number
  size?: number
  sort_by?: 'name' | 'file_count' | 'created_at'
  sort_order?: 'asc' | 'desc'
}

/**
 * 标签列表响应
 */
export interface TagListResponse {
  items: UserTag[]
  total: number
  pagination: {
    page: number
    size: number
  }
}

/**
 * 标签统计信息
 */
export interface TagStats {
  total_tags: number
  user_created_tags: number
  ai_created_tags: number
  total_usage: number
}

/**
 */
export function createTag(data: CreateTagRequest): Promise<ApiResponse<UserTag>> {
  return post<UserTag>('/user-tags/create', data)
}

export function getTagList(params: TagListQuery): Promise<ApiResponse<TagListResponse>> {
  return get<TagListResponse>('/user-tags/list', params)
}

export function getAllTags(): Promise<ApiResponse<UserTag[]>> {
  return get<UserTag[]>('/user-tags/all')
}

export function updateTag(data: UpdateTagRequest): Promise<ApiResponse<UserTag>> {
  return post<UserTag>('/user-tags/update', data)
}

export function deleteTag(data: DeleteTagRequest): Promise<ApiResponse<void>> {
  return post<void>('/user-tags/delete', data)
}

export function batchDeleteTags(data: BatchDeleteTagRequest): Promise<ApiResponse<{ deleted_count: number }>> {
  return post<{ deleted_count: number }>('/user-tags/batch-delete', data)
}

export function mergeTags(data: MergeTagsRequest): Promise<ApiResponse<void>> {
  return post<void>('/user-tags/merge', data)
}

export function getTagStats(): Promise<ApiResponse<TagStats>> {
  return get<TagStats>('/user-tags/stats')
}

export default {
  createTag,
  getTagList,
  getAllTags,
  updateTag,
  deleteTag,
  batchDeleteTags,
  mergeTags,
  getTagStats,
}
