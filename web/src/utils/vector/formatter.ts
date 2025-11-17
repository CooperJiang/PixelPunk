import type { TranslationFunction } from '@/composables/useTexts'

export interface VectorLog {
  id: number
  file_id: string
  action: string
  type: string
  data: Record<string, any>
  model: string
  duration: number
  error_code?: string
  task_id?: string
  created_at: string
}

export function formatVectorLog(log: VectorLog, $t: TranslationFunction): string {
  const { type, data } = log

  try {
    switch (type) {
      case 'vector.reset':
        if (data.model) {
          return $t('admin.vectors.log.reset_with_model', {
            model: String(data.model),
          })
        }
        return $t('admin.vectors.log.reset')

      case 'vector.retry':
        if (data.retry_count) {
          return $t('admin.vectors.log.retry_with_count', {
            retry_count: String(data.retry_count),
          })
        }
        return $t('admin.vectors.log.retry')

      case 'vector.delete':
        return $t('admin.vectors.log.delete')

      case 'vector.retry_all':
        return $t('admin.vectors.log.retry_all', {
          retry_count: String(data.retry_count || 0),
        })

      case 'vector.recover_stuck':
        return $t('admin.vectors.log.recover_stuck', {
          count: String(data.count || 0),
        })

      case 'vector.start':
        return $t('admin.vectors.log.start', {
          model: String(data.model || ''),
        })

      default:
        return `${type}: ${JSON.stringify(data)}`
    }
  } catch (error) {
    return `${type}: ${JSON.stringify(data)}`
  }
}
