<script setup lang="ts">
  const _props = defineProps<{
    filesLength: number
    isLoadingMore: boolean
    hasReachedLimit: boolean
    hasMore: boolean
  }>()
</script>

<template>
  <div v-if="filesLength > 0" class="mb-4 mt-6 flex justify-center">
    <div v-if="isLoadingMore" class="loading-indicator">
      <div class="loading-spinner">
        <i class="fas fa-circle-notch fa-spin" />
      </div>
      <span>{{ $t('explore.footer.loadingMore') }}</span>
    </div>
    <div v-else-if="hasReachedLimit" class="end-message">
      <i class="fas fa-info-circle mr-2" />
      {{ $t('explore.footer.enoughContent') }}
    </div>
    <div v-else-if="!hasMore" class="end-message">
      <i class="fas fa-check-circle mr-2" />
      {{ $t('explore.footer.loadedAll') }}
    </div>
  </div>
</template>

<style scoped></style>
