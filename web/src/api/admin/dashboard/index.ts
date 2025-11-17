import { get } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/* ===== 仪表盘API接口类型定义 ===== */

/* 用户统计响应 */
export interface DashboardUserStatsResponse {
  total_users: number
  active_users_today: number
  new_users_today: number
  growth_rate: number
  banned_users: number
}

/* 文件统计响应 */
export interface DashboardImageStatsResponse {
  total_files: number
  uploaded_today: number
  ai_tagged_percentage: number
  pending_review: number
  nsfw_detected: number
  untagged_count: number
}

/* 存储统计响应 */
export interface DashboardStorageStatsResponse {
  total_storage: number
  formatted_storage: string
  growth_rate: number
  average_file_size: number
  usage_percentage: number
  new_storage_today: number // 今日新增存储量(字节)
  formatted_new_storage: string // 今日新增存储量(格式化)
}

/* 系统资源响应 */
export interface DashboardSystemResourcesResponse {
  cpu_usage: number
  memory_usage: number
  disk_usage: number
  load_average: number[]
  uptime: string
}

/* 上传趋势响应 */
export interface DashboardUploadTrendsResponse {
  timeline: string[]
  upload_counts: number[]
  storage_growth: number[]
}

/* AI服务状态响应 */
export interface AITaggingStats {
  total_count: number
  none_count: number // 未处理
  pending_count: number // 打标中
  done_count: number // 已完成
  failed_count: number // 失败
  skipped_count: number // 跳过
  ignored_count: number // 忽略
}

export interface AIVectorStats {
  total_count: number
  pending_count: number // 待处理
  processing_count: number // 处理中
  completed_count: number // 已完成
  failed_count: number // 失败
  reset_count: number // 待重置
  active_workers: number
  max_workers: number
}

export interface AIHealthStatus {
  tagging_service: 'healthy' | 'degraded' | 'down'
  vector_service: 'healthy' | 'degraded' | 'down'
}

export interface DashboardAIServicesResponse {
  tagging: AITaggingStats
  vectors: AIVectorStats
  health_status: AIHealthStatus
}

/* 分享统计响应 */
export interface DashboardShareStatsResponse {
  total_shares: number
  public_shares: number
  private_shares: number
  total_visits: number
  total_downloads: number
  growth_rate: number
  new_shares_today: number // 今日新增分享数
}

/* 标签统计响应 */
export interface DashboardTagItem {
  name: string
  count: number
}

export interface DashboardTagStatsResponse {
  ai_tags_count: number // AI标签种类数量
  tagged_images: number // 已标记的文件数量
  untagged_images: number // 未标记文件数量
  popular_tags: DashboardTagItem[] // 热门标签列表
  manual_tags_count: number // 手动标签数量（兼容性，始终为0）
}

/* 最近上传响应 */
export interface RecentUploadItem {
  id: string // 修改为string类型，与后端数据库一致
  file_name: string
  display_name?: string
  size: number
  size_formatted: string
  format: string // 文件格式，如 jpg, png, webp
  url: string
  thumbnail_url: string
  full_url: string
  full_thumb_url: string
  created_at: string
  user: {
    id: number
    username: string
  }
  ai_info?: {
    description?: string
    dominant_color?: string
    is_nsfw?: boolean
    nsfw_evaluation?: string
    nsfw_score?: number
    resolution?: string
    tags?: string[]
  }
}

export interface DashboardRecentUploadsResponse {
  recent_uploads: RecentUploadItem[]
}

/* ===== API接口函数 ===== */

const baseUrl = '/admin/stats'

/**
 * 获取用户统计数据
 */
export function getDashboardUserStats(): Promise<ApiResult<DashboardUserStatsResponse>> {
  return get<DashboardUserStatsResponse>(`${baseUrl}/users`)
}

export function getDashboardImageStats(): Promise<ApiResult<DashboardImageStatsResponse>> {
  return get<DashboardImageStatsResponse>(`${baseUrl}/files`)
}

export function getDashboardStorageStats(): Promise<ApiResult<DashboardStorageStatsResponse>> {
  return get<DashboardStorageStatsResponse>(`${baseUrl}/storage`)
}

export function getDashboardSystemResources(): Promise<ApiResult<DashboardSystemResourcesResponse>> {
  return get<DashboardSystemResourcesResponse>(`${baseUrl}/system-resources`)
}

export function getDashboardUploadTrends(days: number = 7): Promise<ApiResult<DashboardUploadTrendsResponse>> {
  return get<DashboardUploadTrendsResponse>(`${baseUrl}/upload-trends`, { days })
}

export function getDashboardShareStats(): Promise<ApiResult<DashboardShareStatsResponse>> {
  return get<DashboardShareStatsResponse>(`${baseUrl}/shares`)
}

export function getDashboardTagStats(): Promise<ApiResult<DashboardTagStatsResponse>> {
  return get<DashboardTagStatsResponse>(`${baseUrl}/tags`)
}

export function getDashboardRecentUploads(limit: number = 10): Promise<ApiResult<RecentUploadItem[]>> {
  return get<RecentUploadItem[]>(`${baseUrl}/latest-files`, { limit })
}

export default {
  getDashboardUserStats,
  getDashboardImageStats,
  getDashboardStorageStats,
  getDashboardSystemResources,
  getDashboardUploadTrends,
  getDashboardShareStats,
  getDashboardTagStats,
  getDashboardRecentUploads,
}
