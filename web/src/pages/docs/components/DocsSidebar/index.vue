<script setup lang="ts">
  import { computed, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'DocsSidebar',
  })

  const props = defineProps<Props>()

  const progressStyle = computed(() => {
    const width = `${Math.max(0, Math.min(100, props.readingProgress || 0))}%`
    return {
      width,
    }
  })

  interface NavItem {
    id: string
    label: string
    iconClass: string
  }

  interface Props {
    navItems: NavItem[]
    activeSection: string
    readingProgress: number
    scrollToSection: (event: Event) => void
  }

  const onScrollToSection = (event: Event) => {
    props.scrollToSection(event)
  }
</script>

<template>
  <aside class="docs-sidebar">
    <div class="sidebar-content">
      <div class="sidebar-header">
        <router-link to="/" class="header-brand-link">
          <div class="brand-section">
            <div class="brand-icon-wrapper">
              <i class="fas fa-arrow-left back-icon" />
              <i class="fas fa-code main-icon" />
            </div>
            <div class="brand-text">
              <span class="brand-title">{{ $t('docs.sidebar.title') }}</span>
              <span class="brand-subtitle">{{ $t('docs.sidebar.backHome') }}</span>
            </div>
          </div>
        </router-link>
      </div>

      <div class="sidebar-nav-section">
        <div class="nav-header">
          <div class="nav-title-row">
            <h3><i class="fas fa-list-ul" />{{ $t('docs.sidebar.navigation') }}</h3>
            <div class="progress-percentage">{{ Math.round(readingProgress) }}%</div>
          </div>
          <div class="progress-container">
            <div class="progress-bar">
              <div class="progress-fill" :style="progressStyle" :data-progress="readingProgress" />
            </div>
          </div>
        </div>

        <nav class="sidebar-nav">
          <a
            v-for="item in navItems"
            :key="item.id"
            :href="`#${item.id}`"
            class="sidebar-link"
            :class="[{ active: activeSection === item.id }]"
            @click="onScrollToSection"
          >
            <div class="link-content">
              <div class="link-icon-wrapper">
                <i :class="item.iconClass" />
              </div>
              <span class="link-text">{{ item.label }}</span>
            </div>
            <div v-if="activeSection === item.id" class="active-indicator" />
          </a>
        </nav>

        <div class="sidebar-footer">
          <div class="api-info">
            <div class="info-item">
              <i class="fas fa-shield-alt" />
              <span>API v1.0</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>

<style scoped>
  .docs-sidebar {
    width: 280px;
    height: 100vh;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.98) 50%,
      rgba(var(--color-background-900-rgb), 0.95) 100%
    );
    border-right: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    position: fixed;
    top: 0;
    left: 0;
    z-index: 100;
    overflow: hidden;
  }

  .sidebar-content {
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    position: relative;
    z-index: 1;
  }

  .sidebar-header {
    flex-shrink: 0;
    padding: 1rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.05) 0%, transparent 100%);
  }

  .header-brand-link {
    display: block;
    text-decoration: none;
    transition: all 0.3s ease;
  }

  .header-brand-link:hover .brand-section {
    transform: translateY(-2px);
  }

  .brand-section {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.2) 0%, rgba(255, 255, 255, 0.02) 100%);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
  }

  .brand-icon-wrapper {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.2) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.2) 100%
    );
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .back-icon {
    position: absolute;
    color: var(--color-brand-500);
    font-size: 0.7rem;
    opacity: 0;
    transform: translateX(-8px);
    transition: all 0.3s ease;
  }

  .main-icon {
    color: var(--color-brand-500);
    font-size: 1rem;
    transition: all 0.3s ease;
  }

  .header-brand-link:hover .back-icon {
    opacity: 1;
    transform: translateX(-4px);
  }

  .header-brand-link:hover .main-icon {
    opacity: 0.3;
    transform: translateX(4px);
  }

  .brand-text {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
  }

  .brand-title {
    font-size: 0.95rem;
    font-weight: 700;
    color: var(--color-text-content);
    transition: color 0.3s ease;
  }

  .brand-subtitle {
    font-size: 0.65rem;
    font-weight: 500;
    color: var(--color-slate-500);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    opacity: 0.8;
    transition: all 0.3s ease;
  }

  .header-brand-link:hover .brand-title {
    color: var(--color-brand-500);
  }

  .header-brand-link:hover .brand-subtitle {
    color: var(--color-brand-500);
    opacity: 1;
  }

  .sidebar-nav-section {
    flex: 1;
    padding: 1rem;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
  }

  .nav-header {
    margin-bottom: 1rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .nav-title-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .nav-header h3 {
    margin: 0;
    font-size: 0.8rem;
    font-weight: 700;
    color: var(--color-text-content);
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .nav-header h3 i {
    font-size: 0.7rem;
    color: var(--color-brand-500);
  }

  .progress-percentage {
    font-size: 0.7rem;
    color: var(--color-brand-500);
    font-weight: 700;
    background: rgba(var(--color-brand-500-rgb), 0.3);
    padding: 0.15rem 0.4rem;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .progress-container {
    background: rgba(var(--color-background-700-rgb), 0.95);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    border-radius: var(--radius-sm);
    padding: 0.5rem;
    box-shadow:
      inset 0 1px 3px rgba(0, 0, 0, 0.1),
      0 1px 2px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .progress-bar {
    width: 100%;
    height: 6px;
    background: rgba(var(--color-background-700-rgb), 0.9);
    border-radius: var(--radius-sm);
    overflow: hidden;
    position: relative;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .progress-fill {
    height: 100%;
    display: block;
    background: linear-gradient(90deg, var(--color-brand-500) 0%, var(--color-error-500) 100%);
    border-radius: var(--radius-sm);
    transition: width 0.3s ease-out;
    min-width: 2px;
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .progress-fill[data-progress]:not([data-progress='0']) {
    min-width: 4px;
    opacity: 1;
  }

  .sidebar-nav {
    flex: 1;
    margin-bottom: 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .sidebar-link {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.6rem 0.75rem;
    color: var(--color-slate-500);
    text-decoration: none;
    border-radius: var(--radius-sm);
    font-size: 0.8rem;
    font-weight: 500;
    position: relative;
    overflow: hidden;
    border: 1px solid transparent;
    background: transparent;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .sidebar-link:hover,
  .sidebar-link.active {
    color: var(--color-text-content);
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.2) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.05) 100%
    );
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
    transform: translateX(4px);
  }

  .sidebar-link.active {
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.3) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.08) 100%
    );
  }

  .link-content {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    z-index: 2;
    position: relative;
  }

  .link-icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    transition: all 0.3s ease;
    flex-shrink: 0;
  }

  .link-icon-wrapper i {
    color: var(--color-brand-500);
    font-size: 0.75rem;
    transition: all 0.3s ease;
  }

  .sidebar-link:hover .link-icon-wrapper,
  .sidebar-link.active .link-icon-wrapper {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: scale(1.05);
  }

  .link-text {
    font-weight: 600;
    transition: color 0.3s ease;
    flex: 1;
  }

  .active-indicator {
    width: 6px;
    height: 6px;
    background: linear-gradient(45deg, var(--color-brand-500), var(--color-error-500));
    border-radius: var(--radius-full);
    flex-shrink: 0;
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.6);
  }

  .sidebar-footer {
    flex-shrink: 0;
    padding-top: 1rem;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .api-info {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .info-item {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    font-size: 0.7rem;
    color: var(--color-gray-200);
    font-weight: 600;
    padding: 0.4rem 0.6rem;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.05) 0%, rgba(var(--color-success-rgb), 0.05) 100%);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    transition: all 0.3s ease;
  }

  .info-item i {
    color: var(--color-brand-500);
    font-size: 0.75rem;
  }

  .sidebar-nav-section::-webkit-scrollbar {
    width: 4px;
  }

  .sidebar-nav-section::-webkit-scrollbar-track {
    background: var(--color-hover-bg-neutral);
    border-radius: var(--radius-sm);
  }

  .sidebar-nav-section::-webkit-scrollbar-thumb {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.4) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.4) 100%
    );
    border-radius: var(--radius-sm);
    transition: background 0.3s ease;
  }

  .sidebar-nav-section::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(
      180deg,
      rgba(var(--color-brand-500-rgb), 0.6) 0%,
      rgba(var(--color-badge-accent-text-rgb), 0.6) 100%
    );
  }
</style>
