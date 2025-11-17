import { get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/* 接口返回的用户信息类型 */
export interface UserItem {
  id: number
  username: string
  email: string
  avatar: string
  avatar_full_path?: string
  status: number
  role: number
  storage_limit: number // 存储空间限制（字节）
  bandwidth_limit: number // 带宽限制（字节）
  used_storage: number // 已使用存储（字节）
  used_bandwidth: number // 已使用带宽（字节）
  total_images: number // 总文件数量
  total_views: number // 总浏览数
  last_activity_at?: string // 最后操作时间
  last_activity_ip?: string // 最后操作IP
  created_at: string
  updated_at: string
}

/* 用户列表查询参数类型 */
export interface UserListParams {
  page?: number
  size?: number
  keyword?: string
  status?: number
  role?: number
  start_date?: string
  end_date?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

/* 用户列表查询结果类型 */
export interface UserListResult {
  total: number
  list: UserItem[]
}

export interface UpdateUserParams {
  id: number
  username: string
  status: number
  role: number
}

/* 用户详情结果类型 - 简化版，直接返回用户数据 */
export interface UserDetailResult extends UserItem {
  last_login_at?: string
}

/* 批量操作参数类型 */
export interface BatchOperationParams {
  user_ids: number[]
  operation: 'enable' | 'disable' | 'delete' | 'set_role'
  role?: number
}

/* 存储设置参数类型 */
export interface StorageSettingsParams {
  user_id: number
  storage_limit: number
  bandwidth_limit: number
}

/* 重置密码结果类型 */
export interface ResetPasswordResult {
  new_password: string
}

/* 发送邮件参数类型 */
export interface SendEmailParams {
  user_id: number
  subject: string
  content: string
}

export interface CreateUserParams {
  username: string
  email: string
  password: string
  role: number
  storage_limit?: number
  bandwidth_limit?: number
}

export function getUserList(params: UserListParams): Promise<ApiResult<UserListResult>> {
  return get<UserListResult>('/admin/user/list', params)
}

export function updateUser(data: UpdateUserParams): Promise<ApiResult<null>> {
  return post<null>('/admin/user/update', data)
}

export function getUserDetail(userId: number): Promise<ApiResult<UserDetailResult>> {
  return get<UserDetailResult>(`/admin/user/detail/${userId}`)
}

export function batchOperateUsers(data: BatchOperationParams): Promise<ApiResult<null>> {
  return post<null>('/admin/user/batch', data)
}

export function updateUserStorage(data: StorageSettingsParams): Promise<ApiResult<null>> {
  return post<null>('/admin/user/storage', data)
}

export function resetUserPassword(userId: number): Promise<ApiResult<ResetPasswordResult>> {
  return post<ResetPasswordResult>(`/admin/user/reset-password/${userId}`, {})
}

export function sendUserEmail(data: SendEmailParams): Promise<ApiResult<null>> {
  return post<null>('/admin/user/send-email', data)
}

export function toggleUserStatus(userId: number, status: number): Promise<ApiResult<null>> {
  return post<null>('/admin/user/toggle-status', { user_id: userId, status })
}

export function createUser(data: CreateUserParams): Promise<ApiResult<UserItem>> {
  return post<UserItem>('/admin/user/create', data)
}

export function deleteUser(userId: number): Promise<ApiResult<null>> {
  return post<null>(`/admin/user/delete/${userId}`, {})
}

export default {
  getUserList,
  updateUser,
  getUserDetail,
  batchOperateUsers,
  updateUserStorage,
  resetUserPassword,
  sendUserEmail,
  toggleUserStatus,
  createUser,
  deleteUser,
}
