import SparkMD5 from 'spark-md5'
import { uploadWorkerManager } from '@/workers/uploadWorkerManager'
import { useTextThemeStore } from '@/store/textTheme'

/* 判断是否应该使用Worker（大于2MB的文件使用Worker，避免主线程阻塞） */
const shouldUseWorker = (fileSize: number): boolean => fileSize > 2 * 1024 * 1024 // 2MB

/**
 * 格式化翻译文本，替换参数
 */
function formatText(template: string, params: Record<string, string | number>): string {
  let result = template
  Object.keys(params).forEach((key) => {
    result = result.replace(new RegExp(`\\{${key}\\}`, 'g'), String(params[key]))
  })
  return result
}

export async function calculateFileMD5(file: File, onProgress?: (progress: number) => void): Promise<string> {
  if (shouldUseWorker(file.size)) {
    try {
      const hash = await uploadWorkerManager.calculateMD5(file, onProgress)
      return hash
    } catch {}
  }

  return calculateFileMD5InMainThread(file, onProgress)
}

function calculateFileMD5InMainThread(file: File, onProgress?: (progress: number) => void): Promise<string> {
  return new Promise((resolve, reject) => {
    const store = useTextThemeStore()
    const spark = new SparkMD5.ArrayBuffer()
    const fileReader = new FileReader()
    const chunkSize = 2097152 // 2MB 分块读取
    let currentChunk = 0
    const chunks = Math.ceil(file.size / chunkSize)
    let retryCount = 0
    const maxRetries = 3

    fileReader.onload = (e) => {
      try {
        if (e.target?.result) {
          spark.append(e.target.result as ArrayBuffer)
          currentChunk++
          retryCount = 0 // 重置重试计数

          if (onProgress) {
            const progress = Math.round((currentChunk / chunks) * 100)
            onProgress(progress)
          }

          if (currentChunk < chunks) {
            loadNext()
          } else {
            resolve(spark.end())
          }
        } else {
          throw new Error(store.getText('utils.file.fileHash.errors.emptyResult'))
        }
      } catch (error) {
        reject(new Error(formatText(store.getText('utils.file.fileHash.errors.md5Failed'), { error: String(error) })))
      }
    }

    fileReader.onerror = (_e) => {
      if (retryCount < maxRetries) {
        retryCount++
        setTimeout(() => loadNext(), 1000 * retryCount) // 递增延迟重试
      } else {
        reject(new Error(formatText(store.getText('utils.file.fileHash.errors.readFailedAfterRetries'), { maxRetries })))
      }
    }

    fileReader.onabort = () => {
      reject(new Error(store.getText('utils.file.fileHash.errors.readAborted')))
    }

    function loadNext() {
      try {
        const start = currentChunk * chunkSize
        const end = Math.min(start + chunkSize, file.size)

        if (start >= file.size) {
          reject(new Error(store.getText('utils.file.fileHash.errors.chunkIndexOutOfRange')))
          return
        }

        const slice = file.slice(start, end)
        fileReader.readAsArrayBuffer(slice)
      } catch (error) {
        reject(new Error(formatText(store.getText('utils.file.fileHash.errors.chunkCreationFailed'), { error: String(error) })))
      }
    }

    loadNext()
  })
}

export async function calculateChunkMD5(chunk: Blob): Promise<string> {
  if (chunk.size > 1024 * 1024) {
    try {
      return await uploadWorkerManager.calculateChunkMD5(chunk)
    } catch {}
  }

  return new Promise((resolve, reject) => {
    const store = useTextThemeStore()
    const spark = new SparkMD5.ArrayBuffer()
    const fileReader = new FileReader()

    fileReader.onload = (e) => {
      if (e.target?.result) {
        spark.append(e.target.result as ArrayBuffer)
        resolve(spark.end())
      }
    }

    fileReader.onerror = (_e) => {
      reject(new Error(store.getText('utils.file.fileHash.errors.chunkReadFailed')))
    }

    fileReader.readAsArrayBuffer(chunk)
  })
}

export function generateId(): string {
  return Math.random().toString(36).substr(2, 9) + Date.now().toString(36)
}
