import { post } from '@/utils/network/http'
import { useAuthStore } from '@/store/auth'
import type { FileInfo } from '@/api/types/file'

/* 秒传相关类型定义 */
export interface CheckDuplicateRequest {
  md5: string
  filename: string
  file_size: number
}

export interface CheckDuplicateFileInfo {
  id: string
  original_name: string
  size: number
  format: string
}

export interface CheckDuplicateResponse {
  exists: boolean
  original_file?: CheckDuplicateFileInfo
}

export interface InstantUploadRequest {
  md5: string
  filename: string
  file_size: number
  folder_id?: string
  access_level?: 'public' | 'private' | 'protected'
  optimize?: boolean
  storage_duration?: string
}

export interface InstantUploadResponse {
  file_info: FileInfo
  is_instant: boolean
  message: string
}

/* 秒传API */
export const instantUploadAPI = {
  async checkDuplicate(data: CheckDuplicateRequest): Promise<CheckDuplicateResponse> {
    const authStore = useAuthStore()
    const endpoint = authStore.isLoggedIn ? '/files/check-duplicate' : '/files/guest/check-duplicate'
    const result = await post<CheckDuplicateResponse>(endpoint, data, { silent: true })
    return result.data
  },

  async instantUpload(data: InstantUploadRequest): Promise<InstantUploadResponse> {
    const authStore = useAuthStore()
    const endpoint = authStore.isLoggedIn ? '/files/instant-upload' : '/files/guest/instant-upload'
    const result = await post<InstantUploadResponse>(endpoint, data, { silent: true })
    return result.data
  },
}

export default instantUploadAPI
