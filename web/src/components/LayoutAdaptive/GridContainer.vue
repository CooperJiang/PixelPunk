<script setup lang="ts">
  import { computed } from 'vue'
  import { useLayoutStyles } from '@/composables/useLayoutStyles'

  defineOptions({
    name: 'LayoutAdaptiveGridContainer',
  })

  interface Props {
    minColumnWidth?: number
    columns?: string
    gap?: string
    autoColumns?: boolean
  }

  const props = withDefaults(defineProps<Props>(), {
    minColumnWidth: 300,
    columns: '',
    gap: '',
    autoColumns: true,
  })

  const { gridContainerStyle, getGridCols, getSpacing } = useLayoutStyles()

  const gridStyle = computed(() => {
    const baseStyle = { ...gridContainerStyle.value }

    if (props.columns) {
      baseStyle.gridTemplateColumns = props.columns
    } else if (props.autoColumns) {
      baseStyle.gridTemplateColumns = getGridCols(props.minColumnWidth)
    }

    if (props.gap) {
      baseStyle.gap = props.gap
    } else {
      baseStyle.gap = getSpacing('card')
    }

    return baseStyle
  })
</script>

<template>
  <div class="adaptive-grid" :style="gridStyle">
    <slot />
  </div>
</template>

<style scoped>
  .adaptive-grid {
    width: 100%;
    transition: all var(--layout-animation-duration, 0.3s) ease;
  }

  .adaptive-grid > * {
    min-width: 0;
  }
</style>
