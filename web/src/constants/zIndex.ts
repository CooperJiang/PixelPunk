/**
 * Z-Index 层级管理
 */

export const Z_INDEX = {
  BASE: 1,
  NORMAL: 1,
  TABLE_HEADER: 100,
  TOOLTIP: 500,
  DROPDOWN: 1000,
  POPOVER: 1500,
  MODAL_BACKDROP: 2000,
  MODAL: 2100,
  MESSAGE: 3000,
  DRAWER: 5000,
  DRAWER_MODAL: 5100,
  FULLSCREEN_MODAL: 9000,
  TOAST: 10000,
  LOADING: 10100,
} as const

export const getZIndex = (base: keyof typeof Z_INDEX, offset = 0) => Z_INDEX[base] + offset

export const getTopZIndex = () => 99999
