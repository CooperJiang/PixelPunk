/**
 * 多选选择器组件类型定义
 */

export interface SelectorOption {
  value: string
  label: string
  icon?: string
}

/**
 * MultiSelector 组件属性接口
 */
export interface MultiSelectorProps {
  modelValue?: string[]

  inputIdPrefix?: string

  disabled?: boolean

  options?: SelectorOption[]

  forcedValue?: string | null

  size?: 'sm' | 'md' | 'lg'

  showIcons?: boolean

  rounded?: 'none' | 'sm' | 'md' | 'lg'

  defaultIcon?: string

  forcedIcon?: string

  editable?: boolean

  addText?: string

  maxOptions?: number

  isGuest?: boolean

  validateFn?: (value: string, isGuest?: boolean) => string | null
}

/**
 * 组件事件类型定义
 */
export interface MultiSelectorEmits {
  'update:modelValue': [value: string[]]
  'add-option': [option: SelectorOption]
  'remove-option': [value: string]
  'validation-error': [errorMessage: string]
}
