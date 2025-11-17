/**
 * 数组处理工具函数
 */

/**
 * @param array 原数组
 * @param keyFn 分组键函数
 * @returns 分组结果
 */
export function groupBy<T, K extends string | number | symbol>(array: T[], keyFn: (item: T) => K): Record<K, T[]> {
  return array.reduce(
    (result, item) => {
      const key = keyFn(item)
      result[key] = result[key] || []
      result[key].push(item)
      return result
    },
    {} as Record<K, T[]>
  )
}
