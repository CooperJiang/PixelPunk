<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { StorageUtil } from '@/utils/storage/storage'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps({
    modelValue: {
      type: Boolean,
      default: false,
    },
    shareKey: {
      type: String,
      default: '',
    },
  })

  const emit = defineEmits(['update:modelValue', 'submit', 'dismiss'])

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  const name = ref('')
  const email = ref('')
  const submitting = ref(false)

  const getStorageKey = () => `visitor_dialog_dismissed_${props.shareKey}`

  const markDialogDismissed = () => {
    if (props.shareKey) {
      StorageUtil.set(getStorageKey(), true, 24)
    }
  }

  const handleSubmit = async () => {
    if (!name.value) {
      return
    }

    try {
      submitting.value = true

      const visitorInfo = {
        name: name.value,
        ...(email.value ? { email: email.value } : {}),
      }

      markDialogDismissed()
      emit('submit', visitorInfo)
    } finally {
      setTimeout(() => {
        submitting.value = false
      }, 500)
    }
  }

  const handleDismiss = () => {
    markDialogDismissed()
    emit('dismiss')
    visible.value = false
  }
</script>

<template>
  <CyberDialog
    v-model="visible"
    :title="$t('share.visitor.dialogTitle')"
    width="400px"
    :append-to-body="true"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :close-on-click-overlay="false"
    :show-close-button="false"
  >
    <div class="visitor-dialog">
      <div class="dialog-header">
        <div class="dialog-icon">
          <i class="fas fa-handshake" />
        </div>
        <h3 class="welcome-title">{{ $t('share.visitor.welcomeTitle') }}</h3>
      </div>

      <div class="dialog-content">
        <p class="dialog-message">
          {{ $t('share.visitor.message') }}
        </p>

        <div class="visitor-form">
          <div class="form-item">
            <label>{{ $t('share.visitor.name.label') }}</label>
            <CyberInput v-model="name" :placeholder="$t('share.visitor.name.placeholder')" />
          </div>
          <div class="form-item">
            <label>{{ $t('share.visitor.email.label') }}</label>
            <CyberInput v-model="email" :placeholder="$t('share.visitor.email.placeholder')" />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer gap-2 p-4">
        <button class="skip-button" @click="handleDismiss">{{ $t('share.visitor.dismiss') }}</button>

        <CyberButton type="primary" :loading="submitting" :disabled="!name" @click="handleSubmit">
          {{ $t('share.visitor.submit') }}
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped>
  .visitor-dialog {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: var(--space-md) 0;
  }

  .dialog-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 1rem;
  }

  .welcome-title {
    color: rgba(var(--color-brand-500-rgb), 0.9);
    margin: var(--space-md) 0 0;
    font-size: 1.2rem;
    font-weight: 600;
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .dialog-icon {
    width: 65px;
    height: 65px;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), rgba(var(--color-error-rgb), 0.1));
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    overflow: hidden;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .dialog-icon i {
    font-size: 1.8rem;
    color: rgba(var(--color-brand-500-rgb), 0.9);
    filter: drop-shadow(0 0 5px rgba(var(--color-brand-500-rgb), 0.4));
  }

  .dialog-content {
    text-align: center;
    width: 100%;
  }

  .dialog-message {
    color: rgba(var(--color-content-rgb), 0.8);
    font-size: 0.9rem;
    margin-bottom: 1.5rem;
    line-height: 1.5;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    padding: var(--space-md);
    border-radius: var(--radius-md);
    border-left: 2px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .visitor-form {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
    margin: var(--space-md) 0;
  }

  .form-item {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
    text-align: left;
  }

  .form-item label {
    color: var(--color-content-default);
    font-size: 0.9rem;
    display: flex;
    align-items: center;
  }

  .dialog-footer {
    display: flex;
    justify-content: center;
    width: 100%;
  }

  .skip-button {
    background: transparent;
    border: none;
    color: var(--color-content-default);
    padding: var(--space-md) 1rem;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
    text-decoration: underline;
  }

  .skip-button:hover {
    color: rgba(var(--color-content-rgb), 0.8);
  }
</style>
