<script setup lang="ts">
  import { computed } from 'vue'
  import { useSettingsStore } from '@/store/settings'
  import type { LogoProps } from './types'

  defineOptions({
    name: 'CyberLogo',
  })

  const _props = withDefaults(defineProps<LogoProps>(), {
    iconOnly: false,
  })

  const settingsStore = useSettingsStore()

  /* 将网站名称分成两部分用于样式显示 */
  const logoTextParts = computed(() => {
    const siteName = settingsStore.siteName
    if (siteName.includes(' ')) {
      const parts = siteName.split(' ', 2)
      return { first: parts[0], second: parts[1] || '' }
    }
    const match = siteName.match(/^([A-Z][a-z]+)([A-Z][a-z]+)$/)
    if (match) {
      return { first: match[1], second: match[2] }
    }
    return { first: siteName, second: '' }
  })
</script>

<template>
  <router-link to="/" class="logo-link">
    <div class="logo-container">
      <div class="logo-icon">
        <img src="/logo.png" :alt="`${settingsStore.siteName} Logo`" class="logo-image" />
      </div>
      <div v-if="!iconOnly" class="logo-text">
        <span class="cyber">{{ logoTextParts.first }}</span>
        <span v-if="logoTextParts.second" class="punk"> {{ logoTextParts.second }}</span>
      </div>
    </div>
  </router-link>
</template>

<style scoped>
  .logo-link {
    @apply flex items-center no-underline;
  }

  .logo-container {
    @apply flex items-center gap-2;
  }

  .logo-icon {
    @apply flex h-10 w-10 items-center justify-center rounded-lg transition-all duration-300 ease-in-out;
  }

  .logo-image {
    @apply h-full w-full object-contain;
  }

  .logo-icon:hover {
    @apply scale-105;
  }

  .logo-text {
    @apply text-xl font-bold tracking-wider;
    text-shadow: 0 0 10px rgba(var(--color-brand-500), 0.3);
  }

  .cyber {
    background: linear-gradient(90deg, var(--color-brand-500), var(--color-brand-500) 70%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .punk {
    background: linear-gradient(90deg, var(--color-brand-500) 30%, var(--color-brand-600));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
</style>
