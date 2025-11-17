<script setup lang="ts">
  import { reactive } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'ReviewFilterPanel',
  })

  const { $t } = useTexts()
  const emit = defineEmits<Emits>()

  interface Filters {
    keyword: string
    sort: string
    nsfw_only: boolean
  }

  interface Emits {
    (e: 'filter-change', filters: Filters): void
  }

  const filters = reactive<Filters>({
    keyword: '',
    sort: 'newest',
    nsfw_only: false,
  })

  const sortOptions = [
    { label: $t('admin.contentReview.filterPanel.sortOptions.newest'), value: 'newest' },
    { label: $t('admin.contentReview.filterPanel.sortOptions.oldest'), value: 'oldest' },
    { label: $t('admin.contentReview.filterPanel.sortOptions.nameAsc'), value: 'name_asc' },
    { label: $t('admin.contentReview.filterPanel.sortOptions.nameDesc'), value: 'name_desc' },
    { label: $t('admin.contentReview.filterPanel.sortOptions.size'), value: 'size' },
  ]

  const handleFilterChange = () => {
    emit('filter-change', { ...filters })
  }

  const resetFilters = () => {
    filters.keyword = ''
    filters.sort = 'newest'
    filters.nsfw_only = false
    handleFilterChange()
  }

  defineExpose({
    resetFilters,
  })
</script>

<template>
  <div class="filter-panel">
    <div class="filter-header">
      <h3 class="filter-title">
        <i class="fas fa-filter mr-2 text-content" />
        {{ $t('admin.contentReview.filterPanel.title') }}
      </h3>
    </div>

    <div class="filter-content">
      <div class="filter-group">
        <label class="filter-label">{{ $t('admin.contentReview.filterPanel.keywordLabel') }}</label>
        <CyberInput
          v-model="filters.keyword"
          :placeholder="$t('admin.contentReview.filterPanel.keywordPlaceholder')"
          size="small"
          @input="handleFilterChange"
        >
          <template #prefix>
            <i class="fas fa-search text-content-disabled" />
          </template>
        </CyberInput>
      </div>

      <div class="filter-group">
        <label class="filter-label">{{ $t('admin.contentReview.filterPanel.sortLabel') }}</label>
        <CyberDropdown
          v-model="filters.sort"
          :options="sortOptions"
          :placeholder="$t('admin.contentReview.filterPanel.sortPlaceholder')"
          size="small"
          @change="handleFilterChange"
        />
      </div>

      <div class="filter-group">
        <label class="filter-label">{{ $t('admin.contentReview.filterPanel.contentLabel') }}</label>
        <div class="flex items-center space-x-3">
          <cyberCheckbox v-model="filters.nsfw_only" @change="handleFilterChange">
            {{ $t('admin.contentReview.filterPanel.nsfwOnly') }}
          </cyberCheckbox>
        </div>
      </div>

      <div class="filter-actions">
        <CyberButton type="secondary" size="small" @click="resetFilters">
          <i class="fas fa-undo mr-1" />
          {{ $t('admin.contentReview.filterPanel.resetButton') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .filter-panel {
    border: 1px solid var(--color-border-subtle);
    overflow: hidden;
    border-radius: var(--radius-sm);
    background: var(--color-background-600);
  }

  .filter-header {
    border-bottom: 1px solid var(--color-border-subtle);
    padding: var(--space-md);
  }

  .filter-title {
    display: flex;
    align-items: center;
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--color-white);
  }

  .filter-content {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    padding: var(--space-md);
  }

  .filter-group {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
  }

  .filter-label {
    display: block;
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: rgba(255, 255, 255, 0.8);
  }

  .cyber-checkbox {
    height: var(--space-md);
    width: var(--space-md);
    border-radius: var(--radius-sm);
    border: 2px solid var(--color-brand-600);
    background: transparent;

    &:checked {
      border-color: var(--color-brand-500);
      background: var(--color-brand-500);
    }

    &:focus {
      ring: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
    }
  }

  .filter-actions {
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    padding-top: var(--space-xs);
  }
</style>
