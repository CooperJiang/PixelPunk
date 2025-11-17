import { ref, type Ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'

/**
 * Loading状态管理Hook
 * 提供统一的loading状态管理，确保状态正确重置，防止UI卡死
 */
export function useLoading(initialValue = false) {
  const loading = ref(initialValue)
  const toast = useToast()
  const { $t } = useTexts()

  const withLoading = async <T>(
    asyncFn: () => Promise<T>,
    options?: {
      showErrorToast?: boolean // 是否自动显示错误toast，默认false
      errorMessage?: string // 自定义错误消息
      onError?: (_error: Error) => void // 错误回调
    }
  ): Promise<T> => {
    loading.value = true
    try {
      const result = await asyncFn()
      return result
    } catch (error) {
      console.error('useLoading execution failed:', error)

      if (options?.onError) {
        options.onError(error)
      } else if (options?.showErrorToast) {
        const message =
          options.errorMessage ||
          (error as any)?.message ||
          (error as any)?.response?.data?.message ||
          $t('common.errors.operationFailed')
        toast.error(message)
      }

      throw error
    } finally {
      loading.value = false // 确保loading状态被重置
    }
  }

  const setLoading = (value: boolean) => {
    loading.value = value
  }

  return {
    loading: loading as Ref<boolean>,
    withLoading,
    setLoading,
  }
}

export function useMultipleLoading<T extends Record<string, boolean>>(initialStates: T) {
  const loadingStates = ref({ ...initialStates }) as Ref<T>
  const toast = useToast()
  const { $t } = useTexts()

  const withLoading = async <K extends keyof T>(
    key: K,
    asyncFn: () => Promise<any>,
    options?: {
      showErrorToast?: boolean
      errorMessage?: string
      onError?: (_error: Error) => void
    }
  ) => {
    loadingStates.value[key] = true as T[K]
    try {
      return await asyncFn()
    } catch (error) {
      if (options?.onError) {
        options.onError(error)
      } else if (options?.showErrorToast) {
        const message = options.errorMessage || $t('common.errors.operationFailed')
        toast.error(message)
      }
      throw error
    } finally {
      loadingStates.value[key] = false as T[K]
    }
  }

  const setLoading = <K extends keyof T>(key: K, value: T[K]) => {
    loadingStates.value[key] = value
  }

  return {
    loading: loadingStates,
    withLoading,
    setLoading,
  }
}

export default useLoading
