<script setup lang="ts">
  import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { PaginationProps, PaginationEmits } from './types'

  defineOptions({
    name: 'CyberPagination',
  })

  const _props = withDefaults(defineProps<PaginationProps>(), {
    pageSize: 20,
    pageSizeOptions: () => [10, 20, 50, 100],
    showPageSizeSelector: true,
    showQuickJumper: true,
    showTotal: true,
    maxVisiblePages: 7,
    autoScrollToTop: true,
    scrollTarget: '',
  })

  const emit = defineEmits<PaginationEmits>()
  const { $t } = useTexts()

  const currentPageSize = ref(_props.pageSize)
  const sizeDropdownVisible = ref(false)
  const jumpPageValue = ref<number | string>('')
  const selectorRef = ref<HTMLElement>()
  const optionsRef = ref<HTMLElement>()

  type Placement = 'up' | 'down'
  const popupLeft = ref(0)
  const popupTop = ref(0)
  const popupWidth = ref(0)
  const popupPlacement = ref<Placement>('down')

  const popupStyle = computed(() => ({
    position: 'fixed',
    left: `${popupLeft.value}px`,
    top: `${popupTop.value}px`,
    width: `${popupWidth.value}px`,
    zIndex: 9999,
  }))

  const totalPages = computed(() => Math.ceil(_props.total / currentPageSize.value))

  const visiblePages = computed(() => {
    const current = _props.currentPage
    const total = totalPages.value
    const maxVisible = _props.maxVisiblePages

    if (total <= maxVisible) {
      return Array.from({ length: total }, (_, i) => i + 1)
    }

    const pages: (number | string)[] = []
    const halfVisible = Math.floor(maxVisible / 2)

    pages.push(1)

    let start = Math.max(2, current - halfVisible)
    let end = Math.min(total - 1, current + halfVisible)

    if (end - start + 1 < maxVisible - 2) {
      if (start === 2) {
        end = Math.min(total - 1, start + maxVisible - 3)
      } else {
        start = Math.max(2, end - maxVisible + 3)
      }
    }

    if (start > 2) {
      pages.push('...')
    }

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    if (end < total - 1) {
      pages.push('...')
    }

    if (total > 1) {
      pages.push(total)
    }

    return pages
  })

  const scrollToTop = () => {
    if (!_props.autoScrollToTop) {
      return
    }

    nextTick(() => {
      let targetElement: HTMLElement | null = null

      if (_props.scrollTarget) {
        targetElement = document.querySelector(_props.scrollTarget)
      }

      if (targetElement) {
        targetElement.scrollIntoView({
          behavior: 'smooth',
          block: 'start',
          inline: 'nearest',
        })
      } else {
        window.scrollTo({
          top: 0,
          behavior: 'smooth',
        })
      }
    })
  }

  const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPages.value && page !== _props.currentPage) {
      emit('update:currentPage', page)
      emit('page-change', page)
      scrollToTop()
    }
  }

  const computePopupPosition = async (preOpen = false) => {
    if (!preOpen) {
      await nextTick()
    }
    const trigger = selectorRef.value
    if (!trigger) {
      return
    }
    const rect = trigger.getBoundingClientRect()
    const estimatedHeight = preOpen ? 160 : optionsRef.value?.offsetHeight || 160
    const spaceBelow = window.innerHeight - rect.bottom
    const spaceAbove = rect.top
    const openUp = spaceBelow < estimatedHeight + 8 && spaceAbove > spaceBelow
    popupPlacement.value = openUp ? 'up' : 'down'
    popupLeft.value = Math.round(rect.left)
    popupWidth.value = Math.round(rect.width)
    popupTop.value = Math.round(openUp ? rect.top - estimatedHeight - 4 : rect.bottom + 4)
  }

  const openSizeDropdown = async () => {
    await computePopupPosition(true)
    sizeDropdownVisible.value = true
    await computePopupPosition(false)
  }

  const toggleSizeDropdown = async () => {
    if (sizeDropdownVisible.value) {
      sizeDropdownVisible.value = false
    } else {
      await openSizeDropdown()
    }
  }

  const selectPageSize = (size: number) => {
    currentPageSize.value = size
    sizeDropdownVisible.value = false
    emit('update:pageSize', size)
    emit('page-size-change', size)

    const newTotalPages = Math.ceil(_props.total / size)
    if (_props.currentPage > newTotalPages) {
      goToPage(newTotalPages)
    } else {
      scrollToTop()
    }
  }

  const handleJump = () => {
    const page = parseInt(jumpPageValue.value as string)
    if (page >= 1 && page <= totalPages.value) {
      goToPage(page)
      jumpPageValue.value = ''
    }
  }

  const handleClickOutside = (event: Event) => {
    const target = event.target as Node
    const insideTrigger = selectorRef.value?.contains(target)
    const insidePopup = optionsRef.value?.contains(target)
    if (!insideTrigger && !insidePopup) {
      sizeDropdownVisible.value = false
    }
  }

  const handleWindowUpdate = () => {
    if (sizeDropdownVisible.value) {
      computePopupPosition()
    }
  }

  onMounted(() => {
    document.addEventListener('click', handleClickOutside)
    window.addEventListener('resize', handleWindowUpdate, { passive: true })
    window.addEventListener('scroll', handleWindowUpdate, { passive: true })
  })

  onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
    window.removeEventListener('resize', handleWindowUpdate)
    window.removeEventListener('scroll', handleWindowUpdate)
  })

  watch(
    () => _props.pageSize,
    (newSize) => {
      currentPageSize.value = newSize
    }
  )
</script>

<template>
  <div v-if="totalPages > 0" class="cyber-pagination">
    <div class="pagination-container">
      <div v-if="_props.showPageSizeSelector" class="page-size-section">
        <span class="size-label">{{ $t('components.pagination.perPage') }}</span>
        <div ref="selectorRef" class="size-selector">
          <div class="size-current" :class="{ active: sizeDropdownVisible }" @click="toggleSizeDropdown">
            <span>{{ currentPageSize }}</span>
            <i class="fas fa-angle-down" :class="{ rotate: sizeDropdownVisible }" />
          </div>
          <teleport to="body">
            <transition name="fade-slide">
              <div v-if="sizeDropdownVisible" ref="optionsRef" class="size-options-popup" :style="popupStyle" @mousedown.stop>
                <div
                  v-for="size in _props.pageSizeOptions"
                  :key="size"
                  class="size-option"
                  :class="{ selected: size === currentPageSize }"
                  @click="selectPageSize(size)"
                >
                  {{ size }}
                </div>
              </div>
            </transition>
          </teleport>
        </div>
        <span class="size-label">{{ $t('components.pagination.items') }}</span>
      </div>

      <div class="page-controls">
        <button
          class="page-btn prev-btn"
          :disabled="_props.currentPage === 1"
          :title="$t('components.pagination.prevPage')"
          @click="goToPage(_props.currentPage - 1)"
        >
          <i class="fas fa-chevron-left" />
        </button>

        <template v-for="page in visiblePages" :key="page">
          <button
            v-if="page !== '...'"
            class="page-btn number-btn"
            :class="{ active: page === _props.currentPage }"
            @click="goToPage(page as number)"
          >
            {{ page }}
          </button>
          <span v-else class="page-ellipsis">•••</span>
        </template>

        <button
          class="page-btn next-btn"
          :disabled="_props.currentPage === totalPages"
          :title="$t('components.pagination.nextPage')"
          @click="goToPage(_props.currentPage + 1)"
        >
          <i class="fas fa-chevron-right" />
        </button>
      </div>

      <div class="page-info">
        <div v-if="_props.showQuickJumper && totalPages > 5" class="page-jump">
          <span class="jump-label">{{ $t('components.pagination.jumpTo') }}</span>
          <input
            v-model="jumpPageValue"
            type="number"
            :min="1"
            :max="totalPages"
            class="jump-input"
            :placeholder="$t('components.pagination.page')"
            @keyup.enter="handleJump"
            @blur="handleJump"
          />
        </div>

        <div v-if="_props.showTotal" class="total-info">
          <span class="total-text">
            {{ $t('components.pagination.total') }} <em class="total-number">{{ _props.total.toLocaleString() }}</em>
            {{ $t('components.pagination.items') }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .cyber-pagination {
    @apply flex justify-center;
  }

  .pagination-container {
    @apply flex items-center;
    gap: var(--space-xl);
  }

  .page-size-section {
    @apply flex items-center text-sm;
    gap: var(--space-xs);
    color: var(--color-content-default);
  }

  .size-label {
    @apply whitespace-nowrap font-medium;
  }

  .size-selector {
    @apply relative;
  }

  .size-current {
    @apply flex cursor-pointer items-center justify-between rounded border px-2 transition-all;
    gap: var(--space-xs);
    height: 30px;
    background: var(--color-background-800);
    border-color: var(--color-border-default);
    color: var(--color-content-default);
    min-width: 50px;
  }

  .size-current:hover,
  .size-current.active {
    border-color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.05);
  }

  .size-current i {
    @apply text-xs transition-transform;
    color: var(--color-content-muted);
  }

  .size-current i.rotate {
    @apply rotate-180;
  }

  .size-options-popup {
    @apply fixed overflow-hidden rounded border;
    background: var(--color-background-800);
    border-color: var(--color-border-default);
    box-shadow:
      0 4px 12px rgba(0, 0, 0, 0.2),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.1);
    transform: none !important;
  }

  .size-option {
    @apply cursor-pointer px-2 py-1 text-sm transition-colors;
  }

  .size-option:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .size-option.selected {
    @apply text-content-heading;
    background: var(--color-brand-500);
  }

  .page-controls {
    @apply flex items-center;
    gap: var(--space-xs);
  }

  .page-btn {
    @apply flex cursor-pointer select-none items-center justify-center rounded border px-2 text-sm font-medium transition-all;
    height: 30px;
    min-width: 30px;
    background: var(--color-background-800);
    border-color: var(--color-border-default);
    color: var(--color-content-muted);
  }

  .page-btn:hover:not(:disabled) {
    @apply -translate-y-px;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: var(--color-brand-500);
    color: var(--color-brand-500);
    box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .page-btn:disabled {
    @apply cursor-not-allowed opacity-40;
  }

  .page-btn.active {
    @apply font-semibold;
    background: var(--color-brand-500);
    border-color: var(--color-brand-500);
    color: var(--color-background-50);
    box-shadow:
      0 2px 8px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 16px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .page-ellipsis {
    @apply flex items-center justify-center font-bold;
    height: 30px;
    min-width: 30px;
    color: var(--color-content-muted);
    letter-spacing: 2px;
  }

  .page-info {
    @apply flex items-center text-sm;
    gap: var(--space-md);
  }

  .page-jump {
    @apply flex items-center;
    gap: var(--space-xs);
    color: var(--color-content-default);
  }

  .jump-label {
    @apply whitespace-nowrap font-medium;
  }

  .jump-input {
    @apply w-12 rounded border px-2 text-center text-sm transition-all;
    height: 30px;
    background: var(--color-background-800);
    border-color: var(--color-border-default);
    color: var(--color-content-default);
  }

  .jump-input::-webkit-outer-spin-button,
  .jump-input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  .jump-input[type='number'] {
    -moz-appearance: textfield;
  }

  .jump-input:focus {
    @apply outline-none;
    border-color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.05);
  }

  .total-info {
    @apply font-medium;
    color: var(--color-content-default);
  }

  .total-text {
    @apply whitespace-nowrap;
  }

  .total-number {
    @apply font-semibold not-italic;
    color: var(--color-brand-500);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .fade-slide-enter-active,
  .fade-slide-leave-active {
    @apply transition-all duration-200;
  }

  .fade-slide-enter-from {
    @apply -translate-y-2 opacity-0;
  }

  .fade-slide-leave-to {
    @apply -translate-y-2 opacity-0;
  }

  @media (max-width: 1024px) {
    .pagination-container {
      @apply p-3;
      gap: var(--space-lg);
    }
  }

  @media (max-width: 768px) {
    .pagination-container {
      @apply flex-col items-center gap-4;
    }

    .page-size-section {
      @apply order-2;
    }

    .page-controls {
      @apply order-1;
    }

    .page-info {
      @apply order-3 flex-col gap-2;
    }

    .page-btn {
      height: 28px;
      min-width: 28px;
      @apply text-xs;
    }
  }

  @media (max-width: 480px) {
    .page-controls {
      gap: 2px;
    }

    .page-btn {
      height: 26px;
      min-width: 26px;
      @apply px-1;
    }

    .pagination-container {
      @apply -mx-2;
    }
  }
</style>
