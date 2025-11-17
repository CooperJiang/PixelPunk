import { del, get, post, put } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

export interface RandomImageAPI {
  id: number
  name: string
  api_key: string
  folder_id: string | null
  folder_name: string
  return_type: 'redirect' | 'direct'
  status: number
  is_active: boolean
  call_count: number
  last_called_at: string | null
  created_at: string
  updated_at: string
}

export interface CreateRandomAPIRequest {
  name: string
  folder_id: string | null
  return_type: 'redirect' | 'direct'
}

export interface UpdateRandomAPIStatusRequest {
  status: number
}

export interface UpdateRandomAPIConfigRequest {
  folder_id: string | null
  return_type: 'redirect' | 'direct'
}

export interface RandomAPIListParams {
  page?: number
  size?: number
  status?: number
  search?: string
}

export interface RandomAPIListResponse {
  items: RandomImageAPI[]
  total: number
  page: number
  size: number
}

export function createRandomAPI(data: CreateRandomAPIRequest): Promise<ApiResult<RandomImageAPI>> {
  return post<RandomImageAPI>('/random-api/create', data)
}

export function getRandomAPIList(params?: RandomAPIListParams): Promise<ApiResult<RandomAPIListResponse>> {
  return get<RandomAPIListResponse>('/random-api/list', params)
}

export function updateRandomAPIStatus(id: number, data: UpdateRandomAPIStatusRequest): Promise<ApiResult<null>> {
  return put<null>(`/random-api/${id}/status`, data)
}

export function updateRandomAPIConfig(id: number, data: UpdateRandomAPIConfigRequest): Promise<ApiResult<null>> {
  return put<null>(`/random-api/${id}/config`, data)
}

export function deleteRandomAPI(id: number): Promise<ApiResult<null>> {
  return del<null>(`/random-api/${id}`)
}
