import { computed, ref, getCurrentInstance } from 'vue'
import type { FileInfo } from '../types'

export function useFileViewer(initialFile: FileInfo | null = null) {
  const instance = getCurrentInstance()
  const $t = instance?.appContext.config.globalProperties.$t

  const isLoading = ref(true)
  const hasError = ref(false)
  const isFirstView = ref(true)
  const isAllContentVisible = ref(false)
  const isImmersiveMode = ref(false)
  const isFullscreen = ref(false)
  const isFillMode = ref(false) // 默认适应模式，暂时不使用智能调整
  const isHorizontalFile = ref(true)
  const showModeIndicator = ref(false)
  const showZoomIndicator = ref(false)
  const isControlsHidden = ref(false)

  const loadingProgress = ref(0)
  const loadingText = ref($t?.('loading').value || 'Loading...')

  const controlsTimer = ref<number | null>(null)

  const shouldUseFillMode = computed(() => isFillMode.value)

  const toggleControls = () => {
    if (isFirstView.value) {
      isFirstView.value = false
      isAllContentVisible.value = true
    } else {
      isAllContentVisible.value = !isAllContentVisible.value
    }
  }

  const toggleFullscreen = async () => {
    if (!document.fullscreenElement) {
      await document.documentElement.requestFullscreen()
      isFullscreen.value = true
    } else {
      await document.exitFullscreen()
      isFullscreen.value = false
    }
  }

  const toggleFitMode = () => {
    isFillMode.value = !isFillMode.value
    showModeIndicatorTemporary()
  }

  const showModeIndicatorTemporary = () => {
    showModeIndicator.value = true
    setTimeout(() => {
      showModeIndicator.value = false
    }, 1500)
  }

  const showZoomIndicatorTemporary = () => {
    showZoomIndicator.value = true
    setTimeout(() => {
      showZoomIndicator.value = false
    }, 1000)
  }

  const resetControlsTimer = () => {
    if (controlsTimer.value) {
      clearTimeout(controlsTimer.value)
    }
    controlsTimer.value = setTimeout(() => {
      isControlsHidden.value = true
    }, 3000)
  }

  let abortController: AbortController | null = null
  let progressInterval: number | null = null

  const startRealImageLoading = async (imageUrl: string): Promise<Blob | null> => {
    loadingProgress.value = 0
    loadingText.value = $t?.('loading').value || 'Loading...'

    if (abortController) {
      abortController.abort()
    }

    abortController = new AbortController()

    try {
      const response = await fetch(imageUrl, {
        signal: abortController.signal,
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const contentLength = response.headers.get('content-length')
      if (!contentLength) {
        loadingText.value = $t?.('processing').value || 'Processing...'
        startLoadingProgressSimulation()
        return response.blob()
      }

      const total = parseInt(contentLength, 10)
      let loaded = 0

      const reader = response.body?.getReader()
      if (!reader) {
        startLoadingProgressSimulation()
        return response.blob()
      }

      loadingText.value = $t?.('processing').value || 'Processing...'
      const chunks: Uint8Array[] = []

      let done = false
      while (!done) {
        const result = await reader.read()
        done = result.done

        if (done) {
          loadingProgress.value = 100
          loadingText.value = $t?.('completed').value || 'Completed'
          break
        }

        const value = result.value
        chunks.push(value)
        loaded += value.length
        const progress = Math.round((loaded / total) * 100)
        loadingProgress.value = Math.min(progress, 99) // 保留1%给文件解析

        const loadedKB = Math.round(loaded / 1024)
        const totalKB = Math.round(total / 1024)
        loadingText.value = `${$t?.('processing').value || 'Processing'} ${loadedKB}KB / ${totalKB}KB`
      }

      const blob = new Blob(chunks)
      return blob
    } catch (error) {
      if (error.name === 'AbortError') {
        return null // 请求被取消
      }
      console.warn('Failed to get real progress, falling back to simulation:', error)
      startLoadingProgressSimulation()
      return null
    }
  }

  const startLoadingProgressSimulation = () => {
    loadingProgress.value = 0
    loadingText.value = $t?.('loading').value || 'Loading...'

    progressInterval = setInterval(() => {
      if (loadingProgress.value < 90) {
        loadingProgress.value += Math.random() * 20
        if (loadingProgress.value > 90) {
          loadingProgress.value = 90
        }
      }
    }, 200)
  }

  const stopLoadingProgressSimulation = () => {
    if (progressInterval) {
      clearInterval(progressInterval)
      progressInterval = null
    }
  }

  const handleFileLoad = (fileRef: HTMLImageElement) => {
    stopLoadingProgressSimulation()

    loadingProgress.value = 100
    loadingText.value = $t?.('completed').value || 'Completed'

    setTimeout(() => {
      isLoading.value = false
      hasError.value = false

      if (fileRef) {
        const { naturalWidth } = fileRef
        const { naturalHeight } = fileRef
        const aspectRatio = naturalWidth / naturalHeight

        isHorizontalFile.value = aspectRatio > 1.2

        analyzeFileBrightness(fileRef)
      }
    }, 300)
  }

  const analyzeFileBrightness = (fileRef: HTMLImageElement) => {
    try {
      const fileUrl = fileRef.src
      const currentOrigin = window.location.origin
      const fileOrigin = new URL(fileUrl).origin

      if (fileOrigin !== currentOrigin) {
        document.documentElement.style.setProperty('--viewer-is-light', '0')
        document.documentElement.style.setProperty('--text-shadow-intensity', '1')
        document.documentElement.style.setProperty('--border-opacity', '0.3')
        return
      }

      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')
      if (!ctx) {
        return
      }

      const sampleSize = 50
      canvas.width = sampleSize
      canvas.height = sampleSize

      ctx.drawImage(fileRef, 0, 0, sampleSize, sampleSize)
      const imageData = ctx.getImageData(0, 0, sampleSize, sampleSize)
      const { data } = imageData

      let totalBrightness = 0
      const pixelCount = data.length / 4

      for (let i = 0; i < data.length; i += 4) {
        const brightness = data[i] * 0.299 + data[i + 1] * 0.587 + data[i + 2] * 0.114
        totalBrightness += brightness
      }

      const averageBrightness = totalBrightness / pixelCount
      const isLightFile = averageBrightness > 128

      document.documentElement.style.setProperty('--viewer-is-light', isLightFile ? '1' : '0')
      document.documentElement.style.setProperty('--text-shadow-intensity', isLightFile ? '3' : '1')
      document.documentElement.style.setProperty('--border-opacity', isLightFile ? '0.7' : '0.3')
    } catch {
      document.documentElement.style.setProperty('--viewer-is-light', '0')
      document.documentElement.style.setProperty('--text-shadow-intensity', '1')
      document.documentElement.style.setProperty('--border-opacity', '0.3')
    }
  }

  const handleFileError = () => {
    loadingProgress.value = 0
    loadingText.value = $t?.('error').value || 'Error'
    isLoading.value = false
    hasError.value = true
  }

  const retryLoading = () => {
    isLoading.value = true
    hasError.value = false
    startLoadingProgressSimulation()
  }

  const handleFullscreenChange = () => {
    isFullscreen.value = Boolean(document.fullscreenElement)
  }

  if (initialFile) {
    startLoadingProgressSimulation()
  }

  return {
    isLoading,
    hasError,
    isFirstView,
    isAllContentVisible,
    isImmersiveMode,
    isFullscreen,
    isFillMode,
    isHorizontalFile,
    showModeIndicator,
    showZoomIndicator,
    isControlsHidden,
    loadingProgress,
    loadingText,

    shouldUseFillMode,

    toggleControls,
    toggleFullscreen,
    toggleFitMode,
    showModeIndicatorTemporary,
    showZoomIndicatorTemporary,
    resetControlsTimer,
    handleFileLoad,
    handleFileError,
    retryLoading,
    startRealImageLoading,
    startLoadingProgressSimulation,
    stopLoadingProgressSimulation,
    handleFullscreenChange,
  }
}
