<script setup lang="ts">
  import { defineProps, defineEmits } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const _props = defineProps<{
    isAllSelected: boolean
    selectedCount: number
  }>()

  const emit = defineEmits<{
    (e: 'toggle-all'): void
    (e: 'invert'): void
  }>()

  function onToggleAll() {
    emit('toggle-all')
  }

  function onInvert() {
    emit('invert')
  }
</script>

<template>
  <div class="selection-toolbar-compact">
    <div class="flex items-center gap-2">
      <button
        class="select-btn select-all-btn"
        :class="{ active: isAllSelected }"
        :title="isAllSelected ? $t('explore.selection.deselectAll') : $t('explore.selection.selectAllPage')"
        @click="onToggleAll"
      >
        <i class="fas" :class="isAllSelected ? 'fa-check-double' : 'fa-check'" />
        {{ isAllSelected ? $t('explore.selection.deselectAll') : $t('explore.selection.selectAllPage') }}
      </button>
      <button class="select-btn invert-btn" :title="$t('explore.selection.invertPage')" @click="onInvert">
        <i class="fas fa-exchange-alt" />
        {{ $t('explore.selection.invertPage') }}
      </button>
    </div>

    <div class="text-xs text-content-muted">{{ $t('explore.selection.selected', { count: selectedCount }) }}</div>
  </div>
</template>

<style scoped>
  .selection-toolbar-compact {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    flex: 1;
    min-width: 0;
  }
</style>
