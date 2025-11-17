<script setup lang="ts">
  import { computed, onMounted, ref, watch } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { getAuthorFolder } from '@/api/author'
  import type { AuthorFolderContent, AuthorFolderInfo, AuthorImageInfo } from '@/api/types/index'
  import { useTexts } from '@/composables/useTexts'

  const route = useRoute()
  const router = useRouter()
  const { $t } = useTexts()

  const folderData = ref<AuthorFolderContent | null>(null)
  const loading = ref(false)
  const error = ref('')
  const currentPage = ref(1)
  const pageSize = ref(20)

  const authorId = route.params.authorId as string
  const folderId = route.params.folderId as string

  const visiblePages = computed(() => {
    if (!folderData.value) {
      return []
    }

    const totalPages = folderData.value.pagination.lastPage
    const current = currentPage.value
    const delta = 2 // 当前页前后显示的页数

    let start = Math.max(1, current - delta)
    let end = Math.min(totalPages, current + delta)

    if (current - start < delta) {
      end = Math.min(totalPages, end + (delta - (current - start)))
    }
    if (end - current < delta) {
      start = Math.max(1, start - (delta - (end - current)))
    }

    const pages = []
    for (let i = start; i <= end; i++) {
      pages.push(i)
    }
    return pages
  })

  const loadFolderData = async () => {
    if (!authorId || !folderId) {
      error.value = 'Invalid parameters'
      return
    }

    loading.value = true
    error.value = ''

    try {
      const data = await getAuthorFolder(authorId, folderId, {
        page: currentPage.value,
        size: pageSize.value,
      })

      folderData.value = data
    } catch (err) {
      error.value = (err as Error).message || 'Failed to load folder content'
    } finally {
      loading.value = false
    }
  }

  const goToPage = (page: number) => {
    if (page < 1 || (folderData.value && page > folderData.value.pagination.lastPage)) {
      return
    }
    currentPage.value = page
  }

  const openSubfolder = (subfolder: AuthorFolderInfo) => {
    router.push(`/author/${authorId}/folder/${subfolder.id}`)
  }

  const openImage = (image: AuthorImageInfo) => {
    window.open(image.fullPath, '_blank')
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) {
      return '0 B'
    }

    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
  }

  const goBack = () => {
    if (window.history.length > 1) {
      window.history.back()
    } else {
      closeWindow()
    }
  }

  const closeWindow = () => {
    try {
      window.close()
      setTimeout(() => {
        if (!window.closed) {
          if (window.history.length > 1) {
            window.history.back()
          } else {
            window.location.href = '/'
          }
        }
      }, 100)
    } catch {
      if (window.history.length > 1) {
        window.history.back()
      } else {
        window.location.href = '/'
      }
    }
  }

  const goToAuthorPage = () => {
    window.open(`/author/${authorId}`, '_blank')
  }

  watch(currentPage, () => {
    loadFolderData()
  })

  onMounted(() => {
    loadFolderData()
  })
</script>

<template>
  <div class="author-folder-page-standalone">
    <div class="standalone-header">
      <div class="header-left">
        <button class="back-btn" @click="goBack">
          <i class="fas fa-arrow-left" />
          <span>{{ $t('author.folder.back') }}</span>
        </button>
        <div class="page-title">
          <i class="fas fa-folder" />
          <span>{{ $t('author.folder.title') }}</span>
        </div>
      </div>
      <div class="header-right">
        <button class="close-btn" @click="closeWindow">
          <i class="fas fa-times" />
          <span>{{ $t('author.folder.close') }}</span>
        </button>
      </div>
    </div>

    <div class="author-folder-page">
      <div v-if="folderData" class="breadcrumb">
        <a class="breadcrumb-link" @click="goToAuthorPage">
          <i class="fas fa-user" />
          {{ $t('author.folder.authorHome') }}
        </a>
        <span class="breadcrumb-separator">/</span>
        <span class="current-folder">{{ folderData.folder.name }}</span>
      </div>

      <div v-if="folderData" class="folder-header">
        <div class="folder-info">
          <h1 class="folder-title">
            <i class="fas fa-folder" />
            {{ folderData.folder.name }}
          </h1>
          <div class="folder-stats">
            <span class="stat-item">
              <i class="fas fa-images" />
              {{ $t('author.folder.resourceCount', { count: folderData.folder.imageCount }) }}
            </span>
            <span v-if="folderData.subFolders.length > 0" class="stat-item">
              <i class="fas fa-folder" />
              {{ $t('author.folder.subFolderCount', { count: folderData.subFolders.length }) }}
            </span>
          </div>
        </div>
      </div>

      <div v-if="folderData && folderData.subFolders.length > 0" class="subfolders-section">
        <h2 class="section-title">
          <i class="fas fa-folder" />
          {{ $t('author.folder.subFolders') }}
        </h2>
        <div class="folders-grid">
          <div
            v-for="subfolder in folderData.subFolders"
            :key="subfolder.id"
            class="folder-card"
            @click="openSubfolder(subfolder)"
          >
            <div class="folder-cover">
              <img
                v-if="subfolder.coverImageFullPath"
                :src="subfolder.coverImageFullPath"
                :alt="subfolder.name"
                class="cover-img"
              />
              <div v-else class="empty-cover">
                <i class="fas fa-folder" />
              </div>
            </div>
            <div class="folder-info">
              <h3 class="folder-name">{{ subfolder.name }}</h3>
            </div>
          </div>
        </div>
      </div>

      <div v-if="folderData && folderData.images.length > 0" class="images-section">
        <h2 class="section-title">
          <i class="fas fa-images" />
          {{ $t('author.folder.files') }}
        </h2>
        <div class="images-grid">
          <div v-for="image in folderData.images" :key="image.id" class="image-card" @click="openImage(image)">
            <div class="image-wrapper">
              <img :src="image.fullPath" :alt="image.originalName" class="image-thumbnail" loading="lazy" />
              <div class="image-overlay">
                <div class="image-info">
                  <span class="image-name">{{ image.originalName }}</span>
                  <span class="image-size">{{ formatFileSize(image.size) }}</span>
                  <span class="image-dimensions">{{ image.width }} × {{ image.height }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="folderData && folderData.pagination.total > folderData.pagination.perPage" class="pagination-section">
        <div class="pagination">
          <button :disabled="currentPage === 1" class="page-btn" @click="goToPage(1)">
            {{ $t('author.folder.pagination.first') }}
          </button>
          <button :disabled="currentPage === 1" class="page-btn" @click="goToPage(currentPage - 1)">
            {{ $t('author.folder.pagination.prev') }}
          </button>

          <div class="page-numbers">
            <button
              v-for="page in visiblePages"
              :key="page"
              class="page-btn"
              :class="[{ active: page === currentPage }]"
              @click="goToPage(page)"
            >
              {{ page }}
            </button>
          </div>

          <button :disabled="currentPage === folderData.pagination.lastPage" class="page-btn" @click="goToPage(currentPage + 1)">
            {{ $t('author.folder.pagination.next') }}
          </button>
          <button
            :disabled="currentPage === folderData.pagination.lastPage"
            class="page-btn"
            @click="goToPage(folderData.pagination.lastPage)"
          >
            {{ $t('author.folder.pagination.last') }}
          </button>
        </div>

        <div class="pagination-info">
          {{
            $t('author.folder.pagination.info', {
              total: folderData.pagination.total,
              current: currentPage,
              last: folderData.pagination.lastPage,
            })
          }}
        </div>
      </div>

      <div v-if="folderData && folderData.images.length === 0 && folderData.subFolders.length === 0" class="empty-state">
        <div class="empty-content">
          <i class="fas fa-folder-open" />
          <h3>{{ $t('author.folder.empty.title') }}</h3>
          <p>{{ $t('author.folder.empty.description') }}</p>
        </div>
      </div>

      <div v-if="loading" class="loading-container">
        <div class="loading-spinner">
          <i class="fas fa-spinner fa-spin" />
          <span>{{ $t('author.folder.loading') }}</span>
        </div>
      </div>

      <div v-if="error && !loading" class="error-container">
        <div class="error-content">
          <i class="fas fa-exclamation-triangle" />
          <h3>{{ $t('author.folder.error.title') }}</h3>
          <p>{{ error }}</p>
          <button class="retry-btn" @click="loadFolderData">{{ $t('author.folder.error.retry') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .author-folder-page-standalone {
    min-height: 100vh;
    background: linear-gradient(135deg, rgba(15, 25, 35, 0.95) 0%, rgba(20, 30, 45, 0.9) 50%, rgba(25, 35, 55, 0.95) 100%);
    position: relative;
  }

  .standalone-header {
    position: sticky;
    top: 0;
    z-index: var(--z-index-dropdown);
    background: rgba(15, 25, 35, 0.95);
    backdrop-filter: blur(20px);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    padding: 1rem 2rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    box-shadow: 0 2px 20px rgba(0, 0, 0, 0.3);
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .header-right {
    display: flex;
    align-items: center;
  }

  .back-btn,
  .close-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: var(--color-brand-500);
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 0.9rem;
    font-weight: 500;
  }

  .back-btn:hover,
  .close-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    transform: translateY(-1px);
    box-shadow: 0 4px 15px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .close-btn {
    background: rgba(var(--color-error-rgb), 0.1);
    border-color: rgba(var(--color-error-rgb), 0.3);
    color: var(--color-error-500);
  }

  .close-btn:hover {
    background: rgba(var(--color-error-rgb), 0.2);
    border-color: rgba(var(--color-error-rgb), 0.5);
    box-shadow: 0 4px 15px rgba(var(--color-error-rgb), 0.3);
  }

  .page-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.1rem;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
  }

  .page-title i {
    color: var(--color-brand-500);
    font-size: 1.2rem;
  }

  .author-folder-page {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    background: transparent;
  }

  .breadcrumb {
    display: flex;
    align-items: center;
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
    color: rgba(255, 255, 255, 0.6);
    background: rgba(15, 25, 35, 0.6);
    padding: 1rem 1.5rem;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    backdrop-filter: blur(10px);
  }

  .breadcrumb-link {
    color: var(--color-brand-500);
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 0.25rem;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .breadcrumb-link:hover {
    color: var(--color-content-heading);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.8);
  }

  .breadcrumb-separator {
    margin: 0 0.5rem;
    color: rgba(255, 255, 255, 0.4);
  }

  .current-folder {
    font-weight: 500;
    color: rgba(255, 255, 255, 0.9);
  }

  .folder-header {
    background: rgba(15, 25, 35, 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 2rem;
    margin-bottom: 2rem;
    backdrop-filter: blur(10px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  }

  .folder-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-brand-500);
    margin: 0 0 1rem 0;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.5);
  }

  .folder-stats {
    display: flex;
    gap: 1.5rem;
    font-size: 0.875rem;
    color: rgba(255, 255, 255, 0.7);
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }

  .stat-item i {
    color: var(--color-brand-500);
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.25rem;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
    margin-bottom: 1rem;
  }

  .section-title i {
    color: var(--color-brand-500);
  }

  .subfolders-section {
    margin-bottom: 2rem;
  }

  .folders-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1.5rem;
  }

  .folder-card {
    background: white;
    border-radius: var(--radius-sm);
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    transition:
      transform 0.2s,
      box-shadow 0.2s;
  }

  .folder-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  }

  .folder-cover {
    aspect-ratio: 16/9;
    position: relative;
    overflow: hidden;
  }

  .cover-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .empty-cover {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-background-700);
    color: var(--color-content-muted);
  }

  .empty-cover i {
    font-size: 2rem;
  }

  .folder-info {
    padding: 1rem;
  }

  .folder-name {
    font-weight: 600;
    color: var(--color-content-heading);
    margin: 0;
    font-size: 1rem;
  }

  .images-section {
    margin-bottom: 2rem;
  }

  .images-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 1rem;
  }

  .image-card {
    cursor: pointer;
    border-radius: var(--radius-sm);
    overflow: hidden;
    transition: transform 0.2s;
  }

  .image-card:hover {
    transform: scale(1.02);
  }

  .image-wrapper {
    position: relative;
    aspect-ratio: 1;
    overflow: hidden;
  }

  .image-thumbnail {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: opacity 0.2s;
  }

  .image-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: linear-gradient(transparent, rgba(var(--color-background-900-rgb), 0.7));
    color: var(--color-content-heading);
    padding: 1rem;
    transform: translateY(100%);
    transition: transform 0.2s;
  }

  .image-card:hover .image-overlay {
    transform: translateY(0);
  }

  .image-info {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .image-name {
    font-weight: 500;
    font-size: 0.875rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .image-size,
  .image-dimensions {
    font-size: 0.75rem;
    opacity: 0.8;
  }

  .pagination-section {
    margin-top: 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .pagination {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .page-btn {
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--color-gray-300);
    background: white;
    color: var(--color-gray-700);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.875rem;
  }

  .page-btn:hover:not(:disabled) {
    background: var(--color-background-700);
    border-color: var(--color-content-muted);
  }

  .page-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .page-btn.active {
    background: var(--color-cyan-500);
    color: var(--color-text-on-brand);
    border-color: var(--color-cyan-500);
  }

  .page-numbers {
    display: flex;
    gap: 0.25rem;
  }

  .pagination-info {
    font-size: 0.875rem;
    color: var(--color-text-quaternary);
  }

  .empty-state {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 300px;
  }

  .empty-content {
    text-align: center;
    color: var(--color-text-quaternary);
  }

  .empty-content i {
    font-size: 4rem;
    margin-bottom: 1rem;
    color: var(--color-gray-300);
  }

  .empty-content h3 {
    margin: 0 0 0.5rem 0;
    color: var(--color-gray-700);
  }

  .loading-container,
  .error-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 300px;
  }

  .loading-spinner {
    text-align: center;
    color: var(--color-text-quaternary);
  }

  .loading-spinner i {
    font-size: 2rem;
    margin-bottom: 0.5rem;
    display: block;
  }

  .error-content {
    text-align: center;
    color: var(--color-text-quaternary);
  }

  .error-content i {
    font-size: 3rem;
    color: var(--color-danger);
    margin-bottom: 1rem;
  }

  .error-content h3 {
    margin: 0 0 0.5rem 0;
    color: var(--color-gray-700);
  }

  .retry-btn {
    margin-top: 1rem;
    padding: 0.5rem 1rem;
    background: var(--color-cyan-500);
    color: var(--color-text-on-brand);
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .retry-btn:hover {
    background: var(--color-brand-500);
  }

  @media (max-width: 768px) {
    .author-folder-page {
      padding: 1rem;
    }

    .folder-header {
      padding: 1.5rem;
    }

    .folder-stats {
      flex-direction: column;
      gap: 0.5rem;
    }

    .folders-grid {
      grid-template-columns: 1fr;
    }

    .images-grid {
      grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    }

    .pagination {
      flex-wrap: wrap;
      justify-content: center;
    }
  }
</style>
