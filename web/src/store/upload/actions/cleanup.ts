/**
 * 负责队列清空和资源清理

 */
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import { regularQueue, chunkedUploadInstance, runningUploads, previewUrls } from '../state'

/* 清空队列 */
export function clearQueue() {
  chunkedUploadInstance.value?.clearQueue()

  regularQueue.value = []

  previewUrls.value.forEach((url) => URL.revokeObjectURL(url))
  previewUrls.value.clear()

  runningUploads.value.clear()
}

export function clearAllSessions() {
  const toast = useToast()
  const { $t } = useTexts()
  chunkedUploadInstance.value?.clearAllSessions()
  toast.success($t('upload.uploadProgress.toast.cacheCleared'))
}

export function cleanup() {
  previewUrls.value.forEach((url) => URL.revokeObjectURL(url))
  previewUrls.value.clear()
}
