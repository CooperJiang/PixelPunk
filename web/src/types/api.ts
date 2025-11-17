/**
 * API 相关类型定义
 * 包含请求响应、分页、搜索等API相关类型

 */
import type { ApiResult } from '@/utils/network/http-types'
import type { Timestamp } from './base'

/* 统一API响应格式 - 使用 ApiResult */
export type ApiResponse<T = any> = ApiResult<T>

/* API错误响应 */
export interface ApiError {
  code: number
  message: string
  details?: any
  field?: string
  timestamp: Timestamp
  request_id?: string
}

/* 分页请求参数 */
export interface PaginationParams {
  page: number
  size: number
  sort?: string
  order?: 'asc' | 'desc'
}

/* 分页响应数据 */
export interface PaginationResponse<T> {
  items: T[]
  total: number
  page: number
  size: number
  pages: number
  has_next: boolean
  has_prev: boolean
}

/* 搜索参数 */
export interface SearchParams extends PaginationParams {
  keyword?: string
  filters?: Record<string, any>
  date_range?: {
    start: string
    end: string
  }
}
