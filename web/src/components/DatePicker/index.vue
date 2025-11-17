<script setup lang="ts">
  import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { DatePickerProps, DatePickerEmits, Day } from './types'

  defineOptions({
    name: 'CyberDatePicker',
  })

  const { $t } = useTexts()

  const props = withDefaults(defineProps<DatePickerProps>(), {
    modelValue: undefined,
    placeholder: undefined,
    startPlaceholder: undefined,
    endPlaceholder: undefined,
    rangeSeparator: undefined,
    disabled: false,
    type: 'date',
    size: 'medium',
    format: 'YYYY-MM-DD',
    disabledDate: () => false,
  })

  const emit = defineEmits<DatePickerEmits>()

  const datepickerRef = ref<HTMLElement | null>(null)
  const pickerVisible = ref(false)
  const showYearMonthPanel = ref(false)
  const currentDate = ref(new Date())
  const panelPosition = ref({ top: 0, left: 0 })

  const isRange = computed(() => props.type === 'daterange')
  const rangeStartDate = ref<Date | null>(null)
  const rangeEndDate = ref<Date | null>(null)
  const rangeSelecting = ref(false)

  const currentYear = ref(currentDate.value.getFullYear())
  const currentMonth = ref(currentDate.value.getMonth())

  const initSelectedDate = () => {
    if (props.modelValue) {
      if (isRange.value && Array.isArray(props.modelValue)) {
        rangeStartDate.value = props.modelValue[0]
        rangeEndDate.value = props.modelValue[1]
        if (rangeStartDate.value) {
          currentYear.value = rangeStartDate.value.getFullYear()
          currentMonth.value = rangeStartDate.value.getMonth()
        }
      } else if (!isRange.value && !Array.isArray(props.modelValue)) {
        currentDate.value = props.modelValue
        currentYear.value = currentDate.value.getFullYear()
        currentMonth.value = currentDate.value.getMonth()
      }
    }
  }

  watch(
    () => props.modelValue,
    () => {
      initSelectedDate()
    },
    { deep: true, immediate: true }
  )

  const weekdayNames = computed(() => $t('components.datePicker.weekdays'))
  const monthNames = computed(() => $t('components.datePicker.months'))
  const yearMonthLabel = computed(() =>
    ($t('components.datePicker.yearMonth') as string)
      .replace('{year}', String(currentYear.value))
      .replace('{month}', String(currentMonth.value + 1))
  )

  const yearList = computed(() => {
    const years = []
    const baseYear = currentYear.value
    for (let i = baseYear - 5; i <= baseYear + 5; i++) {
      years.push(i)
    }
    return years
  })

  const formatDate = (date: Date): string => {
    if (!date) return ''
    const year = date.getFullYear()
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const day = date.getDate().toString().padStart(2, '0')
    return `${year}-${month}-${day}`
  }

  const formattedDate = computed(() => {
    if (!isRange.value && props.modelValue && !Array.isArray(props.modelValue)) {
      return formatDate(props.modelValue)
    }
    return ''
  })

  const formattedStartDate = computed(() => {
    if (isRange.value && rangeStartDate.value) {
      return formatDate(rangeStartDate.value)
    } else if (isRange.value && Array.isArray(props.modelValue) && props.modelValue[0]) {
      return formatDate(props.modelValue[0])
    }
    return ''
  })

  const formattedEndDate = computed(() => {
    if (isRange.value && rangeEndDate.value) {
      return formatDate(rangeEndDate.value)
    } else if (isRange.value && Array.isArray(props.modelValue) && props.modelValue[1]) {
      return formatDate(props.modelValue[1])
    }
    return ''
  })

  const isValidRange = computed(() => isRange.value && rangeStartDate.value && rangeEndDate.value)

  const days = computed(() => {
    const daysArray: Day[] = []
    const firstDay = new Date(currentYear.value, currentMonth.value, 1)
    const lastDay = new Date(currentYear.value, currentMonth.value + 1, 0)
    const firstDayWeekday = firstDay.getDay()
    const prevMonthLastDate = new Date(currentYear.value, currentMonth.value, 0).getDate()

    for (let i = 0; i < firstDayWeekday; i++) {
      const date = prevMonthLastDate - firstDayWeekday + i + 1
      const fullDate = new Date(currentYear.value, currentMonth.value - 1, date)
      daysArray.push({
        date,
        fullDate,
        isPrevMonth: true,
        isNextMonth: false,
        isToday: isSameDate(fullDate, new Date()),
        isDisabled: props.disabledDate(fullDate),
        isSelected: isSameDate(fullDate, currentDate.value),
        isRangeStart: isRangeStartDate(fullDate),
        isRangeEnd: isRangeEndDate(fullDate),
        isInRange: isDateInRange(fullDate),
      })
    }

    for (let i = 1; i <= lastDay.getDate(); i++) {
      const fullDate = new Date(currentYear.value, currentMonth.value, i)
      daysArray.push({
        date: i,
        fullDate,
        isPrevMonth: false,
        isNextMonth: false,
        isToday: isSameDate(fullDate, new Date()),
        isDisabled: props.disabledDate(fullDate),
        isSelected: isSameDate(fullDate, currentDate.value),
        isRangeStart: isRangeStartDate(fullDate),
        isRangeEnd: isRangeEndDate(fullDate),
        isInRange: isDateInRange(fullDate),
      })
    }

    const nextMonthDays = 42 - daysArray.length
    for (let i = 1; i <= nextMonthDays; i++) {
      const fullDate = new Date(currentYear.value, currentMonth.value + 1, i)
      daysArray.push({
        date: i,
        fullDate,
        isPrevMonth: false,
        isNextMonth: true,
        isToday: isSameDate(fullDate, new Date()),
        isDisabled: props.disabledDate(fullDate),
        isSelected: isSameDate(fullDate, currentDate.value),
        isRangeStart: isRangeStartDate(fullDate),
        isRangeEnd: isRangeEndDate(fullDate),
        isInRange: isDateInRange(fullDate),
      })
    }

    return daysArray
  })

  function isSameDate(date1: Date, date2: Date | null): boolean {
    if (!date1 || !date2) return false
    return (
      date1.getFullYear() === date2.getFullYear() && date1.getMonth() === date2.getMonth() && date1.getDate() === date2.getDate()
    )
  }

  function isRangeStartDate(date: Date): boolean {
    return isRange.value && isSameDate(date, rangeStartDate.value)
  }

  function isRangeEndDate(date: Date): boolean {
    return isRange.value && isSameDate(date, rangeEndDate.value)
  }

  function isDateInRange(date: Date): boolean {
    if (!isRange.value || !rangeStartDate.value || !rangeEndDate.value) return false
    const time = date.getTime()
    return time > rangeStartDate.value.getTime() && time < rangeEndDate.value.getTime()
  }

  const togglePicker = () => {
    if (props.disabled) return
    pickerVisible.value = !pickerVisible.value
    if (pickerVisible.value) {
      updatePanelPosition()
      showYearMonthPanel.value = false
    }
  }

  const updatePanelPosition = () => {
    if (!datepickerRef.value) return
    const rect = datepickerRef.value.getBoundingClientRect()
    panelPosition.value = {
      top: rect.bottom + window.scrollY + 5,
      left: rect.left + window.scrollX,
    }
  }

  const prevMonth = () => {
    if (currentMonth.value === 0) {
      currentYear.value--
      currentMonth.value = 11
    } else {
      currentMonth.value--
    }
  }

  const nextMonth = () => {
    if (currentMonth.value === 11) {
      currentYear.value++
      currentMonth.value = 0
    } else {
      currentMonth.value++
    }
  }

  const prevYear = () => {
    currentYear.value--
  }

  const nextYear = () => {
    currentYear.value++
  }

  const selectYear = (year: number) => {
    currentYear.value = year
    showYearMonthPanel.value = false
  }

  const selectMonth = (month: number) => {
    currentMonth.value = month
    showYearMonthPanel.value = false
  }

  const selectDate = (day: Day) => {
    if (day.isDisabled) return

    if (isRange.value) {
      if (!rangeSelecting.value || !rangeStartDate.value) {
        rangeStartDate.value = new Date(day.fullDate)
        rangeEndDate.value = null
        rangeSelecting.value = true
      } else {
        const startTime = rangeStartDate.value.getTime()
        const currentTime = day.fullDate.getTime()

        if (currentTime < startTime) {
          rangeEndDate.value = new Date(rangeStartDate.value)
          rangeStartDate.value = new Date(day.fullDate)
        } else {
          rangeEndDate.value = new Date(day.fullDate)
        }
        rangeSelecting.value = false
      }
    } else {
      currentDate.value = new Date(day.fullDate)
      emitSelectedValue()
      pickerVisible.value = false
    }
  }

  const confirmSelection = () => {
    if (isRange.value && rangeStartDate.value && rangeEndDate.value) {
      emitSelectedValue()
      pickerVisible.value = false
    }
  }

  const clearSelection = () => {
    if (isRange.value) {
      rangeStartDate.value = null
      rangeEndDate.value = null
      rangeSelecting.value = false
      emit('update:modelValue', [null, null])
      emit('change', [null, null])
    } else {
      currentDate.value = new Date()
      emit('update:modelValue', null)
      emit('change', null)
    }
  }

  const emitSelectedValue = () => {
    if (isRange.value) {
      if (rangeStartDate.value && rangeEndDate.value) {
        const value: [Date, Date] = [new Date(rangeStartDate.value), new Date(rangeEndDate.value)]
        emit('update:modelValue', value)
        emit('change', value)
      }
    } else {
      emit('update:modelValue', new Date(currentDate.value))
      emit('change', new Date(currentDate.value))
    }
  }

  const closeOnClickOutside = (e: MouseEvent) => {
    if (pickerVisible.value && datepickerRef.value) {
      const target = e.target as Node
      if (datepickerRef.value.contains(target)) return

      const panel = document.querySelector('.cyber-date-picker__panel')
      if (panel && panel.contains(target)) return

      pickerVisible.value = false
      showYearMonthPanel.value = false

      if (isRange.value && rangeSelecting.value) {
        rangeSelecting.value = false
      }
    }
  }

  onMounted(() => {
    initSelectedDate()
    window.addEventListener('resize', updatePanelPosition)
    document.addEventListener('click', closeOnClickOutside)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('resize', updatePanelPosition)
    document.removeEventListener('click', closeOnClickOutside)
  })
</script>

<template>
  <div
    ref="datepickerRef"
    class="cyber-date-picker"
    :class="[
      `cyber-date-picker--${size}`,
      {
        'cyber-date-picker--disabled': disabled,
        'cyber-date-picker--range': isRange,
        'cyber-date-picker--focused': pickerVisible,
      },
    ]"
  >
    <div class="cyber-date-picker__trigger" :class="{ disabled: disabled }" @click="togglePicker">
      <template v-if="!isRange">
        <CyberInput
          :model-value="formattedDate"
          :placeholder="placeholder || $t('components.datePicker.selectDate')"
          :disabled="disabled"
          :size="size"
          readonly
        >
          <template #suffix>
            <i class="fas fa-calendar-alt" />
          </template>
        </CyberInput>
      </template>

      <template v-else>
        <div class="cyber-date-range-input">
          <CyberInput
            :model-value="formattedStartDate"
            :placeholder="startPlaceholder || $t('components.datePicker.startDate')"
            :disabled="disabled"
            :size="size"
            readonly
            class="range-start-input"
          >
            <template #suffix>
              <i class="fas fa-calendar-alt" />
            </template>
          </CyberInput>

          <span class="range-separator">{{ rangeSeparator || $t('components.datePicker.to') }}</span>

          <CyberInput
            :model-value="formattedEndDate"
            :placeholder="endPlaceholder || $t('components.datePicker.endDate')"
            :disabled="disabled"
            :size="size"
            readonly
            class="range-end-input"
          >
            <template #suffix>
              <i class="fas fa-calendar-alt" />
            </template>
          </CyberInput>
        </div>
      </template>
    </div>

    <Teleport to="body">
      <transition name="cyber-datepicker-fade">
        <div
          v-show="pickerVisible"
          class="cyber-date-picker__panel"
          :style="[{ left: panelPosition.left + 'px' }, { top: panelPosition.top + 'px' }, { zIndex: 2000 }]"
          @click.stop
        >
          <div class="cyber-date-picker__header">
            <button class="header-btn prev-year-btn" :title="$t('components.datePicker.prevYear')" @click.stop="prevYear">
              <i class="fas fa-angle-double-left" />
            </button>
            <button class="header-btn prev-month-btn" :title="$t('components.datePicker.prevMonth')" @click.stop="prevMonth">
              <i class="fas fa-angle-left" />
            </button>

            <div class="current-date" @click.stop="showYearMonthPanel = !showYearMonthPanel">
              {{ yearMonthLabel }}
              <i class="fas ml-1" :class="[showYearMonthPanel ? 'fa-chevron-up' : 'fa-chevron-down']" />
            </div>

            <button class="header-btn next-month-btn" :title="$t('components.datePicker.nextMonth')" @click.stop="nextMonth">
              <i class="fas fa-angle-right" />
            </button>
            <button class="header-btn next-year-btn" :title="$t('components.datePicker.nextYear')" @click.stop="nextYear">
              <i class="fas fa-angle-double-right" />
            </button>
          </div>

          <div v-if="showYearMonthPanel" class="cyber-date-picker__year-month-panel">
            <div class="year-month-grid">
              <button
                v-for="year in yearList"
                :key="'year-' + year"
                class="year-btn"
                :class="{ current: year === currentYear }"
                @click.stop="selectYear(year)"
              >
                {{ year }}
              </button>
            </div>

            <div class="month-grid">
              <button
                v-for="(month, index) in monthNames"
                :key="'month-' + index"
                class="month-btn"
                :class="{ current: index === currentMonth }"
                @click.stop="selectMonth(index)"
              >
                {{ month }}
              </button>
            </div>
          </div>

          <div class="cyber-date-picker__weekdays">
            <div v-for="weekday in weekdayNames" :key="weekday" class="weekday-cell">
              {{ weekday }}
            </div>
          </div>

          <div class="cyber-date-picker__days">
            <div
              v-for="(day, index) in days"
              :key="index"
              class="day-cell"
              :class="{
                'prev-month': day.isPrevMonth,
                'next-month': day.isNextMonth,
                today: day.isToday,
                selected: day.isSelected,
                'range-start': day.isRangeStart,
                'range-end': day.isRangeEnd,
                'in-range': day.isInRange,
                disabled: day.isDisabled,
              }"
              @click.stop="selectDate(day)"
            >
              {{ day.date }}
            </div>
          </div>

          <div v-if="isRange" class="cyber-date-picker__footer">
            <button class="footer-btn clear-btn" @click.stop="clearSelection">
              {{ $t('components.datePicker.clear') }}
            </button>
            <button class="footer-btn confirm-btn" :disabled="!isValidRange" @click.stop="confirmSelection">
              {{ $t('components.datePicker.confirm') }}
            </button>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped>
  .cyber-date-picker {
    @apply relative inline-flex w-full;
  }

  .cyber-date-picker__trigger {
    @apply w-full cursor-pointer;
  }

  .cyber-date-picker--range .cyber-date-range-input {
    @apply flex w-full items-center;
  }

  .range-separator {
    @apply px-1.5 text-content-muted;
  }

  .cyber-date-picker--disabled .cyber-date-picker__trigger {
    @apply cursor-not-allowed;
  }

  .cyber-date-picker__panel {
    @apply absolute w-[300px] overflow-hidden rounded-md p-2.5 text-content;
    background: rgba(var(--color-background-900-rgb), 0.95);
    border: 1px solid var(--color-border-default);
    box-shadow:
      0 4px 12px rgba(var(--color-background-900-rgb), 0.3),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.2);
    backdrop-filter: blur(8px);
    font-family: 'Roboto Mono', monospace;
  }

  .cyber-date-picker--range .cyber-date-picker__panel {
    @apply w-[320px];
  }

  .cyber-date-picker__header {
    @apply mb-1.5 flex items-center justify-between px-1.5 pb-2.5;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .current-date {
    @apply flex cursor-pointer items-center rounded px-2 py-1 font-semibold text-brand-400;
    transition: background-color 0.2s;
  }

  .current-date:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .header-btn {
    @apply flex h-7 w-7 cursor-pointer items-center justify-center rounded border-none bg-transparent text-content transition-all duration-200;
  }

  .header-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    @apply text-brand-500;
  }

  .cyber-date-picker__year-month-panel {
    @apply absolute left-0 top-0 z-10 flex h-full w-full flex-col p-2.5;
    background: rgba(var(--color-background-900-rgb), 0.98);
  }

  .year-month-grid {
    @apply mb-2.5 grid grid-cols-3 gap-1.5;
  }

  .month-grid {
    @apply grid grid-cols-4 gap-1.5;
  }

  .year-btn,
  .month-btn {
    @apply cursor-pointer rounded py-2 text-center text-content;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    transition: all 0.2s;
  }

  .year-btn:hover,
  .month-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .year-btn.current,
  .month-btn.current {
    @apply font-semibold;
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: var(--color-brand-500);
    @apply text-brand-500;
  }

  .cyber-date-picker__weekdays {
    @apply mb-1.5 grid grid-cols-7;
  }

  .weekday-cell {
    @apply py-1.5 text-center text-xs font-semibold text-brand-400;
  }

  .cyber-date-picker__days {
    @apply grid grid-cols-7 gap-0.5;
  }

  .day-cell {
    @apply relative flex h-7 cursor-pointer items-center justify-center rounded text-sm;
    transition: all 0.2s;
  }

  .day-cell:hover:not(.disabled) {
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .day-cell.prev-month,
  .day-cell.next-month {
    @apply text-content-disabled;
  }

  .day-cell.today {
    @apply font-semibold text-brand-400;
  }

  .day-cell.selected {
    @apply font-semibold text-brand-500;
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .day-cell.disabled {
    @apply cursor-not-allowed text-content-disabled;
  }

  .day-cell.range-start,
  .day-cell.range-end {
    @apply z-10 font-semibold;
    background: var(--color-brand-500);
    color: var(--color-text-on-brand);
  }

  .day-cell.in-range {
    @apply rounded-none;
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-date-picker__footer {
    @apply mt-2.5 flex justify-end pt-2.5;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .footer-btn {
    @apply ml-2 cursor-pointer rounded border-none bg-transparent px-3 py-1.5 text-sm text-content transition-all duration-200;
  }

  .footer-btn:hover:not(:disabled) {
    background: rgba(var(--color-brand-500-rgb), 0.1);
    @apply text-brand-500;
  }

  .confirm-btn {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .confirm-btn:hover:not(:disabled) {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-color: var(--color-brand-500);
  }

  .footer-btn:disabled {
    @apply cursor-not-allowed opacity-50;
  }

  .cyber-datepicker-fade-enter-active,
  .cyber-datepicker-fade-leave-active {
    @apply transition-all duration-300;
  }

  .cyber-datepicker-fade-enter-from,
  .cyber-datepicker-fade-leave-to {
    @apply -translate-y-2.5 opacity-0;
  }

  .cyber-date-picker--small .day-cell {
    @apply h-6 text-xs;
  }

  .cyber-date-picker--large .day-cell {
    @apply h-8 text-base;
  }

  @media (max-width: 640px) {
    .cyber-date-picker__panel {
      @apply w-[290px];
    }

    .cyber-date-picker--range .cyber-date-range-input {
      @apply flex-col items-start;
    }

    .range-separator {
      @apply my-1.5;
    }
  }

  .cyber-date-range-input {
    @apply flex w-full items-center gap-1.5;
  }

  .range-start-input,
  .range-end-input {
    @apply flex-1;
  }
</style>
