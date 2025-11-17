/**
 * 瀑布流组件属性
 */
export interface WaterfallLayoutProps {
  items: any[]

  columnCount?: number

  columnWidth?: number

  selectable?: boolean

  selectedIds?: string[]

  lazyLoad?: boolean

  initialLoadCount?: number

  gap?: number

  shareKey?: string

  accessToken?: string
}

/**
 * 瀑布流组件事件
 */
export interface WaterfallLayoutEmits {
  (e: 'image-click', item: any): void

  (e: 'select', itemId: string): void

  (e: 'image-load', item: any): void
}

export interface WaterfallItemProps {
  item: any
  selectable: boolean
  selectedIds: string[]
  columnWidth: number
  gap: number
}

export interface WaterfallItemEmits {
  (e: 'image-click', item: any): void
  (e: 'image-load', item: any): void
  (e: 'image-error', item: any): void
  (e: 'toggle-tags', id: string): void
}
