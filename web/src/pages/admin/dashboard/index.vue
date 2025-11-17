<script setup lang="ts">
  import { ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'AdminDashboard',
  })

  const { $t } = useTexts()

  /* 导入所有组件 */
  import DashboardMetricsRow from './components/DashboardMetricsRow.vue'
  import AIServicesCard from './components/AIServicesCard.vue'
  import UploadTrendsChart from './components/UploadTrendsChart.vue'
  import RecentUploadsCard from './components/RecentUploadsCard.vue'
  import TagStatsCard from './components/TagStatsCard.vue'

  /* 刷新控制 */
  const refreshKey = ref(0)
</script>

<template>
  <div class="admin-dashboard-page">
    <CyberAdminWrapper
      :title="$t('admin.dashboard.title')"
      :subtitle="$t('admin.dashboard.subtitle')"
      icon="fas fa-tachometer-alt"
    >
      <template #content>
        <div class="dashboard-content">
          <DashboardMetricsRow :refresh-key="refreshKey" />

          <div class="charts-analytics">
            <div class="main-chart-section">
              <UploadTrendsChart :key="refreshKey" />
            </div>

            <div class="monitoring-panels">
              <div class="panel-row">
                <div class="panel-item">
                  <AIServicesCard :key="refreshKey" />
                </div>
                <div class="panel-item">
                  <TagStatsCard :key="refreshKey" />
                </div>
              </div>
            </div>
          </div>

          <div class="recent-uploads-section">
            <div class="recent-uploads-full">
              <RecentUploadsCard :key="refreshKey" />
            </div>
          </div>
        </div>
      </template>
    </CyberAdminWrapper>
  </div>
</template>

<style scoped>
  .admin-dashboard-page {
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    color: var(--color-content-default);
  }

  .dashboard-content {
    @apply w-full;
    max-width: 100vw;
    position: relative;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-xl);
  }

  .charts-analytics {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
  }

  .main-chart-section {
    @apply w-full;
    min-height: 420px;
  }

  .monitoring-panels {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
  }

  .panel-row {
    @apply grid grid-cols-1 lg:grid-cols-2;
    gap: var(--space-xl);
  }

  .panel-item {
    @apply w-full;
    min-height: 280px;
  }

  .recent-uploads-section {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    margin-bottom: 20px;
  }

  .recent-uploads-full {
    @apply w-full;
    min-height: 320px;
  }

  @media (max-width: 1024px) {
    .dashboard-content {
      gap: var(--space-lg);
    }

    .panel-row {
      @apply grid-cols-1;
    }

    .recent-uploads-section {
      gap: var(--space-sm);
    }
  }

  @media (max-width: 768px) {
    .dashboard-content {
      gap: var(--space-md);
    }
  }

  @media (max-width: 640px) {
    .dashboard-content {
      gap: var(--space-sm);
    }
  }
</style>
