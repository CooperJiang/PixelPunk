/**
 * API请求Loading状态管理Hook
 * 提供自动化的loading状态管理，防止UI卡死
 */
import { type Ref, ref } from 'vue'
import type { ApiResult, ExtendedRequestConfig } from '@/utils/network/http-types'
import { useTexts } from '@/composables/useTexts'

export interface UseApiLoadingOptions {
  minLoadingTime?: number // 最小loading时间，防止闪烁，默认300ms
  autoToast?: boolean // 是否自动显示错误toast，默认true
  mode?: 'auto' | 'manual' | 'shared' // loading管理模式
  groupId?: string // loading分组标识（shared模式下生效）
}

/**
 * API请求Loading状态管理Hook
 */
export function useApiLoading(options: UseApiLoadingOptions = {}) {
  const { minLoadingTime = 300, autoToast = true, mode = 'auto', groupId } = options

  const loading = ref(false)

  const getApiConfig = (customConfig: ExtendedRequestConfig = {}): ExtendedRequestConfig => ({
    loadingTarget: loading,
    minLoadingTime,
    autoShowError: autoToast,
    loadingMode: mode,
    loadingGroup: groupId,
    ...customConfig,
  })

  const callApi = async <T>(
    apiCall: (config?: ExtendedRequestConfig) => Promise<ApiResult<T>>,
    config: ExtendedRequestConfig = {}
  ): Promise<T> => {
    const mergedConfig = getApiConfig(config)

    const result = await apiCall(mergedConfig)

    if ('success' in result && !result.success) {
      const error = new Error(result.message)
      ;(error as any).code = result.code
      ;(error as any).result = result
      throw error
    }

    return 'success' in result ? result.data : result
  }

  const setLoading = (value: boolean) => {
    loading.value = value
  }

  return {
    loading: loading as Readonly<Ref<boolean>>,
    getApiConfig,
    callApi,
    setLoading,
  }
}

export function useSimpleApiLoading() {
  const loading = ref(false)
  const { $t } = useTexts()

  const execute = async <T>(
    apiPromise: Promise<ApiResult<T>>,
    options: {
      onSuccess?: (data: T) => void
      onError?: (_error: Error) => void
      minLoadingTime?: number
    } = {}
  ) => {
    const { onSuccess, onError, minLoadingTime = 300 } = options
    const startTime = Date.now()

    loading.value = true

    try {
      const result = await apiPromise

      if ('success' in result && result.success) {
        if (onSuccess) {
          onSuccess(result.data)
        }
        return result.data
      } else if ('success' in result) {
        const errorMsg = result.message || $t('utils.network.http.operationFailed')
        if (onError) {
          onError(errorMsg)
        }
        throw new Error(errorMsg)
      }

      return result
    } catch (error: unknown) {
      const errorMsg = (error as any).message || $t('utils.network.http.requestFailed')
      if (onError) {
        onError(errorMsg)
      }
      throw error
    } finally {
      const elapsed = Date.now() - startTime
      if (elapsed < minLoadingTime) {
        setTimeout(() => {
          loading.value = false
        }, minLoadingTime - elapsed)
      } else {
        loading.value = false
      }
    }
  }

  return { loading, execute }
}

export default useApiLoading
