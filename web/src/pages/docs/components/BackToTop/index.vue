<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'BackToTop',
  })

  const props = defineProps<Props>()

  interface Props {
    showBackToTop: boolean
    backToTop: () => void
  }

  const onBackToTop = () => {
    props.backToTop()
  }
</script>

<template>
  <button v-show="showBackToTop" class="back-to-top" :title="$t('docs.backToTop')" @click="onBackToTop">
    <i class="fas fa-arrow-up" />
  </button>
</template>

<style scoped>
  .back-to-top {
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    width: 40px;
    height: 40px;
    background: linear-gradient(135deg, var(--color-brand-500) 0%, var(--color-error-500) 100%);
    border: none;
    border-radius: var(--radius-full);
    color: var(--color-black);
    cursor: pointer;
    font-size: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    z-index: var(--z-index-dropdown);
    box-shadow: 0 4px 12px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .back-to-top:hover {
    transform: translateY(-2px) scale(1.05);
    box-shadow: 0 6px 20px rgba(var(--color-brand-500-rgb), 0.4);
  }

  @media (max-width: 768px) {
    .back-to-top {
      bottom: 1.5rem;
      right: 1.5rem;
      width: 44px;
      height: 44px;
      font-size: 1.125rem;
    }
  }
</style>
