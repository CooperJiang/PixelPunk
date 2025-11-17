<script setup lang="ts">
  import { defineEmits, defineProps, onBeforeUnmount, onMounted, ref, watch } from 'vue'

  interface DocumentWithFullscreen extends Document {
    mozCancelFullScreen?: () => Promise<void>
    webkitExitFullscreen?: () => Promise<void>
    msExitFullscreen?: () => Promise<void>
  }

  interface ElementWithFullscreen extends HTMLElement {
    mozRequestFullScreen?: () => Promise<void>
    webkitRequestFullscreen?: () => Promise<void>
    msRequestFullscreen?: () => Promise<void>
  }

  const props = defineProps({
    enabled: {
      type: Boolean,
      default: false,
    },
  })

  const emit = defineEmits(['update:enabled'])

  const isFullScreen = ref(false)
  const isImmersiveMode = ref(false)

  const toggleFullScreen = async () => {
    try {
      if (!isFullScreen.value) {
        const elem = document.documentElement as ElementWithFullscreen
        if (elem.requestFullscreen) {
          await elem.requestFullscreen()
        } else if (elem.mozRequestFullScreen) {
          await elem.mozRequestFullScreen()
        } else if (elem.webkitRequestFullscreen) {
          await elem.webkitRequestFullscreen()
        } else if (elem.msRequestFullscreen) {
          await elem.msRequestFullscreen()
        }

        isImmersiveMode.value = true
        document.body.classList.add('immersive-mode')
      } else {
        const doc = document as DocumentWithFullscreen
        if (doc.exitFullscreen) {
          await doc.exitFullscreen()
        } else if (doc.mozCancelFullScreen) {
          await doc.mozCancelFullScreen()
        } else if (doc.webkitExitFullscreen) {
          await doc.webkitExitFullscreen()
        } else if (doc.msExitFullscreen) {
          await doc.msExitFullscreen()
        }

        isImmersiveMode.value = false
        document.body.classList.remove('immersive-mode')
      }

      emit('update:enabled', isFullScreen.value)
    } catch {
      isImmersiveMode.value = !isImmersiveMode.value
      if (isImmersiveMode.value) {
        document.body.classList.add('immersive-mode')
      } else {
        document.body.classList.remove('immersive-mode')
      }

      emit('update:enabled', isImmersiveMode.value)
    }
  }

  const handleFullscreenChange = () => {
    isFullScreen.value = document.fullscreenElement !== null
    emit('update:enabled', isFullScreen.value)
  }

  watch(
    () => props.enabled,
    (enabled) => {
      if (enabled !== isFullScreen.value) {
        toggleFullScreen()
      }
    }
  )

  onMounted(() => {
    document.addEventListener('fullscreenchange', handleFullscreenChange)
    document.addEventListener('webkitfullscreenchange', handleFullscreenChange)
    document.addEventListener('mozfullscreenchange', handleFullscreenChange)
    document.addEventListener('MSFullscreenChange', handleFullscreenChange)

    isFullScreen.value = document.fullscreenElement !== null
  })

  onBeforeUnmount(() => {
    document.removeEventListener('fullscreenchange', handleFullscreenChange)
    document.removeEventListener('webkitfullscreenchange', handleFullscreenChange)
    document.removeEventListener('mozfullscreenchange', handleFullscreenChange)
    document.removeEventListener('MSFullscreenChange', handleFullscreenChange)

    if (isFullScreen.value || isImmersiveMode.value) {
      try {
        const doc = document as DocumentWithFullscreen
        if (doc.exitFullscreen) {
          doc.exitFullscreen()
        } else if (doc.mozCancelFullScreen) {
          doc.mozCancelFullScreen()
        } else if (doc.webkitExitFullscreen) {
          doc.webkitExitFullscreen()
        } else if (doc.msExitFullscreen) {
          doc.msExitFullscreen()
        }

        document.body.classList.remove('immersive-mode')
        isImmersiveMode.value = false
      } catch {}
    }
  })

  defineExpose({
    toggleFullScreen,
  })
</script>

<template>
  <div style="display: none"></div>
</template>

<style scoped></style>
