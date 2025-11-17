<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { createMessageTemplate, updateMessageTemplate } from '@/api/message'
  import {
    getMessageTypeConfig,
    type CreateTemplateRequest,
    type MessageTemplate,
    type UpdateTemplateRequest,
  } from '@/api/message/types'
  import { useTexts } from '@/composables/useTexts'
  import Validator, { getValidationRules } from '@/utils/validation/validator'
  import { useToast } from '@/components/Toast/useToast'

  const { $t } = useTexts()
  const toast = useToast()
  const R = getValidationRules($t)

  const messageTypeConfig = getMessageTypeConfig($t)

  interface Props {
    show: boolean
    template?: MessageTemplate | null
  }

  const props = defineProps<Props>()

  const emit = defineEmits<{
    close: []
    saved: []
  }>()

  const loading = ref(false)
  const dialogVisible = ref(false)
  const formData = ref<CreateTemplateRequest>({
    type: '',
    title: '',
    content: '',
    description: '',
    is_enabled: true,
    send_email: false,
    show_toast: true,
    toast_type: 'info',
    default_action_type: '',
    default_action_text: '',
    default_action_style: 'primary',
    action_url_template: '',
  })

  const isEdit = computed(() => Boolean(props.template))

  const isFormValid = computed(() => {
    const typeOk = !!formData.value.type
    const titleOk = Validator.validate(formData.value.title, [R.required, R.maxLength(100)], $t).valid
    const contentOk = Validator.validate(formData.value.content, [R.required], $t).valid
    return typeOk && titleOk && contentOk
  })

  const initFormData = () => {
    if (props.template) {
      formData.value = {
        type: props.template.type,
        title: props.template.title,
        content: props.template.content,
        description: props.template.description || '',
        is_enabled: props.template.is_enabled,
        send_email: props.template.send_email,
        show_toast: props.template.show_toast,
        toast_type: props.template.toast_type,
        default_action_type: props.template.default_action_type || '',
        default_action_text: props.template.default_action_text || '',
        default_action_style: props.template.default_action_style,
        action_url_template: props.template.action_url_template || '',
      }
    } else {
      formData.value = {
        type: '',
        title: '',
        content: '',
        description: '',
        is_enabled: true,
        send_email: false,
        show_toast: true,
        toast_type: 'info',
        default_action_type: '',
        default_action_text: '',
        default_action_style: 'primary',
        action_url_template: '',
      }
    }
  }

  const handleClose = () => {
    if (!loading.value) {
      dialogVisible.value = false
      emit('close')
    }
  }

  const handleSubmit = async () => {
    if (loading.value) return
    const titleCheck = Validator.validate(formData.value.title, [R.required, R.maxLength(100)], $t)
    const contentCheck = Validator.validate(formData.value.content, [R.required], $t)
    if (!formData.value.type || !titleCheck.valid || !contentCheck.valid) {
      const msg =
        (!formData.value.type && $t('validation.errors.required')) ||
        titleCheck.message ||
        contentCheck.message ||
        $t('validation.errors.custom')
      toast?.error(msg)
      return
    }

    loading.value = true
    try {
      if (isEdit.value && props.template) {
        const updateData: UpdateTemplateRequest = {
          title: formData.value.title,
          content: formData.value.content,
          description: formData.value.description || '',
          is_enabled: formData.value.is_enabled ?? true,
          send_email: formData.value.send_email ?? false,
          show_toast: formData.value.show_toast ?? true,
          toast_type: formData.value.toast_type,
          default_action_type: formData.value.default_action_type || '',
          default_action_text: formData.value.default_action_text || '',
          default_action_style: formData.value.default_action_style,
          action_url_template: formData.value.action_url_template || '',
        }

        await updateMessageTemplate(props.template.id, updateData)
      } else {
        await createMessageTemplate(formData.value)
      }

      emit('saved')
      handleClose()
    } catch {
    } finally {
      loading.value = false
    }
  }

  watch(
    () => props.show,
    (newShow) => {
      dialogVisible.value = newShow
      if (newShow) {
        initFormData()
      }
    },
    { immediate: true }
  )

  watch(
    () => props.template,
    () => {
      if (props.show) {
        initFormData()
      }
    }
  )

  watch(dialogVisible, (newValue) => {
    if (!newValue) {
      emit('close')
    }
  })
</script>

<template>
  <CyberDialog
    v-model="dialogVisible"
    :title="isEdit ? $t('components.templateModal.editTitle') : $t('components.templateModal.createTitle')"
    width="800px"
    max-width="95vw"
    max-height="95vh"
    :show-default-footer="false"
    @close="handleClose"
  >
    <form class="template-form" @submit.prevent="handleSubmit">
      <div class="form-content space-y-4">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-2 block text-sm font-medium text-content-muted">
              {{ $t('components.templateModal.type') }} <span class="text-red-500">*</span>
            </label>
            <select v-model="formData.type" :disabled="isEdit" class="cyber-select w-full" required>
              <option value="">{{ $t('components.templateModal.selectType') }}</option>
              <option v-for="(config, type) in messageTypeConfig" :key="type" :value="type">
                {{ config.icon }} {{ config.label }}
              </option>
            </select>
            <p v-if="isEdit" class="mt-1 text-xs text-content-disabled">{{ $t('components.templateModal.cannotModifyType') }}</p>
          </div>

          <div>
            <label class="mb-2 block text-sm font-medium text-content-muted">{{ $t('components.templateModal.status') }}</label>
            <div class="flex items-center space-x-4">
              <cyberCheckbox v-model="formData.is_enabled">{{ $t('components.templateModal.enableTemplate') }}</cyberCheckbox>
            </div>
          </div>
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-content-muted">
            {{ $t('components.templateModal.title') }} <span class="text-red-500">*</span>
          </label>
          <input
            v-model="formData.title"
            type="text"
            class="cyber-input w-full"
            :placeholder="$t('components.templateModal.titlePlaceholder')"
            required
            maxlength="100"
          />
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-content-muted">
            {{ $t('components.templateModal.content') }} <span class="text-red-500">*</span>
          </label>
          <textarea
            v-model="formData.content"
            rows="4"
            class="cyber-textarea w-full"
            :placeholder="$t('components.templateModal.contentPlaceholder')"
            required
          />
          <p class="mt-1 text-xs text-content-disabled">
            {{ $t('components.templateModal.variablesHint') }}
          </p>
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-content-muted">{{
            $t('components.templateModal.description')
          }}</label>
          <input
            v-model="formData.description"
            type="text"
            class="cyber-input w-full"
            :placeholder="$t('components.templateModal.descriptionPlaceholder')"
            maxlength="200"
          />
        </div>

        <div class="cyber-panel p-4">
          <h4 class="mb-3 text-sm font-medium text-cyan-400">{{ $t('components.templateModal.notificationSettings') }}</h4>
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <cyberCheckbox v-model="formData.send_email">{{ $t('components.templateModal.sendEmail') }}</cyberCheckbox>
            </div>

            <div class="flex items-center justify-between">
              <cyberCheckbox v-model="formData.show_toast">{{ $t('components.templateModal.showToast') }}</cyberCheckbox>
              <select v-model="formData.toast_type" :disabled="!formData.show_toast" class="cyber-select ml-3">
                <option value="info">{{ $t('components.templateModal.toastTypes.info') }}</option>
                <option value="success">{{ $t('components.templateModal.toastTypes.success') }}</option>
                <option value="warning">{{ $t('components.templateModal.toastTypes.warning') }}</option>
                <option value="error">{{ $t('components.templateModal.toastTypes.error') }}</option>
              </select>
            </div>
          </div>
        </div>

        <div class="cyber-panel p-4">
          <h4 class="mb-3 text-sm font-medium text-cyan-400">{{ $t('components.templateModal.actionButtonSettings') }}</h4>
          <div class="space-y-3">
            <div>
              <label class="mb-2 block text-sm font-medium text-content-muted">{{
                $t('components.templateModal.actionType')
              }}</label>
              <input
                v-model="formData.default_action_type"
                type="text"
                class="cyber-input w-full"
                :placeholder="$t('components.templateModal.actionTypePlaceholder')"
              />
            </div>

            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-2 block text-sm font-medium text-content-muted">{{
                  $t('components.templateModal.buttonText')
                }}</label>
                <input
                  v-model="formData.default_action_text"
                  type="text"
                  class="cyber-input w-full"
                  :placeholder="$t('components.templateModal.buttonTextPlaceholder')"
                />
              </div>
              <div>
                <label class="mb-2 block text-sm font-medium text-content-muted">{{
                  $t('components.templateModal.buttonStyle')
                }}</label>
                <select v-model="formData.default_action_style" class="cyber-select w-full">
                  <option value="primary">{{ $t('components.templateModal.buttonStyles.primary') }}</option>
                  <option value="secondary">{{ $t('components.templateModal.buttonStyles.secondary') }}</option>
                  <option value="success">{{ $t('components.templateModal.buttonStyles.success') }}</option>
                  <option value="warning">{{ $t('components.templateModal.buttonStyles.warning') }}</option>
                  <option value="danger">{{ $t('components.templateModal.buttonStyles.danger') }}</option>
                  <option value="info">{{ $t('components.templateModal.buttonStyles.info') }}</option>
                </select>
              </div>
            </div>

            <div>
              <label class="mb-2 block text-sm font-medium text-content-muted">{{
                $t('components.templateModal.urlTemplate')
              }}</label>
              <input
                v-model="formData.action_url_template"
                type="text"
                class="cyber-input w-full"
                :placeholder="$t('components.templateModal.urlTemplatePlaceholder')"
              />
              <p class="mt-1 text-xs text-content-disabled">{{ $t('components.templateModal.urlVariablesHint') }}</p>
            </div>
          </div>
        </div>
      </div>
    </form>

    <template #footer>
      <div class="flex justify-end space-x-3 border-t border-cyan-800/30 p-4">
        <CyberButton type="outlined" :disabled="loading" @click="handleClose">
          {{ $t('components.templateModal.cancel') }}
        </CyberButton>
        <CyberButton type="primary" :disabled="loading || !isFormValid" :loading="loading" @click="handleSubmit">
          {{ isEdit ? $t('components.templateModal.updateButton') : $t('components.templateModal.createButton') }}
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped>
  .template-form {
    max-height: 70vh;
    overflow-y: auto;
  }

  .form-content {
    padding: 1rem;
  }

  .cyber-input {
    width: 100%;
    background: rgba(var(--color-background-900-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 0.5rem 0.75rem;
    color: var(--color-content-default);
    font-size: 0.875rem;
    transition: all 0.3s ease;
  }

  .cyber-input:focus {
    outline: none;
    border-color: var(--color-brand-500);
    box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-input::placeholder {
    color: var(--color-content-muted);
  }

  .cyber-input:disabled {
    background: rgba(var(--color-background-900-rgb), 0.1);
    border-color: rgba(75, 85, 99, 0.3);
    color: var(--color-content-muted);
    cursor: not-allowed;
  }

  .cyber-textarea {
    width: 100%;
    background: rgba(var(--color-background-900-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 0.5rem 0.75rem;
    color: var(--color-content-default);
    font-size: 0.875rem;
    resize: vertical;
    min-height: 80px;
    transition: all 0.3s ease;
  }

  .cyber-textarea:focus {
    outline: none;
    border-color: var(--color-brand-500);
    box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-textarea::placeholder {
    color: var(--color-content-muted);
  }

  .cyber-select {
    background: rgba(var(--color-background-900-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 0.5rem 0.75rem;
    color: var(--color-content-default);
    font-size: 0.875rem;
    transition: all 0.3s ease;
    cursor: pointer;
  }

  .cyber-select:focus {
    outline: none;
    border-color: var(--color-brand-500);
    box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-select:disabled {
    background: rgba(var(--color-background-900-rgb), 0.1);
    border-color: rgba(75, 85, 99, 0.3);
    color: var(--color-content-muted);
    cursor: not-allowed;
  }

  .cyber-select option {
    background: rgba(var(--color-background-900-rgb), 1);
    color: var(--color-content-default);
  }

  .cyber-checkbox {
    appearance: none;
    width: 1rem;
    height: 1rem;
    background: rgba(var(--color-background-900-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    position: relative;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .cyber-checkbox:checked {
    background: var(--color-brand-500);
    border-color: var(--color-brand-500);
  }

  .cyber-checkbox:checked::after {
    content: '✓';
    position: absolute;
    top: -1px;
    left: 1px;
    color: var(--color-text-on-brand);
    font-size: 0.75rem;
    font-weight: bold;
  }

  .cyber-checkbox:focus {
    outline: none;
    box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-panel {
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
  }

  .space-y-4 > * + * {
    margin-top: 1rem;
  }

  .space-y-3 > * + * {
    margin-top: 0.75rem;
  }

  .grid {
    display: grid;
  }

  .grid-cols-1 {
    grid-template-columns: repeat(1, minmax(0, 1fr));
  }

  .gap-4 {
    gap: 1rem;
  }

  @media (min-width: 640px) {
    .sm\:grid-cols-2 {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }
</style>
