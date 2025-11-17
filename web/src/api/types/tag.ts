/**
 * 标签相关类型定义
 */
import type { BatchOperation, DateRangeParams, PaginationParams, SearchParams, SortParams, TimeStamps } from './common'

/* ==================== 标签基础信息类型 ==================== */
export interface TagInfo extends TimeStamps {
  id: number
  name: string
  count: number
}

/* ==================== 标签查询参数类型 ==================== */
export interface TagListParams extends PaginationParams, SortParams, SearchParams, DateRangeParams {
  keyword?: string
  sort_by?: 'name' | 'count' | 'created_at'
  sort_order?: 'asc' | 'desc'
  min_count?: number
}

/* ==================== 标签列表响应类型 ==================== */
export interface TagListResponse {
  items: TagInfo[]
  total: number
  pagination: {
    page: number
    size: number
  }
}

/* ==================== 标签创建请求类型 ==================== */
export interface CreateTagRequest {
  name: string
}

/* ==================== 标签更新请求类型 ==================== */
export interface UpdateTagRequest {
  id: number
  name: string
}

/* ==================== 标签批量操作类型 ==================== */
export interface BatchTagOperation extends BatchOperation<number> {
  action: 'delete' | 'merge'
  tag_ids: number[]
  target?: number // 合并目标标签ID
}
