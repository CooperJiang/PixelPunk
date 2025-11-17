<script setup lang="ts">
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'

  const route = useRoute()

  const isFullWidthPage = computed(() => ['/', '/home', '/random', '/explore'].includes(route.path))

  const isDocsPage = computed(() => route.path.startsWith('/docs'))

  const pageContainerClasses = computed(() => ({
    'page-container': true,
    'full-width': isFullWidthPage.value,
    'docs-container': isDocsPage.value,
  }))

  const routeContainerClasses = computed(() => ({
    'page-transition': true,
    'main-route-container': true,
    'container mx-auto': !isFullWidthPage.value && !isDocsPage.value,
  }))
</script>

<template>
  <div :class="pageContainerClasses">
    <div :key="route.path" :class="routeContainerClasses">
      <router-view />
    </div>
  </div>
</template>

<style scoped>
  .page-container {
    width: 100%;
    margin: 0 auto;
    padding: 0 0.75rem;
    max-width: 1536px;
  }

  .page-container.full-width {
    max-width: none;
    padding: 0;
  }

  .page-container.docs-container {
    max-width: none;
    padding: 0;
  }

  .page-transition {
    width: 100%;
    animation: slideUp var(--transition-normal) var(--ease-out);
  }

  @media (min-width: 768px) {
    .page-container {
      padding: 0 1.5rem;
    }

    .page-container.docs-container {
      padding: 0;
    }
  }
</style>
