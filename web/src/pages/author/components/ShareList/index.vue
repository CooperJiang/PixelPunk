<script setup lang="ts">
  import { ref } from 'vue'
  import type { AuthorShareInfo } from '@/api/types'
  import { formatRelativeTime } from '@/utils/formatting/format'

  defineOptions({
    name: 'AuthorShareList',
  })

  interface ShareListTexts {
    sectionTitle: string
    sectionSubtitle: string
    statusBadge: string
    meta: {
      viewsLabel: string
      createdLabel: string
      keyLabel: string
    }
    empty: {
      title: string
      description: string
    }
  }

  interface Props {
    shares: AuthorShareInfo[]
    texts: ShareListTexts
  }

  const props = defineProps<Props>()

  const failedCoverIds = ref(new Set<string>())

  const emit = defineEmits<{
    (e: 'share-click', share: AuthorShareInfo): void
  }>()

  const handleShareClick = (share: AuthorShareInfo) => {
    emit('share-click', share)
  }

  const shouldShowCover = (share: AuthorShareInfo): boolean => {
    const hasCover = Boolean(share.coverImageFullThumbURL || share.coverImageFullPath)
    return hasCover && !failedCoverIds.value.has(share.id)
  }

  const handleCoverError = (share: AuthorShareInfo, event: Event) => {
    const next = new Set(failedCoverIds.value)
    next.add(share.id)
    failedCoverIds.value = next
    const target = event.target as HTMLImageElement | null
    if (target) {
      target.remove()
    }
  }

  const formatViews = (views: number): string => {
    if (views < 1000) {
      return views.toString()
    }
    if (views < 10_000) {
      return `${(views / 1000).toFixed(1)}K`
    }
    if (views < 1_000_000) {
      return `${Math.floor(views / 1000)}K`
    }
    return `${(views / 1_000_000).toFixed(1)}M`
  }

  const getCreatedLabel = (date: string): string => {
    if (!date) {
      return ''
    }

    const formatted = formatRelativeTime(new Date(date).toISOString())
    return props.texts.meta.createdLabel.replace('{time}', formatted)
  }
</script>

<template>
  <div class="section-container">
    <div class="section-header">
      <h2 class="section-title">
        <div class="title-icon">
          <i class="fas fa-share-alt" />
          <div class="icon-glow" />
        </div>
        <div class="title-content">
          <span class="title-text">{{ props.texts.sectionTitle }}</span>
          <span class="title-subtitle">{{ props.texts.sectionSubtitle }}</span>
        </div>
        <div class="item-count">
          <span class="count-number">{{ shares.length }}</span>
          <span class="count-label">{{ props.texts.statusBadge }}</span>
        </div>
      </h2>
      <div class="header-decoration">
        <div class="decoration-line" />
        <div class="decoration-dots" />
      </div>
    </div>

    <div v-if="shares.length > 0" class="shares-list">
      <div
        v-for="(share, index) in shares"
        :key="share.id"
        class="share-item animate-in"
        :style="{ '--delay': `${index * 0.08}s` }"
        @click="handleShareClick(share)"
      >
        <div class="share-cover">
          <div class="cover-wrapper">
            <template v-if="shouldShowCover(share)">
              <img
                :src="share.coverImageFullThumbURL || share.coverImageFullPath"
                :alt="share.name"
                class="cover-img"
                @error="handleCoverError(share, $event)"
              />
            </template>
            <div v-else class="empty-cover">
              <i class="fas fa-images" />
            </div>
            <div class="cover-overlay">
              <i class="fas fa-external-link-alt" />
            </div>
          </div>
        </div>
        <div class="share-content">
          <h3 class="share-name">{{ share.name }}</h3>
          <p v-if="share.description" class="share-description">{{ share.description }}</p>
          <div class="share-meta">
            <div class="meta-item">
              <i class="fas fa-eye" />
              <span>{{ props.texts.meta.viewsLabel.replace('{count}', formatViews(share.views || 0)) }}</span>
            </div>
            <div class="meta-item">
              <i class="fas fa-calendar-alt" />
              <span>{{ getCreatedLabel(share.createdAt) }}</span>
            </div>
            <div class="meta-item">
              <i class="fas fa-key" />
              <span class="share-key">{{ props.texts.meta.keyLabel.replace('{key}', share.shareKey) }}</span>
            </div>
          </div>
        </div>
        <div class="share-actions">
          <div class="status-badge">{{ props.texts.statusBadge }}</div>
        </div>
        <div class="hover-glow" />
      </div>
    </div>

    <div v-else class="empty-state">
      <div class="empty-icon">
        <i class="fas fa-share-alt" />
      </div>
      <h3>{{ props.texts.empty.title }}</h3>
      <p>{{ props.texts.empty.description }}</p>
      <div class="empty-decoration">
        <div class="decoration-line" />
        <div class="decoration-dot" />
        <div class="decoration-line" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss" src="./ShareList.scss"></style>
