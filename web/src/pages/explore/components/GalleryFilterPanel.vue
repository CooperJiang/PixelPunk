<script setup lang="ts">
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import type { ImageListParams } from '@/api/common'
  import { getFilterSortOptions, LAYOUT, getResolutionOptions } from '@/constants'
  import { usePaginatedTags } from '../composables/usePaginatedTags'
  import { usePaginatedColors } from '../composables/usePaginatedColors'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  /* 接收筛选条件作为props */
  const props = defineProps<{
    initialFilters?: Partial<ImageListParams>
    show?: boolean
  }>()

  const emit = defineEmits<{
    (e: 'filter', filters: ImageListParams): void
  }>()

  /* 筛选条件 */
  const filters = reactive<ImageListParams>({
    page: 1,
    size: LAYOUT.DEFAULT_PAGE_SIZE,
    sort: 'newest',
    keyword: '',
    tags: '',
    dominant_color: '',
    resolution: '',
    min_width: undefined,
    max_width: undefined,
    min_height: undefined,
    max_height: undefined,
  })

  /* 单独管理选中的颜色和标签，然后在筛选时转换为字符串 */
  const selectedColors = ref<string[]>([])
  const selectedTags = ref<string[]>([])

  /* 自定义展示颜色的组件 */
  const _formatSelectedColors = computed(() => {
    if (!selectedColors.value?.length) {
      return ''
    }

    return selectedColors.value.map((color) => {
      const colorOption = colorOptions.value?.find((opt) => opt.value === color)
      return colorOption ? colorOption : { value: color, label: color, color }
    })
  })

  const {
    tagOptions,
    loading: tagsLoading,
    loadingMore: tagsLoadingMore,
    hasMore: tagsHasMore,
    loadMore: loadMoreTags,
    searchTags,
    initialize: initializeTags,
  } = usePaginatedTags(20)

  const {
    colorOptions,
    loading: colorsLoading,
    loadingMore: colorsLoadingMore,
    hasMore: colorsHasMore,
    loadMore: loadMoreColors,
    searchColors,
    initialize: initializeColors,
  } = usePaginatedColors(20)

  const sortOptions = computed(() => getFilterSortOptions($t))
  const resolutionOptions = computed(() => getResolutionOptions($t))

  const handleTagLoadMore = () => {
    loadMoreTags()
  }

  const handleTagSearch = (query: string) => {
    searchTags(query)
  }

  const handleColorLoadMore = () => {
    loadMoreColors()
  }

  const handleColorSearch = (query: string) => {
    searchColors(query)
  }

  watch(
    () => props.initialFilters,
    (newFilters) => {
      if (newFilters) {
        Object.assign(filters, newFilters)

        if (newFilters.tags) {
          selectedTags.value = newFilters.tags.split(',')
        }

        if (newFilters.dominant_color) {
          selectedColors.value = newFilters.dominant_color.split(',')
        }
      }
    },
    { immediate: true, deep: true }
  )

  const applyFilters = () => {
    filters.dominant_color = selectedColors.value.join(',')
    filters.tags = selectedTags.value.join(',')
    emit('filter', { ...filters })
  }

  const resetFilters = () => {
    selectedColors.value = []
    selectedTags.value = []

    Object.assign(filters, {
      page: 1,
      size: LAYOUT.DEFAULT_PAGE_SIZE,
      sort: 'newest',
      keyword: '',
      tags: '',
      dominant_color: '',
      resolution: '',
      min_width: undefined,
      max_width: undefined,
      min_height: undefined,
      max_height: undefined,
    })

    emit('filter', { ...filters })
  }

  const isFilterDataLoaded = ref(false)
  const loadFilterData = async () => {
    if (!isFilterDataLoaded.value) {
      isFilterDataLoaded.value = true
      await Promise.all([initializeTags(), initializeColors()])
    }
  }

  watch(
    () => props.show,
    (show) => {
      if (show) {
        loadFilterData()
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    if (props.show) {
      loadFilterData()
    }
  })
</script>

<template>
  <div class="gallery-filter-panel-container">
    <div class="gallery-filter-panel mb-3 rounded-lg border border-subtle bg-background-600 p-2.5">
      <div class="mb-2 flex items-center justify-between">
        <h3 class="text-sm font-bold text-content-heading">
          <i class="fas fa-filter mr-1 text-content" />{{ $t('explore.filter.conditions') }}
        </h3>
        <div class="flex space-x-1">
          <CyberButton type="secondary" class="px-1.5 py-0.5 text-xs" @click="resetFilters">
            <i class="fas fa-undo-alt mr-0.5" />{{ $t('explore.search.reset') }}
          </CyberButton>
          <CyberButton type="primary" class="px-1.5 py-0.5 text-xs" @click="applyFilters">
            <i class="fas fa-search mr-0.5" />{{ $t('explore.actions.filter') }}
          </CyberButton>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-3 md:grid-cols-2 lg:grid-cols-3">
        <div class="filter-item">
          <label class="mb-0.5 block text-xs text-content-muted">{{ $t('explore.filter.keyword') }}</label>
          <CyberInput v-model="filters.keyword" :placeholder="$t('explore.filter.keywordPlaceholder')" />
        </div>

        <div class="filter-item">
          <label class="mb-0.5 block text-xs text-content-muted">{{ $t('explore.filter.tag') }}</label>
          <CyberDropdown
            v-model="selectedTags"
            :options="tagOptions"
            :placeholder="$t('explore.filter.tagPlaceholder')"
            multiple
            searchable
            pagination
            :loading="tagsLoading"
            :loading-more="tagsLoadingMore"
            :has-more="tagsHasMore"
            class="compact"
            @load-more="handleTagLoadMore"
            @search="handleTagSearch"
          />
        </div>

        <div class="filter-item">
          <label class="mb-0.5 block text-xs text-content-muted">{{ $t('explore.filter.color') }}</label>
          <CyberDropdown
            v-model="selectedColors"
            :options="colorOptions"
            :placeholder="$t('explore.filter.colorPlaceholder')"
            is-color-mode
            multiple
            searchable
            pagination
            :loading="colorsLoading"
            :loading-more="colorsLoadingMore"
            :has-more="colorsHasMore"
            class="compact"
            @load-more="handleColorLoadMore"
            @search="handleColorSearch"
          >
            <template #option="{ option, isSelected, onSelect }">
              <div class="cyber-dropdown-option color-option-row" :class="{ 'is-selected': isSelected }" @click="onSelect">
                <div class="checkbox-wrapper">
                  <span class="checkbox">
                    <i v-if="isSelected" class="fas fa-check" />
                  </span>
                </div>
                <div class="color-item h-6 w-6 rounded" :style="{ backgroundColor: option.color }" />
                <span class="ml-2 text-sm">{{ option.label }}</span>
              </div>
            </template>
          </CyberDropdown>
        </div>

        <!-- 排序方式 -->
        <div class="filter-item">
          <label class="mb-0.5 block text-xs text-content-muted">{{ $t('explore.filter.sort') }}</label>
          <CyberDropdown
            v-model="filters.sort"
            :options="sortOptions"
            :placeholder="$t('explore.filter.sortPlaceholder')"
            class="compact"
          />
        </div>

        <!-- 分辨率范围 -->
        <div class="filter-item">
          <label class="mb-0.5 block text-xs text-content-muted">{{ $t('explore.filter.size') }}</label>
          <div class="flex items-center gap-1">
            <CyberInput
              v-model.number="filters.min_width"
              :placeholder="$t('explore.filter.minWidth')"
              type="number"
              class="w-full"
            />
            <span class="text-content-content-disabled mx-0.5">×</span>
            <CyberInput
              v-model.number="filters.min_height"
              :placeholder="$t('explore.filter.minHeight')"
              type="number"
              class="w-full"
            />
          </div>
        </div>

        <!-- 分辨率选择 -->
        <div class="filter-item">
          <label class="mb-0.5 block text-xs text-content-muted">{{ $t('explore.filter.resolution') }}</label>
          <CyberDropdown
            v-model="filters.resolution"
            :options="resolutionOptions"
            :placeholder="$t('explore.filter.resolutionPlaceholder')"
            class="compact"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .cyber-dropdown-option.color-option-row {
    display: flex;
    align-items: center;
    padding: 4px 8px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .cyber-dropdown-option.color-option-row:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-dropdown-option.color-option-row.is-selected {
    background-color: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .checkbox-wrapper {
    display: flex;
    align-items: center;
    margin-right: 6px;
    min-width: 16px;
  }

  .checkbox {
    width: 16px;
    height: 16px;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.5);
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-brand-500);
    font-size: 10px;
  }

  .selected-colors-preview {
    padding: 3px 6px;
    border-radius: var(--radius-sm);
    background: var(--color-background-700);
  }

  .color-preview-item {
    padding: 1px 4px;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    transition: all 0.2s ease;
  }

  .color-preview-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .color-item {
    position: relative;
    overflow: hidden;
    border: 1px solid var(--color-border-default);
    box-shadow: var(--shadow-cyber-sm);
    transition: all 0.2s ease;
  }

  .color-item:hover {
    transform: scale(1.05);
    box-shadow: var(--shadow-cyber-md);
  }

  .color-item::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, transparent 0%, rgba(0, 0, 0, 0.2) 100%);
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-header) {
    min-height: 28px;
    padding: 0.2rem 0.6rem;
    font-size: 0.8rem;
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-selected) {
    font-size: 0.8rem;
    max-width: 100%;
    overflow-x: auto;
    scrollbar-width: thin;
    scrollbar-color: rgba(var(--color-brand-500-rgb), 0.5) transparent;
  }

  :deep(.selected-tags) {
    display: flex;
    flex-wrap: nowrap;
    overflow-x: auto;
    scrollbar-width: thin;
    scrollbar-color: rgba(var(--color-brand-500-rgb), 0.5) transparent;
    padding-bottom: 3px;
    max-width: 100%;
    margin-right: -6px;
  }

  :deep(.selected-tags::-webkit-scrollbar) {
    height: 3px;
  }

  :deep(.selected-tags::-webkit-scrollbar-track) {
    background: transparent;
  }

  :deep(.selected-tags::-webkit-scrollbar-thumb) {
    background-color: rgba(var(--color-brand-500-rgb), 0.5);
    border-radius: var(--radius-sm);
  }

  :deep(.cyber-dropdown.compact .selected-tag) {
    padding: 0px 4px;
    font-size: 0.7rem;
    white-space: nowrap;
    flex-shrink: 0;
    margin-right: 3px;
    display: flex;
    align-items: center;
    border-radius: var(--radius-sm);
    height: 18px;
  }

  :deep(.selected-tag .color-block) {
    width: 10px;
    height: 10px;
    border-radius: var(--radius-sm);
    margin-right: 3px;
    border: 1px solid rgba(var(--color-content-rgb), 0.3);
  }

  :deep(.selected-tag .remove-tag) {
    margin-left: 3px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 0.8rem;
    height: 12px;
    width: 12px;
  }

  :deep(.cyber-dropdown-menu) {
    max-height: 200px;
    overflow-y: auto;
    padding: 4px;
    font-size: 0.8rem;
  }

  :deep(.cyber-dropdown-search) {
    padding: 3px;
  }

  :deep(.cyber-dropdown-search input) {
    padding: 2px 6px;
    font-size: 0.75rem;
    min-height: 22px;
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-option) {
    padding: 3px 6px;
    font-size: 0.75rem;
    display: flex;
    align-items: center;
    width: 100%;
    min-height: 24px;
  }

  :deep(.cyber-dropdown.compact) {
    min-height: 30px;
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown__trigger) {
    min-height: 30px;
    height: 30px;
    padding: 0 0.6rem;
    font-size: 0.8rem;
    display: flex;
    align-items: center;
    box-sizing: border-box;
  }

  :deep(.cyber-input-wrapper) {
    min-height: 30px;
    box-sizing: border-box;
  }

  :deep(.cyber-input) {
    box-sizing: border-box;
    line-height: 1;
  }

  :deep(.filter-item .cyber-input) {
    font-size: 0.8rem;
    padding-top: 0;
    padding-bottom: 0;
    display: flex;
    align-items: center;
  }
</style>
