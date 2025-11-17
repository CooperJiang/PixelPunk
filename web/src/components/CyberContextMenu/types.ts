export interface ContextMenuItem {
  key: string
  label: string
  icon?: string // FontAwesome 图标类名
  disabled?: boolean // 禁用状态
  divided?: boolean // 显示分隔线
  children?: ContextMenuItem[] // 子菜单
  shortcut?: string // 快捷键显示
  danger?: boolean // 危险操作样式
  onClick?: (item: ContextMenuItem, event: MouseEvent) => void
  hasAsyncChildren?: boolean // 是否有异步子菜单
  loadChildren?: () => Promise<ContextMenuItem[]> // 异步加载子菜单函数
  isLoading?: boolean // 加载状态
  loadError?: string // 加载错误信息
}

export interface ContextMenuOptions {
  items: ContextMenuItem[]
  trigger?: 'contextmenu' | 'click' // 触发方式
  placement?: 'auto' | 'top' | 'bottom' | 'left' | 'right'
  offset?: { x: number; y: number } // 偏移量
  zIndex?: number // 层级
  className?: string // 自定义样式类
}

export interface CyberContextMenuProps {
  modelValue: boolean
  items: ContextMenuItem[]
  x?: number // 显示位置 X 坐标
  y?: number // 显示位置 Y 坐标
  trigger?: HTMLElement // 触发元素
  zIndex?: number // 层级，默认 9999
  className?: string // 自定义样式类
  placement?: 'auto' | 'top' | 'bottom' | 'left' | 'right'
}

export interface CyberContextMenuEmits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'item-click', item: ContextMenuItem, event: MouseEvent): void
  (e: 'close'): void
}

export interface MenuPosition {
  x: number
  y: number
  placement: 'top' | 'bottom' | 'left' | 'right'
}
