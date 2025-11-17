<script setup lang="ts">
  interface Props {
    title?: string
    subtitle?: string
    icon?: string
    hoverable?: boolean
    padding?: 'none' | 'small' | 'medium' | 'large'
  }

  const props = withDefaults(defineProps<Props>(), {
    hoverable: true,
    padding: 'medium',
  })

  const paddingClass = computed(() => {
    const map = {
      none: 'p-0',
      small: 'p-3',
      medium: 'p-4',
      large: 'p-6',
    }
    return map[props.padding]
  })
</script>

<template>
  <div
    class="card"
    :class="{
      'card-hoverable': hoverable,
    }"
  >
    <div v-if="title || $slots.header" class="card-header">
      <slot name="header">
        <div class="flex items-center gap-3">
          <i v-if="icon" :class="icon" class="card-icon"></i>
          <div class="flex-1">
            <h3 class="card-title">{{ title }}</h3>
            <p v-if="subtitle" class="card-subtitle">{{ subtitle }}</p>
          </div>
        </div>
      </slot>
    </div>

    <div class="card-content" :class="paddingClass">
      <slot />
    </div>

    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<style scoped>
  .card {
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-cyber-sm);
    overflow: hidden;
  }

  .card-hoverable {
    transition: all var(--transition-normal) var(--ease-in-out);
    cursor: pointer;
  }

  .card-hoverable:hover {
    border-color: var(--color-hover-border);
    box-shadow: var(--shadow-cyber-md);
    transform: translateY(-2px);
  }

  .card-header {
    padding: 1rem 1.25rem;
    border-bottom: 1px solid var(--color-border-subtle);
    background: rgba(var(--color-background-800-rgb), 0.5);
  }

  .card-icon {
    width: 2.5rem;
    height: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-brand-500-rgb), 0.05));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-lg);
    color: var(--color-brand-500);
    font-size: 1.25rem;
  }

  .card-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-content-heading);
    line-height: 1.25;
  }

  .card-subtitle {
    font-size: 0.875rem;
    color: var(--color-content-muted);
    margin-top: 0.25rem;
  }

  .card-content {
  }

  .card-footer {
    padding: 0.75rem 1.25rem;
    border-top: 1px solid var(--color-border-subtle);
    background: rgba(var(--color-background-800-rgb), 0.3);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
</style>
