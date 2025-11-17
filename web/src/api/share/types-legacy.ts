/* 分享项目类型 */
export type ShareItemType = 'folder' | 'file'

/* 分享项目 */
export interface ShareItem {
  item_type: ShareItemType
  item_id: string
}

export interface CreateShareParams {
  name?: string
  description?: string
  password?: string
  expired_days?: number
  max_views?: number
  items: ShareItem[]
  collect_visitor_info?: boolean
  notification_on_access?: boolean
  notification_threshold?: number
}

/* 分享状态 */
export enum ShareStatus {
  Normal = 1,
  Expired = 2,
  Deleted = 3,
  Disabled = 4,
}

/* 分享信息 */
export interface ShareInfo {
  id: string
  user_id: number
  share_key: string
  name: string
  description?: string
  expired_days: number
  expired_at?: string
  max_views: number
  current_views: number
  status: ShareStatus
  created_at: string
  updated_at: string
  collect_visitor_info: boolean
  notification_on_access: boolean
  notification_threshold?: number
}

/* 分享详情中的项目 */
export interface ShareItemDetail {
  id: string
  share_id: string
  item_type: ShareItemType
  item_id: string
  sort_order: number
  created_at: string
}

/* 分享详情 */
export interface ShareDetail {
  share: ShareInfo
  items: ShareItemDetail[]
  url: string
}

/* 需要密码的分享响应 */
export interface RequirePasswordResponse {
  require_password: true
  share_id: string
  name: string
}

/* 密码验证成功响应 */
export interface PasswordVerifyResponse {
  access_token: string
}

/* 公开分享查看响应 */
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
  folders: {
    id: string
    name: string
    description?: string
    created_at: string
  }[]
  images: {
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
  }[]
}

/* 访客信息提交参数 */
export interface VisitorInfo {
  name: string
  email?: string
  access_token: string
}

/* 访客信息详情 */
export interface VisitorInfoDetail {
  id: string
  created_at: string
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

/* 访客信息列表查询参数 */
export interface VisitorListParams {
  page?: number
  size?: number
  keyword?: string
  order_by?: string
}

/* 访客信息列表响应 */
export interface VisitorListResponse {
  list: VisitorInfoDetail[]
  total: number
  current_page?: number
  last_page?: number
}

/* 管理员分享列表查询参数 */
export interface AdminShareListParams {
  page?: number
  size?: number
  status?: ShareStatus
  keyword?: string
  user_id?: number | string
  start_date?: string
  end_date?: string
  order_by?: string
}

/* 管理员分享列表响应 */
export interface AdminShareListResponse {
  list: ShareInfo[]
  total: number
  current_page?: number
  last_page?: number
}

/* 管理员更新分享状态参数 */
export interface UpdateShareStatusParams {
  status: ShareStatus
  reason?: string
}

/* 分享统计数据 */
export interface ShareStats {
  total_shares: number
  active_shares: number
  views_today: number
  popular_shares: {
    id: string
    share_key: string
    name: string
    user_id: number
    username: string
    total_views: number
    recent_views: number
  }[]
  stats_by_date: {
    date: string
    new_shares: number
    total_views: number
  }[]
}

/* 管理员获取所有访客信息参数 */
export interface AdminVisitorListParams {
  page?: number
  size?: number
  keyword?: string
  order_by?: string
  share_id?: string
}

/* 管理员获取所有访客信息响应 */
export interface AdminVisitorListResponse {
  list: VisitorInfoDetail[]
  total: number
  current_page?: number
  last_page?: number
}
