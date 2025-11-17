/**
 * 性能监控Composable
 * 提供性能指标收集和优化建议
 */
import { onMounted, onUnmounted, ref } from 'vue'
import { useTexts } from '@/composables/useTexts'

interface PerformanceMetrics {
  loadTime: number
  renderTime: number
  memoryUsage: number
  componentCount: number
}

interface PerformanceEntry {
  name: string
  duration: number
  startTime: number
  entryType: string
}

export function usePerformance() {
  const { $t } = useTexts()

  const metrics = ref<PerformanceMetrics>({
    loadTime: 0,
    renderTime: 0,
    memoryUsage: 0,
    componentCount: 0,
  })

  const isSupported = ref(typeof window !== 'undefined' && 'performance' in window)

  const measurePageLoad = () => {
    if (!isSupported.value) {
      return
    }

    const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming
    if (navigation) {
      metrics.value.loadTime = navigation.loadEventEnd - navigation.loadEventStart
    }
  }

  const measureRenderTime = (componentName: string) => {
    if (!isSupported.value) {
      return
    }

    const startTime = performance.now()

    return {
      end: () => {
        const endTime = performance.now()
        const duration = endTime - startTime

        if ('mark' in performance) {
          performance.mark(`${componentName}-render-start`, { startTime })
          performance.mark(`${componentName}-render-end`, { startTime: endTime })
          performance.measure(`${componentName}-render`, `${componentName}-render-start`, `${componentName}-render-end`)
        }

        metrics.value.renderTime = Math.max(metrics.value.renderTime, duration)
        return duration
      },
    }
  }

  const measureMemoryUsage = () => {
    if (!isSupported.value) {
      return
    }

    const { memory } = performance as any
    if (memory) {
      metrics.value.memoryUsage = memory.usedJSHeapSize / 1024 / 1024 // MB
    }
  }

  const getPerformanceEntries = (entryType?: string): PerformanceEntry[] => {
    if (!isSupported.value) {
      return []
    }

    const entries = entryType ? performance.getEntriesByType(entryType) : performance.getEntries()

    return entries.map((entry) => ({
      name: entry.name,
      duration: entry.duration,
      startTime: entry.startTime,
      entryType: entry.entryType,
    }))
  }

  const getOptimizationSuggestions = () => {
    const suggestions: string[] = []

    if (metrics.value.loadTime > 3000) {
      suggestions.push($t('composables.performance.suggestions.slowPageLoad'))
    }

    if (metrics.value.renderTime > 100) {
      suggestions.push($t('composables.performance.suggestions.slowRender'))
    }

    if (metrics.value.memoryUsage > 100) {
      suggestions.push($t('composables.performance.suggestions.highMemory'))
    }

    return suggestions
  }

  const clearPerformanceEntries = () => {
    if (isSupported.value && performance.clearMarks) {
      performance.clearMarks()
      performance.clearMeasures()
    }
  }

  let monitorInterval: NodeJS.Timeout | null = null

  const startMonitoring = (interval = 5000) => {
    if (!isSupported.value) {
      return
    }

    monitorInterval = setInterval(() => {
      measureMemoryUsage()
    }, interval)
  }

  const stopMonitoring = () => {
    if (monitorInterval) {
      clearInterval(monitorInterval)
      monitorInterval = null
    }
  }

  onMounted(() => {
    measurePageLoad()
    startMonitoring()
  })

  onUnmounted(() => {
    stopMonitoring()
  })

  return {
    metrics: readonly(metrics),
    isSupported: readonly(isSupported),
    measureRenderTime,
    measureMemoryUsage,
    getPerformanceEntries,
    getOptimizationSuggestions,
    clearPerformanceEntries,
    startMonitoring,
    stopMonitoring,
  }
}

export function withPerformanceMonitoring<T extends Record<string, any>>(component: T, componentName: string): T {
  if (typeof window === 'undefined') {
    return component
  }

  const originalMounted = component.mounted
  const originalUnmounted = component.unmounted

  return {
    ...component,
    mounted() {
      const perf = usePerformance()
      const measure = perf.measureRenderTime(componentName)

      if (originalMounted) {
        originalMounted.call(this)
      }

      measure.end()
    },
    unmounted() {
      if (originalUnmounted) {
        originalUnmounted.call(this)
      }
    },
  }
}
