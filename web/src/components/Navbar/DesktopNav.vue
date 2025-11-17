<script setup lang="ts">
  import { useNavigation } from '@/composables/useNavigation'
  import { useAuthStore } from '@/store/auth'
  import { computed } from 'vue'

  const { navigationItems: allNavItems, isRouteActive } = useNavigation()
  const authStore = useAuthStore()

  /* 过滤显示的导航项（根据认证状态） */
  const navigationItems = computed(() =>
    allNavItems.value.filter((item) => {
      if (item.requiresAuth) {
        return authStore.isLoggedIn
      }
      return true
    })
  )

  const getNavItemClass = (item: { name: string; path: string }) => {
    const classes = []

    if (item.path === '/upload') {
      classes.push('nav-upload')
    } else if (item.path === '/hive') {
      classes.push('nav-hive')
    } else if (item.path === '/docs') {
      classes.push('nav-docs')
    }

    return classes
  }
</script>

<template>
  <div class="navbar-nav desktop-only">
    <div class="nav-items">
      <div class="nav-group primary">
        <template v-for="item in navigationItems" :key="item.path">
          <a
            v-if="item.external"
            :href="item.path"
            :target="item.target || '_self'"
            class="nav-item"
            :class="getNavItemClass(item)"
          >
            <i class="nav-icon" :class="item.icon" />
            <span>{{ item.name }}</span>
            <i v-if="item.external" class="fas fa-external-link-alt external-icon" />
          </a>

          <router-link
            v-else
            :to="item.path"
            class="nav-item"
            :class="[getNavItemClass(item), { active: isRouteActive(item.path) }]"
          >
            <i class="nav-icon" :class="item.icon" />
            <span>{{ item.name }}</span>
          </router-link>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .navbar-nav {
    display: flex;
    justify-content: flex-start;
    margin-left: 3rem;
    flex: 1;
  }

  .nav-items {
    display: flex;
    align-items: center;
  }

  .nav-group {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .nav-item {
    position: relative;
    color: var(--color-content-heading);
    text-decoration: none;
    border-radius: var(--radius-sm);
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    background: transparent;
    border: none;
    cursor: pointer;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    padding: 0.5rem 0.875rem;
    font-size: 0.8rem;
    font-weight: 600;
    letter-spacing: 0.025em;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    white-space: nowrap;
    overflow: hidden;
  }

  .nav-icon {
    font-size: 0.8rem;
    width: 14px;
    text-align: center;
    transition: all 0.25s ease;
    opacity: 0.85;
  }

  .external-icon {
    font-size: 0.65rem;
    margin-left: 0.25rem;
    opacity: 0.5;
  }

  .nav-item::after {
    content: '';
    position: absolute;
    inset: 0;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    opacity: 0;
    transition: opacity 0.25s ease;
    border-radius: var(--radius-sm);
  }

  .nav-item:hover {
    color: var(--color-content-default);
    transform: translateY(-1px);
  }

  .nav-item:hover::after {
    opacity: 1;
  }

  .nav-item:hover .nav-icon {
    opacity: 1;
    transform: scale(1.05);
    color: var(--color-brand-500);
  }

  .nav-item.active {
    color: var(--color-brand-500);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .nav-item.active::after {
    opacity: 1;
    background: rgba(var(--color-brand-500-rgb), 0.08);
  }

  .nav-item.active .nav-icon {
    color: var(--color-brand-500);
    opacity: 1;
    filter: drop-shadow(0 0 4px rgba(var(--color-brand-500-rgb), 0.6));
  }

  .nav-hive {
    color: var(--color-warning-500);
  }

  .nav-hive::after {
    background: rgba(var(--color-warning-rgb), 0.05);
  }

  .nav-hive:hover .nav-icon {
    color: var(--color-warning-500);
  }

  .nav-hive .nav-icon {
    color: var(--color-warning-500);
  }

  .nav-hive:hover .nav-icon {
    filter: drop-shadow(0 0 4px rgba(var(--color-warning-rgb), 0.6));
  }

  .nav-upload {
    color: var(--color-brand-600);
  }

  .nav-upload::after {
    background: rgba(var(--color-brand-600-rgb), 0.05);
  }

  .nav-upload:hover .nav-icon {
    color: var(--color-brand-600);
  }

  .nav-upload.active {
    color: var(--color-brand-600);
    text-shadow: 0 0 8px rgba(var(--color-brand-600-rgb), 0.4);
  }

  .nav-upload.active::after {
    background: rgba(var(--color-brand-600-rgb), 0.08);
  }

  .nav-upload .nav-icon {
    color: var(--color-brand-600);
  }

  .nav-upload:hover .nav-icon,
  .nav-upload.active .nav-icon {
    filter: drop-shadow(0 0 4px rgba(var(--color-brand-600-rgb), 0.6));
  }

  .nav-docs {
    color: var(--color-brand-400);
  }

  .nav-docs::after {
    background: rgba(var(--color-brand-400-rgb), 0.05);
  }

  .nav-docs:hover .nav-icon {
    color: var(--color-brand-400);
  }

  .nav-docs .nav-icon {
    color: var(--color-brand-400);
  }

  .nav-docs:hover .nav-icon {
    filter: drop-shadow(0 0 4px rgba(var(--color-brand-400-rgb), 0.6));
  }

  .desktop-only {
    display: block;
  }

  @media (max-width: 768px) {
    .desktop-only {
      display: none;
    }
  }

  @media (max-width: 1200px) {
    .navbar-nav {
      margin-left: 2rem;
    }

    .nav-item {
      padding: 0.5rem 0.8rem;
      font-size: 0.8rem;
    }
  }

  @media (max-width: 900px) {
    .navbar-nav {
      margin-left: 1.5rem;
    }

    .nav-item {
      padding: 0.4rem 0.6rem;
      font-size: 0.75rem;
    }
  }
</style>
