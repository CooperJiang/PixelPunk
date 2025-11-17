<script setup lang="ts">
  import { onMounted, reactive, ref, watch } from 'vue'
  import { defaultSettings, testMailServer, type Setting } from '@/api/admin/settings'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  const toast = useToast()

  /* 测试邮件相关 */
  const testEmail = ref('')
  const isSending = ref(false)
  const testResult = ref<{ success: boolean; message: string } | null>(null)

  /* 从defaultSettings获取邮件设置的默认值 */
  const mailDefaults = defaultSettings.mail.reduce(
    (acc, setting) => {
      acc[setting.key] = setting.value
      return acc
    },
    {} as Record<string, unknown>
  )

  /* 本地设置对象（扁平化）- 使用导入的默认值 */
  const localSettings = reactive({
    smtp_host: mailDefaults.smtp_host || 'smtp.example.com',
    smtp_port: mailDefaults.smtp_port || 587,
    smtp_encryption: mailDefaults.smtp_encryption || 'tls',
    smtp_username: mailDefaults.smtp_username || '',
    smtp_password: mailDefaults.smtp_password || '',
    smtp_from_address: mailDefaults.smtp_from_address || '',
    smtp_from_name: mailDefaults.smtp_from_name || '',
  })

  /* 将扁平化的设置转换为Setting数组 - 使用defaultSettings中的类型和描述信息 */
  const getSettingsArray = (): Setting[] =>
    defaultSettings.mail.map((defaultSetting) => {
      const settingKey = defaultSetting.key as keyof typeof localSettings
      return {
        key: defaultSetting.key,
        value: localSettings[settingKey],
        type: defaultSetting.type,
        group: 'mail',
        description: defaultSetting.description,
      }
    })

  const applySettings = (settings: Setting[]) => {
    settings.forEach((setting) => {
      const key = setting.key as keyof typeof localSettings

      if (key in localSettings) {
        const settingsRecord = localSettings as Record<string, unknown>
        settingsRecord[key] = setting.value
      }
    })
  }

  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        applySettings(newSettings)
      }
    },
    { deep: true, immediate: true }
  )

  watch(
    localSettings,
    () => {
      emit('update', getSettingsArray())
    },
    { deep: true }
  )

  const sendTestEmail = async () => {
    if (!testEmail.value) {
      testResult.value = { success: false, message: $t('admin.settings.mail.enterTestEmail') }
      return
    }

    isSending.value = true
    testResult.value = null

    try {
      await testMailServer({
        email: testEmail.value,
        smtp_host: localSettings.smtp_host || '',
        smtp_port: localSettings.smtp_port || 587,
        smtp_encryption: localSettings.smtp_encryption || 'tls',
        smtp_username: localSettings.smtp_username || '',
        smtp_password: localSettings.smtp_password || '',
        smtp_from_address: localSettings.smtp_from_address || '',
        smtp_from_name: localSettings.smtp_from_name || '',
      })

      testResult.value = { success: true, message: $t('admin.settings.mail.sendSuccess') }
      toast.success($t('admin.settings.mail.testSuccess'))
    } catch (error: unknown) {
      testResult.value = {
        success: false,
        message: $t('admin.settings.mail.sendFailed', { error: error.message || $t('admin.settings.mail.unknownError') }),
      }
      toast.error(error.message)
    }

    isSending.value = false
  }

  onMounted(() => {
    if (props.settings && props.settings.length > 0) {
      applySettings(props.settings)
    }
  })
</script>

<template>
  <div class="space-y-6">
    <div class="space-y-5">
      <h3 class="mb-4 border-b border-subtle pb-2 text-lg font-medium text-content-heading">
        {{ $t('admin.settings.mail.smtpConfig') }}
      </h3>

      <div class="space-y-4">
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="w-full text-sm text-content-muted md:w-48 md:py-2">{{ $t('admin.settings.mail.smtpHost') }}</label>
          <div class="flex-1">
            <CyberInput v-model="localSettings.smtp_host" placeholder="smtp.example.com" width="500px" />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="w-full text-sm text-content-muted md:w-48 md:py-2">{{ $t('admin.settings.mail.smtpPort') }}</label>
          <div class="flex-1">
            <CyberInput v-model="localSettings.smtp_port" type="number" placeholder="587" width="500px" />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="w-full text-sm text-content-muted md:w-48 md:py-2">{{ $t('admin.settings.mail.encryption') }}</label>
          <div class="flex flex-1 gap-4">
            <CyberRadioGroup
              v-model="localSettings.smtp_encryption"
              :options="[
                { label: $t('admin.settings.mail.encryptionNone'), value: 'none' },
                { label: 'TLS', value: 'tls' },
                { label: 'SSL', value: 'ssl' },
              ]"
            />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="w-full text-sm text-content-muted md:w-48 md:py-2">{{ $t('admin.settings.mail.smtpUsername') }}</label>
          <div class="flex-1">
            <CyberInput v-model="localSettings.smtp_username" placeholder="username@example.com" width="500px" />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="w-full text-sm text-content-muted md:w-48 md:py-2">{{ $t('admin.settings.mail.smtpPassword') }}</label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.smtp_password"
              type="password"
              :placeholder="$t('admin.settings.mail.passwordPlaceholder')"
              width="500px"
            />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="w-full text-sm text-content-muted md:w-48 md:py-2">{{ $t('admin.settings.mail.fromAddress') }}</label>
          <div class="flex-1">
            <CyberInput v-model="localSettings.smtp_from_address" placeholder="noreply@example.com" width="500px" />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <label class="w-full text-sm text-content-muted md:w-48 md:py-2">{{ $t('admin.settings.mail.fromName') }}</label>
          <div class="flex-1">
            <CyberInput
              v-model="localSettings.smtp_from_name"
              :placeholder="$t('admin.settings.mail.fromNamePlaceholder')"
              width="500px"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="space-y-5">
      <h3 class="mb-4 border-b border-subtle pb-2 text-lg font-medium text-content-heading">
        {{ $t('admin.settings.mail.testEmail') }}
      </h3>

      <div class="space-y-4">
        <div class="flex flex-col md:flex-row md:items-center">
          <label class="w-full text-sm text-content-muted md:w-48 md:py-2">{{
            $t('admin.settings.mail.testEmailAddress')
          }}</label>
          <div class="flex-1">
            <CyberInput v-model="testEmail" placeholder="test@example.com" width="500px" />
          </div>
        </div>

        <div class="flex flex-col md:flex-row md:items-center">
          <div class="w-full md:w-48" />
          <div class="flex-1">
            <CyberButton type="primary" :loading="isSending" @click="sendTestEmail">
              {{ isSending ? $t('admin.settings.mail.sending') : $t('admin.settings.mail.sendTestEmail') }}
            </CyberButton>
            <span v-if="testResult" class="ml-3 text-sm" :class="testResult.success ? 'text-green-400' : 'text-red-400'">
              {{ testResult.message }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
