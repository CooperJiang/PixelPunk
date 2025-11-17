/**
 * English language pack entry
 */
import { normalTheme } from './themes/normal'
import { cyberTheme } from './themes/cyber'
import { commonTexts } from './common'

export const enUS = {
  themes: {
    normal: normalTheme,
    cyber: cyberTheme,
  },
  common: commonTexts,
}

/* Export types (for TypeScript type inference) */
export type EnUSLocale = typeof enUS
export type ThemeTexts = typeof normalTheme
export type CommonTexts = typeof commonTexts
export type TextTheme = keyof typeof enUS.themes
export type TextKey = keyof ThemeTexts
