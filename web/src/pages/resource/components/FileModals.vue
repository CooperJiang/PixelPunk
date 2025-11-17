<script setup lang="ts">
  import { computed } from 'vue'

  const props = defineProps({
    showDetailModal: {
      type: Boolean,
      default: false,
    },
    selectedImage: {
      type: Object,
      default: null,
    },
  })

  const emit = defineEmits(['update:showDetailModal', 'download'])

  const _localDetailModal = computed({
    get: () => props.showDetailModal,
    set: (value) => emit('update:showDetailModal', value),
  })
</script>

<template>
  <div>
    <CyberFileDetailModal
      :model-value="showDetailModal"
      :fileInfo="selectedImage"
      :is-admin="false"
      @update:model-value="$emit('update:showDetailModal', $event)"
      @download="$emit('download', selectedImage?.id)"
    />
  </div>
</template>
