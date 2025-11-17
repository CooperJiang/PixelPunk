import { computed, ref, reactive } from 'vue'
import { useToast } from '@/components'
import type { ContextMenuItem } from '@/components/CyberContextMenu/types'
import type { ImageInfo, FolderInfo } from '@/api/types'
import { getFolderList } from '@/api/folder'
import { moveFiles, deleteFile } from '@/api/file'
import { showConfirm } from '@/utils/dialog'
import { useTexts } from '@/composables/useTexts'

export interface ImageContextMenuOptions {
  onImageMoved?: (fileId: string, targetFolderId?: string) => void
  onImageDeleted?: (fileId: string) => void
  onPreview?: (image: ImageInfo) => void
  onCopyLink?: (image: ImageInfo) => void
  onDownload?: (image: ImageInfo) => void
  onDelete?: (image: ImageInfo) => void
}

export function useFileContextMenu(options: ImageContextMenuOptions = {}) {
  const showContextMenu = ref(false)
  const contextMenuPosition = reactive({ x: 0, y: 0 })
  const selectedImage = ref<ImageInfo | null>(null)
  const toast = useToast()
  const { $t } = useTexts()

  const buildFolderMenuItem = (folder: FolderInfo): ContextMenuItem => {
    return {
      key: `folder-${folder.id}`,
      label: folder.name,
      icon: 'fas fa-folder',
      onClick: async () => {
        if (!selectedImage.value) {
          toast.error($t('folders.fileContextMenu.toast.noFileSelected'))
          return
        }
        const { id } = selectedImage.value
        try {
          await moveFiles([id], folder.id)
          toast.success($t('folders.fileContextMenu.toast.moveToFolderSuccess', { name: folder.name }))
          options.onImageMoved?.(id, folder.id)
        } catch {
          toast.error($t('folders.fileContextMenu.toast.moveFailed'))
        }
      },
      hasAsyncChildren: folder.has_children,
      loadChildren: folder.has_children
        ? async () => {
            try {
              const result = await getFolderList(folder.id)
              if (result.success && result.data) {
                const childItems = result.data.map(buildFolderMenuItem)
                return childItems
              }
              return []
            } catch {
              throw new Error($t('folders.fileContextMenu.toast.loadSubfoldersFailed'))
            }
          }
        : undefined,
    }
  }

  const loadRootFolders = async (): Promise<ContextMenuItem[]> => {
    try {
      const result = await getFolderList()
      if (result.success && result.data) {
        const folderItems = result.data.map(buildFolderMenuItem)
        return [
          {
            key: 'root-folder',
            label: $t('folders.fileContextMenu.menu.root'),
            icon: 'fas fa-home',
            onClick: async () => {
              if (!selectedImage.value) {
                toast.error($t('folders.fileContextMenu.toast.noFileSelected'))
                return
              }
              const { id } = selectedImage.value
              try {
                await moveFiles([id], undefined)
                toast.success($t('folders.fileContextMenu.toast.moveToRootSuccess'))
                options.onImageMoved?.(id, undefined)
              } catch {
                toast.error($t('folders.fileContextMenu.toast.moveFailed'))
              }
            },
          },
          ...folderItems,
        ]
      }
      return []
    } catch {
      throw new Error($t('folders.fileContextMenu.toast.loadFoldersFailed'))
    }
  }

  const contextMenuItems = computed<ContextMenuItem[]>(() => [
    {
      key: 'preview',
      label: $t('folders.fileContextMenu.menu.preview'),
      icon: 'fas fa-eye',
      onClick: () => {
        if (!selectedImage.value) return
        options.onPreview?.(selectedImage.value)
      },
    },
    {
      key: 'open-new',
      label: $t('folders.fileContextMenu.menu.openNew'),
      icon: 'fas fa-external-link-alt',
      onClick: () => {
        if (!selectedImage.value) return
        const url = selectedImage.value.full_url || selectedImage.value.url || ''
        if (!url) {
          toast.error($t('folders.fileContextMenu.toast.noValidImageLink'))
          return
        }
        window.open(url, '_blank')
      },
    },
    {
      key: 'divider-preview-copy',
      label: '',
      divided: true,
    },
    {
      key: 'copy-link',
      label: $t('folders.fileContextMenu.menu.copyLink'),
      icon: 'fas fa-link',
      children: [
        {
          key: 'copy-thumb',
          label: $t('folders.fileContextMenu.menu.copyThumb'),
          icon: 'fas fa-image',
          onClick: async () => {
            if (!selectedImage.value) return
            const url = selectedImage.value.full_thumb_url || selectedImage.value.thumb_url || ''
            if (!url) {
              toast.error($t('folders.fileContextMenu.toast.noThumbLink'))
              return
            }
            try {
              await navigator.clipboard.writeText(url)
              toast.success($t('folders.fileContextMenu.toast.thumbLinkCopied'))
            } catch {
              toast.error($t('folders.fileContextMenu.toast.copyFailed'))
            }
          },
        },
        {
          key: 'copy-original',
          label: $t('folders.fileContextMenu.menu.copyOriginal'),
          icon: 'fas fa-photo-film',
          onClick: async () => {
            if (!selectedImage.value) return
            const url = selectedImage.value.full_url || selectedImage.value.url || ''
            if (!url) {
              toast.error($t('folders.fileContextMenu.toast.noOriginalLink'))
              return
            }
            try {
              await navigator.clipboard.writeText(url)
              toast.success($t('folders.fileContextMenu.toast.originalLinkCopied'))
            } catch {
              toast.error($t('folders.fileContextMenu.toast.copyFailed'))
            }
          },
        },
        {
          key: 'copy-short',
          label: $t('folders.fileContextMenu.menu.copyShort'),
          icon: 'fas fa-link',
          onClick: async () => {
            if (!selectedImage.value) return
            const url = selectedImage.value.short_url || ''
            if (!url) {
              toast.error($t('folders.fileContextMenu.toast.noShortLink'))
              return
            }
            try {
              await navigator.clipboard.writeText(url)
              toast.success($t('folders.fileContextMenu.toast.shortLinkCopied'))
            } catch {
              toast.error($t('folders.fileContextMenu.toast.copyFailed'))
            }
          },
        },
      ],
    },
    {
      key: 'divider-after-copy',
      label: '',
      divided: true,
    },
    {
      key: 'download',
      label: $t('folders.fileContextMenu.menu.download'),
      icon: 'fas fa-download',
      onClick: () => {
        if (!selectedImage.value) return
        if (options.onDownload) {
          options.onDownload(selectedImage.value)
          return
        }
        const a = document.createElement('a')
        a.href = selectedImage.value.full_url || selectedImage.value.full_thumb_url || ''
        a.download = selectedImage.value.display_name || selectedImage.value.name || 'download'
        a.target = '_blank'
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
      },
    },
    {
      key: 'move-to',
      label: $t('folders.fileContextMenu.menu.moveTo'),
      icon: 'fas fa-arrows-alt',
      hasAsyncChildren: true,
      loadChildren: loadRootFolders,
    },
    {
      key: 'divider-after-move',
      label: '',
      divided: true,
    },
    {
      key: 'delete',
      label: $t('folders.fileContextMenu.menu.delete'),
      icon: 'fas fa-trash',
      danger: true,
      onClick: async () => {
        if (!selectedImage.value) return
        if (options.onDelete) {
          options.onDelete(selectedImage.value)
          return
        }
        const name = selectedImage.value.display_name || selectedImage.value.name || $t('common.file')
        const confirmed = showConfirm($t('folders.fileContextMenu.confirm.deleteFile', { name }))
        if (confirmed) {
          try {
            const { id } = selectedImage.value
            await deleteFile(id)
            toast.success($t('folders.fileContextMenu.toast.fileDeleted'))
            options.onImageDeleted?.(id)
          } catch {
            toast.error($t('folders.fileContextMenu.toast.deleteFailed'))
          }
        }
      },
    },
  ])

  const showImageContextMenu = (event: MouseEvent, image: ImageInfo) => {
    event.preventDefault()
    event.stopPropagation()

    selectedImage.value = image
    contextMenuPosition.x = event.clientX
    contextMenuPosition.y = event.clientY
    showContextMenu.value = true
  }

  const hideContextMenu = () => {
    showContextMenu.value = false
    selectedImage.value = null
  }

  return {
    showContextMenu,
    contextMenuPosition,
    contextMenuItems,
    selectedImage,
    showImageContextMenu,
    hideContextMenu,
  }
}
