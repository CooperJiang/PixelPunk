/**
 * 访问级别工具函数
 * 支持多语言的访问级别常量生成

 */
import { useTexts } from '@/composables/useTexts'
import type { AccessLevel } from '@/components/AccessLevelToggle/types'

/**
 * 获取访问级别标题映射（支持主题切换）
 * 替代原来的 ACCESS_LEVEL_TITLES 常量
 */
export function getAccessLevelTitles() {
  const { $t } = useTexts()

  return {
    image: {
      private: `${$t('access.levels.private')}${$t('access.types.image')} - ${$t('access.actions.toggle')}`,
      protected: `${$t('access.levels.protected')}${$t('access.types.image')} - ${$t('access.actions.toggle')}`,
      public: `${$t('access.levels.public')}${$t('access.types.image')} - ${$t('access.actions.toggle')}`,
    },
    folder: {
      private: `${$t('access.levels.private')}${$t('access.types.folder')} - ${$t('access.actions.toggle')}`,
      protected: `${$t('access.levels.protected')}${$t('access.types.folder')} - ${$t('access.actions.toggle')}`,
      public: `${$t('access.levels.public')}${$t('access.types.folder')} - ${$t('access.actions.toggle')}`,
    },
  }
}

export function getAccessLevelTitle(type: 'image' | 'folder', level: AccessLevel): string {
  const { $t } = useTexts()
  return `${$t(`access.levels.${level}`)}${$t(`access.types.${type}`)} - ${$t('access.actions.toggle')}`
}

export const ACCESS_LEVEL_ICONS = {
  private: 'fas fa-lock',
  protected: 'fas fa-shield-alt',
  public: 'fas fa-globe',
} as const
