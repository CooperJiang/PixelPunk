import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface BreadcrumbItem {
  id?: string
  name: string
  path?: string
}

export const useBreadcrumbStore = defineStore('breadcrumb', () => {
  const items = ref<BreadcrumbItem[]>([])

  function setItems(newItems: BreadcrumbItem[]) {
    items.value = newItems
  }

  function clear() {
    items.value = []
  }

  return {
    items,
    setItems,
    clear,
  }
})
