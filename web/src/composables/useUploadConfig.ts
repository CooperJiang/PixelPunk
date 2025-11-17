import { computed, ref } from 'vue'
import { get } from '@/utils/network/http'
import { useTexts } from '@/composables/useTexts'

/* 上传配置接口 */
export interface UploadConfig {
  maxFileSize: number
  allowedTypes: string[]
  clientMaxConcurrentUploads: number // 新增：客户端最大并发上传数
  chunkedUpload: {
    enabled: boolean
    isAllowChunkUpload: boolean // 新增：是否允许分片上传（仅local存储允许）
    threshold: number
    chunkSize: number
    maxConcurrency: number // 只保留客户端需要的配置
    retryCount: number
  }
}

/* 能力接口返回 */
interface UploadCapabilities {
  supported_extensions: string[]
  supported_extensions_with_dot: string[]
  mime_map: Record<string, string>
  thumbnail: {
    rasterize_svg: boolean
    webp_conversion: boolean
    transparent_preserve: boolean
  }
}

/* 默认配置（与后端保持一致） */
const defaultConfig: UploadConfig = {
  maxFileSize: 100 * 1024 * 1024, // 100MB
  allowedTypes: [
    'image/jpeg',
    'image/jpg',
    'image/png',
    'image/gif',
    'image/webp',
    'image/bmp',
    'image/svg+xml',
    'image/x-icon',
    'image/vnd.microsoft.icon',
    'image/tiff',
    'image/x-tga',
  ],
  clientMaxConcurrentUploads: 3, // 默认客户端最大并发数
  chunkedUpload: {
    enabled: true,
    isAllowChunkUpload: false, // 默认不允许分片上传，需要服务端确认
    threshold: 10 * 1024 * 1024, // 10MB
    chunkSize: 2 * 1024 * 1024, // 2MB
    maxConcurrency: 4, // 分片并发数
    retryCount: 3,
  },
}

/* 全局配置状态 */
const uploadConfig = ref<UploadConfig>(defaultConfig)
const isConfigLoaded = ref(false)
const isLoading = ref(false) // 添加加载状态，防止重复请求
const capabilities = ref<UploadCapabilities | null>(null)

export function useUploadConfig() {
  const { $t } = useTexts()

  const loadConfig = async () => {
    if (isLoading.value) {
      return
    }
    if (isConfigLoaded.value) {
      return
    }
    isLoading.value = true

    try {
      try {
        const capsResp = await get<any>('/config/upload/capabilities')
        capabilities.value = capsResp?.data as UploadCapabilities
      } catch {
        capabilities.value = null
      }

      const response = await get<any>('/config/upload')
      const serverConfig = response.data

      const capMime = capabilities.value?.mime_map || {}
      const addMime = (set: Set<string>, extOrMime: string) => {
        const v = extOrMime.trim().toLowerCase()
        if (v.startsWith('image/')) {
          set.add(v)
          return
        }
        const ext = v.replace(/^\./, '')
        const mapped = capMime[ext]
        if (mapped) {
          set.add(mapped)
          if (ext === 'ico') {
            set.add('image/vnd.microsoft.icon')
          }
        } else {
          switch (ext) {
            case 'jpg':
            case 'jpeg':
              set.add('image/jpeg')
              break
            case 'png':
              set.add('image/png')
              break
            case 'gif':
              set.add('image/gif')
              break
            case 'webp':
              set.add('image/webp')
              break
            case 'bmp':
              set.add('image/bmp')
              break
            case 'svg':
              set.add('image/svg+xml')
              break
            case 'ico':
              set.add('image/x-icon')
              set.add('image/vnd.microsoft.icon')
              break
            case 'apng':
              set.add('image/apng')
              break
            case 'jp2':
              set.add('image/jp2')
              break
            case 'tiff':
            case 'tif':
              set.add('image/tiff')
              break
            case 'tga':
              set.add('image/x-tga')
              break
            default:
              set.add(`image/${ext}`)
          }
        }
      }

      const rawAllowed = Array.isArray(serverConfig?.allowed_types) ? serverConfig.allowed_types : defaultConfig.allowedTypes

      const allowedTypesSet = new Set<string>()
      for (const t of rawAllowed as any[]) {
        if (typeof t !== 'string') {
          continue
        }
        addMime(allowedTypesSet, t)
      }

      uploadConfig.value = {
        maxFileSize: Math.max((serverConfig?.max_file_size || 100) * 1024 * 1024, 50 * 1024 * 1024),
        allowedTypes: Array.from(allowedTypesSet),
        clientMaxConcurrentUploads: serverConfig?.client_max_concurrent_uploads || 3,
        chunkedUpload: {
          enabled: serverConfig?.chunked_upload?.enabled ?? true,
          isAllowChunkUpload: serverConfig?.chunked_upload?.is_allow_chunk_upload ?? false,
          threshold: (serverConfig?.chunked_upload?.threshold || 20) * 1024 * 1024,
          chunkSize: (serverConfig?.chunked_upload?.chunk_size || 2) * 1024 * 1024,
          maxConcurrency: serverConfig?.chunked_upload?.max_concurrency || 4,
          retryCount: serverConfig?.chunked_upload?.retry_count || 3,
        },
      }

      isConfigLoaded.value = true
    } catch (error) {
      console.warn($t('composables.useUploadConfig.errors.loadFailed'), error)
      uploadConfig.value = defaultConfig
      isConfigLoaded.value = true
    } finally {
      isLoading.value = false
    }
  }

  const getClientMaxConcurrentUploads = () => uploadConfig.value.clientMaxConcurrentUploads
  const getChunkedThreshold = () => uploadConfig.value.chunkedUpload.threshold
  const getChunkSize = () => uploadConfig.value.chunkedUpload.chunkSize
  const getChunkedMaxConcurrency = () => uploadConfig.value.chunkedUpload.maxConcurrency
  const getRetryCount = () => uploadConfig.value.chunkedUpload.retryCount

  const isAllowedFileType = (mimeType: string) => uploadConfig.value.allowedTypes.includes(mimeType)

  const isAllowedFileSize = (fileSize: number) => fileSize <= uploadConfig.value.maxFileSize
  const shouldUseChunkedUpload = (fileSize: number) =>
    uploadConfig.value.chunkedUpload.enabled &&
    uploadConfig.value.chunkedUpload.isAllowChunkUpload &&
    fileSize >= uploadConfig.value.chunkedUpload.threshold

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const chunkedUploadEnabled = computed(() => uploadConfig.value.chunkedUpload.enabled)
  const maxFileSize = computed(() => uploadConfig.value.maxFileSize)
  const allowedTypes = computed(() => uploadConfig.value.allowedTypes)

  return {
    uploadConfig: computed(() => uploadConfig.value),
    isConfigLoaded: computed(() => isConfigLoaded.value),
    chunkedUploadEnabled,
    maxFileSize,
    allowedTypes,

    loadConfig,
    getClientMaxConcurrentUploads,
    getChunkedThreshold,
    getChunkSize,
    getChunkedMaxConcurrency,
    getRetryCount,
    isAllowedFileType,
    isAllowedFileSize,
    shouldUseChunkedUpload,
    formatFileSize,
  }
}
