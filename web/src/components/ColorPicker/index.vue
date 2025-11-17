<template>
  <div class="cyber-color-picker">
    <div class="color-input-wrapper" @click="togglePicker">
      <div class="color-preview" :style="{ backgroundColor: modelValue }"></div>
      <input :value="modelValue" type="text" class="color-text" readonly :placeholder="placeholder" />
      <i class="fas fa-palette color-icon" />
    </div>

    <Teleport to="body">
      <Transition name="color-picker-fade">
        <div v-if="showPicker" class="color-picker-overlay" @click="closePicker">
          <div ref="panelRef" class="color-picker-panel" :style="panelStyle" @click.stop>
            <div class="picker-body-compact">
              <div class="color-canvas-wrapper">
                <canvas
                  ref="canvasRef"
                  class="color-canvas"
                  width="220"
                  height="140"
                  @mousedown="handleCanvasMouseDown"
                  @mousemove="handleCanvasMouseMove"
                  @mouseup="handleCanvasMouseUp"
                  @mouseleave="handleCanvasMouseUp"
                />
                <div
                  class="color-cursor"
                  :style="{
                    left: `${cursorPosition.x}px`,
                    top: `${cursorPosition.y}px`,
                  }"
                />
              </div>

              <div class="hue-slider-wrapper">
                <div class="hue-slider" @mousedown="handleHueMouseDown">
                  <div class="hue-cursor" :style="{ left: `${huePosition}%` }" />
                </div>
              </div>

              <div class="preset-colors-compact">
                <div
                  v-for="color in presetColors"
                  :key="color"
                  class="preset-color-compact"
                  :class="{ active: modelValue.toLowerCase() === color.toLowerCase() }"
                  :style="{ backgroundColor: color }"
                  @click="selectPresetColor(color)"
                />
              </div>

              <div class="color-input-row-compact">
                <div class="current-color-preview-compact" :style="{ backgroundColor: currentColor }"></div>
                <input
                  v-model="hexInput"
                  type="text"
                  class="hex-input-compact"
                  placeholder="#000000"
                  @input="handleHexInput"
                  @blur="validateHexInput"
                />
              </div>

              <div class="picker-footer-compact">
                <CyberButton size="sm" type="secondary" @click="closePicker">
                  <i class="fas fa-times mr-1.5" />
                  {{ $t('actions.cancel') }}
                </CyberButton>
                <CyberButton size="sm" type="primary" @click="confirmColor">
                  <i class="fas fa-check mr-1.5" />
                  {{ $t('actions.confirm') }}
                </CyberButton>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
  import { ref, watch, onMounted, nextTick } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  interface Props {
    modelValue: string
    placeholder?: string
  }

  interface Emits {
    (e: 'update:modelValue', value: string): void
  }

  const props = withDefaults(defineProps<Props>(), {
    placeholder: '#000000',
  })

  const emit = defineEmits<Emits>()

  const showPicker = ref(false)
  const canvasRef = ref<HTMLCanvasElement | null>(null)
  const panelRef = ref<HTMLDivElement | null>(null)
  const currentColor = ref(props.modelValue)
  const hexInput = ref(props.modelValue)
  const isDraggingCanvas = ref(false)
  const isDraggingHue = ref(false)
  const panelStyle = ref<Record<string, string>>({})

  /* HSV 值 */
  const hue = ref(0)
  const saturation = ref(100)
  const value = ref(100)

  const cursorPosition = ref({ x: 280, y: 0 })
  const huePosition = ref(0)

  const presetColors = [
    '#000000',
    '#FFFFFF',
    '#FF0000',
    '#00FF00',
    '#0000FF',
    '#FFFF00',
    '#FF00FF',
    '#00FFFF',
    '#FF6B6B',
    '#4ECDC4',
    '#45B7D1',
    '#FFA07A',
    '#98D8C8',
    '#F7DC6F',
    '#BB8FCE',
    '#85C1E2',
  ]

  onMounted(() => {
    initializeFromHex(props.modelValue)
    drawColorCanvas()
  })

  watch(
    () => props.modelValue,
    (newVal) => {
      currentColor.value = newVal
      hexInput.value = newVal
      initializeFromHex(newVal)
    }
  )

  watch(hue, () => {
    drawColorCanvas()
  })

  function togglePicker(event: MouseEvent) {
    showPicker.value = !showPicker.value
    if (showPicker.value) {
      nextTick(() => {
        calculatePanelPosition(event)
        drawColorCanvas()
      })
    }
  }

  function calculatePanelPosition(event: MouseEvent) {
    const target = event.currentTarget as HTMLElement
    const rect = target.getBoundingClientRect()
    const panelWidth = 260
    const panelHeight = 380
    const padding = 8

    let left = rect.left
    let top = rect.bottom + padding

    if (left + panelWidth > window.innerWidth) {
      left = window.innerWidth - panelWidth - padding
    }

    if (left < padding) {
      left = padding
    }

    if (top + panelHeight > window.innerHeight) {
      top = rect.top - panelHeight - padding
    }

    if (top < padding) {
      top = padding
    }

    panelStyle.value = {
      left: `${left}px`,
      top: `${top}px`,
    }
  }

  function closePicker() {
    showPicker.value = false
  }

  function confirmColor() {
    emit('update:modelValue', currentColor.value)
    closePicker()
  }

  function selectPresetColor(color: string) {
    currentColor.value = color
    hexInput.value = color
    initializeFromHex(color)
  }

  function drawColorCanvas() {
    if (!canvasRef.value) return
    const canvas = canvasRef.value
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const width = canvas.width
    const height = canvas.height

    ctx.fillStyle = `hsl(${hue.value}, 100%, 50%)`
    ctx.fillRect(0, 0, width, height)

    const whiteGradient = ctx.createLinearGradient(0, 0, width, 0)
    whiteGradient.addColorStop(0, 'rgba(255, 255, 255, 1)')
    whiteGradient.addColorStop(1, 'rgba(255, 255, 255, 0)')
    ctx.fillStyle = whiteGradient
    ctx.fillRect(0, 0, width, height)

    const blackGradient = ctx.createLinearGradient(0, 0, 0, height)
    blackGradient.addColorStop(0, 'rgba(0, 0, 0, 0)')
    blackGradient.addColorStop(1, 'rgba(0, 0, 0, 1)')
    ctx.fillStyle = blackGradient
    ctx.fillRect(0, 0, width, height)
  }

  function handleCanvasMouseDown(e: MouseEvent) {
    isDraggingCanvas.value = true
    updateColorFromCanvas(e)
  }

  function handleCanvasMouseMove(e: MouseEvent) {
    if (isDraggingCanvas.value) {
      updateColorFromCanvas(e)
    }
  }

  function handleCanvasMouseUp() {
    isDraggingCanvas.value = false
  }

  function updateColorFromCanvas(e: MouseEvent) {
    if (!canvasRef.value) return
    const rect = canvasRef.value.getBoundingClientRect()
    const x = Math.max(0, Math.min(e.clientX - rect.left, rect.width))
    const y = Math.max(0, Math.min(e.clientY - rect.top, rect.height))

    cursorPosition.value = { x, y }

    saturation.value = (x / rect.width) * 100
    value.value = 100 - (y / rect.height) * 100

    updateColorFromHSV()
  }

  function handleHueMouseDown(e: MouseEvent) {
    isDraggingHue.value = true
    updateHueFromSlider(e)

    const handleMouseMove = (e: MouseEvent) => {
      if (isDraggingHue.value) {
        updateHueFromSlider(e)
      }
    }

    const handleMouseUp = () => {
      isDraggingHue.value = false
      document.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseup', handleMouseUp)
    }

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
  }

  function updateHueFromSlider(e: MouseEvent) {
    const target = e.currentTarget as HTMLElement
    const rect = target.getBoundingClientRect()
    const x = Math.max(0, Math.min(e.clientX - rect.left, rect.width))
    const percentage = (x / rect.width) * 100

    huePosition.value = percentage
    hue.value = (percentage / 100) * 360

    updateColorFromHSV()
  }

  function hsvToHex(h: number, s: number, v: number): string {
    const rgb = hsvToRgb(h, s, v)
    return rgbToHex(rgb.r, rgb.g, rgb.b)
  }

  function hsvToRgb(h: number, s: number, v: number) {
    s = s / 100
    v = v / 100
    const c = v * s
    const x = c * (1 - Math.abs(((h / 60) % 2) - 1))
    const m = v - c

    let r = 0,
      g = 0,
      b = 0

    if (h >= 0 && h < 60) {
      r = c
      g = x
      b = 0
    } else if (h >= 60 && h < 120) {
      r = x
      g = c
      b = 0
    } else if (h >= 120 && h < 180) {
      r = 0
      g = c
      b = x
    } else if (h >= 180 && h < 240) {
      r = 0
      g = x
      b = c
    } else if (h >= 240 && h < 300) {
      r = x
      g = 0
      b = c
    } else if (h >= 300 && h < 360) {
      r = c
      g = 0
      b = x
    }

    return {
      r: Math.round((r + m) * 255),
      g: Math.round((g + m) * 255),
      b: Math.round((b + m) * 255),
    }
  }

  function rgbToHex(r: number, g: number, b: number): string {
    return `#${[r, g, b].map((x) => x.toString(16).padStart(2, '0')).join('')}`
  }

  function hexToHsv(hex: string) {
    const rgb = hexToRgb(hex)
    if (!rgb) return { h: 0, s: 0, v: 0 }

    const r = rgb.r / 255
    const g = rgb.g / 255
    const b = rgb.b / 255

    const max = Math.max(r, g, b)
    const min = Math.min(r, g, b)
    const delta = max - min

    let h = 0
    const s = max === 0 ? 0 : delta / max
    const v = max

    if (delta !== 0) {
      if (max === r) {
        h = ((g - b) / delta + (g < b ? 6 : 0)) / 6
      } else if (max === g) {
        h = ((b - r) / delta + 2) / 6
      } else {
        h = ((r - g) / delta + 4) / 6
      }
    }

    return {
      h: h * 360,
      s: s * 100,
      v: v * 100,
    }
  }

  function hexToRgb(hex: string) {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
    return result
      ? {
          r: parseInt(result[1], 16),
          g: parseInt(result[2], 16),
          b: parseInt(result[3], 16),
        }
      : null
  }

  function initializeFromHex(hex: string) {
    const hsv = hexToHsv(hex)
    hue.value = hsv.h
    saturation.value = hsv.s
    value.value = hsv.v

    huePosition.value = (hue.value / 360) * 100

    if (canvasRef.value) {
      cursorPosition.value = {
        x: (saturation.value / 100) * canvasRef.value.width,
        y: ((100 - value.value) / 100) * canvasRef.value.height,
      }
    }
  }

  function updateColorFromHSV() {
    currentColor.value = hsvToHex(hue.value, saturation.value, value.value)
    hexInput.value = currentColor.value
  }

  function handleHexInput(e: Event) {
    const input = (e.target as HTMLInputElement).value
    if (/^#[0-9A-Fa-f]{6}$/.test(input)) {
      currentColor.value = input
      initializeFromHex(input)
    }
  }

  function validateHexInput() {
    if (!/^#[0-9A-Fa-f]{6}$/.test(hexInput.value)) {
      hexInput.value = currentColor.value
    }
  }
</script>

<style scoped>
  .cyber-color-picker {
    @apply relative;
  }

  .color-input-wrapper {
    @apply flex cursor-pointer items-center gap-2 rounded-lg border-2 px-3 py-2 transition-all duration-200;
    border-color: var(--color-border-default);
    background: var(--color-background-800);
  }

  .color-input-wrapper:hover {
    border-color: var(--color-brand-500);
    background: var(--color-hover-bg);
  }

  .color-preview {
    @apply h-6 w-6 flex-shrink-0 rounded border-2;
    border-color: var(--color-border-subtle);
  }

  .color-text {
    @apply flex-1 bg-transparent font-mono text-sm outline-none;
    color: var(--color-content-default);
  }

  .color-icon {
    @apply flex-shrink-0 text-sm;
    color: var(--color-brand-500);
  }

  .color-picker-overlay {
    @apply fixed inset-0 z-[9999];
    background: transparent;
  }

  .color-picker-panel {
    @apply fixed rounded-lg border shadow-2xl;
    border-color: var(--color-border-default);
    background: var(--color-background-800);
    width: 260px;
    max-width: 90vw;
  }

  .picker-body-compact {
    @apply space-y-2.5 p-3;
  }

  .color-canvas-wrapper {
    @apply relative overflow-hidden rounded-lg;
    border: 2px solid var(--color-border-default);
  }

  .color-canvas {
    @apply block w-full cursor-crosshair;
  }

  .color-cursor {
    @apply pointer-events-none absolute h-4 w-4 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-white;
    box-shadow:
      0 0 0 1px rgba(0, 0, 0, 0.3),
      inset 0 0 0 1px rgba(0, 0, 0, 0.3);
  }

  .hue-slider-wrapper {
    @apply relative;
  }

  .hue-slider {
    @apply relative h-2.5 cursor-pointer rounded-full;
    background: linear-gradient(
      to right,
      #ff0000 0%,
      #ffff00 17%,
      #00ff00 33%,
      #00ffff 50%,
      #0000ff 67%,
      #ff00ff 83%,
      #ff0000 100%
    );
    border: 1px solid var(--color-border-default);
  }

  .hue-cursor {
    @apply pointer-events-none absolute top-1/2 h-4 w-4 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-white;
    box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.3);
  }

  .preset-colors-compact {
    @apply grid grid-cols-8 gap-1.5;
  }

  .preset-color-compact {
    @apply h-6 w-6 cursor-pointer rounded border transition-all duration-200;
    border-color: var(--color-border-default);
  }

  .preset-color-compact:hover {
    border-color: var(--color-brand-500);
    transform: scale(1.1);
  }

  .preset-color-compact.active {
    border-color: var(--color-brand-500);
    box-shadow: 0 0 0 2px var(--color-brand-500);
  }

  .color-input-row-compact {
    @apply flex items-center gap-2;
  }

  .current-color-preview-compact {
    @apply h-7 w-7 flex-shrink-0 rounded border;
    border-color: var(--color-border-default);
  }

  .hex-input-compact {
    @apply flex-1 rounded border px-2 py-1 font-mono text-xs outline-none transition-all duration-200;
    border-color: var(--color-border-default);
    background: var(--color-background-700);
    color: var(--color-content-default);
  }

  .hex-input-compact:focus {
    border-color: var(--color-brand-500);
  }

  .picker-footer-compact {
    @apply flex gap-2 border-t px-3 py-2;
    border-color: var(--color-border-subtle);
  }

  .color-picker-fade-enter-active,
  .color-picker-fade-leave-active {
    transition: opacity 0.2s ease;
  }

  .color-picker-fade-enter-from,
  .color-picker-fade-leave-to {
    opacity: 0;
  }

  .color-picker-fade-enter-active .color-picker-panel {
    animation: color-picker-scale-in 0.2s ease;
  }

  .color-picker-fade-leave-active .color-picker-panel {
    animation: color-picker-scale-out 0.2s ease;
  }

  @keyframes color-picker-scale-in {
    from {
      transform: scale(0.95);
      opacity: 0;
    }
    to {
      transform: scale(1);
      opacity: 1;
    }
  }

  @keyframes color-picker-scale-out {
    from {
      transform: scale(1);
      opacity: 1;
    }
    to {
      transform: scale(0.95);
      opacity: 0;
    }
  }
</style>
