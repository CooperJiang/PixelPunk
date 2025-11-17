import { ref } from 'vue'
import type { FileInfo } from '@/api/types'
import { setCssVariable } from '@/utils/ui/domUtils'
import type { ViewportDimensions } from '../types'
import { logger } from '@/utils/system/logger'

export function useHiveUI() {
  const isFullscreen = ref(false)
  const tipsHidden = ref(false)
  const viewportWidth = ref(0)
  const viewportHeight = ref(0)
  const showPreview = ref(false)
  const previewFile = ref<FileInfo | null>(null)

  const getViewportDimensions = (): ViewportDimensions => ({
    width: window.innerWidth,
    height: window.innerHeight,
  })

  const updateViewportDimensions = () => {
    const { width, height } = getViewportDimensions()
    viewportWidth.value = width
    viewportHeight.value = height
    setCssVariable('viewport-width', `${width}px`)
    setCssVariable('viewport-height', `${height}px`)
  }

  const toggleTips = () => {
    tipsHidden.value = !tipsHidden.value
  }

  const toggleFullscreen = async () => {
    try {
      if (!document.fullscreenElement) {
        await document.documentElement.requestFullscreen()
      } else {
        await document.exitFullscreen()
      }
    } catch (error) {
      logger.error('Fullscreen toggle failed:', error)
    }
  }

  const handleFullscreenChange = () => {
    isFullscreen.value = Boolean(document.fullscreenElement)
  }

  const handleFileClick = (file: FileInfo) => {
    previewFile.value = file
    showPreview.value = true
  }

  const closePreview = () => {
    showPreview.value = false
    previewFile.value = null
  }

  const handleKeydown = (e: KeyboardEvent) => {
    switch (e.key) {
      case 'f':
      case 'F':
        e.preventDefault()
        toggleFullscreen()
        break
      case 'Escape':
        if (document.fullscreenElement) {
          document.exitFullscreen()
        }
        break
    }
  }

  return {
    isFullscreen,
    tipsHidden,
    viewportWidth,
    viewportHeight,
    showPreview,
    previewFile,
    updateViewportDimensions,
    toggleTips,
    toggleFullscreen,
    handleFullscreenChange,
    handleFileClick,
    closePreview,
    handleKeydown,
  }
}
