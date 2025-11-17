/**
 * 标签辅助函数
 */
import { useTexts } from '@/composables/useTexts'

/**
 * 获取来源颜色
 */
export function getSourceColor(source: string): string {
  switch (source) {
    case 'manual':
      return 'purple'
    case 'ai':
      return 'green'
    case 'system':
      return 'blue'
    default:
      return 'gray'
  }
}

export function getSourceText(source: string): string {
  const { $t } = useTexts()
  switch (source) {
    case 'manual':
      return $t('tagManage.source.manual')
    case 'ai':
      return $t('tagManage.source.ai')
    case 'system':
      return $t('tagManage.source.system')
    default:
      return $t('tagManage.source.unknown')
  }
}

export function getSourceIcon(source: string): string {
  switch (source) {
    case 'manual':
      return 'fa-user'
    case 'ai':
      return 'fa-robot'
    case 'system':
      return 'fa-cube'
    default:
      return 'fa-tag'
  }
}
