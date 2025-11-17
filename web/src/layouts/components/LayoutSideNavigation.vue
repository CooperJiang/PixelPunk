<script setup lang="ts">
  import { useRouter } from 'vue-router'
  import SidebarTopbar from '@/components/Sidebar/SidebarTopbar.vue'

  const router = useRouter()

  const redirectToAuth = () => {
    router.push('/auth')
  }
</script>

<template>
  <div class="side-navigation">
    <cyberSidebar />
    <SidebarTopbar @login="redirectToAuth" />
  </div>
</template>

<style scoped>
  .side-navigation {
    display: contents;
  }
</style>
