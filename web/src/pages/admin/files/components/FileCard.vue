<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { FileInfo } from '@/api/admin/files'
  // import { COMMON_ACTIONS } from '@/constants'

  const props = defineProps<{
    fileData: FileInfo
    selectMode?: boolean
    selected?: boolean
  }>()

  const { $t } = useTexts()
  const isFileLoading = ref(false)

  const emit = defineEmits<{
    (e: 'view', id: string): void
    (e: 'delete', id: string): void
    (e: 'recommend', id: string, recommended: boolean): void
    (e: 'details', id: string): void
    (e: 'select', id: string): void
  }>()

  const isNsfw = computed(
    () => props.fileData.ai_info?.is_nsfw
    /* return props.fileData.ai_info?.is_nsfw && props.fileData.ai_info?.nsfw_score > 0.7; */
  )

  /* 是否显示违规内容 */
  const showNsfwContent = ref(false)

  const handleCardClick = () => {
    if (props.selectMode) {
      emit('select', props.fileData.id)
    } else {
      if (isNsfw.value && !showNsfwContent.value) {
        return
      }
      emit('view', props.fileData.id)
    }
  }

  const viewNsfwImage = () => {
    showNsfwContent.value = true
    emit('view', props.fileData.id)
  }

  const deleteImage = () => {
    emit('delete', props.fileData.id)
  }

  const toggleRecommend = () => {
    emit('recommend', props.fileData.id, !props.fileData.is_recommended)
  }

  const viewDetails = () => {
    emit('details', props.fileData.id)
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) {
      return '0 B'
    }

    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    })
  }

  const _handleSelectClick = () => {
    emit('select', props.fileData.id)
  }
</script>

<template>
  <div
    class="image-card hover:shadow-glow group overflow-hidden rounded-lg border border-subtle bg-background-800"
    :class="{
      'is-nsfw': isNsfw,
      selected: selected,
      'select-mode': selectMode,
    }"
    @click="handleCardClick"
  >
    <div v-if="selectMode && selected" class="selected-indicator">
      <i class="fas fa-check"></i>
    </div>

    <div class="relative">
      <div
        v-if="isNsfw && !selectMode"
        class="nsfw-overlay absolute inset-0 z-10 flex flex-col items-center justify-center p-3 backdrop-blur-sm"
      >
        <i class="fas fa-eye-slash mb-2 text-2xl text-error-500" />
        <p class="text-center text-sm text-content-heading">{{ $t('admin.files.fileCard.nsfwOverlay.warning') }}</p>
        <CyberButton type="primary" class="mt-3" @click.stop="viewNsfwImage">
          <i class="fas fa-eye mr-1" />{{ $t('admin.files.fileCard.nsfwOverlay.viewFile') }}
        </CyberButton>
      </div>

      <div class="image-wrapper relative aspect-[4/3]" v-loading="isFileLoading">
        <CyberFile
          :src="fileData.full_thumb_url || fileData.thumb_url"
          :alt="fileData.display_name || fileData.original_name"
          class="h-full w-full object-cover"
          loading="lazy"
          @loading="isFileLoading = $event"
        />

        <div
          class="overlay pointer-events-none absolute inset-0 flex flex-col justify-end bg-gradient-to-t from-background-900 to-transparent p-3 opacity-0 group-hover:opacity-100"
        >
          <div class="mb-1 flex items-center justify-between">
            <span class="truncate text-xs text-content">
              {{ formatFileSize(fileData.size) }}
            </span>
            <span class="text-xs text-content"> {{ fileData.width }} × {{ fileData.height }} </span>
          </div>
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-1">
              <i class="fas fa-tag text-xs text-brand-700" />
              <span class="text-xs text-content">{{ fileData.ai_info?.tags?.length || 0 }}</span>
            </div>
            <div
              class="h-5 w-5 rounded-full"
              :style="{ backgroundColor: fileData.ai_info?.dominant_color || '#333' }"
              :title="fileData.ai_info?.dominant_color"
            />
          </div>
        </div>

        <div class="pointer-events-none absolute left-2 top-2 flex flex-col items-start space-y-1">
          <cyber-file-expiry-tag
            v-if="fileData.is_time_limited"
            :expires-at="fileData.expires_at"
            :storage-duration="fileData.storage_duration"
            :is-time-limited="fileData.is_time_limited"
            position="top-left"
            mode="both"
            :show-icon="true"
          />
        </div>

        <div class="pointer-events-none absolute right-2 top-2 flex flex-col items-end space-y-1">
          <div
            v-if="fileData.ai_info?.is_nsfw"
            class="quality-badge rounded bg-error-800 px-2 py-0.5 text-xs text-content-heading backdrop-blur-sm"
          >
            NSFW
          </div>

          <div
            v-if="fileData.is_recommended"
            class="quality-badge rounded bg-brand-800 px-2 py-0.5 text-xs text-content-heading backdrop-blur-sm"
          >
            <i class="fas fa-star mr-1" />{{ $t('admin.files.fileCard.badges.recommended') }}
          </div>

          <div
            v-if="fileData.is_duplicate"
            class="quality-badge rounded bg-warning-800 px-2 py-0.5 text-xs text-content-heading backdrop-blur-sm"
          >
            <i class="fas fa-copy mr-1" />{{ $t('admin.files.fileCard.badges.duplicate') }}
          </div>
        </div>
      </div>
    </div>

    <div class="p-2.5">
      <div class="mb-1.5">
        <h3 class="truncate text-sm font-medium text-content-heading" :title="fileData.display_name || fileData.original_name">
          {{ fileData.display_name || fileData.original_name }}
        </h3>
        <p class="truncate text-xs text-content-muted">
          {{ fileData.format.toUpperCase() }} | {{ $t('admin.files.fileCard.uploader') }}: {{ fileData.user_name }}
        </p>
      </div>

      <div v-if="fileData.ai_info?.tags?.length" class="tags-list mb-1.5">
        <div class="tags-container">
          <CyberTag
            v-for="(tag, index) in fileData.ai_info.tags.slice(0, 2)"
            :key="index"
            variant="primary"
            size="small"
            truncate
          >
            {{ tag }}
          </CyberTag>
          <CyberTag v-if="fileData.ai_info.tags.length > 2" variant="secondary" size="small">
            +{{ fileData.ai_info.tags.length - 2 }}
          </CyberTag>
        </div>
      </div>

      <div class="flex items-center justify-between">
        <div class="text-xs text-content-muted">
          {{ formatDate(fileData.created_at) }}
        </div>
        <div v-if="!selectMode" class="action-buttons flex gap-1.5">
          <button
            class="cyber-action-btn"
            :class="{ active: fileData.is_recommended }"
            :title="$t('admin.files.fileCard.actions.recommend')"
            @click.stop="toggleRecommend"
          >
            <i class="fas fa-star" />
            <div class="cyber-btn-glow" />
          </button>
          <button
            v-if="fileData.ai_info"
            class="cyber-action-btn info-btn"
            :title="$t('admin.files.fileCard.actions.details')"
            @click.stop="viewDetails"
          >
            <i class="fas fa-info-circle" />
            <div class="cyber-btn-glow" />
          </button>
          <button
            class="cyber-action-btn danger-btn"
            :title="$t('admin.files.fileCard.actions.delete')"
            @click.stop="deleteImage"
          >
            <i class="fas fa-trash-alt" />
            <div class="cyber-btn-glow" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .image-card {
    position: relative;
    cursor: pointer;
    height: fit-content;
    content-visibility: auto;
    contain: content;
    contain-intrinsic-size: 220px 250px;
    border-radius: var(--radius-md);
    overflow: visible;
    transition: all 0.2s ease;
  }

  .image-card::after {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    pointer-events: none;
    z-index: 99;
    opacity: 0;
    transition:
      opacity var(--transition-normal) var(--ease-out),
      box-shadow var(--transition-normal) var(--ease-out),
      border-color var(--transition-normal) var(--ease-out);
  }

  .image-card:not(.select-mode):hover::after {
    opacity: 1;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow:
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.2),
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .image-card.selected::after {
    opacity: 0;
    border: none;
    box-shadow: none;
  }

  .image-wrapper {
    cursor: pointer;
    overflow: hidden;
    backface-visibility: hidden;
    -webkit-backface-visibility: hidden;
  }

  .action-buttons {
    gap: 0.375rem;
  }

  .cyber-action-btn {
    position: relative;
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(145deg, rgba(var(--color-background-900-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.9));
    border: 1px solid rgba(var(--color-border-default-rgb), 0.5);
    border-radius: var(--radius-sm);
    color: var(--color-content-muted);
    font-size: 11px;
    cursor: pointer;
    transition: all var(--transition-normal) var(--ease-out);
    overflow: hidden;
    backdrop-filter: var(--backdrop-blur-sm);
  }

  .cyber-btn-glow {
    position: absolute;
    top: -1px;
    left: -1px;
    right: -1px;
    bottom: -1px;
    border-radius: var(--radius-md);
    opacity: 0;
    transition: opacity var(--transition-normal) var(--ease-out);
    pointer-events: none;
  }

  .cyber-action-btn:hover {
    color: var(--color-content-heading);
    border-color: var(--color-border-default);
    background: var(--color-hover-bg);
    box-shadow: var(--shadow-sm);
  }

  .cyber-action-btn:hover .cyber-btn-glow {
    opacity: 0.5;
  }

  .cyber-action-btn:active {
    transform: scale(0.95);
    box-shadow: none;
  }

  .cyber-action-btn.active {
    color: var(--color-warning-500);
    border-color: rgba(var(--color-warning-rgb), 0.4);
    background: rgba(var(--color-warning-rgb), 0.1);
  }

  .cyber-action-btn.active:hover {
    color: var(--color-warning-400);
    border-color: rgba(var(--color-warning-rgb), 0.6);
    background: rgba(var(--color-warning-rgb), 0.15);
  }

  .cyber-action-btn.info-btn:hover {
    color: var(--color-info-500);
    border-color: rgba(var(--color-info-rgb), 0.5);
    background: rgba(var(--color-info-rgb), 0.1);
  }

  .cyber-action-btn.danger-btn:hover {
    color: var(--color-error-500);
    border-color: rgba(var(--color-error-rgb), 0.5);
    background: rgba(var(--color-error-rgb), 0.1);
  }

  .quality-badge {
    font-size: var(--text-xs);
    padding: var(--space-xs) var(--space-sm);
    border-radius: var(--radius-sm);
  }

  .tags-list {
    height: auto;
    overflow: hidden;
    display: flex;
    align-items: flex-start;
  }

  .tags-container {
    display: flex;
    flex-wrap: nowrap;
    gap: var(--space-xs);
    align-items: center;
    width: 100%;
    overflow: hidden;
  }

  .p-2\.5 {
    padding: 0.625rem;
  }
  .mb-1\.5 {
    margin-bottom: 0.375rem;
  }
  .text-sm {
    font-size: var(--text-sm);
  }
  .space-x-1 {
    margin-left: var(--space-xs);
    margin-right: var(--space-xs);
  }
  .rounded-lg {
    border-radius: var(--radius-md);
  }

  .is-nsfw .image-wrapper img {
    filter: blur(10px);
  }

  .overlay {
    z-index: 5;

    backface-visibility: hidden;
    -webkit-backface-visibility: hidden;
    will-change: opacity;

    transform: translateZ(0);
    -webkit-transform: translateZ(0);
  }

  .nsfw-overlay {
    background: rgba(var(--color-background-900-rgb), 0.95);
    border: 1px solid rgba(var(--color-error-rgb), 0.2);
  }

  .hover\:shadow-glow:hover {
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.12);
  }

  .image-card.select-mode {
    cursor: pointer;
  }

  .image-card.select-mode:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
  }

  .image-card.select-mode .image-wrapper img {
    opacity: 1;
  }

  .image-card.select-mode.selected .image-wrapper img {
    filter: brightness(0.8);
  }

  .image-card,
  .image-card *,
  .image-card:hover,
  .image-card:hover * {
    transform: none !important;
    animation: none !important;
  }

  .image-card .overlay {
    backface-visibility: visible;
    -webkit-backface-visibility: visible;
    will-change: auto;
    transform: none;
    -webkit-transform: none;
  }
</style>
<style src="@/styles/selection.css"></style>
