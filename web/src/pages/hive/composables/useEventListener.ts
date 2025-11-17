import { onMounted, onUnmounted } from 'vue'

/* 事件监听 composable，自动管理事件绑定和解绑 */
export function useEventListener(target: 'window' | 'document', event: string, handler: EventListenerOrEventListenerObject) {
  onMounted(() => {
    if (target === 'window') {
      window.addEventListener(event, handler)
    } else {
      document.addEventListener(event, handler)
    }
  })

  onUnmounted(() => {
    if (target === 'window') {
      window.removeEventListener(event, handler)
    } else {
      document.removeEventListener(event, handler)
    }
  })
}
