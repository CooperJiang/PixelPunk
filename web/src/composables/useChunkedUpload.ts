import { computed, readonly, ref, triggerRef } from 'vue'
import { chunkedUploadAPI } from '@/api/upload/chunked'
import { calculateChunkMD5, calculateFileMD5, generateId } from '@/utils/file/fileHash'
import { type UploadSessionData, UploadStorageManager } from '@/utils/storage/uploadStorage'
import { InstantUploadUtil } from '@/utils/business/instantUpload'
import { getFileDimensions } from '@/utils/file/watermarkUtils'
import { useTexts } from '@/composables/useTexts'
import { useSettingsStore } from '@/store/settings'
import { UploadStatus, type UploadItem, type ChunkedUploadOptions } from './upload/types'

export function useChunkedUpload(options: ChunkedUploadOptions = {}) {
  const { $t } = useTexts()
  const {
    chunkSize = 2 * 1024 * 1024, // 2MB (可通过配置覆盖)
    concurrency = 3, // 并发数 (可通过配置覆盖)
    retryCount = 3, // 重试次数 (可通过配置覆盖)
    folderId,
    accessLevel = 'private',
    optimize = true,
    watermarkEnabled,
    watermarkConfig,
  } = options

  const uploadQueue = ref<UploadItem[]>([])

  const isUploading = ref(false)

  const globalProgress = computed(() => {
    if (uploadQueue.value.length === 0) {
      return 0
    }

    const totalProgress = uploadQueue.value.reduce((sum, item) => sum + item.progress, 0)
    return Math.round(totalProgress / uploadQueue.value.length)
  })

  const globalSpeed = computed(() => uploadQueue.value.reduce((sum, item) => sum + item.speed, 0))

  const addFiles = async (files: FileList | File[]) => {
    const fileArray = Array.from(files)
    const uploadItems: UploadItem[] = []

    for (const file of fileArray) {
      if (!isValidImageFile(file)) {
        console.warn(`⚠️ [ChunkedUpload] 跳过无效文件: ${file.name}`)
        continue
      }

      const totalChunks = Math.ceil(file.size / chunkSize)

      const uploadItem: UploadItem = {
        id: generateId(),
        file,
        status: UploadStatus.ANALYZING,
        statusMessage: $t('upload.smartUpload.status.analyzing'),
        progress: 0,
        totalChunks,
        uploadedChunks: 0,
        uploadedChunksList: [],
        failedChunks: new Set(),
        speed: 0,
        remainingTime: 0,
        chunkSize,
        folderId,
        accessLevel,
        optimize,
        watermarkEnabled,
        watermarkConfig,
        retryCount: 0,
        maxRetries: retryCount,
        totalBytesUploaded: 0,
        averageSpeed: 0,
      }

      uploadQueue.value.push(uploadItem)
      uploadItems.push(uploadItem)
    }

    const analysisPromises = uploadItems.map(async (uploadItem) => {
      try {
        const analysisTimeout = new Promise((_, reject) => {
          setTimeout(() => reject(new Error($t('upload.smartUpload.error.analysisTimeout'))), 30000) // 30秒超时
        })

        const onMD5Progress = (progress: number) => {
          uploadItem.statusMessage = $t('upload.smartUpload.status.analyzingProgress', { progress })
          triggerRef(uploadQueue)
        }

        const fileMD5 = (await Promise.race([calculateFileMD5(uploadItem.file, onMD5Progress), analysisTimeout])) as string

        uploadItem.fileMD5 = fileMD5

        if (uploadItem.file.type.startsWith('image/')) {
          try {
            const dimensions = await getFileDimensions(uploadItem.file)
            uploadItem.dimensions = dimensions
          } catch (error) {
            console.warn(`⚠️ [ChunkedUpload] 获取文件尺寸失败: ${uploadItem.file.name}`, error)
          }
        }

        const settingsStore = useSettingsStore()
        const instantUploadEnabled = settingsStore.rawSettings?.upload?.instant_upload_enabled ?? false

        if (instantUploadEnabled) {
          try {
            uploadItem.statusMessage = $t('upload.smartUpload.status.checkingInstant')
            triggerRef(uploadQueue)

            const instantResult = await InstantUploadUtil.attemptInstantUpload(
              uploadItem.file,
              {
                folder_id: uploadItem.folderId,
                access_level: uploadItem.accessLevel,
                optimize: uploadItem.optimize,
              },
              (progress) => {
                uploadItem.progress = Math.floor(progress)
                uploadItem.statusMessage = $t('upload.smartUpload.status.checkingInstantProgress', {
                  progress: uploadItem.progress,
                })
                triggerRef(uploadQueue)
              }
            )

            if (instantResult) {
              uploadItem.status = UploadStatus.COMPLETED
              uploadItem.progress = 100
              uploadItem.statusMessage = $t('upload.smartUpload.status.instantCompleted')
              uploadItem.result = instantResult.file_info || instantResult
              uploadItem.endTime = Date.now()

              triggerRef(uploadQueue)
              return
            }
          } catch (error) {
            // Continue with chunked upload
          }
        }

        const existingSession = UploadStorageManager.findExistingSession(fileMD5, uploadItem.file.size)

        if (existingSession) {
          uploadItem.sessionId = existingSession.sessionId
          uploadItem.status = UploadStatus.PAUSED
          uploadItem.statusMessage = $t('upload.smartUpload.status.foundResumable')
          uploadItem.progress = Math.round((existingSession.uploadedChunks.length / uploadItem.totalChunks) * 100)
          uploadItem.uploadedChunks = existingSession.uploadedChunks.length
          uploadItem.uploadedChunksList = existingSession.uploadedChunks
          uploadItem.folderId = existingSession.folderId || folderId
          uploadItem.accessLevel = existingSession.accessLevel || accessLevel
          uploadItem.optimize = existingSession.optimize ?? optimize
          triggerRef(uploadQueue)
        } else {
          uploadItem.status = UploadStatus.PENDING
          uploadItem.statusMessage = $t('upload.smartUpload.status.analyzingCompleted')
        }

        triggerRef(uploadQueue)
      } catch (error) {
        console.error(`❌ [ChunkedUpload] 文件分析失败: ${uploadItem.file.name}`, error)
        uploadItem.status = UploadStatus.FAILED
        uploadItem.error =
          error instanceof Error ? (error as any).message : $t('upload.smartUpload.error.analysisFailed')
        uploadItem.statusMessage = $t('upload.smartUpload.status.analyzingFailed', { error: uploadItem.error })

        triggerRef(uploadQueue)
      }
    })

    return Promise.allSettled(analysisPromises).then((results) => {
      const _successful = results.filter((r) => r.status === 'fulfilled').length
      const _failed = results.filter((r) => r.status === 'rejected').length
      return results
    })
  }

  const startUpload = async (itemId?: string) => {
    if (itemId) {
      const item = uploadQueue.value.find((item) => item.id === itemId)
      if (item && item.status === UploadStatus.PENDING) {
        await uploadSingleFile(item)
      }
    } else {
      isUploading.value = true

      const pendingItems = uploadQueue.value.filter(
        (item) => item.status === UploadStatus.PENDING || item.status === UploadStatus.PAUSED
      )

      for (const item of pendingItems) {
        if (!isUploading.value) {
          break
        } // 检查是否被停止
        await uploadSingleFile(item)
      }

      isUploading.value = false
    }
  }

  const pauseUpload = (itemId: string) => {
    const item = uploadQueue.value.find((item) => item.id === itemId)
    if (item && item.status === UploadStatus.UPLOADING) {
      item.status = UploadStatus.PAUSED
    }
  }

  const resumeUpload = async (itemId: string) => {
    const item = uploadQueue.value.find((item) => item.id === itemId)
    if (item && item.status === UploadStatus.PAUSED) {
      await uploadSingleFile(item)
    }
  }

  const cancelUpload = async (itemId: string) => {
    const item = uploadQueue.value.find((item) => item.id === itemId)
    if (!item) {
      return
    }

    if (item.sessionId) {
      try {
        await chunkedUploadAPI.cancel(item.sessionId)
      } catch (error) {
        console.error('取消上传失败:', error)
      }
    }

    item.status = UploadStatus.CANCELLED
  }

  const retryUpload = async (itemId: string) => {
    const item = uploadQueue.value.find((item) => item.id === itemId)
    if (!item || (item.status !== UploadStatus.FAILED && item.status !== UploadStatus.CANCELLED)) {
      return
    }

    item.retryCount = 0
    item.status = UploadStatus.RETRYING
    item.statusMessage = $t('upload.smartUpload.status.retryingUpload')
    item.error = undefined
    item.failedChunks.clear()

    await uploadSingleFile(item)
  }

  const removeUploadItem = (itemId: string) => {
    const index = uploadQueue.value.findIndex((item) => item.id === itemId)
    if (index > -1) {
      UploadStorageManager.removeSession(itemId)
      uploadQueue.value.splice(index, 1)
    }
  }

  const clearQueue = () => {
    uploadQueue.value.forEach((item) => {
      UploadStorageManager.removeSession(item.id)
    })
    uploadQueue.value = []
  }

  const getResumableUploads = () => UploadStorageManager.getResumableSessions()

  const cleanExpiredSessions = () => {
    UploadStorageManager.cleanExpiredSessions()
  }

  const clearAllSessions = () => {
    localStorage.removeItem('chunked_upload_sessions')
  }

  const uploadSingleFile = async (item: UploadItem) => {
    try {
      item.startTime = Date.now()
      let sessionData: UploadSessionData

      if (item.sessionId && (item.status === UploadStatus.PAUSED || item.status === UploadStatus.RETRYING)) {
        item.status = UploadStatus.UPLOADING
        item.statusMessage = $t('upload.smartUpload.status.resumingUpload')

        try {
          const serverStatus = await chunkedUploadAPI.getStatus(item.sessionId)

          item.uploadedChunks = serverStatus.uploaded_chunks
          item.progress = serverStatus.progress
          item.statusMessage = $t('upload.smartUpload.status.continuingUpload')

          item.uploadedChunksList = Array.from({ length: serverStatus.uploaded_chunks }, (_, i) => i)

          sessionData = createSessionData(item)
        } catch (error) {
          console.error('获取后端状态失败，重新开始上传:', error)
          item.statusMessage = $t('upload.smartUpload.status.restartingUpload')
          sessionData = await createNewSession(item)
        }
      } else {
        item.status = UploadStatus.PREPARING
        item.statusMessage = $t('upload.smartUpload.status.preparingSession')
        sessionData = await createNewSession(item)
      }

      UploadStorageManager.saveSession(sessionData)

      item.status = UploadStatus.UPLOADING
      item.statusMessage = $t('upload.smartUpload.status.uploading')
      await uploadChunksWithConcurrency(item)

      item.statusMessage = $t('upload.smartUpload.status.verifyingIntegrity')

      if (!item.sessionId) {
        throw new Error($t('composables.chunkedUpload.errors.sessionNotFound'))
      }

      try {
        const finalStatus = await chunkedUploadAPI.getStatus(item.sessionId)

        if (finalStatus.uploaded_chunks < item.totalChunks) {
          throw new Error(`Upload incomplete, uploaded ${finalStatus.uploaded_chunks}/${item.totalChunks} chunks`)
        }

        item.uploadedChunks = finalStatus.uploaded_chunks
        item.progress = finalStatus.progress
      } catch (statusError) {
        console.error('最终状态验证失败:', statusError)
        throw statusError
      }

      item.statusMessage = $t('upload.smartUpload.status.completingUpload')
      const completeResponse = await chunkedUploadAPI.complete(item.sessionId)

      let responseData = completeResponse
      if (completeResponse.data) {
        responseData = completeResponse.data
      }

      item.status = UploadStatus.COMPLETED
      item.progress = 100
      item.endTime = Date.now()
      item.statusMessage = $t('upload.smartUpload.status.completed')
      item.result = responseData

      if (item.startTime && item.endTime) {
        const totalTime = (item.endTime - item.startTime) / 1000
        item.averageSpeed = Math.round(item.file.size / totalTime / 1024) // KB/s
      }

      UploadStorageManager.removeSession(item.id)
    } catch (error: unknown) {
      item.status = UploadStatus.FAILED
      item.error = (error as any).message || $t('upload.smartUpload.error.uploadFailed')
      item.statusMessage = $t('upload.smartUpload.status.failed', { error: item.error })
      item.retryCount++

      console.error('上传失败:', error)

      const sessionData = createSessionData(item)
      UploadStorageManager.saveSession(sessionData)

      if (item.retryCount < item.maxRetries) {
        setTimeout(() => {
          if (item.status === UploadStatus.FAILED) {
            retryUpload(item.id)
          }
        }, 2000 * item.retryCount) // 递增等待时间
      }
    }
  }

  const createNewSession = async (item: UploadItem): Promise<UploadSessionData> => {
    if (!item.fileMD5) {
      item.fileMD5 = await calculateFileMD5(item.file)
    }

    const resolveValue = (value: any) => {
      if (value && typeof value === 'object' && 'value' in value) {
        return value.value
      }
      return value
    }

    let watermarkConfigForUpload: any = undefined
    if (resolveValue(item.watermarkEnabled) && item.watermarkConfig) {
      try {
        const config = { ...item.watermarkConfig }

        let finalScale = 1.0
        if (item.dimensions && config.generatedWidth) {
          const targetImageWidth = item.dimensions.width
          const watermarkOriginalWidth = config.generatedWidth
          const relativeSize = config.relativeSize || 0.15 // 默认占图片宽度的15%
          const calculatedScale = (targetImageWidth * relativeSize) / watermarkOriginalWidth
          finalScale = Math.max(0.1, Math.min(5.0, calculatedScale)) // 限制在合理范围内
        } else if (config.scale) {
          finalScale = config.scale
        }

        watermarkConfigForUpload = {
          enabled: true,
          type: config.type || 'image', // 前后端统一使用 'image' 或 'text'
          fileBase64: config.fileBase64 || '',
          position: config.position || 'bottom-right',
          offsetX: config.offsetX ?? 20,
          offsetY: config.offsetY ?? 20,
          offsetUnit: 'px',
          opacity: config.opacity ?? 0.7,
          scale: finalScale, // 水印在目标图片上的缩放比例
        }
      } catch (error) {
        console.warn('[ChunkedUpload] 水印配置准备失败:', error)
      }
    }

    const initResponse = await chunkedUploadAPI.init({
      file_name: item.file.name,
      file_size: item.file.size,
      file_md5: item.fileMD5,
      mime_type: item.file.type,
      chunk_size: item.chunkSize,
      folder_id: resolveValue(item.folderId) || undefined, // 确保null/empty转为undefined
      access_level: resolveValue(item.accessLevel),
      optimize: resolveValue(item.optimize),
      watermark_config: watermarkConfigForUpload,
    })

    item.sessionId = initResponse.session_id
    item.status = UploadStatus.UPLOADING

    return createSessionData(item)
  }

  const createSessionData = (item: UploadItem): UploadSessionData => {
    const resolveValue = (value: any) => {
      if (value && typeof value === 'object' && 'value' in value) {
        return value.value
      }
      return value
    }

    if (!item.sessionId || !item.fileMD5) {
      throw new Error($t('composables.chunkedUpload.errors.sessionDataMissing'))
    }

    return {
      id: item.id,
      sessionId: item.sessionId,
      fileName: item.file.name,
      fileSize: item.file.size,
      fileMD5: item.fileMD5,
      mimeType: item.file.type,
      chunkSize: item.chunkSize,
      totalChunks: item.totalChunks,
      uploadedChunks: item.uploadedChunksList,
      folderId: resolveValue(item.folderId),
      accessLevel: resolveValue(item.accessLevel),
      optimize: resolveValue(item.optimize),
      createdAt: Date.now(),
      lastActivity: Date.now(),
    }
  }

  const uploadChunksWithConcurrency = async (item: UploadItem) => {
    if (!item.sessionId) {
      throw new Error($t('composables.chunkedUpload.errors.sessionIdMissing'))
    }

    const startTime = Date.now()
    const chunksToUpload: number[] = []

    for (let chunkNumber = 0; chunkNumber < item.totalChunks; chunkNumber++) {
      if (!item.uploadedChunksList.includes(chunkNumber)) {
        chunksToUpload.push(chunkNumber)
      }
    }

    if (chunksToUpload.length === 0) {
      return
    }

    let currentConcurrency = concurrency
    let failedAttempts = 0
    const maxFailedAttempts = 3

    for (let i = 0; i < chunksToUpload.length; i += currentConcurrency) {
      if (item.status !== UploadStatus.UPLOADING) {
        break
      }

      const batch = chunksToUpload.slice(i, i + currentConcurrency)

      const batchStartTime = Date.now()
      const batchPromises = batch.map((chunkNumber) =>
        uploadSingleChunk(item, chunkNumber, startTime).catch((error) => {
          console.error(`分片 ${chunkNumber} 批次上传失败:`, error)
          failedAttempts++
          throw error
        })
      )

      try {
        await Promise.allSettled(batchPromises)

        const batchTime = Date.now() - batchStartTime
        const avgTimePerChunk = batchTime / batch.length

        if (avgTimePerChunk > 10000 && currentConcurrency > 1) {
          currentConcurrency = Math.max(1, currentConcurrency - 1)
        } else if (avgTimePerChunk < 3000 && currentConcurrency < concurrency) {
          currentConcurrency = Math.min(concurrency, currentConcurrency + 1)
        }

        failedAttempts = 0
      } catch (error) {
        console.error(`批次上传失败:`, error)

        if (failedAttempts >= maxFailedAttempts && currentConcurrency > 1) {
          currentConcurrency = Math.max(1, Math.floor(currentConcurrency / 2))
          failedAttempts = 0
        }
      }

      if (i + currentConcurrency < chunksToUpload.length) {
        await new Promise((resolve) => setTimeout(resolve, 100))
      }
    }
  }

  const uploadSingleChunk = async (item: UploadItem, chunkNumber: number, _startTime: number) => {
    if (item.status !== UploadStatus.UPLOADING) {
      return
    }

    let attempts = 0
    let lastError: Error | null = null

    while (attempts < retryCount) {
      try {
        const start = chunkNumber * item.chunkSize
        const end = Math.min(start + item.chunkSize, item.file.size)
        const chunk = item.file.slice(start, end)

        let chunkMD5: string
        if (attempts === 0) {
          chunkMD5 = await calculateChunkMD5(chunk)
        } else {
          chunkMD5 = await calculateChunkMD5(chunk)
        }

        const chunkStartTime = Date.now()

        const _networkCheckStart = Date.now()

        await chunkedUploadAPI.uploadChunk(
          {
            session_id: item.sessionId,
            chunk_number: chunkNumber,
            chunk_md5: chunkMD5,
          },
          chunk
        )

        item.uploadedChunks++
        item.uploadedChunksList.push(chunkNumber)
        item.uploadedChunksList.sort((a, b) => a - b)
        item.progress = Math.round((item.uploadedChunks / item.totalChunks) * 100)

        item.totalBytesUploaded += chunk.size

        const chunkEndTime = Date.now()
        const uploadTime = chunkEndTime - chunkStartTime
        if (uploadTime > 0) {
          const instantSpeed = chunk.size / (uploadTime / 1000) / 1024 // KB/s
          item.speed = item.speed === 0 ? instantSpeed : item.speed * 0.7 + instantSpeed * 0.3

          const remainingBytes = item.file.size - item.totalBytesUploaded
          item.remainingTime = item.speed > 0 ? Math.round(remainingBytes / (item.speed * 1024)) : 0
        }

        const speedText = item.speed > 1024 ? `${(item.speed / 1024).toFixed(1)} MB/s` : `${Math.round(item.speed)} KB/s`

        const remainingText =
          item.remainingTime > 60
            ? $t('upload.smartUpload.time.minutes', {
                minutes: Math.floor(item.remainingTime / 60),
                seconds: item.remainingTime % 60,
              })
            : $t('upload.smartUpload.time.seconds', { seconds: item.remainingTime })

        item.statusMessage = $t('upload.smartUpload.status.uploadingWithProgress', {
          progress: item.progress,
          speed: speedText,
          remaining: remainingText,
        })

        item.failedChunks.delete(chunkNumber)

        UploadStorageManager.addUploadedChunk(item.id, chunkNumber)

        return
      } catch (error: unknown) {
        attempts++
        lastError = error
        console.error(`分片 ${chunkNumber} 上传失败 (尝试 ${attempts}/${retryCount}):`, error)

        const isNetworkError =
          error.name === 'NetworkError' ||
          (error as any).message?.includes('fetch') ||
          (error as any).message?.includes('timeout') ||
          error.code === 'NETWORK_ERROR'

        const isServerError = error.status >= 500 && error.status < 600
        const isRateLimited = error.status === 429

        if (attempts >= retryCount) {
          item.failedChunks.add(chunkNumber)

          const errorMsg = error.message || $t('composables.chunkedUpload.errors.unknownError')
          let fullError = errorMsg

          if (isNetworkError) {
            fullError += $t('composables.chunkedUpload.errors.networkError')
          } else if (isServerError) {
            fullError += $t('composables.chunkedUpload.errors.serverError')
          } else if (isRateLimited) {
            fullError += $t('composables.chunkedUpload.errors.rateLimited')
          }

          item.statusMessage = $t('upload.smartUpload.status.chunkUploadFailed', {
            chunk: chunkNumber,
            error: fullError,
          })

          throw error
        }

        let retryDelay = 1000 * Math.pow(2, attempts - 1) // 基础指数退避

        if (isRateLimited) {
          retryDelay = Math.max(retryDelay, 5000) // 速率限制至少等待5秒
        } else if (isServerError) {
          retryDelay = Math.max(retryDelay, 2000) // 服务器错误至少等待2秒
        } else if (isNetworkError) {
          retryDelay = Math.max(retryDelay, 3000) // 网络错误至少等待3秒
        }

        retryDelay += Math.random() * 1000

        item.statusMessage = $t('upload.smartUpload.status.chunkRetrying', {
          chunk: chunkNumber,
          attempt: attempts,
          maxRetry: retryCount,
          delay: Math.round(retryDelay / 1000),
        })

        await new Promise((resolve) => setTimeout(resolve, retryDelay))
      }
    }

    throw lastError || new Error(`Chunk ${chunkNumber} upload failed after ${retryCount} retries`)
  }

  const isValidImageFile = (file: File): boolean => {
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp', 'image/bmp']
    return validTypes.includes(file.type)
  }

  return {
    uploadQueue: readonly(uploadQueue),
    isUploading: readonly(isUploading),
    globalProgress,
    globalSpeed,

    addFiles,
    startUpload,
    pauseUpload,
    resumeUpload,
    cancelUpload,
    retryUpload,
    removeUploadItem,
    clearQueue,
    getResumableUploads,
    cleanExpiredSessions,
    clearAllSessions,
  }
}
