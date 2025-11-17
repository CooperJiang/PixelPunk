/**
 * 文件相关类型定义
 */
import type { AccessLevel, PaginatedResponse, PaginationParams, SortParams, TimeStamps } from './common'
import type { BaseUserInfo } from './user'

/* ==================== 文件基础信息类型 ==================== */
export interface FileInfo extends TimeStamps {
  id: string
  url: string
  thumb_url: string
  full_url: string
  full_thumb_url: string
  original_name: string
  display_name: string
  size: number
  width?: number
  height?: number
  format: string
  access_level?: AccessLevel
  folder_id?: string
  views?: number
  downloads?: number
  is_duplicate?: boolean
  md5_hash?: string
  is_recommended?: boolean
  storage_provider_id?: string
  expires_at?: string
  storage_duration?: string
  is_time_limited?: boolean
  is_guest_upload?: boolean
  guest_fingerprint?: string
  expiry_notification_sent?: boolean
  user_id?: number
  user_name?: string
  user_info?: FileUserInfo
  ai_info?: FileAIInfo
  thumbnail_generation_failed?: boolean
  thumbnail_failure_reason?: string
  watermark_applied?: boolean
  watermark_failure_reason?: string
}

/* ==================== 文件关联用户信息类型 ==================== */
export interface FileUserInfo extends BaseUserInfo {
  email: string
  created_at: string
  total_files: number
  total_views: number
  days_joined: number
  other_files: FileThumbnail[]
  top_tags: string[]
}

/* ==================== 文件缩略图信息类型 ==================== */
export interface FileThumbnail {
  id: string
  display_name: string
  full_url: string
  full_thumb_url: string
  width: number
  height: number
}

/* ==================== 文件AI信息类型 ==================== */
export interface FileAIInfo {
  description: string
  tags: string[]
  dominant_color: string
  resolution: string
  is_nsfw: boolean
  nsfw_score: number
  nsfw_evaluation: string
}

/* ==================== 文件查询参数类型 ==================== */
export interface FileListParams extends PaginationParams, SortParams {
  folder_id?: string
  access_level?: AccessLevel | ''
  sort?: 'newest' | 'oldest' | 'name' | 'size'
}

/* ==================== 文件列表响应类型 ==================== */
export type FileListResponse = PaginatedResponse<FileInfo>

/* ==================== 文件更新请求类型 ==================== */
export interface UpdateFileRequest {
  folder_id?: string
  access_level?: AccessLevel
  name?: string
}

/* ==================== 文件统计类型 ==================== */
export interface FileStatsResponse {
  views: number
  downloads: number
  bandwidth: number
  referrers: Array<{
    domain: string
    views: number
  }>
  countries?: Array<{
    code: string
    views: number
  }>
}

/* ==================== 文件批量操作类型 ==================== */
export interface BatchFileOperation {
  action: 'delete' | 'move' | 'update_access' | 'update_folder'
  file_ids: string[]
  folder_id?: string
  access_level?: AccessLevel
}

/* ==================== 作者页面相关类型 ==================== */
export interface AuthorInfo extends BaseUserInfo {
  avatarFullPath: string
  bio: string
  website: string
  createdAt: string
  daysJoined: number
}

export interface AuthorFileInfo {
  id: string
  fileName: string
  originalName: string
  url: string
  fullPath: string
  thumbURL: string
  fullThumbURL: string
  size: number
  width: number
  height: number
  views: number
  createdAt: string
}
