<script setup lang="ts">
  import type { FileInfo, FileTag } from '../types'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    isVisible: boolean
    currentFile: FileInfo | null
    isFullscreen: boolean
    hasAIInfo?: boolean
    isLightBackground?: boolean
  }

  withDefaults(defineProps<Props>(), {
    hasAIInfo: false,
    isLightBackground: false,
  })

  defineEmits<{
    'toggle-fullscreen': []
    download: []
    close: []
    'tag-click': [tag: FileTag]
    'find-similar': []
  }>()

  const { $t } = useTexts()
</script>

<template>
  <div class="navigation-bar" :class="{ 'is-hidden': !isVisible, 'is-light-bg': isLightBackground }">
    <div class="nav-title">
      <span class="image-name" :class="{ 'is-light-bg': isLightBackground }">{{
        currentFile?.display_name || $t('components.navigationBar.filePreview')
      }}</span>
      <div class="image-info">
        <span v-if="currentFile?.resolution" class="resolution-tag" :class="{ 'is-light-bg': isLightBackground }">{{
          currentFile.resolution
        }}</span>
        <span class="dimension-tag" :class="{ 'is-light-bg': isLightBackground }"
          >{{ currentFile?.width || 0 }} × {{ currentFile?.height || 0 }}</span
        >
        <span v-if="currentFile?.size_formatted" class="size-tag" :class="{ 'is-light-bg': isLightBackground }">{{
          currentFile.size_formatted
        }}</span>
        <span v-if="currentFile?.format" class="format-tag" :class="{ 'is-light-bg': isLightBackground }">{{
          currentFile.format.toUpperCase()
        }}</span>
      </div>

      <div v-if="currentFile?.tags && currentFile.tags.length > 0" class="image-tags-header">
        <div
          v-for="tag in currentFile.tags"
          :key="tag.id"
          class="tag"
          :class="{ 'is-light-bg': isLightBackground }"
          @click.stop="$emit('tag-click', tag)"
        >
          {{ tag.name }}
        </div>
      </div>
    </div>

    <div class="nav-actions">
      <button
        v-if="hasAIInfo"
        class="action-btn tooltip-btn similar-btn-icon"
        :title="$t('components.navigationBar.findSimilar')"
        @click="$emit('find-similar')"
      >
        <i class="fas fa-layer-group" />
      </button>
      <button
        class="action-btn tooltip-btn"
        :title="$t('components.navigationBar.fullscreenToggle')"
        @click="$emit('toggle-fullscreen')"
      >
        <i :class="isFullscreen ? 'fas fa-compress' : 'fas fa-expand'" />
      </button>
      <button class="action-btn tooltip-btn" :title="$t('components.navigationBar.download')" @click="$emit('download')">
        <i class="fas fa-download" />
      </button>
      <button class="action-btn tooltip-btn" :title="$t('components.navigationBar.close')" @click="$emit('close')">
        <i class="fas fa-times" />
      </button>
    </div>
  </div>
</template>

<style scoped>
  .navigation-bar {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: transparent;
    padding: 12px 24px;
    color: rgba(var(--color-white-rgb), 0.9);
    transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    z-index: 100;
  }

  .navigation-bar.is-hidden {
    opacity: 0;
    transform: translateY(-20px);
    pointer-events: none;
  }

  .navigation-bar.is-light-bg {
  }

  .nav-title {
    flex: 1;
    min-width: 0;
  }

  .image-name {
    font-size: 15px;
    font-weight: 700;
    color: var(--color-brand-400);
    display: block;
    margin-bottom: 6px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    text-shadow: 0 0 calc(20px * var(--text-shadow-intensity, 1)) rgba(var(--color-info-rgb), 0.8);
    letter-spacing: 0.5px;
  }

  .image-name.is-light-bg {
    color: var(--color-background-900);
    font-weight: 900;
    font-size: 16px;
    text-shadow:
      0 0 3px rgba(var(--color-content-rgb), 0.9),
      0 0 6px rgba(var(--color-content-rgb), 0.8),
      0 0 9px rgba(var(--color-content-rgb), 0.7),
      1px 1px 2px rgba(var(--color-content-rgb), 0.6),
      -1px -1px 2px rgba(var(--color-content-rgb), 0.6);
  }

  .image-info {
    display: flex;
    align-items: center;
    gap: 5px;
    margin-bottom: 4px;
    flex-wrap: wrap;
  }

  .resolution-tag,
  .dimension-tag,
  .size-tag,
  .format-tag {
    background: rgba(var(--color-brand-500-rgb), 0.14);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.35);
    border-radius: var(--radius-sm);
    padding: 1px 6px;
    font-size: 11px;
    color: rgba(var(--color-white-rgb), 0.9);
    font-weight: 500;
  }

  .resolution-tag.is-light-bg,
  .dimension-tag.is-light-bg,
  .size-tag.is-light-bg,
  .format-tag.is-light-bg {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: rgba(var(--color-content-rgb), 0.95);
    backdrop-filter: blur(25px) saturate(110%);
    box-shadow:
      0 12px 32px rgba(var(--color-background-900-rgb), 0.6),
      0 6px 16px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 0 1px rgba(var(--color-content-rgb), 0.08),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.15);
    text-shadow:
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 6px rgba(var(--color-content-rgb), 0.3),
      2px 2px 4px rgba(var(--color-background-900-rgb), 0.8);
  }

  .image-tags-header {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
  }

  .tag {
    background: rgba(var(--color-error-rgb), 0.18);
    border: 1px solid rgba(var(--color-error-rgb), 0.35);
    border-radius: var(--radius-sm);
    padding: 2px 8px;
    font-size: 11px;
    color: rgba(var(--color-white-rgb), 0.9);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .tag:hover {
    background: rgba(var(--color-error-rgb), 0.28);
    border-color: rgba(var(--color-error-rgb), 0.55);
    transform: translateY(-1px);
  }

  .tag.is-light-bg {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: rgba(var(--color-content-rgb), 0.95);
    backdrop-filter: blur(25px) saturate(110%);
    box-shadow:
      0 12px 32px rgba(var(--color-background-900-rgb), 0.6),
      0 6px 16px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 0 1px rgba(var(--color-content-rgb), 0.08),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.15);
    text-shadow:
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 6px rgba(var(--color-content-rgb), 0.3),
      2px 2px 4px rgba(var(--color-background-900-rgb), 0.8);
  }

  .tag.is-light-bg:hover {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.8) 0%,
      rgba(var(--color-background-900-rgb), 0.75) 50%,
      rgba(var(--color-background-900-rgb), 0.85) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    transform: translateY(-1px);
    box-shadow:
      0 16px 40px rgba(var(--color-background-900-rgb), 0.7),
      0 8px 20px rgba(var(--color-background-900-rgb), 0.5),
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 0 1px rgba(var(--color-content-rgb), 0.12),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.25);
  }

  .nav-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .action-btn {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-content-heading);
    width: 36px;
    height: 36px;
    border-radius: var(--radius-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all var(--transition-normal) var(--ease-in-out);
    font-size: var(--text-base);
    backdrop-filter: var(--backdrop-blur-md);
    position: relative;
    overflow: hidden;
    box-shadow:
      0 2px 8px rgba(var(--color-background-900-rgb), 0.3),
      inset 0 1px 0 rgba(var(--color-content-rgb), 0.08);
  }

  .action-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.35);
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    color: var(--color-brand-400);
    transform: translateY(-2px);
    box-shadow:
      0 4px 15px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 16px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 1px 0 rgba(var(--color-content-rgb), 0.12);
  }

  .action-btn:active {
    transform: translateY(0);
    box-shadow:
      0 2px 8px rgba(var(--color-info-rgb), calc(0.2 + var(--viewer-is-light, 0) * 0.2)),
      inset 0 1px 0 rgba(var(--color-content-rgb), 0.1);
  }

  .similar-btn-icon {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25) 0%, rgba(var(--color-brand-500-rgb), 0.15) 100%);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .similar-btn-icon:hover {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.35) 0%, rgba(var(--color-brand-500-rgb), 0.25) 100%);
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    color: var(--color-brand-400);
  }

  .tooltip-btn {
    position: relative;
  }

  .navigation-bar.is-light-bg .action-btn {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(25px) saturate(110%);
    box-shadow:
      0 12px 32px rgba(var(--color-background-900-rgb), 0.6),
      0 6px 16px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.08),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.15);
    text-shadow:
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 6px rgba(255, 255, 255, 0.3),
      2px 2px 4px rgba(var(--color-background-900-rgb), 0.8);
  }

  .navigation-bar.is-light-bg .action-btn:hover {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.8) 0%,
      rgba(var(--color-background-900-rgb), 0.75) 50%,
      rgba(var(--color-background-900-rgb), 0.85) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    color: rgba(255, 255, 255, 0.98);
    transform: translateY(-2px);
    box-shadow:
      0 16px 40px rgba(var(--color-background-900-rgb), 0.7),
      0 8px 20px rgba(var(--color-background-900-rgb), 0.5),
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 0 1px rgba(255, 255, 255, 0.12),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.25);
  }

  .navigation-bar.is-light-bg .similar-btn-icon {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
  }

  .navigation-bar.is-light-bg .similar-btn-icon:hover {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.8) 0%,
      rgba(var(--color-background-900-rgb), 0.75) 50%,
      rgba(var(--color-background-900-rgb), 0.85) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    color: rgba(255, 255, 255, 0.98);
  }

  @media (max-width: 768px) {
    .navigation-bar {
      padding: 10px 16px;
    }

    .image-name {
      font-size: 14px;
    }

    .image-info {
      gap: 4px;
    }

    .resolution-tag,
    .dimension-tag,
    .size-tag,
    .format-tag {
      font-size: 10px;
      padding: 1px 5px;
    }

    .action-btn {
      width: 32px;
      height: 32px;
      font-size: 13px;
    }

    .tag {
      font-size: 9px;
      padding: 1px 5px;
    }
  }

  @media (max-width: 480px) {
    .navigation-bar {
      padding: 6px 12px;
    }

    .image-name {
      font-size: 12px;
      margin-bottom: 4px;
    }

    .image-info {
      gap: 3px;
    }

    .nav-actions {
      gap: 4px;
    }

    .action-btn {
      width: 28px;
      height: 28px;
      font-size: 11px;
    }
  }
</style>
