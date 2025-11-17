/**
 * 上传状态定义
 * 集中管理所有响应式状态

 */
import { ref, computed } from 'vue'
import type { GlobalUploadOptions, UploadItem } from './types'
import type { TranslationFunction } from '@/composables/useTexts'

/* 全局上传选项 */
export const globalOptions = ref<GlobalUploadOptions>({
  folderId: null,
  accessLevel: 'public',
  optimize: true,
  autoRemove: false,
})

/* 普通上传队列 */
export const regularQueue = ref<UploadItem[]>([])

/* 分片上传实例（单例） */
export const chunkedUploadInstance = ref<any>(null)
export const isChunkedReady = ref(false)

/* 全局并发控制 */
export const runningUploads = ref(new Set<string>())

/* 存储创建的预览URL，用于清理 */
export const previewUrls = ref(new Set<string>())

/* 翻译函数引用 */
export const $t = ref<TranslationFunction>()

/* 翻译文本 */
export const translations = computed(() => {
  if (!$t.value) return null
  const t = $t.value
  return {
    control: {
      paused: t('store.upload.control.paused'),
      preparingResume: t('store.upload.control.preparingResume'),
      userCancelled: t('store.upload.control.userCancelled'),
      cancelled: t('store.upload.control.cancelled'),
      preparingRetry: t('store.upload.control.preparingRetry'),
    },
    queue: {
      preparingAnalysis: t('store.upload.queue.preparingAnalysis'),
      analyzing: t('store.upload.queue.analyzing'),
      checkingInstant: t('store.upload.queue.checkingInstant'),
      calculatingMD5: t('store.upload.queue.calculatingMD5'),
      checkingServerDuplicate: t('store.upload.queue.checkingServerDuplicate'),
      preparingUpload: t('store.upload.queue.preparingUpload'),
      instantComplete: t('store.upload.queue.instantComplete'),
      analysisComplete: t('store.upload.queue.analysisComplete'),
      waitingUpload: t('store.upload.queue.waitingUpload'),
    },
    upload: {
      analyzingFile: t('store.upload.upload.analyzingFile'),
      analyzingProgress: t('store.upload.upload.analyzingProgress'),
      instantProgress: t('store.upload.upload.instantProgress'),
      instantComplete: t('store.upload.upload.instantComplete'),
      uploading: t('store.upload.upload.uploading'),
      mergingChunks: t('store.upload.upload.mergingChunks'),
      uploadComplete: t('store.upload.upload.uploadComplete'),
      uploadFailed: t('store.upload.upload.uploadFailed'),
      instantFailed: t('store.upload.upload.instantFailed'),
      preparingRetry: t('store.upload.upload.preparingRetry'),
      retrying: t('store.upload.upload.retrying'),
      postProcessing: t('store.upload.upload.postProcessing'),
      analysisError: t('store.upload.upload.analysisError'),
    },
  }
})
