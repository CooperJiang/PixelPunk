import { computed, nextTick, ref } from 'vue'
import { onClickOutside, useElementBounding, useEventListener } from '@vueuse/core'
import type { DropdownOption, DropdownProps, DropdownEmits } from '../types'
import { Z_INDEX } from '@/constants/zIndex'

export function useDropdown(props: DropdownProps, emit: DropdownEmits, createLabelText?: string) {
  const isOpen = ref(false)
  const searchQuery = ref('')
  const highlightedIndex = ref(-1)
  const searchInput = ref<HTMLInputElement | null>(null)
  const dropdownRef = ref<HTMLElement | null>(null)
  const dropdownMenuRef = ref<HTMLElement | null>(null)
  const showClearButton = ref(false)
  const hoveredTag = ref<string | number | boolean | null>(null)
  const customOptions = ref<DropdownOption[]>([])

  const { left, width, bottom } = useElementBounding(dropdownRef)

  const getSmartZIndex = (): number => {
    let element = dropdownRef.value?.parentElement
    let maxZIndex = Z_INDEX.DROPDOWN

    while (element && element !== document.body) {
      const computedStyle = getComputedStyle(element)
      const zIndex = parseInt(computedStyle.zIndex)

      if (!isNaN(zIndex) && zIndex > 0) {
        if (zIndex >= 10000) {
          maxZIndex = zIndex + 100
          break
        } else if (zIndex > maxZIndex) {
          maxZIndex = zIndex + 10
        }
      }
      element = element.parentElement
    }

    return maxZIndex
  }

  const dropdownPosition = computed(() => ({
    top: bottom.value,
    left: left.value,
    width: width.value,
    zIndex: getSmartZIndex(),
  }))

  onClickOutside(dropdownRef, (event) => {
    if (!isOpen.value) {
      return
    }

    const target = event.target as HTMLElement
    const isClickInMenu = dropdownMenuRef.value && dropdownMenuRef.value.contains(target)

    if (!isClickInMenu) {
      close()
    }
  })

  const mergedOptions = computed(() => [...props.options, ...customOptions.value])

  const filteredOptions = computed(() => {
    const query = searchQuery.value.toLowerCase()

    if (!query) {
      return mergedOptions.value
    }

    const filtered = mergedOptions.value.filter((option) => option.label.toLowerCase().includes(query))

    const canCreate = props.allowCreate && query && !mergedOptions.value.some((option) => option.label.toLowerCase() === query)

    if (canCreate) {
      return [
        ...filtered,
        {
          value: `create-option-${query}`,
          label: `${createLabelText || props.createLabel} "${query}"`,
          isCreateOption: true,
          originalValue: query,
        },
      ]
    }

    return filtered
  })

  const selectedOption = computed(() => {
    if (props.multiple) {
      return null
    }

    return mergedOptions.value.find((option) => String(option.value) === String(props.modelValue))
  })

  const selectedOptions = computed(() => {
    if (!props.multiple || !Array.isArray(props.modelValue)) {
      return []
    }

    return mergedOptions.value.filter((option) => props.modelValue.some((value) => String(value) === String(option.value)))
  })

  const isSelected = (value: string | number | boolean) => {
    if (props.multiple && Array.isArray(props.modelValue)) {
      return props.modelValue.some((selectedValue) => String(selectedValue) === String(value))
    }

    return String(props.modelValue) === String(value)
  }

  const selectOption = (option: DropdownOption) => {
    if (props.disabled || option.disabled) {
      return
    }

    if (option.isCreateOption && props.allowCreate) {
      const newOption = {
        value: option.originalValue,
        label: option.originalValue,
      }

      customOptions.value.push(newOption)

      if (props.multiple) {
        const newValue = Array.isArray(props.modelValue) ? [...props.modelValue] : []
        newValue.push(newOption.value)
        emit('update:modelValue', newValue)
        emit('change', newValue)
      } else {
        emit('update:modelValue', newOption.value)
        emit('change', newOption.value)
        close()
      }

      searchQuery.value = ''

      emit('create', newOption)

      return
    }

    if (props.multiple) {
      const newValue = Array.isArray(props.modelValue) ? [...props.modelValue] : []
      const index = newValue.indexOf(option.value)

      if (index === -1) {
        newValue.push(option.value)
      } else {
        newValue.splice(index, 1)
      }

      emit('update:modelValue', newValue)
      emit('change', newValue)
    } else {
      emit('update:modelValue', option.value)
      emit('change', option.value)
      close()
    }
  }

  const removeSelected = (value: string | number | boolean) => {
    if (props.disabled) {
      return
    }

    if (props.multiple && Array.isArray(props.modelValue)) {
      const newValue = [...props.modelValue]
      const index = newValue.indexOf(value)

      if (index !== -1) {
        newValue.splice(index, 1)
        emit('update:modelValue', newValue)
        emit('change', newValue)
      }
    }
  }

  const clearAllSelected = () => {
    if (props.disabled) {
      return
    }

    if (props.multiple) {
      emit('update:modelValue', [])
      emit('change', [])
    } else {
      emit('update:modelValue', '')
      emit('change', '')
    }
  }

  const toggleDropdown = () => {
    if (props.disabled) {
      return
    }

    isOpen.value = !isOpen.value

    if (isOpen.value) {
      searchQuery.value = ''
      highlightedIndex.value = -1

      if (props.searchable) {
        nextTick(() => {
          searchInput.value?.focus()
        })
      }
    }
  }

  const close = () => {
    isOpen.value = false
  }

  const handleKeyDown = (e: KeyboardEvent) => {
    if (!isOpen.value) {
      return
    }

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault()
        if (highlightedIndex.value < filteredOptions.value.length - 1) {
          highlightedIndex.value++
        } else {
          highlightedIndex.value = 0
        }
        break
      case 'ArrowUp':
        e.preventDefault()
        if (highlightedIndex.value > 0) {
          highlightedIndex.value--
        } else {
          highlightedIndex.value = filteredOptions.value.length - 1
        }
        break
      case 'Enter':
        e.preventDefault()
        if (highlightedIndex.value >= 0) {
          selectOption(filteredOptions.value[highlightedIndex.value])
        }
        break
      case 'Escape':
        e.preventDefault()
        close()
        break
    }
  }

  useEventListener(window, 'keydown', handleKeyDown)

  return {
    isOpen,
    searchQuery,
    highlightedIndex,
    searchInput,
    dropdownRef,
    dropdownMenuRef,
    showClearButton,
    hoveredTag,

    dropdownPosition,
    filteredOptions,
    selectedOption,
    selectedOptions,

    isSelected,
    selectOption,
    removeSelected,
    clearAllSelected,
    toggleDropdown,
    close,
  }
}
