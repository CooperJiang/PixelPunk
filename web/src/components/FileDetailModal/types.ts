/**
 * FileDetailModal组件类型定义
 */
import type { FileInfo } from '@/api/types'

/**
 * FileDetailModal组件Props
 */
export interface FileDetailModalProps {
  modelValue: boolean

  fileInfo: FileInfo | null

  isAdmin?: boolean
}

/**
 * FileDetailModal组件事件定义
 */
export interface FileDetailModalEmits {
  'update:modelValue': [value: boolean]

  delete: [id?: string]

  recommend: [id?: string, recommended?: boolean]
}

/**
 * 质量等级类型（仅使用规范化英文键）
 * UI 文案请通过 i18n：components.fileDetailModal.quality.{high|medium|low}
 */
export type QualityLevel = 'high' | 'medium' | 'low'

/**
 */
export type VisibilityType = 'public' | 'private' | 'protected' | 'link'

/**
 * 可见性翻译映射（已废弃，请使用 $t('components.fileDetailModal.visibility.xxx')）
 * 保留仅作为类型兼容
 */
export const VisibilityLabelMap = {
  public: 'public',
  private: 'private',
  protected: 'protected',
  link: 'link',
} as const
