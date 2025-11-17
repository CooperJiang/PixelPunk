import { get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

export interface VectorListParams {
  page: number
  page_size: number
  status?: 'pending' | 'processing' | 'completed' | 'failed' | 'reset'
  model?: string
  keyword?: string
}

export interface VectorStatsInfo {
  min: number
  max: number
  mean: number
  stddev: number
}

export interface VectorItem {
  id: number
  file_id: string
  status: string
  model: string
  description: string
  dimension: number
  retry_count: number
  error_message?: string
  last_retry_at?: string
  processing_duration: number
  created_at: string
  updated_at: string
  file_name?: string
  file_url?: string
  thumbnail_url?: string
  vector_preview?: number[]
  vector_norm?: number
  vector_stats?: VectorStatsInfo
}

export interface PaginationInfo {
  page: number
  page_size: number
  total: number
  total_page: number
}

export interface VectorStats {
  total_count: number
  pending_count: number
  running_count: number
  completed_count: number
  failed_count: number
  reset_count: number
  active_workers: number
  max_workers: number
}

export interface VectorListResponse {
  data: VectorItem[]
  pagination: PaginationInfo
}

export interface VectorBatchActionParams {
  file_ids: string[]
  action: 'reset' | 'retry' | 'delete'
}

export interface VectorRetryParams {
  file_id: string
}

export interface VectorLogParams {
  page: number
  page_size: number
  file_id?: string
  action?: string
  model?: string
}

export interface VectorLogItem {
  id: number
  file_id: string
  action: string
  type: string
  data: Record<string, any>
  message?: string
  model: string
  duration: number
  status?: string
  error_code?: string
  error_message?: string
  task_id?: string
  metadata?: any
  processing_time?: number
  queue_time?: number
  created_at: string
}

export interface VectorLogResponse {
  data: VectorLogItem[]
  pagination: PaginationInfo
}

export interface VectorRegenerateAllParams {
  force?: boolean
}

export interface AvailableModelsResponse {
  models: string[]
}

export interface QdrantRealStats {
  qdrant_vector_count: number
  qdrant_indexed_count: number
  mysql_total_count: number
  mysql_completed_count: number
  sync_ratio: number
  collection_status: string
  is_healthy: boolean
  last_checked: string
}

export function getVectorList(params: VectorListParams): Promise<ApiResult<VectorListResponse>> {
  const result = get<VectorListResponse>('/admin/vector/list', params)
  return result
}

export function getVectorStats(): Promise<ApiResult<VectorStats>> {
  return get<VectorStats>('/admin/vector/stats')
}

export function getAutoProcessing(): Promise<ApiResult<{ enabled: boolean }>> {
  return get<{ enabled: boolean }>('/admin/vector/auto-processing')
}

export function setAutoProcessing(enabled: boolean): Promise<ApiResult<{ enabled: boolean }>> {
  return post('/admin/vector/auto-processing', { enabled })
}

export function setConcurrency(concurrency: number): Promise<ApiResult<{ configured: number }>> {
  return post('/admin/vector/concurrency', { concurrency })
}

export function getQdrantRealStats(): Promise<ApiResult<QdrantRealStats>> {
  return get<QdrantRealStats>('/admin/vector/qdrant-stats')
}

export function getAvailableModels(): Promise<ApiResult<AvailableModelsResponse>> {
  return get<AvailableModelsResponse>('/admin/vector/models')
}

export function getVectorDetail(fileId: string): Promise<ApiResult<VectorItem>> {
  return get<VectorItem>(`/admin/vector/detail/${fileId}`)
}

export function batchVectorAction(params: VectorBatchActionParams): Promise<ApiResult<{ affected_count: number }>> {
  return post('/admin/vector/batch', params)
}

export function retryVector(params: VectorRetryParams): Promise<ApiResult<{ success: boolean }>> {
  return post('/admin/vector/retry', params)
}

export function regenerateAllVectors(params: VectorRegenerateAllParams): Promise<ApiResult<{ affected_count: number }>> {
  return post('/admin/vector/regenerate-all', params)
}

export function retryAllFailedVectors(): Promise<ApiResult<{ affected_count: number }>> {
  return post('/admin/vector/retry-all-failed', {})
}

export function recoverStuckTasks(): Promise<ApiResult<{ affected_count: number }>> {
  return post('/admin/vector/recover-stuck', {})
}

export function reconcileMissing(params: {
  limit?: number
  dry_run?: boolean
}): Promise<ApiResult<{ found: number; enqueued: number; dry_run: boolean }>> {
  return post('/admin/vector/reconcile/missing', params)
}

export function cleanOrphans(params: {
  limit?: number
  dry_run?: boolean
}): Promise<ApiResult<{ found: number; removed: number; dry_run: boolean }>> {
  return post('/admin/vector/reconcile/orphans', params)
}

export function rebuildStale(params: { limit?: number }): Promise<ApiResult<{ enqueued: number }>> {
  return post('/admin/vector/rebuild/stale', params)
}

export function getVectorLogs(params: VectorLogParams): Promise<ApiResult<VectorLogResponse>> {
  return get<VectorLogResponse>('/admin/vector/logs', params)
}

export default {
  getVectorList,
  getVectorStats,
  getAutoProcessing,
  setAutoProcessing,
  setConcurrency,
  getQdrantRealStats,
  getAvailableModels,
  getVectorDetail,
  batchVectorAction,
  retryVector,
  regenerateAllVectors,
  retryAllFailedVectors,
  recoverStuckTasks,
  reconcileMissing,
  cleanOrphans,
  rebuildStale,
  getVectorLogs,
}
