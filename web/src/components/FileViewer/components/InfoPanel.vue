<script setup lang="ts">
  import { onMounted, ref, watch } from 'vue'
  import type { FileInfo } from '../types'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    currentFile: FileInfo | null
    isVisible: boolean
    showKeyboardTips: boolean
    shouldUseFillMode: boolean
    isLightBackground: boolean
  }

  const props = defineProps<Props>()
  const { $t } = useTexts()

  /* 快捷键提示显示控制 */
  const showKeyboardTipsVisible = ref(false)
  let keyboardTipsTimer: number | null = null

  /* 监听isVisible变化，控制快捷键提示的显示 */
  watch(
    () => props.isVisible,
    (newValue) => {
      if (newValue && props.showKeyboardTips) {
        showKeyboardTipsVisible.value = true

        if (keyboardTipsTimer) {
          clearTimeout(keyboardTipsTimer)
        }

        keyboardTipsTimer = setTimeout(() => {
          showKeyboardTipsVisible.value = false
        }, 5000)
      } else {
        showKeyboardTipsVisible.value = false
        if (keyboardTipsTimer) {
          clearTimeout(keyboardTipsTimer)
          keyboardTipsTimer = null
        }
      }
    }
  )

  onMounted(() => {
    if (props.isVisible && props.showKeyboardTips) {
      showKeyboardTipsVisible.value = true
      keyboardTipsTimer = setTimeout(() => {
        showKeyboardTipsVisible.value = false
      }, 5000)
    }
  })
</script>

<template>
  <div class="info-panel">
    <div
      v-if="showKeyboardTips && isVisible && showKeyboardTipsVisible"
      class="keyboard-tips"
      :class="{ 'is-light-bg': props.isLightBackground }"
    >
      <div class="tips-header">
        <i class="fas fa-keyboard" />
        <span>{{ $t('components.infoPanel.shortcuts.title') }}</span>
      </div>
      <div class="tips-content">
        <div class="tip-item">
          <kbd>←→</kbd>
          <span>{{ $t('components.infoPanel.shortcuts.switchFile') }}</span>
        </div>
        <div class="tip-item">
          <kbd>{{ $t('components.infoPanel.shortcuts.space') }}</kbd>
          <span>{{
            $t(shouldUseFillMode ? 'components.infoPanel.shortcuts.fitMode' : 'components.infoPanel.shortcuts.fillMode')
          }}</span>
        </div>
        <div class="tip-item">
          <kbd>±</kbd>
          <span>{{ $t('components.infoPanel.shortcuts.zoom') }}</span>
        </div>
        <div class="tip-item">
          <kbd>R/L</kbd>
          <span>{{ $t('components.infoPanel.shortcuts.rotate') }}</span>
        </div>
        <div class="tip-item">
          <kbd>F</kbd>
          <span>{{ $t('components.infoPanel.shortcuts.fullscreen') }}</span>
        </div>
        <div class="tip-item">
          <kbd>ESC</kbd>
          <span>{{ $t('components.infoPanel.shortcuts.exit') }}</span>
        </div>
      </div>
    </div>

    <div
      v-if="(currentFile?.ai_info?.description || currentFile?.description) && isVisible"
      class="image-description"
      :class="{ 'is-light-bg': props.isLightBackground }"
    >
      <div class="description-content">
        <div class="description-text">
          {{ currentFile?.ai_info?.description || currentFile.description }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .info-panel {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 50;
    pointer-events: none;
  }

  .keyboard-tips {
    position: absolute;
    bottom: 30px;
    left: 30px;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.92) 50%,
      rgba(var(--color-background-800-rgb), 0.9) 100%
    );
    backdrop-filter: blur(25px);
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.6);
    border-radius: var(--radius-lg);
    padding: 12px 14px;
    color: var(--color-content-default);
    animation: slideInLeft 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    max-width: 240px;
    pointer-events: auto;
    box-shadow:
      0 8px 32px rgba(var(--color-background-900-rgb), 0.6),
      0 4px 16px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.2),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.1);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .tips-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 8px;
    font-size: 12px;
    font-weight: 600;
    color: var(--color-brand-500);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.5);
    padding-bottom: 6px;
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.6);
    letter-spacing: 0.3px;
  }

  .tips-content {
    display: grid;
    grid-template-columns: 1fr;
    gap: 6px;
  }

  .tip-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    font-size: 11px;
    color: var(--color-content-muted);
    font-weight: 500;
    text-shadow: 0 0 4px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .tip-item kbd {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25) 0%, rgba(var(--color-brand-500-rgb), 0.15) 100%);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.6);
    border-radius: var(--radius-sm);
    padding: 3px 8px;
    font-size: 10px;
    font-family: inherit;
    color: var(--color-content-default);
    min-width: 28px;
    text-align: center;
    font-weight: 600;
    backdrop-filter: blur(8px);
    box-shadow:
      0 2px 8px rgba(var(--color-background-900-rgb), 0.3),
      0 0 4px rgba(var(--color-brand-500-rgb), 0.4),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.2);
    text-shadow: 0 0 4px rgba(var(--color-brand-500-rgb), 0.5);
    letter-spacing: 0.2px;
  }

  .keyboard-tips.is-light-bg {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    backdrop-filter: blur(25px) saturate(110%);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow:
      0 12px 32px rgba(var(--color-background-900-rgb), 0.6),
      0 6px 16px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.08),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.15);
  }

  .keyboard-tips.is-light-bg .tips-header {
    color: rgba(255, 255, 255, 0.95);
    text-shadow:
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 6px rgba(255, 255, 255, 0.3),
      2px 2px 4px rgba(var(--color-background-900-rgb), 0.8);
    border-bottom-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .keyboard-tips.is-light-bg .tip-item {
    color: rgba(255, 255, 255, 0.95);
    text-shadow:
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 6px rgba(255, 255, 255, 0.3),
      2px 2px 4px rgba(var(--color-background-900-rgb), 0.8);
  }

  .keyboard-tips.is-light-bg .tip-item kbd {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    color: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(25px) saturate(110%);
    box-shadow:
      0 12px 32px rgba(var(--color-background-900-rgb), 0.6),
      0 6px 16px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.08),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.15);
    text-shadow:
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 6px rgba(255, 255, 255, 0.3),
      2px 2px 4px rgba(var(--color-background-900-rgb), 0.8);
  }

  .image-description {
    position: absolute;
    bottom: 100px;
    left: 50%;
    transform: translateX(-50%);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.92) 0%,
      rgba(var(--color-background-800-rgb), 0.9) 50%,
      rgba(var(--color-background-800-rgb), 0.86) 100%
    );
    backdrop-filter: blur(25px);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.35);
    border-radius: var(--radius-sm);
    padding: 0;
    color: rgba(var(--color-white-rgb), 0.9);
    animation: slideInUp 0.5s cubic-bezier(0.16, 1, 0.3, 1);
    max-width: min(85vw, 800px);
    min-width: 350px;
    pointer-events: auto;
    box-shadow:
      0 8px 32px rgba(var(--color-background-900-rgb), 0.45),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.2),
      0 0 0 1px rgba(255, 255, 255, 0.05),
      inset 0 1px 0 rgba(255, 255, 255, 0.06);
    text-shadow: 0 0 calc(8px * var(--text-shadow-intensity)) rgba(var(--color-brand-500-rgb), 0.3);
    overflow: hidden;
  }

  .description-content {
    position: relative;
    padding: 18px 24px;
  }

  .description-text {
    font-size: 15px;
    line-height: 1.8;
    color: rgba(var(--color-white-rgb), 0.95);
    text-align: center;
    word-wrap: break-word;
    position: relative;
    font-weight: 500;
    letter-spacing: 0.4px;
    text-shadow:
      0 0 calc(12px * var(--text-shadow-intensity, 1)) rgba(var(--color-brand-500-rgb), 0.35),
      0 0 calc(6px * var(--text-shadow-intensity, 1)) rgba(255, 255, 255, 0.6),
      2px 2px calc(4px * var(--text-shadow-intensity, 1)) rgba(var(--color-background-900-rgb), 0.7);
    filter: drop-shadow(0 0 calc(6px * var(--text-shadow-intensity, 1)) rgba(148, 163, 184, 0.3));
  }

  .description-text::after {
    content: '"';
    position: absolute;
    right: -12px;
    bottom: -16px;
    font-size: 28px;
    color: rgba(var(--color-brand-500-rgb), 0.45);
    font-family: Georgia, serif;
    line-height: 1;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.35);
  }

  .image-description.is-light-bg {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.75) 0%,
      rgba(var(--color-background-900-rgb), 0.7) 50%,
      rgba(var(--color-background-900-rgb), 0.78) 100%
    );
    backdrop-filter: blur(25px) saturate(110%);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow:
      0 12px 32px rgba(var(--color-background-900-rgb), 0.6),
      0 6px 16px rgba(var(--color-background-900-rgb), 0.4),
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.08),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.15);
  }

  .image-description.is-light-bg .description-text {
    color: rgba(255, 255, 255, 0.95);
    text-shadow:
      0 0 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 6px rgba(255, 255, 255, 0.3),
      2px 2px 4px rgba(var(--color-background-900-rgb), 0.8);
  }

  .image-description.is-light-bg .description-text::before,
  .image-description.is-light-bg .description-text::after {
    color: rgba(var(--color-brand-500-rgb), 0.6);
    text-shadow:
      0 0 8px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 4px rgba(255, 255, 255, 0.2);
  }

  @keyframes slideInLeft {
    from {
      opacity: 0;
      transform: translateX(-20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  @keyframes slideInUp {
    from {
      opacity: 0;
      transform: translateX(-50%) translateY(30px) scale(0.95);
    }
    to {
      opacity: 1;
      transform: translateX(-50%) translateY(0) scale(1);
    }
  }

  @media (max-width: 1024px) {
    .keyboard-tips {
      max-width: 250px;
    }

    .image-description {
      max-width: min(85vw, 500px);
      min-width: 260px;
    }
  }

  @media (max-width: 768px) {
    .keyboard-tips {
      left: 20px;
      bottom: 25px;
      max-width: 200px;
      padding: 10px 12px;
    }

    .image-description {
      bottom: 90px;
      max-width: min(90vw, 450px);
      min-width: 280px;
    }

    .description-content {
      padding: 16px 20px;
    }

    .description-text {
      font-size: 13px;
      line-height: 1.6;
    }

    .tips-header {
      font-size: 11px;
      margin-bottom: 6px;
    }

    .tip-item {
      font-size: 10px;
    }

    .tip-item kbd {
      padding: 2px 5px;
      font-size: 9px;
      min-width: 24px;
    }
  }

  @media (max-width: 480px) {
    .keyboard-tips {
      left: 15px;
      right: auto;
      max-width: none;
      bottom: 20px;
      width: calc(100vw - 30px);
      transform: none;
      padding: 8px 10px;
    }

    .image-description {
      bottom: 140px;
      left: 15px;
      right: 15px;
      transform: none;
      max-width: none;
      min-width: auto;
      width: calc(100vw - 30px);
    }

    .description-content {
      padding: 14px 18px;
    }

    .description-text {
      font-size: 12px;
      line-height: 1.6;
    }

    .description-text::before,
    .description-text::after {
      font-size: 18px;
    }
  }
</style>
