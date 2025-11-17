<script setup lang="ts">
  import { computed } from 'vue'
  import { useLayoutStore } from '@/store/layout'
  import { useLayoutNavigation } from '@/composables/useLayoutNavigation'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'LayoutToggleButton',
  })

  const { $t } = useTexts()

  interface Props {
    size?: 'sm' | 'md' | 'lg'
    variant?: 'icon' | 'text' | 'both'
    tooltip?: boolean
    tooltipPlacement?: 'top' | 'bottom' | 'left' | 'right'
  }

  const props = withDefaults(defineProps<Props>(), {
    size: 'md',
    variant: 'icon',
    tooltip: true,
    tooltipPlacement: 'bottom',
  })

  const layoutStore = useLayoutStore()
  const { toggleLayoutWithNavigation } = useLayoutNavigation()

  const buttonClasses = computed(() => ['layout-toggle-btn', `size-${props.size}`, `variant-${props.variant}`])

  const iconClass = computed(() => (layoutStore.isTopLayout ? 'fas fa-grip-vertical' : 'fas fa-grip-horizontal'))

  const tooltipText = computed(() =>
    layoutStore.isLeftLayout ? $t('components.layoutToggleButton.switchToTop') : $t('components.layoutToggleButton.switchToLeft')
  )

  const buttonText = computed(() =>
    layoutStore.isTopLayout ? $t('components.layoutToggleButton.sidebarLayout') : $t('components.layoutToggleButton.topLayout')
  )

  const handleToggle = async () => {
    await toggleLayoutWithNavigation()
  }
</script>

<template>
  <CyberIconButton
    v-if="variant === 'icon'"
    type="cyber"
    :size="props.size === 'sm' ? 'small' : 'medium'"
    :tooltip="tooltip ? tooltipText : undefined"
    :tooltip-placement="tooltipPlacement"
    @click="handleToggle"
  >
    <i :class="iconClass" />
  </CyberIconButton>

  <component
    v-else
    :is="tooltip ? 'CyberTooltip' : 'div'"
    v-bind="tooltip ? { content: tooltipText, placement: tooltipPlacement } : {}"
  >
    <button :class="buttonClasses" @click="handleToggle" :title="!tooltip ? tooltipText : undefined">
      <i :class="iconClass" />
      <span v-if="variant === 'text' || variant === 'both'" class="button-text">
        {{ buttonText }}
      </span>
    </button>
  </component>
</template>

<style scoped>
  .layout-toggle-btn {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08), rgba(var(--color-background-800-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-content-muted);
    @apply flex items-center justify-center;
    transition: all var(--transition-normal) var(--ease-in-out);
    border-radius: var(--radius-lg);
    cursor: pointer;
  }

  .layout-toggle-btn:hover {
    color: var(--color-brand-500);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-background-700-rgb), 0.8));
    transform: translateY(-1px);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .size-sm {
    width: var(--size-button-sm);
    height: var(--size-button-sm);
    padding: 0;
    font-size: var(--size-icon-sm);
  }

  .size-md {
    width: var(--size-button-sm);
    height: var(--size-button-sm);
    padding: 0;
    font-size: var(--size-icon-sm);
  }

  .size-lg {
    width: var(--size-button-md);
    height: var(--size-button-md);
    padding: 0;
    font-size: var(--size-icon-md);
  }

  .variant-text,
  .variant-both {
    padding: 0 var(--space-md);
    width: auto;
    min-width: auto;
  }

  .variant-text.size-sm,
  .variant-both.size-sm {
    height: var(--size-button-sm);
    padding: 0 var(--space-md);
    gap: var(--space-xs);
  }

  .variant-text.size-md,
  .variant-both.size-md {
    height: var(--size-button-sm);
    padding: 0 var(--space-md);
    gap: var(--space-sm);
  }

  .variant-text.size-lg,
  .variant-both.size-lg {
    height: var(--size-button-md);
    padding: 0 var(--space-lg);
    gap: var(--space-sm);
  }

  .button-text {
    @apply whitespace-nowrap;
    font-weight: var(--font-medium);
    font-size: inherit;
  }

  .size-sm .button-text {
    font-size: var(--text-xs);
  }

  .size-md .button-text {
    font-size: var(--text-sm);
  }

  .size-lg .button-text {
    font-size: var(--text-base);
  }
</style>
