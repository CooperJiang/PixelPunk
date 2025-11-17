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
      status: 'analyzing',
      progress: 0,
      speed: 0,
      remainingTime: 0,
      statusMessage: translations.value?.queue.preparingAnalysis || 'Preparing analysis...',
      type: 'regular',
    })

    validFiles.push({ file, id })
  }

  for (const { file, id } of validFiles) {
    processFileAsync(file, id, uploadConfig, toast).catch(() => {})
  }
}

async function processFileAsync(file: File, itemId: string, uploadConfig: any, toast: any) {
  const item = regularQueue.value.find((i) => i.id === itemId)
  if (!item) {
    return
  }

  try {
    item.statusMessage = translations.value?.queue.analyzing || 'Analyzing file...'
    const uploadItem = await createUploadItemWithDimensions(file, itemId, 'regular')

    Object.assign(item, uploadItem)
    item.status = 'analyzing'
    item.statusMessage = translations.value?.queue.checkingInstant || 'Checking instant upload...'

    const instantResult = await InstantUploadUtil.attemptInstantUpload(
      file,
      {
        folder_id: globalOptions.value.folderId || undefined,
        access_level: globalOptions.value.accessLevel,
        optimize: globalOptions.value.optimize,
      },
      (progress) => {
        item.progress = progress
        if (progress < 50) {
          item.statusMessage = translations.value?.queue.calculatingMD5?.replace('{progress}', String(progress)) || `Calculating MD5... ${progress}%`
        } else if (progress < 75) {
          item.statusMessage = translations.value?.queue.checkingServerDuplicate || 'Checking duplicates...'
        } else {
          item.statusMessage = translations.value?.queue.preparingUpload || 'Instant uploading...'
        }
      }
    )

    if (instantResult) {
      const resultData = instantResult.file_info || instantResult.data || instantResult
      item.status = 'completed'
      item.progress = 100
      item.statusMessage = translations.value?.queue.instantComplete || 'Instant upload completed'
      item.result = resultData
      toast.success($t('upload.uploadProgress.toast.instantUploadSuccess', { name: file.name }))

      if (globalOptions.value.autoRemove) {
        setTimeout(() => {
          removeUploadItem(itemId)
        }, 500)
      }
    } else {
      item.status = 'pending'
      item.progress = 0
      item.statusMessage = translations.value?.queue.waitingUpload || 'Waiting for upload'

      if (classifyFile(file)) {
        const index = regularQueue.value.findIndex((i) => i.id === itemId)
        if (index !== -1) {
          regularQueue.value.splice(index, 1)
        }
        if (chunkedUploadInstance.value) {
          await chunkedUploadInstance.value.addFiles([file])
        }
      }
    }
  } catch (error) {
    item.status = 'pending'
    item.progress = 0
    item.statusMessage = translations.value?.queue.waitingUpload || 'Waiting for upload'

    if (classifyFile(file)) {
      const index = regularQueue.value.findIndex((i) => i.id === itemId)
      if (index !== -1) {
        regularQueue.value.splice(index, 1)
      }
      if (chunkedUploadInstance.value) {
        await chunkedUploadInstance.value.addFiles([file])
      }
    }
  }
}

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
