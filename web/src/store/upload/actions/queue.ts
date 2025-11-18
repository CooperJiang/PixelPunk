/**
 * 队列管理操作
 * 负责文件添加、秒传等队列相关逻辑

 */
import { useToast } from '@/components/Toast/useToast'
import { useUploadConfig } from '@/composables/useUploadConfig'
import { useTexts } from '@/composables/useTexts'
import { generateId } from '@/utils/file/fileHash'
import { InstantUploadUtil } from '@/utils/business/instantUpload'
import { regularQueue, chunkedUploadInstance, isChunkedReady, globalOptions, previewUrls, translations } from '../state'
import { initChunkedUpload, classifyFile } from './init'
import { createUploadItemWithDimensions } from '../utils/watermark'

export async function addFiles(files: FileList | File[]) {
  const uploadConfig = useUploadConfig()
  const toast = useToast()
  const { $t } = useTexts()

  if (!uploadConfig.isConfigLoaded.value) {
    await uploadConfig.loadConfig()
  }

  if (!isChunkedReady.value) {
    initChunkedUpload()
  }

  const fileArray = Array.from(files)
  const validFiles: Array<{ file: File; id: string }> = []

  for (const file of fileArray) {
    if (!uploadConfig.isAllowedFileType(file.type)) {
      toast.error($t('upload.uploadProgress.toast.fileTypeNotSupported', { name: file.name }))
      continue
    }

    if (!uploadConfig.isAllowedFileSize(file.size)) {
      toast.error(
        $t('upload.uploadProgress.toast.fileSizeExceeded', {
          name: file.name,
          size: uploadConfig.formatFileSize(file.size),
        })
      )
      continue
    }

    const id = generateId()
    regularQueue.value.push({
      id,
      file,
      status: 'pending',
      progress: 0,
      speed: 0,
      remainingTime: 0,
      statusMessage: translations.value?.queue.waitingUpload || 'Ready to upload',
      type: 'regular',
    })

    validFiles.push({ file, id })
  }

  // 异步创建预览图和读取图片尺寸（并发控制5个）
  const MAX_CONCURRENT_PREVIEW = 5
  const previewQueue = [...validFiles]
  const processingPreviews: Promise<void>[] = []

  while (previewQueue.length > 0 || processingPreviews.length > 0) {
    while (processingPreviews.length < MAX_CONCURRENT_PREVIEW && previewQueue.length > 0) {
      const { file, id } = previewQueue.shift()!

      const promise = createPreviewForFile(file, id).finally(() => {
        const index = processingPreviews.indexOf(promise)
        if (index > -1) processingPreviews.splice(index, 1)
      })

      processingPreviews.push(promise)
    }

    if (processingPreviews.length > 0) {
      await Promise.race(processingPreviews)
    }
  }
}

async function createPreviewForFile(file: File, itemId: string) {
  const item = regularQueue.value.find((i) => i.id === itemId)
  if (!item) return

  try {
    const uploadItem = await createUploadItemWithDimensions(file, itemId, 'regular')
    Object.assign(item, uploadItem)
  } catch (error) {
    // 预览创建失败不影响上传
  }
}

// processFileAsync 函数已移除，秒传检测移至 startUpload 阶段

export function removeUploadItem(itemId: string) {
  if (chunkedUploadInstance.value) {
    const chunkedItem = chunkedUploadInstance.value.uploadQueue.find((item: any) => item.id === itemId)
    if (chunkedItem) {
      chunkedUploadInstance.value.removeUploadItem(itemId)
      return
    }
  }

  const index = regularQueue.value.findIndex((item) => item.id === itemId)
  if (index > -1) {
    const item = regularQueue.value[index]
    if (item.preview && previewUrls.value.has(item.preview)) {
      URL.revokeObjectURL(item.preview)
      previewUrls.value.delete(item.preview)
    }
    regularQueue.value.splice(index, 1)
  }
}
