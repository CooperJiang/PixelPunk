import { get } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'

/* 活动日志类型定义 */
export interface ActivityLog {
  id: number
  type: string
  // ❌ 已删除 content 字段 - 多语言由前端动态生成
  data: Record<string, any> // 结构化参数数据
  module: string
  entity_type: string
  entity_id: string
  is_visible: boolean
  tags: string
  created_at: string
}

/* 统计数据 */
export interface TodayStats {
  today_uploads: number
  total_views: number
}

/* 活动日志列表响应 */
export interface ActivityLogsResponse {
  list: ActivityLog[]
  total: number
  page: number
  size: number
  today_stats: TodayStats
}

export interface GetActivitiesParams {
  page?: number
  size?: number
}

/**
 * 获取当前用户的活动日志列表
 */
export function getUserActivities(params?: GetActivitiesParams): Promise<ApiResult<ActivityLogsResponse>> {
  const queryParams = {
    page: params?.page || 1,
    size: params?.size || 20,
  }

  return get<ActivityLogsResponse>('/personal/activities', queryParams)
}

export default {
  getUserActivities,
}
