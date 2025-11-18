/**
 * 上传执行操作
 * 负责实际的文件上传逻辑

 */
import { useToast } from '@/components/Toast/useToast'
import { guestUpload, uploadFile } from '@/api/file'
import { useAuthStore } from '@/store/auth'
import generateFingerprint from '@/utils/system/fingerprint'
import { useTexts } from '@/composables/useTexts'
import { globalOptions, translations } from '../state'
import { prepareWatermarkConfig } from '../utils/watermark'
import { removeUploadItem } from './queue'
import type { UploadItem } from '../types'

export async function uploadRegularFile(item: UploadItem) {
  const toast = useToast()
  const { $t } = useTexts()
  const authStore = useAuthStore()

  if (item.status === 'failed') {
    return
  }

  const currentSessionId = item.uploadSessionId

  try {
    item.status = 'uploading'
    item.progress = 0
    item.statusMessage = translations.value?.upload.uploading || 'Uploading...'

    const startTime = Date.now()

    let result: any

    if (!authStore.isLoggedIn) {
      const watermarkConfig = await prepareWatermarkConfig(item)
      result = await guestUpload(
        item.file,
        {
          access_level: 'public', // 游客默认公开
          optimize: globalOptions.value.optimize,
          storage_duration: globalOptions.value.storageDuration || '7d',
          fingerprint: generateFingerprint(),
          ...watermarkConfig,
        },
        (progressEvent) => {
          if (item.uploadSessionId !== currentSessionId || item.status === 'failed') {
            return
          }
          if (progressEvent.lengthComputable) {
            item.progress = Math.round((progressEvent.loaded / progressEvent.total) * 100)

            const elapsed = (Date.now() - startTime) / 1000
            if (elapsed > 0) {
              item.speed = Math.round(progressEvent.loaded / elapsed / 1024)

              const remainingBytes = progressEvent.total - progressEvent.loaded
              item.remainingTime = item.speed > 0 ? Math.round(remainingBytes / (item.speed * 1024)) : 0

              const speedText = item.speed > 1024 ? `${(item.speed / 1024).toFixed(1)} MB/s` : `${item.speed} KB/s`

              const remainingText =
                item.remainingTime > 60
                  ? $t('store.upload.upload.remainingTimeMinutes', { minutes: Math.floor(item.remainingTime / 60), seconds: item.remainingTime % 60 })
                  : $t('store.upload.upload.remainingTimeSeconds', { seconds: item.remainingTime })

              item.statusMessage = $t('store.upload.upload.uploadingProgress', { progress: item.progress, speed: speedText, remaining: remainingText })
            }
          }
        }
      )

      item.result = result.data.file_info

      if (result.data.remaining_count !== undefined) {
        toast.success(
          $t('upload.uploadProgress.toast.uploadSuccessWithRemainingCount', { count: result.data.remaining_count })
        )
      }

      const guestFileData = item.result
      if (guestFileData && globalOptions.value.watermarkEnabled) {
        if (guestFileData.watermark_applied === false && guestFileData.watermark_failure_reason) {
          toast.warning(
            $t('upload.uploadProgress.toast.watermarkApplyFailed', {
              reason: guestFileData.watermark_failure_reason,
            })
          )
        }
      }
    } else {
      const watermarkConfig = await prepareWatermarkConfig(item)
      result = await uploadFile(
        item.file,
        {
          folder_id: globalOptions.value.folderId || undefined,
          access_level: globalOptions.value.accessLevel,
          optimize: globalOptions.value.optimize,
          storage_duration: globalOptions.value.storageDuration,
          ...watermarkConfig,
        },
        (progressEvent) => {
          if (item.uploadSessionId !== currentSessionId || item.status === 'failed') {
            return
          }
          if (progressEvent.lengthComputable) {
            item.progress = Math.round((progressEvent.loaded / progressEvent.total) * 100)

            const elapsed = (Date.now() - startTime) / 1000
            if (elapsed > 0) {
              item.speed = Math.round(progressEvent.loaded / elapsed / 1024)

              const remainingBytes = progressEvent.total - progressEvent.loaded
              item.remainingTime = item.speed > 0 ? Math.round(remainingBytes / (item.speed * 1024)) : 0

              const speedText = item.speed > 1024 ? `${(item.speed / 1024).toFixed(1)} MB/s` : `${item.speed} KB/s`

              const remainingText =
                item.remainingTime > 60
                  ? $t('store.upload.upload.remainingTimeMinutes', { minutes: Math.floor(item.remainingTime / 60), seconds: item.remainingTime % 60 })
                  : $t('store.upload.upload.remainingTimeSeconds', { seconds: item.remainingTime })

              item.statusMessage = $t('store.upload.upload.uploadingProgress', { progress: item.progress, speed: speedText, remaining: remainingText })
            }
          }
        },
        { silent: true } // 静默上传，不显示默认toast
      )

      item.result = result.data || result
    }

    if (item.status === 'failed') {
      return
    }

    const fileData = item.result
    if (fileData && globalOptions.value.watermarkEnabled) {
      if (fileData.watermark_applied === false && fileData.watermark_failure_reason) {
        toast.warning(
          $t('upload.uploadProgress.toast.watermarkApplyFailed', { reason: fileData.watermark_failure_reason })
        )
      }
    }

    if (fileData && fileData.thumbnail_generation_failed && fileData.thumbnail_failure_reason) {
      toast.warning(
        $t('upload.uploadProgress.toast.thumbnailGenerationFailed', { reason: fileData.thumbnail_failure_reason })
      )
    }

    item.status = 'completed'
    item.progress = 100
    item.statusMessage = translations.value?.upload.uploadComplete || 'Upload completed'

    if (globalOptions.value.autoRemove) {
      setTimeout(() => {
        removeUploadItem(item.id)
      }, 500)
    }
  } catch (error: unknown) {
    item.status = 'failed'
    item.error = (error as any).message || translations.value?.upload.uploadFailed || 'Upload failed'
    item.statusMessage = $t('store.upload.upload.uploadFailedWithError', { error: item.error })
  }
}
