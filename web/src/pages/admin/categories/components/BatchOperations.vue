<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { deleteTemplate, updateTemplate, batchUpdateSortOrder } from '@/api/admin/category'
  import type { CategoryTemplate } from '@/api/types/category'
  import { showConfirm } from '@/utils/dialog'

  defineOptions({
    name: 'BatchOperations',
  })

  const props = defineProps<{
    selectedTemplates: number[]
    allTemplates: CategoryTemplate[]
  }>()

  const emit = defineEmits<{
    'clear-selection': []
    'batch-completed': []
  }>()

  const { $t } = useTexts()
  const toast = useToast()
  const loading = ref(false)

  const selectedCount = computed(() => props.selectedTemplates.length)
  const hasSelection = computed(() => selectedCount.value > 0)

  const selectedTemplateData = computed(() => {
    return props.allTemplates.filter((template) => props.selectedTemplates.includes(template.id))
  })

  /* 批量删除 */
  const handleBatchDelete = async () => {
    if (!hasSelection.value) {
      return
    }

    const confirmText = $t('admin.categories.confirm.deleteMultiple').replace('{count}', String(selectedCount.value))
    if (!showConfirm(confirmText)) {
      return
    }

    loading.value = true

    try {
      const deletePromises = props.selectedTemplates.map((id) => deleteTemplate({ id }))

      const results = await Promise.allSettled(deletePromises)

      const successCount = results.filter((result) => result.status === 'fulfilled' && result.value.success).length

      const failCount = results.length - successCount

      if (successCount > 0) {
        toast.success($t('admin.categories.batchOperations.deleteSuccess').replace('{count}', String(successCount)))
      }

      if (failCount > 0) {
        toast.warning($t('admin.categories.batchOperations.deleteFailed').replace('{count}', String(failCount)))
      }

      emit('clear-selection')
      emit('batch-completed')
    } catch {
      toast.error($t('admin.categories.batchOperations.deleteError'))
    }

    loading.value = false
  }

  const handleBatchSetPopular = async (isPopular: boolean) => {
    if (!hasSelection.value) {
      return
    }

    const actionText = isPopular
      ? $t('admin.categories.batchOperations.actionSetPopular')
      : $t('admin.categories.batchOperations.actionCancelPopular')
    const confirmText = $t('admin.categories.confirm.setPopular')
      .replace('{count}', String(selectedCount.value))
      .replace('{action}', actionText)

    if (!showConfirm(confirmText)) {
      return
    }

    loading.value = true

    try {
      const updatePromises = props.selectedTemplates.map((id) =>
        updateTemplate({
          id,
          is_popular: isPopular,
        })
      )

      const results = await Promise.allSettled(updatePromises)

      const successCount = results.filter((result) => result.status === 'fulfilled' && result.value.success).length

      const failCount = results.length - successCount

      if (successCount > 0) {
        toast.success(
          $t('admin.categories.batchOperations.setPopularSuccess')
            .replace('{action}', actionText)
            .replace('{count}', String(successCount))
        )
      }

      if (failCount > 0) {
        toast.warning(
          $t('admin.categories.batchOperations.setPopularFailed')
            .replace('{count}', String(failCount))
            .replace('{action}', actionText)
        )
      }

      emit('clear-selection')
      emit('batch-completed')
    } catch {
      toast.error($t('admin.categories.batchOperations.setPopularError').replace('{action}', actionText))
    }

    loading.value = false
  }

  const handleBatchReorder = async () => {
    if (!hasSelection.value) {
      return
    }

    const sortOrders = selectedTemplateData.value.map((template, index) => ({
      id: template.id,
      sort_order: (index + 1) * 10, // 设置为10, 20, 30...
    }))

    loading.value = true

    try {
      const result = await batchUpdateSortOrder({ sort_orders: sortOrders })

      if (result.success) {
        toast.success($t('admin.categories.batchOperations.reorderSuccess'))
        emit('batch-completed')
      }
    } catch {
      toast.error($t('admin.categories.batchOperations.reorderError'))
    }

    loading.value = false
  }
</script>

<template>
  <div v-if="hasSelection" class="batch-operations">
    <div class="batch-info">
      <div class="selection-info">
        <i class="fas fa-check-circle mr-2 text-cyan-400" />
        <span>
          {{ $t('admin.categories.batchOperations.selected') }} <strong>{{ selectedCount }}</strong>
          {{ $t('admin.categories.batchOperations.selectedCount') }}
        </span>
      </div>

      <div class="batch-actions">
        <CyberButton size="sm" type="secondary" :disabled="loading" @click="handleBatchSetPopular(true)">
          <i class="fas fa-star mr-1.5" />
          {{ $t('admin.categories.batchOperations.setPopular') }}
        </CyberButton>

        <CyberButton size="sm" type="secondary" :disabled="loading" @click="handleBatchSetPopular(false)">
          <i class="fas fa-star-half-alt mr-1.5" />
          {{ $t('admin.categories.batchOperations.cancelPopular') }}
        </CyberButton>

        <CyberButton size="sm" type="secondary" :disabled="loading" @click="handleBatchReorder">
          <i class="fas fa-sort-numeric-down mr-1.5" />
          {{ $t('admin.categories.batchOperations.reorder') }}
        </CyberButton>

        <CyberButton size="sm" type="danger" :loading="loading" @click="handleBatchDelete">
          <i class="fas fa-trash mr-1.5" />
          {{ $t('admin.categories.batchOperations.batchDelete') }}
        </CyberButton>

        <CyberButton size="sm" type="secondary" :disabled="loading" @click="emit('clear-selection')">
          <i class="fas fa-times mr-1.5" />
          {{ $t('admin.categories.batchOperations.cancelSelection') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .batch-operations {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    padding: var(--space-md);
    margin-bottom: var(--space-md);

    .batch-info {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: var(--space-md);
    }

    .selection-info {
      display: flex;
      align-items: center;
      color: var(--color-gray-300);
      font-size: var(--text-sm);

      strong {
        color: var(--color-cyan-400);
      }
    }

    .batch-actions {
      display: flex;
      gap: var(--space-sm);
      flex-wrap: wrap;
    }
  }

  @media (max-width: 768px) {
    .batch-operations {
      .batch-info {
        flex-direction: column;
        align-items: stretch;
        gap: var(--space-sm);
      }

      .selection-info {
        justify-content: center;
      }

      .batch-actions {
        justify-content: center;

        :deep(.cyber-button) {
          flex: 1;
          min-width: 0;
        }
      }
    }
  }
</style>
