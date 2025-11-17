import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useUploadStore } from '@/store/upload'
import { useTexts } from '@/composables/useTexts'

/**
 * 全局上传管理Hook
 * 提供与原useSmartUpload兼容的接口，但底层使用全局Store
 */
export function useGlobalUpload(options?: { folderId?: any; accessLevel?: any; optimize?: any }) {
  const uploadStore = useUploadStore()
  const { $t } = useTexts()

  const { allUploads, globalProgress, globalSpeed, isUploading, statistics, maxConcurrentUploads } = storeToRefs(uploadStore)

  if (options) {
    const globalOptions: any = {}

    if (options.folderId) {
      globalOptions.folderId =
        typeof options.folderId === 'object' && 'value' in options.folderId ? options.folderId.value : options.folderId
    }

    if (options.accessLevel) {
      globalOptions.accessLevel =
        typeof options.accessLevel === 'object' && 'value' in options.accessLevel
          ? options.accessLevel.value
          : options.accessLevel
    }

    if (options.optimize !== undefined) {
      globalOptions.optimize =
        typeof options.optimize === 'object' && 'value' in options.optimize ? options.optimize.value : options.optimize
    }

    uploadStore.setGlobalOptions(globalOptions)
  }

  const uploadQueue = computed(() =>
    allUploads.value.map((item) => {
      let status = 'pending'
      if (item.status === 'completed') {
        status = 'success'
      } else if (item.status === 'failed') {
        status = 'error'
      } else if (item.status === 'uploading' || item.status === 'preparing') {
        status = 'uploading'
      } else if (item.status === 'analyzing') {
        status = 'analyzing'
      } else if (item.status === 'instant') {
        status = 'instant'
      } else if (item.status === 'retrying') {
        status = 'retrying'
      } else if (item.status === 'paused') {
        status = 'paused'
      } else if (item.status === 'pending') {
        status = 'pending'
      }

      let resultUrls = {}
      if (item.status === 'completed' && item.result) {
        const resultData = item.result.data || item.result
        resultUrls = {
          full_url: resultData?.full_url,
          full_thumb_url: resultData?.full_thumb_url,
          url: resultData?.url,
          id: resultData?.id,
        }
      }

      const isInstantUpload =
        item.statusMessage === $t('composables.useGlobalUpload.statusMessages.instantComplete') ||
        item.statusMessage === $t('composables.useGlobalUpload.statusMessages.instantSuccess') ||
        item.statusMessage?.includes($t('composables.useGlobalUpload.statusMessages.instantKeyword')) ||
        status === 'instant'

      return {
        file: item.file,
        name: item.file.name,
        size: item.file.size,
        progress: item.progress,
        status,
        ...resultUrls,
        preview: item.preview || '',
        error: item.error,
        id: item.id || '',
        uniqueId: item.id,
        previewLoaded: true,
        itemLoaded: true,
        statusMessage: item.statusMessage,
        speed: item.speed,
        remainingTime: item.remainingTime,
        is_duplicate: isInstantUpload, // 标记为秒传文件
        type: item.type,
      }
    })
  )

  const addFiles = (files: FileList | File[]) => uploadStore.addFiles(files)

  const startUpload = () => uploadStore.startUpload()

  const pauseUpload = (itemId?: string) => {
    uploadStore.pauseUpload(itemId)
  }

  const resumeUpload = (itemId: string) => {
    const item = allUploads.value[parseInt(itemId)]
    if (item) {
      return uploadStore.resumeUpload(item.id)
    }
  }

  const cancelUpload = (itemId?: string) => {
    uploadStore.cancelUpload(itemId)
  }

  const retryUpload = (index: number) => {
    const item = allUploads.value[index]
    if (item) {
      return uploadStore.retryUpload(item.id)
    }
  }

  const removeUploadItem = (itemId: string) => {
    uploadStore.removeUploadItem(itemId)
  }

  const removeFile = (index: number) => {
    const item = allUploads.value[index]
    if (item) {
      uploadStore.removeUploadItem(item.id)
    }
  }

  const clearQueue = () => {
    uploadStore.clearQueue()
  }

  const clearAllSessions = () => {
    uploadStore.clearAllSessions()
  }

  const pendingCount = computed(() => statistics.value.pending)
  const uploadingCount = computed(() => statistics.value.uploading)
  const successCount = computed(() => statistics.value.completed)
  const errorCount = computed(() => statistics.value.failed)

  const hasPendingFiles = computed(() => pendingCount.value > 0)
  const hasUploadingFiles = computed(() => uploadingCount.value > 0)
  const hasSuccessFiles = computed(() => successCount.value > 0)

  const totalFileSize = computed(() => allUploads.value.reduce((total, file) => total + file.file.size, 0))

  const copyAllUrls = async () => {
    const urls = allUploads.value
      .filter((file) => file.status === 'completed' && file.result?.full_url)
      .map((file) => file.result?.full_url)
      .filter((url): url is string => !!url)
      .join('\n')

    if (urls) {
      try {
        await navigator.clipboard.writeText(urls)
        return true
      } catch (_err) {
        return false
      }
    }
    return false
  }

  return {
    allUploads,
    uploadQueue, // 兼容原格式
    globalProgress,
    globalSpeed,
    isUploading,
    runningUploads: computed(() => uploadStore.runningUploads),
    maxConcurrentUploads,

    pendingCount,
    uploadingCount,
    successCount,
    errorCount,
    hasPendingFiles,
    hasUploadingFiles,
    hasSuccessFiles,
    totalFileSize,

    addFiles,
    startUpload,
    pauseUpload,
    resumeUpload,
    cancelUpload,
    retryUpload,
    removeUploadItem,
    removeFile, // 兼容原方法名
    clearQueue,
    clearAllSessions,
    copyAllUrls,

    formatFileSize: uploadStore.formatFileSize,
    isValidImageFile: uploadStore.isValidImageFile,

    uploadStore,
  }
}
