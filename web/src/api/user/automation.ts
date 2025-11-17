/**
 * 用户自动任务 API
 * 用于查询当前用户的 AI 打标和向量化任务状态
 */
import { get } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/**
 * 打标任务统计信息
 */
export interface TaggingTaskStats {
  none_count: number
  pending_count: number
  processing_count: number
  done_count: number
  failed_count: number
  ignored_count: number
  total_count: number
  queue_position?: number
}

/**
 * 向量化任务统计信息
 */
export interface VectorTaskStats {
  pending_count: number
  processing_count: number
  completed_count: number
  failed_count: number
  reset_count: number
  total_count: number
  queue_position?: number
}

/**
 * 自动任务总览信息
 */
export interface AutomationOverview {
  tagging: TaggingTaskStats
  vector: VectorTaskStats
  system_status: {
    tagging_enabled: boolean
    vector_enabled: boolean
    tagging_concurrency: number
    vector_concurrency: number
    total_queue_length: number
  }
}

/**
 * 打标任务详情项
 */
export interface TaggingTaskItem {
  id: string
  file_id: string
  file_name: string
  thumbnail_url?: string
  status: 'none' | 'pending' | 'processing' | 'done' | 'failed' | 'ignored'
  tries: number
  error_message?: string
  format: string
  size_formatted: string
  resolution: string
  updated_at: string
  created_at: string
}

/**
 * 向量任务详情项
 */
export interface VectorTaskItem {
  id: number
  file_id: string
  file_name: string
  thumbnail_url?: string
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'reset'
  model: string
  retry_count: number
  error_message?: string
  updated_at: string
  created_at: string
}

/**
 * 获取用户自动任务总览
 */
export function getUserAutomationOverview(): Promise<ApiResult<AutomationOverview>> {
  return get<AutomationOverview>('/user/automation/overview')
}

export function getUserTaggingTasks(params?: { status?: string; page?: number; limit?: number }): Promise<
  ApiResult<{
    data: TaggingTaskItem[]
    total: number
    page: number
    limit: number
  }>
> {
  return get('/user/automation/tagging/tasks', params)
}

export function getUserVectorTasks(params?: { status?: string; page?: number; limit?: number }): Promise<
  ApiResult<{
    data: VectorTaskItem[]
    total: number
    page: number
    limit: number
  }>
> {
  return get('/user/automation/vector/tasks', params)
}

export default {
  getUserAutomationOverview,
  getUserTaggingTasks,
  getUserVectorTasks,
}
