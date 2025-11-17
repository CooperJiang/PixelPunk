<script setup lang="ts">
  import { computed } from 'vue'
  import { useLayoutStore } from '@/store/layout'
  import { useTexts } from '@/composables/useTexts'
  import type { FolderInfo } from '@/api/types/index'
  import ShareButton from '@/components/Share/ShareButton.vue'
  import FolderBreadcrumb from './FolderBreadcrumb.vue'

  defineOptions({
    name: 'PageHeader',
  })

  const { $t } = useTexts()

  interface Props {
    breadcrumbItems: FolderInfo[]
    folders: any[]
    images: any[]
    currentFolderId?: string | null
    selectMode: boolean
  }

  const props = defineProps<Props>()

  const layoutStore = useLayoutStore()

  /* 是否为 Left 布局 */
  const isLeftLayout = computed(() => layoutStore.mode === 'left')

  /* Top 布局的面包屑数据（只显示文件夹层级，不加前缀） */
  /* 直接使用 props.breadcrumbItems，确保"根目录"始终存在 */
  const topLayoutBreadcrumbItems = computed(() => {
    if (props.breadcrumbItems.length === 0 || props.breadcrumbItems[0].id !== null) {
      return [{ id: null, name: $t('folders.rootFolder') }, ...props.breadcrumbItems]
    }
    return props.breadcrumbItems
  })

  interface Emits {
    (e: 'breadcrumb-click', item: FolderInfo): void
    (e: 'navigate-to-upload'): void
    (e: 'show-create-dialog'): void
    (e: 'update:select-mode', value: boolean): void
    (e: 'update:share-dialog-visible', value: boolean): void
  }

  const emit = defineEmits<Emits>()

  const shareButtonRef = ref()

  const uploadButtonText = computed(() => {
    const currentFolder = props.breadcrumbItems
    if (currentFolder.length <= 1) {
      return $t('folders.pageHeader.upload.default')
    }
    const currentFolderName = currentFolder[currentFolder.length - 1]?.name
    if (currentFolderName && currentFolderName.length > 8) {
      return $t('folders.pageHeader.upload.toFolder', { name: `${currentFolderName.substring(0, 6)}...` })
    }
    return currentFolderName
      ? $t('folders.pageHeader.upload.toFolder', { name: currentFolderName })
      : $t('folders.pageHeader.upload.default')
  })

  const handleBreadcrumbClick = (item: FolderInfo) => {
    emit('breadcrumb-click', item)
  }

  const handleNavigateToUpload = () => {
    emit('navigate-to-upload')
  }

  const handleShowCreateDialog = () => {
    emit('show-create-dialog')
  }

  const handleSelectModeChange = (value: boolean) => {
    emit('update:select-mode', value)
  }

  const handleShareDialogVisibleChange = (value: boolean) => {
    emit('update:share-dialog-visible', value)
  }

  defineExpose({
    shareButtonRef,
  })
</script>

<template>
  <div class="page-header py-2">
    <div v-if="!isLeftLayout" class="header-left">
      <h1 class="page-title">{{ $t('folders.pageHeader.title') }}</h1>
      <FolderBreadcrumb :items="topLayoutBreadcrumbItems" :show-icon="true" class="breadcrumb" @click="handleBreadcrumbClick" />
    </div>

    <div></div>

    <div class="header-actions">
      <ShareButton
        ref="shareButtonRef"
        :folders="folders"
        :images="images"
        @update:select-mode="handleSelectModeChange"
        @update:share-dialog-visible="handleShareDialogVisibleChange"
      />
      <CyberButton type="primary" icon="upload" custom-class="upload-btn" @click="handleNavigateToUpload">
        {{ uploadButtonText }}
      </CyberButton>
      <CyberButton type="outlined" icon="folder-plus" custom-class="new-folder-btn" @click="handleShowCreateDialog">
        {{ $t('folders.pageHeader.actions.newFolder') }}
      </CyberButton>
    </div>
  </div>
</template>

<style scoped>
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 0.4rem;
    flex-shrink: 0;
    min-height: 56px;
    gap: 1rem;
  }

  .header-left {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
    flex: 1;
    min-width: 0;
    overflow: hidden;
  }

  .header-actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    flex-shrink: 0;
  }

  .page-title {
    margin: 0;
    color: var(--color-content-heading);
    font-size: 1.75rem;
    font-weight: 700;
    background: linear-gradient(90deg, var(--color-brand-500), var(--color-error-500));
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    text-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.3);
    letter-spacing: 0.5px;
  }

  .breadcrumb {
    margin-top: 0.2rem;
    width: 100%;
    min-width: 0;
  }
</style>
