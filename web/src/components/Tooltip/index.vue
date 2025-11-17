<script setup lang="ts">
  import { nextTick, onBeforeUnmount, ref } from 'vue'
  import type { TooltipProps } from './types'

  defineOptions({
    name: 'CyberTooltip',
  })

  const props = withDefaults(defineProps<TooltipProps>(), {
    placement: 'top',
    trigger: 'hover',
    disabled: false,
    offset: 8,
    showDelay: 0,
    hideDelay: 200,
    theme: 'dark',
    maxWidth: '300px',
  })

  const visible = ref(false)
  const trigger = ref<HTMLElement>()
  const tooltip = ref<HTMLElement>()
  const tooltipStyle = ref<Record<string, string | number>>({})
  const arrowStyle = ref<Record<string, string | number>>({})

  let showTimer: number | null = null
  let hideTimer: number | null = null

  const clearTimers = () => {
    if (showTimer) {
      clearTimeout(showTimer)
      showTimer = null
    }
    if (hideTimer) {
      clearTimeout(hideTimer)
      hideTimer = null
    }
  }

  const show = () => {
    if (props.disabled) {
      return
    }

    clearTimers()
    showTimer = window.setTimeout(() => {
      visible.value = true
      nextTick(() => {
        updatePosition()
      })
    }, props.showDelay)
  }

  const hide = () => {
    clearTimers()
    hideTimer = window.setTimeout(() => {
      visible.value = false
    }, props.hideDelay)
  }

  const handleMouseEnter = () => {
    if (props.trigger === 'hover') {
      show()
    }
  }

  const handleMouseLeave = () => {
    if (props.trigger === 'hover') {
      hide()
    }
  }

  const handleClick = () => {
    if (props.trigger === 'click') {
      visible.value = !visible.value
      if (visible.value) {
        nextTick(() => {
          updatePosition()
        })
      }
    }
  }

  const handleTooltipMouseEnter = () => {
    if (props.trigger === 'hover') {
      clearTimers()
    }
  }

  const handleTooltipMouseLeave = () => {
    if (props.trigger === 'hover') {
      hide()
    }
  }

  const updatePosition = () => {
    if (!trigger.value || !tooltip.value) {
      return
    }

    const triggerRect = trigger.value.getBoundingClientRect()
    const tooltipEl = tooltip.value

    const tempStyle = {
      position: 'fixed',
      visibility: 'hidden',
      display: 'block',
      top: '0px',
      left: '0px',
      maxWidth: typeof props.maxWidth === 'number' ? `${props.maxWidth}px` : props.maxWidth,
    }

    Object.assign(tooltipEl.style, tempStyle)

    const { offsetWidth: tooltipWidth, offsetHeight: tooltipHeight } = tooltipEl

    let left = 0
    let top = 0

    const { placement } = props

    let offsetX = 0
    let offsetY = 0
    if (Array.isArray(props.offset)) {
      offsetX = props.offset[0] || 0
      offsetY = props.offset[1] || 0
    } else {
      offsetY = props.offset || 8
    }

    const _expectedLeft = 0
    switch (placement) {
      case 'top':
      case 'top-start':
      case 'top-end':
        top = triggerRect.top - tooltipHeight - offsetY
        if (placement === 'top-start') {
          left = triggerRect.left + offsetX
          const _expectedLeft = left
        } else if (placement === 'top-end') {
          left = triggerRect.right - tooltipWidth + offsetX
          const _expectedLeft = left
        } else {
          left = triggerRect.left + triggerRect.width / 2 - tooltipWidth / 2 + offsetX
          const _expectedLeft = left
        }
        break

      case 'bottom':
      case 'bottom-start':
      case 'bottom-end':
        top = triggerRect.bottom + offsetY
        if (placement === 'bottom-start') {
          left = triggerRect.left + offsetX
          const _expectedLeft = left
        } else if (placement === 'bottom-end') {
          left = triggerRect.right - tooltipWidth + offsetX
          const _expectedLeft = left
        } else {
          left = triggerRect.left + triggerRect.width / 2 - tooltipWidth / 2 + offsetX
          const _expectedLeft = left
        }
        break

      case 'left':
        left = triggerRect.left - tooltipWidth - offsetY
        top = triggerRect.top + triggerRect.height / 2 - tooltipHeight / 2 + offsetX
        break

      case 'right':
        left = triggerRect.right + offsetY
        top = triggerRect.top + triggerRect.height / 2 - tooltipHeight / 2 + offsetX
        break
    }

    const viewportWidth = window.innerWidth
    const viewportHeight = window.innerHeight

    if (left < 10) {
      left = 10
    } else if (left + tooltipWidth > viewportWidth - 10) {
      left = viewportWidth - tooltipWidth - 10
    }

    if (top < 10) {
      top = 10
    } else if (top + tooltipHeight > viewportHeight - 10) {
      top = viewportHeight - tooltipHeight - 10
    }

    tooltipStyle.value = {
      position: 'fixed',
      left: `${left}px`,
      top: `${top}px`,
      maxWidth: typeof props.maxWidth === 'number' ? `${props.maxWidth}px` : props.maxWidth,
      zIndex: 99999,
      visibility: 'visible',
      display: 'block',
    }

    arrowStyle.value = {}
    if (placement === 'top' || placement === 'bottom') {
      const triggerCenter = triggerRect.left + triggerRect.width / 2
      const tooltipStart = left
      const arrowLeft = triggerCenter - tooltipStart - 4 // 4px是箭头宽度的一半

      const minLeft = 8
      const maxLeft = tooltipWidth - 16
      const clampedLeft = Math.max(minLeft, Math.min(maxLeft, arrowLeft))

      arrowStyle.value = {
        left: `${clampedLeft}px`,
      }
    } else if (placement === 'left' || placement === 'right') {
      const triggerCenter = triggerRect.top + triggerRect.height / 2
      const tooltipStart = top
      const arrowTop = triggerCenter - tooltipStart - 4 // 4px是箭头高度的一半

      const minTop = 8
      const maxTop = tooltipHeight - 16
      const clampedTop = Math.max(minTop, Math.min(maxTop, arrowTop))

      arrowStyle.value = {
        top: `${clampedTop}px`,
      }
    }
  }

  const handleResize = () => {
    if (visible.value) {
      updatePosition()
    }
  }

  const handleScroll = () => {
    if (visible.value) {
      updatePosition()
    }
  }

  const handleClickOutside = (event: Event) => {
    if (props.trigger === 'click' && visible.value) {
      const target = event.target as Node
      if (!trigger.value?.contains(target) && !tooltip.value?.contains(target)) {
        visible.value = false
      }
    }
  }

  onBeforeUnmount(() => {
    clearTimers()
    window.removeEventListener('resize', handleResize)
    window.removeEventListener('scroll', handleScroll, true)
    document.removeEventListener('click', handleClickOutside, true)
  })

  window.addEventListener('resize', handleResize)
  window.addEventListener('scroll', handleScroll, true)
  document.addEventListener('click', handleClickOutside, true)

  defineExpose({
    show,
    hide,
    updatePosition,
  })
</script>

<template>
  <div ref="trigger" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave" @click="handleClick">
    <slot />
  </div>

  <Teleport to="body">
    <Transition name="cyber-tooltip-fade">
      <div
        v-if="visible"
        ref="tooltip"
        class="cyber-tooltip"
        :class="[`cyber-tooltip--${placement}`, { 'cyber-tooltip--dark': theme === 'dark' }]"
        :style="tooltipStyle"
        @mouseenter="handleTooltipMouseEnter"
        @mouseleave="handleTooltipMouseLeave"
      >
        <div class="cyber-tooltip__arrow" :class="`cyber-tooltip__arrow--${placement}`" :style="arrowStyle" />
        <div class="cyber-tooltip__content">
          <slot name="content">{{ content }}</slot>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .cyber-tooltip {
    @apply pointer-events-none fixed max-w-xs select-none whitespace-normal rounded-md text-xs leading-relaxed;
    z-index: 99999;
    padding: 8px 12px;
    word-wrap: break-word;
    word-break: break-word;
    transform-origin: center center;
  }

  .cyber-tooltip--dark {
    background: var(--color-background-800);
    color: var(--color-content-default);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow:
      0 4px 12px rgba(0, 0, 0, 0.4),
      0 0 8px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 1px 0 rgba(var(--color-content-default-rgb), 0.05);
    backdrop-filter: blur(8px);
  }

  .cyber-tooltip:not(.cyber-tooltip--dark) {
    background: rgba(var(--color-content-default-rgb), 0.95);
    color: var(--color-background-900);
    border: 1px solid rgba(var(--color-background-900-rgb), 0.15);
    box-shadow:
      0 4px 12px rgba(0, 0, 0, 0.15),
      0 2px 4px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(8px);
  }

  .cyber-tooltip__arrow {
    @apply absolute h-2 w-2;
    background: var(--color-background-800);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    transform: rotate(45deg);
  }

  .cyber-tooltip--top .cyber-tooltip__arrow--top {
    bottom: -5px;
    border-bottom-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-right-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-top-color: transparent !important;
    border-left-color: transparent !important;
  }

  .cyber-tooltip--top-start .cyber-tooltip__arrow--top-start,
  .cyber-tooltip--top-end .cyber-tooltip__arrow--top-end {
    bottom: -5px;
    border-bottom-color: transparent !important;
    border-right-color: transparent !important;
  }

  .cyber-tooltip--top-start .cyber-tooltip__arrow--top-start {
    left: 16px;
  }

  .cyber-tooltip--top-end .cyber-tooltip__arrow--top-end {
    right: 16px;
  }

  .cyber-tooltip--bottom .cyber-tooltip__arrow--bottom {
    top: -5px;
    border-top-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-left-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-bottom-color: transparent !important;
    border-right-color: transparent !important;
  }

  .cyber-tooltip--bottom-start .cyber-tooltip__arrow--bottom-start,
  .cyber-tooltip--bottom-end .cyber-tooltip__arrow--bottom-end {
    top: -5px;
    border-top-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-left-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-bottom-color: transparent !important;
    border-right-color: transparent !important;
  }

  .cyber-tooltip--bottom-start .cyber-tooltip__arrow--bottom-start {
    left: 16px;
  }

  .cyber-tooltip--bottom-end .cyber-tooltip__arrow--bottom-end {
    right: 16px;
  }

  .cyber-tooltip--left .cyber-tooltip__arrow--left {
    right: -5px;
    border-bottom-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-right-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-top-color: transparent !important;
    border-left-color: transparent !important;
  }

  .cyber-tooltip--right .cyber-tooltip__arrow--right {
    left: -5px;
    border-bottom-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-left-color: rgba(var(--color-brand-500-rgb), 0.3) !important;
    border-top-color: transparent !important;
    border-right-color: transparent !important;
  }

  .cyber-tooltip__content {
    @apply relative z-10;
    word-wrap: break-word;
    word-break: break-word;
    overflow-wrap: break-word;
    hyphens: auto;
    line-break: anywhere;
  }

  .cyber-tooltip-fade-enter-active {
    transition:
      opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1),
      transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .cyber-tooltip-fade-leave-active {
    transition: opacity 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .cyber-tooltip-fade-enter-from {
    opacity: 0;
    transform: translateY(-4px) scale(0.95);
  }

  .cyber-tooltip-fade-leave-to {
    opacity: 0;
  }

  .cyber-tooltip--top,
  .cyber-tooltip--top-start,
  .cyber-tooltip--top-end {
    transform-origin: center bottom;
  }

  .cyber-tooltip--top.cyber-tooltip-fade-enter-from,
  .cyber-tooltip--top-start.cyber-tooltip-fade-enter-from,
  .cyber-tooltip--top-end.cyber-tooltip-fade-enter-from {
    transform: translateY(8px) scale(0.95);
  }

  .cyber-tooltip--bottom,
  .cyber-tooltip--bottom-start,
  .cyber-tooltip--bottom-end {
    transform-origin: center top;
  }

  .cyber-tooltip--bottom.cyber-tooltip-fade-enter-from,
  .cyber-tooltip--bottom-start.cyber-tooltip-fade-enter-from,
  .cyber-tooltip--bottom-end.cyber-tooltip-fade-enter-from {
    transform: translateY(-8px) scale(0.95);
  }

  .cyber-tooltip--left {
    transform-origin: right center;
  }

  .cyber-tooltip--left.cyber-tooltip-fade-enter-from {
    transform: translateX(8px) scale(0.95);
  }

  .cyber-tooltip--right {
    transform-origin: left center;
  }

  .cyber-tooltip--right.cyber-tooltip-fade-enter-from {
    transform: translateX(-8px) scale(0.95);
  }

  @media (max-width: 768px) {
    .cyber-tooltip {
      max-width: 250px;
      font-size: 12px;
      padding: 6px 10px;
    }
  }
</style>
