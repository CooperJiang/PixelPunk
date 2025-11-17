<script setup lang="ts">
  import { computed } from 'vue'
  import { useLayoutStyles } from '@/composables/useLayoutStyles'

  defineOptions({
    name: 'LayoutAdaptivePageContainer',
  })

  interface Props {
    fullWidth?: boolean
    maxWidth?: string
    padding?: string
    pageType?: 'home' | 'admin' | 'docs' | 'gallery' | 'share' | 'default'
  }

  const props = withDefaults(defineProps<Props>(), {
    fullWidth: false,
    maxWidth: '',
    padding: '',
    pageType: 'default',
  })

  const { containerClasses, pageContainerStyle, getPageSpecificStyles } = useLayoutStyles()

  /* 合并容器样式 */
  const containerStyle = computed(() => {
    const baseStyle = { ...pageContainerStyle.value }
    const pageStyles = getPageSpecificStyles()

    if (props.fullWidth) {
      baseStyle.maxWidth = 'none'
    }

    if (props.maxWidth) {
      baseStyle.maxWidth = props.maxWidth
    }

    if (props.padding) {
      baseStyle.padding = props.padding
    }

    return {
      ...baseStyle,
      ...pageStyles,
    }
  })

  const containerClassList = computed(() => ({
    ...containerClasses.value,
    'full-width': props.fullWidth,
    [`page-${props.pageType}`]: true,
  }))
</script>

<template>
  <div :class="containerClassList" :style="containerStyle">
    <slot />
  </div>
</template>

<style scoped>
  .layout-container {
    box-sizing: border-box;
    transition: all var(--layout-animation-duration, 0.3s) ease;
  }

  .layout-container.full-width {
    width: 100%;
    max-width: none !important;
  }

  .layout-top .layout-container {
    margin: 0 auto;
  }

  .layout-left .layout-container {
    width: 100%;
    margin: 0;
  }
</style>
