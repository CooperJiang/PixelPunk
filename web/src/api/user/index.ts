import { get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type {
  AccessControlConfig,
  LoginResponse,
  ResetPasswordRequest,
  SendCodeRequest,
  SendCodeResponse,
  UserInfo,
  UserLoginRequest,
  UserRegisterRequest,
} from '../types'

/* 用户注册 */
export function register(data: UserRegisterRequest): Promise<ApiResult<UserInfo>> {
  return post<UserInfo>('/user/register', data)
}

export function login(data: UserLoginRequest): Promise<ApiResult<LoginResponse>> {
  return post<LoginResponse>('/user/login', data, {
    autoShowError: false,
    useResultMode: true,
  })
}

export function sendRegistrationCode(data: SendCodeRequest): Promise<ApiResult<SendCodeResponse>> {
  return post<SendCodeResponse>('/user/send-registration-code', data, {
    autoShowError: false,
  })
}

export function sendResetPasswordCode(data: SendCodeRequest): Promise<ApiResult<SendCodeResponse>> {
  return post<SendCodeResponse>('/user/send-reset-password-code', data, {
    autoShowError: false,
  })
}

export function resetPassword(data: ResetPasswordRequest): Promise<ApiResult<{ email: string }>> {
  return post<{ email: string }>('/user/reset-password', data)
}

export function getUserProfile(): Promise<ApiResult<UserInfo>> {
  return get<UserInfo>('/user/personal/profile')
}

export function updateUserProfile(data: UserInfo): Promise<ApiResult<UserInfo>> {
  return post<UserInfo>('/user/personal/profile', data)
}

export function updatePassword(data: { oldPassword: string; newPassword: string }): Promise<ApiResult<void>> {
  return post<void>('/user/personal/update-password', data)
}

export function getUserFolders(params?: {
  user_id?: number
  public_only?: boolean
  parent_id?: string
  with_stats?: boolean
}): Promise<ApiResult<Record<string, unknown>>> {
  return get<Record<string, unknown>>('/user/personal/folders', params)
}

export function sendChangeEmailCode(data: { new_email: string }): Promise<ApiResult<Record<string, unknown>>> {
  return post<Record<string, unknown>>('/user/personal/send-change-email-code', data)
}

export function changeEmail(data: { new_email: string; code: string }): Promise<ApiResult<Record<string, unknown>>> {
  return post<Record<string, unknown>>('/user/personal/change-email', data)
}

export function getAccessControlConfig(): Promise<ApiResult<AccessControlConfig>> {
  return get<AccessControlConfig>('/user/personal/access-control')
}

export function createOrUpdateAccessControl(data: AccessControlConfig): Promise<ApiResult<AccessControlConfig>> {
  return post<AccessControlConfig>('/user/personal/access-control/createOrUpdate', data)
}

export function resetAccessControlConfig(): Promise<ApiResult<null>> {
  return post<null>('/user/personal/access-control/reset')
}

export default {
  register,
  login,
  sendRegistrationCode,
  sendResetPasswordCode,
  resetPassword,
  sendChangeEmailCode,
  changeEmail,
  getAccessControlConfig,
  createOrUpdateAccessControl,
  resetAccessControlConfig,
}
