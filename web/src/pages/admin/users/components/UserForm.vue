<script setup lang="ts">
  import { reactive, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { UserItem } from '@/api/admin/user'

  const { $t } = useTexts()

  const props = defineProps<{
    modelValue: UserItem | null
  }>()

  const emit = defineEmits<{
    (e: 'update:modelValue', value: UserItem | null): void
  }>()

  const formData = reactive<UserItem>({
    id: 0,
    username: '',
    email: '',
    avatar: '',
    status: 1,
    role: 3,
    created_at: '',
    updated_at: '',
  })

  /* 当props.modelValue变化时，更新表单数据 */
  watch(
    () => props.modelValue,
    (newVal) => {
      if (newVal) {
        Object.assign(formData, newVal)
      }
    },
    { immediate: true, deep: true }
  )

  watch(
    formData,
    (newVal) => {
      emit('update:modelValue', { ...newVal })
    },
    { deep: true }
  )
</script>

<template>
  <div v-if="modelValue" class="space-y-4">
    <div>
      <label class="mb-1 block text-sm text-content-muted">{{ $t('admin.users.form.username') }}</label>
      <cyberInput v-model="formData.username" :placeholder="$t('admin.users.form.usernamePlaceholder')" />
    </div>
    <div>
      <label class="mb-1 block text-sm text-content-muted">{{ $t('admin.users.form.email') }}</label>
      <cyberInput v-model="formData.email" :placeholder="$t('admin.users.form.emailPlaceholder')" disabled />
    </div>
    <div>
      <label class="mb-1 block text-sm text-content-muted">{{ $t('admin.users.form.status') }}</label>
      <div class="flex gap-4">
        <cyberRadio v-model="formData.status" :value="1">{{ $t('admin.users.form.statusActive') }}</cyberRadio>
        <cyberRadio v-model="formData.status" :value="2">{{ $t('admin.users.form.statusDisabled') }}</cyberRadio>
      </div>
    </div>
    <div>
      <label class="mb-1 block text-sm text-content-muted">{{ $t('admin.users.form.role') }}</label>
      <div class="flex gap-4">
        <cyberRadio v-model="formData.role" :value="3">{{ $t('admin.users.form.roleUser') }}</cyberRadio>
        <!-- <cyberRadio v-model="formData.role" :value="2">{{ $t('admin.users.form.roleAdmin') }}</cyberRadio> -->
        <cyberRadio v-model="formData.role" :value="1">{{ $t('admin.users.form.roleSuperAdmin') }}</cyberRadio>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .edit-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-group label {
    font-size: 0.9rem;
    color: rgba(255, 255, 255, 0.7);
  }

  .radio-group {
    display: flex;
    gap: 1rem;
  }
</style>
