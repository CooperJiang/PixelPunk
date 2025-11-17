/**
 * 负责管理并发上传任务的调度

 */
import { runningUploads, chunkedUploadInstance } from '../state'
import { allUploads, maxConcurrentUploads } from '../getters'
import type { UploadItem } from '../types'

/* 调度下一个上传任务 */
export function scheduleNextUploads(uploadRegularFile: (item: UploadItem) => Promise<void>) {
  const currentRunning = runningUploads.value.size
  const maxConcurrent = maxConcurrentUploads.value

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
      await uploadRegularFile(file)
    }
  } finally {
    runningUploads.value.delete(file.id)
    scheduleNextUploads(uploadRegularFile)
  }
}
