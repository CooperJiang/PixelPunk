import { del, get, post, put } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { AdminShareListParams, AdminVisitorListParams, UpdateShareStatusParams, VisitorListParams } from './types'

const baseUrl = '/shares'

export const createShare = (params: Record<string, unknown>): Promise<ApiResult<{ id: string; share_key: string }>> =>
  post(`${baseUrl}`, params)

export const getShareList = (params?: {
  page?: number
  size?: number
  status?: number
  keyword?: string
}): Promise<ApiResult<{ shares: unknown[]; total: number }>> => get(`${baseUrl}`, params)

export const getShareDetail = (id: string): Promise<ApiResult<Record<string, unknown>>> => get(`${baseUrl}/${id}`)

export const deleteShare = (id: string, force?: boolean): Promise<ApiResult<Record<string, unknown>>> => {
  const params = force ? { force: 'true' } : undefined
  return del(`${baseUrl}/${id}`, params)
}

export const getPublicShare = (key: string, params?: Record<string, unknown>): Promise<ApiResult<Record<string, unknown>>> =>
  get(`${baseUrl}/public/${key}`, params)

export const verifySharePassword = (key: string, password: string): Promise<ApiResult<Record<string, unknown>>> =>
  post(`${baseUrl}/public/${key}/verify`, { password })

export const submitVisitorInfo = (
  key: string,
  params: { name: string; email?: string }
): Promise<ApiResult<Record<string, unknown>>> => post(`${baseUrl}/public/${key}/visitor`, params)

export const getVisitorList = (shareId: string, params?: VisitorListParams): Promise<ApiResult<Record<string, unknown>>> =>
  get(`${baseUrl}/${shareId}/visitors`, params)

export const deleteVisitor = (shareId: string, visitorId: string): Promise<ApiResult<Record<string, unknown>>> =>
  del(`${baseUrl}/${shareId}/visitors/${visitorId}`)

/* 批量下载文件 */
export const downloadFiles = (fileIds: string[], accessToken?: string): Promise<ApiResult<Record<string, unknown>>> => {
  const params: Record<string, unknown> = { file_ids: fileIds }

  if (accessToken) {
    params.access_token = accessToken
  }

  return post(`${baseUrl}/download-files`, params)
}

export const adminGetShareList = (params?: AdminShareListParams): Promise<ApiResult<Record<string, unknown>>> =>
  get(`/admin/shares/list`, params)

export const adminGetShareDetail = (id: string): Promise<ApiResult<Record<string, unknown>>> => get(`/admin/shares/${id}`)

export const adminUpdateShareStatus = (
  id: string,
  params: UpdateShareStatusParams
): Promise<ApiResult<Record<string, unknown>>> => put(`/admin/shares/${id}/status`, params)

export const adminDeleteShare = (id: string, force?: boolean): Promise<ApiResult<Record<string, unknown>>> => {
  const params = force ? { force: 'true' } : undefined
  return del(`/admin/shares/${id}`, params)
}

export const getShareStats = (): Promise<ApiResult<Record<string, unknown>>> => get(`/admin/shares/stats`)

export const adminGetVisitorList = (params?: AdminVisitorListParams): Promise<ApiResult<Record<string, unknown>>> =>
  get(`/admin/shares/visitors`, params)

export const adminDeleteVisitor = (id: string): Promise<ApiResult<Record<string, unknown>>> => del(`/admin/shares/visitors/${id}`)

export const adminGetShareAccessToken = (id: string): Promise<ApiResult<Record<string, unknown>>> =>
  get(`/admin/shares/${id}/access-token`)

export const shareApi = {
  createShare,
  getShareList,
  getShareDetail,
  deleteShare,
  getPublicShare,
  verifySharePassword,
  submitVisitorInfo,
  getVisitorList,
  deleteVisitor,
  downloadFiles,
  adminGetShareList,
  adminGetShareDetail,
  adminUpdateShareStatus,
  adminDeleteShare,
  getShareStats,
  adminGetVisitorList,
  adminDeleteVisitor,
  adminGetShareAccessToken,
}
