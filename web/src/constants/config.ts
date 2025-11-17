/**
 * 全局配置文件
 * @description 包含应用程序的基本配置信息

/* 应用配置 */
export const APP_CONFIG = {
  name: 'PIXELPUNK',
  description: '文件分享与管理平台',

  brand: {
    name: 'PIXELPUNK',
    pixel: 'PIXEL',
    punk: 'PUNK',
    logo: '/logo.png',
    favicon: '/favicon.ico',
  },
} as const

/* 配置类型定义 */
export type AppConfig = typeof APP_CONFIG
