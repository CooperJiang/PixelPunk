import { ref } from 'vue'

export function useVectorDialogs() {
  const showDeleteConfirm = ref(false)
  const showBatchConfirm = ref(false)
  const showRegenerateConfirm = ref(false)
  const showRetryFailedConfirm = ref(false)

  const deleteConfirmData = ref({
    count: 0,
    expectedText: '',
    fileIds: [] as string[],
  })

  const batchConfirmData = ref({
    action: 'reset' as 'reset' | 'retry' | 'delete',
    message: '',
    fileIds: [] as string[],
  })

  const regenerateConfirmData = ref({
    message: [] as string[],
  })

  const retryFailedConfirmData = ref({
    message: '',
    count: 0,
  })

  return {
    showDeleteConfirm,
    showBatchConfirm,
    showRegenerateConfirm,
    showRetryFailedConfirm,
    deleteConfirmData,
    batchConfirmData,
    regenerateConfirmData,
    retryFailedConfirmData,
  }
}
