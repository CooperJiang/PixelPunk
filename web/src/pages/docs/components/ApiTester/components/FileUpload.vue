<template>
  <div class="cyber-section">
    <div class="section-header">
      <i class="fas fa-file-image"></i>
      <h4>{{ $t('docs.fileUpload.title') }}</h4>
      <div class="section-line"></div>
    </div>
    <div
      class="cyber-upload-zone"
      :class="{
        'has-files': selectedFiles.length > 0,
        'drag-over': isDragOver,
        error: !selectedFiles.length && showValidation,
      }"
      @click="triggerFileInput"
      @drop="handleDrop"
      @dragover.prevent="handleDragOver"
      @dragleave="handleDragLeave"
    >
      <input ref="fileInput" type="file" multiple accept="image/*" @change="handleFileSelect" style="display: none" />

      <div v-if="!selectedFiles.length" class="upload-placeholder">
        <div class="upload-icon">
          <i class="fas fa-images"></i>
        </div>
        <div class="upload-text">
          <p class="primary-text">{{ $t('docs.fileUpload.dragOrClick') }}</p>
          <p class="secondary-text">{{ $t('docs.fileUpload.supportedFormats') }}</p>
        </div>
      </div>

      <div v-else class="files-list">
        <div class="files-header">
          <span class="files-count">{{ $t('docs.fileUpload.selectedCount', { count: selectedFiles.length }) }}</span>
          <span class="files-size">{{ $t('docs.fileUpload.totalSize', { size: formatFileSize(totalFileSize) }) }}</span>
        </div>
        <div class="files-container">
          <div v-for="(file, index) in selectedFiles" :key="index" class="file-cyber-item">
            <div class="file-icon">
              <i class="fas fa-file-image"></i>
            </div>
            <div class="file-info">
              <div class="file-name">{{ file.name }}</div>
              <div class="file-size">{{ formatFileSize(file.size) }}</div>
            </div>
            <button @click.stop="removeFile(index)" class="remove-file-btn">
              <i class="fas fa-times"></i>
            </button>
          </div>
        </div>
        <button @click.stop="triggerFileInput" class="add-more-files-btn">
          <i class="fas fa-plus"></i>
          <span>{{ $t('docs.fileUpload.addMore') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed } from 'vue'
  import { useFileUpload } from '../composables/useFileUpload'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'FileUpload',
  })

  interface Props {
    selectedFiles: File[]
    totalFileSize: number
    showValidation: boolean
    isDragOver: boolean
  }

  const props = defineProps<Props>()

  const emit = defineEmits<{
    'update:selectedFiles': [files: File[]]
    'update:showValidation': [value: boolean]
    'update:isDragOver': [value: boolean]
  }>()

  /* 创建响应式引用以便于 composable 使用 */
  const selectedFiles = computed({
    get: () => props.selectedFiles,
    set: (value) => emit('update:selectedFiles', value),
  })

  const showValidation = computed({
    get: () => props.showValidation,
    set: (value) => emit('update:showValidation', value),
  })

  const isDragOver = computed({
    get: () => props.isDragOver,
    set: (value) => emit('update:isDragOver', value),
  })

  const {
    fileInput,
    formatFileSize,
    triggerFileInput,
    handleFileSelect,
    handleDrop,
    handleDragOver,
    handleDragLeave,
    removeFile,
  } = useFileUpload(selectedFiles, showValidation, isDragOver)
</script>

<style scoped>
  .cyber-section {
    margin-bottom: 1.5rem;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.75rem;
    padding-bottom: 0.4rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .section-header i {
    color: var(--color-brand-500);
    font-size: 1rem;
  }

  .section-header h4 {
    color: var(--color-text-content);
    font-size: 1rem;
    font-weight: 600;
    margin: 0;
  }

  .section-line {
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.3), transparent);
  }

  .cyber-upload-zone {
    border: 2px dashed rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 1.5rem;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s ease;
    background: rgba(var(--color-brand-500-rgb), 0.02);
    position: relative;
    overflow: hidden;
  }

  .cyber-upload-zone:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    background: rgba(var(--color-brand-500-rgb), 0.05);
  }

  .cyber-upload-zone.has-files {
    border-color: rgba(var(--color-success-rgb), 0.4);
    background: rgba(var(--color-success-rgb), 0.05);
  }

  .cyber-upload-zone.drag-over {
    border-color: rgba(var(--color-success-rgb), 0.6);
    background: rgba(var(--color-success-rgb), 0.1);
    transform: scale(1.02);
  }

  .cyber-upload-zone.error {
    border-color: rgba(var(--color-badge-accent-text-rgb), 0.4);
    background: rgba(var(--color-badge-accent-text-rgb), 0.05);
  }

  .upload-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    position: relative;
    z-index: 1;
  }

  .upload-icon {
    font-size: 2.5rem;
    color: var(--color-brand-500);
    margin-bottom: 0.5rem;
    animation: float 3s ease-in-out infinite;
  }

  .upload-text {
    text-align: center;
  }

  .primary-text {
    color: var(--color-text-content);
    font-size: 1rem;
    font-weight: 500;
    margin: 0 0 0.5rem 0;
  }

  .secondary-text {
    color: var(--color-slate-500);
    font-size: 0.875rem;
    margin: 0;
  }

  .files-list {
    text-align: left;
    position: relative;
    z-index: 1;
  }

  .files-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding: 0.5rem 0;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .files-count {
    color: var(--color-brand-500);
    font-weight: 600;
    font-size: 0.875rem;
  }

  .files-size {
    color: var(--color-slate-500);
    font-size: 0.8rem;
  }

  .files-container {
    margin-bottom: 1rem;
    max-height: 200px;
    overflow-y: auto;
  }

  .file-cyber-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    background: var(--color-background-700);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    margin-bottom: 0.5rem;
    transition: all 0.3s ease;
  }

  .file-cyber-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
    transform: translateX(4px);
  }

  .file-icon {
    color: var(--color-brand-500);
    font-size: 1rem;
    width: 20px;
    text-align: center;
  }

  .file-info {
    flex: 1;
  }

  .file-name {
    color: var(--color-text-content);
    font-size: 0.875rem;
    font-weight: 500;
    margin-bottom: 0.25rem;
    word-break: break-all;
  }

  .file-size {
    color: var(--color-slate-500);
    font-size: 0.75rem;
  }

  .remove-file-btn {
    background: rgba(var(--color-badge-accent-text-rgb), 0.1);
    border: 1px solid rgba(var(--color-badge-accent-text-rgb), 0.2);
    color: var(--color-error-500);
    width: 28px;
    height: 28px;
    border-radius: var(--radius-sm);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    transition: all 0.3s ease;
  }

  .remove-file-btn:hover {
    background: rgba(var(--color-badge-accent-text-rgb), 0.2);
    transform: scale(1.1);
  }

  .add-more-files-btn {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    padding: 0.6rem 1rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: all 0.3s ease;
    width: 100%;
    justify-content: center;
  }

  .add-more-files-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateY(-1px);
  }

  @media (max-width: 768px) {
    .cyber-upload-zone {
      padding: 1.2rem;
    }
  }

  .files-container::-webkit-scrollbar {
    width: 6px;
  }

  .files-container::-webkit-scrollbar-track {
    background: var(--color-hover-bg-neutral);
    border-radius: var(--radius-sm);
  }

  .files-container::-webkit-scrollbar-thumb {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.6) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.6) 100%
    );
    border-radius: var(--radius-sm);
  }

  .files-container::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.8) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.8) 100%
    );
  }
</style>
