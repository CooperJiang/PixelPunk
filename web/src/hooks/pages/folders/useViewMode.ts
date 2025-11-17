import { ref } from 'vue'
import { StorageUtil } from '@/utils/storage/storage'

/* 常量 */
const VIEW_MODE_KEY = 'folder_view_mode'

/**
 * 视图模式管理 Hook
 * 负责管理网格/列表视图切换及持久化存储
 */
export function useViewMode() {
  const viewMode = ref<'grid' | 'list'>(StorageUtil.get<'grid' | 'list'>(VIEW_MODE_KEY) || 'grid')

  const toggleViewMode = (mode: 'grid' | 'list') => {
    viewMode.value = mode
    StorageUtil.set<'grid' | 'list'>(VIEW_MODE_KEY, mode, 24 * 7)
  }

  return {
    viewMode,
    toggleViewMode,
  }
}
