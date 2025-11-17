<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'
  import type { FileInfo } from '@/api/admin/files'

  defineOptions({
    name: 'AIDescription',
  })

  const { $t } = useTexts()

  defineProps<{
    fileInfo: FileInfo
  }>()
</script>

<template>
  <div v-if="fileInfo.ai_info?.description" class="description-section">
    <div class="info-section">
      <h3 class="section-title-compact mb-2 flex items-center">
        <i class="fas fa-comment-alt mr-2" />{{ $t('components.fileDetailModal.ai.description') }}
      </h3>
      <div class="description-content rounded-lg border border-subtle bg-background-300 p-3">
        <p class="text-sm leading-relaxed text-content">{{ fileInfo.ai_info.description }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .info-section {
    @apply rounded-xl border p-3 backdrop-blur-md;
    background: var(--color-background-800);
    border-color: var(--color-border-default);
    box-shadow: var(--shadow-cyber-sm);
  }

  .section-title-compact {
    @apply mb-0 pb-1 text-sm font-semibold;
    color: var(--color-brand-500);
  }

  .description-content {
    @apply relative overflow-hidden;
    background: var(--color-background-700);
    backdrop-filter: blur(5px);
  }
</style>
