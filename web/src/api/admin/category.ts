/**
 * 分类模板管理 API

 */
import { post, get } from '@/utils/network/http'
import type {
  CategoryTemplate,
  CreateTemplateRequest,
  UpdateTemplateRequest,
  DeleteTemplateRequest,
  TemplateListQuery,
  TemplateListResponse,
  BatchSortOrderRequest,
  PopularTemplatesResponse,
  AllTemplatesResponse,
} from '@/api/types/category'
import type { ApiResponse } from '@/api/types/common'

/**
 * 创建分类模板
 */
export function createTemplate(data: CreateTemplateRequest): Promise<ApiResponse<CategoryTemplate>> {
  return post<CategoryTemplate>('/admin/category-templates/create', data)
}

export function getTemplate(id: number): Promise<ApiResponse<CategoryTemplate>> {
  return get<CategoryTemplate>(`/admin/category-templates/detail/${id}`)
}

export function updateTemplate(data: UpdateTemplateRequest): Promise<ApiResponse<CategoryTemplate>> {
  return post<CategoryTemplate>('/admin/category-templates/update', data)
}

export function deleteTemplate(data: DeleteTemplateRequest): Promise<ApiResponse<void>> {
  return post<void>('/admin/category-templates/delete', data)
}

export function getTemplateList(params: TemplateListQuery): Promise<ApiResponse<TemplateListResponse>> {
  return get<TemplateListResponse>('/admin/category-templates/list', params)
}

export function batchUpdateSortOrder(data: BatchSortOrderRequest): Promise<ApiResponse<void>> {
  return post<void>('/admin/category-templates/batch-sort', data)
}

export function getPopularTemplates(limit?: number): Promise<ApiResponse<PopularTemplatesResponse>> {
  const params = limit ? { limit } : {}
  return get<PopularTemplatesResponse>('/category-templates/popular', params)
}

export function getAllTemplates(): Promise<ApiResponse<AllTemplatesResponse>> {
  return get<AllTemplatesResponse>('/category-templates/all')
}

export default {
  createTemplate,
  getTemplate,
  updateTemplate,
  deleteTemplate,
  getTemplateList,
  batchUpdateSortOrder,
  getPopularTemplates,
  getAllTemplates,
}
