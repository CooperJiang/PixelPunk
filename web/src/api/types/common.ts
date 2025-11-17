/**
 * 通用API类型定义
 * 包含所有模块共享的基础类型

 */
import type { ApiResult } from '@/utils/network/http-types'

/* ==================== 基础响应类型 ==================== */
/* 使用统一的 ApiResult 类型 */
export type ApiResponse<T> = ApiResult<T>

/* ==================== 分页相关类型 ==================== */
export interface PaginationParams {
  page?: number
  size?: number
}

export interface PaginationInfo {
  total: number
  per_page: number
  current_page: number
  last_page: number
}

export interface PaginatedResponse<T> {
  items: T[]
  pagination: PaginationInfo
}

/* ==================== 排序相关类型 ==================== */
export type SortOrder = 'asc' | 'desc'

export interface SortParams {
  sort?: string
  order?: SortOrder
}

/* ==================== 通用查询参数 ==================== */
export interface SearchParams {
  keyword?: string
  search?: string
}

export interface DateRangeParams {
  start_date?: string
  end_date?: string
}

/* ==================== 状态相关类型 ==================== */
export interface StatusInfo {
  id: number | string
  status: number
  status_text?: string
  is_active?: boolean
}

/* ==================== 时间戳类型 ==================== */
export interface TimeStamps {
  created_at: string
  updated_at?: string
}

/* ==================== 用户基础信息 ==================== */
export interface BaseUserInfo {
  id: number
  username: string
  avatar: string
}

/* ==================== 文件信息类型 ==================== */
export interface FileInfo {
  name: string
  size: number
  type: string
}

/* ==================== 访问级别类型 ==================== */
export type AccessLevel = 'public' | 'private' | 'protected'

/* ==================== 权限相关类型 ==================== */
export type Permission = 'read' | 'write' | 'delete' | 'admin'

/* ==================== 操作结果类型 ==================== */
export interface OperationResult {
  success: boolean
  message?: string
  affected_count?: number
}

/* ==================== 批量操作类型 ==================== */
export interface BatchOperation<T = string | number> {
  action: string
  ids: T[]
  params?: Record<string, unknown>
}
