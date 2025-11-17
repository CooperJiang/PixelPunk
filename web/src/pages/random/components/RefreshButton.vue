<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    loading: boolean
  }

  interface Emits {
    (e: 'refresh'): void
  }

  defineProps<Props>()
  const emit = defineEmits<Emits>()
  const { $t } = useTexts()

  const handleRefresh = () => {
    emit('refresh')
  }
</script>

<template>
  <div class="fixed bottom-6 right-6 z-30">
    <button :disabled="loading" class="refresh-btn" :title="$t('random.refreshButton.title')" @click="handleRefresh">
      <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }" />
    </button>
  </div>
</template>

<style scoped>
  .refresh-btn {
    width: 40px;
    height: 40px;
    border-radius: var(--radius-full);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1rem;
    transition: all 0.3s ease;
    backdrop-filter: blur(5px);
    box-shadow: 0 2px 10px rgba(var(--color-background-900-rgb), 0.2);
  }

  .refresh-btn:hover:not(:disabled) {
    transform: scale(1.1);
    background: rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow: 0 4px 15px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .refresh-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  @media (max-width: 768px) {
    .refresh-btn {
      width: 36px;
      height: 36px;
      bottom: 20px;
      right: 20px;
    }
  }
</style>
