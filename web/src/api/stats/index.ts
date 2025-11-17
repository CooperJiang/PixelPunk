import { get } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { UserStatsResponse } from '../types'

export function getUserStats(): Promise<ApiResult<UserStatsResponse>> {
  return get<UserStatsResponse>('/stats')
}

export default {
  getUserStats,
}
