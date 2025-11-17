<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'SelectionToolbar',
  })

  defineProps<Props>()

  defineEmits<{
    'toggle-select-all': []
    'select-invert': []
  }>()

  interface Props {
    selectMode: boolean
    isAllSelected: boolean
  }
</script>

<template>
  <div v-if="selectMode" class="selection-toolbar">
    <div class="selection-buttons">
      <button
        class="select-btn select-all-btn"
        :class="{ active: isAllSelected }"
        :title="isAllSelected ? $t('admin.files.batch.deselectAll') : $t('admin.files.batch.selectAll')"
        @click="$emit('toggle-select-all')"
      >
        <i class="fas" :class="isAllSelected ? 'fa-check-double' : 'fa-check'" />
        {{ isAllSelected ? $t('admin.files.batch.deselectAll') : $t('admin.files.batch.selectAll') }}
      </button>

      <button class="select-btn invert-btn" :title="$t('admin.files.batch.invertSelection')" @click="$emit('select-invert')">
        <i class="fas fa-exchange-alt" />
        {{ $t('admin.files.batch.invertSelection') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
  .selection-toolbar {
    @apply mb-4 border-b border-subtle pb-3;
  }

  .selection-buttons {
    @apply flex gap-2;
  }

  .select-btn {
    @apply flex items-center gap-1.5;
    @apply rounded border px-3 py-1.5;
    @apply border-default bg-background-700 text-content;
    @apply text-xs font-medium;
    @apply cursor-pointer;
    @apply transition-all duration-200;
  }

  .select-btn:hover {
    @apply border-brand-500 bg-background-600 text-brand-400;
  }

  .select-btn.active,
  .select-all-btn.active {
    @apply border-brand-500 bg-background-600 text-brand-400;
  }

  .select-btn i {
    @apply text-xs transition-transform duration-200;
  }

  .select-btn:hover i {
    @apply scale-110;
  }

  @media (max-width: 768px) {
    .selection-buttons {
      @apply flex-wrap justify-center;
    }

    .select-btn {
      @apply flex-1 justify-center;
      min-width: 120px;
    }
  }
</style>
