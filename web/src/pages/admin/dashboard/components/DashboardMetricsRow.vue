<script setup lang="ts">
  import { defineProps } from 'vue'
  import UserStatsCard from './UserStatsCard.vue'
  import FileStatsCard from './FileStatsCard.vue'
  import StorageStatsCard from './StorageStatsCard.vue'
  import ShareStatsCard from './ShareStatsCard.vue'

  defineOptions({ name: 'DashboardMetricsRow' })

  /* 可传入一个 key 触发内部子组件刷新 */
  defineProps<{ refreshKey?: number }>()
</script>

<template>
  <div class="metrics-overview">
    <div class="metrics-row">
      <UserStatsCard :key="refreshKey" />
      <FileStatsCard :key="refreshKey" />
      <StorageStatsCard :key="refreshKey" />
      <ShareStatsCard :key="refreshKey" />
    </div>
  </div>
</template>

<style scoped>
  .metrics-overview {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
  }

  .metrics-row {
    @apply grid grid-cols-2 lg:grid-cols-4;
    gap: var(--space-lg);
    align-items: stretch;
  }

  .metrics-row > * {
    height: auto;
    min-height: 165px;
  }

  @media (max-width: 640px) {
    .metrics-row {
      gap: var(--space-sm);
    }
  }
</style>
