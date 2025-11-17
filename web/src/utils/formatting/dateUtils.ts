/**
 * 时间格式化工具函数
 */
import { getCurrentLocale } from '@/utils/locale'
import { useTextThemeStore } from '@/store/textTheme'

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

/**
 * 格式化时间为相对时间显示
 * @param date 时间字符串或Date对象
 * @returns 格式化后的相对时间字符串
 */
export function formatDistanceToNow(date: string | Date): string {
  if (!date) {
    return '-'
  }

  try {
    const targetDate = typeof date === 'string' ? new Date(date) : date
    const now = new Date()
    const diffInSeconds = Math.floor((now.getTime() - targetDate.getTime()) / 1000)

    if (diffInSeconds < 0) {
      return targetDate.toLocaleString(getCurrentLocale())
    }

    const store = useTextThemeStore()

    if (diffInSeconds < 60) {
      return formatText(store.getText('utils.formatting.time.secondsAgo'), { seconds: diffInSeconds })
    }

    const diffInMinutes = Math.floor(diffInSeconds / 60)
    if (diffInMinutes < 60) {
      return formatText(store.getText('utils.formatting.time.minutesAgo'), { minutes: diffInMinutes })
    }

    const diffInHours = Math.floor(diffInMinutes / 60)
    if (diffInHours < 24) {
      return formatText(store.getText('utils.formatting.time.hoursAgo'), { hours: diffInHours })
    }

    const diffInDays = Math.floor(diffInHours / 24)
    if (diffInDays < 7) {
      return formatText(store.getText('utils.formatting.time.daysAgo'), { days: diffInDays })
    }

    if (diffInDays < 30) {
      const diffInWeeks = Math.floor(diffInDays / 7)
      return formatText(store.getText('utils.formatting.time.weeksAgo'), { weeks: diffInWeeks })
    }

    if (diffInDays < 365) {
      const diffInMonths = Math.floor(diffInDays / 30)
      return formatText(store.getText('utils.formatting.time.monthsAgo'), { months: diffInMonths })
    }

    const diffInYears = Math.floor(diffInDays / 365)
    return formatText(store.getText('utils.formatting.time.yearsAgo'), { years: diffInYears })
  } catch {
    if (typeof date === 'string') {
      try {
        return new Date(date).toLocaleString(getCurrentLocale())
      } catch {
        return date
      }
    }
    return '-'
  }
}

export function formatDateTime(date: string | Date): string {
  if (!date) {
    return '-'
  }

  try {
    const targetDate = typeof date === 'string' ? new Date(date) : date
    return targetDate.toLocaleString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  } catch {
    return typeof date === 'string' ? date : '-'
  }
}

export function formatDate(date: string | Date): string {
  if (!date) {
    return '-'
  }

  try {
    const targetDate = typeof date === 'string' ? new Date(date) : date
    return targetDate.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    })
  } catch {
    return typeof date === 'string' ? date : '-'
  }
}

export function formatTime(date: string | Date): string {
  if (!date) {
    return '-'
  }

  try {
    const targetDate = typeof date === 'string' ? new Date(date) : date
    return targetDate.toLocaleTimeString(getCurrentLocale(), {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  } catch {
    return typeof date === 'string' ? date : '-'
  }
}
