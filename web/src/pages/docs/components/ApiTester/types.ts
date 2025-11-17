export interface ApiTestParams {
  folderId: string
  filePath: string
  accessLevel: 'public' | 'private' | 'protected'
  optimize: boolean
}

export interface ApiResponse {
  status: number
  statusText: string
  data: any
  duration: number
}

export interface UploadedImage {
  id: string
  url: string
  thumb_url: string
  original_name: string
  display_name?: string
  width: number
  height: number
  size: number
  format: string
  access_level: string
}

export interface ApiTestState {
  apiKey: string
  selectedFiles: File[]
  params: ApiTestParams
  response: ApiResponse | null
  isLoading: boolean
  responseCopied: boolean
  showValidation: boolean
  isDragOver: boolean
}

export interface ErrorData {
  oversized_files?: string[]
  unsupported_files?: string[]
  invalid_files?: string[]
  upload_errors?: string[]
  size_limit?: string
}
