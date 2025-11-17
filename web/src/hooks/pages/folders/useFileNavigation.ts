import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { folderApi } from '@/api'
import { useFolderPath } from '@/hooks/useFolderPath'
import type { FolderInfo } from '@/api/types/index'

/**
 * 文件导航管理 Hook
 * 负责面包屑导航、URL路径管理和文件夹路径构建
 */
export function useFileNavigation() {
  const router = useRouter()
  const { getFolderPath } = useFolderPath()

  const currentFolderId = ref<string | null>(null)
  const breadcrumbItems = ref<FolderInfo[]>([])

  const getRouteNames = () => ({
    base: 'folders',
    withPath: 'foldersWithPath',
  })

  const buildCompleteBreadcrumbPath = async (folder: FolderInfo): Promise<FolderInfo[]> => await getFolderPath(folder.id)

  const buildBreadcrumbFromPath = async (path: string) => {
    if (!path || path === '/') {
      breadcrumbItems.value = []
      currentFolderId.value = null
      return
    }

    try {
      const folderIds = path.split('/').filter((id) => id)

      try {
        const lastFolderId = folderIds[folderIds.length - 1]
        const newBreadcrumbs = await getFolderPath(lastFolderId)

        breadcrumbItems.value = newBreadcrumbs
        currentFolderId.value = lastFolderId
        return
      } catch (error) {
        console.warn('[useFileNavigation] Unified path hook failed, falling back to legacy method:', error)
      }

      const newBreadcrumbs: FolderInfo[] = []

      for (const folderId of folderIds) {
        try {
          const result = await folderApi.getFolderDetail(folderId)
          if (result.success) {
            newBreadcrumbs.push(result.data)
          } else {
            console.warn(`文件夹 ${folderId} 获取失败`)
            break
          }
        } catch {
          console.warn(`文件夹 ${folderId} 不存在或无权访问`)
          break
        }
      }

      breadcrumbItems.value = newBreadcrumbs
      currentFolderId.value = newBreadcrumbs.length > 0 ? newBreadcrumbs[newBreadcrumbs.length - 1].id : null
    } catch {
      breadcrumbItems.value = []
      currentFolderId.value = null
      updateURL()
    }
  }

  const updateURL = () => {
    const pathSegments = breadcrumbItems.value.map((item) => item.id)
    const routeNames = getRouteNames()

    if (pathSegments.length === 0) {
      router.replace({ name: routeNames.base })
    } else {
      const folderPath = pathSegments.join('/')
      router.replace({
        name: routeNames.withPath,
        params: { folderPath },
      })
    }
  }

  const navigateToFolder = async (folder: FolderInfo | null) => {
    if (folder === null) {
      currentFolderId.value = null
      breadcrumbItems.value = []
    } else {
      currentFolderId.value = folder.id

      const existingIndex = breadcrumbItems.value.findIndex((item) => item.id === folder.id)

      if (existingIndex !== -1) {
        breadcrumbItems.value = breadcrumbItems.value.slice(0, existingIndex + 1)
      } else {
        breadcrumbItems.value = await buildCompleteBreadcrumbPath(folder)
      }
    }
    updateURL() // 只更新URL，由路由监听器统一处理数据加载
  }

  const handleBreadcrumbClick = (folder: FolderInfo | null) => {
    if (folder === null || folder.id === null) {
      currentFolderId.value = null
      breadcrumbItems.value = []
    } else {
      const index = breadcrumbItems.value.findIndex((item) => item.id === folder.id)
      if (index !== -1) {
        currentFolderId.value = folder.id
        breadcrumbItems.value = breadcrumbItems.value.slice(0, index + 1)
      }
    }
    updateURL() // 只更新URL，由路由监听器统一处理数据加载
  }

  return {
    currentFolderId,
    breadcrumbItems,
    buildBreadcrumbFromPath,
    updateURL,
    navigateToFolder,
    handleBreadcrumbClick,
  }
}
