<script setup lang="ts">
  import { inject, computed, type Ref } from 'vue'

  defineOptions({
    name: 'CyberTabPane',
  })

  interface Props {
    name: string
    label: string
    disabled?: boolean
  }

  const props = defineProps<Props>()

  const activeTab = inject<Ref<string>>('activeTab')

  const isActive = computed(() => activeTab?.value === props.name)
</script>

<template>
  <div v-show="isActive" class="cyber-tab-pane">
    <slot />
  </div>
</template>

<style scoped lang="scss">
  .cyber-tab-pane {
    animation: fadeIn 0.3s ease-in-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
