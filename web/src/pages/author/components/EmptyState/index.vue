<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { onMounted, onUnmounted, ref, computed } from 'vue'

  defineOptions({
    name: 'AuthorEmptyState',
  })

  interface Props {
    iconClass?: string
    title: string
    description: string
    showRetryButton?: boolean
    retryText?: string
    statusCode?: string
    suggestion?: string
  }

  interface Emits {
    (e: 'retry'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    iconClass: 'fas fa-folder-open',
    showRetryButton: false,
    retryText: 'Retry',
    statusCode: 'INFO_100',
    suggestion: '',
  })

  const emit = defineEmits<Emits>()

  const currentTime = ref('')
  let timeInterval: ReturnType<typeof setInterval> | null = null

  const renderedSuggestion = computed(() => props.suggestion ?? '')

  const updateTime = () => {
    const now = new Date()
    currentTime.value = now.toLocaleTimeString(getCurrentLocale(), {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  }

  const handleRetry = () => {
    emit('retry')
  }

  onMounted(() => {
    updateTime()
    timeInterval = setInterval(updateTime, 1000)
  })

  onUnmounted(() => {
    if (timeInterval) {
      clearInterval(timeInterval)
      timeInterval = null
    }
  })
</script>

<template>
  <div class="author-empty-state">
    <div class="empty-icon-wrapper">
      <i :class="iconClass" />
    </div>
    <h3 class="empty-title">{{ title }}</h3>
    <p class="empty-description">{{ description }}</p>
    <p v-if="renderedSuggestion" class="empty-suggestion">{{ renderedSuggestion }}</p>
    <button v-if="showRetryButton" class="retry-button" @click="handleRetry">
      <i class="fas fa-sync-alt" />
      <span>{{ retryText }}</span>
    </button>
  </div>
</template>

<style scoped lang="scss" src="./EmptyState.scss"></style>
