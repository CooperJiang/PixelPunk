<script setup lang="ts">
  import type { RandomImageAPI } from '@/api/openapi'
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'ApiCard',
  })

  const { $t } = useTexts()

  const props = defineProps<{
    api: RandomImageAPI
    url: string
  }>()

  const emit = defineEmits<{
    'toggle-status': [api: RandomImageAPI]
    edit: [api: RandomImageAPI]
    delete: [api: RandomImageAPI]
    'copy-url': [url: string]
  }>()

  const openApi = () => {
    window.open(props.url, '_blank')
  }

  const statusColor = computed(() => (props.api.is_active ? 'success' : 'disabled'))

  const formatCallCount = (count: number) => {
    if (count >= 1000000) {
      return `${(count / 1000000).toFixed(1)}M`
    }
    if (count >= 1000) {
      return `${(count / 1000).toFixed(1)}K`
    }
    return count.toString()
  }

  const formatTime = (time: string | null) => {
    if (!time) return $t('openApi.card.lastCall.never')
    const date = new Date(time)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (minutes < 1) return $t('openApi.card.lastCall.justNow')
    if (minutes < 60) return $t('openApi.card.lastCall.minutesAgo', { minutes })
    if (hours < 24) return $t('openApi.card.lastCall.hoursAgo', { hours })
    if (days < 30) return $t('openApi.card.lastCall.daysAgo', { days })
    return date.toLocaleDateString()
  }

  const formatCreatedTime = (time: string | null) => {
    if (!time) return '-'
    const date = new Date(time)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}`
  }

  const returnTypeLabel = computed(() => {
    return props.api.return_type === 'direct' ? $t('openApi.card.returnType.direct') : $t('openApi.card.returnType.redirect')
  })

  const returnTypeIcon = computed(() => {
    return props.api.return_type === 'direct' ? 'fas fa-file-image' : 'fas fa-share'
  })

  const scopeLabel = computed(() => {
    if (props.api.folder_id) {
      return props.api.folder_name
    }
    return $t('openApi.card.scopeAll')
  })
</script>

<template>
  <div class="api-card">
    <div class="api-card-header">
      <div class="api-info">
        <i class="fas fa-random api-icon"></i>
        <h3 class="api-name">{{ api.name }}</h3>
      </div>
      <div class="api-actions">
        <button class="action-btn open-btn" @click="openApi" :title="$t('openApi.card.actions.open')">
          <i class="fas fa-external-link-alt"></i>
        </button>
        <button class="action-btn edit-btn" @click="emit('edit', api)" :title="$t('openApi.card.actions.edit')">
          <i class="fas fa-edit"></i>
        </button>
        <button
          class="action-btn toggle-btn"
          :class="{ active: api.is_active }"
          @click="emit('toggle-status', api)"
          :title="api.is_active ? $t('openApi.form.status.disabled') : $t('openApi.form.status.enabled')"
        >
          <i :class="api.is_active ? 'fas fa-pause' : 'fas fa-play'"></i>
        </button>
        <button class="action-btn delete-btn" @click="emit('delete', api)" :title="$t('openApi.card.actions.delete')">
          <i class="fas fa-trash"></i>
        </button>
      </div>
    </div>

    <div class="api-card-body">
      <div class="api-details-row">
        <div class="api-detail api-detail-scope">
          <span class="detail-label">{{ $t('openApi.card.scope') }}:</span>
          <span class="detail-value detail-value-ellipsis">{{ scopeLabel }}</span>
        </div>
        <div class="api-detail api-detail-return">
          <span class="detail-label">{{ $t('openApi.card.returnMethod') }}:</span>
          <span class="detail-value">
            <i :class="returnTypeIcon" class="mr-1"></i>
            {{ returnTypeLabel }}
          </span>
        </div>
      </div>

      <div class="api-url">
        <code class="url-text">{{ url }}</code>
        <button class="copy-btn" @click="emit('copy-url', url)">
          <i class="fas fa-copy"></i>
        </button>
      </div>

      <div class="api-stats">
        <div class="stat-item">
          <i class="fas fa-chart-line stat-icon"></i>
          <span class="stat-value">{{ formatCallCount(api.call_count) }}{{ $t('openApi.card.times') }}</span>
        </div>
        <div class="stat-item">
          <i class="fas fa-clock stat-icon"></i>
          <span class="stat-value">{{ formatTime(api.last_called_at) }}</span>
        </div>
      </div>

      <div class="api-footer">
        <div class="created-time">
          <i class="fas fa-calendar-plus"></i>
          <span>{{ formatCreatedTime(api.created_at) }}</span>
        </div>
        <span class="status-badge" :class="`status-${statusColor}`">
          <i class="fas fa-circle"></i>
          {{ api.is_active ? $t('openApi.card.status.enabled') : $t('openApi.card.status.disabled') }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .api-card {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.6) 0%,
      rgba(var(--color-background-700-rgb), 0.8) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    overflow: hidden;
    transition: all 0.3s ease;
  }

  .api-card:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    transform: translateY(-2px);
    box-shadow:
      0 6px 20px rgba(0, 0, 0, 0.3),
      0 0 16px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .api-card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.875rem 1rem;
    background: linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 0.08) 0%, transparent 100%);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .api-info {
    display: flex;
    align-items: center;
    gap: 0.625rem;
    flex: 1;
    min-width: 0;
  }

  .api-icon {
    color: var(--color-brand-500);
    font-size: 1.125rem;
    flex-shrink: 0;
  }

  .api-name {
    font-size: 0.9375rem;
    font-weight: 600;
    color: var(--color-text-content);
    margin: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .api-actions {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    margin-left: 0.75rem;
    flex-shrink: 0;
  }

  .action-btn {
    width: 28px;
    height: 28px;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
    background: rgba(var(--color-background-900-rgb), 0.3);
    color: var(--color-content-muted);
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
  }

  .action-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-brand-500);
  }

  .action-btn.open-btn:hover {
    background: rgba(var(--color-info-rgb), 0.15);
    border-color: rgba(var(--color-info-rgb), 0.3);
    color: var(--color-info-500);
  }

  .action-btn.edit-btn:hover {
    color: var(--color-brand-500);
  }

  .action-btn.toggle-btn.active {
    background: rgba(var(--color-success-rgb), 0.15);
    border-color: rgba(var(--color-success-rgb), 0.3);
    color: var(--color-success-500);
  }

  .action-btn.toggle-btn.active:hover {
    background: rgba(var(--color-success-rgb), 0.25);
    border-color: rgba(var(--color-success-rgb), 0.4);
  }

  .action-btn.delete-btn:hover {
    background: rgba(var(--color-error-rgb), 0.15);
    border-color: rgba(var(--color-error-rgb), 0.3);
    color: var(--color-error-500);
  }

  .api-card-body {
    padding: 0.875rem 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .api-details-row {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    justify-content: space-between;
  }

  .api-detail {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .api-detail-scope {
    flex: 1;
    min-width: 0;
  }

  .api-detail-return {
    flex-shrink: 0;
  }

  .detail-label {
    color: var(--color-content-muted);
    font-size: 0.8125rem;
    white-space: nowrap;
  }

  .detail-value {
    color: var(--color-brand-500);
    font-size: 0.8125rem;
    font-weight: 600;
  }

  .detail-value-ellipsis {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .api-url {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: rgba(var(--color-background-900-rgb), 0.4);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 0.5rem 0.625rem;
  }

  .url-text {
    flex: 1;
    font-family: 'SF Mono', Monaco, monospace;
    font-size: 0.6875rem;
    color: var(--color-brand-500);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .copy-btn {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    padding: 0.1875rem 0.4375rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.6875rem;
    flex-shrink: 0;
  }

  .copy-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .api-stats {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    flex-wrap: wrap;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .stat-icon {
    color: var(--color-content-muted);
    font-size: 0.8125rem;
  }

  .stat-value {
    color: var(--color-text-content);
    font-size: 0.8125rem;
    font-weight: 600;
  }

  .api-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 0.5rem;
    margin-top: 0.25rem;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .created-time {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.6875rem;
    color: var(--color-content-muted);
  }

  .created-time i {
    font-size: 0.6875rem;
    color: var(--color-content-subtle);
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.3125rem;
    padding: 0.1875rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.6875rem;
    font-weight: 600;
  }

  .status-badge.status-success {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-500);
    border: 1px solid rgba(var(--color-success-rgb), 0.3);
  }

  .status-badge.status-disabled {
    background: rgba(var(--color-content-disabled-rgb), 0.1);
    color: var(--color-content-disabled);
    border: 1px solid rgba(var(--color-content-disabled-rgb), 0.2);
  }

  .status-badge i {
    font-size: 0.4375rem;
  }
</style>
