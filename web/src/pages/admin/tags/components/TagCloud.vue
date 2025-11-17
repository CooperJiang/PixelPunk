<script setup lang="ts">
  import { ref, computed, watch, onMounted, onUnmounted, type PropType } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { TagInfo } from '@/api/types/index'
  import { useRouter } from 'vue-router'
  import { groupTagsByLetter, getAvailableLetters } from '@/utils/pinyin'

  const { $t } = useTexts()
  const props = defineProps({
    tags: {
      type: Array as PropType<TagInfo[]>,
      required: true,
    },
    isLoading: {
      type: Boolean,
      default: false,
    },
    selectedTagIds: {
      type: Array as PropType<number[]>,
      default: () => [],
    },
    totalTags: {
      type: Number,
      default: 0,
    },
    isBatchMode: {
      type: Boolean,
      default: false,
    },
  })

  const emit = defineEmits(['select', 'edit', 'delete'])

  const router = useRouter()

  const highlightedTagId = ref<number | null>(null)
  const activeLetter = ref<string>('')
  const scrollContainer = ref<HTMLElement | null>(null)

  const isMultiSelectMode = ref(false)
  const tempSelectedTags = ref<number[]>([])

  const isMac = ref(false)
  const modifierKeyName = ref('Ctrl')

  const detectOS = () => {
    const userAgent = window.navigator.userAgent.toLowerCase()
    const platform = window.navigator.platform.toLowerCase()

    isMac.value = /mac|iphone|ipad|ipod/.test(userAgent) || /mac/.test(platform)
    modifierKeyName.value = isMac.value ? 'Command' : 'Ctrl'
  }

  const groupedTags = computed(() => {
    return groupTagsByLetter(props.tags, 'count')
  })

  const availableLetters = computed(() => {
    return getAvailableLetters(groupedTags.value)
  })

  const allLetters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ#'.split('')

  const getGroupInfo = (letter: string) => {
    const tags = groupedTags.value[letter] || []
    const totalCount = tags.reduce((sum, tag) => sum + (tag.count || 0), 0)
    return {
      count: tags.length,
      totalUsage: totalCount,
    }
  }

  const scrollToLetter = (letter: string) => {
    const element = document.querySelector(`[data-letter-group="${letter}"]`)
    if (element && scrollContainer.value) {
      const containerTop = scrollContainer.value.getBoundingClientRect().top
      const elementTop = element.getBoundingClientRect().top
      const offset = elementTop - containerTop - 20

      scrollContainer.value.scrollTo({
        top: scrollContainer.value.scrollTop + offset,
        behavior: 'smooth',
      })
    }
  }

  const handleScroll = () => {
    if (!scrollContainer.value) return

    const containerTop = scrollContainer.value.getBoundingClientRect().top
    const groups = scrollContainer.value.querySelectorAll('[data-letter-group]')

    let currentLetter = ''
    groups.forEach((group) => {
      const rect = group.getBoundingClientRect()
      const groupTop = rect.top - containerTop

      if (groupTop <= 100) {
        currentLetter = group.getAttribute('data-letter-group') || ''
      }
    })

    if (currentLetter) {
      activeLetter.value = currentLetter
    }
  }

  const toggleTagSelection = (tag: TagInfo) => {
    emit('select', tag)
  }

  const goToImagesByTag = (tags: TagInfo | TagInfo[]) => {
    const tagArray = Array.isArray(tags) ? tags : [tags]
    const tagIds = tagArray.map((t) => t.id).join(',')
    const tagNames = tagArray.map((t) => t.name).join(',')

    router.push({
      name: 'adminFiles',
      query: {
        tag_id: tagIds,
        tag_name: tagNames,
      },
    })
  }

  const handleTagClick = (tag: TagInfo) => {
    if (props.isBatchMode) {
      toggleTagSelection(tag)
      return
    }

    if (isMultiSelectMode.value) {
      const index = tempSelectedTags.value.indexOf(tag.id)
      if (index === -1) {
        tempSelectedTags.value.push(tag.id)
      } else {
        tempSelectedTags.value.splice(index, 1)
      }
      return
    }

    goToImagesByTag(tag)
  }

  const checkModifierKey = (event: KeyboardEvent | MouseEvent) => {
    return event.altKey || event.ctrlKey || event.metaKey
  }

  const handleKeyDown = (event: KeyboardEvent) => {
    if (checkModifierKey(event) && !props.isBatchMode) {
      isMultiSelectMode.value = true
    }
  }

  const handleKeyUp = (event: KeyboardEvent) => {
    if (!event.altKey && !event.ctrlKey && !event.metaKey) {
      isMultiSelectMode.value = false

      if (tempSelectedTags.value.length > 0) {
        const selectedTags = props.tags.filter((t) => tempSelectedTags.value.includes(t.id))
        goToImagesByTag(selectedTags)
        tempSelectedTags.value = []
      }
    }
  }

  const handleWindowBlur = () => {
    isMultiSelectMode.value = false
    tempSelectedTags.value = []
  }

  onMounted(() => {
    detectOS()
    window.addEventListener('keydown', handleKeyDown)
    window.addEventListener('keyup', handleKeyUp)
    window.addEventListener('blur', handleWindowBlur)

    if (scrollContainer.value) {
      scrollContainer.value.addEventListener('scroll', handleScroll)
    }
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
    window.removeEventListener('keyup', handleKeyUp)
    window.removeEventListener('blur', handleWindowBlur)

    if (scrollContainer.value) {
      scrollContainer.value.removeEventListener('scroll', handleScroll)
    }
  })

  const positionCache = ref<Map<number, string>>(new Map())

  const estimatedTagsPerRow = ref(Math.floor(window.innerWidth / 200))

  watch(
    () => props.tags.length,
    () => {
      positionCache.value.clear()
      estimatedTagsPerRow.value = Math.floor(window.innerWidth / 200)
    }
  )

  const getActionPosition = (tagId: number) => {
    const cached = positionCache.value.get(tagId)
    if (cached) return cached

    const tagIndex = props.tags.findIndex((tag) => tag.id === tagId)
    if (tagIndex === -1) return 'actions-below'

    const currentRow = Math.floor(tagIndex / estimatedTagsPerRow.value)
    const totalRows = Math.ceil(props.tags.length / estimatedTagsPerRow.value)

    const position = currentRow < totalRows / 2 ? 'actions-below' : 'actions-above'

    positionCache.value.set(tagId, position)
    return position
  }
</script>

<template>
  <div class="tag-cloud-container">
    <div
      v-if="isLoading"
      v-loading="{ loading: true, text: $t('admin.tags.cloud.loading'), type: 'default' }"
      class="loading-state"
    ></div>

    <div v-else-if="tags.length === 0" class="empty-state">
      <div class="empty-content animate-slide-up">
        <div class="empty-icon-wrapper">
          <i class="fas fa-tags empty-icon" />
        </div>
        <h3 class="empty-title">{{ $t('admin.tags.cloud.emptyTitle') }}</h3>
        <p class="empty-desc">{{ $t('admin.tags.cloud.emptyDesc') }}</p>
      </div>
    </div>

    <div v-else class="tags-wrapper">
      <div ref="scrollContainer" class="tag-groups-container">
        <div v-for="letter in availableLetters" :key="letter" :data-letter-group="letter" class="letter-group">
          <div class="letter-group-header">
            <div class="letter-badge">{{ letter }}</div>
            <div class="letter-stats">
              <span class="stat-item">
                <i class="fas fa-tags mr-1"></i>
                {{ getGroupInfo(letter).count }} {{ $t('admin.tags.cloud.tagsCount') }}
              </span>
              <span class="stat-divider">·</span>
              <span class="stat-item">
                <i class="fas fa-chart-line mr-1"></i>
                {{ $t('admin.tags.cloud.usage') }} {{ getGroupInfo(letter).totalUsage }} {{ $t('admin.tags.cloud.times') }}
              </span>
            </div>
          </div>

          <div class="tags-grid">
            <div
              v-for="tag in groupedTags[letter]"
              :key="tag.id"
              :data-tag-id="tag.id"
              class="tag-item"
              :class="[
                { selected: selectedTagIds.includes(tag.id) || tempSelectedTags.includes(tag.id) },
                { highlight: highlightedTagId === tag.id },
                { 'multi-select-mode': isMultiSelectMode },
              ]"
              @mouseenter="highlightedTagId = tag.id"
              @mouseleave="highlightedTagId = null"
              @click="handleTagClick(tag)"
            >
              <div class="tag-content">
                <span class="tag-name">{{ tag.name }}</span>
                <span class="tag-count">{{ tag.count }}</span>
              </div>

              <div
                v-if="highlightedTagId === tag.id && !isBatchMode && !isMultiSelectMode"
                class="tag-actions"
                :class="getActionPosition(tag.id)"
              >
                <button class="action-btn edit-btn" :title="$t('admin.tags.cloud.editTitle')" @click.stop="$emit('edit', tag)">
                  <i class="fas fa-edit" />
                </button>
                <button
                  class="action-btn delete-btn"
                  :title="$t('admin.tags.cloud.deleteTitle')"
                  @click.stop="$emit('delete', tag)"
                >
                  <i class="fas fa-trash-alt" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="letter-navigation">
        <div class="letter-nav-inner">
          <div
            v-for="letter in allLetters"
            :key="letter"
            class="letter-nav-item"
            :class="{
              active: activeLetter === letter,
              disabled: !availableLetters.includes(letter),
            }"
            :title="
              availableLetters.includes(letter)
                ? `${letter} (${getGroupInfo(letter).count})`
                : `${letter} (${$t('admin.tags.cloud.letterEmpty')})`
            "
            @click="availableLetters.includes(letter) && scrollToLetter(letter)"
          >
            {{ letter }}
          </div>
        </div>
      </div>

      <div v-if="!isBatchMode" class="keyboard-hint-fixed">
        <i class="fas fa-keyboard mr-1.5" />
        <span class="hint-text">
          {{ $t('admin.tags.cloud.keyboardHint') }} <kbd class="kbd">{{ modifierKeyName }}</kbd> / <kbd class="kbd">Alt</kbd>
          {{ $t('admin.tags.cloud.keyboardHintAction') }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .tag-cloud-container {
    @apply flex h-full flex-col;
    flex: 1;
    min-height: 0;
  }

  .loading-state {
    @apply flex flex-1 items-center justify-center;
    min-height: 100%;
    height: 100%;
  }

  .empty-state {
    @apply flex flex-1 items-center justify-center;
    padding: var(--space-3xl) var(--space-md);
    min-height: 100%;
    height: 100%;
  }

  .empty-content {
    @apply text-center;
    max-width: 400px;
  }

  .empty-icon-wrapper {
    margin-bottom: var(--space-md);
    position: relative;
    display: inline-block;
  }

  .empty-icon {
    font-size: var(--text-5xl);
    color: var(--color-brand-500);
    filter: drop-shadow(var(--shadow-glow-md));
    animation: float 3s var(--ease-in-out) infinite;
  }

  .empty-title {
    margin-bottom: var(--space-sm);
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    color: var(--color-content-default);
    letter-spacing: var(--tracking-tight);
  }

  .empty-desc {
    font-size: var(--text-sm);
    color: var(--color-content-muted);
    line-height: var(--leading-relaxed);
  }

  .tags-wrapper {
    @apply relative flex h-full;
    flex: 1;
    min-height: 0;
    gap: 0;
  }

  .tag-groups-container {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    overflow-x: hidden;
    padding-right: 30px;

    scrollbar-width: thin;
    scrollbar-color: rgba(var(--color-brand-500-rgb), 0.4) transparent;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(var(--color-brand-500-rgb), 0.4);
      border-radius: var(--radius-full);

      &:hover {
        background: rgba(var(--color-brand-500-rgb), 0.6);
      }
    }
  }

  .letter-group {
    margin-bottom: var(--space-lg);

    &:last-child {
      margin-bottom: 0;
    }
  }

  .letter-group-header {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    margin-bottom: var(--space-sm);
    padding: 0 var(--space-lg) var(--space-sm) var(--space-md);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .letter-badge {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.08));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-md);
    font-size: var(--text-lg);
    font-weight: var(--font-bold);
    color: var(--color-brand-500);
    box-shadow: var(--shadow-sm);
    flex-shrink: 0;
  }

  .letter-stats {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    font-size: var(--text-xs);
    color: var(--color-content-muted);

    .stat-item {
      display: flex;
      align-items: center;

      i {
        opacity: 0.7;
      }
    }

    .stat-divider {
      opacity: 0.4;
    }
  }

  .tags-grid {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-xs);
    align-items: flex-start;
    padding: 0 var(--space-lg) 0 var(--space-md);
  }

  .letter-navigation {
    position: sticky;
    top: var(--space-md);
    right: var(--space-lg);
    height: fit-content;
    max-height: calc(100vh - 180px);
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: var(--space-xs) 0;
    z-index: 20;
  }

  .letter-nav-inner {
    display: flex;
    flex-direction: column;
    gap: 1px;
    background: rgba(var(--color-background-800-rgb), 0.95);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-md);
    padding: 4px;
    box-shadow: var(--shadow-md);
    backdrop-filter: var(--backdrop-blur-md);
  }

  .letter-nav-item {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 26px;
    height: 22px;
    font-size: 11px;
    font-weight: var(--font-semibold);
    color: var(--color-content-muted);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition:
      background var(--transition-fast) var(--ease-out),
      color var(--transition-fast) var(--ease-out),
      transform var(--transition-fast) var(--ease-out),
      border-color var(--transition-fast) var(--ease-out);
    user-select: none;

    &:hover:not(.disabled) {
      background: rgba(var(--color-brand-500-rgb), 0.15);
      color: var(--color-brand-400);
      transform: translateX(-2px);
    }

    &.active {
      background: rgba(var(--color-brand-500-rgb), 0.25);
      color: var(--color-brand-500);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
      box-shadow: var(--shadow-xs);
      font-weight: var(--font-bold);
    }

    &.disabled {
      opacity: 0.25;
      cursor: not-allowed;
    }
  }

  .tag-item {
    @apply relative inline-flex cursor-pointer items-center justify-between text-content;
    border-radius: var(--radius-sm);
    padding: 4px 8px;
    background: var(--color-hover-bg-neutral);
    border: 1px solid var(--color-border-subtle);
    box-shadow: none;
    transition:
      border-color var(--transition-normal) var(--ease-out),
      box-shadow var(--transition-normal) var(--ease-out),
      background var(--transition-normal) var(--ease-out),
      transform var(--transition-fast) var(--ease-out);
  }

  .tag-item:hover {
    border-color: var(--color-hover-border);
    background: var(--color-hover-bg);
    transform: translateY(-1px);
    box-shadow: var(--shadow-sm);
  }

  .tag-content {
    @apply flex items-center;
    gap: 6px;
  }

  .tag-name {
    font-weight: var(--font-medium);
    @apply text-content;
    letter-spacing: var(--tracking-tight);
    font-size: var(--text-sm);
  }

  .tag-count {
    font-size: 11px;
    font-weight: var(--font-semibold);
    color: var(--color-brand-400);
    opacity: 0.8;
    font-family: var(--font-mono);
  }

  .multi-select-mode {
    cursor: pointer;
  }

  .selected {
    position: relative;
    border-color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.2) !important;
    box-shadow:
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.3),
      var(--shadow-md) !important;

    .tag-name {
      @apply text-content-heading;
      font-weight: var(--font-semibold);
    }

    .tag-count {
      color: var(--color-brand-300);
      opacity: 1;
    }
  }

  .selected::after {
    content: '';
    position: absolute;
    top: calc(var(--space-xs) * -1);
    right: calc(var(--space-xs) * -1);
    width: var(--space-sm);
    height: var(--space-sm);
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    border: 2px solid var(--color-background-900);
    box-shadow: var(--shadow-glow-sm);
  }

  .highlight {
    @apply z-30;
    box-shadow:
      var(--shadow-lg),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .tag-actions {
    @apply absolute right-0 flex;
    gap: 4px;
    border-radius: var(--radius-md);
    padding: 4px;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.98) 0%,
      rgba(var(--color-background-700-rgb), 0.95) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow:
      0 4px 12px rgba(0, 0, 0, 0.3),
      0 2px 6px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    z-index: 50;
    opacity: 0;
    transform: translateY(-4px) scale(0.92);
    pointer-events: none;
    transition:
      opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1),
      transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .tag-item:hover .tag-actions,
  .tag-actions:hover {
    opacity: 1;
    transform: translateY(0) scale(1);
    pointer-events: auto;
  }

  .actions-below {
    top: calc(100% + var(--space-xs));
  }

  .actions-above {
    bottom: calc(100% + var(--space-xs));
  }

  .action-btn {
    @apply relative;
    border-radius: var(--radius-sm);
    padding: 6px;
    font-size: 12px;
    background: rgba(var(--color-background-600-rgb), 0.3);
    border: 1px solid transparent;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    cursor: pointer;
    transition:
      background-color 0.15s cubic-bezier(0.4, 0, 0.2, 1),
      border-color 0.15s cubic-bezier(0.4, 0, 0.2, 1),
      color 0.15s cubic-bezier(0.4, 0, 0.2, 1),
      transform 0.15s cubic-bezier(0.4, 0, 0.2, 1),
      box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      transform: translateY(-1px);
    }

    &:active {
      transform: translateY(0) scale(0.96);
    }

    i {
      transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1);
    }

    &:hover i {
      transform: scale(1.1);
    }
  }

  .edit-btn {
    color: var(--color-content-muted);

    &:hover {
      color: var(--color-brand-400);
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15) 0%, rgba(var(--color-brand-500-rgb), 0.08) 100%);
      border-color: rgba(var(--color-brand-500-rgb), 0.3);
      box-shadow:
        0 2px 8px rgba(var(--color-brand-500-rgb), 0.2),
        inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.1);
    }

    &:active {
      box-shadow:
        0 1px 4px rgba(var(--color-brand-500-rgb), 0.15),
        inset 0 1px 2px rgba(0, 0, 0, 0.1);
    }
  }

  .delete-btn {
    color: var(--color-content-muted);

    &:hover {
      color: var(--color-error-400);
      background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.15) 0%, rgba(var(--color-error-rgb), 0.08) 100%);
      border-color: rgba(var(--color-error-rgb), 0.3);
      box-shadow:
        0 2px 8px rgba(var(--color-error-rgb), 0.2),
        inset 0 1px 0 rgba(var(--color-error-rgb), 0.1);
    }

    &:active {
      box-shadow:
        0 1px 4px rgba(var(--color-error-rgb), 0.15),
        inset 0 1px 2px rgba(0, 0, 0, 0.1);
    }
  }

  .keyboard-hint-fixed {
    position: fixed;
    bottom: var(--space-lg);
    right: var(--space-lg);
    display: flex;
    align-items: center;
    padding: var(--space-xs) var(--space-md);
    background: rgba(var(--color-background-800-rgb), 0.95);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-lg);
    backdrop-filter: var(--backdrop-blur-md);
    color: var(--color-content-muted);
    opacity: 0.7;
    transition: opacity var(--transition-normal) var(--ease-out);
    z-index: 30;

    &:hover {
      opacity: 1;
    }

    i {
      font-size: var(--text-xs);
      opacity: 0.7;
    }

    .hint-text {
      font-size: var(--text-xs);
      letter-spacing: var(--tracking-tight);
      line-height: var(--leading-normal);
    }

    .kbd {
      margin: 0 var(--space-xs);
      @apply inline-block;
      border-radius: var(--radius-sm);
      padding: 2px 6px;
      font-size: var(--text-xs);
      font-weight: var(--font-semibold);
      font-family: var(--font-mono);
      background: rgba(var(--color-background-600-rgb), 0.5);
      border: 1px solid var(--color-border-subtle);
      box-shadow: var(--shadow-xs);
      color: var(--color-content);
      min-width: var(--space-lg);
      text-align: center;
      letter-spacing: var(--tracking-tight);
    }
  }

  @media (max-width: 768px) {
    .tag-cloud {
      gap: var(--space-xs);
      padding: var(--space-lg) var(--space-sm) var(--space-sm);
    }

    .tag-groups-container {
      padding: var(--space-sm);
      padding-right: 40px;
    }

    .letter-navigation {
      right: var(--space-xs);
    }

    .letter-nav-item {
      width: 22px;
      height: 20px;
      font-size: 10px;
    }

    .keyboard-hint-fixed {
      display: none;
    }

    .tag-item {
      padding: 3px 6px;
    }

    .tag-content {
      gap: 4px;
    }

    .tag-actions {
      padding: 3px;
      gap: 3px;
    }

    .action-btn {
      width: 26px;
      height: 26px;
      padding: 5px;
      font-size: 11px;
    }
  }
</style>
