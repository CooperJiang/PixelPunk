/**
 * API相关常量定义
 *
 * 注意: API 端点采用就近原则管理,直接在各个 API 模块中定义
 * 此文件只保留全局技术配置常量
 */

import type { Composer } from '@/composables/useTexts'

/* ===== HTTP方法常量 ===== */
export const HTTP_METHODS = {
  GET: 'GET',
  POST: 'POST',
  PUT: 'PUT',
  DELETE: 'DELETE',
  PATCH: 'PATCH',
} as const

/* ===== HTTP状态码常量 ===== */
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,

  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  METHOD_NOT_ALLOWED: 405,
  CONFLICT: 409,
  PAYLOAD_TOO_LARGE: 413,
  UNPROCESSABLE_ENTITY: 422,
  TOO_MANY_REQUESTS: 429,

  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
  GATEWAY_TIMEOUT: 504,
} as const

/* ===== 错误码映射 ===== */

/**
 * 获取错误码消息映射
 */
export const getErrorCodeMessages = ($t: Composer['$t']) => ({
  [HTTP_STATUS.BAD_REQUEST]: $t('constants.api.errors.badRequest'),
  [HTTP_STATUS.UNAUTHORIZED]: $t('constants.api.errors.unauthorized'),
  [HTTP_STATUS.FORBIDDEN]: $t('constants.api.errors.forbidden'),
  [HTTP_STATUS.NOT_FOUND]: $t('constants.api.errors.notFound'),
  [HTTP_STATUS.METHOD_NOT_ALLOWED]: $t('constants.api.errors.methodNotAllowed'),
  [HTTP_STATUS.CONFLICT]: $t('constants.api.errors.conflict'),
  [HTTP_STATUS.PAYLOAD_TOO_LARGE]: $t('constants.api.errors.payloadTooLarge'),
  [HTTP_STATUS.UNPROCESSABLE_ENTITY]: $t('constants.api.errors.unprocessableEntity'),
  [HTTP_STATUS.TOO_MANY_REQUESTS]: $t('constants.api.errors.tooManyRequests'),
  [HTTP_STATUS.INTERNAL_SERVER_ERROR]: $t('constants.api.errors.internalServerError'),
  [HTTP_STATUS.BAD_GATEWAY]: $t('constants.api.errors.badGateway'),
  [HTTP_STATUS.SERVICE_UNAVAILABLE]: $t('constants.api.errors.serviceUnavailable'),
  [HTTP_STATUS.GATEWAY_TIMEOUT]: $t('constants.api.errors.gatewayTimeout'),
})

/* ===== 请求超时配置 ===== */
export const REQUEST_TIMEOUT = {
  DEFAULT: 5000, // 5秒 - 默认超时时间
  UPLOAD: 60000, // 60秒 - 文件上传超时
  DOWNLOAD: 30000, // 30秒 - 文件下载超时
  AUTH: 8000, // 8秒 - 认证请求超时
  SEARCH: 10000, // 10秒 - 搜索请求超时
} as const

/* ===== 重试配置 ===== */
export const RETRY_CONFIG = {
  MAX_ATTEMPTS: 3, // 最多重试次数
  INITIAL_DELAY: 1000, // 初始重试延迟 (ms)
  MAX_DELAY: 10000, // 最大重试延迟 (ms)
  BACKOFF_FACTOR: 2, // 退避因子

  RETRYABLE_STATUS: [
    HTTP_STATUS.INTERNAL_SERVER_ERROR,
    HTTP_STATUS.BAD_GATEWAY,
    HTTP_STATUS.SERVICE_UNAVAILABLE,
    HTTP_STATUS.GATEWAY_TIMEOUT,
  ],
} as const

/* ===== 内容类型常量 ===== */
export const CONTENT_TYPES = {
  JSON: 'application/json',
  FORM_DATA: 'multipart/form-data',
  URL_ENCODED: 'application/x-www-form-urlencoded',
  TEXT_PLAIN: 'text/plain',
  OCTET_STREAM: 'application/octet-stream',
} as const

/* ===== 缓存配置 ===== */
export const CACHE_CONFIG = {
  PREFIX: {
    USER: 'user:',
    IMAGE: 'image:',
    FOLDER: 'folder:',
    SHARE: 'share:',
    SETTINGS: 'settings:',
  },

  TTL: {
    SHORT: 300, // 5分钟
    MEDIUM: 1800, // 30分钟
    LONG: 3600, // 1小时
    VERY_LONG: 86400, // 24小时
  },

  STRATEGY: {
    CACHE_FIRST: 'cache-first',
    NETWORK_FIRST: 'network-first',
    CACHE_ONLY: 'cache-only',
    NETWORK_ONLY: 'network-only',
  },
} as const

/* ===== 分页配置 ===== */
export const PAGINATION_CONFIG = {
  DEFAULT_PAGE: 1,
  DEFAULT_SIZE: 20,
  MAX_SIZE: 100,
  SIZE_OPTIONS: [10, 20, 50, 100],
} as const

/* ===== 上传配置 ===== */
export const UPLOAD_CONFIG = {
  CHUNK_SIZE: 2 * 1024 * 1024, // 2MB 分片大小
  MAX_CONCURRENT_CHUNKS: 5, // 最大并发分片数

  PROGRESS_UPDATE_INTERVAL: 100, // 100ms 更新一次进度

  MAX_QUEUE_SIZE: 100, // 最大队列大小
  AUTO_START: true, // 自动开始上传
} as const

/* ===== 搜索配置 ===== */
export const SEARCH_CONFIG = {
  MIN_QUERY_LENGTH: 2, // 最小搜索长度
  DEBOUNCE_DELAY: 300, // 防抖延迟 (ms)
  MAX_SUGGESTIONS: 10, // 最大建议数量
  HIGHLIGHT_TAG: 'mark', // 高亮标签
} as const

/* ===== WebSocket配置 ===== */
export const WEBSOCKET_CONFIG = {
  RECONNECT_INTERVAL: 5000, // 重连间隔 (ms)
  MAX_RECONNECT_ATTEMPTS: 5, // 最大重连次数
  HEARTBEAT_INTERVAL: 30000, // 心跳间隔 (ms)

  MESSAGE_TYPES: {
    PING: 'ping',
    PONG: 'pong',
    UPLOAD_PROGRESS: 'upload_progress',
    UPLOAD_COMPLETE: 'upload_complete',
    NOTIFICATION: 'notification',
  },
} as const

/* ===== 类型定义 ===== */
export type HttpMethod = (typeof HTTP_METHODS)[keyof typeof HTTP_METHODS]
export type HttpStatus = (typeof HTTP_STATUS)[keyof typeof HTTP_STATUS]
export type ContentType = (typeof CONTENT_TYPES)[keyof typeof CONTENT_TYPES]
export type CacheStrategy = (typeof CACHE_CONFIG.STRATEGY)[keyof typeof CACHE_CONFIG.STRATEGY]
export type WebSocketMessageType = (typeof WEBSOCKET_CONFIG.MESSAGE_TYPES)[keyof typeof WEBSOCKET_CONFIG.MESSAGE_TYPES]
