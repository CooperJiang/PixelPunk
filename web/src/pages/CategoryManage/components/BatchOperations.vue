<script setup lang="ts">
  import type { ImageCategory } from '@/api/types/category'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'BatchOperations',
  })

  defineProps<{
    selectedCategories: number[]
    allCategories: ImageCategory[]
  }>()

  const emit = defineEmits<{
    clearSelection: []
    batchDelete: []
    batchCompleted: []
  }>()

  const { $t } = useTexts()

  const handleClearSelection = () => {
    emit('clearSelection')
  }

  const handleBatchDelete = () => {
    emit('batchDelete')
  }
</script>

<template>
  <div class="batch-operations mb-4 p-3">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <span class="text-sm text-content">
          {{ $t('category.batch.selected', { count: selectedCategories.length }) }}
        </span>
        <button class="text-xs text-content-muted transition-colors hover:text-content" @click="handleClearSelection">
          <i class="fas fa-times mr-1" />
          {{ $t('category.batch.clearSelection') }}
        </button>
      </div>

      <div class="flex gap-2">
        <CyberButton type="danger" size="small" @click="handleBatchDelete">
          <i class="fas fa-trash mr-2" />
          {{ $t('category.batch.delete') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .batch-operations {
    background: rgba(var(--color-brand-500-rgb), 0.12);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    backdrop-filter: var(--backdrop-blur-md);
    box-shadow: var(--shadow-sm);
    animation: slideDown var(--transition-slow) var(--ease-out);
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
