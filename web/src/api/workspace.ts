import { get } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/* 工作台数据接口（匹配后端DTO） */

/* 仪表盘数据类型定义 */
export interface DashboardStats {
  user_id: number
  username: string
  email: string
  created_at: string
  total_files: number
  total_views: number
  used_storage: number
  used_storage_formatted: string
  total_storage: number
  total_storage_formatted: string
  storage_usage_percent: number
  total_folders: number
  files_uploaded_today: number
  files_uploaded_this_week: number
  files_uploaded_this_month: number
  views_today: number
  views_this_week: number
  views_this_month: number
  recent_files: RecentFile[]
  folders: FolderInfo[]
}

/* 最近文件类型（匹配后端DTO） */
export interface RecentFile {
  id: string
  name: string
  url: string
  thumb_url: string
  full_url: string
  full_thumb_url: string
  size: number
  size_formatted: string
  views: number
  folder_id?: string
  folder_name?: string
  created_at: string
  storage_type: string
}

/* 文件夹信息类型（匹配后端DTO） */
export interface FolderInfo {
  id: string
  name: string
  file_count: number
  total_size: number
  total_size_formatted: string
  description?: string
  permission: string
  created_at: string
}

/* 用户资料类型 */
export interface UserProfile {
  id: number
  username: string
  email: string
  avatar_url?: string
  created_at: string
  role: number
  status: number
  used_storage: number
  total_storage: number
}

export function getDashboardStats(recentFilesLimit: number = 6): Promise<DashboardStats> {
  return get<DashboardStats>(`/user/personal/dashboard?recent_files_limit=${recentFilesLimit}`)
}

export function getUserProfile(): Promise<UserProfile> {
  return get<UserProfile>('/user/personal/profile')
}

export function getRecentImages(limit: number = 10): Promise<RecentFile[]> {
  return getDashboardStats(limit).then((data) => data.recent_files || [])
}

export function getFoldersList(): Promise<{ folders: FolderInfo[] }> {
  return get('/user/personal/folders').then((data) => ({
    folders: data.folders || [],
  }))
}

export interface WorkspaceStats {
  total_files: number // 文件总数
  used_storage: number // 已用存储(字节)
  used_storage_formatted: string // 已用存储格式化
  total_storage: number // 总存储(字节)
  total_storage_formatted: string // 总存储格式化
  total_views: number // 总访问量
  total_shares: number // 总分享数
  used_bandwidth: number // 已用带宽(字节)
  used_bandwidth_formatted: string // 已用带宽格式化
  total_bandwidth: number // 总带宽(字节)
  total_bandwidth_formatted: string // 总带宽格式化
}

export function getWorkspaceStats(): Promise<ApiResult<WorkspaceStats>> {
  return get<WorkspaceStats>('/user/personal/workspace/stats')
}

export default {
  getDashboardStats,
  getUserProfile,
  getRecentImages,
  getFoldersList,
  getWorkspaceStats,
}
