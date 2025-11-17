<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { StorageUtil } from '@/utils/storage/storage'
  import { TOKEN_KEY } from '@/constants'
  import type { ShareButtonProps, ShareButtonEmits, ShareItem } from './types'

  defineOptions({
    name: 'ShareButton',
  })

  const props = defineProps<ShareButtonProps>()
  const emit = defineEmits<ShareButtonEmits>()

  const { $t } = useTexts()
  const toast = useToast()

  const selectMode = ref(false)
  const selectedFolders = ref<ShareItem[]>([])
  const selectedImages = ref<ShareItem[]>([])

  const startSelectMode = () => {
    if (props.requireLogin) {
      const token = StorageUtil.get<string>(TOKEN_KEY)
      if (!token) {
        toast.warning($t('components.share.messages.loginRequired'))
        return
      }
    }

    selectMode.value = true
    emit('update:selectMode', true)
    selectedFolders.value = []
    selectedImages.value = []
  }

  const cancelSelectMode = () => {
    selectMode.value = false
    emit('update:selectMode', false)
    selectedFolders.value = []
    selectedImages.value = []
  }

  const isFolderSelected = (folder: ShareItem) => selectedFolders.value.some((f) => f.id === folder.id)

  const isImageSelected = (image: ShareItem) => selectedImages.value.some((img) => img.id === image.id)

  const toggleFolderSelect = (folder: ShareItem) => {
    if (props.requireLogin) {
      const token = StorageUtil.get<string>(TOKEN_KEY)
      if (!token) {
        toast.warning($t('components.share.messages.loginRequired'))
        return
      }
    }

    const index = selectedFolders.value.findIndex((f) => f.id === folder.id)
    if (index === -1) {
      selectedFolders.value.push(folder)
    } else {
      selectedFolders.value.splice(index, 1)
    }
  }

  const toggleImageSelect = (image: ShareItem) => {
    if (props.requireLogin) {
      const token = StorageUtil.get<string>(TOKEN_KEY)
      if (!token) {
        toast.warning($t('components.share.messages.loginRequired'))
        return
      }
    }

    const index = selectedImages.value.findIndex((img) => img.id === image.id)
    if (index === -1) {
      selectedImages.value.push(image)
    } else {
      selectedImages.value.splice(index, 1)
    }
  }

  const openShareDialog = () => {
    if (selectedFolders.value.length > 0 || selectedImages.value.length > 0) {
      emit('update:shareDialogVisible', true)
    } else {
      toast.info($t('components.share.messages.noItemsSelected'))
    }
  }

  const handleShareCreated = (_share: ShareItem) => {
    cancelSelectMode()
    emit('update:shareDialogVisible', false)
  }

  const toggleSelectAllFolders = (folders: ShareItem[]) => {
    if (!folders || folders.length === 0) {
      return
    }

    const allSelected = folders.every((folder) => isFolderSelected(folder))

    if (allSelected) {
      selectedFolders.value = []
    } else {
      selectedFolders.value = [...folders]
    }
  }

  const invertFolderSelection = (folders: ShareItem[]) => {
    if (!folders || folders.length === 0) {
      return
    }

    const currentlySelected = selectedFolders.value.map((f) => f.id)
    selectedFolders.value = folders.filter((folder) => !currentlySelected.includes(folder.id))
  }

  const toggleSelectAllImages = (images: ShareItem[]) => {
    if (!images || images.length === 0) {
      return
    }

    const allSelected = images.every((image) => isImageSelected(image))

    if (allSelected) {
      selectedImages.value = []
    } else {
      selectedImages.value = [...images]
    }
  }

  const invertImageSelection = (images: ShareItem[]) => {
    if (!images || images.length === 0) {
      return
    }

    const currentlySelected = selectedImages.value.map((img) => img.id)
    selectedImages.value = images.filter((image) => !currentlySelected.includes(image.id))
  }

  const isAllFoldersSelected = computed(() => {
    if (!props.folders || props.folders.length === 0) {
      return false
    }
    return props.folders.every((folder) => isFolderSelected(folder))
  })

  defineExpose({
    selectMode,
    selectedFolders,
    selectedImages,
    selectedFiles: selectedImages, // 别名，用于兼容 ShareButtonRef 接口
    isFolderSelected,
    isImageSelected,
    isFileSelected: isImageSelected, // 别名，用于兼容 ShareButtonRef 接口
    toggleFolderSelect,
    toggleImageSelect,
    toggleFileSelect: toggleImageSelect, // 别名，用于兼容 ShareButtonRef 接口
    cancelSelectMode,
    toggleSelectAllFolders,
    toggleSelectAllImages,
    invertFolderSelection,
    invertImageSelection,
    isAllFoldersSelected,
    handleShareCreated,
  })
</script>

<template>
  <div class="share-controls">
    <div v-if="selectMode" class="selection-info">
      <span v-if="selectedFolders.length > 0" class="selection-count folder-count">
        <i class="fas fa-folder" /> {{ selectedFolders.length }}{{ $t('components.share.selectionInfo.folders') }}
      </span>
      <span v-if="selectedImages.length > 0" class="selection-count image-count">
        <i class="fas fa-image" /> {{ selectedImages.length }}{{ $t('components.share.selectionInfo.images') }}
      </span>
    </div>
    <CyberButton v-if="!selectMode" type="outlined" icon="share-alt" @click="startSelectMode">
      {{ $t('components.share.createShare') }}
    </CyberButton>
    <div v-else class="select-mode-controls">
      <CyberButton
        type="primary"
        icon="share-alt"
        :disabled="selectedFolders.length === 0 && selectedImages.length === 0"
        @click="openShareDialog"
      >
        {{ $t('components.share.shareSelected') }}({{ selectedFolders.length + selectedImages.length }})
      </CyberButton>
      <CyberButton type="outlined" icon="times" @click="cancelSelectMode">{{ $t('components.share.cancel') }}</CyberButton>
    </div>
  </div>
</template>

<style scoped>
  .share-controls {
    @apply flex flex-wrap items-center gap-2;
    min-height: 32px;
  }

  .selection-info {
    @apply flex gap-2;
    align-items: center;
    height: 32px;
  }

  .selection-count {
    @apply flex items-center gap-1 rounded border px-2 text-xs font-medium;
    height: 24px;
    line-height: 1;
  }

  .folder-count {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .image-count {
    @apply whitespace-nowrap;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .select-mode-controls {
    @apply flex flex-wrap gap-2;
  }
</style>
