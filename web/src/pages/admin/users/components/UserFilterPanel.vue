<script setup lang="ts">
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  interface UserFilterParams {
    keyword: string
    status: string
    role: string
    sort_by: string
    start_date: string
    end_date: string
  }

  const props = defineProps<{
    initialFilters?: Partial<UserFilterParams>
  }>()

  const emit = defineEmits<{
    (e: 'filter', filters: UserFilterParams): void
  }>()

  /* Filter conditions */
  const filters = reactive<UserFilterParams>({
    keyword: '',
    status: '0',
    role: '0',
    sort_by: '',
    start_date: '',
    end_date: '',
  })

  /* Date range */
  const dateRange = ref<[Date | null, Date | null]>([null, null])

  /* Multi-select state management */
  const selectedStatuses = ref<string[]>([])
  const selectedRoles = ref<string[]>([])

  /* Status options */
  const statusOptions = computed(() => [
    { label: $t('admin.users.filter.status.all'), value: '0' },
    { label: $t('admin.users.filter.status.normal'), value: '1' },
    { label: $t('admin.users.filter.status.disabled'), value: '2' },
    { label: $t('admin.users.filter.status.deleted'), value: '3' },
  ])

  /* Role options */
  const roleOptions = computed(() => [
    { label: $t('admin.users.filter.role.all'), value: '0' },
    { label: $t('admin.users.filter.role.superAdmin'), value: '1' },
    // { label: $t('admin.users.filter.role.admin'), value: '2' },
    { label: $t('admin.users.filter.role.user'), value: '3' },
  ])

  /* Sort options */
  const sortOptions = computed(() => [
    { label: $t('admin.users.filter.sort.default'), value: '' },
    { label: $t('admin.users.filter.sort.newest'), value: 'created_at_desc' },
    { label: $t('admin.users.filter.sort.oldest'), value: 'created_at_asc' },
    { label: $t('admin.users.filter.sort.lastActivityNewest'), value: 'last_activity_desc' },
    { label: $t('admin.users.filter.sort.lastActivityOldest'), value: 'last_activity_asc' },
    { label: $t('admin.users.filter.sort.usernameAZ'), value: 'username_asc' },
    { label: $t('admin.users.filter.sort.usernameZA'), value: 'username_desc' },
  ])

  const handleDateRangeChange = (dates: [Date | null, Date | null]) => {
    if (dates && dates.length === 2 && dates[0] && dates[1]) {
      filters.start_date = formatDateToString(dates[0])
      filters.end_date = formatDateToString(dates[1])
    } else {
      filters.start_date = ''
      filters.end_date = ''
    }
  }

  const formatDateToString = (date: Date): string => {
    const year = date.getFullYear()
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const day = date.getDate().toString().padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  const toggleStatusFilter = (value: string, checked: boolean) => {
    if (checked) {
      if (!selectedStatuses.value.includes(value)) {
        selectedStatuses.value.push(value)
      }
    } else {
      const index = selectedStatuses.value.indexOf(value)
      if (index > -1) {
        selectedStatuses.value.splice(index, 1)
      }
    }

    if (selectedStatuses.value.length === 0) {
      filters.status = '0' // All
    } else if (selectedStatuses.value.length === 1) {
      filters.status = selectedStatuses.value[0]
    } else {
      filters.status = selectedStatuses.value.join(',')
    }
  }

  const toggleRoleFilter = (value: string, checked: boolean) => {
    if (checked) {
      if (!selectedRoles.value.includes(value)) {
        selectedRoles.value.push(value)
      }
    } else {
      const index = selectedRoles.value.indexOf(value)
      if (index > -1) {
        selectedRoles.value.splice(index, 1)
      }
    }

    if (selectedRoles.value.length === 0) {
      filters.role = '0' // All
    } else if (selectedRoles.value.length === 1) {
      filters.role = selectedRoles.value[0]
    } else {
      filters.role = selectedRoles.value.join(',')
    }
  }

  const applyFilters = () => {
    emit('filter', { ...filters })
  }

  const resetFilters = () => {
    Object.assign(filters, {
      keyword: '',
      status: '0',
      role: '0',
      sort_by: '',
      start_date: '',
      end_date: '',
    })
    dateRange.value = [null, null]
    selectedStatuses.value = []
    selectedRoles.value = []
    applyFilters()
  }

  onMounted(() => {
    if (props.initialFilters) {
      Object.assign(filters, props.initialFilters)

      if (filters.start_date && filters.end_date) {
        dateRange.value = [new Date(filters.start_date), new Date(filters.end_date)]
      }
    }
  })
</script>

<template>
  <div class="user-filter-panel mb-4 rounded-xl border border-subtle bg-background-800 p-3">
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-base font-bold text-content-heading">
        <i class="fas fa-filter mr-1.5 text-content" />{{ $t('admin.users.filter.title') }}
      </h3>
      <div class="flex space-x-1.5">
        <CyberButton type="secondary" class="px-2 py-1 text-xs" @click="resetFilters">
          <i class="fas fa-undo-alt mr-1" />{{ $t('admin.users.buttons.resetFilter') }}
        </CyberButton>
        <CyberButton type="primary" class="px-2 py-1 text-xs" @click="applyFilters">
          <i class="fas fa-search mr-1" />{{ $t('admin.users.buttons.applyFilter') }}
        </CyberButton>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-x-3 gap-y-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.users.filter.keyword') }}</label>
        <CyberInput v-model="filters.keyword" :placeholder="$t('admin.users.filter.keywordPlaceholder')" class="compact" />
      </div>

      <div class="filter-item">
        <label class="mb-1 block text-xs text-content">{{ $t('admin.users.filter.status.label') }}</label>
        <div class="checkbox-group">
          <CyberCheckbox
            v-for="option in statusOptions.slice(1)"
            :key="option.value"
            :model-value="selectedStatuses.includes(option.value)"
            class="checkbox-item"
            @update:model-value="(checked) => toggleStatusFilter(option.value, checked)"
          >
            {{ option.label }}
          </CyberCheckbox>
        </div>
      </div>

      <div class="filter-item">
        <label class="mb-1 block text-xs text-content">{{ $t('admin.users.filter.role.label') }}</label>
        <div class="checkbox-group">
          <CyberCheckbox
            v-for="option in roleOptions.slice(1)"
            :key="option.value"
            :model-value="selectedRoles.includes(option.value)"
            class="checkbox-item"
            @update:model-value="(checked) => toggleRoleFilter(option.value, checked)"
          >
            {{ option.label }}
          </CyberCheckbox>
        </div>
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.users.filter.sort.label') }}</label>
        <CyberDropdown
          v-model="filters.sort_by"
          :options="sortOptions"
          :placeholder="$t('admin.users.filter.sort.placeholder')"
          class="compact"
        />
      </div>

      <div class="filter-item">
        <label class="mb-0.5 block text-xs text-content">{{ $t('admin.users.filter.registrationTime') }}</label>
        <CyberDatePicker
          v-model="dateRange"
          type="daterange"
          :placeholder="$t('admin.users.filter.dateRangePlaceholder')"
          :start-placeholder="$t('admin.users.filter.startDatePlaceholder')"
          :end-placeholder="$t('admin.users.filter.endDatePlaceholder')"
          class="compact"
          @change="handleDateRangeChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  :deep(.cyber-dropdown.compact .cyber-dropdown-header) {
    min-height: 32px;
    height: 32px;
    padding: 0 0.6rem;
    font-size: 0.8rem;
    display: flex;
    align-items: center;
    box-sizing: border-box;
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-selected) {
    font-size: 0.8rem;
  }

  :deep(.cyber-input.compact) {
    height: 32px;
    padding: 0 0.75rem;
    box-sizing: border-box;
    font-size: 0.8rem;
    line-height: 1;
  }

  :deep(.cyber-datepicker.compact) {
    height: 32px;
    font-size: 0.8rem;
  }

  :deep(.cyber-datepicker.compact input) {
    height: 32px;
    padding: 0 0.75rem;
    font-size: 0.8rem;
    box-sizing: border-box;
  }

  .filter-item {
    margin-bottom: 0;
  }

  :deep(.cyber-dropdown-menu) {
    max-height: 200px;
    overflow-y: auto;
    padding: 4px;
    font-size: 0.8rem;
  }

  :deep(.cyber-dropdown.compact .cyber-dropdown-option) {
    padding: 3px 6px;
    font-size: 0.75rem;
    min-height: 24px;
  }

  .checkbox-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    padding: 0.5rem;
    background: rgba(var(--color-background-800-rgb), 0.4);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    min-height: 80px;
  }

  .checkbox-item {
    font-size: 0.75rem;

    :deep(.cyber-checkbox__inner) {
      width: 14px;
      height: 14px;
    }

    :deep(.cyber-checkbox__label) {
      font-size: 0.75rem;
      line-height: 1.2;
    }
  }
</style>
