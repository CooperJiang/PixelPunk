<script setup lang="ts">
  import { ref, watch, computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { AnnouncementListQuery } from '@/api/types/announcement'
  import type { DropdownOption } from '@/components/Dropdown/types'

  interface Props {
    initialFilters: AnnouncementListQuery
  }

  interface Emits {
    (e: 'filter', filters: Partial<AnnouncementListQuery>): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const { $t } = useTexts()

  /* 状态选项 */
  const statusOptions = computed<DropdownOption[]>(() => [
    { label: $t('admin.announcements.filter.status.all'), value: 'all' },
    { label: $t('admin.announcements.status.draft'), value: 'draft' },
    { label: $t('admin.announcements.status.published'), value: 'published' },
    { label: $t('admin.announcements.status.archived'), value: 'archived' },
  ])

  /* 置顶选项 */
  const pinnedOptions = computed<DropdownOption[]>(() => [
    { label: $t('admin.announcements.filter.pinned.all'), value: 'all' },
    { label: $t('admin.announcements.filter.pinned.onlyPinned'), value: 'true' },
    { label: $t('admin.announcements.filter.pinned.notPinned'), value: 'false' },
  ])

  /* 内部过滤器状态 - 使用字符串值 */
  const internalStatus = ref<string>(props.initialFilters.status || 'all')
  const internalIsPinned = ref<string>(
    props.initialFilters.is_pinned === undefined ? 'all' : String(props.initialFilters.is_pinned)
  )
  const keyword = ref(props.initialFilters.keyword || '')

  /* 转换为实际过滤器值 */
  const localFilters = computed(() => {
    return {
      status: internalStatus.value === 'all' ? undefined : internalStatus.value,
      is_pinned: internalIsPinned.value === 'all' ? undefined : internalIsPinned.value === 'true',
      keyword: keyword.value,
    }
  })

  /* 应用筛选 */
  function applyFilters() {
    emit('filter', { ...localFilters.value })
  }

  function resetFilters() {
    internalStatus.value = 'all'
    internalIsPinned.value = 'all'
    keyword.value = ''
    applyFilters()
  }

  watch(
    () => props.initialFilters,
    (newFilters) => {
      internalStatus.value = newFilters.status || 'all'
      internalIsPinned.value = newFilters.is_pinned === undefined ? 'all' : String(newFilters.is_pinned)
      keyword.value = newFilters.keyword || ''
    },
    { deep: true }
  )
</script>

<template>
  <div class="filter-panel">
    <div class="filter-row">
      <div class="filter-item">
        <label class="filter-label">{{ $t('admin.announcements.filter.status.label') }}</label>
        <CyberDropdown
          v-model="internalStatus"
          :options="statusOptions"
          :placeholder="$t('admin.announcements.filter.status.all')"
          width="150px"
          :clearable="false"
        />
      </div>

      <div class="filter-item">
        <label class="filter-label">{{ $t('admin.announcements.filter.pinned.label') }}</label>
        <CyberDropdown
          v-model="internalIsPinned"
          :options="pinnedOptions"
          :placeholder="$t('admin.announcements.filter.pinned.all')"
          width="150px"
          :clearable="false"
        />
      </div>

      <div class="filter-item flex-1">
        <label class="filter-label">{{ $t('admin.announcements.filter.search.label') }}</label>
        <CyberInput
          v-model="keyword"
          :placeholder="$t('admin.announcements.filter.search.placeholder')"
          @keypress.enter="applyFilters"
        />
      </div>

      <div class="filter-actions">
        <CyberButton type="primary" size="small" @click="applyFilters">
          <i class="fas fa-search mr-1.5" />
          {{ $t('admin.announcements.filter.actions.query') }}
        </CyberButton>
        <CyberButton type="secondary" size="small" @click="resetFilters">
          <i class="fas fa-redo mr-1.5" />
          {{ $t('admin.announcements.filter.actions.reset') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .filter-panel {
    background: var(--color-background-container);
    border-radius: var(--radius-lg);
    padding: var(--space-md);
    margin-bottom: 0;
  }

  .filter-row {
    display: flex;
    gap: var(--space-md);
    align-items: flex-end;
    flex-wrap: wrap;
  }

  .filter-item {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
    min-width: 150px;
  }

  .filter-item.flex-1 {
    flex: 1;
    min-width: 200px;
  }

  .filter-label {
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    color: var(--color-content-muted);
  }

  .filter-actions {
    display: flex;
    gap: var(--space-sm);
    align-items: center;
  }
</style>
