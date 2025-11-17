<script setup lang="ts">
  import { ref, watch, computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    modelValue: string
    showFilter: boolean
    searchMode?: 'normal' | 'vector'
  }>()

  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
    (e: 'update:showFilter', value: boolean): void
    (e: 'update:searchMode', value: 'normal' | 'vector'): void
    (e: 'search'): void
  }>()

  /* 搜索模式选项 */
  const searchModeOptions = computed(() => [
    {
      value: 'normal',
      label: $t('explore.search.mode.normal'),
    },
    {
      value: 'vector',
      label: $t('explore.search.mode.ai'),
    },
  ])

  const keyword = ref(props.modelValue)
  const searchMode = ref<'normal' | 'vector'>(props.searchMode || 'normal')

  /* 监听外部变化 */
  watch(
    () => props.modelValue,
    (newVal) => {
      keyword.value = newVal
    }
  )

  watch(
    () => props.searchMode,
    (newVal) => {
      if (newVal) {
        searchMode.value = newVal
      }
    }
  )

  watch(keyword, (newVal) => {
    emit('update:modelValue', newVal)
  })

  watch(searchMode, (newVal) => {
    emit('update:searchMode', newVal)
  })

  const search = () => {
    emit('search')
  }

  const toggleFilter = () => {
    emit('update:showFilter', !props.showFilter)
  }
</script>

<template>
  <div class="flex items-center gap-2">
    <CyberInput
      v-model="keyword"
      :placeholder="searchMode === 'vector' ? $t('explore.search.aiSearchPlaceholder') : $t('explore.search.searchPlaceholder')"
      class="flex-grow"
      @enter="search"
    >
      <template #prepend>
        <CyberDropdown
          v-model="searchMode"
          :options="searchModeOptions"
          :clearable="false"
          :searchable="false"
          class="search-mode-dropdown"
          style="min-width: 100px; width: 100px"
        />
      </template>

      <!-- 后置搜索按钮 -->
      <template #append>
        <button
          class="search-btn"
          :title="searchMode === 'vector' ? $t('explore.search.mode.ai') : $t('explore.search.mode.normal')"
          @click="search"
        >
          <i class="fas" :class="searchMode === 'vector' ? 'fa-brain' : 'fa-search'" />
        </button>
      </template>
    </CyberInput>

    <!-- 筛选按钮 -->
    <CyberButton
      type="secondary"
      class="flex-shrink-0 text-xs"
      :class="{ 'border-error-500 text-error-500': showFilter }"
      @click="toggleFilter"
    >
      <i class="fas fa-filter mr-1" />{{ $t('explore.actions.filter') }}
    </CyberButton>
  </div>
</template>

<style scoped>
  .search-mode-dropdown {
    min-width: 100px;
  }

  .search-mode-dropdown :deep(.cyber-dropdown-header) {
    border: none;
    background: transparent;
    border-radius: 0;
    padding: 0 8px;
    height: 100%;
    min-height: auto;
  }

  .search-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 12px;
    background: transparent;
    border: none;
    color: var(--color-content-default);
    cursor: pointer;
    transition: color 0.2s ease;
    font-size: 14px;
  }

  .search-btn:hover {
    color: var(--color-error-500);
    background: transparent;
  }

  .search-btn i {
    font-size: 14px;
  }
</style>
