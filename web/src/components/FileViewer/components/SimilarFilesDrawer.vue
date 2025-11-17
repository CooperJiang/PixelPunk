<script setup lang="ts">
  import { computed, toRefs } from 'vue'
  import type { FileInfo } from '@/api/admin/files'
  import { useTexts } from '@/composables/useTexts'

  const props = defineProps<{
    modelValue: boolean
    loading: boolean
    error: Error | null
    similarFiles: FileInfo[]
  }>()

  const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'retry'): void
    (e: 'select', file: FileInfo): void
    (e: 'download', file: FileInfo): void
    (e: 'view', file: FileInfo): void
  }>()

  const { $t } = useTexts()
  const { modelValue, loading, error, similarFiles } = toRefs(props)

  const visible = computed({
    get: () => modelValue.value,
    set: (value: boolean) => emit('update:modelValue', value),
  })

  const maxSimilarity = computed(() => {
    if (!similarFiles.value.length) {
      return 0
    }
    return Math.max(...similarFiles.value.map((file) => file.similarity || 0))
  })

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) {
      return `${bytes} B`
    }
    if (bytes < 1024 * 1024) {
      return `${(bytes / 1024).toFixed(1)} KB`
    }
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  const handleFileError = (event: Event) => {
    const img = event.target as HTMLImageElement
    img.src = '/placeholder-image.png'
  }
</script>

<template>
  <CyberDrawer v-model="visible" :title="$t('components.similarFilesDrawer.title')" width="420px" :mask-closable="true">
    <template #header="{ close }">
      <div class="custom-header">
        <div class="header-bg-pattern" />

        <div class="header-content">
          <div class="header-left">
            <div class="header-icon-wrapper">
              <i class="fas fa-layer-group" />
              <div class="icon-glow" />
            </div>
            <div class="header-text">
              <h2 class="header-title">{{ $t('components.similarFilesDrawer.title') }}</h2>
              <p class="header-subtitle">{{ $t('components.similarFilesDrawer.subtitle') }}</p>
            </div>
          </div>

          <div class="header-right">
            <button class="custom-close-btn" :title="$t('components.similarFilesDrawer.close')" @click="close">
              <i class="fas fa-times" />
            </button>
          </div>
        </div>

        <div v-if="!loading && !error && similarFiles.length > 0" class="header-stats">
          <div class="stat-item">
            <span class="stat-label">{{ $t('components.similarFilesDrawer.stats.matchResults') }}</span>
            <span class="stat-value">{{ similarFiles.length }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">{{ $t('components.similarFilesDrawer.stats.maxSimilarity') }}</span>
            <span class="stat-value">{{ Math.round(maxSimilarity * 100) }}%</span>
          </div>
        </div>
      </div>
    </template>

    <div class="similar-images-container">
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner">
          <i class="fas fa-circle-notch fa-spin" />
        </div>
        <p>{{ $t('components.similarFilesDrawer.loading') }}</p>
      </div>

      <div v-else-if="error" class="error-state">
        <i class="fas fa-exclamation-triangle" />
        <p>{{ error.message }}</p>
        <button class="retry-btn" @click="$emit('retry')">
          <i class="fas fa-redo" />
          {{ $t('components.similarFilesDrawer.retry') }}
        </button>
      </div>

      <div v-else-if="!similarFiles.length" class="empty-state">
        <i class="fas fa-search" />
        <p>{{ $t('components.similarFilesDrawer.notFound') }}</p>
      </div>

      <div v-else class="files-list">
        <div v-for="file in similarFiles" :key="file.id" class="similar-file-item" @click="$emit('select', file)">
          <div class="similarity-badge">
            <i class="fas fa-percentage" />
            {{ Math.round(file.similarity * 100) }}
          </div>

          <div class="file-container">
            <img
              :src="file.full_thumb_url || file.thumb_url"
              :alt="file.display_name || file.original_name"
              loading="lazy"
              @error="handleFileError"
            />
          </div>

          <div class="file-info">
            <p class="file-name">{{ file.display_name || file.original_name }}</p>
            <div class="file-meta">
              <span class="meta-item">
                <i class="fas fa-expand-arrows-alt" />
                {{ file.width }} × {{ file.height }}
              </span>
              <span class="meta-item">
                <i class="fas fa-file" />
                {{ formatFileSize(file.size) }}
              </span>
            </div>

            <div v-if="file.ai_info?.tags?.length" class="file-tags">
              <span v-for="(tag, idx) in file.ai_info.tags.slice(0, 3)" :key="idx" class="tag-chip">
                {{ tag }}
              </span>
              <span v-if="file.ai_info.tags.length > 3" class="tag-more"> +{{ file.ai_info.tags.length - 3 }} </span>
            </div>
          </div>

          <div class="hover-actions">
            <button
              class="action-btn"
              :title="$t('components.similarFilesDrawer.actions.download')"
              @click.stop="$emit('download', file)"
            >
              <i class="fas fa-download" />
            </button>
            <button
              class="action-btn"
              :title="$t('components.similarFilesDrawer.actions.view')"
              @click.stop="$emit('view', file)"
            >
              <i class="fas fa-eye" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </CyberDrawer>
</template>

<style scoped>
  .custom-header {
    position: relative;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.98) 0%,
      rgba(var(--color-background-900-rgb), 0.98) 100%
    );
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    overflow: hidden;
  }

  .header-bg-pattern {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    opacity: 0.1;
    background-image:
      repeating-linear-gradient(
        45deg,
        transparent,
        transparent 10px,
        rgba(var(--color-brand-500-rgb), 0.1) 10px,
        rgba(var(--color-brand-500-rgb), 0.1) 20px
      ),
      repeating-linear-gradient(
        -45deg,
        transparent,
        transparent 10px,
        rgba(var(--color-error-rgb), 0.05) 10px,
        rgba(var(--color-error-rgb), 0.05) 20px
      );
    pointer-events: none;
  }

  .header-content {
    position: relative;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1.5rem;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .header-icon-wrapper {
    position: relative;
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2) 0%, rgba(var(--color-brand-500-rgb), 0.1) 100%);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    border-radius: var(--radius-sm);
    font-size: 16px;
    color: var(--color-brand-400);
  }

  .icon-glow {
    position: absolute;
    inset: -2px;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.4) 0%, transparent 70%);
    border-radius: var(--radius-sm);
    animation: pulse 2s infinite;
    pointer-events: none;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 0.5;
      transform: scale(1);
    }
    50% {
      opacity: 0.8;
      transform: scale(1.05);
    }
  }

  .header-text {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .header-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-white);
    margin: 0;
    text-shadow: 0 2px 4px rgba(var(--color-background-900-rgb), 0.3);
  }

  .header-subtitle {
    font-size: 0.7rem;
    color: rgba(var(--color-brand-500-rgb), 0.8);
    margin: 0;
    letter-spacing: 0.05em;
    text-transform: uppercase;
  }

  .custom-close-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-hover-bg-neutral);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    color: rgba(var(--color-content-rgb), 0.6);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .custom-close-btn:hover {
    background: rgba(var(--color-error-rgb), 0.2);
    border-color: rgba(var(--color-error-rgb), 0.4);
    color: var(--color-error-500);
    transform: scale(1.05);
  }

  .header-stats {
    display: flex;
    gap: 1.5rem;
    padding: 0.5rem 1.5rem;
    background: rgba(var(--color-background-900-rgb), 0.3);
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .stat-label {
    font-size: 0.7rem;
    color: var(--color-content-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .stat-value {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-brand-400);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .similar-images-container {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 0;
  }

  .loading-state,
  .error-state,
  .empty-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    color: rgba(var(--color-content-rgb), 0.6);
  }

  .loading-spinner {
    font-size: 2rem;
    color: var(--color-brand-400);
    animation: pulse 1.5s ease-in-out infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 0.6;
    }
    50% {
      opacity: 1;
    }
  }

  .error-state > i:first-child,
  .empty-state > i:first-child {
    font-size: 2.5rem;
    color: rgba(var(--color-content-rgb), 0.3);
  }

  .retry-btn {
    padding: 8px 20px;
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: var(--color-brand-400);
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .retry-btn i {
    font-size: 0.875rem;
  }

  .retry-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateY(-2px);
  }

  .files-list {
    flex: 1;
    overflow-y: auto;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    scroll-behavior: smooth;
    min-height: 0;
    max-height: calc(100vh - 120px);
  }

  .similar-file-item {
    position: relative;
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s ease;
    backdrop-filter: blur(10px);
    flex-shrink: 0;
  }

  .similar-file-item:hover {
    background: rgba(var(--color-background-800-rgb), 0.8);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 8px 24px rgba(var(--color-brand-500-rgb), 0.25);
  }

  .similarity-badge {
    position: absolute;
    top: 8px;
    right: 8px;
    background: linear-gradient(135deg, var(--color-brand-400), var(--color-brand-500));
    color: var(--color-text-on-brand);
    padding: 4px 10px;
    border-radius: var(--radius-sm);
    font-size: 12px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 4px;
    z-index: 2;
    box-shadow: 0 2px 8px rgba(var(--color-background-900-rgb), 0.3);
  }

  .similarity-badge i {
    font-size: 10px;
  }

  .file-container {
    width: 100%;
    overflow: hidden;
    background:
      linear-gradient(
        to bottom,
        rgba(var(--color-error-rgb), 0.05) 0%,
        transparent 50%,
        rgba(var(--color-brand-500-rgb), 0.05) 100%
      ),
      repeating-linear-gradient(
        to right,
        transparent 0,
        transparent 10px,
        rgba(var(--color-error-rgb), 0.1) 10px,
        rgba(var(--color-error-rgb), 0.1) 12px
      ),
      repeating-linear-gradient(
        to bottom,
        transparent 0,
        transparent 10px,
        rgba(var(--color-brand-500-rgb), 0.1) 10px,
        rgba(var(--color-brand-500-rgb), 0.1) 12px
      );
    background-size:
      100% 100%,
      20px 100%,
      100% 20px;
    background-position: center;
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    position: relative;
  }

  .file-container::after {
    content: '';
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    width: 1px;
    background: linear-gradient(to bottom, transparent, rgba(var(--color-error-rgb), 0.5), transparent);
    box-shadow: 0 0 10px rgba(var(--color-error-rgb), 0.5);
    animation: cyber-scan-vert 4s linear infinite;
  }

  .file-container img {
    width: 100%;
    height: 180px;
    object-fit: contain;
    transition: transform 0.3s ease;
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    display: block;
  }

  .similar-file-item:hover .file-container img {
    transform: scale(1.1);
  }

  .file-info {
    padding: 0.75rem;
  }

  .file-name {
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--color-content-heading);
    margin: 0 0 0.5rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .file-meta {
    display: flex;
    gap: 1rem;
    margin-bottom: 0.5rem;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.75rem;
    color: rgba(var(--color-content-rgb), 0.6);
  }

  .meta-item i {
    font-size: 0.7rem;
    color: rgba(var(--color-brand-500-rgb), 0.6);
  }

  .file-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-top: 0.5rem;
  }

  .tag-chip {
    padding: 2px 8px;
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-400);
    border-radius: var(--radius-sm);
    font-size: 0.7rem;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .tag-more {
    padding: 2px 8px;
    background: var(--color-hover-bg-neutral);
    color: var(--color-content-muted);
    border-radius: var(--radius-sm);
    font-size: 0.7rem;
  }

  .hover-actions {
    position: absolute;
    bottom: 0.75rem;
    right: 0.75rem;
    display: flex;
    gap: 0.5rem;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .similar-file-item:hover .hover-actions {
    opacity: 1;
  }

  .action-btn {
    width: 32px;
    height: 32px;
    background: rgba(var(--color-background-900-rgb), 0.7);
    backdrop-filter: blur(10px);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-full);
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .action-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    transform: scale(1.1);
  }

  .action-btn i {
    font-size: 14px;
  }

  @keyframes cyber-scan-vert {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(200%);
    }
  }
</style>
