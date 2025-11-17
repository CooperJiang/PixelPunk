<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useToast } from '@/components/Toast/useToast'
  import { useAuthApi } from '@/api/auth'
  import { useTexts } from '@/composables/useTexts'
  import Validator, { getValidationRules } from '@/utils/validation/validator'

  defineOptions({
    name: 'ResetPassword',
  })

  const route = useRoute()
  const router = useRouter()
  const toast = useToast()
  const authApi = useAuthApi()
  const { $t } = useTexts()
  const R = getValidationRules($t)

  const token = ref('')
  const newPassword = ref('')
  const confirmPassword = ref('')
  const loading = ref(false)
  const verifying = ref(true)
  const tokenValid = ref(false)
  const email = ref('')
  const errorMessage = ref('')

  const verifyToken = async () => {
    const tokenParam = route.query.token as string
    if (!tokenParam) {
      errorMessage.value = $t('resetPassword.error.missingToken')
      verifying.value = false
      return
    }

    token.value = tokenParam

    try {
      const response = await authApi.verifyResetToken(tokenParam)
      tokenValid.value = true
      email.value = response.email || ''
    } catch (error: any) {
      const message = error?.message || $t('resetPassword.error.invalidToken')
      errorMessage.value = message
    } finally {
      verifying.value = false
    }
  }

  const handleSubmit = async () => {
    const pwdCheck = Validator.validate(newPassword.value || '', [R.required, R.minLength(6), R.maxLength(50)], $t)
    if (!pwdCheck.valid) {
      toast.error(pwdCheck.message || $t('utils.validation.errors.password'))
      return
    }
    const confirmCheck = Validator.validate(confirmPassword.value || '', [R.required], $t)
    if (!confirmCheck.valid) {
      toast.error(confirmCheck.message || $t('utils.validation.errors.required'))
      return
    }
    if ((confirmPassword.value || '') !== (newPassword.value || '')) {
      toast.error($t('resetPassword.validation.passwordMismatch'))
      return
    }

    loading.value = true
    try {
      await authApi.resetPassword(token.value, newPassword.value)
      setTimeout(() => {
        router.push('/auth')
      }, 2000)
    } catch (_error: any) {
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    verifyToken()
  })
</script>

<template>
  <div class="reset-password-page">
    <CyberHomeBackground :show-enhanced-effects="true" :show-digital-rain="false" :show-red-orb="true" :show-aurora="true" />

    <div class="floating-particles">
      <div v-for="i in 8" :key="i" class="particle" :style="{ '--delay': i * 0.5 + 's', '--duration': 8 + i * 2 + 's' }" />
    </div>

    <div class="page-content">
      <div class="main-content">
        <div class="reset-panel-container">
          <div class="panel-glow-effect" />

          <div class="corner-lines">
            <div class="corner-line top-left" />
            <div class="corner-line top-right" />
            <div class="corner-line bottom-left" />
            <div class="corner-line bottom-right" />
          </div>

          <div class="scan-line-container">
            <div class="scan-line" />
          </div>

          <div v-if="verifying" class="reset-panel">
            <div class="panel-header">
              <h1 class="panel-title">{{ $t('resetPassword.verifying.title') }}</h1>
            </div>
            <div class="panel-body">
              <div class="loading-container">
                <i class="fas fa-spinner fa-spin loading-icon" />
                <p class="loading-text">{{ $t('resetPassword.verifying.loading') }}</p>
              </div>
            </div>
          </div>

          <div v-else-if="!tokenValid" class="reset-panel">
            <div class="panel-header">
              <h1 class="panel-title">{{ $t('resetPassword.title') }}</h1>
            </div>
            <div class="panel-body">
              <div class="error-container">
                <div class="error-icon">
                  <i class="fas fa-exclamation-triangle" />
                </div>
                <h3 class="error-title">{{ $t('resetPassword.error.invalidOrExpired') }}</h3>
                <p class="error-message">{{ errorMessage }}</p>
                <p class="error-hint">{{ $t('resetPassword.error.hint') }}</p>
                <button class="action-button" @click="router.push('/auth')">
                  <i class="fas fa-arrow-left mr-2" />
                  {{ $t('resetPassword.actions.backToLogin') }}
                </button>
              </div>
            </div>
          </div>

          <div v-else class="reset-panel">
            <div class="panel-header">
              <h1 class="panel-title">{{ $t('resetPassword.title') }}</h1>
              <p class="panel-subtitle">
                {{ $t('resetPassword.form.subtitle', { email }) }}
              </p>
            </div>

            <div class="panel-body">
              <form @submit.prevent="handleSubmit">
                <div class="form-group">
                  <label class="form-label">
                    <i class="fas fa-lock mr-2" />
                    {{ $t('resetPassword.form.newPassword.label') }}
                  </label>
                  <input
                    v-model="newPassword"
                    type="password"
                    class="form-input"
                    :placeholder="$t('resetPassword.form.newPassword.placeholder')"
                    :disabled="loading"
                    autocomplete="new-password"
                  />
                </div>

                <div class="form-group">
                  <label class="form-label">
                    <i class="fas fa-lock mr-2" />
                    {{ $t('resetPassword.form.confirmPassword.label') }}
                  </label>
                  <input
                    v-model="confirmPassword"
                    type="password"
                    class="form-input"
                    :placeholder="$t('resetPassword.form.confirmPassword.placeholder')"
                    :disabled="loading"
                    autocomplete="new-password"
                    @keyup.enter="handleSubmit"
                  />
                </div>

                <button type="submit" class="submit-button" :disabled="loading">
                  <i v-if="loading" class="fas fa-spinner fa-spin mr-2" />
                  <i v-else class="fas fa-check mr-2" />
                  {{ loading ? $t('resetPassword.form.submitting') : $t('resetPassword.form.submit') }}
                </button>

                <div class="back-to-login">
                  <button type="button" class="link-button" @click="router.push('/auth')">
                    <i class="fas fa-arrow-left mr-2" />
                    {{ $t('resetPassword.actions.backToLogin') }}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .reset-password-page {
    position: relative;
    min-height: 100vh;
    overflow: hidden;
    background: linear-gradient(135deg, var(--color-background-900) 0%, var(--color-background-800) 100%);
  }

  .floating-particles {
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 1;
  }

  .particle {
    position: absolute;
    width: 4px;
    height: 4px;
    background: var(--color-brand-500);
    border-radius: var(--radius-full);
    opacity: 0;
    animation: float var(--duration, 10s) infinite;
    animation-delay: var(--delay, 0s);
  }

  @keyframes float {
    0%,
    100% {
      opacity: 0;
      transform: translate(0, 0);
    }
    25% {
      opacity: 0.5;
    }
    50% {
      opacity: 0.8;
      transform: translate(100px, -100px);
    }
    75% {
      opacity: 0.5;
    }
  }

  .particle:nth-child(1) {
    top: 10%;
    left: 10%;
  }
  .particle:nth-child(2) {
    top: 20%;
    left: 80%;
  }
  .particle:nth-child(3) {
    top: 40%;
    left: 20%;
  }
  .particle:nth-child(4) {
    top: 60%;
    left: 70%;
  }
  .particle:nth-child(5) {
    top: 80%;
    left: 30%;
  }
  .particle:nth-child(6) {
    top: 30%;
    left: 60%;
  }
  .particle:nth-child(7) {
    top: 50%;
    left: 90%;
  }
  .particle:nth-child(8) {
    top: 70%;
    left: 15%;
  }

  .page-content {
    position: relative;
    z-index: 2;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
  }

  .main-content {
    width: 100%;
    max-width: 500px;
  }

  .reset-panel-container {
    position: relative;
    background: linear-gradient(135deg, rgba(var(--color-background-900-rgb), 0.98), rgba(var(--color-background-800-rgb), 0.95));
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow:
      0 20px 60px rgba(0, 0, 0, 0.5),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.1) inset,
      0 0 50px rgba(var(--color-brand-500-rgb), 0.15);
    overflow: hidden;
    animation: slideIn 0.5s ease-out;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(-20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .panel-glow-effect {
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.1) 0%, transparent 70%);
    animation: rotate 20s linear infinite;
    pointer-events: none;
  }

  @keyframes rotate {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .corner-lines {
    position: absolute;
    inset: 0;
    pointer-events: none;
  }

  .corner-line {
    position: absolute;
    width: 20px;
    height: 20px;
    border: 2px solid var(--color-brand-500);
  }

  .corner-line.top-left {
    top: 10px;
    left: 10px;
    border-right: none;
    border-bottom: none;
  }

  .corner-line.top-right {
    top: 10px;
    right: 10px;
    border-left: none;
    border-bottom: none;
  }

  .corner-line.bottom-left {
    bottom: 10px;
    left: 10px;
    border-right: none;
    border-top: none;
  }

  .corner-line.bottom-right {
    bottom: 10px;
    right: 10px;
    border-left: none;
    border-top: none;
  }

  .scan-line-container {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 100%;
    overflow: hidden;
    pointer-events: none;
  }

  .scan-line {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, transparent 0%, var(--color-brand-500) 50%, transparent 100%);
    opacity: 0.5;
    animation: scan 3s ease-in-out infinite;
  }

  @keyframes scan {
    0%,
    100% {
      top: 0;
      opacity: 0;
    }
    50% {
      opacity: 0.5;
    }
    100% {
      top: 100%;
    }
  }

  .reset-panel {
    position: relative;
    z-index: 1;
  }

  .panel-header {
    padding: 2rem 2rem 1rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.05), transparent);
  }

  .panel-title {
    font-size: 1.75rem;
    font-weight: 700;
    text-align: center;
    background: linear-gradient(90deg, var(--color-brand-500), var(--color-error-500));
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    margin-bottom: 0.5rem;
  }

  .panel-subtitle {
    text-align: center;
    color: var(--color-content-muted);
    font-size: 0.875rem;
  }

  .panel-subtitle strong {
    color: var(--color-brand-500);
    font-weight: 600;
  }

  .panel-body {
    padding: 2rem;
  }

  .loading-container,
  .error-container {
    text-align: center;
    padding: 2rem 0;
  }

  .loading-icon {
    font-size: 3rem;
    color: var(--color-brand-500);
    margin-bottom: 1rem;
  }

  .loading-text {
    color: var(--color-content-muted);
    font-size: 0.9375rem;
  }

  .error-icon {
    font-size: 4rem;
    color: var(--color-warning-500);
    margin-bottom: 1rem;
  }

  .error-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--color-content-heading);
    margin-bottom: 1rem;
  }

  .error-message {
    color: var(--color-content);
    font-size: 0.9375rem;
    margin-bottom: 0.75rem;
  }

  .error-hint {
    color: var(--color-content-muted);
    font-size: 0.875rem;
    margin-bottom: 2rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-label {
    display: block;
    margin-bottom: 0.5rem;
    color: var(--color-content);
    font-size: 0.875rem;
    font-weight: 500;
  }

  .form-input {
    width: 100%;
    padding: 0.875rem 1rem;
    background: rgba(var(--color-background-700-rgb), 0.5);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    color: var(--color-content);
    font-size: 0.875rem;
    transition: all 0.2s;
  }

  .form-input:focus {
    outline: none;
    border-color: var(--color-brand-500);
    background: rgba(var(--color-background-700-rgb), 0.8);
    box-shadow: 0 0 0 3px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .form-input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .submit-button,
  .action-button {
    width: 100%;
    padding: 0.875rem 1.5rem;
    background: linear-gradient(90deg, var(--color-brand-500), var(--color-error-500));
    border: none;
    border-radius: var(--radius-sm);
    color: var(--color-text-on-brand);
    font-size: 0.9375rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .submit-button:hover:not(:disabled),
  .action-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 20px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .submit-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }

  .back-to-login {
    margin-top: 1.5rem;
    text-align: center;
  }

  .link-button {
    background: transparent;
    border: none;
    color: var(--color-brand-500);
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .link-button:hover {
    color: var(--color-error-500);
    text-decoration: underline;
  }

  .mr-2 {
    margin-right: 0.5rem;
  }

  @media (max-width: 640px) {
    .page-content {
      padding: 1rem;
    }

    .reset-panel-container {
      border-radius: var(--radius-sm);
    }

    .panel-header {
      padding: 1.5rem 1.5rem 1rem;
    }

    .panel-title {
      font-size: 1.5rem;
    }

    .panel-body {
      padding: 1.5rem;
    }
  }
</style>
