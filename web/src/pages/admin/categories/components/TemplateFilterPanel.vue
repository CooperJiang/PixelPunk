<script setup lang="ts">
  import { computed, reactive } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { CATEGORY_SORT_BY_OPTIONS, SORT_ORDER_OPTIONS, CATEGORY_TYPE_OPTIONS, DEFAULT_CATEGORY_SORT } from '@/constants'

  defineOptions({
    name: 'TemplateFilterPanel',
  })

  const { $t } = useTexts()

  interface FilterParams {
    keyword: string
    is_popular?: boolean
    sort_by: string
    sort_order: string
  }

  const props = defineProps<{
    initialFilters: FilterParams
  }>()

  const emit = defineEmits<{
    filter: [filters: FilterParams]
  }>()

  const filters = reactive({
    keyword: props.initialFilters.keyword || '',
    is_popular: props.initialFilters.is_popular,
    sort_by: props.initialFilters.sort_by || DEFAULT_CATEGORY_SORT.sortBy,
    sort_order: props.initialFilters.sort_order || DEFAULT_CATEGORY_SORT.sortOrder,
  })

  const isPopularOptions = CATEGORY_TYPE_OPTIONS

  const sortByOptions = CATEGORY_SORT_BY_OPTIONS

  const sortOrderOptions = SORT_ORDER_OPTIONS

  const handleFilter = () => {
    emit('filter', { ...filters })
  }

  const handleReset = () => {
    filters.keyword = ''
    filters.is_popular = undefined
    filters.sort_by = DEFAULT_CATEGORY_SORT.sortBy
    filters.sort_order = DEFAULT_CATEGORY_SORT.sortOrder
    handleFilter()
  }
</script>

<template>
  <div class="template-filter-panel">
    <div class="filter-container">
      <div class="filter-row">
        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.categories.filterPanel.keyword') }}</label>
          <CyberInput
            v-model="filters.keyword"
            :placeholder="$t('admin.categories.filterPanel.keywordPlaceholder')"
            clearable
            @keyup.enter="handleFilter"
          />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.categories.filterPanel.type') }}</label>
          <CyberDropdown
            v-model="filters.is_popular"
            :placeholder="$t('admin.categories.filterPanel.typePlaceholder')"
            :options="isPopularOptions"
            clearable
          />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.categories.filterPanel.sortBy') }}</label>
          <CyberDropdown v-model="filters.sort_by" :options="sortByOptions" />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.categories.filterPanel.sortOrder') }}</label>
          <CyberDropdown v-model="filters.sort_order" :options="sortOrderOptions" />
        </div>
      </div>

      <div class="filter-actions">
        <CyberButton type="secondary" @click="handleReset">
          <i class="fas fa-undo mr-1.5" />
          {{ $t('admin.categories.filterPanel.reset') }}
        </CyberButton>
        <CyberButton type="primary" @click="handleFilter">
          <i class="fas fa-search mr-1.5" />
          {{ $t('admin.categories.filterPanel.filter') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .template-filter-panel {
    background: rgba(var(--color-background-800-rgb), 0.4);
    border-radius: var(--radius-md);
    padding: var(--space-sm) var(--space-md);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    margin-bottom: var(--space-sm);

    .filter-container {
      display: flex;
      flex-direction: column;
      gap: var(--space-sm);
    }

    .filter-row {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      gap: var(--space-sm);
      align-items: end;
    }

    .filter-item {
      display: flex;
      flex-direction: column;
      gap: var(--space-xs);
    }

    .filter-label {
      color: var(--color-cyber-light-600);
      font-size: var(--text-xs);
      font-weight: var(--font-medium);
    }

    .filter-actions {
      display: flex;
      gap: var(--space-sm);
      justify-content: flex-end;
    }
  }

  @media (max-width: 768px) {
    .template-filter-panel {
      .filter-row {
        grid-template-columns: 1fr;
      }
    }
  }
</style>
