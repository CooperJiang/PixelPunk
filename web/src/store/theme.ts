import { defineStore } from 'pinia'

/**
 * 网站模式 Store (已简化为固定 website 模式)
 * 注意：此 Store 已简化，仅保留必要的接口以保持向后兼容
 * 实际主题管理请使用 visualTheme.ts 和 textTheme.ts
 * 推荐使用统一入口：useTheme() composable
 */
export const useThemeStore = defineStore('theme', () => {
  const getSiteMode = () => 'website' as const

  const isWebsiteMode = () => true

  return {
    getSiteMode,
    isWebsiteMode,
  }
})
