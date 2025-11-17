<script setup lang="ts">
  import { computed, ref, onMounted } from 'vue'
  import { useAuthStore } from '@/store/auth'
  import { useLayoutStore } from '@/store/layout'
  import { useTexts } from '@/composables/useTexts'
  import { useToast } from '@/components/Toast/useToast'
  import { useAppearanceSettings } from '@/composables/useAppearanceSettings'
  import { useRouter } from 'vue-router'
  import { onClickOutside } from '@vueuse/core'
  import UserAvatar from '@/components/UserAvatar/index.vue'
  import MessageNotification from '@/components/Message/MessageNotification.vue'
  import AnnouncementButton from '@/components/AnnouncementButton/index.vue'
  import LocaleSwitch from '@/components/LocaleSwitch/index.vue'
  import type { NavbarEmits } from './types'
  import type { CommunityDialogExpose } from '@/components/CommunityDialog/types'

  defineOptions({
    name: 'UserTerminal',
  })

  defineProps<{
    mobileMenuOpen: boolean
  }>()

  /* Props and Emits */
  const emit = defineEmits<NavbarEmits>()

  const authStore = useAuthStore()
  const layoutStore = useLayoutStore()
  const { $t } = useTexts()
  const toast = useToast()
  const router = useRouter()
  const { socialLinks, layoutSettings, loadSettings } = useAppearanceSettings()

  /* UI状态 */
  const userMenuOpen = ref<boolean>(false)
  const userMenuRef = ref<HTMLElement | null>(null)
  const terminalDisplayRef = ref<HTMLElement | null>(null)
  const communityDialogRef = ref<CommunityDialogExpose>()
  const qqGroupDialogRef = ref<CommunityDialogExpose>()

  /* 点击菜单外部时关闭菜单(排除触发按钮) */
  onClickOutside(userMenuRef, (event) => {
    if (terminalDisplayRef.value?.contains(event.target as Node)) {
      return
    }
    if (userMenuOpen.value) {
      userMenuOpen.value = false
    }
  })

  const userName = computed<string>(() => authStore.userInfo?.username || '')

  const isAdmin = computed<boolean>(() => authStore.isAdmin)

  const toggleUserMenu = (): void => {
    userMenuOpen.value = !userMenuOpen.value
  }

  const closeUserMenu = (): void => {
    userMenuOpen.value = false
  }

  const handleLogout = (): void => {
    authStore.logout()
    userMenuOpen.value = false
    toast.info($t('components.navbar.user.logoutSuccess'))
    router.push('/')
  }

  const toggleMobileMenu = (): void => {
    emit('toggleMobileMenu')
  }

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
  <div class="navbar-actions">
    <div v-if="!authStore.isLoggedIn" class="auth-section desktop-only">
      <div class="announcement-button-wrapper">
        <AnnouncementButton />
      </div>

      <LocaleSwitch size="small" />

      <CyberThemeToggle size="small" />

      <div v-if="layoutStore.isTopLayout" class="social-links-unauth">
        <CyberIconButton
          v-if="socialLinks.official.show"
          type="cyber"
          size="small"
          :tooltip="$t('components.navbar.social.visitWebsite')"
          tooltip-placement="bottom"
          @click="openOfficialSite"
        >
          <i class="fas fa-globe" />
        </CyberIconButton>

        <CyberIconButton v-if="socialLinks.github.show" type="cyber" size="small" @click="openGitHub">
          <i class="fab fa-github" />
        </CyberIconButton>

        <CyberIconButton
          v-if="socialLinks.wechat.show && socialLinks.wechat.hasConfig"
          type="cyber"
          size="small"
          @click="showCommunityDialog"
        >
          <i class="fab fa-weixin" />
        </CyberIconButton>

        <CyberIconButton
          v-if="socialLinks.qq.show && socialLinks.qq.hasConfig"
          type="cyber"
          size="small"
          @click="showQQGroupDialog"
        >
          <i class="fab fa-qq" />
        </CyberIconButton>
      </div>

      <button class="cyber-login-btn" @click="emit('login')">
        <i class="fas fa-plug login-icon" />
        <span class="btn-text">{{ $t('components.navbar.user.login') }}</span>
      </button>
    </div>

    <div v-if="authStore.isLoggedIn" class="user-terminal desktop-only">
      <div class="message-notification-wrapper">
        <MessageNotification />
      </div>

      <div class="announcement-button-wrapper">
        <AnnouncementButton />
      </div>

      <LocaleSwitch size="small" />

      <CyberThemeToggle size="small" />

      <CyberLayoutToggleButton v-if="layoutStore.isTopLayout && layoutSettings.multiLayoutEnabled" size="sm" />

      <div v-if="layoutStore.isTopLayout" class="social-links-top">
        <CyberIconButton
          v-if="socialLinks.official.show"
          type="cyber"
          size="small"
          :tooltip="$t('components.navbar.social.visitWebsite')"
          tooltip-placement="bottom"
          @click="openOfficialSite"
        >
          <i class="fas fa-globe" />
        </CyberIconButton>

        <CyberIconButton v-if="socialLinks.github.show" type="cyber" size="small" @click="openGitHub">
          <i class="fab fa-github" />
        </CyberIconButton>

        <CyberIconButton
          v-if="socialLinks.wechat.show && socialLinks.wechat.hasConfig"
          type="cyber"
          size="small"
          @click="showCommunityDialog"
        >
          <i class="fab fa-weixin" />
        </CyberIconButton>

        <CyberIconButton
          v-if="socialLinks.qq.show && socialLinks.qq.hasConfig"
          type="cyber"
          size="small"
          @click="showQQGroupDialog"
        >
          <i class="fab fa-qq" />
        </CyberIconButton>
      </div>

      <div ref="terminalDisplayRef" class="terminal-display" @click="toggleUserMenu">
        <div class="user-info-line">
          <UserAvatar :avatar-url="authStore.userAvatar" :username="userName" size="sm" class="user-avatar-small" />
          <div class="user-details">
            <span class="user-name">{{ userName }}</span>
            <span class="user-role" :class="{ 'admin-role': isAdmin }">{{
              isAdmin ? $t('components.navbar.user.roles.admin') : $t('components.navbar.user.roles.user')
            }}</span>
          </div>
          <i class="fas fa-chevron-down terminal-chevron" :class="{ rotated: userMenuOpen }" />
        </div>
      </div>

      <Transition name="terminal-menu">
        <div v-show="userMenuOpen" ref="userMenuRef" class="terminal-dropdown">
          <div class="compact-header">
            <div class="user-info">
              <UserAvatar :avatar-url="authStore.userAvatar" :username="userName" size="sm" class="header-avatar" />
              <div class="user-text">
                <span class="user-name">{{ userName }}</span>
                <span class="user-badge" :class="{ admin: isAdmin }">
                  {{ isAdmin ? $t('components.navbar.user.roles.admin') : $t('components.navbar.user.roles.user') }}
                </span>
              </div>
            </div>
            <button class="close-btn" @click="closeUserMenu">
              <i class="fas fa-times" />
            </button>
          </div>

          <div class="compact-content">
            <router-link v-if="isAdmin" to="/admin" class="compact-item admin-item" @click="closeUserMenu">
              <i class="fas fa-crown item-icon admin-icon" />
              <span class="item-text">{{ $t('components.navbar.navigation.admin') }}</span>
              <i class="fas fa-chevron-right arrow" />
            </router-link>

            <div class="function-grid">
              <router-link to="/dashboard" class="grid-item" @click="closeUserMenu">
                <i class="fas fa-chart-pie item-icon dashboard" />
                <span>{{ $t('components.navbar.user.dashboard') }}</span>
              </router-link>

              <router-link to="/resource" class="grid-item" @click="closeUserMenu">
                <i class="fas fa-images item-icon images" />
                <span>{{ $t('components.navbar.user.myWorks') }}</span>
              </router-link>

              <router-link to="/folders" class="grid-item" @click="closeUserMenu">
                <i class="fas fa-folder item-icon folders" />
                <span>{{ $t('components.navbar.user.folders') }}</span>
              </router-link>

              <router-link to="/shares" class="grid-item" @click="closeUserMenu">
                <i class="fas fa-share-alt item-icon shares" />
                <span>{{ $t('components.navbar.user.shares') }}</span>
              </router-link>
            </div>

            <div class="bottom-actions">
              <router-link to="/settings" class="compact-item settings-item" @click="closeUserMenu">
                <i class="fas fa-cog item-icon" />
                <span class="item-text">{{ $t('components.navbar.user.settings') }}</span>
                <i class="fas fa-chevron-right arrow" />
              </router-link>

              <button class="compact-item logout-item" @click="handleLogout">
                <i class="fas fa-sign-out-alt item-icon" />
                <span class="item-text">{{ $t('components.navbar.user.exit') }}</span>
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </div>

    <button class="mobile-menu-btn mobile-only" @click="toggleMobileMenu">
      <div class="hamburger" :class="{ active: mobileMenuOpen }">
        <span />
        <span />
        <span />
      </div>
    </button>
  </div>

  <CyberCommunityDialog
    v-if="layoutStore.isTopLayout && socialLinks.wechat.show && socialLinks.wechat.hasConfig"
    type="wechat"
    :qr-image="socialLinks.wechat.qrImage"
    :contact-info="socialLinks.wechat.account"
    ref="communityDialogRef"
  />
  <CyberCommunityDialog
    v-if="layoutStore.isTopLayout && socialLinks.qq.show && socialLinks.qq.hasConfig"
    type="qq"
    :qr-image="socialLinks.qq.qrImage"
    :contact-info="socialLinks.qq.groupNumber"
    ref="qqGroupDialogRef"
  />
</template>

<style scoped>
  .navbar-actions {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-shrink: 0;
    min-width: 0;
    flex: 1;
    justify-content: flex-end;
    position: relative;
  }

  .auth-section {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .cyber-login-btn {
    background: transparent;
    color: var(--color-brand-500);
    border-radius: var(--radius-sm);
    padding: 0.4rem 0.7rem;
    font-size: 0.75rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.25s ease;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    position: relative;
    overflow: hidden;
    letter-spacing: 0.025em;
    text-transform: none;
    display: flex;
    align-items: center;
    gap: 0.4rem;
    backdrop-filter: blur(8px);
    border: none;
  }

  .login-icon {
    font-size: 0.7rem;
    opacity: 0.9;
  }

  .cyber-login-btn:hover {
    color: var(--color-white);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-login-btn:hover .login-icon {
    transform: rotate(15deg);
  }

  .user-terminal {
    position: relative;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    max-width: fit-content;
    flex-shrink: 0;
  }

  .layout-toggle-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    background: rgba(var(--color-content-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    color: var(--color-brand-500);
    cursor: pointer;
    transition: all 0.25s ease;
    font-size: 0.85rem;
    flex-shrink: 0;
  }

  .layout-toggle-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-white);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .announcement-button-wrapper {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    z-index: 1001;
  }

  .message-notification-wrapper {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    z-index: 1001;
  }

  .social-links-top {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }

  .social-links-unauth {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }

  .terminal-display {
    background: transparent;
    border: none;
    border-radius: var(--radius-sm);
    padding: 0.5rem 0.8rem;
    cursor: pointer;
    transition: all 0.3s ease;
    color: var(--color-text-secondary);
    min-width: 140px;
    max-width: 200px;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }

  .terminal-display:hover {
    background: rgba(var(--color-brand-500-rgb), 0.08);
    transform: translateY(-1px);
  }

  .user-info-line {
    display: flex;
    align-items: center;
    gap: 0.6rem;
  }

  .user-avatar-small {
    flex-shrink: 0;
  }

  .user-details {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    flex: 1;
    min-width: 0;
    overflow: hidden;
  }

  .user-name {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--color-white);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .user-role {
    font-size: 0.7rem;
    font-weight: 500;
    color: var(--color-content-muted);
    background: rgba(var(--color-content-rgb), 0.08);
    padding: 0.1rem 0.4rem;
    border-radius: var(--radius-sm);
    flex-shrink: 0;
  }

  .user-role.admin-role {
    color: var(--color-success);
    background: rgba(var(--color-success-rgb), 0.15);
  }

  .terminal-chevron {
    font-size: 0.7rem;
    transition: transform 0.3s ease;
    color: var(--color-content-muted);
    flex-shrink: 0;
  }

  .terminal-chevron.rotated {
    transform: rotate(180deg);
  }

  .terminal-dropdown {
    position: absolute;
    top: calc(100% + 0.5rem);
    right: -20px;
    width: 320px;
    background: linear-gradient(145deg, rgba(var(--color-background-700-rgb), 0.98), rgba(var(--color-background-900-rgb), 0.96));
    backdrop-filter: blur(20px) saturate(180%);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    box-shadow:
      0 20px 40px rgba(0, 0, 0, 0.6),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.15),
      inset 0 1px 0 rgba(var(--color-content-rgb), 0.1);
    overflow: hidden;
    z-index: var(--z-index-dropdown);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }

  .compact-header {
    padding: 1rem;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.08), rgba(var(--color-background-700-rgb), 0.8));
    border-bottom: 1px solid var(--color-border-default);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .header-avatar {
    flex-shrink: 0;
  }

  .user-text {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .user-name {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--color-white);
    line-height: 1.2;
  }

  .user-badge {
    font-size: 0.65rem;
    font-weight: 500;
    color: var(--color-content-muted);
    background: rgba(var(--color-content-rgb), 0.08);
    padding: 0.125rem 0.375rem;
    border-radius: var(--radius-sm);
    width: fit-content;
  }

  .user-badge.admin {
    color: var(--color-success);
    background: rgba(var(--color-success-rgb), 0.15);
  }

  .close-btn {
    width: 24px;
    height: 24px;
    background: rgba(var(--color-content-rgb), 0.05);
    border: 1px solid rgba(var(--color-content-rgb), 0.1);
    border-radius: var(--radius-sm);
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all 0.25s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
  }

  .close-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
  }

  .compact-content {
    padding: 1rem;
  }

  .compact-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    color: var(--color-text-secondary);
    text-decoration: none;
    border-radius: var(--radius-sm);
    transition: all 0.3s ease;
    background: rgba(var(--color-content-rgb), 0.02);
    border: 1px solid rgba(var(--color-content-rgb), 0.05);
    cursor: pointer;
    width: 100%;
    text-align: left;
    font-family: inherit;
    margin-bottom: 0.5rem;
  }

  .compact-item:hover {
    background: rgba(var(--color-background-600-rgb), 0.6);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-white);
    transform: translateY(-1px);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .compact-item .item-icon {
    font-size: 1rem;
    color: var(--color-brand-500);
    width: 18px;
    text-align: center;
  }

  .item-text {
    font-size: 0.8rem;
    font-weight: 500;
    flex: 1;
  }

  .arrow {
    font-size: 0.65rem;
    color: var(--color-text-quaternary);
    transition: all 0.25s ease;
  }

  .compact-item:hover .arrow {
    color: var(--color-brand-500);
    transform: translateX(2px);
  }

  .admin-item {
    background: rgba(var(--color-success-rgb), 0.05);
    border-color: rgba(var(--color-success-rgb), 0.1);
    margin-bottom: 0.5rem;
  }

  .admin-item:hover {
    background: rgba(var(--color-success-rgb), 0.1) !important;
    border-color: rgba(var(--color-success-rgb), 0.3) !important;
  }

  .admin-icon {
    color: var(--color-success) !important;
  }

  .admin-item:hover .arrow {
    color: var(--color-success) !important;
  }

  .function-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.5rem;
    margin-bottom: 0.75rem;
  }

  .grid-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.375rem;
    padding: 0.625rem;
    background: rgba(var(--color-content-rgb), 0.02);
    border: 1px solid rgba(var(--color-content-rgb), 0.05);
    border-radius: var(--radius-sm);
    text-decoration: none;
    color: var(--color-text-secondary);
    transition: all 0.25s ease;
    cursor: pointer;
    min-height: 60px;
    justify-content: center;
    font-size: 0.7rem;
    font-weight: 500;
  }

  .grid-item:hover {
    background: rgba(var(--color-background-600-rgb), 0.6);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-white);
    transform: translateY(-1px);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .grid-item .item-icon {
    font-size: 1rem;
    width: auto;
  }

  .grid-item .item-icon.dashboard {
    color: var(--color-cyan-500);
  }

  .grid-item .item-icon.messages {
    color: var(--color-info-500);
  }

  .grid-item .item-icon.images {
    color: var(--color-brand-500);
  }

  .grid-item .item-icon.folders {
    color: var(--color-success);
  }

  .grid-item .item-icon.shares {
    color: var(--color-warning-500);
  }

  .bottom-actions {
    border-top: 1px solid rgba(var(--color-content-rgb), 0.05);
    padding-top: 0.75rem;
    margin-top: 0.75rem;
  }

  .bottom-actions .compact-item {
    margin-bottom: 0.25rem;
  }

  .bottom-actions .compact-item:last-child {
    margin-bottom: 0;
  }

  .settings-item .item-icon {
    color: var(--color-content-muted);
  }

  .settings-item:hover .item-icon {
    color: var(--color-brand-500);
  }

  .logout-item {
    background: rgba(var(--color-error-rgb), 0.05);
    border-color: rgba(var(--color-error-rgb), 0.1);
  }

  .logout-item:hover {
    background: rgba(var(--color-error-rgb), 0.1) !important;
    border-color: rgba(var(--color-error-rgb), 0.3) !important;
  }

  .logout-item .item-icon {
    color: var(--color-danger) !important;
  }

  .logout-item .item-text {
    color: var(--color-danger);
  }

  .mobile-menu-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.3s ease;
    background: transparent;
    border: none;
  }

  .mobile-menu-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .hamburger {
    width: 18px;
    height: 14px;
    position: relative;
    cursor: pointer;
  }

  .hamburger span {
    display: block;
    width: 100%;
    height: 2px;
    background: var(--color-brand-500);
    border-radius: var(--radius-sm);
    transition: all 0.3s ease;
    position: absolute;
    box-shadow: 0 0 6px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .hamburger span:nth-child(1) {
    top: 0;
  }

  .hamburger span:nth-child(2) {
    top: 6px;
  }

  .hamburger span:nth-child(3) {
    top: 12px;
  }

  .hamburger.active span:nth-child(1) {
    transform: rotate(45deg);
    top: 6px;
    background: var(--color-error-500);
    box-shadow: 0 0 8px rgba(var(--color-error-rgb), 0.6);
  }

  .hamburger.active span:nth-child(2) {
    opacity: 0;
    transform: scale(0);
  }

  .hamburger.active span:nth-child(3) {
    transform: rotate(-45deg);
    top: 6px;
    background: var(--color-error-500);
    box-shadow: 0 0 8px rgba(var(--color-error-rgb), 0.6);
  }

  .desktop-only {
    display: flex;
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

  @media (max-width: 840px) {
    .terminal-display {
      min-width: 100px;
      max-width: 140px;
      padding: 0.4rem 0.5rem;
    }

    .user-name {
      max-width: 60px;
      font-size: 0.8rem;
    }

    .user-role {
      font-size: 0.65rem;
      padding: 0.1rem 0.3rem;
    }

    .navbar-actions {
      gap: 0.4rem;
      overflow: visible;
    }

    .user-terminal {
      gap: 0.5rem;
      flex-shrink: 0;
    }
  }

  @media (max-width: 900px) {
    .terminal-display {
      min-width: 120px;
      max-width: 160px;
      padding: 0.4rem 0.6rem;
    }

    .terminal-dropdown {
      right: -10px;
      width: 260px;
    }

    .function-grid {
      grid-template-columns: 1fr 1fr;
      gap: 0.375rem;
    }

    .grid-item {
      min-height: 55px;
      padding: 0.5rem;
      font-size: 0.65rem;
    }

    .navbar-actions {
      gap: 0.5rem;
    }

    .user-terminal {
      gap: 0.5rem;
    }

    .user-name {
      font-size: 0.8rem;
      max-width: 80px;
    }
  }

  @media (max-width: 480px) {
    .terminal-dropdown {
      right: 0;
      width: 240px;
    }

    .compact-header {
      padding: 0.625rem;
    }

    .user-name {
      font-size: 0.8rem;
    }

    .user-badge {
      font-size: 0.6rem;
      padding: 0.1rem 0.3rem;
    }

    .compact-content {
      padding: 0.625rem;
    }

    .function-grid {
      grid-template-columns: 1fr;
      gap: 0.375rem;
      margin-bottom: 0.625rem;
    }

    .grid-item {
      flex-direction: row;
      text-align: left;
      min-height: auto;
      gap: 0.5rem;
      justify-content: flex-start;
      padding: 0.5rem;
      font-size: 0.7rem;
    }

    .grid-item .item-icon {
      font-size: 0.875rem;
    }

    .compact-item {
      padding: 0.45rem 0.5rem;
      margin-bottom: 0.25rem;
    }

    .compact-item .item-icon {
      font-size: 0.8rem;
      width: 14px;
    }

    .item-text {
      font-size: 0.75rem;
    }

    .arrow {
      font-size: 0.6rem;
    }

    .bottom-actions {
      padding-top: 0.625rem;
      margin-top: 0.625rem;
    }
  }

  .terminal-menu-enter-active,
  .terminal-menu-leave-active {
    transition: all 0.3s ease;
  }

  .terminal-menu-enter-from,
  .terminal-menu-leave-to {
    opacity: 0;
    transform: translateY(-10px) scale(0.95);
  }
</style>
