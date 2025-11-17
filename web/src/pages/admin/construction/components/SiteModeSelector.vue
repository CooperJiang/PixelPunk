<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  interface ModeOption {
    value: string
    title: string
    icon: string
    description: string
    features: string[]
    disabled?: boolean
  }

  const props = defineProps<{
    modelValue: string
  }>()

  const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
  }>()

  const { $t } = useTexts()

  const selectedMode = ref(props.modelValue || 'website')

  const modeOptions = computed<ModeOption[]>(() => [
    {
      value: 'website',
      title: $t('admin.construction.siteMode.website.title'),
      icon: 'fas fa-globe',
      description: $t('admin.construction.siteMode.website.description'),
      features: [
        $t('admin.construction.siteMode.website.features.0'),
        $t('admin.construction.siteMode.website.features.1'),
        $t('admin.construction.siteMode.website.features.2'),
        $t('admin.construction.siteMode.website.features.3'),
      ],
      disabled: false,
    },
  ])

  const selectMode = (mode: string) => {
    selectedMode.value = mode
    emit('update:modelValue', mode)
  }

  onMounted(() => {
    if (!props.modelValue || props.modelValue !== 'website') {
      emit('update:modelValue', 'website')
    }
  })

  watch(
    () => props.modelValue,
    (newValue) => {
      if (newValue) {
        selectedMode.value = newValue
      }
    }
  )

  watch(selectedMode, (newValue) => {
    emit('update:modelValue', newValue)
  })
</script>

<template>
  <div class="site-mode-selector">
    <div class="mode-options-list">
      <div
        v-for="mode in modeOptions"
        :key="mode.value"
        class="mode-option"
        :class="{
          active: selectedMode === mode.value,
          disabled: mode.disabled,
        }"
        @click="!mode.disabled && selectMode(mode.value)"
      >
        <div class="mode-radio">
          <input
            :id="`mode-${mode.value}`"
            v-model="selectedMode"
            type="radio"
            :value="mode.value"
            :disabled="mode.disabled"
            class="radio-input"
          />
          <label :for="`mode-${mode.value}`" class="radio-label">
            <div class="radio-circle">
              <div class="radio-dot" />
            </div>
          </label>
        </div>

        <div class="mode-content">
          <div class="mode-title-row">
            <i :class="mode.icon" class="mode-icon" />
            <span class="mode-title">{{ mode.title }}</span>
            <span v-if="mode.disabled" class="disabled-badge">{{ $t('admin.construction.siteMode.comingSoon') }}</span>
          </div>
          <div class="mode-description">{{ mode.description }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .site-mode-selector {
    width: 100%;
  }

  .mode-options-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }

  .mode-option {
    display: flex;
    align-items: flex-start;
    padding: var(--space-md);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-lg);
    background: rgba(var(--color-background-800-rgb), 0.3);
    cursor: pointer;
    transition: all var(--transition-normal) var(--ease-out);
  }

  .mode-option:hover:not(.disabled) {
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    background: rgba(var(--color-background-800-rgb), 0.5);
  }

  .mode-option.active {
    border-color: rgba(var(--color-brand-500-rgb), 0.8);
    background: rgba(var(--color-background-800-rgb), 0.6);
    box-shadow: var(--shadow-glow-md);
  }

  .mode-option.disabled {
    opacity: var(--opacity-disabled);
    cursor: not-allowed;
    border-color: rgba(128, 128, 128, 0.2);
    background: rgba(var(--color-background-800-rgb), 0.2);
  }

  .mode-option.disabled .mode-icon,
  .mode-option.disabled .mode-title {
    color: rgba(255, 255, 255, 0.5);
  }

  .mode-radio {
    margin-right: var(--space-sm);
    margin-top: var(--space-xs);
  }

  .radio-input {
    display: none;
  }

  .radio-label {
    cursor: pointer;
  }

  .mode-option.disabled .radio-label {
    cursor: not-allowed;
  }

  .radio-circle {
    width: var(--space-md);
    height: var(--space-md);
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.5);
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all var(--transition-normal) var(--ease-out);
  }

  .mode-option.active .radio-circle {
    border-color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .mode-option.disabled .radio-circle {
    border-color: rgba(128, 128, 128, 0.3);
  }

  .radio-dot {
    width: var(--space-xs);
    height: var(--space-xs);
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    transform: scale(0);
    transition: transform var(--transition-fast) var(--ease-out);
  }

  .mode-option.active .radio-dot {
    transform: scale(1);
  }

  .mode-option.disabled .radio-dot {
    background: rgba(128, 128, 128, 0.5);
  }

  .mode-content {
    flex: 1;
  }

  .mode-title-row {
    display: flex;
    align-items: center;
    margin-bottom: var(--space-sm);
  }

  .mode-icon {
    color: var(--color-brand-500);
    font-size: var(--text-base);
    margin-right: var(--space-sm);
    width: var(--space-md);
    text-align: center;
  }

  .mode-title {
    font-size: var(--text-sm);
    font-weight: var(--font-semibold);
    color: var(--color-content-heading);
    margin-right: var(--space-sm);
  }

  .disabled-badge {
    display: inline-block;
    padding: var(--space-xs) var(--space-sm);
    background: rgba(255, 152, 0, 0.2);
    color: var(--color-warning-400);
    font-size: 0.7rem;
    font-weight: var(--font-medium);
    border-radius: var(--radius-full);
    border: 1px solid rgba(255, 152, 0, 0.3);
  }

  .mode-description {
    color: rgba(255, 255, 255, 0.7);
    font-size: var(--text-xs);
    line-height: 1.4;
    margin: 0;
  }

  @media (max-width: 768px) {
    .mode-option {
      padding: var(--space-sm);
    }

    .mode-title {
      font-size: var(--text-sm);
    }

    .mode-description {
      font-size: 0.75rem;
    }

    .disabled-badge {
      font-size: 0.65rem;
      padding: var(--space-xs) var(--space-sm);
    }
  }
</style>
