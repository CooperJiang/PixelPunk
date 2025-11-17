<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { ref } from 'vue'
  import { useDevice } from '@/hooks/useDevice'

  const _props = defineProps({
    share: {
      type: Object,
      default: () => ({}),
    },
    user: {
      type: Object,
      default: () => ({}),
    },
  })

  const avatarError = ref(false)
  const { isMobile } = useDevice()
  const showDescription = ref(false)

  const handleAvatarError = () => {
    avatarError.value = true
  }

  const toggleDescription = () => {
    showDescription.value = !showDescription.value
  }

  const getUserInitial = (username: string) => {
    if (!username) {
      return '?'
    }
    return username.charAt(0).toUpperCase()
  }

  const formatDate = (dateString: string) => {
    if (!dateString) {
      return ''
    }

    const date = new Date(dateString)
    return date
      .toLocaleDateString(getCurrentLocale(), {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      })
      .replace(/\//g, '-')
  }

  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const formatDateShort = (dateString: string) => {
    if (!dateString) {
      return $t('share.header.permanent')
    }

    const date = new Date(dateString)
    return date
      .toLocaleDateString(getCurrentLocale(), {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
      })
      .replace(/\//g, '-')
  }
</script>

<template>
  <div class="share-header-container" :class="{ 'is-mobile': isMobile }">
    <div v-if="isMobile" class="mobile-fixed-header">
      <div class="mobile-header-top">
        <div class="mobile-title-area">
          <h1 class="mobile-title">
            <span v-if="share && share.name">{{ share.name }}</span>
            <span v-else>{{ $t('share.header.title') }}</span>
          </h1>
          <div v-if="share && share.has_password" class="mobile-badges">
            <i class="fas fa-lock" />
          </div>
        </div>

        <div v-if="user" class="mobile-user-info">
          <div v-if="user.avatar && !avatarError" class="mobile-avatar">
            <img :src="user.avatar" :alt="$t('share.header.avatar')" @error="handleAvatarError" />
          </div>
          <div v-else class="mobile-avatar-fallback h-full w-full">
            {{ getUserInitial(user.username) }}
          </div>
        </div>
      </div>

      <slot name="breadcrumb" />

      <div v-if="share" class="mobile-quick-info">
        <div v-if="share.current_views !== undefined" class="quick-info-item">
          <i class="fas fa-eye" /> {{ share.current_views }}
        </div>
        <div v-if="share && share.max_views > 0" class="quick-info-item">
          <i class="fas fa-stop-circle" /> {{ share.max_views }}
        </div>
        <div class="quick-info-item"><i class="fas fa-clock" /> {{ formatDateShort(share.expired_at) }}</div>
        <div v-if="user" class="quick-info-item"><i class="fas fa-user" /> {{ user.username }}</div>
      </div>
    </div>

    <div v-if="!isMobile" class="share-header-card">
      <div class="share-header-main">
        <div class="share-title-section">
          <h1 class="share-title">
            <span v-if="share && share.name">{{ share.name }}</span>
            <span v-else>{{ $t('share.header.title') }}</span>
          </h1>

          <div v-if="share" class="share-badges">
            <div v-if="share.has_password" class="badge">
              <i class="fas fa-lock" />
              <span>{{ $t('share.header.passwordProtected') }}</span>
            </div>
          </div>
        </div>

        <div class="share-info-row">
          <div v-if="user" class="share-user-section">
            <div class="user-info">
              <div v-if="user.avatar && !avatarError" class="user-avatar">
                <img :src="user.avatar" :alt="$t('share.header.userAvatar')" @error="handleAvatarError" />
              </div>
              <div v-else class="user-avatar">
                <div class="avatar-fallback">
                  {{ getUserInitial(user.username) }}
                </div>
              </div>
              <div class="user-details">
                <span class="username">{{ user.username }}</span>
                <span v-if="share && share.created_at" class="share-date">
                  <i class="fas fa-calendar-alt" /> {{ formatDate(share.created_at) }}
                </span>
              </div>
            </div>
          </div>

          <div v-if="share" class="share-metrics-section">
            <div class="metrics-grid">
              <div v-if="share.current_views !== undefined" class="metric-item">
                <div class="metric-icon">
                  <i class="fas fa-eye" />
                </div>
                <div class="metric-data">
                  <div class="metric-value">{{ share.current_views }}</div>
                  <div class="metric-label">{{ $t('share.header.viewCount') }}</div>
                </div>
              </div>

              <div v-if="share && share.max_views > 0" class="metric-item">
                <div class="metric-icon">
                  <i class="fas fa-stop-circle" />
                </div>
                <div class="metric-data">
                  <div class="metric-value">{{ share.max_views }}</div>
                  <div class="metric-label">{{ $t('share.header.maxViews') }}</div>
                </div>
              </div>

              <div v-if="share" class="metric-item">
                <div class="metric-icon">
                  <i class="fas fa-clock" />
                </div>
                <div class="metric-data">
                  <div class="metric-value">{{ formatDateShort(share.expired_at) }}</div>
                  <div class="metric-label">{{ $t('share.header.validUntil') }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="share && share.description" class="share-description-section">
          <div class="description-container">
            <i class="fas fa-quote-left description-icon" />
            <p class="share-description">{{ share.description }}</p>
          </div>
        </div>
      </div>

      <div class="decoration corner-tl" />
      <div class="decoration corner-tr" />
      <div class="decoration corner-bl" />
      <div class="decoration corner-br" />
      <div class="glow-effect" />
    </div>

    <div v-if="isMobile && share && share.description" class="mobile-description-section">
      <button class="mobile-description-toggle" @click="toggleDescription">
        <span>{{ $t('share.header.description') }}</span>
        <i :class="showDescription ? 'fas fa-chevron-up' : 'fas fa-chevron-down'" />
      </button>
      <div v-show="showDescription" class="mobile-description-content">
        <p>{{ share.description }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss" src="./ShareHeader.scss"></style>
