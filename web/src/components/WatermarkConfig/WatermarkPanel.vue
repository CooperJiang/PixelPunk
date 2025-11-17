<script setup lang="ts">
  import { ref, watch } from 'vue'
  import type { WatermarkConfig, WatermarkPanelProps } from './types'
  import WatermarkTypeSelector from './components/WatermarkTypeSelector.vue'
  import WatermarkTextConfig from './components/WatermarkTextConfig.vue'
  import WatermarkFileConfig from './components/WatermarkFileConfig.vue'
  import WatermarkPositionConfig from './components/WatermarkPositionConfig.vue'
  import WatermarkStyleConfig from './components/WatermarkStyleConfig.vue'
  import WatermarkEffectsConfig from './components/WatermarkEffectsConfig.vue'

  defineOptions({
    name: 'WatermarkPanel',
  })

  const props = defineProps<WatermarkPanelProps>()

  const emit = defineEmits<{
    (e: 'config-change', config: WatermarkConfig): void
  }>()

  const localConfig = ref<WatermarkConfig>({ ...props.config })

  watch(
    localConfig,
    (newConfig) => {
      emit('config-change', { ...newConfig })
    },
    { deep: true }
  )

  watch(
    () => props.config,
    (newConfig) => {
      if (JSON.stringify(newConfig) !== JSON.stringify(localConfig.value)) {
        localConfig.value = { ...newConfig }
      }
    },
    { deep: true }
  )

  const handleConfigUpdate = (newConfig: WatermarkConfig) => {
    localConfig.value = { ...newConfig }
  }
</script>

<template>
  <div class="watermark-panel">
    <div class="panel-sections">
      <WatermarkTypeSelector :config="localConfig" @update:config="handleConfigUpdate" />

      <WatermarkTextConfig v-if="localConfig.type === 'text'" :config="localConfig" @update:config="handleConfigUpdate" />

      <WatermarkFileConfig v-if="localConfig.type === 'image'" :config="localConfig" @update:config="handleConfigUpdate" />

      <WatermarkPositionConfig :config="localConfig" @update:config="handleConfigUpdate" />

      <WatermarkStyleConfig :config="localConfig" @update:config="handleConfigUpdate" />

      <WatermarkEffectsConfig :config="localConfig" @update:config="handleConfigUpdate" />
    </div>
  </div>
</template>

<style scoped>
  .watermark-panel {
    @apply space-y-2;
  }

  .panel-sections {
    @apply space-y-3;
  }

  .panel-sections > * {
    animation: slideUp var(--transition-normal) var(--ease-out);
  }
</style>
