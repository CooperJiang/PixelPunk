/* Drawer component props */
export interface DrawerProps {
  modelValue: boolean
  title?: string
  width?: string | number
  position?: 'left' | 'right'
  maskClosable?: boolean
  escClosable?: boolean
  showHeader?: boolean
  showKeyboardTip?: boolean
}

export interface DrawerEmits {
  'update:modelValue': [value: boolean]
  close: []
}
