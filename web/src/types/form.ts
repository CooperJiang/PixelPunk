/**
 * 表单相关类型定义
 * 包含表单字段、验证、错误等类型

/* 表单字段类型 */
export type FormFieldType =
  | 'text'
  | 'textarea'
  | 'password'
  | 'email'
  | 'number'
  | 'select'
  | 'radio'
  | 'checkbox'
  | 'switch'
  | 'date'
  | 'file'

/* 表单字段定义 */
export interface FormField {
  name: string
  label: string
  type: FormFieldType
  required?: boolean
  placeholder?: string
  default_value?: any
  options?: Array<{ label: string; value: any }>
  validation?: any[]
  props?: Record<string, any>
}

/* 表单数据 */
export type FormData = Record<string, any>

/* 表单错误 */
export type FormErrors = Record<string, string[]>
