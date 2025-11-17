<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { formatDate, formatSize } from '@/utils/formatting/format'
  import { downloadSharedFileQuick } from '@/utils/file/downloader'
  import { useTexts } from '@/composables/useTexts'
  import type { ShareFileProps, ShareFileEmits, ShareFileData } from './types'

  defineOptions({
    name: 'CyberShareFile',
  })

  const props = withDefaults(defineProps<ShareFileProps>(), {
    selectable: false,
    selected: false,
    accessToken: '',
  })

  const emit = defineEmits<ShareFileEmits>()
  const texts = useTexts()

  const isLoading = ref(true)
  const loadError = ref(false)
  const showAllTags = ref(false)
  const maxTagsToShow = 2 // 只显示2个标签，其余显示 +X
  let loadingTimeout: NodeJS.Timeout | null = null

  const handleImageLoadStart = () => {
    setLoadingTimeout()
  }

  const handleImageLoaded = () => {
    if (loadingTimeout) {
      clearTimeout(loadingTimeout)
      loadingTimeout = null
    }
    isLoading.value = false
    loadError.value = false
  }

  const handleImageError = () => {
    if (loadingTimeout) {
      clearTimeout(loadingTimeout)
      loadingTimeout = null
    }
    isLoading.value = false
    loadError.value = true
  }

  const setLoadingTimeout = () => {
    if (loadingTimeout) {
      clearTimeout(loadingTimeout)
    }
    loadingTimeout = setTimeout(() => {
      if (import.meta.env.DEV) {
      }
      isLoading.value = false
      loadError.value = true
    }, 20000)
  }

  const displayedTags = computed(() => {
    if (!props.file.tags || props.file.tags.length === 0) {
      return []
    }
    return showAllTags.value ? props.file.tags : props.file.tags.slice(0, maxTagsToShow)
  })

  const showMoreTags = computed(() => props.file.tags && props.file.tags.length > maxTagsToShow && !showAllTags.value)

  const toggleShowAllTags = (e: Event) => {
    e.stopPropagation()
    showAllTags.value = !showAllTags.value
  }

  const getFileUrl = (file: ShareFileData) => {
    const url = file.full_thumb_url || file.thumb_url || file.full_url || file.url || ''
    if (!url && import.meta.env.DEV) {
      console.warn('Empty file URL:', file)
    }
    return url
  }

  const handleClick = () => {
    if (props.selectable) {
      emit('select', props.file.id)
    } else {
      emit('click', props.file)
    }
  }

  const handleSelectClick = () => {
    emit('select', props.file.id)
  }

  const handleDownload = async (e: Event) => {
    e.stopPropagation()
    try {
      const fileName = props.file.display_name || props.file.original_name || `file_${props.file.id}`
      await downloadSharedFileQuick(props.file.id, props.shareKey, props.accessToken, fileName)
    } catch (error: unknown) {
      if (import.meta.env.DEV) {
        console.error('ShareImage download failed:', error)
      }
      if (window.toast) {
        window.toast.error(texts.share.image.downloadError)
      }
    }
  }

  onMounted(() => {
    if (getFileUrl(props.file)) {
      isLoading.value = true
    } else {
      isLoading.value = false
      loadError.value = true
    }
  })

  onUnmounted(() => {
    if (loadingTimeout) {
      clearTimeout(loadingTimeout)
      loadingTimeout = null
    }
  })
</script>

<template>
  <div
    class="share-image-card"
    :class="{ 'is-selectable': props.selectable, 'is-selected': props.selected }"
    @click="handleClick"
  >
    <div class="image-preview" :class="{ 'is-loading': isLoading }">
      <div v-if="isLoading" class="loading-indicator">
        <div class="loading-spinner">
          <i class="fas fa-spin fa-circle-notch" />
        </div>
      </div>

      <img
        :src="getFileUrl(props.file)"
        :alt="props.file.display_name || texts.share.image.alt"
        loading="lazy"
        :style="{ visibility: isLoading ? 'hidden' : 'visible' }"
        @load="handleImageLoaded"
        @error="handleImageError"
        @loadstart="handleImageLoadStart"
      />

      <div class="image-overlay">
        <div class="overlay-content">
          <div class="action-buttons">
            <button class="action-button download-button" @click.stop="handleDownload">
              <i class="fas fa-download" />
            </button>
          </div>
        </div>
      </div>

      <div v-if="props.selectable" class="selection-indicator" @click.stop="handleSelectClick">
        <i :class="props.selected ? 'fas fa-check-circle' : 'far fa-circle'" />
      </div>

      <div class="image-info-tags">
        <span v-if="props.file.resolution" class="info-tag">{{ props.file.resolution }}</span>
        <span v-if="props.file.format" class="info-tag">{{ props.file.format.toUpperCase() }}</span>
        <span v-if="props.file.size_formatted" class="info-tag">{{ formatSize(props.file.size) }}</span>
      </div>

      <div v-if="props.file.width && props.file.height" class="image-info-tags bottom-tags">
        <span class="info-tag dimension-tag">{{ props.file.width }} × {{ props.file.height }}</span>
      </div>
    </div>

    <div class="image-info">
      <div class="image-name" :title="props.file.display_name || props.file.original_name">
        {{ props.file.display_name || props.file.original_name }}
      </div>
      <div class="image-meta">
        <span v-if="props.file.size" class="image-size">{{ formatSize(props.file.size) }}</span>
        <span v-if="props.file.created_at" class="image-date">{{ formatDate(props.file.created_at) }}</span>
      </div>

      <div v-if="props.file.tags && props.file.tags.length > 0" class="image-tags">
        <div class="tags-container">
          <span v-for="tag in displayedTags" :key="tag.id" class="image-tag" :title="tag.name">
            {{ tag.name }}
          </span>
          <span v-if="showMoreTags" class="more-tags" @click.stop="toggleShowAllTags">
            +{{ props.file.tags.length - maxTagsToShow }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .share-image-card {
    @apply relative flex h-full cursor-pointer flex-col overflow-hidden rounded-lg transition-all duration-300;
    background: var(--color-background-800);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
    transition-timing-function: cubic-bezier(0.3, 0, 0.2, 1);
  }

  .share-image-card:hover,
  .share-image-card:active {
    @apply -translate-y-0.5;
    box-shadow:
      0 6px 14px rgba(0, 0, 0, 0.2),
      0 0 8px rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .share-image-card:active {
    @apply translate-y-0;
  }

  .share-image-card.is-selectable {
    @apply relative;
  }

  .share-image-card.is-selected {
    @apply -translate-y-0.5;
    border-color: rgba(var(--color-brand-500-rgb), 0.7);
    box-shadow:
      0 0 0 2px rgba(var(--color-brand-500-rgb), 0.4),
      0 6px 14px rgba(0, 0, 0, 0.2);
  }

  .selection-indicator {
    @apply absolute right-2 top-2 z-[5] flex h-7 w-7 items-center justify-center rounded-full border opacity-100 transition-all duration-200;
    background: var(--color-background-700);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    visibility: visible !important;
  }

  .share-image-card.is-selected .selection-indicator {
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .selection-indicator i {
    @apply text-base;
    color: var(--color-content-muted);
  }

  .share-image-card.is-selected .selection-indicator i {
    color: var(--color-brand-500);
  }

  .image-preview {
    @apply relative aspect-square overflow-hidden border-b;
    background: var(--color-background-900);
    border-color: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .image-preview.is-loading {
    background: var(--color-background-900);
  }

  .loading-indicator {
    @apply absolute inset-0 z-[2] flex items-center justify-center;
    background: rgba(var(--color-background-900-rgb), 0.2);
  }

  .loading-spinner {
    @apply text-2xl;
    color: rgba(var(--color-brand-500-rgb), 0.8);
    animation: pulse 1.5s infinite ease-in-out;
  }

  .image-preview img {
    @apply h-full w-full object-cover transition-transform duration-500;
  }

  .share-image-card:hover .image-preview img {
    @apply scale-105;
  }

  .image-overlay {
    @apply absolute inset-0 z-[3] flex flex-col justify-end p-3 opacity-0 transition-all duration-300;
    background: rgba(0, 0, 0, 0);
  }

  .share-image-card:hover .image-overlay {
    @apply opacity-100;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.2));
  }

  .overlay-content {
    @apply flex justify-end;
  }

  .action-buttons {
    @apply flex translate-y-2.5 gap-2 opacity-0 transition-all duration-300;
  }

  .share-image-card:hover .action-buttons {
    @apply translate-y-0 opacity-100;
  }

  .action-button {
    @apply flex h-[30px] w-[30px] cursor-pointer items-center justify-center rounded-full border transition-all duration-200;
    background: rgba(var(--color-brand-500-rgb), 0.25);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    color: var(--color-content-default);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }

  .action-button:hover {
    @apply scale-110;
    background: rgba(var(--color-brand-500-rgb), 0.6);
    border-color: rgba(var(--color-brand-500-rgb), 0.8);
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .action-button:active {
    @apply scale-95;
  }

  .download-button i {
    @apply text-sm;
    color: var(--color-content-default);
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

  .image-info {
    @apply flex flex-1 flex-col p-2;
    background: var(--color-background-700);
  }

  .image-name {
    @apply mb-1 truncate text-sm font-medium;
    color: rgba(var(--color-content-default-rgb), 0.9);
  }

  .image-meta {
    @apply mb-2 flex justify-between text-xs;
    color: rgba(var(--color-content-default-rgb), 0.6);
  }

  .image-size,
  .image-date {
    @apply truncate;
  }

  .image-tags {
    @apply mt-auto overflow-hidden;
  }

  .tags-container {
    @apply flex gap-1 overflow-x-hidden whitespace-nowrap;
    flex-wrap: nowrap;
  }

  .image-tag {
    @apply flex-shrink-0 truncate whitespace-nowrap rounded px-1.5 py-0.5 text-xs transition-all duration-200;
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    max-width: 80px;
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

  @media (max-width: 768px) {
    .share-image-card {
      @apply rounded-md;
    }

    .action-button {
      @apply h-[26px] w-[26px];
    }

    .download-button i {
      @apply text-xs;
    }

    .selection-indicator {
      @apply h-6 w-6;
    }

    .selection-indicator i {
      @apply text-sm;
    }

    .image-info {
      @apply p-1.5;
    }

    .image-name {
      @apply text-sm;
    }

    .image-meta {
      @apply text-xs;
    }

    .info-tag {
      @apply px-1 py-0.5 text-xs;
    }

    .image-tag,
    .more-tags {
      @apply px-1 py-px text-xs;
    }
  }
</style>
