import { onMounted, onUnmounted, watch, type Ref } from 'vue'

export function useKeyboardShortcuts(
  isVisible: Ref<boolean>,
  handlers: {
    onEscape?: () => void
    onSpace?: () => void
    onArrowLeft?: () => void
    onArrowRight?: () => void
    onPlus?: () => void
    onMinus?: () => void
    onR?: () => void
    onL?: () => void
    onF?: () => void
    onReset?: () => void
  }
) {
  const handleKeyDown = (e: KeyboardEvent) => {
    if (!isVisible.value) {
      return
    }

    const target = e.target as HTMLElement
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
      return
    }

    if (e.key === ' ') {
      e.preventDefault()
    }

    switch (e.key) {
      case 'Escape':
        handlers.onEscape?.()
        break
      case ' ':
        handlers.onSpace?.()
        break
      case 'ArrowLeft':
        handlers.onArrowLeft?.()
        break
      case 'ArrowRight':
        handlers.onArrowRight?.()
        break
      case '+':
      case '=':
        handlers.onPlus?.()
        break
      case '-':
        handlers.onMinus?.()
        break
      case 'r':
      case 'R':
        handlers.onR?.()
        break
      case 'l':
      case 'L':
        handlers.onL?.()
        break
      case 'f':
      case 'F':
        handlers.onF?.()
        break
      case '0':
        handlers.onReset?.()
        break
    }
  }

  const updateEventListener = (visible: boolean) => {
    if (visible) {
      document.addEventListener('keydown', handleKeyDown)
    } else {
      document.removeEventListener('keydown', handleKeyDown)
    }
  }

  watch(isVisible, (newVal) => {
    updateEventListener(newVal)
  })

  onMounted(() => {
    if (isVisible.value) {
      document.addEventListener('keydown', handleKeyDown)
    }
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeyDown)
  })

  return {
    handleKeyDown,
  }
}
