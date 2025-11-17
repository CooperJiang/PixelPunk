<template>
  <div class="vector-filter-panel">
    <div class="filter-form">
      <div class="filter-row">
        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.vectors.columns.status') }}</label>
          <CyberDropdown
            v-model="localFilters.status"
            :placeholder="$t('admin.vectors.filter.options.all')"
            :options="statusOptions"
            class="filter-input"
          />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.vectors.qdrant.version') }}</label>
          <CyberDropdown
            v-model="localFilters.model"
            :placeholder="$t('admin.vectors.filter.options.all')"
            :options="modelOptions"
            class="filter-input"
          />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.vectors.filter.keyword') }}</label>
          <CyberInput
            v-model="localFilters.keyword"
            :placeholder="$t('admin.vectors.filter.options.keywordPlaceholder')"
            clearable
            class="filter-input"
          >
            <template #prefix>
              <i class="fas fa-search"></i>
            </template>
          </CyberInput>
        </div>
      </div>

      <div class="filter-row">
        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.vectors.pagination.perPage') }}</label>
          <CyberDropdown v-model="localFilters.page_size" :options="pageSizeOptions" class="filter-input" />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ $t('admin.vectors.filter.sortBy') }}</label>
          <CyberDropdown v-model="sortField" :options="sortOptions" class="filter-input" />
        </div>

        <div class="filter-item filter-actions">
          <CyberButton type="secondary" @click="handleReset" class="action-btn">
            <i class="fas fa-undo mr-1"></i>
            {{ $t('admin.vectors.filter.reset') }}
          </CyberButton>

          <CyberButton type="primary" @click="handleApply" class="action-btn">
            <i class="fas fa-filter mr-1"></i>
            {{ $t('admin.vectors.filter.apply') }}
          </CyberButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, watch, computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { VectorListParams } from '@/api/admin/vectors'

  interface Props {
    initialFilters: VectorListParams
    availableModels?: string[]
  }

  interface Emits {
    (e: 'filter', filters: VectorListParams): void
  }

  const props = withDefaults(defineProps<Props>(), {
    availableModels: () => [],
  })

  const emit = defineEmits<Emits>()
  const { $t } = useTexts()

  const localFilters = reactive<VectorListParams>({
    page: 1,
    page_size: 10,
    status: undefined,
    model: undefined,
    keyword: undefined,
  })

  const sortField = ref('created_at_desc')

  const statusOptions = computed(() => [
    { label: $t('admin.vectors.status.pending'), value: 'pending' },
    { label: $t('admin.vectors.status.running'), value: 'running' },
    { label: $t('admin.vectors.status.completed'), value: 'completed' },
    { label: $t('admin.vectors.status.failed'), value: 'failed' },
    { label: $t('admin.vectors.status.reset'), value: 'reset' },
  ])

  const pageSizeOptions = computed(() => [
    { label: $t('admin.vectors.filter.pageSize.size10'), value: 10 },
    { label: $t('admin.vectors.filter.pageSize.size20'), value: 20 },
    { label: $t('admin.vectors.filter.pageSize.size50'), value: 50 },
    { label: $t('admin.vectors.filter.pageSize.size100'), value: 100 },
  ])

  const sortOptions = computed(() => [
    { label: $t('admin.vectors.filter.sort.createdAtDesc'), value: 'created_at_desc' },
    { label: $t('admin.vectors.filter.sort.createdAtAsc'), value: 'created_at_asc' },
    { label: $t('admin.vectors.filter.sort.updatedAtDesc'), value: 'updated_at_desc' },
    { label: $t('admin.vectors.filter.sort.updatedAtAsc'), value: 'updated_at_asc' },
    { label: $t('admin.vectors.filter.sort.retryCountDesc'), value: 'retry_count_desc' },
    { label: $t('admin.vectors.filter.sort.durationDesc'), value: 'duration_desc' },
  ])

  const modelOptions = computed(() => {
    return props.availableModels.map((model) => ({
      label: model,
      value: model,
    }))
  })

  const initializeFilters = () => {
    Object.assign(localFilters, {
      page: props.initialFilters.page || 1,
      page_size: props.initialFilters.page_size || 10,
      status: props.initialFilters.status,
      model: props.initialFilters.model,
      keyword: props.initialFilters.keyword,
    })
  }

  watch(() => props.initialFilters, initializeFilters, { immediate: true })

  const handleApply = () => {
    const filters: VectorListParams = {
      ...localFilters,
      page: 1, // 重置到第一页
    }

    emit('filter', filters)
  }

  const handleReset = () => {
    localFilters.page = 1
    localFilters.page_size = 10
    localFilters.status = undefined
    localFilters.model = undefined
    localFilters.keyword = undefined
    sortField.value = 'created_at_desc'

    handleApply()
  }
</script>

<style scoped lang="scss">
  .vector-filter-panel {
    margin: 15px;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.9) 0%,
      rgba(var(--color-background-700-rgb), 0.7) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 20px;
    margin-bottom: 16px;
    backdrop-filter: blur(10px);
    box-shadow:
      0 4px 20px rgba(0, 0, 0, 0.2),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.08);
  }

  .filter-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .filter-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
    align-items: end;
  }

  .filter-item {
    display: flex;
    flex-direction: column;
    gap: 6px;

    &.filter-actions {
      flex-direction: row;
      align-items: center;
      gap: 8px;
      margin-top: 20px;

      .action-btn {
        flex: 1;
        min-width: 100px;
      }
    }
  }

  .filter-label {
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--color-slate-400);
    margin-bottom: 4px;
  }

  .filter-input {
    width: 100%;
  }

  @media (max-width: 768px) {
    .filter-row {
      grid-template-columns: 1fr;
    }

    .filter-item.filter-actions {
      flex-direction: column;

      .action-btn {
        width: 100%;
      }
    }
  }
</style>
