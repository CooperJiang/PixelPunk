/**
 * ID生成工具函数
 */

/**
 * 生成唯一ID
 * @param prefix 可选前缀
 * @returns 唯一ID
 */
export function generateUniqueId(prefix: string = ''): string {
  const timestamp = new Date().getTime()
  const randomPart = Math.floor(Math.random() * 10000)
    .toString()
    .padStart(4, '0')
  return `${prefix}${timestamp}-${randomPart}`
}
