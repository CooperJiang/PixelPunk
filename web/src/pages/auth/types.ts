/**
 * 认证页面类型定义
 */

export interface LoginFormState {
  account: string
  password: string
  remember: boolean
}

/**
 * 记住的登录信息
 */
export interface RememberedLogin {
  account: string
  password: string
}

/**
 * 注册表单状态
 */
export interface RegisterFormState {
  username: string
  email: string
  password: string
  confirmPassword: string
  code: string
}
