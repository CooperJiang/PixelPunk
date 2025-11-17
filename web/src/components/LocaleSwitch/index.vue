<script setup lang="ts">
  import { ref } from 'vue'
  import { useLocale } from '@/composables/useLocale'
  import { onClickOutside } from '@vueuse/core'
  import type { LocaleSwitchProps, LocaleOption } from './types'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'LocaleSwitch',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<LocaleSwitchProps>(), {
    size: 'small',
    showDropdown: true,
    showText: false,
    placement: 'bottom',
    align: 'right',
  })

  const { currentLocale, setLocale, localeOptions, currentLocaleInfo } = useLocale()

  const dropdownOpen = ref(false)
  const dropdownRef = ref<HTMLElement>()
  const buttonRef = ref<HTMLElement>()
  const dropdownStyle = ref<{ top?: string; bottom?: string; left?: string; right?: string }>({})

  const updateDropdownPosition = () => {
    if (!buttonRef.value) return

    const rect = buttonRef.value.getBoundingClientRect()
    const style: typeof dropdownStyle.value = {}

    if (props.placement === 'top') {
      style.bottom = `${window.innerHeight - rect.top + 8}px`
    } else {
      style.top = `${rect.bottom + 8}px`
    }

    if (props.align === 'left') {
      style.left = `${rect.left}px`
    } else {
      style.right = `${window.innerWidth - rect.right}px`
    }

    dropdownStyle.value = style
  }

  const toggleDropdown = () => {
    if (props.showDropdown) {
      dropdownOpen.value = !dropdownOpen.value
      if (dropdownOpen.value) {
        updateDropdownPosition()
      }
    }
  }

  const selectLocale = (locale: LocaleOption['value']) => {
    setLocale(locale)
    dropdownOpen.value = false
  }

  onClickOutside(
    dropdownRef,
    (event) => {
      if (buttonRef.value?.contains(event.target as Node)) {
        return
      }
      dropdownOpen.value = false
    },
    { ignore: [buttonRef] }
  )
</script>

<template>
  <div class="locale-switch">
    <button
      ref="buttonRef"
      class="locale-switch-btn"
      :class="[`locale-switch-btn--${size}`, { 'locale-switch-btn--active': dropdownOpen }]"
      @click="toggleDropdown"
    >
      <i class="fas fa-language locale-icon" />
      <span v-if="showText" class="locale-text">{{ currentLocaleInfo.nativeName }}</span>
    </button>

    <Teleport to="body">
      <Transition name="locale-dropdown">
        <div v-show="dropdownOpen && showDropdown" ref="dropdownRef" class="locale-dropdown" :style="dropdownStyle">
          <div class="dropdown-header">
            <i class="fas fa-globe header-icon" />
            <span class="header-title">{{ $t('locale.selectLanguage') }}</span>
          </div>

          <div class="locale-options">
            <button
              v-for="option in localeOptions"
              :key="option.value"
              class="locale-option"
              :class="{ active: currentLocale === option.value }"
              @click="selectLocale(option.value)"
            >
              <div class="option-icon">
                {{ option.emoji }}
              </div>
              <div class="option-content">
                <div class="option-label">{{ option.nativeName }}</div>
                <div class="option-description">{{ option.description }}</div>
              </div>
              <div v-if="currentLocale === option.value" class="option-check">
                <i class="fas fa-check-circle" />
              </div>
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
  .locale-switch {
    @apply relative flex-shrink-0;
  }

  .locale-switch-btn {
    @apply relative flex items-center justify-center gap-2 border;
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal) var(--ease-in-out);
    background: rgba(var(--color-background-700-rgb), 0.1);
    border-color: var(--color-border-default);
    color: var(--color-content-muted);
    backdrop-filter: var(--backdrop-blur-md) saturate(120%);
    box-shadow: var(--shadow-sm);
  }

  .locale-switch-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: var(--color-hover-border);
    color: var(--color-brand-500);
    transform: translateY(-2px);
    box-shadow:
      var(--shadow-lg),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .locale-switch-btn--active {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: var(--color-active-border);
    color: var(--color-brand-500);
  }

  .locale-switch-btn:focus-visible {
    @apply outline-none;
    box-shadow:
      0 0 0 2px var(--color-focus-ring),
      0 0 25px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .locale-switch-btn--small {
    width: var(--size-button-sm);
    height: var(--size-button-sm);
    font-size: var(--text-sm);
    padding: 0;
    border-radius: var(--radius-sm);
  }

  .locale-switch-btn--small .locale-icon {
    font-size: var(--size-icon-sm);
  }

  .locale-switch-btn--medium {
    height: var(--size-button-md);
    padding: 0 var(--space-md);
    font-size: var(--text-base);
    border-radius: var(--radius-sm);
  }

  .locale-switch-btn--medium .locale-icon {
    font-size: var(--size-icon-sm);
  }

  .locale-switch-btn--large {
    height: var(--size-button-lg);
    padding: 0 var(--space-lg);
    font-size: var(--text-lg);
    border-radius: var(--radius-lg);
  }

  .locale-switch-btn--large .locale-icon {
    font-size: var(--size-icon-md);
  }

  .locale-icon {
    transition: transform var(--transition-normal) ease;
  }

  .locale-text {
    font-weight: var(--font-medium);
  }

  .locale-dropdown {
    position: fixed;
    z-index: 9999;
    width: 256px;
    overflow: hidden;
    border-radius: var(--radius-sm);
    background: linear-gradient(145deg, rgba(var(--color-background-700-rgb), 0.98), rgba(var(--color-background-800-rgb), 0.96));
    backdrop-filter: var(--backdrop-blur-xl) saturate(180%);
    border: 1px solid var(--color-border-default);
    box-shadow: var(--shadow-xl);
  }

  .dropdown-header {
    @apply flex items-center border-b;
    gap: var(--space-xs);
    padding: var(--space-sm) var(--space-md);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08), rgba(var(--color-background-700-rgb), 0.8));
    border-color: var(--color-border-default);
  }

  .header-icon {
    font-size: 11px;
    color: var(--color-brand-500);
  }

  .header-title {
    font-size: 11px;
    font-weight: var(--font-medium);
    color: var(--color-content-heading);
  }

  .locale-options {
    padding: var(--space-sm);
    max-height: 400px;
    overflow-y: auto;
    overflow-x: hidden;
  }

  .locale-options::-webkit-scrollbar {
    width: 6px;
  }

  .locale-options::-webkit-scrollbar-track {
    background: rgba(var(--color-background-800-rgb), 0.5);
    border-radius: var(--radius-sm);
  }

  .locale-options::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
  }

  .locale-options::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .locale-option {
    @apply flex w-full cursor-pointer items-center border text-left;
    gap: var(--space-xs);
    padding: var(--space-sm);
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal) ease;
    background: var(--color-hover-bg-neutral);
    border-color: var(--color-border-subtle);
    margin-bottom: var(--space-sm);
  }

  .locale-option:last-child {
    margin-bottom: 0;
  }

  .locale-option:hover {
    background: var(--color-hover-bg);
    border-color: var(--color-hover-border);
    transform: translateY(-1px);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .locale-option.active {
    background: var(--color-active-bg);
    border-color: var(--color-active-border);
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .option-icon {
    @apply flex flex-shrink-0 items-center justify-center;
    width: 28px;
    height: 28px;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.05));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    font-size: 16px;
  }

  .locale-option.active .option-icon {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25), rgba(var(--color-brand-500-rgb), 0.15));
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .option-content {
    @apply flex flex-1 flex-col gap-0;
  }

  .option-label {
    font-size: var(--text-xs);
    font-weight: var(--font-medium);
    line-height: var(--leading-tight);
    color: var(--color-content-heading);
  }

  .option-description {
    font-size: 10px;
    line-height: var(--leading-tight);
    color: var(--color-content-muted);
  }

  .option-check {
    @apply flex-shrink-0;
    font-size: 13px;
    color: var(--color-brand-500);
  }

  .locale-dropdown-enter-active,
  .locale-dropdown-leave-active {
    transition: all var(--transition-normal) ease;
  }

  .locale-dropdown-enter-from,
  .locale-dropdown-leave-to {
    opacity: 0;
    transform: scale(0.95);
  }

  @media (max-width: 768px) {
    .locale-dropdown {
      width: 240px;
    }

    .locale-option {
      padding: 0.375rem;
    }

    .option-icon {
      @apply h-6 w-6;
      font-size: 10px;
    }

    .option-label {
      font-size: 11px;
    }

    .option-description {
      font-size: 9px;
    }
  }
</style>
