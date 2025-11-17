<script setup lang="ts">
  import { computed } from 'vue'
  import type { RadioGroupProps, RadioGroupEmits } from './types'

  defineOptions({
    name: 'CyberRadioGroup',
  })

  const props = withDefaults(defineProps<RadioGroupProps>(), {
    disabled: false,
    layout: 'horizontal',
  })

  const emit = defineEmits<RadioGroupEmits>()

  const localValue = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })
</script>

<template>
  <div class="cyber-radio-group" :class="[`layout-${props.layout}`]">
    <CyberRadio
      v-for="option in props.options"
      :key="option.value"
      v-model="localValue"
      :value="option.value"
      :disabled="props.disabled"
    >
      {{ option.label }}
    </CyberRadio>
  </div>
</template>

<style scoped>
  .cyber-radio-group {
    @apply flex gap-4;
  }

  .cyber-radio-group.layout-vertical {
    @apply flex-col gap-2;
  }

  .cyber-radio-group.layout-horizontal {
    @apply flex-row flex-wrap;
  }
</style>
