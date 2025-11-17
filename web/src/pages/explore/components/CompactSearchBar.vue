<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import { getGlobalSettings } from '@/api/admin/settings/common'

  interface SearchHistory {
    keyword: string
    mode: 'normal' | 'vector'
    timestamp: number
  }

  const props = defineProps<{
    modelValue: string
    searchMode?: 'normal' | 'vector'
    showFilter: boolean
    filterCount?: number
    isSearching?: boolean
  }>()

  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
    (e: 'update:searchMode', value: 'normal' | 'vector'): void
    (e: 'update:showFilter', value: boolean): void
    (e: 'search', payload: { keyword: string; mode: 'normal' | 'vector' }): void
  }>()

  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const searchModes = computed(() => [
    {
      value: 'vector' as const,
      label: $t('explore.search.mode.ai'),
      icon: 'fas fa-brain',
      description: $t('explore.search.mode.aiDesc'),
    },
    {
      value: 'normal' as const,
      label: $t('explore.search.mode.normal'),
      icon: 'fas fa-search',
      description: $t('explore.search.mode.normalDesc'),
    },
  ])

  const keyword = ref(props.modelValue)
  const vectorEnabled = ref(false)
  const currentSearchMode = ref<'normal' | 'vector'>(props.searchMode || 'vector')
  const isFocused = ref(false)
  const showSuggestions = ref(false)
  const searchInput = ref<HTMLInputElement>()
  const searchHistory = ref<SearchHistory[]>([])
  const suggestions = ref<string[]>([])

  let suggestionTimer: NodeJS.Timeout | null = null
  const currentModePlaceholder = computed(() => {
    if (!vectorEnabled.value) {
      return $t('explore.search.placeholder')
    }
    return currentSearchMode.value === 'vector' ? $t('explore.search.aiPlaceholder') : $t('explore.search.placeholder')
  })
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
        currentSearchMode.value = newVal
      }
    }
  )

  watch(keyword, (newVal) => {
    emit('update:modelValue', newVal)
  })

  watch(currentSearchMode, (newVal) => {
    emit('update:searchMode', newVal)
    suggestions.value = []
  })
  const switchMode = (mode: 'normal' | 'vector') => {
    if (currentSearchMode.value !== mode) {
      currentSearchMode.value = mode
      nextTick(() => {
        searchInput.value?.focus()
      })
    }
  }

  const handleFocus = () => {
    isFocused.value = true
    showSuggestions.value = true

    if (!keyword.value.trim()) {
      loadSearchHistory()
    }
  }

  const handleBlur = () => {
    isFocused.value = false
    setTimeout(() => {
      showSuggestions.value = false
    }, 150)
  }

  const handleInput = () => {
    if (suggestionTimer) {
      clearTimeout(suggestionTimer)
    }

    suggestionTimer = setTimeout(() => {
      getSuggestions()
    }, 300)
  }

  const handleSearch = () => {
    const searchKeyword = keyword.value.trim()

    if (searchKeyword) {
      saveSearchHistory(searchKeyword, currentSearchMode.value)
    }

    emit('search', {
      keyword: searchKeyword,
      mode: currentSearchMode.value,
    })

    showSuggestions.value = false
    searchInput.value?.blur()
  }

  const clearSearch = () => {
    keyword.value = ''
    suggestions.value = []
    searchInput.value?.focus()
  }

  const selectSuggestion = (suggestion: string) => {
    keyword.value = suggestion
    showSuggestions.value = false
    nextTick(() => {
      handleSearch()
    })
  }

  const toggleFilter = () => {
    emit('update:showFilter', !props.showFilter)
  }
  const getSuggestions = () => {
    const query = keyword.value.trim().toLowerCase()
    if (!query || query.length < 2) {
      suggestions.value = []
      return
    }

    if (currentSearchMode.value === 'vector') {
      const aiSuggestions = $t('explore.search.suggestions.ai').filter((item: string) => item.toLowerCase().includes(query))

      suggestions.value = aiSuggestions
    } else {
      const normalSuggestions = ['screenshot', 'photo', 'image', 'picture', 'wallpaper'].filter((item) =>
        item.toLowerCase().includes(query)
      )

      suggestions.value = normalSuggestions
    }
  }
  const loadSearchHistory = () => {
    try {
      const saved = localStorage.getItem('gallery_search_history')
      if (saved) {
        const parsed = JSON.parse(saved) as SearchHistory[]
        const thirtyDaysAgo = Date.now() - 30 * 24 * 60 * 60 * 1000
        searchHistory.value = parsed.filter((item) => item.timestamp > thirtyDaysAgo).sort((a, b) => b.timestamp - a.timestamp)
      }
    } catch (_error) {
      searchHistory.value = []
    }
  }

  const saveSearchHistory = (searchKeyword: string, mode: 'normal' | 'vector') => {
    try {
      const filtered = searchHistory.value.filter((item) => item.keyword !== searchKeyword || item.mode !== mode)

      const newHistory: SearchHistory = {
        keyword: searchKeyword,
        mode,
        timestamp: Date.now(),
      }

      const updated = [newHistory, ...filtered].slice(0, 10)
      searchHistory.value = updated

      localStorage.setItem('gallery_search_history', JSON.stringify(updated))
    } catch (_error) {}
  }

  const clearHistory = () => {
    searchHistory.value = []
    localStorage.removeItem('gallery_search_history')
  }
  const fetchGlobalSettings = async () => {
    try {
      const result = await getGlobalSettings()
      if (result.success && result.data) {
        const vectorConfig = result.data.vector as Record<string, unknown> | undefined
        if (vectorConfig && typeof vectorConfig.vector_enabled === 'boolean') {
          vectorEnabled.value = vectorConfig.vector_enabled
          if (vectorEnabled.value) {
            currentSearchMode.value = props.searchMode || 'vector'
          } else {
            currentSearchMode.value = 'normal'
          }
        }
      }
    } catch (_error) {
      vectorEnabled.value = false
      currentSearchMode.value = 'normal'
    }
  }
  onMounted(async () => {
    await fetchGlobalSettings()
    loadSearchHistory()
  })

  onUnmounted(() => {
    if (suggestionTimer) {
      clearTimeout(suggestionTimer)
    }
  })
  defineExpose({
    focus: () => searchInput.value?.focus(),
    clear: clearSearch,
  })
</script>

<template>
  <div class="compact-search-bar">
    <div class="search-container">
      <div class="search-input-wrapper" :class="{ focused: isFocused, loading: isSearching }">
        <div v-if="vectorEnabled" class="mode-selector">
          <button
            v-for="mode in searchModes"
            :key="mode.value"
            class="mode-btn"
            :class="{ active: currentSearchMode === mode.value }"
            :title="mode.description"
            @click="switchMode(mode.value)"
          >
            <i :class="mode.icon" />
          </button>
        </div>

        <div v-if="vectorEnabled" class="separator" />

        <input
          ref="searchInput"
          v-model="keyword"
          type="text"
          class="search-input"
          :placeholder="currentModePlaceholder"
          @focus="handleFocus"
          @blur="handleBlur"
          @keyup.enter="handleSearch"
          @input="handleInput"
        />

        <button v-if="keyword" class="clear-btn" :title="$t('explore.search.clear')" @click="clearSearch">
          <i class="fas fa-times" />
        </button>

        <button
          class="search-btn"
          :disabled="isSearching"
          :title="`${currentSearchMode === 'vector' ? $t('explore.search.mode.ai') : $t('explore.search.mode.normalText')}`"
          @click="handleSearch"
          v-loading="isSearching"
        >
          <i class="fas fa-search" />
        </button>
      </div>

      <transition name="suggestions-fade">
        <div
          v-if="showSuggestions && (suggestions.length > 0 || searchHistory.length > 0)"
          class="search-suggestions"
          @click.stop
        >
          <div v-if="suggestions.length > 0" class="suggestions-group">
            <div class="group-header">{{ $t('explore.search.suggestions.header') }}</div>
            <div
              v-for="(suggestion, index) in suggestions.slice(0, 5)"
              :key="`suggestion-${index}`"
              class="suggestion-item"
              @click="selectSuggestion(suggestion)"
            >
              <i class="fas fa-search" />
              <span>{{ suggestion }}</span>
            </div>
          </div>

          <div v-if="searchHistory.length > 0 && !keyword.trim()" class="suggestions-group">
            <div class="group-header">
              {{ $t('explore.search.history.title') }}
              <button class="clear-history" :title="$t('actions.clear')" @click="clearHistory">
                <i class="fas fa-trash-alt" />
              </button>
            </div>
            <div
              v-for="(item, index) in searchHistory.slice(0, 3)"
              :key="`history-${index}`"
              class="suggestion-item history-item"
              @click="selectSuggestion(item.keyword)"
            >
              <i :class="item.mode === 'vector' ? 'fas fa-brain' : 'fas fa-search'" />
              <span>{{ item.keyword }}</span>
            </div>
          </div>
        </div>
      </transition>
    </div>

    <button class="filter-btn" :class="{ active: showFilter }" :title="$t('explore.actions.filter')" @click="toggleFilter">
      <i class="fas fa-sliders-h" />
      <span>{{ $t('explore.actions.filter') }}</span>
      <span v-if="filterCount > 0" class="filter-badge">{{ filterCount }}</span>
    </button>
  </div>
</template>

<style scoped>
  .compact-search-bar {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .search-container {
    position: relative;
    flex: 1;
    min-width: 320px;
    max-width: 400px;
    z-index: 100;
  }

  .search-input-wrapper {
    display: flex;
    align-items: center;
    min-height: 32px;
    height: 32px;
    background: var(--color-background-700);
    border: 1.5px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    transition:
      border-color 0.25s ease,
      box-shadow 0.25s ease,
      transform 0.25s ease;
    overflow: hidden;
    box-sizing: border-box;
    box-shadow:
      2px 2.5px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 0 rgba(var(--color-brand-500-rgb), 0);
  }

  .search-input-wrapper:hover {
    border-color: var(--color-border-strong);
    box-shadow:
      2.5px 3px 0 rgba(var(--color-brand-500-rgb), 0.3),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.08);
    transform: translate(-0.3px, -0.3px);
  }

  .search-input-wrapper.focused {
    border-color: var(--color-brand-500);
    box-shadow:
      3.5px 4px 0 rgba(var(--color-brand-500-rgb), 0.4),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.2);
    transform: translate(-0.5px, -0.5px);
  }
  .search-input-wrapper:active {
    transform: translate(0.5px, 0.5px);
    box-shadow:
      1.5px 2px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .search-input-wrapper.loading {
    border-color: var(--color-brand-500);
    box-shadow:
      3.5px 4px 0 rgba(var(--color-brand-500-rgb), 0.4),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
  }
  .mode-selector {
    display: flex;
    flex-shrink: 0;
    gap: 0;
    padding: 4px 4px 4px 0;
  }

  .mode-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 28px;
    background: transparent;
    border: none;
    border-radius: var(--radius-sm);
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 13px;
  }

  .mode-btn:hover {
    color: var(--color-content-default);
    background: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .mode-btn.active {
    color: var(--color-brand-400);
    background: rgba(var(--color-brand-500-rgb), 0.25);
    font-weight: 500;
  }
  .separator {
    width: 1px;
    height: 24px;
    background: rgba(var(--color-border-default-rgb), 0.3);
    flex-shrink: 0;
    margin: 0 4px;
  }
  .search-input {
    flex: 1;
    height: 100%;
    padding: 0 8px;
    background: transparent;
    border: none;
    color: var(--color-white);
    font-size: 0.875rem;
    outline: none;
  }

  .search-input::placeholder {
    color: var(--color-content-disabled);
  }
  .clear-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: transparent;
    border: none;
    border-radius: var(--radius-sm);
    color: var(--color-content-disabled);
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 12px;
    margin-right: 2px;
  }

  .clear-btn:hover {
    color: var(--color-error-500);
    background: rgba(var(--color-error-rgb), 0.1);
  }
  .search-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    background: var(--color-brand-500);
    border: none;
    border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
    color: var(--color-text-on-brand);
    cursor: pointer;
    transition: all 0.2s ease;
    margin: 0;
    flex-shrink: 0;
    font-size: 13px;
    position: relative;
    overflow: hidden;
  }
  .search-btn :deep(.v-loading-mask) {
    background: transparent;
    backdrop-filter: none;
  }

  .search-btn :deep(.loading-container) {
    padding: 0;
  }

  .search-btn :deep(.simple-loading-core) {
    width: 20px !important;
    height: 20px !important;
  }

  .search-btn :deep(.spinner-ring),
  .search-btn :deep(.spinner-ring-middle) {
    border-width: 2px !important;
  }

  .search-btn :deep(.loading-text) {
    display: none;
  }

  .search-btn:hover:not(:disabled) {
    background: rgba(var(--color-brand-500-rgb), 0.85);
  }

  .search-btn:active:not(:disabled) {
    transform: scale(0.95);
  }

  .search-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
  .filter-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 0 12px;
    min-height: 32px;
    height: 32px;
    background: var(--color-background-700);
    border: 1.5px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    color: var(--color-content-default);
    cursor: pointer;
    transition:
      border-color 0.25s ease,
      box-shadow 0.25s ease,
      transform 0.25s ease;
    font-size: 0.875rem;
    white-space: nowrap;
    position: relative;
    box-sizing: border-box;
    box-shadow:
      2px 2.5px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 0 rgba(var(--color-brand-500-rgb), 0);
  }

  .filter-btn:hover {
    border-color: var(--color-border-strong);
    box-shadow:
      2.5px 3px 0 rgba(var(--color-brand-500-rgb), 0.3),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.08);
    transform: translate(-0.3px, -0.3px);
  }

  .filter-btn.active {
    border-color: var(--color-brand-500);
    color: var(--color-brand-500);
    box-shadow:
      3.5px 4px 0 rgba(var(--color-brand-500-rgb), 0.4),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.25);
    transform: translate(-0.5px, -0.5px);
  }
  .filter-btn:active {
    transform: translate(0.5px, 0.5px);
    box-shadow:
      1.5px 2px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .filter-badge {
    position: absolute;
    top: -6px;
    right: -6px;
    background: var(--color-error-500);
    color: var(--color-text-on-brand);
    font-size: 0.65rem;
    font-weight: 600;
    width: 18px;
    height: 18px;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 18px;
  }
  .search-suggestions {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: 9999;
    background: var(--color-background-800);
    border: 1px solid rgba(var(--color-border-default-rgb), 0.3);
    border-radius: var(--radius-sm);
    margin-top: 4px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
    max-height: 240px;
    overflow-y: auto;
  }

  .suggestions-group {
    padding: 4px 0;
  }

  .suggestions-group + .suggestions-group {
    border-top: 1px solid rgba(var(--color-border-default-rgb), 0.2);
  }

  .group-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 12px;
    font-size: 0.75rem;
    color: var(--color-content-muted);
    font-weight: 500;
  }

  .clear-history {
    background: transparent;
    border: none;
    border-radius: var(--radius-sm);
    color: var(--color-content-disabled);
    cursor: pointer;
    padding: 4px 6px;
    font-size: 0.7rem;
    transition: all 0.2s ease;
  }

  .clear-history:hover {
    color: var(--color-error-500);
    background: rgba(var(--color-error-rgb), 0.1);
  }

  .suggestion-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 12px;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.875rem;
    color: var(--color-content-default);
  }

  .suggestion-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
  }

  .suggestion-item i {
    width: 14px;
    font-size: 11px;
    color: var(--color-content-muted);
    flex-shrink: 0;
  }

  .suggestion-item:hover i {
    color: var(--color-brand-500);
  }

  .history-item i {
    color: var(--color-brand-500);
  }
  .suggestions-fade-enter-active,
  .suggestions-fade-leave-active {
    transition: all 0.2s ease;
  }

  .suggestions-fade-enter-from,
  .suggestions-fade-leave-to {
    opacity: 0;
    transform: translateY(-4px);
  }
  .search-suggestions::-webkit-scrollbar {
    width: 6px;
  }

  .search-suggestions::-webkit-scrollbar-track {
    background: transparent;
  }

  .search-suggestions::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
  }

  .search-suggestions::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }
  @media (max-width: 768px) {
    .search-input-wrapper {
      min-height: 36px;
      height: 36px;
    }

    .search-btn {
      width: 36px;
      height: 36px;
    }

    .filter-btn {
      min-height: 36px;
      height: 36px;
      padding: 0 12px;
    }

    .filter-btn span {
      display: none;
    }

    .search-input {
      font-size: 16px;
    }

    .mode-btn {
      width: 32px;
      height: 32px;
    }
  }

  @media (max-width: 480px) {
    .search-container {
      min-width: 240px;
    }
  }
</style>
