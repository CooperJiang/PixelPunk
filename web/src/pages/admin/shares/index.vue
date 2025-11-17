<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import ShareList from './components/ShareList.vue'
  import VisitorList from './components/VisitorList.vue'
  import Statistics from './components/Statistics.vue'

  defineOptions({
    name: 'AdminShares',
  })

  const { $t } = useTexts()

  interface TabItem {
    key: string
    name: string
    icon?: string
  }

  /* 选项卡数据 */
  const tabs = computed<TabItem[]>(() => [
    { key: 'list', name: $t('admin.shares.tabs.list'), icon: 'fas fa-list-ul' },
    { key: 'visitors', name: $t('admin.shares.tabs.visitors'), icon: 'fas fa-users' },
    { key: 'stats', name: $t('admin.shares.tabs.stats'), icon: 'fas fa-chart-line' },
  ])

  const activeTab = ref('list')

  const setActiveTab = (tab: string) => {
    activeTab.value = tab
  }
</script>

<template>
  <div class="admin-shares-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.shares.title')"
      :subtitle="$t('admin.shares.subtitle')"
      icon="fas fa-share-alt"
      :sidebar-layout="true"
    >
      <template #sidebar>
        <CyberSidebarNav :tabs="tabs" :active-tab="activeTab" @tab-change="setActiveTab" />
      </template>

      <template #content>
        <ShareList v-if="activeTab === 'list'" />
        <VisitorList v-if="activeTab === 'visitors'" />
        <Statistics v-if="activeTab === 'stats'" />
      </template>
    </CyberAdminWrapper>
  </div>
</template>

<style scoped lang="scss">
  .admin-shares-page {
    color: var(--color-content);
  }
</style>
