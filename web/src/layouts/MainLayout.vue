<script setup lang="ts">
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'
  import { useLayoutManager } from '@/composables/useLayoutManager'
  import LayoutTopNavigation from './components/LayoutTopNavigation.vue'
  import LayoutSideNavigation from './components/LayoutSideNavigation.vue'
  import AdminContentContainer from './components/AdminContentContainer.vue'
  import TopContentContainer from './components/TopContentContainer.vue'

  const route = useRoute()
  const { isTopLayout, isLeftLayout, layoutClasses, isFullWidthPage, isDocsPage, mainContentStyle } = useLayoutManager()

  const mainContentClasses = computed(() => ({
    'full-width-content': isFullWidthPage.value,
    'docs-content': isDocsPage.value,
    'content-with-sidebar': isLeftLayout.value,
    'admin-style-content': isLeftLayout.value,
  }))

  const contentComponent = computed(() => {
    if (isLeftLayout.value) {
      return AdminContentContainer
    }
    return TopContentContainer
  })

  const backgroundProps = computed(() => {
    const path = route.path
    const name = route.name

    if (name === 'home' || path === '/' || path === '/home') {
      return {
        showEnhancedEffects: true,
        showDigitalRain: true,
        showRedOrb: true,
        showAurora: true,
      }
    }

    return {
      showEnhancedEffects: false,
      showDigitalRain: false,
      showRedOrb: false,
      showAurora: false,
    }
  })
</script>

<template>
  <div class="main-layout" :class="layoutClasses">
    <LayoutTopNavigation v-if="isTopLayout" />
    <LayoutSideNavigation v-if="isLeftLayout" />

    <main class="main-content" :class="mainContentClasses" :style="mainContentStyle">
      <component :is="contentComponent" />
    </main>

    <CyberHomeBackground v-bind="backgroundProps" />
  </div>
</template>

<style src="./styles/main-layout.css" scoped></style>
