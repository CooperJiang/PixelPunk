<script setup lang="ts">
  import { ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    keyword: string
    status: string
    role: string
    startDate?: string
    endDate?: string
    sortBy?: string
  }>()

  const emit = defineEmits<{
    (
      e: 'search',
      params: {
        keyword: string
        status: string
        role: string
        start_date?: string
        end_date?: string
        sort_by?: string
        sort_order?: string
      }
    ): void
    (e: 'reset'): void
  }>()

  /* 基础搜索状态 */
  const searchText = ref(props.keyword || '')
  const statusValue = ref(props.status || '0')
  const roleValue = ref(props.role || '0')

  /* 高级搜索状态 */
  const showAdvanced = ref(false)
  const startDate = ref(props.startDate || '')
  const endDate = ref(props.endDate || '')
  const sortBy = ref(props.sortBy || 'created_at_desc')
  const activeQuickFilter = ref('')

  /* 状态选项 */
  const statusOptions = [
    { label: t('admin.users.search.status.all'), value: '0' },
    { label: t('admin.users.search.status.normal'), value: '1' },
    { label: t('admin.users.search.status.disabled'), value: '2' },
    { label: t('admin.users.search.status.deleted'), value: '3' },
  ]

  /* 角色选项 */
  const roleOptions = [
    { label: t('admin.users.search.role.all'), value: '0' },
    { label: t('admin.users.search.role.super_admin'), value: '1' },
    // { label: t('admin.users.search.role.admin'), value: '2' },
    { label: t('admin.users.search.role.user'), value: '3' },
  ]

  /* 排序选项 */
  const sortOptions = [
    { label: t('admin.users.search.sort.created_newest'), value: 'created_at_desc' },
    { label: t('admin.users.search.sort.created_oldest'), value: 'created_at_asc' },
    { label: t('admin.users.search.sort.username_asc'), value: 'username_asc' },
    { label: t('admin.users.search.sort.username_desc'), value: 'username_desc' },
    { label: t('admin.users.search.sort.last_active'), value: 'last_activity_desc' },
  ]

  /* 快速筛选标签 */
  const quickFilterTags = [
    { key: 'today', label: t('admin.users.search.quick_filter.today'), icon: 'fas fa-calendar-day' },
    { key: 'week', label: t('admin.users.search.quick_filter.week'), icon: 'fas fa-calendar-week' },
    { key: 'month', label: t('admin.users.search.quick_filter.month'), icon: 'fas fa-calendar-alt' },
    { key: 'active', label: t('admin.users.search.quick_filter.active'), icon: 'fas fa-user-check' },
    { key: 'inactive', label: t('admin.users.search.quick_filter.inactive'), icon: 'fas fa-user-times' },
    { key: 'admins', label: t('admin.users.search.quick_filter.admins'), icon: 'fas fa-user-shield' },
  ]

  /* 监听props变化 */
  watch(
    () => props.keyword,
    (newVal) => {
      searchText.value = newVal
    },
    { immediate: true }
  )

  watch(
    () => props.status,
    (newVal) => {
      statusValue.value = newVal
    },
    { immediate: true }
  )

  watch(
    () => props.role,
    (newVal) => {
      roleValue.value = newVal
    },
    { immediate: true }
  )

  const toggleAdvanced = () => {
    showAdvanced.value = !showAdvanced.value
  }

  const handleSearch = () => {
    const [sortField, sortOrder] = sortBy.value.split('_')
    emit('search', {
      keyword: searchText.value,
      status: statusValue.value,
      role: roleValue.value,
      start_date: startDate.value || undefined,
      end_date: endDate.value || undefined,
      sort_by: sortField,
      sort_order: sortOrder,
    })
  }

  const handleFilterChange = () => {
    handleSearch()
  }

  const handleQuickFilter = (key: string) => {
    if (activeQuickFilter.value === key) {
      activeQuickFilter.value = ''
      handleReset()
      return
    }

    activeQuickFilter.value = key
    const today = new Date()

    switch (key) {
      case 'today':
        startDate.value = today.toISOString().split('T')[0]
        endDate.value = today.toISOString().split('T')[0]
        break
      case 'week':
        const weekStart = new Date(today.setDate(today.getDate() - today.getDay()))
        startDate.value = weekStart.toISOString().split('T')[0]
        endDate.value = new Date().toISOString().split('T')[0]
        break
      case 'month':
        const monthStart = new Date(today.getFullYear(), today.getMonth(), 1)
        startDate.value = monthStart.toISOString().split('T')[0]
        endDate.value = new Date().toISOString().split('T')[0]
        break
      case 'active':
        sortBy.value = 'last_activity_desc'
        break
      case 'inactive':
        sortBy.value = 'last_activity_asc'
        break
      case 'admins':
        roleValue.value = '1' // 超级管理员
        break
    }

    handleSearch()
  }

  const handleReset = () => {
    searchText.value = ''
    statusValue.value = '0'
    roleValue.value = '0'
    startDate.value = ''
    endDate.value = ''
    sortBy.value = 'created_at_desc'
    activeQuickFilter.value = ''
    showAdvanced.value = false
    emit('reset')
  }
</script>

<template>
  <div class="search-bar space-y-4">
    <div class="flex items-center gap-4">
      <div class="relative flex-1">
        <cyberInput
          v-model="searchText"
          :placeholder="$t('admin.users.search.placeholder')"
          suffix-icon="search"
          @enter="handleFilterChange"
        />
      </div>

      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2">
          <span class="whitespace-nowrap text-sm text-content">{{ $t('admin.users.search.status_label') }}：</span>
          <cyberDropdown
            v-model="statusValue"
            :options="statusOptions"
            :placeholder="$t('admin.users.search.status_placeholder')"
            class="w-28"
            @change="handleFilterChange"
          />
        </div>

        <div class="flex items-center gap-2">
          <span class="whitespace-nowrap text-sm text-content">{{ $t('admin.users.search.role_label') }}：</span>
          <cyberDropdown
            v-model="roleValue"
            :options="roleOptions"
            :placeholder="$t('admin.users.search.role_placeholder')"
            class="w-28"
            @change="handleFilterChange"
          />
        </div>

        <button
          class="bg-secondary flex h-10 w-10 items-center justify-center rounded-lg border border-brand-300 text-content-heading transition-colors hover:bg-background-600"
          :class="{ 'bg-brand-200': showAdvanced }"
          @click="toggleAdvanced"
        >
          <i class="fas fa-filter" />
        </button>

        <button
          class="bg-secondary flex h-10 w-10 items-center justify-center rounded-lg border border-brand-300 text-content-heading transition-colors hover:bg-background-600"
          @click="handleSearch"
        >
          <i class="fas fa-search" />
        </button>

        <button
          class="bg-secondary flex h-10 w-10 items-center justify-center rounded-lg border border-brand-300 text-content-heading transition-colors hover:bg-background-600"
          @click="handleReset"
        >
          <i class="fas fa-redo-alt" />
        </button>
      </div>
    </div>

    <Transition name="advanced-search">
      <div v-if="showAdvanced" class="space-y-4 rounded-lg border border-subtle bg-background-300 p-4">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          <div>
            <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.search.start_date_label') }}：</label>
            <cyberDatePicker
              v-model="startDate"
              :placeholder="$t('admin.users.search.start_date_placeholder')"
              @change="handleFilterChange"
            />
          </div>

          <div>
            <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.search.end_date_label') }}：</label>
            <cyberDatePicker
              v-model="endDate"
              :placeholder="$t('admin.users.search.end_date_placeholder')"
              @change="handleFilterChange"
            />
          </div>

          <div>
            <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.search.sort_label') }}：</label>
            <cyberDropdown
              v-model="sortBy"
              :options="sortOptions"
              :placeholder="$t('admin.users.search.sort_placeholder')"
              @change="handleFilterChange"
            />
          </div>
        </div>

        <div>
          <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.search.quick_filter_label') }}：</label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="tag in quickFilterTags"
              :key="tag.key"
              class="flex items-center gap-1 rounded-lg border-2 px-3 py-2 text-xs font-medium transition-all duration-200"
              :class="[
                activeQuickFilter === tag.key
                  ? 'shadow-cyber-blue-250 scale-105 border-brand-500 bg-brand-200 text-content shadow-lg'
                  : 'border-brand-300 bg-background-500 text-content hover:scale-105 hover:border-brand-500 hover:bg-background-600',
              ]"
              @click="handleQuickFilter(tag.key)"
            >
              <i class="text-xs" :class="[tag.icon, activeQuickFilter === tag.key ? 'text-content' : '']" />
              {{ tag.label }}
              <div
                v-if="activeQuickFilter === tag.key"
                class="border-cyber-gray absolute -right-1 -top-1 flex h-3 w-3 items-center justify-center rounded-full border-2 bg-brand-500"
              >
                <i class="fas fa-check text-[8px] text-content-heading" />
              </div>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped lang="scss">
  .advanced-search-enter-active,
  .advanced-search-leave-active {
    transition: all 0.3s ease;
  }

  .advanced-search-enter-from,
  .advanced-search-leave-to {
    opacity: 0;
    transform: translateY(-10px);
  }

  .scale-105 {
    transform: scale(1.05);
  }

  .scale-102 {
    transform: scale(1.02);
  }

  .hover\:scale-102:hover {
    transform: scale(1.02);
  }
</style>
