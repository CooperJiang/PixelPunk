/**
 * MultiSelector 时长相关工具函数
 */

/**
 * 时长格式验证正则表达式
 * 支持格式：30m, 2h, 7d 等 (不再支持秒级)
 */
const DURATION_REGEX = /^(\d+)([mhd])$/i

/**
 * 验证时长格式
 * @param duration 时长字符串，如 '1h', '30m', '7d'
 * @param isGuest 是否游客模式
 * @returns 错误信息，null 表示验证通过
 */
export function validateDuration(duration: string, isGuest = false): string | null {
  if (!duration || duration.trim() === '') {
    return 'empty'
  }

  duration = duration.trim().toLowerCase()

  if (duration === 'permanent') {
    if (isGuest) {
      return 'guestNoPermanent'
    }
    return null
  }

  const match = duration.match(DURATION_REGEX)
  if (!match) {
    return 'invalidFormat'
  }

  const [, valueStr, unit] = match
  const value = parseInt(valueStr)

  if (value <= 0) {
    return 'mustBePositive'
  }

  switch (unit) {
    case 'm': // 分钟
      if (value < 1) {
        return 'minutesMin'
      }
      if (value > 1440 * 365) {
        return 'minutesMax'
      }
      break
    case 'h': // 小时
      if (value > 24 * 365) {
        return 'hoursMax'
      }
      if (isGuest && value > 24 * 30) {
        return 'guestHoursMax'
      }
      break
    case 'd': // 天
      if (value > 365) {
        return 'daysMax'
      }
      if (isGuest && value > 30) {
        return 'guestDaysMax'
      }
      break
  }

  return null
}

export function formatDuration(duration: string, translateFn?: (key: string) => string): string {
  if (!duration) {
    return ''
  }

  duration = duration.trim().toLowerCase()

  if (duration === 'permanent') {
    return translateFn ? translateFn('components.multiSelector.duration.permanent') : 'Permanent'
  }

  const match = duration.match(DURATION_REGEX)
  if (!match) {
    return duration
  }

  const [, valueStr, unit] = match
  const value = parseInt(valueStr)

  if (translateFn) {
    switch (unit) {
      case 'm':
        return `${value}${translateFn('components.multiSelector.duration.minute')}`
      case 'h':
        return `${value}${translateFn('components.multiSelector.duration.hour')}`
      case 'd':
        return `${value}${translateFn('components.multiSelector.duration.day')}`
      default:
        return duration
    }
  } else {
    switch (unit) {
      case 'm':
        return `${value}m`
      case 'h':
        return `${value}h`
      case 'd':
        return `${value}d`
      default:
        return duration
    }
  }
}

export function normalizeDuration(duration: string): string {
  if (!duration) {
    return ''
  }

  duration = duration.trim().toLowerCase()

  if (duration === 'permanent') {
    return 'permanent'
  }

  const match = duration.match(DURATION_REGEX)
  if (!match) {
    return duration
  }

  const [, valueStr, unit] = match
  return `${valueStr}${unit.toLowerCase()}`
}

export function isDurationExists(duration: string, existingDurations: string[]): boolean {
  const normalized = normalizeDuration(duration)
  return existingDurations.some((existing) => normalizeDuration(existing) === normalized)
}

export function getDurationIcon(duration: string): string {
  if (!duration) {
    return 'fas fa-tag'
  }

  duration = duration.trim().toLowerCase()

  if (duration === 'permanent') {
    return 'fas fa-infinity'
  }

  const match = duration.match(DURATION_REGEX)
  if (!match) {
    return 'fas fa-tag'
  }

  const [, , unit] = match

  switch (unit) {
    case 'm':
      return 'fas fa-clock'
    case 'h':
      return 'fas fa-hourglass-half'
    case 'd':
      return 'fas fa-calendar-day'
    default:
      return 'fas fa-tag'
  }
}

export function compareDurations(a: string, b: string): number {
  if (a === 'permanent' && b !== 'permanent') {
    return 1
  }
  if (b === 'permanent' && a !== 'permanent') {
    return -1
  }
  if (a === 'permanent' && b === 'permanent') {
    return 0
  }

  const getSeconds = (duration: string): number => {
    const match = duration.match(DURATION_REGEX)
    if (!match) {
      return 0
    }

    const [, valueStr, unit] = match
    const value = parseInt(valueStr)

    switch (unit) {
      case 'm':
        return value * 60
      case 'h':
        return value * 3600
      case 'd':
        return value * 86400
      default:
        return 0
    }
  }

  return getSeconds(a) - getSeconds(b)
}

export function getSuggestedDurations(
  isGuest = false,
  translateFn?: (key: string) => string
): Array<{ value: string; label: string }> {
  const suggestions = translateFn
    ? [
        { value: '1m', label: `1${translateFn('components.multiSelector.duration.minute')}` },
        { value: '5m', label: `5${translateFn('components.multiSelector.duration.minute')}` },
        { value: '30m', label: `30${translateFn('components.multiSelector.duration.minute')}` },
        { value: '1h', label: `1${translateFn('components.multiSelector.duration.hour')}` },
        { value: '3h', label: `3${translateFn('components.multiSelector.duration.hour')}` },
        { value: '6h', label: `6${translateFn('components.multiSelector.duration.hour')}` },
        { value: '12h', label: `12${translateFn('components.multiSelector.duration.hour')}` },
        { value: '1d', label: `1${translateFn('components.multiSelector.duration.day')}` },
        { value: '3d', label: `3${translateFn('components.multiSelector.duration.day')}` },
        { value: '7d', label: `7${translateFn('components.multiSelector.duration.day')}` },
        { value: '15d', label: `15${translateFn('components.multiSelector.duration.day')}` },
        { value: '30d', label: `30${translateFn('components.multiSelector.duration.day')}` },
      ]
    : [
        { value: '1m', label: '1m' },
        { value: '5m', label: '5m' },
        { value: '30m', label: '30m' },
        { value: '1h', label: '1h' },
        { value: '3h', label: '3h' },
        { value: '6h', label: '6h' },
        { value: '12h', label: '12h' },
        { value: '1d', label: '1d' },
        { value: '3d', label: '3d' },
        { value: '7d', label: '7d' },
        { value: '15d', label: '15d' },
        { value: '30d', label: '30d' },
      ]

  if (!isGuest) {
    suggestions.push({
      value: 'permanent',
      label: translateFn ? translateFn('components.multiSelector.duration.permanent') : 'Permanent',
    })
  }

  return suggestions
}
