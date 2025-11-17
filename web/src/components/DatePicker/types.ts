export type DatePickerType = 'date' | 'daterange'
export type DatePickerSize = 'small' | 'medium' | 'large'

export interface DatePickerProps {
  modelValue?: Date | [Date | null, Date | null]
  type?: DatePickerType
  size?: DatePickerSize
  placeholder?: string
  startPlaceholder?: string
  endPlaceholder?: string
  rangeSeparator?: string
  disabled?: boolean
  format?: string
  disabledDate?: (date: Date) => boolean
}

export interface DatePickerEmits {
  (e: 'update:modelValue', value: Date | [Date | null, Date | null] | null): void
  (e: 'change', value: Date | [Date | null, Date | null]): void
}

export interface Day {
  date: number
  fullDate: Date
  isPrevMonth: boolean
  isNextMonth: boolean
  isToday: boolean
  isDisabled: boolean
  isSelected: boolean
  isRangeStart: boolean
  isRangeEnd: boolean
  isInRange: boolean
}
