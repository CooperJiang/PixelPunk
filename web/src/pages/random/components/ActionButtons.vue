<script setup lang="ts">
  import type { ImageInfo } from '@/api/types'
  import { downloadFileQuick } from '@/utils/file/downloader'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    imageData: ImageInfo
  }

  interface Emits {
    (e: 'refresh'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const toast = useToast()
  const { $t } = useTexts()

  const handleRefresh = () => {
    emit('refresh')
  }

  const handleDownload = async () => {
    try {
      const fileName = props.imageData.display_name || props.imageData.original_name || `image_${props.imageData.id}`
      await downloadFileQuick(props.imageData.id, fileName)
      toast.success($t('random.actionButtons.downloadSuccess', { fileName }))
    } catch (error) {
      console.error($t('random.actionButtons.downloadFailedLog'), error)
      toast.error($t('random.actionButtons.downloadFailed'))
    }
  }

  const handleOpenInNewWindow = () => {
    window.open(props.imageData.full_url, '_blank')
  }
</script>

<template>
  <div class="action-section">
    <CyberButton type="primary" icon="download" @click="handleDownload" size="large">
      {{ $t('random.actionButtons.downloadOriginal') }}
    </CyberButton>

    <CyberButton type="secondary" icon="external-link-alt" @click="handleOpenInNewWindow" size="large">
      {{ $t('random.actionButtons.openInNewWindow') }}
    </CyberButton>

    <CyberButton type="outlined" icon="dice" @click="handleRefresh" size="large">
      {{ $t('random.actionButtons.nextOne') }}
    </CyberButton>
  </div>
</template>

<style scoped>
  .action-section {
    margin-top: 2rem;
    display: flex;
    gap: 1rem;
    justify-content: center;
    flex-wrap: wrap;
  }

  @media (max-width: 768px) {
    .action-section {
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
    }
  }
</style>
