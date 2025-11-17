import { type App, createApp } from 'vue'
import DownloaderFloat from '@/components/DownloaderFloat/index.vue'
import type { DownloadOptions } from '@/components/DownloaderFloat/types'
import { StorageUtil } from '@/utils/storage'

interface DownloaderInstance {
  app: App
  container: HTMLDivElement
  close: () => void
}

interface QueueItem {
  options: DownloadOptions
  resolve: (success: boolean) => void
}

/* 存储当前活动的下载器实例 */
const activeDownloaders = new Set<DownloaderInstance>()

/* 下载队列管理 */
let currentDownloader: DownloaderInstance | null = null
const downloadQueue: QueueItem[] = []
let isProcessingQueue = false

/* 队列变化回调 */
const queueChangeCallbacks: ((size: number) => void)[] = []

/* 注册队列变化监听器 */
export function onQueueChange(callback: (size: number) => void) {
  queueChangeCallbacks.push(callback)
  callback(downloadQueue.length)

  return () => {
    const index = queueChangeCallbacks.indexOf(callback)
    if (index > -1) {
      queueChangeCallbacks.splice(index, 1)
    }
  }
}

function notifyQueueChange() {
  queueChangeCallbacks.forEach((callback) => {
    try {
      callback(downloadQueue.length)
    } catch (error) {
      console.error('Queue change callback error:', error)
    }
  })
}

function createDownloaderInstance(options: DownloadOptions, queueSize: number): Promise<boolean> {
  return new Promise((resolve) => {
    const container = document.createElement('div')
    document.body.appendChild(container)

    const app = createApp(DownloaderFloat, {
      options: { ...options, queueSize },
      onComplete: (success: boolean) => {
        resolve(success)
        processNextInQueue()
      },
      onClose: () => {
        setTimeout(() => destroyDownloader(instance), 300)
        currentDownloader = null
      },
    })

    app.mount(container)

    const instance: DownloaderInstance = {
      app,
      container,
      close: () => destroyDownloader(instance),
    }

    activeDownloaders.add(instance)
    currentDownloader = instance
  })
}

function processNextInQueue() {
  if (downloadQueue.length === 0) {
    isProcessingQueue = false
    notifyQueueChange()
    return
  }

  const nextItem = downloadQueue.shift()
  if (!nextItem) {
    isProcessingQueue = false
    return
  }

  notifyQueueChange()
  createDownloaderInstance(nextItem.options, downloadQueue.length).then(nextItem.resolve)
}

export function createDownloader(options: DownloadOptions): Promise<boolean> {
  return new Promise((resolve) => {
    if (!currentDownloader && !isProcessingQueue) {
      isProcessingQueue = true
      createDownloaderInstance(options, downloadQueue.length).then(resolve)
    } else {
      downloadQueue.push({ options, resolve })

      showQueueToast(downloadQueue.length)

      notifyQueueChange()
    }
  })
}

function showQueueToast(queueLength: number) {
  Promise.all([import('@/components/Toast/useToast'), import('@/composables/useTexts')])
    .then(([{ useToast }, { useTexts }]) => {
      const toast = useToast()
      const { $t } = useTexts()
      toast.info($t('components.downloader.queueAdded').replace('{count}', queueLength.toString()))
    })
    .catch(() => {})
}

function destroyDownloader(instance: DownloaderInstance) {
  try {
    instance.app.unmount()
    instance.container.parentNode?.removeChild(instance.container)
    activeDownloaders.delete(instance)
  } catch (error) {
    console.error('Downloader instance destroy failed:', error)
  }
}

export function downloadFileQuick(fileId: string, fileName?: string): Promise<boolean> {
  return createDownloader({
    fileId,
    fileName,
    quality: 'original',
    position: { side: 'right', y: '50%' },
    autoClose: true,
    closeDelay: 3000,
  })
}

export function downloadSharedFileQuick(
  fileId: string,
  shareKey: string,
  accessToken?: string,
  fileName?: string
): Promise<boolean> {
  return createDownloader({
    fileId,
    shareKey,
    accessToken,
    fileName,
    quality: 'original',
    position: { side: 'right', y: '50%' },
    autoClose: true,
    closeDelay: 3000,
  })
}

export function downloadFileWithPosition(
  fileId: string,
  position: DownloadOptions['position'],
  options?: Partial<Omit<DownloadOptions, 'fileId' | 'position'>>
): Promise<boolean> {
  return createDownloader({
    fileId,
    position,
    quality: 'original',
    autoClose: true,
    closeDelay: 3000,
    ...options,
  })
}

export function closeAllDownloaders() {
  activeDownloaders.forEach((instance) => instance.close())
  activeDownloaders.clear()
  downloadQueue.length = 0
  currentDownloader = null
  isProcessingQueue = false
}

export function getActiveDownloadersCount(): number {
  return activeDownloaders.size
}

export function getQueueSize(): number {
  return downloadQueue.length
}

export function clearDownloaderPosition(): void {
  StorageUtil.remove('cyber-downloader-position')
}

export function getDownloaderPosition(): { x: number; y: number } | null {
  return StorageUtil.get<{ x: number; y: number }>('cyber-downloader-position')
}

export const Downloader = {
  create: createDownloader,
  quick: downloadFileQuick,
  quickShared: downloadSharedFileQuick,
  withPosition: downloadFileWithPosition,
  closeAll: closeAllDownloaders,
  getActiveCount: getActiveDownloadersCount,
  getQueueSize,
  clearPosition: clearDownloaderPosition,
  getPosition: getDownloaderPosition,
}

if (typeof window !== 'undefined') {
  ;(window as any).CyberDownloader = Downloader
}

export default Downloader
