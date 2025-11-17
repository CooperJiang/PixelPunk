<script setup lang="ts">
  import { computed, ref, onMounted } from 'vue'
  import { useLayoutStore } from '@/store/layout'
  import { useAppearanceSettings } from '@/composables/useAppearanceSettings'
  import { useTexts } from '@/composables/useTexts'
  import BrandLogo from '@/components/BrandLogo/index.vue'
  import SidebarNav from './SidebarNav.vue'
  import type { CommunityDialogExpose } from '@/components/CommunityDialog/types'

  defineOptions({
    name: 'CyberSidebar',
  })

  const { $t } = useTexts()
  const layoutStore = useLayoutStore()
  const { socialLinks, layoutSettings, loadSettings } = useAppearanceSettings()
  const communityDialogRef = ref<CommunityDialogExpose>()
  const qqGroupDialogRef = ref<CommunityDialogExpose>()

  const sidebarStyle = computed(() => ({
    width: `${layoutStore.currentSidebarWidth}px`,
    transform: layoutStore.isLeftLayout ? 'translateX(0)' : 'translateX(-100%)',
  }))

  const isCollapsed = computed(() => layoutStore.sidebarCollapsed)

  const buttonLayout = computed(() => {
    let count = 1 // 折叠/展开按钮（始终存在）

    if (layoutSettings.value.multiLayoutEnabled) count++

    if (socialLinks.value.official.show) count++
    if (socialLinks.value.github.show) count++
    if (socialLinks.value.wechat.show && socialLinks.value.wechat.hasConfig) count++
    if (socialLinks.value.qq.show && socialLinks.value.qq.hasConfig) count++

    if (count <= 5) {
      return {
        totalCount: count,
        rows: 1,
        firstRowCount: count,
        secondRowCount: 0,
        gridColumns: count,
      }
    } else {
      const firstRowCount = Math.ceil(count / 2) // 第一排多一个（如果奇数）
      const secondRowCount = count - firstRowCount
      return {
        totalCount: count,
        rows: 2,
        firstRowCount,
        secondRowCount,
        gridColumns: Math.max(firstRowCount, secondRowCount),
      }
    }
  })

  const openGitHub = () => {
    window.open(socialLinks.value.github.url, '_blank')
  }

  const openOfficialSite = () => {
    window.open(socialLinks.value.official.url, '_blank')
  }

  const showCommunityDialog = () => {
    communityDialogRef.value?.show()
  }

  const showQQGroupDialog = () => {
    qqGroupDialogRef.value?.show()
  }

  onMounted(async () => {
    await loadSettings()
  })
</script>

<template>
  <div
    class="cyber-sidebar"
    :class="{
      'sidebar-collapsed': isCollapsed,
      'sidebar-expanded': !isCollapsed,
    }"
    :style="sidebarStyle"
  >
    <div class="sidebar-background">
      <div class="sidebar-glow" />
      <div class="sidebar-grid" />
    </div>

    <div class="sidebar-content">
      <div class="sidebar-header">
        <BrandLogo :show-version="!isCollapsed" :collapsed="isCollapsed" />
      </div>

      <div class="sidebar-nav">
        <SidebarNav :collapsed="isCollapsed" />
      </div>

      <div class="sidebar-footer">
        <div
          class="sidebar-tools"
          :style="{
            gridTemplateColumns: isCollapsed ? '1fr' : `repeat(${buttonLayout.gridColumns}, 1fr)`,
          }"
        >
          <div v-if="layoutSettings.multiLayoutEnabled" class="layout-toggle-wrapper">
            <CyberLayoutToggleButton
              size="sm"
              variant="icon"
              :tooltip="true"
              :tooltip-placement="isCollapsed ? 'right' : 'top'"
            />
          </div>

          <div class="social-links" :class="{ 'social-collapsed': isCollapsed }">
            <CyberIconButton
              v-if="socialLinks.official.show"
              type="cyber"
              size="small"
              :tooltip="$t('components.sidebar.social.visitWebsite')"
              :tooltip-placement="isCollapsed ? 'right' : 'top'"
              @click="openOfficialSite"
            >
              <i class="fas fa-globe" />
            </CyberIconButton>

            <CyberIconButton
              v-if="socialLinks.github.show"
              type="cyber"
              size="small"
              :tooltip="$t('components.sidebar.social.visitGitHub')"
              :tooltip-placement="isCollapsed ? 'right' : 'top'"
              @click="openGitHub"
            >
              <i class="fab fa-github" />
            </CyberIconButton>

            <CyberIconButton
              v-if="socialLinks.wechat.show && socialLinks.wechat.hasConfig"
              type="cyber"
              size="small"
              :tooltip="$t('components.sidebar.social.joinWechat')"
              :tooltip-placement="isCollapsed ? 'right' : 'top'"
              @click="showCommunityDialog"
            >
              <i class="fab fa-weixin" />
            </CyberIconButton>

            <CyberIconButton
              v-if="socialLinks.qq.show && socialLinks.qq.hasConfig"
              type="cyber"
              size="small"
              :tooltip="$t('components.sidebar.social.joinQQ')"
              :tooltip-placement="isCollapsed ? 'right' : 'top'"
              @click="showQQGroupDialog"
            >
              <i class="fab fa-qq" />
            </CyberIconButton>
          </div>

          <CyberIconButton
            type="cyber"
            size="small"
            :tooltip="isCollapsed ? $t('components.sidebar.toggle.expand') : $t('components.sidebar.toggle.collapse')"
            :tooltip-placement="isCollapsed ? 'right' : 'top'"
            @click="layoutStore.toggleSidebar"
          >
            <i :class="isCollapsed ? 'fas fa-angle-right' : 'fas fa-angle-left'" />
          </CyberIconButton>
        </div>
      </div>
    </div>
  </div>

  <CyberCommunityDialog
    v-if="socialLinks.wechat.show && socialLinks.wechat.hasConfig"
    type="wechat"
    :qr-image="socialLinks.wechat.qrImage"
    :contact-info="socialLinks.wechat.account"
    ref="communityDialogRef"
  />
  <CyberCommunityDialog
    v-if="socialLinks.qq.show && socialLinks.qq.hasConfig"
    type="qq"
    :qr-image="socialLinks.qq.qrImage"
    :contact-info="socialLinks.qq.groupNumber"
    ref="qqGroupDialogRef"
  />
</template>

<style scoped>
  .cyber-sidebar {
    @apply fixed left-0 top-0 h-full text-content;
    background: rgba(var(--color-background-800-rgb), 0.95);
    backdrop-filter: blur(10px);
    z-index: 1000;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    border-right: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    box-shadow:
      2px 0 20px rgba(var(--color-brand-500-rgb), 0.15),
      inset -1px 0 0 rgba(var(--color-brand-500-rgb), 0.1);
    box-sizing: border-box;
  }

  .sidebar-background {
    @apply absolute inset-0;
    pointer-events: none;
    z-index: 0;
  }

  .sidebar-glow {
    @apply absolute inset-0;
    background: radial-gradient(ellipse at 50% 0%, rgba(var(--color-brand-500-rgb), 0.12) 0%, transparent 60%);
    animation: sidebarGlowPulse 4s ease-in-out infinite alternate;
  }

  @keyframes sidebarGlowPulse {
    0% {
      opacity: 0.6;
    }
    100% {
      opacity: 1;
    }
  }

  .sidebar-grid {
    @apply absolute inset-0;
    background-image:
      linear-gradient(rgba(var(--color-brand-500-rgb), 0.08) 1px, transparent 1px),
      linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.08) 1px, transparent 1px);
    background-size: 30px 30px;
    opacity: 0.4;
    animation: sidebarGridMove 20s linear infinite;
  }

  @keyframes sidebarGridMove {
    0% {
      background-position: 0 0;
    }
    100% {
      background-position: 30px 30px;
    }
  }

  .sidebar-content {
    @apply relative flex h-full flex-col;
    z-index: 1;
  }

  .sidebar-header {
    @apply flex items-center justify-center px-4;
    height: 64px;
    position: relative;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    background: linear-gradient(180deg, rgba(var(--color-brand-500-rgb), 0.05) 0%, transparent 100%);
  }

  .sidebar-header::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--color-brand-500), transparent);
    opacity: 0.3;
  }

  .sidebar-nav {
    @apply flex-1 py-4;
    overflow-y: auto;
  }

  .sidebar-collapsed .sidebar-nav {
    @apply py-2;
  }

  .sidebar-footer {
    @apply px-3 py-3;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    background: linear-gradient(180deg, transparent 0%, rgba(var(--color-brand-500-rgb), 0.03) 100%);
    position: relative;
  }

  .sidebar-tools {
    display: grid;
    gap: 6px;
    width: 100%;
    justify-items: center;
    align-items: center;
  }

  .social-links {
    display: contents;
  }

  .social-collapsed {
    display: contents;
  }

  .layout-toggle-wrapper {
    display: flex;
    justify-content: center;
  }

  .sidebar-collapsed .sidebar-footer {
    @apply px-2;
  }

  .sidebar-collapsed .sidebar-tools {
    display: grid;
    grid-template-columns: 1fr;
    gap: 6px;
    align-items: center;
    justify-items: center;
  }

  @media (max-width: 768px) {
    .sidebar-footer {
      @apply px-2 py-2;
    }

    .sidebar-tools {
      @apply flex-col gap-1;
    }
  }

  .sidebar-collapsed {
    .sidebar-header {
      @apply px-2;
    }
  }

  @media (max-width: 768px) {
    .cyber-sidebar {
      transform: translateX(-100%);
      box-shadow: none;
    }

    .cyber-sidebar.mobile-open {
      transform: translateX(0);

      box-shadow: none;
    }
  }
</style>
