<script setup lang="ts">
  import { onMounted, reactive, ref, computed } from 'vue'
  import { useRouter, useRoute } from 'vue-router'

  import { useAuthStore } from '@/store/auth'
  import { useToast } from '@/components/Toast/useToast'
  import { useLoading } from '@/hooks/useLoading'
  import { useTexts } from '@/composables/useTexts'

  import { StorageUtil } from '@/utils/storage/storage'
  import { REMEMBERED_LOGIN_KEY, REMEMBER_PASSWORD_HOURS } from '@/constants/auth'
  import type { LoginFormState, RememberedLogin } from '../../types'
  import ForgotPasswordModal from '@/components/ForgotPasswordModal/index.vue'
  import LegalDocumentModal from '@/components/LegalDocumentModal/index.vue'
  import {
    getLegalDocuments,
    getOAuthConfig,
    githubOAuthLogin,
    googleOAuthLogin,
    linuxdoOAuthLogin,
    type OAuthConfigResponse,
  } from '@/api/admin/settings'

  defineOptions({
    name: 'AuthLoginForm',
  })

  const emit = defineEmits<{
    (e: 'login-success'): void
    (e: 'switch-to-register'): void
    (e: 'password-focus', isFocused: boolean): void
  }>()

  const router = useRouter()
  const route = useRoute()
  const toast = useToast()
  const authStore = useAuthStore()
  const { $t } = useTexts()
  const { loading: isLoading, withLoading } = useLoading()

  const showForgotPasswordModal = ref(false)
  const showLegalModal = ref(false)
  const legalModalType = ref<'privacy' | 'terms'>('privacy')
  const privacyContent = ref('')
  const termsContent = ref('')
  const oauthConfig = ref<OAuthConfigResponse | null>(null)
  const isOAuthLoading = ref(false)

  interface OAuthState {
    provider: string
    random: string
    timestamp: number
  }

  const generateOAuthState = (provider: string): string => {
    const random = crypto.randomUUID()
    const stateData: OAuthState = {
      provider,
      random,
      timestamp: Date.now(),
    }
    sessionStorage.setItem('oauth_state', JSON.stringify(stateData))
    return btoa(JSON.stringify({ provider, random }))
  }

  const validateOAuthState = (stateParam: string): { valid: boolean; provider?: string } => {
    try {
      const savedStateStr = sessionStorage.getItem('oauth_state')
      if (!savedStateStr) {
        return { valid: false }
      }

      const savedState: OAuthState = JSON.parse(savedStateStr)
      const receivedState = JSON.parse(atob(stateParam))

      if (savedState.provider !== receivedState.provider || savedState.random !== receivedState.random) {
        return { valid: false }
      }

      const tenMinutes = 10 * 60 * 1000
      if (Date.now() - savedState.timestamp > tenMinutes) {
        sessionStorage.removeItem('oauth_state')
        return { valid: false }
      }

      sessionStorage.removeItem('oauth_state')
      return { valid: true, provider: savedState.provider }
    } catch (_error) {
      return { valid: false }
    }
  }

  const form = reactive<LoginFormState>({
    account: '',
    password: '',
    remember: false,
  })

  const isGithubEnabled = computed(() => {
    return oauthConfig.value?.github?.enabled ?? false
  })

  const isGoogleEnabled = computed(() => {
    return oauthConfig.value?.google?.enabled ?? false
  })

  const isLinuxdoEnabled = computed(() => {
    return oauthConfig.value?.linuxdo?.enabled ?? false
  })

  const enabledProvidersCount = computed(() => {
    let count = 0
    if (isGithubEnabled.value) count++
    if (isGoogleEnabled.value) count++
    if (isLinuxdoEnabled.value) count++
    return count
  })

  const hasMultipleProviders = computed(() => {
    return enabledProvidersCount.value > 1
  })

  const rememberedLogin = StorageUtil.get<RememberedLogin>(REMEMBERED_LOGIN_KEY)
  if (rememberedLogin) {
    form.account = rememberedLogin.account
    form.password = rememberedLogin.password
    form.remember = true
  }

  const handleLogin = async () => {
    if (isLoading.value) {
      return
    }

    await withLoading(
      async () => {
        await authStore.login({
          account: form.account,
          password: form.password,
        })

        if (form.remember) {
          StorageUtil.set<RememberedLogin>(
            REMEMBERED_LOGIN_KEY,
            {
              account: form.account,
              password: form.password,
            },
            REMEMBER_PASSWORD_HOURS
          )
        } else {
          StorageUtil.remove(REMEMBERED_LOGIN_KEY)
        }

        const fallbackName = $t('auth.login.successFallbackName')
        const username = authStore.userInfo?.username || fallbackName
        toast.success($t('auth.login.success').replace('{username}', username))

        if (!form.remember) {
          form.account = ''
          form.password = ''
          form.remember = false
        }

        emit('login-success')
        router.push('/')
      },
      {
        onError: (caughtError: unknown) => {
          const message = getErrorMessage(caughtError, $t('auth.login.failure'))
          toast.error(message)
        },
      }
    )
  }

  const forgotPassword = () => {
    showForgotPasswordModal.value = true
  }

  const switchToRegister = () => {
    emit('switch-to-register')
  }

  const openPrivacyPolicy = (e: Event) => {
    e.preventDefault()
    legalModalType.value = 'privacy'
    showLegalModal.value = true
  }

  const openTermsOfService = (e: Event) => {
    e.preventDefault()
    legalModalType.value = 'terms'
    showLegalModal.value = true
  }

  const loadLegalDocuments = async () => {
    try {
      const result = await getLegalDocuments()
      if (result.success && result.data) {
        privacyContent.value = result.data.privacy_policy || ''
        termsContent.value = result.data.terms_of_service || ''
      }
    } catch (_error) {}
  }

  const loadOAuthConfig = async () => {
    try {
      const result = await getOAuthConfig()
      if (result.success && result.data) {
        oauthConfig.value = result.data
      }
    } catch (_error) {}
  }

  type OAuthProvider = 'github' | 'google' | 'linuxdo'

  const oauthApiMap = {
    github: githubOAuthLogin,
    google: googleOAuthLogin,
    linuxdo: linuxdoOAuthLogin,
  }

  const handleOAuthCallback = async (provider: OAuthProvider, code: string) => {
    try {
      isOAuthLoading.value = true
      const result = await oauthApiMap[provider](code)
      if (result.success && result.data) {
        authStore.setToken(result.data.token)
        authStore.setUserInfo(result.data.userInfo)
        toast.success($t('auth.login.oauth.successWelcome').replace('{username}', result.data.userInfo.username))
        emit('login-success')
        router.replace('/')
      }
    } catch (error) {
      const message = getErrorMessage(error, $t('auth.login.oauth.failure'))
      toast.error(message)
      router.replace('/auth')
    } finally {
      isOAuthLoading.value = false
    }
  }

  const loginWithGithub = () => {
    if (!oauthConfig.value?.github) return

    const { client_id, redirect_uri, scope } = oauthConfig.value.github
    const state = generateOAuthState('github')
    const githubAuthUrl = `https://github.com/login/oauth/authorize?client_id=${client_id}&redirect_uri=${encodeURIComponent(redirect_uri)}&scope=${scope}&state=${state}`

    window.location.href = githubAuthUrl
  }

  const loginWithGoogle = () => {
    if (!oauthConfig.value?.google) return

    const { client_id, redirect_uri, scope } = oauthConfig.value.google
    const state = generateOAuthState('google')
    const googleAuthUrl = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${client_id}&redirect_uri=${encodeURIComponent(redirect_uri)}&response_type=code&scope=${encodeURIComponent(scope)}&access_type=offline&prompt=consent&state=${state}`

    window.location.href = googleAuthUrl
  }

  const loginWithLinuxdo = () => {
    if (!oauthConfig.value?.linuxdo) return

    const { client_id, redirect_uri, scope } = oauthConfig.value.linuxdo
    const state = generateOAuthState('linuxdo')
    const linuxdoAuthUrl = `https://connect.linux.do/oauth2/authorize?client_id=${client_id}&redirect_uri=${encodeURIComponent(redirect_uri)}&response_type=code&scope=${scope}&state=${state}`

    window.location.href = linuxdoAuthUrl
  }

  onMounted(async () => {
    loadLegalDocuments()
    loadOAuthConfig()

    const code = route.query.code as string
    const state = route.query.state as string

    if (code && state) {
      const validation = validateOAuthState(state)
      if (validation.valid && validation.provider) {
        await handleOAuthCallback(validation.provider as OAuthProvider, code)
      } else {
        toast.error($t('auth.login.oauth.stateValidationFailed'))
        router.replace('/auth')
      }
    }
  })

  function getErrorMessage(error: unknown, fallback: string): string {
    if (typeof error === 'string') {
      return error
    }

    if (error instanceof Error) {
      const axiosError = error as Error & {
        response?: {
          data?: {
            message?: string
          }
        }
      }
      return axiosError.response?.data?.message || axiosError.message || fallback
    }

    if (error && typeof error === 'object' && 'message' in error) {
      const message = (error as { message?: string }).message
      if (message) {
        return message
      }
    }

    return fallback
  }
</script>

<template>
  <div class="relative">
    <div v-if="isOAuthLoading" class="oauth-loading-overlay">
      <div class="loading-spinner">
        <div class="spinner-ring"></div>
        <p class="loading-text">{{ $t('auth.login.oauth.loading') }}</p>
      </div>
    </div>

    <form class="flex flex-col" @submit.prevent="handleLogin">
      <div class="input-container mb-2">
        <label class="cyber-label mb-1 block text-sm text-content" for="login-account">
          <span class="cyber-icon">⌘</span> {{ $t('auth.login.accountLabel') }}
        </label>
        <div class="cyber-input-wrapper">
          <input
            id="login-account"
            v-model="form.account"
            type="text"
            class="cyber-input w-full rounded border border-input-border bg-input-bg px-3 py-1.5 text-sm text-content outline-none focus:border-input-border-focus"
            :placeholder="$t('auth.login.accountPlaceholder')"
            required
            :disabled="isLoading"
          />
          <span class="input-highlight" />
        </div>
      </div>

      <div class="input-container mb-3">
        <label class="cyber-label mb-1 block text-sm text-content" for="login-password">
          <span class="cyber-icon">⌥</span> {{ $t('auth.login.passwordLabel') }}
        </label>
        <div class="cyber-input-wrapper">
          <input
            id="login-password"
            v-model="form.password"
            type="password"
            class="cyber-input w-full rounded border border-input-border bg-input-bg px-3 py-1.5 text-sm text-content outline-none focus:border-input-border-focus"
            :placeholder="$t('auth.login.passwordPlaceholder')"
            required
            :disabled="isLoading"
            @focus="emit('password-focus', true)"
            @blur="emit('password-focus', false)"
          />
          <span class="input-highlight" />
        </div>
        <div class="mt-3 flex justify-between">
          <CyberCheckbox v-model="form.remember" size="small">
            {{ $t('auth.login.rememberMe') }}
          </CyberCheckbox>
          <button type="button" class="cyber-link text-xs text-content hover:text-brand-400" @click="forgotPassword">
            {{ $t('auth.login.forgotPassword') }}
          </button>
        </div>
      </div>

      <div class="hint-text-area mb-3 mt-5 text-center text-xs">
        <p class="mb-2">
          {{ $t('auth.login.registerHint.prefix') }}
          <button type="button" class="register-link text-content" @click="switchToRegister">
            {{ $t('auth.login.registerHint.highlight') }}
          </button>
          {{ $t('auth.login.registerHint.suffix') }}
        </p>
        <p>
          {{ $t('auth.login.agreementPrefix') }}
          <a href="#" @click="openTermsOfService">{{ $t('auth.login.termsOfService') }}</a>
          {{ $t('auth.login.agreementAnd') }}
          <a href="#" @click="openPrivacyPolicy">{{ $t('auth.login.privacyPolicy') }}</a>
        </p>
      </div>

      <button
        type="submit"
        class="cyber-btn w-full rounded bg-brand-500 py-2 text-sm font-semibold transition-colors hover:bg-brand-400"
        :disabled="isLoading"
      >
        <span class="btn-glitch-effect" />
        {{ isLoading ? $t('auth.login.submitting') : $t('auth.login.submit') }}
      </button>

      <div v-if="isGithubEnabled || isGoogleEnabled || isLinuxdoEnabled" class="oauth-section">
        <div class="oauth-divider">
          <span>{{ $t('auth.login.oauth.divider') }}</span>
        </div>

        <div v-if="hasMultipleProviders" class="oauth-compact-layout">
          <CyberTooltip v-if="isGithubEnabled" :content="$t('auth.login.oauth.githubTooltip')" placement="top">
            <button type="button" class="oauth-icon-btn github" :disabled="isLoading" @click="loginWithGithub">
              <i class="fab fa-github" />
            </button>
          </CyberTooltip>

          <CyberTooltip v-if="isGoogleEnabled" :content="$t('auth.login.oauth.googleTooltip')" placement="top">
            <button type="button" class="oauth-icon-btn google" :disabled="isLoading" @click="loginWithGoogle">
              <i class="fab fa-google" />
            </button>
          </CyberTooltip>

          <CyberTooltip v-if="isLinuxdoEnabled" :content="$t('auth.login.oauth.linuxdoTooltip')" placement="top">
            <button type="button" class="oauth-icon-btn linuxdo" :disabled="isLoading" @click="loginWithLinuxdo">
              <i class="fab fa-linux" />
            </button>
          </CyberTooltip>
        </div>

        <div v-else>
          <button v-if="isGithubEnabled" type="button" class="github-login-btn" :disabled="isLoading" @click="loginWithGithub">
            <i class="fab fa-github" />
            <span>{{ $t('auth.login.oauth.githubButton') }}</span>
          </button>

          <button v-if="isGoogleEnabled" type="button" class="google-login-btn" :disabled="isLoading" @click="loginWithGoogle">
            <i class="fab fa-google" />
            <span>{{ $t('auth.login.oauth.googleButton') }}</span>
          </button>

          <button v-if="isLinuxdoEnabled" type="button" class="linuxdo-login-btn" :disabled="isLoading" @click="loginWithLinuxdo">
            <i class="fab fa-linux" />
            <span>{{ $t('auth.login.oauth.linuxdoButton') }}</span>
          </button>
        </div>
      </div>

      <ForgotPasswordModal v-model:visible="showForgotPasswordModal" />

      <LegalDocumentModal
        v-model="showLegalModal"
        :type="legalModalType"
        :content="legalModalType === 'privacy' ? privacyContent : termsContent"
      />
    </form>
  </div>
</template>

<style scoped lang="scss" src="./LoginForm.scss"></style>
