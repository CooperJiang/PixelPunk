#!/usr/bin/env node
/**
 * 多语言完整性检查工具
 *
 * 功能：
 * 1. 检查所有语言和主题的翻译key是否一致
 * 2. 找出缺失的翻译
 * 3. 找出多余的翻译
 * 4. 生成检查报告
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
    if (typeof obj[key] === 'object' && obj[key] !== null && !Array.isArray(obj[key])) {
      keys.push(...getObjectKeys(obj[key], prefix ? `${prefix}.${key}` : key))
    } else {
      keys.push(prefix ? `${prefix}.${key}` : key)
    }
  }

  return keys
}

/**
 * 动态导入语言包
 */
async function loadLocaleData(language, theme) {
  try {
    const modulePath = `../src/locales/${language}/themes/${theme}/index.ts`
    const module = await import(modulePath)
    return module[theme === 'normal' ? 'normalTheme' : 'cyberTheme']
  } catch (error) {
    console.error(`❌ 无法加载 ${language}/${theme}:`, error.message)
    return null
  }
}

/**
 * 比较两个key数组，找出差异
 */
function compareKeys(keys1, keys2, label1, label2) {
  const set1 = new Set(keys1)
  const set2 = new Set(keys2)

  const missing = [...set2].filter((k) => !set1.has(k))
  const extra = [...set1].filter((k) => !set2.has(k))

  return { missing, extra }
}

/**
 * 主检查函数
 */
async function checkI18nCompleteness() {
  console.log('🔍 开始检查多语言配置完整性...\n')

  const allData = {}
  const allKeys = {}

  // 加载所有语言和主题的数据
  for (const lang of LANGUAGES) {
    allData[lang] = {}
    allKeys[lang] = {}

    for (const theme of THEMES) {
      const data = await loadLocaleData(lang, theme)
      if (data) {
        allData[lang][theme] = data
        allKeys[lang][theme] = getObjectKeys(data).sort()
      }
    }
  }

  let hasIssues = false

  // 1. 检查同一语言下，两个主题的key是否一致
  console.log('📋 检查 1: 同一语言下，normal 和 cyber 主题的key一致性\n')
  for (const lang of LANGUAGES) {
    if (allKeys[lang].normal && allKeys[lang].cyber) {
      const { missing, extra } = compareKeys(allKeys[lang].cyber, allKeys[lang].normal, 'cyber', 'normal')

      if (missing.length > 0 || extra.length > 0) {
        hasIssues = true
        console.log(`⚠️  ${lang}: normal 和 cyber 主题不一致`)
        if (missing.length > 0) {
          console.log(`   Cyber主题缺失的key (${missing.length}个):`)
          missing.forEach((k) => console.log(`     - ${k}`))
        }
        if (extra.length > 0) {
          console.log(`   Cyber主题多余的key (${extra.length}个):`)
          extra.forEach((k) => console.log(`     - ${k}`))
        }
        console.log()
      } else {
        console.log(`✅ ${lang}: normal 和 cyber 主题一致`)
      }
    }
  }

  // 2. 检查不同语言下，相同主题的key是否一致
  console.log('\n📋 检查 2: 不同语言下，相同主题的key一致性\n')
  for (const theme of THEMES) {
    console.log(`检查 ${theme} 主题:`)

    // 以中文为基准
    const baseKeys = allKeys['zh-CN']?.[theme] || []

    for (const lang of LANGUAGES) {
      if (lang === 'zh-CN') continue

      const currentKeys = allKeys[lang]?.[theme] || []
      const { missing, extra } = compareKeys(currentKeys, baseKeys, lang, 'zh-CN')

      if (missing.length > 0 || extra.length > 0) {
        hasIssues = true
        console.log(`  ⚠️  ${lang} 与 zh-CN 不一致`)
        if (missing.length > 0) {
          console.log(`     ${lang}缺失的key (${missing.length}个):`)
          missing.slice(0, 10).forEach((k) => console.log(`       - ${k}`))
          if (missing.length > 10) {
            console.log(`       ... 还有 ${missing.length - 10} 个`)
          }
        }
        if (extra.length > 0) {
          console.log(`     ${lang}多余的key (${extra.length}个):`)
          extra.slice(0, 10).forEach((k) => console.log(`       - ${k}`))
          if (extra.length > 10) {
            console.log(`       ... 还有 ${extra.length - 10} 个`)
          }
        }
      } else {
        console.log(`  ✅ ${lang} 与 zh-CN 一致`)
      }
    }
    console.log()
  }

  // 3. 统计信息
  console.log('\n📊 统计信息:\n')
  for (const lang of LANGUAGES) {
    for (const theme of THEMES) {
      const count = allKeys[lang]?.[theme]?.length || 0
      console.log(`${lang} - ${theme}: ${count} 个翻译key`)
    }
  }

  if (!hasIssues) {
    console.log('\n✨ 太棒了！所有多语言配置都是完整且一致的！')
  } else {
    console.log('\n⚠️  发现了一些不一致的地方，请检查并修复。')
    process.exit(1)
  }
}

// 运行检查
checkI18nCompleteness().catch((error) => {
  console.error('❌ 检查过程中出错:', error)
  process.exit(1)
})
