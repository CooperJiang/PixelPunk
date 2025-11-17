<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { RandomImageAPI, UpdateRandomAPIConfigRequest } from '@/api/openapi'
  import { useToast } from '@/components/Toast/useToast'
  import Validator, { getValidationRules } from '@/utils/validation/validator'

  defineOptions({
    name: 'EditApiModal',
  })

  const { $t } = useTexts()
  const toast = useToast()
  const R = getValidationRules($t)

  const props = defineProps<{
    visible: boolean
    isSubmitting: boolean
    api: RandomImageAPI | null
  }>()

  const emit = defineEmits<{
    'update:visible': [value: boolean]
    submit: [data: UpdateRandomAPIConfigRequest]
  }>()

  const formData = ref<UpdateRandomAPIConfigRequest>({
    folder_id: null,
    return_type: 'redirect',
  })

  const scopeType = ref<'all' | 'folder'>('all')

  const scopeOptions = computed(() => [
    { label: $t('openApi.card.scopeAllPublic'), value: 'all' },
    { label: $t('openApi.card.scopeFolder'), value: 'folder' },
  ])

  const isFormValid = computed(() => {
    if (scopeType.value === 'folder') {
      const folderCheck = Validator.validate(formData.value.folder_id, [R.required], $t)
      if (!folderCheck.valid) return false
    }
    return true
  })

  const handleFolderSelected = (folder: any) => {
    if (folder && folder.id) {
      formData.value.folder_id = folder.id
    }
  }

  watch(
    () => props.visible,
    (newVal) => {
      if (newVal && props.api) {
        formData.value = {
          folder_id: props.api.folder_id,
          return_type: props.api.return_type,
        }
        scopeType.value = props.api.folder_id ? 'folder' : 'all'
      }
    }
  )

  watch(scopeType, (newVal) => {
    if (newVal === 'all') {
      formData.value.folder_id = null
    }
  })

  const handleSubmit = () => {
    if (scopeType.value === 'folder') {
      const folderCheck = Validator.validate(formData.value.folder_id, [R.required], $t)
      if (!folderCheck.valid) {
        toast.error(folderCheck.message || $t('validation.errors.required'))
        return
      }
    }
    emit('submit', formData.value)
  }

  const handleDialogUpdate = (val: boolean) => {
    emit('update:visible', val)
  }
</script>

<template>
  <CyberDialog
    :model-value="visible"
    :title="$t('openApi.dialog.edit.title')"
    width="500px"
    :loading="isSubmitting"
    :show-default-footer="true"
    @confirm="handleSubmit"
    @update:model-value="handleDialogUpdate"
  >
    <div class="form-container">
      <div class="api-name-display">
        <span class="label">{{ $t('openApi.dialog.edit.form.nameLabel') }}</span>
        <span class="value">{{ api?.name }}</span>
      </div>

      <div class="form-group">
        <label class="form-label">{{ $t('openApi.dialog.edit.form.scope.label') }}</label>
        <CyberRadioGroup v-model="scopeType" :options="scopeOptions" :disabled="isSubmitting" layout="horizontal" />
      </div>

      <div v-if="scopeType === 'folder'" class="form-group">
        <label class="form-label">
          {{ $t('openApi.dialog.edit.form.folder.label') }}
          <span class="required">*</span>
        </label>
        <cyberFolderTree
          v-model="formData.folder_id"
          class="compact-dropdown text-sm"
          :disabled="isSubmitting"
          :dropdown-z-index="6000"
          @folder-selected="handleFolderSelected"
        />
        <p class="form-hint">{{ $t('openApi.dialog.edit.form.folder.hint') }}</p>
      </div>

      <div class="form-group">
        <label class="form-label">{{ $t('openApi.dialog.edit.form.returnType.label') }}</label>
        <div class="return-type-options">
          <label class="return-type-option">
            <CyberRadio v-model="formData.return_type" value="redirect" :disabled="isSubmitting" />
            <div class="option-content">
              <span class="option-label">{{ $t('openApi.dialog.edit.form.returnType.redirect.label') }}</span>
              <span class="option-desc">{{ $t('openApi.dialog.edit.form.returnType.redirect.desc') }}</span>
            </div>
          </label>
          <label class="return-type-option">
            <CyberRadio v-model="formData.return_type" value="direct" :disabled="isSubmitting" />
            <div class="option-content">
              <span class="option-label">{{ $t('openApi.dialog.edit.form.returnType.direct.label') }}</span>
              <span class="option-desc">{{ $t('openApi.dialog.edit.form.returnType.direct.desc') }}</span>
            </div>
          </label>
        </div>
      </div>
    </div>
  </CyberDialog>
</template>

<style scoped>
  .form-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .api-name-display {
    padding: 0.75rem 1rem;
    background: rgba(var(--color-background-800-rgb), 0.4);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .api-name-display .label {
    color: var(--color-content-muted);
    font-size: 0.875rem;
  }

  .api-name-display .value {
    color: var(--color-brand-500);
    font-size: 0.875rem;
    font-weight: 600;
    flex: 1;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-label {
    color: var(--color-text-content);
    font-size: 0.875rem;
    font-weight: 600;
  }

  .required {
    color: var(--color-error-500);
    margin-left: 0.25rem;
  }

  .form-hint {
    color: var(--color-content-muted);
    font-size: 0.75rem;
  }

  .return-type-options {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .return-type-option {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 0.75rem;
    background: rgba(var(--color-background-800-rgb), 0.4);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .return-type-option:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    background: rgba(var(--color-background-800-rgb), 0.6);
  }

  .option-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .option-label {
    color: var(--color-text-content);
    font-size: 0.875rem;
    font-weight: 600;
  }

  .option-desc {
    color: var(--color-content-muted);
    font-size: 0.75rem;
    line-height: 1.4;
  }
</style>
