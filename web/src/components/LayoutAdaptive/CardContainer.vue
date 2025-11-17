<script setup lang="ts">
  import { computed } from 'vue'
  import { useLayoutStyles } from '@/composables/useLayoutStyles'

  defineOptions({
    name: 'LayoutAdaptiveCardContainer',
  })

  interface Props {
    tag?: string
    hoverable?: boolean
    padding?: string
    borderRadius?: string
    variant?: 'default' | 'admin' | 'gallery' | 'stats'
  }

  const props = withDefaults(defineProps<Props>(), {
    tag: 'div',
    hoverable: false,
    padding: '',
    borderRadius: '',
    variant: 'default',
  })

  const { cardStyle, isLeftLayout } = useLayoutStyles()

  const computedCardStyle = computed(() => {
    const baseStyle = { ...cardStyle.value }

    if (props.padding) {
      baseStyle.padding = props.padding
    }

    if (props.borderRadius) {
      baseStyle.borderRadius = props.borderRadius
    }

    switch (props.variant) {
      case 'admin':
        baseStyle.padding = isLeftLayout.value ? '1rem' : '1.5rem'
        baseStyle.borderRadius = isLeftLayout.value ? '0.5rem' : '0.75rem'
        break
      case 'gallery':
        baseStyle.padding = isLeftLayout.value ? '0.5rem' : '1rem'
        baseStyle.borderRadius = 'var(--layout-card-radius)'
        break
      case 'stats':
        baseStyle.padding = isLeftLayout.value ? '1.25rem' : '1.75rem'
        break
    }

    return baseStyle
  })

  const cardClasses = computed(() => ({
    'adaptive-card': true,
    hoverable: props.hoverable,
    [`variant-${props.variant}`]: true,
  }))
</script>

<template>
  <component :is="tag" :class="cardClasses" :style="computedCardStyle">
    <slot />
  </component>
</template>

<style scoped>
  .adaptive-card {
    background: var(--color-background-800);
    border: 1px solid var(--color-border-default);
    box-sizing: border-box;
    position: relative;
    overflow: hidden;
  }

  .adaptive-card.hoverable {
    cursor: pointer;
    transition: all var(--layout-animation-duration, 0.3s) ease;
  }

  .adaptive-card.hoverable:hover {
    transform: scale(var(--layout-hover-scale, 1.02));
    border-color: var(--color-brand-500);
    box-shadow: 0 8px 32px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .variant-admin {
    background: rgba(var(--color-background-900-rgb), 0.9);
    border-color: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .variant-gallery {
    background: rgba(var(--color-background-900-rgb), 0.95);
    border: none;
  }

  .variant-stats {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.9) 0%,
      rgba(var(--color-background-800-rgb), 0.8) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
  }

  @media (max-width: 768px) {
    .adaptive-card {
      padding: calc(var(--layout-card-padding) * 0.75) !important;
      border-radius: calc(var(--layout-card-radius) * 0.75) !important;
    }
  }
</style>
