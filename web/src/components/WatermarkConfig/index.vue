<script setup lang="ts">
  import { ref, watch } from 'vue'
  import { DEFAULT_WATERMARK_CONFIG, type WatermarkConfig } from './types'
  import WatermarkPanel from './WatermarkPanel.vue'
  import WatermarkPreview from './WatermarkPreview.vue'
  import { type CanvasWatermarkConfig, generateWatermark } from '@/utils/file/canvasWatermarkGenerator'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'WatermarkConfig',
  })

  const { $t } = useTexts()

  const toast = useToast()

  const props = withDefaults(defineProps<Props>(), {
    config: () => ({ ...DEFAULT_WATERMARK_CONFIG }),
    previewImage: undefined,
  })

  const emit = defineEmits<Emits>()

  interface Props {
    visible: boolean
    config?: WatermarkConfig
    previewImage?: string
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'confirm', config: WatermarkConfig): void
    (e: 'cancel'): void
  }

  const localVisible = ref(props.visible)
  const localConfig = ref<WatermarkConfig>({ ...props.config, enabled: true })

  watch(
    () => props.visible,
    (newValue) => {
      localVisible.value = newValue
      if (newValue) {
        localConfig.value = { ...props.config, enabled: true }
      }
    }
  )

  watch(localVisible, (newValue) => {
    emit('update:visible', newValue)
  })

  watch(
    () => props.config,
    (newConfig) => {
      localConfig.value = { ...newConfig }
    },
    { deep: true }
  )

  const handleConfigChange = (newConfig: WatermarkConfig) => {
    localConfig.value = { ...newConfig }
  }

  const handleConfirm = async () => {
    const finalConfig = { ...localConfig.value }

    let pureBase64 = ''

    if (finalConfig.type === 'text' && finalConfig.text && finalConfig.text.trim()) {
      try {
        const canvasConfig: CanvasWatermarkConfig = {
          ...finalConfig,
          targetWidth: 1000,
          targetHeight: 1000,
        }

        const result = await generateWatermark(canvasConfig)

        if (result.success && result.base64Data) {
          pureBase64 = result.base64Data.replace(/^data:image\/\w+;base64,/, '')
          finalConfig.fileBase64 = pureBase64
        } else {
          const errorMsg = result.error || $t('components.watermark.errors.generateFailed')
          toast.error(errorMsg)
          return
        }
      } catch (error) {
        const errorMsg = error instanceof Error ? error.message : $t('components.watermark.errors.generateException')
        toast.error(errorMsg)
        return
      }
    } else if (finalConfig.type === 'image' && finalConfig.fileBase64) {
      pureBase64 = finalConfig.fileBase64
    }

    emit('confirm', finalConfig)
    localVisible.value = false
  }

  const handleClose = () => {
    emit('cancel')
    localVisible.value = false
  }

  const clearPreviewImage = () => {}
</script>

<template>
  <cyberDialog
    v-model="localVisible"
    :title="$t('components.watermark.title')"
    width="1400px"
    height="800px"
    :close-on-click-overlay="false"
    @close="handleClose"
  >
    <div class="watermark-config-main">
      <div class="config-content">
        <div class="config-sidebar">
          <div class="sidebar-header">
            <h3 class="sidebar-title">
              <i class="fas fa-sliders-h mr-2" />
              {{ $t('components.watermark.configOptions') }}
            </h3>
          </div>
          <div class="sidebar-content">
            <WatermarkPanel :config="localConfig" @config-change="handleConfigChange" />
          </div>
        </div>

        <div class="preview-main">
          <div class="preview-header">
            <div class="preview-title-row">
              <h3 class="preview-title">
                <i class="fas fa-eye mr-2" />
                {{ $t('components.watermark.livePreview') }}
              </h3>
              <cyberTooltip placement="bottom" :max-width="350">
                <button class="help-btn">
                  <i class="fas fa-question-circle text-sm" />
                  <span class="help-text">{{ $t('components.watermark.help.title') }}</span>
                </button>
                <template #content>
                  <div class="help-tooltip-content">
                    <div class="help-section">
                      <div class="help-title">
                        <i class="fas fa-map-pin mr-1" />
                        {{ $t('components.watermark.help.positionTitle') }}
                      </div>
                      <div class="help-desc">
                        <span class="help-highlight">{{ $t('components.watermark.help.positionFormula') }}</span>
                        <br />
                        {{ $t('components.watermark.help.positionDesc') }}
                      </div>
                    </div>
                    <div class="help-section">
                      <div class="help-title">
                        <i class="fas fa-keyboard mr-1" />
                        {{ $t('components.watermark.help.shortcutsTitle') }}
                      </div>
                      <div class="help-desc">
                        <kbd>{{ $t('components.watermark.help.spaceKey') }}</kbd> + {{ $t('components.watermark.help.dragCanvas') }}<br />
                        <kbd>{{ $t('components.watermark.help.scrollWheel') }}</kbd
                        >：{{ $t('components.watermark.help.zoomCanvas') }}<br />
                        {{ $t('components.watermark.help.dragWatermark') }}<br />
                        {{ $t('components.watermark.help.dragCorners') }}
                      </div>
                    </div>
                  </div>
                </template>
              </cyberTooltip>
            </div>
            <div class="preview-actions">
              <CyberButton size="sm" type="secondary" @click="clearPreviewImage">
                <i class="fas fa-times mr-1.5" />
                {{ $t('components.watermark.clearPreview') }}
              </CyberButton>
            </div>
          </div>
          <div class="preview-content">
            <WatermarkPreview :config="localConfig" :preview-image="previewImage" @config-change="handleConfigChange" />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="config-footer">
        <div class="footer-hints">
          <div class="hint-item">
            <i class="fas fa-hand-pointer" />
            <span>{{ $t('components.watermark.footer.hints.dragWatermark') }}</span>
          </div>
          <div class="hint-item">
            <i class="fas fa-arrows-alt" />
            <span>{{ $t('components.watermark.footer.hints.dragCorners') }}</span>
          </div>
          <div class="hint-item">
            <i class="fas fa-upload" />
            <span>{{ $t('components.watermark.footer.hints.dropImage') }}</span>
          </div>
        </div>
        <div class="footer-actions">
          <CyberButton type="secondary" @click="handleClose">
            <i class="fas fa-times mr-1.5" />
            {{ $t('actions.cancel') }}
          </CyberButton>
          <CyberButton type="primary" @click="handleConfirm">
            <i class="fas fa-check mr-1.5" />
            {{ $t('components.watermark.footer.apply') }}
          </CyberButton>
        </div>
      </div>
    </template>
  </cyberDialog>
</template>

<style scoped>
  .watermark-config-main {
    @apply flex flex-col;
    height: 650px;
    overflow: hidden;
  }

  .config-content {
    @apply flex;
    height: 650px;
    overflow: hidden;
  }

  .config-sidebar {
    @apply flex w-72 flex-none flex-col border-r;
    background: var(--color-background-800);
    border-color: var(--color-border-subtle);
    height: 650px;
  }

  .sidebar-header {
    @apply flex-none border-b p-3;
    background: var(--color-background-700);
    border-color: var(--color-border-subtle);
    height: 50px;
  }

  .sidebar-title {
    @apply flex items-center text-sm font-semibold;
    color: var(--color-content-heading);
  }

  .sidebar-content {
    @apply flex-1 overflow-y-auto p-4;
    height: calc(650px - 50px);
  }

  .preview-main {
    @apply flex flex-1 flex-col overflow-hidden;
    background: var(--color-background-900);
    height: 650px;
  }

  .preview-header {
    @apply flex items-center justify-between gap-3 border-b p-3;
    background: var(--color-background-700);
    border-color: var(--color-border-subtle);
    height: 50px;
  }

  .preview-actions {
    @apply flex-shrink-0;
  }

  .preview-title-row {
    @apply flex items-center gap-3;
  }

  .preview-title {
    @apply flex items-center text-sm font-semibold;
    color: var(--color-content-heading);
  }

  .help-btn {
    @apply flex items-center gap-1.5 rounded px-2.5 py-1 text-xs font-medium transition-all duration-200;
    background: var(--color-background-800);
    color: var(--color-content-subtle);
    border: 1px solid var(--color-border-default);
    cursor: pointer;
  }

  .help-btn:hover {
    background: var(--color-hover-bg);
    color: var(--color-brand-500);
    border-color: var(--color-brand-500);
  }

  .help-btn i {
    color: var(--color-brand-400);
  }

  .help-text {
    @apply select-none;
  }

  .help-tooltip-content {
    @apply space-y-3;
  }

  .help-section {
    @apply space-y-1.5;
  }

  .help-title {
    @apply flex items-center text-sm font-semibold;
    color: var(--color-brand-400);
  }

  .help-desc {
    @apply text-xs leading-relaxed;
    color: var(--color-content-subtle);
  }

  .help-desc kbd {
    @apply rounded border px-1.5 py-0.5 font-mono text-xs;
    border-color: var(--color-border-default);
    background: var(--color-background-700);
    color: var(--color-content-default);
  }

  .help-highlight {
    @apply font-semibold;
    color: var(--color-brand-400);
  }

  .preview-content {
    @apply flex items-center justify-center overflow-hidden p-4;
    height: calc(650px - 50px);
  }

  .config-footer {
    @apply flex items-center justify-between gap-4 border-t p-4;
    background: var(--color-background-700);
    border-color: var(--color-border-subtle);
  }

  .footer-hints {
    @apply flex flex-1 flex-wrap items-center gap-4;
  }

  .footer-hints .hint-item {
    @apply flex items-center gap-2 text-xs;
    color: var(--color-content-secondary);
  }

  .footer-hints .hint-item i {
    color: var(--color-brand-500);
    font-size: 0.75rem;
  }

  .footer-actions {
    @apply flex flex-shrink-0 items-center space-x-3;
  }

  .sidebar-content::-webkit-scrollbar {
    width: 6px;
  }

  .sidebar-content::-webkit-scrollbar-track {
    background: var(--color-background-800);
  }

  .sidebar-content::-webkit-scrollbar-thumb {
    @apply rounded-full;
    background: var(--color-brand-500);
  }

  .sidebar-content::-webkit-scrollbar-thumb:hover {
    background: var(--color-brand-500);
    opacity: 0.8;
  }

  .sidebar-header,
  .preview-header {
    backdrop-filter: blur(8px);
  }

  @media (max-width: 1440px) {
    .config-sidebar {
      @apply w-64;
    }

    .preview-hints-group {
      @apply gap-1.5;
    }

    .usage-hint-inline,
    .control-hint-inline {
      @apply px-1.5 py-0.5 text-xs;
    }

    .control-hint-inline kbd {
      @apply px-1 py-0.5;
      font-size: 0.65rem;
    }
  }

  @media (max-width: 1200px) {
    .config-sidebar {
      @apply w-56;
    }

    .preview-hints-group {
      @apply gap-1;
    }

    .usage-hint-inline,
    .control-hint-inline {
      font-size: 0.7rem;
      @apply px-1.5 py-0.5;
    }

    .control-hint-inline kbd {
      font-size: 0.6rem;
    }

    .footer-hints {
      @apply gap-2;
    }

    .footer-hints .hint-item {
      @apply text-[0.7rem];
    }
  }

  @media (max-width: 1024px) {
    .watermark-config-main,
    .config-content,
    .config-sidebar,
    .preview-main {
      height: auto;
    }

    .config-content {
      @apply flex-col;
      max-height: calc(100vh - 200px);
    }

    .config-sidebar {
      @apply w-full border-b border-r-0;
      max-height: 300px;
    }

    .sidebar-content {
      max-height: 250px;
    }

    .preview-main {
      flex: 1;
      min-height: 400px;
    }

    .preview-hints-group {
      @apply flex-wrap gap-1;
    }

    .control-hint-inline {
      @apply hidden;
    }

    .config-footer {
      @apply flex-col items-stretch gap-3;
    }

    .footer-hints {
      @apply justify-center;
    }

    .footer-actions {
      @apply w-full justify-end;
    }
  }

  @media (max-width: 768px) {
    .footer-hints {
      @apply flex-col items-start gap-2;
    }

    .footer-actions {
      @apply w-full flex-col-reverse;
    }

    .footer-actions :deep(.cyber-btn) {
      @apply w-full;
    }
  }
</style>
