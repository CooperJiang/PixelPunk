import { folderApi } from '@/api'
import type { FolderInfo } from '@/api/types/index'

/**
 * 统一的文件夹路径获取Hook
 * 提供优化的路径查询功能，包含降级机制
 */
export function useFolderPath() {
  const getFolderPath = async (folderId: string): Promise<FolderInfo[]> => {
    try {
      const pathChainResult = await folderApi.getFolderPathChain(folderId)
      if (pathChainResult.success) {
        const breadcrumbs = pathChainResult.data.path_chain.map(
          (chain) =>
            ({
              id: chain.id,
              name: chain.name,
              parent_id: chain.parent_id,
              created_at: '',
            }) as FolderInfo
        )

        return breadcrumbs
      }
    } catch (error) {
      console.warn('[useFolderPath] Path API failed, using legacy method:', error)
    }

    const result = await getFolderPathLegacy(folderId)
    return result
  }

  const getBatchFolderPaths = async (folderIds: string[]): Promise<{ [folderId: string]: FolderInfo[] }> => {
    try {
      const batchResult = await folderApi.getBatchFolderPathChains(folderIds)
      if (batchResult.success) {
        const result: { [folderId: string]: FolderInfo[] } = {}

        batchResult.data.forEach((pathChain) => {
          result[pathChain.folder_id] = pathChain.path_chain.map(
            (chain) =>
              ({
                id: chain.id,
                name: chain.name,
                parent_id: chain.parent_id,
                created_at: '',
              }) as FolderInfo
          )
        })

        return result
      }
    } catch (error) {
      console.warn('[useFolderPath] Batch path chain API failed, falling back to individual requests:', error)
    }

    const result: { [folderId: string]: FolderInfo[] } = {}
    for (const folderId of folderIds) {
      try {
        result[folderId] = await getFolderPath(folderId)
      } catch (error) {
        console.error(`[useFolderPath] Failed to get path for folder ${folderId}:`, error)
      }
    }
    return result
  }

  const getFolderPathString = async (folderId: string): Promise<string> => {
    try {
      const pathChainResult = await folderApi.getFolderPathChain(folderId)
      if (pathChainResult.success) {
        return pathChainResult.data.full_path
      }
    } catch (error) {
      console.warn('[useFolderPath] Path string API failed, using fallback:', error)
    }

    const pathArray = await getFolderPath(folderId)
    const result = pathArray.map((folder) => folder.name).join(' / ')
    return result
  }

  const getFolderPathLegacy = async (folderId: string): Promise<FolderInfo[]> => {
    const path: FolderInfo[] = []
    let currentId = folderId

    while (currentId) {
      try {
        const result = await folderApi.getFolderDetail(currentId)
        if (result.success) {
          path.unshift(result.data) // 添加到数组开头
          currentId = result.data.parent_id || ''
        } else {
          break
        }
      } catch (error) {
        console.error('[useFolderPath] Legacy method failed:', error)
        break
      }
    }

    return path
  }

  return {
    getFolderPath,
    getBatchFolderPaths,
    getFolderPathString,
  }
}
