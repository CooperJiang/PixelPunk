import { del, get, post } from '@/utils/network/http'
import type { ApiResponse, PaginationInfo } from '@/api/types/index'

/* 向量验证任务类型 */
export interface VectorVerificationTask {
  id: number
  task_id: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  task_type: string
  total_count: number
  processed_count: number
  verified_count: number
  missing_count: number
  error_count: number
  batch_size: number
  created_at: string
  updated_at: string
  started_at?: string
  completed_at?: string
  creator_id?: number
  error_message?: string
  filter_conditions?: Record<string, unknown>
}

/* 向量验证统计信息 */
export interface VectorVerificationStats {
  total_vectors: number
  verified_count: number
  missing_count: number
  unknown_count: number
  last_verification_time?: string
  running_task?: {
    task_id: string
    progress: number
  }
}

/* 启动验证任务请求 */
export interface StartVerificationRequest {
  task_type: 'manual' | 'scheduled' | 'partial'
  actual_status?: 'unknown' | 'verified' | 'missing'
  needs_verification?: boolean
  file_ids?: string[]
}

/* 修复缺失向量请求 */
export interface RepairMissingVectorRequest {
  file_id: string
}

/* 批量修复问题向量请求（不再需要参数） */
export interface BatchRepairMissingVectorsRequest {
  [key: string]: never
}

/* 智能修复请求 */
export interface SmartFixRequest {
  [key: string]: never
}

/* 任务列表参数 */
export interface TaskListParams {
  page?: number
  size?: number
  task_type?: string
}

/* 任务列表响应 */
export interface TaskListResponse {
  items: VectorVerificationTask[]
  pagination: PaginationInfo
}

/* 任务进度响应 */
export interface TaskProgressResponse {
  is_running: boolean
  progress: number
  task?: VectorVerificationTask
}

/* 启动向量验证任务 */
export function startVectorVerification(data: StartVerificationRequest): Promise<ApiResponse<VectorVerificationTask>> {
  return post<VectorVerificationTask>('/admin/vector-verification/start', data)
}

export function getTaskStatus(taskId: string): Promise<ApiResponse<VectorVerificationTask>> {
  return get<VectorVerificationTask>(`/admin/vector-verification/task/${taskId}`)
}

export function getTaskList(params?: TaskListParams): Promise<ApiResponse<TaskListResponse>> {
  return get<TaskListResponse>('/admin/vector-verification/tasks', params)
}

export function getVerificationStatistics(): Promise<ApiResponse<VectorVerificationStats>> {
  return get<VectorVerificationStats>('/admin/vector-verification/statistics')
}

export function getVerificationProgress(): Promise<ApiResponse<TaskProgressResponse>> {
  return get<TaskProgressResponse>('/admin/vector-verification/progress')
}

export function repairMissingVector(data: RepairMissingVectorRequest): Promise<ApiResponse<null>> {
  return post<null>('/admin/vector-verification/repair', data)
}

export function batchRepairMissingVectors(
  data: BatchRepairMissingVectorsRequest
): Promise<ApiResponse<{ affected_count: number }>> {
  return post<{ affected_count: number }>('/admin/vector-verification/batch-repair', data)
}

export function verifySpecificVectors(fileIds: string[]): Promise<ApiResponse<VectorVerificationTask>> {
  const fileIdsStr = fileIds.join(',')
  return post<VectorVerificationTask>(`/admin/vector-verification/verify/${fileIdsStr}`)
}

export function cancelVerificationTask(taskId: string): Promise<ApiResponse<null>> {
  return del<null>(`/admin/vector-verification/task/${taskId}`)
}

export function smartFix(data: SmartFixRequest): Promise<ApiResponse<{ message: string }>> {
  return post<{ message: string }>('/admin/vector-verification/smart-fix', data)
}

export default {
  startVectorVerification,
  getTaskStatus,
  getTaskList,
  getVerificationStatistics,
  getVerificationProgress,
  repairMissingVector,
  batchRepairMissingVectors,
  verifySpecificVectors,
  cancelVerificationTask,
  smartFix,
}
