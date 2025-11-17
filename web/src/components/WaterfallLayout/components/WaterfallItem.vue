<template>
  <div class="waterfall-item" :style="getItemStyle(item)">
    <div
      class="waterfall-image-container"
      :class="{ 'is-selectable': selectable, 'is-selected': isSelected(item.id) }"
      :style="getFileContainerStyle(item)"
      @click="handleImageClick(item)"
    >
      <img
        v-if="item.visible"
        class="waterfall-image"
        :src="getImageSrc(item)"
        :alt="item.name || 'Image'"
        @load="handleImageLoad(item)"
        @error="handleImageError(item)"
        loading="lazy"
      />

      <div v-if="!item.loaded && item.visible" class="image-loading">
        <div class="loading-spinner" />
      </div>

      <div class="image-info-tags">
        <span v-if="item.resolution" class="info-tag">{{ item.resolution }}</span>
        <span v-if="item.format" class="info-tag">{{ item.format.toUpperCase() }}</span>
        <span v-if="item.size" class="info-tag">{{ formatSize(item.size) }}</span>
      </div>

      <div v-if="item.width && item.height" class="image-info-tags bottom-tags">
        <span class="info-tag dimension-tag">{{ item.width }} × {{ item.height }}</span>
      </div>

      <div v-if="selectable" class="selection-indicator">
        <i :class="isSelected(item.id) ? 'fas fa-check-circle' : 'far fa-circle'" />
      </div>

      <div v-if="item.is_duplicate" class="duplicate-indicator">
        <i class="fas fa-copy" />
      </div>

      <div v-if="item.hasError" class="image-error">
        <i class="fas fa-exclamation-triangle" />
        <span>{{ $t('waterfallLayout.loadFailed') }}</span>
      </div>
    </div>

    <div class="image-info">
      <div class="image-name" :title="item.name">{{ item.name }}</div>

      <div class="image-meta">
        <span class="image-date">
          {{ formatDate(item.created_at) }}
        </span>
        <span class="image-size"> {{ item.width }}x{{ item.height }} </span>
      </div>

      <div v-if="item.tags && item.tags.length > 0" class="image-tags">
        <div class="tags-container">
          <template v-if="!item.showAllTags">
            <span
              v-for="tag in item.tags.slice(0, maxTagsToShow)"
              :key="getTagKey(tag)"
              class="image-tag"
              :title="getTagName(tag)"
            >
              {{ getTagName(tag) }}
            </span>
            <span v-if="item.tags.length > maxTagsToShow" class="more-tags" @click.stop="handleToggleTags(item.id)">
              +{{ item.tags.length - maxTagsToShow }}
            </span>
          </template>
          <template v-else>
            <span v-for="tag in item.tags" :key="getTagKey(tag)" class="image-tag" :title="getTagName(tag)">
              {{ getTagName(tag) }}
            </span>
            <span class="more-tags" @click.stop="handleToggleTags(item.id)">
              {{ $t('waterfallLayout.collapse') }}
            </span>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { formatDate, formatSize } from '@/utils/formatting/format'
  import { LAYOUT } from '@/constants'
  import { useTexts } from '@/composables/useTexts'
  import type { WaterfallItemProps, WaterfallItemEmits } from '../types'

  const props = defineProps<WaterfallItemProps>()
  const emit = defineEmits<WaterfallItemEmits>()
  const { $t } = useTexts()

  const maxTagsToShow = LAYOUT.WATERFALL.MAX_TAGS_TO_SHOW

  const getItemStyle = (item) => {
    return {
      marginBottom: `${props.gap}px`,
      opacity: item.visible ? 1 : 0,
      transition: 'opacity 0.3s ease',
    }
  }

  const getFileContainerStyle = (item) => {
    const style = {
      position: 'relative',
      borderRadius: '8px',
      overflow: 'hidden',
      background: 'rgba(var(--color-background-800-rgb), 0.5)',
      border: '1px solid var(--color-border-subtle)',
      transition: 'all 0.25s ease',
      cursor: 'pointer',
      width: '100%',
    }

    if (!item.loaded && item.width && item.height) {
      const aspectRatio = item.width / item.height
      const effectiveColumnWidth = Math.max(props.columnWidth * 1.2, 180)
      const calculatedHeight = Math.round(effectiveColumnWidth / aspectRatio)

      style.minHeight = `${Math.max(calculatedHeight, 100)}px`
      style.background = 'rgba(var(--color-background-800-rgb), 0.7)'
    } else if (!item.loaded) {
      style.minHeight = '200px'
      style.background = 'rgba(var(--color-background-800-rgb), 0.7)'
    }

    return style
  }

  const getImageSrc = (item) => {
    const url = item.full_thumb_url || item.thumb_url || item.full_url || item.url || ''
    return url
  }

  const isSelected = (id) => {
    return props.selectedIds.includes(id)
  }

  const handleImageClick = (item) => {
    emit('image-click', item)
  }

  const handleImageLoad = (item) => {
    emit('image-load', item)
  }

  const handleImageError = (item) => {
    emit('image-error', item)
  }

  const handleToggleTags = (id) => {
    emit('toggle-tags', id)
  }

  const getTagName = (tag) => {
    if (typeof tag === 'string') {
      return tag
    }
    if (typeof tag === 'object' && tag !== null) {
      return tag.name || tag.toString()
    }
    return String(tag)
  }

  const getTagKey = (tag) => {
    if (typeof tag === 'object' && tag !== null && tag.id) {
      return tag.id
    }
    return getTagName(tag)
  }
</script>

<style scoped>
  .waterfall-item {
    @apply w-full;
    break-inside: avoid;
  }

  .waterfall-image-container {
    @apply relative cursor-pointer overflow-hidden rounded-lg transition-all duration-300;
    background: var(--color-background-800);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
  }

  .waterfall-image-container:hover {
    @apply -translate-y-0.5;
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow:
      0 6px 14px rgba(0, 0, 0, 0.3),
      0 0 8px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .waterfall-image-container.is-selected {
    border-color: rgba(var(--color-brand-500-rgb), 0.7);
    box-shadow:
      0 0 0 2px rgba(var(--color-brand-500-rgb), 0.4),
      0 6px 14px rgba(0, 0, 0, 0.3);
  }

  .waterfall-image {
    @apply block h-auto w-full transition-all duration-300;
  }

  .image-info-tags {
    @apply absolute left-2 top-2 z-[4] flex gap-1;
  }

  .bottom-tags {
    @apply bottom-2 top-auto;
  }

  .info-tag {
    @apply rounded border px-1.5 py-0.5 text-xs backdrop-blur;
    background: var(--color-background-800);
    color: var(--color-content-default);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .dimension-tag {
    background: var(--color-background-900);
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .image-loading {
    @apply absolute inset-0 flex items-center justify-center;
    background: rgba(var(--color-background-800-rgb), 0.8);
  }

  .loading-spinner {
    @apply h-8 w-8 rounded-full border-2 border-t-transparent;
    border-color: var(--color-brand-500);
    animation: spin 1s linear infinite;
  }

  .selection-indicator {
    @apply absolute right-2 top-2 z-10 flex h-6 w-6 cursor-pointer items-center justify-center rounded-full text-sm transition-all duration-200;
    background: var(--color-background-700);
    color: var(--color-content-default);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
  }

  .waterfall-image-container.is-selected .selection-indicator {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-content-default);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.6);
  }

  .duplicate-indicator {
    @apply absolute left-2 top-2 flex h-6 w-6 items-center justify-center rounded-full text-xs;
    background: rgba(var(--color-info-rgb), 0.85);
    color: var(--color-text-on-brand);
  }

  .image-error {
    @apply absolute inset-0 flex flex-col items-center justify-center text-sm text-error-400;
    background: rgba(var(--color-background-900-rgb), 0.9);
  }

  .image-info {
    @apply p-3;
    background: var(--color-background-700);
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .image-name {
    @apply mb-1 overflow-hidden text-ellipsis whitespace-nowrap text-sm font-medium;
    color: rgba(var(--color-content-default-rgb), 0.9);
  }

  .image-meta {
    @apply mb-2 flex justify-between overflow-hidden text-ellipsis whitespace-nowrap text-xs;
    color: rgba(var(--color-content-default-rgb), 0.6);
  }

  .image-tags {
    @apply mt-2 overflow-hidden;
  }

  .tags-container {
    @apply flex flex-wrap gap-1 overflow-x-hidden whitespace-nowrap;
  }

  .image-tag {
    @apply max-w-20 flex-shrink-0 overflow-hidden text-ellipsis whitespace-nowrap rounded px-1.5 py-0.5 text-xs transition-all duration-200;
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .image-tag:hover {
    background: rgba(var(--color-brand-500-rgb), 0.25);
    color: var(--color-content-default);
  }

  .more-tags {
    @apply flex-shrink-0 cursor-pointer whitespace-nowrap rounded px-1.5 py-0.5 text-xs transition-all duration-200;
    background: var(--color-background-600);
    color: rgba(var(--color-content-default-rgb), 0.7);
  }

  .more-tags:hover {
    background: var(--color-background-500);
    color: var(--color-content-default);
  }
</style>
