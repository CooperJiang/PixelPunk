<script setup lang="ts">
  import { computed, ref } from 'vue'

  const props = defineProps({
    modelValue: {
      type: Boolean,
      default: false,
    },
    shareName: {
      type: String,
      default: '',
    },
  })

  const emit = defineEmits(['update:modelValue', 'verify'])

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  const password = ref('')
  const verifying = ref(false)

  const handleVerify = async () => {
    if (!password.value) {
      return
    }

    verifying.value = true
    emit('verify', password.value)

    setTimeout(() => {
      verifying.value = false
    }, 500)
  }
</script>

<template>
  <CyberDialog
    v-model="visible"
    :title="$t('share.passwordDialog.title')"
    width="400px"
    :append-to-body="true"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :close-on-click-overlay="false"
    :show-close-button="false"
  >
    <div class="password-dialog">
      <div class="dialog-icon">
        <i class="fas fa-lock" />
      </div>

      <div class="dialog-content">
        <div v-if="shareName" class="share-name">{{ shareName }}</div>
        <p class="dialog-message">{{ $t('share.passwordDialog.message') }}</p>

        <div class="password-form">
          <CyberInput
            v-model="password"
            type="password"
            :placeholder="$t('share.passwordDialog.placeholder')"
            @keyup.enter="handleVerify"
          />
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer p-4">
        <CyberButton type="secondary" :loading="verifying" :disabled="!password" @click="handleVerify">
          {{ $t('share.passwordDialog.verify') }}
        </CyberButton>
      </div>
    </template>
  </CyberDialog>
</template>

<style scoped lang="scss">
  .password-dialog {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: var(--space-md) 0;
  }

  .dialog-icon {
    width: 65px;
    height: 65px;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: var(--space-lg);
    position: relative;
    overflow: hidden;
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.2);

    i {
      font-size: var(--text-3xl);
      color: var(--color-brand-400);
      filter: drop-shadow(0 0 5px rgba(var(--color-brand-500-rgb), 0.4));
    }
  }

  .dialog-content {
    text-align: center;
    width: 100%;
  }

  .share-name {
    font-size: var(--text-lg);
    font-weight: 600;
    color: var(--color-content-heading);
    margin-bottom: var(--space-sm);
  }

  .dialog-message {
    color: var(--color-content-muted);
    font-size: var(--text-sm);
    margin-bottom: var(--space-md);
  }

  .password-form {
    margin: var(--space-sm) 0;
  }

  .dialog-footer {
    display: flex;
    justify-content: center;
  }
</style>
