<script setup lang="ts">
  import { computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useBreadcrumbStore } from '@/store/breadcrumb'
  import { useTexts } from '@/composables/useTexts'
  import type { FolderInfo } from '@/api/types/index'

  const { $t } = useTexts()

  interface BreadcrumbItem {
    label: string
    path?: string
  }

  interface Props {
    items?: BreadcrumbItem[] | FolderInfo[] // 支持两种格式
    autoGenerate?: boolean // 是否自动根据路由生成面包屑
    maxItems?: number // 最大显示项数
    startVisible?: number // 开始显示的项数
    endVisible?: number // 结束显示的项数
  }

  const props = withDefaults(defineProps<Props>(), {
    items: () => [],
    autoGenerate: true,
    maxItems: 6,
    startVisible: 1,
    endVisible: 2,
  })

  const route = useRoute()
  const router = useRouter()
  const breadcrumbStore = useBreadcrumbStore()

  /* 标准化面包屑数据格式 */
  const normalizeItems = (items: BreadcrumbItem[] | FolderInfo[]): BreadcrumbItem[] => {
    return items.map((item, _index) => {
      if ('name' in item && 'id' in item) {
        return {
          label: item.name,
          path: item.id ? `/folders/${item.id}` : '/folders',
        }
      }
      return item as BreadcrumbItem
    })
  }

  const routeNameMap = computed<Record<string, string>>(() => ({
    dashboard: $t('components.breadcrumb.routes.dashboard'),
    folders: $t('components.breadcrumb.routes.folders'),
    explore: $t('components.breadcrumb.routes.explore'),
    'category-manage': $t('components.breadcrumb.routes.categoryManage'),
    'tag-manage': $t('components.breadcrumb.routes.tagManage'),
    automation: $t('components.breadcrumb.routes.automation'),
    shares: $t('components.breadcrumb.routes.shares'),
    settings: $t('components.breadcrumb.routes.settings'),
  }))

  const breadcrumbItems = computed<BreadcrumbItem[]>(() => {
    if (props.items.length > 0) {
      const normalized = normalizeItems(props.items)
      return normalized.map((item, index) => ({
        ...item,
        path: index === normalized.length - 1 ? undefined : item.path,
      }))
    }

    if (!props.autoGenerate) {
      return []
    }

    const items: BreadcrumbItem[] = [{ label: $t('components.breadcrumb.home'), path: '/dashboard' }]

    if ((route.name === 'folders' || route.name === 'foldersWithPath') && breadcrumbStore.items.length > 0) {
      items.push({ label: $t('components.breadcrumb.routes.folders'), path: '/folders' })

      for (let i = 0; i < breadcrumbStore.items.length; i++) {
        const item = breadcrumbStore.items[i]

        if (i === 0) continue

        if (i === breadcrumbStore.items.length - 1) {
          items.push({ label: item.name })
        } else {
          items.push({ label: item.name, path: item.path })
        }
      }

      return items
    }

    const pathSegments = route.path.split('/').filter(Boolean)
    let currentPath = ''
    pathSegments.forEach((segment, index) => {
      currentPath += `/${segment}`
      const label = routeNameMap.value[segment] || segment

      if (index === pathSegments.length - 1) {
        items.push({ label })
      } else {
        items.push({ label, path: currentPath })
      }
    })

    return items
  })

  const handleClick = (item: BreadcrumbItem) => {
    if (item.path) {
      router.push(item.path)
    }
  }
</script>

<template>
  <nav class="breadcrumb-nav" :aria-label="$t('components.breadcrumb.ariaLabel')">
    <ol class="breadcrumb-list">
      <li v-for="(item, index) in breadcrumbItems" :key="index" class="breadcrumb-item">
        <button v-if="item.path" class="breadcrumb-link" @click="handleClick(item)">
          <span>{{ item.label }}</span>
        </button>

        <span v-else class="breadcrumb-current">
          <span>{{ item.label }}</span>
        </span>

        <span v-if="index < breadcrumbItems.length - 1" class="breadcrumb-separator"> / </span>
      </li>
    </ol>
  </nav>
</template>

<style scoped>
  .breadcrumb-nav {
    padding: 0;
  }

  .breadcrumb-list {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .breadcrumb-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .breadcrumb-link,
  .breadcrumb-current {
    display: inline-flex;
    align-items: center;
    font-size: 0.875rem;
    line-height: 1.25rem;
    transition: all 0.2s ease;
  }

  .breadcrumb-link {
    color: var(--color-content-muted);
    background: transparent;
    border: none;
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
  }

  .breadcrumb-link:hover {
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .breadcrumb-current {
    color: var(--color-content-heading);
    font-weight: 500;
    padding: 0.25rem 0.5rem;
  }

  .breadcrumb-separator {
    margin: 0 0.25rem;
    color: var(--color-content-disabled);
    user-select: none;
  }
</style>
