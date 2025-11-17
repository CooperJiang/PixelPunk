export type UserAvatarSize = 'sm' | 'md' | 'lg'

export interface UserAvatarProps {
  avatarUrl?: string
  username?: string
  size?: UserAvatarSize
  showBorder?: boolean
  showGlow?: boolean
  customStyles?: Record<string, string | number>
}
