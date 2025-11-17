/**
 * 分类辅助函数
 */
import { useTexts } from '@/composables/useTexts'

interface SourceInfo {
  text: string
  class: string
  icon: string
  style: {
    backgroundColor?: string
    color?: string
    borderColor?: string
  }
}

/**
 * 获取来源信息映射
 */
function getSourceMap($t: (key: string) => string): Record<string, SourceInfo> {
  return {
    user: {
      text: $t('category.sourceInfo.manual'),
      class: 'bg-purple-500/10 text-purple-400 border border-purple-500/30',
      icon: 'fas fa-user',
      style: {
        backgroundColor: 'var(--color-badge-accent-bg)',
        color: 'var(--color-badge-accent-text)',
        borderColor: 'var(--color-badge-accent-border)',
      },
    },
    ai_suggestion: {
      text: $t('category.sourceInfo.aiSuggestion'),
      class: 'bg-success-500/10 text-success-400 border border-success-500/30',
      icon: 'fas fa-robot',
      style: {
        backgroundColor: 'var(--color-badge-success-bg)',
        color: 'var(--color-badge-success-text)',
        borderColor: 'var(--color-badge-success-border)',
      },
    },
    system_template: {
      text: $t('category.sourceInfo.system'),
      class: 'bg-brand-500/10 text-brand-400 border border-brand-500/30',
      icon: 'fas fa-cube',
      style: {
        backgroundColor: 'var(--color-badge-primary-bg)',
        color: 'var(--color-badge-primary-text)',
        borderColor: 'var(--color-badge-primary-border)',
      },
    },
    imported: {
      text: $t('category.sourceInfo.imported'),
      class: 'bg-info-500/10 text-info-400 border border-info-500/30',
      icon: 'fas fa-file-import',
      style: {
        backgroundColor: 'var(--color-badge-info-bg)',
        color: 'var(--color-badge-info-text)',
        borderColor: 'var(--color-badge-info-border)',
      },
    },
  }
}

export function getSourceInfo(source: string): SourceInfo {
  const { $t } = useTexts()
  const sourceMap = getSourceMap($t)
  return (
    sourceMap[source] || {
      text: source,
      class: 'bg-background-700 text-content-muted border border-border-default',
      icon: 'fas fa-question',
      style: {},
    }
  )
}

export function getStatusStyle(status: string): {
  containerStyle: {
    backgroundColor: string
    color: string
    borderColor: string
  }
  dotStyle: {
    backgroundColor: string
  }
  text: string
} {
  const { $t } = useTexts()
  const isActive = status === 'active'
  return {
    containerStyle: isActive
      ? {
          backgroundColor: 'var(--color-badge-success-bg)',
          color: 'var(--color-badge-success-text)',
          borderColor: 'var(--color-badge-success-border)',
        }
      : {
          backgroundColor: 'var(--color-badge-neutral-bg)',
          color: 'var(--color-badge-neutral-text)',
          borderColor: 'var(--color-badge-neutral-border)',
        },
    dotStyle: {
      backgroundColor: isActive ? 'var(--color-badge-success-text)' : 'var(--color-badge-neutral-text)',
    },
    text: $t(`category.statusInfo.${isActive ? 'active' : 'archived'}`),
  }
}
