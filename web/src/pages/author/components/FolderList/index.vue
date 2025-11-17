<script setup lang="ts">
  import { computed } from 'vue'
  import type { AuthorFolderInfo } from '@/api/types/index'
  import { formatRelativeTime } from '@/utils/formatting/format'
  import FolderBreadcrumb from '@/pages/folders/components/FolderBreadcrumb.vue'

  defineOptions({
    name: 'AuthorFolderList',
  })

  interface FolderListTexts {
    rootTitle: string
    rootSubtitle: string
    childTitle: string
    childSubtitle: string
    countLabel: string
    createdLabel: string
    unknownTime?: string
    rootFolder?: string
  }

  interface Props {
    folders: AuthorFolderInfo[]
    currentFolderId?: string
    folderLoading: boolean
    breadcrumbLoading: boolean
    loading: boolean
    texts: FolderListTexts
    breadcrumbTrail?: AuthorFolderInfo[]
    authorName?: string
  }

  interface Emits {
    (e: 'folder-click', folder: AuthorFolderInfo): void
    (e: 'breadcrumb-click', item: { id: string | null; name: string; isRoot?: boolean } | null): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const handleFolderClick = (folder: AuthorFolderInfo) => {
    emit('folder-click', folder)
  }

  const getFormattedDate = (dateString: string | Date): string => {
    if (!dateString) {
      console.warn('Date is empty')
      return props.texts.unknownTime || props.texts.unknownTime
    }
    try {
      const formatted = formatRelativeTime(new Date(dateString).toISOString())
      return formatted
    } catch (e) {
      console.error('Date formatting error:', e, dateString)
      return String(dateString)
    }
  }

  const handleBreadcrumbClick = (item: { id: string | null; name: string; isRoot?: boolean } | null) => {
    emit('breadcrumb-click', item)
  }

  const breadcrumbItems = computed(() => {
    const rootItem = {
      id: null,
      name: props.authorName || props.texts.rootFolder || props.texts.rootFolder,
    } as any

    if (!props.currentFolderId) {
      return [rootItem]
    }

    const trail = props.breadcrumbTrail || []
    const items = [rootItem, ...trail] as any[]
    return items
  })
</script>

<template>
  <div class="section-container">
    <div class="section-header">
      <div class="header-left">
        <h2 class="section-title">
          <div class="title-icon">
            <i class="fas fa-folder-open" />
            <div class="icon-glow" />
          </div>
          <div class="title-content">
            <span class="title-text">{{ currentFolderId ? props.texts.childTitle : props.texts.rootTitle }}</span>
            <span class="title-subtitle">{{ currentFolderId ? props.texts.childSubtitle : props.texts.rootSubtitle }}</span>
          </div>
        </h2>
      </div>
    </div>

    <FolderBreadcrumb
      v-if="breadcrumbItems.length > 0"
      class="folder-breadcrumb-nav"
      :items="breadcrumbItems"
      :show-icon="true"
      :auto-add-root="false"
      @click="handleBreadcrumbClick"
    />

    <div v-if="folders.length > 0" class="folder-cards">
      <div
        v-for="folder in folders"
        :key="folder.id"
        class="folder-card"
        :class="{
          loading: folderLoading && currentFolderId === folder.id,
          disabled: folderLoading || breadcrumbLoading || loading,
          current: currentFolderId === folder.id,
        }"
        @click="handleFolderClick(folder)"
      >
        <div class="folder-card-content">
          <div class="folder-icon-wrapper">
            <i v-if="!folderLoading || currentFolderId !== folder.id" class="fas fa-folder" />
            <i v-else class="fas fa-spinner fa-spin" />
          </div>
          <div class="folder-details">
            <div class="folder-name" :title="folder.name">{{ folder.name }}</div>
            <div class="folder-meta">
              <div class="folder-meta-item">
                <i class="fas fa-clock" />
                <span>{{ props.texts.createdLabel.replace('{time}', getFormattedDate(folder.createdAt)) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="folder-highlight" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss" src="./FolderList.scss"></style>
