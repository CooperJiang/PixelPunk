<script setup lang="ts">
  import { computed, defineEmits, defineProps } from 'vue'
  import { formatDate } from '@/utils/formatting/format'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps({
    showBackButton: {
      type: Boolean,
      default: false,
    },
    viewCount: {
      type: Number,
      default: 0,
    },
    createdAt: {
      type: String,
      required: true,
    },
  })

  const emit = defineEmits(['back-click'])

  const formattedDate = computed(() => formatDate(props.createdAt, 'YYYY-MM-DD'))

  const handleBackClick = () => {
    emit('back-click')
  }
</script>

<template>
  <div class="mobile-share-info">
    <button v-if="showBackButton" class="back-button" @click="handleBackClick">
      <i class="fas fa-arrow-left" />
      <span>{{ $t('share.mobile.backToRoot') }}</span>
    </button>
    <div v-else class="info-item">
      <i class="fas fa-eye" />
      <span>{{ $t('share.mobile.viewCount', { count: viewCount || 0 }) }}</span>
    </div>
    <div class="info-item">
      <i class="fas fa-calendar-alt" />
      <span>{{ formattedDate }}</span>
    </div>
  </div>
</template>

<style scoped>
  .mobile-share-info {
    background: rgba(var(--color-background-800-rgb), 0.5);
    border-radius: var(--radius-md);
    padding: var(--space-md);
    margin-bottom: 1rem;
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-sm);
  }

  .info-item,
  .back-button {
    display: flex;
    align-items: center;
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.8);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    padding: var(--space-md) 0.5rem;
  }

  .info-item i,
  .back-button i {
    margin-right: 0.25rem;
    color: var(--color-brand-500);
  }

  .back-button {
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .back-button:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }
</style>
