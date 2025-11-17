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
 * 日期格式化函数
 * @param dateString 日期字符串
 * @param format 格式化模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns 格式化后的日期字符串
 */
export function formatDate(dateString: string, format: string = 'YYYY-MM-DD HH:mm:ss'): string {
  if (!dateString) {
    return ''
  }

  const date = new Date(dateString)

  if (isNaN(date.getTime())) {
    return 'Invalid date'
  }

  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hours = date.getHours()
  const minutes = date.getMinutes()
  const seconds = date.getSeconds()

  const formatMap: Record<string, number | string> = {
    YYYY: year,
    MM: padZero(month),
    DD: padZero(day),
    HH: padZero(hours),
    mm: padZero(minutes),
    ss: padZero(seconds),
    M: month,
    D: day,
    H: hours,
    m: minutes,
    s: seconds,
  }

  let result = format
  for (const [key, value] of Object.entries(formatMap)) {
    result = result.replace(key, String(value))
  }

  return result
}

export function formatDateTime(dateString: string): string {
  return formatDate(dateString, 'YYYY-MM-DD HH:mm:ss')
}

export function formatTime(dateString: string): string {
  return formatDate(dateString, 'HH:mm:ss')
}

export function formatSize(bytes: number, decimals: number = 2): string {
  if (bytes === 0) {
    return '0 Bytes'
  }

  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(decimals))} ${sizes[i]}`
}

export const formatFileSize = formatSize

function padZero(num: number): string {
  return num < 10 ? `0${num}` : `${num}`
}

export function handleBooleanParam<T extends Record<string, any>>(params: T, key: keyof T): T {
  if (params[key] === false || (params[key] === undefined && key in params)) {
    const newParams = { ...params }
    delete newParams[key]
    return newParams
  }
  return params
}

export function handleBooleanParams<T extends Record<string, any>>(params: T, keys: Array<keyof T>): T {
  let newParams = { ...params }
  for (const key of keys) {
    newParams = handleBooleanParam(newParams, key)
  }
  return newParams
}

export function processParams<T extends Record<string, any>>(
  params: T,
  options: {
    booleanFields?: Array<keyof T>
    removeEmpty?: boolean
    removeUndefined?: boolean
    removeNull?: boolean
  } = {}
): Record<string, any> {
  const { booleanFields = [], removeEmpty = true, removeUndefined = true, removeNull = false } = options

  let processedParams = { ...params }

  if (booleanFields.length > 0) {
    processedParams = handleBooleanParams(processedParams, booleanFields)
  }

  Object.keys(processedParams).forEach((key) => {
    const value = processedParams[key]

    if ((removeUndefined && value === undefined) || (removeNull && value === null) || (removeEmpty && value === '')) {
      delete processedParams[key]
    }
  })

  return processedParams
}

export function formatNumber(num: number, digits: number = 1): string {
  if (num === undefined || num === null || isNaN(num)) {
    return '0'
  }
  if (num === 0) {
    return '0'
  }

  const safeNum = Number(num) || 0

  const lookup = [
    { value: 1e9, symbol: 'B' },
    { value: 1e6, symbol: 'M' },
    { value: 1e3, symbol: 'K' },
  ]

  for (const item of lookup) {
    if (safeNum >= item.value) {
      return (safeNum / item.value).toFixed(digits) + item.symbol
    }
  }

  return safeNum.toString()
}

export function formatPercentage(value: number, total: number, decimals: number = 1): string {
  if (total === 0) {
    return '0%'
  }
  const percentage = (value / total) * 100
  return `${percentage.toFixed(decimals)}%`
}

export function formatRelativeTime(dateString: string): string {
  const now = new Date()
  const date = new Date(dateString)
  const diffInMs = now.getTime() - date.getTime()

  const diffInMinutes = Math.floor(diffInMs / (1000 * 60))
  const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60))
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24))

  const store = useTextThemeStore()

  if (diffInMinutes < 1) {
    return store.getText('utils.formatting.time.justNow')
  }
  if (diffInMinutes < 60) {
    return formatText(store.getText('utils.formatting.time.minutesAgo'), { minutes: diffInMinutes })
  }
  if (diffInHours < 24) {
    return formatText(store.getText('utils.formatting.time.hoursAgo'), { hours: diffInHours })
  }
  if (diffInDays < 30) {
    return formatText(store.getText('utils.formatting.time.daysAgo'), { days: diffInDays })
  }

  return date.toLocaleDateString(getCurrentLocale())
}

export function getStorageStatusClass(percentage: number): string {
  if (percentage >= 90) {
    return 'storage-critical'
  }
  if (percentage >= 75) {
    return 'storage-warning'
  }
  if (percentage >= 50) {
    return 'storage-moderate'
  }
  return 'storage-safe'
}
