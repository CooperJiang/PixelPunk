import { computed, onUnmounted, ref, watch } from 'vue'
import { useChunkedUpload } from './useChunkedUpload'
import { useUploadConfig } from './useUploadConfig'
import { generateId } from '@/utils/fileHash'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'

/* 智能上传配置 */
interface SmartUploadOptions {
  folderId?: any
  accessLevel?: any
  optimize?: any
}

/* 普通上传项接口 */
interface RegularUploadItem {
  id: string
  file: File
  status: 'pending' | 'uploading' | 'completed' | 'failed'
  progress: number
  error?: string
  result?: any
  speed: number
  remainingTime: number
  statusMessage?: string
}

export function useSmartUpload(options: SmartUploadOptions = {}) {
  const uploadConfig = useUploadConfig()

  const toast = useToast()
  const { $t } = useTexts()

  const chunkedUpload = ref<any>(null)
  const isChunkedUploadReady = ref(false)

  const initializeChunkedUpload = () => {
    if (!uploadConfig.isConfigLoaded.value) {
      return
    }

    if (chunkedUpload.value) {
      return
    }

    chunkedUpload.value = useChunkedUpload({
      chunkSize: uploadConfig.getChunkSize(),
      concurrency: uploadConfig.getChunkedMaxConcurrency(),
      retryCount: uploadConfig.getRetryCount(),
      folderId: options.folderId,
      accessLevel: options.accessLevel,
      optimize: options.optimize,
    })
    isChunkedUploadReady.value = true
  }

  watch(
    () => uploadConfig.isConfigLoaded.value,
    (loaded) => {
      if (loaded) {
        initializeChunkedUpload()
      }
    },
    { immediate: true }
  )

  const regularUploadQueue = ref<RegularUploadItem[]>([])

  const createdUrls = new Set<string>()

  const runningUploads = ref(new Set<string>()) // 当前正在上传的文件ID集合

  const classifyFile = (file: File) => {
    if (!uploadConfig.isConfigLoaded.value) {
      console.warn('⚠️ [SmartUpload] 配置未加载，使用默认阈值 20MB')
      const shouldUseChunked = file.size >= 20 * 1024 * 1024
      return shouldUseChunked
    }
    const shouldUseChunked = uploadConfig.shouldUseChunkedUpload(file.size)
    return shouldUseChunked
  }

  const allUploads = computed(() => {
    const regular = regularUploadQueue.value.map((item) => {
      let previewUrl: string | undefined
      if (item.file.type.startsWith('image/')) {
        previewUrl = URL.createObjectURL(item.file)
        createdUrls.add(previewUrl)
      }

      return {
        ...item,
        type: 'regular' as const,
        totalChunks: 1,
        uploadedChunks: item.status === 'completed' ? 1 : 0,
        preview: previewUrl,
      }
    })

    let chunked: any[] = []
    if (isChunkedUploadReady.value && chunkedUpload.value?.uploadQueue) {
      const chunkedQueue = chunkedUpload.value.uploadQueue
      chunked = chunkedQueue.map((item: any) => {
        let previewUrl: string | undefined
        if (item.file.type.startsWith('image/')) {
          previewUrl = URL.createObjectURL(item.file)
          createdUrls.add(previewUrl)
        }

        return {
          ...item,
          type: 'chunked' as const,
          preview: previewUrl,
        }
      })
    }

    const totalUploads = [...regular, ...chunked]

    return totalUploads
  })

  const globalProgress = computed(() => {
    const total = allUploads.value.length
    if (total === 0) {
      return 0
    }

    const totalProgress = allUploads.value.reduce((sum, item) => sum + item.progress, 0)
    return Math.round(totalProgress / total)
  })

  const globalSpeed = computed(() => allUploads.value.reduce((sum, item) => sum + (item.speed || 0), 0))

  const isUploading = computed(
    () =>
      runningUploads.value.size > 0 || allUploads.value.some((item) => item.status === 'uploading' || item.status === 'preparing')
  )

  const addFiles = async (files: FileList | File[]) => {
    if (!uploadConfig.isConfigLoaded.value) {
      await new Promise((resolve) => {
        const unwatch = watch(
          () => uploadConfig.isConfigLoaded.value,
          (loaded) => {
            if (loaded) {
              unwatch()
              resolve(true)
            }
          }
        )
      })
    }

    if (!isChunkedUploadReady.value) {
      initializeChunkedUpload()

      await new Promise((resolve) => setTimeout(resolve, 100))
    }

    const fileArray = Array.from(files)
    const chunkedFiles: File[] = []
    const regularFiles: File[] = []

    for (const file of fileArray) {
      if (!uploadConfig.isAllowedFileType(file.type)) {
        const errorMsg = $t('upload.upload.smartUpload.errors.fileTypeNotSupported', {
          name: file.name,
          type: file.type,
        })
        toast.error(errorMsg)
        continue
      }

      if (!uploadConfig.isAllowedFileSize(file.size)) {
        const errorMsg = $t('upload.upload.smartUpload.errors.fileSizeExceeded', {
          name: file.name,
          size: uploadConfig.formatFileSize(file.size),
        })
        toast.error(errorMsg)
        continue
      }

      const shouldUseChunked = classifyFile(file)
      if (shouldUseChunked) {
        chunkedFiles.push(file)
      } else {
        regularFiles.push(file)
      }
    }

    if (chunkedFiles.length > 0) {
      if (!chunkedUpload.value) {
        toast.error($t('upload.upload.smartUpload.errors.chunkedUploadNotReady'))
      } else {
        await chunkedUpload.value.addFiles(chunkedFiles)
      }
    }

    for (const file of regularFiles) {
      regularUploadQueue.value.push({
        id: generateId(),
        file,
        status: 'pending',
        progress: 0,
        speed: 0,
        remainingTime: 0,
        statusMessage: $t('upload.upload.smartUpload.status.preparing'),
      })
    }
  }

  const uploadRegularFile = async (item: RegularUploadItem) => {
    try {
      item.status = 'uploading'
      item.progress = 0
      item.statusMessage = $t('upload.upload.smartUpload.status.uploading')

      const startTime = Date.now()

      const { uploadFile } = await import('@/api/file')

      const result = await uploadFile(
        item.file,
        {
          folder_id:
            typeof options.folderId === 'object' && 'value' in options.folderId ? options.folderId.value : options.folderId,
          access_level:
            typeof options.accessLevel === 'object' && 'value' in options.accessLevel
              ? options.accessLevel.value
              : options.accessLevel,
          optimize:
            typeof options.optimize === 'object' && 'value' in options.optimize ? options.optimize.value : options.optimize,
        },
        (progressEvent) => {
          if (progressEvent.lengthComputable) {
            item.progress = Math.round((progressEvent.loaded / progressEvent.total) * 100)

            const elapsed = (Date.now() - startTime) / 1000
            if (elapsed > 0) {
              item.speed = Math.round(progressEvent.loaded / elapsed / 1024) // KB/s

              const remainingBytes = progressEvent.total - progressEvent.loaded
              item.remainingTime = item.speed > 0 ? Math.round(remainingBytes / (item.speed * 1024)) : 0

              const speedText = item.speed > 1024 ? `${(item.speed / 1024).toFixed(1)} MB/s` : `${item.speed} KB/s`

              const remainingText =
                item.remainingTime > 60
                  ? $t('upload.upload.smartUpload.time.minutes', {
                      minutes: Math.floor(item.remainingTime / 60),
                      seconds: item.remainingTime % 60,
                    })
                  : $t('upload.upload.smartUpload.time.seconds', { seconds: item.remainingTime })

              item.statusMessage = $t('upload.upload.smartUpload.status.uploadingWithProgress', {
                progress: item.progress,
                speed: speedText,
                remaining: remainingText,
              })
            }
          }
        }
      )

      item.status = 'completed'
      item.progress = 100
      item.statusMessage = $t('upload.upload.smartUpload.status.completed')

      let resultData = result
      if (result.data) {
        resultData = result.data
      }
      item.result = resultData
    } catch (error: unknown) {
      item.status = 'failed'
      item.error = (error as any).message || $t('upload.upload.smartUpload.errors.uploadFailed')
      item.statusMessage = $t('upload.upload.smartUpload.status.failed', { error: item.error })
      console.error(`❌ [SmartUpload] Regular upload failed: ${item.file.name}`, error)
    }
  }

  const maxConcurrentUploads = computed(() => uploadConfig.getClientMaxConcurrentUploads())

  const startUpload = async () => {
    if (!uploadConfig.isConfigLoaded.value) {
      await new Promise((resolve) => {
        const unwatch = watch(
          () => uploadConfig.isConfigLoaded.value,
          (loaded) => {
            if (loaded) {
              unwatch()
              resolve(true)
            }
          }
        )
      })
    }

    scheduleNextUploads()
  }

  const scheduleNextUploads = () => {
    const currentRunning = runningUploads.value.size
    const maxConcurrent = maxConcurrentUploads.value

    if (!uploadConfig.isConfigLoaded.value) {
      return
    }

    if (currentRunning >= maxConcurrent) {
      return
    }

    const pendingFiles = allUploads.value.filter((item) => item.status === 'pending' && !runningUploads.value.has(item.id))

    const canStart = Math.min(pendingFiles.length, maxConcurrent - currentRunning)

    if (canStart === 0) {
      return
    }
    for (let i = 0; i < canStart; i++) {
      const file = pendingFiles[i]
      startSingleFileUpload(file)
    }
  }

  const startSingleFileUpload = async (file: Event) => {
    runningUploads.value.add(file.id)
    try {
      if (file.type === 'chunked') {
        await chunkedUpload.value.startUpload(file.id)
      } else {
        const regularItem = regularUploadQueue.value.find((item) => item.id === file.id)
        if (regularItem) {
          await uploadRegularFile(regularItem)
        }
      }
    } catch {
    } finally {
      runningUploads.value.delete(file.id)
      scheduleNextUploads()
    }
  }

  const pauseUpload = (itemId: string) => {
    if (chunkedUpload.value) {
      const chunkedItem = chunkedUpload.value.uploadQueue.find((item: any) => item.id === itemId)
      if (chunkedItem) {
        chunkedUpload.value.pauseUpload(itemId)
        return
      }
    }

    const regularItem = regularUploadQueue.value.find((item) => item.id === itemId)
    if (regularItem && regularItem.status === 'uploading') {
      regularItem.status = 'failed'
      regularItem.error = $t('upload.upload.smartUpload.status.userPaused')
      regularItem.statusMessage = $t('upload.upload.smartUpload.status.paused')
    }
  }

  const resumeUpload = async (itemId: string) => {
    if (chunkedUpload.value) {
      const chunkedItem = chunkedUpload.value.uploadQueue.find((item: any) => item.id === itemId)
      if (chunkedItem) {
        await chunkedUpload.value.resumeUpload(itemId)
        return
      }
    }

    const regularItem = regularUploadQueue.value.find((item) => item.id === itemId)
    if (regularItem && regularItem.status === 'failed') {
      await uploadRegularFile(regularItem)
    }
  }

  const cancelUpload = (itemId?: string) => {
    if (itemId) {
      if (chunkedUpload.value) {
        const chunkedItem = chunkedUpload.value.uploadQueue.find((item: any) => item.id === itemId)
        if (chunkedItem) {
          chunkedUpload.value.cancelUpload(itemId)
          return
        }
      }

      const regularItem = regularUploadQueue.value.find((item) => item.id === itemId)
      if (regularItem) {
        regularItem.status = 'failed'
        regularItem.error = $t('upload.upload.smartUpload.status.userCancelled')
        regularItem.statusMessage = $t('upload.upload.smartUpload.status.cancelled')
      }
    } else {
      if (chunkedUpload.value) {
        chunkedUpload.value.uploadQueue.forEach((item: any) => {
          if (item.status === 'uploading') {
            chunkedUpload.value.cancelUpload(item.id)
          }
        })
      }

      regularUploadQueue.value.forEach((item) => {
        if (item.status === 'uploading') {
          item.status = 'failed'
          item.error = $t('upload.upload.smartUpload.status.userCancelled')
          item.statusMessage = $t('upload.upload.smartUpload.status.cancelled')
        }
      })
    }
  }

  const retryUpload = async (itemId: string) => {
    if (chunkedUpload.value) {
      const chunkedItem = chunkedUpload.value.uploadQueue.find((item: any) => item.id === itemId)
      if (chunkedItem) {
        await chunkedUpload.value.retryUpload(itemId)
        return
      }
    }

    const regularItem = regularUploadQueue.value.find((item) => item.id === itemId)
    if (regularItem && regularItem.status === 'failed') {
      regularItem.status = 'pending'
      regularItem.progress = 0
      regularItem.error = undefined
      regularItem.statusMessage = $t('upload.upload.smartUpload.status.retrying')
      await uploadRegularFile(regularItem)
    }
  }

  const removeUploadItem = (itemId: string) => {
    if (chunkedUpload.value) {
      const chunkedItem = chunkedUpload.value.uploadQueue.find((item: any) => item.id === itemId)
      if (chunkedItem) {
        chunkedUpload.value.removeUploadItem(itemId)
        return
      }
    }

    const index = regularUploadQueue.value.findIndex((item) => item.id === itemId)
    if (index > -1) {
      regularUploadQueue.value.splice(index, 1)
    }
  }

  const clearQueue = () => {
    if (chunkedUpload.value) {
      chunkedUpload.value.clearQueue()
    }
    regularUploadQueue.value = []

    createdUrls.forEach((url) => {
      URL.revokeObjectURL(url)
    })
    createdUrls.clear()
  }

  const clearAllSessions = () => {
    if (chunkedUpload.value) {
      chunkedUpload.value.clearAllSessions()
    }
  }

  onUnmounted(() => {
    createdUrls.forEach((url) => {
      URL.revokeObjectURL(url)
    })
    createdUrls.clear()
  })

  const { formatFileSize } = uploadConfig

  const isValidImageFile = (file: File) => uploadConfig.isAllowedFileType(file.type)

  return {
    allUploads,
    globalProgress,
    globalSpeed,
    isUploading,
    runningUploads: computed(() => runningUploads.value), // 当前运行的上传
    maxConcurrentUploads, // 最大并发数

    regularUploadQueue,
    chunkedUploadQueue: computed(() => chunkedUpload.value?.uploadQueue || []),

    addFiles,
    startUpload,
    pauseUpload,
    resumeUpload,
    cancelUpload,
    retryUpload,
    removeUploadItem,
    clearQueue,
    clearAllSessions,
    formatFileSize,
    isValidImageFile,

    shouldUseChunkedUpload: uploadConfig.shouldUseChunkedUpload,
    uploadConfig: uploadConfig.uploadConfig,

    classifyFile,
  }
}
