<script setup lang="ts">
  import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
  import { onClickOutside, useElementBounding, useEventListener } from '@vueuse/core'
  import type { PopconfirmEmits, PopconfirmProps } from './types'
  import { Z_INDEX } from '@/constants'
  import { useTexts } from '@/composables/useTexts'

  const props = withDefaults(defineProps<PopconfirmProps>(), {
    title: '',
    description: '',
    placement: 'top',
    trigger: 'click',
    visible: false,
    showArrow: true,
    icon: 'question-circle',
    type: 'primary',
    cancelText: '',
    confirmText: '',
    confirmLoading: false,
    hideAfterConfirm: true,
    zIndex: Z_INDEX.DROPDOWN,
    width: 'auto',
    offset: 10,
  })

  const { $t } = useTexts()

  const emit = defineEmits<PopconfirmEmits>()

  const isVisible = ref(props.visible)
  const triggerRef = ref<HTMLElement | null>(null)
  const popupRef = ref<HTMLElement | null>(null)

  const triggerBounding = useElementBounding(triggerRef)

  const popupStyle = computed(() => ({
    width: typeof props.width === 'number' ? `${props.width}px` : props.width,
    zIndex: props.zIndex,
  }))

  const showOnHover = computed(() => props.trigger === 'hover')

  /* 监听props.visible变化 */
  watch(
    () => props.visible,
    (val) => {
      if (val !== isVisible.value) {
        isVisible.value = val
        if (val) {
          nextTick(() => updatePosition())
        }
      }
    }
  )

  watch(isVisible, (val) => {
    emit('update:visible', val)

    if (val) {
      nextTick(() => updatePosition())
    }
  })

  let resizeObserver: ResizeObserver | null = null
  const setupResizeObserver = () => {
    if (!triggerRef.value) {
      return
    }

    if ('ResizeObserver' in window) {
      resizeObserver = new ResizeObserver(() => {
        if (isVisible.value) {
          updatePosition()
        }
      })

      resizeObserver.observe(triggerRef.value)
    }
  }

  const show = async () => {
    isVisible.value = true
    emit('show')

    nextTick(() => updatePosition())
  }

  const hide = () => {
    isVisible.value = false
    emit('hide')
  }

  const updatePosition = () => {
    if (!popupRef.value || !triggerRef.value || !document.body.contains(triggerRef.value)) {
      if (isVisible.value) {
        hide()
      }
      return
    }

    triggerBounding.update()

    const popupRect = popupRef.value.getBoundingClientRect()
    const { placement, offset } = props

    let left = 0
    let top = 0

    switch (placement) {
      case 'top':
        left = triggerBounding.left.value + (triggerBounding.width.value - popupRect.width) / 2
        top = triggerBounding.top.value - popupRect.height - offset
        break
      case 'top-start':
        left = triggerBounding.left.value
        top = triggerBounding.top.value - popupRect.height - offset
        break
      case 'top-end':
        left = triggerBounding.right.value - popupRect.width
        top = triggerBounding.top.value - popupRect.height - offset
        break
      case 'bottom':
        left = triggerBounding.left.value + (triggerBounding.width.value - popupRect.width) / 2
        top = triggerBounding.bottom.value + offset
        break
      case 'bottom-start':
        left = triggerBounding.left.value
        top = triggerBounding.bottom.value + offset
        break
      case 'bottom-end':
        left = triggerBounding.right.value - popupRect.width
        top = triggerBounding.bottom.value + offset
        break
      case 'left':
        left = triggerBounding.left.value - popupRect.width - offset
        top = triggerBounding.top.value + (triggerBounding.height.value - popupRect.height) / 2
        break
      case 'right':
        left = triggerBounding.right.value + offset
        top = triggerBounding.top.value + (triggerBounding.height.value - popupRect.height) / 2
        break
    }

    left = Math.max(10, Math.min(left, window.innerWidth - popupRect.width - 10))
    top = Math.max(10, Math.min(top, window.innerHeight - popupRect.height - 10))

    if (popupRef.value) {
      popupRef.value.style.left = `${left}px`
      popupRef.value.style.top = `${top}px`
    }
  }

  const debounce = (fn: () => void, delay: number) => {
    let timer: number | null = null
    return function () {
      if (timer) {
        window.clearTimeout(timer)
      }
      timer = window.setTimeout(() => {
        fn()
      }, delay)
    }
  }

  const debouncedUpdatePosition = debounce(updatePosition, 50)

  const handleTriggerClick = (e: MouseEvent) => {
    if (props.trigger === 'click') {
      e.stopPropagation() // 防止冒泡触发外部点击事件
      e.preventDefault() // 阻止默认行为

      if (isVisible.value) {
        hide()
      } else {
        show()
      }
    }
  }

  const handleCancel = () => {
    hide()
    emit('cancel')
  }

  const handleConfirm = () => {
    emit('confirm')

    if (props.hideAfterConfirm) {
      hide()
    }
  }

  onClickOutside(
    popupRef,
    (event) => {
      if (isVisible.value && triggerRef.value && !triggerRef.value.contains(event.target as Node)) {
        hide()
      }
    },
    {
      ignore: [],
      capture: true,
      detectIframe: true,
    }
  )

  useEventListener(window, 'resize', () => {
    if (isVisible.value) {
      debouncedUpdatePosition()
    }
  })

  useEventListener(
    window,
    'scroll',
    () => {
      if (isVisible.value) {
        debouncedUpdatePosition()
      }
    },
    { capture: true }
  )

  let mutationObserver: MutationObserver | null = null
  const setupMutationObserver = () => {
    if (!triggerRef.value) {
      return
    }

    let listContainer = triggerRef.value.parentElement
    while (listContainer && !['UL', 'OL', 'DIV'].includes(listContainer.tagName)) {
      listContainer = listContainer.parentElement
    }

    if (listContainer) {
      mutationObserver = new MutationObserver(() => {
        if (isVisible.value) {
          debouncedUpdatePosition()
        }
      })

      mutationObserver.observe(listContainer, {
        childList: true,
        subtree: false,
        attributes: false,
      })
    }
  }

  onMounted(() => {
    if (isVisible.value) {
      nextTick(() => updatePosition())
    }

    setupResizeObserver()
    setupMutationObserver()
  })

  onBeforeUnmount(() => {
    if (resizeObserver) {
      resizeObserver.disconnect()
      resizeObserver = null
    }

    if (mutationObserver) {
      mutationObserver.disconnect()
      mutationObserver = null
    }
  })

  defineExpose({
    show,
    hide,
    updatePosition,
  })
</script>

<template>
  <div class="cyber-popconfirm">
    <div
      ref="triggerRef"
      class="cyber-popconfirm__trigger"
      @click="handleTriggerClick"
      @mouseenter="showOnHover && show()"
      @mouseleave="showOnHover && hide()"
    >
      <slot name="reference" />
    </div>

    <Teleport to="body">
      <Transition name="cyber-popconfirm-fade">
        <div
          v-show="isVisible"
          ref="popupRef"
          class="cyber-popconfirm__popup"
          :class="[`cyber-popconfirm__popup--${placement}`, { 'cyber-popconfirm__popup--with-arrow': showArrow }]"
          :style="popupStyle"
        >
          <div v-if="showArrow" class="cyber-popconfirm__arrow" />

          <div class="cyber-popconfirm__content">
            <div class="cyber-popconfirm__header">
              <i v-if="icon" class="cyber-popconfirm__icon" :class="[`fas fa-${icon}`, `cyber-popconfirm__icon--${type}`]" />
              <div class="cyber-popconfirm__title">{{ title || $t('popconfirm.title') }}</div>
            </div>

            <div v-if="description" class="cyber-popconfirm__description">
              {{ description }}
            </div>

            <div v-if="$slots.content" class="cyber-popconfirm__custom-content">
              <slot name="content" />
            </div>

            <div class="cyber-popconfirm__buttons">
              <CyberButton type="secondary" size="small" @click="handleCancel">
                {{ cancelText || $t('actions.cancel') }}
              </CyberButton>
              <CyberButton
                :type="type === 'warning' ? 'danger' : type"
                :loading="confirmLoading"
                size="small"
                @click="handleConfirm"
              >
                {{ confirmText || $t('popconfirm.confirm') }}
              </CyberButton>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
  .cyber-popconfirm {
    @apply relative inline-block;
  }

  .cyber-popconfirm__trigger {
    @apply inline-block cursor-pointer;
  }

  .cyber-popconfirm__popup {
    @apply pointer-events-auto fixed max-w-[300px] rounded-lg p-4;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.98) 0%,
      rgba(var(--color-background-900-rgb), 0.98) 100%
    );
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.5),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.25);
    backdrop-filter: blur(12px);
    z-index: 9999;
  }

  .cyber-popconfirm__arrow {
    @apply pointer-events-none absolute h-3 w-3;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.98) 0%,
      rgba(var(--color-background-900-rgb), 0.98) 100%
    );
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    transform: rotate(45deg);
  }

  .cyber-popconfirm__popup--top .cyber-popconfirm__arrow {
    @apply -bottom-1.5 left-1/2 -ml-1.5 border-l-0 border-t-0;
  }

  .cyber-popconfirm__popup--top-start .cyber-popconfirm__arrow {
    bottom: -6px;
    left: 16px;
    border-top: none;
    border-left: none;
  }

  .cyber-popconfirm__popup--top-end .cyber-popconfirm__arrow {
    bottom: -6px;
    right: 16px;
    border-top: none;
    border-left: none;
  }

  .cyber-popconfirm__popup--bottom .cyber-popconfirm__arrow {
    top: -6px;
    left: 50%;
    margin-left: -6px;
    border-bottom: none;
    border-right: none;
  }

  .cyber-popconfirm__popup--bottom-start .cyber-popconfirm__arrow {
    top: -6px;
    left: 16px;
    border-bottom: none;
    border-right: none;
  }

  .cyber-popconfirm__popup--bottom-end .cyber-popconfirm__arrow {
    top: -6px;
    right: 16px;
    border-bottom: none;
    border-right: none;
  }

  .cyber-popconfirm__popup--left .cyber-popconfirm__arrow {
    right: -6px;
    top: 50%;
    margin-top: -6px;
    border-left: none;
    border-bottom: none;
  }

  .cyber-popconfirm__popup--right .cyber-popconfirm__arrow {
    left: -6px;
    top: 50%;
    margin-top: -6px;
    border-right: none;
    border-top: none;
  }

  .cyber-popconfirm__content {
    position: relative;
    z-index: 1;
  }

  .cyber-popconfirm__header {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    min-height: 24px;
  }

  .cyber-popconfirm__icon {
    font-size: 18px;
    margin-right: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border-radius: var(--radius-sm);
    flex-shrink: 0;
  }

  .cyber-popconfirm__icon--primary {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-400);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-popconfirm__icon--success {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-400);
    border: 1px solid rgba(var(--color-success-rgb), 0.3);
  }

  .cyber-popconfirm__icon--warning {
    background: rgba(var(--color-warning-rgb), 0.15);
    color: var(--color-warning-400);
    border: 1px solid rgba(var(--color-warning-rgb), 0.3);
  }

  .cyber-popconfirm__icon--danger {
    background: rgba(var(--color-error-rgb), 0.15);
    color: var(--color-error-400);
    border: 1px solid rgba(var(--color-error-rgb), 0.3);
  }

  .cyber-popconfirm__title {
    font-weight: 600;
    font-size: 15px;
    color: var(--color-content-heading);
    line-height: 1.4;
    letter-spacing: 0.01em;
  }

  .cyber-popconfirm__description {
    font-size: 13px;
    color: var(--color-content-muted);
    margin-bottom: 16px;
    line-height: 1.6;
    padding-left: 42px;
  }

  .cyber-popconfirm__custom-content {
    margin-bottom: 16px;
    padding-left: 42px;
  }

  .cyber-popconfirm__buttons {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 16px;
    padding-top: 12px;
    border-top: 1px solid rgba(var(--color-border-default-rgb), 0.2);
  }

  .cyber-popconfirm-fade-enter-active,
  .cyber-popconfirm-fade-leave-active {
    transition:
      opacity 0.25s cubic-bezier(0.4, 0, 0.2, 1),
      transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .cyber-popconfirm-fade-enter-from {
    opacity: 0;
    transform: scale(0.92) translateY(-4px);
  }

  .cyber-popconfirm-fade-leave-to {
    opacity: 0;
    transform: scale(0.95) translateY(2px);
  }

  .cyber-popconfirm__button {
    padding: 0.25rem 0.75rem;
    font-size: 12px;
    font-weight: 500;
    min-width: 60px;
  }
</style>
