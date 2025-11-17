/**
 * 公告公开 API（用户端）

 */
import { get } from '@/utils/network/http'
import type { PublicAnnouncementListResponse, AnnouncementDetail, AnnouncementSettings } from '@/api/types/announcement'
import type { ApiResponse } from '@/api/types/common'

/**
 * 获取公开的公告列表
 */
export function getPublicAnnouncementList(): Promise<ApiResponse<PublicAnnouncementListResponse>> {
  return get<PublicAnnouncementListResponse>('/announcements')
}

export function getPublicAnnouncementDetail(id: number): Promise<ApiResponse<AnnouncementDetail>> {
  return get<AnnouncementDetail>(`/announcements/${id}`)
}

export function getPublicAnnouncementSettings(): Promise<ApiResponse<AnnouncementSettings>> {
  return get<AnnouncementSettings>('/announcements/settings')
}

export default {
  getPublicAnnouncementList,
  getPublicAnnouncementDetail,
  getPublicAnnouncementSettings,
}
