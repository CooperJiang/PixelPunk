/**
 * Canvas 主题颜色工具
 * 从 CSS 变量中读取主题颜色，用于 Canvas 绘制

 */
export interface CanvasThemeColors {
  brand: {
    primary: string
    light: string
    lighter: string
  }
  success: string
  danger: string
  white: string
  border: string
  background: string
}

/**
 * 从 CSS 变量获取主题颜色（用于 Canvas 绘制）
 */
export function getCanvasThemeColors(): CanvasThemeColors {
  const root = document.documentElement
  const computedStyle = getComputedStyle(root)

  const getCssVar = (varName: string): string => {
    return computedStyle.getPropertyValue(varName).trim()
  }

  const rgbToRgba = (rgb: string, alpha: number = 1): string => {
    if (rgb.startsWith('rgb')) return rgb

    const values = rgb.split(',').map((v) => v.trim())
    if (values.length === 3) {
      return `rgba(${values[0]}, ${values[1]}, ${values[2]}, ${alpha})`
    }

    return `rgba(139, 92, 246, ${alpha})`
  }

  const brandRgb = getCssVar('--color-brand-500-rgb') || '139, 92, 246'

  return {
    brand: {
      primary: rgbToRgba(brandRgb, 1),
      light: rgbToRgba(brandRgb, 0.8),
      lighter: rgbToRgba(brandRgb, 0.4),
    },
    success: getCssVar('--color-success-500') || 'rgb(34, 197, 94)',
    danger: getCssVar('--color-danger-500') || 'rgb(239, 68, 68)',
    white: '#ffffff',
    border: rgbToRgba(brandRgb, 0.1),
    background: rgbToRgba(brandRgb, 0.05),
  }
}

export function getBrandColor(alpha: number = 1): string {
  const root = document.documentElement
  const computedStyle = getComputedStyle(root)
  const brandRgb = computedStyle.getPropertyValue('--color-brand-500-rgb').trim() || '139, 92, 246'

  const values = brandRgb.split(',').map((v) => v.trim())
  if (values.length === 3) {
    return `rgba(${values[0]}, ${values[1]}, ${values[2]}, ${alpha})`
  }

  return `rgba(139, 92, 246, ${alpha})`
}

export function getSuccessColor(alpha: number = 1): string {
  const root = document.documentElement
  const computedStyle = getComputedStyle(root)
  const successRgb = computedStyle.getPropertyValue('--color-success-500-rgb').trim() || '34, 197, 94'

  const values = successRgb.split(',').map((v) => v.trim())
  if (values.length === 3) {
    return `rgba(${values[0]}, ${values[1]}, ${values[2]}, ${alpha})`
  }

  return `rgba(34, 197, 94, ${alpha})`
}

export function getDangerColor(alpha: number = 1): string {
  const root = document.documentElement
  const computedStyle = getComputedStyle(root)
  const dangerRgb = computedStyle.getPropertyValue('--color-danger-500-rgb').trim() || '239, 68, 68'

  const values = dangerRgb.split(',').map((v) => v.trim())
  if (values.length === 3) {
    return `rgba(${values[0]}, ${values[1]}, ${values[2]}, ${alpha})`
  }

  return `rgba(239, 68, 68, ${alpha})`
}
