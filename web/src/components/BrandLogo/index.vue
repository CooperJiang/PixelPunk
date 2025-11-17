<script setup lang="ts">
  import { computed } from 'vue'
  import { useConfig } from '@/composables/useConfig'
  import { useSettingsStore } from '@/store/settings'

  defineOptions({
    name: 'BrandLogo',
  })

  interface Props {
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
    showVersion?: boolean
    showIcon?: boolean
    clickable?: boolean
    collapsed?: boolean
    containerClass?: string
    shape?: 'rectangle' | 'square' | 'circle'
  }

  const props = withDefaults(defineProps<Props>(), {
    size: 'md',
    showVersion: true,
    showIcon: true,
    clickable: true,
    collapsed: false,
    containerClass: '',
    shape: 'rectangle',
  })

  const { appInfo, brand } = useConfig()
  const settingsStore = useSettingsStore()

  const versionDisplay = computed(() => {
    const version = settingsStore.currentVersion
    if (version && version !== '1.0.0') {
      return `VERSION v${version}`
    }
    return 'VERSION v1.0.0'
  })

  const dynamicSiteName = computed(() => {
    const siteName = settingsStore.siteName
    if (siteName) {
      const words = siteName.split(' ').filter((word) => word.trim())
      if (words.length >= 2) {
        return {
          firstWord: words[0],
          otherWords: words.slice(1),
        }
      }
      return {
        firstWord: siteName,
        otherWords: [],
      }
    }
    return {
      firstWord: brand.pixel,
      otherWords: [brand.punk],
    }
  })

  const dynamicLogo = computed(() => {
    return settingsStore.siteLogoUrl || brand.logo
  })

  const sizeClasses = computed(() => {
    const sizeMap = {
      xs: {
        container: 'brand-xs',
        icon: 'w-5 h-5',
        title: 'text-base',
        version: 'text-[0.5rem]',
      },
      sm: {
        container: 'brand-sm',
        icon: 'w-6 h-6',
        title: 'text-lg',
        version: 'text-[0.55rem]',
      },
      md: {
        container: 'brand-md',
        icon: 'w-7 h-7',
        title: 'text-xl',
        version: 'text-[0.6rem]',
      },
      lg: {
        container: 'brand-lg',
        icon: 'w-9 h-9',
        title: 'text-2xl',
        version: 'text-xs',
      },
      xl: {
        container: 'brand-xl',
        icon: 'w-12 h-12',
        title: 'text-3xl',
        version: 'text-sm',
      },
    }
    return sizeMap[props.size]
  })

  const shapeClasses = computed(() => {
    const shapeMap = {
      rectangle: 'brand-rectangle',
      square: 'brand-square',
      circle: 'brand-circle rounded-full',
    }
    return shapeMap[props.shape]
  })
</script>

<template>
  <div
    class="brand-logo"
    :class="[
      sizeClasses.container,
      shapeClasses,
      containerClass,
      {
        'brand-collapsed': collapsed,
        'brand-clickable': clickable,
      },
    ]"
  >
    <component :is="clickable ? 'router-link' : 'div'" :to="clickable ? '/' : undefined" class="brand-link">
      <img
        v-if="showIcon && !collapsed"
        :src="dynamicLogo"
        :alt="settingsStore.siteName || appInfo.name"
        class="brand-icon"
        :class="sizeClasses.icon"
      />

      <img
        v-if="collapsed"
        :src="dynamicLogo"
        :alt="settingsStore.siteName || appInfo.name"
        class="brand-icon-only"
        :class="sizeClasses.icon"
      />

      <div v-if="!collapsed" class="brand-text">
        <h1 class="brand-title" :class="sizeClasses.title">
          <span class="brand-word brand-color-1">{{ dynamicSiteName.firstWord }}</span>
          <span
            v-for="(word, index) in dynamicSiteName.otherWords"
            :key="index"
            :class="`brand-word brand-color-${((index + 1) % 5) + 1}`"
          >
            {{ word }}
          </span>
        </h1>
        <span v-if="showVersion" class="brand-version" :class="sizeClasses.version">
          {{ versionDisplay }}
        </span>
      </div>
    </component>
  </div>
</template>

<style scoped>
  .brand-logo {
    flex-shrink: 0;
    display: inline-flex;
    align-items: center;
  }

  .brand-rectangle {
    border-radius: var(--radius-sm);
    padding: var(--space-sm) var(--space-md);
  }

  .brand-square {
    border-radius: var(--radius-sm);
    padding: var(--space-sm);
    aspect-ratio: 1;
  }

  .brand-circle {
    padding: var(--space-sm);
    aspect-ratio: 1;
  }

  .brand-xs {
    gap: var(--space-xs);
  }

  .brand-sm {
    gap: var(--space-sm);
  }

  .brand-md {
    gap: var(--space-md);
  }

  .brand-lg {
    gap: var(--space-lg);
  }

  .brand-xl {
    gap: var(--space-xl);
  }

  .brand-link {
    display: flex;
    align-items: center;
    gap: inherit;
    text-decoration: none;
    transition: all var(--transition-normal) var(--ease-in-out);
    width: 100%;
    height: 100%;
  }

  .brand-clickable .brand-link {
    cursor: pointer;
  }

  .brand-clickable:hover {
    background: rgba(var(--color-brand-500-rgb), 0.05);
    transform: scale(1.02);
  }

  .brand-icon {
    transition: all var(--transition-normal) var(--ease-in-out);
    filter: drop-shadow(0 0 10px rgba(var(--color-brand-500-rgb), 0.6));
    flex-shrink: 0;
  }

  .brand-icon-only {
    transition: all var(--transition-normal) var(--ease-in-out);
    filter: drop-shadow(0 0 10px rgba(var(--color-brand-500-rgb), 0.6));
    margin: 0 auto;
  }

  .brand-text {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
    min-width: 0;
  }

  .brand-title {
    font-weight: 900;
    line-height: var(--leading-none);
    letter-spacing: var(--tracking-wide);
    margin: 0;
    font-family: 'Courier New', monospace;
    display: flex;
    align-items: baseline;
    gap: var(--space-xs);
    flex-wrap: wrap;
  }

  .brand-word {
    white-space: nowrap;
    font-weight: 900;
  }

  .brand-color-1 {
    color: var(--color-brand-500);
    text-shadow:
      0 0 10px rgba(var(--color-brand-500-rgb), 0.8),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .brand-color-2 {
    color: var(--color-brand-600);
    text-shadow:
      0 0 10px rgba(var(--color-brand-600-rgb), 0.8),
      0 0 20px rgba(var(--color-brand-600-rgb), 0.4);
  }

  .brand-color-3 {
    color: var(--color-success-500);
    text-shadow:
      0 0 10px rgba(var(--color-success-rgb), 0.8),
      0 0 20px rgba(var(--color-success-rgb), 0.4);
  }

  .brand-color-4 {
    color: var(--color-brand-400);
    text-shadow:
      0 0 10px rgba(var(--color-brand-400-rgb), 0.8),
      0 0 20px rgba(var(--color-brand-400-rgb), 0.4);
  }

  .brand-color-5 {
    color: var(--color-warning-500);
    text-shadow:
      0 0 10px rgba(var(--color-warning-rgb), 0.8),
      0 0 20px rgba(var(--color-warning-rgb), 0.4);
  }

  .brand-version {
    color: var(--color-success-500);
    font-weight: var(--font-semibold);
    text-transform: uppercase;
    letter-spacing: var(--tracking-wide);
    font-family: 'Courier New', monospace;
    white-space: nowrap;
  }

  .brand-collapsed {
    justify-content: center;
    padding: var(--space-sm);
  }

  .brand-collapsed .brand-link {
    justify-content: center;
  }

  .brand-collapsed .brand-text {
    display: none;
  }

  @media (max-width: 768px) {
    .brand-title {
      font-size: var(--text-base);
    }

    .brand-rectangle,
    .brand-square,
    .brand-circle {
      padding: var(--space-xs) var(--space-sm);
    }
  }
</style>
