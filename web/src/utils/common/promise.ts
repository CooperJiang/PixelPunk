/**
 * Promise相关工具函数
 * 包含等待、重试等异步操作工具
 */
import { useTextThemeStore } from '@/store/textTheme'

/**
 * 等待指定时间
 * @param ms 等待毫秒数
 * @returns Promise
 */
export function wait(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export async function retry<T>(fn: () => Promise<T>, maxAttempts: number = 3, delay: number = 300): Promise<T> {
  let lastError: Error

  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      return await fn()
    } catch (error) {
      lastError = error as Error

      if (attempt < maxAttempts) {
        await wait(delay)
        delay *= 1.5 // 增加延迟
      }
    }
  }

  const store = useTextThemeStore()
  throw lastError || new Error(store.getText('utils.common.promise.operationFailed'))
}
