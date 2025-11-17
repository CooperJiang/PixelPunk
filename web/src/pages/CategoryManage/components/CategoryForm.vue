<script setup lang="ts">
  import type { ImageCategory } from '@/api/types/category'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'CategoryForm',
  })

  const { $t } = useTexts()
  const category = defineModel<ImageCategory | null>({ required: true })
</script>

<template>
  <div v-if="category" class="category-form space-y-4 p-4">
    <div class="form-group">
      <label class="mb-2 block text-sm font-medium text-content">
        {{ $t('category.form.name.label') }} <span class="text-error-400">*</span>
      </label>
      <CyberInput v-model="category.name" :placeholder="$t('category.placeholders.name')" :maxlength="50" class="w-full" />
    </div>

    <div class="form-group">
      <label class="mb-2 block text-sm font-medium text-content">
        {{ $t('category.form.description.label') }}
      </label>
      <CyberInput
        v-model="category.description"
        type="textarea"
        :placeholder="$t('category.placeholders.description')"
        :maxlength="500"
        :rows="4"
        class="w-full"
      />
    </div>

    <div class="form-group">
      <label class="mb-2 block text-sm font-medium text-content">
        {{ $t('category.form.sortOrder.label') }}
      </label>
      <CyberInput
        v-model.number="category.sort_order"
        type="number"
        :placeholder="$t('category.placeholders.sortOrder')"
        class="w-full"
      />
      <p class="mt-1 text-xs text-content-muted">{{ $t('category.form.sortOrder.hint') }}</p>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .category-form {
    background: var(--color-background-900);
  }

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
</style>
