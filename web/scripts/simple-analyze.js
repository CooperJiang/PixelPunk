#!/usr/bin/env node

import { fileURLToPath } from 'url'
import { dirname, join } from 'path'
import { existsSync } from 'fs'
import { exec } from 'child_process'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)
const rootDir = join(__dirname, '..')
const statsPath = join(rootDir, 'dist', 'stats.html')

console.log('🔍 检查打包分析文件...')

if (!existsSync(statsPath)) {
  console.log('❌ 未找到 stats.html 文件')
  console.log('💡 请先运行: pnpm build')
  process.exit(1)
}

console.log('✅ 找到 stats.html')
console.log('🌐 正在浏览器中打开...')

// 跨平台打开文件
const openCommand = process.platform === 'darwin'
  ? 'open'
  : process.platform === 'win32'
  ? 'start'
  : 'xdg-open'

exec(`${openCommand} "${statsPath}"`, (error) => {
  if (error) {
    console.error('❌ 打开失败:', error.message)
    console.log(`📁 文件位置: ${statsPath}`)
    console.log('💡 请手动打开该文件')
    process.exit(1)
  }
  console.log('✅ 已在浏览器中打开分析报告')
  console.log(`📊 文件位置: ${statsPath}`)
})
