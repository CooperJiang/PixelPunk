<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'BatchOperationToolbar',
  })

  defineProps<Props>()

  defineEmits<{
    'batch-recommend': [isRecommend: boolean]
    'batch-delete': []
    'cancel-select-mode': []
  }>()

  interface Props {
    selectedCount: number
    processing?: boolean
  }
</script>

<template>
  <div class="batch-operation-toolbar">
    <div class="selection-info">
      <span class="selection-count">
        <i class="fas fa-image mr-1" />
        {{ $t('admin.files.batch.selected', { count: selectedCount }) }}
      </span>
    </div>

    <div class="batch-actions">
      <CyberButton
        type="success"
        :disabled="selectedCount === 0 || processing"
        :loading="processing"
        size="small"
        @click="$emit('batch-recommend', true)"
      >
        <i class="fas fa-star mr-1" />{{ $t('admin.files.batch.recommend') }}
      </CyberButton>

      <CyberButton
        type="secondary"
        :disabled="selectedCount === 0 || processing"
        :loading="processing"
        size="small"
        @click="$emit('batch-recommend', false)"
      >
        <i class="far fa-star mr-1" />{{ $t('admin.files.batch.cancelRecommend') }}
      </CyberButton>

      <CyberButton
        type="danger"
        :disabled="selectedCount === 0 || processing"
        :loading="processing"
        size="small"
        @click="$emit('batch-delete')"
      >
        <i class="fas fa-trash mr-1" />{{ $t('admin.files.batch.batchDelete') }}
      </CyberButton>

      <CyberButton type="outlined" size="small" @click="$emit('cancel-select-mode')">
        <i class="fas fa-times mr-1" />{{ $t('common.cancel') }}
      </CyberButton>
    </div>
  </div>
</template>

<style scoped>
  .batch-operation-toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-md);
    padding: var(--space-md) var(--space-lg);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.05) 0%, rgba(var(--color-error-rgb), 0.03) 100%);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-md);
    backdrop-filter: blur(8px);
    margin-bottom: var(--space-lg);
  }

  .selection-info {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }

  .selection-count {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--color-content-default);
    display: flex;
    align-items: center;
    padding: var(--space-xs) var(--space-sm);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
  }

  .batch-actions {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }

  @media (max-width: 768px) {
    .batch-operation-toolbar {
      flex-direction: column;
      align-items: stretch;
      gap: var(--space-md);
    }

    .selection-info {
      justify-content: center;
    }

    .batch-actions {
      justify-content: center;
      flex-wrap: wrap;
    }
  }
</style>
