<script setup lang="ts">
  defineProps({
    visible: {
      type: Boolean,
      default: false,
    },
    selectedCount: {
      type: Number,
      default: 0,
    },
    totalCount: {
      type: Number,
      default: 0,
    },
  })

  const emit = defineEmits(['select-all', 'deselect-all', 'download'])

  const selectAll = () => {
    emit('select-all')
  }

  const deselectAll = () => {
    emit('deselect-all')
  }

  const downloadSelected = () => {
    emit('download')
  }
</script>

<template>
  <div v-if="visible" class="selection-toolbar">
    <div class="selection-info">
      <span
        >{{ $t('share.selectionToolbar.selectedPrefix') }} {{ selectedCount }}
        {{ $t('share.selectionToolbar.selectedSuffix') }}</span
      >
    </div>
    <div class="selection-actions">
      <button class="action-btn" :disabled="totalCount === 0" @click="selectAll">
        <i class="fas fa-check-double" />
        <span>{{ $t('share.selectionToolbar.selectAll') }}</span>
      </button>
      <button class="action-btn" :disabled="selectedCount === 0" @click="deselectAll">
        <i class="fas fa-times" />
        <span>{{ $t('actions.cancel') }}</span>
      </button>
      <button class="action-btn download-btn" :disabled="selectedCount === 0" @click="downloadSelected">
        <i class="fas fa-download" />
        <span v-if="selectedCount <= 1">{{ $t('actions.download') }}</span>
        <span v-else>{{ $t('share.selectionToolbar.batchDownload') }} ({{ selectedCount }})</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
  .selection-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 1rem;
    background: rgba(var(--color-background-800-rgb), 0.6);
    border-radius: var(--radius-md);
    margin-top: var(--space-md);
  }

  .selection-info {
    color: rgba(var(--color-content-rgb), 0.8);
    font-size: 0.85rem;
  }

  .selection-actions {
    display: flex;
    gap: var(--space-sm);
  }

  .action-btn {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-md);
    color: var(--color-content-default);
    font-size: 0.75rem;
    padding: var(--space-md) 0.5rem;
    transition: all var(--transition-fast) var(--ease-in-out);
    cursor: pointer;
  }

  .action-btn:hover:not(:disabled) {
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .action-btn:active:not(:disabled) {
    transform: translateY(1px);
  }

  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .action-btn i {
    margin-right: 0.25rem;
  }

  .download-btn {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .download-btn:hover:not(:disabled) {
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .download-btn:active:not(:disabled) {
    transform: translateY(1px);
  }

  @media (max-width: 767px) {
    .selection-toolbar {
      flex-direction: column;
      gap: var(--space-sm);
      align-items: flex-start;
    }

    .selection-actions {
      width: 100%;
      justify-content: space-between;
    }
  }
</style>
