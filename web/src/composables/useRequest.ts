/**
 * 增强的请求Hook
 * 支持重试、缓存、取消等功能
 */
import { onUnmounted, ref, type Ref } from 'vue'
import { useCache } from './useCache'

interface RequestOptions<T> {
  cacheKey?: string
  cacheTTL?: number
  enableCache?: boolean

  maxRetries?: number
  retryDelay?: number

  immediate?: boolean
  onSuccess?: (data: T) => void
  onError?: (_error: Error) => void
}

interface UseRequestReturn<T> {
  data: Ref<T | null>
  loading: Ref<boolean>
  error: Ref<Error | null>
  execute: (forceRefresh?: boolean) => Promise<T | null>
  cancel: () => void
  refresh: () => Promise<T | null>
}

export function useRequest<T = any>(requestFn: () => Promise<T>, options: RequestOptions<T> = {}): UseRequestReturn<T> {
  const {
    cacheKey,
    cacheTTL = 5 * 60 * 1000,
    enableCache = Boolean(cacheKey),
    maxRetries = 3,
    retryDelay = 1000,
    immediate = false,
    onSuccess,
    onError,
  } = options

  const data = ref<T | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const { withCache } = useCache()
  let cancelToken: { cancelled: boolean } | null = null

  const requestWithRetry = async (fn: () => Promise<T>, retriesLeft: number = maxRetries): Promise<T> => {
    try {
      return await fn()
    } catch (err) {
      if (retriesLeft > 0 && !cancelToken?.cancelled) {
        await new Promise((resolve) => setTimeout(resolve, retryDelay * (maxRetries - retriesLeft + 1)))
        return requestWithRetry(fn, retriesLeft - 1)
      }
      throw err
    }
  }

  const execute = async (forceRefresh = false): Promise<T | null> => {
    if (loading.value) {
      return data.value
    }

    error.value = null
    loading.value = true
    cancelToken = { cancelled: false }

    try {
      let result: T

      if (enableCache && cacheKey) {
        result = await withCache(
          cacheKey,
          () =>
            requestWithRetry(() => {
              if (cancelToken?.cancelled) {
                throw new Error('Request cancelled')
              }
              return requestFn()
            }),
          { ttl: cacheTTL, forceRefresh }
        )
      } else {
        result = await requestWithRetry(() => {
          if (cancelToken?.cancelled) {
            throw new Error('Request cancelled')
          }
          return requestFn()
        })
      }

      if (!cancelToken?.cancelled) {
        data.value = result
        onSuccess?.(result)
      }

      return result
    } catch (err) {
      if (!cancelToken?.cancelled) {
        error.value = err as Error
        onError?.(err as Error)
      }
      return null
    } finally {
      if (!cancelToken?.cancelled) {
        loading.value = false
      }
    }
  }

  const cancel = () => {
    if (cancelToken) {
      cancelToken.cancelled = true
    }
    loading.value = false
  }

  const refresh = () => execute(true)

  if (immediate) {
    execute()
  }

  onUnmounted(() => {
    cancel()
  })

  return {
    data,
    loading,
    error,
    execute,
    cancel,
    refresh,
  }
}

export function useConcurrentRequests<T extends Record<string, any>>(
  requests: { [K in keyof T]: () => Promise<T[K]> },
  _options: {
    enableCache?: boolean
    cacheTTL?: number
  } = {}
) {
  const data = ref<Partial<T>>({})
  const loading = ref(false)
  const errors = ref<Partial<Record<keyof T, Error>>>({})

  const execute = async (): Promise<Partial<T>> => {
    loading.value = true
    errors.value = {}

    const promises = Object.entries(requests).map(async ([key, requestFn]) => {
      try {
        const result = await (requestFn as () => Promise<any>)()
        return { key, result, error: null }
      } catch (_error) {
        return { key, result: null, error: _error as Error }
      }
    })

    const results = await Promise.allSettled(promises)

    results.forEach((promiseResult) => {
      if (promiseResult.status === 'fulfilled') {
        const { key, result, error } = promiseResult.value
        if (error) {
          errors.value[key as keyof T] = error
        } else {
          data.value[key as keyof T] = result
        }
      }
    })

    loading.value = false
    return data.value
  }

  return {
    data,
    loading,
    errors,
    execute,
  }
}
