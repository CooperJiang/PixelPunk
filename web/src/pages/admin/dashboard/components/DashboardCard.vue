<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  const { $t } = useTexts()

  interface Props {
    title: string
    badge?: string
    badgeType?: 'success' | 'warning' | 'danger' | 'info' | 'normal'
    iconGradient?: string[]
    refreshable?: boolean
    isRefreshing?: boolean
    lastUpdated?: string
    showUpdateTime?: boolean
    size?: 'small' | 'medium' | 'large'
  }

  const props = withDefaults(defineProps<Props>(), {
    badgeType: 'normal',
    iconGradient: () => ['rgba(var(--color-brand-500-rgb), 0.15)', 'rgba(var(--color-brand-500-rgb), 0.25)'],
    refreshable: true,
    showUpdateTime: true,
    size: 'medium',
  })

  defineEmits<{
    refresh: []
  }>()

  const cardClass = computed(() => ({
    [`card-${props.size}`]: true,
  }))

  const badgeClass = computed(() => ({
    [`badge-${props.badgeType}`]: true,
  }))

  const iconStyle = computed(() => {
    const [from, to] = props.iconGradient
    return {
      background: `linear-gradient(135deg, ${from} 0%, ${to} 100%)`,
    }
  })
</script>

<template>
  <div class="dashboard-card" :class="cardClass">
    <div class="card-header">
      <div class="header-icon" :style="iconStyle">
        <slot name="icon">
          <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
            />
          </svg>
        </slot>
      </div>
      <div class="header-info">
        <h3 class="card-title">{{ title }}</h3>
        <div v-if="badge" class="header-badge" :class="badgeClass">
          {{ badge }}
        </div>
      </div>
      <div v-if="$slots.extra" class="header-extra">
        <slot name="extra" />
      </div>
    </div>

    <div class="card-content">
      <slot />
    </div>

    <div v-if="$slots.footer || showUpdateTime" class="card-footer">
      <slot name="footer">
        <span v-if="lastUpdated" class="update-time">{{ $t('admin.dashboard.common.updatedAt') }} {{ lastUpdated }}</span>
        <button
          v-if="refreshable"
          class="refresh-btn"
          :disabled="isRefreshing"
          :aria-busy="isRefreshing ? 'true' : 'false'"
          @click="$emit('refresh')"
        >
          <svg class="h-4 w-4" :class="{ 'animate-spin': isRefreshing }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
        </button>
      </slot>
    </div>
  </div>
</template>

<style scoped>
  .dashboard-card {
    @apply flex flex-col;
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-cyber-sm);
    overflow: hidden;
    transition:
      border-color var(--transition-normal) var(--ease-out),
      box-shadow var(--transition-normal) var(--ease-out),
      transform var(--transition-normal) var(--ease-out);
    position: relative;
    height: auto;
  }

  .dashboard-card:hover {
    border-color: var(--color-hover-border);
    box-shadow: var(--shadow-cyber-md);
    transform: translateY(-2px);
  }

  .card-small {
    min-height: 140px;
  }

  .card-medium {
    min-height: 165px;
  }

  .card-large {
    min-height: 240px;
  }

  .card-header {
    @apply flex items-center gap-3 px-4 py-2.5;
    border-bottom: 1px solid var(--color-border-subtle);
    min-height: 52px;
  }

  .header-icon {
    @apply flex h-10 w-10 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-badge-accent-text-rgb), 0.08));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    box-shadow: var(--shadow-sm);
  }

  .header-icon svg {
    @apply h-5 w-5;
  }

  .header-info {
    @apply flex flex-1 items-center justify-between;
  }

  .card-title {
    color: var(--color-content-heading);
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    line-height: 1.2;
    letter-spacing: var(--tracking-tight);
  }

  .header-badge {
    @apply px-2 py-0.5;
    border-radius: var(--radius-full);
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
  }

  .badge-success {
    background: var(--color-badge-success-bg);
    color: var(--color-badge-success-text);
    border: 1px solid var(--color-badge-success-border);
  }

  .badge-warning {
    background: var(--color-badge-warning-bg);
    color: var(--color-badge-warning-text);
    border: 1px solid var(--color-badge-warning-border);
  }

  .badge-danger {
    background: var(--color-badge-error-bg);
    color: var(--color-badge-error-text);
    border: 1px solid var(--color-badge-error-border);
  }

  .badge-info {
    background: var(--color-badge-primary-bg);
    color: var(--color-badge-primary-text);
    border: 1px solid var(--color-badge-primary-border);
  }

  .badge-normal {
    background: var(--color-badge-neutral-bg);
    color: var(--color-badge-neutral-text);
    border: 1px solid var(--color-badge-neutral-border);
  }

  .header-extra {
    @apply flex-shrink-0;
  }

  .card-content {
    @apply flex-1 px-4 pt-2.5;
    font-variant-numeric: tabular-nums;
  }

  .card-footer {
    @apply flex items-center justify-between px-4 py-2;
    border-top: 1px solid var(--color-border-subtle);
    background: rgba(var(--color-background-800-rgb), 0.5);
  }

  .update-time {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    font-family: var(--font-mono);
  }

  .refresh-btn {
    @apply flex h-5 w-5 items-center justify-center;
    border-radius: var(--radius-full);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    transition: all var(--transition-normal) var(--ease-out);
  }

  .refresh-btn svg {
    @apply h-3 w-3;
  }

  .refresh-btn:hover:not(:disabled) {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: var(--shadow-glow-sm);
  }

  .refresh-btn:disabled {
    opacity: var(--opacity-disabled);
    cursor: not-allowed;
    box-shadow: none;
    background: rgba(var(--color-brand-500-rgb), 0.08);
  }
</style>
