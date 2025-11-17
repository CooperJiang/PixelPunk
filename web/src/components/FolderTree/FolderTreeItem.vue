<script lang="ts">
  import { computed, defineComponent, ref, type PropType } from 'vue'
  import { getFolderList } from '@/api/folder'
  import type { FolderInfo } from '@/api/types'

  export default defineComponent({
    name: 'CyberFolderTreeItem',
    props: {
      folder: {
        type: Object as PropType<FolderInfo>,
        required: true,
      },
      level: {
        type: Number,
        default: 0,
      },
      selectedId: {
        type: String,
        default: '',
      },
      borderColor: {
        type: String,
        default: 'var(--color-brand-500)',
      },
      isExpanded: {
        type: Boolean,
        default: false,
      },
      expandedFolders: {
        type: Set as PropType<Set<string>>,
        default: () => new Set(),
      },
      cachedChildren: {
        type: Array as PropType<FolderInfo[]>,
        default: () => [],
      },
      getFolderChildren: {
        type: Function as PropType<(folderId: string) => FolderInfo[]>,
        required: true,
      },
    },
    emits: ['select', 'toggle', 'children-loaded', 'subfolder-request-start', 'subfolder-request-end'],
    setup(props, { emit }) {
      const isLoading = ref(false)
      const isExpanded = computed(() => props.isExpanded)
      const children = computed(() => props.cachedChildren)
      const hasChildren = computed(() => {
        if (props.folder.has_children === false) {
          return false
        }
        return true
      })

      const toggleExpand = async (e: Event) => {
        e.stopPropagation()

        const newExpandedState = !isExpanded.value
        if (newExpandedState && children.value.length === 0 && !isLoading.value) {
          isLoading.value = true
          emit('subfolder-request-start')

          try {
            const result = await getFolderList(props.folder.id)
            if (result.success) {
              const folders = result.data || []
              const processedFolders = folders.map((folder) => ({
                ...folder,
                fullPath: `${props.folder.fullPath} / ${folder.name}`,
              }))
              emit('children-loaded', props.folder.id, processedFolders)
            } else {
              console.error('Failed to fetch subfolders:', result.message)
            }
          } catch (error) {
            console.error('Failed to fetch subfolders:', error)
          } finally {
            isLoading.value = false
            emit('subfolder-request-end')
          }
        }
        emit('toggle', props.folder.id, newExpandedState)
      }

      const handleToggleClick = (e: Event) => {
        if (hasChildren.value) {
          toggleExpand(e)
        }
      }

      const selectFolderItem = () => {
        emit('select', props.folder)
      }

      return {
        isLoading,
        children,
        hasChildren,
        toggleExpand,
        handleToggleClick,
        selectFolderItem,
      }
    },
  })
</script>

<template>
  <div
    class="cyber-folder-tree-item"
    :class="{ 'cyber-folder-tree-item--selected': selectedId === folder.id }"
    :style="{ paddingLeft: `${level * 16 + 8}px` }"
    @click.stop="selectFolderItem"
  >
    <div class="cyber-folder-tree-item__content">
      <div
        class="cyber-folder-tree-item__toggle"
        :class="{ 'cyber-folder-tree-item__toggle--has-children': hasChildren }"
        @click.stop="handleToggleClick"
      >
        <i v-if="hasChildren && isLoading" class="fas fa-spinner fa-spin cyber-folder-tree-item__toggle-icon" />
        <i
          v-else-if="hasChildren && isExpanded"
          class="fas fa-chevron-down cyber-folder-tree-item__toggle-icon cyber-folder-tree-item__toggle-icon--expanded"
        />
        <i v-else-if="hasChildren" class="fas fa-chevron-right cyber-folder-tree-item__toggle-icon" />
      </div>

      <div class="cyber-folder-tree-item__icon">
        <i class="fas fa-folder" :style="{ color: borderColor }" />
      </div>

      <span class="cyber-folder-tree-item__name">{{ folder.name }}</span>
    </div>
  </div>

  <transition name="cyber-folder-tree-item-expand">
    <div v-if="isExpanded" class="cyber-folder-tree-item__children" :style="{ '--guide-color': `${borderColor}40` }">
      <cyber-folder-tree-item
        v-for="child in children"
        :key="child.id"
        :folder="child"
        :level="level + 1"
        :selected-id="selectedId"
        :border-color="borderColor"
        :is-expanded="expandedFolders.has(child.id)"
        :expanded-folders="expandedFolders"
        :cached-children="getFolderChildren(child.id)"
        :get-folder-children="getFolderChildren"
        @select="$emit('select', $event)"
        @toggle="(folderId, isExpanded) => $emit('toggle', folderId, isExpanded)"
        @children-loaded="(folderId, children) => $emit('children-loaded', folderId, children)"
        @subfolder-request-start="$emit('subfolder-request-start')"
        @subfolder-request-end="$emit('subfolder-request-end')"
      />
    </div>
  </transition>
</template>

<style scoped>
  .cyber-folder-tree-item {
    position: relative;
    cursor: pointer;
    transition: all 0.2s ease;
    padding: 2px 0;
  }

  .cyber-folder-tree-item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-folder-tree-item--selected {
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-folder-tree-item__content {
    display: flex;
    align-items: center;
    height: 28px;
    padding-right: 8px;
  }

  .cyber-folder-tree-item__toggle {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    margin-right: 4px;
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
    color: rgba(var(--color-cyber-light-rgb), 0.7);
    cursor: default;
    opacity: 0.3;
  }

  .cyber-folder-tree-item__toggle--has-children {
    cursor: pointer;
    opacity: 1;
  }

  .cyber-folder-tree-item__toggle--has-children:hover {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: rgba(var(--color-cyber-light-rgb), 0.9);
  }

  .cyber-folder-tree-item__toggle-icon {
    font-size: 10px;
    transition: transform 0.2s ease;
  }

  .cyber-folder-tree-item__toggle-icon--expanded {
    transform: rotate(0);
  }

  .cyber-folder-tree-item__icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    margin-right: 6px;
    font-size: 12px;
  }

  .cyber-folder-tree-item__name {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 13px;
    color: var(--color-cyber-light);
  }

  .cyber-folder-tree-item__children {
    position: relative;
    --guide-color: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .cyber-folder-tree-item-expand-enter-active,
  .cyber-folder-tree-item-expand-leave-active {
    transition: all 0.3s ease;
    max-height: 1000px;
    overflow: hidden;
    opacity: 1;
  }

  .cyber-folder-tree-item-expand-enter-from,
  .cyber-folder-tree-item-expand-leave-to {
    max-height: 0;
    opacity: 0;
  }
</style>
