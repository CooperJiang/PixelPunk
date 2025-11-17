import { post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/* 文件上传响应接口 */
export interface FileUploadResponse {
  file_url: string // 相对路径
  full_url: string // 完整URL
}

/* 管理员文件上传 */
export async function uploadAdminFile(file: File): Promise<ApiResult<FileUploadResponse>> {
  const formData = new FormData()
  formData.append('file', file)

  return post<FileUploadResponse>('/admin/file/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

export default {
  uploadAdminFile,
}
