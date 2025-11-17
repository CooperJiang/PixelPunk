import { del, get, post } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { CreateFolderRequest, FolderInfo, FolderPathChainResponse, FolderTreeNode, UpdateFolderRequest } from '../types'

export function createFolder(data: CreateFolderRequest): Promise<ApiResult<FolderInfo>> {
  return post<FolderInfo>('/folders/create', data)
}

export function getFolderList(parentId?: string): Promise<ApiResult<FolderInfo[]>> {
  return get<FolderInfo[]>('/folders/list', parentId ? { parent_id: parentId } : undefined)
}

export function getFolderListWithImages(parentId?: string): Promise<ApiResult<FolderInfo[]>> {
  return get<FolderInfo[]>('/folders/contents', parentId ? { parent_id: parentId } : undefined)
}

export function getFolderDetail(folderId: string): Promise<ApiResult<FolderInfo>> {
  return get<FolderInfo>(`/folders/${folderId}`)
}

export function updateFolder(data: UpdateFolderRequest): Promise<ApiResult<FolderInfo>> {
  return post<FolderInfo>(`/folders/update`, data)
}

export function deleteFolder(folderId: string): Promise<ApiResult<{ id: string }>> {
  return del<{ id: string }>(`/folders/${folderId}`)
}

export function visibilityFolder(folderId: string): Promise<ApiResult<{ id: string }>> {
  return post<{ id: string }>(`/folders/${folderId}/toggle-access-level`)
}

export function reorderFolders(data: { parent_id?: string; folder_ids: string[] }): Promise<ApiResult<void>> {
  return post<void>('/folders/reorder', data)
}

export function getFolderTree(): Promise<ApiResult<FolderTreeNode[]>> {
  return get<FolderTreeNode[]>('/folders/tree')
}

export function getFolderPathChain(folderId: string): Promise<ApiResult<FolderPathChainResponse>> {
  return get<FolderPathChainResponse>(`/folders/${folderId}/path-chain`)
}

export function getBatchFolderPathChains(folderIds: string[]): Promise<ApiResult<FolderPathChainResponse[]>> {
  return post<FolderPathChainResponse[]>('/folders/batch-path-chains', { folder_ids: folderIds })
}

export function searchFolders(keyword: string): Promise<ApiResult<FolderInfo[]>> {
  return get<FolderInfo[]>('/folders/search', { keyword })
}

export function moveFolders(folderIds: string[], newParentId?: string): Promise<ApiResult<void>> {
  return post<void>('/folders/move', {
    folder_ids: folderIds,
    new_parent_id: newParentId || '',
  })
}

export function moveFolder(folderId: string, newParentId?: string): Promise<ApiResult<void>> {
  return moveFolders([folderId], newParentId)
}

export default {
  createFolder,
  getFolderList,
  getFolderDetail,
  updateFolder,
  deleteFolder,
  getFolderListWithImages,
  visibilityFolder,
  reorderFolders,
  getFolderTree,
  getFolderPathChain,
  getBatchFolderPathChains,
  searchFolders,
  moveFolders,
  moveFolder,
}
