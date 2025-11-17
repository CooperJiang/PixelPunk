import type { FeatureItem } from './types'
import { useTexts } from '@/composables/useTexts'

/* ==================== 功能特性 ==================== */

/* API 测试功能特性 - 使用函数以支持i18n */
export const getApiTesterFeatures = (): FeatureItem[] => {
  const { $t } = useTexts()
  return [
    {
      id: 'file-upload',
      icon: 'fas fa-upload',
      title: $t('docs.apiTester.features.fileUpload.title'),
      description: $t('docs.apiTester.features.fileUpload.description'),
    },
    {
      id: 'params-config',
      icon: 'fas fa-cogs',
      title: $t('docs.apiTester.features.paramsConfig.title'),
      description: $t('docs.apiTester.features.paramsConfig.description'),
    },
    {
      id: 'real-time-response',
      icon: 'fas fa-eye',
      title: $t('docs.apiTester.features.realTimeResponse.title'),
      description: $t('docs.apiTester.features.realTimeResponse.description'),
    },
    {
      id: 'result-preview',
      icon: 'fas fa-image',
      title: $t('docs.apiTester.features.resultPreview.title'),
      description: $t('docs.apiTester.features.resultPreview.description'),
    },
  ]
}

/* 保持向后兼容 */
export const API_TESTER_FEATURES = getApiTesterFeatures()

/* ==================== 站点配置 ==================== */

export const SITE_DOMAIN = import.meta.env.VITE_SITE_DOMAIN || 'http://localhost:5173'

/* ==================== 滚动配置 ==================== */

/**
 * 返回顶部按钮显示阈值(px)
 */
export const BACK_TO_TOP_THRESHOLD = 300

/**
 * 滚动到区域时的偏移量(px)
 */
export const SCROLL_OFFSET = 100

/**
 * 章节激活判定偏移量(px)
 */
export const SECTION_ACTIVE_OFFSET = 150

/* ==================== UI交互配置 ==================== */

/**
 * 复制成功提示持续时间(ms)
 */
export const COPY_SUCCESS_DURATION = 2000

/**
 * DOM初始化延迟(ms)
 */
export const DOM_INIT_DELAY = 100

/**
 * 二次滚动初始化延迟(ms)
 */
export const SCROLL_REINIT_DELAY = 200
