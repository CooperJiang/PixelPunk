<script setup lang="ts">
  import { computed, defineEmits, defineProps } from 'vue'

  import { useTexts } from '@/composables/useTexts'

  const props = defineProps({
    title: {
      type: String,
      default: '',
    },
    count: {
      type: Number,
      default: 0,
    },
    currentSort: {
      type: String,
      default: 'date',
    },
    sortDirection: {
      type: String,
      default: 'desc',
    },
    sortOptions: {
      type: Array,
      default: () => [],
    },
    layoutMode: {
      type: String,
      default: 'grid',
    },
    isFullScreen: {
      type: Boolean,
      default: false,
    },
    isSelectionMode: {
      type: Boolean,
      default: false,
    },
  })

  const emit = defineEmits(['sort-change', 'toggle-fullscreen', 'toggle-selection', 'open-grid-size', 'open-layout'])

  const { $t } = useTexts()

  /* Computed properties for translated defaults */
  const displayTitle = computed(() => props.title || $t('share.file.title'))

  const displaySortOptions = computed(() => {
    if (props.sortOptions.length > 0) {
      return props.sortOptions
    }
    return [
      { value: 'date', label: $t('share.sort.date') },
      { value: 'name', label: $t('share.sort.name') },
      { value: 'size', label: $t('share.sort.size') },
    ]
  })

  const handleSortChange = (sortType: string) => {
    emit('sort-change', sortType)
  }

  const toggleFullScreen = () => {
    emit('toggle-fullscreen')
  }

  const toggleSelectionMode = () => {
    emit('toggle-selection')
  }

  const openGridSizePopover = (e: Event) => {
    e.stopPropagation()
    emit('open-grid-size', e)
  }

  const openLayoutPopover = (e: Event) => {
    e.stopPropagation()
    emit('open-layout', e)
  }

  const getLayoutIcon = computed(() => {
    switch (props.layoutMode) {
      case 'waterfall':
        return 'fas fa-stream'
      case 'masonry':
        return 'fas fa-th-large'
      case 'large':
        return 'fas fa-image'
      default:
        return 'fas fa-th'
    }
  })
</script>

<template>
  <div class="section-header">
    <div class="section-title">
      <i class="fas fa-images" />
      <span>{{ displayTitle }}</span>
      <span class="count-badge">{{ count }}</span>
    </div>

    <div class="sort-options">
      <button
        v-for="option in displaySortOptions"
        :key="option.value"
        :class="{ active: currentSort === option.value }"
        class="sort-btn"
        @click="handleSortChange(option.value)"
      >
        {{ option.label }}
        <i v-if="currentSort === option.value" :class="sortDirection === 'asc' ? 'fas fa-sort-up' : 'fas fa-sort-down'" />
      </button>

      <button class="sort-btn grid-size-btn" :title="$t('share.sort.adjustGridSize')" @click="openGridSizePopover">
        <i class="fas fa-sliders-h" />
      </button>

      <button
        class="sort-btn fullscreen-btn"
        :title="isFullScreen ? $t('share.sort.exitFullscreen') : $t('share.sort.fullscreen')"
        @click="toggleFullScreen"
      >
        <i :class="isFullScreen ? 'fas fa-compress' : 'fas fa-expand'" />
      </button>

      <button class="sort-btn layout-btn" :title="$t('share.sort.selectLayout')" @click="openLayoutPopover">
        <i :class="getLayoutIcon" />
      </button>

      <button class="sort-btn selection-btn" :class="{ active: isSelectionMode }" @click="toggleSelectionMode">
        <i class="fas fa-check-square" />
        {{ isSelectionMode ? $t('share.sort.exitSelect') : $t('share.sort.multiSelect') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    position: relative;
    width: 100%;
  }

  .section-title {
    display: flex;
    align-items: center;
    color: var(--color-content-default);
    font-size: 1rem;
    font-weight: 500;
  }

  .section-title i {
    margin-right: 0.5rem;
    color: var(--color-brand-500);
  }

  .count-badge {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    border-radius: var(--radius-md);
    padding: 0.375rem 0.75rem;
    font-size: 0.75rem;
    margin-left: 0.5rem;
  }

  .sort-options {
    display: flex;
    gap: var(--space-sm);
    flex-wrap: wrap;
    align-items: center;
    justify-content: flex-end;
    min-width: 60%;
  }

  .sort-btn {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    color: rgba(255, 255, 255, 0.8);
    padding: 0.375rem 0.75rem;
    font-size: 0.75rem;
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    transition: all 0.2s ease;
    cursor: pointer;
  }

  .sort-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .sort-btn.active {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .sort-btn i {
    font-size: 0.7rem;
  }

  .sort-btn i.fa-th-large,
  .sort-btn i.fa-stream {
    color: rgba(var(--color-brand-500-rgb), 0.75);
  }

  @media (max-width: 767px) {
    .section-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .sort-options {
      margin-top: 0.5rem;
      width: 100%;
      justify-content: flex-start;
    }
  }
</style>
