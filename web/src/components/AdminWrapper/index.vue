<script setup lang="ts">
  import type { CyberAdminWrapperProps } from './types'
  import AdminWrapperHeader from './components/AdminWrapperHeader.vue'
  import AdminWrapperContent from './components/AdminWrapperContent.vue'

  defineOptions({ name: 'CyberAdminWrapper' })

  const {
    title,
    subtitle,
    icon,
    compact = false,
    contentPadding = '1.5rem',
    sidebarLayout = false,
    showTopbar = false,
  } = defineProps<CyberAdminWrapperProps>()
</script>

<template>
  <div class="cyber-admin-wrapper" :data-compact="compact">
    <AdminWrapperHeader :title="title" :subtitle="subtitle" :icon="icon" :compact="compact" :has-topbar="showTopbar">
      <template #icon>
        <slot name="icon" />
      </template>
      <template #title>
        <slot name="title" />
      </template>
      <template #subtitle>
        <slot name="subtitle" />
      </template>
      <template #stats>
        <slot name="stats" />
      </template>
      <template #actions>
        <slot name="actions" />
      </template>
      <template v-if="showTopbar" #topbar>
        <slot name="topbar" />
      </template>
    </AdminWrapperHeader>

    <AdminWrapperContent :sidebar-layout="sidebarLayout" :content-padding="contentPadding" :compact="compact">
      <template #sidebar>
        <slot name="sidebar" />
      </template>
      <template #toolbar>
        <slot name="toolbar" />
      </template>
      <template #default>
        <slot name="default" />
      </template>
      <template #content>
        <slot name="content" />
      </template>
    </AdminWrapperContent>
  </div>
</template>

<style scoped>
  .cyber-admin-wrapper {
    @apply relative flex flex-col overflow-hidden;
    flex: 1;
    min-height: 0;
    animation: layoutFadeIn 0.4s ease-out;
  }
</style>

<style>
  .cyber-admin-wrapper[data-compact='true'] .header-actions .cyber-button {
    @apply min-h-8 px-3.5 py-1.5 text-xs;
  }
</style>
