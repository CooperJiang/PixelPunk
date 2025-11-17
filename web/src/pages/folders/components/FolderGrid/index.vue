<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { ref, watch, reactive } from 'vue'
  import { VueDraggable } from 'vue-draggable-plus'
  import type { FolderGridEmits, FolderGridProps } from './types'
  import { useToast } from '@/components'
  import { getFolderList, moveFolder } from '@/api/folder'
  import type { ContextMenuItem } from '@/components/CyberContextMenu/types'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'FolderGrid',
  })

  const props = defineProps<FolderGridProps>()
  const emit = defineEmits<FolderGridEmits>()
  const toast = useToast()
  const { $t } = useTexts()

  const mutableFolders = ref([...props.folders])

  watch(
    () => props.folders,
    (newFolders) => {
      mutableFolders.value = [...newFolders]
    },
    { deep: true }
  )

  const handleFolderClick = (folder: FolderInfo) => {
    if (props.selectMode && !props.preview) {
      emit('toggle-folder-select', folder)
    } else {
      emit('folder-click', folder)
    }
  }

  const showContextMenu = ref(false)
  const contextMenuPosition = reactive({ x: 0, y: 0 })
  const selectedFolder = ref<FolderInfo | null>(null)
  const contextMenuItems = ref<ContextMenuItem[]>([])

  const buildFolderMenuItem = (folder: FolderInfo): ContextMenuItem => ({
    key: `folder-${folder.id}`,
    label: folder.name,
    icon: 'fas fa-folder',
    onClick: async () => {
      if (!selectedFolder.value) return
      try {
        await moveFolder(selectedFolder.value.id, folder.id)
        toast.success($t('folders.folderGrid.toast.moveSuccess', { name: folder.name }))
        emit('folder-moved', selectedFolder.value.id, folder.id)
      } catch {
        toast.error($t('folders.folderGrid.toast.moveFailed'))
      }
    },
    hasAsyncChildren: folder.has_children,
    loadChildren: folder.has_children
      ? async () => {
          const res = await getFolderList(folder.id)
          if (res.success && res.data) return res.data.map(buildFolderMenuItem)
          return []
        }
      : undefined,
  })

  const loadRootAsMenu = async (): Promise<ContextMenuItem[]> => {
    const res = await getFolderList()
    const list = res.success && res.data ? res.data.map(buildFolderMenuItem) : []
    return [
      {
        key: 'root',
        label: $t('folders.folderGrid.contextMenu.root'),
        icon: 'fas fa-home',
        onClick: async () => {
          if (!selectedFolder.value) return
          try {
            await moveFolder(selectedFolder.value.id, undefined)
            toast.success($t('folders.folderGrid.toast.moveToRootSuccess'))
            emit('folder-moved', selectedFolder.value.id, undefined)
          } catch {
            toast.error($t('folders.folderGrid.toast.moveFailed'))
          }
        },
      },
      ...list,
    ]
  }

  const handleFolderContextMenu = async (event: MouseEvent, folder: FolderInfo) => {
    event.preventDefault()
    event.stopPropagation()
    selectedFolder.value = folder
    contextMenuItems.value = [
      {
        key: 'move-to',
        label: $t('folders.folderGrid.contextMenu.moveTo'),
        icon: 'fas fa-arrows-alt',
        hasAsyncChildren: true,
        loadChildren: loadRootAsMenu,
      },
      { key: 'divider-1', label: '', divided: true },
      {
        key: 'edit',
        label: $t('folders.folderGrid.contextMenu.rename'),
        icon: 'fas fa-edit',
        onClick: () => {
          if (!selectedFolder.value) return
          emit('edit-folder', selectedFolder.value)
        },
      },
      {
        key: 'delete',
        label: $t('folders.folderGrid.contextMenu.delete'),
        icon: 'fas fa-trash',
        onClick: () => {
          if (!selectedFolder.value) return
          emit('delete-folder', selectedFolder.value)
        },
        danger: true,
      },
    ]
    contextMenuPosition.x = event.clientX
    contextMenuPosition.y = event.clientY
    showContextMenu.value = true
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString(getCurrentLocale(), {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }
</script>

<template>
  <div v-if="folders.length > 0" class="section-container">
    <div class="section-header">
      <h2 class="section-title">
        <i class="fas fa-folder" /> {{ $t('folders.folderGrid.sectionTitle') }}
        <span class="item-count">{{ folders.length }}</span>
        <span v-if="!selectMode && !preview && folders.length > 1" class="drag-tip">
          <i class="fas fa-arrows-alt" /> {{ $t('folders.folderGrid.dragTip') }}
        </span>
      </h2>
      <div v-if="selectMode && !preview" class="section-actions">
        <div class="select-actions">
          <CyberButton
            type="ghost"
            :icon="isAllFoldersSelected ? 'check-double' : 'check'"
            :title="
              isAllFoldersSelected
                ? $t('folders.folderGrid.actions.deselectAllTitle')
                : $t('folders.folderGrid.actions.selectAllTitle')
            "
            :class="{ active: isAllFoldersSelected }"
            @click="$emit('toggle-select-all-folders')"
          >
            {{ isAllFoldersSelected ? $t('folders.folderGrid.actions.deselectAll') : $t('folders.folderGrid.actions.selectAll') }}
          </CyberButton>
          <CyberButton
            type="ghost"
            icon="exchange-alt"
            :title="$t('folders.folderGrid.actions.invertTitle')"
            @click="$emit('invert-folder-selection')"
          >
            {{ $t('folders.folderGrid.actions.invert') }}
          </CyberButton>
        </div>
      </div>
    </div>
    <VueDraggable
      v-model="mutableFolders"
      animation="150"
      ghost-class="folder-ghost"
      chosen-class="folder-chosen"
      drag-class="folder-drag"
      :disabled="selectMode || preview"
      class="folder-cards"
      @start="$emit('folder-drag-start')"
      @end="(event) => $emit('folder-drag-end', { event, sortedFolders: mutableFolders })"
    >
      <div
        v-for="folder in mutableFolders"
        :key="folder.id"
        class="folder-card selectable-item"
        :class="{
          selected: selectMode && isFolderSelected?.(folder),
          'is-selected': selectMode && isFolderSelected?.(folder),
          selectable: selectMode,
        }"
        @click="handleFolderClick(folder)"
        @contextmenu="!preview ? handleFolderContextMenu($event, folder) : undefined"
      >
        <div v-if="selectMode && isFolderSelected?.(folder)" class="selected-indicator selectable-check">
          <i class="fas fa-check" />
        </div>

        <div
          v-if="!preview && !selectMode"
          class="folder-visibility"
          :class="folder.permission === 'private' ? 'private' : 'public'"
          :title="$t('folders.folderGrid.visibility.toggleTitle')"
          :aria-label="
            folder.permission === 'private'
              ? $t('folders.folderGrid.visibility.setPublic')
              : $t('folders.folderGrid.visibility.setPrivate')
          "
          @click="$emit('toggle-folder-visibility', folder, $event)"
        >
          <i :class="folder.permission === 'private' ? 'fas fa-lock' : 'fas fa-globe'" />
        </div>

        <div class="folder-card-content">
          <div class="folder-icon-container">
            <div class="folder-icon-wrapper">
              <i class="fas fa-folder" />
            </div>
          </div>
          <div class="folder-details">
            <div class="folder-name" :title="folder.name">{{ folder.name }}</div>
            <div class="folder-meta">
              <div class="folder-meta-item">
                <i class="fas fa-clock" />
                <span class="meta-text">{{ formatDate(folder.created_at) }}</span>
              </div>
              <div v-if="folder.file_count !== undefined" class="folder-meta-item">
                <i class="fas fa-image" />
                <span class="meta-text">{{ $t('folders.folderGrid.fileCount', { count: folder.file_count || 0 }) }}</span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="!preview && !selectMode" class="folder-card-actions">
          <button
            class="folder-action-btn delete-btn"
            :title="$t('folders.folderGrid.actions.deleteTitle')"
            @click.stop="$emit('delete-folder', folder)"
          >
            <i class="fas fa-trash" />
          </button>
          <button class="folder-action-btn edit-btn" @click.stop="$emit('edit-folder', folder)">
            <i class="fas fa-edit" />
          </button>
        </div>

        <div class="folder-highlight" />
      </div>
    </VueDraggable>

    <CyberContextMenu
      v-model="showContextMenu"
      :items="contextMenuItems"
      :x="contextMenuPosition.x"
      :y="contextMenuPosition.y"
      @close="showContextMenu = false"
    />
  </div>
</template>

<style scoped>
  .section-container {
    background: rgba(var(--color-background-800-rgb), 0.95);
    backdrop-filter: blur(10px);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    overflow: visible;
    padding: 1.25rem;
    box-shadow:
      0 4px 12px var(--color-overlay-medium),
      0 2px 6px var(--color-overlay-light);
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding-bottom: 0.6rem;
    border-bottom: 1.5px solid rgba(var(--color-brand-500-rgb), 0.25);
  }

  .section-title {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 1rem;
  }

  .item-count {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    padding: 0.15rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
    margin-left: 0.4rem;
  }

  .drag-tip {
    background: rgba(var(--color-error-rgb), 0.1);
    color: rgba(var(--color-error-rgb), 0.8);
    padding: 0.1rem 0.4rem;
    border-radius: var(--radius-sm);
    font-size: 0.7rem;
    margin-left: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.2rem;
    opacity: 0.8;
    transition: all 0.2s ease;
  }

  .drag-tip:hover {
    opacity: 1;
    background: var(--color-hover-bg);
  }

  .drag-tip i {
    font-size: 0.65rem;
  }

  .section-actions {
    display: flex;
    gap: 0.4rem;
    align-items: center;
  }

  .select-actions {
    display: flex;
    gap: 0.5rem;
    margin-right: 0.5rem;
  }

  .select-mode-tip {
    display: flex;
    align-items: center;
    padding: 0.4rem 0.8rem;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: rgba(var(--color-brand-500-rgb), 0.9);
    font-size: 0.85rem;
    font-weight: 500;
    white-space: nowrap;
    margin-right: 0.8rem;
    animation: tipPulse 2s ease-in-out infinite;
  }

  .select-mode-tip i {
    font-size: 0.9rem;
    color: rgba(var(--color-brand-500-rgb), 0.8);
  }

  @keyframes tipPulse {
    0%,
    100% {
      opacity: 0.8;
      transform: scale(1);
    }
    50% {
      opacity: 1;
      transform: scale(1.02);
    }
  }

  .folder-cards {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(230px, 1fr));
    gap: 15px;
    perspective: 1000px;
  }

  .folder-card {
    position: relative;
    background: var(--color-background-700);
    border-radius: var(--radius-md);
    cursor: grab;
    box-shadow: var(--shadow-cyber-sm);
    border: 1px solid var(--color-border-default);
    min-height: 80px;
    width: 100%;
    height: auto;
    min-width: 0;
    transform-style: preserve-3d;
    transition: all 0.3s ease;
  }

  .folder-card:hover {
    box-shadow: var(--shadow-cyber-lg);
    border-color: var(--color-hover-border);
    background: var(--color-background-700);
  }

  .folder-card:active {
    cursor: grabbing;
  }

  .folder-card.selectable {
    cursor: pointer;
  }

  .folder-card.selected {
    background: var(--color-selected-bg);
    border: 2px solid var(--color-selected-border);
    box-shadow: var(--shadow-glow-lg);
  }

  @keyframes selectedGlow {
    0% {
      transform: scale(1);
    }
    100% {
      transform: scale(1.02);
    }
  }

  .selected-indicator {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 28px;
    height: 28px;
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-badge-accent-text));
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 3;
    box-shadow: var(--shadow-glow-md);
    animation: checkmark-appear 0.3s ease-out;
  }

  .selected-indicator i {
    color: var(--color-content-heading);
    font-size: 14px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  }

  @keyframes checkmark-appear {
    0% {
      transform: scale(0);
      opacity: 0;
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }

  .folder-card-content {
    padding: 0.8rem;
    padding-right: 2.8rem;
    padding-bottom: 2.2rem;
    display: flex;
    align-items: center;
    gap: 0.8rem;
  }

  .folder-icon-container {
    flex-shrink: 0;
  }

  .folder-icon-wrapper {
    width: 44px;
    height: 44px;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15) 0%, rgba(var(--color-brand-500-rgb), 0) 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    transition: all 0.3s ease;
    box-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .folder-icon-wrapper i {
    font-size: 20px;
    color: var(--color-brand-500);
    transition: all 0.3s ease;
    filter: drop-shadow(0 0 6px rgba(var(--color-brand-500-rgb), 0.6));
  }

  .folder-card:hover .folder-icon-wrapper {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.3) 0%, rgba(var(--color-brand-500-rgb), 0.1) 100%);
    transform: scale(1.05);
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .folder-card:hover .folder-icon-wrapper i {
    transform: scale(1.1);
    filter: drop-shadow(0 0 10px rgba(var(--color-brand-500-rgb), 0.8));
  }

  .folder-details {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 0.4rem;
  }

  .folder-name {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--color-content-heading);
    margin-bottom: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-family: 'Orbitron', sans-serif;
    letter-spacing: 0.5px;
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.4);
    line-height: 1.2;
  }

  .folder-meta {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }

  .folder-meta-item {
    display: flex;
    align-items: center;
    font-size: 0.75rem;
    color: var(--color-content-default);
    line-height: 1.3;
  }

  .folder-meta-item i {
    margin-right: 0.4rem;
    color: rgba(var(--color-brand-500-rgb), 0.8);
    font-size: 0.7rem;
    width: 10px;
    text-align: center;
    flex-shrink: 0;
  }

  .folder-meta-item .meta-text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
  }

  .folder-visibility {
    position: absolute;
    right: 8px;
    top: 8px;
    width: 26px;
    height: 26px;
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
    z-index: 2;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
  }

  .folder-visibility.private {
    background: rgba(var(--color-warning-rgb), 0.15);
    color: var(--color-warning-500);
    border: 1px solid rgba(var(--color-warning-rgb), 0.3);
  }

  .folder-visibility.public {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-500);
    border: 1px solid rgba(var(--color-success-rgb), 0.3);
  }

  .folder-visibility.private:hover {
    background: rgba(var(--color-warning-rgb), 0.25);
    color: var(--color-warning-500);
    border-color: rgba(var(--color-warning-rgb), 0.5);
    box-shadow: 0 0 10px rgba(var(--color-warning-rgb), 0.3);
    transform: scale(1.1);
  }

  .folder-visibility.public:hover {
    background: rgba(var(--color-success-rgb), 0.25);
    color: var(--color-success-500);
    border-color: rgba(var(--color-success-rgb), 0.5);
    box-shadow: 0 0 10px rgba(var(--color-success-rgb), 0.3);
    transform: scale(1.1);
  }

  .folder-visibility i {
    font-size: 0.75rem;
  }

  .folder-card-actions {
    position: absolute;
    bottom: 8px;
    right: 8px;
    display: flex;
    gap: 0.4rem;
    opacity: 1;
    z-index: 3;
  }

  .folder-action-btn {
    background: rgba(var(--color-background-700-rgb), 0.8);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-content-muted);
    width: 24px;
    height: 24px;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    font-size: 0.7rem;
  }

  .folder-action-btn.delete-btn {
    color: var(--color-error-500);
    background: rgba(var(--color-error-rgb), 0.1);
    border-color: rgba(var(--color-error-rgb), 0.3);
  }

  .folder-action-btn.delete-btn:hover {
    background: rgba(var(--color-error-rgb), 0.2);
    color: var(--color-error-500);
    border-color: rgba(var(--color-error-rgb), 0.5);
    box-shadow: 0 0 10px rgba(var(--color-error-rgb), 0.4);
    transform: scale(1.05);
  }

  .folder-action-btn.edit-btn {
    color: rgba(var(--color-brand-500-rgb), 0.8);
    border-color: rgba(var(--color-brand-500-rgb), 0.2);
  }

  .folder-action-btn.edit-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: scale(1.05);
  }

  .folder-highlight {
    position: absolute;
    top: 0;
    left: 0;
    width: 3px;
    height: 0;
    background: linear-gradient(to bottom, var(--color-brand-500), var(--color-pink-400));
    transition: height 0.3s ease;
  }

  .folder-card:hover .folder-highlight {
    height: 100%;
  }

  .folder-ghost {
    opacity: 0.4 !important;
    background: rgba(var(--color-brand-500-rgb), 0.15) !important;
    border: 2px dashed rgba(var(--color-brand-500-rgb), 0.6) !important;
    border-radius: var(--radius-sm) !important;
    transform: scale(0.95) !important;
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.3) !important;
  }

  .folder-chosen {
    opacity: 1 !important;
    background: rgba(var(--color-brand-500-rgb), 0.1) !important;
    border-color: rgba(var(--color-brand-500-rgb), 0.5) !important;
    transform: scale(1.02) !important;
    box-shadow: 0 5px 15px rgba(var(--color-brand-500-rgb), 0.2) !important;
    z-index: 10 !important;
  }
</style>
