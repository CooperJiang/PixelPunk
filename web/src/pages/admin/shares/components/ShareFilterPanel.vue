<script setup lang="ts">
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { AdminShareListParams } from '@/api/share/types'

  const { $t } = useTexts()

  const props = defineProps<{
    initialFilters?: Partial<AdminShareListParams>
  }>()

  const emit = defineEmits<{
    (e: 'filter', filters: AdminShareListParams): void
  }>()

  /* 筛选条件 */
  const filters = reactive<AdminShareListParams>({
    keyword: '',
    status: undefined,
    start_date: '',
    end_date: '',
    page: 1,
    size: 10,
    order_by: 'created_at DESC',
  })

  /* 日期状态 */
  const startDate = ref<Date | null>(null)
  const endDate = ref<Date | null>(null)

  /* 状态选项 */
  const statusOptions = computed(() => [
    { label: $t('admin.shares.filter.statusOptions.all'), value: undefined },
    { label: $t('admin.shares.filter.statusOptions.normal'), value: 1 },
    { label: $t('admin.shares.filter.statusOptions.expired'), value: 2 },
    { label: $t('admin.shares.filter.statusOptions.deleted'), value: 3 },
    { label: $t('admin.shares.filter.statusOptions.disabled'), value: 4 },
  ])

  /* 每页显示数量选项 */
  const pageSizeOptions = [
    { label: '10', value: '10' },
    { label: '20', value: '20' },
    { label: '50', value: '50' },
    { label: '100', value: '100' },
  ]

  /* 使用计算属性处理页面大小转换 */
  const pageSizeStr = computed({
    get: () => filters.size?.toString() || '10',
    set: (val: string) => {
      filters.size = parseInt(val)
    },
  })

  const handleDateChange = () => {
    if (startDate.value) {
      filters.start_date = startDate.value.toISOString()
    } else {
      filters.start_date = ''
    }

    if (endDate.value) {
      filters.end_date = endDate.value.toISOString()
    } else {
      filters.end_date = ''
    }
  }

  const applyFilters = () => {
    filters.page = 1
    emit('filter', { ...filters })
  }

  const resetFilters = () => {
    Object.assign(filters, {
      keyword: '',
      status: undefined,
      start_date: '',
      end_date: '',
      page: 1,
      size: 10,
      order_by: 'created_at DESC',
    })

    startDate.value = null
    endDate.value = null

    applyFilters()
  }

  onMounted(() => {
    if (props.initialFilters) {
      Object.assign(filters, props.initialFilters)

      if (props.initialFilters.start_date) {
        startDate.value = new Date(props.initialFilters.start_date)
      }
      if (props.initialFilters.end_date) {
        endDate.value = new Date(props.initialFilters.end_date)
      }
    }
  })
</script>

<template>
  <div class="share-filter-panel mb-4 rounded-xl border border-subtle bg-background-600 p-3">
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-base font-bold text-content-heading">
        <i class="fas fa-filter mr-1.5 text-content" />{{ $t('admin.shares.filter.title') }}
      </h3>
      <div class="flex space-x-1.5">
        <CyberButton type="secondary" class="px-2 py-1 text-xs" @click="resetFilters">
          <i class="fas fa-undo-alt mr-1" />{{ $t('admin.shares.filter.buttons.reset') }}
        </CyberButton>
        <CyberButton type="primary" class="px-2 py-1 text-xs" @click="applyFilters">
          <i class="fas fa-search mr-1" />{{ $t('admin.shares.filter.buttons.apply') }}
        </CyberButton>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-x-3 gap-y-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.shares.filter.fields.keyword') }}</label>
        <CyberInput
          v-model="filters.keyword"
          :placeholder="$t('admin.shares.filter.fields.keywordPlaceholder')"
          class="compact"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.shares.filter.fields.status') }}</label>
        <CyberDropdown
          v-model="filters.status"
          :options="statusOptions"
          :placeholder="$t('admin.shares.filter.fields.statusPlaceholder')"
          class="compact"
          clearable
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.shares.filter.fields.startDate') }}</label>
        <CyberDatePicker
          v-model="startDate"
          type="date"
          :placeholder="$t('admin.shares.filter.fields.startDatePlaceholder')"
          class="compact"
          @change="handleDateChange"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.shares.filter.fields.endDate') }}</label>
        <CyberDatePicker
          v-model="endDate"
          type="date"
          :placeholder="$t('admin.shares.filter.fields.endDatePlaceholder')"
          class="compact"
          @change="handleDateChange"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.shares.filter.fields.pageSize') }}</label>
        <CyberDropdown
          v-model="pageSizeStr"
          :options="pageSizeOptions"
          :placeholder="$t('admin.shares.filter.fields.pageSizePlaceholder')"
          class="compact"
          clearable
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  :deep(.cyber-dropdown.compact .cyber-dropdown-header) {
    min-height: var(--space-3xl);
    height: var(--space-3xl);
    padding: 0 var(--space-sm);
    font-size: var(--text-sm);
    display: flex;
    align-items: center;
    box-sizing: border-box;
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-selected) {
    font-size: var(--text-sm);
  }

  :deep(.cyber-input.compact) {
    height: var(--space-3xl);
    padding: 0 var(--space-sm);
    box-sizing: border-box;
    font-size: var(--text-sm);
    line-height: 1;
  }

  :deep(.cyber-date-picker.compact .cyber-date-picker-input) {
    height: var(--space-3xl);
    padding: 0 var(--space-sm);
    font-size: var(--text-sm);
    box-sizing: border-box;
  }

  .filter-item {
    margin-bottom: 0;
  }

  :deep(.cyber-dropdown-menu) {
    max-height: 200px;
    overflow-y: auto;
    padding: var(--space-xs);
    font-size: var(--text-sm);
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-option) {
    padding: var(--space-xs) var(--space-sm);
    font-size: var(--text-xs);
    min-height: var(--space-xl);
  }
</style>
