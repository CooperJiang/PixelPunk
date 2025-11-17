<script setup lang="ts">
  import { computed, nextTick, ref, watch } from 'vue'
  import type { MultiSelectorEmits, MultiSelectorProps, SelectorOption } from './types'
  import {
    compareDurations,
    formatDuration,
    getDurationIcon,
    isDurationExists,
    normalizeDuration,
    validateDuration,
  } from './utils'
  import { useTexts } from '@/composables/useTexts'

  /* 🚨 必须：定义组件名 */
  defineOptions({
    name: 'MultiSelector',
  })

  /* 使用类型文件中的接口 */
  const props = withDefaults(defineProps<MultiSelectorProps>(), {
    modelValue: () => [],
    inputIdPrefix: 'selector',
    disabled: false,
    options: () => [],
    forcedValue: null,
    size: 'md',
    showIcons: true,
    rounded: 'sm',
    defaultIcon: 'fas fa-clock',
    forcedIcon: 'fas fa-infinity',
    editable: false,
    addText: '',
    maxOptions: 10,
    isGuest: false,
  })

  const { $t } = useTexts()

  const emit = defineEmits<MultiSelectorEmits>()

  /* 输入状态 */
  const isAddingOption = ref(false)
  const inputValue = ref('')
  const inputError = ref('')
  const inputErrorCode = ref('')

  const displayAddText = computed(() => props.addText || $t('components.multiSelector.addOption'))
  const hasPermanentOption = computed(() => (props.modelValue || []).includes('permanent'))
  const placeholderText = computed(() =>
    props.isGuest || hasPermanentOption.value
      ? $t('components.multiSelector.placeholders.basic')
      : $t('components.multiSelector.placeholders.withPermanent')
  )

  /* 错误提示文案 */
  const inputErrorText = computed(() => {
    if (!inputErrorCode.value) return ''
    const known = [
      'empty',
      'guestNoPermanent',
      'invalidFormat',
      'mustBePositive',
      'minutesMin',
      'minutesMax',
      'hoursMax',
      'guestHoursMax',
      'daysMax',
      'guestDaysMax',
      'exists',
    ]
    return known.includes(inputErrorCode.value)
      ? $t(`components.multiSelector.errors.${inputErrorCode.value}`)
      : inputErrorCode.value
  })

  /* 本地选中的值 */
  const selectedValues = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  /* 尺寸类 */
  const sizeClasses = {
    sm: 'text-sm px-2.5 py-1',
    md: 'text-sm px-3 py-1.5',
    lg: 'text-base px-4 py-2',
  }

  /* 圆角类 */
  const roundedClasses = {
    none: 'rounded-none',
    sm: 'rounded',
    md: 'rounded-md',
    lg: 'rounded-lg',
  }

  /* 显示的选项（优先使用传入的options，否则使用默认选项） */
  const displayOptions = computed(() => {
    if (props.options && props.options.length > 0) {
      return [...props.options].sort((a, b) => compareDurations(a.value, b.value))
    }

    if (props.modelValue && props.modelValue.length > 0) {
      return props.modelValue
        .map((value) => ({
          value,
          label: formatDuration(value, $t),
          icon: getDurationIcon(value),
        }))
        .sort((a, b) => compareDurations(a.value, b.value))
    }

    return [
      { value: '1h', label: formatDuration('1h', $t), icon: 'fas fa-clock text-xs' },
      { value: '3d', label: formatDuration('3d', $t), icon: 'fas fa-calendar-day text-xs' },
      { value: '7d', label: formatDuration('7d', $t), icon: 'fas fa-calendar-week text-xs' },
      { value: '30d', label: formatDuration('30d', $t), icon: 'fas fa-calendar-alt text-xs' },
      { value: 'permanent', label: formatDuration('permanent', $t), icon: 'fas fa-infinity text-xs' },
    ].sort((a, b) => compareDurations(a.value, b.value))
  })

  const inputRef = ref<HTMLInputElement>()

  const canAddMore = computed(() => displayOptions.value.length < props.maxOptions)

  const startAddingOption = () => {
    isAddingOption.value = true
    inputValue.value = ''
    inputError.value = ''
    inputErrorCode.value = ''
    nextTick(() => {
      inputRef.value?.focus()
    })
  }

  const cancelAddingOption = () => {
    isAddingOption.value = false
    inputValue.value = ''
    inputError.value = ''
    inputErrorCode.value = ''
  }

  const confirmAddOption = () => {
    const value = inputValue.value.trim()

    if (!value) {
      inputErrorCode.value = 'empty'
      inputError.value = inputErrorText.value
      return
    }

    const validateFn = props.validateFn || validateDuration
    const error = validateFn(value, props.isGuest)

    if (error) {
      inputErrorCode.value = error
      inputError.value = inputErrorText.value
      emit('validation-error', inputError.value)
      return
    }

    const normalizedValue = normalizeDuration(value)

    if (isDurationExists(normalizedValue, props.modelValue || [])) {
      inputErrorCode.value = 'exists'
      inputError.value = inputErrorText.value
      return
    }

    const newOption: SelectorOption = {
      value: normalizedValue,
      label: formatDuration(normalizedValue, $t),
      icon: getDurationIcon(normalizedValue),
    }

    const newValues = [...(props.modelValue || []), normalizedValue]
    const sortedValues = newValues.sort((a, b) => compareDurations(a, b))
    emit('update:modelValue', sortedValues)
    emit('add-option', newOption)

    cancelAddingOption()
  }

  const removeOption = (value: string) => {
    const newValues = (props.modelValue || []).filter((v) => v !== value)
    emit('update:modelValue', newValues)
    emit('remove-option', value)
  }

  watch(
    () => props.forcedValue,
    (forcedValue) => {
      if (forcedValue) {
        if (!selectedValues.value?.includes(forcedValue)) {
          emit('update:modelValue', [...(selectedValues.value || []), forcedValue])
        }
      }
    },
    { immediate: true }
  )
</script>

<template>
  <div class="multi-selector">
    <div class="flex flex-wrap gap-1.5">
      <div v-for="option in displayOptions" :key="option.value" class="group relative">
        <input
          :id="`${inputIdPrefix}-${option.value}`"
          v-model="selectedValues"
          type="checkbox"
          :value="option.value"
          :disabled="disabled || (forcedValue && option.value === forcedValue)"
          class="sr-only"
        />
        <label
          :for="`${inputIdPrefix}-${option.value}`"
          class="relative inline-flex cursor-pointer select-none items-center gap-1 px-2 py-0.5 text-xs font-medium transition-all"
          :class="[
            selectedValues?.includes(option.value) || (forcedValue && option.value === forcedValue)
              ? 'bg-brand-primary text-content shadow-sm'
              : 'bg-hover-bg text-content-muted hover:bg-background-600',
            disabled || (forcedValue && option.value === forcedValue) ? 'cursor-not-allowed opacity-60' : '',
            sizeClasses[size],
            roundedClasses[rounded],
          ]"
        >
          <span v-if="showIcons !== false" class="flex-shrink-0">
            <slot :name="`icon-${option.value}`" :option="option">
              <i v-if="option.icon" class="text-sm" :class="[option.icon]" />
              <i v-else-if="forcedValue && option.value === forcedValue" class="text-sm text-error-500" :class="[forcedIcon]" />
              <i v-else class="text-brand-primary text-sm" :class="[getDurationIcon(option.value)]" />
            </slot>
          </span>

          <span class="leading-tight">{{ option.label || formatDuration(option.value, $t) }}</span>

          <span v-if="forcedValue && option.value === forcedValue" class="ml-1 flex-shrink-0">
            <i class="fas fa-lock text-sm text-warning-500" />
          </span>

          <button
            v-if="editable && !disabled && !(forcedValue && option.value === forcedValue)"
            class="bg-error-500/90 absolute -right-1.5 -top-1.5 flex h-4 w-4 items-center justify-center rounded-full text-content opacity-0 shadow-lg transition-all duration-300 hover:scale-110 hover:bg-error-600 group-hover:opacity-100"
            :title="$t('components.multiSelector.deleteOption')"
            @click.prevent="removeOption(option.value)"
          >
            <i
              class="fas fa-times text-xs"
              style="
                line-height: 1;
                margin: 0;
                padding: 0;
                position: absolute;
                top: 50%;
                left: 50%;
                transform: translate(-50%, -50%);
              "
            />
          </button>
        </label>
      </div>

      <div v-if="editable && !disabled && canAddMore" class="relative">
        <button
          v-if="!isAddingOption"
          class="hover:border-brand-primary hover:text-brand-primary inline-flex cursor-pointer select-none items-center gap-1 border-2 border-dashed border-subtle px-2 py-0.5 text-xs font-medium text-content-muted transition-all"
          :class="[sizeClasses[size], roundedClasses[rounded]]"
          @click="startAddingOption"
        >
          <i class="fas fa-plus text-sm" />
          <span class="leading-tight">{{ displayAddText }}</span>
        </button>

        <div v-else class="inline-flex items-center gap-1" @click.stop>
          <input
            ref="inputRef"
            v-model="inputValue"
            :placeholder="placeholderText"
            class="border border-input-border bg-input-bg px-2 py-0.5 text-xs text-input-text placeholder-placeholder focus:border-input-border-focus focus:outline-none"
            :class="[sizeClasses[size], roundedClasses[rounded]]"
            @keyup.enter="confirmAddOption"
            @keyup.escape="cancelAddingOption"
            @blur="confirmAddOption"
            @click.stop
          />
          <button
            class="bg-brand-primary rounded px-1.5 py-0.5 text-xs text-content transition-colors hover:bg-brand-600"
            :title="$t('components.multiSelector.confirmAdd')"
            @click.stop="confirmAddOption"
          >
            <i class="fas fa-check" />
          </button>
          <button
            class="rounded bg-hover-bg px-1.5 py-0.5 text-xs text-content transition-colors hover:bg-background-600"
            :title="$t('components.multiSelector.cancel')"
            @click.stop="cancelAddingOption"
          >
            <i class="fas fa-times" />
          </button>
        </div>
      </div>
    </div>

    <div v-if="inputErrorText" class="mt-1 text-xs text-error-500">
      {{ inputErrorText }}
    </div>

    <div v-if="editable && !disabled" class="mt-1 text-xs text-content-muted">
      <span>{{ $t('components.multiSelector.hints.format') }}</span>
      <span v-if="!isGuest && !hasPermanentOption" class="text-error-500">{{
        $t('components.multiSelector.hints.permanent')
      }}</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .multi-selector {
    .flex {
      align-items: flex-start;
    }
    label {
      position: relative;
      user-select: none;
      transition: all 0.2s ease-in-out;
    }
    label:hover:not(.cursor-not-allowed) {
      transform: translateY(-0.5px);
    }
    label.cursor-not-allowed {
      transform: none !important;
    }
    i {
      line-height: 1;
    }
  }
</style>
