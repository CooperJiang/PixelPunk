<script setup lang="ts">
  import type { AuthorInfo, AuthorStats } from '@/api/types'

  defineOptions({
    name: 'AuthorHeader',
  })

  interface AuthorHeaderTexts {
    joinedLabel: string
    websiteLabel: string
    stats: {
      views: string
      shares: string
      images: string
    }
  }

  interface Props {
    authorInfo: AuthorInfo | null
    stats: AuthorStats
    texts: AuthorHeaderTexts
  }

  const props = defineProps<Props>()

  const formatJoinedDays = (days: number): string => props.texts.joinedLabel.replace('{days}', String(days))

  const formatNumber = (num: number): string => {
    if (num < 1000) {
      return num.toString()
    }
    if (num < 10_000) {
      return `${(num / 1000).toFixed(1)}K`
    }
    if (num < 1_000_000) {
      return `${Math.floor(num / 1000)}K`
    }
    if (num < 10_000_000) {
      return `${(num / 1_000_000).toFixed(1)}M`
    }
    return `${Math.floor(num / 1_000_000)}M`
  }
</script>

<template>
  <div v-if="authorInfo" class="author-header">
    <div class="author-profile">
      <div class="profile-main">
        <div class="avatar-section">
          <CyberUserAvatar
            :avatar-url="authorInfo.avatarFullPath"
            :username="authorInfo.username"
            size="md"
            :show-border="true"
          />
        </div>
        <div class="info-section">
          <h1 class="author-name">{{ authorInfo.username }}</h1>
          <p v-if="authorInfo.bio" class="author-bio">{{ authorInfo.bio }}</p>
          <div class="author-meta">
            <span class="meta-item">
              <i class="fas fa-calendar-alt" />
              {{ formatJoinedDays(authorInfo.daysJoined) }}
            </span>
            <a v-if="authorInfo.website" :href="authorInfo.website" target="_blank" rel="noopener noreferrer" class="meta-link">
              <i class="fas fa-external-link-alt" />
              {{ props.texts.websiteLabel }}
            </a>
          </div>
        </div>
      </div>

      <div class="stats-section">
        <div class="stat-item">
          <div class="stat-value">{{ formatNumber(stats.totalViews) }}</div>
          <div class="stat-label">{{ props.texts.stats.views }}</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ stats.totalShares }}</div>
          <div class="stat-label">{{ props.texts.stats.shares }}</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ stats.totalFiles || 0 }}</div>
          <div class="stat-label">{{ props.texts.stats.images }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss" src="./AuthorHeader.scss"></style>
