<script setup lang="ts">
  import { defineExpose, ref } from 'vue'
  import ApiKeyManager from './ApiKeyManager.vue'

  const apiKeyManagerRef = ref(null)

  defineExpose({
    apiKeyManagerRef,
  })
</script>

<template>
  <div class="apikey-section-card border">
    <ApiKeyManager ref="apiKeyManagerRef" :lazy-load="true" />
  </div>
</template>

<style scoped>
  .apikey-section-card {
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.9));
    border-color: var(--color-border-subtle);
    box-shadow:
      0 4px 16px rgba(0, 0, 0, 0.08),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .apikey-section-card:hover {
    border-color: var(--color-border-default);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.03);
    transform: translateY(-1px);
  }
</style>
