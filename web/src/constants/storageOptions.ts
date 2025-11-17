import type { Composer } from '@/composables/useTexts'

/**
 * 存储相关选项常量
 * 用于存储渠道、存储类型配置等
 */

/* ==================== 存储类型选项 ==================== */

export const getStorageTypeOptions = ($t: Composer['$t']) => [
  { label: $t('constants.storageOptions.storageTypes.local'), value: 'local' as const },
  { label: $t('constants.storageOptions.storageTypes.s3'), value: 's3' as const },
  { label: $t('constants.storageOptions.storageTypes.aliyun'), value: 'aliyun' as const },
  { label: $t('constants.storageOptions.storageTypes.tencent'), value: 'tencent' as const },
]

export type StorageTypeOptionI18n = ReturnType<typeof getStorageTypeOptions>[number]

/* ==================== 渠道类型选项 ==================== */

export const getChannelTypeOptions = ($t: Composer['$t']) => [
  { label: $t('constants.storageOptions.channelTypes.oss'), value: 'oss' as const },
  { label: $t('constants.storageOptions.channelTypes.cos'), value: 'cos' as const },
  { label: $t('constants.storageOptions.channelTypes.rainyun'), value: 'rainyun' as const },
  { label: $t('constants.storageOptions.channelTypes.local'), value: 'local' as const },
]

export type ChannelTypeOptionI18n = ReturnType<typeof getChannelTypeOptions>[number]

/* ==================== 存储类别选项 ==================== */

export const getStorageClassOptions = ($t: Composer['$t']) => ({
  cos: [
    { label: $t('constants.storageOptions.storageClasses.standard'), value: 'STANDARD' as const },
    { label: $t('constants.storageOptions.storageClasses.standardMAZ'), value: 'MAZ_STANDARD' as const },
    { label: $t('constants.storageOptions.storageClasses.standardIA'), value: 'STANDARD_IA' as const },
    { label: $t('constants.storageOptions.storageClasses.standardIAMAZ'), value: 'MAZ_STANDARD_IA' as const },
    { label: $t('constants.storageOptions.storageClasses.intelligentTiering'), value: 'INTELLIGENT_TIERING' as const },
    {
      label: $t('constants.storageOptions.storageClasses.intelligentTieringMAZ'),
      value: 'MAZ_INTELLIGENT_TIERING' as const,
    },
    { label: $t('constants.storageOptions.storageClasses.archive'), value: 'ARCHIVE' as const },
    { label: $t('constants.storageOptions.storageClasses.deepArchive'), value: 'DEEP_ARCHIVE' as const },
  ],
  oss: [
    { label: $t('constants.storageOptions.storageClasses.standard'), value: 'STANDARD' as const },
    { label: $t('constants.storageOptions.storageClasses.ia'), value: 'IA' as const },
    { label: $t('constants.storageOptions.storageClasses.archive'), value: 'ARCHIVE' as const },
    { label: $t('constants.storageOptions.storageClasses.coldArchive'), value: 'COLD_ARCHIVE' as const },
  ],
  rainyun: [
    { label: $t('constants.storageOptions.storageClasses.standard'), value: 'STANDARD' as const },
    { label: $t('constants.storageOptions.storageClasses.standardIA'), value: 'STANDARD_IA' as const },
    { label: $t('constants.storageOptions.storageClasses.glacier'), value: 'GLACIER' as const },
    { label: $t('constants.storageOptions.storageClasses.deepArchive'), value: 'DEEP_ARCHIVE' as const },
  ],
})

export type StorageClassOptionsI18n = ReturnType<typeof getStorageClassOptions>

/* ==================== 存储大小选项 ==================== */

export const getStorageSizeOptions = ($t: Composer['$t']) => [
  { label: $t('constants.storageOptions.unlimited'), value: 0 },
  { label: '200 MB', value: 200 * 1024 * 1024 },
  { label: '500 MB', value: 500 * 1024 * 1024 },
  { label: '1 GB', value: 1 * 1024 * 1024 * 1024 },
  { label: '5 GB', value: 5 * 1024 * 1024 * 1024 },
  { label: '10 GB', value: 10 * 1024 * 1024 * 1024 },
  { label: '20 GB', value: 20 * 1024 * 1024 * 1024 },
  { label: '50 GB', value: 50 * 1024 * 1024 * 1024 },
  { label: '1 TB', value: 1024 * 1024 * 1024 * 1024 },
] as const

export type StorageSizeOptionI18n = ReturnType<typeof getStorageSizeOptions>[number]

/* ==================== 文件大小选项 ==================== */

export const getFileSizeOptions = ($t: Composer['$t']) => [
  { label: $t('constants.storageOptions.unlimited'), value: 0 },
  { label: '1 MB', value: 1 * 1024 * 1024 },
  { label: '5 MB', value: 5 * 1024 * 1024 },
  { label: '10 MB', value: 10 * 1024 * 1024 },
  { label: '20 MB', value: 20 * 1024 * 1024 },
] as const

export type FileSizeOptionI18n = ReturnType<typeof getFileSizeOptions>[number]
