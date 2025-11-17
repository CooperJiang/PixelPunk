<script setup lang="ts">
  import ShareButton from './ShareButton.vue'
  import type { ShareButtonProps, ShareButtonEmits } from './types'

  defineOptions({
    name: 'Share',
  })

  defineProps<ShareButtonProps>()
  defineEmits<ShareButtonEmits>()
</script>

<template>
  <ShareButton v-bind="$attrs" />
</template>
