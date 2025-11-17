<script setup lang="ts">
  import type { SettingCardProps } from './types'

  defineOptions({
    name: 'SettingCard',
  })

  withDefaults(defineProps<SettingCardProps>(), {
    show: true,
  })
</script>

<template>
  <div v-if="show" class="setting-card">
    <div class="card-header">
      <div class="card-title">
        <i :class="['card-icon', icon]" />
        <span>{{ title }}</span>
      </div>
    </div>

    <div class="card-content">
      <div v-if="hint" class="form-group">
        <div class="form-help">{{ hint }}</div>
      </div>

      <div class="setting-body">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped>
  .setting-card {
    @apply relative overflow-hidden;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.9));
    border: 1px solid var(--color-border-subtle);
    box-shadow:
      0 4px 16px rgba(0, 0, 0, 0.08),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .setting-card:hover {
    border-color: var(--color-border-default);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.03);
    transform: translateY(-1px);
  }

  .card-header {
    @apply border-b px-6 py-4;
    border-color: var(--color-border-subtle);
    background: linear-gradient(135deg, rgba(var(--color-background-600-rgb), 0.3), rgba(var(--color-background-700-rgb), 0.5));
  }

  .card-title {
    @apply flex items-center gap-3 text-base font-medium;
    color: var(--color-content-heading);
  }

  .card-icon {
    @apply text-sm;
    color: var(--color-brand-400);
  }

  .card-content {
    @apply space-y-4 p-6;
  }

  .form-group {
    @apply space-y-2;
  }

  .form-help {
    @apply text-xs leading-relaxed;
    color: var(--color-content-muted);
  }

  .setting-body {
    @apply py-2;
  }
</style>
