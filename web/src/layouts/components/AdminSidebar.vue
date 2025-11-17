<script setup lang="ts">
  import { useRoute } from 'vue-router'
  import { useTexts } from '@/composables/useTexts'
  import BrandLogo from '@/components/BrandLogo/index.vue'
  import ThemeToggle from '@/components/ThemeToggle/index.vue'
  import type { MenuItem } from '@/types'

  const { $t } = useTexts()

  interface Props {
    menuItems: MenuItem[]
    isCollapsed: boolean
    isMobile: boolean
    isMobileOpen: boolean
    siteName: string
  }

  defineProps<Props>()

  const emit = defineEmits<{
    'toggle-sidebar': []
    'close-mobile': []
    'route-change': []
  }>()

  const route = useRoute()
</script>

<template>
  <aside
    class="admin-sidebar"
    :class="{
      collapsed: isCollapsed,
      'mobile-open': isMobile && isMobileOpen,
    }"
  >
    <div class="sidebar-header">
      <BrandLogo size="sm" :collapsed="isCollapsed && !isMobile" />

      <button v-if="isMobile" class="mobile-close-btn" @click="emit('close-mobile')">
        <i class="fas fa-times" />
      </button>
    </div>

    <nav class="sidebar-menu">
      <template v-for="(item, index) in menuItems" :key="index">
        <div v-if="isCollapsed && !isMobile" class="menu-wrapper">
          <CyberTooltip :content="item.label" placement="right" :offset="20">
            <router-link
              :to="item.path"
              class="menu-item"
              :class="{ active: route.path.includes(item.path) }"
              @click="emit('route-change')"
            >
              <i class="menu-icon fas" :class="item.icon" />
            </router-link>
          </CyberTooltip>
        </div>

        <div v-else class="menu-wrapper">
          <router-link
            :to="item.path"
            class="menu-item"
            :class="{ active: route.path.includes(item.path) }"
            @click="emit('route-change')"
          >
            <i class="menu-icon fas" :class="item.icon" />
            <span class="menu-label">{{ item.label }}</span>
          </router-link>
        </div>
      </template>
    </nav>

    <!-- 底部工具栏 -->
    <div v-if="!isMobile" class="sidebar-footer">
      <!-- 折叠/展开按钮 -->
      <div class="footer-actions">
        <!-- 展开菜单按钮 - 折叠时显示 -->
        <CyberTooltip v-if="isCollapsed" :content="$t('layouts.sidebar.expandMenu')" placement="right" :offset="20">
          <button class="footer-icon-btn" @click="emit('toggle-sidebar')">
            <i class="fas fa-angles-right" />
          </button>
        </CyberTooltip>

        <!-- 收起菜单按钮 - 展开时显示 -->
        <CyberTooltip v-if="!isCollapsed" :content="$t('layouts.sidebar.collapseMenu')" placement="top" :offset="8">
          <button class="footer-icon-btn" @click="emit('toggle-sidebar')">
            <i class="fas fa-angles-left" />
          </button>
        </CyberTooltip>

        <!-- 主题切换 - 仅在展开状态显示 -->
        <div v-if="!isCollapsed" class="theme-toggle-wrapper">
          <ThemeToggle size="small" placement="top" align="left" />
        </div>

        <!-- AI配置按钮 - 仅在展开状态显示 -->
        <CyberTooltip v-if="!isCollapsed" :content="$t('layouts.sidebar.aiConfig')" placement="top" :offset="8">
          <button class="footer-icon-btn" @click="$router.push('/admin/ai')">
            <i class="fas fa-robot" />
          </button>
        </CyberTooltip>

        <!-- 网站配置按钮 - 仅在展开状态显示 -->
        <CyberTooltip v-if="!isCollapsed" :content="$t('layouts.sidebar.websiteConfig')" placement="top" :offset="8">
          <button class="footer-icon-btn" @click="$router.push('/admin/construction?tab=website')">
            <i class="fas fa-globe" />
          </button>
        </CyberTooltip>

        <!-- 设置按钮 - 仅在展开状态显示 -->
        <CyberTooltip v-if="!isCollapsed" :content="$t('layouts.sidebar.globalSettings')" placement="top" :offset="8">
          <button class="footer-icon-btn" @click="$router.push('/admin/settings')">
            <i class="fas fa-cog" />
          </button>
        </CyberTooltip>
      </div>
    </div>
  </aside>
</template>

<style scoped lang="scss">
  .admin-sidebar {
    width: 240px;
    background: var(--color-background-900);
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 1050;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    border-right: 1px solid var(--color-border-default);
    box-shadow:
      2px 0 8px rgba(0, 0, 0, 0.08),
      4px 0 16px rgba(0, 0, 0, 0.04);

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 1px;
      background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.3), transparent);
      opacity: 0.5;
    }
  }

  .admin-sidebar.collapsed {
    width: 64px;

    .sidebar-footer {
      padding: 0.75rem 0.5rem;
    }

    .footer-actions {
      flex-direction: column;
      gap: 0.5rem;
    }

    .sidebar-menu {
      gap: 0.375rem;
    }

    .menu-wrapper {
      display: flex;
      justify-content: center;
    }

    .menu-item {
      width: 42px;
      height: 42px;
      justify-content: center;
      padding: 0;
      border-radius: var(--radius-sm);
      border: none;

      .menu-icon {
        margin: 0;
        font-size: 0.875rem;
      }

      &:hover {
        background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.08));
        border: none;
        transform: translateY(-1px);

        &::before {
          width: 0;
          height: 0;
        }

        .menu-icon {
          color: var(--color-brand-400);
        }
      }

      &.active {
        background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.12));
        border: none;

        &::before {
          width: 0;
          height: 0;
        }

        .menu-icon {
          color: var(--color-brand-500);
        }
      }
    }
  }

  @media (max-width: 768px) {
    .admin-sidebar {
      transform: translateX(-100%);
      width: 280px;
      max-width: 280px;
      box-shadow: var(--shadow-cyber-lg);
    }

    .admin-sidebar.mobile-open {
      transform: translateX(0);
    }
  }

  .sidebar-header {
    padding: 0.5rem 0.75rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-bottom: 1px solid var(--color-border-subtle);
    background: linear-gradient(
      180deg,
      rgba(var(--color-background-900-rgb), 0.5) 0%,
      rgba(var(--color-background-800-rgb), 0.3) 100%
    );
    min-height: 48px;
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 1px;
      background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.2), transparent);
    }
  }

  @media (max-width: 768px) {
    .sidebar-header {
      justify-content: space-between;
      padding: 1rem 1.25rem;
    }
  }

  .mobile-close-btn {
    background: transparent;
    border: 1px solid var(--color-border-default);
    color: var(--color-content-default);
    cursor: pointer;
    padding: 0.5rem;
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;

    &:hover {
      background: rgba(var(--color-error-rgb), 0.1);
      border-color: rgba(var(--color-error-rgb), 0.3);
      color: var(--color-error-500);
    }

    &:active {
      transform: scale(0.95);
    }
  }

  .sidebar-menu {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding: 1rem 0.75rem;
    overflow-y: auto;
    flex: 1;
    scrollbar-width: thin;
    scrollbar-color: rgba(var(--color-brand-500-rgb), 0.3) transparent;

    &::-webkit-scrollbar {
      width: 4px;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      background: rgba(var(--color-brand-500-rgb), 0.3);
      border-radius: var(--radius-sm);

      &:hover {
        background: rgba(var(--color-brand-500-rgb), 0.5);
      }
    }
  }

  .menu-wrapper {
    width: 100%;
    position: relative;
  }

  .menu-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    border-radius: var(--radius-sm);
    color: var(--color-content-default);
    text-decoration: none;
    font-size: 0.875rem;
    font-weight: 500;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
    background: transparent;
    border: none;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 0;
      height: 0;
      background: linear-gradient(180deg, var(--color-brand-500), var(--color-brand-400));
      border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    &:hover {
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08), rgba(var(--color-brand-500-rgb), 0.04));
      color: var(--color-brand-400);
      border: none;
      transform: translateX(2px);

      &::before {
        width: 3px;
        height: 60%;
      }

      .menu-icon {
        color: var(--color-brand-400);
        transform: scale(1.1);
      }
    }

    &.active {
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-brand-500-rgb), 0.08));
      color: var(--color-brand-500);
      border: none;
      font-weight: 600;

      &::before {
        width: 3px;
        height: 70%;
        box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.4);
      }

      .menu-icon {
        color: var(--color-brand-500);
      }
    }

    &:focus-visible {
      outline: none;
      box-shadow: 0 0 0 3px rgba(var(--color-brand-500-rgb), 0.15);
    }
  }

  .menu-icon {
    width: 18px;
    text-align: center;
    color: var(--color-content-muted);
    font-size: 1rem;
    flex-shrink: 0;
    transition: all 0.2s ease;
  }

  .menu-label {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    letter-spacing: 0.01em;
  }

  .sidebar-footer {
    padding: 0.75rem;
    border-top: 1px solid var(--color-border-subtle);
    background: linear-gradient(
      180deg,
      rgba(var(--color-background-800-rgb), 0.3) 0%,
      rgba(var(--color-background-900-rgb), 0.5) 100%
    );
    position: relative;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 1px;
      background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.2), transparent);
    }
  }

  .footer-actions {
    display: flex;
    justify-content: space-between;
    gap: 0.5rem;
    align-items: center;
  }

  .footer-icon-btn {
    background: transparent;
    border: 1px solid var(--color-border-default);
    color: var(--color-content-default);
    cursor: pointer;
    padding: 0.5rem;
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    flex-shrink: 0;

    &:hover {
      background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-brand-500-rgb), 0.05));
      border-color: rgba(var(--color-brand-500-rgb), 0.3);
      color: var(--color-brand-400);
      transform: translateY(-1px);
      box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.15);
    }

    &:active {
      transform: translateY(0);
    }

    i {
      font-size: 0.875rem;
    }
  }

  .theme-toggle-wrapper {
    display: flex;
    align-items: center;
    flex-shrink: 0;
  }
</style>
