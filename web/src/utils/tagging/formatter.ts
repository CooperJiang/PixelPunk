import type { TranslationFunction } from '@/composables/useTexts'

export interface TaggingLog {
  id: number
  file_id: string
  file_type: string
  status: string
  action: string
  type: string
  data: Record<string, any>
  operator_id: number
  duration: number
  created_at: string
  updated_at: string
}

/**
 * 格式化打标日志消息
 * @param log - 日志对象
 * @param $t - 翻译函数
 * @returns 格式化后的消息字符串
 */
export function formatTaggingLog(log: TaggingLog, $t: TranslationFunction): string {
  const { type, data } = log

  try {
    switch (type) {
      case 'tagging.reset_stuck':
        return $t('admin.tagging.log.reset_stuck', {
          count: String(data.count || 0),
          time_threshold: String(data.time_threshold || 0),
        })

      case 'tagging.scheduled_scan':
        if (data.errors) {
          return $t('admin.tagging.log.scheduled_scan_with_errors', {
            duration: String(data.duration?.toFixed?.(2) || data.duration || 0),
            reset_count: String(data.reset_count || 0),
            errors: String(data.errors || ''),
          })
        }
        return $t('admin.tagging.log.scheduled_scan', {
          duration: String(data.duration?.toFixed?.(2) || data.duration || 0),
          reset_count: String(data.reset_count || 0),
        })

      case 'tagging.ignore':
        if (data.reason) {
          return $t('admin.tagging.log.ignore_with_reason', {
            reason: String(data.reason),
          })
        }
        return $t('admin.tagging.log.ignore')

      case 'tagging.unignore':
        return $t('admin.tagging.log.unignore')

      case 'tagging.retry_failed_all':
        return $t('admin.tagging.log.retry_failed_all')

      case 'tagging.retry':
        return $t('admin.tagging.log.retry')

      case 'tagging.ai_completed':
        return $t('admin.tagging.log.ai_completed', {
          prompt_tokens: String(data.prompt_tokens || 0),
          completion_tokens: String(data.completion_tokens || 0),
          total_tokens: String(data.total_tokens || 0),
        })

      default:
        // 如果类型未知，返回类型和数据
        return `${type}: ${JSON.stringify(data)}`
    }
  } catch (error) {
    // 出错时返回类型和数据
    return `${type}: ${JSON.stringify(data)}`
  }
}
