<script setup lang="ts">
  import { ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { useToast } from '@/components/Toast/useToast'
  import type { WatermarkConfig } from '../types'

  defineOptions({
    name: 'WatermarkFileConfig',
  })

  interface Props {
    config: WatermarkConfig
  }

  interface Emits {
    (e: 'update:config', config: WatermarkConfig): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const { $t } = useTexts()
  const toast = useToast()
  const imageInput = ref<HTMLInputElement>()

  const handleImageUpload = (event: Event) => {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]

    if (file) {
      if (!file.type.startsWith('image/')) {
        toast.error($t('components.watermark.messages.fileTypeError'))
        return
      }

      if (file.size > 2 * 1024 * 1024) {
        toast.error($t('components.watermark.messages.fileSizeError'))
        return
      }

      const reader = new FileReader()
      reader.onload = (e) => {
        const result = e.target?.result as string

        if (result && result.includes(',')) {
          const base64Data = result.split(',')[1]
          const newConfig = {
            ...props.config,
            fileBase64: base64Data,
            imageFile: file,
          }
          emit('update:config', newConfig)
        }
      }
      reader.readAsDataURL(file)
    }
  }

  const useDefaultImage = async () => {
    try {
      const resp = await fetch('/pixelpunk.png')
      const blob = await resp.blob()
      const reader = new FileReader()
      reader.onload = () => {
        const dataUrl = reader.result as string
        const base64 = dataUrl.split(',')[1] || ''
        const newConfig = {
          ...props.config,
          fileBase64: base64,
          imageFile: undefined,
        }
        emit('update:config', newConfig)
      }
      reader.readAsDataURL(blob)
    } catch (e) {
      // Silently fail
    }
  }

  const handleImageError = () => {
    // Silently fail
  }
</script>

<template>
  <div class="section-group">
    <div class="section-header">
      <h4 class="section-title">{{ $t('components.watermark.imageConfig.title') }}</h4>
    </div>
    <div class="section-content">
      <div class="config-item">
        <label class="config-label">{{ $t('components.watermark.imageConfig.image') }}</label>
        <div class="upload-area">
          <div class="upload-buttons">
            <input ref="imageInput" type="file" accept="image/*" class="hidden" @change="handleImageUpload" />
            <cyberButton variant="outline" size="sm" class="upload-btn" @click="imageInput?.click()">
              <i class="fas fa-upload mr-2" />
              {{ $t('components.watermark.imageConfig.upload') }}
            </cyberButton>
            <cyberButton variant="outline" size="sm" class="upload-btn" @click="useDefaultImage">
              <i class="fas fa-image mr-2" />
              {{ $t('components.watermark.imageConfig.useDefault') }}
            </cyberButton>
          </div>

          <div v-if="config.fileBase64" class="image-preview">
            <div class="preview-container">
              <img
                :src="'data:image/png;base64,' + config.fileBase64"
                :alt="$t('components.watermark.imageConfig.previewAlt')"
                class="preview-image"
                @error="handleImageError"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  @import '../styles/common.css';

  .upload-area {
    @apply space-y-3;
  }

  .upload-buttons {
    @apply flex space-x-3;
  }

  .upload-btn {
    @apply flex-1;
  }

  .image-preview {
    @apply flex items-center justify-center rounded-lg border-2 border-dashed p-4;
    border-color: var(--color-border-default);
    background: var(--color-background-800);
  }

  .preview-container {
    @apply relative;
  }

  .preview-image {
    @apply max-h-20 max-w-full rounded-md shadow-lg;
  }

  @media (max-width: 640px) {
    .upload-buttons {
      @apply flex-col space-x-0 space-y-2;
    }
  }
</style>
