<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import { initBaiduAnalytics } from './utils/system/baidu-analytics'
  import { initGoogleAnalytics } from './utils/system/google-analytics'
  import { useSettingsStore } from './store/settings'

  const showUploadDrawer = ref(false)
  const route = useRoute()
  const settingsStore = useSettingsStore()

  const shouldShowUploadFloat = computed(() => {
    const excludeRoutes = ['/auth', '/setup', '/docs', '/refuse']
    return !excludeRoutes.some((r) => route.path.startsWith(r))
  })

  // 监听百度统计配置变化，自动初始化
  watch(
    () => [settingsStore.baiduEnabled, settingsStore.baiduSiteId],
    ([enabled, siteId]) => {
      if (enabled && siteId) {
        initBaiduAnalytics(siteId as string)
      }
    },
    { immediate: true }
  )

  // 监听Google Analytics配置变化，自动初始化
  watch(
    () => [settingsStore.googleEnabled, settingsStore.googleMeasurementId],
    ([enabled, measurementId]) => {
      if (enabled && measurementId) {
        initGoogleAnalytics(measurementId as string)
      }
    },
    { immediate: true }
  )
</script>

<template>
  <router-view />
  <UploadFloatButton v-if="shouldShowUploadFloat" />
  <CyberGlobalUploadFloat v-if="shouldShowUploadFloat" @open-drawer="showUploadDrawer = true" />
  <CyberDrawer
    v-if="shouldShowUploadFloat"
    v-model="showUploadDrawer"
    :title="$t('app.uploadQueue')"
    width="480px"
    :mask-closable="true"
  >
    <template #title-icon>
      <i class="fas fa-cloud-upload-alt text-content" />
    </template>
    <CyberGlobalUploadDrawer @close="showUploadDrawer = false" />
  </CyberDrawer>
</template>

<style>
  @import url('/fonts/fonts.css');
  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  html,
  body {
    overflow-x: hidden;
    position: relative;
    width: 100%;
  }

  ::-webkit-scrollbar {
    width: 6px;
  }

  ::-webkit-scrollbar-track {
    background: transparent;
  }

  ::-webkit-scrollbar-thumb {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.3) 0%, rgba(var(--color-brand-400-rgb), 0.25) 100%);
    border-radius: var(--radius-sm);
    transition: all 0.3s ease;
  }

  ::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(135deg, rgba(var(--color-brand-600-rgb), 0.5) 0%, rgba(var(--color-brand-500-rgb), 0.45) 100%);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  * {
    scrollbar-width: thin;
    scrollbar-color: rgba(var(--color-brand-500-rgb), 0.3) transparent;
  }

  .fade-slide-enter-active,
  .fade-slide-leave-active,
  .fade-scale-enter-active,
  .fade-scale-leave-active,
  .slide-right-enter-active,
  .slide-right-leave-active {
    position: relative;
    width: 100%;
    overflow-x: hidden;
  }
</style>
