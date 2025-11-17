# LocaleSwitch 多语言切换组件

参考 `ThemeToggle` 组件实现的多语言切换组件。

## 功能特点

- 🌍 支持多语言切换（简体中文、英语、日语）
- 🎨 与主题切换组件风格一致
- 📱 响应式设计，支持移动端
- ✨ 流畅的动画效果
- 🎯 灵活的尺寸和位置配置

## 使用方式

### 基础用法

```vue
<script setup>
import LocaleSwitch from '@/components/LocaleSwitch'
</script>

<template>
  <LocaleSwitch />
</template>
```

### 带文字标签

```vue
<template>
  <LocaleSwitch show-text />
</template>
```

### 不同尺寸

```vue
<template>
  <LocaleSwitch size="small" />
  <LocaleSwitch size="medium" />
  <LocaleSwitch size="large" />
</template>
```

### 自定义位置

```vue
<template>
  <!-- 下拉菜单在上方 -->
  <LocaleSwitch placement="top" />

  <!-- 下拉菜单在下方 -->
  <LocaleSwitch placement="bottom" />

  <!-- 下拉菜单左对齐 -->
  <LocaleSwitch align="left" />

  <!-- 下拉菜单右对齐 -->
  <LocaleSwitch align="right" />
</template>
```

## Props

| 参数 | 说明 | 类型 | 可选值 | 默认值 |
|------|------|------|--------|--------|
| size | 按钮尺寸 | string | small / medium / large | small |
| showDropdown | 是否显示下拉菜单 | boolean | - | true |
| showText | 是否显示文字 | boolean | - | false |
| placement | 下拉菜单位置 | string | top / bottom | bottom |
| align | 下拉菜单对齐方式 | string | left / right | right |

## Composable: useLocale

如果需要在代码中操作语言切换，可以使用 `useLocale` composable：

```vue
<script setup>
import { useLocale } from '@/composables/useLocale'

const {
  currentLocale,        // 当前语言
  setLocale,            // 设置语言
  localeOptions,        // 所有语言选项
  currentLocaleInfo,    // 当前语言信息
  getLocaleLabel,       // 获取语言标签
  getLocaleDescription, // 获取语言描述
  getLocaleNativeName,  // 获取语言原生名称
} = useLocale()

// 切换到英语
const switchToEnglish = () => {
  setLocale('en-US')
}
</script>
```

## 文件结构

```
components/LocaleSwitch/
├── index.vue         # 组件主文件
└── types.ts          # TypeScript 类型定义

composables/
└── useLocale.ts      # 语言管理 Composable

locales/
├── zh-CN/common.ts   # 中文翻译
├── en-US/common.ts   # 英文翻译
└── ja-JP/common.ts   # 日文翻译
```

## 与 ThemeToggle 的对比

| 特性 | ThemeToggle | LocaleSwitch |
|------|-------------|--------------|
| 切换对象 | 视觉主题 | 界面语言 |
| 选项数量 | 12个主题 | 3种语言 |
| 图标 | 主题相关图标 | 语言相关图标 |
| Composable | useTheme | useLocale |
| Store | visualThemeStore | textThemeStore |

## 样式特点

- 使用 CSS 变量保持主题一致性
- 支持 hover、active 等交互状态
- 流畅的过渡动画
- 响应式设计，移动端优化
- 与系统主题无缝集成
