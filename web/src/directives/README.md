# 自定义指令文档

## v-loading 指令

基于现有CyberLoading UI的加载指令，支持多种赛博朋克风格的loading效果。

### 基本用法

```vue
<template>
  <!-- 布尔值控制 -->
  <div v-loading="loading">内容区域</div>

  <!-- 字符串文本 -->
  <div v-loading="'正在处理...'">内容区域</div>

  <!-- 对象配置 -->
  <div v-loading="{ loading: isLoading, text: '数据加载中...', type: 'progress' }">内容区域</div>
</template>

<script setup>
  import { ref } from 'vue'

  const loading = ref(true)
  const isLoading = ref(false)
</script>
```

### Loading类型

#### 1. 默认类型 (default)

经典的六边形 + 圆环 + 扫描线效果

```vue
<div v-loading="{ loading: true, type: 'default', text: '加载中...' }">
  内容区域
</div>
```

#### 2. 进度条类型 (progress)

赛博朋克风格的进度条，支持真实进度和模拟进度

```vue
<!-- 模拟进度 -->
<div v-loading="{ loading: true, type: 'progress', text: '数据处理中...' }">
  内容区域
</div>

<!-- 真实进度 -->
<div v-loading="{ loading: true, type: 'progress', progress: 75, text: '上传进度' }">
  内容区域
</div>
```

#### 3. 脉冲类型 (pulse)

多层脉冲圆环效果，使用霓虹青、霓虹粉、神秘紫三色

```vue
<div v-loading="{ loading: true, type: 'pulse', text: '连接中...' }">
  内容区域
</div>
```

#### 4. 矩阵雨类型 (matrix)

矩阵代码雨效果，使用日文假名和数字0/1

```vue
<div v-loading="{ loading: true, type: 'matrix', text: 'ACCESSING DATA...' }">
  内容区域
</div>
```

### 完整配置选项

```vue
<div
  v-loading="{
    loading: true,
    type: 'progress',
    text: '加载中...',
    progress: 50,
    background: 'rgba(0, 0, 0, 0.8)',
    textColor: 'var(--color-cyber-pink)',
  }"
>
  内容区域
</div>
```

### 配置参数

| 参数       | 类型    | 默认值                                  | 说明                                       |
| ---------- | ------- | --------------------------------------- | ------------------------------------------ |
| loading    | Boolean | false                                   | 是否显示loading                            |
| visible    | Boolean | false                                   | 是否显示loading（与loading等效）           |
| type       | String  | 'default'                               | loading类型：default/progress/pulse/matrix |
| text       | String  | '加载中...'                             | loading文本                                |
| progress   | Number  | 0                                       | 进度值(0-100)，仅progress类型有效          |
| background | String  | rgba(var(--color-cyber-gray-rgb), 0.85) | 背景颜色                                   |
| textColor  | String  | var(--color-cyber-blue)                 | 文本颜色                                   |

### 特性

- ✅ 自动为目标元素添加相对定位
- ✅ 完整的CyberPunk风格动画（六边形、扫描线、故障效果）
- ✅ 响应式更新
- ✅ 自动清理资源
- ✅ 支持多种使用方式
- ✅ TypeScript支持

### 动画效果

- 六边形脉冲动画
- 圆形旋转动画
- 扫描线动画
- 故障文字效果

## v-lazy 指令

文件懒加载指令，使用IntersectionObserver API实现高性能懒加载。

### 基本用法

```vue
<template>
  <!-- 基础懒加载 -->
  <img v-lazy src="/api/image/large.jpg" alt="文件" />

  <!-- 条件懒加载 -->
  <img v-lazy="shouldLazy" src="/api/image/large.jpg" alt="文件" />
</template>
```

### 特性

- ✅ 基于IntersectionObserver API
- ✅ 高性能，自动管理观察器
- ✅ 支持条件懒加载
- ✅ 优雅的渐入动画
- ✅ 自动清理资源
