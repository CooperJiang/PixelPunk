<script setup lang="ts">
  import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'

  import { userApi } from '@/api'
  import { useAuthStore } from '@/store/auth'
  import { useToast } from '@/components/Toast/useToast'
  import { useMultipleLoading } from '@/hooks/useLoading'
  import { useTexts } from '@/composables/useTexts'

  import { StorageUtil } from '@/utils/storage/storage'
  import { VERIFY_CODE_COUNTDOWN_KEY, VERIFY_CODE_COUNTDOWN_SECONDS } from '@/constants/auth'
  import { EMAIL_REGEX } from '@/constants/validation'
  import type { RegisterFormState } from '../../types'

  defineOptions({
    name: 'AuthRegisterForm',
  })

  const props = defineProps<{
    emailVerificationEnabled?: boolean
  }>()

  const emit = defineEmits<{
    (e: 'register-success', username: string): void
    (e: 'password-focus', isFocused: boolean): void
  }>()

  const authStore = useAuthStore()
  const toast = useToast()
  const { $t } = useTexts()

  const { loading: loadingStates, withLoading } = useMultipleLoading({
    isLoading: false,
    isSendingCode: false,
  })

  const isLoading = computed(() => loadingStates.value.isLoading)
  const isSendingCode = computed(() => loadingStates.value.isSendingCode)

  const registerForm = reactive<RegisterFormState>({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    code: '',
  })

  const verifyCodeSent = ref(false)
  const countdown = ref(0)
  let countdownTimer: ReturnType<typeof setInterval> | null = null

  onMounted(() => {
    const storedCountdownTime = StorageUtil.get<number>(VERIFY_CODE_COUNTDOWN_KEY)
    if (!storedCountdownTime) {
      return
    }

    const now = Date.now()
    const remainingTime = Math.max(0, Math.floor((storedCountdownTime - now) / 1000))

    if (remainingTime > 0) {
      countdown.value = remainingTime
      startCountdown()
    } else {
      StorageUtil.remove(VERIFY_CODE_COUNTDOWN_KEY)
    }
  })

  onBeforeUnmount(() => {
    if (countdownTimer) {
      clearInterval(countdownTimer)
    }
  })

  const startCountdown = () => {
    if (countdownTimer) {
      clearInterval(countdownTimer)
    }

    countdownTimer = setInterval(() => {
      countdown.value -= 1

      if (countdown.value <= 0) {
        if (countdownTimer) {
          clearInterval(countdownTimer)
          countdownTimer = null
        }

        StorageUtil.remove(VERIFY_CODE_COUNTDOWN_KEY)
        countdown.value = 0
      }
    }, 1000)
  }

  const sendVerifyCode = async () => {
    if (!props.emailVerificationEnabled || countdown.value > 0) {
      return
    }

    if (!EMAIL_REGEX.test(registerForm.email)) {
      toast.error($t('auth.register.toastInvalidEmail'))
      return
    }

    await withLoading(
      'isSendingCode',
      async () => {
        const result = await userApi.sendRegistrationCode({
          email: registerForm.email,
        })

        if (!result.success) {
          throw new Error(result.message || $t('auth.register.toastSendCodeFailed'))
        }

        toast.success($t('auth.register.toastCodeSent'))
        verifyCodeSent.value = true
        countdown.value = VERIFY_CODE_COUNTDOWN_SECONDS

        const expiryTime = Date.now() + countdown.value * 1000
        StorageUtil.set<number>(VERIFY_CODE_COUNTDOWN_KEY, expiryTime, 1)
        startCountdown()
      },
      {
        onError: (error: Error) => {
          const message = error.message || $t('auth.register.toastSendCodeFailed')
          toast.error(message, 30000)
        },
      }
    )
  }

  const handleRegister = async () => {
    if (isLoading.value) {
      return
    }

    if (registerForm.password !== registerForm.confirmPassword) {
      toast.error($t('auth.register.toastPasswordMismatch'))
      return
    }

    if (!registerForm.code) {
      toast.error(
        props.emailVerificationEnabled ? $t('auth.register.toastCodeRequired') : $t('auth.register.toastDefaultCodeRequired')
      )
      return
    }

    if (!props.emailVerificationEnabled && registerForm.code !== '999999') {
      toast.error($t('auth.register.toastDefaultCodeInvalid'))
      return
    }

    await withLoading(
      'isLoading',
      async () => {
        await authStore.register({
          username: registerForm.username,
          email: registerForm.email,
          password: registerForm.password,
          code: registerForm.code,
        })

        try {
          await authStore.login({
            account: registerForm.username,
            password: registerForm.password,
          })

          toast.success($t('auth.register.toastSuccessAutoLogin'))
          window.location.href = '/'
          return
        } catch {
          toast.success($t('auth.register.toastSuccessFallback'))
          emit('register-success', registerForm.username)

          registerForm.username = ''
          registerForm.email = ''
          registerForm.password = ''
          registerForm.confirmPassword = ''
          registerForm.code = ''
          verifyCodeSent.value = false
          countdown.value = 0

          if (countdownTimer) {
            clearInterval(countdownTimer)
            countdownTimer = null
          }

          StorageUtil.remove(VERIFY_CODE_COUNTDOWN_KEY)
        }
      },
      {
        onError: (error: Error) => {
          const message = error.message || $t('auth.register.toastFailure')
          toast.error(message)
        },
      }
    )
  }
</script>

<template>
  <form class="form-small" @submit.prevent="handleRegister">
    <div class="input-container mb-2">
      <label class="cyber-label mb-1 block text-sm text-content" for="register-username">
        <span class="cyber-icon">⌦</span> {{ $t('auth.register.usernameLabel') }}
      </label>
      <div class="cyber-input-wrapper">
        <input
          id="register-username"
          v-model="registerForm.username"
          type="text"
          class="cyber-input w-full rounded border border-input-border bg-input-bg px-3 py-1.5 text-sm text-content outline-none focus:border-input-border-focus"
          :placeholder="$t('auth.register.usernamePlaceholder')"
          required
          :disabled="isLoading"
        />
        <span class="input-highlight" />
      </div>
    </div>

    <div class="mb-2">
      <div class="input-container mb-2">
        <label class="cyber-label mb-1 block text-sm text-content" for="register-email">
          <span class="cyber-icon">✉</span> {{ $t('auth.register.emailLabel') }}
        </label>
        <div class="flex items-stretch">
          <div class="cyber-input-wrapper flex-1">
            <input
              id="register-email"
              v-model="registerForm.email"
              type="email"
              class="cyber-input w-full rounded-l border border-input-border bg-input-bg px-3 py-1.5 text-sm text-content outline-none focus:border-input-border-focus"
              :placeholder="$t('auth.register.emailPlaceholder')"
              required
              :disabled="isLoading || verifyCodeSent"
            />
            <span class="input-highlight" />
          </div>
          <button
            v-if="props.emailVerificationEnabled"
            type="button"
            class="verify-code-btn min-w-16 flex-shrink-0 whitespace-nowrap rounded-r px-1 text-xs"
            :disabled="isSendingCode || countdown > 0"
            @click="sendVerifyCode"
          >
            <span v-if="countdown > 0"> {{ countdown }}{{ $t('auth.register.countdownUnit') }} </span>
            <span v-else-if="isSendingCode">{{ $t('auth.register.sendingCode') }}</span>
            <span v-else>{{ $t('auth.register.sendCode') }}</span>
          </button>
          <div v-else class="verify-code-btn disabled-verification-notice min-w-16 flex-shrink-0 whitespace-nowrap px-1 text-xs">
            {{ $t('auth.register.verificationDisabled') }}
          </div>
        </div>
      </div>

      <div class="input-container">
        <label class="cyber-label mb-1 block text-sm text-content" for="register-code">
          <span class="cyber-icon">#</span> {{ $t('auth.register.codeLabel') }}
        </label>
        <div class="cyber-input-wrapper">
          <input
            id="register-code"
            v-model="registerForm.code"
            type="text"
            class="cyber-input w-full rounded border border-input-border bg-input-bg px-3 py-1.5 text-sm text-content outline-none focus:border-input-border-focus"
            :placeholder="
              props.emailVerificationEnabled ? $t('auth.register.codePlaceholder') : $t('auth.register.codePlaceholderDisabled')
            "
            required
            :disabled="isLoading"
          />
          <span class="input-highlight" />
        </div>
        <div v-if="!props.emailVerificationEnabled" class="verification-notice mt-1">
          <span class="text-accent text-xs">
            <i class="fas fa-info-circle mr-1" />
            {{ $t('auth.register.verificationNotice') }}
          </span>
        </div>
      </div>
    </div>

    <div class="mb-3 grid grid-cols-2 gap-2">
      <div class="input-container">
        <label class="cyber-label mb-1 block text-sm text-content" for="register-password">
          <span class="cyber-icon">⌥</span> {{ $t('auth.register.passwordLabel') }}
        </label>
        <div class="cyber-input-wrapper">
          <input
            id="register-password"
            v-model="registerForm.password"
            type="password"
            class="cyber-input w-full rounded border border-input-border bg-input-bg px-3 py-1.5 text-sm text-content outline-none focus:border-input-border-focus"
            :placeholder="$t('auth.register.passwordPlaceholder')"
            required
            :disabled="isLoading"
            @focus="emit('password-focus', true)"
            @blur="emit('password-focus', false)"
          />
          <span class="input-highlight" />
        </div>
      </div>

      <div class="input-container">
        <label class="cyber-label mb-1 block text-sm text-content" for="register-confirm-password">
          <span class="cyber-icon">⍗</span> {{ $t('auth.register.confirmPasswordLabel') }}
        </label>
        <div class="cyber-input-wrapper">
          <input
            id="register-confirm-password"
            v-model="registerForm.confirmPassword"
            type="password"
            class="cyber-input w-full rounded border border-input-border bg-input-bg px-3 py-1.5 text-sm text-content outline-none focus:border-input-border-focus"
            :placeholder="$t('auth.register.confirmPasswordPlaceholder')"
            required
            :disabled="isLoading"
            @focus="emit('password-focus', true)"
            @blur="emit('password-focus', false)"
          />
          <span class="input-highlight" />
        </div>
      </div>
    </div>

    <button
      type="submit"
      class="cyber-btn w-full rounded bg-brand-500 py-2 text-sm text-background-900 hover:bg-error-500"
      :disabled="isLoading"
    >
      <span class="btn-glitch-effect" />
      {{ isLoading ? $t('auth.register.submitting') : $t('auth.register.submit') }}
    </button>
  </form>
</template>

<style scoped lang="scss" src="./RegisterForm.scss"></style>
