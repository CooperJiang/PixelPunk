<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useUploadStore } from '@/store/upload'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'UploadEntry',
  })

  const router = useRouter()
  const uploadStore = useUploadStore()
  const { $t } = useTexts()
  const fileInput = ref<HTMLInputElement>()

  const uploadTitle = computed(() => {
    return uploadStore.statistics.total > 0
      ? $t('dashboard.uploadEntry.title.withCount', { count: uploadStore.statistics.total })
      : $t('dashboard.uploadEntry.title.default')
  })
  const isDragOver = ref(false)

  const handleDragOver = (e: DragEvent) => {
    if (e.dataTransfer?.types.includes('Files')) {
      isDragOver.value = true
    }
  }

  const handleDragLeave = () => {
    isDragOver.value = false
  }

  const handleDrop = async (e: DragEvent) => {
    isDragOver.value = false
    const files = e.dataTransfer?.files
    if (files && files.length > 0) {
      await uploadStore.addFiles(files)
      router.push('/upload')
    }
  }

  const triggerFileSelect = (e: Event) => {
    e.stopPropagation()
    fileInput.value?.click()
  }

  const handleFileSelect = async (event: Event) => {
    const { files } = event.target as HTMLInputElement
    if (files && files.length > 0) {
      await uploadStore.addFiles(files)
      router.push('/upload')
    }
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }
</script>

<template>
  <div
    class="upload-entry"
    :class="{
      'drag-over': isDragOver,
      'has-files': uploadStore.statistics.total > 0,
    }"
    :title="uploadTitle"
    @dragover.prevent="handleDragOver"
    @dragleave.prevent="handleDragLeave"
    @drop.prevent="handleDrop"
    @click="triggerFileSelect"
  >
    <div class="entry-content">
      <div class="upload-icon">
        <i class="fas fa-cloud-upload-alt" />
      </div>

      <div v-if="uploadStore.statistics.total > 0" class="queue-badge">
        {{ uploadStore.statistics.total }}
      </div>
    </div>

    <input ref="fileInput" type="file" multiple accept="image/*" style="display: none" @change="handleFileSelect" />
  </div>
</template>

<style scoped>
  .upload-entry {
    position: fixed;
    bottom: 24px;
    right: 24px;
    z-index: 1000;
    width: 64px;
    height: 64px;
    background: linear-gradient(135deg, var(--color-brand-500) 0%, var(--color-brand-600) 50%, var(--color-brand-700) 100%);
    border: 2px solid var(--color-brand-400);
    border-radius: var(--radius-full);
    cursor: pointer;
    transition: all var(--transition-slow) var(--ease-in-out);
    backdrop-filter: blur(8px);
    box-shadow:
      0 0 20px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 40px rgba(var(--color-brand-500-rgb), 0.3),
      0 4px 12px rgba(0, 0, 0, 0.3);
    animation: float-breath 3s ease-in-out infinite;
  }

  .upload-entry::before {
    content: '';
    position: absolute;
    inset: -2px;
    border-radius: var(--radius-full);
    padding: 2px;
    background: linear-gradient(135deg, var(--color-brand-300), transparent, var(--color-brand-500));
    -webkit-mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    opacity: 0.5;
    pointer-events: none;
  }

  .upload-entry:hover {
    transform: translateY(-4px) scale(1.08);
    border-color: var(--color-brand-300);
    box-shadow:
      0 0 30px rgba(var(--color-brand-400-rgb), 0.8),
      0 0 60px rgba(var(--color-brand-500-rgb), 0.5),
      0 8px 20px rgba(0, 0, 0, 0.4);
  }

  .upload-entry:hover::before {
    opacity: 0.8;
  }

  .upload-entry:active {
    transform: translateY(-2px) scale(1.04);
  }

  .upload-entry.drag-over {
    background: linear-gradient(135deg, var(--color-success-500) 0%, var(--color-success-600) 50%, var(--color-brand-500) 100%);
    border-color: var(--color-success-400);
    transform: scale(1.12);
    box-shadow:
      0 0 40px rgba(var(--color-success-rgb), 0.7),
      0 0 80px rgba(var(--color-success-rgb), 0.4),
      0 8px 20px rgba(0, 0, 0, 0.4);
  }

  .upload-entry.has-files {
    border-color: var(--color-brand-300);
    animation:
      float-breath 3s ease-in-out infinite,
      has-files-pulse 2s ease-in-out infinite;
  }

  .entry-content {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .upload-icon {
    color: var(--color-text-on-brand);
    font-size: 26px;
    transition: all var(--transition-slow) var(--ease-in-out);
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
  }

  .upload-entry:hover .upload-icon {
    transform: scale(1.15);
    filter: drop-shadow(0 0 8px rgba(255, 255, 255, 0.6));
  }

  .queue-badge {
    position: absolute;
    top: -4px;
    right: -4px;
    min-width: 22px;
    height: 22px;
    background: var(--color-error-500);
    color: var(--color-text-on-error);
    font-size: 11px;
    font-weight: 700;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--color-background-900);
    font-family: var(--font-mono);
    box-shadow:
      0 0 15px rgba(var(--color-error-rgb), 0.7),
      0 0 30px rgba(var(--color-error-rgb), 0.4);
    animation: badge-pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }

  @keyframes float-breath {
    0%,
    100% {
      box-shadow:
        0 0 20px rgba(var(--color-brand-500-rgb), 0.6),
        0 0 40px rgba(var(--color-brand-500-rgb), 0.3),
        0 4px 12px rgba(0, 0, 0, 0.3);
    }
    50% {
      box-shadow:
        0 0 30px rgba(var(--color-brand-500-rgb), 0.8),
        0 0 60px rgba(var(--color-brand-500-rgb), 0.5),
        0 6px 16px rgba(0, 0, 0, 0.3);
    }
  }

  @keyframes has-files-pulse {
    0%,
    100% {
      border-color: var(--color-brand-300);
    }
    50% {
      border-color: var(--color-brand-200);
    }
  }

  @keyframes badge-pulse {
    0%,
    100% {
      transform: scale(1);
      box-shadow:
        0 0 15px rgba(var(--color-error-rgb), 0.7),
        0 0 30px rgba(var(--color-error-rgb), 0.4);
    }
    50% {
      transform: scale(1.05);
      box-shadow:
        0 0 20px rgba(var(--color-error-rgb), 0.9),
        0 0 40px rgba(var(--color-error-rgb), 0.6);
    }
  }

  @media (max-width: 768px) {
    .upload-entry {
      bottom: 20px;
      right: 20px;
      width: 56px;
      height: 56px;
    }

    .upload-icon {
      font-size: 20px;
    }

    .queue-badge {
      min-width: 20px;
      height: 20px;
      font-size: 10px;
      top: -3px;
      right: -3px;
    }
  }

  @media (max-width: 480px) {
    .upload-entry {
      bottom: 16px;
      right: 16px;
      width: 52px;
      height: 52px;
    }

    .upload-icon {
      font-size: 18px;
    }

    .queue-badge {
      min-width: 18px;
      height: 18px;
      font-size: 9px;
    }
  }
</style>
