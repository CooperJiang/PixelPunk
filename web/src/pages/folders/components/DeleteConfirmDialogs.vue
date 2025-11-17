<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'DeleteConfirmDialogs',
  })

  const { $t } = useTexts()

  interface Props {
    deleteDialogVisible: boolean
    deleteType: 'folder' | 'image' | null
    deleteTarget: any | null
    isDeleting?: boolean

    batchDeleteDialogVisible: boolean
    selectedBatchImagesCount: number
    isBatchDeleting?: boolean
  }

  defineProps<Props>()

  interface Emits {
    (e: 'update:delete-dialog-visible', value: boolean): void
    (e: 'update:batch-delete-dialog-visible', value: boolean): void
    (e: 'confirm-delete'): void
    (e: 'confirm-batch-delete'): void
    (e: 'cancel-delete'): void
    (e: 'cancel-batch-delete'): void
  }

  const emit = defineEmits<Emits>()

  /* Event handlers */
  const handleDeleteDialogUpdate = (value: boolean) => {
    emit('update:delete-dialog-visible', value)
  }

  const handleBatchDeleteDialogUpdate = (value: boolean) => {
    emit('update:batch-delete-dialog-visible', value)
  }

  const handleConfirmDelete = () => {
    emit('confirm-delete')
  }

  const handleConfirmBatchDelete = () => {
    emit('confirm-batch-delete')
  }

  const handleCancelDelete = () => {
    emit('cancel-delete')
  }

  const handleCancelBatchDelete = () => {
    emit('cancel-batch-delete')
  }
</script>

<template>
  <CyberDialog
    :model-value="deleteDialogVisible"
    :title="
      deleteType === 'folder' ? $t('folders.deleteConfirmDialogs.title.folder') : $t('folders.deleteConfirmDialogs.title.file')
    "
    width="420px"
    max-width="90vw"
    :show-default-footer="true"
    :confirm-text="$t('folders.deleteConfirmDialogs.confirmText')"
    :cancel-text="$t('actions.cancel')"
    :loading="isDeleting"
    @update:model-value="handleDeleteDialogUpdate"
    @confirm="handleConfirmDelete"
    @cancel="handleCancelDelete"
  >
    <div class="delete-confirm-content">
      <div class="warning-icon">
        <i v-if="deleteType === 'folder'" class="fas fa-folder-minus" />
        <i v-else class="fas fa-image" />
      </div>
      <div class="warning-text">
        <p v-if="deleteType === 'folder'" class="delete-message">
          {{ $t('folders.deleteConfirmDialogs.message.folder', { name: deleteTarget?.name }) }}
        </p>
        <p v-else class="delete-message">
          {{ $t('folders.deleteConfirmDialogs.message.file', { name: deleteTarget?.display_name }) }}
        </p>
        <p class="warning-note">
          <i class="fas fa-exclamation-triangle warning-note-icon"></i>
          {{ $t('folders.deleteConfirmDialogs.warning') }}
        </p>
      </div>
    </div>
  </CyberDialog>

  <CyberDialog
    :model-value="batchDeleteDialogVisible"
    :title="$t('folders.deleteConfirmDialogs.title.batchFile')"
    width="420px"
    max-width="90vw"
    :show-default-footer="true"
    :confirm-text="$t('folders.deleteConfirmDialogs.confirmText')"
    :cancel-text="$t('actions.cancel')"
    :loading="isBatchDeleting"
    @update:model-value="handleBatchDeleteDialogUpdate"
    @confirm="handleConfirmBatchDelete"
    @cancel="handleCancelBatchDelete"
  >
    <div class="delete-confirm-content">
      <div class="warning-icon">
        <i class="fas fa-images" />
      </div>
      <div class="warning-text">
        <p class="delete-message">
          {{ $t('folders.deleteConfirmDialogs.message.batchFile', { count: selectedBatchImagesCount }) }}
        </p>
        <p class="warning-note">
          <i class="fas fa-exclamation-triangle warning-note-icon"></i>
          {{ $t('folders.deleteConfirmDialogs.warning') }}
        </p>
      </div>
    </div>
  </CyberDialog>
</template>

<style scoped>
  .delete-confirm-content {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    padding: 0.5rem 0.75rem;
    background: rgba(var(--color-brand-500-rgb), 0.04);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    position: relative;
    overflow: hidden;
  }

  .warning-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.5rem;
    height: 2.5rem;
    flex-shrink: 0;
  }

  .warning-icon i {
    font-size: 1.25rem;
    color: var(--color-status-warning);
  }

  .warning-text {
    flex: 1;
    min-width: 0;
  }

  .delete-message {
    margin: 0 0 0.75rem 0;
    color: var(--color-content-default);
    line-height: 1.5;
    font-size: 1rem;
    font-weight: 500;
  }

  .target-name {
    color: var(--color-error-500);
    font-weight: 600;
    text-shadow: 0 0 4px rgba(var(--color-error-rgb), 0.3);
    background: rgba(var(--color-error-rgb), 0.15);
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    display: inline-block;
    max-width: 100%;
    word-break: break-all;
    line-height: 1.4;
  }

  .warning-note {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.825rem;
    color: var(--color-status-warning);
    font-weight: 500;
    margin: 0;
    padding: 0.5rem;
    background: rgba(var(--color-status-warning-rgb), 0.08);
    border: 1px solid rgba(var(--color-status-warning-rgb), 0.2);
    border-radius: var(--radius-sm);
  }

  .warning-note-icon {
    color: var(--color-status-warning);
    font-size: 0.75rem;
    flex-shrink: 0;
  }

  @media (max-width: 480px) {
    .delete-confirm-content {
      gap: 0.75rem;
      padding: 0.5rem;
    }

    .warning-icon {
      width: 2rem;
      height: 2rem;
    }

    .warning-icon i {
      font-size: 1rem;
    }

    .delete-message {
      font-size: 0.9rem;
      margin-bottom: 0.5rem;
    }

    .target-name {
      padding: 0.1rem 0.25rem;
    }

    .warning-note {
      font-size: 0.75rem;
      padding: 0.375rem;
    }

    .warning-note-icon {
      font-size: 0.7rem;
    }
  }
</style>
