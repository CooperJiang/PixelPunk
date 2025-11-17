<script setup lang="ts">
  import { ref } from 'vue'
  import { useTheme } from '@/composables/useTheme'
  import { onClickOutside } from '@vueuse/core'
  import type { ThemeToggleProps, ThemeOption } from './types'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'CyberThemeToggle',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<ThemeToggleProps>(), {
    size: 'small',
    showDropdown: true,
    showText: false,
    placement: 'bottom',
    align: 'right',
  })

  const { selectedTheme, appliedTheme, setTheme, themeOptions, currentThemeInfo } = useTheme()

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

  const selectTheme = (theme: ThemeOption['value']) => {
    setTheme(theme)
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
  <div class="theme-toggle">
    <button
      ref="buttonRef"
      class="theme-toggle-btn"
      :class="[`theme-toggle-btn--${size}`, { 'theme-toggle-btn--active': dropdownOpen }]"
      @click="toggleDropdown"
    >
      <i :class="['fas', `fa-${currentThemeInfo.icon}`, 'theme-icon']" />
      <span v-if="showText" class="theme-text">{{ currentThemeInfo.label }}</span>
    </button>

    <Teleport to="body">
      <Transition name="theme-dropdown">
        <div v-show="dropdownOpen && showDropdown" ref="dropdownRef" class="theme-dropdown" :style="dropdownStyle">
          <div class="dropdown-header">
            <i class="fas fa-palette header-icon" />
            <span class="header-title">{{ $t('theme.selectTheme') }}</span>
          </div>

          <div class="theme-options">
            <button
              v-for="option in themeOptions"
              :key="option.value"
              class="theme-option"
              :class="{ active: selectedTheme === option.value }"
              @click="selectTheme(option.value)"
            >
              <div class="option-icon">
                <i :class="['fas', `fa-${option.icon}`]" />
              </div>
              <div class="option-content">
                <div class="option-label">{{ option.label }}</div>
                <div class="option-description">{{ option.description }}</div>
              </div>
              <div v-if="selectedTheme === option.value" class="option-check">
                <i class="fas fa-check-circle" />
              </div>
            </button>
          </div>

          <div v-if="selectedTheme === 'system'" class="system-hint">
            <i class="fas fa-info-circle hint-icon" />
            <span class="hint-text">{{
              $t('theme.currentSystem', { theme: appliedTheme === 'light' ? $t('theme.light') : $t('theme.dark') })
            }}</span>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
  .theme-toggle {
    @apply relative flex-shrink-0;
  }

  .theme-toggle-btn {
    @apply relative flex items-center justify-center gap-2 border;
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal) var(--ease-in-out);
    background: rgba(var(--color-background-700-rgb), 0.1);
    border-color: var(--color-border-default);
    color: var(--color-content-muted);
    backdrop-filter: var(--backdrop-blur-md) saturate(120%);
    box-shadow: var(--shadow-sm);
  }

  .theme-toggle-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: var(--color-hover-border);
    color: var(--color-brand-500);
    transform: translateY(-2px);
    box-shadow:
      var(--shadow-lg),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .theme-toggle-btn--active {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: var(--color-active-border);
    color: var(--color-brand-500);
  }

  .theme-toggle-btn:focus-visible {
    @apply outline-none;
    box-shadow:
      0 0 0 2px var(--color-focus-ring),
      0 0 25px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .theme-toggle-btn--small {
    width: var(--size-button-sm);
    height: var(--size-button-sm);
    font-size: var(--text-sm);
    padding: 0;
    border-radius: var(--radius-sm);
  }

  .theme-toggle-btn--small .theme-icon {
    font-size: var(--size-icon-sm);
  }

  .theme-toggle-btn--medium {
    height: var(--size-button-md);
    padding: 0 var(--space-md);
    font-size: var(--text-base);
    border-radius: var(--radius-sm);
  }

  .theme-toggle-btn--medium .theme-icon {
    font-size: var(--size-icon-sm);
  }

  .theme-toggle-btn--large {
    height: var(--size-button-lg);
    padding: 0 var(--space-lg);
    font-size: var(--text-lg);
    border-radius: var(--radius-lg);
  }

  .theme-toggle-btn--large .theme-icon {
    font-size: var(--size-icon-md);
  }

  .theme-icon {
    transition: transform var(--transition-normal) ease;
  }

  .theme-text {
    font-weight: var(--font-medium);
  }

  .theme-dropdown {
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

  .theme-options {
    padding: var(--space-sm);
    max-height: 400px;
    overflow-y: auto;
    overflow-x: hidden;
  }

  .theme-options::-webkit-scrollbar {
    width: 6px;
  }

  .theme-options::-webkit-scrollbar-track {
    background: rgba(var(--color-background-800-rgb), 0.5);
    border-radius: var(--radius-sm);
  }

  .theme-options::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
  }

  .theme-options::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .theme-option {
    @apply flex w-full cursor-pointer items-center border text-left;
    gap: var(--space-xs);
    padding: var(--space-sm);
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal) ease;
    background: var(--color-hover-bg-neutral);
    border-color: var(--color-border-subtle);
    margin-bottom: var(--space-sm);
  }

  .theme-option:last-child {
    margin-bottom: 0;
  }

  .theme-option:hover {
    background: var(--color-hover-bg);
    border-color: var(--color-hover-border);
    transform: translateY(-1px);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .theme-option.active {
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
    color: var(--color-brand-500);
    font-size: 11px;
  }

  .theme-option.active .option-icon {
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

  .system-hint {
    @apply flex items-center border-t;
    gap: var(--space-sm);
    padding: var(--space-sm) var(--space-md);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-color: var(--color-border-default);
  }

  .hint-icon {
    font-size: 10px;
    color: var(--color-brand-500);
  }

  .hint-text {
    font-size: 10px;
    color: var(--color-content-muted);
  }

  .theme-dropdown-enter-active,
  .theme-dropdown-leave-active {
    transition: all var(--transition-normal) ease;
  }

  .theme-dropdown-enter-from,
  .theme-dropdown-leave-to {
    opacity: 0;
    transform: scale(0.95);
  }

  @media (max-width: 768px) {
    .theme-dropdown {
      width: 240px;
    }

    .theme-option {
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
