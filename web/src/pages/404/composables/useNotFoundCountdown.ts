import { computed, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { NotFoundCountdownOptions } from '../types'

/* 404页面倒计时逻辑 */
export function useNotFoundCountdown(options: NotFoundCountdownOptions = {}) {
  const router = useRouter()
  const redirectPath = options.redirectPath ?? '/'
  const totalSeconds = Math.max(options.seconds ?? 5, 1)
  const countdown = ref(totalSeconds)
  let intervalTimer: number | null = null

  const progress = computed(() => {
    const consumed = totalSeconds - countdown.value
    return (consumed / totalSeconds) * 100
  })

  function clear() {
    if (intervalTimer) {
      clearInterval(intervalTimer)
      intervalTimer = null
    }
  }

  function tick() {
    countdown.value -= 1
    if (countdown.value <= 0) {
      clear()
      router.push(redirectPath)
    }
  }

  function start() {
    clear()
    countdown.value = totalSeconds
    intervalTimer = window.setInterval(() => {
      tick()
    }, 1000)
  }

  onUnmounted(() => {
    clear()
  })

  return {
    countdown,
    progress,
    start,
    clear,
  }
}
