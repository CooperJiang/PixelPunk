<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { storeToRefs } from 'pinia'
  import { useUploadStore } from '@/store/upload'
  import { useRoute } from 'vue-router'
  import { useDrag } from './composables/useDrag'
  import ProgressRing from './components/ProgressRing.vue'
  import CenterContent from './components/CenterContent.vue'
  import FileTooltip from './components/FileTooltip.vue'
  import type { GlobalUploadFloatEmits, GlobalUploadFloatExposed, GlobalUploadFloatProps } from './types'

  defineOptions({
    name: 'GlobalUploadFloat',
  })

  const props = withDefaults(defineProps<GlobalUploadFloatProps>(), {
    hideOnRoutes: () => ['/upload', '/dashboard'],
    draggable: true,
    showFileInfo: true,
  })

  const emit = defineEmits<GlobalUploadFloatEmits>()

  /* ✨ 优化：useUploadStore 会自动触发 initialize()
   * initialize() 中已经包含了懒加载配置的逻辑
   * 只有在非排除页面才会渲染此组件（见 App.vue） */
  const uploadStore = useUploadStore()
  const route = useRoute()
  const { statistics, globalProgress, allUploads } = storeToRefs(uploadStore)

  const isHidden = ref(false)
  const showFileInfo = ref(false)

  const { floatStyle, handleMouseDown, handleMouseUp } = useDrag(props.initialPosition)

  const visible = computed(() => {
    const hasUnfinishedTasks = statistics.value.pending > 0 || statistics.value.uploading > 0 || statistics.value.failed > 0
    const shouldHideOnCurrentRoute = props.hideOnRoutes?.some((hiddenRoute) => route.path.startsWith(hiddenRoute))
    return hasUnfinishedTasks && !isHidden.value && !shouldHideOnCurrentRoute
  })

  const status = computed(() => {
    if (statistics.value.failed > 0 && statistics.value.uploading === 0) {
      return 'error'
    }
    if (statistics.value.uploading > 0) {
      return 'uploading'
    }
    if (statistics.value.completed === statistics.value.total && statistics.value.total > 0) {
      return 'success'
    }
    if (statistics.value.pending > 0) {
      return 'preparing'
    }
    return 'preparing'
  })

  const queueSize = computed(() => statistics.value.pending + statistics.value.uploading + statistics.value.failed)

  const currentFileName = computed(() => {
    const uploadingFiles = allUploads.value.filter((f) => f.status === 'uploading')
    if (uploadingFiles.length > 0) {
      return uploadingFiles[0].file.name
    }

    const preparingFiles = allUploads.value.filter((f) => f.status === 'preparing')
    if (preparingFiles.length > 0) {
      return preparingFiles[0].file.name
    }

    return null
  })

  const handleMouseDownWrapper = (e: MouseEvent) => {
    handleMouseDown(e)
  }

  const handleMouseUpWrapper = (e: MouseEvent) => {
    const wasClick = handleMouseUp(e)
    if (wasClick) {
      emit('open-drawer')
    }
  }

  const handleMouseEnter = () => {
    showFileInfo.value = true
  }

  const handleMouseLeave = () => {
    showFileInfo.value = false
  }

  const hide = () => {
    isHidden.value = true
  }

  const show = () => {
    isHidden.value = false
  }

  defineExpose<GlobalUploadFloatExposed>({
    hide,
    show,
  })
</script>

<template>
  <Teleport to="body">
    <Transition name="cyber-float">
      <div
        v-if="visible"
        class="cyber-uploader-float"
        :style="floatStyle"
        @mousedown="handleMouseDownWrapper"
        @mouseup="handleMouseUpWrapper"
        @mouseenter="handleMouseEnter"
        @mouseleave="handleMouseLeave"
      >
        <div class="cyber-ring-outer">
          <div class="cyber-ring-main">
            <ProgressRing :progress="globalProgress" :active="status === 'uploading'" />
            <CenterContent :status="status" :progress="globalProgress" :queue-size="queueSize" />
          </div>
        </div>

        <FileTooltip :file-name="currentFileName" :statistics="statistics" :visible="showFileInfo" />
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
  .cyber-uploader-float {
    @apply pointer-events-auto fixed z-[10000] select-none transition-transform duration-200 ease-in-out;
  }

  .cyber-uploader-float:active {
    @apply scale-[0.98];
  }

  .cyber-ring-outer {
    @apply relative h-[90px] w-[90px];
  }

  .cyber-ring-main {
    @apply relative flex h-[90px] w-[90px] items-center justify-center rounded-full border-2 transition-all duration-300 ease-in-out;
    background: linear-gradient(135deg, rgba(15, 23, 42, 1) 0%, rgba(30, 41, 59, 0.98) 50%, rgba(15, 23, 42, 1) 100%);
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    backdrop-filter: blur(16px) saturate(180%);
    box-shadow:
      0 0 30px rgba(var(--color-brand-500-rgb), 0.4),
      0 0 60px rgba(var(--color-brand-500-rgb), 0.15),
      inset 0 0 30px rgba(0, 0, 0, 0.5),
      inset 0 2px 4px rgba(var(--color-brand-500-rgb), 0.1),
      0 8px 16px rgba(0, 0, 0, 0.6);
  }

  .cyber-ring-main::before {
    content: '';
    position: absolute;
    inset: -2px;
    border-radius: var(--radius-full);
    padding: 2px;
    background: linear-gradient(135deg, var(--color-brand-500), transparent, var(--color-brand-400));
    -webkit-mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    opacity: 0.3;
    pointer-events: none;
  }

  .cyber-ring-main:hover {
    @apply scale-105;
    border-color: rgba(var(--color-brand-400-rgb), 0.8);
    box-shadow:
      0 0 40px rgba(var(--color-brand-500-rgb), 0.6),
      0 0 80px rgba(var(--color-brand-500-rgb), 0.25),
      inset 0 0 30px rgba(0, 0, 0, 0.5),
      inset 0 2px 4px rgba(var(--color-brand-500-rgb), 0.2),
      0 12px 24px rgba(0, 0, 0, 0.7);
  }

  .cyber-ring-main:hover::before {
    opacity: 0.6;
  }

  .cyber-float-enter-active,
  .cyber-float-leave-active {
    @apply transition-all duration-300 ease-in-out;
  }

  .cyber-float-enter-from {
    @apply scale-50 opacity-0;
  }

  .cyber-float-leave-to {
    @apply scale-50 opacity-0;
  }
</style>
