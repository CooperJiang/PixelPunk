/**
 * Japanese language pack entry
 */
import { normalTheme } from './themes/normal'
import { cyberTheme } from './themes/cyber'
import { commonTexts } from './common'

export const jaJP = {
  themes: {
    normal: normalTheme,
    cyber: cyberTheme,
  },
  common: commonTexts,
}

/* Export types (for TypeScript type inference) */
export type JaJPLocale = typeof jaJP
export type ThemeTexts = typeof normalTheme
export type CommonTexts = typeof commonTexts
export type TextTheme = keyof typeof jaJP.themes
export type TextKey = keyof ThemeTexts
