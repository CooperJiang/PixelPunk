<script setup lang="ts">
  import { computed } from 'vue'
  import { useAuthStore } from '@/store/auth'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'ActionButtons',
  })

  const authStore = useAuthStore()
  const { $t } = useTexts()

  const isLoggedIn = computed(() => authStore.isLoggedIn)
</script>

<template>
  <div class="action-buttons">
    <div class="desktop-buttons hidden flex-row items-center justify-center gap-4 md:flex">
      <router-link to="/explore" class="cyber-btn btn-blue">
        <i class="bi bi-search mr-2" />
        <span>{{ $t('home.actions.explore') }}</span>
      </router-link>
      <router-link to="/upload" class="cyber-btn btn-pink">
        <i class="bi bi-cloud-upload mr-2" />
        <span>{{ $t('home.actions.upload') }}</span>
      </router-link>
      <router-link :to="isLoggedIn ? '/dashboard' : '/auth'" class="cyber-btn btn-purple">
        <i class="bi bi-person mr-2" />
        <span>{{ isLoggedIn ? $t('home.actions.dashboard') : $t('home.actions.login') }}</span>
      </router-link>
    </div>

    <div class="mobile-buttons flex w-full flex-col items-center gap-3 md:hidden">
      <div class="mobile-row flex w-full max-w-sm flex-row gap-3">
        <router-link to="/explore" class="cyber-btn btn-blue flex-1">
          <i class="bi bi-search mr-1" />
          <span>{{ $t('home.actions.exploreShort') }}</span>
        </router-link>
        <router-link to="/upload" class="cyber-btn btn-pink flex-1">
          <i class="bi bi-cloud-upload mr-1" />
          <span>{{ $t('home.actions.uploadShort') }}</span>
        </router-link>
      </div>
      <div class="mobile-row w-full max-w-sm">
        <router-link :to="isLoggedIn ? '/dashboard' : '/auth'" class="cyber-btn btn-purple w-full">
          <i class="bi bi-person mr-2" />
          <span>{{ isLoggedIn ? $t('home.actions.dashboard') : $t('home.actions.login') }}</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .cyber-btn {
    @apply relative inline-flex min-w-32 items-center justify-center overflow-hidden rounded-md px-6 py-3 text-sm font-semibold uppercase tracking-wide no-underline;
    border: 2px solid;
    transition: all 0.15s ease;
    position: relative;
    z-index: 1;
    font-family: 'Inter', 'SF Pro Display', system-ui, sans-serif;
    letter-spacing: 0.05em;
    backdrop-filter: blur(10px);
    box-shadow: 0.15em 0.15em 0;
  }

  .cyber-btn:hover {
    transform: translate(-0.05em, -0.05em);
    filter: brightness(1.1);
  }

  .cyber-btn::after {
    content: '';
    position: absolute;
    inset: -2px;
    padding: 2px;
    background: linear-gradient(45deg, transparent, transparent);
    border-radius: inherit;
    mask:
      linear-gradient(var(--color-white) 0 0) content-box,
      linear-gradient(var(--color-white) 0 0);
    mask-composite: exclude;
    opacity: 0;
    transition: opacity 0.3s ease;
    z-index: -2;
  }

  .btn-blue {
    border: 2px solid var(--color-home-btn-blue);
    color: var(--color-home-btn-blue);
    background: rgba(var(--color-home-btn-blue-rgb), 0.1);
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-home-btn-blue-rgb), 0.6);
  }

  .btn-blue::after {
    background: linear-gradient(
      45deg,
      var(--color-home-btn-blue),
      var(--color-brand-600),
      var(--color-home-btn-blue),
      var(--color-brand-600)
    );
  }

  .btn-blue:hover {
    background: var(--color-home-btn-blue);
    color: var(--color-background-900);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-home-btn-blue-rgb), 0.7),
      0 0 25px rgba(var(--color-home-btn-blue-rgb), 0.5),
      inset 0 0 15px rgba(var(--color-content-rgb), 0.2);
    border-color: var(--color-home-btn-blue);
  }

  .btn-blue:hover::after {
    opacity: 1;
  }

  .btn-blue:active {
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-home-btn-blue-rgb), 0.5),
      0 0 15px rgba(var(--color-home-btn-blue-rgb), 0.2);
  }

  .btn-pink {
    border: 2px solid var(--color-home-btn-pink);
    color: var(--color-home-btn-pink);
    background: rgba(var(--color-home-btn-pink-rgb), 0.1);
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-home-btn-pink-rgb), 0.6);
  }

  .btn-pink::after {
    background: linear-gradient(
      45deg,
      var(--color-home-btn-pink),
      var(--color-error-600),
      var(--color-home-btn-pink),
      var(--color-error-600)
    );
  }

  .btn-pink:hover {
    background: var(--color-home-btn-pink);
    color: var(--color-background-900);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-home-btn-pink-rgb), 0.7),
      0 0 25px rgba(var(--color-home-btn-pink-rgb), 0.5),
      inset 0 0 15px rgba(var(--color-content-rgb), 0.2);
    border-color: var(--color-home-btn-pink);
  }

  .btn-pink:hover::after {
    opacity: 1;
  }

  .btn-pink:active {
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-home-btn-pink-rgb), 0.5),
      0 0 15px rgba(var(--color-home-btn-pink-rgb), 0.2);
  }

  .btn-purple {
    border: 2px solid var(--color-home-btn-purple);
    color: var(--color-home-btn-purple);
    background: rgba(var(--color-home-btn-purple-rgb), 0.08);
    box-shadow: 0.15em 0.15em 0 rgba(var(--color-home-btn-purple-rgb), 0.6);
  }

  .btn-purple::after {
    background: linear-gradient(
      45deg,
      var(--color-home-btn-purple),
      var(--color-func-accent),
      var(--color-home-btn-purple),
      var(--color-func-accent)
    );
  }

  .btn-purple:hover {
    background: var(--color-home-btn-purple);
    color: var(--color-background-900);
    transform: translate(-0.05em, -0.05em);
    box-shadow:
      0.2em 0.2em 0 rgba(var(--color-home-btn-purple-rgb), 0.7),
      0 0 25px rgba(var(--color-home-btn-purple-rgb), 0.5),
      inset 0 0 15px rgba(var(--color-content-rgb), 0.2);
    border-color: var(--color-home-btn-purple);
  }

  .btn-purple:hover::after {
    opacity: 1;
  }

  .btn-purple:active {
    transform: translate(0.05em, 0.05em);
    box-shadow:
      0.1em 0.1em 0 rgba(var(--color-home-btn-purple-rgb), 0.5),
      0 0 15px rgba(var(--color-home-btn-purple-rgb), 0.2);
  }

  .cyber-btn i {
    transition: all 0.3s ease;
  }

  .cyber-btn:hover i {
    transform: scale(1.1);
    filter: drop-shadow(0 0 5px currentColor);
  }

  .action-buttons {
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .mobile-buttons .cyber-btn {
    min-width: auto;
    padding: 12px 16px;
    font-size: 14px;
    font-weight: 600;
  }

  .mobile-buttons .mobile-row {
    padding: 0 1rem;
  }

  .mobile-buttons .mobile-row .cyber-btn {
    text-align: center;
    justify-content: center;
  }

  .mobile-buttons .cyber-btn i {
    font-size: 16px;
  }

  @media (max-width: 768px) {
    .mobile-buttons {
      margin-top: 0.5rem;
    }

    .mobile-buttons .mobile-row {
      padding: 0 0.5rem;
    }
  }

  @media (max-width: 480px) {
    .mobile-buttons .cyber-btn {
      padding: 10px 12px;
      font-size: 13px;
      min-height: 44px;
    }

    .mobile-buttons .mobile-row {
      gap: 8px;
    }

    .mobile-buttons .mobile-row .cyber-btn i {
      font-size: 14px;
    }
  }

  @media (max-width: 360px) {
    .mobile-buttons .cyber-btn {
      padding: 8px 10px;
      font-size: 12px;
    }

    .mobile-buttons .cyber-btn span {
      font-size: 12px;
    }

    .mobile-buttons .mobile-row {
      gap: 6px;
    }
  }
</style>
