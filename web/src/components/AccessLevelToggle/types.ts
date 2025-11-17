export type AccessLevel = 'public' | 'private' | 'protected'
export type ToggleType = 'image' | 'folder'
export type ToggleSize = 'small' | 'medium' | 'large'

export interface AccessLevelToggleProps {
  accessLevel: AccessLevel
  type?: ToggleType
  size?: ToggleSize
  disabled?: boolean
}

export interface AccessLevelToggleEmits {
  (e: 'toggle'): void
}
