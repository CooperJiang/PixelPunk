/**
 * 负责管理并发上传任务的调度
 * 优化版：增加秒传检测，并发控制为5个
 */
import { runningUploads, chunkedUploadInstance, globalOptions, translations } from '../state'
import { allUploads, maxConcurrentUploads } from '../getters'
import { InstantUploadUtil } from '@/utils/business/instantUpload'
import { removeUploadItem } from '../actions/queue'
import { useSettingsStore } from '@/store/settings'
import type { UploadItem } from '../types'

// 秒传检测并发控制
const MAX_CONCURRENT_INSTANT_CHECK = 5

/* 调度下一个上传任务 */
export function scheduleNextUploads(uploadRegularFile: (item: UploadItem) => Promise<void>) {
  const currentRunning = runningUploads.value.size
  const maxConcurrent = Math.max(maxConcurrentUploads.value, MAX_CONCURRENT_INSTANT_CHECK)

  if (currentRunning >= maxConcurrent) {
    return
  }

  const pendingFiles = allUploads.value.filter((item) => item.status === 'pending' && !runningUploads.value.has(item.id))

  const canStart = Math.min(pendingFiles.length, maxConcurrent - currentRunning)

  if (canStart === 0) {
    return
  }

  for (let i = 0; i < canStart; i++) {
    startSingleFileUpload(pendingFiles[i], uploadRegularFile)
  }
}

export async function startSingleFileUpload(file: UploadItem, uploadRegularFile: (item: UploadItem) => Promise<void>) {
  runningUploads.value.add(file.id)

  try {
    if (file.type === 'chunked') {
      await chunkedUploadInstance.value.startUpload(file.id)
    } else {
      // 先尝试秒传检测
      const instantSuccess = await attemptInstantUploadForItem(file)

      // 如果秒传失败，进行正常上传
      if (!instantSuccess) {
        await uploadRegularFile(file)
      }
    }
  } finally {
    runningUploads.value.delete(file.id)
    scheduleNextUploads(uploadRegularFile)
  }
}

/**
 * 尝试对单个文件进行秒传
 * @returns true: 秒传成功, false: 秒传失败需要正常上传
 */
async function attemptInstantUploadForItem(item: UploadItem): Promise<boolean> {
  const settingsStore = useSettingsStore()
  const instantUploadEnabled = settingsStore.rawSettings?.upload?.instant_upload_enabled ?? false

  if (!instantUploadEnabled) {
    return false
  }

  try {
    item.status = 'analyzing'
    item.progress = 0
    item.statusMessage = translations.value?.queue.checkingInstant || 'Checking instant upload...'

    const instantResult = await InstantUploadUtil.attemptInstantUpload(
      item.file,
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
      // 秒传成功
      const resultData = instantResult.file_info || instantResult.data || instantResult
      item.status = 'completed'
      item.progress = 100
      item.statusMessage = translations.value?.queue.instantComplete || 'Instant upload completed'
      item.result = resultData

      if (globalOptions.value.autoRemove) {
        setTimeout(() => {
          removeUploadItem(item.id)
        }, 500)
      }

      return true
    }

    return false
  } catch (error) {
    return false
  }
}
