import { computed, ref } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useSettingsStore } from '@/store/settings'
import { useSmartUpload } from '@/composables/useSmartUpload'
import { useTexts } from '@/composables/useTexts'

/**
 * 通用上传Hook - 抽离所有上传页面的共同逻辑
 * 支持默认布局和主题模式下的不同UI布局
 */
export function useUpload() {
  const toast = useToast()
  const { $t } = useTexts()
  const settingsStore = useSettingsStore()

  const fileInput = ref<HTMLInputElement | null>(null)
  const isDragging = ref(false)

  const folderId = ref<string | null>(null)
  const accessLevel = ref<'public' | 'private' | 'protected'>('public')
  const optimize = ref<boolean>(true)
  const autoRemove = ref<boolean>(false)

  const smartUpload = useSmartUpload({
    folderId: computed(() => folderId.value),
    accessLevel: computed(() => accessLevel.value),
    optimize: computed(() => optimize.value),
  })

  const {
    allUploads,
    globalProgress,
    globalSpeed,
    isUploading,
    addFiles,
    startUpload: startSmartUpload,
    pauseUpload: _pauseUpload,
    resumeUpload,
    cancelUpload: cancelSmartUpload,
    retryUpload: retrySmartUpload,
    removeUploadItem,
    clearQueue,
    clearAllSessions,
    formatFileSize: smartFormatFileSize,
    isValidImageFile: _isValidImageFile,
  } = smartUpload

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
      }
    })
  )

  const hasPendingFiles = computed(() => allUploads.value.some((file) => file.status === 'pending' || file.status === 'paused'))

  const hasUploadingFiles = computed(() => isUploading.value)

  const hasSuccessFiles = computed(() => allUploads.value.some((file) => file.status === 'completed'))

  const pendingCount = computed(
    () => allUploads.value.filter((file) => file.status === 'pending' || file.status === 'paused').length
  )

  const uploadingCount = computed(() => allUploads.value.filter((file) => file.status === 'uploading').length)

  const successCount = computed(() => allUploads.value.filter((file) => file.status === 'completed').length)

  const errorCount = computed(() => allUploads.value.filter((file) => file.status === 'failed').length)

  const totalFileSize = computed(() => allUploads.value.reduce((total, file) => total + file.file.size, 0))

  const acceptFormats = computed(() => {
    const formats = settingsStore.allowedImageFormats
    if (formats.length === 0) {
      return 'image/*'
    }
    return formats.map((format) => `image/${format.toLowerCase()}`).join(',')
  })

  const pasteShortcut = computed(() => {
    const isMac = /Mac|iPod|iPhone|iPad/.test(navigator.userAgent)
    return {
      key: isMac ? 'Cmd' : 'Ctrl',
      isMac,
    }
  })

  const triggerFileInput = () => {
    fileInput.value?.click()
  }

  const onDragOver = () => {
    isDragging.value = true
  }

  const onDragLeave = () => {
    isDragging.value = false
  }

  const onDrop = async (event: DragEvent) => {
    isDragging.value = false
    const files = event.dataTransfer?.files
    if (files) {
      await addFiles(files)
    }
  }

  const handleFileChange = async (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files) {
      await addFiles(target.files)
      target.value = ''
    }
  }

  const handlePaste = async (e: ClipboardEvent) => {
    const items = e.clipboardData?.items
    if (!items) {
      return
    }

    const files: File[] = []
    for (let i = 0; i < items.length; i++) {
      const item = items[i]
      if (item.type.indexOf('image') !== -1) {
        const file = item.getAsFile()
        if (file) {
          files.push(file)
        }
      }
    }

    if (files.length > 0) {
      const dt = new DataTransfer()
      files.forEach((file) => dt.items.add(file))
      await addFiles(dt.files)
    }
  }

  const startUpload = async () => {
    if (!hasPendingFiles.value) {
      showNoFilesMessage()
      return
    }

    await startSmartUpload()
  }

  const cancelUpload = () => {
    cancelSmartUpload()
    toast.info($t('upload.upload.toast.uploadCancelled'))
  }

  const retryUpload = (index: number) => {
    const item = allUploads.value[index]
    if (item) {
      retrySmartUpload(item.id)
    }
  }

  const resumeUploadWrapper = (index: number) => {
    const item = allUploads.value[index]
    if (item) {
      resumeUpload(item.id)
    }
  }

  const removeFile = (index: number) => {
    const item = allUploads.value[index]
    if (item) {
      removeUploadItem(item.id)
    }
  }

  const showNoFilesMessage = () => {
    toast.warning($t('upload.upload.toast.noFilesInQueue'))
  }

  const copyAllUrls = async () => {
    const urls = uploadQueue.value
      .filter((item) => item.status === 'success' && item.full_url)
      .map((item) => item.full_url)
      .join('\n')

    if (urls) {
      try {
        await navigator.clipboard.writeText(urls)
        toast.success($t('upload.upload.toast.allLinksCopied'))
      } catch (err) {
        console.error('Copy failed:', err)
        toast.error($t('upload.upload.uploadProgress.toast.copyFailed'))
      }
    } else {
      toast.warning($t('upload.upload.toast.noLinksToCopy'))
    }
  }

  const onSettingsChange = (_settings: any) => {}

  const setupEventListeners = async () => {
    document.addEventListener('paste', handlePaste)
  }

  const cleanupEventListeners = () => {
    document.removeEventListener('paste', handlePaste)
  }

  const formatFileSize = smartFormatFileSize

  return {
    fileInput,
    isDragging,
    folderId,
    accessLevel,
    optimize,
    autoRemove,

    smartUpload,
    uploadQueue,
    globalProgress,
    globalSpeed,

    hasPendingFiles,
    hasUploadingFiles,
    hasSuccessFiles,
    pendingCount,
    uploadingCount,
    successCount,
    errorCount,
    totalFileSize,
    acceptFormats,
    pasteShortcut,

    triggerFileInput,
    onDragOver,
    onDragLeave,
    onDrop,
    handleFileChange,

    startUpload,
    cancelUpload,
    retryUpload,
    resumeUploadWrapper,
    removeFile,
    showNoFilesMessage,
    copyAllUrls,
    onSettingsChange,
    formatFileSize,

    setupEventListeners,
    cleanupEventListeners,

    settingsStore,

    clearQueue,
    clearAllSessions,
  }
}
