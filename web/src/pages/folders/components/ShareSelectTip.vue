<script setup lang="ts">
  defineOptions({
    name: 'ShareSelectTip',
  })

  interface Props {
    visible: boolean
  }

  defineProps<Props>()
</script>

<template>
  <div v-if="visible" class="share-select-tip">
    <i class="fas fa-hand-pointer mr-2" />
    {{ $t('folders.shareSelectTip.message') }}
  </div>
</template>

<style scoped>
  .share-select-tip {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.75rem 1.5rem;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: rgba(var(--color-brand-500-rgb), 0.9);
    font-size: 0.9rem;
    font-weight: 500;
    letter-spacing: 0.3px;
    margin: 0 1rem;
    animation: tipPulse 2s ease-in-out infinite;
  }

  .share-select-tip i {
    font-size: 1rem;
    color: rgba(var(--color-brand-500-rgb), 0.8);
  }

  @keyframes tipPulse {
    0%,
    100% {
      background: rgba(var(--color-brand-500-rgb), 0.1);
      box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.2);
    }
    50% {
      background: rgba(var(--color-brand-500-rgb), 0.15);
      box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.4);
    }
  }
</style>
