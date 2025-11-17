export interface CommunityDialogProps {
  type?: 'wechat' | 'qq'
  qrImage?: string
  contactInfo?: string
}

export interface CommunityDialogExpose {
  show: () => void
  hide: () => void
}
