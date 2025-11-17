import { get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/* 分片上传管理API */
export const chunkedUploadManagementAPI = {
  manualCleanup(): Promise<ApiResult<void>> {
    return post<void>('/files/chunked/admin/cleanup')
  },

  getStats(): Promise<
    ApiResult<{
      activeSessions: number
      expiredSessions: number
      totalSessions: number
      cleanupHistory: Array<{
        date: string
        cleanedCount: number
      }>
    }>
  > {
    return get('/files/chunked/admin/stats')
  },
}

export default chunkedUploadManagementAPI
