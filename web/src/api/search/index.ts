import { get, post } from '@/utils/network/http'
import type { ImageInfo } from '@/types/global'

export type { ImageInfo }

/* 向量搜索请求参数 */
export interface VectorSearchRequest {
  query: string
  limit?: number
  threshold?: number
}

/* 用户个人文件向量搜索请求参数 */
export interface UserVectorSearchRequest {
  query: string
  page?: number
}

/* Gallery公开文件向量搜索请求参数 */
export interface GalleryVectorSearchRequest {
  query: string
  page?: number
}

/* 向量搜索结果 */
export interface VectorSearchResult {
  file_id: string
  similarity: number
  description: string
  file_info?: ImageInfo
}

/* 向量搜索响应 */
export interface VectorSearchResponse {
  query: string
  total: number
  results: VectorSearchResult[]
  process_time: string
  used_cache: boolean
}

/* 带分页的向量搜索结果（兼容现有文件接口格式） */
export interface PaginatedVectorSearchResponse {
  items: Array<{
    [key: string]: unknown // Image模型字段
    similarity: number // 相似度字段
  }>
  pagination: {
    total: number
    size: number
    current_page: number
    last_page: number
  }
  search_info: {
    query: string
    threshold: number
    process_time: string
    used_cache: boolean
  }
}

/* 向量统计信息 */
export interface VectorStatsResponse {
  total_vectors: number
  completed_vectors: number
  pending_vectors: number
  failed_vectors: number
  user_vectors: number
  storage_size: number
  avg_dimension: number
}

/* 向量引擎健康状态 */
export interface VectorHealthResponse {
  status: string
  enabled: boolean
  provider: string
  model: string
  database_ok: boolean
  embedding_ok: boolean
  message: string
}

/* 向量搜索文件（测试用） */
export function vectorSearch(data: VectorSearchRequest): Promise<VectorSearchResponse> {
  return post<VectorSearchResponse>('/search/vector/search', data)
}

export function getVectorStats(): Promise<VectorStatsResponse> {
  return get<VectorStatsResponse>('/search/vector/stats')
}

export function getVectorHealth(): Promise<VectorHealthResponse> {
  return get<VectorHealthResponse>('/search/vector/health')
}

export function searchSimilarImages(fileId: string): Promise<PaginatedVectorSearchResponse> {
  return get<PaginatedVectorSearchResponse>(`/search/vector/similar/${fileId}`, undefined, { autoShowError: false })
}

export function gallerySimilarImages(fileId: string): Promise<PaginatedVectorSearchResponse> {
  return get<PaginatedVectorSearchResponse>(`/search/gallery/similar/${fileId}`, undefined, { autoShowError: false })
}

export function userSimilarImages(fileId: string): Promise<PaginatedVectorSearchResponse> {
  return get<PaginatedVectorSearchResponse>(`/search/user/similar/${fileId}`, undefined, { autoShowError: false })
}

export function adminSimilarImages(fileId: string): Promise<PaginatedVectorSearchResponse> {
  return get<PaginatedVectorSearchResponse>(`/search/admin/similar/${fileId}`, undefined, { autoShowError: false })
}

export function userVectorSearch(data: UserVectorSearchRequest): Promise<PaginatedVectorSearchResponse> {
  return post<PaginatedVectorSearchResponse>('/search/user/vector/search', data)
}

export function galleryVectorSearch(data: GalleryVectorSearchRequest): Promise<PaginatedVectorSearchResponse> {
  return post<PaginatedVectorSearchResponse>('/search/gallery/vector/search', data)
}

export default {
  vectorSearch,
  getVectorStats,
  getVectorHealth,
  searchSimilarImages,
  gallerySimilarImages,
  userSimilarImages,
  adminSimilarImages,
  userVectorSearch,
  galleryVectorSearch,
}
