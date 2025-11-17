<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { FileDetailModalProps, FileDetailModalEmits } from './types'
  import FilePreview from './components/FilePreview.vue'
  import BasicInfo from './components/BasicInfo.vue'
  import VisualAnalysis from './components/VisualAnalysis.vue'
  import AIDescription from './components/AIDescription.vue'

  defineOptions({
    name: 'FileDetailModal',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<FileDetailModalProps>(), {
    isAdmin: false,
  })

  const fileInfo = computed(() => props.fileInfo)

  defineEmits<FileDetailModalEmits>()

  const getDialogTitle = () => {
    const baseTitle = $t('components.fileDetailModal.title')
    const fileName = fileInfo.value?.display_name || fileInfo.value?.original_name || $t('components.fileDetailModal.untitled')
    return `${baseTitle} - ${fileName}`
  }
</script>

<template>
  <CyberDialog
    :model-value="modelValue"
    :title="getDialogTitle()"
    width="80vw"
    max-width="1400px"
    max-height="90vh"
    :show-default-footer="false"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <div v-if="fileInfo" class="file-detail-content">
      <div class="main-info-grid">
        <div class="preview-section">
          <FilePreview :file-info="fileInfo" />
        </div>
        <div class="info-section-wrapper">
          <BasicInfo :file-info="fileInfo" />
        </div>
        <div class="analysis-section">
          <VisualAnalysis :file-info="fileInfo" />
        </div>
      </div>

      <div class="ai-description-wrapper">
        <AIDescription :file-info="fileInfo" />
      </div>
    </div>
  </CyberDialog>
</template>

<style scoped>
  .file-detail-content {
    @apply flex flex-col gap-3;
  }

  .main-info-grid {
    @apply grid grid-cols-12 gap-3;
  }

  .preview-section {
    @apply col-span-4 flex;
  }

  .preview-section > * {
    @apply flex-1;
  }

  .info-section-wrapper {
    @apply col-span-4 flex;
  }

  .info-section-wrapper > * {
    @apply flex-1;
  }

  .analysis-section {
    @apply col-span-4 flex;
  }

  .analysis-section > * {
    @apply flex-1;
  }

  .ai-description-wrapper {
    @apply w-full;
  }

  @media (max-width: 1200px) {
    .main-info-grid {
      @apply grid-cols-1 gap-3;
    }

    .preview-section,
    .info-section-wrapper,
    .analysis-section {
      @apply col-span-12;
    }
  }

  @media (max-width: 768px) {
    .file-detail-content {
      @apply gap-3;
    }
  }

  @media (max-width: 768px) {
    :deep(.dialog-container) {
      width: 90vw !important;
    }
  }
</style>
