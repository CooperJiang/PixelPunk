export interface HexagonStyle {
  width: string
  height: string
  left: string
  top: string
  animationDelay: string
}

export interface LoaderHexStyle {
  transform: string
  animationDelay: string
}

export interface ViewportDimensions {
  width: number
  height: number
}

export interface HiveDataState<T = unknown> {
  images: T[]
  isLoading: boolean
  isInitialLoading: boolean
  isLoadingFadingOut: boolean
  error: string
  currentPage: number
  hasMorePages: boolean
  loadingProgress: number
  loadingStartTime: number
}

export interface HiveUIState<T = unknown> {
  isFullscreen: boolean
  tipsHidden: boolean
  viewportWidth: number
  viewportHeight: number
  showPreview: boolean
  previewFile: T | null
}

export interface LayoutParams {
  horizontalGap: number
  verticalGap: number
  strength3D: number
  clipGap: number
}
