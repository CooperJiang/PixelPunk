<script setup lang="ts">
  import { defineEmits, defineProps, onBeforeUnmount, ref, watchEffect } from 'vue'
  import { StorageUtil } from '@/utils/storage/storage'
  import { STORAGE_KEYS, Z_INDEX } from '@/constants'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps({
    visible: {
      type: Boolean,
      default: false,
    },
    modelValue: {
      type: String,
      default: 'grid',
    },
    layoutOptions: {
      type: Array,
      default: () => [],
    },
    triggerRef: {
      type: Object,
      default: null,
    },
  })

  const emit = defineEmits(['update:modelValue', 'close'])

  const dropdownStyle = ref({
    position: 'fixed',
    top: '170px',
    right: '100px',
    zIndex: Z_INDEX.DROPDOWN,
  })

  const handleOutsideClick = (e) => {
    const { target } = e
    const isLayoutBtn = target.closest('.layout-btn')
    const isDropdown = target.closest('.layout-dropdown')

    if (!isLayoutBtn && !isDropdown) {
      close()
    }
  }

  const close = () => {
    emit('close')
  }

  const selectLayout = (value) => {
    emit('update:modelValue', value)

    StorageUtil.set(STORAGE_KEYS.LAYOUT_MODE, value)

    close()
  }

  const updatePosition = () => {
    const layoutBtn = document.querySelector('.layout-btn')
    if (layoutBtn) {
      const rect = layoutBtn.getBoundingClientRect()
      dropdownStyle.value = {
        position: 'fixed',
        top: `${rect.bottom + 20}px`,
        left: `${rect.left - 222}px`,
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

  onBeforeUnmount(() => {
    document.removeEventListener('click', handleOutsideClick)
    window.removeEventListener('scroll', updatePosition)
    window.removeEventListener('resize', updatePosition)
  })
</script>

<template>
  <Teleport to="body">
    <transition name="dropdown-fade">
      <div v-if="visible" class="layout-dropdown" :style="dropdownStyle" @click.stop>
        <div class="layout-header">
          <span>{{ $t('share.layout.title') }}</span>
          <button class="close-btn" @click="close">
            <i class="fas fa-times" />
          </button>
        </div>
        <div class="layout-options">
          <button
            v-for="option in layoutOptions"
            :key="option.value"
            class="layout-option-btn"
            :class="{ active: modelValue === option.value }"
            @click="selectLayout(option.value)"
          >
            <div class="option-icon">
              <i :class="option.icon" />
            </div>
            <div class="option-content">
              <div class="option-label">{{ option.label }}</div>
              <div class="option-desc">{{ option.description }}</div>
            </div>
          </button>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
  .layout-dropdown {
    background-color: rgba(var(--color-background-900-rgb), 0.95);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-md);
    width: 280px;
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(4px);
  }

  .layout-header {
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

  .layout-options {
    padding: 6px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .layout-option-btn {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 10px;
    border-radius: var(--radius-md);
    background-color: rgba(var(--color-background-800-rgb), 0.4);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-content-default);
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
  }

  .layout-option-btn:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-1px);
  }

  .layout-option-btn.active {
    background-color: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .option-icon {
    width: 28px;
    height: 28px;
    background-color: rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    color: var(--color-brand-500);
    flex-shrink: 0;
  }

  .option-content {
    flex: 1;
    min-width: 0;
  }

  .option-label {
    font-weight: 500;
    margin-bottom: 2px;
    font-size: 13px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .option-desc {
    font-size: 11px;
    color: var(--color-content-default);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
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
    .layout-dropdown {
      width: 260px;
      max-width: 90vw;
      position: fixed;
      left: 50% !important;
      transform: translateX(-50%);
    }
  }
</style>
