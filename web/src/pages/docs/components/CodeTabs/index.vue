<script setup lang="ts">
  import { ref, watch } from 'vue'
  import type { CodeTabsProps, Tab } from './types'

  defineOptions({
    name: 'CodeTabs',
  })

  const props = withDefaults(defineProps<CodeTabsProps>(), {
    default: 'curl',
  })

  const emit = defineEmits<{
    (e: 'update:tab', value: string): void
  }>()

  const activeTab = ref(props.default)

  function changeTab(tab: string) {
    activeTab.value = tab
    emit('update:tab', tab)
  }

  watch(
    () => props.default,
    (newValue) => {
      activeTab.value = newValue
    }
  )
</script>

<template>
  <div class="code-tabs">
    <div class="code-tabs-header">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-button"
        :class="[{ active: activeTab === tab.value }]"
        @click="changeTab(tab.value)"
      >
        {{ tab.label }}
      </button>
    </div>
    <div class="code-tabs-content">
      <slot v-if="activeTab === 'curl'" name="curl" />
      <slot v-if="activeTab === 'python'" name="python" />
      <slot v-if="activeTab === 'nodejs'" name="nodejs" />
      <slot v-if="activeTab === 'javascript'" name="javascript" />
      <slot v-if="activeTab === 'java'" name="java" />
      <slot v-if="activeTab === 'go'" name="go" />
    </div>
  </div>
</template>

<style scoped>
  .code-tabs {
    margin-bottom: 0;
    width: 100%;
    max-width: 100%;
    overflow-x: hidden;
  }

  .code-tabs-header {
    display: flex;
    overflow-x: auto;
    background: rgba(var(--color-background-800-rgb), 0.6);
    border-top-left-radius: 8px;
    border-top-right-radius: 8px;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-bottom: none;
  }

  .code-tabs-header::-webkit-scrollbar {
    height: 4px;
  }

  .code-tabs-header::-webkit-scrollbar-track {
    background: rgba(var(--color-background-700-rgb), 0.5);
  }

  .code-tabs-header::-webkit-scrollbar-thumb {
    background-color: rgba(var(--color-brand-500-rgb), 0.5);
    border-radius: var(--radius-sm);
  }

  .tab-button {
    padding: 0.75rem 1.25rem;
    background: transparent;
    border: none;
    color: var(--color-content-heading);
    cursor: pointer;
    font-family: 'Fira Code', monospace;
    font-size: 0.875rem;
    white-space: nowrap;
    transition: all 0.3s ease;
    position: relative;
  }

  .tab-button.active {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15) 0%, rgba(var(--color-brand-500-rgb), 0.08) 100%);
    color: var(--color-brand-500);
    border-bottom: 2px solid var(--color-brand-500);
  }

  .tab-button:hover:not(.active) {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
  }

  .code-tabs-content {
    width: 100%;
    max-width: 100%;
  }

  @media (max-width: 768px) {
    .tab-button {
      padding: 0.5rem 0.75rem;
      font-size: 0.8rem;
    }
  }
</style>
