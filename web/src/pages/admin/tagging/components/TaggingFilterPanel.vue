<script setup lang="ts">
  import { reactive, watch } from 'vue'
  import type { TaggingStatusParams } from '@/api/admin/tagging'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  interface Props {
    initialFilters?: Partial<TaggingStatusParams>
  }

  const props = withDefaults(defineProps<Props>(), {
    initialFilters: () => ({}),
  })

  const emit = defineEmits<{
    (e: 'filter', filters: TaggingStatusParams): void
  }>()

  const statusOptions = computed(() => [
    { value: '', label: $t('admin.tagging.status.all') },
    { value: 'none', label: $t('admin.tagging.status.none') },
    { value: 'pending', label: $t('admin.tagging.status.pending') },
    { value: 'done', label: $t('admin.tagging.status.done') },
    { value: 'failed', label: $t('admin.tagging.status.failed') },
    { value: 'skipped', label: $t('admin.tagging.status.skipped') },
    { value: 'ignored', label: $t('admin.tagging.status.ignored') },
  ])

  const orderByOptions = computed(() => [
    { value: 'created_at', label: $t('admin.tagging.filter.createdAt') },
    { value: 'updated_at', label: $t('admin.tagging.filter.updatedAt') },
    { value: 'ai_tagging_tries', label: $t('admin.tagging.filter.tries') },
    { value: 'ai_tagging_status', label: $t('admin.tagging.filter.taggingStatus') },
  ])

  const orderOptions = computed(() => [
    { value: 'desc', label: $t('admin.tagging.filter.desc') },
    { value: 'asc', label: $t('admin.tagging.filter.asc') },
  ])

  const formState = reactive<TaggingStatusParams>({
    status: '',
    page: 1,
    limit: 20,
    order_by: 'created_at',
    order: 'desc',
    ...props.initialFilters,
  })

  watch(
    () => props.initialFilters,
    (newVal) => {
      if (newVal) {
        Object.assign(formState, newVal)
      }
    },
    { deep: true }
  )

  const applyFilter = () => {
    emit('filter', { ...formState })
  }

  const resetForm = () => {
    Object.assign(formState, {
      status: '',
      page: 1,
      limit: 20,
      order_by: 'created_at',
      order: 'desc',
    })

    emit('filter', { ...formState })
  }
</script>

<template>
  <div class="filter-panel">
    <div class="filter-grid">
      <div class="filter-item">
        <label class="filter-label">{{ $t('admin.tagging.filter.status') }}</label>
        <CyberDropdown v-model="formState.status" class="w-full" :options="statusOptions" />
      </div>

      <div class="filter-item">
        <label class="filter-label">{{ $t('admin.tagging.filter.orderBy') }}</label>
        <CyberDropdown v-model="formState.order_by" class="w-full" :options="orderByOptions" />
      </div>

      <div class="filter-item">
        <label class="filter-label">{{ $t('admin.tagging.filter.order') }}</label>
        <CyberDropdown v-model="formState.order" class="w-full" :options="orderOptions" />
      </div>

      <div class="filter-actions">
        <CyberButton type="primary" class="apply-btn" @click="applyFilter">
          <i class="fas fa-check mr-1.5" />{{ $t('admin.tagging.filter.apply') }}
        </CyberButton>
        <CyberButton type="text" class="reset-btn" @click="resetForm">
          <i class="fas fa-sync-alt mr-1.5" />{{ $t('admin.tagging.filter.reset') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .filter-panel {
    border: 1px solid var(--color-border-subtle);
    margin-bottom: var(--space-sm);
    border-radius: var(--radius-lg);
    background: var(--color-background-600);
    padding: var(--space-md);
    box-shadow: var(--shadow-md);
    backdrop-filter: var(--backdrop-blur-md);
  }

  .filter-grid {
    display: grid;
    gap: var(--space-md);
    grid-template-columns: repeat(1, 1fr);

    @media (min-width: 768px) {
      grid-template-columns: repeat(4, 1fr);
    }
  }

  .filter-item {
    display: flex;
    flex-direction: column;
  }

  .filter-label {
    margin-bottom: var(--space-xs);
    display: block;
    font-size: var(--text-xs);
    color: var(--color-brand-900);
  }

  .filter-actions {
    display: flex;
    align-items: flex-end;

    .apply-btn {
      margin-right: var(--space-xs);
      flex: 1;
      font-size: var(--text-xs);
    }

    .reset-btn {
      font-size: var(--text-xs);
    }
  }
</style>
