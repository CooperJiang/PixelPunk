<script setup lang="ts">
  import { computed, nextTick, ref, watch } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { folderApi } from '@/api'
  import type { CreateFolderRequest, FolderInfo, UpdateFolderRequest } from '@/api/types/index'
  import { useTexts } from '@/composables/useTexts'

  /* 定义Props和Emits */
  interface Props {
    modelValue: boolean
    parentId?: string
    mode?: 'create' | 'edit'
    folder?: FolderInfo
  }

  const props = defineProps<Props>()

  const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'created', folder: FolderInfo): void
    (e: 'updated', folder: FolderInfo): void
  }>()

  const isEditMode = computed(() => props.mode === 'edit')
  const { $t } = useTexts()

  /* 表单状态 */
  const form = ref<CreateFolderRequest | UpdateFolderRequest>({
    name: '',
    permission: 'public',
    description: '',
    parent_id: props.parentId,
  })

  /* 错误信息 */
  const errors = ref({
    name: '',
  })

  const confirmLoading = ref(false)

  /* 表单引用 */
  const nameInput = ref<HTMLInputElement | null>(null)

  /* 表单验证 - 统一使用 Validator */
  import Validator, { getValidationRules } from '@/utils/validation/validator'
  const R = getValidationRules($t)
  const validateName = () => {
    const check = Validator.validate(form.value?.name || '', [R.required, R.maxLength(30)], $t)
    errors.value.name = check.valid ? '' : check.message || $t('validation.errors.custom')
    return check.valid
  }

  const isFormValid = computed(() => form.value?.name?.trim() !== '' && !errors.value.name)

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  const toast = useToast()

  const resetForm = () => {
    form.value = {
      name: '',
      permission: 'public', // 确保每次都默认为公开
      description: '',
      parent_id: props.parentId,
    }
    errors.value = {
      name: '',
    }
  }

  const initEditForm = () => {
    if (props.folder) {
      form.value = {
        name: props.folder.name || '',
        parent_id: props.folder.parent_id,
        permission: props.folder.permission || 'public',
        description: props.folder.description || '',
      }
    }
  }

  const handleConfirm = async () => {
    if (!validateName() || confirmLoading.value) {
      return
    }

    confirmLoading.value = true

    try {
      if (isEditMode.value && props.folder?.id) {
        const result = await folderApi.updateFolder({
          folderId: props.folder.id,
          ...form.value,
        })

        if (result.success) {
          const updatedFolder = {
            ...result.data,
          }
          emit('updated', updatedFolder)
          resetForm()
          visible.value = false
          toast.success($t('folders.createFolderDialog.toast.updateSuccess'))
        }
      } else {
        const result = await folderApi.createFolder(form.value as CreateFolderRequest)

        if (result.success) {
          const newFolder = {
            ...result.data,
            ...form.value,
          }

          emit('created', newFolder)
          resetForm()
          visible.value = false
          toast.success($t('folders.createFolderDialog.toast.createSuccess'))
        }
      }
    } catch (error) {
      console.error($t('folders.createFolderDialog.toast.operationFailed'), error)
    }

    confirmLoading.value = false
  }

  const handleCancel = () => {
    resetForm()
    visible.value = false
  }

  watch(
    () => visible.value,
    (newVal) => {
      if (newVal) {
        nextTick(() => {
          nameInput.value?.focus()

          if (isEditMode.value) {
            initEditForm()
          } else {
            resetForm() // 新建模式时重置表单，确保默认权限为公开
          }
        })
      }
    }
  )

  watch(
    () => props.parentId,
    (newVal) => {
      if (!isEditMode.value && form.value) {
        form.value.parent_id = newVal
      }
    }
  )

  watch(
    () => props.folder,
    (newVal) => {
      if (isEditMode.value && newVal) {
        initEditForm()
      }
    }
  )
</script>

<template>
  <CyberDialog
    v-model="visible"
    :title="isEditMode ? $t('folders.createFolderDialog.title.edit') : $t('folders.createFolderDialog.title.create')"
    :loading="confirmLoading"
    width="450px"
    @confirm="handleConfirm"
    @cancel="handleCancel"
  >
    <form class="create-folder-form" @submit.prevent="handleConfirm">
      <div class="form-item">
        <label class="form-label"> <span class="required">*</span> {{ $t('folders.createFolderDialog.form.name.label') }} </label>
        <div class="form-control">
          <CyberInput
            ref="nameInput"
            v-model="form.name"
            :placeholder="$t('folders.createFolderDialog.form.name.placeholder')"
            :error="!!errors.name"
            :error-message="errors.name"
            clearable
            @input="validateName"
            @enter="handleConfirm"
          />
        </div>
      </div>

      <div class="form-item">
        <label class="form-label">{{ $t('folders.createFolderDialog.form.permission.label') }}</label>
        <div class="form-control">
          <div class="permission-options">
            <label class="permission-option" :class="{ active: form.permission === 'public' }">
              <input v-model="form.permission" type="radio" name="permission" value="public" />
              <div class="permission-icon">
                <i class="fas fa-globe" />
              </div>
              <div class="permission-details">
                <div class="permission-name">{{ $t('folders.createFolderDialog.form.permission.public.name') }}</div>
                <div class="permission-desc">{{ $t('folders.createFolderDialog.form.permission.public.desc') }}</div>
              </div>
            </label>

            <label class="permission-option" :class="{ active: form.permission === 'private' }">
              <input v-model="form.permission" type="radio" name="permission" value="private" />
              <div class="permission-icon">
                <i class="fas fa-lock" />
              </div>
              <div class="permission-details">
                <div class="permission-name">{{ $t('folders.createFolderDialog.form.permission.private.name') }}</div>
                <div class="permission-desc">{{ $t('folders.createFolderDialog.form.permission.private.desc') }}</div>
              </div>
            </label>
          </div>
        </div>
      </div>

      <div class="form-item">
        <label class="form-label">{{ $t('folders.createFolderDialog.form.description.label') }}</label>
        <div class="form-control">
          <CyberInput
            v-model="form.description"
            type="textarea"
            :placeholder="$t('folders.createFolderDialog.form.description.placeholder')"
            :rows="3"
          />
        </div>
      </div>
    </form>

    <template #footer>
      <div class="flex justify-end gap-3 p-4">
        <CyberButton type="outlined" size="medium" @click="handleCancel"> {{ $t('actions.cancel') }} </CyberButton>
        <CyberButton
          type="primary"
          size="medium"
          :icon="isEditMode ? 'edit' : 'folder-plus'"
          :loading="confirmLoading"
          :disabled="!isFormValid"
          @click="handleConfirm"
        >
          {{ isEditMode ? $t('folders.createFolderDialog.actions.save') : $t('folders.createFolderDialog.actions.create') }}
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped>
  .create-folder-form {
    margin-bottom: 1rem;
    padding: 0.5rem;
  }

  .form-item {
    margin-bottom: 1.5rem;
  }

  .form-label {
    display: block;
    margin-bottom: 0.75rem;
    font-weight: 600;
    color: var(--color-content-heading);
    font-size: 0.95rem;
    letter-spacing: 0.025em;
  }

  .required {
    color: var(--color-error-500);
    margin-right: 4px;
    font-weight: 700;
    text-shadow: 0 0 4px rgba(var(--color-error-rgb), 0.4);
  }

  .form-control {
    position: relative;
  }

  .cyber-input,
  .cyber-textarea {
    width: 100%;
    padding: 0.75rem;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: var(--color-input-text);
    font-family: inherit;
    transition: all 0.2s ease;
  }

  .cyber-input:focus,
  .cyber-textarea:focus {
    border-color: var(--color-brand-500);
    box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.2);
    outline: none;
  }

  .cyber-input.has-error {
    border-color: var(--color-error-500);
  }

  .error-message {
    color: var(--color-error-500);
    font-size: 0.875rem;
    margin-top: 0.25rem;
  }

  .permission-options {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .permission-option {
    flex: 1;
    display: flex;
    align-items: center;
    padding: 1rem;
    background: rgba(var(--color-background-700-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
    backdrop-filter: blur(8px);
  }

  .permission-option input {
    position: absolute;
    opacity: 0;
  }

  .permission-option:hover {
    background: rgba(var(--color-background-600-rgb), 0.8);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .permission-option.active {
    background: rgba(var(--color-brand-500-rgb), 0.12);
    border-color: var(--color-brand-500);
    box-shadow:
      0 0 25px rgba(var(--color-brand-500-rgb), 0.25),
      inset 0 1px 0 rgba(var(--color-brand-500-rgb), 0.3),
      0 4px 15px rgba(var(--color-brand-500-rgb), 0.15);
    transform: translateY(-1px);
  }

  .permission-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    margin-right: 0.75rem;
    transition: all 0.3s ease;
    font-size: 1rem;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    flex-shrink: 0;
  }

  .permission-option.active .permission-icon {
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-error-500));
    color: var(--color-content-heading);
    box-shadow:
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
    transform: scale(1.05);
    border-color: var(--color-brand-500);
  }

  .permission-details {
    flex: 1;
  }

  .permission-name {
    font-weight: 600;
    margin-bottom: 0.25rem;
    color: var(--color-content-heading);
    font-size: 0.95rem;
    letter-spacing: 0.015em;
  }

  .permission-desc {
    font-size: 0.8rem;
    color: var(--color-content-muted);
    line-height: 1.3;
    opacity: 0.85;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  @media (max-width: 640px) {
    .permission-options {
      flex-direction: column;
      gap: 0.75rem;
    }

    .permission-option {
      padding: 1rem;
    }

    .permission-icon {
      width: 32px;
      height: 32px;
      border-radius: var(--radius-sm);
      font-size: 0.9rem;
      margin-right: 0.5rem;
    }

    .form-label {
      font-size: 0.9rem;
      margin-bottom: 0.5rem;
    }
  }

  .permission-option:hover {
    animation: subtle-float 2s ease-in-out infinite;
  }

  @keyframes subtle-float {
    0%,
    100% {
      transform: translateY(-2px);
    }
    50% {
      transform: translateY(-3px);
    }
  }

  .permission-option.active .permission-icon {
    animation: icon-pulse 2s ease-in-out infinite;
  }

  @keyframes icon-pulse {
    0%,
    100% {
      box-shadow:
        0 4px 12px rgba(var(--color-brand-500-rgb), 0.4),
        0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
    }
    50% {
      box-shadow:
        0 6px 16px rgba(var(--color-brand-500-rgb), 0.6),
        0 0 30px rgba(var(--color-brand-500-rgb), 0.5);
    }
  }

  @keyframes gradient-flow {
    0% {
      background-position: 0% 50%;
    }
    50% {
      background-position: 100% 50%;
    }
    100% {
      background-position: 0% 50%;
    }
  }
</style>
