<script setup lang="ts">
  import { computed } from 'vue'
  import { useAuthStore } from '@/store/auth'
  import { useTexts } from '@/composables/useTexts'
  import { useToast } from '@/components/Toast/useToast'
  import { useRoute, useRouter } from 'vue-router'
  import { useConfig } from '@/composables/useConfig'
  import UserAvatar from '@/components/UserAvatar/index.vue'

  defineOptions({
    name: 'MobileMenu',
  })

  defineProps<{
    isOpen: boolean
  }>()

  const emit = defineEmits<{
    (e: 'close'): void
    (e: 'login'): void
  }>()

  const authStore = useAuthStore()
  const { $t } = useTexts()
  const toast = useToast()
  const router = useRouter()
  const route = useRoute()

  const { appInfo, brand } = useConfig()

  const userName = computed<string>(() => authStore.userInfo?.username || '')
  const isAdmin = computed<boolean>(() => authStore.isAdmin)

  const isRouteActive = (path: string): boolean => {
    if (path === '/') {
      return route.path === '/'
    }
    return route.path.startsWith(path)
  }

  const closeMobileMenu = (): void => {
    emit('close')
  }

  const handleLogin = (): void => {
    emit('login')
    closeMobileMenu()
  }

  const handleLogout = (): void => {
    authStore.logout()
    closeMobileMenu()
    toast.info($t('components.navbar.user.logoutSuccess'))
    router.push('/')
  }
</script>

<template>
  <Transition name="mobile-menu">
    <div v-show="isOpen" class="mobile-menu-overlay" @click="closeMobileMenu">
      <div class="mobile-menu" @click.stop>
        <div class="mobile-menu-header">
          <div class="mobile-brand">
            <div class="brand-icon">
              <img :src="brand.logo" :alt="appInfo.name" class="brand-img-logo" />
            </div>
            <div class="brand-text">
              <h2 class="brand-title">
                <span class="brand-pixel">{{ brand.pixel }}</span>
                <span class="brand-punk">{{ brand.punk }}</span>
              </h2>
            </div>
          </div>
          <button class="mobile-close-btn" @click="closeMobileMenu">
            <i class="fas fa-times" />
          </button>
        </div>

        <div class="mobile-menu-content">
          <div class="mobile-nav-section">
            <router-link to="/" class="mobile-nav-item" :class="{ active: isRouteActive('/') }" @click="closeMobileMenu">
              <i class="mobile-nav-icon fas fa-home" />
              <span>{{ $t('components.navbar.navigation.home') }}</span>
            </router-link>

            <router-link
              to="/upload"
              class="mobile-nav-item mobile-upload"
              :class="{ active: isRouteActive('/upload') }"
              @click="closeMobileMenu"
            >
              <i class="mobile-nav-icon fas fa-upload" />
              <span>{{ $t('components.navbar.navigation.upload') }}</span>
            </router-link>

            <router-link
              to="/explore"
              class="mobile-nav-item"
              :class="{ active: isRouteActive('/explore') }"
              @click="closeMobileMenu"
            >
              <i class="mobile-nav-icon fas fa-images" />
              <span>{{ $t('components.navbar.navigation.explore') }}</span>
            </router-link>

            <router-link
              to="/random"
              class="mobile-nav-item"
              :class="{ active: isRouteActive('/random') }"
              @click="closeMobileMenu"
            >
              <i class="mobile-nav-icon fas fa-random" />
              <span>{{ $t('components.navbar.navigation.random') }}</span>
            </router-link>

            <a href="/docs" target="_blank" class="mobile-nav-item" @click="closeMobileMenu">
              <i class="mobile-nav-icon fas fa-code" />
              <span>{{ $t('components.navbar.navigation.apiDocs') }}</span>
            </a>

            <router-link
              v-if="isAdmin"
              to="/admin"
              class="mobile-nav-item mobile-admin"
              :class="{ active: isRouteActive('/admin') }"
              @click="closeMobileMenu"
            >
              <i class="mobile-nav-icon fas fa-cog" />
              <span>{{ $t('components.navbar.navigation.admin') }}</span>
            </router-link>
          </div>

          <div v-if="authStore.isLoggedIn" class="mobile-user-section">
            <div class="mobile-user-info">
              <UserAvatar :avatar-url="authStore.userAvatar" :username="userName" size="lg" />
              <div class="mobile-user-details">
                <p class="mobile-user-name">{{ userName }}</p>
                <p class="mobile-user-role">
                  {{ isAdmin ? $t('components.navbar.user.roles.admin') : $t('components.navbar.user.roles.user') }}
                </p>
              </div>
            </div>

            <div class="mobile-user-actions">
              <router-link to="/dashboard" class="mobile-nav-item" @click="closeMobileMenu">
                <span>{{ $t('components.navbar.user.dashboard') }}</span>
              </router-link>

              <router-link to="/my-files" class="mobile-nav-item" @click="closeMobileMenu">
                <span>{{ $t('components.navbar.user.myWorks') }}</span>
              </router-link>

              <router-link to="/settings" class="mobile-nav-item" @click="closeMobileMenu">
                <span>{{ $t('components.navbar.user.settings') }}</span>
              </router-link>

              <button class="mobile-nav-item mobile-logout" @click="handleLogout">
                <span>{{ $t('components.navbar.user.logout') }}</span>
              </button>
            </div>
          </div>

          <div v-else class="mobile-login-section">
            <button class="mobile-login-btn" @click="handleLogin">
              <i class="fas fa-plug login-icon" />
              <span>{{ $t('components.navbar.user.login') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
  .mobile-menu-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.9);
    backdrop-filter: blur(20px);
    z-index: var(--z-index-modal);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
  }

  .mobile-menu {
    width: 100%;
    max-width: 340px;
    max-height: 80vh;
    background: rgba(0, 0, 0, 0.95);
    backdrop-filter: blur(25px);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    overflow: hidden;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.8);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }

  .mobile-menu-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem;
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .mobile-brand {
    display: flex;
    align-items: center;
    gap: 0.6rem;
  }

  .brand-icon {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .brand-img-logo {
    width: 32px;
    height: 32px;
    transition: all 0.3s ease;
    filter: drop-shadow(0 0 10px rgba(var(--color-brand-500-rgb), 0.6));
  }

  .brand-text {
    display: flex;
    flex-direction: column;
  }

  .brand-title {
    font-size: 1rem;
    font-weight: 800;
    line-height: 1;
    letter-spacing: 0.5px;
    margin: 0;
    font-family: 'Courier New', monospace;
  }

  .brand-pixel {
    color: var(--color-brand-500);
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.6);
  }

  .brand-punk {
    color: var(--color-error-500);
    text-shadow: 0 0 10px rgba(255, 121, 198, 0.6);
  }

  .mobile-close-btn {
    width: 32px;
    height: 32px;
    background: rgba(255, 121, 198, 0.15);
    border: none;
    border-radius: var(--radius-sm);
    color: var(--color-error-500);
    cursor: pointer;
    transition: all 0.25s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.9rem;
  }

  .mobile-close-btn:hover {
    background: rgba(255, 121, 198, 0.25);
    transform: scale(1.05);
  }

  .mobile-menu-content {
    padding: 0.75rem;
    overflow-y: auto;
    max-height: calc(80vh - 80px);
  }

  .mobile-nav-section {
    margin-bottom: 1.25rem;
  }

  .mobile-nav-item {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.75rem 1rem;
    color: var(--color-content-muted);
    text-decoration: none;
    border-radius: var(--radius-sm);
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    margin-bottom: 0.4rem;
    font-size: 0.85rem;
    font-weight: 500;
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid rgba(255, 255, 255, 0.03);
    cursor: pointer;
    width: 100%;
    text-align: left;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    position: relative;
    overflow: hidden;
  }

  .mobile-nav-icon {
    font-size: 0.9rem;
    width: 18px;
    text-align: center;
    transition: all 0.25s ease;
    opacity: 0.7;
    flex-shrink: 0;
  }

  .mobile-nav-item::after {
    content: '';
    position: absolute;
    inset: 0;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    opacity: 0;
    transition: opacity 0.25s ease;
    border-radius: var(--radius-sm);
  }

  .mobile-nav-item:hover {
    color: var(--color-white);
    border-color: rgba(var(--color-brand-500-rgb), 0.15);
    transform: translateX(2px);
  }

  .mobile-nav-item:hover::after {
    opacity: 1;
  }

  .mobile-nav-item:hover .mobile-nav-icon {
    opacity: 1;
    transform: scale(1.05);
    color: var(--color-brand-500);
  }

  .mobile-nav-item.active {
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .mobile-nav-item.active::after {
    opacity: 1;
    background: rgba(var(--color-brand-500-rgb), 0.08);
  }

  .mobile-nav-item.active .mobile-nav-icon {
    color: var(--color-brand-500);
    opacity: 1;
    filter: drop-shadow(0 0 4px rgba(var(--color-brand-500-rgb), 0.6));
  }

  .mobile-upload {
    color: var(--color-error-500) !important;
  }

  .mobile-upload::after {
    background: rgba(255, 121, 198, 0.05) !important;
  }

  .mobile-upload:hover .mobile-nav-icon {
    color: var(--color-error-500) !important;
  }

  .mobile-upload.active {
    color: var(--color-error-500) !important;
    border-color: rgba(255, 121, 198, 0.2) !important;
    text-shadow: 0 0 8px rgba(255, 121, 198, 0.4) !important;
  }

  .mobile-upload.active::after {
    background: rgba(255, 121, 198, 0.08) !important;
  }

  .mobile-upload .mobile-nav-icon {
    color: var(--color-error-500) !important;
  }

  .mobile-upload:hover .mobile-nav-icon,
  .mobile-upload.active .mobile-nav-icon {
    filter: drop-shadow(0 0 4px rgba(255, 121, 198, 0.6)) !important;
  }

  .mobile-admin {
    color: var(--color-success) !important;
  }

  .mobile-admin::after {
    background: rgba(126, 231, 135, 0.05) !important;
  }

  .mobile-admin:hover .mobile-nav-icon {
    color: var(--color-success) !important;
  }

  .mobile-admin.active {
    color: var(--color-success) !important;
    border-color: rgba(126, 231, 135, 0.2) !important;
    text-shadow: 0 0 8px rgba(126, 231, 135, 0.4) !important;
  }

  .mobile-admin.active::after {
    background: rgba(126, 231, 135, 0.08) !important;
  }

  .mobile-admin .mobile-nav-icon {
    color: var(--color-success) !important;
  }

  .mobile-admin:hover .mobile-nav-icon,
  .mobile-admin.active .mobile-nav-icon {
    filter: drop-shadow(0 0 4px rgba(126, 231, 135, 0.6)) !important;
  }

  .mobile-user-section {
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    padding-top: 0.75rem;
    margin-top: 0.75rem;
  }

  .mobile-user-info {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.75rem;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-radius: var(--radius-sm);
    margin-bottom: 0.75rem;
  }

  .mobile-user-details {
    flex: 1;
    min-width: 0;
  }

  .mobile-user-name {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin: 0 0 0.15rem 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .mobile-user-role {
    font-size: 0.7rem;
    color: var(--color-success);
    margin: 0;
    font-weight: 500;
    text-transform: uppercase;
  }

  .mobile-logout {
    color: var(--color-danger) !important;
  }

  .mobile-logout::after {
    background: rgba(248, 113, 113, 0.05) !important;
  }

  .mobile-logout:hover {
    color: var(--color-danger) !important;
    border-color: rgba(248, 113, 113, 0.2) !important;
    text-shadow: 0 0 8px rgba(248, 113, 113, 0.4) !important;
  }

  .mobile-logout:hover::after {
    background: rgba(248, 113, 113, 0.08) !important;
  }

  .mobile-login-section {
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    padding-top: 0.75rem;
    margin-top: 0.75rem;
  }

  .mobile-login-btn {
    width: 100%;
    background: transparent;
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 0.75rem 1rem;
    font-size: 0.85rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.25s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    letter-spacing: 0.025em;
    position: relative;
    overflow: hidden;
  }

  .mobile-login-btn:hover {
    color: var(--color-white);
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .login-icon {
    font-size: 0.8rem;
  }

  .mobile-menu-enter-active,
  .mobile-menu-leave-active {
    transition: all 0.3s ease;
  }

  .mobile-menu-enter-from,
  .mobile-menu-leave-to {
    opacity: 0;
    transform: scale(0.9);
  }

  .mobile-menu-enter-active .mobile-menu,
  .mobile-menu-leave-active .mobile-menu {
    transition: all 0.3s ease;
  }

  .mobile-menu-enter-from .mobile-menu,
  .mobile-menu-leave-to .mobile-menu {
    transform: translateY(-20px) scale(0.9);
  }
</style>
