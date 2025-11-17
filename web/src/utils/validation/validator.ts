/**
 * 提供常用的数据验证和格式化功能（支持可选 i18n）
 */
type TranslationFunction = (key: string, params?: Record<string, string | number>) => string

export interface ValidationResult {
  valid: boolean
  message?: string
  value?: any
}

export interface ValidationRule {
  required?: boolean
  type?: 'string' | 'number' | 'email' | 'url' | 'phone' | 'password' | 'custom'
  min?: number
  max?: number
  pattern?: RegExp
  custom?: (value: any) => ValidationResult
  message?: string
}

/* 批量验证结果 */
export interface BatchValidationResult {
  valid: boolean
  errors: Record<string, string[]>
  values: Record<string, any>
}

class Validator {
  private static readonly PATTERNS = {
    email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
    phone: /^1[3-9]\d{9}$/,
    url: /^https?:\/\/([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/,
    username: /^[a-zA-Z0-9_-]{3,20}$/,
    password: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{8,}$/,
    ipv4: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
    mac: /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/,
    uuid: /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i,
    hex: /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/,
    base64: /^[A-Za-z0-9+/]+=*$/,
  }

  /* 获取验证消息（带翻译） */
  private static getMessages($t: TranslationFunction) {
    return {
      required: $t('utils.errors.required'),
      type: $t('utils.errors.type'),
      min: $t('utils.errors.min', { min: '{min}' }),
      max: $t('utils.errors.max', { max: '{max}' }),
      length: $t('utils.errors.length'),
      pattern: $t('utils.errors.pattern'),
      email: $t('utils.errors.email'),
      phone: $t('utils.errors.phone'),
      url: $t('utils.errors.url'),
      username: $t('utils.errors.username'),
      password: $t('utils.errors.password'),
      custom: $t('utils.errors.custom'),
    }
  }

  /* 向后兼容 - 默认英文消息 */
  private static readonly MESSAGES = {
    required: 'This field is required',
    type: 'Invalid field type',
    min: 'Value is too small, min is {min}',
    max: 'Value is too large, max is {max}',
    length: 'Invalid length',
    pattern: 'Invalid format',
    email: 'Please enter a valid email address',
    phone: 'Please enter a valid phone number',
    url: 'Please enter a valid URL',
    username: 'Username must be 3-20 chars: letters, numbers, _ or -',
    password: 'Password must contain upper/lowercase letters and numbers, min 8 chars',
    custom: 'Validation failed',
  }

  static validate(value: any, rules: ValidationRule[], $t?: TranslationFunction): ValidationResult {
    for (const rule of rules) {
      const result = this.validateSingle(value, rule, $t)
      if (!result.valid) {
        return result
      }
    }
    return { valid: true, value }
  }

  private static validateSingle(value: any, rule: ValidationRule, $t?: TranslationFunction): ValidationResult {
    const msgs = $t ? this.getMessages($t) : this.MESSAGES

    if (rule.required && this.isEmpty(value)) {
      return {
        valid: false,
        message: rule.message || msgs.required,
      }
    }

    if (!rule.required && this.isEmpty(value)) {
      return { valid: true, value }
    }

    if (rule.type) {
      const typeResult = this.validateType(value, rule.type, $t)
      if (!typeResult.valid) {
        return {
          valid: false,
          message: rule.message || typeResult.message || msgs.type,
        }
      }
    }

    if (rule.min !== undefined || rule.max !== undefined) {
      const rangeResult = this.validateRange(value, rule.min, rule.max, $t)
      if (!rangeResult.valid) {
        return {
          valid: false,
          message: rule.message || rangeResult.message,
        }
      }
    }

    if (rule.pattern) {
      if (typeof value === 'string' && !rule.pattern.test(value)) {
        return {
          valid: false,
          message: rule.message || msgs.pattern,
        }
      }
    }

    if (rule.custom) {
      const customResult = rule.custom(value)
      if (!customResult.valid) {
        return {
          valid: false,
          message: rule.message || customResult.message || msgs.custom,
        }
      }
    }

    return { valid: true, value }
  }

  private static isEmpty(value: any): boolean {
    if (value === null || value === undefined) {
      return true
    }
    if (typeof value === 'string') {
      return value.trim() === ''
    }
    if (Array.isArray(value)) {
      return value.length === 0
    }
    if (typeof value === 'object') {
      return Object.keys(value).length === 0
    }
    return false
  }

  private static validateType(value: any, type: string, $t?: TranslationFunction): ValidationResult {
    const msgs = $t ? this.getMessages($t) : this.MESSAGES

    switch (type) {
      case 'string':
        if (typeof value !== 'string') {
          return {
            valid: false,
            message: $t ? $t('utils.errors.mustBeString') : 'Must be a string',
          }
        }
        break

      case 'number':
        if (typeof value !== 'number' || isNaN(value)) {
          return {
            valid: false,
            message: $t ? $t('utils.errors.mustBeNumber') : 'Must be a number',
          }
        }
        break

      case 'email':
        if (!this.PATTERNS.email.test(String(value))) {
          return { valid: false, message: msgs.email }
        }
        break

      case 'phone':
        if (!this.PATTERNS.phone.test(String(value))) {
          return { valid: false, message: msgs.phone }
        }
        break

      case 'url':
        if (!this.PATTERNS.url.test(String(value))) {
          return { valid: false, message: msgs.url }
        }
        break

      case 'password':
        if (!this.PATTERNS.password.test(String(value))) {
          return { valid: false, message: msgs.password }
        }
        break
    }
    return { valid: true }
  }

  private static validateRange(value: any, min?: number, max?: number, $t?: TranslationFunction): ValidationResult {
    const length =
      typeof value === 'string' ? value.length : typeof value === 'number' ? value : Array.isArray(value) ? value.length : 0

    if (min !== undefined && length < min) {
      const message = $t ? $t('utils.errors.min', { min: String(min) }) : this.MESSAGES.min.replace('{min}', String(min))
      return {
        valid: false,
        message,
      }
    }

    if (max !== undefined && length > max) {
      const message = $t ? $t('utils.errors.max', { max: String(max) }) : this.MESSAGES.max.replace('{max}', String(max))
      return {
        valid: false,
        message,
      }
    }

    return { valid: true }
  }

  static batchValidate(
    data: Record<string, unknown>,
    rules: Record<string, ValidationRule[]>,
    $t?: TranslationFunction
  ): BatchValidationResult {
    const result: BatchValidationResult = {
      valid: true,
      errors: {},
      values: {},
    }

    for (const [field, fieldRules] of Object.entries(rules)) {
      const value = data[field]
      const validation = this.validate(value, fieldRules, $t)

      if (!validation.valid) {
        result.valid = false
        const defaultMsg = $t ? $t('utils.errors.custom') : 'Validation failed'
        result.errors[field] = [validation.message || defaultMsg]
      } else {
        result.values[field] = validation.value
      }
    }

    return result
  }

  static isEmail(email: string): boolean {
    return this.PATTERNS.email.test(email)
  }

  static isPhone(phone: string): boolean {
    return this.PATTERNS.phone.test(phone)
  }

  static isUrl(url: string): boolean {
    return this.PATTERNS.url.test(url)
  }

  static isUsername(username: string): boolean {
    return this.PATTERNS.username.test(username)
  }

  static isPassword(password: string): boolean {
    return this.PATTERNS.password.test(password)
  }

  static isIPv4(ip: string): boolean {
    return this.PATTERNS.ipv4.test(ip)
  }

  static isMacAddress(mac: string): boolean {
    return this.PATTERNS.mac.test(mac)
  }

  static isUUID(uuid: string): boolean {
    return this.PATTERNS.uuid.test(uuid)
  }

  static isHexColor(color: string): boolean {
    return this.PATTERNS.hex.test(color)
  }

  static isBase64(str: string): boolean {
    return this.PATTERNS.base64.test(str)
  }

  static sanitize = {
    html(str: string): string {
      return str.replace(/<[^>]*>/g, '')
    },

    sql(str: string): string {
      return str.replace(/['"\\;]/g, '')
    },

    xss(str: string): string {
      return str
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#x27;')
        .replace(/\//g, '&#x2F;')
    },

    whitespace(str: string): string {
      return str.trim().replace(/\s+/g, ' ')
    },

    special(str: string): string {
      return str.replace(/[^\w\s-]/gi, '')
    },

    number(str: string): string {
      return str.replace(/[^\d.-]/g, '')
    },

    email(str: string): string {
      return str.toLowerCase().trim()
    },

    phone(str: string): string {
      return str.replace(/[^\d]/g, '')
    },

    url(str: string): string {
      try {
        const url = new URL(str)
        return url.toString()
      } catch {
        return str
      }
    },
  }

  static format = {
    fileSize(bytes: number, decimals = 2): string {
      if (bytes === 0) {
        return '0 Bytes'
      }
      const k = 1024
      const dm = decimals < 0 ? 0 : decimals
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
    },

    number(num: number, locale = 'zh-CN'): string {
      return new Intl.NumberFormat(locale).format(num)
    },

    currency(amount: number, currency = 'CNY', locale = 'zh-CN'): string {
      return new Intl.NumberFormat(locale, {
        style: 'currency',
        currency,
      }).format(amount)
    },

    percentage(value: number, decimals = 2): string {
      return `${(value * 100).toFixed(decimals)}%`
    },

    duration(milliseconds: number, $t?: TranslationFunction): string {
      const seconds = Math.floor(milliseconds / 1000)
      const minutes = Math.floor(seconds / 60)
      const hours = Math.floor(minutes / 60)
      const days = Math.floor(hours / 24)

      if ($t) {
        if (days > 0) {
          return $t('utils.time.daysHours', { d: days, h: hours % 24 })
        }
        if (hours > 0) {
          return $t('utils.time.hoursMinutes', { h: hours, m: minutes % 60 })
        }
        if (minutes > 0) {
          return $t('utils.time.minutesSeconds', { m: minutes, s: seconds % 60 })
        }
        return $t('utils.time.seconds', { n: seconds })
      }

      // Backward compatibility - default English
      if (days > 0) {
        return `${days}d ${hours % 24}h`
      }
      if (hours > 0) {
        return `${hours}h ${minutes % 60}m`
      }
      if (minutes > 0) {
        return `${minutes}m ${seconds % 60}s`
      }
      return `${seconds}s`
    },

    relative(date: Date, locale = 'zh-CN'): string {
      const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })
      const now = new Date()
      const diff = date.getTime() - now.getTime()
      const seconds = Math.floor(diff / 1000)
      const minutes = Math.floor(seconds / 60)
      const hours = Math.floor(minutes / 60)
      const days = Math.floor(hours / 24)

      if (Math.abs(days) > 0) {
        return rtf.format(days, 'day')
      }
      if (Math.abs(hours) > 0) {
        return rtf.format(hours, 'hour')
      }
      if (Math.abs(minutes) > 0) {
        return rtf.format(minutes, 'minute')
      }
      return rtf.format(seconds, 'second')
    },
  }

  static checkPasswordStrength(
    password: string,
    $t?: TranslationFunction
  ): {
    score: number
    level: 'weak' | 'medium' | 'strong' | 'very-strong'
    feedback: string[]
  } {
    let score = 0
    const feedback: string[] = []

    if (password.length >= 8) {
      score += 1
    } else {
      feedback.push($t ? $t('utils.password.feedback.minLength') : 'Password must be at least 8 characters')
    }

    if (password.length >= 12) {
      score += 1
    }

    if (/[a-z]/.test(password)) {
      score += 1
    } else {
      feedback.push($t ? $t('utils.password.feedback.requireLowercase') : 'Include at least one lowercase letter')
    }

    if (/[A-Z]/.test(password)) {
      score += 1
    } else {
      feedback.push($t ? $t('utils.password.feedback.requireUppercase') : 'Include at least one uppercase letter')
    }

    if (/[0-9]/.test(password)) {
      score += 1
    } else {
      feedback.push($t ? $t('utils.password.feedback.requireNumber') : 'Include at least one number')
    }

    if (/[^a-zA-Z0-9]/.test(password)) {
      score += 1
    } else {
      feedback.push($t ? $t('utils.password.feedback.suggestSpecial') : 'Consider adding special characters')
    }

    if (!/(.)\1{2,}/.test(password)) {
      score += 1
    } else {
      feedback.push($t ? $t('utils.password.feedback.avoidRepeat') : 'Avoid repeating characters')
    }

    const levels = ['weak', 'weak', 'medium', 'medium', 'strong', 'strong', 'very-strong', 'very-strong'] as const

    return {
      score,
      level: levels[Math.min(score, levels.length - 1)],
      feedback,
    }
  }
}

/* 向后兼容 - 默认中文 */
export const ValidationRules = {
  required: { required: true },
  email: { type: 'email' as const, required: true },
  phone: { type: 'phone' as const, required: true },
  url: { type: 'url' as const },
  password: { type: 'password' as const, required: true },
  username: {
    required: true,
    pattern: Validator.PATTERNS.username,
    message: 'Username must be 3-20 chars: letters, numbers, _ or -',
  },
  minLength: (length: number) => ({ min: length, message: `At least ${length} characters` }),
  maxLength: (length: number) => ({ max: length, message: `At most ${length} characters` }),
  range: (min: number, max: number) => ({ min, max }),
  numeric: { type: 'number' as const },
  positive: {
    type: 'number' as const,
    custom: (value: number) => (value > 0 ? { valid: true } : { valid: false, message: 'Must be a positive number' }),
  },
}

export function getValidationRules($t: TranslationFunction) {
  return {
    required: { required: true, message: $t('utils.errors.required') },
    email: { type: 'email' as const, required: true, message: $t('utils.errors.email') },
    phone: { type: 'phone' as const, required: true, message: $t('utils.errors.phone') },
    url: { type: 'url' as const, message: $t('utils.errors.url') },
    password: { type: 'password' as const, required: true, message: $t('utils.errors.password') },
    username: {
      required: true,
      pattern: Validator.PATTERNS.username,
      message: $t('utils.errors.username'),
    },
    minLength: (length: number) => ({ min: length, message: $t('utils.errors.minLength', { n: String(length) }) }),
    maxLength: (length: number) => ({ max: length, message: $t('utils.errors.maxLength', { n: String(length) }) }),
    range: (min: number, max: number) => ({ min, max }),
    numeric: { type: 'number' as const, message: $t('utils.errors.mustBeNumber') },
    positive: {
      type: 'number' as const,
      custom: (value: number) =>
        value > 0 ? { valid: true } : { valid: false, message: $t('utils.errors.mustBePositive') },
    },
  }
}

export default Validator
