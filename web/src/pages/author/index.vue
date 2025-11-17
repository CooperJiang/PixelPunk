<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import { getAuthorFolder, getAuthorHomepage } from '@/api/author'
  import type { AuthorFolderInfo, AuthorImageInfo, AuthorInfo, AuthorShareInfo, AuthorStats } from '@/api/types'
  import type { ImageInfo } from '@/components/FileViewer/types'
  import NavbarBrand from '@/components/Navbar/NavbarBrand.vue'
  import Loading from '@/components/Loading/index.vue'
  import { useToast } from '@/components/Toast/useToast'
  import { downloadFileQuick } from '@/utils/file/downloader'
  import { useTexts } from '@/composables/useTexts'
  import AuthorHeader from './components/AuthorHeader/index.vue'
  import AuthorFolderList from './components/FolderList/index.vue'
  import FilesGridView from '@/pages/folders/components/FileGrid/FilesGridView.vue'
  import AuthorShareList from './components/ShareList/index.vue'
  import AuthorEmptyState from './components/EmptyState/index.vue'

  defineOptions({
    name: 'AuthorPage',
  })

  const route = useRoute()
  const toast = useToast()
  const { $t } = useTexts()

  interface BreadcrumbItem {
    id: string | null
    name: string
    isRoot?: boolean
  }

  interface PaginationState {
    currentPage: number
    perPage: number
    total: number
    lastPage: number
  }

  const loading = ref(true)
  const imageLoading = ref(false)
  const folderLoading = ref(false)
  const breadcrumbLoading = ref(false)
  const isFullscreen = ref(false)
  const error = ref('')

  const authorInfo = ref<AuthorInfo | null>(null)
  const stats = ref<AuthorStats>({
    totalFiles: 0,
    totalFolders: 0,
    totalShares: 0,
    totalViews: 0,
  })
  const folders = ref<AuthorFolderInfo[]>([])
  const images = ref<AuthorImageInfo[]>([])
  const shares = ref<AuthorShareInfo[]>([])
  const pagination = ref<PaginationState | null>(null)
  const currentFolderId = ref<string | undefined>()
  const breadcrumbTrail = ref<AuthorFolderInfo[]>([])

  const searchQuery = ref('')

  const previewDialogVisible = ref(false)
  const previewImage = ref<ImageInfo | null>(null)
  const previewImages = ref<ImageInfo[]>([])
  const previewIndex = ref(0)

  const authorId = computed(() => route.params.id as string | undefined)

  const headerTexts = computed(() => ({
    joinedLabel: $t('author.header.joinedLabel'),
    websiteLabel: $t('author.header.websiteLabel'),
    stats: {
      views: $t('author.header.stats.views'),
      shares: $t('author.header.stats.shares'),
      images: $t('author.header.stats.images'),
    },
  }))

  const folderTexts = computed(() => ({
    rootTitle: $t('author.folders.rootTitle'),
    rootSubtitle: $t('author.folders.rootSubtitle'),
    childTitle: $t('author.folders.childTitle'),
    childSubtitle: $t('author.folders.childSubtitle'),
    countLabel: $t('author.folders.countLabel'),
    createdLabel: $t('author.folders.createdLabel'),
    unknownTime: $t('author.folders.unknownTime'),
    rootFolder: $t('author.folders.rootFolder'),
  }))

  const shareTexts = computed(() => ({
    sectionTitle: $t('author.shares.sectionTitle'),
    sectionSubtitle: $t('author.shares.sectionSubtitle'),
    statusBadge: $t('author.shares.statusBadge'),
    meta: {
      viewsLabel: $t('author.shares.meta.viewsLabel'),
      createdLabel: $t('author.shares.meta.createdLabel'),
      keyLabel: $t('author.shares.meta.keyLabel'),
    },
    empty: {
      title: $t('author.shares.empty.title'),
      description: $t('author.shares.empty.description'),
    },
  }))

  const emptyTexts = computed(() => ({
    title: $t('author.empty.title'),
    description: $t('author.empty.description'),
    suggestion: $t('author.empty.suggestion'),
    statusCode: 'INFO_204',
    retryText: $t('author.empty.retry'),
  }))

  const errorTexts = computed(() => ({
    title: $t('author.error.title'),
    retry: $t('author.error.retry'),
  }))

  const toolbarTexts = computed(() => ({
    searchPlaceholder: $t('author.toolbar.searchPlaceholder'),
    refresh: $t('author.toolbar.refresh'),
    fullscreenEnter: $t('author.toolbar.fullscreenEnter'),
    fullscreenExit: $t('author.toolbar.fullscreenExit'),
  }))

  const filteredFolders = computed(() => {
    if (!searchQuery.value) {
      return folders.value
    }
    const keyword = searchQuery.value.toLowerCase()
    return folders.value.filter((folder) => folder.name?.toLowerCase().includes(keyword))
  })

  const filteredImages = computed(() => {
    if (!searchQuery.value) {
      return images.value
    }
    const keyword = searchQuery.value.toLowerCase()
    return images.value.filter((image) => image.originalName?.toLowerCase().includes(keyword))
  })

  const filteredShares = computed(() => {
    if (!searchQuery.value) {
      return shares.value
    }
    const keyword = searchQuery.value.toLowerCase()
    return shares.value.filter((share) => share.name?.toLowerCase().includes(keyword))
  })

  const isAllEmpty = computed(
    () => filteredFolders.value.length === 0 && filteredImages.value.length === 0 && filteredShares.value.length === 0
  )

  const isFolderContentEmpty = computed(() => filteredFolders.value.length === 0 && filteredImages.value.length === 0)

  const loadAuthorHomepage = async () => {
    if (!authorId.value) {
      error.value = $t('author.error.invalidAuthor')
      loading.value = false
      return
    }

    loading.value = true
    imageLoading.value = false
    error.value = ''

    try {
      const response = await getAuthorHomepage(authorId.value)
      const data = response.data
      authorInfo.value = data.author ? normalizeAuthorInfo(data.author) : null
      stats.value = data.stats
        ? normalizeStats(data.stats)
        : {
            totalFiles: 0,
            totalFolders: 0,
            totalShares: 0,
            totalViews: 0,
          }
      folders.value = (data.folders || []).map(normalizeFolder)
      shares.value = (data.shares || []).map(normalizeShare)
      images.value = (data.files || data.images || []).map(normalizeImage)
      pagination.value = data.pagination ? normalizePagination(data.pagination) : null
      breadcrumbTrail.value = []
      currentFolderId.value = undefined
    } catch (err) {
      error.value = extractErrorMessage(err, $t('author.error.loadFailed'))
    } finally {
      loading.value = false
      folderLoading.value = false
      breadcrumbLoading.value = false
    }
  }

  const loadFolderContent = async (folderId: string, page?: number) => {
    if (!authorId.value) {
      return
    }

    imageLoading.value = true
    error.value = ''

    try {
      const response = await getAuthorFolder(authorId.value, folderId, {
        page: page ?? pagination.value?.currentPage ?? 1,
        size: pagination.value?.perPage ?? 50,
      })
      const data = response.data
      folders.value = (data.subFolders || []).map(normalizeFolder)
      images.value = (data.files || data.images || []).map(normalizeImage)
      pagination.value = data.pagination ? normalizePagination(data.pagination) : null
    } catch (err) {
      error.value = extractErrorMessage(err, $t('author.error.loadFolderFailed'))
    } finally {
      imageLoading.value = false
      folderLoading.value = false
      breadcrumbLoading.value = false
      loading.value = false
    }
  }

  const handleFolderClick = async (folder: AuthorFolderInfo) => {
    if (loading.value || folderLoading.value || breadcrumbLoading.value || !folder.id) {
      return
    }
    if (currentFolderId.value === folder.id) {
      return
    }

    folderLoading.value = true
    currentFolderId.value = folder.id

    const existingIndex = breadcrumbTrail.value.findIndex((item) => item.id === folder.id)
    if (existingIndex === -1) {
      breadcrumbTrail.value = [...breadcrumbTrail.value, folder]
    } else {
      breadcrumbTrail.value = breadcrumbTrail.value.slice(0, existingIndex + 1)
    }

    await loadFolderContent(folder.id)
  }

  const handleBreadcrumbClick = async (item: BreadcrumbItem | null) => {
    if (loading.value || folderLoading.value || breadcrumbLoading.value) {
      return
    }

    if (!item || item.id === null) {
      if (!currentFolderId.value) {
        return
      }
      breadcrumbLoading.value = true
      await loadAuthorHomepage()
      return
    }

    breadcrumbLoading.value = true
    currentFolderId.value = item.id
    const index = breadcrumbTrail.value.findIndex((breadcrumb) => breadcrumb.id === item.id)
    breadcrumbTrail.value = breadcrumbTrail.value.slice(0, index + 1)
    await loadFolderContent(item.id)
  }

  const handleRefresh = () => {
    if (currentFolderId.value) {
      loadFolderContent(currentFolderId.value)
    } else {
      loadAuthorHomepage()
    }
  }

  const handleRetry = () => {
    handleRefresh()
  }

  const toggleFullscreen = async () => {
    try {
      if (!document.fullscreenElement) {
        await document.documentElement.requestFullscreen()
        isFullscreen.value = true
      } else {
        await document.exitFullscreen()
        isFullscreen.value = false
      }
    } catch (err) {
      console.error(err)
    }
  }

  const handleImageClick = (image: ImageInfo) => {
    if (!image) {
      console.warn('Image data is empty')
      return
    }

    previewImages.value = filteredImages.value
    const index = filteredImages.value.findIndex((item) => item.id === image.id)
    previewIndex.value = index >= 0 ? index : 0
    previewImage.value = image
    previewDialogVisible.value = true
  }

  const handleImageDownload = (image: ImageInfo) => {
    downloadFileQuick(image.id, image.original_name || image.display_name)
  }

  const handleCopyLink = async (image: ImageInfo) => {
    const url = image.full_url || image.url
    if (!url) {
      toast.error($t('author.toast.copyFailed'))
      return
    }
    try {
      await navigator.clipboard.writeText(url)
      toast.success($t('author.toast.copySuccess'))
    } catch (err) {
      console.error(err)
      toast.error($t('author.toast.copyFailed'))
    }
  }

  const handleShareClick = (share: AuthorShareInfo) => {
    if (!share?.shareKey) {
      return
    }
    const shareUrl = `/share/${share.shareKey}`
    window.open(shareUrl, '_blank')
  }

  const handlePreviewClose = () => {
    previewDialogVisible.value = false
    previewImage.value = null
  }

  const normalizeAuthorInfo = (payload: any): AuthorInfo => ({
    id: Number(payload.id),
    username: payload.username ?? '',
    avatar: payload.avatar ?? payload.avatar_url ?? '',
    avatarFullPath: payload.avatarFullPath ?? payload.avatar_full_path ?? payload.avatar ?? '',
    bio: payload.bio ?? '',
    website: payload.website ?? '',
    createdAt: payload.createdAt ?? payload.created_at ?? '',
    daysJoined: payload.daysJoined ?? payload.days_joined ?? 0,
  })

  const normalizeStats = (payload: any): AuthorStats => ({
    totalFiles: payload.totalFiles ?? payload.total_files ?? 0,
    totalFolders: payload.totalFolders ?? payload.total_folders ?? 0,
    totalShares: payload.totalShares ?? payload.total_shares ?? 0,
    totalViews: payload.totalViews ?? payload.total_views ?? 0,
  })

  const normalizeFolder = (payload: any): AuthorFolderInfo => ({
    id: String(payload.id),
    name: payload.name ?? payload.folderName ?? payload.folder_name ?? '',
    imageCount: payload.fileCount ?? payload.file_count ?? payload.imageCount ?? payload.image_count ?? 0,
    coverImage: payload.coverFile ?? payload.cover_file ?? payload.coverImage ?? payload.cover_image ?? '',
    coverImageFullPath:
      payload.coverFileFullPath ??
      payload.cover_file_full_path ??
      payload.coverImageFullPath ??
      payload.cover_image_full_path ??
      '',
    coverImageThumbURL:
      payload.coverFileThumbURL ??
      payload.cover_file_thumb_url ??
      payload.coverImageThumbURL ??
      payload.cover_image_thumb_url ??
      '',
    coverImageFullThumbURL:
      payload.coverFileFullThumbURL ??
      payload.cover_file_full_thumb_url ??
      payload.coverImageFullThumbURL ??
      payload.cover_image_full_thumb_url ??
      '',
    createdAt: payload.createdAt ?? payload.created_at ?? '',
    totalSize: payload.totalSize ?? payload.total_size ?? 0,
    totalSizeFormatted: payload.totalSizeFormatted ?? payload.total_size_formatted ?? '',
  })

  const normalizeImage = (payload: any): ImageInfo => ({
    id: String(payload.id),
    url: payload.url ?? '',
    thumb_url: payload.thumb_url ?? '',
    full_url: payload.full_url ?? '',
    full_thumb_url: payload.full_thumb_url ?? '',
    original_name: payload.original_name ?? payload.originalName ?? payload.file_name ?? payload.fileName ?? '',
    display_name: payload.display_name ?? payload.displayName ?? payload.original_name ?? payload.originalName ?? '',
    size: payload.size ?? 0,
    width: payload.width ?? 0,
    height: payload.height ?? 0,
    format: payload.format ?? '',
    access_level: payload.access_level ?? payload.accessLevel ?? 'public',
    views: payload.views ?? 0,
    description: payload.description ?? payload.ai_info?.description ?? undefined,
    ai_info: payload.ai_info ?? undefined,
    created_at: payload.created_at ?? payload.createdAt ?? '',
    updated_at: payload.updated_at ?? payload.updatedAt ?? '',
  })

  const normalizeShare = (payload: any): AuthorShareInfo => ({
    id: String(payload.id),
    name: payload.name ?? '',
    description: payload.description ?? '',
    shareKey: payload.shareKey ?? payload.share_key ?? '',
    views: payload.views ?? 0,
    createdAt: payload.createdAt ?? payload.created_at ?? '',
    coverImage: payload.coverFile ?? payload.cover_file ?? payload.coverImage ?? payload.cover_image ?? '',
    coverImageFullPath:
      payload.coverFileFullPath ??
      payload.cover_file_full_path ??
      payload.coverImageFullPath ??
      payload.cover_image_full_path ??
      '',
    coverImageThumbURL:
      payload.coverFileThumbURL ??
      payload.cover_file_thumb_url ??
      payload.coverImageThumbURL ??
      payload.cover_image_thumb_url ??
      '',
    coverImageFullThumbURL:
      payload.coverFileFullThumbURL ??
      payload.cover_file_full_thumb_url ??
      payload.coverImageFullThumbURL ??
      payload.cover_image_full_thumb_url ??
      '',
  })

  const normalizePagination = (payload: any): PaginationState => ({
    currentPage: payload.currentPage ?? payload.current_page ?? 1,
    perPage: payload.perPage ?? payload.per_page ?? payload.size ?? 20,
    total: payload.total ?? 0,
    lastPage: payload.lastPage ?? payload.last_page ?? 1,
  })

  const extractErrorMessage = (err: unknown, fallback: string): string => {
    if (typeof err === 'string') {
      return err
    }
    if (err instanceof Error) {
      const axiosError = err as Error & {
        response?: {
          data?: {
            message?: string
          }
        }
      }
      return axiosError.response?.data?.message || axiosError.message || fallback
    }
    if (err && typeof err === 'object' && 'message' in err) {
      return (err as { message?: string }).message || fallback
    }
    return fallback
  }

  watch(
    () => authorId.value,
    async (nextId) => {
      if (!nextId) {
        error.value = $t('author.error.invalidAuthor')
        loading.value = false
        return
      }
      await loadAuthorHomepage()
    },
    { immediate: true }
  )

  onMounted(() => {
    if (!authorId.value) {
      loading.value = false
      error.value = $t('author.error.invalidAuthor')
    }
  })
</script>

<template>
  <div class="author-page">
    <nav class="author-page__nav">
      <div class="author-page__nav-brand">
        <NavbarBrand />
      </div>
      <div class="author-page__nav-actions">
        <button
          class="nav-action"
          :title="toolbarTexts.refresh"
          :disabled="loading || folderLoading || imageLoading"
          @click="handleRefresh"
        >
          <i class="fas fa-sync-alt" />
        </button>
        <button
          class="nav-action"
          :title="isFullscreen ? toolbarTexts.fullscreenExit : toolbarTexts.fullscreenEnter"
          @click="toggleFullscreen"
        >
          <i :class="isFullscreen ? 'fas fa-compress-arrows-alt' : 'fas fa-expand-arrows-alt'" />
        </button>
      </div>
      <div class="author-page__nav-progress" :class="{ 'is-active': loading || imageLoading }">
        <div class="author-page__nav-progress-bar" />
      </div>
    </nav>

    <main class="author-page__main">
      <Loading :visible="loading && !authorInfo" full-screen :text="$t('author.loading')" />

      <cyberFileViewer
        v-model="previewDialogVisible"
        :file="previewImage"
        :files="previewImages"
        :initial-index="previewIndex"
        :show-side-nav="true"
        @close="handlePreviewClose"
      />

      <section v-if="authorInfo" class="author-page__section">
        <AuthorHeader :author-info="authorInfo" :stats="stats" :texts="headerTexts" />
      </section>

      <section class="author-page__toolbar">
        <CyberInput
          v-model="searchQuery"
          type="text"
          :placeholder="toolbarTexts.searchPlaceholder"
          prefix-icon="search"
          :clearable="true"
          width="350px"
          class="author-page__search"
        />
      </section>

      <section class="author-page__content">
        <AuthorFolderList
          :folders="filteredFolders"
          :current-folder-id="currentFolderId"
          :folder-loading="folderLoading"
          :breadcrumb-loading="breadcrumbLoading"
          :loading="loading"
          :texts="folderTexts"
          :breadcrumb-trail="breadcrumbTrail"
          :author-name="authorInfo?.username"
          @folder-click="handleFolderClick"
          @breadcrumb-click="handleBreadcrumbClick"
        />

        <div v-if="filteredImages.length > 0" class="files-section">
          <FilesGridView
            :images="filteredImages"
            :select-mode="false"
            :batch-mode="false"
            :preview="true"
            @preview-image="handleImageClick"
            @copy-link="handleCopyLink"
            @download-image="handleImageDownload"
          />
        </div>

        <AuthorEmptyState
          v-if="!loading && !imageLoading && !error && currentFolderId && isFolderContentEmpty"
          :title="$t('author.empty.folderEmpty.title')"
          :description="$t('author.empty.folderEmpty.description')"
          :suggestion="$t('author.empty.folderEmpty.suggestion')"
          status-code=""
          :retry-text="emptyTexts.retryText"
          :show-retry-button="false"
          @retry="handleRetry"
        />

        <AuthorShareList :shares="filteredShares" :texts="shareTexts" @share-click="handleShareClick" />

        <AuthorEmptyState
          v-if="!loading && !imageLoading && !error && !currentFolderId && isAllEmpty"
          :title="emptyTexts.title"
          :description="emptyTexts.description"
          :suggestion="emptyTexts.suggestion"
          :status-code="emptyTexts.statusCode"
          :retry-text="emptyTexts.retryText"
          :show-retry-button="false"
          @retry="handleRetry"
        />

        <AuthorEmptyState
          v-if="error"
          :title="errorTexts.title"
          :description="error"
          status-code="ERROR_500"
          :retry-text="errorTexts.retry"
          :show-retry-button="true"
          @retry="handleRetry"
        />
      </section>
    </main>
  </div>
</template>

<style scoped lang="scss" src="./index.scss"></style>
