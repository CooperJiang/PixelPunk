<script setup lang="ts">
  import { ref, watch, provide, useSlots, VNode } from 'vue'

  defineOptions({
    name: 'CyberTabs',
  })

  interface Props {
    modelValue: string
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: '',
  })

  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
    (e: 'change', value: string): void
  }>()

  const activeTab = ref(props.modelValue)
  const slots = useSlots()

  watch(
    () => props.modelValue,
    (newVal) => {
      activeTab.value = newVal
    }
  )

  const getTabs = () => {
    const defaultSlot = slots.default?.() || []
    return defaultSlot
      .filter((vnode: VNode) => vnode.type && typeof vnode.type === 'object' && (vnode.type as any).name === 'CyberTabPane')
      .map((vnode: VNode) => ({
        name: (vnode.props as any)?.name || '',
        label: (vnode.props as any)?.label || '',
        disabled: (vnode.props as any)?.disabled || false,
      }))
  }

  const tabs = getTabs()

  const handleTabClick = (tabName: string, disabled: boolean) => {
    if (disabled) return
    activeTab.value = tabName
    emit('update:modelValue', tabName)
    emit('change', tabName)
  }

  provide('activeTab', activeTab)
</script>

<template>
  <div class="cyber-tabs">
    <div class="tabs-header">
      <div
        v-for="tab in tabs"
        :key="tab.name"
        class="tab-item"
        :class="{
          active: activeTab === tab.name,
          disabled: tab.disabled,
        }"
        @click="handleTabClick(tab.name, tab.disabled)"
      >
        <span class="tab-label">{{ tab.label }}</span>
        <div v-if="activeTab === tab.name" class="tab-indicator" />
      </div>
    </div>

    <div class="tabs-content">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .cyber-tabs {
    display: flex;
    flex-direction: column;
    gap: var(--space-xl);
  }

  .tabs-header {
    display: flex;
    gap: var(--space-sm);
    border-bottom: 1px solid var(--color-border-default);
    padding-bottom: var(--space-xs);
    overflow-x: auto;
    scrollbar-width: thin;

    &::-webkit-scrollbar {
      height: 4px;
    }

    &::-webkit-scrollbar-track {
      background: var(--color-background-800);
    }

    &::-webkit-scrollbar-thumb {
      background: var(--color-background-600);
      border-radius: var(--radius-full);

      &:hover {
        background: var(--color-brand-500);
      }
    }
  }

  .tab-item {
    position: relative;
    padding: var(--space-md) var(--space-xl);
    cursor: pointer;
    transition: all var(--transition-fast);
    border-radius: var(--radius-md) var(--radius-md) 0 0;
    white-space: nowrap;
    user-select: none;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover:not(.disabled) {
      background: var(--color-background-700);
    }

    &.active {
      background: var(--color-background-700);
      color: var(--color-brand-400);
      box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);

      .tab-label {
        color: var(--color-brand-400);
        text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.5);
      }
    }

    &.disabled {
      opacity: 0.5;
      cursor: not-allowed;
      color: var(--color-content-muted);

      &:hover {
        background: transparent;
      }
    }
  }

  .tab-label {
    font-size: var(--text-base);
    font-weight: var(--font-medium);
    color: var(--color-content-default);
    transition: all var(--transition-fast);
  }

  .tab-indicator {
    position: absolute;
    bottom: -1px;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, var(--color-brand-500), var(--color-brand-400));
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.6);
    animation: glow 2s ease-in-out infinite;
  }

  @keyframes glow {
    0%,
    100% {
      box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.6);
    }
    50% {
      box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.8);
    }
  }

  .tabs-content {
    flex: 1;
    min-height: 0;
  }

  @media (max-width: 768px) {
    .tab-item {
      padding: var(--space-sm) var(--space-lg);
    }

    .tab-label {
      font-size: var(--text-sm);
    }
  }
</style>
