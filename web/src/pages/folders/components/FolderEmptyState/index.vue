<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'FolderEmptyState',
  })

  defineEmits<{
    'create-folder': []
  }>()

  const { $t } = useTexts()
</script>

<template>
  <div class="empty-state">
    <div class="empty-icon">
      <i class="fas fa-folder-open" />
    </div>
    <h3 class="empty-title">{{ $t('folders.emptyState.title') }}</h3>
    <p class="empty-desc">{{ $t('folders.emptyState.description') }}</p>
    <CyberButton type="secondary" icon="folder-plus" @click="$emit('create-folder')">
      {{ $t('folders.emptyState.createButton') }}
    </CyberButton>
  </div>
</template>

<style scoped>
  .empty-state {
    @apply relative mt-14 flex min-h-[200px] w-full flex-col items-center justify-center overflow-hidden rounded-xl border border-default px-6 py-10 text-center backdrop-blur-sm;
    background-color: rgba(var(--color-background-800-rgb), 0.6);
  }

  .empty-icon {
    @apply relative z-[2] mb-5 flex h-[70px] w-[70px] items-center justify-center rounded-xl border border-default shadow-sm transition-all duration-300;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-brand-600-rgb), 0.1));
  }

  .empty-state:hover .empty-icon {
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .empty-icon i {
    @apply relative z-[3] text-3xl text-brand-400;
    filter: drop-shadow(0 0 6px rgba(var(--color-brand-500-rgb), 0.3));
  }

  .empty-title {
    @apply relative z-[2] mb-3 text-lg font-semibold tracking-wide text-content-heading;
    font-family: 'Orbitron', 'SF Mono', monospace;
  }

  .empty-desc {
    @apply relative z-[2] mb-6 max-w-md text-sm leading-relaxed text-content-muted;
  }
</style>
