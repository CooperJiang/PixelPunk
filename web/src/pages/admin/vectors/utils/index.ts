import { formatDistanceToNow } from '@/utils/formatting/dateUtils'
import { STATUS_TEXT_MAP, STATUS_ICON_MAP, COLLECTION_STATUS_TEXT_MAP, MODEL_DESCRIPTIONS } from '../constants'
import { useTexts } from '@/composables/useTexts'

const { $t } = useTexts()

export const createDebounce = <T extends unknown[]>(func: (...args: T) => Promise<unknown> | unknown, wait: number) => {
  let timeout: number | null = null
  return (...args: T) => {
    return new Promise((resolve, reject) => {
      if (timeout !== null) {
        clearTimeout(timeout)
      }
      timeout = window.setTimeout(async () => {
        try {
          const result = await func(...args)
          resolve(result)
        } catch (error) {
          reject(error)
        }
      }, wait)
    })
  }
}

export const getStatusText = (status: string): string => {
  return STATUS_TEXT_MAP[status] || status
}

export const getStatusIcon = (status: string): string => {
  return STATUS_ICON_MAP[status] || 'fas fa-question-circle'
}

export const getCollectionStatusText = (status: string): string => {
  return COLLECTION_STATUS_TEXT_MAP[status] || status
}

export const getModelDescription = (model: string): string => {
  return MODEL_DESCRIPTIONS[model] || $t('admin.vectors.utils.noDescription')
}

export const formatTime = (time: string): string => {
  return formatDistanceToNow(time)
}

export const formatDuration = (duration: number): string => {
  if (!duration || duration === 0) return '-'
  if (duration < 60) return $t('admin.vectors.utils.seconds', { duration })
  const minutes = Math.floor(duration / 60)
  const seconds = duration % 60
  return seconds > 0
    ? $t('admin.vectors.utils.minutesSeconds', { minutes, seconds })
    : $t('admin.vectors.utils.minutes', { minutes })
}

export const canRetry = (vector: { status: string }): boolean => {
  return vector.status === 'failed' || vector.status === 'reset'
}

export const handleError = (error: unknown, context: string, toast: { error: (msg: string) => void }) => {
  const err = error as Record<string, unknown>
  const errorCode =
    (err.code as number) || (err.status as number) || ((err.response as Record<string, unknown>)?.status as number)
  const errorMessage =
    (err.message as string) ||
    (((err.response as Record<string, unknown>)?.data as Record<string, unknown>)?.message as string) ||
    $t('admin.vectors.utils.unknownError')

  if (errorCode === 404) {
    toast.error($t('admin.vectors.utils.error404', { context }))
  } else if (errorCode >= 500) {
    toast.error($t('admin.vectors.utils.error500', { context }))
  } else if (errorCode === 403) {
    toast.error($t('admin.vectors.utils.error403', { context }))
  } else {
    toast.error($t('admin.vectors.utils.errorGeneral', { context, message: errorMessage }))
  }
}
