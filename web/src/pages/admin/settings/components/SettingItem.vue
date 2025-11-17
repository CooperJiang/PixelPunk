<script setup lang="ts">
  interface Props {
    label: string
    icon?: string
    description?: string
  }

  defineProps<Props>()
</script>

<template>
  <div class="flex flex-col md:flex-row md:items-start">
    <label class="flex w-full items-center text-sm text-content md:w-48 md:py-2">
      <i v-if="icon" :class="`fas fa-${icon} mr-2 text-brand-600`" />
      {{ label }}
    </label>
    <div class="flex-1">
      <slot />
      <p v-if="description" class="text-content-subtle mt-2 text-xs">{{ description }}</p>
    </div>
  </div>
</template>
