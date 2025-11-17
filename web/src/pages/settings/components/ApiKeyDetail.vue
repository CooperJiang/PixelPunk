<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { defineEmits, defineProps, ref } from 'vue'
  import type { ApiKeyInfo } from '@/api/types/index'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  const props = defineProps<{
    apiKey: ApiKeyInfo | null
    apiKeyValue?: string
  }>()

  const emit = defineEmits(['copy'])

  const toast = useToast()
  const { $t } = useTexts()
  const copied = ref(false)

  const formatDate = (dateString: string) => {
    if (!dateString) {
      return ''
    }

    const date = new Date(dateString)
    return date.toLocaleString(getCurrentLocale(), {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const formatSize = (bytes: number) => {
    if (bytes === 0) {
      return '0 B'
    }

    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const getStoragePercentageText = () => {
    if (!props.apiKey) {
      return '0%'
    }

    if (props.apiKey.storage_limit === 0) {
      return $t('settings.apiKeys.text.unlimited')
    }

    const percentage = Math.round(((props.apiKey.storage_used || 0) / props.apiKey.storage_limit) * 100)
    return `${percentage}%`
  }

  const getStorageProgressWidth = () => {
    if (!props.apiKey || props.apiKey.storage_limit === 0) {
      return 0
    }

    return Math.min(((props.apiKey.storage_used || 0) / props.apiKey.storage_limit) * 100, 100)
  }

  const getStorageLimitText = () => {
    if (!props.apiKey) {
      return '0 B'
    }

    if (props.apiKey.storage_limit === 0) {
      return $t('settings.apiKeys.text.unlimited')
    }

    return formatSize(props.apiKey.storage_limit)
  }

  const getUploadPercentageText = () => {
    if (!props.apiKey) {
      return '0%'
    }

    if (props.apiKey.upload_count_limit === 0) {
      return $t('settings.apiKeys.text.unlimited')
    }

    const percentage = Math.round(((props.apiKey.upload_count_used || 0) / props.apiKey.upload_count_limit) * 100)
    return `${percentage}%`
  }

  const getUploadProgressWidth = () => {
    if (!props.apiKey || props.apiKey.upload_count_limit === 0) {
      return 0
    }

    return Math.min(((props.apiKey.upload_count_used || 0) / props.apiKey.upload_count_limit) * 100, 100)
  }

  const getUploadLimitText = () => {
    if (!props.apiKey) {
      return '0'
    }

    if (props.apiKey.upload_count_limit === 0) {
      return $t('settings.apiKeys.text.unlimited')
    }

    return props.apiKey.upload_count_limit.toString()
  }

  const copyApiKey = async () => {
    if (!props.apiKeyValue) {
      return
    }

    const response = await navigator.clipboard.writeText(props.apiKeyValue)
    if (response !== undefined) {
      copied.value = true

      setTimeout(() => {
        copied.value = false
      }, 3000)

      toast.success($t('settings.apiKeys.toasts.detailCopySuccess'))
      emit('copy')
    }
  }
</script>

<template>
  <div v-if="apiKey" class="apikey-details">
    <div v-if="apiKeyValue" class="mb-6 radius-sm border border-subtle bg-background-700 p-4">
      <div class="mb-2 flex items-center justify-between">
        <h3 class="flex items-center text-content">
          <i class="fas fa-key mr-2" />
          <span>{{ $t('settings.apiKeys.fields.apiKey') }}</span>
        </h3>
        <button
          class="flex items-center rounded border border-subtle bg-background-600 p-1 px-2 text-xs text-content transition-colors hover:border-default hover:bg-background-500 hover:text-brand-400"
          :title="$t('settings.apiKeys.detail.copyCta')"
          @click="copyApiKey"
        >
          <i class="fas" :class="copied ? 'fa-check' : 'fa-copy'" />
          <span class="ml-1">
            {{ copied ? $t('settings.apiKeys.detail.copied') : $t('settings.apiKeys.detail.copy') }}
          </span>
        </button>
      </div>
      <div class="overflow-x-auto radius-sm border border-subtle bg-background-800 p-3 font-mono text-sm text-content">
        {{ apiKeyValue }}
      </div>
      <p class="mt-2 flex items-center text-xs text-warning-400">
        <i class="fas fa-exclamation-triangle mr-1" />
        {{ $t('settings.apiKeys.detail.copyNotice') }}
      </p>
    </div>

    <div class="apikey-details-content">
      <div class="mb-5 radius-sm border border-subtle bg-background-800 p-4">
        <h4 class="mb-3 flex items-center text-base font-medium text-content">
          <i class="fas fa-info-circle mr-2" />
          <span>{{ $t('settings.apiKeys.detail.baseInfoTitle') }}</span>
        </h4>

        <div class="grid grid-cols-2 gap-x-6 gap-y-4">
          <div class="info-item">
            <div class="info-label">{{ $t('settings.apiKeys.form.nameLabel') }}</div>
            <div class="info-value">{{ apiKey.name }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">{{ $t('settings.apiKeys.fields.status') }}</div>
            <div class="info-value">
              <span class="status-badge" :class="apiKey.is_active ? 'status-active' : 'status-inactive'">
                {{ apiKey.status === 1 ? $t('settings.apiKeys.status.active') : $t('settings.apiKeys.status.banned') }}
              </span>
              <span v-if="apiKey.is_expired" class="status-badge status-expired ml-2">
                {{ $t('settings.apiKeys.status.expired') }}
              </span>
            </div>
          </div>
          <div class="info-item">
            <div class="info-label">{{ $t('settings.apiKeys.fields.createdAt') }}</div>
            <div class="info-value">{{ formatDate(apiKey.created_at) }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">{{ $t('settings.apiKeys.fields.expiresAt') }}</div>
            <div class="info-value">
              {{ apiKey.expires_at ? `${apiKey.expires_at}` : $t('settings.apiKeys.text.neverExpires') }}
            </div>
          </div>
          <div class="info-item">
            <div class="info-label">{{ $t('settings.apiKeys.fields.lastUsedAt') }}</div>
            <div class="info-value">
              {{ apiKey.last_used_at ? formatDate(apiKey.last_used_at) : $t('settings.apiKeys.text.neverUsed') }}
            </div>
          </div>
          <div class="info-item">
            <div class="info-label">{{ $t('settings.apiKeys.fields.folder') }}</div>
            <div class="info-value">
              <span class="folder-path">
                <i class="fas fa-folder-tree mr-1 text-brand-400" />
                {{ apiKey.folder_path || $t('settings.apiKeys.text.rootFolder') }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="mb-5 radius-sm border border-subtle bg-background-800 p-4">
        <h4 class="mb-3 flex items-center text-base font-medium text-content">
          <i class="fas fa-chart-pie mr-2" />
          <span>{{ $t('settings.apiKeys.detail.usageTitle') }}</span>
        </h4>

        <div class="grid grid-cols-1 gap-5 md:grid-cols-3">
          <div class="usage-card">
            <div class="usage-header">
              <div class="usage-title">
                <i class="fas fa-hdd mr-2 text-content" />
                {{ $t('settings.apiKeys.fields.storage') }}
              </div>
              <div class="usage-percentage text-content">
                {{ getStoragePercentageText() }}
              </div>
            </div>

            <div class="progress-bar">
              <div class="progress-fill progress-storage" :style="{ width: `${getStorageProgressWidth()}%` }" />
            </div>

            <div class="usage-details">
              <div class="usage-value">{{ formatSize(apiKey.storage_used || 0) }}</div>
              <div class="usage-max">{{ getStorageLimitText() }}</div>
            </div>
          </div>

          <div class="usage-card">
            <div class="usage-header">
              <div class="usage-title">
                <i class="fas fa-upload mr-2 text-error-400" />
                {{ $t('settings.apiKeys.fields.uploadCount') }}
              </div>
              <div class="usage-percentage text-error-400">
                {{ getUploadPercentageText() }}
              </div>
            </div>

            <div class="progress-bar">
              <div class="progress-fill progress-upload" :style="{ width: `${getUploadProgressWidth()}%` }" />
            </div>

            <div class="usage-details">
              <div class="usage-value">{{ apiKey.upload_count_used || 0 }}</div>
              <div class="usage-max">{{ getUploadLimitText() }}</div>
            </div>
          </div>

          <div class="usage-card">
            <div class="usage-header">
              <div class="usage-title">
                <i class="fas fa-file mr-2 text-success-400" />
                {{ $t('settings.apiKeys.fields.singleFileLimit') }}
              </div>
            </div>

            <div class="file-size-limit">
              <div class="font-mono text-3xl text-success-400">
                {{ formatSize(apiKey.single_file_limit) }}
              </div>
              <div class="mt-1 text-xs text-brand-600">
                {{ $t('settings.apiKeys.detail.perFileHint') }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="mb-5 radius-sm border border-subtle bg-background-800 p-4">
        <h4 class="mb-3 flex items-center text-base font-medium text-content">
          <i class="fas fa-file-alt mr-2" />
          <span>{{ $t('settings.apiKeys.detail.allowedTypesTitle') }}</span>
        </h4>

        <div v-if="apiKey.allowed_types && apiKey.allowed_types.length > 0" class="flex flex-wrap gap-2">
          <div
            v-for="type in apiKey.allowed_types"
            :key="type"
            class="flex items-center rounded bg-brand-100 px-3 py-1 text-content"
          >
            <i class="fas fa-file mr-2" />
            {{ type }}
          </div>
        </div>
        <div v-else class="italic text-brand-600">
          {{ $t('settings.apiKeys.text.allFileTypes') }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .apikey-details-content {
    display: flex;
    flex-direction: column;
  }

  .info-item {
    margin-bottom: 0.5rem;
  }

  .info-label {
    font-size: 0.75rem;
    color: var(--color-content-muted);
    opacity: 0.85;
    margin-bottom: 0.25rem;
  }

  .info-value {
    font-size: 0.875rem;
  }

  .status-badge {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    font-size: 0.75rem;
    border-radius: var(--radius-sm);
    border: 1px solid transparent;
  }

  .status-active {
    background-color: rgba(var(--color-success-rgb), 0.16);
    color: var(--color-success-400);
    border-color: rgba(var(--color-success-rgb), 0.32);
  }

  .status-inactive {
    background-color: rgba(var(--color-error-rgb), 0.16);
    color: var(--color-error-400);
    border-color: rgba(var(--color-error-rgb), 0.28);
  }

  .status-expired {
    background-color: rgba(var(--color-warning-rgb), 0.18);
    color: var(--color-warning-400);
    border-color: rgba(var(--color-warning-rgb), 0.32);
  }

  .folder-path {
    display: inline-flex;
    align-items: center;
    font-family: monospace;
    background-color: rgba(var(--color-background-900-rgb), 0.4);
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    color: var(--color-content);
  }

  .usage-card {
    background-color: rgba(var(--color-background-800-rgb), 0.55);
    border-radius: var(--radius-sm);
    padding: 1rem;
    border: 1px solid var(--color-border-subtle);
  }

  .usage-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
  }

  .usage-title {
    display: flex;
    align-items: center;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .usage-percentage {
    font-size: 0.875rem;
    font-weight: 600;
  }

  .progress-bar {
    height: 0.5rem;
    background-color: rgba(var(--color-background-600-rgb), 0.65);
    border-radius: var(--radius-sm);
    overflow: hidden;
    margin-bottom: 0.5rem;
  }

  .progress-fill {
    height: 100%;
    border-radius: var(--radius-sm);
    transition: width 0.3s ease;
    position: relative;
  }

  .progress-fill::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.25) 50%, transparent 100%);
    pointer-events: none;
  }

  .usage-details {
    display: flex;
    justify-content: space-between;
    font-size: 0.75rem;
    color: var(--color-content-muted);
  }

  .progress-storage {
    background: linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 1) 0%, rgba(var(--color-brand-600-rgb), 1) 100%);
  }

  .progress-upload {
    background: linear-gradient(90deg, rgba(var(--color-error-rgb), 1) 0%, rgba(var(--color-error-rgb), 1) 100%);
  }

  .file-size-limit {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 70px;
  }
</style>
