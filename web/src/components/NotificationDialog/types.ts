export interface NotificationDialogProps {
  title: string
  width?: string | number
  showCloseButton?: boolean
}

export interface NotificationDialogExpose {
  show: () => void
  hide: () => void
}
