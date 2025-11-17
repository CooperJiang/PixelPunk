<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, reactive, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { ALL_SIZE_UNITS, QUICK_SIZE_PRESETS, SIZE_UNIT_OPTIONS } from '@/constants'
  import { useTexts } from '@/composables/useTexts'

  const toast = useToast()
  const { $t } = useTexts()

  const quickPresets = QUICK_SIZE_PRESETS
  const sizeUnitOptions = SIZE_UNIT_OPTIONS
  const allUnits = ALL_SIZE_UNITS

  const converter = reactive({
    value: 100,
    unit: 'MB',
  })

  const bytesValue = computed(() => {
    if (!converter.value || converter.value <= 0) {
      return 0
    }
    return unitToBytes(converter.value, converter.unit)
  })

  const unitToBytes = (value: number, unit: string): number => {
    const multipliers = {
      B: 1,
      KB: 1024,
      MB: 1024 * 1024,
      GB: 1024 * 1024 * 1024,
      TB: 1024 * 1024 * 1024 * 1024,
    }
    return Math.round(value * (multipliers[unit as keyof typeof multipliers] || 1))
  }

  const formatNumber = (num: number): string => num.toLocaleString(getCurrentLocale())

  const formatUnitValue = (bytes: number, divisor: number, precision: number): string => {
    if (bytes === 0) {
      return '0'
    }
    if (divisor === 1) {
      return formatNumber(bytes)
    }
    const value = bytes / divisor
    return value.toFixed(precision)
  }

  const setPreset = (preset: any) => {
    converter.value = preset.value
    converter.unit = preset.unit
  }

  const isPresetActive = (preset: any): boolean => converter.value === preset.value && converter.unit === preset.unit

  const bytesCopied = ref(false)

  const copyBytesValue = () => {
    copyToClipboard(bytesValue.value.toString(), $t('settings.sizeConverter.bytesLabel'))
  }

  const copyUnitValue = (unit: any) => {
    const value = formatUnitValue(bytesValue.value, unit.divisor, unit.precision)
    copyToClipboard(value, unit.label)
  }

  const copyAllValues = () => {
    const allValues = allUnits
      .map((unit) => {
        const value = formatUnitValue(bytesValue.value, unit.divisor, unit.precision)
        return `${unit.label}: ${value}`
      })
      .join('\n')

    copyToClipboard(allValues, $t('settings.sizeConverter.copyAllLabel'))
  }

  const copyToClipboard = async (text: string, label: string) => {
    const response = await navigator.clipboard.writeText(text)
    if (response) {
      toast.success(
        $t('settings.sizeConverter.copySuccess', {
          label,
        })
      )
    }
  }

  const onValueChange = () => {}

  const onUnitChange = () => {}
</script>

<template>
  <div class="size-converter-panel">
    <div class="panel-header">
      <div class="header-content">
        <div class="title-wrapper">
          <div class="title-icon">
            <i class="fas fa-exchange-alt" />
          </div>
          <div class="title-text">
            <h3 class="title">{{ $t('settings.sizeConverter.title') }}</h3>
            <p class="subtitle">{{ $t('settings.sizeConverter.subtitle') }}</p>
          </div>
        </div>
      </div>
    </div>

    <div class="panel-content">
      <div class="presets-section">
        <label class="section-label">{{ $t('settings.sizeConverter.presets') }}</label>
        <div class="presets-grid">
          <button
            v-for="preset in quickPresets"
            :key="preset.label"
            class="preset-btn"
            :class="{ 'preset-active': isPresetActive(preset) }"
            @click="setPreset(preset)"
          >
            {{ preset.label }}
          </button>
        </div>
      </div>

      <div class="conversion-section">
        <div class="flex items-center gap-4">
          <div class="input-group">
            <label class="input-label">{{ $t('settings.sizeConverter.valueLabel') }}</label>
            <CyberInput
              v-model.number="converter.value"
              type="number"
              :placeholder="$t('settings.sizeConverter.valueLabel')"
              min="0"
              step="0.01"
              class="w-32"
              @input="onValueChange"
            />
          </div>

          <div class="input-group">
            <label class="input-label">{{ $t('settings.sizeConverter.unitLabel') }}</label>
            <CyberDropdown v-model="converter.unit" :options="sizeUnitOptions" class="w-20" @change="onUnitChange" />
          </div>

          <div class="equals-sign">
            <i class="fas fa-equals text-content" />
          </div>

          <div class="result-group">
            <label class="input-label">{{ $t('settings.sizeConverter.bytesLabel') }}</label>
            <div class="result-display">
              <span class="result-number">{{ formatNumber(bytesValue) }}</span>
              <span class="result-unit">{{ $t('settings.sizeConverter.bytes') }}</span>
              <button class="copy-btn" :title="$t('settings.sizeConverter.copyBytes')" @click="copyBytesValue">
                <i class="fas" :class="bytesCopied ? 'fa-check' : 'fa-copy'" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="table-section">
        <div class="table-header">
          <span class="table-title">{{ $t('settings.sizeConverter.detailTitle') }}</span>
          <button class="copy-all-btn" :title="$t('settings.sizeConverter.copyAll')" @click="copyAllValues">
            <i class="fas fa-clipboard mr-1" />
            {{ $t('settings.sizeConverter.copyAll') }}
          </button>
        </div>
        <div class="table-grid">
          <div v-for="unit in allUnits" :key="unit.key" class="table-item">
            <div class="unit-label">{{ unit.label }}</div>
            <div class="unit-value">
              {{ formatUnitValue(bytesValue, unit.divisor, unit.precision) }}
              <button
                class="unit-copy-btn"
                :title="$t('settings.sizeConverter.copyUnit', { label: unit.label })"
                @click="copyUnitValue(unit)"
              >
                <i class="fas fa-copy" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="info-section">
        <div class="info-grid">
          <div class="info-item">
            <i class="fas fa-info-circle info-icon" />
            <span>{{ $t('settings.sizeConverter.info.usage') }}</span>
          </div>
          <div class="info-item">
            <i class="fas fa-lightbulb info-icon" />
            <span>{{ $t('settings.sizeConverter.info.calculation') }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .size-converter-panel {
    @apply space-y-6;
  }

  .panel-header {
    @apply relative overflow-hidden;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-800-rgb), 0.9), rgba(var(--color-background-900-rgb), 0.95));
    border: 1px solid var(--color-border-subtle);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
  }

  .header-content {
    @apply relative p-6;
    z-index: 1;
  }

  .title-wrapper {
    @apply flex items-start gap-4;
  }

  .title-icon {
    @apply flex h-12 w-12 items-center justify-center;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.18), rgba(var(--color-brand-500-rgb), 0.08));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.24);
    color: var(--color-brand-400);
    font-size: 18px;
    box-shadow:
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.18),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
  }

  .title-text {
    @apply flex flex-col gap-1;
  }

  .title {
    @apply text-xl font-semibold;
    color: var(--color-content-heading);
    margin: 0;
  }

  .subtitle {
    @apply text-sm leading-relaxed;
    color: var(--color-content-muted);
    margin: 0;
  }

  .panel-content {
    @apply space-y-6;
  }

  .presets-section {
    @apply space-y-3;
  }

  .section-label {
    @apply block text-sm font-medium;
    color: var(--color-content);
  }

  .presets-grid {
    @apply flex flex-wrap gap-2;
  }

  .preset-btn {
    @apply px-4 py-2 text-sm font-medium transition-all duration-200;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-500-rgb), 0.06));
    color: var(--color-brand-400);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.22);
  }

  .preset-btn:hover {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.1));
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateY(-1px);
  }

  .preset-btn.preset-active {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.85), rgba(var(--color-brand-500-rgb), 0.72));
    color: var(--color-content-inverse);
    border-color: var(--color-brand-400);
  }

  .conversion-section {
    @apply p-6;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.9));
    border: 1px solid var(--color-border-subtle);
    box-shadow:
      0 4px 16px rgba(0, 0, 0, 0.08),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
  }

  .input-group {
    @apply flex flex-col gap-2;
  }

  .input-label {
    @apply text-sm font-medium;
    color: var(--color-content);
  }

  .equals-sign {
    @apply flex items-center justify-center px-3;
    height: 42px;
    font-size: 18px;
    color: var(--color-brand-400);
  }

  .result-group {
    @apply flex-1;
    min-width: 200px;
  }

  .result-display {
    @apply flex items-center px-4 py-2;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-500-rgb), 0.06));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.22);
    height: 42px;
    font-family: ui-monospace, SFMono-Regular, 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace;
  }

  .result-number {
    @apply flex-1 text-lg font-semibold;
    color: var(--color-brand-400);
  }

  .result-unit {
    @apply ml-2 text-sm;
    color: var(--color-content-muted);
  }

  .copy-btn {
    @apply ml-2 p-1 transition-all duration-200;
    color: var(--color-content-muted);
    opacity: 0.6;
  }

  .copy-btn:hover {
    opacity: 1;
    color: var(--color-brand-400);
  }

  .table-section {
    @apply overflow-hidden;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.5), rgba(var(--color-background-800-rgb), 0.7));
    border: 1px solid var(--color-border-subtle);
  }

  .table-header {
    @apply flex items-center justify-between px-6 py-4;
    background: linear-gradient(135deg, rgba(var(--color-background-600-rgb), 0.3), rgba(var(--color-background-700-rgb), 0.5));
    border-bottom: 1px solid var(--color-border-subtle);
  }

  .table-title {
    @apply text-base font-medium;
    color: var(--color-content-heading);
  }

  .copy-all-btn {
    @apply text-sm font-medium transition-opacity duration-200;
    color: var(--color-brand-400);
    opacity: 0.8;
  }

  .copy-all-btn:hover {
    opacity: 1;
  }

  .table-grid {
    @apply grid gap-px;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    background: rgba(var(--color-background-700-rgb), 0.35);
  }

  .table-item {
    @apply flex items-center justify-between px-6 py-4;
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.6), rgba(var(--color-background-800-rgb), 0.8));
  }

  .unit-label {
    @apply text-sm font-medium;
    color: var(--color-content);
  }

  .unit-value {
    @apply flex items-center gap-2 font-semibold;
    color: var(--color-success-500);
    font-family: ui-monospace, SFMono-Regular, 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace;
  }

  .unit-copy-btn {
    @apply text-xs transition-opacity duration-200;
    color: var(--color-content-muted);
    opacity: 0;
  }

  .table-item:hover .unit-copy-btn {
    opacity: 0.6;
  }

  .unit-copy-btn:hover {
    opacity: 1 !important;
    color: var(--color-brand-400);
  }

  .info-section {
    @apply pt-6;
    border-top: 1px solid var(--color-border-subtle);
  }

  .info-grid {
    @apply flex flex-wrap gap-6;
  }

  .info-item {
    @apply flex items-center gap-3 text-sm;
    color: var(--color-content-muted);
    margin-bottom: 20px;
    min-width: 0;
  }

  .info-icon {
    @apply flex-shrink-0 text-base;
    color: var(--color-brand-400);
  }

  .info-item span {
    line-height: 1.6;
  }

  @media (max-width: 768px) {
    .conversion-section .flex {
      @apply flex-col gap-4;
    }

    .equals-sign {
      @apply justify-center py-2;
      height: auto;
    }

    .table-grid {
      grid-template-columns: 1fr;
    }

    .info-grid {
      @apply flex-col gap-4;
    }

    .info-item {
      flex: 1 1 auto;
    }

    .presets-grid {
      @apply grid grid-cols-2 gap-2;
    }
  }
</style>
