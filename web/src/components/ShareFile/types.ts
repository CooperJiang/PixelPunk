export interface ShareFileTag {
  id: string | number
  name: string
}

export interface ShareFileData {
  id: string | number
  display_name?: string
  original_name?: string
  format?: string
  size?: number
  size_formatted?: string
  width?: number
  height?: number
  resolution?: string
  full_thumb_url?: string
  thumb_url?: string
  full_url?: string
  url?: string
  created_at?: string
  tags?: ShareFileTag[]
}

export interface ShareFileProps {
  file: ShareFileData
  selectable?: boolean
  selected?: boolean
  shareKey: string
  accessToken?: string
}

export interface ShareFileEmits {
  click: [file: ShareFileData]
  select: [fileId: string | number]
}

/* 保持向后兼容的类型别名 */
export type ShareImageTag = ShareFileTag
export type ShareImageData = ShareFileData
export type ShareImageProps = ShareFileProps
export type ShareImageEmits = ShareFileEmits
