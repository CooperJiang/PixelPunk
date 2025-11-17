/**
 * 分享相关类型定义
 */
import type { AccessLevel, PaginationParams, TimeStamps } from './common'
import type { AuthorFolderInfo } from './folder'

/* ==================== 分享项目类型 ==================== */
export type ShareItemType = 'folder' | 'file'

export interface ShareItem {
  item_type: ShareItemType
  item_id: string
}

/* ==================== 分享状态枚举 ==================== */
export enum ShareStatus {
  Normal = 1,
  Expired = 2,
  Deleted = 3,
  Disabled = 4,
}

/* ==================== 分享基础信息类型 ==================== */
export interface ShareInfo extends TimeStamps {
  id: string
  user_id?: number
  share_key: string
  name: string
  description?: string
  expired_days?: number
  expired_at?: string
  max_views?: number
  current_views?: number
  status?: ShareStatus
  collect_visitor_info?: boolean
  notification_on_access?: boolean
  notification_threshold?: number
  shareKey?: string
  views?: number
  access_level?: AccessLevel
  password?: string
  expires_at?: string
  is_active?: boolean
  folder_id?: string
}

/* ==================== 分享详情类型 ==================== */
export interface ShareItemDetail extends TimeStamps {
  id: string
  share_id: string
  item_type: ShareItemType
  item_id: string
  sort_order: number
}

export interface ShareDetail {
  share: ShareInfo
  items: ShareItemDetail[]
  url: string
}

/* ==================== 作者分享信息类型 ==================== */
export interface AuthorShareInfo {
  id: string
  name: string
  description: string
  shareKey: string
  views: number
  createdAt: string
  coverImage: string
  coverImageFullPath: string
  coverImageThumbURL: string
  coverImageFullThumbURL: string
}

/* ==================== 分享创建请求类型 ==================== */
export interface CreateShareRequest {
  name?: string
  description?: string
  password?: string
  expired_days?: number
  max_views?: number
  items: ShareItem[]
  collect_visitor_info?: boolean
  notification_on_access?: boolean
  notification_threshold?: number
  folder_id?: string
  expires_at?: string
  access_level?: AccessLevel
}

/* ==================== 分享更新请求类型 ==================== */
export interface UpdateShareRequest {
  name?: string
  description?: string
  password?: string
  expires_at?: string
  access_level?: AccessLevel
  is_active?: boolean
  status?: ShareStatus
}

/* ==================== 分享访问相关类型 ==================== */
export interface ShareAccessRequest {
  shareKey: string
  password?: string
}

export interface RequirePasswordResponse {
  require_password: true
  share_id: string
  name: string
}

export interface PasswordVerifyResponse {
  access_token: string
}

/* ==================== 公开分享响应类型 ==================== */
export interface PublicShareResponse {
  share: {
    id: string
    name: string
    description?: string
    created_at: string
    expired_at?: string
    current_views: number
    max_views: number
    has_password: boolean
    collect_visitor_info?: boolean
  }
  user: {
    username: string
    avatar: string
  }
  folders: Array<{
    id: string
    name: string
    description?: string
    created_at: string
  }>
  files: Array<{
    id: string
    file_name: string
    display_name: string
    url: string
    thumb_url: string
    full_url: string
    full_thumb_url: string
    width?: number
    height?: number
    size: number
    size_formatted: string
    format: string
    created_at: string
  }>
}

/* ==================== 分享内容响应类型 ==================== */
export interface ShareContentResponse {
  share: ShareInfo
  folders?: AuthorFolderInfo[]
  files?: unknown[]
  pagination?: {
    currentPage: number
    perPage: number
    total: number
    lastPage: number
  }
}

/* ==================== 访客信息类型 ==================== */
export interface VisitorInfo {
  name: string
  email?: string
  access_token: string
}

export interface VisitorInfoDetail extends TimeStamps {
  id: string
  share_id: string
  share_key: string
  visitor_name: string
  visitor_email?: string
  ip_address?: string
  user_agent?: string
  referer?: string
  visit_count: number
  last_visit_at: string
}

/* ==================== 访客列表查询类型 ==================== */
export interface VisitorListParams extends PaginationParams {
  keyword?: string
  order_by?: string
}

export interface VisitorListResponse {
  list: VisitorInfoDetail[]
  total: number
  current_page?: number
  last_page?: number
}

/* ==================== 管理员分享管理类型 ==================== */
export interface AdminShareListParams extends PaginationParams {
  status?: ShareStatus
  keyword?: string
  user_id?: number | string
  start_date?: string
  end_date?: string
  order_by?: string
}

export interface AdminShareListResponse {
  list: ShareInfo[]
  total: number
  current_page?: number
  last_page?: number
}

export interface UpdateShareStatusParams {
  status: ShareStatus
  reason?: string
}

export interface AdminVisitorListParams extends PaginationParams {
  keyword?: string
  order_by?: string
  share_id?: string
}

export interface AdminVisitorListResponse {
  list: VisitorInfoDetail[]
  total: number
  current_page?: number
  last_page?: number
}

/* ==================== 分享统计类型 ==================== */
export interface ShareStats {
  total_shares: number
  active_shares: number
  views_today: number
  popular_shares: Array<{
    id: string
    share_key: string
    name: string
    user_id: number
    username: string
    total_views: number
    recent_views: number
  }>
  stats_by_date: Array<{
    date: string
    new_shares: number
    total_views: number
  }>
}

export interface ShareStatsResponse {
  total_views: number
  unique_visitors: number
  recent_views: Array<{
    date: string
    views: number
  }>
}
