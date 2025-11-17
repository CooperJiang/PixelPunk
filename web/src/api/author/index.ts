import { get } from '@/utils/network/http'
import type { ApiResult } from '@/utils/network/http-types'
import type { AuthorFolderContent, AuthorHomepage } from '../types'

export function getAuthorHomepage(authorId: string): Promise<ApiResult<AuthorHomepage>> {
  return get<AuthorHomepage>(`/authors/${authorId}`)
}

export function getAuthorFolder(
  authorId: string,
  folderId: string,
  params?: {
    page?: number
    size?: number
  }
): Promise<ApiResult<AuthorFolderContent>> {
  return get<AuthorFolderContent>(`/authors/${authorId}/folders/${folderId}`, params)
}
