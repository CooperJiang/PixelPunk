<script setup lang="ts">
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { useUploadStore } from '@/store/upload'
  import { useAuthStore } from '@/store/auth'
  import { StorageUtil } from '@/utils/storage/storage'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'UploadFloatButton',
  })

  const { $t } = useTexts()

  const router = useRouter()
  const route = useRoute()
  const uploadStore = useUploadStore()
  const authStore = useAuthStore()
  const fileInput = ref<HTMLInputElement>()

  /* 拖动相关状态 */
  const isDragging = ref(false)
  const dragOffset = ref({ x: 0, y: 0 })
  const currentPosition = ref({ x: 0, y: 0 })
  const dragStartTime = ref(0)

  const getSavedPosition = () => StorageUtil.get<{ x: number; y: number }>('upload-float-button-position')

  const savePosition = (x: number, y: number) => {
    StorageUtil.set('upload-float-button-position', { x, y })
  }

  const initializePosition = () => {
    const savedPos = getSavedPosition()
    const { innerWidth: windowWidth, innerHeight: windowHeight } = window
    const elementSize = 70

    let x, y

    if (savedPos && savedPos.x !== 0 && savedPos.y !== 0) {
      const isInViewport =
        savedPos.x >= 0 && savedPos.y >= 0 && savedPos.x + elementSize <= windowWidth && savedPos.y + elementSize <= windowHeight

      if (isInViewport) {
        x = savedPos.x
        y = savedPos.y
      } else {
        x = windowWidth - elementSize - 100
        y = windowHeight - elementSize - 150
      }
    } else {
      x = windowWidth - elementSize - 100
      y = windowHeight - elementSize - 150
    }

    currentPosition.value.x = Math.max(0, Math.min(x, windowWidth - elementSize))
    currentPosition.value.y = Math.max(0, Math.min(y, windowHeight - elementSize))
  }

  const floatStyle = computed(() => ({
    position: 'fixed',
    zIndex: 1000,
    pointerEvents: 'auto',
    userSelect: 'none',
    cursor: isDragging.value ? 'grabbing' : 'grab',
    left: `${currentPosition.value.x}px`,
    top: `${currentPosition.value.y}px`,
    transform: 'none',
  }))

  const handleMouseDown = (e: MouseEvent) => {
    dragStartTime.value = Date.now()
    const element = e.currentTarget as HTMLElement
    const rect = element.getBoundingClientRect()

    dragOffset.value = { x: e.clientX - rect.left, y: e.clientY - rect.top }
    currentPosition.value = { x: rect.left, y: rect.top }
    isDragging.value = true

    document.addEventListener('mousemove', handleMouseMove, { passive: false })
    document.addEventListener('mouseup', handleMouseUp, { passive: false })

    e.preventDefault()
    e.stopPropagation()
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) {
      return
    }

    const elementSize = 70
    const maxX = window.innerWidth - elementSize
    const maxY = window.innerHeight - elementSize

    currentPosition.value.x = Math.max(0, Math.min(e.clientX - dragOffset.value.x, maxX))
    currentPosition.value.y = Math.max(0, Math.min(e.clientY - dragOffset.value.y, maxY))

    e.preventDefault()
  }

  const handleMouseUp = (e: MouseEvent) => {
    if (!isDragging.value) {
      return
    }

    isDragging.value = false
    savePosition(currentPosition.value.x, currentPosition.value.y)

    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)

    e.preventDefault()
    e.stopPropagation()
  }

  const handleClick = () => {
    const dragDuration = Date.now() - dragStartTime.value
    if (dragDuration < 200) {
      fileInput.value?.click()
    }
  }

  const handleFileSelect = async (event: Event) => {
    const { files } = event.target as HTMLInputElement
    if (files && files.length > 0) {
      await uploadStore.addFiles(files)
      router.push('/upload')
    }
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }

  const isDragOver = ref(false)

  const handleDragOver = (e: DragEvent) => {
    if (e.dataTransfer?.types.includes('Files')) {
      isDragOver.value = true
    }
  }

  const handleDragLeave = () => {
    isDragOver.value = false
  }

  const handleDrop = async (e: DragEvent) => {
    isDragOver.value = false
    const files = e.dataTransfer?.files
    if (files && files.length > 0) {
      await uploadStore.addFiles(files)
      router.push('/upload')
    }
  }

  const visible = computed(() => {
    if (!authStore.isLoggedIn) {
      return false
    }

    const allowedRoutes = ['/dashboard', '/folders', '/resource']
    return allowedRoutes.some((allowedRoute) => route.path.startsWith(allowedRoute))
  })

  onMounted(() => {
    initializePosition()
    window.addEventListener('resize', initializePosition)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', initializePosition)
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
  })
</script>

<template>
  <Teleport to="body">
    <Transition name="cyber-float">
      <div
        v-if="visible"
        class="upload-float-button"
        :class="{ 'drag-over': isDragOver, 'has-files': uploadStore.statistics.total > 0, dragging: isDragging }"
        :style="floatStyle"
        :title="
          uploadStore.statistics.total > 0
            ? $t('upload.queueHasFiles', { count: uploadStore.statistics.total })
            : $t('upload.uploadFiles')
        "
        @mousedown="handleMouseDown"
        @click="handleClick"
        @dragover.prevent="handleDragOver"
        @dragleave.prevent="handleDragLeave"
        @drop.prevent="handleDrop"
      >
        <div class="button-content">
          <div class="upload-icon">
            <i class="fas fa-cloud-upload-alt" />
          </div>

          <div v-if="uploadStore.statistics.total > 0" class="queue-badge">
            {{ uploadStore.statistics.total }}
          </div>
        </div>

        <input ref="fileInput" type="file" multiple accept="image/*" style="display: none" @change="handleFileSelect" />
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .upload-float-button {
    width: 64px;
    height: 64px;
    background: linear-gradient(135deg, var(--color-brand-500) 0%, var(--color-brand-600) 50%, var(--color-brand-700) 100%);
    border: 2px solid var(--color-brand-400);
    border-radius: var(--radius-full);
    transition: all var(--transition-slow) var(--ease-in-out);
    backdrop-filter: blur(8px);
    box-shadow:
      0 0 20px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 40px rgba(var(--color-brand-500-rgb), 0.3),
      0 4px 12px rgba(0, 0, 0, 0.3);
    animation: float-breath 3s ease-in-out infinite;
  }

  .upload-float-button::before {
    content: '';
    position: absolute;
    inset: -2px;
    border-radius: var(--radius-full);
    padding: 2px;
    background: linear-gradient(135deg, var(--color-brand-300), transparent, var(--color-brand-500));
    -webkit-mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    opacity: 0.5;
    pointer-events: none;
    transition: opacity var(--transition-slow);
  }

  .upload-float-button:hover {
    transform: scale(1.08);
    border-color: var(--color-brand-300);
    box-shadow:
      0 0 30px rgba(var(--color-brand-400-rgb), 0.8),
      0 0 60px rgba(var(--color-brand-500-rgb), 0.5),
      0 8px 20px rgba(0, 0, 0, 0.4);
  }

  .upload-float-button:hover::before {
    opacity: 0.8;
  }

  .upload-float-button.dragging {
    transition: none !important;
    animation: none !important;
    transform: scale(0.98);
    box-shadow:
      0 0 15px rgba(var(--color-brand-500-rgb), 0.5),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.3),
      0 2px 8px rgba(0, 0, 0, 0.3);
  }

  .upload-float-button.dragging::before {
    transition: none !important;
  }

  .upload-float-button.dragging .upload-icon {
    transition: none !important;
  }

  .upload-float-button:active {
    transform: scale(1.04);
  }

  .upload-float-button.drag-over {
    background: linear-gradient(135deg, var(--color-success-500) 0%, var(--color-success-600) 50%, var(--color-brand-500) 100%);
    border-color: var(--color-success-400);
    transform: scale(1.12);
    box-shadow:
      0 0 40px rgba(var(--color-success-rgb), 0.7),
      0 0 80px rgba(var(--color-success-rgb), 0.4),
      0 8px 20px rgba(0, 0, 0, 0.4);
  }

  .upload-float-button.has-files {
    border-color: var(--color-brand-300);
    animation:
      float-breath 3s ease-in-out infinite,
      has-files-pulse 2s ease-in-out infinite;
  }

  .button-content {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .upload-icon {
    color: var(--color-text-on-brand);
    font-size: 26px;
    transition: all var(--transition-slow) var(--ease-in-out);
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
  }

  .upload-float-button:hover .upload-icon {
    transform: scale(1.1);
    filter: drop-shadow(0 0 8px rgba(255, 255, 255, 0.6));
  }

  .queue-badge {
    position: absolute;
    top: -4px;
    right: -4px;
    min-width: 22px;
    height: 22px;
    background: var(--color-error-500);
    color: var(--color-text-on-error);
    font-size: 11px;
    font-weight: 700;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--color-background-900);
    font-family: var(--font-mono);
    box-shadow:
      0 0 15px rgba(var(--color-error-rgb), 0.7),
      0 0 30px rgba(var(--color-error-rgb), 0.4);
    animation: badge-pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }

  @keyframes float-breath {
    0%,
    100% {
      box-shadow:
        0 0 20px rgba(var(--color-brand-500-rgb), 0.6),
        0 0 40px rgba(var(--color-brand-500-rgb), 0.3),
        0 4px 12px rgba(0, 0, 0, 0.3);
    }
    50% {
      box-shadow:
        0 0 30px rgba(var(--color-brand-500-rgb), 0.8),
        0 0 60px rgba(var(--color-brand-500-rgb), 0.5),
        0 6px 16px rgba(0, 0, 0, 0.3);
    }
  }

  @keyframes has-files-pulse {
    0%,
    100% {
      border-color: var(--color-brand-300);
    }
    50% {
      border-color: var(--color-brand-200);
    }
  }

  @keyframes badge-pulse {
    0%,
    100% {
      transform: scale(1);
      box-shadow:
        0 0 15px rgba(var(--color-error-rgb), 0.7),
        0 0 30px rgba(var(--color-error-rgb), 0.4);
    }
    50% {
      transform: scale(1.05);
      box-shadow:
        0 0 20px rgba(var(--color-error-rgb), 0.9),
        0 0 40px rgba(var(--color-error-rgb), 0.6);
    }
  }

  .cyber-float-enter-active,
  .cyber-float-leave-active {
    transition: all var(--transition-slow) var(--ease-in-out);
  }

  .cyber-float-enter-from {
    transform: scale(0.5);
    opacity: 0;
  }

  .cyber-float-leave-to {
    transform: scale(0.5);
    opacity: 0;
  }

  @media (max-width: 768px) {
    .upload-float-button {
      width: 56px;
      height: 56px;
    }

    .upload-icon {
      font-size: 22px;
    }

    .queue-badge {
      min-width: 20px;
      height: 20px;
      font-size: 10px;
    }
  }

  @media (max-width: 480px) {
    .upload-float-button {
      width: 52px;
      height: 52px;
    }

    .upload-icon {
      font-size: 20px;
    }

    .queue-badge {
      min-width: 18px;
      height: 18px;
      font-size: 9px;
    }
  }
</style>
