/**
 * 验证相关常量
 */

export const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
export const USERNAME_REGEX = /^[\w\u4e00-\u9fa5]+$/
export const PASSWORD_STRONG_REGEX = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]{6,}$/

export const USERNAME_MIN_LENGTH = 2
export const USERNAME_MAX_LENGTH = 20
export const PASSWORD_MIN_LENGTH = 6
export const PASSWORD_MAX_LENGTH = 20
export const VERIFY_CODE_LENGTH = 6
