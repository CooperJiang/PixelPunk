/**
 * 全局配置 Composable
 * @description 提供获取应用信息和品牌信息的便捷方法
 */
import { computed } from 'vue'
import { APP_CONFIG } from '@/constants/config'

/**
 * 全局配置 Hook
 * @returns 配置相关的方法和属性
 */
export function useConfig() {
  const appInfo = computed(() => ({
    name: APP_CONFIG.name,
    description: APP_CONFIG.description,
  }))

  const brand = computed(() => APP_CONFIG.brand)

  return {
    appInfo,
    brand,
  }
}
