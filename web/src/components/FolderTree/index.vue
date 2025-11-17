<script setup lang="ts">
  import { computed, nextTick, onMounted, ref, watch, shallowRef } from 'vue'
  import { getFolderList, searchFolders } from '@/api/folder'
  import type { FolderInfo } from '@/api/types'
  import FolderTreeItem from './FolderTreeItem.vue'
  import { onClickOutside, useDebounceFn, useElementBounding, useEventListener, useWindowSize } from '@vueuse/core'
  import { Z_INDEX } from '@/constants'
  import { useAuthStore } from '@/store/auth'
  import { useFolderPath } from '@/hooks/useFolderPath'
  import { useTexts } from '@/composables/useTexts'
  import type { FolderTreeProps, FolderTreeEmits } from './types'

  defineOptions({
    name: 'CyberFolderTree',
  })

  const props = withDefaults(defineProps<FolderTreeProps>(), {
    modelValue: '',
    borderColor: '',
    presetPath: '',
    maxPathLength: 50,
    width: '',
    height: '30px',
    dropdownZIndex: 0,
  })

  const emit = defineEmits<FolderTreeEmits>()
  const { $t } = useTexts()
  const authStore = useAuthStore()
  const { getFolderPathString } = useFolderPath()

  /* 核心状态 - 使用更少的响应式变量 */
  const state = ref({
    isOpen: false,
    isLoading: false,
    isRequestingSubfolders: false,
    selectedId: props.modelValue || '',
    selectedPath: '',
    searchQuery: '',
  })

  /* 使用shallowRef优化大数组的响应性能 */
  const folders = shallowRef<{
    topLevel: FolderInfo[]
    all: FolderInfo[]
    searchResults: FolderInfo[]
  }>({
    topLevel: [],
    all: [],
    searchResults: [],
  })

  /* 文件夹展开状态和缓存 */
  const expandedFolders = ref<Set<string>>(new Set())
  const foldersChildrenCache = ref<Map<string, FolderInfo[]>>(new Map())

  /* DOM引用 */
  const dropdownRef = ref<HTMLElement | null>(null)
  const searchInputRef = ref<HTMLInputElement | null>(null)

  const { update: updateBounding, top, left, width, bottom } = useElementBounding(dropdownRef, { windowScroll: true })

  const { height: windowHeight } = useWindowSize()

  /* 计算属性 - 优化和简化 */
  const wrapperStyle = computed(() => (props.width ? { width: props.width } : {}))

  const triggerStyle = computed(() => (props.height ? { height: props.height } : {}))

  const isDropdownAbove = computed(() => {
    const dropdownHeight = 330
    const bottomSpace = windowHeight.value - bottom.value
    return bottomSpace < dropdownHeight
  })

  const dropdownStyle = computed(() => {
    const baseStyle = {
      position: 'fixed' as const,
      left: `${left.value}px`,
      width: `${width.value}px`,
      zIndex: props.dropdownZIndex || Z_INDEX.DROPDOWN,
    }

    return isDropdownAbove.value
      ? { ...baseStyle, bottom: `${windowHeight.value - top.value}px` }
      : { ...baseStyle, top: `${bottom.value}px` }
  })

  const effectiveBorderColor = computed(() => props.borderColor || 'var(--color-brand-500)')

  /* 路径截断逻辑优化 */
  const displayPath = computed(() => {
    if (props.presetPath) return props.presetPath
    if (!state.value.selectedPath) return $t('components.folderTree.selectFolder')

    const path = state.value.selectedPath
    if (path.length <= props.maxPathLength) return path

    const parts = path.split(' / ')
    if (parts.length <= 1) {
      return `${path.substring(0, props.maxPathLength - 3)}...`
    }

    const firstPart = parts[0]
    const lastPart = parts[parts.length - 1]
    const truncated = `${firstPart} / ... / ${lastPart}`

    if (truncated.length <= props.maxPathLength) return truncated

    const maxLastPartLength = props.maxPathLength - firstPart.length - 7
    if (maxLastPartLength > 0) {
      const shortLastPart = lastPart.length > maxLastPartLength ? `${lastPart.substring(0, maxLastPartLength - 3)}...` : lastPart
      return `${firstPart} / ... / ${shortLastPart}`
    }

    return `${firstPart.substring(0, props.maxPathLength - 13)}... / ...`
  })

  const filteredFolders = computed(() => (state.value.searchQuery.trim() ? folders.value.searchResults : folders.value.topLevel))

  const createDefaultFolder = (): FolderInfo => ({
    id: '',
    name: $t('components.folderTree.rootDirectory'),
    fullPath: '/',
    children: [],
  })

  const buildFolderPaths = (folderList: FolderInfo[], parentPath = ''): FolderInfo[] =>
    folderList.map((folder) => ({
      ...folder,
      fullPath: parentPath ? `${parentPath} / ${folder.name}` : folder.name,
    }))

  const updatePosition = () => {
    if (state.value.isOpen) updateBounding()
  }

  const toggleDropdown = () => {
    state.value.isOpen = !state.value.isOpen
    if (state.value.isOpen) {
      nextTick(() => {
        searchInputRef.value?.focus()
        updateBounding()
      })
    }
  }

  const closeDropdown = () => {
    state.value.isOpen = false
  }

  const selectFolder = (folder: FolderInfo) => {
    state.value.selectedId = folder.id
    state.value.selectedPath = folder.fullPath || folder.name

    emit('update:modelValue', folder.id)
    emit('folderSelected', {
      id: folder.id,
      name: folder.name,
      path: folder.fullPath || folder.name,
    })

    closeDropdown()

    setTimeout(() => {
      state.value.searchQuery = ''
      folders.value.searchResults = []
    }, 200)
  }

  const loadFolders = async () => {
    try {
      state.value.isLoading = true
      const result = await getFolderList()

      if (result.success && Array.isArray(result.data)) {
        const processedFolders = buildFolderPaths(result.data)
        folders.value = {
          topLevel: processedFolders,
          all: [...processedFolders],
          searchResults: [],
        }

        if (state.value.selectedId) {
          await loadSelectedFolderPath()
        }
      } else {
        const defaultFolder = createDefaultFolder()
        folders.value = {
          topLevel: [defaultFolder],
          all: [defaultFolder],
          searchResults: [],
        }
      }
    } catch (error) {
      console.error('[FolderTree] Failed to load folders:', error)
      const defaultFolder = createDefaultFolder()
      folders.value = {
        topLevel: [defaultFolder],
        all: [defaultFolder],
        searchResults: [],
      }
    } finally {
      state.value.isLoading = false
    }
  }

  const loadSelectedFolderPath = async () => {
    if (!state.value.selectedId) return

    if (props.presetPath) {
      state.value.selectedPath = props.presetPath
      return
    }

    try {
      state.value.selectedPath = await getFolderPathString(state.value.selectedId)
    } catch (error) {
      console.error('[FolderTree] Failed to load folder path:', error)
    }
  }

  const searchFoldersAPI = async () => {
    const keyword = state.value.searchQuery.trim()
    if (!keyword) {
      folders.value.searchResults = []
      return
    }

    if (!authStore.isAuthenticated) {
      folders.value.searchResults = []
      return
    }

    try {
      const response = await searchFolders(keyword)
      if (response.success && response.data) {
        folders.value.searchResults = response.data.map((folder) => ({
          ...folder,
          fullPath: folder.description || folder.name,
        }))
      } else {
        folders.value.searchResults = []
      }
    } catch (error) {
      console.error('[FolderTree] Search folders failed:', error)
      folders.value.searchResults = []
    }
  }

  const debouncedSearch = useDebounceFn(searchFoldersAPI, 300)

  const handleFolderToggle = (folderId: string, isExpanded: boolean) => {
    if (isExpanded) {
      expandedFolders.value.add(folderId)
    } else {
      expandedFolders.value.delete(folderId)
    }
  }

  const handleChildrenLoaded = (folderId: string, children: FolderInfo[]) => {
    foldersChildrenCache.value.set(folderId, children)
    folders.value.all = [...folders.value.all, ...children]
  }

  const getFolderChildren = (folderId: string): FolderInfo[] => foldersChildrenCache.value.get(folderId) || []

  useEventListener(window, 'scroll', updatePosition, { passive: true })
  useEventListener(window, 'resize', updatePosition, { passive: true })

  onClickOutside(dropdownRef, (event) => {
    const target = event.target as HTMLElement
    const isClickInDropdown = target.closest('.cyber-folder-tree__dropdown') !== null

    if (state.value.isOpen && !state.value.isRequestingSubfolders && !isClickInDropdown) {
      closeDropdown()
    }
  })

  watch(
    () => props.modelValue,
    (newValue) => {
      if (newValue !== state.value.selectedId) {
        state.value.selectedId = newValue
        loadSelectedFolderPath()
      }
    }
  )

  watch(
    () => props.presetPath,
    (newPath) => {
      if (newPath && state.value.selectedId) {
        state.value.selectedPath = newPath
      }
    }
  )

  watch(() => state.value.searchQuery, debouncedSearch)

  watch(
    () => state.value.isOpen,
    (isOpen) => {
      if (isOpen) {
        nextTick(updateBounding)
      }
    }
  )

  onMounted(async () => {
    await loadFolders()
  })
</script>

<template>
  <div class="cyber-folder-tree" :class="{ 'cyber-folder-tree--open': state.isOpen }" :style="wrapperStyle" ref="dropdownRef">
    <div
      class="cyber-folder-tree__trigger"
      :class="{ 'cyber-folder-tree__trigger--active': state.isOpen }"
      :style="[triggerStyle, { '--border-color': effectiveBorderColor }]"
      @click="toggleDropdown"
    >
      <div class="cyber-folder-tree__display">
        <i class="fas fa-folder cyber-folder-tree__icon" :style="{ color: effectiveBorderColor }" />
        <span
          class="cyber-folder-tree__text"
          :class="{ 'cyber-folder-tree__text--placeholder': !state.selectedPath }"
          :title="state.selectedPath || displayPath"
        >
          {{ displayPath }}
        </span>
      </div>
      <i class="fas fa-chevron-down cyber-folder-tree__arrow" :class="{ 'cyber-folder-tree__arrow--rotated': state.isOpen }" />
    </div>

    <Teleport to="body">
      <Transition name="cyber-folder-tree-dropdown" @enter="updateBounding" @after-enter="updateBounding">
        <div
          v-if="state.isOpen"
          class="cyber-folder-tree__dropdown"
          :class="{ 'cyber-folder-tree__dropdown--above': isDropdownAbove }"
          :style="dropdownStyle"
        >
          <div class="cyber-folder-tree__search">
            <div class="cyber-folder-tree__search-wrapper">
              <i class="fas fa-search cyber-folder-tree__search-icon" />
              <input
                ref="searchInputRef"
                v-model="state.searchQuery"
                type="text"
                class="cyber-folder-tree__search-input"
                :placeholder="$t('components.folderTree.searchPlaceholder')"
                @keydown.escape="closeDropdown"
              />
              <i v-if="state.searchQuery" class="fas fa-times cyber-folder-tree__search-clear" @click="state.searchQuery = ''" />
            </div>
          </div>

          <div class="cyber-folder-tree__list">
            <div v-if="state.isLoading" class="cyber-folder-tree__loading">
              <i class="fas fa-spinner fa-spin" />
              <span>{{ $t('common.loading') }}...</span>
            </div>

            <div v-else-if="!filteredFolders.length" class="cyber-folder-tree__empty">
              <span>{{
                state.searchQuery.trim() ? $t('components.folderTree.noMatches') : $t('components.folderTree.noFolders')
              }}</span>
            </div>

            <div v-else class="cyber-folder-tree__items">
              <div
                v-if="!state.searchQuery.trim()"
                class="cyber-folder-tree__item cyber-folder-tree__item--root"
                :class="{ 'cyber-folder-tree__item--selected': state.selectedId === '' }"
                @click="selectFolder(createDefaultFolder())"
              >
                <div class="cyber-folder-tree__item-content">
                  <div class="cyber-folder-tree__item-icon">
                    <i class="fas fa-home" />
                  </div>
                  <span class="cyber-folder-tree__item-name">
                    {{ $t('components.folderTree.rootDirectory') }}
                  </span>
                </div>
              </div>

              <template v-if="state.searchQuery.trim()">
                <div
                  v-for="folder in filteredFolders"
                  :key="folder.id"
                  class="cyber-folder-tree__item cyber-folder-tree__item--search"
                  :class="{ 'cyber-folder-tree__item--selected': state.selectedId === folder.id }"
                  @click="selectFolder(folder)"
                >
                  <div class="cyber-folder-tree__item-content">
                    <div class="cyber-folder-tree__item-icon">
                      <i class="fas fa-folder" :style="{ color: effectiveBorderColor }" />
                    </div>
                    <div class="cyber-folder-tree__item-info">
                      <span class="cyber-folder-tree__item-name">{{ folder.name }}</span>
                      <span class="cyber-folder-tree__item-path" :title="folder.fullPath">
                        {{ folder.fullPath }}
                      </span>
                    </div>
                  </div>
                </div>
              </template>

              <!-- 文件夹树 -->
              <template v-else>
                <FolderTreeItem
                  v-for="folder in filteredFolders"
                  :key="folder.id"
                  :folder="folder"
                  :level="0"
                  :selected-id="state.selectedId"
                  :border-color="effectiveBorderColor"
                  :is-expanded="expandedFolders.has(folder.id)"
                  :expanded-folders="expandedFolders"
                  :cached-children="getFolderChildren(folder.id)"
                  :get-folder-children="getFolderChildren"
                  @select="selectFolder"
                  @toggle="handleFolderToggle"
                  @children-loaded="handleChildrenLoaded"
                  @subfolder-request-start="state.isRequestingSubfolders = true"
                  @subfolder-request-end="state.isRequestingSubfolders = false"
                />
              </template>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
  .cyber-folder-tree {
    position: relative;
    display: inline-block;
    width: 100%;
  }

  .cyber-folder-tree__trigger {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    height: 30px;
    min-height: 30px;
    padding: 0 8px;
    border: 1.5px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    background: var(--color-background-700);
    cursor: pointer;
    transition:
      border-color 0.25s ease,
      box-shadow 0.25s ease,
      transform 0.25s ease;
    font-size: 0.875rem;
    color: var(--color-content);
    box-sizing: border-box;

    box-shadow:
      2px 2.5px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 0 rgba(var(--color-brand-500-rgb), 0);
  }

  .cyber-folder-tree__trigger:hover {
    border-color: var(--color-border-strong);

    box-shadow:
      2.5px 3px 0 rgba(var(--color-brand-500-rgb), 0.3),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.08);
    transform: translate(-0.3px, -0.3px);
  }

  .cyber-folder-tree__trigger--active {
    border-color: var(--color-brand-500);

    box-shadow:
      3.5px 4px 0 rgba(var(--color-brand-500-rgb), 0.4),
      0 0 12px rgba(var(--color-brand-500-rgb), 0.2);
    transform: translate(-0.5px, -0.5px);
  }

  .cyber-folder-tree__trigger:active {
    transform: translate(0.5px, 0.5px);
    box-shadow:
      1.5px 2px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 6px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-folder-tree__display {
    display: flex;
    align-items: center;
    min-width: 0;
    flex: 1;
    gap: 6px;
  }

  .cyber-folder-tree__icon {
    font-size: 12px;
    flex-shrink: 0;
  }

  .cyber-folder-tree__text {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 0.875rem;
  }

  .cyber-folder-tree__text--placeholder {
    color: rgba(var(--color-content-rgb), 0.6);
    font-style: italic;
  }

  .cyber-folder-tree__arrow {
    font-size: 10px;
    color: rgba(var(--color-content-rgb), 0.6);
    transition: transform 0.2s ease;
    flex-shrink: 0;
    margin-left: 4px;
  }

  .cyber-folder-tree__arrow--rotated {
    transform: rotate(180deg);
  }

  .cyber-folder-tree__dropdown {
    background: var(--color-background-900);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-2xl);
    overflow: hidden;
    z-index: 9999;
    max-height: 330px;
    min-width: 200px;
  }

  .cyber-folder-tree__search {
    padding: 8px;
    border-bottom: 1px solid var(--color-border-subtle);
    background: rgba(var(--color-background-800-rgb), 0.5);
  }

  .cyber-folder-tree__search-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }

  .cyber-folder-tree__search-icon {
    position: absolute;
    left: 8px;
    font-size: 12px;
    color: rgba(var(--color-content-rgb), 0.5);
    z-index: 1;
  }

  .cyber-folder-tree__search-input {
    width: 100%;
    height: 30px;
    padding: 0 32px 0 28px;
    border: 1.5px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    background: var(--color-background-700);
    color: var(--color-content);
    font-size: 0.8125rem;
    outline: none;
    transition:
      border-color 0.25s ease,
      box-shadow 0.25s ease,
      transform 0.25s ease;

    box-shadow:
      1.5px 2px 0 rgba(var(--color-brand-500-rgb), 0.2),
      0 0 0 rgba(var(--color-brand-500-rgb), 0);
  }

  .cyber-folder-tree__search-input:hover {
    border-color: var(--color-border-strong);
    box-shadow:
      2px 2.5px 0 rgba(var(--color-brand-500-rgb), 0.25),
      0 0 4px rgba(var(--color-brand-500-rgb), 0.08);
    transform: translate(-0.3px, -0.3px);
  }

  .cyber-folder-tree__search-input:focus {
    border-color: var(--color-brand-500);
    box-shadow:
      3px 3.5px 0 rgba(var(--color-brand-500-rgb), 0.35),
      0 0 10px rgba(var(--color-brand-500-rgb), 0.2);
    transform: translate(-0.5px, -0.5px);
  }

  .cyber-folder-tree__search-input::placeholder {
    color: rgba(var(--color-content-rgb), 0.5);
  }

  .cyber-folder-tree__search-clear {
    position: absolute;
    right: 8px;
    font-size: 10px;
    color: rgba(var(--color-content-rgb), 0.5);
    cursor: pointer;
    z-index: 1;
    transition: color 0.2s ease;
  }

  .cyber-folder-tree__search-clear:hover {
    color: var(--color-error-500);
  }

  .cyber-folder-tree__list {
    max-height: 280px;
    overflow-y: auto;
  }

  .cyber-folder-tree__loading,
  .cyber-folder-tree__empty {
    padding: 16px;
    text-align: center;
    color: rgba(var(--color-content-rgb), 0.5);
    font-size: 0.8125rem;
  }

  .cyber-folder-tree__loading {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .cyber-folder-tree__items {
    padding: 4px 0;
  }

  .cyber-folder-tree__item {
    cursor: pointer;
    transition: background-color 0.2s ease;
  }

  .cyber-folder-tree__item:hover {
    background: rgba(var(--color-brand-500-rgb), 0.08);
  }

  .cyber-folder-tree__item--selected {
    background: rgba(var(--color-brand-500-rgb), 0.15);
  }

  .cyber-folder-tree__item--selected:hover {
    background: rgba(var(--color-brand-500-rgb), 0.22);
  }

  .cyber-folder-tree__item--root {
    border-bottom: 1px solid var(--color-border-subtle);
  }

  .cyber-folder-tree__item-content {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    gap: 8px;
  }

  .cyber-folder-tree__item-icon {
    font-size: 12px;
    width: 16px;
    text-align: center;
    flex-shrink: 0;
  }

  .cyber-folder-tree__item-name {
    font-size: 0.875rem;
    color: var(--color-content);
  }

  .cyber-folder-tree__item-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
    flex: 1;
  }

  .cyber-folder-tree__item-path {
    font-size: 0.75rem;
    color: rgba(var(--color-content-rgb), 0.6);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .cyber-folder-tree-dropdown-enter-active,
  .cyber-folder-tree-dropdown-leave-active {
    transition: all 0.2s ease;
  }

  .cyber-folder-tree-dropdown-enter-from {
    opacity: 0;
    transform: translateY(-8px);
  }

  .cyber-folder-tree-dropdown-leave-to {
    opacity: 0;
    transform: translateY(-8px);
  }

  .cyber-folder-tree__dropdown--above .cyber-folder-tree-dropdown-enter-from {
    transform: translateY(8px);
  }

  .cyber-folder-tree__dropdown--above .cyber-folder-tree-dropdown-leave-to {
    transform: translateY(8px);
  }

  .cyber-folder-tree__list::-webkit-scrollbar {
    width: 6px;
  }

  .cyber-folder-tree__list::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-folder-tree__list::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
  }

  .cyber-folder-tree__list::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.5);
  }
</style>
