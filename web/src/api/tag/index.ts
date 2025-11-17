import { del, get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { BatchTagOperation, CreateTagRequest, TagInfo, TagListParams, TagListResponse, UpdateTagRequest } from '../types'

export function getTagList(params: TagListParams = {}): Promise<ApiResult<TagListResponse>> {
  return get<TagListResponse>('/tags/list', params)
}

export function createTag(data: CreateTagRequest): Promise<ApiResult<TagInfo>> {
  return post<TagInfo>('/tags/admin/create', data)
}

export function updateTag(data: UpdateTagRequest): Promise<ApiResult<TagInfo>> {
  return post<TagInfo>('/tags/admin/update', data)
}

export function deleteTag(tagId: number): Promise<ApiResult<{ id: number }>> {
  return del<{ id: number }>(`/tags/admin/${tagId}`)
}

export function batchOperateTags(data: BatchTagOperation): Promise<ApiResult<{ success: boolean }>> {
  return post<{ success: boolean }>('/tags/admin/batch', data)
}

export default {
  getTagList,
  createTag,
  updateTag,
  deleteTag,
  batchOperateTags,
}
