<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'TesterHeader',
  })

  interface Props {
    currentDomain: string
    canSend: boolean
    isLoading: boolean
  }

  defineProps<Props>()

  defineEmits<{
    'send-request': []
  }>()

  const { $t } = useTexts()
</script>

<template>
  <div class="cyber-tester-header">
    <div class="header-left">
      <div class="endpoint-display">
        <div class="method-badge">POST</div>
        <div class="endpoint-url">
          <code>{{ currentDomain }}/api/v1/external/upload</code>
        </div>
      </div>
    </div>
    <div class="header-right">
      <button
        class="cyber-send-btn"
        @click="$emit('send-request')"
        :disabled="!canSend || isLoading"
        :class="{ loading: isLoading }"
      >
        <div class="btn-content">
          <i :class="isLoading ? 'fas fa-spinner fa-spin' : 'fas fa-paper-plane'"></i>
          <span>{{ isLoading ? $t('docs.tester.executing') : $t('docs.tester.send') }}</span>
        </div>
        <div class="btn-glow"></div>
      </button>
    </div>
  </div>
</template>

<style scoped>
  .cyber-tester-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.6rem 1rem;
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.08) 0%,
      rgba(var(--color-background-800-rgb), 0.6) 100%
    );
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    position: relative;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .endpoint-display {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .method-badge {
    background: linear-gradient(135deg, var(--color-warning-500) 0%, rgba(var(--color-warning-rgb), 0.9) 100%);
    color: var(--color-text-on-warning);
    padding: 0.3rem 0.6rem;
    border-radius: var(--radius-sm);
    font-weight: 700;
    font-size: 0.7rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    box-shadow: 0 2px 8px rgba(var(--color-warning-rgb), 0.4);
    border: 1px solid rgba(var(--color-warning-rgb), 0.6);
  }

  .endpoint-url {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 0.4rem 0.8rem;
    font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
    position: relative;
    overflow: hidden;
  }

  @keyframes shimmer {
    0% {
      left: -100%;
    }
    100% {
      left: 100%;
    }
  }

  .endpoint-url code {
    color: var(--color-brand-500);
    font-size: 0.8rem;
    font-weight: 500;
  }

  .cyber-send-btn {
    background: linear-gradient(135deg, var(--color-brand-500) 0%, var(--color-brand-500-light) 100%);
    color: var(--color-black);
    border: none;
    padding: 0.4rem 1rem;
    border-radius: var(--radius-sm);
    font-weight: 600;
    font-size: 0.8rem;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
    box-shadow: 0 2px 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-send-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none !important;
  }

  .cyber-send-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(var(--color-brand-500-rgb), 0.4);
    filter: brightness(1.1);
  }

  .cyber-send-btn.loading {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.3) 0%, rgba(14, 165, 233, 0.3) 100%);
  }

  .btn-content {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    position: relative;
    z-index: 1;
  }

  .btn-glow {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(var(--color-content-default-rgb), 0.2) 0%, transparent 70%);
    transform: translate(-50%, -50%) scale(0);
    transition: transform 0.3s ease;
    border-radius: var(--radius-full);
  }

  .cyber-send-btn:hover:not(:disabled) .btn-glow {
    transform: translate(-50%, -50%) scale(1);
  }

  @media (max-width: 768px) {
    .cyber-tester-header {
      flex-direction: column;
      gap: 0.75rem;
      align-items: stretch;
      padding: 0.75rem;
    }

    .header-left {
      flex-direction: column;
      gap: 0.5rem;
      align-items: flex-start;
    }

    .endpoint-display {
      flex-direction: column;
      gap: 0.5rem;
      align-items: flex-start;
      width: 100%;
    }

    .endpoint-url {
      width: 100%;
    }

    .endpoint-url code {
      font-size: 0.75rem;
      word-break: break-all;
    }

    .cyber-send-btn {
      align-self: center;
      padding: 0.5rem 1.2rem;
    }
  }
</style>
