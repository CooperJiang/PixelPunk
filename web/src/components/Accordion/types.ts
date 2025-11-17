export interface AccordionItem {
  id: string | number
  title: string
  content?: string
  icon?: string
  disabled?: boolean
  data?: unknown
}

export interface AccordionProps {
  items: AccordionItem[]
  modelValue?: string | number | null
  multiple?: boolean
  expandedItems?: (string | number)[]
  bordered?: boolean
  gap?: string
  animationDuration?: number
  card?: boolean
  size?: 'small' | 'medium' | 'large'
}

export interface AccordionEmits {
  (e: 'update:modelValue', value: string | number | null): void
  (e: 'update:expandedItems', value: (string | number)[]): void
  (e: 'item-click', item: AccordionItem, expanded: boolean): void
  (e: 'item-expand', item: AccordionItem): void
  (e: 'item-collapse', item: AccordionItem): void
}
