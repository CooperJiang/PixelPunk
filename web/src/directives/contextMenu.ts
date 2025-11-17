import { h, render, type App, type DirectiveBinding, type VNode } from 'vue'
import CyberContextMenu from '@/components/CyberContextMenu/index.vue'
import type { ContextMenuOptions, ContextMenuItem } from '@/components/CyberContextMenu/types'

interface ContextMenuElement extends HTMLElement {
  __contextMenuOptions?: ContextMenuOptions
  __contextMenuHandler?: (event: MouseEvent) => void
}

/* 全局菜单容器 */
let globalMenuContainer: HTMLElement | null = null
let globalMenuVNode: VNode | null = null

const getGlobalMenuContainer = (): HTMLElement => {
  if (!globalMenuContainer) {
    globalMenuContainer = document.createElement('div')
    globalMenuContainer.id = 'global-context-menu-container'
    globalMenuContainer.style.position = 'absolute'
    globalMenuContainer.style.top = '0'
    globalMenuContainer.style.left = '0'
    globalMenuContainer.style.pointerEvents = 'auto'
    globalMenuContainer.style.zIndex = '9999'
    document.body.appendChild(globalMenuContainer)
  }
  return globalMenuContainer
}

const showContextMenu = (options: {
  items: ContextMenuItem[]
  x: number
  y: number
  zIndex?: number
  className?: string
  placement?: 'auto' | 'top' | 'bottom' | 'left' | 'right'
  trigger?: HTMLElement
}) => {
  const container = getGlobalMenuContainer()

  const menuVNode = h(CyberContextMenu, {
    modelValue: true,
    items: options.items,
    x: options.x,
    y: options.y,
    zIndex: options.zIndex,
    className: options.className,
    placement: options.placement,
    trigger: options.trigger ?? container,
    'onUpdate:modelValue': (value: boolean) => {
      if (!value) {
        hideContextMenu()
      }
    },
    onClose: () => {
      hideContextMenu()
    },
    onItemClick: (_item: ContextMenuItem, _event: MouseEvent) => {
      hideContextMenu()
    },
  })

  render(menuVNode, container)
  globalMenuVNode = menuVNode
}

const hideContextMenu = () => {
  if (globalMenuContainer && globalMenuVNode) {
    render(null, globalMenuContainer)
    globalMenuVNode = null
  }
}

const parseDirectiveValue = (value: ContextMenuItem[] | ContextMenuOptions): ContextMenuOptions => {
  if (Array.isArray(value)) {
    return {
      items: value,
      trigger: 'contextmenu',
    }
  }
  return {
    trigger: 'contextmenu',
    ...value,
  }
}

const createEventHandler = (options: ContextMenuOptions) => {
  return (event: MouseEvent) => {
    event.preventDefault()
    event.stopPropagation()

    hideContextMenu()

    let x = event.clientX
    let y = event.clientY

    if (options.offset) {
      x += options.offset.x
      y += options.offset.y
    }

    setTimeout(() => {
      showContextMenu({
        items: options.items,
        x,
        y,
        zIndex: options.zIndex,
        className: options.className,
        placement: options.placement,
        trigger: (event.currentTarget as HTMLElement) || undefined,
      })
    }, 0)
  }
}

let globalClickHandler: ((event: MouseEvent) => void) | null = null

const setupGlobalClickHandler = () => {
  if (!globalClickHandler) {
    globalClickHandler = (event: MouseEvent) => {
      const target = event.target as Element
      const inMenu =
        target &&
        (target.closest('.cyber-context-menu') ||
          target.closest('.cyber-context-submenu') ||
          target.closest('.cyber-context-menu-overlay'))
      if (!inMenu) hideContextMenu()
    }

    document.addEventListener('click', globalClickHandler, true)
    document.addEventListener('contextmenu', globalClickHandler, true)
  }
}

export const vContextMenu = {
  mounted(el: ContextMenuElement, binding: DirectiveBinding<ContextMenuItem[] | ContextMenuOptions>) {
    if (!binding.value) return

    setupGlobalClickHandler()

    const options = parseDirectiveValue(binding.value)
    el.__contextMenuOptions = options

    const handler = createEventHandler(options)
    el.__contextMenuHandler = handler

    const eventType = options.trigger === 'click' ? 'click' : 'contextmenu'
    el.addEventListener(eventType, handler, { passive: false })

    el.style.cursor = eventType === 'contextmenu' ? 'context-menu' : 'pointer'
  },

  updated(el: ContextMenuElement, binding: DirectiveBinding<ContextMenuItem[] | ContextMenuOptions>) {
    if (!binding.value) {
      if (el.__contextMenuHandler && el.__contextMenuOptions) {
        const eventType = el.__contextMenuOptions.trigger === 'click' ? 'click' : 'contextmenu'
        el.removeEventListener(eventType, el.__contextMenuHandler)
        el.__contextMenuHandler = undefined
        el.__contextMenuOptions = undefined
        el.style.cursor = ''
      }
      return
    }

    const options = parseDirectiveValue(binding.value)
    const oldOptions = el.__contextMenuOptions

    if (el.__contextMenuHandler && oldOptions) {
      const oldEventType = oldOptions.trigger === 'click' ? 'click' : 'contextmenu'
      el.removeEventListener(oldEventType, el.__contextMenuHandler)
    }

    const handler = createEventHandler(options)
    el.__contextMenuHandler = handler
    const eventType = options.trigger === 'click' ? 'click' : 'contextmenu'
    el.addEventListener(eventType, handler, { passive: false })

    el.style.cursor = eventType === 'contextmenu' ? 'context-menu' : 'pointer'
    el.__contextMenuOptions = options
  },

  unmounted(el: ContextMenuElement) {
    if (el.__contextMenuHandler && el.__contextMenuOptions) {
      const eventType = el.__contextMenuOptions.trigger === 'click' ? 'click' : 'contextmenu'
      el.removeEventListener(eventType, el.__contextMenuHandler)
      el.__contextMenuHandler = undefined
      el.__contextMenuOptions = undefined
      el.style.cursor = ''
    }
  },
}

export default {
  install(app: App) {
    app.directive('context-menu', vContextMenu)
  },
}

export { showContextMenu, hideContextMenu }
