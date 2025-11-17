<script setup lang="ts">
  import { defineEmits, defineProps, onBeforeUnmount, onMounted, ref, watchEffect } from 'vue'
  import { StorageUtil } from '@/utils/storage/storage'
  import { LAYOUT, STORAGE_KEYS, Z_INDEX } from '@/constants'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps({
    visible: {
      type: Boolean,
      default: false,
    },
    modelValue: {
      type: [Number, String],
      default: LAYOUT.GRID.DEFAULT_SIZE,
    },
    triggerRef: {
      type: Object,
      default: null,
    },
  })

  const emit = defineEmits(['update:modelValue', 'close'])

  const localValue = ref(Number(props.modelValue))

  watchEffect(() => {
    localValue.value = Number(props.modelValue)
  })

  const dropdownStyle = ref({
    position: 'fixed',
    top: '170px',
    right: '100px',
    zIndex: Z_INDEX.DROPDOWN,
  })

  const handleOutsideClick = (e) => {
    const { target } = e
    const isGridSizeBtn = target.closest('.grid-size-btn')
    const isDropdown = target.closest('.grid-size-dropdown')

    if (!isGridSizeBtn && !isDropdown) {
      close()
    }
  }

  const close = () => {
    emit('close')
  }

  const updateValue = () => {
    emit('update:modelValue', localValue.value)

    StorageUtil.set(STORAGE_KEYS.GRID_SIZE, localValue.value)
  }

  const updatePosition = () => {
    const triggerBtn = document.querySelector('.grid-size-btn')
    if (triggerBtn) {
      const rect = triggerBtn.getBoundingClientRect()
      dropdownStyle.value = {
        position: 'fixed',
        top: `${rect.bottom + 10}px`,
        right: `${window.innerWidth - rect.right}px`,
        zIndex: Z_INDEX.DROPDOWN,
      }
    }
  }

  watchEffect(() => {
    if (props.visible) {
      setTimeout(updatePosition, 0)
      document.addEventListener('click', handleOutsideClick)
      window.addEventListener('scroll', updatePosition)
      window.addEventListener('resize', updatePosition)
    } else {
      document.removeEventListener('click', handleOutsideClick)
      window.removeEventListener('scroll', updatePosition)
      window.removeEventListener('resize', updatePosition)
    }
  })

  onMounted(() => {
    localValue.value = Number(props.modelValue)
  })

  onBeforeUnmount(() => {
    document.removeEventListener('click', handleOutsideClick)
    window.removeEventListener('scroll', updatePosition)
    window.removeEventListener('resize', updatePosition)
  })
</script>

<template>
  <Teleport to="body">
    <transition name="dropdown-fade">
      <div v-if="visible" class="grid-size-dropdown grid-size-popover" :style="dropdownStyle" @click.stop>
        <div class="grid-size-header">
          <span>{{ $t('share.gridSize.title') }}</span>
          <button class="close-btn" @click="close">
            <i class="fas fa-times" />
          </button>
        </div>
        <div class="size-options">
          <div class="size-slider">
            <input
              v-model="localValue"
              type="range"
              :min="LAYOUT.GRID.MIN_SIZE"
              :max="LAYOUT.GRID.MAX_SIZE"
              :step="LAYOUT.GRID.STEP"
              @input="updateValue"
            />
            <div class="size-preview">
              <div class="small-icon"><i class="fas fa-th" /></div>
              <div class="large-icon"><i class="fas fa-th-large" /></div>
            </div>
          </div>
          <div class="size-value">{{ localValue }}px</div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
  .grid-size-dropdown {
    background-color: rgba(10, 15, 25, 0.95);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-md);
    width: var(--grid-selector-width);
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(4px);
  }

  .grid-size-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    font-weight: 500;
    font-size: 14px;
  }

  .close-btn {
    background: transparent;
    border: none;
    color: var(--color-content-default);
    cursor: pointer;
    transition: color 0.2s;
    font-size: 12px;
    padding: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .close-btn:hover {
    color: var(--color-brand-500);
  }

  .size-options {
    padding: var(--space-md);
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }

  .size-slider {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }

  .size-preview {
    display: flex;
    justify-content: space-between;
    color: var(--color-content-default);
    font-size: 0.75rem;
  }

  .small-icon,
  .large-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
  }

  .small-icon {
    font-size: 0.7rem;
    color: var(--color-content-default);
  }

  .large-icon {
    font-size: 0.9rem;
    color: rgba(var(--color-brand-500-rgb), 0.8);
  }

  .size-value {
    text-align: center;
    color: var(--color-brand-500);
    font-weight: 500;
    font-size: 0.9rem;
  }

  input[type='range'] {
    width: 100%;
    -webkit-appearance: none;
    height: 6px;
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    outline: none;
  }

  input[type='range']::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: 18px;
    height: 18px;
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    cursor: pointer;
    transition: all 0.15s ease;
  }

  input[type='range']::-webkit-slider-thumb:hover {
    box-shadow: 0 0 0 6px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .dropdown-fade-enter-active,
  .dropdown-fade-leave-active {
    transition:
      opacity 0.2s,
      transform 0.2s;
  }

  .dropdown-fade-enter-from,
  .dropdown-fade-leave-to {
    opacity: 0;
    transform: translateY(-10px);
  }

  @media (max-width: 767px) {
    .grid-size-dropdown {
      width: 260px;
      max-width: 90vw;
      position: fixed;
      left: 50% !important;
      right: auto !important;
      transform: translateX(-50%);
    }
  }
</style>
