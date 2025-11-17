<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import type { FolderInfo } from '@/api/types/index'

  const { $t } = useTexts()

  interface Props {
    items: FolderInfo[]
    maxItems?: number
    startVisible?: number
    endVisible?: number
    showIcon?: boolean // 是否显示根目录图标
    autoAddRoot?: boolean // 是否自动添加根目录
  }

  const props = withDefaults(defineProps<Props>(), {
    maxItems: 6,
    startVisible: 1,
    endVisible: 2,
    showIcon: true, // 默认显示图标
    autoAddRoot: true, // 默认自动添加根目录
  })

  const emit = defineEmits<{
    (e: 'click', item: FolderInfo): void
  }>()

  /* 根据 autoAddRoot 参数决定是否自动添加根目录 */
  const displayItems = computed<FolderInfo[]>(() => {
    if (!props.autoAddRoot) {
      return props.items
    }

    if (props.items.length === 0 || props.items[0].id !== null) {
      return [{ id: null, name: $t('folders.rootFolder') } as FolderInfo, ...props.items]
    }
    return props.items
  })

  const handleClick = (item: FolderInfo) => {
    emit('click', item)
  }
</script>

<template>
  <nav class="folder-breadcrumb" :aria-label="$t('folders.breadcrumb.label')">
    <ol class="breadcrumb-list">
      <li v-for="(item, index) in displayItems" :key="item.id || 'root'" class="breadcrumb-item">
        <button v-if="index < displayItems.length - 1" class="breadcrumb-link" @click="handleClick(item)">
          <i v-if="showIcon && item.id === null" class="fas fa-home root-icon"></i>
          <span>{{ item.name }}</span>
        </button>

        <button v-else-if="item.id === null" class="breadcrumb-link breadcrumb-current" @click="handleClick(item)">
          <i v-if="showIcon" class="fas fa-home root-icon"></i>
          <span>{{ item.name }}</span>
        </button>

        <span v-else class="breadcrumb-current">
          <span>{{ item.name }}</span>
        </span>

        <span v-if="index < displayItems.length - 1" class="breadcrumb-separator"> / </span>
      </li>
    </ol>
  </nav>
</template>

<style scoped>
  .folder-breadcrumb {
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

  .root-icon {
    margin-right: 0.375rem;
    font-size: 1rem;
    color: var(--color-brand-500);
    transition: all 0.2s ease;
  }

  .breadcrumb-link:hover .root-icon {
    color: var(--color-brand-500);
    filter: brightness(1.2);
  }

  .breadcrumb-current .root-icon {
    color: var(--color-content-heading);
  }
</style>
