import { del, get, getApiBaseUrl, post, put, upload, uploadBatch } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import { StorageUtil } from '@/utils/storage/storage'

import type { FileInfo, FileListParams, FileListResponse, FileStatsResponse, UpdateFileRequest } from '../types'

/* 导入秒传API */
import instantUploadAPI from './instant'

/* 上传单张文件 - 新接口名 */
export function uploadFile(
  file: File,
  options?: {
    folder_id?: string
    access_level?: string
    optimize?: boolean
    storage_duration?: string
    watermark?: string
  },
  onUploadProgress?: (progressEvent: ProgressEvent<XMLHttpRequestEventTarget>) => void,
  config?: { silent?: boolean }
): Promise<ApiResult<FileInfo>> {
  return upload<FileInfo>('/files/upload', file, options, onUploadProgress, config)
}

export function uploadAvatar(file: File): Promise<ApiResult<{ url: string; id: string }>> {
  return upload<{ url: string; id: string }>('/avatar/upload', file)
}

export function uploadFiles(
  files: File[],
  options?: {
    folder_id?: string
    access_level?: string
    optimize?: boolean
    storage_duration?: string
    watermark?: string
  },
  onUploadProgress?: (progressEvent: ProgressEvent<XMLHttpRequestEventTarget>) => void
): Promise<ApiResult<FileInfo[]>> {
  return uploadBatch<FileInfo[]>('/files/batch-upload', files, options, onUploadProgress)
}

export function getFileList(params?: FileListParams): Promise<ApiResult<FileListResponse>> {
  return get<FileListResponse>('/files/list', params)
}

export function getFileDetail(fileId: string): Promise<ApiResult<FileInfo>> {
  return get<FileInfo>(`/files/${fileId}`)
}

export function getFileStats(fileId: string, period?: string): Promise<ApiResult<FileStatsResponse>> {
  return get<FileStatsResponse>(`/files/${fileId}/stats`, { period })
}

export function updateFile(fileId: string, data: UpdateFileRequest): Promise<ApiResult<FileInfo>> {
  return put<FileInfo>(`/files/${fileId}`, data)
}

export function deleteFile(fileId: string): Promise<ApiResult<{ id: string }>> {
  return del<{ id: string }>(`/files/${fileId}`)
}

export function batchDeleteFiles(fileIds: string[]): Promise<ApiResult<{ deleted: number; errors: string[] }>> {
  return post<{ deleted: number; errors: string[] }>('/files/batch-delete', { file_ids: fileIds })
}

export function moveFiles(fileIds: string[], folderId?: string): Promise<ApiResult<{ moved: number }>> {
  return post<{ moved: number }>('/files/move', { file_ids: fileIds, target_folder_id: folderId })
}

export function guestUpload(
  file: File,
  options: {
    access_level?: string
    optimize?: boolean
    storage_duration: string
    fingerprint: string
    watermark?: string
  },
  onUploadProgress?: (progressEvent: ProgressEvent<XMLHttpRequestEventTarget>) => void
): Promise<
  ApiResult<{
    file_info: FileInfo
    remaining_count: number
  }>
> {
  return upload<{
    file_info: FileInfo
    remaining_count: number
  }>('/files/guest/upload', file, options, onUploadProgress)
}

export function toggleAccessLevel(fileId: string): Promise<ApiResult<{ id: string }>> {
  return post<{ id: string }>(`/files/${fileId}/toggle-access-level`)
}

export function getRandomRecommendedFile(): Promise<ApiResult<FileInfo>> {
  return get<FileInfo>('/files/guest/random')
}

export function downloadSharedFile(
  fileId: string,
  shareKey: string,
  accessToken?: string,
  options?: {
    quality?: 'original' | 'compressed'
    format?: 'jpg' | 'png' | 'webp'
  },
  onProgress?: (progress: { loaded: number; total: number; percent: number }) => void
): Promise<{ blob: Blob; filename?: string }> {
  const params = new URLSearchParams()
  if (options?.quality) {
    params.append('quality', options.quality)
  }
  if (options?.format) {
    params.append('format', options.format)
  }

  if (accessToken) {
    params.append('access_token', accessToken)
  }

  const url = `/shares/public/${shareKey}/files/${fileId}/download${params.toString() ? `?${params.toString()}` : ''}`

  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.open('GET', getApiBaseUrl() + url, true)
    xhr.responseType = 'blob'

    xhr.onprogress = (event) => {
      if (event.lengthComputable && onProgress) {
        const percent = Math.round((event.loaded / event.total) * 100)
        onProgress({
          loaded: event.loaded,
          total: event.total,
          percent,
        })
      }
    }

    xhr.onload = () => {
      if (xhr.status === 200) {
        const contentDisposition = xhr.getResponseHeader('Content-Disposition')
        let filename = ''
        if (contentDisposition) {
          const filenameStarMatch = contentDisposition.match(/filename\*=UTF-8''([^;]+)/)
          if (filenameStarMatch) {
            try {
              filename = decodeURIComponent(filenameStarMatch[1])
            } catch {}
          }

          if (!filename) {
            const filenameMatch = contentDisposition.match(/filename="?([^"]+)"?/)
            if (filenameMatch) {
              filename = filenameMatch[1]
            }
          }
        }

        const result = {
          blob: xhr.response,
          filename,
        }
        resolve(result)
      } else {
        reject(new Error(`Download failed: ${xhr.status}`))
      }
    }

    xhr.onerror = () => {
      reject(new Error('Download request failed'))
    }

    xhr.send()
  })
}

export function downloadFile(
  fileId: string,
  options?: {
    quality?: 'original' | 'compressed'
    format?: 'jpg' | 'png' | 'webp'
  },
  onProgress?: (progress: { loaded: number; total: number; percent: number }) => void
): Promise<{ blob: Blob; filename?: string }> {
  const params = new URLSearchParams()
  if (options?.quality) {
    params.append('quality', options.quality)
  }
  if (options?.format) {
    params.append('format', options.format)
  }

  const url = `/files/${fileId}/download${params.toString() ? `?${params.toString()}` : ''}`

  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.open('GET', getApiBaseUrl() + url, true)
    xhr.responseType = 'blob'

    const tokenStr = StorageUtil.get<string>('token')
    if (tokenStr) {
      try {
        const tokenObj = JSON.parse(tokenStr)
        const token = tokenObj.data || tokenStr // 如果是对象格式则取data字段，否则直接使用
        xhr.setRequestHeader('Authorization', `Bearer ${token}`)
      } catch {
        xhr.setRequestHeader('Authorization', `Bearer ${tokenStr}`)
      }
    }

    xhr.onprogress = (event) => {
      if (event.lengthComputable && onProgress) {
        const percent = Math.round((event.loaded / event.total) * 100)
        onProgress({
          loaded: event.loaded,
          total: event.total,
          percent,
        })
      }
    }

    xhr.onload = () => {
      if (xhr.status === 200) {
        const contentDisposition = xhr.getResponseHeader('Content-Disposition')
        let filename = ''
        if (contentDisposition) {
          const filenameStarMatch = contentDisposition.match(/filename\*=UTF-8''([^;]+)/)
          if (filenameStarMatch) {
            try {
              filename = decodeURIComponent(filenameStarMatch[1])
            } catch {}
          }

          if (!filename) {
            const filenameMatch = contentDisposition.match(/filename="?([^"]+)"?/)
            if (filenameMatch) {
              filename = filenameMatch[1]
            }
          }
        }

        const result = {
          blob: xhr.response,
          filename,
        }
        resolve(result)
      } else {
        reject(new Error(`Download failed: ${xhr.status}`))
      }
    }

    xhr.onerror = () => {
      reject(new Error('Download request failed'))
    }

    xhr.send()
  })
}

export function reorderFiles(data: { folder_id?: string; file_ids: string[] }): Promise<ApiResult<void>> {
  return post<void>('/files/reorder', data)
}

export default {
  uploadFile,
  uploadFiles,
  uploadAvatar,
  guestUpload,
  getFileList,
  getFileDetail,
  getFileStats,
  updateFile,
  deleteFile,
  batchDeleteFiles,
  toggleAccessLevel,
  getRandomRecommendedFile,
  downloadFile,
  downloadSharedFile,
  reorderFiles,
  moveFiles,
  ...instantUploadAPI,
}
