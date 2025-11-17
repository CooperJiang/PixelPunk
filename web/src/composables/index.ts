/**
 * Composables统一导出
 * 提供所有组合式函数的统一入口

/* 缓存相关 */
export { useCache } from './useCache'

/* 请求相关 */
export { useRequest, useConcurrentRequests } from './useRequest'

/* 性能监控 */
export { usePerformance, withPerformanceMonitoring } from './usePerformance'

/* 工具函数 - 可选使用，不与现有功能冲突 */
export { useLazyload, createLazyLoader, getDefaultLazyLoader } from '@/utils/lazyload'
export { useVirtualScroll, createVirtualScrollProps } from '@/utils/virtual-scroll'

/* 类型导出 */
export type { RequestOptions, UseRequestReturn } from './useRequest'
export type { PerformanceMetrics } from './usePerformance'

/* 常用组合 */
export function useImageOptimization() {
  const { observe, unobserve } = useLazyload({
    placeholder:
      'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iI2Y1ZjVmNSIvPjx0ZXh0IHg9IjUwIiB5PSI1MCIgZm9udC1mYW1pbHk9IkFyaWFsIiBmb250LXNpemU9IjE0IiBmaWxsPSIjOTk5IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkeT0iMC4zZW0iPkxvYWRpbmcuLi48L3RleHQ+PC9zdmc+',
    fadeIn: true,
    retryTimes: 3,
  })

  const { measureRenderTime } = usePerformance()

  return {
    observe,
    unobserve,
    measureRenderTime,
  }
}

export function useDataFetching<T>(requestFn: () => Promise<T>, cacheKey?: string) {
  return useRequest(requestFn, {
    cacheKey,
    enableCache: Boolean(cacheKey),
    maxRetries: 3,
    immediate: true,
  })
}

export function usePaginatedData<T>(
  requestFn: (page: number, size: number) => Promise<{ data: T[]; total: number }>,
  initialPage = 1,
  pageSize = 20
) {
  const page = ref(initialPage)
  const size = ref(pageSize)
  const total = ref(0)
  const data = ref<T[]>([])

  const { loading, error, execute } = useRequest(() => requestFn(page.value, size.value), {
    cacheKey: `paginated_${page.value}_${size.value}`,
    immediate: true,
    onSuccess: (result) => {
      data.value = result.data
      total.value = result.total
    },
  })

  const nextPage = () => {
    if (page.value * size.value < total.value) {
      page.value++
      execute()
    }
  }

  const prevPage = () => {
    if (page.value > 1) {
      page.value--
      execute()
    }
  }

  const goToPage = (targetPage: number) => {
    if (targetPage >= 1 && targetPage <= Math.ceil(total.value / size.value)) {
      page.value = targetPage
      execute()
    }
  }

  return {
    data: readonly(data),
    loading: readonly(loading),
    _error: error(error),
    page: readonly(page),
    size: readonly(size),
    total: readonly(total),
    nextPage,
    prevPage,
    goToPage,
    refresh: execute,
  }
}
