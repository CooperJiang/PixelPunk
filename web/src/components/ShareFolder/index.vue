<script setup lang="ts">
  import { computed } from 'vue'
  import { formatDate as formatDateUtil } from '@/utils/formatting/format'
  import type { ShareFolderProps, ShareFolderEmits } from './types'
  import { useTexts } from '@/composables/useTexts'

  const props = defineProps<ShareFolderProps>()
  const emit = defineEmits<ShareFolderEmits>()
  const { $t } = useTexts()

  const formatDate = (dateString?: string) => {
    if (!dateString) {
      return ''
    }
    return formatDateUtil(dateString)
  }

  const getItemCount = computed(() => {
    const folderCount = props.folder.folder_count || 0
    const fileCount = props.folder.file_count || 0

    if (
      folderCount === 0 &&
      fileCount === 0 &&
      props.folder.folder_count === undefined &&
      props.folder.file_count === undefined
    ) {
      return ''
    }

    const parts = []
    if (folderCount > 0) {
      parts.push(`${folderCount} ${$t('shareFolder.folders')}`)
    }
    if (fileCount > 0) {
      parts.push(`${fileCount} ${$t('shareFolder.images')}`)
    }

    return parts.length > 0 ? parts.join('，') : $t('shareFolder.emptyFolder')
  })
</script>

<template>
  <div class="share-folder-card" @click="emit('click', folder)">
    <div class="folder-icon">
      <div class="icon-inner">
        <i class="fas fa-folder" />
      </div>
    </div>
    <div class="folder-info">
      <div class="folder-name" :title="folder.name">{{ folder.name }}</div>
      <div class="folder-meta">
        <span v-if="getItemCount" class="folder-count">{{ getItemCount }}</span>
        <span v-if="folder.created_at" class="folder-date">
          <i class="fas fa-calendar-alt" />
          {{ formatDate(folder.created_at) }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .share-folder-card {
    @apply relative mb-3 flex cursor-pointer items-center gap-3 overflow-hidden rounded-lg p-3;
    background: var(--color-background-700);
    transition: all 0.3s ease;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }

  .share-folder-card:hover {
    @apply -translate-y-0.5;
    background: var(--color-background-600);
    box-shadow:
      0 4px 12px rgba(0, 0, 0, 0.4),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .share-folder-card:active {
    @apply translate-y-0 scale-[0.98];
  }

  .folder-icon {
    @apply flex h-12 w-12 flex-shrink-0 items-center justify-center;
  }

  .icon-inner {
    @apply flex h-full w-full items-center justify-center rounded-lg text-lg;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.08));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
    transition: all 0.3s ease;
  }

  .share-folder-card:hover .icon-inner {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.3), rgba(var(--color-brand-500-rgb), 0.12));
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .folder-info {
    @apply flex-1 overflow-hidden;
  }

  .folder-name {
    @apply mb-1.5 overflow-hidden text-ellipsis whitespace-nowrap text-base font-semibold;
    color: rgba(var(--color-content-default-rgb), 0.95);
    letter-spacing: 0.01em;
  }

  .folder-meta {
    @apply flex items-center gap-2.5;
    min-width: 0;
  }

  .folder-date {
    @apply flex items-center gap-1.5 rounded-md px-2 py-1 text-xs;
    color: rgba(var(--color-content-default-rgb), 0.7);
    background: var(--color-background-800);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
    flex-shrink: 1;
  }

  .folder-date i {
    @apply flex-shrink-0 text-xs;
    color: var(--color-brand-500);
  }

  .folder-count {
    @apply flex-shrink-0 text-xs font-medium;
    color: rgba(var(--color-content-default-rgb), 0.7);
    white-space: nowrap;
  }

  @media (max-width: 480px) {
    .share-folder-card {
      padding: 0.5rem;
      gap: 0.5rem;
      margin-bottom: 0.35rem;
    }

    .folder-icon {
      width: 36px;
      height: 36px;
    }

    .icon-inner {
      width: 30px;
      height: 30px;
      font-size: 0.8rem;
    }

    .folder-name {
      font-size: 0.85rem;
      margin-bottom: 0.15rem;
    }

    .folder-date,
    .folder-count {
      font-size: 0.65rem;
      padding: 0.1rem 0.3rem;
    }
  }
</style>
