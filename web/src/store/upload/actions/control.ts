/**
 * 上传控制操作
 * 负责启动、暂停、继续、取消、重试上传

 */
import { chunkedUploadInstance, runningUploads, translations } from '../state'
import { allUploads } from '../getters'
import { scheduleNextUploads } from '../utils/scheduler'
import { uploadRegularFile } from './upload'

export function startUpload() {
  scheduleNextUploads(uploadRegularFile)
}

export function pauseUpload(itemId?: string) {
  if (itemId) {
    const item = allUploads.value.find((i) => i.id === itemId)
    if (item) {
      if (item.type === 'chunked') {
        chunkedUploadInstance.value?.pauseUpload(itemId)
      } else {
        item.status = 'paused'
        item.statusMessage = translations.value?.control.paused || 'Paused'
      }
    }
  } else {
    allUploads.value.forEach((item) => {
      if (item.status === 'uploading') {
        pauseUpload(item.id)
      }
    })
  }
}

export async function resumeUpload(itemId: string) {
  const item = allUploads.value.find((i) => i.id === itemId)
  if (!item) {
    return
  }

  if (item.type === 'chunked') {
    await chunkedUploadInstance.value?.resumeUpload(itemId)
  } else {
    item.status = 'pending'
    item.statusMessage = translations.value?.control.preparingResume || 'Preparing to resume'
    scheduleNextUploads(uploadRegularFile)
  }
}

export function cancelUpload(itemId?: string) {
  if (itemId) {
    const item = allUploads.value.find((i) => i.id === itemId)
    if (item) {
      if (item.type === 'chunked') {
        chunkedUploadInstance.value?.cancelUpload(itemId)
      } else {
        item.status = 'failed'
        item.error = translations.value?.control.userCancelled || 'Cancelled by user'
        item.statusMessage = translations.value?.control.cancelled || 'Cancelled'
      }
      runningUploads.value.delete(itemId)
    }
  } else {
    const itemsToCancel = allUploads.value.filter(
      (item) =>
        item.status === 'uploading' ||
        item.status === 'pending' ||
        item.status === 'analyzing' ||
        item.status === 'preparing' ||
        item.status === 'retrying'
    )

    itemsToCancel.forEach((item) => {
      cancelUpload(item.id)
    })
  }
}

export async function retryUpload(itemId: string) {
  const item = allUploads.value.find((i) => i.id === itemId)
  if (!item || item.status !== 'failed') {
    return
  }

  if (item.type === 'chunked') {
    await chunkedUploadInstance.value?.retryUpload(itemId)
  } else {
    item.uploadSessionId = `${Date.now()}-${Math.random().toString(36).substring(2, 9)}`
    item.status = 'pending'
    item.progress = 0
    item.error = undefined
    item.statusMessage = translations.value?.control.preparingRetry || 'Preparing to retry'
    scheduleNextUploads(uploadRegularFile)
  }
}
