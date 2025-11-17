<script setup lang="ts">
  import { computed, ref, watch, type PropType } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { BatchTagOperation, TagInfo } from '@/api/types/index'

  const { $t } = useTexts()
  const props = defineProps({
    visible: {
      type: Boolean,
      default: false,
    },
    selectedTags: {
      type: Array as PropType<TagInfo[]>,
      default: () => [],
    },
    availableTags: {
      type: Array as PropType<TagInfo[]>,
      default: () => [],
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

  const maxDisplayTags = 6

  const operationOptions = computed(() => [
    { label: $t('admin.tags.batchOperation.deleteOption'), value: 'delete' },
    { label: $t('admin.tags.batchOperation.mergeOption'), value: 'merge' },
  ])

  const formData = ref<{
    action: 'delete' | 'merge'
    targetId: number | null
  }>({
    action: 'delete',
    targetId: null,
  })

  const targetError = ref('')
  // i18n validation
  import Validator, { getValidationRules } from '@/utils/validation/validator'
  const R = getValidationRules($t)

  const isFormValid = computed(() => {
    if (formData.value.action === 'merge') {
      return formData.value.targetId !== null && !targetError.value
    }
    return true
  })

  const submitButtonText = computed(() => {
    if (props.isSubmitting) {
      return $t('status.processing')
    }

    if (formData.value.action === 'delete') {
      return $t('admin.tags.batchOperation.deleteButton').replace('{count}', String(props.selectedTags.length))
    }
    return $t('admin.tags.batchOperation.mergeButton')
  })

  const targetOptions = computed(() => {
    const selectedIds = props.selectedTags.map((tag) => tag.id)
    const validTags = props.availableTags.filter((tag) => !selectedIds.includes(tag.id))

    return validTags.map((tag) => ({
      label: `${tag.name} (${$t('admin.tags.batchOperation.usageTimesSuffix').replace('{count}', String(tag.count))})`,
      value: tag.id,
    }))
  })

  const displayedTags = computed(() => props.selectedTags.slice(0, maxDisplayTags))

  const warningText = computed(() => {
    if (formData.value.action === 'delete') {
      return $t('admin.tags.batchOperation.deleteWarning').replace('{count}', String(getTotalImagesCount()))
    }
    return $t('admin.tags.batchOperation.mergeWarning').replace('{count}', String(getTotalImagesCount()))
  })

  const warningClass = computed(() => (formData.value.action === 'delete' ? 'warning-delete' : 'warning-merge'))

  const getTotalImagesCount = () => props.selectedTags.reduce((sum, tag) => sum + tag.count, 0)

  const resetForm = () => {
    formData.value = {
      action: 'delete',
      targetId: null,
    }
    targetError.value = ''
  }

  const onClose = () => {
    resetForm()
    emit('close')
  }

  const submitForm = () => {
    if (formData.value.action === 'merge') {
      const check = Validator.validate(formData.value.targetId, [R.required], $t)
      if (!check.valid) {
        targetError.value = check.message || $t('admin.tags.batchOperation.targetError')
        return
      }
    }

    const operation: BatchTagOperation = {
      action: formData.value.action,
      tag_ids: props.selectedTags.map((tag) => tag.id),
    }

    if (formData.value.action === 'merge') {
      operation.target = formData.value.targetId!
    }

    emit('submit', operation)
  }

  watch(
    () => props.selectedTags,
    (newTags) => {
      const selectedIds = newTags.map((tag) => tag.id)
      if (formData.value.targetId && selectedIds.includes(formData.value.targetId)) {
        formData.value.targetId = null
      }
    },
    { deep: true }
  )
</script>

<template>
  <CyberDialog v-model="dialogVisible" :title="$t('admin.tags.batchOperation.title')" width="480" @close="onClose">
    <div class="batch-operation-form">
      <div class="selected-tags-info mb-3 rounded-md p-3">
        <div class="flex items-center justify-between">
          <span class="text-content-heading">{{ $t('admin.tags.batchOperation.selectedLabel') }}</span>
          <span class="text-accent font-medium"
            >{{ selectedTags.length }} {{ $t('admin.tags.batchOperation.selectedCount') }}</span
          >
        </div>

        <div v-if="selectedTags.length > 0" class="selected-tags-list mt-2">
          <div class="flex flex-wrap gap-1">
            <div v-for="tag in displayedTags" :key="tag.id" class="tag-chip">
              {{ tag.name }}
              <span class="ml-1 text-xs text-content">({{ tag.count }})</span>
            </div>

            <div v-if="selectedTags.length > maxDisplayTags" class="tag-chip more-tag">
              +{{ selectedTags.length - maxDisplayTags }} {{ $t('admin.tags.batchOperation.moreCount') }}
            </div>
          </div>
        </div>
      </div>

      <div class="form-section">
        <div class="cyber-form-item">
          <label class="mb-2 block text-sm text-content-heading">{{ $t('admin.tags.batchOperation.operationType') }}</label>
          <div class="operation-type-selector mb-1">
            <CyberRadioGroup v-model="formData.action" :options="operationOptions" layout="horizontal" />
          </div>
        </div>

        <div v-if="formData.action === 'merge'" class="cyber-form-item mt-4">
          <label class="mb-2 block text-sm text-content-heading">{{ $t('admin.tags.batchOperation.targetLabel') }}</label>
          <CyberDropdown
            v-model="formData.targetId"
            :options="targetOptions"
            :placeholder="$t('admin.tags.batchOperation.targetPlaceholder')"
            class="w-full"
            :status="targetError ? 'error' : ''"
            @update:model-value="targetError = ''"
          />
          <div v-if="targetError" class="text-accent mt-1 text-xs">{{ targetError }}</div>
          <p v-else class="text-content-content-muted mt-1 text-xs">{{ $t('admin.tags.batchOperation.mergeDesc') }}</p>
        </div>
      </div>

      <div class="warning-section mt-4 rounded-md p-3" :class="warningClass">
        <div class="flex items-start">
          <i class="fas fa-exclamation-triangle mr-2 mt-0.5 text-sm" />
          <div>
            <p class="warning-title text-sm font-medium">{{ $t('admin.tags.batchOperation.warningTitle') }}</p>
            <p class="warning-desc mt-1 text-xs">{{ warningText }}</p>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2 p-4">
        <CyberButton type="secondary" @click="onClose">{{ $t('admin.tags.batchOperation.cancel') }}</CyberButton>
        <CyberButton
          :type="formData.action === 'delete' ? 'danger' : 'primary'"
          :loading="isSubmitting"
          :disabled="!isFormValid || isSubmitting"
          @click="submitForm"
        >
          {{ submitButtonText }}
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped lang="scss">
  .batch-operation-form {
    width: 100%;
  }

  .selected-tags-info {
    background: rgba(var(--color-background-900-rgb), 0.5);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .tag-chip {
    @apply rounded-md px-2 py-0.5 text-xs;
    background: rgba(var(--color-background-600-rgb), 0.7);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-content-heading);
    transition: all var(--transition-normal) var(--ease-out);
  }

  .more-tag {
    @apply px-2 py-0.5 text-xs;
    background: rgba(var(--color-background-400-rgb), 0.5);
    border: 1px solid rgba(var(--color-border-default-rgb), 0.2);
    color: var(--color-content-muted);
  }

  .warning-section {
    border: 1px solid rgba(var(--color-border-default-rgb), 0.3);
    transition: all var(--transition-normal) var(--ease-out);
  }

  .warning-section.warning-delete {
    background: rgba(var(--color-error-rgb), 0.1);
    border-color: rgba(var(--color-error-rgb), 0.3);
  }

  .warning-section.warning-delete .warning-title {
    color: var(--color-error-500);
  }

  .warning-section.warning-merge {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .warning-section.warning-merge .warning-title {
    color: var(--color-brand-500);
  }

  .warning-desc {
    color: var(--color-content-default);
  }
</style>
