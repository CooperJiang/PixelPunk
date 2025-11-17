import { onBeforeUnmount, type Ref } from 'vue'

/**
 * 点击元素外部的钩子函数
 * @param elementRefs 需要排除的元素引用数组
 * @param callback 点击外部时的回调函数
 */
export function useClickOutside(elementRefs: Ref<HTMLElement | null>[] | HTMLElement[], callback: () => void) {
  const clickHandler = (e: MouseEvent) => {
    const target = e.target as HTMLElement

    const isOutside = !elementRefs.some((elementRef) => {
      const el = 'value' in elementRef ? elementRef.value : elementRef
      return el && el.contains(target)
    })

    if (isOutside) {
      callback()
    }
  }

  const addClickOutside = () => {
    document.addEventListener('mousedown', clickHandler)
  }

  const removeClickOutside = () => {
    document.removeEventListener('mousedown', clickHandler)
  }

  onBeforeUnmount(removeClickOutside)

  return {
    addClickOutside,
    removeClickOutside,
  }
}
