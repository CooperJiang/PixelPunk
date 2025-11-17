<template>
  <div
    class="cyber-context-menu-item"
    :class="{
      'cyber-context-menu-item--disabled': item.disabled,
      'cyber-context-menu-item--danger': item.danger && !item.disabled,
      'cyber-context-menu-item--has-children': hasChildren,
      'cyber-context-menu-item--active': isActive,
      'cyber-context-menu-item--has-active-submenu': hasChildren && showSubMenu,
      'cyber-context-menu-item--hover': isHovered && !item.disabled,
    }"
    @click="handleClick"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div v-if="item.divided && !isFirst" class="cyber-context-menu-item__divider" />

    <div v-if="!item.divided || item.label" class="cyber-context-menu-item__content">
      <span v-if="item.icon" class="cyber-context-menu-item__icon">
        <i :class="item.icon" />
      </span>

      <span class="cyber-context-menu-item__label">
        {{ item.label }}
      </span>

      <span v-if="item.shortcut" class="cyber-context-menu-item__shortcut">
        {{ item.shortcut }}
      </span>

      <span v-if="hasChildren" class="cyber-context-menu-item__arrow">
        <i class="fas fa-chevron-right" />
      </span>
    </div>

    <transition name="cyber-context-submenu" appear>
      <div
        v-if="hasChildren && showSubMenu"
        ref="submenuRef"
        class="cyber-context-submenu"
        :class="{ 'cyber-context-submenu--left': submenuOpensLeft }"
        :style="submenuStyle"
        @click.stop
        @mouseenter="handleSubmenuMouseEnter"
        @mouseleave="handleSubmenuMouseLeave"
      >
        <div v-if="isLoading" class="cyber-context-menu-item__loading">
          <i class="fas fa-spinner fa-spin cyber-context-menu-item__loading-icon" />
          <span class="cyber-context-menu-item__loading-text">{{ $t('components.cyberContextMenu.loading') }}</span>
        </div>

        <div v-else-if="loadError" class="cyber-context-menu-item__error">
          <i class="fas fa-exclamation-triangle cyber-context-menu-item__error-icon" />
          <span class="cyber-context-menu-item__error-text">{{ loadError }}</span>
        </div>

        <CyberContextMenuItem
          v-else
          v-for="(childItem, index) in localChildren"
          :key="childItem.key"
          :item="childItem"
          :is-first="index === 0"
          @item-click="(item, event) => $emit('item-click', item, event)"
        />
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref, nextTick, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { ContextMenuItem } from './types'

  const { $t } = useTexts()

  defineOptions({
    name: 'CyberContextMenuItem',
  })

  interface Props {
    item: ContextMenuItem
    isFirst?: boolean
  }

  interface Emits {
    (e: 'item-click', item: ContextMenuItem, event: MouseEvent): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const showSubMenu = ref(false)
  const isActive = ref(false)
  const isHovered = ref(false)
  const submenuRef = ref<HTMLElement>()
  const submenuStyle = ref<Record<string, string>>({})
  const localChildren = ref<ContextMenuItem[]>(props.item.children ?? [])
  const isLoading = ref(false)
  const loadError = ref<string | undefined>(undefined)
  const submenuOpensLeft = ref(false)
  let openTimer: number | null = null

  const hasChildren = computed(() => {
    return (localChildren.value && localChildren.value.length > 0) || props.item.hasAsyncChildren
  })

  /* 与外部 children 同步（若父侧更新） */
  watch(
    () => props.item.children,
    (val) => {
      if (Array.isArray(val)) localChildren.value = val
    }
  )

  /* 子菜单内容变化（首次加载完成、错误切换等）时，展开状态下重新计算位置 */
  watch([localChildren, isLoading, loadError], async () => {
    if (showSubMenu.value) {
      await nextTick()
      calculateSubmenuPosition()
    }
  })

  const handleClick = async (event: MouseEvent) => {
    event.stopPropagation()

    if (props.item.divided && !props.item.label) {
      return
    }

    if (props.item.disabled) {
      return
    }

    const target = event.target as Element
    const clickedArrow = !!target.closest('.cyber-context-menu-item__arrow')

    if (clickedArrow && hasChildren.value) {
      showSubMenu.value = !showSubMenu.value
      if (showSubMenu.value) {
        nextTick(() => {
          calculateSubmenuPosition()
        })
      }
      return
    }

    if (props.item.onClick) {
      await Promise.resolve(props.item.onClick(props.item, event))
      emit('item-click', props.item, event)
      return
    }

    if (hasChildren.value) {
      showSubMenu.value = !showSubMenu.value
      if (showSubMenu.value) {
        nextTick(() => {
          calculateSubmenuPosition()
        })
      }
    } else {
      emit('item-click', props.item, event)
    }
  }

  const handleMouseEnter = async (event: MouseEvent) => {
    if (props.item.disabled) return

    const target = event.target as Element
    const currentTarget = event.currentTarget as Element

    if (
      target.classList.contains('cyber-context-menu-item__content') ||
      target.closest('.cyber-context-menu-item') === currentTarget
    ) {
      if (!target.closest('.cyber-context-submenu')) {
        isHovered.value = true
      }
    }

    isActive.value = true
    if (hasChildren.value) {
      if (openTimer) {
        clearTimeout(openTimer)
        openTimer = null
      }
      openTimer = window.setTimeout(async () => {
        if (props.item.hasAsyncChildren && localChildren.value.length === 0 && props.item.loadChildren) {
          isLoading.value = true
          showSubMenu.value = true
          await nextTick()
          calculateSubmenuPosition()
          try {
            const children = await props.item.loadChildren()
            localChildren.value = Array.isArray(children) ? children : []
            loadError.value = undefined
          } catch (error) {
            console.error('Failed to load submenu children:', error)
            loadError.value = $t('components.cyberContextMenu.loadFailed')
          } finally {
            isLoading.value = false
            await nextTick()
            calculateSubmenuPosition()
          }
        } else {
          showSubMenu.value = true
          await nextTick()
          calculateSubmenuPosition()
        }
      }, 140)
    }
  }

  const handleMouseLeave = (event: MouseEvent) => {
    const rect = (event.currentTarget as Element).getBoundingClientRect()
    const isLeavingItem =
      event.clientX < rect.left || event.clientX > rect.right || event.clientY < rect.top || event.clientY > rect.bottom

    if (isLeavingItem) {
      isHovered.value = false
    }

    isActive.value = false
    if (openTimer) {
      clearTimeout(openTimer)
      openTimer = null
    }
    if (hasChildren.value) {
      const delay = 250
      setTimeout(() => {
        const hoveringCurrent = submenuRef.value ? submenuRef.value.matches(':hover') : false
        const hoveringDescendant = submenuRef.value ? !!submenuRef.value.querySelector('.cyber-context-submenu:hover') : false
        if (!isActive.value && !isSubmenuHovered.value && !hoveringCurrent && !hoveringDescendant) {
          showSubMenu.value = false
        }
      }, delay)
    }
  }

  const isSubmenuHovered = ref(false)

  const handleSubmenuMouseEnter = () => {
    isSubmenuHovered.value = true
    isActive.value = true
    isHovered.value = false
  }

  const handleSubmenuMouseLeave = () => {
    isSubmenuHovered.value = false
    isActive.value = false
    setTimeout(() => {
      if (!isActive.value && !isSubmenuHovered.value) {
        showSubMenu.value = false
      }
    }, 100)
  }

  const calculateSubmenuPosition = () => {
    if (!submenuRef.value) return

    const rect = submenuRef.value.getBoundingClientRect()
    const parentRect = (submenuRef.value.parentElement as HTMLElement).getBoundingClientRect()
    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight

    let left: string | undefined = '100%'
    let right: string | undefined = 'auto'
    let topOffset = 0
    const margin = 10

    if (parentRect.right + rect.width > viewportWidth - 10) {
      submenuOpensLeft.value = true
      left = 'auto'
      right = '100%'
    } else {
      submenuOpensLeft.value = false
      left = '100%'
      right = 'auto'
    }

    const minTopOffset = margin - parentRect.top // 不超过视口上边缘
    const maxTopOffset = viewportHeight - rect.height - margin - parentRect.top // 不超过视口下边缘
    topOffset = Math.min(Math.max(0, minTopOffset), maxTopOffset)

    submenuStyle.value = {
      left: left as string,
      right: right as string,
      top: `${Math.round(topOffset)}px`,
    }
  }
</script>

<style scoped>
  .cyber-context-menu-item {
    position: relative;
    user-select: none;
  }

  .cyber-context-menu-item__divider {
    height: 1px;
    background: var(--color-border-default);
    margin: 6px 12px;
    opacity: 0.6;
  }

  .cyber-context-menu-item__content {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    min-height: 28px;
    cursor: pointer;
    transition: all 0.15s ease;
    border-radius: var(--radius-sm);
    margin: 1px 6px;
    position: relative;
  }

  .cyber-context-menu-item--hover > .cyber-context-menu-item__content {
    background: var(--color-background-600);
    border-left: 2px solid var(--color-brand-400);
    padding-left: 10px;
    transform: translateX(2px);
  }

  .cyber-context-menu-item:not(.cyber-context-menu-item--disabled):hover > .cyber-context-menu-item__content {
    background: var(--color-background-600);
    border-left: 2px solid var(--color-brand-400);
    padding-left: 10px;
    transform: translateX(2px);
  }

  .cyber-context-menu-item--has-active-submenu > .cyber-context-menu-item__content {
    background: var(--color-background-600);
    border-left: 2px solid var(--color-brand-500);
    padding-left: 10px;
  }

  .cyber-context-menu-item--danger.cyber-context-menu-item--hover > .cyber-context-menu-item__content {
    background: rgba(var(--color-error-rgb), 0.15) !important;
    color: var(--color-error-500);
    border-left-color: var(--color-error-500) !important;
  }

  .cyber-context-menu-item--danger:not(.cyber-context-menu-item--disabled):hover > .cyber-context-menu-item__content {
    background: rgba(var(--color-error-rgb), 0.15) !important;
    color: var(--color-error-500);
    border-left-color: var(--color-error-500) !important;
  }

  .cyber-context-menu-item--disabled .cyber-context-menu-item__content {
    opacity: 0.5;
    cursor: not-allowed;
    color: var(--color-content-disabled);
  }

  .cyber-context-menu-item__icon {
    min-width: 14px;
    text-align: center;
    color: var(--color-content-muted);
    font-size: 0.875rem;
  }

  .cyber-context-menu-item--danger .cyber-context-menu-item__icon {
    color: var(--color-error-500);
  }

  .cyber-context-menu-item__label {
    flex: 1;
    font-size: 0.8125rem;
    color: var(--color-content-default);
    white-space: nowrap;
    font-weight: 500;
  }

  .cyber-context-menu-item--danger .cyber-context-menu-item__label {
    color: var(--color-error-500);
  }

  .cyber-context-menu-item__shortcut {
    font-size: 0.6875rem;
    color: var(--color-content-muted);
    opacity: 0.75;
    font-weight: 400;
  }

  .cyber-context-menu-item__arrow {
    min-width: 12px;
    text-align: center;
    color: var(--color-content-muted);
    font-size: 0.75rem;
  }

  .cyber-context-submenu {
    position: absolute;
    left: 100%;
    top: 0;
    min-width: 140px;
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    box-shadow: 0 6px 24px rgba(0, 0, 0, 0.35);
    backdrop-filter: blur(12px);
    padding: 6px 0;
    z-index: 10001;
  }

  .cyber-context-submenu.cyber-context-submenu--left {
    left: auto;
    right: 100%;
  }

  .cyber-context-submenu-enter-active,
  .cyber-context-submenu-leave-active {
    transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .cyber-context-submenu-enter-from {
    opacity: 0;
    transform: scale(0.95) translateX(-10px);
  }

  .cyber-context-submenu--left.cyber-context-submenu-enter-from {
    transform: scale(0.95) translateX(10px);
  }

  .cyber-context-submenu-leave-to {
    opacity: 0;
    transform: scale(0.95) translateX(-10px);
  }
  .cyber-context-submenu--left.cyber-context-submenu-leave-to {
    transform: scale(0.95) translateX(10px);
  }

  .cyber-context-menu-item__loading,
  .cyber-context-menu-item__error {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    min-height: 32px;
    color: var(--color-content-muted);
    font-size: 0.8125rem;
  }

  .cyber-context-menu-item__loading-icon,
  .cyber-context-menu-item__error-icon {
    font-size: 0.75rem;
  }

  .cyber-context-menu-item__error {
    color: var(--color-error-500);
  }
</style>
