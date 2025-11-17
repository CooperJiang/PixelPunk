import { get } from '@/utils/network/http'

export interface SystemStatsResponse {
  totalImages: number
  totalStorage?: number
  todayUploads?: number
  totalUsers?: number
  activeUsers?: number
}

export interface PublicStatsResponse {
  total: number
  timestamp: number
}

/**
 * 获取系统统计数据（用于首页展示）
 */
export function getSystemStats(): Promise<{ data: PublicStatsResponse }> {
  return get<PublicStatsResponse>('/pb/stats/files/count')
}

export default {
  getSystemStats,
}
