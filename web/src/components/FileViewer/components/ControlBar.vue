<script setup lang="ts">
  import { reactive } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    isVisible: boolean
    scale: number
    shouldUseFillMode: boolean
    hasMultipleFiles: boolean
    hasPreviousFile: boolean
    hasNextFile: boolean
    currentIndex: number
    totalFiles: number
    isLightBackground?: boolean
  }

  const props = defineProps<Props>()

  defineEmits<{
    'zoom-in': []
    'zoom-out': []
    'rotate-left': []
    'rotate-right': []
    'reset-transform': []
    'toggle-fit-mode': []
    'prev-image': []
    'next-image': []
  }>()

  const { $t } = useTexts()

  /* 工具提示状态 */
  const tooltip = reactive({
    visible: false,
    text: '',
    style: {},
  })

  const showTooltip = (event: MouseEvent, text: string) => {
    const button = event.currentTarget as HTMLElement
    const rect = button.getBoundingClientRect()

    tooltip.text = text
    tooltip.visible = true

    const tooltipLeft = rect.left + rect.width / 2
    const tooltipTop = rect.top - 40 // 在按钮上方40px

    tooltip.style = {
      position: 'fixed',
      left: `${tooltipLeft}px`,
      top: `${tooltipTop}px`,
      transform: 'translateX(-50%)',
      zIndex: '10001',
    }
  }

  const hideTooltip = () => {
    tooltip.visible = false
  }
</script>

<template>
  <div class="control-bar" :class="{ 'is-hidden': !isVisible, 'is-light-bg': props.isLightBackground }">
    <div class="transform-controls">
      <button
        class="control-btn tooltip-btn"
        @click="$emit('zoom-out')"
        @mouseenter="showTooltip($event, $t('components.controlBar.zoomOut'))"
        @mouseleave="hideTooltip"
      >
        <i class="fas fa-search-minus" />
      </button>
      <div class="zoom-level">{{ Math.round((isNaN(scale) ? 1 : scale) * 100) }}%</div>
      <button
        class="control-btn tooltip-btn"
        @click="$emit('zoom-in')"
        @mouseenter="showTooltip($event, $t('components.controlBar.zoomIn'))"
        @mouseleave="hideTooltip"
      >
        <i class="fas fa-search-plus" />
      </button>
      <button
        class="control-btn tooltip-btn"
        @click="$emit('rotate-left')"
        @mouseenter="showTooltip($event, $t('components.controlBar.rotateLeft'))"
        @mouseleave="hideTooltip"
      >
        <i class="fas fa-undo" />
      </button>
      <button
        class="control-btn tooltip-btn"
        @click="$emit('rotate-right')"
        @mouseenter="showTooltip($event, $t('components.controlBar.rotateRight'))"
        @mouseleave="hideTooltip"
      >
        <i class="fas fa-redo" />
      </button>
      <button
        class="control-btn tooltip-btn"
        @click="$emit('reset-transform')"
        @mouseenter="showTooltip($event, $t('components.controlBar.reset'))"
        @mouseleave="hideTooltip"
      >
        <i class="fas fa-expand-arrows-alt" />
      </button>
      <button
        class="control-btn tooltip-btn"
        @click="$emit('toggle-fit-mode')"
        @mouseenter="
          showTooltip($event, $t(shouldUseFillMode ? 'components.controlBar.switchToFit' : 'components.controlBar.switchToFill'))
        "
        @mouseleave="hideTooltip"
      >
        <i :class="shouldUseFillMode ? 'fas fa-compress-arrows-alt' : 'fas fa-expand-arrows-alt'" />
      </button>
    </div>

    <div v-if="totalFiles > 1" class="navigation-controls">
      <button
        class="nav-btn prev-btn tooltip-btn"
        :disabled="!hasPreviousFile"
        @click="$emit('prev-image')"
        @mouseenter="showTooltip($event, $t('components.controlBar.previous'))"
        @mouseleave="hideTooltip"
      >
        <i class="fas fa-chevron-left" />
      </button>

      <div class="image-counter">{{ currentIndex + 1 }} / {{ totalFiles }}</div>

      <button
        class="nav-btn next-btn tooltip-btn"
        :disabled="!hasNextFile"
        @click="$emit('next-image')"
        @mouseenter="showTooltip($event, $t('components.controlBar.next'))"
        @mouseleave="hideTooltip"
      >
        <i class="fas fa-chevron-right" />
      </button>
    </div>

    <div v-if="tooltip.visible" class="dynamic-tooltip" :style="tooltip.style">
      {{ tooltip.text }}
      <div class="tooltip-arrow" />
    </div>
  </div>
</template>

<style scoped lang="scss" src="./ControlBar.scss"></style>
