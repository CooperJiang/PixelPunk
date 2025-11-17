/**
 * 分类模板相关常量
 */
import type { TranslationFunction } from '@/composables/useTexts'

/* 分类排序字段选项值 */
export const CATEGORY_SORT_BY_VALUES = ['sort_order', 'name', 'usage_count', 'created_at'] as const

/* 排序方向选项值 */
export const SORT_ORDER_VALUES = ['asc', 'desc'] as const

/* 分类类型值 */
export const CATEGORY_TYPE_VALUES = {
  ALL: undefined,
  HOT: true,
  NORMAL: false,
} as const

export const getCategorySortByOptions = ($t: TranslationFunction) => [
  { label: $t('constants.categorySortBy.sort_order'), value: 'sort_order' },
  { label: $t('constants.categorySortBy.name'), value: 'name' },
  { label: $t('constants.categorySortBy.usage_count'), value: 'usage_count' },
  { label: $t('constants.categorySortBy.created_at'), value: 'created_at' },
] as const

export const getSortOrderOptions = ($t: TranslationFunction) => [
  { label: $t('constants.sortOrder.asc'), value: 'asc' },
  { label: $t('constants.sortOrder.desc'), value: 'desc' },
] as const

export const getCategoryTypeOptions = ($t: TranslationFunction) => [
  { label: $t('constants.categoryType.all'), value: undefined },
  { label: $t('constants.categoryType.hot'), value: true },
  { label: $t('constants.categoryType.normal'), value: false },
] as const

export const CATEGORY_SORT_BY_OPTIONS = CATEGORY_SORT_BY_VALUES.map((value) => ({
  label: `constants.categorySortBy.${value}`,
  value,
}))

export const SORT_ORDER_OPTIONS = SORT_ORDER_VALUES.map((value) => ({
  label: `constants.sortOrder.${value}`,
  value,
}))

export const CATEGORY_TYPE_OPTIONS = [
  { label: 'constants.categoryType.all', value: undefined },
  { label: 'constants.categoryType.hot', value: true },
  { label: 'constants.categoryType.normal', value: false },
]

/* 默认排序设置 */
export const DEFAULT_CATEGORY_SORT = {
  sortBy: 'sort_order',
  sortOrder: 'asc',
} as const

/* 类型定义 */
export type CategorySortByOption = (typeof CATEGORY_SORT_BY_OPTIONS)[number]
export type SortOrderOption = (typeof SORT_ORDER_OPTIONS)[number]
export type CategoryTypeOption = (typeof CATEGORY_TYPE_OPTIONS)[number]
