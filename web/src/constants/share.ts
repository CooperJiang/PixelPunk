import type { Composer } from '@/composables/useTexts'

/**
 * 分享模块相关常量
 */

/* ==================== 排序选项 ==================== */

/**
 * 获取分享排序选项
 */
export const getShareSortOptions = ($t: Composer['$t']) => [
  { value: 'date' as const, label: $t('constants.share.sortOptions.date') },
  { value: 'name' as const, label: $t('constants.share.sortOptions.name') },
  { value: 'size' as const, label: $t('constants.share.sortOptions.size') },
]

export type ShareSortOptionI18n = ReturnType<typeof getShareSortOptions>[number]

/* ==================== 排序方向（技术常量，不需要翻译） ==================== */

export const SORT_DIRECTIONS = {
  ASC: 'asc',
  DESC: 'desc',
} as const

/* ==================== 布局选项 ==================== */

/**
 * 获取布局选项
 */
export const getLayoutOptions = ($t: Composer['$t']) => [
  {
    value: 'grid' as const,
    label: $t('constants.share.layoutOptions.grid.label'),
    icon: 'fas fa-th',
    description: $t('constants.share.layoutOptions.grid.description'),
  },
  {
    value: 'waterfall' as const,
    label: $t('constants.share.layoutOptions.waterfall.label'),
    icon: 'fas fa-stream',
    description: $t('constants.share.layoutOptions.waterfall.description'),
  },
  {
    value: 'masonry' as const,
    label: $t('constants.share.layoutOptions.masonry.label'),
    icon: 'fas fa-th-large',
    description: $t('constants.share.layoutOptions.masonry.description'),
  },
  {
    value: 'large' as const,
    label: $t('constants.share.layoutOptions.large.label'),
    icon: 'fas fa-image',
    description: $t('constants.share.layoutOptions.large.description'),
  },
]

export type LayoutOptionI18n = ReturnType<typeof getLayoutOptions>[number]

/* ==================== 布局模式（技术常量，不需要翻译） ==================== */

export const LAYOUT_MODES = {
  GRID: 'grid' as const,
  WATERFALL: 'waterfall' as const,
  MASONRY: 'masonry' as const,
  LARGE: 'large' as const,
}

/* ==================== 分页和加载相关（技术常量，不需要翻译） ==================== */

export const INITIAL_LOAD_COUNT = 24
export const PAGE_SIZE = 24

/* ==================== 存储键（技术常量，不需要翻译） ==================== */

export const STORAGE_KEYS = {
  GRID_SIZE: 'grid_size_preference',
  LAYOUT_MODE: 'share_layout_mode',
  DYNAMIC_BG: 'dynamic_bg_enabled',
  VISITOR_DIALOG_DISMISSED: 'visitor_dialog_dismissed_',
} as const

/* ==================== 默认值（技术常量，不需要翻译） ==================== */

export const DEFAULTS = {
  GRID_SIZE: 250,
} as const

/* ==================== 选择相关常量（技术常量，不需要翻译） ==================== */

export const SELECTION = {
  INDICATOR_SIZE: '22px',
  INDICATOR_BORDER_COLOR: 'rgba(255, 255, 255, 0.7)',
  SELECTED_BG_COLOR: 'rgba(5, 217, 232, 0.8)',
  SELECTED_BORDER_COLOR: '#ffffff',
} as const

/* ==================== 样式相关常量（技术常量，不需要翻译） ==================== */

export const STYLE = {
  PRIMARY_COLOR: '#05d9e8',
  PRIMARY_COLOR_DARK: '#04b4c0',
  PRIMARY_COLOR_LIGHT: 'rgba(5, 217, 232, 0.2)',
  SECONDARY_COLOR: '#ff2a6d',
  BACKGROUND_COLOR: 'rgba(15, 25, 35, 0.5)',
  CARD_BORDER_COLOR: 'rgba(255, 255, 255, 0.1)',
} as const
