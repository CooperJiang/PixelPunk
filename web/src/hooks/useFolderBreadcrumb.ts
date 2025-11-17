import { useRouter } from 'vue-router'

/**
 * 文件夹面包屑导航 Hook
 * 统一管理面包屑点击逻辑，避免页面刷新
 * 使用 router.replace 而不是 push，避免历史记录堆积
 */
export function useFolderBreadcrumb() {
  const router = useRouter()

  const handleBreadcrumbClick = (item: any) => {
    if (item.isSpecial && item.id === 'home') {
      router.replace('/dashboard')
      return
    }

    if (item.id === null) {
      router.replace('/folders')
      return
    }

    if (item.id) {
      router.replace(`/folders/${item.id}`)
      return
    }

    if (item.path) {
      router.replace(item.path)
    }
  }

  return {
    handleBreadcrumbClick,
  }
}
