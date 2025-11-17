<script setup lang="ts">
  import { ref, computed } from 'vue'
  import type { CommunityDialogProps } from './types'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'CommunityDialog',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<CommunityDialogProps>(), {
    type: 'wechat',
    qrImage: '',
    contactInfo: '',
  })

  const notificationDialogRef = ref()

  /* 根据类型和配置计算显示内容 */
  const dialogConfig = computed(() => {
    if (props.type === 'qq') {
      const hasQrImage = !!props.qrImage
      const hasGroupNumber = !!props.contactInfo
      return {
        title: $t('components.communityDialog.qq.title'),
        imageUrl: props.qrImage,
        imageAlt: $t('components.communityDialog.qq.qrCodeAlt'),
        mainTitle: hasQrImage ? $t('components.communityDialog.qq.scanToJoin') : '',
        subtitle: hasQrImage
          ? props.contactInfo
            ? $t('components.communityDialog.qq.groupNumber', { number: props.contactInfo })
            : ''
          : '',
        groupNumber: hasGroupNumber && !hasQrImage ? props.contactInfo : '',
        showImage: hasQrImage,
        showGroupNumberOnly: hasGroupNumber && !hasQrImage,
      }
    } else {
      const hasQrImage = !!props.qrImage
      return {
        title: $t('components.communityDialog.wechat.title'),
        imageUrl: props.qrImage,
        imageAlt: $t('components.communityDialog.wechat.qrCodeAlt'),
        mainTitle: hasQrImage
          ? $t('components.communityDialog.wechat.scanToJoin')
          : $t('components.communityDialog.wechat.groupInfo'),
        subtitle: props.contactInfo ? $t('components.communityDialog.wechat.addContact', { contact: props.contactInfo }) : '',
        showImage: hasQrImage,
      }
    }
  })

  const show = () => {
    notificationDialogRef.value?.show()
  }

  const hide = () => {
    notificationDialogRef.value?.hide()
  }

  defineExpose({
    show,
    hide,
  })
</script>

<template>
  <CyberNotificationDialog ref="notificationDialogRef" :title="dialogConfig.title" :width="320">
    <div class="community-content">
      <div v-if="dialogConfig.showImage" class="qr-code-container">
        <img
          :src="dialogConfig.imageUrl"
          :alt="dialogConfig.imageAlt"
          class="qr-code"
          @error="$event.target.style.display = 'none'"
        />
      </div>

      <div v-if="dialogConfig.showGroupNumberOnly" class="group-number-only">
        <div class="qq-icon">
          <i class="fab fa-qq"></i>
        </div>
        <div class="group-info">
          <div class="group-label">{{ $t('components.communityDialog.qq.groupNumberLabel') }}</div>
          <div class="large-group-number">{{ dialogConfig.groupNumber }}</div>
          <div class="join-hint">{{ $t('components.communityDialog.qq.joinHint') }}</div>
        </div>
      </div>

      <div v-else class="description">
        <h4 v-if="dialogConfig.mainTitle" class="main-title">{{ dialogConfig.mainTitle }}</h4>
        <p v-if="dialogConfig.subtitle" class="subtitle">{{ dialogConfig.subtitle }}</p>
      </div>
    </div>
  </CyberNotificationDialog>
</template>

<style scoped>
  .community-content {
    @apply flex flex-col items-center py-2 text-center;
  }

  .qr-code-container {
    @apply mb-4 rounded-2xl border border-default bg-background-700 p-3 shadow-cyber-md;
  }

  .qr-code {
    @apply block h-40 w-40 rounded-xl;
  }

  .description {
    @apply max-w-sm space-y-3;
  }

  .main-title {
    @apply mb-0 text-xl font-bold leading-tight text-content-heading;
  }

  .subtitle {
    @apply m-0 text-sm leading-relaxed text-content-muted opacity-80;
  }

  .group-number-only {
    @apply flex flex-col items-center justify-center space-y-3 py-3;
  }

  .qq-icon {
    @apply flex h-12 w-12 items-center justify-center rounded-full;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1) 0%, rgba(var(--color-brand-500-rgb), 0.05) 100%);
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .qq-icon i {
    @apply text-2xl;
    color: var(--color-brand-500);
  }

  .group-info {
    @apply flex flex-col items-center space-y-1 text-center;
  }

  .group-label {
    @apply text-xs font-medium uppercase tracking-wider text-content-muted;
    opacity: 0.7;
  }

  .large-group-number {
    @apply rounded-lg px-3 py-1 text-xl font-bold tracking-wider text-content-heading;
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .join-hint {
    @apply text-xs text-content-muted;
    opacity: 0.6;
  }
</style>
