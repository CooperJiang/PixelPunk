<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { IconPickerProps, IconPickerEmits, IconOption } from './types'

  defineOptions({
    name: 'CyberIconPicker',
  })

  const props = withDefaults(defineProps<IconPickerProps>(), {
    modelValue: '',
    disabled: false,
    options: undefined,
  })

  const emit = defineEmits<IconPickerEmits>()
  const { $t } = useTexts()

  const showPicker = ref(false)

  const tIcon = (key: string) => $t(`components.iconPicker.${key}`)

  const iconOptions = computed<IconOption[]>(() => [
    { value: 'fas fa-star', label: tIcon('star') },
    { value: 'fas fa-bell', label: tIcon('bell') },
    { value: 'fas fa-bullhorn', label: tIcon('bullhorn') },
    { value: 'fas fa-info-circle', label: tIcon('infoCircle') },
    { value: 'fas fa-exclamation-circle', label: tIcon('exclamationCircle') },
    { value: 'fas fa-question-circle', label: tIcon('questionCircle') },
    { value: 'fas fa-check-circle', label: tIcon('checkCircle') },
    { value: 'fas fa-times-circle', label: tIcon('timesCircle') },
    { value: 'fas fa-fire', label: tIcon('fire') },
    { value: 'fas fa-bolt', label: tIcon('bolt') },
    { value: 'fas fa-gem', label: tIcon('gem') },
    { value: 'fas fa-crown', label: tIcon('crown') },
    { value: 'fas fa-heart', label: tIcon('heart') },
    { value: 'fas fa-diamond', label: tIcon('diamond') },
    { value: 'fas fa-snowflake', label: tIcon('snowflake') },
    { value: 'fas fa-sun', label: tIcon('sun') },
    { value: 'fas fa-moon', label: tIcon('moon') },
    { value: 'fas fa-cloud', label: tIcon('cloud') },
    { value: 'fas fa-gift', label: tIcon('gift') },
    { value: 'fas fa-trophy', label: tIcon('trophy') },
    { value: 'fas fa-medal', label: tIcon('medal') },
    { value: 'fas fa-flag', label: tIcon('flag') },
    { value: 'fas fa-key', label: tIcon('key') },
    { value: 'fas fa-lock', label: tIcon('lock') },
    { value: 'fas fa-shield-alt', label: tIcon('shieldAlt') },
    { value: 'fas fa-bookmark', label: tIcon('bookmark') },
    { value: 'fas fa-rocket', label: tIcon('rocket') },
    { value: 'fas fa-magic', label: tIcon('magic') },
    { value: 'fas fa-atom', label: tIcon('atom') },
    { value: 'fas fa-cpu', label: tIcon('cpu') },
    { value: 'fas fa-robot', label: tIcon('robot') },
    { value: 'fas fa-satellite', label: tIcon('satellite') },
    { value: 'fas fa-microchip', label: tIcon('microchip') },
    { value: 'fas fa-wifi', label: tIcon('wifi') },
    { value: 'fas fa-play', label: tIcon('play') },
    { value: 'fas fa-pause', label: tIcon('pause') },
    { value: 'fas fa-stop', label: tIcon('stop') },
    { value: 'fas fa-upload', label: tIcon('upload') },
    { value: 'fas fa-download', label: tIcon('download') },
    { value: 'fas fa-share', label: tIcon('share') },
    { value: 'fas fa-sync', label: tIcon('sync') },
    { value: 'fas fa-power-off', label: tIcon('powerOff') },
    { value: 'fas fa-smile', label: tIcon('smile') },
    { value: 'fas fa-laugh', label: tIcon('laugh') },
    { value: 'fas fa-angry', label: tIcon('angry') },
    { value: 'fas fa-surprise', label: tIcon('surprise') },
    { value: 'fas fa-thumbs-up', label: tIcon('thumbsUp') },
    { value: 'fas fa-thumbs-down', label: tIcon('thumbsDown') },
  ])

  const currentIconOptions = computed(() => props.options || iconOptions.value)

  const togglePicker = () => {
    if (!props.disabled) {
      showPicker.value = !showPicker.value
    }
  }

  const selectIcon = (iconValue: string) => {
    emit('update:modelValue', iconValue)
    showPicker.value = false
  }

  const hasSelectedIcon = computed(() => Boolean(props.modelValue))

  const getIconName = (iconValue: string) => {
    if (!iconValue) {
      return tIcon('selectPlaceholder')
    }

    const option = currentIconOptions.value.find((opt) => opt.value === iconValue)
    return option?.label || tIcon('custom')
  }

  const handleClickOutside = (event: MouseEvent) => {
    const target = event.target as Element
    if (!target.closest('.cyber-icon-picker')) {
      showPicker.value = false
    }
  }

  onMounted(() => {
    document.addEventListener('click', handleClickOutside)
  })

  onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
  })
</script>

<template>
  <div class="cyber-icon-picker">
    <div class="selected-icon" :class="{ focused: showPicker }" @click="togglePicker">
      <i v-if="hasSelectedIcon" :class="modelValue" class="icon-preview" />
      <span class="icon-text" :class="{ 'placeholder-text': !hasSelectedIcon }">{{ getIconName(modelValue) }}</span>
      <i class="fas fa-chevron-down dropdown-arrow" :class="{ open: showPicker }" />
    </div>

    <div v-if="showPicker" class="icon-panel">
      <div class="icon-grid">
        <div
          v-for="icon in currentIconOptions"
          :key="icon.value"
          class="icon-option"
          :class="{ active: modelValue === icon.value }"
          :title="icon.label"
          @click="selectIcon(icon.value)"
        >
          <i :class="icon.value" class="option-icon" />
          <span class="option-text">{{ icon.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .cyber-icon-picker {
    position: relative;
    width: 100%;
  }

  .selected-icon {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.2s ease;
    height: 32px;
  }

  .selected-icon:hover {
    border-color: var(--color-border-strong);
    background: rgba(var(--color-brand-500-rgb), 0.06);
  }

  .selected-icon:focus,
  .selected-icon.focused {
    border-color: var(--color-brand-400);
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .icon-preview {
    color: var(--color-brand-400);
    font-size: 1rem;
    min-width: 16px;
    text-align: center;
  }

  .icon-text {
    flex: 1;
    color: var(--color-content-default);
    font-size: 0.875rem;
  }

  .placeholder-text {
    color: var(--color-content-muted);
  }

  .dropdown-arrow {
    color: var(--color-content-muted);
    font-size: 0.75rem;
    transition: transform 0.2s ease;
  }

  .dropdown-arrow.open {
    transform: rotate(180deg);
  }

  .icon-panel {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    min-width: 140px;
    z-index: 1000;
    background: rgba(var(--color-background-900-rgb), 0.98);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    backdrop-filter: blur(15px);
    max-height: 250px;
    box-shadow: 0 12px 40px rgba(var(--color-background-900-rgb), 0.6);
    overflow: hidden;
  }

  .icon-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(64px, 1fr));
    gap: 2px;
    padding: 12px;
    max-height: 250px;
    overflow-y: auto;
  }

  .icon-option {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding: 12px 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    border-radius: var(--radius-sm);
    min-height: 64px;
    justify-content: center;
    border: 1px solid transparent;
  }

  .icon-option:hover {
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border-color: var(--color-border-default);
    transform: translateY(-1px);
  }

  .icon-option.active {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: var(--color-brand-400);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .option-icon {
    color: var(--color-brand-400);
    font-size: 1.25rem;
    transition: all 0.2s ease;
  }

  .icon-option:hover .option-icon {
    color: var(--color-brand-500);
    transform: scale(1.1);
  }

  .icon-option.active .option-icon {
    color: var(--color-brand-500);
  }

  .option-text {
    color: var(--color-content-muted);
    font-size: 0.6875rem;
    font-weight: 500;
    text-align: center;
    line-height: 1.2;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .icon-option:hover .option-text {
    color: var(--color-content-default);
  }

  .icon-option.active .option-text {
    color: var(--color-brand-400);
  }

  .icon-grid::-webkit-scrollbar {
    @apply w-1.5;
  }
  .icon-grid::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    @apply m-1 rounded-sm;
  }
  .icon-grid::-webkit-scrollbar-thumb {
    background: var(--color-brand-500);
    @apply rounded-sm;
    box-shadow: 0 0 4px rgba(var(--color-brand-500-rgb), 0.3);
  }
  .icon-grid::-webkit-scrollbar-thumb:hover {
    background: var(--color-brand-500);
    box-shadow: 0 0 6px rgba(var(--color-brand-500-rgb), 0.35);
  }

  @media (max-width: 768px) {
    .icon-panel {
      max-height: 220px;
    }

    .icon-grid {
      grid-template-columns: repeat(auto-fill, minmax(48px, 1fr));
      gap: 1px;
      padding: 8px;
      max-height: 220px;
    }

    .icon-option {
      padding: 8px 4px;
      min-height: 48px;
    }

    .option-icon {
      font-size: 1rem;
    }

    .option-text {
      font-size: 0.625rem;
    }
  }

  @media (max-width: 480px) {
    .icon-panel {
      max-height: 200px;
    }

    .icon-grid {
      grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
      padding: 6px;
      max-height: 200px;
    }

    .icon-option {
      padding: 6px 2px;
      min-height: 40px;
    }

    .option-icon {
      font-size: 0.875rem;
    }

    .option-text {
      font-size: 0.5rem;
    }
  }
</style>
