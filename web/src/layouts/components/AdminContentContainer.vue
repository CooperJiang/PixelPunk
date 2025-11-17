<script setup lang="ts">
  import { useRoute } from 'vue-router'

  const route = useRoute()
</script>

<template>
  <div class="admin-content-container">
    <div :key="route.path" class="page-transition admin-page-content">
      <router-view />
    </div>
  </div>
</template>

<style scoped>
  .admin-content-container {
    width: 100%;
    height: 100%;
    overflow: hidden;
    position: relative;
    padding: 1rem;
    box-sizing: border-box;
  }

  .admin-page-content {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    border-radius: var(--radius-sm);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.95) 0%,
      rgba(var(--color-background-900-rgb), 0.9) 50%,
      rgba(var(--color-background-800-rgb), 0.95) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    box-sizing: border-box;
    padding: 1.5rem;
  }

  .page-transition {
    animation: slideUp var(--transition-normal) var(--ease-out);
  }
</style>
