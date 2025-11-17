<script setup lang="ts">
  import { ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useAuthApi } from '@/api/auth'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'ForgotPasswordModal',
  })

  const { $t } = useTexts()

  withDefaults(
    defineProps<{
      visible: boolean
    }>(),
    {
      visible: false,
    }
  )

  const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void
  }>()

  const toast = useToast()
  const authApi = useAuthApi()

  const email = ref('')
  const loading = ref(false)
  const sent = ref(false)

  const closeModal = () => {
    emit('update:visible', false)
    setTimeout(() => {
      email.value = ''
      sent.value = false
    }, 300)
  }

  const handleSubmit = async () => {
    if (!email.value) {
      toast.warning($t('auth.forgotPassword.errors.enterEmail'))
      return
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(email.value)) {
      toast.warning($t('auth.forgotPassword.errors.invalidEmail'))
      return
    }

    loading.value = true
    try {
      await authApi.forgotPassword(email.value)
      sent.value = true
    } catch (_error: any) {
    } finally {
      loading.value = false
    }
  }
</script>

<template>
  <CyberDialog
    :model-value="visible"
    :title="$t('auth.forgotPassword.title')"
    width="400px"
    :show-footer="false"
    @update:model-value="(val) => emit('update:visible', val)"
    @close="closeModal"
  >
    <div v-if="!sent" class="forgot-password-form">
      <p class="description">{{ $t('auth.forgotPassword.description') }}</p>

      <div class="form-group">
        <label class="form-label">
          <i class="fas fa-envelope" />
          {{ $t('auth.forgotPassword.emailLabel') }}
        </label>
        <CyberInput
          v-model="email"
          type="email"
          :placeholder="$t('auth.forgotPassword.emailPlaceholder')"
          :disabled="loading"
          @keyup.enter="handleSubmit"
        />
      </div>

      <CyberButton type="primary" block :loading="loading" @click="handleSubmit">
        <i v-if="!loading" class="fas fa-paper-plane mr-2" />
        {{ loading ? $t('auth.forgotPassword.sending') : $t('auth.forgotPassword.sendButton') }}
      </CyberButton>
    </div>

    <div v-else class="success-message">
      <div class="success-icon">
        <i class="fas fa-check-circle" />
      </div>
      <h3 class="success-title">{{ $t('auth.forgotPassword.success.title') }}</h3>
      <p class="success-text">
        {{ $t('auth.forgotPassword.success.message', { email }) }}
      </p>
      <p class="success-hint">{{ $t('auth.forgotPassword.success.hint') }}</p>
      <p class="success-note">{{ $t('auth.forgotPassword.success.note') }}</p>

      <CyberButton type="primary" block @click="closeModal">
        <i class="fas fa-check mr-2" />
        {{ $t('auth.forgotPassword.success.button') }}
      </CyberButton>
    </div>
  </CyberDialog>
</template>

<style scoped>
  .forgot-password-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .description {
    margin: 0 0 0.5rem 0;
    color: var(--color-content-muted);
    font-size: 0.8125rem;
    line-height: 1.5;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-label {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    color: var(--color-content);
    font-size: 0.8125rem;
    font-weight: 500;
  }

  .form-label i {
    font-size: 0.75rem;
    color: var(--color-brand-500);
  }

  .success-message {
    text-align: center;
    padding: 1rem 0;
  }

  .success-icon {
    margin-bottom: 0.75rem;
    font-size: 3rem;
    color: var(--color-success-500);
  }

  .success-title {
    margin: 0 0 0.75rem 0;
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-content-heading);
  }

  .success-text {
    margin: 0 0 0.5rem 0;
    color: var(--color-content);
    font-size: 0.8125rem;
    line-height: 1.5;
  }

  .success-text strong {
    color: var(--color-brand-500);
  }

  .success-hint {
    margin: 0 0 0.5rem 0;
    color: var(--color-content-muted);
    font-size: 0.75rem;
  }

  .success-note {
    margin: 0 0 1.5rem 0;
    padding: 0.5rem;
    background: rgba(var(--color-warning-rgb), 0.1);
    border-radius: var(--radius-sm);
    color: var(--color-warning-400);
    font-size: 0.6875rem;
    line-height: 1.4;
  }
</style>
