/**
 * 中文语言包入口
 */
import { normalTheme } from './themes/normal'
import { cyberTheme } from './themes/cyber'
import { commonTexts } from './common'

export const zhCN = {
  themes: {
    normal: normalTheme,
    cyber: cyberTheme,
  },

  common: commonTexts,
}

/* 导出类型（用于TypeScript类型推断） */
export type ZhCNLocale = typeof zhCN
export type ThemeTexts = typeof normalTheme
export type CommonTexts = typeof commonTexts
export type TextTheme = keyof typeof zhCN.themes
export type TextKey = keyof ThemeTexts
