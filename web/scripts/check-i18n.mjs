#!/usr/bin/env node
/**
 * 多语言完整性检查工具
 *
 * 使用方法: node scripts/check-i18n.mjs
 */

import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const LOCALES_DIR = path.resolve(__dirname, '../src/locales')
const LANGUAGES = ['zh-CN', 'en-US', 'ja-JP']
const THEMES = ['normal', 'cyber']

/**
 * 递归获取对象的所有key路径
 */
function getObjectKeys(obj, prefix = '') {
  const keys = []

  for (const key in obj) {
    if (obj[key] && typeof obj[key] === 'object' && !Array.isArray(obj[key])) {
      keys.push(...getObjectKeys(obj[key], prefix ? `${prefix}.${key}` : key))
    } else {
      keys.push(prefix ? `${prefix}.${key}` : key)
    }
  }

  return keys.sort()
}

/**
 * 读取并解析TypeScript文件中的导出对象
 */
function parseExportedObject(filePath) {
  try {
    const content = fs.readFileSync(filePath, 'utf-8')

    // 简单的正则匹配导出的对象（这不是完美的解析，但对于我们的用例足够了）
    const exportMatch = content.match(/export\s+const\s+\w+\s*=\s*(\{[\s\S]*\})\s*$/m)
    if (!exportMatch) {
      return null
    }

    // 使用 eval 来解析对象（注意：这只适用于可信的代码）
    // 在生产环境中应该使用更安全的方法
    try {
      const objStr = exportMatch[1]
      const obj = eval(`(${objStr})`)
      return obj
    } catch (e) {
      console.warn(`⚠️  无法解析文件: ${filePath}`)
      return null
    }
  } catch (error) {
    console.error(`❌ 读取文件失败: ${filePath}`, error.message)
    return null
  }
}

/**
 * 递归读取目录下的所有.ts文件
 */
function getAllTsFiles(dir, fileList = []) {
  const files = fs.readdirSync(dir)

  files.forEach((file) => {
    const filePath = path.join(dir, file)
    const stat = fs.statSync(filePath)

    if (stat.isDirectory()) {
      getAllTsFiles(filePath, fileList)
    } else if (file.endsWith('.ts') && !file.endsWith('.d.ts')) {
      fileList.push(filePath)
    }
  })

  return fileList
}

/**
 * 收集某个主题下所有的翻译key
 */
function collectKeysFromTheme(themePath) {
  const allKeys = new Set()
  const files = getAllTsFiles(themePath)

  files.forEach((file) => {
    const obj = parseExportedObject(file)
    if (obj) {
      const keys = getObjectKeys(obj)
      keys.forEach((k) => allKeys.add(k))
    }
  })

  return Array.from(allKeys).sort()
}

/**
 * 比较两个key数组
 */
function compareKeys(keys1, keys2) {
  const set1 = new Set(keys1)
  const set2 = new Set(keys2)

  const missing = [...set2].filter((k) => !set1.has(k))
  const extra = [...set1].filter((k) => !set2.has(k))

  return { missing, extra }
}

/**
 * 主检查函数
 */
function checkI18nCompleteness() {
  console.log('🔍 开始检查多语言配置完整性...\n')

  const allKeys = {}

  // 收集所有语言和主题的keys
  for (const lang of LANGUAGES) {
    allKeys[lang] = {}

    for (const theme of THEMES) {
      const themePath = path.join(LOCALES_DIR, lang, 'themes', theme)
      if (fs.existsSync(themePath)) {
        allKeys[lang][theme] = collectKeysFromTheme(themePath)
        console.log(`✓ 已加载 ${lang}/${theme}: ${allKeys[lang][theme].length} 个key`)
      } else {
        console.warn(`⚠️  路径不存在: ${themePath}`)
      }
    }
  }

  console.log('\n' + '='.repeat(60) + '\n')

  let hasIssues = false

  // 检查1: 同一语言下，两个主题的key是否一致
  console.log('📋 检查 1: 同一语言下，normal 和 cyber 主题的一致性\n')
  for (const lang of LANGUAGES) {
    if (allKeys[lang].normal && allKeys[lang].cyber) {
      const { missing, extra } = compareKeys(allKeys[lang].cyber, allKeys[lang].normal)

      if (missing.length > 0 || extra.length > 0) {
        hasIssues = true
        console.log(`⚠️  ${lang}: normal 和 cyber 主题不一致`)
        if (missing.length > 0) {
          console.log(`   Cyber主题缺失 (${missing.length}个):`)
          missing.slice(0, 5).forEach((k) => console.log(`     - ${k}`))
          if (missing.length > 5) console.log(`     ... 还有 ${missing.length - 5} 个`)
        }
        if (extra.length > 0) {
          console.log(`   Cyber主题多余 (${extra.length}个):`)
          extra.slice(0, 5).forEach((k) => console.log(`     - ${k}`))
          if (extra.length > 5) console.log(`     ... 还有 ${extra.length - 5} 个`)
        }
        console.log()
      } else {
        console.log(`✅ ${lang}: 主题一致`)
      }
    }
  }

  // 检查2: 不同语言下，相同主题的key是否一致
  console.log('\n📋 检查 2: 不同语言间的一致性 (以 zh-CN 为基准)\n')
  for (const theme of THEMES) {
    console.log(`${theme} 主题:`)

    const baseKeys = allKeys['zh-CN']?.[theme] || []

    for (const lang of LANGUAGES) {
      if (lang === 'zh-CN') continue

      const currentKeys = allKeys[lang]?.[theme] || []
      const { missing, extra } = compareKeys(currentKeys, baseKeys)

      if (missing.length > 0 || extra.length > 0) {
        hasIssues = true
        console.log(`  ⚠️  ${lang} 与 zh-CN 不一致`)
        if (missing.length > 0) {
          console.log(`     缺失 (${missing.length}个):`)
          missing.slice(0, 5).forEach((k) => console.log(`       - ${k}`))
          if (missing.length > 5) console.log(`       ... 还有 ${missing.length - 5} 个`)
        }
        if (extra.length > 0) {
          console.log(`     多余 (${extra.length}个):`)
          extra.slice(0, 5).forEach((k) => console.log(`       - ${k}`))
          if (extra.length > 5) console.log(`       ... 还有 ${extra.length - 5} 个`)
        }
      } else {
        console.log(`  ✅ ${lang}`)
      }
    }
    console.log()
  }

  // 统计信息
  console.log('📊 统计信息:\n')
  for (const lang of LANGUAGES) {
    for (const theme of THEMES) {
      const count = allKeys[lang]?.[theme]?.length || 0
      console.log(`  ${lang.padEnd(8)} ${theme.padEnd(8)} ${count} 个key`)
    }
  }

  console.log('\n' + '='.repeat(60) + '\n')

  if (!hasIssues) {
    console.log('✨ 太棒了！所有多语言配置都是完整且一致的！\n')
    return true
  } else {
    console.log('⚠️  发现了一些不一致的地方，请检查并修复。\n')
    return false
  }
}

// 运行检查
const success = checkI18nCompleteness()
process.exit(success ? 0 : 1)
