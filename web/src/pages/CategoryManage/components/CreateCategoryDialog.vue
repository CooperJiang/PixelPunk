<script setup lang="ts">
  import { ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { createCategory } from '@/api/user/category'
  import type { ImageCategory } from '@/api/types/category'
  import { useTexts } from '@/composables/useTexts'
  import Validator, { getValidationRules } from '@/utils/validation/validator'

  defineOptions({
    name: 'CreateCategoryDialog',
  })

  const { $t } = useTexts()
  const visible = defineModel<boolean>({ required: true })
  const R = getValidationRules($t)

  const emit = defineEmits<{
    categoryCreated: [category: ImageCategory]
  }>()

  const toast = useToast()

  const newCategory = ref({
    name: '',
    description: '',
    sort_order: 0,
  })

  const loading = ref(false)

  const handleCreate = async () => {
    /* 统一校验 */
    const nameCheck = Validator.validate(newCategory.value.name || '', [R.required, R.maxLength(50)], $t)
    if (!nameCheck.valid) {
      toast.error(nameCheck.message || $t('utils.validation.errors.required'))
      return
    }

    loading.value = true

    try {
      const result = await createCategory(newCategory.value)
      if (result.success) {
        emit('categoryCreated', result.data)
        visible.value = false
        resetForm()
      }
    } catch {
      toast.error($t('category.toast.createError'))
    }

    loading.value = false
  }

  const resetForm = () => {
    newCategory.value = {
      name: '',
      description: '',
      sort_order: 0,
    }
  }

  const handleClose = () => {
    visible.value = false
    resetForm()
  }
</script>

<template>
  <CyberDialog v-model="visible" :title="$t('category.dialog.create.title')" width="500px">
    <div class="create-category-form space-y-4 p-4">
      <div class="form-group">
        <label class="mb-2 block text-sm font-medium text-content">
          {{ $t('category.form.name.label') }} <span class="text-error-400">*</span>
        </label>
        <CyberInput v-model="newCategory.name" :placeholder="$t('category.placeholders.name')" maxlength="50" class="w-full" />
      </div>

      <div class="form-group">
        <label class="mb-2 block text-sm font-medium text-content">
          {{ $t('category.form.description.label') }}
        </label>
        <CyberInput
          v-model="newCategory.description"
          type="textarea"
          :placeholder="$t('category.placeholders.description')"
          maxlength="500"
          :rows="4"
          class="w-full"
        />
      </div>

      <div class="form-group">
        <label class="mb-2 block text-sm font-medium text-content">
          {{ $t('category.form.sortOrder.label') }}
        </label>
        <CyberInput
          v-model.number="newCategory.sort_order"
          type="number"
          :placeholder="$t('category.placeholders.sortOrder')"
          class="w-full"
        />
        <p class="mt-1 text-xs text-content-muted">{{ $t('category.form.sortOrder.hint') }}</p>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2 p-4">
        <CyberButton type="secondary" @click="handleClose">{{ $t('category.dialog.create.cancel') }}</CyberButton>
        <CyberButton type="primary" :loading="loading" @click="handleCreate">{{
          $t('category.dialog.create.create')
        }}</CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped lang="scss">
  .create-category-form {
    background: var(--color-background-900);

    .form-group {
      :deep(.cyber-input) {
        background: var(--color-background-800);
        border-color: var(--color-border-default);
        color: var(--color-text-content);

        &:focus {
          border-color: var(--color-brand-500);
        }
      }

      :deep(textarea.cyber-input) {
        resize: vertical;
        min-height: 100px;
      }
    }
  }
</style>
