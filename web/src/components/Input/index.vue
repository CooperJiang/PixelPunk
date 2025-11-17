<template>
  <div
    class="cyber-input-wrapper"
    :class="[
      {
        'cyber-input--disabled': disabled,
        'cyber-input--error': error === true || !!errorMessage,
        'cyber-input--with-prepend': $slots.prepend,
        'cyber-input--with-prefix': $slots.prefix || prefixIcon,
        'cyber-input--with-suffix': $slots.suffix || suffixIcon || clearable,
        'cyber-input--with-append': $slots.append,
        'cyber-input--with-unit': hasUnitSlot,
        'cyber-input--focused': focused,
        'cyber-input--no3d': no3d,
      },
    ]"
    :style="wrapperStyle"
  >
    <div v-if="$slots.prepend" class="cyber-input__prepend">
      <slot name="prepend" />
    </div>

    <div class="cyber-input__inner">
      <div v-if="$slots.prefix || prefixIcon" class="cyber-input__prefix">
        <slot name="prefix">
          <i v-if="prefixIcon" :class="`fas fa-${prefixIcon}`" />
        </slot>
      </div>

      <textarea
        v-if="type === 'textarea'"
        :id="inputId || undefined"
        ref="inputRef"
        :value="modelValue"
        class="cyber-textarea"
        :class="{ 'has-error': error === true || !!errorMessage }"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :rows="rows"
        :maxlength="maxlength"
        :name="name"
        v-bind="filteredAttrs"
        @input="onInput"
        @change="onChange"
        @focus="onFocus"
        @blur="onBlur"
        @keyup="onKeyup"
        @keydown="onKeydown"
      />

      <input
        v-else
        :id="inputId || undefined"
        ref="inputRef"
        :value="modelValue"
        :type="inputType"
        class="cyber-input"
        :class="{ 'has-error': error || !!errorMessage }"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :maxlength="maxlength"
        :name="name"
        :autocomplete="autocomplete"
        v-bind="filteredAttrs"
        @input="onInput"
        @change="onChange"
        @focus="onFocus"
        @blur="onBlur"
        @keyup="onKeyup"
        @keydown="onKeydown"
      />

      <div v-if="$slots.suffix || suffixIcon || clearable" class="cyber-input__suffix">
        <i v-if="clearable && modelValue && !disabled" class="cyber-input__clear fas fa-times-circle" @click.stop="onClear" />
        <slot name="suffix">
          <i v-if="suffixIcon" :class="`fas fa-${suffixIcon}`" />
        </slot>
      </div>

      <div v-if="$slots.unit" class="cyber-input__unit">
        <slot name="unit" />
      </div>

      <div v-if="type === 'password'" class="cyber-input__password-toggle">
        <i :class="`fas fa-${showPassword ? 'eye-slash' : 'eye'}`" @click.stop="togglePasswordVisibility" />
      </div>
    </div>

    <div v-if="$slots.append" class="cyber-input__append">
      <slot name="append" />
    </div>
  </div>

  <div v-if="errorMessage" class="cyber-input__error-message">
    {{ errorMessage }}
  </div>
</template>

<script setup lang="ts">
  import { computed, onMounted, ref, useAttrs, useSlots } from 'vue'
  import type { InputEmits, InputProps } from './types'

  defineOptions({
    name: 'CyberInput',
    inheritAttrs: false,
  })

  const props = withDefaults(defineProps<InputProps>(), {
    modelValue: '',
    type: 'text',
    placeholder: '',
    disabled: false,
    readonly: false,
    clearable: false,
    autofocus: false,
    autocomplete: 'off',
    name: '',
    rows: 3,
    error: false,
    errorMessage: '',
    prefixIcon: '',
    suffixIcon: '',
    inputId: '',
    width: undefined,
    height: undefined,
    no3d: false,
  })

  const emit = defineEmits<InputEmits>()

  const inputRef = ref<HTMLInputElement | HTMLTextAreaElement | null>(null)
  const focused = ref(false)
  const showPassword = ref(false)
  const slots = useSlots()
  const attrs = useAttrs()

  const filteredAttrs = computed(() => {
    const filtered = { ...attrs }
    delete filtered.class
    delete filtered.style
    delete filtered.size // Remove size prop to prevent HTML input size attribute warning
    return filtered
  })

  const hasUnitSlot = computed(() => Boolean(slots.unit))

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

  const inputType = computed(() => {
    if (props.type === 'password') {
      return showPassword.value ? 'text' : 'password'
    }
    return props.type === 'textarea' ? 'text' : props.type
  })

  const togglePasswordVisibility = () => {
    showPassword.value = !showPassword.value
  }

  const onInput = (e: Event) => {
    const target = e.target as HTMLInputElement | HTMLTextAreaElement
    const { value } = target

    if (props.type === 'number') {
      if (value === '') {
        emit('update:modelValue', '')
      } else {
        const numberValue = parseFloat(value)
        if (!isNaN(numberValue)) {
          emit('update:modelValue', numberValue)
        }
      }
    } else {
      emit('update:modelValue', value)
    }

    emit('input', e)
  }

  const onChange = (e: Event) => {
    emit('change', e)
  }

  const onFocus = (e: FocusEvent) => {
    focused.value = true
    emit('focus', e)
  }

  const onBlur = (e: FocusEvent) => {
    focused.value = false
    emit('blur', e)
  }

  const onKeyup = (e: KeyboardEvent) => {
    if (e.key === 'Enter') {
      emit('enter', e)
    }
    emit('keyup', e)
  }

  const onKeydown = (e: KeyboardEvent) => {
    emit('keydown', e)
  }

  const onClear = () => {
    emit('update:modelValue', '')
    emit('clear')
    inputRef.value?.focus()
  }

  const focus = () => {
    inputRef.value?.focus()
  }

  const blur = () => {
    inputRef.value?.blur()
  }

  const select = () => {
    inputRef.value?.select()
  }

  defineExpose({
    focus,
    blur,
    select,
    inputRef,
  })

  onMounted(() => {
    if (props.autofocus && !props.disabled) {
      focus()
    }
  })
</script>

<style scoped>
  .cyber-input-wrapper {
    position: relative;
    display: flex;
    width: 100%;
    align-items: stretch;
    overflow: hidden;
    border-radius: var(--radius-sm);
    border: 1.5px solid var(--color-border-default);
    background: var(--color-background-700);
    transition:
      border-color 0.25s ease,
      box-shadow 0.25s ease,
      transform 0.25s ease;
    box-sizing: border-box;
    min-height: 32px;
    box-shadow:
      2px 2.5px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 0 rgba(var(--color-brand-500-rgb), 0);
  }

  .cyber-input__inner {
    display: flex;
    min-width: 0;
    flex: 1;
    align-items: center;
  }

  .cyber-input,
  .cyber-textarea {
    min-width: 0;
    flex: 1;
    border: none;
    background: transparent;
    padding: 0 12px;
    font-size: 0.875rem;
    outline: none;
    color: var(--color-content);
    font-family: inherit;
    box-sizing: border-box;
  }

  .cyber-input__prepend,
  .cyber-input__append {
    display: flex;
    flex-shrink: 0;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    padding: 0;
    font-size: 0.875rem;
    color: var(--color-content-muted);
  }

  .cyber-input__prepend {
    border-radius: var(--radius-sm) 0 0 var(--radius-sm);
  }

  .cyber-input__append {
    border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  }

  .cyber-input__prefix {
    pointer-events: none;
    display: flex;
    flex-shrink: 0;
    align-items: center;
    justify-content: center;
    padding: 0 6px;
    color: var(--color-content-muted);
  }

  .cyber-input__suffix,
  .cyber-input__unit {
    pointer-events: none;
    display: flex;
    flex-shrink: 0;
    align-items: center;
    justify-content: center;
    padding: 0 6px;
    color: var(--color-content-muted);
  }

  .cyber-input__unit {
    padding-right: 8px;
    font-size: 0.75rem;
    font-weight: normal;
    color: var(--color-content-muted);
  }

  .cyber-input__password-toggle {
    pointer-events: auto;
    display: flex;
    flex-shrink: 0;
    cursor: pointer;
    align-items: center;
    justify-content: center;
    padding: 0 6px;
    color: var(--color-content-muted);
    transition: color 0.2s ease;
  }

  .cyber-input__password-toggle:hover {
    color: var(--color-brand-500);
  }

  .cyber-input__clear {
    pointer-events: auto;
    margin-right: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .cyber-input__clear:hover {
    color: var(--color-error-500);
  }

  .cyber-input[type='number']::-webkit-inner-spin-button,
  .cyber-input[type='number']::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  .cyber-input[type='number'] {
    -moz-appearance: textfield;
  }

  .cyber-input::placeholder,
  .cyber-textarea::placeholder {
    color: var(--color-content-muted) !important;
    font-style: italic;
    font-weight: 400;
  }

  .cyber-input::-webkit-input-placeholder,
  .cyber-textarea::-webkit-input-placeholder {
    color: var(--color-content-muted) !important;
    font-style: italic;
    font-weight: 400;
  }

  .cyber-input::-moz-placeholder,
  .cyber-textarea::-moz-placeholder {
    color: var(--color-content-muted) !important;
    font-style: italic;
    font-weight: 400;
    opacity: 1;
  }

  .cyber-input:-ms-input-placeholder,
  .cyber-textarea:-ms-input-placeholder {
    color: var(--color-content-muted) !important;
    font-style: italic;
    font-weight: 400;
  }

  .cyber-input-wrapper:hover:not(.cyber-input--disabled):not(.cyber-input--focused) {
    border-color: var(--color-border-strong);

    box-shadow:
      2.5px 3px 0 rgba(var(--color-brand-500-rgb), 0.3),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.08);
    transform: translate(-0.3px, -0.3px);
  }

  .cyber-input-wrapper:focus-within,
  .cyber-input--focused {
    border-color: var(--color-brand-500);

    box-shadow:
      3.5px 4px 0 rgba(var(--color-brand-500-rgb), 0.4),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.2);
    transform: translate(-0.5px, -0.5px);
  }

  .cyber-input-wrapper:active:not(.cyber-input--disabled) {
    transform: translate(0.5px, 0.5px);
    box-shadow:
      1.5px 2px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-input--error {
    border-color: var(--color-error-500);

    box-shadow:
      2px 2.5px 0 rgba(var(--color-error-rgb), 0.25),
      0 0 6px rgba(var(--color-error-rgb), 0.15);
  }

  .cyber-input--error:hover:not(.cyber-input--disabled) {
    box-shadow:
      2.5px 3px 0 rgba(var(--color-error-rgb), 0.3),
      0 0 8px rgba(var(--color-error-rgb), 0.2);
    transform: translate(-0.3px, -0.3px);
  }

  .cyber-input--error:focus-within {
    box-shadow:
      3.5px 4px 0 rgba(var(--color-error-rgb), 0.4),
      0 0 12px rgba(var(--color-error-rgb), 0.25);
    transform: translate(-0.5px, -0.5px);
  }

  .cyber-input--disabled {
    background: rgba(var(--color-background-800-rgb), 0.3);
    border-color: var(--color-border-default);
    opacity: 0.6;
    cursor: not-allowed;

    box-shadow: none !important;
    transform: none !important;
  }

  .cyber-input--no3d,
  .cyber-input--no3d:hover,
  .cyber-input--no3d:focus-within,
  .cyber-input--no3d:active {
    box-shadow: none !important;
    transform: none !important;
  }

  .cyber-input--no3d:focus-within {
    border-color: var(--color-brand-500);
    box-shadow: 0 0 0 2px rgba(var(--color-brand-500-rgb), 0.1) !important;
  }

  .cyber-input--disabled .cyber-input,
  .cyber-input--disabled .cyber-textarea {
    color: var(--color-content-disabled);
    cursor: not-allowed;
  }

  .cyber-input--readonly {
    background: rgba(var(--color-brand-500), 0.02);
    border-color: rgba(var(--color-brand-500), 0.3);
  }

  .cyber-input--readonly .cyber-input,
  .cyber-input--readonly .cyber-textarea {
    cursor: default;
  }

  .cyber-textarea {
    resize: vertical;
    min-height: 80px;
    padding: 8px 12px;
    line-height: 1.4;
    transition: all 0.25s ease;
  }

  .cyber-input-wrapper:has(.cyber-textarea) {
    min-height: 80px;
  }

  .cyber-input-wrapper:has(.cyber-textarea):focus-within {
    box-shadow:
      6px 7px 0 rgba(var(--color-brand-500-rgb), 0.5),
      0 0 25px rgba(var(--color-brand-500-rgb), 0.35);
    transform: translate(-1.5px, -1.5px);
  }

  .cyber-input__error-message {
    margin-top: 4px;
    font-size: 0.8125rem;
    color: var(--color-error-500);
  }

  .cyber-input--error:not(:focus-within) {
    border-color: var(--color-error-500);
  }

  .cyber-input--with-prepend .cyber-input__inner {
    border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  }

  .cyber-input--with-append .cyber-input__inner {
    border-radius: var(--radius-sm) 0 0 var(--radius-sm);
  }

  @media (max-width: 768px) {
    .cyber-input,
    .cyber-textarea {
      padding: 4px 8px;
      font-size: 0.8125rem;
    }

    .cyber-input-wrapper {
      min-height: 30px;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .cyber-input-wrapper,
    .cyber-input__clear,
    .cyber-input__password-toggle {
      transition: none !important;
    }
  }
</style>
