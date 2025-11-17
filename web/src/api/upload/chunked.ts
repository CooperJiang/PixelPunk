import { del, get, post } from '@/utils/network/http'

/* 分片上传相关类型定义 */
export interface InitChunkedUploadRequest {
  file_name: string
  file_size: number
  file_md5: string
  mime_type: string
  chunk_size: number
  folder_id?: string
  access_level?: 'public' | 'private' | 'protected'
  optimize?: boolean
  watermark_config?: any
}

export interface ChunkedUploadResponse {
  session_id: string
  status: string
  progress: number
  total_chunks: number
  uploaded_chunks: number
  message: string
}

export interface ChunkedUploadCompleteResponse {
  id: string
  url: string
  thumb_url: string
  full_url: string
  full_thumb_url: string
  original_name: string
  display_name: string
  size: number
  width: number
  height: number
  format: string
  access_level: string
  folder_id?: string
  created_at: string
  updated_at: string
  views?: number
  is_duplicate?: boolean
  md5_hash?: string
  is_recommended: boolean
  storage_provider_id?: string
}

export interface UploadChunkRequest {
  session_id: string
  chunk_number: number
  chunk_md5: string
}

/* 分片上传API */
export const chunkedUploadAPI = {
  async init(data: InitChunkedUploadRequest): Promise<ChunkedUploadResponse> {
    const result = await post<ChunkedUploadResponse>('/files/chunked/init', data)
    return result.data
  },

  async uploadChunk(data: UploadChunkRequest, chunkFile: Blob): Promise<ChunkedUploadResponse> {
    const formData = new FormData()
    formData.append('session_id', data.session_id)
    formData.append('chunk_number', data.chunk_number.toString())
    formData.append('chunk_md5', data.chunk_md5)
    formData.append('file', chunkFile)

    const result = await post<ChunkedUploadResponse>('/files/chunked/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return result.data
  },

  async complete(sessionId: string): Promise<ChunkedUploadCompleteResponse> {
    const result = await post<ChunkedUploadCompleteResponse>('/files/chunked/complete', {
      session_id: sessionId,
    })
    return result.data
  },

  async getStatus(sessionId: string): Promise<ChunkedUploadResponse> {
    const result = await get<ChunkedUploadResponse>('/files/chunked/status', {
      session_id: sessionId,
    })
    return result.data
  },

  async cancel(sessionId: string): Promise<{ session_id: string; status: string }> {
    const result = await del<{ session_id: string; status: string }>('/files/chunked/cancel', {
      session_id: sessionId,
    })
    return result.data
  },
}

export default chunkedUploadAPI
