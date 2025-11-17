import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse, type CancelTokenSource } from 'axios'
import { StorageUtil } from '../storage'
import { HTTP_STATUS, REQUEST_TIMEOUT, TOKEN_KEY } from '@/constants'
// Avoid importing router here to prevent circular dependencies with '@/router'
import {
  ErrorCodes,
  getErrorCategory,
  shouldShowError,
  shouldThrowError,
  type ApiError,
  type ApiResult,
  type ApiSuccess,
  type ExtendedRequestConfig,
} from './http-types'
import { useTextThemeStore } from '@/store/textTheme'
import {
  DEFAULT_LOCALE,
  DEFAULT_THEME,
  getLocale,
  isLocaleSupported,
  type SupportedLocale,
} from '@/locales'
import type { TextTheme } from '@/locales/zh-CN'

/* 动态toast辅助函数 */
const showToast = {
  success: (message: string) => {
    import('@/components/Toast/useToast')
      .then(({ useToast }) => {
        const toast = useToast()
        toast.success(message)
      })
      .catch(() => {})
  },
  error: (message: string) => {
    import('@/components/Toast/useToast')
      .then(({ useToast }) => {
        const toast = useToast()
        toast.error(message)
      })
      .catch(() => {})
  },
}

/* 获取翻译文本：优先使用 store，缺失时回退到默认语言包 */
const getTranslation = async (key: string, params?: Record<string, string | number>): Promise<string> => {
  let text: string | undefined

  try {
    const textThemeStore = useTextThemeStore()
    text = textThemeStore.getText(key)
  } catch (_error) {
    // Store 尚未初始化或 Pinia 未就绪，稍后回退
  }

  if (!text || text === key) {
    text = await getFallbackTranslation(key)
  }

  return applyParams(text || key, params)
}

async function getFallbackTranslation(key: string): Promise<string> {
  try {
    const locale = await getLocale(getPreferredLocale())
    const theme = getPreferredTheme()
    const themeTexts = locale.themes?.[theme]
    const themeValue = getNestedValue(themeTexts, key)
    if (typeof themeValue === 'string') {
      return themeValue
    }

    const commonValue = getNestedValue(locale.common, key)
    if (typeof commonValue === 'string') {
      return commonValue
    }
  } catch (_error) {
    // ignore and fall through to returning key
  }

  return key
}

function getNestedValue(source: any, path: string): unknown {
  if (!source || typeof source !== 'object') {
    return undefined
  }

  return path.split('.').reduce<unknown>((current, segment) => {
    if (current && typeof current === 'object') {
      return (current as Record<string, unknown>)[segment]
    }
    return undefined
  }, source)
}

function applyParams(text: string, params?: Record<string, string | number>): string {
  if (!params) {
    return text
  }

  let result = text
  Object.entries(params).forEach(([key, value]) => {
    const placeholder = `{${key}}`
    result = result.split(placeholder).join(String(value))
  })
  return result
}

function getPreferredLocale(): SupportedLocale {
  if (typeof window === 'undefined') {
    return DEFAULT_LOCALE
  }

  const savedLocale = localStorage.getItem('locale')
  if (savedLocale && isLocaleSupported(savedLocale)) {
    return savedLocale as SupportedLocale
  }

  return DEFAULT_LOCALE
}

function getPreferredTheme(): TextTheme {
  if (typeof window === 'undefined') {
    return DEFAULT_THEME
  }

  const savedTheme = localStorage.getItem('text-theme')
  if (savedTheme === 'normal' || savedTheme === 'cyber') {
    return savedTheme as TextTheme
  }

  return DEFAULT_THEME
}

const IP_NOT_IN_WHITELIST = 6001
const IP_IN_BLACKLIST = 6002
const DOMAIN_NOT_IN_WHITELIST = 6004
const DOMAIN_IN_BLACKLIST = 6005
const USER_ACCOUNT_DISABLED = 1002
const IP_RESTRICTED_KEY = 'ip_restricted'
const IP_ADDRESS_KEY = 'ip_restricted_address'
const IP_ERROR_CODE_KEY = 'ip_restricted_code'
const IP_ERROR_MSG_KEY = 'ip_restricted_message'
const DOMAIN_KEY = 'domain_restricted_name'
const USER_DISABLED_KEY = 'user_disabled'

interface PendingRequest {
  url: string
  method: string
  params: string
  cancelToken: CancelTokenSource
  timestamp: number
}

const REQUEST_DEDUP_WHITELIST = [
  '/files/upload',
  '/files/guest/upload',
  '/user/personal/profile',
  '/files/chunked/upload',
  '/config/upload',
  '/settings',
  '/common/settings',
  '/folders/',
  '/admin/file/upload',
]

const isUrlInWhitelist = (url: string): boolean =>
  REQUEST_DEDUP_WHITELIST.some((whitelistUrl) => url === whitelistUrl || url.startsWith(whitelistUrl))

const pendingRequests = new Map<string, PendingRequest>()

const generateRequestKey = (config: AxiosRequestConfig): string => {
  const { method, url, params, data } = config
  const paramsStr = params ? JSON.stringify(params) : ''
  const dataStr = data ? JSON.stringify(data) : ''
  return `${method}:${url}:${paramsStr}:${dataStr}`
}

const addPendingRequest = (config: AxiosRequestConfig): void => {
  const requestKey = generateRequestKey(config)

  if (isUrlInWhitelist(config.url || '')) {
    const cancelToken = axios.CancelToken.source()
    config.cancelToken = cancelToken.token

    const uniqueKey = `${requestKey}_${Date.now()}_${Math.random()}`
    pendingRequests.set(uniqueKey, {
      url: config.url || '',
      method: config.method || '',
      params: generateRequestKey(config),
      cancelToken,
      timestamp: Date.now(),
    })
    return
  }

  const existingRequest = pendingRequests.get(requestKey)
  if (existingRequest) {
    getTranslation('utils.http.cancelDuplicate', {
      method: config.method?.toUpperCase() || '',
      url: config.url || '',
    }).then((msg) => existingRequest.cancelToken.cancel(msg))
    pendingRequests.delete(requestKey)
  }

  const cancelToken = axios.CancelToken.source()
  config.cancelToken = cancelToken.token

  pendingRequests.set(requestKey, {
    url: config.url || '',
    method: config.method || '',
    params: generateRequestKey(config),
    cancelToken,
    timestamp: Date.now(),
  })
}

const removePendingRequest = (config: AxiosRequestConfig): void => {
  const requestKey = generateRequestKey(config)

  if (isUrlInWhitelist(config.url || '')) {
    const keysToRemove: string[] = []
    for (const [key, request] of pendingRequests.entries()) {
      if (request.url === config.url && key.startsWith(requestKey)) {
        keysToRemove.push(key)
      }
    }
    keysToRemove.forEach((key) => pendingRequests.delete(key))
    return
  }

  if (pendingRequests.has(requestKey)) {
    pendingRequests.delete(requestKey)
  }
}

const cleanupTimeoutRequests = (): void => {
  const now = Date.now()
  const timeout = REQUEST_TIMEOUT.UPLOAD // 上传超时时间

  for (const [key, request] of pendingRequests.entries()) {
    if (now - request.timestamp > timeout) {
      getTranslation('utils.http.requestTimeout').then((msg) => request.cancelToken.cancel(msg))
      pendingRequests.delete(key)
    }
  }
}

setInterval(cleanupTimeoutRequests, REQUEST_TIMEOUT.DEFAULT) // 定期清理超时请求

let isRedirecting = false

export interface ExtendedAxiosRequestConfig extends AxiosRequestConfig, ExtendedRequestConfig {}

export const getApiBaseUrl = (): string => {
  try {
    if (typeof import.meta !== 'undefined' && import.meta.env) {
      return import.meta.env.VITE_API_BASE_URL || '/api/v1'
    }

    if (typeof window !== 'undefined' && (window as any).__VITE_API_BASE_URL__) {
      return (window as any).__VITE_API_BASE_URL__
    }
  } catch (_error) {
    console.warn('Unable to read environment variable VITE_API_BASE_URL:', error)
  }

  return '/api/v1'
}

const apiBaseUrl = getApiBaseUrl()

const instance: AxiosInstance = axios.create({
  baseURL: apiBaseUrl,
  timeout: REQUEST_TIMEOUT.UPLOAD, // 上传请求超时时间
  headers: {
    'Content-Type': 'application/json',
  },
})

function createApiResult<T = any>(success: boolean, code: number, message: string, data: T, request_id?: string): ApiResult<T> {
  return {
    success,
    code,
    message,
    data,
    request_id,
    timestamp: Date.now(),
  }
}

async function createSuccessResult<T = any>(data: T, message: string = '', request_id?: string): Promise<ApiSuccess<T>> {
  const msg = message || (await getTranslation('utils.http.operationSuccess'))
  return createApiResult(true, ErrorCodes.SUCCESS, msg, data, request_id) as ApiSuccess<T>
}

function createErrorResult(code: number, message: string, request_id?: string): ApiError {
  return createApiResult(false, code, message, null, request_id) as ApiError
}

instance.interceptors.request.use(
  (config: ExtendedAxiosRequestConfig) => {
    addPendingRequest(config)

    const token = StorageUtil.get<string>(TOKEN_KEY)
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    try {
      const textThemeStore = useTextThemeStore()
      const currentLocale = textThemeStore.currentLocale
      if (currentLocale) {
        config.headers['Accept-Language'] = currentLocale
      }
    } catch (_error) {
      config.headers['Accept-Language'] = 'zh-CN'
    }

    if (config.autoShowError === undefined) {
      config.autoShowError = true // 默认显示业务错误
    }
    if (config.silent === undefined) {
      config.silent = false // 默认不静默
    }
    if (config.showLoading === undefined) {
      config.showLoading = false // 默认不显示loading
    }
    if (config.autoShowSuccess === undefined) {
      config.autoShowSuccess = false // 默认不显示成功提示
    }
    if (config.throwOnError === undefined) {
      config.throwOnError = false // 默认不抛异常
    }
    if (config.useResultMode === undefined) {
      config.useResultMode = false // 默认使用异常模式
    }
    if (config.minLoadingTime === undefined) {
      config.minLoadingTime = 300 // 默认最小loading时间300ms
    }
    if (config.smartLoading === undefined) {
      config.smartLoading = true // 默认开启智能loading
    }

    if (config.smartLoading && !config.loadingTarget) {
      const detectedRefs = smartLoadingManager.detectLoadingRefs()
      if (detectedRefs.length > 0) {
        config._detectedLoadingRefs = detectedRefs
        smartLoadingManager.startLoading(detectedRefs, config.minLoadingTime)
      }
    }

    if (config.loadingTarget) {
      const loadingMode = config.loadingMode || 'auto'

      if (loadingMode === 'shared') {
        loadingManager.incrementLoading(config.loadingTarget as any, {
          minLoadingTime: config.minLoadingTime,
          groupId: config.loadingGroup,
        })
      } else if (loadingMode === 'auto') {
        config.loadingTarget.value = true
        config._startTime = Date.now()
      }
    }

    return config
  },
  (error) => Promise.reject(error)
)

const handleAccessRestriction = async (error: Error): Promise<boolean> => {
  if (!error.response || !error.response.data) {
    return false
  }

  const { code, message, ip, domain } = error.response.data
  let restrictionMsgKey = ''
  let addressValue = ''
  const unknownText = await getTranslation('utils.http.unknown')

  switch (code) {
    case IP_NOT_IN_WHITELIST:
      restrictionMsgKey = 'utils.http.ipNotInWhitelist'
      addressValue = ip || unknownText
      break
    case IP_IN_BLACKLIST:
      restrictionMsgKey = 'utils.http.ipInBlacklist'
      addressValue = ip || unknownText
      break
    case DOMAIN_NOT_IN_WHITELIST:
      restrictionMsgKey = 'utils.http.domainNotInWhitelist'
      addressValue = domain || unknownText
      StorageUtil.set(DOMAIN_KEY, domain || unknownText, 1)
      break
    case DOMAIN_IN_BLACKLIST:
      restrictionMsgKey = 'utils.http.domainInBlacklist'
      addressValue = domain || unknownText
      StorageUtil.set(DOMAIN_KEY, domain || unknownText, 1)
      break
    case USER_ACCOUNT_DISABLED:
      restrictionMsgKey = 'utils.http.userAccountDisabled'
      const restrictionMsg = await getTranslation(restrictionMsgKey)
      StorageUtil.set(USER_DISABLED_KEY, 'true', 1)
      StorageUtil.set(IP_ERROR_CODE_KEY, code.toString(), 1)
      StorageUtil.set(IP_ERROR_MSG_KEY, message || restrictionMsg, 1)
      if (typeof window !== 'undefined') {
        if (window.location.pathname !== '/refuse') {
          window.location.href = '/refuse'
        }
      }
      return true
    default:
      return false
  }

  const restrictionMsg = await getTranslation(restrictionMsgKey)

  StorageUtil.set(IP_RESTRICTED_KEY, 'true', 1)
  StorageUtil.set(IP_ADDRESS_KEY, addressValue, 1)
  StorageUtil.set(IP_ERROR_CODE_KEY, code.toString(), 1)
  StorageUtil.set(IP_ERROR_MSG_KEY, message || restrictionMsg, 1)

  if (typeof window !== 'undefined') {
    if (window.location.pathname !== '/refuse') {
      window.location.href = '/refuse'
    }
  }
  return true
}

import { loadingManager } from '../ui/loading-manager'
import { smartLoadingManager } from '../business/smart-loading'

const resetLoadingState = (config: ExtendedAxiosRequestConfig) => {
  if (config._detectedLoadingRefs) {
    smartLoadingManager.stopLoading(config._detectedLoadingRefs, config.minLoadingTime)
    return
  }

  if (!config.loadingTarget) {
    return
  }

  const loadingMode = config.loadingMode || 'auto'

  if (loadingMode === 'shared') {
    loadingManager.decrementLoading(config.loadingTarget as any, config.loadingGroup)
  } else if (loadingMode === 'auto') {
    if (config._startTime) {
      const elapsed = Date.now() - config._startTime
      const minTime = config.minLoadingTime || 300

      if (elapsed < minTime) {
        setTimeout(() => {
          if (config.loadingTarget) {
            config.loadingTarget.value = false
          }
        }, minTime - elapsed)
      } else {
        config.loadingTarget.value = false
      }
    }
  }
}

instance.interceptors.response.use(
  async (response: AxiosResponse): Promise<ApiResult<any>> => {
    removePendingRequest(response.config)

    resetLoadingState(response.config as ExtendedAxiosRequestConfig)

    const responseData: ApiResponse = response.data
    const config = response.config as ExtendedAxiosRequestConfig
    const requestId = responseData.request_id || response.headers?.['x-request-id']
    if (responseData.code === ErrorCodes.SUCCESS) {
      if (config.autoShowSuccess && responseData.message && !config.silent) {
        showToast.success(responseData.message)
      }
      return await createSuccessResult(responseData.data, responseData.message, requestId)
    } else if (responseData.code === ErrorCodes.SYSTEM_NOT_INSTALLED) {
      const currentPath = typeof window !== 'undefined' ? window.location.pathname : ''
      if (currentPath !== '/setup') {
        setTimeout(() => {
          if (window.location.pathname !== '/setup') {
            window.location.href = '/setup'
          }
        }, 0)
      }
      return await createSuccessResult(responseData.data, responseData.message, requestId)
    }

    if (responseData.code === USER_ACCOUNT_DISABLED) {
      const msg = responseData.message || (await getTranslation('utils.http.userAccountDisabled'))
      StorageUtil.set(USER_DISABLED_KEY, 'true', 1)
      StorageUtil.set(IP_ERROR_CODE_KEY, responseData.code.toString(), 1)
      StorageUtil.set(IP_ERROR_MSG_KEY, msg, 1)
      if (typeof window !== 'undefined') {
        if (window.location.pathname !== '/refuse') {
          window.location.href = '/refuse'
        }
      }
      const errorResult = createErrorResult(responseData.code, msg, requestId)
      return Promise.reject(errorResult)
    }

    const errorResult = createErrorResult(responseData.code, responseData.message, requestId)

    if (shouldShowError(responseData.code, config)) {
      const msg = responseData.message || (await getTranslation('utils.http.operationFailed'))
      showToast.error(msg)
    }

    if (config.useResultMode === true) {
      if (shouldThrowError(responseData.code, config)) {
        return Promise.reject(errorResult)
      }
      return errorResult
    }
    return Promise.reject(errorResult)
  },
  async (error): Promise<ApiError> => {
    if (!axios.isCancel(error) && error.config) {
      removePendingRequest(error.config)
      resetLoadingState(error.config as ExtendedAxiosRequestConfig)
    }

    const config = error.config as ExtendedAxiosRequestConfig

    if (axios.isCancel(error)) {
      const msg = await getTranslation('utils.http.requestCancelled')
      const cancelError = createErrorResult(-1, msg)
      return Promise.reject(cancelError)
    }

    let message = ''
    let errorCode = 999 // 默认网络错误码

    if (error.response) {
      errorCode = error.response.status

      switch (error.response.status) {
        case HTTP_STATUS.UNAUTHORIZED:
          message = await getTranslation('constants.api.errors.unauthorized')
          if (!config?.silent) {
            showToast.error(message)
          }
          StorageUtil.remove(TOKEN_KEY)
          const currentPath = window.location.pathname
          if (!isRedirecting && currentPath !== '/auth' && currentPath !== '/login') {
            isRedirecting = true
            import('@/store/auth').then((module) => {
              const authStore = module.useAuthStore()
              authStore.logout()
              setTimeout(() => {
                isRedirecting = false
              }, 1000)
            })
          }
          break
        case HTTP_STATUS.FORBIDDEN:
          message = await getTranslation('constants.api.errors.forbidden')
          break
        case HTTP_STATUS.NOT_FOUND:
          message = await getTranslation('constants.api.errors.notFound')
          break
        case 451:
          message = await getTranslation('constants.api.errors.unavailableForLegalReasons')
          break
        case HTTP_STATUS.INTERNAL_SERVER_ERROR:
          message = await getTranslation('constants.api.errors.internalServerError')
          break
        case HTTP_STATUS.BAD_GATEWAY:
          message = await getTranslation('constants.api.errors.badGateway')
          break
        case HTTP_STATUS.SERVICE_UNAVAILABLE:
          message = await getTranslation('constants.api.errors.serviceUnavailable')
          break
        case HTTP_STATUS.GATEWAY_TIMEOUT:
          message = await getTranslation('constants.api.errors.gatewayTimeout')
          break
        default:
          message = await getTranslation('utils.http.requestFailedWithStatus', { status: error.response.status })
      }

      if (await handleAccessRestriction(error)) {
        const accessError = createErrorResult(
          error.response.data.code,
          error.response.data.message,
          error.response.data.request_id
        )
        return Promise.reject(accessError)
      }

      if (error.response.data && error.response.data.code && error.response.data.message) {
        const serverError = error.response.data
        const errorResult = createErrorResult(serverError.code, serverError.message, serverError.request_id)

        if (shouldShowError(serverError.code, config)) {
          showToast.error(serverError.message)
        }

        return Promise.reject(errorResult)
      }
    } else if (error.request) {
      message = await getTranslation('utils.http.networkError')
    } else {
      message = await getTranslation('utils.http.configError')
    }

    if (!config?.silent) {
      showToast.error(message)
    }

    const errorResult = createErrorResult(errorCode, message, error.response?.headers?.['x-request-id'] || '')

    const category = getErrorCategory(errorCode)
    console.error(`[HTTP ${category.toUpperCase()}]`, {
      url: config?.url || 'unknown',
      method: config?.method || 'unknown',
      _error: error,
      originalError: error,
    })

    return Promise.reject(errorResult)
  }
)

export function get<T = any>(url: string, params?: any, config?: ExtendedAxiosRequestConfig): Promise<ApiResult<T>> {
  return instance.get(url, { params, ...config })
}

export function post<T = any>(url: string, data?: any, config?: ExtendedAxiosRequestConfig): Promise<ApiResult<T>> {
  return instance.post(url, data, config)
}

export function put<T = any>(url: string, data?: any, config?: ExtendedAxiosRequestConfig): Promise<ApiResult<T>> {
  return instance.put(url, data, config)
}

export function patch<T = any>(url: string, data?: any, config?: ExtendedAxiosRequestConfig): Promise<ApiResult<T>> {
  return instance.patch(url, data, config)
}

export function del<T = any>(url: string, params?: any, config?: ExtendedAxiosRequestConfig): Promise<ApiResult<T>> {
  return instance.delete(url, { params, ...config })
}

export function upload<T = any>(
  url: string,
  file: File,
  data?: Record<string, any>,
  onUploadProgress?: (progressEvent: any) => void,
  config?: ExtendedAxiosRequestConfig
): Promise<ApiResult<T>> {
  const formData = new FormData()
  formData.append('file', file)

  if (data) {
    Object.keys(data).forEach((key) => {
      const value = data[key]
      if (value !== null && value !== undefined) {
        if (typeof value === 'boolean') {
          formData.append(key, value.toString())
        } else if (typeof value === 'object' && value !== null) {
          if ('value' in value) {
            const actualValue = value.value
            formData.append(key, actualValue !== null && actualValue !== undefined ? String(actualValue) : '')
          } else {
            formData.append(key, String(value))
          }
        } else {
          formData.append(key, String(value))
        }
      } else {
        formData.append(key, '')
      }
    })
  }

  return instance.post(url, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress,
    ...config,
  })
}

export function uploadBatch<T = any>(
  url: string,
  files: File[],
  data?: Record<string, any>,
  onUploadProgress?: (progressEvent: any) => void,
  config?: ExtendedAxiosRequestConfig
): Promise<ApiResult<T>> {
  const formData = new FormData()

  files.forEach((file, _index) => {
    formData.append(`files[]`, file)
  })

  if (data) {
    Object.keys(data).forEach((key) => {
      const value = data[key]
      if (value !== null && value !== undefined) {
        if (typeof value === 'boolean') {
          formData.append(key, value.toString())
        } else if (typeof value === 'object' && value !== null) {
          if ('value' in value) {
            const actualValue = value.value
            formData.append(key, actualValue !== null && actualValue !== undefined ? String(actualValue) : '')
          } else {
            formData.append(key, String(value))
          }
        } else {
          formData.append(key, String(value))
        }
      } else {
        formData.append(key, '')
      }
    })
  }

  return instance.post(url, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress,
    ...config,
  })
}

export async function cancelAllPendingRequests(reason?: string): Promise<void> {
  const msg = reason || (await getTranslation('utils.http.pageSwitch'))
  for (const [_key, request] of pendingRequests.entries()) {
    request.cancelToken.cancel(msg)
  }

  pendingRequests.clear()
}

export async function cancelRequestsByUrl(url: string, reason?: string): Promise<void> {
  const msg = reason || (await getTranslation('utils.http.cancelSpecificUrl'))
  const keysToDelete: string[] = []

  for (const [key, request] of pendingRequests.entries()) {
    if (request.url.includes(url)) {
      request.cancelToken.cancel(msg)
      keysToDelete.push(key)
    }
  }

  keysToDelete.forEach((key) => pendingRequests.delete(key))
}

export function getPendingRequestsCount(): number {
  return pendingRequests.size
}

export function addToRequestWhitelist(url: string): void {
  if (!REQUEST_DEDUP_WHITELIST.includes(url)) {
    REQUEST_DEDUP_WHITELIST.push(url)
  }
}

export function removeFromRequestWhitelist(url: string): void {
  const index = REQUEST_DEDUP_WHITELIST.indexOf(url)
  if (index > -1) {
    REQUEST_DEDUP_WHITELIST.splice(index, 1)
  }
}

export function getRequestWhitelist(): string[] {
  return [...REQUEST_DEDUP_WHITELIST]
}

export const AccessRestrictionCodes = {
  IP_NOT_IN_WHITELIST,
  IP_IN_BLACKLIST,
  DOMAIN_NOT_IN_WHITELIST,
  DOMAIN_IN_BLACKLIST,
  USER_ACCOUNT_DISABLED,
}

export const AccessRestrictionKeys = {
  IP_RESTRICTED_KEY,
  IP_ADDRESS_KEY,
  IP_ERROR_CODE_KEY,
  IP_ERROR_MSG_KEY,
  DOMAIN_KEY,
  USER_DISABLED_KEY,
}

export default {
  get,
  post,
  put,
  patch,
  delete: del,
  upload,
  uploadBatch,
  cancelAllPendingRequests,
  cancelRequestsByUrl,
  getPendingRequestsCount,
  addToRequestWhitelist,
  removeFromRequestWhitelist,
  getRequestWhitelist,
}
