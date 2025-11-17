<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'MobileHeader',
  })

  const props = defineProps<Props>()

  interface NavItem {
    id: string
    label: string
    iconClass: string
  }

  interface Props {
    navItems: NavItem[]
    activeSection: string
    mobileMenuOpen: boolean
    toggleMobileMenu: () => void
    handleMobileNavClick: (event: Event) => void
  }

  const onToggleMobileMenu = () => {
    props.toggleMobileMenu()
  }

  const onMobileNavClick = (event: Event) => {
    props.handleMobileNavClick(event)
  }
</script>

<template>
  <div class="mobile-header mobile-only">
    <div class="mobile-header-content">
      <router-link to="/" class="back-home-btn">
        <i class="fas fa-arrow-left" />
        <span>{{ $t('docs.mobile.backHome') }}</span>
      </router-link>
      <div class="header-brand">
        <i class="fas fa-code" />
        <span>{{ $t('docs.mobile.title') }}</span>
      </div>
      <button class="mobile-menu-toggle" @click="onToggleMobileMenu">
        <i :class="mobileMenuOpen ? 'fas fa-times' : 'fas fa-bars'" />
      </button>
    </div>

    <div class="mobile-menu" :class="[{ open: mobileMenuOpen }]">
      <a
        v-for="item in navItems"
        :key="item.id"
        :href="`#${item.id}`"
        class="mobile-nav-link"
        :class="[{ active: activeSection === item.id }]"
        @click="onMobileNavClick"
      >
        <i :class="item.iconClass" />
        <span>{{ item.label }}</span>
      </a>
    </div>
  </div>
</template>

<style scoped>
  .mobile-header {
    position: sticky;
    top: 0;
    background: rgba(var(--color-background-900-rgb), 0.95);
    backdrop-filter: blur(25px);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    z-index: 90;
    box-shadow: 0 2px 10px var(--color-background-700);
  }

  .mobile-header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem;
  }

  .back-home-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.8rem;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    color: var(--color-content-muted);
    text-decoration: none;
    font-size: 0.8rem;
    font-weight: 500;
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
  }

  .back-home-btn:hover {
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-1px);
    box-shadow: 0 4px 15px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .back-home-btn i {
    font-size: 0.8rem;
    transition: transform 0.3s ease;
  }

  .back-home-btn:hover i {
    transform: translateX(-2px);
  }

  .header-brand {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 1rem;
    font-weight: 700;
    color: var(--color-brand-500);
    text-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .header-brand i {
    font-size: 1rem;
    animation: pulse-rotate 4s ease-in-out infinite;
  }

  @keyframes pulse-rotate {
    0%,
    100% {
      transform: rotate(0deg) scale(1);
    }
    50% {
      transform: rotate(180deg) scale(1.1);
    }
  }

  .mobile-menu-toggle {
    display: flex;
    background: none;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    padding: 0.5rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .mobile-menu-toggle:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .mobile-menu {
    background: rgba(var(--color-background-900-rgb), 0.98);
    backdrop-filter: blur(12px);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    max-height: 0;
    overflow: hidden;
    transition: max-height 0.3s ease;
    padding: 0 1rem;
  }

  .mobile-menu.open {
    max-height: 400px;
    padding: 1rem;
  }

  .mobile-nav-link {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    color: var(--color-content-muted);
    text-decoration: none;
    transition: all 0.3s ease;
    border-radius: var(--radius-sm);
    margin-bottom: 0.25rem;
    font-size: 0.9rem;
    font-weight: 500;
  }

  .mobile-nav-link:hover,
  .mobile-nav-link.active {
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.2);
    transform: translateX(8px);
  }

  .mobile-nav-link i {
    font-size: 1rem;
    width: 20px;
    text-align: center;
  }

  .desktop-only {
    display: block;
  }

  .mobile-only {
    display: none;
  }

  @media (max-width: 768px) {
    .desktop-only {
      display: none;
    }

    .mobile-only {
      display: block;
    }
  }
</style>
