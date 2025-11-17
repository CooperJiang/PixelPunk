<script setup lang="ts">
  import { computed } from 'vue'
  import { useAuthStore } from '@/store/auth'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'SidebarUser',
  })

  interface Props {
    collapsed?: boolean
  }

  withDefaults(defineProps<Props>(), {
    collapsed: false,
  })

  const emit = defineEmits<{
    (e: 'login'): void
  }>()

  const { $t } = useTexts()
  const authStore = useAuthStore()

  const isLoggedIn = computed(() => authStore.isLoggedIn)
  const userInfo = computed(() => authStore.userInfo)

  const handleLogin = () => {
    emit('login')
  }
</script>

<template>
  <div class="sidebar-user">
    <div v-if="isLoggedIn" class="user-logged-in">
      <div class="simple-user-area">
        <div class="user-display" :class="{ 'display-collapsed': collapsed }">
          <div class="user-avatar">
            <img
              v-if="authStore.userAvatar"
              :src="authStore.userAvatar"
              :alt="userInfo?.username || 'User'"
              class="avatar-image"
            />
            <div v-else class="avatar-fallback">
              <i class="fas fa-user" />
            </div>
            <div class="online-dot" />
          </div>

          <div v-if="!collapsed" class="user-info">
            <div class="username">{{ userInfo?.username || 'User' }}</div>
            <div class="user-status">
              <span class="status-text">{{ $t('components.sidebar.user.online') }}</span>
              <span v-if="authStore.isAdmin" class="admin-badge">{{ $t('components.sidebar.user.admin') }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="user-not-logged">
      <button class="login-button" :class="{ 'button-collapsed': collapsed }" @click="handleLogin">
        <i class="fas fa-sign-in-alt" />
        <span v-if="!collapsed">{{ $t('components.sidebar.user.login') }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
  .sidebar-user {
    @apply relative;
  }

  .simple-user-area {
    @apply relative p-3 transition-colors;
    background: rgba(var(--color-background-900-rgb), 0.35);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-sm);
  }

  .simple-user-area:hover {
    background: rgba(var(--color-brand-500-rgb), 0.08);
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .user-display {
    @apply flex items-center;
  }

  .display-collapsed {
    @apply justify-center;
  }

  .user-avatar {
    @apply relative h-10 w-10 flex-shrink-0 overflow-hidden rounded-full;
    background: var(--color-background-600);
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .avatar-image {
    @apply h-full w-full object-cover;
  }

  .avatar-fallback {
    @apply flex h-full w-full items-center justify-center;
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    font-size: 14px;
  }

  .online-dot {
    @apply absolute -bottom-0.5 -right-0.5 h-3 w-3 rounded-full;
    background: var(--color-success-500);
    border: 2px solid var(--color-background-400);
    box-shadow: 0 0 6px rgba(var(--color-success-rgb), 0.6);
  }

  .user-info {
    @apply ml-3 flex-1 overflow-hidden;
  }

  .username {
    @apply truncate text-sm font-semibold;
    color: var(--color-content-default);
    margin-bottom: 2px;
  }

  .user-status {
    @apply flex items-center gap-2;
  }

  .status-text {
    @apply text-xs;
    color: var(--color-success-500);
    font-weight: 500;
  }

  .admin-badge {
    @apply rounded-full px-2 py-0.5 text-xs;
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    font-weight: 500;
  }

  .login-button {
    @apply flex w-full items-center justify-center px-4 py-3 font-medium transition-all;
    background: var(--color-brand-500);
    color: var(--color-text-on-brand);
    border: 1px solid var(--color-brand-500);
    border-radius: var(--radius-sm);
  }

  .login-button:hover {
    background: var(--color-brand-400);
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.35);
    transform: translateY(-1px);
  }

  .button-collapsed {
    @apply p-2;
  }

  .login-button i {
    @apply mr-2;
  }

  .button-collapsed i {
    @apply mr-0;
  }
</style>
