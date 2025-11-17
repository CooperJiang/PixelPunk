<template>
  <div class="stats-section">
    <div class="stats-grid">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
  defineOptions({
    name: 'CyberStatsSection',
  })
</script>

<style scoped>
  .stats-section {
    @apply mb-6;
  }

  .stats-grid {
    @apply grid gap-4;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  }

  @media (min-width: 768px) {
    .stats-grid {
      grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
      @apply gap-5;
    }
  }

  @media (min-width: 1024px) {
    .stats-grid {
      grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
      @apply gap-6;
    }
  }

  @media (min-width: 1400px) {
    .stats-grid {
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      max-width: 1200px;
    }
  }
</style>
