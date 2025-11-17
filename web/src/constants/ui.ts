/**
 * UI相关常量定义

/* ===== 页面布局常量 ===== */
export const LAYOUT = {
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
  DEFAULT_PAGE_SIZE: 20,

  GRID: {
    MIN_SIZE: 120,
    MAX_SIZE: 300,
    DEFAULT_SIZE: 150,
    STEP: 5,
  },

  WATERFALL: {
    DEFAULT_COLUMNS: 4,
    DEFAULT_COLUMN_WIDTH: 250,
    DEFAULT_GAP: 15,
    INITIAL_LOAD_COUNT: 20,
    MIN_LOAD_COUNT: 12,
    MAX_LOAD_COUNT: 24,
    MAX_TAGS_TO_SHOW: 2,
    DEFAULT_IMAGE_HEIGHT: 300,
    INFO_AREA_HEIGHT: 140,
    MIN_COLUMN_WIDTH_RATIO: 1.2,
    MIN_COLUMN_WIDTH: 180,
  },

  PAGINATION: {
    MAX_VISIBLE_PAGES: 7,
    FRONT_PAGES_THRESHOLD: 4,
    BACK_PAGES_THRESHOLD: 3,
    DROPDOWN_ITEM_HEIGHT: 40,
  },
} as const

/* ===== 时间相关常量 ===== */
export const TIMING = {
  TOAST: {
    DEFAULT_DURATION: 3000, // 3秒
    ANIMATION_DELAY: 300, // 300ms
    SHOW_DELAY: 50, // 50ms
    MARGIN: 16, // 16px
    INITIAL_TOP: 20, // 20px
  },

  SCROLL: {
    DEBOUNCE_TIME: 50, // 50ms
    UPLOAD_DEBOUNCE_TIME: 500, // 500ms
  },

  ANIMATION: {
    HIGHLIGHT_DURATION: 1500, // 1.5s
    PRELOAD_DELAY: 50, // 50ms
    LOADING_DELAY: 300, // 300ms
  },

  REQUEST: {
    SHORT_TIMEOUT: 5000, // 5秒
    MEDIUM_TIMEOUT: 8000, // 8秒
    LONG_TIMEOUT: 10000, // 10秒
  },
} as const

/* ===== Z-Index层级常量 ===== */
export const Z_INDEX = {
  BASE: 1,
  CONTENT: 10,
  DROPDOWN: 1000,
  MODAL: 2000,
  TOAST: 9999,
  OVERLAY: 5000,
  LOADING: 100,
  IMAGE_INFO: 5,
  INDICATOR: 2,
} as const

/* ===== 尺寸限制常量 ===== */
export const SIZE_LIMITS = {
  MAX_FILE_SIZE: 20 * 1024 * 1024, // 20MB
  MAX_IMAGE_COUNT: 500, // 最大文件数量

  IMAGE: {
    MIN_WIDTH: 100,
    MAX_WIDTH: 8192,
    MIN_HEIGHT: 100,
    MAX_HEIGHT: 8192,
    THUMBNAIL_SIZE: 200,
  },

  COMPONENT: {
    UPLOAD_SCROLL_HEIGHT: 500, // 上传滚动器高度
    GRID_SELECTOR_WIDTH: 280, // 网格选择器宽度
    PRELOAD_BATCH_SIZE: 3, // 预加载批次大小
  },
} as const

/* ===== 响应式断点常量 ===== */
export const BREAKPOINTS = {
  SM: 640,
  MD: 768,
  LG: 1024,
  XL: 1280,
  '2XL': 1536,
} as const

/* ===== 常用样式类名 ===== */
export const COMMON_CLASSES = {
  LAYOUT: {
    FLEX: 'flex',
    GRID: 'grid',
    CONTAINER: 'container',
    CENTER: 'mx-auto',
  },

  RESPONSIVE: {
    MD_FLEX_ROW: 'md:flex-row',
    SM_GRID_COLS_2: 'sm:grid-cols-2',
    LG_MAX_W_4XL: 'lg:max-w-4xl',
  },

  SPACING: {
    P_4: 'p-4',
    M_4: 'm-4',
    GAP_4: 'gap-4',
  },
} as const
