<script setup lang="ts">
  import { onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
  import { useAuthStore } from '@/store/auth'
  import { useToast } from '@/components/Toast/useToast'
  import { changeEmail, sendChangeEmailCode, updatePassword } from '@/api/user'
  import { useRouter } from 'vue-router'
  import { StorageUtil } from '@/utils/storage/storage'
  import { useTexts } from '@/composables/useTexts'

  const EMAIL_CODE_EXPIRE_KEY = 'email_code_expire_time'
  const EMAIL_SENDING_KEY = 'email_sending_target'

  const router = useRouter()
  const authStore = useAuthStore()
  const toast = useToast()
  const { $t } = useTexts()

  const currentEmail = ref('')

  const security = reactive({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  })
  const isChangingPassword = ref(false)

  const emailChange = reactive({
    newEmail: '',
    code: '',
  })
  const isChangingEmail = ref(false)
  const isSendingCode = ref(false)
  const countDown = ref(0)
  let timer: number | null = null

  watch(
    () => authStore.user,
    (newUserInfo) => {
      if (newUserInfo?.email) {
        currentEmail.value = newUserInfo.email
      }
    },
    { immediate: true }
  )

  onMounted(() => {
    if (authStore.user?.email) {
      currentEmail.value = authStore.user.email
    }

    checkExistingCodeTimer()
  })

  onBeforeUnmount(() => {
    clearCountdownTimer()
  })

  const checkExistingCodeTimer = () => {
    const expireTime = StorageUtil.get<number>(EMAIL_CODE_EXPIRE_KEY)
    const targetEmail = StorageUtil.get<string>(EMAIL_SENDING_KEY)

    if (!expireTime || !targetEmail) {
      return
    }

    const now = Date.now()

    if (expireTime > now) {
      const remainingSeconds = Math.ceil((expireTime - now) / 1000)

      if (remainingSeconds > 0) {
        countDown.value = remainingSeconds
        startCountdown()
      }
    } else {
      StorageUtil.remove(EMAIL_CODE_EXPIRE_KEY)
      StorageUtil.remove(EMAIL_SENDING_KEY)
    }
  }

  const startCountdown = () => {
    clearCountdownTimer()

    timer = window.setInterval(() => {
      countDown.value -= 1
      if (countDown.value <= 0) {
        clearCountdownTimer()
        StorageUtil.remove(EMAIL_CODE_EXPIRE_KEY)
        StorageUtil.remove(EMAIL_SENDING_KEY)
      }
    }, 1000)
  }

  const clearCountdownTimer = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }

  const changePassword = async () => {
    if (!security.currentPassword) {
      toast.error($t('settings.security.password.messages.currentRequired'))
      return
    }

    if (!security.newPassword) {
      toast.error($t('settings.security.password.messages.nextRequired'))
      return
    }

    if (security.newPassword !== security.confirmPassword) {
      toast.error($t('settings.security.password.messages.confirmMismatch'))
      return
    }

    if (security.newPassword.length < 6) {
      toast.error($t('settings.security.password.messages.minLength'))
      return
    }

    try {
      isChangingPassword.value = true

      const result = await updatePassword({
        oldPassword: security.currentPassword,
        newPassword: security.newPassword,
      })

      if (result.success) {
        toast.success($t('settings.security.password.messages.success'))

        security.currentPassword = ''
        security.newPassword = ''
        security.confirmPassword = ''

        authStore.logout()

        router.push('/auth')
      }
    } catch (_error) {
    } finally {
      isChangingPassword.value = false
    }
  }

  const sendVerificationCode = async () => {
    if (!emailChange.newEmail) {
      toast.error($t('settings.security.email.messages.nextRequired'))
      return
    }

    if (emailChange.newEmail === currentEmail.value) {
      toast.error($t('settings.security.email.messages.sameAsCurrent'))
      return
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(emailChange.newEmail)) {
      toast.error($t('settings.security.email.messages.invalidEmail'))
      return
    }

    try {
      isSendingCode.value = true

      const result = await sendChangeEmailCode({
        new_email: emailChange.newEmail,
      })

      if (result.success) {
        toast.success($t('settings.security.email.messages.codeSent'))

        countDown.value = 60

        const expireTime = Date.now() + 60 * 1000 // 当前时间 + 60秒
        StorageUtil.set(EMAIL_CODE_EXPIRE_KEY, expireTime)
        StorageUtil.set(EMAIL_SENDING_KEY, emailChange.newEmail)

        startCountdown()
      }
    } catch (_error) {
    } finally {
      isSendingCode.value = false
    }
  }

  const handleChangeEmail = async () => {
    if (!emailChange.newEmail) {
      toast.error($t('settings.security.email.messages.nextRequired'))
      return
    }

    if (!emailChange.code) {
      toast.error($t('settings.security.email.messages.codeRequired'))
      return
    }

    if (emailChange.newEmail === currentEmail.value) {
      toast.error($t('settings.security.email.messages.sameAsCurrent'))
      return
    }

    try {
      isChangingEmail.value = true

      const result = await changeEmail({
        new_email: emailChange.newEmail,
        code: emailChange.code,
      })

      if (result.success) {
        if (authStore.user) {
          authStore.updateUserInfo({ email: emailChange.newEmail })
          currentEmail.value = emailChange.newEmail
        }

        toast.success($t('settings.security.email.messages.success'))

        emailChange.newEmail = ''
        emailChange.code = ''

        clearCountdownTimer()
        StorageUtil.remove(EMAIL_CODE_EXPIRE_KEY)
        StorageUtil.remove(EMAIL_SENDING_KEY)
      }
    } catch (_error) {
    } finally {
      isChangingEmail.value = false
    }
  }
</script>

<template>
  <div class="grid gap-6 md:grid-cols-2">
    <div class="security-card border p-5">
      <h3 class="text-lg font-semibold text-content-heading">
        {{ $t('settings.security.password.title') }}
      </h3>

      <form class="mt-4 space-y-4" @submit.prevent="changePassword">
        <div class="space-y-2">
          <label class="block text-sm font-medium text-content">
            {{ $t('settings.security.password.labels.current') }}
          </label>
          <CyberInput
            v-model="security.currentPassword"
            type="password"
            :placeholder="$t('settings.security.password.placeholders.current')"
            prefix-icon="key"
            required
          />
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-content">
            {{ $t('settings.security.password.labels.next') }}
          </label>
          <CyberInput
            v-model="security.newPassword"
            type="password"
            :placeholder="$t('settings.security.password.placeholders.next')"
            prefix-icon="lock"
            required
            minlength="6"
          />
          <p class="text-xs text-content-muted">{{ $t('settings.security.password.hint') }}</p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-content">
            {{ $t('settings.security.password.labels.confirm') }}
          </label>
          <CyberInput
            v-model="security.confirmPassword"
            type="password"
            :placeholder="$t('settings.security.password.placeholders.confirm')"
            prefix-icon="check-circle"
            required
          />
        </div>

        <div class="flex justify-end">
          <CyberButton type="secondary" icon="key" :loading="isChangingPassword">
            {{
              isChangingPassword
                ? $t('settings.security.password.actions.submitting')
                : $t('settings.security.password.actions.submit')
            }}
          </CyberButton>
        </div>
      </form>
    </div>

    <div class="security-card border p-5">
      <h3 class="text-lg font-semibold text-content-heading">
        {{ $t('settings.security.email.title') }}
      </h3>

      <form class="mt-4 space-y-4" @submit.prevent="handleChangeEmail">
        <div class="space-y-2">
          <label class="block text-sm font-medium text-content">
            {{ $t('settings.security.email.labels.current') }}
          </label>
          <CyberInput
            v-model="currentEmail"
            disabled
            type="email"
            :placeholder="$t('settings.security.email.placeholders.current')"
            prefix-icon="envelope"
          />
          <p class="text-xs text-content-muted">{{ $t('settings.security.email.currentHint') }}</p>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-content">
            {{ $t('settings.security.email.labels.next') }}
          </label>
          <CyberInput
            v-model="emailChange.newEmail"
            type="email"
            :placeholder="$t('settings.security.email.placeholders.next')"
            prefix-icon="envelope"
            required
          />
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-content">
            {{ $t('settings.security.email.labels.code') }}
          </label>
          <div class="flex gap-2">
            <CyberInput
              v-model="emailChange.code"
              type="text"
              :placeholder="$t('settings.security.email.placeholders.code')"
              prefix-icon="shield-alt"
              required
            />
            <CyberButton
              type="outlined"
              :disabled="isSendingCode || countDown > 0"
              :loading="isSendingCode"
              @click.prevent="sendVerificationCode"
            >
              {{
                countDown > 0
                  ? $t('settings.security.email.actions.sendCodeCountdown', { seconds: countDown })
                  : $t('settings.security.email.actions.sendCode')
              }}
            </CyberButton>
          </div>
          <p class="text-xs text-content-muted">{{ $t('settings.security.email.codeHint') }}</p>
        </div>

        <div class="flex justify-end">
          <CyberButton type="secondary" icon="exchange-alt" :loading="isChangingEmail">
            {{
              isChangingEmail ? $t('settings.security.email.actions.submitting') : $t('settings.security.email.actions.submit')
            }}
          </CyberButton>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
  .security-card {
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.9));
    border-color: var(--color-border-subtle);
    box-shadow:
      0 4px 16px rgba(0, 0, 0, 0.08),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .security-card:hover {
    border-color: var(--color-border-default);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.03);
    transform: translateY(-1px);
  }
</style>
