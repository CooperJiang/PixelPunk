/**
 * 分类模块相关类型定义

/* ==================== 分类模板类型 ==================== */

export interface CategoryTemplate {
  id: number
  name: string
  description: string
  icon: string
  sort_order: number
  is_popular: boolean
  usage_count: number
  created_at: string
  updated_at: string
}

/* ==================== 请求参数类型 ==================== */

export interface CreateTemplateRequest {
  name: string
  description?: string
  icon?: string
  is_popular?: boolean
  sort_order?: number
}

export interface UpdateTemplateRequest {
  id: number
  name?: string
  description?: string
  icon?: string
  is_popular?: boolean
  sort_order?: number
}

export interface DeleteTemplateRequest {
  id: number
}

export interface TemplateListQuery {
  keyword?: string
  is_popular?: boolean
  page?: number
  size?: number
  sort_by?: 'name' | 'sort_order' | 'usage_count' | 'created_at'
  sort_order?: 'asc' | 'desc'
}

export interface BatchSortOrderRequest {
  sort_orders: Array<{
    id: number
    sort_order: number
  }>
}

/* ==================== 响应类型 ==================== */

export interface TemplateListResponse {
  templates: CategoryTemplate[]
  total: number
  page: number
  size: number
}

export interface PopularTemplatesResponse {
  templates: CategoryTemplate[]
  count: number
}

export interface AllTemplatesResponse {
  templates: CategoryTemplate[]
  count: number
}

/* ==================== 用户分类类型 ==================== */

export interface ImageCategory {
  id: number
  name: string
  description: string
  user_id: number
  sort_order: number
  source: 'user' | 'system_template' | 'ai_suggestion'
  template_id?: number
  file_count: number
  status: 'active' | 'archived'
  created_at: string
  updated_at: string
}

export interface ImageCategoryRelation {
  file_id: string
  category_id: number
  user_id: number
  source: 'manual' | 'ai'
  created_at: string
  category?: ImageCategory
}
