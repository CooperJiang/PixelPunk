/**
 * 公告管理 API

 */
import { post, get, put, del } from '@/utils/network/http'
import type {
  Announcement,
  CreateAnnouncementRequest,
  UpdateAnnouncementRequest,
  AnnouncementListQuery,
  AnnouncementListResponse,
  AnnouncementSettings,
  AnnouncementSettingsUpdate,
} from '@/api/types/announcement'
import type { ApiResponse } from '@/api/types/common'

/**
 */
export function createAnnouncement(data: CreateAnnouncementRequest): Promise<ApiResponse<Announcement>> {
  return post<Announcement>('/admin/announcements', data)
}

export function getAnnouncement(id: number): Promise<ApiResponse<Announcement>> {
  return get<Announcement>(`/admin/announcements/${id}`)
}

export function updateAnnouncement(id: number, data: Omit<UpdateAnnouncementRequest, 'id'>): Promise<ApiResponse<Announcement>> {
  return put<Announcement>(`/admin/announcements/${id}`, data)
}

export function deleteAnnouncement(id: number): Promise<ApiResponse<void>> {
  return del<void>(`/admin/announcements/${id}`)
}

export function getAnnouncementList(params: AnnouncementListQuery): Promise<ApiResponse<AnnouncementListResponse>> {
  return get<AnnouncementListResponse>('/admin/announcements', params)
}

export function getAnnouncementSettings(): Promise<ApiResponse<AnnouncementSettings>> {
  return get<AnnouncementSettings>('/admin/announcements/settings')
}

export function updateAnnouncementSettings(data: AnnouncementSettingsUpdate): Promise<ApiResponse<{ message: string }>> {
  return put<{ message: string }>('/admin/announcements/settings', data)
}

export function toggleAnnouncementPin(id: number): Promise<ApiResponse<Announcement>> {
  return put<Announcement>(`/admin/announcements/${id}/toggle-pin`, {})
}

export default {
  createAnnouncement,
  getAnnouncement,
  updateAnnouncement,
  deleteAnnouncement,
  getAnnouncementList,
  getAnnouncementSettings,
  updateAnnouncementSettings,
  toggleAnnouncementPin,
}
