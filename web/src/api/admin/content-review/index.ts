import { del, get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/* 审核队列查询参数 */
export interface ReviewQueueQuery {
  page?: number
  size?: number
  sort?: string
  keyword?: string
  nsfw_only?: boolean
}

/* 审核操作参数 */
export interface ReviewAction {
  file_id: string
  action: 'approve' | 'reject'
  reason?: string
  hard_delete?: boolean
}

/* 批量审核参数 */
export interface BatchReviewAction {
  file_ids: string[]
  action: 'approve' | 'reject'
  reason?: string
  hard_delete?: boolean
}

/* 审核记录查询参数 */
export interface ReviewLogQuery {
  page?: number
  size?: number
  action?: 'approve' | 'reject' | ''
  auditor_id?: number
  keyword?: string
  date_from?: string
  date_to?: string
}

/* AI信息接口 - 与 admin/images 保持一致 */
export interface ImageAIInfo {
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

/* 审核队列文件信息 */
export interface ReviewImage {
  id: string
  original_name: string
  display_name: string
  url: string
  thumb_url: string
  full_url: string // 完整签名URL
  full_thumb_url: string // 完整缩略图签名URL
  size_formatted: string
  width: number
  height: number
  format: string
  nsfw: boolean
  nsfw_score?: number // NSFW评分
  created_at: string
  user_id: number
  uploader?: ReviewUser // 上传者信息
  ai_info: ImageAIInfo // 添加 AI 信息
}

/* 审核队列响应 */
export interface ReviewQueueResponse {
  data: ReviewImage[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_page: number
  }
}

/* 审核统计信息 */
export interface ReviewStats {
  pending_count: number
  approved_today: number
  rejected_today: number
}

/* 批量审核结果 */
export interface BatchReviewResult {
  success_count: number
  fail_count: number
  results: string[]
  delete_type?: string
}

/* 用户信息 */
export interface ReviewUser {
  id: number
  username: string
  email: string
}

/* 审核记录文件信息 */
export interface ReviewLogImage {
  id: string
  original_name: string
  display_name?: string
  url?: string
  thumb_url?: string
  full_url?: string
  full_thumb_url?: string
  size_formatted?: string
  width?: number
  height?: number
  format?: string
  status: string
  created_at?: string
}

/* 审核记录 */
export interface ReviewLog {
  id: number
  file_id: string
  action: 'approve' | 'reject'
  delete_type?: string
  reason?: string
  nsfw_score?: number
  nsfw_threshold?: number
  is_nsfw?: boolean
  created_at: string
  updated_at: string
  auditor?: ReviewUser
  uploader?: ReviewUser
  file?: ReviewLogImage
}

/* 审核记录响应 */
export interface ReviewLogResponse {
  data: ReviewLog[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_page: number
  }
}

export function getReviewQueue(params: ReviewQueueQuery): Promise<ApiResult<ReviewQueueResponse>> {
  return get<ReviewQueueResponse>('/admin/content-review/queue', params)
}

export function getReviewLogs(params: ReviewLogQuery): Promise<ApiResult<ReviewLogResponse>> {
  return get<ReviewLogResponse>('/admin/content-review/logs', params)
}

export function getReviewStats(): Promise<ApiResult<ReviewStats>> {
  return get<ReviewStats>('/admin/content-review/stats')
}

export function getFileDetail(fileId: string): Promise<ApiResult<ReviewImage>> {
  return get(`/admin/content-review/files/${fileId}`)
}

export function reviewImage(data: ReviewAction): Promise<ApiResult<void>> {
  return post('/admin/content-review/review', data)
}

export function batchReview(data: BatchReviewAction): Promise<ApiResult<BatchReviewResult>> {
  return post<BatchReviewResult>('/admin/content-review/batch-review', data)
}

export function hardDeleteReviewedImage(fileId: string): Promise<ApiResult<void>> {
  return del(`/admin/content-review/files/${fileId}/hard-delete`)
}

export function batchHardDeleteReviewedImages(fileIds: string[]): Promise<ApiResult<BatchReviewResult>> {
  return post<BatchReviewResult>('/admin/content-review/batch-hard-delete', { file_ids: fileIds })
}

export function restoreReviewedImage(fileId: string): Promise<ApiResult<void>> {
  return post(`/admin/content-review/files/${fileId}/restore`, {})
}

export function batchRestoreReviewedImages(fileIds: string[]): Promise<ApiResult<BatchReviewResult>> {
  return post<BatchReviewResult>('/admin/content-review/batch-restore', { file_ids: fileIds })
}

export default {
  getReviewQueue,
  getReviewLogs,
  getReviewStats,
  getFileDetail,
  reviewImage,
  batchReview,
  hardDeleteReviewedImage,
  batchHardDeleteReviewedImages,
  restoreReviewedImage,
  batchRestoreReviewedImages,
}
