/**
 * 区域设置工具
 * 提供统一的 locale 管理
 */

/**
 * 获取当前应用的区域设置
 * @returns 当前区域设置代码 (例如: 'zh-CN', 'en-US')
 * 注意: 目前返回固定的 'zh-CN'，
 * 后续如需支持多语言，可以从 store 或其他配置中读取
 */
export function getCurrentLocale(): string {
  return 'zh-CN'
}

export const getLocale = getCurrentLocale
