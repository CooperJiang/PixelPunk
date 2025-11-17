/**
 * 业务领域类型定义
 * 包含用户、文件、文件夹、分享、标签等业务相关类型

 */
import type { ID, Timestamp } from './base'

/* ===== 用户相关类型 ===== */

/* 用户基础信息 */
export interface User {
  id: ID
  username: string
  email: string
  avatar?: string
  role: UserRole
  status: UserStatus
  created_at: string
  updated_at: string
  last_login?: string
}

/* 用户角色 */
export enum UserRole {
  ADMIN = 1,
  USER = 2,
  GUEST = 3,
}

/* 用户状态 */
export enum UserStatus {
  ACTIVE = 1,
  INACTIVE = 0,
  BANNED = -1,
}

/* 登录凭据 */
export interface LoginCredentials {
  username: string
  password: string
  remember?: boolean
  captcha?: string
}

/* 注册信息 */
export interface RegisterInfo {
  username: string
  email: string
  password: string
  confirm_password: string
  verification_code: string
  agree_terms: boolean
}

/* 认证Token */
export interface AuthToken {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  scope?: string
}

/* ===== 文件相关类型 ===== */

/* 文件基础信息 */
export interface FileInfo {
  id: ID
  name: string
  size: number
  type: string
  extension: string
  url: string
  thumbnail?: string
  hash: string
  created_at: string
  updated_at: string
}

/* 文件信息（图片类型） */
export interface ImageInfo extends FileInfo {
  width: number
  height: number
  format: string
  quality?: number
  metadata?: ImageMetadata
  tags?: string[]
  description?: string
  alt_text?: string
  full_thumb_url?: string
  thumb_url?: string
  display_name?: string
  original_name?: string
  full_url?: string
  views?: number
  similarity?: number
  ai_info?: {
    description?: string
    [key: string]: any
  }
}

/* 文件元数据（图片类型） */
export interface ImageMetadata {
  camera?: string
  lens?: string
  iso?: number
  aperture?: string
  shutter_speed?: string
  focal_length?: string
  gps?: {
    latitude: number
    longitude: number
  }
  exif?: Record<string, any>
}

/* 上传状态 */
export enum UploadStatus {
  PENDING = 'pending',
  UPLOADING = 'uploading',
  PROCESSING = 'processing',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
}

/* 上传任务 */
export interface UploadTask {
  id: ID
  file: File
  progress: number
  status: UploadStatus
  error?: string
  result?: FileInfo
  created_at: Timestamp
  updated_at: Timestamp
}

/* ===== 文件夹相关类型 ===== */

/* 文件夹信息 */
export interface Folder {
  id: ID
  name: string
  description?: string
  parent_id?: ID
  path: string
  level: number
  children?: Folder[]
  file_count: number
  created_at: string
  updated_at: string
}

/* 文件夹树节点 */
export interface FolderTreeNode extends Folder {
  expanded?: boolean
  selected?: boolean
  loading?: boolean
  children?: FolderTreeNode[]
}

/* ===== 分享相关类型 ===== */

/* 分享类型 */
export enum ShareType {
  FOLDER = 'folder',
  FILE = 'file',
  ALBUM = 'album',
}

/* 分享权限 */
export enum SharePermission {
  VIEW = 'view',
  DOWNLOAD = 'download',
  COMMENT = 'comment',
}

/* 分享信息 */
export interface Share {
  id: ID
  name: string
  type: ShareType
  target_id: ID
  permissions: SharePermission[]
  password?: string
  expired_at?: string
  visit_count: number
  download_count: number
  is_public: boolean
  is_active: boolean
  created_at: string
  updated_at: string
  creator: Pick<User, 'id' | 'username' | 'avatar'>
}

/* 分享访问记录 */
export interface ShareVisit {
  id: ID
  share_id: ID
  ip: string
  user_agent: string
  location?: string
  visited_at: string
}

/* ===== 标签相关类型 ===== */

/* 标签信息 */
export interface Tag {
  id: ID
  name: string
  color?: string
  description?: string
  usage_count: number
  created_at: string
  updated_at: string
}

/* 标签类别 */
export interface TagCategory {
  id: ID
  name: string
  color: string
  tags: Tag[]
  created_at: string
}

/* ===== 系统相关类型 ===== */

/* 系统设置 */
export interface SystemSettings {
  site_name: string
  site_description: string
  site_keywords: string
  logo_url?: string
  favicon_url?: string
  allow_registration: boolean
  require_email_verification: boolean
  max_file_size: number
  allowed_file_types: string[]
  storage_provider: string
  storage_config: Record<string, any>
  email_config: Record<string, any>
  ai_config: Record<string, any>
}

/* 系统统计 */
export interface SystemStats {
  users: {
    total: number
    active: number
    new_today: number
  }
  files: {
    total: number
    size: number
    uploaded_today: number
  }
  storage: {
    used: number
    available: number
    percentage: number
  }
  shares: {
    total: number
    public: number
    private: number
  }
}

/* 系统日志 */
export interface SystemLog {
  id: ID
  level: 'debug' | 'info' | 'warn' | 'error' | 'fatal'
  message: string
  context?: Record<string, any>
  user_id?: ID
  ip?: string
  user_agent?: string
  created_at: string
}
