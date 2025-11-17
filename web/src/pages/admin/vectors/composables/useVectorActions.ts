import { ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { batchVectorAction, retryVector } from '@/api/admin/vectors'
import { useTexts } from '@/composables/useTexts'

export function useVectorActions() {
  const toast = useToast()
  const { $t } = useTexts()

  const isProcessing = ref(false)

  const handleRetryVector = async (fileId: string, onSuccess?: () => void) => {
    if (!fileId) {
      toast.error($t('admin.vectors.toast.fileIdRequired'))
      return
    }

    try {
      isProcessing.value = true

      const result = await retryVector({ file_id: fileId })

      if (result && result.success) {
        toast.success($t('admin.vectors.toast.retrySubmitted'))
        setTimeout(() => {
          onSuccess?.()
        }, 1000)
      } else {
        throw new Error(result?.message || $t('admin.vectors.errors.retryFailed'))
      }
    } catch (error: unknown) {
      toast.error($t('admin.vectors.toast.retryFailed').replace('{message}', (error as Error).message))
    } finally {
      isProcessing.value = false
    }
  }

  const handleBatchActionExecute = async (fileIds: string[], action: 'reset' | 'retry' | 'delete', onSuccess?: () => void) => {
    try {
      isProcessing.value = true

      const actionText = {
        reset: $t('admin.vectors.actions.reset'),
        retry: $t('admin.vectors.actions.retry'),
        delete: $t('admin.vectors.actions.delete'),
      }[action]

      const result = await batchVectorAction({ file_ids: fileIds, action })

      if (result && result.success) {
        toast.success(
          $t('admin.vectors.toast.batchActionSuccess')
            .replace('{action}', actionText)
            .replace('{count}', fileIds.length.toString())
        )
        setTimeout(() => {
          onSuccess?.()
        }, 1000)
      } else {
        throw new Error(result?.message || $t('admin.vectors.errors.batchActionFailed', { action: actionText }))
      }
    } catch (error: unknown) {
      toast.error($t('admin.vectors.toast.batchActionFailed').replace('{message}', (error as Error).message))
    } finally {
      isProcessing.value = false
    }
  }

  return {
    isProcessing,

    handleRetryVector,
    handleBatchActionExecute,
  }
}
