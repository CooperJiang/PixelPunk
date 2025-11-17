<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { ConfirmDialogProps, ConfirmDialogEmits } from './types'

  defineOptions({
    name: 'CyberConfirmDialog',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<ConfirmDialogProps>(), {
    title: '',
    type: 'warning',
    width: '500px',
    confirmText: '',
    cancelText: '',
    loading: false,
    closable: true,
    requireInput: false,
    inputLabel: '',
    inputPlaceholder: '',
    expectedInput: undefined,
  })

  const emit = defineEmits<ConfirmDialogEmits>()

  const inputValue = ref('')

  const isVisible = computed({
    get: () => props.modelValue,
    set: (value: boolean) => emit('update:modelValue', value),
  })

  const confirmButtonType = computed(() => {
    switch (props.type) {
      case 'danger':
        return 'danger'
      case 'warning':
        return 'warning'
      case 'success':
        return 'success'
      default:
        return 'primary'
    }
  })

  const dialogTitle = computed(() => props.title || $t('confirmDialog.title'))
  const dialogConfirmText = computed(() => props.confirmText || $t('actions.confirm'))
  const dialogCancelText = computed(() => props.cancelText || $t('actions.cancel'))
  const dialogInputLabel = computed(() => props.inputLabel || $t('confirmDialog.inputLabel'))
  const dialogInputPlaceholder = computed(() => props.inputPlaceholder || $t('confirmDialog.inputPlaceholder'))

  const getIconClass = (type: string) => {
    const iconMap = {
      warning: 'fas fa-exclamation-triangle',
      danger: 'fas fa-exclamation-circle',
      info: 'fas fa-info-circle',
      success: 'fas fa-check-circle',
    }
    return iconMap[type as keyof typeof iconMap] || iconMap.warning
  }

  const getIconColor = (type: string) => {
    const colorMap = {
      warning: 'var(--color-warning-500)',
      danger: 'var(--color-error-500)',
      info: 'var(--color-brand-500)',
      success: 'var(--color-success-500)',
    }
    return colorMap[type as keyof typeof colorMap] || colorMap.warning
  }

  const formatMessage = (msg: string) => msg.replace(/\n/g, '<br>')

  const handleConfirm = () => {
    if (props.requireInput) {
      emit('confirm', inputValue.value)
    } else {
      emit('confirm')
    }
  }

  const handleCancel = () => {
    inputValue.value = ''
    emit('cancel')
  }

  watch(
    () => props.modelValue,
    (newVal) => {
      if (newVal) {
        inputValue.value = ''
      }
    }
  )
</script>

<template>
  <CyberDialog v-model="isVisible" :title="dialogTitle" :width="width" :closable="closable" @close="handleCancel">
    <div class="confirm-dialog-content">
      <div v-if="type" class="confirm-icon">
        <i :class="getIconClass(type)" :style="{ color: getIconColor(type) }" />
      </div>

      <div class="confirm-text">
        <div class="confirm-message">
          <div v-if="typeof message === 'string'" v-html="formatMessage(message)" />
          <div v-else>
            <div v-for="(line, index) in message" :key="index" class="message-line">
              {{ line }}
            </div>
          </div>
        </div>

        <div v-if="requireInput" class="confirm-input-section">
          <div class="input-label">{{ dialogInputLabel }}</div>
          <CyberInput v-model="inputValue" :placeholder="dialogInputPlaceholder" class="confirm-input" />
          <div v-if="inputValue && inputValue !== expectedInput && expectedInput" class="input-error">
            {{ $t('confirmDialog.inputError') }}
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="confirm-dialog-footer">
        <CyberButton type="secondary" :disabled="loading" @click="handleCancel">
          {{ dialogCancelText }}
        </CyberButton>
        <CyberButton
          :type="confirmButtonType"
          :loading="loading"
          :disabled="requireInput && (!inputValue || (expectedInput && inputValue !== expectedInput))"
          @click="handleConfirm"
        >
          {{ dialogConfirmText }}
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style lang="scss" scoped>
  .confirm-dialog-content {
    @apply flex min-h-[60px] items-start gap-4 py-5;
  }

  .confirm-icon {
    @apply mt-0.5 flex-shrink-0 text-2xl;
  }

  .confirm-text {
    @apply flex-1;
  }

  .confirm-message {
    @apply text-base leading-relaxed text-content;
  }

  .message-line {
    @apply mb-2 last:mb-0;
  }

  .confirm-input-section {
    @apply mt-5;
  }

  .input-label {
    @apply mb-2 text-sm text-content-muted;
  }

  .confirm-input {
    @apply w-full;
  }

  .input-error {
    @apply mt-1 text-xs text-error-500;
  }

  .confirm-dialog-footer {
    @apply flex justify-end gap-3 border-t border-subtle px-4 py-5;
  }
</style>
