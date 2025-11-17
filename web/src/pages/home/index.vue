<script setup lang="ts">
  import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import HeroTitle from './components/HeroTitle.vue'
  import StatsDisplay from './components/StatsDisplay.vue'
  import ActionButtons from './components/ActionButtons.vue'

  defineOptions({
    name: 'ToolHome',
  })

  const showNavAndFooter = ref(false)
  const hasScrolled = ref(false)

  const handleScroll = (event: Event) => {
    const target = event.target as HTMLElement
    const scrolled = target.scrollTop > 10
    if (scrolled !== hasScrolled.value) {
      hasScrolled.value = scrolled
      showNavAndFooter.value = scrolled
    }
  }

  const toggleNavAndFooter = (show: boolean) => {
    nextTick(() => {
      const navbar = document.querySelector('.cyber-navbar') as HTMLElement
      const footer = document.querySelector('.footer-container') as HTMLElement

      if (navbar) {
        if (show) {
          navbar.style.setProperty('opacity', '1', 'important')
          navbar.style.setProperty('transform', 'translateY(0)', 'important')
          navbar.style.setProperty('pointer-events', 'auto', 'important')
          navbar.style.setProperty('transition', 'all 0.5s ease-in-out', 'important')
        } else {
          navbar.style.setProperty('opacity', '0', 'important')
          navbar.style.setProperty('transform', 'translateY(-100%)', 'important')
          navbar.style.setProperty('pointer-events', 'none', 'important')
          navbar.style.setProperty('transition', 'all 0.5s ease-in-out', 'important')
        }
      }

      if (footer) {
        if (show) {
          footer.style.setProperty('opacity', '1', 'important')
          footer.style.setProperty('transform', 'translateY(0)', 'important')
          footer.style.setProperty('pointer-events', 'auto', 'important')
          footer.style.setProperty('transition', 'all 0.5s ease-in-out', 'important')
        } else {
          footer.style.setProperty('opacity', '0', 'important')
          footer.style.setProperty('transform', 'translateY(100%)', 'important')
          footer.style.setProperty('pointer-events', 'none', 'important')
          footer.style.setProperty('transition', 'all 0.5s ease-in-out', 'important')
        }
      }
    })
  }

  watch(showNavAndFooter, (newVal) => {
    toggleNavAndFooter(newVal)
  })

  onMounted(() => {
    const immediateHide = () => {
      const navbar = document.querySelector('.cyber-navbar') as HTMLElement
      const footer = document.querySelector('.footer-container') as HTMLElement

      if (navbar) {
        navbar.style.setProperty('opacity', '0', 'important')
        navbar.style.setProperty('transform', 'translateY(-100%)', 'important')
        navbar.style.setProperty('pointer-events', 'none', 'important')
        navbar.style.setProperty('transition', 'none', 'important')
      }

      if (footer) {
        footer.style.setProperty('opacity', '0', 'important')
        footer.style.setProperty('transform', 'translateY(100%)', 'important')
        footer.style.setProperty('pointer-events', 'none', 'important')
        footer.style.setProperty('transition', 'none', 'important')
      }
    }

    immediateHide()

    nextTick(() => {
      immediateHide()
    })

    setTimeout(() => {
      const navbar = document.querySelector('.cyber-navbar') as HTMLElement
      const footer = document.querySelector('.footer-container') as HTMLElement
      const mainContent = document.querySelector('.main-content')

      if (navbar) {
        navbar.style.setProperty('transition', 'all 0.5s ease-in-out', 'important')
      }
      if (footer) {
        footer.style.setProperty('transition', 'all 0.5s ease-in-out', 'important')
      }

      if (mainContent) {
        mainContent.addEventListener('scroll', handleScroll, { passive: true })
      }
    }, 100)
  })

  onUnmounted(() => {
    const mainContent = document.querySelector('.main-content')
    if (mainContent) {
      mainContent.removeEventListener('scroll', handleScroll)
    }

    const navbar = document.querySelector('.cyber-navbar') as HTMLElement
    const footer = document.querySelector('.footer-container') as HTMLElement

    if (navbar) {
      navbar.style.removeProperty('opacity')
      navbar.style.removeProperty('transform')
      navbar.style.removeProperty('pointer-events')
      navbar.style.removeProperty('transition')
    }

    if (footer) {
      footer.style.removeProperty('opacity')
      footer.style.removeProperty('transform')
      footer.style.removeProperty('pointer-events')
      footer.style.removeProperty('transition')
    }
  })
</script>

<template>
  <div class="tool-home relative h-screen overflow-hidden">
    <div class="main-container relative z-10 flex min-h-screen items-center justify-center">
      <div class="hero-content mx-auto max-w-4xl px-4 text-center">
        <HeroTitle />

        <StatsDisplay />

        <ActionButtons />
      </div>
    </div>

    <div class="footer-container fixed bottom-0 left-0 right-0 z-10">
      <CyberCopyright />
    </div>
  </div>
</template>

<style scoped>
  .tool-home {
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .tool-home::-webkit-scrollbar {
    display: none;
  }

  .tool-home *::-webkit-scrollbar {
    opacity: 0;
    width: 0;
    height: 0;
    background: transparent;
  }

  .tool-home * {
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .tool-home::-webkit-scrollbar-track {
    background: transparent;
    opacity: 0;
  }

  .tool-home::-webkit-scrollbar-thumb {
    background: transparent;
    opacity: 0;
  }

  .tool-home::-webkit-scrollbar-corner {
    background: transparent;
    opacity: 0;
  }

  .main-container {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  @media (max-width: 768px) {
    .hero-content {
      padding: 0 1.5rem;
    }
  }

  @media (max-width: 480px) {
    .hero-content {
      padding: 0 1rem;
    }
  }
</style>
