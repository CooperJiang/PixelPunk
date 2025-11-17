import { ref, computed } from 'vue'
import type { UploadSettings } from './types'
import { useTexts } from '@/composables/useTexts'

/**
 * 上传设置子模块
 * 管理文件上传相关的配置
 */
export function useUploadSettingsModule() {
  const { $t } = useTexts()
  const settings = ref<UploadSettings>({})

  const allowedImageFormats = computed(
    () =>
      settings.value.allowed_image_formats || [
        'jpg',
        'jpeg',
        'png',
        'gif',
        'webp',
        'bmp',
        'svg',
        'ico',
        'apng',
        'tiff',
        'tif',
        'heic',
        'heif',
      ]
  )

  const maxFileSize = computed(() => {
    const size = settings.value.max_file_size
    return typeof size === 'number' ? size : 20
  })

  const maxBatchSize = computed(() => {
    const size = settings.value.max_batch_size
    return typeof size === 'number' ? size : 100
  })

  const allowedFormatsText = computed(() => {
    const formats = settings.value.allowed_image_formats || ['JPG', 'PNG', 'GIF', 'WebP']
    return formats.map((f) => f.toUpperCase()).join('、')
  })

  const uploadLimitText = computed(() => {
    const maxSize = maxFileSize.value
    const formatsText = allowedFormatsText.value

    if (maxSize === 0) {
      return $t('store.settings.upload.limitTextUnlimited', { formats: formatsText })
    }
    return $t('store.settings.upload.limitTextWithMax', { formats: formatsText, maxSize: String(maxSize) })
  })

  const thumbnailQuality = computed(() => settings.value.thumbnail_quality ?? 85)

  const thumbnailMaxWidth = computed(() => settings.value.thumbnail_max_width ?? 800)

  const thumbnailMaxHeight = computed(() => settings.value.thumbnail_max_height ?? 800)

  const preserveExif = computed(() => settings.value.preserve_exif ?? false)

  const dailyUploadLimit = computed(() => settings.value.daily_upload_limit ?? 0)

  const updateSettings = (newSettings: UploadSettings) => {
    settings.value = { ...settings.value, ...newSettings }
  }

  const reset = () => {
    settings.value = {}
  }

  return {
    uploadSettings: settings,

    allowedImageFormats,
    maxFileSize,
    maxBatchSize,
    allowedFormatsText,
    uploadLimitText,
    thumbnailQuality,
    thumbnailMaxWidth,
    thumbnailMaxHeight,
    preserveExif,
    dailyUploadLimit,

    updateUploadSettings: updateSettings,
    resetUploadSettings: reset,
  }
}
