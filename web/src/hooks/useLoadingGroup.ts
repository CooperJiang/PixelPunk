/**
 * Loading分组管理器
 * 解决多API调用需要使用返回结果的场景
 */
import { type Ref, ref } from 'vue'
import type { ExtendedRequestConfig } from '@/utils/network/http-types'

export interface LoadingGroupOptions {
  minLoadingTime?: number
  autoToast?: boolean
}

/**
 * Loading分组管理Hook
 * 适用于：多个API调用，需要各自返回结果，共享一个loading状态
 */
export function useLoadingGroup(options: LoadingGroupOptions = {}) {
  const { minLoadingTime = 300, autoToast = true } = options

  const loading = ref(false)
  const activeRequests = ref(0)
  const startTime = ref(0)

  const startRequest = () => {
    if (activeRequests.value === 0) {
      loading.value = true
      startTime.value = Date.now()
    }
    activeRequests.value++
  }

  const finishRequest = () => {
    activeRequests.value--

    if (activeRequests.value <= 0) {
      activeRequests.value = 0

      const elapsed = Date.now() - startTime.value
      if (elapsed < minLoadingTime) {
        setTimeout(() => {
          loading.value = false
        }, minLoadingTime - elapsed)
      } else {
        loading.value = false
      }
    }
  }

  const wrapApiCall = async <T>(apiCall: Promise<T>): Promise<T> => {
    startRequest()
    try {
      return await apiCall
    } finally {
      finishRequest()
    }
  }

  const getApiConfig = (customConfig: ExtendedRequestConfig = {}): ExtendedRequestConfig => ({
    loadingMode: 'manual',
    autoShowError: autoToast,
    ...customConfig,
  })

  const executeGroup = async <T extends readonly unknown[]>(
    ...apiCalls: readonly [...{ [K in keyof T]: () => Promise<T[K]> }]
  ): Promise<T> => {
    const wrappedCalls = apiCalls.map((call) => wrapApiCall(call()))
    return Promise.all(wrappedCalls) as Promise<T>
  }

  const forceReset = () => {
    loading.value = false
    activeRequests.value = 0
  }

  return {
    loading: loading as Readonly<Ref<boolean>>,
    wrapApiCall,
    getApiConfig,
    executeGroup,
    forceReset,
  }
}

export default useLoadingGroup
