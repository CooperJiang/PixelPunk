import { get, post } from '@/utils/network/http'

export interface TaggingStatusParams {
  status?: string
  page: number
  limit: number
  order_by?: string
  order?: 'asc' | 'desc'
}

export interface ImageTaggingInfo {
  id: string
  original_name: string
  size_formatted: string
  format: string
  user_id: number
  ai_tagging_status: string
  ai_tagging_tries: number
  ai_tagging_duration_formatted?: string
  ai_http_duration_formatted?: string
  updated_at: string
  width?: number
  height?: number
  size?: number
  full_thumb_url?: string
  is_nsfw?: boolean
}

export interface TaggingStatusResponse {
  total: number
  page: number
  limit: number
  data: ImageTaggingInfo[]
}

export interface StatusStatsItem {
  status: string
  count: number
}

export interface TaggingLogItem {
  id: number
  file_id: string
  file_type: string
  status: string
  action: string
  type: string
  data: Record<string, any>
  message?: string
  operator_id: number
  duration: number
  created_at: string
  updated_at: string
}

export interface TaggingStatsResponse {
  status_stats: StatusStatsItem[]
  recent_logs: TaggingLogItem[]
}

export interface QueueStatsResponse {
  queue_stats: Record<string, number>
  config: Record<string, number>
  performance: Record<string, number>
}

export interface ApiResult<T> {
  code: number
  message: string
  data: T
}

export function getTaggingStatusList(params: TaggingStatusParams): Promise<ApiResult<TaggingStatusResponse>> {
  return get<TaggingStatusResponse>('/admin/ai/tagging/status', params)
}

export function getTaggingStats(): Promise<ApiResult<TaggingStatsResponse>> {
  return get<TaggingStatsResponse>('/admin/ai/tagging/stats')
}

export interface TaggingLogsParams {
  file_id?: string
  page: number
  limit: number
}

export interface TaggingLogsResponse {
  logs: TaggingLogItem[]
  pagination: {
    page: number
    limit: number
    total: number
  }
}

export function getTaggingLogs(params: TaggingLogsParams): Promise<ApiResult<TaggingLogsResponse>> {
  return get<TaggingLogsResponse>('/admin/ai/tagging/logs', params)
}

export function getQueueStats(): Promise<ApiResult<QueueStatsResponse>> {
  return get<QueueStatsResponse>('/admin/ai/tagging/queue-stats')
}

export function retryTagging(file_ids: string[]): Promise<ApiResult<{ requested: number; enqueued: number; skipped: number }>> {
  return post<{ requested: number; enqueued: number; skipped: number }>('/admin/ai/tagging/retry', { file_ids })
}

export function triggerTagging(max_images?: number): Promise<ApiResult<{ submitted_count: number; status: string }>> {
  return post<{ submitted_count: number; status: string }>('/admin/ai/tagging/trigger', { max_images })
}

export function resetStuckTagging(time_threshold_minutes: number): Promise<ApiResult<{ reset_count: number }>> {
  return post<{ reset_count: number }>('/admin/ai/tagging/reset-stuck', { time_threshold_minutes })
}

export function ignoreTagging(file_ids: string[], reason?: string): Promise<ApiResult<{ count: number; ignored_ids: string[] }>> {
  return post<{ count: number; ignored_ids: string[] }>('/admin/ai/tagging/ignore', { file_ids, reason })
}

export function unignoreTagging(
  file_ids: string[]
): Promise<ApiResult<{ count: number; updated_ids: string[]; enqueued: number; skipped: number }>> {
  return post<{ count: number; updated_ids: string[]; enqueued: number; skipped: number }>('/admin/ai/tagging/unignore', {
    file_ids,
  })
}

export function retryFailedAll(limit?: number): Promise<ApiResult<{ selected: number; enqueued: number; skipped: number }>> {
  return post<{ selected: number; enqueued: number; skipped: number }>('/admin/ai/tagging/retry-failed-all', { limit })
}

export function getAutoProcessingStatus(): Promise<ApiResult<{ enabled: boolean; recent_failures: number }>> {
  return get<{ enabled: boolean; recent_failures: number }>('/admin/ai/tagging/auto-processing')
}

export function setAutoProcessing(enabled: boolean): Promise<ApiResult<{ enabled: boolean }>> {
  return post<{ enabled: boolean }>('/admin/ai/tagging/auto-processing', { enabled })
}

export function setConcurrency(concurrency: number): Promise<ApiResult<{ configured: number }>> {
  return post<{ configured: number }>('/admin/ai/tagging/concurrency', { concurrency })
}

export function getTaggingDiagnosis(): Promise<
  ApiResult<{
    queue_initialized: boolean
    service_running: boolean
    paused: boolean
    config_issues: string[]
    recommendations: string[]
    worker_count: number
    active_workers: number
  }>
> {
  return get('/admin/ai/tagging/diagnosis')
}

export default {
  getTaggingStatusList,
  getTaggingStats,
  getQueueStats,
  getTaggingDiagnosis,
  retryTagging,
  triggerTagging,
  resetStuckTagging,
  ignoreTagging,
  unignoreTagging,
  retryFailedAll,
  getAutoProcessingStatus,
  setAutoProcessing,
  setConcurrency,
}
