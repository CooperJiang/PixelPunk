<script setup lang="ts">
  import type { AccessLevel, ToggleType, ToggleSize, AccessLevelToggleProps, AccessLevelToggleEmits } from './types'
  import { ACCESS_LEVEL_ICONS, getAccessLevelTitle } from '@/utils/accessLevel'

  const props = withDefaults(defineProps<AccessLevelToggleProps>(), {
    type: 'image' as ToggleType,
    size: 'medium' as ToggleSize,
    disabled: false,
  })

  const emit = defineEmits<AccessLevelToggleEmits>()

  defineOptions({ name: 'AccessLevelToggle' })

  const handleToggle = () => {
    if (!props.disabled) emit('toggle')
  }

  const sizeBox: Record<ToggleSize, string> = {
    small: 'w-5 h-5',
    medium: 'w-6 h-6',
    large: 'w-7 h-7',
  }
  const sizeIcon: Record<ToggleSize, string> = {
    small: 'text-[9px]',
    medium: 'text-[10px]',
    large: 'text-[11px]',
  }
  const levelBase: Record<AccessLevel, string> = {
    private: 'border border-[rgba(var(--color-error-rgb),0.4)] bg-[rgba(var(--color-error-rgb),0.2)]',
    public: 'border border-[rgba(var(--color-success-rgb),0.4)] bg-[rgba(var(--color-success-rgb),0.2)]',
    protected: 'border border-[rgba(var(--color-warning-rgb),0.4)] bg-[rgba(var(--color-warning-rgb),0.2)]',
  }
  const levelHover: Record<AccessLevel, string> = {
    private:
      'hover:border-[rgba(var(--color-error-rgb),0.9)] hover:bg-[rgba(var(--color-error-rgb),0.25)] hover:shadow-[0_0_8px_rgba(var(--color-error-rgb),0.3)]',
    public:
      'hover:border-[rgba(var(--color-success-rgb),0.9)] hover:bg-[rgba(var(--color-success-rgb),0.25)] hover:shadow-[0_0_8px_rgba(var(--color-success-rgb),0.3)]',
    protected:
      'hover:border-[rgba(var(--color-warning-rgb),0.9)] hover:bg-[rgba(var(--color-warning-rgb),0.25)] hover:shadow-[0_0_8px_rgba(var(--color-warning-rgb),0.3)]',
  }
  const iconColor: Record<AccessLevel, string> = {
    private: 'text-[rgb(var(--color-error-rgb))]',
    public: 'text-[rgb(var(--color-success-rgb))]',
    protected: 'text-[rgb(var(--color-warning-rgb))]',
  }
  const iconHover: Record<AccessLevel, string> = {
    private: 'group-hover:[text-shadow:0_0_8px_rgba(var(--color-error-rgb),0.8)]',
    public: 'group-hover:[text-shadow:0_0_8px_rgba(var(--color-success-rgb),0.8)]',
    protected: 'group-hover:[text-shadow:0_0_8px_rgba(var(--color-warning-rgb),0.8)]',
  }
</script>

<template>
  <div
    class="group absolute right-2 top-2 z-[5] flex items-center justify-center rounded-full backdrop-blur-sm transition-all duration-200 ease-in-out"
    :class="[
      sizeBox[size],
      'bg-[rgba(var(--color-background-900-rgb),0.4)]',
      levelBase[accessLevel],
      disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer hover:scale-105',
      !disabled && levelHover[accessLevel],
    ]"
    :title="getAccessLevelTitle(type, accessLevel)"
    @click.stop="handleToggle"
  >
    <i
      :class="[
        ACCESS_LEVEL_ICONS[accessLevel],
        sizeIcon[size],
        iconColor[accessLevel],
        'transition-all duration-200 [text-shadow:0_1px_2px_rgba(0,0,0,0.5)]',
        !disabled && iconHover[accessLevel],
      ]"
    />
  </div>
</template>
