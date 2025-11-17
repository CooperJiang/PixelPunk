<template>
  <header class="layout-header" :class="{ compact: compact }">
    <div class="header-main" :class="{ compact: compact }">
      <div class="header-left">
        <div class="title-section">
          <h1 class="page-title" :class="{ compact: compact }">
            <slot name="icon">
              <i v-if="icon" :class="[icon, 'title-icon', { compact: compact }]" />
            </slot>
            <span class="title-text">
              <slot name="title">{{ title }}</slot>
            </span>
          </h1>
          <div v-if="$slots.subtitle || subtitle" class="subtitle">
            <slot name="subtitle">{{ subtitle }}</slot>
          </div>
        </div>
      </div>

      <div v-if="$slots.actions" class="header-actions" :class="{ compact: compact }">
        <slot name="actions" />
      </div>
    </div>

    <div v-if="hasTopbar" class="header-topbar" :class="{ compact: compact }">
      <slot name="topbar" />
    </div>
  </header>
</template>

<script setup lang="ts">
  interface AdminWrapperHeaderProps {
    title?: string
    subtitle?: string
    icon?: string
    compact?: boolean
    hasTopbar?: boolean
  }

  defineProps<AdminWrapperHeaderProps>()
</script>

<style scoped>
  .layout-header {
    @apply relative z-[100] flex-shrink-0;
    background: var(--color-background-900);
    border-bottom: 1px solid var(--color-border-subtle);
  }

  .header-main {
    @apply flex min-h-14 items-center justify-between gap-6 px-6 py-3.5;
    background: var(--color-background-900);
  }

  .header-left {
    @apply flex min-w-0 flex-1 items-center gap-5;
  }

  .title-section {
    @apply min-w-0 shrink-0;
  }

  .page-title {
    @apply m-0 flex items-center gap-3 text-xl font-bold leading-tight text-[var(--color-cyber-light)];
  }

  .title-icon {
    @apply shrink-0 text-lg text-brand-400;
  }

  .title-text {
    @apply overflow-hidden text-ellipsis whitespace-nowrap text-content-heading;
    letter-spacing: 0.025em;
  }

  .subtitle {
    @apply mt-0.5 text-xs font-normal text-[var(--color-content-muted)] opacity-85;
  }

  .header-actions {
    @apply flex shrink-0 flex-wrap items-center gap-3;
  }

  .header-topbar {
    @apply px-6 py-3;
    border-top: 1px solid var(--color-border-subtle);
    background: var(--color-background-800);
  }

  .header-main.compact {
    @apply min-h-12 px-6 py-2.5;
  }

  .page-title.compact {
    @apply text-base;
  }

  .title-icon.compact {
    @apply text-sm;
  }

  .header-topbar.compact {
    @apply px-6 py-2;
  }

  @media (max-width: 768px) {
    .header-main {
      @apply flex-col items-stretch gap-4 px-4 py-3;
    }

    .header-left {
      @apply flex-col items-start gap-3;
    }

    .header-actions {
      @apply justify-start gap-2;
    }

    .header-topbar {
      @apply px-4 py-3;
    }
  }
</style>
