<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { useTexts } from '@/composables/useTexts'
  import { useSettingsStore } from '@/store/settings'
  import { useLayoutStore } from '@/store/layout'
  import LoginForm from '../LoginForm/index.vue'
  import RegisterForm from '../RegisterForm/index.vue'

  defineOptions({
    name: 'AuthPanel',
  })

  const emit = defineEmits<{
    (e: 'password-focus', isFocused: boolean): void
  }>()

  const router = useRouter()
  const settingsStore = useSettingsStore()
  const layoutStore = useLayoutStore()
  const { $t } = useTexts()

  const activeTab = ref<'login' | 'register'>('login')
  const slideDirection = ref<'left' | 'right'>('right')

  const isRegistrationEnabled = computed(() => settingsStore.isRegistrationEnabled)
  const isEmailVerificationEnabled = computed(() => settingsStore.isEmailVerificationEnabled)

  const indicatorStyle = computed(() => ({
    transform: activeTab.value === 'login' ? 'translateX(0)' : 'translateX(100%)',
  }))

  watch(activeTab, (nextTab, previousTab) => {
    if (nextTab === 'register' && previousTab === 'login') {
      slideDirection.value = 'right'
    } else if (nextTab === 'login' && previousTab === 'register') {
      slideDirection.value = 'left'
    }
  })

  const setActiveTab = (tab: 'login' | 'register') => {
    if (tab === activeTab.value) {
      return
    }
    activeTab.value = tab
  }

  const onRegisterSuccess = () => {
    activeTab.value = 'login'
  }

  const onLoginSuccess = () => {
    return
  }

  const goToHome = () => {
    router.replace('/')
  }

  const showBackHomeButton = computed(() => layoutStore.isTopLayout)
</script>

<template>
  <div class="auth-panel cyber-panel">
    <div class="corner-decoration top-left" />
    <div class="corner-decoration top-right" />
    <div class="corner-decoration bottom-left" />
    <div class="corner-decoration bottom-right" />

    <div class="scan-line" />

    <div class="cyber-header">
      <div class="logo-wrapper mb-3">
        <div class="logo-glow" />
        <img src="/logo.png" :alt="$t('auth.panel.brandName')" class="h-10 w-10" />
        <h1 class="cyber-text ml-2 text-xl font-bold text-content">
          {{ $t('auth.panel.brandName') }}
        </h1>
      </div>

      <div v-if="!isRegistrationEnabled" class="registration-disabled-notice mb-4">
        <div class="notice-icon">
          <i class="fas fa-lock text-brand-400" />
        </div>
        <h2 class="cyber-text-glow mb-2 text-lg font-bold text-content">
          {{ $t('auth.panel.registrationDisabledTitle') }}
        </h2>
        <p class="text-sm text-content">
          {{ $t('auth.panel.registrationDisabledHint') }}
        </p>
      </div>

      <div v-else>
        <div class="tabs-wrapper mb-4">
          <div class="tabs-container">
            <button class="tab-button" :class="{ active: activeTab === 'login' }" @click="setActiveTab('login')">
              {{ $t('auth.panel.loginHeading') }}
            </button>
            <button class="tab-button" :class="{ active: activeTab === 'register' }" @click="setActiveTab('register')">
              {{ $t('auth.panel.registerHeading') }}
            </button>

            <div class="active-indicator" :style="indicatorStyle" />
          </div>
        </div>
      </div>
    </div>

    <div class="forms-container">
      <div v-if="!isRegistrationEnabled">
        <LoginForm class="form-transition" @login-success="onLoginSuccess" @password-focus="emit('password-focus', $event)" />
      </div>

      <div v-else>
        <transition :name="'slide-' + slideDirection" mode="out-in">
          <LoginForm
            v-if="activeTab === 'login'"
            key="login"
            class="form-transition"
            @login-success="onLoginSuccess"
            @switch-to-register="setActiveTab('register')"
            @password-focus="emit('password-focus', $event)"
          />

          <RegisterForm
            v-else
            key="register"
            :email-verification-enabled="isEmailVerificationEnabled"
            class="form-transition"
            @register-success="onRegisterSuccess"
            @password-focus="emit('password-focus', $event)"
          />
        </transition>
      </div>
    </div>

    <div v-if="showBackHomeButton" class="back-home-wrapper">
      <button type="button" class="back-home-button" @click="goToHome">
        <span class="button-icon">
          <i class="fas fa-arrow-left" />
        </span>
        <span class="button-text">{{ $t('auth.page.backHome') }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped lang="scss" src="./AuthPanel.scss"></style>
