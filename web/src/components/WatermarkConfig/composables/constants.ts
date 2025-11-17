export const DEFAULT_PREVIEW_SIZE = {
  width: 2000,
  height: 1200,
}

export const CANVAS_CONFIG = {
  minPreviewSize: 500,
  maxPreviewSize: 10000,
  minZoom: 0.1,
  maxZoom: 5,
  zoomSpeed: 0.9,
  inverseZoomSpeed: 1.1,
  fitMargin: 80,
  baseGridSize: 50,
  handleSize: 12,
  handleHitThreshold: 4,
} as const

export const RULER_CONFIG = {
  intervals: {
    zoom2x: 50,
    zoom1x: 100,
    zoom05x: 200,
    default: 500,
  },
  zoomThresholds: {
    high: 2,
    medium: 1,
    low: 0.5,
  },
} as const

export const PRESET_SIZES = [
  { name: '1080p', width: 1920, height: 1080 },
  { name: '2K', width: 2560, height: 1440 },
  { name: '4K', width: 3840, height: 2160 },
  { name: '8K', width: 7680, height: 4320 },
] as const
