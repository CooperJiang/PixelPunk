<script setup lang="ts">
  import { ref } from 'vue'
  import type { FileContainerProps, FileContainerEmits } from '../types'

  defineOptions({
    name: 'FileContainer',
  })

  defineProps<FileContainerProps>()
  defineEmits<FileContainerEmits>()

  const imageRef = ref<HTMLImageElement>()

  defineExpose({
    imageRef,
  })
</script>

<template>
  <div
    class="enhanced-image-container"
    @wheel.prevent="$emit('wheel', $event)"
    @mousedown="$emit('mousedown', $event)"
    @mousemove="$emit('mousemove', $event)"
    @mouseup="$emit('mouseup', $event)"
    @mouseleave="$emit('mouseleave', $event)"
    @touchstart="$emit('touchstart', $event)"
    @touchmove="$emit('touchmove', $event)"
    @touchend="$emit('touchend', $event)"
  >
    <img
      ref="imageRef"
      :src="fileUrl"
      :alt="fileName"
      class="enhanced-fullscreen-image"
      :class="{ 'fill-mode': shouldUseFillMode }"
      :style="fileStyle"
      draggable="false"
      @click.stop="$emit('fileClick', $event)"
      @load="$emit('fileLoad', $event)"
      @error="$emit('fileError', $event)"
    />
  </div>
</template>

<style scoped>
  .enhanced-image-container {
    @apply absolute inset-0 flex items-center justify-center overflow-hidden;
  }

  .enhanced-fullscreen-image {
    @apply block cursor-grab select-none;
  }

  .enhanced-fullscreen-image:active {
    @apply cursor-grabbing;
  }

  .enhanced-fullscreen-image.fill-mode {
    @apply h-full w-full object-cover;
  }
</style>
