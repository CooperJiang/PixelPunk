import { ref } from 'vue'
import { restoreScrollPosition, saveScrollPosition } from '@/utils/ui/domUtils'

/**
 * 分享页面导航功能
 * 处理面包屑导航和文件夹导航
 */
export function useShareNavigation() {
  const currentFolderId = ref<string | undefined>()

  const breadcrumbItems = ref<{ id: string; name: string }[]>([])

  const scrollPositions = ref<Record<string, number>>({
    root: 0, // 根目录的滚动位置
  })

  const getCurrentFolderKey = (): string => currentFolderId.value || 'root'

  const saveCurrentScrollPosition = (): void => {
    const key = getCurrentFolderKey()
    scrollPositions.value = saveScrollPosition(key, scrollPositions.value)
  }

  const restoreCurrentScrollPosition = (): void => {
    const key = getCurrentFolderKey()
    restoreScrollPosition(key, scrollPositions.value)
  }

  const handleFolderClick = (folder: any, callback?: () => void): void => {
    saveCurrentScrollPosition()

    currentFolderId.value = folder.id

    breadcrumbItems.value.push({
      id: folder.id,
      name: folder.name,
    })

    if (callback) {
      callback()
    }
  }

  const handleBreadcrumbClick = (item: { id: string; name: string } | null, callback?: () => void): void => {
    saveCurrentScrollPosition()

    if (item === null) {
      currentFolderId.value = undefined
      breadcrumbItems.value = []
    } else {
      const index = breadcrumbItems.value.findIndex((bread) => bread.id === item.id)
      if (index !== -1) {
        currentFolderId.value = item.id
        breadcrumbItems.value = breadcrumbItems.value.slice(0, index + 1)
      }
    }

    if (callback) {
      callback()
    }
  }

  return {
    currentFolderId,
    breadcrumbItems,
    scrollPositions,
    getCurrentFolderKey,
    saveCurrentScrollPosition,
    restoreCurrentScrollPosition,
    handleFolderClick,
    handleBreadcrumbClick,
  }
}
