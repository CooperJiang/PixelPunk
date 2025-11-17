<script setup lang="ts">
  import { computed, ref, watch, type PropType } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { TagInfo } from '@/api/types/index'
  import { formatDate as formatDateUtil } from '@/utils/formatting/format'

  const { $t } = useTexts()
  const props = defineProps({
    visible: {
      type: Boolean,
      default: false,
    },
    isEditing: {
      type: Boolean,
      default: false,
    },
    editingTag: {
      type: Object as PropType<TagInfo>,
      default: null,
    },
    isSubmitting: {
      type: Boolean,
      default: false,
    },
  })

  const emit = defineEmits(['close', 'submit'])

  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => {
      if (!value) {
        emit('close')
      }
    },
  })

  const formData = ref({
    name: '',
  })

  const nameError = ref('')
  import Validator, { getValidationRules } from '@/utils/validation/validator'
  const R = getValidationRules($t)

  const isFormValid = computed(() => formData.value.name.trim().length > 0 && !nameError.value)
  watch(
    () => props.editingTag,
    (newTag) => {
      if (newTag) {
        formData.value.name = newTag.name
      }
    },
    { immediate: true }
  )

  const resetForm = () => {
    formData.value.name = props.editingTag ? props.editingTag.name : ''
    nameError.value = ''
  }

  const onClose = () => {
    resetForm()
    emit('close')
  }

  const submitForm = () => {
    const check = Validator.validate(formData.value.name.trim(), [R.required, R.maxLength(50)], $t)
    if (!check.valid) {
      nameError.value = check.message || $t('admin.tags.form.nameEmpty')
      return
    }
    emit('submit', { ...props.editingTag, name: formData.value.name.trim() })
  }

  const formatDate = (dateStr: string) => formatDateUtil(dateStr, 'YYYY-MM-DD HH:mm')
</script>

<template>
  <CyberDialog
    v-model="dialogVisible"
    :title="isEditing ? $t('admin.tags.form.editTitle') : $t('admin.tags.form.createTitle')"
    width="450"
    @close="onClose"
  >
    <div class="tag-form">
      <div class="form-section">
        <div class="cyber-form-item">
          <label class="mb-1 block text-sm text-content-heading">{{ $t('admin.tags.form.nameLabel') }}</label>
          <CyberInput
            v-model="formData.name"
            :placeholder="$t('admin.tags.form.namePlaceholder')"
            :maxlength="50"
            show-character-count
            class="my-2 w-full"
            :status="nameError ? 'error' : ''"
            @update:model-value="nameError = ''"
          />
          <div v-if="nameError" class="text-accent mt-1 text-xs">{{ nameError }}</div>
          <p v-else class="mt-1 text-xs text-content-muted">{{ $t('admin.tags.form.nameDesc') }}</p>
        </div>
      </div>

      <div v-if="isEditing && editingTag" class="mt-4 rounded-md bg-background-900 p-3">
        <div class="flex items-center justify-between">
          <span class="text-sm text-content">{{ $t('admin.tags.form.usageCount') }}</span>
          <span class="font-medium text-content">{{ editingTag.count }}</span>
        </div>
        <div v-if="editingTag.created_at" class="mt-2 flex items-center justify-between">
          <span class="text-sm text-content">{{ $t('admin.tags.form.createdAt') }}</span>
          <span class="text-sm text-content">{{ formatDate(editingTag.created_at) }}</span>
        </div>
        <div v-if="editingTag.updated_at" class="mt-2 flex items-center justify-between">
          <span class="text-sm text-content">{{ $t('admin.tags.form.updatedAt') }}</span>
          <span class="text-sm text-content">{{ formatDate(editingTag.updated_at) }}</span>
        </div>

        <div class="mt-3 border-t border-default pt-2 text-xs text-content-muted">
          <p>{{ $t('admin.tags.form.editWarning') }}</p>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2 p-4">
        <CyberButton type="secondary" @click="onClose">{{ $t('admin.tags.form.cancel') }}</CyberButton>
        <CyberButton type="primary" :loading="isSubmitting" :disabled="!isFormValid || isSubmitting" @click="submitForm">
          {{ isEditing ? $t('admin.tags.form.update') : $t('admin.tags.form.create') }}
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped lang="scss">
  .tag-form {
    width: 100%;
  }

  .form-section {
    margin-bottom: var(--space-md);
  }

  .cyber-form-item + .cyber-form-item {
    margin-top: var(--space-md);
  }
</style>
