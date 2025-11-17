/**
 * 环境变量配置
 */

export const isDevelopment = import.meta.env.MODE === 'development'
export const isProduction = import.meta.env.MODE === 'production'

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

export const DEBUG_CONFIG = {
  enableNetworkLogging: isDevelopment,
} as const
