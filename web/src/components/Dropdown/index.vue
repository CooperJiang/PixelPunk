<template>
  <div
    class="cyber-dropdown"
    :class="[{ 'cyber-dropdown--open': isOpen, 'cyber-dropdown--disabled': disabled }]"
    :style="wrapperStyle"
    ref="dropdownRef"
  >
    <div class="cyber-dropdown__trigger" @click.stop="handleTriggerClick">
      <div class="cyber-dropdown__content">
        <!-- 可搜索模式且处于输入状态时显示输入框 -->
        <input
          v-if="searchable && (isOpen || searchQuery)"
          ref="triggerInput"
          v-model="searchQuery"
          class="cyber-dropdown__trigger-input"
          :placeholder="placeholder"
          @click.stop
          @keydown.enter.prevent="handleEnterKey"
          @keydown.esc.prevent="handleEscKey"
          @keydown.down.prevent="handleArrowDown"
          @keydown.up.prevent="handleArrowUp"
          @blur="handleInputBlur"
        />
        <!-- 非搜索模式或未输入时显示选中值 -->
        <template v-else>
          <template v-if="multiple && selectedOptions.length">
            <div class="cyber-dropdown__tags">
              <div
                v-for="option in selectedOptions"
                :key="option.value"
                class="cyber-dropdown__tag"
                @mouseenter="hoveredTag = option.value"
                @mouseleave="hoveredTag = null"
              >
                <span
                  v-if="isColorMode && option.color"
                  class="cyber-dropdown__color-block"
                  :style="{ backgroundColor: option.color }"
                ></span>
                {{ option.label }}
                <span
                  class="cyber-dropdown__tag-remove"
                  :class="{ 'cyber-dropdown__tag-remove--visible': hoveredTag === option.value }"
                  @click.stop="removeSelected(option.value)"
                  >×</span
                >
              </div>
            </div>
          </template>
          <template v-else-if="!multiple && selectedOption">
            <div class="cyber-dropdown__selected">
              <span v-if="$slots['selected-icon'] || (isColorMode && selectedOption.color)" class="cyber-dropdown__selected-icon">
                <slot name="selected-icon" :option="selectedOption" :is-permanent="selectedOption.value === 'permanent'">
                  <span
                    v-if="isColorMode && selectedOption.color"
                    class="cyber-dropdown__color-block"
                    :style="{ backgroundColor: selectedOption.color }"
                  ></span>
                </slot>
              </span>
              <span class="cyber-dropdown__selected-label">{{ selectedOption.label }}</span>
            </div>
          </template>
          <template v-else>
            <span class="cyber-dropdown__placeholder">{{ placeholder }}</span>
          </template>
        </template>
      </div>
      <div class="cyber-dropdown__suffix" @mouseenter="showClearButton = true" @mouseleave="showClearButton = false">
        <span
          v-if="(multiple ? selectedOptions.length > 0 : selectedOption) && showClearButton && clearable"
          class="cyber-dropdown__clear"
          @click.stop="clearAllSelected"
          :title="$t('components.dropdown.clearAll')"
        >
          <i class="fas fa-times-circle"></i>
        </span>
        <i v-else class="fas fa-chevron-down cyber-dropdown__arrow" :class="{ 'cyber-dropdown__arrow--open': isOpen }"></i>
      </div>
    </div>

    <Teleport to="body">
      <transition name="cyber-dropdown">
        <div
          v-if="isOpen"
          class="cyber-dropdown__menu"
          :style="{
            maxHeight: maxHeight,
            position: 'fixed',
            top: dropdownPosition.top + 'px',
            left: dropdownPosition.left + 'px',
            width: dropdownPosition.width + 'px',
            zIndex: dropdownPosition.zIndex,
          }"
          ref="dropdownMenuRef"
          @click.stop
        >
          <div class="cyber-dropdown__options" @scroll="handleScroll">
            <div v-if="loading && filteredOptions.length === 0" class="cyber-dropdown__loading">
              <i class="fas fa-spinner fa-spin"></i>
              {{ $t('components.dropdown.loading') }}
            </div>

            <template v-else-if="filteredOptions.length">
              <template v-for="option in filteredOptions" :key="option.value">
                <slot
                  v-if="$slots.option"
                  name="option"
                  :option="option"
                  :is-selected="isSelected(option.value)"
                  :is-color-mode="isColorMode"
                  :on-select="() => selectOption(option)"
                  :is-highlighted="highlightedIndex === filteredOptions.indexOf(option)"
                  :is-create-option="option.isCreateOption"
                ></slot>

                <div
                  v-else
                  :key="`default-${option.value}`"
                  class="cyber-dropdown__option"
                  :class="{
                    'cyber-dropdown__option--selected': isSelected(option.value),
                    'cyber-dropdown__option--highlighted': highlightedIndex === filteredOptions.indexOf(option),
                    'cyber-dropdown__option--color': isColorMode && option.color,
                    'cyber-dropdown__option--create': option.isCreateOption,
                  }"
                  @mousedown.prevent="() => {
                    if (inputBlurTimeout) {
                      clearTimeout(inputBlurTimeout)
                      inputBlurTimeout = null
                    }
                  }"
                  @click.stop="selectOption(option)"
                  @mouseover="highlightedIndex = filteredOptions.indexOf(option)"
                >
                  <template v-if="multiple">
                    <div class="cyber-dropdown__checkbox-wrapper" @click.stop="selectOption(option)">
                      <span class="cyber-dropdown__checkbox">
                        <i v-if="isSelected(option.value)" class="fas fa-check"></i>
                      </span>
                    </div>
                  </template>

                  <span
                    v-if="$slots['option-icon'] || (isColorMode && option.color) || option.isCreateOption"
                    class="cyber-dropdown__option-icon"
                  >
                    <slot
                      name="option-icon"
                      :option="option"
                      :is-selected="isSelected(option.value)"
                      :is-permanent="option.value === 'permanent'"
                    >
                      <span
                        v-if="isColorMode && option.color"
                        class="cyber-dropdown__color-block"
                        :style="{ backgroundColor: option.color }"
                      ></span>
                      <i v-else-if="option.isCreateOption" class="fas fa-plus-circle cyber-dropdown__create-icon"></i>
                    </slot>
                  </span>

                  <span class="cyber-dropdown__option-label">{{ option.label }}</span>
                </div>
              </template>

              <div v-if="pagination && loadingMore" class="cyber-dropdown__loading-more">
                <i class="fas fa-spinner fa-spin"></i>
                {{ $t('components.dropdown.loadingMore') }}
              </div>

              <div v-else-if="pagination && !hasMore" class="cyber-dropdown__no-more">{{ $t('components.dropdown.noMore') }}</div>
            </template>

            <div v-else class="cyber-dropdown__empty">
              {{ $t('components.dropdown.empty') }}
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, watch, nextTick } from 'vue'
  import { onClickOutside, useEventListener } from '@vueuse/core'
  import type { DropdownProps, DropdownEmits, DropdownOption } from './types'
  import { Z_INDEX } from '@/constants/zIndex'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'CyberDropdown',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<DropdownProps>(), {
    modelValue: '',
    options: () => [],
    placeholder: undefined,
    disabled: false,
    multiple: false,
    searchable: false,
    maxHeight: '250px',
    borderColor: 'var(--color-brand-500)',
    textColor: 'var(--color-content-default)',
    bgColor: 'var(--color-background-800)',
    isColorMode: false,
    clearable: true,
    allowCreate: false,
    createLabel: undefined,
    pagination: false,
    loading: false,
    loadingMore: false,
    hasMore: true,
    pageSize: 20,
    width: '100%',
    height: undefined,
  })

  const emit = defineEmits<DropdownEmits>()

  const isOpen = ref(false)
  const searchQuery = ref('')
  const highlightedIndex = ref(-1)
  const searchInput = ref<HTMLInputElement | null>(null)
  const triggerInput = ref<HTMLInputElement | null>(null)
  const dropdownRef = ref<HTMLElement | null>(null)
  const dropdownMenuRef = ref<HTMLElement | null>(null)
  const showClearButton = ref(false)
  const hoveredTag = ref<any>(null)
  const customOptions = ref<DropdownOption[]>([])
  const shouldOpenUpward = ref(false)
  const menuTopPosition = ref(0)
  const menuLeftPosition = ref(0)
  const menuWidth = ref(0)
  const instanceId = Math.random().toString(36).slice(2)
  const inputBlurTimeout = ref<number | null>(null)

  const placeholder = computed(() => props.placeholder || $t('components.dropdown.placeholder'))
  // const createLabel = computed(() => props.createLabel || $t('components.dropdown.create'))

  const wrapperStyle = computed(() => {
    const style: Record<string, string> = {}
    if (props.width) {
      style.width = props.width
    }
    if (props.height) {
      style.height = props.height
      style.minHeight = props.height
    }
    return style
  })

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

  const dropdownPosition = computed(() => {
    return {
      top: menuTopPosition.value,
      left: menuLeftPosition.value,
      width: menuWidth.value,
      zIndex: getSmartZIndex(),
    }
  })

  onClickOutside(dropdownRef, (event) => {
    if (!isOpen.value) return

    const target = event.target as HTMLElement
    const isClickInMenu = dropdownMenuRef.value && dropdownMenuRef.value.contains(target)

    if (!isClickInMenu) {
      close()
    }
  })

  const mergedOptions = computed(() => {
    return [...props.options, ...customOptions.value]
  })

  const filteredOptions = computed(() => {
    const query = searchQuery.value.toLowerCase()

    if (!query) return mergedOptions.value

    const filtered = mergedOptions.value.filter((option) => option.label.toLowerCase().includes(query))

    const canCreate = props.allowCreate && query && !mergedOptions.value.some((option) => option.label.toLowerCase() === query)

    if (canCreate) {
      return [
        ...filtered,
        {
          value: `create-option-${query}`,
          label: `${props.createLabel} "${query}"`,
          isCreateOption: true,
          originalValue: query,
        },
      ]
    }

    return filtered
  })

  const selectedOption = computed(() => {
    if (props.multiple) return null

    return mergedOptions.value.find(
      (option) => option.value === props.modelValue || String(option.value) === String(props.modelValue)
    )
  })

  const selectedOptions = computed(() => {
    if (!props.multiple || !Array.isArray(props.modelValue)) return []

    return mergedOptions.value.filter((option) =>
      props.modelValue.some((value) => value === option.value || String(value) === String(option.value))
    )
  })

  const isSelected = (value: any) => {
    if (props.multiple && Array.isArray(props.modelValue)) {
      return props.modelValue.some((selectedValue) => selectedValue === value || String(selectedValue) === String(value))
    }

    return props.modelValue === value || String(props.modelValue) === String(value)
  }

  const selectOption = (option: DropdownOption) => {
    if (props.disabled || option.disabled) return

    // 取消 blur 延迟定时器
    if (inputBlurTimeout.value) {
      clearTimeout(inputBlurTimeout.value)
      inputBlurTimeout.value = null
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

  const removeSelected = (value: any) => {
    if (props.disabled) return

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
    if (props.disabled) return

    if (props.multiple) {
      emit('update:modelValue', [])
      emit('change', [])
    } else {
      emit('update:modelValue', '')
      emit('change', '')
    }
  }

  const calculateMenuPosition = () => {
    if (!dropdownRef.value) return

    const rect = dropdownRef.value.getBoundingClientRect()
    const triggerTop = rect.top
    const triggerBottom = rect.bottom
    const triggerLeft = rect.left
    const triggerWidth = rect.width

    menuLeftPosition.value = triggerLeft
    menuWidth.value = triggerWidth

    const maxHeightValue = parseInt(props.maxHeight) || 250
    const estimatedMenuHeight = maxHeightValue + 80 // 加上搜索框、padding、边框等额外空间
    const viewportHeight = window.innerHeight
    const spaceBelow = viewportHeight - triggerBottom
    const spaceAbove = triggerTop

    shouldOpenUpward.value = spaceBelow < estimatedMenuHeight && spaceAbove > spaceBelow

    if (shouldOpenUpward.value) {
      menuTopPosition.value = triggerTop - estimatedMenuHeight
    } else {
      menuTopPosition.value = triggerBottom + 2
    }

    if (shouldOpenUpward.value) {
      nextTick(() => {
        if (dropdownMenuRef.value && dropdownRef.value) {
          const actualMenuHeight = dropdownMenuRef.value.offsetHeight
          const currentRect = dropdownRef.value.getBoundingClientRect()
          menuTopPosition.value = currentRect.top - actualMenuHeight - 8
        }
      })
    }
  }

  const handleTriggerClick = () => {
    if (props.disabled) return

    if (!isOpen.value) {
      calculateMenuPosition()
      isOpen.value = true
      window.dispatchEvent(new CustomEvent('cyber-dropdown-open', { detail: { id: instanceId } }))

      nextTick(() => {
        calculateMenuPosition()
        if (props.searchable && triggerInput.value) {
          triggerInput.value.focus()
        }
      })

      highlightedIndex.value = -1
    }
  }

  const toggleDropdown = () => {
    if (props.disabled) return

    if (!isOpen.value) {
      calculateMenuPosition()
    }

    isOpen.value = !isOpen.value

    if (isOpen.value) {
      window.dispatchEvent(new CustomEvent('cyber-dropdown-open', { detail: { id: instanceId } }))

      nextTick(() => {
        calculateMenuPosition()
      })

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
    searchQuery.value = ''
  }

  const handleEnterKey = () => {
    if (highlightedIndex.value >= 0 && filteredOptions.value[highlightedIndex.value]) {
      selectOption(filteredOptions.value[highlightedIndex.value])
    } else if (props.allowCreate && searchQuery.value.trim()) {
      // 创建新选项
      const newOption = {
        value: searchQuery.value.trim(),
        label: searchQuery.value.trim(),
        isCreateOption: true,
        originalValue: searchQuery.value.trim(),
      }
      selectOption(newOption)
    }
  }

  const handleEscKey = () => {
    close()
    triggerInput.value?.blur()
  }

  const handleArrowDown = () => {
    if (!isOpen.value) {
      handleTriggerClick()
    } else if (highlightedIndex.value < filteredOptions.value.length - 1) {
      highlightedIndex.value++
    } else {
      highlightedIndex.value = 0
    }
  }

  const handleArrowUp = () => {
    if (!isOpen.value) {
      handleTriggerClick()
    } else if (highlightedIndex.value > 0) {
      highlightedIndex.value--
    } else {
      highlightedIndex.value = filteredOptions.value.length - 1
    }
  }

  const handleInputBlur = () => {
    // 延迟关闭，以便点击选项时不会因为失焦而关闭
    inputBlurTimeout.value = window.setTimeout(() => {
      close()
    }, 200)
  }

  const handleScroll = (event: Event) => {
    if (!props.pagination || props.loadingMore || !props.hasMore) return

    const target = event.target as HTMLElement
    const scrollTop = target.scrollTop
    const scrollHeight = target.scrollHeight
    const clientHeight = target.clientHeight

    if (scrollHeight - scrollTop - clientHeight < 50) {
      emit('load-more')
    }
  }

  watch(searchQuery, (newQuery, oldQuery) => {
    if (props.pagination && newQuery !== oldQuery) {
      setTimeout(() => {
        if (searchQuery.value === newQuery) {
          emit('search', newQuery)
        }
      }, 300)
    }
  })

  const handleKeyDown = (e: KeyboardEvent) => {
    if (!isOpen.value) return

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

  useEventListener(window, 'cyber-dropdown-open', (e: Event) => {
    const detail = (e as CustomEvent).detail as { id?: string } | undefined
    if (detail?.id !== instanceId && isOpen.value) {
      close()
    }
  })

  useEventListener(window, 'resize', () => {
    if (isOpen.value) calculateMenuPosition()
  })
  useEventListener(
    window as any,
    'scroll',
    () => {
      if (isOpen.value) calculateMenuPosition()
    },
    { passive: true }
  )
  useEventListener(
    document,
    'scroll',
    () => {
      if (isOpen.value) calculateMenuPosition()
    },
    { capture: true, passive: true }
  )
</script>

<style scoped>
  .cyber-dropdown {
    position: relative;
    width: 100%;
    font-family: inherit;
    min-height: 32px;
  }

  .cyber-dropdown__trigger {
    box-sizing: border-box;
    display: flex;
    min-height: 32px;
    cursor: pointer;
    align-items: center;
    justify-content: space-between;
    border-radius: var(--radius-sm);
    border: 2px solid var(--color-border-default);
    background: var(--color-background-700);
    padding: 4px 12px;
    font-size: 0.875rem;
    color: var(--color-content);
    transition:
      border-color 0.25s ease,
      box-shadow 0.25s ease,
      transform 0.25s ease;
    box-shadow:
      2px 2.5px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 0 rgba(var(--color-brand-500-rgb), 0);
  }

  .cyber-dropdown__trigger:hover:not(.cyber-dropdown--disabled .cyber-dropdown__trigger) {
    border-color: var(--color-border-strong);

    box-shadow:
      2.5px 3px 0 rgba(var(--color-brand-500-rgb), 0.3),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.08);
    transform: translate(-0.3px, -0.3px);
  }

  .cyber-dropdown--open .cyber-dropdown__trigger {
    border-color: var(--color-brand-500);

    box-shadow:
      3.5px 4px 0 rgba(var(--color-brand-500-rgb), 0.4),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.2);
    transform: translate(-0.5px, -0.5px);
  }

  .cyber-dropdown__trigger:active:not(.cyber-dropdown--disabled .cyber-dropdown__trigger) {
    transform: translate(0.5px, 0.5px);
    box-shadow:
      1.5px 2px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-dropdown__content {
    flex: 1;
    margin-right: 8px;
    display: flex;
    align-items: center;
    min-height: 0;
  }

  .cyber-dropdown__trigger-input {
    flex: 1;
    width: 100%;
    border: none;
    background: transparent;
    outline: none;
    color: var(--color-content);
    font-size: 0.875rem;
    font-family: inherit;
    padding: 0;
    margin: 0;
  }

  .cyber-dropdown__trigger-input::placeholder {
    color: var(--color-placeholder);
  }

  .cyber-dropdown__suffix {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    height: 18px;
  }

  .cyber-dropdown__arrow {
    @apply text-xs text-content-muted transition-transform duration-200;
  }

  .cyber-dropdown__arrow--open {
    transform: rotate(180deg);
  }

  .cyber-dropdown__clear {
    @apply flex cursor-pointer items-center justify-center text-sm text-content-muted transition-colors duration-200 hover:text-error-500;
  }

  .cyber-dropdown__placeholder {
    @apply text-placeholder;
  }

  .cyber-dropdown__selected {
    display: flex;
    align-items: center;
    width: 100%;
  }

  .cyber-dropdown__selected-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    margin-right: 8px;
    min-width: 16px;
    flex-shrink: 0;
  }

  .cyber-dropdown__selected-label {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .cyber-dropdown__tags {
    display: flex;
    flex-wrap: wrap;
    width: 100%;
    align-items: flex-start;
    gap: 4px;
    row-gap: 4px;
  }

  .cyber-dropdown__tag {
    @apply relative inline-flex h-[18px] min-w-max max-w-[120px] items-center overflow-hidden text-ellipsis whitespace-nowrap rounded-sm border border-badge-primary-border bg-badge-primary-bg px-1.5 py-0.5 text-xs text-badge-primary-text transition-all duration-200;
    flex-shrink: 0;
  }

  .cyber-dropdown__tag:hover {
    @apply border-error-500;

    box-shadow: 0 0 0 1px rgba(var(--color-error-rgb), 0.2);
  }

  .cyber-dropdown__tag-remove {
    @apply ml-1 cursor-pointer text-xs font-bold text-content opacity-0 transition-opacity duration-200;
  }

  .cyber-dropdown__tag-remove--visible {
    @apply opacity-70;
  }

  .cyber-dropdown__tag-remove--visible:hover {
    @apply text-error-500 opacity-100;
  }

  .cyber-dropdown__color-block {
    @apply relative mr-1.5 inline-block h-4 w-4 overflow-hidden rounded border border-subtle;
  }

  .cyber-dropdown__color-block::after {
    @apply absolute inset-0 content-[''];
    background: linear-gradient(to bottom right, transparent, var(--color-overlay-light));
  }

  .cyber-dropdown__menu {
    @apply mt-0.5 flex flex-col overflow-hidden rounded-md bg-background-900 backdrop-blur-md;
    background: rgba(var(--color-background-900-rgb), 0.98);
    border: 2px solid var(--color-border-default);

    box-shadow:
      3px 4px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 6px 20px var(--color-overlay-heavy);
    transition:
      box-shadow 0.25s ease,
      transform 0.25s ease;
  }

  .cyber-dropdown__options {
    @apply min-h-0 flex-grow overflow-y-auto py-1;
    scrollbar-width: thin;
    scrollbar-color: var(--color-brand-500) transparent;
  }

  .cyber-dropdown__options::-webkit-scrollbar {
    @apply w-1.5;
  }

  .cyber-dropdown__options::-webkit-scrollbar-track {
    @apply bg-transparent;
  }

  .cyber-dropdown__options::-webkit-scrollbar-thumb {
    @apply rounded-sm;
    background: var(--color-brand-500);
  }

  .cyber-dropdown__options::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.6);
  }

  .cyber-dropdown__option {
    @apply flex cursor-pointer items-center px-3 py-1.5 text-sm text-content transition-colors duration-200;
  }

  .cyber-dropdown__option:hover,
  .cyber-dropdown__option--highlighted {
    background: rgba(var(--color-brand-500-rgb), 0.08);
  }

  .cyber-dropdown__option--selected {
    background: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-dropdown__option--selected:hover,
  .cyber-dropdown__option--selected.cyber-dropdown__option--highlighted {
    background: rgba(var(--color-brand-500-rgb), 0.22);
  }

  .cyber-dropdown__option--create {
    @apply mt-1 border-t border-dashed pt-2.5;
    border-top-color: rgba(var(--color-brand-500-rgb), 0.3);
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.05) 0%,
      rgba(var(--color-brand-500-rgb), 0.02) 100%
    );
    position: relative;
    font-weight: 500;
    color: var(--color-brand-500);
    border-radius: var(--radius-sm);
    margin: 4px 4px 2px 4px;
    padding: 6px 12px;
    box-shadow: 0 0 0 1px rgba(var(--color-brand-500-rgb), 0.1);
    transition: all 0.2s ease;
  }

  .cyber-dropdown__option--create:hover,
  .cyber-dropdown__option--create.cyber-dropdown__option--highlighted {
    background: linear-gradient(
      135deg,
      rgba(var(--color-brand-500-rgb), 0.15) 0%,
      rgba(var(--color-brand-500-rgb), 0.08) 100%
    );
    border-top-color: rgba(var(--color-brand-500-rgb), 0.5);
    box-shadow:
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.15);
    transform: translateX(2px);
  }

  .cyber-dropdown__option--create::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: linear-gradient(
      180deg,
      var(--color-brand-500) 0%,
      rgba(var(--color-brand-500-rgb), 0.3) 100%
    );
    opacity: 0;
    transition: opacity 0.2s ease;
  }

  .cyber-dropdown__option--create:hover::before,
  .cyber-dropdown__option--create.cyber-dropdown__option--highlighted::before {
    opacity: 1;
  }

  .cyber-dropdown__checkbox-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    margin-right: 8px;
  }

  .cyber-dropdown__checkbox {
    @apply flex h-4 w-4 items-center justify-center rounded-sm border bg-input-bg transition-all duration-200;
    border-color: var(--color-brand-500);
  }

  .cyber-dropdown__option--selected .cyber-dropdown__checkbox {
    @apply border-error-500 bg-selected-bg;
  }

  .cyber-dropdown__checkbox i {
    @apply text-xs text-error-500;
  }

  .cyber-dropdown__option-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    margin-right: 8px;
    min-width: 16px;
    flex-shrink: 0;
  }

  .cyber-dropdown__create-icon {
    @apply text-sm;
    color: var(--color-brand-500);
    animation: pulse-glow 2s ease-in-out infinite;
  }

  @keyframes pulse-glow {
    0%,
    100% {
      opacity: 1;
      filter: drop-shadow(0 0 2px var(--color-brand-500));
    }
    50% {
      opacity: 0.8;
      filter: drop-shadow(0 0 6px var(--color-brand-500));
    }
  }

  .cyber-dropdown__option--create:hover .cyber-dropdown__create-icon,
  .cyber-dropdown__option--create.cyber-dropdown__option--highlighted .cyber-dropdown__create-icon {
    animation: none;
    filter: drop-shadow(0 0 4px var(--color-brand-500));
    transform: scale(1.1);
    transition: all 0.2s ease;
  }

  .cyber-dropdown__option-label {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .cyber-dropdown__option--create .cyber-dropdown__option-label {
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
    letter-spacing: 0.02em;
  }

  .cyber-dropdown__loading,
  .cyber-dropdown__loading-more,
  .cyber-dropdown__no-more {
    @apply p-3 text-center text-sm text-content-muted;
  }

  .cyber-dropdown__loading-more {
    @apply border-t;
    border-top-color: rgba(var(--color-brand-500-rgb), 0.1);
    background: rgba(var(--color-error-rgb), 0.05);
  }

  .cyber-dropdown__no-more {
    @apply border-t italic;
    border-top-color: rgba(var(--color-brand-500-rgb), 0.1);
    background: rgba(var(--color-error-rgb), 0.03);
  }

  .cyber-dropdown__loading i,
  .cyber-dropdown__loading-more i {
    @apply mr-1.5 text-brand-500;
  }

  .cyber-dropdown__empty {
    @apply p-4 text-center text-sm text-content-muted;
  }

  .cyber-dropdown--disabled .cyber-dropdown__trigger {
    @apply cursor-not-allowed border-disabled-border bg-disabled-bg opacity-60;

    box-shadow: none !important;
    transform: none !important;
  }

  .cyber-dropdown--disabled .cyber-dropdown__trigger:hover {
    @apply border-disabled-border;
    box-shadow: none !important;
    transform: none !important;
  }

  .cyber-dropdown-enter-active,
  .cyber-dropdown-leave-active {
    transition:
      opacity 0.2s,
      transform 0.2s;
  }

  .cyber-dropdown-enter-from,
  .cyber-dropdown-leave-to {
    opacity: 0;
    transform: translateY(-8px);
  }

  @media (max-width: 768px) {
    .cyber-dropdown__trigger {
      padding: 3px 8px;
      font-size: 0.8125rem;
      height: 32px;
    }

    .cyber-dropdown__option {
      padding: 4px 8px;
      font-size: 0.8125rem;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .cyber-dropdown__trigger,
    .cyber-dropdown__arrow,
    .cyber-dropdown__clear,
    .cyber-dropdown__tag,
    .cyber-dropdown__tag-remove,
    .cyber-dropdown__checkbox,
    .cyber-dropdown__option,
    .cyber-dropdown-enter-active,
    .cyber-dropdown-leave-active {
      transition: none !important;
    }
  }
</style>
