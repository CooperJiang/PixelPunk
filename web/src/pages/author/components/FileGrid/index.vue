<script setup lang="ts">
  import { ref } from 'vue'
  import type { AuthorImageInfo } from '@/api/types'
  import { formatFileSize } from '@/utils/formatting/format'

  defineOptions({
    name: 'AuthorImageGrid',
  })

  type PaginationInfo = {
    currentPage: number
    perPage: number
    total: number
    lastPage: number
  }

  interface ActionTexts {
    label: string
    tooltip: string
  }

  interface ImageGridTexts {
    sectionTitle: string
    sectionSubtitle: string
    countLabel: string
    loadingText: string
    actions: {
      preview: ActionTexts
      download: ActionTexts
      detail: ActionTexts
    }
    paginationInfo: string
  }

  interface Props {
    images: AuthorImageInfo[]
    pagination: PaginationInfo | null
    imageLoading: boolean
    currentFolderId?: string
    texts: ImageGridTexts
    viewMode?: 'grid' | 'list'
  }

  const props = withDefaults(defineProps<Props>(), {
    viewMode: 'grid',
  })

  const emit = defineEmits<{
    (e: 'image-click', image: AuthorImageInfo): void
    (e: 'image-download', image: AuthorImageInfo): void
    (e: 'image-detail', image: AuthorImageInfo): void
    (e: 'page-change', page: number): void
  }>()

  const imageLoadStates = ref<Record<string, boolean>>({})

  const handleImageClick = (image: AuthorImageInfo) => {
    emit('image-click', image)
  }

  const handleDownload = (image: AuthorImageInfo) => {
    emit('image-download', image)
  }

  const handleImageDetail = (image: AuthorImageInfo) => {
    emit('image-detail', image)
  }

  const goToPage = (page: number) => {
    if (!props.pagination) {
      return
    }
    if (page < 1 || page > props.pagination.lastPage || page === props.pagination.currentPage) {
      return
    }
    emit('page-change', page)
  }

  const handleImageLoad = (event: Event) => {
    const img = event.target as HTMLImageElement
    const imageData = props.images.find((image) => img.alt === image.originalName)
    if (!imageData) {
      return
    }
    imageLoadStates.value = {
      ...imageLoadStates.value,
      [imageData.id]: true,
    }
    img.setAttribute('data-loaded', 'true')
  }

  const handleImageError = (event: Event) => {
    const img = event.target as HTMLImageElement
    const imageData = props.images.find((image) => img.alt === image.originalName)
    if (!imageData) {
      return
    }
    imageLoadStates.value = {
      ...imageLoadStates.value,
      [imageData.id]: true,
    }
    img.src = imageData.fullPath
    img.setAttribute('data-loaded', 'true')
  }

  const getFormattedSize = (bytes: number): string => {
    if (!bytes) {
      return '0 B'
    }
    return formatFileSize(bytes)
  }

  const getPaginationSummary = () => {
    if (!props.pagination) {
      return ''
    }
    return props.texts.paginationInfo
      .replace('{total}', String(props.pagination.total))
      .replace('{current}', String(props.pagination.currentPage))
      .replace('{last}', String(props.pagination.lastPage))
  }
</script>

<template>
  <div v-if="(images.length > 0 && !imageLoading) || (imageLoading && currentFolderId)" class="section-container">
    <div class="section-header">
      <h2 class="section-title">
        <div class="title-icon">
          <i class="fas fa-images" />
          <div class="icon-glow" />
        </div>
        <div class="title-content">
          <span class="title-text">{{ props.texts.sectionTitle }}</span>
          <span class="title-subtitle">{{ props.texts.sectionSubtitle }}</span>
        </div>
      </h2>
      <div class="header-decoration">
        <div class="decoration-line" />
        <div class="decoration-dots" />
      </div>
    </div>

    <div v-if="imageLoading" class="section-loading" v-loading="{ loading: true, text: props.texts.loadingText }" />

    <div v-else-if="images.length > 0" class="images-grid" :class="[`view-${props.viewMode}`]">
      <div
        v-for="(image, index) in images"
        :key="image.id"
        class="image-item animate-in"
        :style="{ '--delay': `${index * 0.1}s` }"
      >
        <div class="image-thumbnail">
          <div class="image-wrapper">
            <img
              :src="image.fullThumbURL || image.fullPath"
              :alt="image.originalName"
              class="thumbnail-img"
              @load="handleImageLoad"
              @error="handleImageError"
            />
            <div v-if="!imageLoadStates[image.id]" class="image-loading">
              <div class="loading-pulse" />
              <i class="fas fa-image" />
            </div>
            <div class="image-overlay">
              <div class="overlay-actions">
                <button
                  class="action-btn preview-btn"
                  :title="props.texts.actions.preview.tooltip"
                  @click.stop="handleImageClick(image)"
                >
                  <i class="fas fa-search-plus" />
                  <span>{{ props.texts.actions.preview.label }}</span>
                </button>
                <button
                  class="action-btn download-btn"
                  :title="props.texts.actions.download.tooltip"
                  @click.stop="handleDownload(image)"
                >
                  <i class="fas fa-download" />
                  <span>{{ props.texts.actions.download.label }}</span>
                </button>
                <button
                  v-if="image.ai_info"
                  class="action-btn detail-btn"
                  :title="props.texts.actions.detail.tooltip"
                  @click.stop="handleImageDetail(image)"
                >
                  <i class="fas fa-info-circle" />
                  <span>{{ props.texts.actions.detail.label }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>
        <div class="image-info">
          <div class="image-name" :title="image.originalName">{{ image.originalName }}</div>
          <div class="image-meta">
            <div class="meta-item">
              <i class="fas fa-hdd" />
              <span>{{ getFormattedSize(image.size) }}</span>
            </div>
            <div class="meta-item">
              <i class="fas fa-eye" />
              <span>{{ image.views || 0 }}</span>
            </div>
          </div>
        </div>
        <div class="image-glow" />
      </div>
    </div>

    <div v-if="pagination && pagination.total > pagination.perPage && !imageLoading" class="pagination-container">
      <div class="pagination-info">
        {{ getPaginationSummary() }}
      </div>
      <div class="pagination-controls">
        <button class="pagination-btn" :disabled="pagination.currentPage <= 1" @click="goToPage(pagination.currentPage - 1)">
          <i class="fas fa-chevron-left" />
        </button>
        <span class="page-indicator"> {{ pagination.currentPage }} / {{ pagination.lastPage }} </span>
        <button
          class="pagination-btn"
          :disabled="pagination.currentPage >= pagination.lastPage"
          @click="goToPage(pagination.currentPage + 1)"
        >
          <i class="fas fa-chevron-right" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss" src="./FileGrid.scss"></style>
