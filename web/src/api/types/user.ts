/**
 * 用户相关类型定义
 */
import type { BaseUserInfo, TimeStamps } from './common'

/* ==================== 用户信息类型 ==================== */
export interface UserInfo extends BaseUserInfo, TimeStamps {
  email: string
  avatarFullPath?: string
  bio: string
  website?: string
  status?: number
  role: number
}

/* ==================== 用户认证类型 ==================== */
export interface UserLoginRequest {
  account: string
  password: string
}

export interface UserRegisterRequest {
  username: string
  email: string
  password: string
  code: string
}

export interface LoginResponse {
  token: string
  userInfo: UserInfo
}

/* ==================== 验证码相关类型 ==================== */
export interface SendCodeRequest {
  email: string
}

export interface SendCodeResponse {
  email: string
  expires_in: number
}

/* ==================== 密码重置类型 ==================== */
export interface ResetPasswordRequest {
  email: string
  code: string
  new_password: string
}

export interface UpdatePasswordRequest {
  oldPassword: string
  newPassword: string
}

/* ==================== 用户资料类型 ==================== */
export interface UpdateProfileRequest {
  username?: string
  avatar?: string
  bio?: string
  website?: string
}

/* ==================== 邮箱变更类型 ==================== */
export interface ChangeEmailRequest {
  new_email: string
  code: string
}

/* ==================== 用户统计类型 ==================== */
export interface UserStatsResponse {
  total_images: number
  total_size: number
  total_views: number
  bandwidth_used: number
  storage_limit: number
  storage_used_percent: number
}

/* ==================== 作者统计类型 ==================== */
export interface AuthorStats {
  totalFiles: number
  totalViews: number
  totalFolders: number
  totalShares: number
}

/* ==================== 访问控制类型 ==================== */
export type AccessControlMode = 'whitelist' | 'blacklist'
export type RestrictionMode = 'strict' | 'moderate' | 'loose'
export type BlockAction = 'block' | 'redirect' | 'watermark'

export interface AccessControlConfig {
  id?: number
  user_id?: number
  enabled: boolean
  ip_mode: AccessControlMode
  domain_mode: AccessControlMode
  restriction_mode: RestrictionMode
  block_action: BlockAction
  redirect_url?: string
  custom_error_message?: string
  ip_list: string
  domain_list: string
  created_at?: string
  updated_at?: string
}
