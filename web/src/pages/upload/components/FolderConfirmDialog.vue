<script setup lang="ts">
import { computed } from 'vue'
import type { FolderStats } from '@/utils/file/folderReader'
import { useTexts } from '@/composables/useTexts'

const { $t } = useTexts()

const props = defineProps<{
  visible: boolean
  stats: FolderStats
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value),
})

const handleConfirm = () => {
  emit('confirm')
  dialogVisible.value = false
}

const handleCancel = () => {
  emit('cancel')
  dialogVisible.value = false
}
</script>

<template>
  <CyberDialog v-model="dialogVisible" :title="$t('upload.folder.dialog.title')" width="480" @close="handleCancel">
    <div class="folder-confirm-content">
      <div class="folder-stats">
        <div class="stat-item">
          <div class="stat-icon folder-icon">
            <i class="fas fa-folder" />
          </div>
          <div class="stat-text">
            <span class="stat-label">{{ $t('upload.folder.dialog.totalFiles') }}</span>
            <strong class="stat-value">{{ stats.totalFiles }}</strong>
          </div>
        </div>

        <div class="stat-item success-item">
          <div class="stat-icon success-icon">
            <i class="fas fa-check-circle" />
          </div>
          <div class="stat-text">
            <span class="stat-label">{{ $t('upload.folder.dialog.validFiles') }}</span>
            <strong class="stat-value success-value">{{ stats.validFiles }}</strong>
          </div>
        </div>

        <div v-if="stats.invalidFiles > 0" class="stat-item error-item">
          <div class="stat-icon error-icon">
            <i class="fas fa-times-circle" />
          </div>
          <div class="stat-text">
            <span class="stat-label">{{ $t('upload.folder.dialog.invalidFiles') }}</span>
            <strong class="stat-value error-value">{{ stats.invalidFiles }}</strong>
          </div>
        </div>
      </div>

      <div class="folder-message">
        <i class="fas fa-info-circle message-icon" />
        <p class="message-text" v-html="$t('upload.folder.dialog.confirmMessage', { count: stats.validFiles })"></p>
      </div>

      <div v-if="stats.invalidFiles > 0" class="folder-warning">
        <i class="fas fa-exclamation-triangle warning-icon" />
        <span class="warning-text">{{ $t('upload.folder.dialog.warningMessage', { count: stats.invalidFiles }) }}</span>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <CyberButton type="secondary" @click="handleCancel">{{ $t('upload.folder.dialog.cancel') }}</CyberButton>
        <CyberButton type="primary" :disabled="stats.validFiles === 0" @click="handleConfirm">{{
          $t('upload.folder.dialog.confirm')
        }}</CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped lang="scss">
.folder-confirm-content {
  padding: var(--space-sm);
}

.folder-stats {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  margin-bottom: var(--space-md);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  background: rgba(var(--color-background-700-rgb), 0.5);
  border: 1px solid rgba(var(--color-border-default-rgb), 0.3);
  border-radius: var(--radius-md);
  transition: all var(--transition-normal) var(--ease-out);
}

.stat-item:hover {
  background: rgba(var(--color-background-700-rgb), 0.7);
  border-color: rgba(var(--color-brand-500-rgb), 0.3);
}

.stat-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  font-size: var(--text-lg);
  flex-shrink: 0;
}

.folder-icon {
  background: rgba(var(--color-brand-500-rgb), 0.15);
  color: var(--color-brand-500);
  border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
}

.success-icon {
  background: rgba(var(--color-success-rgb), 0.15);
  color: var(--color-success-500);
  border: 1px solid rgba(var(--color-success-rgb), 0.3);
}

.error-icon {
  background: rgba(var(--color-error-rgb), 0.15);
  color: var(--color-error-500);
  border: 1px solid rgba(var(--color-error-rgb), 0.3);
}

.stat-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.stat-label {
  font-size: var(--text-sm);
  color: var(--color-content-muted);
  font-weight: 500;
}

.stat-value {
  font-size: var(--text-lg);
  color: var(--color-content-heading);
  font-weight: 600;
}

.success-value {
  color: var(--color-success-500);
}

.error-value {
  color: var(--color-error-500);
}

.folder-message {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md);
  background: rgba(var(--color-brand-500-rgb), 0.1);
  border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-sm);
}

.message-icon {
  font-size: var(--text-lg);
  color: var(--color-brand-500);
  flex-shrink: 0;
}

.message-text {
  margin: 0;
  font-size: var(--text-sm);
  color: var(--color-content-default);
  line-height: 1.5;
}

.highlight-text {
  color: var(--color-brand-500);
  font-weight: 600;
  font-size: var(--text-base);
}

.folder-warning {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  padding: var(--space-sm) var(--space-md);
  background: rgba(var(--color-warning-rgb), 0.1);
  border: 1px solid rgba(var(--color-warning-rgb), 0.2);
  border-radius: var(--radius-md);
}

.warning-icon {
  color: var(--color-warning-500);
  font-size: var(--text-sm);
  flex-shrink: 0;
}

.warning-text {
  font-size: var(--text-xs);
  color: var(--color-content-default);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-sm);
  padding: var(--space-md);
}
</style>
