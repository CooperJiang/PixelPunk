import { ref, watch } from 'vue'

const STORAGE_KEY_COMMON = 'dashboard_common_actions_order'
const STORAGE_KEY_SETTINGS = 'dashboard_settings_order'

/**
 * Dashboard 快速操作拖动排序 Hook (纯前端本地存储)
 */
export function useQuickActionsDrag() {
  const loadOrder = <T extends { id: number }>(storageKey: string, defaultItems: T[]): T[] => {
    try {
      const savedOrder = localStorage.getItem(storageKey)
      if (!savedOrder) {
        return defaultItems
      }

      const orderIds = JSON.parse(savedOrder) as number[]

      const orderedItems: T[] = []
      const itemsMap = new Map(defaultItems.map((item) => [item.id, item]))

      orderIds.forEach((id) => {
        const item = itemsMap.get(id)
        if (item) {
          orderedItems.push(item)
          itemsMap.delete(id)
        }
      })

      itemsMap.forEach((item) => {
        orderedItems.push(item)
      })

      return orderedItems
    } catch (error) {
      console.error('Failed to load order:', error)
      return defaultItems
    }
  }

  const saveOrder = <T extends { id: number }>(storageKey: string, items: T[]) => {
    try {
      const orderIds = items.map((item) => item.id)
      localStorage.setItem(storageKey, JSON.stringify(orderIds))
    } catch (error) {
      console.error('Failed to save order:', error)
    }
  }

  const createDraggableActions = <T extends { id: number }>(storageKey: string, defaultItems: T[]) => {
    const items = ref<T[]>(loadOrder(storageKey, defaultItems))

    watch(
      items,
      (newItems) => {
        saveOrder(storageKey, newItems)
      },
      { deep: true }
    )

    const handleDragStart = () => {}

    const handleDragEnd = () => {}

    return {
      items,
      handleDragStart,
      handleDragEnd,
    }
  }

  return {
    createDraggableActions,
    STORAGE_KEY_COMMON,
    STORAGE_KEY_SETTINGS,
  }
}
