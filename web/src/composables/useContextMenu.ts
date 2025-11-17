import { ref, readonly } from 'vue'
import {
  showContextMenu as directiveShowContextMenu,
  hideContextMenu as directiveHideContextMenu,
} from '@/directives/contextMenu'
import type { ContextMenuItem } from '@/components/CyberContextMenu/types'

interface ShowContextMenuOptions {
  items: ContextMenuItem[]
  x: number
  y: number
  zIndex?: number
  className?: string
  placement?: 'auto' | 'top' | 'bottom' | 'left' | 'right'
}

/* 全局菜单状态 */
const isVisible = ref(false)

/**
 * 右键菜单组合式函数
 */
export function useContextMenu() {
  const showContextMenu = (options: ShowContextMenuOptions) => {
    directiveShowContextMenu({
      items: options.items,
      x: options.x,
      y: options.y,
      zIndex: options.zIndex,
      className: options.className,
      placement: options.placement,
    })

    isVisible.value = true
  }

  const hideContextMenu = () => {
    directiveHideContextMenu()
    isVisible.value = false
  }

  const showContextMenuAtElement = (
    element: HTMLElement,
    items: ContextMenuItem[],
    options?: Partial<Omit<ShowContextMenuOptions, 'items' | 'x' | 'y'>>
  ) => {
    const rect = element.getBoundingClientRect()
    showContextMenu({
      items,
      x: rect.left,
      y: rect.bottom,
      ...options,
    })
  }

  const showContextMenuAtMouse = (
    event: MouseEvent,
    items: ContextMenuItem[],
    options?: Partial<Omit<ShowContextMenuOptions, 'items' | 'x' | 'y'>>
  ) => {
    event.preventDefault()
    event.stopPropagation()

    showContextMenu({
      items,
      x: event.clientX,
      y: event.clientY,
      ...options,
    })
  }

  const createContextMenuHandler = (
    items: ContextMenuItem[] | ((event: MouseEvent) => ContextMenuItem[]),
    options?: Partial<Omit<ShowContextMenuOptions, 'items' | 'x' | 'y'>>
  ) => {
    return (event: MouseEvent) => {
      const menuItems = typeof items === 'function' ? items(event) : items
      showContextMenuAtMouse(event, menuItems, options)
    }
  }

  const cleanup = () => {
    hideContextMenu()
  }

  return {
    isVisible: readonly(isVisible),

    showContextMenu,
    hideContextMenu,
    showContextMenuAtElement,
    showContextMenuAtMouse,
    createContextMenuHandler,
    cleanup,
  }
}

let globalContextMenu: ReturnType<typeof useContextMenu> | null = null

export function useGlobalContextMenu() {
  if (!globalContextMenu) {
    globalContextMenu = useContextMenu()
  }
  return globalContextMenu
}

export default useContextMenu
