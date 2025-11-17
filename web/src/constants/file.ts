/**
 * 文件相关常量定义
 */

import type { Composer } from '@/composables/useTexts'

/* ==================== 文件类型定义 ==================== */

export type FileType = 'image' | 'video' | 'document' | 'archive' | 'audio' | 'other'

/* 文件类型常量 */
export const FILE_TYPES = {
  IMAGE: 'image' as const,
  VIDEO: 'video' as const,
  DOCUMENT: 'document' as const,
  ARCHIVE: 'archive' as const,
  AUDIO: 'audio' as const,
  OTHER: 'other' as const,
} as const

/* ==================== 文件类型显示映射 ==================== */

export const getFileTypeNames = ($t: Composer['$t']): Record<FileType, string> => ({
  image: $t('constants.fileTypes.image'),
  video: $t('constants.fileTypes.video'),
  document: $t('constants.fileTypes.document'),
  archive: $t('constants.fileTypes.archive'),
  audio: $t('constants.fileTypes.audio'),
  other: $t('constants.fileTypes.other'),
})

/* 文件类型图标 */
export const FILE_TYPE_ICONS: Record<FileType, string> = {
  image: 'fas fa-image',
  video: 'fas fa-video',
  document: 'fas fa-file-alt',
  archive: 'fas fa-file-archive',
  audio: 'fas fa-file-audio',
  other: 'fas fa-file',
}

/* 文件类型颜色主题 */
export const FILE_TYPE_COLORS: Record<FileType, string> = {
  image: 'blue',
  video: 'purple',
  document: 'green',
  archive: 'orange',
  audio: 'pink',
  other: 'gray',
}

/* 文件类型CSS类名 */
export const FILE_TYPE_CLASSES: Record<FileType, string> = {
  image: 'file-type-image',
  video: 'file-type-video',
  document: 'file-type-document',
  archive: 'file-type-archive',
  audio: 'file-type-audio',
  other: 'file-type-other',
}

/* ==================== 支持的文件格式 ==================== */

/* 支持的文件格式映射 */
export const SUPPORTED_FILE_FORMATS = {
  image: ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp', 'ico'],
  video: ['mp4', 'avi', 'mov', 'wmv', 'flv', 'mkv', 'webm', 'm4v'],
  document: ['pdf', 'doc', 'docx', 'txt', 'rtf', 'odt', 'xls', 'xlsx', 'ppt', 'pptx'],
  archive: ['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz'],
  audio: ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma'],
} as const

/* 所有支持的格式（扁平化） */
export const ALL_SUPPORTED_FORMATS = Object.values(SUPPORTED_FILE_FORMATS).flat()

/* ==================== 文件大小限制 ==================== */

/* 文件大小限制 (bytes) */
export const FILE_SIZE_LIMITS: Record<FileType, number> = {
  image: 50 * 1024 * 1024, // 50MB
  video: 500 * 1024 * 1024, // 500MB
  document: 100 * 1024 * 1024, // 100MB
  archive: 200 * 1024 * 1024, // 200MB
  audio: 100 * 1024 * 1024, // 100MB
  other: 50 * 1024 * 1024, // 50MB
}

/* 默认文件大小限制 */
export const DEFAULT_FILE_SIZE_LIMIT = 50 * 1024 * 1024 // 50MB

/* ==================== 文件状态常量 ==================== */

export const FILE_STATUS = {
  ACTIVE: 'active',
  DELETED: 'deleted',
  PENDING: 'pending',
  PROCESSING: 'processing',
  FAILED: 'failed',
} as const

export type FileStatus = (typeof FILE_STATUS)[keyof typeof FILE_STATUS]

/* ==================== 访问级别常量 ==================== */

export const ACCESS_LEVELS = {
  PRIVATE: 'private',
  PUBLIC: 'public',
  PROTECTED: 'protected',
} as const

export type AccessLevel = (typeof ACCESS_LEVELS)[keyof typeof ACCESS_LEVELS]

export const getAccessLevelNames = ($t: Composer['$t']): Record<AccessLevel, string> => ({
  private: $t('constants.accessLevels.private'),
  public: $t('constants.accessLevels.public'),
  protected: $t('constants.accessLevels.protected'),
})

/* 访问级别图标 */
/**
 * @deprecated 已迁移到 @/utils/accessLevel，请从那里导入
 */
export const ACCESS_LEVEL_ICONS: Record<AccessLevel, string> = {
  private: 'fas fa-lock',
  public: 'fas fa-globe',
  protected: 'fas fa-shield-alt',
}

/* ==================== 排序选项 ==================== */

export const getSortOptions = ($t: Composer['$t']) => [
  { label: $t('constants.sortOptions.created_at'), value: 'created_at' as const },
  { label: $t('constants.sortOptions.updated_at'), value: 'updated_at' as const },
  { label: $t('constants.sortOptions.size'), value: 'size' as const },
  { label: $t('constants.sortOptions.name'), value: 'name' as const },
  { label: $t('constants.sortOptions.views'), value: 'views' as const },
  { label: $t('constants.sortOptions.downloads'), value: 'downloads' as const },
]

export type SortOption = ReturnType<typeof getSortOptions>[number]['value']

export const SORT_DIRECTIONS = {
  ASC: 'asc',
  DESC: 'desc',
} as const

export type SortDirection = (typeof SORT_DIRECTIONS)[keyof typeof SORT_DIRECTIONS]

/* ==================== 批量操作类型 ==================== */

/* 批量操作类型 */
export const BATCH_OPERATIONS = {
  DELETE: 'delete',
  MOVE: 'move',
  UPDATE_ACCESS: 'update_access',
  UPDATE_FOLDER: 'update_folder',
  DOWNLOAD: 'download',
} as const

export type BatchOperation = (typeof BATCH_OPERATIONS)[keyof typeof BATCH_OPERATIONS]

export const getBatchOperationNames = ($t: Composer['$t']): Record<BatchOperation, string> => ({
  delete: $t('constants.batchOperations.delete'),
  move: $t('constants.batchOperations.move'),
  update_access: $t('constants.batchOperations.update_access'),
  update_folder: $t('constants.batchOperations.update_folder'),
  download: $t('constants.batchOperations.download'),
})

/* ==================== 工具函数 ==================== */

/**
 * 根据文件扩展名获取文件类型
 */
export function getFileTypeByExtension(extension: string): FileType {
  const ext = extension.toLowerCase().replace('.', '')

  for (const [type, formats] of Object.entries(SUPPORTED_FILE_FORMATS)) {
    if (formats.includes(ext)) {
      return type as FileType
    }
  }

  return FILE_TYPES.OTHER
}

export function isFileFormatSupported(extension: string): boolean {
  const ext = extension.toLowerCase().replace('.', '')
  return ALL_SUPPORTED_FORMATS.includes(ext)
}

export function getFileSizeLimit(fileType: FileType): number {
  return FILE_SIZE_LIMITS[fileType] || DEFAULT_FILE_SIZE_LIMIT
}

export function isFileSizeValid(size: number, fileType: FileType): boolean {
  return size <= getFileSizeLimit(fileType)
}

/**
 * 获取文件类型信息
 * @param fileType 文件类型
 * @param $t 翻译函数
 */
export function getFileTypeInfo(fileType: FileType, $t: Composer['$t']) {
  const names = getFileTypeNames($t)

  return {
    name: names[fileType],
    icon: FILE_TYPE_ICONS[fileType],
    color: FILE_TYPE_COLORS[fileType],
    className: FILE_TYPE_CLASSES[fileType],
    sizeLimit: FILE_SIZE_LIMITS[fileType],
    supportedFormats: SUPPORTED_FILE_FORMATS[fileType as keyof typeof SUPPORTED_FILE_FORMATS] || [],
  }
}
