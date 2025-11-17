<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { UserAvatarProps } from './types'

  defineOptions({
    name: 'CyberUserAvatar',
  })

  const props = withDefaults(defineProps<UserAvatarProps>(), {
    avatarUrl: '',
    username: '',
    size: 'md',
    showBorder: true,
    showGlow: true,
    customStyles: () => ({}),
  })

  const { $t } = useTexts()

  const hasImageError = ref(false)
  const handleImageError = () => {
    hasImageError.value = true
  }

  const initials = computed(() => {
    if (hasImageError.value || !props.avatarUrl) {
      const username = props.username || ''
      return username ? username.charAt(0).toUpperCase() : '?'
    }
    return ''
  })

  const sizeClass = computed(() => {
    switch (props.size) {
      case 'sm':
        return 'size-sm'
      case 'lg':
        return 'size-lg'
      default:
        return 'size-md'
    }
  })
</script>

<template>
  <div class="cyber-avatar" :class="[sizeClass, { 'with-border': showBorder }, { 'with-glow': showGlow }]" :style="customStyles">
    <img v-if="avatarUrl" :src="avatarUrl" class="avatar-image" :alt="$t('theme.user.avatarAlt')" @error="handleImageError" />
    <div v-else class="avatar-placeholder">
      <span>{{ initials }}</span>
    </div>
  </div>
</template>

<style scoped>
  .cyber-avatar {
    @apply relative flex items-center justify-center overflow-hidden rounded;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08), rgba(var(--color-background-800-rgb), 0.6));
    transition: all 0.3s ease;
  }

  .size-sm {
    @apply h-8 w-8 text-base;
  }

  .size-md {
    @apply h-10 w-10;
    font-size: 1.2rem;
  }

  .size-lg {
    @apply h-16 w-16;
    font-size: 1.8rem;
  }

  .with-border {
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .with-glow:hover {
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateY(-1px);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .avatar-image {
    @apply h-full w-full object-cover;
  }

  .avatar-placeholder {
    @apply flex h-full w-full items-center justify-center font-semibold uppercase;
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }
</style>
