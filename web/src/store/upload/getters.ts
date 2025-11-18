/**
 * 上传计算属性
 * 所有 computed 逻辑集中在此

 */
import { computed } from 'vue'
import { useAuthStore } from '@/store/auth'
import { useUploadConfig } from '@/composables/useUploadConfig'
import { regularQueue, chunkedUploadInstance, isChunkedReady, runningUploads, globalOptions, previewUrls } from './state'
import type { UploadItem } from './types'
import { isHEICFile } from '@/utils/file/heicConverter'

/* 生成 HEIC 占位图 */
function getHEICPlaceholder(): string {
  const svg = `<svg width="200" height="200" xmlns="http://www.w3.org/2000/svg">
    <defs>
      <linearGradient id="bg" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#1e1e1e"/>
        <stop offset="100%" style="stop-color:#2d2d2d"/>
      </linearGradient>
      <linearGradient id="icon" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#6366f1"/>
        <stop offset="100%" style="stop-color:#8b5cf6"/>
      </linearGradient>
    </defs>
    <rect width="200" height="200" fill="url(#bg)" rx="8"/>
    <g transform="translate(100, 100)">
      <rect x="-45" y="-35" width="90" height="70" rx="6" fill="none" stroke="url(#icon)" stroke-width="4"/>
      <circle cx="-20" cy="-15" r="10" fill="url(#icon)" opacity="0.8"/>
      <path d="M -45 35 L -15 -5 L 10 20 L 45 -10 L 45 35 Z" fill="url(#icon)" opacity="0.6"/>
      <path d="M -45 35 L -10 5 L 20 35 Z" fill="url(#icon)"/>
    </g>
  </svg>`

  return `data:image/svg+xml;base64,${btoa(unescape(encodeURIComponent(svg)))}`
}

export const allUploads = computed(() => {
  const regular = regularQueue.value.map((item) => {
    if (!item.preview && item.file.type.startsWith('image/')) {
      if (isHEICFile(item.file)) {
        item.preview = getHEICPlaceholder()
        item.previewLoaded = true
      } else {
        const url = URL.createObjectURL(item.file)
        previewUrls.value.add(url)
        item.preview = url
        item.previewLoaded = true
      }
    }
    return item
  })

  let chunked: UploadItem[] = []
  if (isChunkedReady.value && chunkedUploadInstance.value?.uploadQueue) {
    chunked = chunkedUploadInstance.value.uploadQueue.map((item: any) => {
      let previewUrl = item.preview
      let previewLoaded = item.previewLoaded
      if (!previewUrl && item.file.type.startsWith('image/')) {
        if (isHEICFile(item.file)) {
          previewUrl = getHEICPlaceholder()
          previewLoaded = true
        } else {
          previewUrl = URL.createObjectURL(item.file)
          previewUrls.value.add(previewUrl)
          previewLoaded = true
        }
      }

      return {
        ...item,
        type: 'chunked' as const,
        preview: previewUrl,
        previewLoaded,
      }
    })
  }

  return [...regular, ...chunked]
})

export const globalProgress = computed(() => {
  const total = allUploads.value.length
  if (total === 0) {
    return 0
  }

  const totalProgress = allUploads.value.reduce((sum, item) => sum + item.progress, 0)
  return Math.round(totalProgress / total)
})

export const globalSpeed = computed(() => allUploads.value.reduce((sum, item) => sum + (item.speed || 0), 0))

export const isUploading = computed(
  () =>
    runningUploads.value.size > 0 || allUploads.value.some((item) => item.status === 'uploading' || item.status === 'preparing')
)

export const statistics = computed(() => ({
  pending: allUploads.value.filter((f) => f.status === 'pending' || f.status === 'paused').length,
  uploading: allUploads.value.filter(
    (f) =>
      f.status === 'uploading' ||
      f.status === 'analyzing' ||
      f.status === 'preparing' ||
      f.status === 'retrying'
  ).length,
  completed: allUploads.value.filter((f) => f.status === 'completed').length,
  failed: allUploads.value.filter((f) => f.status === 'failed').length,
  total: allUploads.value.length,
}))

export const maxConcurrentUploads = computed(() => {
  const authStore = useAuthStore()
  const uploadConfig = useUploadConfig()

  if (!authStore.isLoggedIn) {
    return 1
  }
  return uploadConfig.getClientMaxConcurrentUploads()
})

export const getters = {
  allUploads,
  globalProgress,
  globalSpeed,
  isUploading,
  statistics,
  maxConcurrentUploads,
  globalOptions: computed(() => globalOptions.value),
  runningUploads: computed(() => runningUploads.value),
}
