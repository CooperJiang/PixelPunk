/**
 * 公告模块相关类型定义

/* ==================== 公告基础类型 ==================== */

export interface Announcement {
  id: number
  title: string
  content: string // Markdown格式
  summary: string
  is_pinned: boolean
  status: 'draft' | 'published' | 'archived'
  view_count: number
  created_by: number
  created_at: string
  updated_at: string
}

/* ==================== 请求参数类型 ==================== */

export interface CreateAnnouncementRequest {
  title: string
  content: string // Markdown格式
  summary?: string
  is_pinned?: boolean
  status: 'draft' | 'published' | 'archived'
}

export interface UpdateAnnouncementRequest {
  id: number
  title?: string
  content?: string // Markdown格式
  summary?: string
  is_pinned?: boolean
  status?: 'draft' | 'published' | 'archived'
}

export interface AnnouncementListQuery {
  status?: 'draft' | 'published' | 'archived'
  is_pinned?: boolean
  keyword?: string
  page?: number
  page_size?: number
}

/* ==================== 响应类型 ==================== */

export interface PaginationResponse {
  page: number
  page_size: number
  total: number
  total_page: number
}

export interface AnnouncementListResponse {
  data: Announcement[]
  pagination: PaginationResponse
}

export interface AnnouncementSimple {
  id: number
  title: string
  summary: string
  is_pinned: boolean
  view_count: number
  created_at: string
}

export interface AnnouncementDetail {
  id: number
  title: string
  content: string // Markdown格式
  summary: string
  is_pinned: boolean
  view_count: number
  created_at: string
}

export interface PublicAnnouncementListResponse {
  announcements: AnnouncementSimple[]
  total: number
}

/* ==================== 公告设置类型 ==================== */

export interface AnnouncementSettings {
  announcement_enabled?: boolean
  announcement_drawer_position?: 'left' | 'right'
  announcement_drawer_width?: string | number // 支持 px、vw、% 等单位字符串，或纯数字（默认px）
  announcement_display_limit?: number
  announcement_auto_show_delay?: number // 自动弹窗延迟时间（秒）
}

export type AnnouncementSettingsUpdate = Partial<AnnouncementSettings>
