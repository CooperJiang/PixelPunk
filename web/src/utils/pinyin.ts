import { pinyin } from 'pinyin-pro'

/**
 * 获取字符串的首字母（用于字母分组）
 * @param str 输入字符串
 * @returns 大写字母 A-Z 或 #（特殊字符）
 */
export function getFirstLetter(str: string): string {
  if (!str || str.length === 0) {
    return '#'
  }

  const firstChar = str.charAt(0)

  if (/[a-zA-Z]/.test(firstChar)) {
    return firstChar.toUpperCase()
  }

  if (/[\u4e00-\u9fa5]/.test(firstChar)) {
    const pinyinStr = pinyin(firstChar, {
      pattern: 'first', // 只取首字母
      toneType: 'none', // 不要声调
    })
    return pinyinStr.toUpperCase()
  }

  return '#'
}

export function getFullPinyin(str: string): string {
  if (!str) {
    return ''
  }

  return pinyin(str, {
    toneType: 'none',
    type: 'array',
  }).join(' ')
}

export function groupTagsByLetter<T extends { name: string; count?: number }>(
  tags: T[],
  sortBy: 'count' | 'name' = 'count'
): Record<string, T[]> {
  const groups: Record<string, T[]> = {}

  const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ#'.split('')
  letters.forEach((letter) => {
    groups[letter] = []
  })

  tags.forEach((tag) => {
    const firstLetter = getFirstLetter(tag.name)
    if (groups[firstLetter]) {
      groups[firstLetter].push(tag)
    } else {
      groups['#'].push(tag)
    }
  })

  Object.keys(groups).forEach((letter) => {
    if (sortBy === 'count') {
      groups[letter].sort((a, b) => (b.count || 0) - (a.count || 0))
    } else {
      groups[letter].sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'))
    }
  })

  return groups
}

export function getAvailableLetters(groupedTags: Record<string, any[]>): string[] {
  return Object.keys(groupedTags).filter((letter) => groupedTags[letter].length > 0)
}

export function getGroupStats(groupedTags: Record<string, any[]>): {
  totalGroups: number
  totalTags: number
  averagePerGroup: number
} {
  const availableLetters = getAvailableLetters(groupedTags)
  const totalGroups = availableLetters.length
  const totalTags = availableLetters.reduce((sum, letter) => sum + groupedTags[letter].length, 0)
  const averagePerGroup = totalGroups > 0 ? Math.round(totalTags / totalGroups) : 0

  return {
    totalGroups,
    totalTags,
    averagePerGroup,
  }
}
