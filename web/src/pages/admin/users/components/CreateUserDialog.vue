<script setup lang="ts">
  import { reactive, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { useToast } from '@/components/Toast/useToast'
  import { createUser, type CreateUserParams, type UserItem } from '@/api/admin/user'
  import Validator, { getValidationRules } from '@/utils/validation/validator'

  const props = defineProps<{
    modelValue: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    userCreated: [user: UserItem]
  }>()

  const { $t } = useTexts()
  const toast = useToast()
  const R = getValidationRules($t)

  /* 存储单位选项 */
  const storageUnitOptions = [
    { label: 'MB', value: 'MB' },
    { label: 'GB', value: 'GB' },
    { label: 'TB', value: 'TB' },
  ]

  const visible = ref(false)
  const submitting = ref(false)

  /* 表单数据 */
  const formData = reactive<CreateUserParams & { storage_unit: string; bandwidth_unit: string }>({
    username: '',
    email: '',
    password: '',
    role: 3,
    storage_limit: 5,
    bandwidth_limit: 100,
    storage_unit: 'GB',
    bandwidth_unit: 'GB',
  })

  /* 错误状态 */
  const errors = reactive({
    username: null,
    email: null,
    password: null,
  })

  /* 监听props变化 */
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

  const convertToBytes = (value: number, unit: string): number => {
    switch (unit) {
      case 'TB':
        return value * 1024 * 1024 * 1024 * 1024
      case 'GB':
        return value * 1024 * 1024 * 1024
      case 'MB':
        return value * 1024 * 1024
      default:
        return value
    }
  }

  const resetForm = () => {
    formData.username = ''
    formData.email = ''
    formData.password = ''
    formData.role = 3
    formData.storage_limit = 5
    formData.bandwidth_limit = 100
    formData.storage_unit = 'GB'
    formData.bandwidth_unit = 'GB'

    Object.keys(errors).forEach((key) => {
      errors[key] = null
    })
  }

  const clearError = (field: string) => {
    errors[field] = null
  }

  const validateForm = (): boolean => {
    let isValid = true
    // 用户名：必填、长度 2-20、格式
    const usernameCheck = Validator.validate(formData.username, [R.required, R.minLength(2), R.maxLength(20), R.username], $t)
    errors.username = usernameCheck.valid ? null : usernameCheck.message || $t('validation.errors.custom')
    if (!usernameCheck.valid) isValid = false

    // 邮箱：必填、格式
    const emailCheck = Validator.validate(formData.email, [R.required, R.email], $t)
    errors.email = emailCheck.valid ? null : emailCheck.message || $t('validation.errors.email')
    if (!emailCheck.valid) isValid = false

    // 密码：必填、长度 6-20
    const passwordCheck = Validator.validate(formData.password, [R.required, R.minLength(6), R.maxLength(20)], $t)
    errors.password = passwordCheck.valid ? null : passwordCheck.message || $t('validation.errors.custom')
    if (!passwordCheck.valid) isValid = false

    return isValid
  }

  const handleSubmit = async () => {
    if (!validateForm()) {
      return
    }

    submitting.value = true
    try {
      const requestData: CreateUserParams = {
        username: formData.username,
        email: formData.email,
        password: formData.password,
        role: formData.role,
        storage_limit: formData.storage_limit ? convertToBytes(formData.storage_limit, formData.storage_unit) : undefined, // 转换为字节
        bandwidth_limit: formData.bandwidth_limit ? convertToBytes(formData.bandwidth_limit, formData.bandwidth_unit) : undefined, // 转换为字节
      }

      const result = await createUser(requestData)

      if (result.success) {
        toast.success($t('admin.users.toast.createSuccess'))

        emit('userCreated', result.data)

        visible.value = false
      }
    } catch {
      toast.error($t('admin.users.toast.createFailed'))
    }
    submitting.value = false
  }

  const handleCancel = () => {
    visible.value = false
  }
</script>

<template>
  <cyberDialog v-model="visible" :title="$t('admin.users.form.createTitle')" width="500px">
    <div class="space-y-4 p-6">
      <div>
        <label class="mb-2 block text-sm text-brand-700">
          <span class="text-red-500">*</span> {{ $t('admin.users.form.username') }}
        </label>
        <cyberInput
          v-model="formData.username"
          :placeholder="$t('admin.users.form.usernamePlaceholder')"
          :error="errors.username"
          @input="clearError('username')"
        />
        <div v-if="errors.username" class="mt-1 text-xs text-red-500">{{ errors.username }}</div>
      </div>

      <div>
        <label class="mb-2 block text-sm text-brand-700">
          <span class="text-red-500">*</span> {{ $t('admin.users.form.email') }}
        </label>
        <cyberInput
          v-model="formData.email"
          type="email"
          :placeholder="$t('admin.users.form.emailPlaceholder')"
          :error="errors.email"
          @input="clearError('email')"
        />
        <div v-if="errors.email" class="mt-1 text-xs text-red-500">{{ errors.email }}</div>
      </div>

      <div>
        <label class="mb-2 block text-sm text-brand-700">
          <span class="text-red-500">*</span> {{ $t('admin.users.form.password') }}
        </label>
        <cyberInput
          v-model="formData.password"
          type="password"
          :placeholder="$t('admin.users.form.passwordPlaceholder')"
          :error="errors.password"
          @input="clearError('password')"
        />
        <div v-if="errors.password" class="mt-1 text-xs text-red-500">{{ errors.password }}</div>
        <div class="mt-1 text-xs text-brand-500">{{ $t('admin.users.form.passwordHint') }}</div>
      </div>

      <div>
        <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.form.role') }}</label>
        <div class="flex gap-4">
          <cyberRadio v-model="formData.role" :value="3">{{ $t('admin.users.form.roleUser') }}</cyberRadio>
          <!-- <cyberRadio v-model="formData.role" :value="2">{{ $t('admin.users.form.roleAdmin') }}</cyberRadio> -->
          <cyberRadio v-model="formData.role" :value="1">{{ $t('admin.users.form.roleSuperAdmin') }}</cyberRadio>
        </div>
      </div>

      <div class="border-t border-subtle pt-4">
        <h4 class="mb-3 font-medium text-content">{{ $t('admin.users.form.storageConfig') }}</h4>
        <div class="grid grid-cols-1 gap-4">
          <div>
            <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.form.storageLimit') }}</label>
            <div class="flex items-center gap-2">
              <cyberInput
                v-model="formData.storage_limit"
                type="number"
                :min="0"
                :max="99999"
                :placeholder="$t('admin.users.form.storageLimitDefault')"
                class="flex-1"
              />
              <cyberDropdown v-model="formData.storage_unit" :options="storageUnitOptions" class="w-20" />
            </div>
          </div>
          <div>
            <label class="mb-2 block text-sm text-brand-700">{{ $t('admin.users.form.bandwidthLimit') }}</label>
            <div class="flex items-center gap-2">
              <cyberInput
                v-model="formData.bandwidth_limit"
                type="number"
                :min="0"
                :max="99999"
                :placeholder="$t('admin.users.form.bandwidthLimitDefault')"
                class="flex-1"
              />
              <cyberDropdown v-model="formData.bandwidth_unit" :options="storageUnitOptions" class="w-20" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2 p-4">
        <cyberButton type="secondary" @click="handleCancel">{{ $t('admin.users.buttons.cancel') }}</cyberButton>
        <cyberButton type="primary" :loading="submitting" @click="handleSubmit">{{
          $t('admin.users.buttons.create')
        }}</cyberButton>
      </div>
    </template>
  </cyberDialog>
</template>

<style scoped lang="scss"></style>
