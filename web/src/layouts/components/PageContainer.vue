<script setup lang="ts">
  import { computed } from 'vue'
  import { useLayoutStore } from '@/store/layout'

  interface Props {
    title?: string
    showTitle?: boolean // 是否显示标题（可手动控制）
    noPadding?: boolean // 是否禁用 padding
  }

  const props = withDefaults(defineProps<Props>(), {
    title: '',
    showTitle: true,
    noPadding: false,
  })

  const layoutStore = useLayoutStore()

  /* 是否为 Left 布局 */
  const isLeftLayout = computed(() => layoutStore.layout === 'left')

  /* 是否显示标题：Top 布局显示，Left 布局不显示 */
  const shouldShowTitle = computed(() => {
    if (!props.showTitle) return false
    return !isLeftLayout.value && props.title
  })

  /* 是否应用 padding：Top 布局有 padding，Left 布局无 padding（由容器统一管理） */
  const shouldApplyPadding = computed(() => {
    if (props.noPadding) return false
    return !isLeftLayout.value
  })
</script>

<template>
  <div class="page-container" :class="{ 'has-padding': shouldApplyPadding, 'left-layout': isLeftLayout }">
    <div v-if="shouldShowTitle" class="page-header">
      <h1 class="page-title">{{ title }}</h1>
      <slot name="header-actions"></slot>
    </div>

    <div class="page-content">
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
  .page-container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .page-container.has-padding {
    padding: 1.5rem;
  }

  .page-container.left-layout {
    padding: 0;
  }

  .page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid var(--color-border-subtle);
  }

  .page-title {
    font-size: 1.75rem;
    font-weight: 600;
    color: var(--color-content-heading);
    margin: 0;
  }

  .page-content {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
  }

  .left-layout .page-content {
    overflow-y: visible;
  }
</style>
