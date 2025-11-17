<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

  interface SearchHistory {
    keyword: string
    mode: 'normal' | 'vector'
    timestamp: number
  }

  interface SearchStatus {
    type: 'success' | 'error' | 'info'
    message: string
    icon: string
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
    (e: 'suggestion', keyword: string): void
  }>()

  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  /* 搜索模式定义 */
  const searchModes = computed(() => [
    {
      value: 'normal' as const,
      label: $t('explore.search.mode.normal'),
      icon: 'fas fa-search',
      description: $t('explore.search.mode.normalDescDetail'),
    },
    {
      value: 'vector' as const,
      label: $t('explore.search.mode.ai'),
      icon: 'fas fa-brain',
      description: $t('explore.search.mode.aiDescDetail'),
    },
  ])

  /* 组件状态 */
  const keyword = ref(props.modelValue)
  const searchMode = ref<'normal' | 'vector'>(props.searchMode || 'normal')
  const isFocused = ref(false)
  const showSuggestions = ref(false)
  const searchInput = ref<HTMLInputElement>()
  const searchHistory = ref<SearchHistory[]>([])
  const suggestions = ref<string[]>([])
  const searchStatus = ref<SearchStatus | null>(null)

  /* 防抖定时器 */
  const debounceTimer: NodeJS.Timeout | null = null
  let suggestionTimer: NodeJS.Timeout | null = null

  const currentModeIcon = computed(() => {
    const mode = searchModes.value.find((m) => m.value === searchMode.value)
    return mode?.icon || 'fas fa-search'
  })

  const currentModePlaceholder = computed(() =>
    searchMode.value === 'vector' ? $t('explore.search.aiPlaceholderDetail') : $t('explore.search.normalPlaceholder')
  )

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
    suggestions.value = []
  })

  const switchMode = (mode: 'normal' | 'vector') => {
    if (searchMode.value !== mode) {
      searchMode.value = mode
      showStatus({
        type: 'info',
        message: $t('explore.search.switchedTo', {
          mode: mode === 'vector' ? $t('explore.search.mode.aiIntelligent') : $t('explore.search.mode.normalText'),
        }),
        icon: mode === 'vector' ? 'fas fa-brain' : 'fas fa-search',
      })

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
    }, 200)
  }

  const handleInput = () => {
    if (suggestionTimer) {
      clearTimeout(suggestionTimer)
    }

    suggestionTimer = setTimeout(() => {
      getSuggestions()
    }, 300)

    emit('suggestion', keyword.value)
  }

  const handleSearch = () => {
    const searchKeyword = keyword.value.trim()

    // 允许空搜索
    if (searchKeyword) {
      saveSearchHistory(searchKeyword, searchMode.value)
    }

    emit('search', {
      keyword: searchKeyword,
      mode: searchMode.value,
    })

    showSuggestions.value = false

    searchInput.value?.blur()

    showStatus({
      type: 'info',
      message: $t('explore.search.searching', {
        mode: searchMode.value === 'vector' ? $t('explore.search.mode.aiIntelligent') : $t('explore.search.mode.normalText'),
      }),
      icon: 'fas fa-search',
    })
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

    if (searchMode.value === 'vector') {
      const aiSuggestions = $t('explore.search.suggestions.ai').filter((item: string) => item.toLowerCase().includes(query))

      suggestions.value = aiSuggestions.slice(0, 5)
    } else {
      const normalSuggestions = ['screenshot', 'photo', 'image', 'picture', 'wallpaper'].filter((item) =>
        item.toLowerCase().includes(query)
      )

      suggestions.value = normalSuggestions.slice(0, 5)
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

      const updated = [newHistory, ...filtered].slice(0, 20)
      searchHistory.value = updated

      localStorage.setItem('gallery_search_history', JSON.stringify(updated))
    } catch (_error) {}
  }

  const clearHistory = () => {
    searchHistory.value = []
    localStorage.removeItem('gallery_search_history')
    showStatus({
      type: 'success',
      message: $t('explore.search.clearHistory'),
      icon: 'fas fa-check',
    })
  }

  const showStatus = (status: SearchStatus) => {
    searchStatus.value = status
    setTimeout(() => {
      searchStatus.value = null
    }, 3000)
  }

  watch(
    () => props.isSearching,
    (isSearching) => {
      if (
        !isSearching &&
        searchStatus.value?.message.includes($t('explore.search.searching', { mode: '' }).split('{')[0].trim())
      ) {
        searchStatus.value = null
      }
    }
  )

  const handleKeyDown = (e: KeyboardEvent) => {
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      e.preventDefault()
      searchInput.value?.focus()
    }

    if (e.key === 'Escape' && showSuggestions.value) {
      showSuggestions.value = false
      searchInput.value?.blur()
    }
  }

  onMounted(() => {
    loadSearchHistory()
    window.addEventListener('keydown', handleKeyDown)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
    if (suggestionTimer) {
      clearTimeout(suggestionTimer)
    }
  })

  defineExpose({
    focus: () => searchInput.value?.focus(),
    blur: () => searchInput.value?.blur(),
    clear: clearSearch,
  })
</script>

<template>
  <div class="advanced-search-bar">
    <div class="search-mode-tabs">
      <button
        v-for="mode in searchModes"
        :key="mode.value"
        class="mode-tab"
        :class="{ active: searchMode === mode.value }"
        :title="mode.description"
        @click="switchMode(mode.value)"
      >
        <i :class="mode.icon" />
        <span>{{ mode.label }}</span>
      </button>
    </div>

    <div class="search-input-container">
      <div class="search-input-wrapper" :class="{ focused: isFocused, loading: isSearching }">
        <div class="search-icon">
          <i v-if="!isSearching" :class="currentModeIcon" />
          <div v-else class="loading-spinner">
            <i class="fas fa-circle-notch fa-spin" />
          </div>
        </div>

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

        <button v-if="keyword" class="clear-button" :title="$t('explore.search.clear')" @click="clearSearch">
          <i class="fas fa-times" />
        </button>

        <button
          class="search-button"
          :disabled="!keyword.trim() || isSearching"
          :title="`${searchMode === 'vector' ? $t('explore.search.mode.ai') : $t('explore.search.mode.normal')}`"
          @click="handleSearch"
        >
          <i class="fas fa-search" />
        </button>
      </div>

      <transition name="suggestions">
        <div v-if="showSuggestions && (suggestions.length > 0 || searchHistory.length > 0)" class="search-suggestions">
          <div v-if="suggestions.length > 0" class="suggestions-section">
            <div class="suggestions-header">
              <i class="fas fa-lightbulb" />
              <span>{{ $t('explore.search.suggestions.title') }}</span>
            </div>
            <div
              v-for="(suggestion, index) in suggestions"
              :key="`suggestion-${index}`"
              class="suggestion-item"
              @click="selectSuggestion(suggestion)"
            >
              <i class="fas fa-search suggestion-icon" />
              <span class="suggestion-text">{{ suggestion }}</span>
            </div>
          </div>

          <div v-if="searchHistory.length > 0 && !keyword.trim()" class="suggestions-section">
            <div class="suggestions-header">
              <i class="fas fa-history" />
              <span>{{ $t('explore.search.history.title') }}</span>
              <button class="clear-history-btn" @click="clearHistory">
                <i class="fas fa-trash-alt" />
              </button>
            </div>
            <div
              v-for="(item, index) in searchHistory.slice(0, 5)"
              :key="`history-${index}`"
              class="suggestion-item history-item"
              @click="selectSuggestion(item.keyword)"
            >
              <i :class="item.mode === 'vector' ? 'fas fa-brain' : 'fas fa-search'" class="suggestion-icon" />
              <span class="suggestion-text">{{ item.keyword }}</span>
              <span class="suggestion-mode">{{
                item.mode === 'vector' ? $t('explore.search.mode.ai') : $t('explore.search.mode.normalText')
              }}</span>
            </div>
          </div>
        </div>
      </transition>
    </div>

    <transition name="status">
      <div v-if="searchStatus" class="search-status" :class="searchStatus.type">
        <i :class="searchStatus.icon" />
        <span>{{ searchStatus.message }}</span>
      </div>
    </transition>

    <button
      class="filter-button"
      :class="{ active: showFilter }"
      :title="$t('explore.actions.advancedFilter')"
      @click="toggleFilter"
    >
      <i class="fas fa-sliders-h" />
      <span>{{ $t('explore.search.filter') }}</span>
      <span v-if="filterCount > 0" class="filter-count">{{ filterCount }}</span>
    </button>
  </div>
</template>

<style scoped>
  .advanced-search-bar {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    position: relative;
  }

  .search-mode-tabs {
    display: flex;
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-md);
    padding: 2px;
    gap: 2px;
    flex-shrink: 0;
  }

  .mode-tab {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    background: transparent;
    border: none;
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all 0.3s ease;
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    white-space: nowrap;
  }

  .mode-tab:hover {
    color: var(--color-content-heading);
    background: var(--color-hover-bg);
  }

  .mode-tab.active {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-error-rgb), 0.1));
    color: var(--color-brand-500);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .mode-tab i {
    font-size: 14px;
  }

  .search-input-container {
    position: relative;
    flex: 1;
    min-width: 300px;
  }

  .search-input-wrapper {
    display: flex;
    align-items: center;
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-md);
    transition: all 0.3s ease;
    overflow: hidden;
    backdrop-filter: blur(8px);
  }

  .search-input-wrapper.focused {
    border-color: var(--color-focus-border);
    box-shadow: var(--shadow-glow-md);
    background: var(--color-background-700);
  }

  .search-input-wrapper.loading {
    border-color: var(--color-error-500);
    box-shadow: 0 0 15px rgba(var(--color-error-rgb), 0.3);
  }

  .search-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    color: rgba(var(--color-brand-500-rgb), 0.8);
    flex-shrink: 0;
  }

  .loading-spinner i {
    color: var(--color-error-500);
  }

  .search-input {
    flex: 1;
    height: 40px;
    padding: 0 12px;
    background: transparent;
    border: none;
    color: var(--color-white);
    font-size: 0.9rem;
    outline: none;
    font-family: inherit;
  }

  .search-input::placeholder {
    color: var(--color-content-muted);
  }

  .clear-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    background: transparent;
    border: none;
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all 0.2s ease;
    border-radius: var(--radius-sm);
    margin-right: 4px;
  }

  .clear-button:hover {
    color: var(--color-error-500);
    background: rgba(var(--color-error-rgb), 0.1);
  }

  .search-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-error-500));
    border: none;
    color: var(--color-text-on-brand);
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 14px;
  }

  .search-button:hover:not(:disabled) {
    background: linear-gradient(135deg, var(--color-cyan-600), var(--color-pink-600));
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .search-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    background: rgba(var(--color-content-rgb), 0.1);
  }

  .search-suggestions {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: 1000;
    background: var(--color-background-700);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    border-radius: var(--radius-sm);
    margin-top: 4px;
    box-shadow: 0 8px 32px var(--color-overlay-heavy);
    backdrop-filter: blur(12px);
    max-height: 300px;
    overflow-y: auto;
  }

  .suggestions-section {
    padding: 8px 0;
  }

  .suggestions-section + .suggestions-section {
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .suggestions-header {
    display: flex;
    align-items: center;
    padding: 8px 16px;
    font-size: 0.8rem;
    color: rgba(var(--color-brand-500-rgb), 0.8);
    font-weight: 500;
  }

  .suggestions-header i {
    margin-right: 8px;
  }

  .clear-history-btn {
    margin-left: auto;
    background: transparent;
    border: none;
    color: var(--color-content-muted);
    cursor: pointer;
    padding: 2px;
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
  }

  .clear-history-btn:hover {
    color: var(--color-error-500);
    background: rgba(var(--color-error-rgb), 0.1);
  }

  .suggestion-item {
    display: flex;
    align-items: center;
    padding: 8px 16px;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.875rem;
  }

  .suggestion-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
  }

  .suggestion-icon {
    width: 16px;
    margin-right: 12px;
    color: var(--color-content-muted);
    font-size: 12px;
  }

  .suggestion-text {
    flex: 1;
  }

  .suggestion-mode {
    font-size: 0.75rem;
    color: var(--color-content-muted);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    padding: 2px 6px;
    border-radius: var(--radius-sm);
    margin-left: 8px;
  }

  .history-item .suggestion-icon {
    color: rgba(var(--color-error-rgb), 0.6);
  }

  .search-status {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: 999;
    display: flex;
    align-items: center;
    padding: 8px 12px;
    margin-top: 4px;
    border-radius: var(--radius-sm);
    font-size: 0.85rem;
    backdrop-filter: blur(8px);
  }

  .search-status.success {
    background: rgba(var(--color-success-rgb), 0.15);
    border: 1px solid rgba(var(--color-success-rgb), 0.3);
    color: var(--color-success-500);
  }

  .search-status.error {
    background: rgba(var(--color-error-rgb), 0.15);
    border: 1px solid rgba(var(--color-error-rgb), 0.3);
    color: var(--color-error-500);
  }

  .search-status.info {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
  }

  .search-status i {
    margin-right: 8px;
  }

  .filter-button {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    background: rgba(var(--color-background-900-rgb), 0.8);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: var(--color-content-default);
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 0.875rem;
    white-space: nowrap;
    position: relative;
    flex-shrink: 0;
    height: 44px;
  }

  .filter-button:hover {
    border-color: var(--color-brand-500);
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .filter-button.active {
    border-color: var(--color-brand-500);
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .filter-count {
    position: absolute;
    top: -6px;
    right: -6px;
    background: linear-gradient(135deg, var(--color-error-500), var(--color-brand-500));
    color: var(--color-text-on-brand);
    font-size: 0.7rem;
    font-weight: bold;
    width: 18px;
    height: 18px;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 18px;
  }

  .suggestions-enter-active,
  .suggestions-leave-active {
    transition: all 0.2s ease;
  }

  .suggestions-enter-from,
  .suggestions-leave-to {
    opacity: 0;
    transform: translateY(-8px);
  }

  .status-enter-active,
  .status-leave-active {
    transition: all 0.3s ease;
  }

  .status-enter-from,
  .status-leave-to {
    opacity: 0;
    transform: translateY(-10px);
  }

  .search-suggestions::-webkit-scrollbar {
    width: 6px;
  }

  .search-suggestions::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.2);
    border-radius: var(--radius-sm);
  }

  .search-suggestions::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.4);
    border-radius: var(--radius-sm);
  }

  .search-suggestions::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.6);
  }

  @media (max-width: 768px) {
    .advanced-search-bar {
      flex-direction: column;
      gap: 8px;
    }

    .search-mode-tabs {
      width: 100%;
      justify-content: center;
    }

    .search-input-container {
      min-width: auto;
    }

    .mode-tab span {
      display: none;
    }

    .filter-button span {
      display: none;
    }
  }

  @media (max-width: 480px) {
    .search-input {
      font-size: 16px;
    }

    .search-input::placeholder {
      font-size: 14px;
    }
  }
</style>
