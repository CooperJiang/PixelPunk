<script setup lang="ts">
  import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { getAdminColorList, getAdminTagList, type ImageListParams } from '@/api/admin/files'
  import { getAdminSortOptions, LAYOUT, getResolutionOptions, getStorageTypeOptions } from '@/constants'

  const props = defineProps({
    initialFilters: {
      type: Object,
      default: () => ({}),
    },
  })

  const emit = defineEmits<{
    (e: 'filter', filters: ImageListParams): void
  }>()

  const { $t } = useTexts()

  /* 筛选条件 */
  const filters = reactive<ImageListParams>({
    page: 1,
    size: LAYOUT.DEFAULT_PAGE_SIZE,
    sort: 'newest',
    keyword: '',
    tags: '',
    dominant_color: '',
    resolution: '',
    is_nsfw: false,
    is_recommended: false,
    storage_type: '',
    min_width: undefined,
    max_width: undefined,
    min_height: undefined,
    max_height: undefined,
    user_id: undefined,
    status: undefined,
  })

  /* 单独管理选中的颜色和标签，然后在筛选时转换为字符串 */
  const selectedColors = ref<string[]>([])
  const selectedTags = ref<string[]>([])

  /* 监听选中颜色和标签变化，自动滚动到最新位置 */
  watch(
    selectedColors,
    (newVal, oldVal) => {
      if (newVal.length > oldVal.length) {
        nextTick(() => {
          scrollToLatestSelection('.selected-tags')
        })
      }
    },
    { deep: true }
  )

  watch(
    selectedTags,
    (newVal, oldVal) => {
      if (newVal.length > oldVal.length) {
        nextTick(() => {
          scrollToLatestSelection('.selected-tags')
        })
      }
    },
    { deep: true }
  )

  const scrollToLatestSelection = (selector: string) => {
    const containers = document.querySelectorAll(selector)
    containers.forEach((container) => {
      if (container) {
        container.scrollLeft = container.scrollWidth
      }
    })
  }

  const _formatSelectedColors = computed(() => {
    if (!selectedColors.value.length) {
      return ''
    }

    return selectedColors.value.map((color) => {
      const colorOption = colorOptions.value.find((opt) => opt.value === color)
      return colorOption ? colorOption : { value: color, label: color, color }
    })
  })

  const showAdvanced = ref(false)

  const tags = ref<{ id: number; name: string; count: number }[]>([])
  const tagOptions = computed(() =>
    (tags.value || []).map((tag) => ({
      label: `${tag.name} (${tag.count})`,
      value: tag.id.toString(), // 使用ID作为值
    }))
  )

  const colors = ref<string[]>([])
  const colorOptions = computed(() =>
    (colors.value || []).map((color) => ({
      label: color,
      value: color,
      color,
    }))
  )

  const sortOptions = computed(() => getAdminSortOptions($t))
  const resolutionOptions = computed(() => getResolutionOptions($t))
  const storageOptions = computed(() => getStorageTypeOptions($t))

  const statusOptions = computed(() => [
    { label: $t('admin.files.filter.statusOptions.active'), value: 'active' },
    { label: $t('admin.files.filter.statusOptions.pendingReview'), value: 'pending_review' },
    { label: $t('admin.files.filter.statusOptions.deleted'), value: 'deleted' },
  ])

  const loadTags = async () => {
    try {
      const result = await getAdminTagList()
      if (result.success) {
        tags.value = result.data?.items || []
      }
    } catch (_error) {
      tags.value = []
    }
  }

  const loadColors = async () => {
    try {
      const result = await getAdminColorList()
      if (result.success) {
        colors.value = result.data?.items || []
      }
    } catch (_error) {
      colors.value = []
    }
  }

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
      is_nsfw: false,
      is_recommended: false,
      storage_type: '',
      min_width: undefined,
      max_width: undefined,
      min_height: undefined,
      max_height: undefined,
      user_id: undefined,
    })
    applyFilters()
  }

  onMounted(async () => {
    await Promise.all([loadTags(), loadColors()])

    if (props.initialFilters) {
      Object.assign(filters, props.initialFilters)

      if (filters.tags) {
        const tagIds = filters.tags.split(',')
        selectedTags.value = tagIds
      }

      if (filters.dominant_color) {
        selectedColors.value = filters.dominant_color.split(',')
      }
    }
  })
</script>

<template>
  <div class="image-filter-panel mb-4 rounded-xl border border-subtle bg-background-600 p-3">
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-base font-bold text-content-heading">
        <i class="fas fa-filter mr-1.5 text-content" />{{ $t('admin.files.filter.title') }}
      </h3>
      <div class="flex space-x-1.5">
        <CyberButton type="secondary" class="px-2 py-1 text-xs" @click="resetFilters">
          <i class="fas fa-undo-alt mr-1" />{{ $t('admin.files.buttons.reset') }}
        </CyberButton>
        <CyberButton type="primary" class="px-2 py-1 text-xs" @click="applyFilters">
          <i class="fas fa-search mr-1" />{{ $t('admin.files.buttons.applyFilter') }}
        </CyberButton>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-x-3 gap-y-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6">
      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.keyword') }}</label>
        <CyberInput v-model="filters.keyword" :placeholder="$t('admin.files.filter.keywordPlaceholder')" />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.tags') }}</label>
        <CyberDropdown
          v-model="selectedTags"
          :options="tagOptions"
          :placeholder="$t('admin.files.filter.tagsPlaceholder')"
          multiple
          searchable
          class="compact"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.dominantColor') }}</label>
        <CyberDropdown
          v-model="selectedColors"
          :options="colorOptions"
          :placeholder="$t('admin.files.filter.colorPlaceholder')"
          is-color-mode
          multiple
          class="compact"
        >
          <template #option="{ option, isSelected, onSelect }">
            <div class="cyber-dropdown-option color-option-row" :class="{ 'is-selected': isSelected }" @click="onSelect">
              <div class="checkbox-wrapper">
                <span class="checkbox">
                  <i v-show="isSelected" class="fas fa-check" />
                </span>
              </div>
              <div class="color-item h-6 w-6 rounded" :style="{ backgroundColor: option.color }" />
              <span class="ml-2 text-sm">{{ option.label }}</span>
            </div>
          </template>
        </CyberDropdown>
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.sortBy') }}</label>
        <CyberDropdown
          v-model="filters.sort"
          :options="sortOptions"
          :placeholder="$t('admin.files.filter.sortPlaceholder')"
          class="compact"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.resolution') }}</label>
        <CyberDropdown
          v-model="filters.resolution"
          :options="resolutionOptions"
          :placeholder="$t('admin.files.filter.resolutionPlaceholder')"
          class="compact"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.status') }}</label>
        <CyberDropdown
          v-model="filters.status"
          :options="statusOptions"
          :placeholder="$t('admin.files.filter.statusPlaceholder')"
          class="compact"
          clearable
        />
      </div>

      <div class="filter-item col-span-1 flex sm:col-span-2 md:col-span-1">
        <div class="flex h-full flex-col justify-center">
          <label class="mb-1.5 block text-xs text-content">{{ $t('admin.files.filter.quickOptions') }}</label>
          <div class="flex space-x-4">
            <div class="flex items-center">
              <label class="mr-2 text-xs text-content">{{ $t('admin.files.filter.isNsfw') }}</label>
              <CyberSwitch v-model="filters.is_nsfw" class="scale-75" />
            </div>
            <div class="flex items-center">
              <label class="mr-2 text-xs text-content">{{ $t('admin.files.filter.isRecommended') }}</label>
              <CyberSwitch v-model="filters.is_recommended" class="scale-75" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-2">
      <div
        class="flex cursor-pointer items-center text-xs text-content transition-colors hover:text-error-500"
        @click="showAdvanced = !showAdvanced"
      >
        <i class="fas mr-1" :class="[showAdvanced ? 'fa-chevron-up' : 'fa-chevron-down']" />
        <span>{{ showAdvanced ? $t('admin.files.filter.hideAdvanced') : $t('admin.files.filter.showAdvanced') }}</span>
      </div>

      <div
        v-if="showAdvanced"
        class="mt-2 grid grid-cols-1 gap-x-3 gap-y-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6"
      >
        <div class="filter-item">
          <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.advanced.userId') }}</label>
          <CyberInput
            v-model.number="filters.user_id"
            :placeholder="$t('admin.files.filter.advanced.userIdPlaceholder')"
            type="number"
          />
        </div>

        <div class="filter-item col-span-1 sm:col-span-2 md:col-span-1">
          <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.advanced.size') }}</label>
          <div class="flex items-center gap-1">
            <CyberInput
              v-model.number="filters.min_width"
              :placeholder="$t('admin.files.filter.advanced.minWidth')"
              type="number"
              class="w-full"
            />
            <span class="text-content-content-disabled mx-1">×</span>
            <CyberInput
              v-model.number="filters.min_height"
              :placeholder="$t('admin.files.filter.advanced.minHeight')"
              type="number"
              class="w-full"
            />
          </div>
        </div>

        <div class="filter-item">
          <label class="mb-0.5 block text-xs text-content">{{ $t('admin.files.filter.advanced.storageType') }}</label>
          <CyberDropdown
            v-model="filters.storage_type"
            :options="storageOptions"
            :placeholder="$t('admin.files.filter.advanced.storageTypePlaceholder')"
            class="compact"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .color-option-row {
    display: flex;
    align-items: center;
    padding: var(--space-xs) var(--space-sm);
    cursor: pointer;
    transition: all var(--transition-normal) var(--ease-out);
  }

  .color-option-row:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .color-option-row.is-selected {
    background-color: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .checkbox-wrapper {
    display: flex;
    align-items: center;
    margin-right: var(--space-sm);
    min-width: var(--space-md);
  }

  .checkbox {
    display: flex;
    justify-content: center;
    align-items: center;
    width: var(--space-md);
    height: var(--space-md);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.5);
    background-color: rgba(var(--color-brand-500-rgb), 0.1);
    font-size: var(--text-xs);
    color: var(--color-white);
  }

  :deep(.cyber-dropdown-option) {
    display: flex;
    align-items: center;
    width: 100%;
  }

  :deep(.cyber-dropdown-menu) {
    max-height: 220px;
    overflow-y: auto;
  }

  .color-item {
    position: relative;
    overflow: hidden;
    border: 1px solid rgba(255, 255, 255, 0.3);
    box-shadow: var(--shadow-sm);
    transition: all var(--transition-normal) var(--ease-out);
  }

  .color-item:hover {
    transform: scale(1.05);
    box-shadow: var(--shadow-md);
  }

  .color-item::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, transparent 0%, rgba(0, 0, 0, 0.2) 100%);
  }

  .selected-colors-preview {
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-sm);
    background: rgba(0, 0, 0, 0.2);
  }

  .color-preview-item {
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    transition: all var(--transition-normal) var(--ease-out);
  }

  .color-preview-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .hover\:shadow-glow:hover {
    box-shadow: var(--shadow-glow-md);
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-header) {
    min-height: var(--space-2xl);
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-sm);
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-selected) {
    font-size: var(--text-sm);
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
    padding-bottom: var(--space-xs);
    max-width: 100%;
    margin-right: calc(var(--space-sm) * -1);
  }

  :deep(.selected-tags::-webkit-scrollbar) {
    height: var(--space-xs);
  }

  :deep(.selected-tags::-webkit-scrollbar-track) {
    background: transparent;
  }

  :deep(.selected-tags::-webkit-scrollbar-thumb) {
    background-color: rgba(var(--color-brand-500-rgb), 0.5);
    border-radius: var(--radius-sm);
  }

  :deep(.cyber-dropdown.compact .selected-tag) {
    padding: 0 var(--space-xs);
    font-size: var(--text-xs);
    white-space: nowrap;
    flex-shrink: 0;
    margin-right: var(--space-xs);
    display: flex;
    align-items: center;
    border-radius: var(--radius-xs);
    height: var(--space-lg);
  }

  :deep(.selected-tag .color-block) {
    width: var(--space-sm);
    height: var(--space-sm);
    border-radius: var(--radius-xs);
    margin-right: var(--space-xs);
    border: 1px solid rgba(255, 255, 255, 0.3);
  }

  :deep(.selected-tag .remove-tag) {
    margin-left: var(--space-xs);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: var(--text-sm);
    height: var(--space-sm);
    width: var(--space-sm);
  }

  :deep(.cyber-dropdown-menu) {
    max-height: 200px;
    overflow-y: auto;
    padding: var(--space-xs);
    font-size: var(--text-sm);
  }

  :deep(.cyber-dropdown-search) {
    padding: var(--space-xs);
  }

  :deep(.cyber-dropdown-search input) {
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-xs);
    min-height: var(--space-xl);
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-option) {
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-xs);
    display: flex;
    align-items: center;
    width: 100%;
    min-height: var(--space-xl);
  }

  :deep(.cyber-dropdown-option .checkbox-wrapper) {
    margin-right: var(--space-sm);
  }

  :deep(.cyber-dropdown-option .checkbox) {
    width: var(--space-md);
    height: var(--space-md);
    font-size: var(--text-xs);
  }

  :deep(.cyber-dropdown-option .color-block) {
    width: var(--space-md);
    height: var(--space-md);
    margin-right: var(--space-xs);
  }

  :deep(.cyber-dropdown-option .option-label) {
    font-size: var(--text-xs);
  }

  .filter-item {
    margin-bottom: 0;
  }

  :deep(.cyber-dropdown.compact) {
    height: var(--space-3xl);
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-header) {
    min-height: var(--space-3xl);
    height: var(--space-3xl);
    padding: 0 var(--space-sm);
    font-size: var(--text-sm);
    display: flex;
    align-items: center;
    box-sizing: border-box;
  }

  :deep(.cyber-input) {
    padding: 0 var(--space-sm);
    box-sizing: border-box;
    line-height: 1;
  }

  :deep(.filter-item .cyber-input) {
    font-size: var(--text-sm);
    padding-top: 0;
    padding-bottom: 0;
    display: flex;
    align-items: center;
  }
</style>
