<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { useRoute } from 'vue-router'
  import { useTheme } from '@/composables/useTheme'
  import NavbarBrand from './NavbarBrand.vue'
  import DesktopNav from './DesktopNav.vue'
  import UserTerminal from './UserTerminal.vue'
  import MobileMenu from './MobileMenu.vue'
  import type { NavbarEmits } from './types'

  defineOptions({
    name: 'Navbar',
  })

  const emit = defineEmits<NavbarEmits>()
  const route = useRoute()
  const { isLight } = useTheme()

  const mobileMenuOpen = ref<boolean>(false)
  const scrollY = ref(0)

  /**
   * 导航栏滚动状态判断逻辑
   * 1. 如果是主页 (/) -> 始终返回 false (不显示背景)
   * 2. 如果是亮色主题 -> 始终返回 true (始终显示背景)
   * 3. 其他情况 -> 根据滚动距离判断 (滚动超过 20px 显示背景)
   */
  const isScrolled = computed(() => {
    if (route.path === '/') {
      return false
    }

    if (isLight.value) {
      return true
    }

    return scrollY.value > 20
  })

  const handleScroll = (event: Event) => {
    const target = event.target as HTMLElement
    scrollY.value = target.scrollTop
  }

  onMounted(() => {
    const scrollContainer = document.querySelector('.main-content') as HTMLElement
    if (scrollContainer) {
      scrollContainer.addEventListener('scroll', handleScroll)
    }
  })

  onUnmounted(() => {
    const scrollContainer = document.querySelector('.main-content') as HTMLElement
    if (scrollContainer) {
      scrollContainer.removeEventListener('scroll', handleScroll)
    }
  })

  const toggleMobileMenu = (): void => {
    mobileMenuOpen.value = !mobileMenuOpen.value
  }
</script>

<template>
  <div class="navbar-wrapper">
    <nav class="cyber-navbar" :class="{ scrolled: isScrolled }">
      <div class="navbar-container">
        <NavbarBrand />
        <DesktopNav />
        <UserTerminal :mobile-menu-open="mobileMenuOpen" @login="emit('login')" @toggle-mobile-menu="toggleMobileMenu" />
      </div>

      <div v-if="isScrolled" class="scan-line" />
    </nav>

    <MobileMenu :is-open="mobileMenuOpen" @close="toggleMobileMenu" @login="emit('login')" />
  </div>
</template>

<style scoped>
  .navbar-wrapper {
    @apply fixed left-0 right-0 top-0 h-16;
    z-index: var(--z-index-dropdown);
  }

  .cyber-navbar {
    @apply relative h-full w-full border-b border-transparent bg-transparent;
    backdrop-filter: none;
    transition: all 0.3s ease;
  }

  .cyber-navbar.scrolled {
    background: rgba(var(--color-background-900-rgb), 0.9);
    backdrop-filter: blur(25px);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 4px 20px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .navbar-container {
    @apply relative flex h-full w-full items-center justify-between;
    padding: 0 max(2rem, calc((100vw - 1600px) / 2));
    overflow: visible;
  }

  .scan-line {
    @apply absolute bottom-0 left-0 h-px w-24;
    background: linear-gradient(90deg, transparent, var(--color-brand-500), transparent);
    animation: scan-across 3s infinite;
  }

  @keyframes scan-across {
    0% {
      left: 0;
      opacity: 0;
    }
    10% {
      opacity: 1;
    }
    90% {
      opacity: 1;
    }
    100% {
      left: calc(100% - 6rem);
      opacity: 0;
    }
  }

  @media (max-width: 768px) {
    .navbar-container {
      @apply px-4;
    }

    .navbar-wrapper {
      height: 60px;
    }
  }
</style>
