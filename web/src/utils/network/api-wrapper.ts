/**
 * API调用包装器 - 解决Result模式和异常模式混用的问题
 */
import type { ApiResult } from './http-types'

/**
 * 将Result模式转换为异常模式
 * 使用场景：需要try-catch处理的地方
 */
export async function apiCall<T>(apiPromise: Promise<ApiResult<T>>): Promise<T> {
  const result = await apiPromise
  if (result.success) {
    return result.data
  }
  const error = new Error(result.message)
  ;(error as any).code = result.code
  ;(error as any).result = result
  throw error
}

export async function safeApiCall<T>(
  apiPromise: Promise<ApiResult<T>>,
  options: {
    onSuccess?: (data: T) => void
    onError?: (_error: Error) => void
    finallyCallback?: () => void
  } = {}
): Promise<ApiResult<T>> {
  try {
    const result = await apiPromise

    if (result.success && options.onSuccess) {
      options.onSuccess(result.data)
    } else if (!result.success && options.onError) {
      options.onError(result.message)
    }

    return result
  } catch (error: unknown) {
    const { useTexts } = await import('@/composables/useTexts')
    const { $t } = useTexts()
    const errorMessage = (error as any).message || $t('http.requestFailed')
    if (options.onError) {
      options.onError(errorMessage)
    }

    return {
      success: false,
      code: error.code || 999,
      message: errorMessage,
      data: null as T,
      timestamp: Date.now(),
    }
  } finally {
    if (options.finallyCallback) {
      options.finallyCallback()
    }
  }
}

export async function apiCallWithLoading<T>(
  apiPromise: Promise<ApiResult<T>>,
  loadingRef: { value: boolean },
  options: {
    onSuccess?: (data: T) => void
    onError?: (_error: Error) => void
    useToast?: boolean // 是否自动显示toast，默认true
  } = {}
): Promise<ApiResult<T>> {
  loadingRef.value = true

  const { useToast = true } = options

  let toast: any = null
  if (useToast) {
    try {
      const { useToast: useToastFn } = await import('@/components/Toast/useToast')
      toast = useToastFn()
    } catch {
      console.warn('Unable to load toast component')
    }
  }

  return safeApiCall(apiPromise, {
    onSuccess: (data: T) => {
      if (options.onSuccess) {
        options.onSuccess(data)
      }
    },
    onError: (_error: Error) => {
      if (useToast && toast) {
        toast.error(error)
      }
      if (options.onError) {
        options.onError(error)
      }
    },
    finallyCallback: () => {
      loadingRef.value = false
    },
  })
}

export default {
  apiCall,
  safeApiCall,
  apiCallWithLoading,
}
