/**
 * 用户分类管理 API

 */
import { post, get } from '@/utils/network/http'
import type { ImageCategory } from '@/api/types/category'
import type { ApiResponse } from '@/api/types/common'

/**
 * 创建分类请求
 */
export interface CreateCategoryRequest {
  name: string
  description?: string
  sort_order?: number
}

/**
 * 更新分类请求
 */
export interface UpdateCategoryRequest {
  id: number
  name: string
  description?: string
  sort_order?: number
}

/**
 * 删除分类请求
 */
export interface DeleteCategoryRequest {
  id: number
}

/**
 * 批量删除请求
 */
export interface BatchDeleteRequest {
  ids: number[]
}

/**
 * 更新状态请求
 */
export interface UpdateStatusRequest {
  id: number
  status: 'active' | 'archived'
}

/**
 * 批量更新排序请求
 */
export interface BatchSortOrderRequest {
  sort_orders: Array<{
    id: number
    sort_order: number
  }>
}

/**
 * 分类列表查询参数
 */
export interface CategoryListQuery {
  keyword?: string
  status?: string
  page?: number
  size?: number
  sort_by?: 'name' | 'sort_order' | 'file_count' | 'created_at' | 'updated_at'
  sort_order_dir?: 'asc' | 'desc'
}

/**
 * 分类列表响应
 */
export interface CategoryListResponse {
  categories: ImageCategory[]
  total: number
  page: number
  size: number
}

/**
 * 所有分类响应
 */
export interface AllCategoriesResponse {
  categories: ImageCategory[]
  count: number
}

/**
 */
export function createCategory(data: CreateCategoryRequest): Promise<ApiResponse<ImageCategory>> {
  return post<ImageCategory>('/user-categories/create', data)
}

export function getCategoryList(params: CategoryListQuery): Promise<ApiResponse<CategoryListResponse>> {
  return get<CategoryListResponse>('/user-categories/list', params)
}

export function getAllCategories(): Promise<ApiResponse<AllCategoriesResponse>> {
  return get<AllCategoriesResponse>('/user-categories/all')
}

export function getCategoryDetail(id: number): Promise<ApiResponse<ImageCategory>> {
  return get<ImageCategory>(`/user-categories/${id}`)
}

export function updateCategory(data: UpdateCategoryRequest): Promise<ApiResponse<ImageCategory>> {
  return post<ImageCategory>('/user-categories/update', data)
}

export function deleteCategory(data: DeleteCategoryRequest): Promise<ApiResponse<void>> {
  return post<void>('/user-categories/delete', data)
}

export function batchDelete(data: BatchDeleteRequest): Promise<ApiResponse<{ deleted_count: number }>> {
  return post<{ deleted_count: number }>('/user-categories/batch-delete', data)
}

export function updateCategoryStatus(data: UpdateStatusRequest): Promise<ApiResponse<void>> {
  return post<void>('/user-categories/update-status', data)
}

export function batchUpdateSortOrder(data: BatchSortOrderRequest): Promise<ApiResponse<void>> {
  return post<void>('/user-categories/batch-sort', data)
}

export default {
  createCategory,
  getCategoryList,
  getAllCategories,
  getCategoryDetail,
  updateCategory,
  deleteCategory,
  batchDelete,
  updateCategoryStatus,
  batchUpdateSortOrder,
}
