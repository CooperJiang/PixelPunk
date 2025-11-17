<script setup lang="ts">
  import type { FloatUploadStatistics } from '../types'

  defineOptions({
    name: 'FileTooltip',
  })

  interface FileTooltipProps {
    fileName: string | null
    statistics: FloatUploadStatistics
    visible: boolean
  }

  defineProps<FileTooltipProps>()
</script>

<template>
  <div v-if="visible && fileName" class="file-tooltip">
    <div class="tooltip-content">
      <div class="file-name">{{ fileName }}</div>
      <div class="file-stats">
        <span class="stat-item"> <i class="fas fa-cloud-upload-alt" /> {{ statistics.uploading }} </span>
        <span class="stat-item"> <i class="fas fa-clock" /> {{ statistics.pending }} </span>
        <span class="stat-item"> <i class="fas fa-check" /> {{ statistics.completed }} </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .file-tooltip {
    @apply pointer-events-none absolute -top-16 left-1/2 -translate-x-1/2 opacity-100 transition-opacity duration-300;
  }

  .tooltip-content {
    @apply whitespace-nowrap rounded-lg border border-brand-400 bg-background-900 px-3 py-2;
    box-shadow: 0 0 15px var(--color-brand-500);
  }

  .file-name {
    @apply mb-0.5 text-xs font-medium text-brand-400;
    text-shadow: 0 0 5px var(--color-brand-500);
  }

  .file-stats {
    @apply flex gap-3 font-mono text-[10px] text-content;
  }

  .stat-item {
    @apply flex items-center gap-1;
  }

  .file-stats i {
    @apply text-[10px] text-brand-400;
    text-shadow: 0 0 3px var(--color-brand-500);
  }
</style>
