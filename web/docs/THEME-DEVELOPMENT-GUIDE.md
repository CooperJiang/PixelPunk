# PixelPunk 主题开发指南

> 🎨 3 步完成主题开发：复制模板 → 改颜色 → 注册主题
>
> ⚡ **性能优化**：主题采用动态懒加载，无需手动导入到 main.css

---

## 🚀 快速开始

### 第 1 步: 复制模板文件

```bash
cd web/src/styles/design-system/
cp theme-template.css theme-ocean-blue.css  # 改成你的主题名
```

### 第 2 步: 修改颜色

打开 `theme-ocean-blue.css`，修改核心颜色变量：

```css
/* 1. 修改主题名称 */
:root[data-theme='ocean-blue'] {
  /* 改成你的主题名 */

  /* 2. 修改品牌色（按钮、链接等） - 必需完整色阶 */
  --color-brand-300: #bfdbfe;
  --color-brand-400: #93c5fd;
  --color-brand-500: #0ea5e9; /* ⭐️ 主色调 */
  --color-brand-500-rgb: 14, 165, 233; /* ⭐️ RGB 格式 */
  --color-brand-600: #3b82f6;
  --color-brand-700: #2563eb;

  /* 3. 修改背景色（页面背景） */
  --color-background-900: #0c1222; /* ⭐️ 主背景 */
  --color-background-900-rgb: 12, 18, 34;
  --color-background-800: #151b26; /* 容器背景 */
  --color-background-700: #1a202c; /* 卡片背景 */

  /* 4. 修改文本色（确保与背景对比清晰） */
  --color-content-default: #cbd5e1; /* ⭐️ 默认文本 */
  --color-content-heading: #f9fafb; /* 标题文本 */

  /* 5. 状态色（必需完整色阶 400/500/600） */
  --color-success-400: #6ee7b7;
  --color-success-500: #34d399; /* ⭐️ 成功色 */
  --color-success-600: #10b981;
  --color-warning-400: #fcd34d;
  --color-warning-500: #fbbf24; /* ⭐️ 警告色 */
  --color-warning-600: #f59e0b;
  --color-error-400: #fca5a5;
  --color-error-500: #f87171; /* ⭐️ 错误色 */
  --color-error-600: #ef4444;

  /* ⚠️ 重要：所有颜色变量都需要对应的 RGB 格式 */
  /* 详见模板文件中的完整注释 */
}
```

💡 **颜色转换**: [RGB.to](https://rgb.to/) 可将 HEX 转换为 RGB

⚠️ **重要提示**：

- 状态色必须定义完整的 400/500/600 三个等级（消息徽章等组件依赖 600 等级）
- 品牌色建议定义完整的 300-700 色阶
- 所有主要颜色都需要对应的 `-rgb` 变量（用于透明度计算）
- 亮色主题需要覆盖 `--color-white` 和 `--color-black` 变量

### 第 3 步: 注册主题

**3.1 在 visualTheme.ts 中注册类型**

```typescript
// web/src/store/visualTheme.ts

// 在类型定义中添加你的主题
export type VisualTheme =
  | 'light'
  | 'dark'
  | 'cyberpunk-cyan'
  // ... 其他主题
  | 'ocean-blue' // ← 添加你的主题

// 在 isLight computed 中注册亮色主题（如果是亮色）
const isLight = computed(() => {
  const lightThemes: VisualTheme[] = [
    'light',
    'lavender-dream',
    // ... 其他亮色主题
    'ocean-blue', // ← 如果是亮色主题，添加到这里
  ]
  return lightThemes.includes(appliedTheme.value)
})
```

**3.2 在 useTheme.ts 中添加主题配置**

```typescript
// web/src/composables/useTheme.ts

// 在 THEME_CONFIGS 数组中添加一个对象即可
const THEME_CONFIGS: ThemeConfig[] = [
  // ... 其他主题 ...

  // ← 添加你的主题（一个对象搞定所有信息！）
  {
    label: 'Ocean Blue · 海洋蓝', // 显示名称
    value: 'ocean-blue', // 主题值（与 CSS 文件名一致）
    description: '深邃的海洋蓝主题', // 描述
    icon: 'water', // FontAwesome 图标名
    isDark: true, // 是否暗色主题
  },
]
```

💡 **图标选择**: 访问 [FontAwesome 图标库](https://fontawesome.com/icons) 搜索图标，复制图标名（不含 `fa-` 前缀）

**常用图标推荐**:

- 天气: `sun`, `moon`, `cloud`, `snowflake`
- 自然: `leaf`, `tree`, `water`, `fire`, `seedling`
- 食物: `lemon`, `ice-cream`, `candy-cane`
- 装饰: `heart`, `star`, `sparkles`, `gem`, `crown`
- 其他: `bolt`, `palette`, `wand-magic-sparkles`

**3.3 (可选) 添加到 main.css**

⚠️ **重要**：由于主题采用**动态懒加载**机制，大部分主题**无需**在 `main.css` 中手动导入。

**只有**以下默认主题需要在 `main.css` 中预先导入（见 `themeLoader.ts` 的 `DEFAULT_THEMES`）：
- `theme-dark.css`
- `theme-light.css`
- `theme-cyberpunk-cyan.css`

其他主题会在用户首次切换时自动加载，无需手动导入。

如果你希望将新主题设为**默认主题**（预先加载），需要：
1. 在 `main.css` 中添加 `@import`
2. 在 `themeLoader.ts` 的 `DEFAULT_THEMES` 数组中添加主题名

---

## ✅ 完成！

```typescript
import { useTheme } from '@/composables/useTheme'
const theme = useTheme()
theme.setVisualTheme('ocean-blue') // 切换到新主题
```

---

## 📝 注意事项

### 主题命名

- ✅ 使用 `kebab-case`: `ocean-blue`, `forest-green`
- ❌ 避免版本号: `theme-v2`

### RGB 值格式

```css
/* ✅ 正确 */
--color-brand-500-rgb: 14, 165, 233;

/* ❌ 错误 */
--color-brand-500-rgb: rgb(14, 165, 233); /* 不要加 rgb() */
```

### 颜色选择建议

| 类型   | 说明                                 | 必需变量                       |
| ------ | ------------------------------------ | ------------------------------ |
| 品牌色 | 按钮、链接等强调元素                 | 300/400/500/600/700 + RGB      |
| 背景色 | 页面背景，至少 3 个主要层级          | 700/800/900 + RGB              |
| 文本色 | 与背景对比度要足够                   | default/heading/muted/disabled |
| 状态色 | 成功/警告/错误，完整色阶             | 400/500/600 + RGB（每种状态）  |
| 功能色 | 多彩功能区分（可选）                 | primary/secondary/accent/info  |
| 通用色 | 亮色主题必需覆盖（暗色主题使用默认） | white/black + RGB              |

⚠️ **必需变量清单**（缺少会导致组件显示异常）：

- `--color-error-600` - 消息徽章使用
- `--color-success-600` - 成功提示使用
- `--color-warning-600` - 警告提示使用
- `--color-brand-300` - 浅色品牌元素使用
- 所有 `-rgb` 变量 - 透明度计算必需

---

## 🔧 调试技巧

### 查看当前主题

```typescript
const theme = useTheme()
console.log(theme.appliedTheme.value) // 当前主题名
console.log(theme.isDark.value) // 是否暗色
```

### 检查 CSS 变量

```javascript
// 浏览器控制台
getComputedStyle(document.documentElement).getPropertyValue('--color-brand-500')
```

### 快速测试所有主题

```typescript
const themes = ['light', 'dark', 'ocean-blue']
let i = 0
setInterval(() => {
  theme.setVisualTheme(themes[i++ % themes.length])
}, 2000)
```

---

## 💡 常见问题

**Q: 主题颜色不生效？**

检查：

1. 主题名称是否一致（CSS 文件名、visualTheme.ts 类型定义、useTheme.ts 配置）
2. 是否在 visualTheme.ts 的 `VisualTheme` 类型中添加
3. 是否在 useTheme.ts 的 `THEME_CONFIGS` 中添加配置
4. 主题 CSS 文件是否在 `web/src/styles/design-system/` 目录下
5. 是否重启开发服务器（TypeScript 类型更改需要重启）

**Q: 如何选择配色？**

推荐工具：

- [Coolors](https://coolors.co/) - 配色生成器
- [Adobe Color](https://color.adobe.com/) - 专业配色

**Q: HEX 转 RGB？**

```javascript
// 浏览器控制台
function hexToRgb(hex) {
  const r = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  return r ? `${parseInt(r[1], 16)}, ${parseInt(r[2], 16)}, ${parseInt(r[3], 16)}` : null
}
hexToRgb('#0ea5e9') // "14, 165, 233"
```

---

## 📦 文件位置

```
web/src/styles/design-system/
├── theme-template.css          ⭐️ 主题模板（复制这个）
├── theme-light.css             亮色主题（默认预加载）
├── theme-dark.css              暗色主题（默认预加载）
├── theme-cyberpunk-cyan.css    赛博朋克主题（默认预加载）
└── theme-你的主题名.css         你的新主题（懒加载）

web/src/store/
└── visualTheme.ts              注册主题类型

web/src/composables/
└── useTheme.ts                 主题配置和管理

web/src/utils/
└── themeLoader.ts              动态加载主题（懒加载）

web/src/styles/
└── main.css                    只导入默认主题
```

---

## 🎨 主题系统架构

```
useTheme() Composable (统一入口)
    ↓
visualTheme Store (视觉主题状态管理)
    ↓
themeLoader (动态懒加载主题 CSS)
    ↓
CSS Variables (应用到 :root[data-theme])
```

**性能优化特性**：
- ✅ 默认主题预加载（`dark`, `light`, `cyberpunk-cyan`）
- ✅ 其他主题按需懒加载，减少初始包体积
- ✅ 使用 `requestIdleCallback` 在浏览器空闲时预加载
- ✅ 防止重复加载，每个主题只加载一次

### 使用主题

```typescript
import { useTheme } from '@/composables/useTheme'

const theme = useTheme()

// 设置主题
theme.setVisualTheme('ocean-blue')

// 判断主题
theme.isDark // 是否暗色
theme.isLight // 是否亮色

// 切换主题
theme.toggleVisualTheme()

// 获取主题列表（用于主题切换组件）
theme.themeOptions.value // 所有主题的完整信息
theme.currentThemeInfo.value // 当前主题信息
theme.allThemes // 所有主题名称数组
```

### 主题切换组件示例

```vue
<script setup lang="ts">
  import { useTheme } from '@/composables/useTheme'

  const { themeOptions, currentThemeInfo, setTheme } = useTheme()
</script>

<template>
  <!-- 当前主题 -->
  <div>{{ currentThemeInfo.label }}</div>

  <!-- 主题列表 -->
  <button v-for="option in themeOptions" :key="option.value" @click="setTheme(option.value)">
    <i :class="`fa-${option.icon}`" />
    {{ option.label }}
  </button>
</template>
```

---

**就这么简单！** 🎉
