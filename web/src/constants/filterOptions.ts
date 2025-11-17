/**
 * 筛选选项相关常量
 * 用于文件筛选、排序、搜索等功能
 */

import type { Composer } from '@/composables/useTexts'

/* 排序选项 */
export const getFilterSortOptions = ($t: Composer['$t']) => [
  { label: $t('constants.filterOptions.sortOptions.newest'), value: 'newest' as const },
  { label: $t('constants.filterOptions.sortOptions.oldest'), value: 'oldest' as const },
  { label: $t('constants.filterOptions.sortOptions.name'), value: 'name' as const },
  { label: $t('constants.filterOptions.sortOptions.size'), value: 'size' as const },
  { label: $t('constants.filterOptions.sortOptions.width'), value: 'width' as const },
  { label: $t('constants.filterOptions.sortOptions.height'), value: 'height' as const },
  { label: $t('constants.filterOptions.sortOptions.quality'), value: 'quality' as const },
]

/* 管理员排序选项（包含敏感度评分） */
export const getAdminSortOptions = ($t: Composer['$t']) => [
  ...getFilterSortOptions($t),
  { label: $t('constants.filterOptions.sortOptions.nsfw_score'), value: 'nsfw_score' as const },
]

/* 分辨率选项 */
export const getResolutionOptions = ($t: Composer['$t']) => [
  { label: $t('constants.filterOptions.resolutionOptions.low'), value: '低分辨率' as const }, // value 为后端API枚举值，不可修改
  { label: $t('constants.filterOptions.resolutionOptions.720p'), value: '720p' as const },
  { label: $t('constants.filterOptions.resolutionOptions.1080p'), value: '1080p' as const },
  { label: $t('constants.filterOptions.resolutionOptions.2k'), value: '2k' as const },
  { label: $t('constants.filterOptions.resolutionOptions.4k'), value: '4k' as const },
  { label: $t('constants.filterOptions.resolutionOptions.8k'), value: '8k' as const },
]

/* 文件格式选项 */
export const getImageFormatOptions = ($t: Composer['$t']) => [
  { label: $t('constants.filterOptions.imageFormatOptions.jpg'), value: 'jpg' as const },
  { label: $t('constants.filterOptions.imageFormatOptions.jpeg'), value: 'jpeg' as const },
  { label: $t('constants.filterOptions.imageFormatOptions.png'), value: 'png' as const },
  { label: $t('constants.filterOptions.imageFormatOptions.gif'), value: 'gif' as const },
  { label: $t('constants.filterOptions.imageFormatOptions.webp'), value: 'webp' as const },
  { label: $t('constants.filterOptions.imageFormatOptions.bmp'), value: 'bmp' as const },
  { label: $t('constants.filterOptions.imageFormatOptions.svg'), value: 'svg' as const },
  { label: $t('constants.filterOptions.imageFormatOptions.ico'), value: 'ico' as const },
]

/* 文件类型选项 */
export const getFileTypeOptions = ($t: Composer['$t']) => [
  { label: $t('constants.filterOptions.fileTypeOptions.all'), value: '*' as const },
  { label: $t('constants.filterOptions.fileTypeOptions.jpeg'), value: 'image/jpeg' as const },
  { label: $t('constants.filterOptions.fileTypeOptions.png'), value: 'image/png' as const },
  { label: $t('constants.filterOptions.fileTypeOptions.gif'), value: 'image/gif' as const },
  { label: $t('constants.filterOptions.fileTypeOptions.webp'), value: 'image/webp' as const },
  { label: $t('constants.filterOptions.fileTypeOptions.svg'), value: 'image/svg+xml' as const },
  { label: $t('constants.filterOptions.fileTypeOptions.bmp'), value: 'image/bmp' as const },
  { label: $t('constants.filterOptions.fileTypeOptions.ico'), value: 'image/x-icon' as const },
]

/* 类型定义 */
export type FilterSortOption = ReturnType<typeof getFilterSortOptions>[number]
export type AdminSortOption = ReturnType<typeof getAdminSortOptions>[number]
export type ResolutionOption = ReturnType<typeof getResolutionOptions>[number]
export type ImageFormatOption = ReturnType<typeof getImageFormatOptions>[number]
export type FileTypeOption = ReturnType<typeof getFileTypeOptions>[number]
