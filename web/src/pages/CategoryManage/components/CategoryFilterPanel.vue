<script setup lang="ts">
  import { ref, computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'CategoryFilterPanel',
  })

  const { $t } = useTexts()

  const props = defineProps<{
    initialFilters: {
      keyword: string
      status: string
      sort_by: string
      sort_order_dir: string
    }
  }>()

  const emit = defineEmits<{
    filter: [filters: typeof props.initialFilters]
  }>()

  const filters = ref({ ...props.initialFilters })

  const statusOptions = computed(() => [
    { label: $t('category.filter.status.all'), value: '' },
    { label: $t('category.filter.status.active'), value: 'active' },
    { label: $t('category.filter.status.archived'), value: 'archived' },
  ])

  const sortOptions = computed(() => [
    { label: $t('category.table.columns.sortOrder'), value: 'sort_order' },
    { label: $t('category.table.columns.name'), value: 'name' },
    { label: $t('category.table.columns.fileCount'), value: 'file_count' },
    { label: $t('category.table.columns.createdAt'), value: 'created_at' },
  ])

  const handleSearch = () => {
    emit('filter', { ...filters.value })
  }

  const handleReset = () => {
    filters.value = {
      keyword: '',
      status: '',
      sort_by: 'sort_order',
      sort_order_dir: 'asc',
    }
    emit('filter', { ...filters.value })
  }
</script>

<template>
  <div class="filter-panel mb-4 p-3">
    <div class="flex flex-wrap items-center gap-3">
      <div class="filter-item min-w-[200px] max-w-[300px] flex-1">
        <CyberInput
          v-model="filters.keyword"
          :placeholder="$t('category.placeholders.search')"
          prefix-icon="search"
          @keyup.enter="handleSearch"
        />
      </div>

      <div class="filter-item min-w-[150px]">
        <CyberDropdown v-model="filters.status" :options="statusOptions" :placeholder="$t('category.placeholders.status')" />
      </div>

      <div class="filter-item min-w-[150px]">
        <CyberDropdown v-model="filters.sort_by" :options="sortOptions" :placeholder="$t('category.filter.sortBy')" />
      </div>

      <div class="filter-actions ml-auto flex items-center gap-2">
        <CyberButton type="text" size="small" @click="handleReset">
          <i class="fas fa-redo mr-1" />
          {{ $t('category.filter.reset') }}
        </CyberButton>
        <CyberButton type="primary" size="small" @click="handleSearch">
          <i class="fas fa-check mr-1" />
          {{ $t('category.filter.apply') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  :global(.filter-slide-enter-active),
  :global(.filter-slide-leave-active) {
    transition: all var(--transition-slow) var(--ease-in-out);
  }

  :global(.filter-slide-enter-from) {
    opacity: 0;
    transform: translateY(-10px);
  }

  :global(.filter-slide-leave-to) {
    opacity: 0;
    transform: translateY(-10px);
  }

  .filter-panel {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-sm);
    backdrop-filter: var(--backdrop-blur-md);

    :deep(.cyber-input),
    :deep(.cyber-dropdown) {
      width: 100%;
    }
  }
</style>
