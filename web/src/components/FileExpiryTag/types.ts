export interface ImageExpiryTagProps {
  expiresAt?: string | null
  storageDuration?: string
  isTimeLimited?: boolean
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right'
  mode?: 'countdown' | 'duration' | 'both'
  showIcon?: boolean
  className?: string
}

export interface ImageExpiryData {
  remaining: number
  status: 'expired' | 'expiring-soon' | 'active'
  remainingText: string
  durationText: string
  visible: boolean
}

export interface ExpiryStatusConfig {
  color: string
  bgColor: string
  borderColor: string
  icon: string
  text: string
}
