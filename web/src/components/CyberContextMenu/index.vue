<template>
  <teleport v-if="useDirectTeleport" to="body">
    <transition name="cyber-context-menu" appear>
      <div
        v-if="modelValue"
        ref="menuRef"
        class="cyber-context-menu"
        :class="[className, `cyber-context-menu--${finalPlacement}`]"
        :style="menuStyle"
        tabindex="-1"
        @click.stop
        @keydown="handleKeydown"
      >
        <div class="cyber-context-menu__content">
          <CyberContextMenuItem
            v-for="(item, index) in items"
            :key="item.key"
            :item="item"
            :is-first="index === 0"
            @item-click="handleItemClick"
          />
        </div>
      </div>
    </transition>

    <div
      v-if="modelValue"
      class="cyber-context-menu-overlay"
      @click="handleOverlayClick"
      @contextmenu.prevent="handleOverlayClick"
    />
  </teleport>

  <transition v-else name="cyber-context-menu" appear>
    <div
      v-if="modelValue"
      ref="menuRef"
      class="cyber-context-menu"
      :class="[className, `cyber-context-menu--${finalPlacement}`]"
      :style="menuStyle"
      tabindex="-1"
      @click.stop
      @keydown="handleKeydown"
    >
      <div class="cyber-context-menu__content">
        <CyberContextMenuItem
          v-for="(item, index) in items"
          :key="item.key"
          :item="item"
          :is-first="index === 0"
          @item-click="handleItemClick"
        />
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import CyberContextMenuItem from './CyberContextMenuItem.vue'
  import type { CyberContextMenuProps, CyberContextMenuEmits, ContextMenuItem, MenuPosition } from './types'

  defineOptions({
    name: 'CyberContextMenu',
  })

  const props = withDefaults(defineProps<CyberContextMenuProps>(), {
    modelValue: false,
    x: 0,
    y: 0,
    zIndex: 9999,
    placement: 'auto',
  })

  /* 检测是否作为独立组件使用（有父组件）或通过指令渲染（无父组件，直接渲染到 body） */
  const useDirectTeleport = computed(() => {
    return !props.trigger
  })

  const emit = defineEmits<CyberContextMenuEmits>()

  const menuRef = ref<HTMLElement>()
  const menuPosition = ref<MenuPosition>({ x: props.x, y: props.y, placement: 'bottom' })
  const currentFocusIndex = ref(-1)

  const finalPlacement = computed(() => {
    return props.placement === 'auto' ? menuPosition.value.placement : props.placement
  })

  const menuStyle = computed(() => {
    return {
      left: `${menuPosition.value.x}px`,
      top: `${menuPosition.value.y}px`,
      zIndex: props.zIndex,
    }
  })

  const flattenedItems = computed(() => {
    const items: ContextMenuItem[] = []
    const flatten = (menuItems: ContextMenuItem[]) => {
      menuItems.forEach((item) => {
        if (!item.divided) {
          items.push(item)
        }
        if (item.children) {
          flatten(item.children)
        }
      })
    }
    flatten(props.items)
    return items.filter((item) => !item.disabled)
  })

  const calculatePosition = async () => {
    await nextTick()

    if (!menuRef.value) return

    const menuRect = menuRef.value.getBoundingClientRect()
    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight

    let x = props.x
    let y = props.y
    let placement: 'top' | 'bottom' | 'left' | 'right' = 'bottom'

    if (props.placement === 'auto') {
      const spaceBottom = viewportHeight - props.y

      if (spaceBottom >= menuRect.height + 8) {
        placement = 'bottom'
        y = props.y
      } else {
        placement = 'top'
        y = props.y - menuRect.height
      }

      x = props.x
      if (x + menuRect.width > viewportWidth - 10) x = viewportWidth - menuRect.width - 10
      if (x < 10) x = 10

      if (y < 10) y = 10
      if (y + menuRect.height > viewportHeight - 10) y = viewportHeight - menuRect.height - 10
    } else {
      switch (props.placement) {
        case 'top':
          y = props.y - menuRect.height
          placement = 'top'
          break
        case 'left':
          x = props.x - menuRect.width
          placement = 'left'
          break
        case 'right':
          x = props.x
          placement = 'right'
          break
        default:
          placement = 'bottom'
      }
    }

    x = Math.max(10, Math.min(x, viewportWidth - menuRect.width - 10))
    y = Math.max(10, Math.min(y, viewportHeight - menuRect.height - 10))

    menuPosition.value = { x, y, placement }
  }

  const handleItemClick = (item: ContextMenuItem, event: MouseEvent) => {
    emit('item-click', item, event)
    closeMenu()
  }

  const handleOverlayClick = () => {
    closeMenu()
  }

  const closeMenu = () => {
    emit('update:modelValue', false)
    emit('close')
    currentFocusIndex.value = -1
  }

  const handleKeydown = (event: KeyboardEvent) => {
    switch (event.key) {
      case 'Escape':
        event.preventDefault()
        closeMenu()
        break
      case 'ArrowDown':
        event.preventDefault()
        navigateMenu(1)
        break
      case 'ArrowUp':
        event.preventDefault()
        navigateMenu(-1)
        break
      case 'Enter':
      case ' ':
        event.preventDefault()
        if (currentFocusIndex.value >= 0) {
          const item = flattenedItems.value[currentFocusIndex.value]
          if (item && !item.disabled) {
            handleItemClick(item, event as unknown as MouseEvent)
          }
        }
        break
    }
  }

  const navigateMenu = (direction: 1 | -1) => {
    const items = flattenedItems.value
    if (items.length === 0) return

    currentFocusIndex.value = Math.max(0, Math.min(items.length - 1, currentFocusIndex.value + direction))
  }

  watch(
    () => props.modelValue,
    (newValue) => {
      if (newValue) {
        menuPosition.value.x = props.x
        menuPosition.value.y = props.y

        calculatePosition()
        nextTick(() => {
          menuRef.value?.focus()
        })
      } else {
        currentFocusIndex.value = -1
      }
    }
  )

  watch([() => props.x, () => props.y], ([newX, newY]) => {
    menuPosition.value.x = newX
    menuPosition.value.y = newY

    if (props.modelValue) {
      calculatePosition()
    }
  })

  onMounted(() => {
    if (props.modelValue) {
      calculatePosition()
    }
  })

  const handleGlobalKeydown = (event: KeyboardEvent) => {
    if (!props.modelValue) return

    if (event.key === 'Escape') {
      closeMenu()
      return
    }

    const shortcutKey = getShortcutKey(event)
    if (shortcutKey) {
      const matchingItem = findItemByShortcut(props.items, shortcutKey)
      if (matchingItem && !matchingItem.disabled) {
        event.preventDefault()
        event.stopPropagation()
        handleItemClick(matchingItem, event as unknown as MouseEvent)
      }
    }
  }

  const getShortcutKey = (event: KeyboardEvent): string => {
    const modifiers = []
    if (event.ctrlKey || event.metaKey) modifiers.push('Ctrl')
    if (event.altKey) modifiers.push('Alt')
    if (event.shiftKey) modifiers.push('Shift')

    let key = event.key
    if (key === ' ') key = 'Space'

    return modifiers.length > 0 ? `${modifiers.join('+')}+${key}` : key
  }

  const findItemByShortcut = (items: ContextMenuItem[], shortcut: string): ContextMenuItem | null => {
    for (const item of items) {
      if (item.shortcut === shortcut) {
        return item
      }
      if (item.children) {
        const found = findItemByShortcut(item.children, shortcut)
        if (found) return found
      }
    }
    return null
  }

  onMounted(() => {
    document.addEventListener('keydown', handleGlobalKeydown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleGlobalKeydown)
  })
</script>

<style scoped>
  .cyber-context-menu {
    position: fixed;
    min-width: 160px;
    max-width: 240px;
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(12px);
    padding: 6px 0;
    outline: none;
    user-select: none;
  }

  .cyber-context-menu__content {
    width: 100%;
  }

  .cyber-context-menu-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 9998;
    background: transparent;
  }

  .cyber-context-menu-enter-active {
    transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .cyber-context-menu-leave-active {
    transition: all 0.15s ease-in;
  }

  .cyber-context-menu-enter-from {
    opacity: 0;
    transform: scale(0.95) translateY(-10px);
  }

  .cyber-context-menu-leave-to {
    opacity: 0;
    transform: scale(0.95) translateY(-5px);
  }

  .cyber-context-menu--top.cyber-context-menu-enter-from {
    transform: scale(0.95) translateY(10px);
  }

  .cyber-context-menu--left.cyber-context-menu-enter-from {
    transform: scale(0.95) translateX(10px);
  }

  .cyber-context-menu--right.cyber-context-menu-enter-from {
    transform: scale(0.95) translateX(-10px);
  }
</style>
