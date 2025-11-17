import { nextTick, type Ref } from 'vue'
import type { FileInfo } from '@/api/types'
import {
  LOADING_MIN_DISPLAY_TIME,
  LOADING_FADEOUT_DURATION,
  COMPONENT_RENDER_WAIT_TIME,
  ELEMENT_CHECK_INTERVAL,
  RENDER_RETRY_INTERVAL,
  MAX_RENDER_ATTEMPTS,
} from '../constants'

/* 等待组件渲染完成 */
async function waitForComponentRender(): Promise<void> {
  await nextTick()
  await new Promise((resolve) => setTimeout(resolve, COMPONENT_RENDER_WAIT_TIME))

  let attempts = 0
  while (attempts < MAX_RENDER_ATTEMPTS) {
    const honeycombElements = document.querySelectorAll('.honeycomb-item')
    if (honeycombElements.length > 0) {
      await new Promise((resolve) => setTimeout(resolve, ELEMENT_CHECK_INTERVAL))
      return
    }
    attempts++
    await new Promise((resolve) => setTimeout(resolve, RENDER_RETRY_INTERVAL))
  }
}

export function useHiveInit(
  loadFiles: (page: number, append: boolean) => Promise<void>,
  preloadVisibleFiles: (fileList: FileInfo[]) => Promise<void>,
  images: Ref<FileInfo[]>,
  isInitialLoading: Ref<boolean>,
  isLoadingFadingOut: Ref<boolean>,
  loadingStartTime: Ref<number>
) {
  const initializeHoneycomb = async () => {
    try {
      await loadFiles(1, false)

      if (images.value.length > 0) {
        await preloadVisibleFiles(images.value)
        await waitForComponentRender()
      }

      const elapsedTime = Date.now() - loadingStartTime.value
      const remainingTime = Math.max(0, LOADING_MIN_DISPLAY_TIME - elapsedTime)

      if (remainingTime > 0) {
        await new Promise((resolve) => setTimeout(resolve, remainingTime))
      }

      isLoadingFadingOut.value = true

      await new Promise((resolve) => setTimeout(resolve, LOADING_FADEOUT_DURATION))
      isInitialLoading.value = false
      isLoadingFadingOut.value = false
    } catch {
      isInitialLoading.value = false
      isLoadingFadingOut.value = false
    }
  }

  return {
    initializeHoneycomb,
  }
}
