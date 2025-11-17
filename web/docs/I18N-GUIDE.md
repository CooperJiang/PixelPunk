# PixelPunk 多语言接入指南

## 概述

**双维度多语言系统**: 语言(中/英) × 主题(普通/赛博) = 4 种文案风格

## 目录结构

```
src/locales/
├── zh-CN/
│   ├── common.ts              # 通用文案（所有主题共享）
│   └── themes/
│       ├── normal/            # 普通主题
│       │   ├── pages/         # 页面文案
│       │   ├── admin/         # 管理后台
│       │   ├── components/    # 组件文案
│       │   ├── api/           # API相关
│       │   └── constants.ts   # 常量选项
│       └── cyber/             # 赛博朋克主题（结构同 normal）
└── en-US/                      # 英文（结构同 zh-CN）
```

## 快速开始

### 1. 组件中使用

```vue
<script setup lang="ts">
import { useTexts } from '@/composables/useTexts'
const { $t } = useTexts()
</script>

<template>
  <h1>{{ $t('pages.dashboard.title') }}</h1>
  <p>{{ $t('pages.dashboard.welcome', { name: 'User' }) }}</p>
</template>
```

### 2. 添加新文案

**步骤**：在 4 个语言包中添加相同路径的文案

```typescript
// 1. zh-CN/themes/normal/pages/dashboard.ts
export default {
  title: '仪表盘',
  welcome: '欢迎回来, {name}',
}

// 2. zh-CN/themes/cyber/pages/dashboard.ts
export default {
  title: '█ 数据面板 █',
  welcome: '>>> 系统在线 {name} <<<',
}

// 3. en-US/themes/normal/pages/dashboard.ts
export default {
  title: 'Dashboard',
  welcome: 'Welcome Back, {name}',
}

// 4. en-US/themes/cyber/pages/dashboard.ts
export default {
  title: '█ DATA PANEL █',
  welcome: '>>> SYSTEM ONLINE {name} <<<',
}
```

**在 index.ts 中导出**：
```typescript
import dashboard from './dashboard'
export default { dashboard }
```

### 3. 常量选项（下拉框等）

```typescript
// src/constants/myOptions.ts
import type { Composer } from '@/composables/useTexts'

export const getMyOptions = ($t: Composer['$t']) => [
  { label: $t('constants.myOptions.option1'), value: 'value1' },
  { label: $t('constants.myOptions.option2'), value: 'value2' },
]

export type MyOption = ReturnType<typeof getMyOptions>[number]
```

```vue
<script setup>
import { useTexts } from '@/composables/useTexts'
import { getMyOptions } from '@/constants/myOptions'

const { $t } = useTexts()
const options = getMyOptions($t)
</script>

<template>
  <select>
    <option v-for="opt in options" :key="opt.value" :value="opt.value">
      {{ opt.label }}
    </option>
  </select>
</template>
```

## 主题切换

```typescript
import { useTexts } from '@/composables/useTexts'

const { $toggleTheme, $setTheme } = useTexts()

$toggleTheme()              // 切换主题（normal ⇄ cyber）
$setTheme('cyber')          // 设置指定主题
```

## 命名规范

### 文案 Key 命名

- 使用小驼峰命名 (camelCase)
- 层级用点号分隔
- 保持语义清晰

```typescript
// ✅ 好的命名
$t('pages.upload.dragToUpload')
$t('components.fileCard.actions.delete')

// ❌ 不好的命名
$t('pages.upload.text1')
$t('components.fileCard.btn')
```

### 文案分类

| 类型 | 路径 | 说明 |
|------|------|------|
| 页面专属 | `pages/{pageName}.ts` | 仅在特定页面使用 |
| 管理后台 | `admin/{moduleName}.ts` | 后台管理页面 |
| 组件 | `components/{componentName}.ts` | 可复用组件 |
| API | `api/{moduleName}.ts` | API相关文案 |
| 通用 | `common.ts` | 所有主题共享 |
| 常量 | `constants.ts` | 枚举选项等 |

## 重要规则

### ⚠️ 禁止硬编码文本

```typescript
// ❌ 错误
const message = '上传成功'

// ✅ 正确
const { $t } = useTexts()
const message = $t('messages.uploadSuccess')
```

### ⚠️ 必须在 4 个语言包同步

添加新文案时，务必在 4 个位置都添加：
1. `zh-CN/themes/normal/` ✓
2. `zh-CN/themes/cyber/` ✓
3. `en-US/themes/normal/` ✓
4. `en-US/themes/cyber/` ✓

### ⚠️ value 不使用多语言

API 枚举值、数据库字段保持原样，仅 label 使用多语言：

```typescript
// ✅ 正确
export const getStatusOptions = ($t) => [
  { label: $t('status.active'), value: 'ACTIVE' },
]

// ❌ 错误
export const getStatusOptions = ($t) => [
  { label: $t('status.active'), value: $t('status.active') },
]
```

**特殊情况**：如果后端 API 返回中文枚举值，需要保留：

```typescript
export const getResolutionOptions = ($t) => [
  { label: $t('constants.resolution.low'), value: '低分辨率' }, // 后端API枚举值
]
```

## 赛博朋克风格参考

赛博风格特点：
- 特殊字符：`█ ◢ ◣ ▓ ░ >>>`
- 全大写英文
- 技术感词汇：矩阵、系统、终端

```typescript
// Normal
title: '上传文件'
success: '上传成功'

// Cyber
title: '█ 数据传输 █'
success: '>>> 传输完成 <<<'
```

## 相关文件

- **核心实现**: `src/store/textTheme.ts`
- **Composable**: `src/composables/useTexts.ts`
- **全局插件**: `src/plugins/textPlugin.ts`
- **类型定义**: `src/types/global.d.ts`
