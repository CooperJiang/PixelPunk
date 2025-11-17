/**
 * UI 相关类型定义
 * 包含主题、布局、菜单、表格等UI组件相关类型

/* 主题类型 */
export type Theme = 'light' | 'dark' | 'auto'

/* 语言类型 */
export type Language = 'zh-CN' | 'en-US' | 'ja-JP'

/* 布局类型 */
export type Layout = 'grid' | 'list' | 'masonry'

/* 排序选项 */
export interface SortOption {
  field: string
  order: 'asc' | 'desc'
  label: string
}

/* 过滤选项 */
export interface FilterOption {
  key: string
  value: any
  label: string
  type: 'select' | 'input' | 'date' | 'range'
  options?: Array<{ label: string; value: any }>
}

/* 表格列定义 */
export interface TableColumn<T = any> {
  key: keyof T
  title: string
  width?: number | string
  align?: 'left' | 'center' | 'right'
  sortable?: boolean
  filterable?: boolean
  render?: (value: any, record: T, index: number) => any
}

/* 菜单项 */
export interface MenuItem {
  key: string
  label: string
  icon?: string
  path?: string
  children?: MenuItem[]
  permissions?: string[]
  meta?: Record<string, any>
}

/* 面包屑项 */
export interface BreadcrumbItem {
  label: string
  path?: string
  icon?: string
}
