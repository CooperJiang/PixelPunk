<script setup lang="ts">
  import { ref } from 'vue'
  import type { ImageInfo } from '@/api/types'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    imageData: ImageInfo
  }

  interface Emits {
    (e: 'go-to-author'): void
    (e: 'view-work', work: any): void
  }

  defineProps<Props>()
  const emit = defineEmits<Emits>()
  const { $t } = useTexts()

  const workLoadingStates = ref<Record<string, boolean>>({})

  const formatNumber = (num: number | undefined): string => {
    if (num === undefined || num === null || isNaN(num)) {
      return '0'
    }
    if (num >= 10000) {
      return `${(num / 10000).toFixed(1)}w`
    } else if (num >= 1000) {
      return `${(num / 1000).toFixed(1)}k`
    }
    return num.toString()
  }

  const getInitials = (username: string): string => {
    if (!username) {
      return '?'
    }
    return username.charAt(0).toUpperCase()
  }

  const goToAuthorPage = () => {
    emit('go-to-author')
  }

  const viewWork = (work: any) => {
    emit('view-work', work)
  }
</script>

<template>
  <div class="creator-card">
    <h3 class="section-title"><i class="fas fa-user mr-2" />{{ $t('random.creator.title') }}</h3>

    <div class="creator-info">
      <div class="creator-avatar">
        <img
          v-if="imageData.user_info?.avatar"
          :src="imageData.user_info.avatar"
          :alt="imageData.user_info.username"
          class="avatar-img"
        />
        <div v-else class="avatar-placeholder">
          {{ getInitials(imageData.user_name) }}
        </div>
      </div>
      <div class="creator-details">
        <div class="creator-name-row">
          <h4 class="creator-name">{{ imageData.user_name }}</h4>
          <button class="creator-link-btn" @click="goToAuthorPage">
            <i class="fas fa-external-link-alt" />
          </button>
        </div>
        <div class="creator-stats">
          <span class="stat-item">
            <i class="fas fa-images" />
            {{ $t('random.creator.worksCount', { count: imageData.user_info?.total_images || 0 }) }}
          </span>
          <span v-if="imageData.user_info?.total_views" class="stat-item">
            <i class="fas fa-eye" />
            {{ $t('random.creator.viewsCount', { count: formatNumber(imageData.user_info.total_views) }) }}
          </span>
          <span v-if="imageData.user_info?.days_joined" class="stat-item">
            <i class="fas fa-calendar" />
            {{ $t('random.creator.daysJoined', { days: imageData.user_info.days_joined }) }}
          </span>
        </div>
      </div>
    </div>

    <div v-if="imageData.user_info?.top_tags && imageData.user_info.top_tags.length" class="creator-preferences">
      <h5 class="preferences-title">{{ $t('random.creator.preferences') }}</h5>
      <div class="preference-tags">
        <CyberTag v-for="tag in imageData.user_info.top_tags" :key="tag" variant="primary" size="small">
          {{ tag }}
        </CyberTag>
      </div>
    </div>

    <div v-if="imageData.user_info?.other_images && imageData.user_info.other_images.length" class="other-works">
      <h5 class="works-title">{{ $t('random.creator.otherWorks') }}</h5>
      <div class="works-grid">
        <div
          v-for="work in imageData.user_info.other_images"
          :key="work.id"
          class="work-item"
          v-loading="workLoadingStates[work.id]"
          @click="viewWork(work)"
        >
          <CyberFile
            :src="work.full_thumb_url"
            :alt="work.display_name"
            class="work-image"
            fit-mode="cover"
            :retry-count="2"
            :is-nsfw="work.is_nsfw"
            background-pattern="none"
            @loading="workLoadingStates[work.id] = $event"
          />
          <div class="work-overlay">
            <span class="work-size">{{ work.width }}×{{ work.height }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .creator-card {
    background: rgba(var(--color-background-700-rgb), 0.35);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    backdrop-filter: blur(12px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    transition: all 0.3s ease;
  }

  .creator-card:hover {
    background: rgba(var(--color-background-700-rgb), 0.45);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }

  .section-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin-bottom: 1rem;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .creator-info {
    display: flex;
    gap: 0.75rem;
    align-items: flex-start;
    margin-bottom: 1rem;
  }

  .creator-avatar {
    flex-shrink: 0;
  }

  .avatar-img,
  .avatar-placeholder {
    width: 45px;
    height: 45px;
    border-radius: var(--radius-full);
    border: 2px solid var(--color-border-strong);
  }

  .avatar-img {
    object-fit: cover;
  }

  .avatar-placeholder {
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-error-500));
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.1rem;
    font-weight: bold;
    color: var(--color-text-on-brand);
  }

  .creator-details {
    flex: 1;
  }

  .creator-name-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .creator-name {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-content-heading);
  }

  .creator-link-btn {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border: 1px solid var(--color-border-strong);
    padding: 0.25rem;
    border-radius: var(--radius-full);
    font-size: 0.7rem;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
  }

  .creator-link-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.25);
    transform: scale(1.1);
  }

  .creator-stats {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .stat-item {
    color: rgba(var(--color-content-rgb), 0.7);
    font-size: 0.75rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .stat-item i {
    color: var(--color-brand-500);
    width: 12px;
    font-size: 0.7rem;
  }

  .creator-preferences {
    margin-bottom: 1rem;
  }

  .preferences-title {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin-bottom: 0.5rem;
  }

  .preference-tags {
    display: flex;
    gap: 0.4rem;
    overflow: hidden;
    max-width: 100%;
  }

  .other-works {
    margin-top: 1rem;
  }

  .works-title {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin-bottom: 0.5rem;
  }

  .works-grid {
    display: flex;
    gap: 0.75rem;
    max-width: 100%;
    padding: 0.25rem;
    overflow: hidden;
    justify-content: center;
  }

  .work-item {
    position: relative;
    cursor: pointer;
    border-radius: var(--radius-sm);
    overflow: hidden;
    background: rgba(var(--color-background-900-rgb), 0.2);
    border: 1px solid var(--color-border-subtle);
    width: 80px;
    height: 80px;
    flex-shrink: 0;
    transition: all 0.3s ease;
    box-shadow: var(--shadow-cyber-sm);
  }

  .work-item:hover {
    transform: translateY(-2px);
    border-color: var(--color-border-strong);
    box-shadow: var(--shadow-cyber-md);
  }

  .work-image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center;
    transition: transform 0.3s ease;
    display: block;
  }

  .work-item:hover .work-image {
    transform: scale(1.1);
  }

  .work-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: linear-gradient(transparent, rgba(var(--color-background-900-rgb), 0.8));
    padding: 0.25rem;
    opacity: 0;
    transition: opacity 0.3s ease;
    z-index: 10;
  }

  .work-item:hover .work-overlay {
    opacity: 1;
  }

  .work-size {
    color: var(--color-content-heading);
    font-size: 0.6rem;
    text-shadow: 0 1px 3px rgba(var(--color-background-900-rgb), 0.5);
  }

  @media (max-width: 768px) {
    .creator-info {
      flex-direction: column;
      text-align: center;
      align-items: center;
    }

    .creator-name-row {
      justify-content: center;
    }

    .works-grid {
      gap: 0.6rem;
      padding: 0.2rem;
    }

    .work-item {
      width: 70px;
      height: 70px;
    }
  }
</style>
