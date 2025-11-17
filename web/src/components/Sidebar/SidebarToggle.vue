<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'SidebarToggle',
  })

  interface Props {
    collapsed?: boolean
  }

  withDefaults(defineProps<Props>(), {
    collapsed: false,
  })

  const { $t } = useTexts()
  const emit = defineEmits<{
    (e: 'toggle'): void
  }>()

  const handleToggle = () => {
    emit('toggle')
  }
</script>

<template>
  <button
    class="sidebar-toggle"
    :class="{ 'toggle-collapsed': collapsed }"
    @click="handleToggle"
    :title="collapsed ? $t('components.sidebar.toggle.expand') : $t('components.sidebar.toggle.collapse')"
  >
    <div class="toggle-arrows">
      <i class="fas fa-chevron-left" :class="{ 'arrow-hidden': collapsed }" />
      <i class="fas fa-chevron-left" :class="{ 'arrow-hidden': collapsed }" />
    </div>

    <div v-if="collapsed" class="expand-icon">
      <i class="fas fa-bars" />
    </div>

    <div class="toggle-glow" />
  </button>
</template>

<style scoped>
  .sidebar-toggle {
    @apply relative flex h-7 w-7 items-center justify-center;
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
    color: var(--color-brand-500);
    transition: all 0.2s ease;
    overflow: hidden;
    border-radius: var(--radius-sm);
  }

  .sidebar-toggle:hover {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: scale(1.05);
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .sidebar-toggle:active {
    transform: scale(0.95);
  }

  .toggle-arrows {
    @apply relative flex items-center;
    transition: all 0.2s ease;
  }

  .toggle-arrows i {
    font-size: 10px;
    transition: all 0.2s ease;
    position: relative;
  }

  .toggle-arrows i:first-child {
    margin-right: -3px;
    opacity: 0.7;
  }

  .toggle-arrows i:last-child {
    opacity: 1;
  }

  .arrow-hidden {
    opacity: 0 !important;
    transform: translateX(10px);
  }

  .expand-icon {
    @apply absolute inset-0 flex items-center justify-center;
    font-size: 11px;
    animation: scaleIn var(--transition-fast) var(--ease-out);
  }

  .toggle-glow {
    @apply absolute inset-0;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.15) 0%, transparent 70%);
    opacity: 0;
    transition: opacity 0.2s ease;
    border-radius: var(--radius-sm);
  }

  .sidebar-toggle:hover .toggle-glow {
    opacity: 1;
  }

  .sidebar-toggle:hover .toggle-arrows {
    transform: translateX(-1px);
  }

  .toggle-collapsed .sidebar-toggle:hover .toggle-arrows {
    transform: translateX(1px);
  }
</style>
