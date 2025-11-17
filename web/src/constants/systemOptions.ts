import type { Composer } from '@/composables/useTexts'

/**
 * 系统配置选项常量
 * 用于系统设置、内容处理等功能配置
 */

/* ==================== 敏感内容处理选项 ==================== */

/**
 * 获取敏感内容处理选项
 */
export const getSensitiveContentOptions = ($t: Composer['$t']) => [
  { label: $t('constants.systemOptions.sensitiveContent.markOnly'), value: 'mark_only' as const },
  { label: $t('constants.systemOptions.sensitiveContent.pendingReview'), value: 'pending_review' as const },
  { label: $t('constants.systemOptions.sensitiveContent.autoDelete'), value: 'auto_delete' as const },
]

export type SensitiveContentOptionI18n = ReturnType<typeof getSensitiveContentOptions>[number]

/* ==================== 大小单位选项（技术单位，不需要翻译） ==================== */

export const SIZE_UNIT_OPTIONS = [
  { label: 'B', value: 'B' },
  { label: 'KB', value: 'KB' },
  { label: 'MB', value: 'MB' },
  { label: 'GB', value: 'GB' },
  { label: 'TB', value: 'TB' },
] as const

export type SizeUnitOption = (typeof SIZE_UNIT_OPTIONS)[number]

/* ==================== 快速大小预设（技术单位，不需要翻译） ==================== */

export const QUICK_SIZE_PRESETS = [
  { label: '1 MB', value: 1, unit: 'MB' },
  { label: '5 MB', value: 5, unit: 'MB' },
  { label: '10 MB', value: 10, unit: 'MB' },
  { label: '50 MB', value: 50, unit: 'MB' },
  { label: '100 MB', value: 100, unit: 'MB' },
  { label: '500 MB', value: 500, unit: 'MB' },
  { label: '1 GB', value: 1, unit: 'GB' },
  { label: '10 GB', value: 10, unit: 'GB' },
] as const

export type QuickSizePreset = (typeof QUICK_SIZE_PRESETS)[number]

/* ==================== 完整的大小单位定义（技术单位，不需要翻译） ==================== */

export const ALL_SIZE_UNITS = [
  { key: 'B', label: 'Bytes (B)', divisor: 1, precision: 0 },
  { key: 'KB', label: 'Kilobytes (KB)', divisor: 1024, precision: 2 },
  { key: 'MB', label: 'Megabytes (MB)', divisor: 1024 * 1024, precision: 2 },
  { key: 'GB', label: 'Gigabytes (GB)', divisor: 1024 * 1024 * 1024, precision: 4 },
  { key: 'TB', label: 'Terabytes (TB)', divisor: 1024 * 1024 * 1024 * 1024, precision: 6 },
] as const

export type AllSizeUnit = (typeof ALL_SIZE_UNITS)[number]

/* ==================== 用户角色选项 ==================== */

/**
 * 获取用户角色选项
 */
export const getUserRoleOptions = ($t: Composer['$t']) => [
  { label: $t('constants.systemOptions.userRoles.regular'), value: 3 },
  { label: $t('constants.systemOptions.userRoles.advanced'), value: 2 },
  { label: $t('constants.systemOptions.userRoles.admin'), value: 1 },
  { label: $t('constants.systemOptions.userRoles.superAdmin'), value: 0 },
]

export type UserRoleOptionI18n = ReturnType<typeof getUserRoleOptions>[number]

/* ==================== 用户状态选项 ==================== */

/**
 * 获取用户状态选项
 */
export const getUserStatusOptions = ($t: Composer['$t']) => [
  { label: $t('constants.systemOptions.userStatus.active'), value: 1 },
  { label: $t('constants.systemOptions.userStatus.disabled'), value: 0 },
  { label: $t('constants.systemOptions.userStatus.pending'), value: 2 },
]

export type UserStatusOptionI18n = ReturnType<typeof getUserStatusOptions>[number]

/* ==================== 文件状态选项 ==================== */

/**
 * 获取文件状态选项
 */
export const getImageStatusOptions = ($t: Composer['$t']) => [
  { label: $t('constants.systemOptions.imageStatus.active'), value: 1 },
  { label: $t('constants.systemOptions.imageStatus.hidden'), value: 0 },
  { label: $t('constants.systemOptions.imageStatus.pending'), value: 2 },
  { label: $t('constants.systemOptions.imageStatus.deleted'), value: -1 },
]

export type ImageStatusOptionI18n = ReturnType<typeof getImageStatusOptions>[number]

/* ==================== 分享状态选项 ==================== */

/**
 * 获取分享状态选项
 */
export const getShareStatusOptions = ($t: Composer['$t']) => [
  { label: $t('constants.systemOptions.shareStatus.active'), value: 1 },
  { label: $t('constants.systemOptions.shareStatus.disabled'), value: 0 },
  { label: $t('constants.systemOptions.shareStatus.expired'), value: -1 },
]

export type ShareStatusOptionI18n = ReturnType<typeof getShareStatusOptions>[number]
