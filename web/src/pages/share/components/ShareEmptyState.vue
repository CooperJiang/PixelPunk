<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()
</script>

<template>
  <div class="empty-state">
    <div class="empty-icon">
      <i class="fas fa-folder-open" />
    </div>
    <h3>{{ $t('share.emptyState.title') }}</h3>
    <p>{{ $t('share.emptyState.description') }}</p>
  </div>
</template>

<style scoped lang="scss">
  .empty-state {
    text-align: center;
    padding: var(--space-3xl) var(--space-md);
    color: var(--color-content-muted);
    background: rgba(var(--color-background-800-rgb), 0.3);
    border-radius: var(--radius-md);
    border: 1.5px dashed rgba(var(--color-border-subtle-rgb), 0.2);
  }

  .empty-icon {
    font-size: var(--text-4xl);
    color: var(--color-content-muted);
    opacity: 0.2;
    margin-bottom: var(--space-md);
  }

  .empty-state h3 {
    font-size: var(--text-xl);
    margin-bottom: var(--space-sm);
    font-weight: 500;
    color: var(--color-content-heading);
  }

  .empty-state p {
    font-size: var(--text-sm);
    color: var(--color-content-muted);
  }
</style>
