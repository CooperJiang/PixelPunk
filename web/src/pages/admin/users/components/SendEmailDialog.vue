<script setup lang="ts">
  import { computed, reactive, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { useToast } from '@/components/Toast/useToast'
  import type { UserItem } from '@/api/admin/user'
  import Validator, { getValidationRules } from '@/utils/validation/validator'

  const props = defineProps<{
    modelValue: boolean
    user: UserItem | null
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    emailSent: [data: { userId: number; subject: string; content: string }]
  }>()

  const { $t } = useTexts()
  const toast = useToast()
  const R = getValidationRules($t)

  const visible = ref(false)
  const submitting = ref(false)

  /* 用户信息显示 */
  const userInfo = computed(() => {
    if (!props.user) {
      return ''
    }
    return `${props.user.username} (${props.user.email})`
  })

  const formData = reactive({
    subject: '',
    content: '',
  })

  const errors = reactive({
    subject: '',
    content: '',
  })

  const templates = {
    welcome: {
      subject: $t('admin.users.sendEmail.templatesContent.welcome.subject'),
      content: $t('admin.users.sendEmail.templatesContent.welcome.content'),
    },
    warning: {
      subject: $t('admin.users.sendEmail.templatesContent.warning.subject'),
      content: $t('admin.users.sendEmail.templatesContent.warning.content'),
    },
    notice: {
      subject: $t('admin.users.sendEmail.templatesContent.notice.subject'),
      content: $t('admin.users.sendEmail.templatesContent.notice.content'),
    },
    custom: {
      subject: '',
      content: $t('admin.users.sendEmail.templatesContent.custom.content'),
    },
  }

  watch(
    () => props.modelValue,
    (newVal) => {
      visible.value = newVal
      if (newVal) {
        resetForm()
      }
    },
    { immediate: true }
  )

  watch(visible, (newVal) => {
    emit('update:modelValue', newVal)
  })

  const resetForm = () => {
    formData.subject = ''
    formData.content = ''

    Object.keys(errors).forEach((key) => {
      errors[key] = ''
    })
  }

  const clearError = (field: string) => {
    errors[field] = ''
  }

  const useTemplate = (templateType: keyof typeof templates) => {
    const template = templates[templateType]
    formData.subject = template.subject
    formData.content = template.content
  }

  const validateForm = (): boolean => {
    let isValid = true

    const subjectCheck = Validator.validate(formData.subject, [R.required, R.maxLength(100)], $t)
    errors.subject = subjectCheck.valid ? '' : subjectCheck.message || $t('validation.errors.custom')
    if (!subjectCheck.valid) isValid = false

    const contentCheck = Validator.validate(formData.content, [R.required, R.maxLength(5000)], $t)
    errors.content = contentCheck.valid ? '' : contentCheck.message || $t('validation.errors.custom')
    if (!contentCheck.valid) isValid = false

    return isValid
  }

  const handleSubmit = async () => {
    if (!props.user) {
      return
    }

    if (!validateForm()) {
      return
    }

    submitting.value = true
    try {
      emit('emailSent', {
        userId: props.user.id,
        subject: formData.subject,
        content: formData.content,
      })

      toast.success($t('admin.users.sendEmail.messages.sent'))

      visible.value = false
    } catch {
      toast.error($t('admin.users.sendEmail.messages.failed'))
    }
    submitting.value = false
  }

  const handleCancel = () => {
    visible.value = false
  }
</script>

<template>
  <cyberDialog v-model="visible" :title="$t('admin.users.sendEmail.title')" width="500px">
    <div class="space-y-4 p-6">
      <div>
        <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.sendEmail.recipients') }}</label>
        <cyberInput :value="userInfo" readonly class="bg-background-200" />
      </div>

      <div>
        <label class="mb-2 block text-sm text-brand-700">
          <span class="text-red-500">*</span> {{ $t('admin.users.sendEmail.subject') }}
        </label>
        <cyberInput
          v-model="formData.subject"
          :placeholder="$t('admin.users.sendEmail.subjectPlaceholder')"
          :error="!!errors.subject"
          :error-message="errors.subject"
          @input="clearError('subject')"
        />
      </div>

      <div>
        <label class="mb-2 block text-sm text-brand-700">
          <span class="text-red-500">*</span> {{ $t('admin.users.sendEmail.content') }}
        </label>
        <cyberInput
          v-model="formData.content"
          type="textarea"
          :rows="6"
          :placeholder="$t('admin.users.sendEmail.contentPlaceholder')"
          :error="!!errors.content"
          :error-message="errors.content"
          @input="clearError('content')"
        />
      </div>

      <div>
        <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.sendEmail.templates') }}</label>
        <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
          <cyberButton type="secondary" @click="useTemplate('welcome')">{{
            $t('admin.users.sendEmail.templateButtons.welcome')
          }}</cyberButton>
          <cyberButton type="secondary" @click="useTemplate('warning')">{{
            $t('admin.users.sendEmail.templateButtons.warning')
          }}</cyberButton>
          <cyberButton type="secondary" @click="useTemplate('notice')">{{
            $t('admin.users.sendEmail.templateButtons.notice')
          }}</cyberButton>
          <cyberButton type="secondary" @click="useTemplate('custom')">{{
            $t('admin.users.sendEmail.templateButtons.custom')
          }}</cyberButton>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2 p-4">
        <cyberButton type="secondary" @click="handleCancel">{{ $t('admin.users.sendEmail.cancel') }}</cyberButton>
        <cyberButton type="primary" :loading="submitting" @click="handleSubmit">
          {{ submitting ? $t('admin.users.sendEmail.sending') : $t('admin.users.sendEmail.send') }}
        </cyberButton>
      </div>
    </template>
  </cyberDialog>
</template>

<style scoped lang="scss"></style>
