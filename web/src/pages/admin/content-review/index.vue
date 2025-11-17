<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import ReviewQueue from './components/ReviewQueue.vue'
  import ReviewLogs from './components/ReviewLogs.vue'

  defineOptions({
    name: 'ContentReviewPage',
  })

  const { $t } = useTexts()

  interface TabItem {
    key: string
    name: string
    icon?: string
    badge?: string | number
  }

  const tabs = computed<TabItem[]>(() => [
    {
      key: 'queue',
      name: $t('admin.contentReview.tabs.queue'),
      icon: 'fas fa-clock',
      badge: '0',
    },
    {
      key: 'logs',
      name: $t('admin.contentReview.tabs.logs'),
      icon: 'fas fa-history',
    },
  ])

  const activeTab = ref('queue')
  const reviewQueueRef = ref<InstanceType<typeof ReviewQueue>>()

  const handleTabChange = (tabId: string) => {
    activeTab.value = tabId
  }

  const refreshStats = (stats: { pending_count?: number }) => {
    const queueTab = tabs.value.find((tab) => tab.key === 'queue')
    if (queueTab && stats?.pending_count !== undefined) {
      queueTab.badge = stats.pending_count.toString()
    }
  }
</script>

<template>
  <div class="admin-content-review-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.contentReview.title')"
      :subtitle="$t('admin.contentReview.subtitle')"
      icon="fas fa-gavel"
      :sidebar-layout="true"
    >
      <template #sidebar>
        <CyberSidebarNav :tabs="tabs" :active-tab="activeTab" @tab-change="handleTabChange" />
      </template>

      <template #content>
        <div v-show="activeTab === 'queue'" class="tab-panel">
          <ReviewQueue ref="reviewQueueRef" :is-active="activeTab === 'queue'" @refresh-stats="refreshStats" />
        </div>

        <div v-show="activeTab === 'logs'" class="tab-panel">
          <ReviewLogs :is-active="activeTab === 'logs'" />
        </div>
      </template>
    </CyberAdminWrapper>
  </div>
</template>

<style scoped lang="scss">
  .admin-content-review-page {
    color: var(--color-content);
  }

  .tab-panel {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    overflow-x: hidden;
  }

  .tab-panel::-webkit-scrollbar {
    width: var(--space-sm);
    background: transparent;
  }

  .tab-panel::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
  }

  .tab-panel::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal);
  }

  .tab-panel::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }
</style>
