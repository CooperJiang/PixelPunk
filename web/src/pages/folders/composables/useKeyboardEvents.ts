import type { Ref } from 'vue'
import type { ShareButtonRef } from '../types'

/* 键盘事件处理逻辑 */
export function useKeyboardEvents(selectMode: Ref<boolean>, shareButtonRef: Ref<ShareButtonRef | undefined>) {
  const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Escape' && selectMode.value) {
      shareButtonRef.value?.cancelSelectMode()
      event.preventDefault()
    }
  }

  const setupKeyboardEvents = () => {
    window.addEventListener('keydown', handleKeydown)
  }

  const cleanupKeyboardEvents = () => {
    window.removeEventListener('keydown', handleKeydown)
  }

  return {
    handleKeydown,
    setupKeyboardEvents,
    cleanupKeyboardEvents,
  }
}
