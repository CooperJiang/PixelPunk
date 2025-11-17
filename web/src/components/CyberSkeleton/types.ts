export interface CyberSkeletonProps {
  type: 'card' | 'table' | 'list'

  size?: 'small' | 'normal' | 'large'

  count?: number

  loading?: boolean

  animation?: 'wave' | 'pulse' | 'none'

  simple?: boolean
}
