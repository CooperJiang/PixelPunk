/**
 * 环境变量和全局类型声明
 */
import type { PerformanceConfig } from './performance'

/* 环境变量 */
declare global {
  interface ImportMeta {
    readonly env: ImportMetaEnv
  }

  interface ImportMetaEnv {
    readonly VITE_API_BASE_URL: string
    readonly VITE_SITE_DOMAIN: string
    readonly VITE_ENABLE_ANALYTICS: string
    readonly VITE_ENABLE_ERROR_REPORT: string
    readonly VITE_ENABLE_PERFORMANCE_MONITOR: string
    readonly VITE_ENABLE_DEBUG: string
    readonly VITE_ERROR_REPORT_URL: string
    readonly VITE_ANALYTICS_ID: string
    readonly VITE_SENTRY_DSN: string
    readonly VITE_API_TIMEOUT: string
    readonly VITE_MAX_FILE_SIZE: string
    readonly VITE_CHUNK_SIZE: string
    readonly VITE_CACHE_TTL: string
    readonly VITE_MAX_CACHE_SIZE: string
    readonly VITE_TOKEN_EXPIRE_TIME: string
    readonly VITE_REFRESH_TOKEN_THRESHOLD: string
    readonly VITE_DEFAULT_THEME: string
    readonly VITE_DEFAULT_LANGUAGE: string
    readonly VITE_PAGE_SIZE: string
    readonly VITE_MOCK_API: string
    readonly VITE_HOT_RELOAD: string
    readonly VITE_BUILD_TIME: string
    readonly VITE_APP_VERSION: string
    readonly VITE_GIT_COMMIT: string
    readonly VITE_LOG_ENDPOINT: string
  }

  /* Window 扩展 */
  interface Window {
    __APP_CONFIG__?: any
    __INITIAL_STATE__?: any
    __PERFORMANCE_OBSERVER__?: PerformanceObserver
  }
}

export {}
