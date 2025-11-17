<script setup lang="ts">
  import { defineProps, defineEmits } from 'vue'

  type SortOption = { value: string; label: string; icon: string }
  type Density = 'compact' | 'normal' | 'comfortable'

  const _props = defineProps<{
    filtersSort: string
    sortOptions: SortOption[]
    currentDensity: Density
    densityOptions: { value: Density; label: string; icon: string }[]
  }>()

  const emit = defineEmits<{
    (e: 'sort', value: string): void
    (e: 'change-density', value: Density): void
  }>()

  function onSort(value: string) {
    emit('sort', value)
  }

  function setDensity(density: Density) {
    emit('change-density', density)
  }
</script>

<template>
  <div class="toolbar-group">
    <div class="sort-switcher">
      <button
        v-for="sort in sortOptions"
        :key="sort.value"
        class="sort-btn"
        :class="{ active: filtersSort === sort.value }"
        :title="sort.label"
        @click="onSort(sort.value)"
      >
        <i :class="sort.icon" />
      </button>
    </div>

    <div class="density-switcher">
      <button
        v-for="density in densityOptions"
        :key="density.value"
        class="density-btn"
        :class="{ active: currentDensity === density.value }"
        :title="density.label"
        @click="setDensity(density.value)"
      >
        <i :class="density.icon" />
      </button>
    </div>
  </div>
</template>

<style scoped></style>
