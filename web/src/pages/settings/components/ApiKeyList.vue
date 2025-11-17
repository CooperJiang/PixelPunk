<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { defineEmits, ref } from 'vue'
  import type { ApiKeyInfo } from '@/api/types/index'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  const _props = defineProps<{
    apiKeys: ApiKeyInfo[]
    isLoading: boolean
    pagination: {
      page: number
      size: number
      total: number
      total_pages: number
    }
  }>()

  const emit = defineEmits(['create', 'edit', 'detail', 'delete', 'toggle', 'change-page'])

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

  const { $t } = useTexts()

  const formatSize = (bytes: number) => {
    if (bytes === 0) {
      return $t('settings.apiKeys.text.unlimited')
    }

    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  const getStatusClass = (key: ApiKeyInfo) => (key.is_active ? 'status-active' : 'status-inactive')

  const getStorageText = (key: ApiKeyInfo) => {
    const used = key.storage_used === 0 ? '0 B' : formatSize(key.storage_used || 0)
    const limit = key.storage_limit === 0 ? $t('settings.apiKeys.text.unlimited') : formatSize(key.storage_limit)
    return `${used} / ${limit}`
  }

  const getStorageProgress = (key: ApiKeyInfo) => {
    if (key.storage_limit === 0) {
      return 0
    }
    return Math.min(((key.storage_used || 0) / key.storage_limit) * 100, 100)
  }

  const getStorageProgressText = (key: ApiKeyInfo) => {
    if (key.storage_limit === 0) {
      return $t('settings.apiKeys.text.unlimited')
    }
    const percentage = Math.round(((key.storage_used || 0) / key.storage_limit) * 100)
    return `${percentage}%`
  }

  const getUploadText = (key: ApiKeyInfo) => {
    const used = key.upload_count_used || 0
    const limit = key.upload_count_limit === 0 ? $t('settings.apiKeys.text.unlimited') : key.upload_count_limit.toString()
    return `${used} / ${limit}`
  }

  const getUploadProgress = (key: ApiKeyInfo) => {
    if (key.upload_count_limit === 0) {
      return 0
    }
    return Math.min(((key.upload_count_used || 0) / key.upload_count_limit) * 100, 100)
  }

  const getUploadProgressText = (key: ApiKeyInfo) => {
    if (key.upload_count_limit === 0) {
      return $t('settings.apiKeys.text.unlimited')
    }
    const percentage = Math.round(((key.upload_count_used || 0) / key.upload_count_limit) * 100)
    return `${percentage}%`
  }

  const getSingleFileLimitText = (key: ApiKeyInfo) =>
    key.single_file_limit === 0 ? $t('settings.apiKeys.text.unlimited') : formatSize(key.single_file_limit)

  const getFileTypesText = (key: ApiKeyInfo) => {
    if (!key.allowed_types || key.allowed_types.length === 0) {
      return $t('settings.apiKeys.text.unlimitedFiles')
    }
    return key.allowed_types.length > 2
      ? `${key.allowed_types.slice(0, 2).join(', ')}${$t('settings.apiKeys.text.fileTypeMore', {
          count: key.allowed_types.length,
        })}`
      : key.allowed_types.join(', ')
  }

  const onCreateClick = () => {
    emit('create')
  }

  const onEditClick = (key: ApiKeyInfo) => {
    emit('edit', key)
  }

  const onDetailClick = (key: ApiKeyInfo) => {
    emit('detail', key)
  }

  const onDeleteClick = (key: ApiKeyInfo) => {
    emit('delete', key)
  }

  const onToggleClick = (key: ApiKeyInfo) => {
    emit('toggle', key)
  }

  const changePage = (page: number) => {
    emit('change-page', page)
  }

  const toast = useToast()
  const visibleKeys = ref<Set<string>>(new Set())

  const toggleKeyVisibility = (keyId: string) => {
    if (visibleKeys.value.has(keyId)) {
      visibleKeys.value.delete(keyId)
    } else {
      visibleKeys.value.add(keyId)
    }
  }

  const isKeyVisible = (keyId: string) => visibleKeys.value.has(keyId)

  const maskApiKey = (key: string) => {
    if (!key) {
      return ''
    }
    const visibleLength = 8 // 显示前8个字符
    if (key.length <= visibleLength) {
      return '*'.repeat(key.length)
    }
    return key.substring(0, visibleLength) + '*'.repeat(Math.max(12, key.length - visibleLength))
  }

  const copyApiKey = async (key: string, keyName: string) => {
    if (!key) {
      toast.error($t('settings.apiKeys.toasts.emptyKey'))
      return
    }

    try {
      await navigator.clipboard.writeText(key)
      toast.success(
        $t('settings.apiKeys.toasts.copySuccess', {
          name: keyName,
        })
      )
    } catch (_error) {
      try {
        const textArea = document.createElement('textarea')
        textArea.value = key
        textArea.style.position = 'fixed'
        textArea.style.opacity = '0'
        document.body.appendChild(textArea)
        textArea.select()
        document.execCommand('copy')
        document.body.removeChild(textArea)
        toast.success(
          $t('settings.apiKeys.toasts.copySuccess', {
            name: keyName,
          })
        )
      } catch (_error) {
        toast.error($t('settings.apiKeys.toasts.copyFallbackError'))
      }
    }
  }
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h3 class="text-lg font-semibold">
          {{ $t('settings.apiKeys.header.title') }}
        </h3>
        <p class="text-sm text-brand-700">
          {{ $t('settings.apiKeys.header.subtitle') }}
        </p>
      </div>
      <CyberButton type="secondary" icon="plus" @click="onCreateClick">
        {{ $t('settings.apiKeys.header.create') }}
      </CyberButton>
    </div>

    <div v-if="isLoading" class="flex items-center justify-center py-12">
      <div class="loading-spinner">
        <i class="fas fa-spinner fa-spin text-3xl text-content" />
      </div>
    </div>

    <div
      v-else-if="!apiKeys.length"
      class="mb-4 radius-sm flex flex-col items-center justify-center border border-subtle bg-background-800 p-8 text-center"
    >
      <div class="empty-icon-container mb-4">
        <i class="fas fa-key" />
      </div>
      <h4 class="mb-2 text-lg font-semibold">
        {{ $t('settings.apiKeys.empty.title') }}
      </h4>
      <p class="mb-4 max-w-sm text-center text-sm text-brand-700">
        {{ $t('settings.apiKeys.empty.description') }}
      </p>
      <CyberButton type="secondary" icon="plus" @click="onCreateClick">
        {{ $t('settings.apiKeys.empty.action') }}
      </CyberButton>
    </div>

    <div v-else>
      <div v-for="key in apiKeys" :key="key.id" class="api-key-card mb-2  p-3 transition-all duration-200">
        <div class="mb-2 flex items-start justify-between">
          <div class="flex-1">
            <div class="mb-1 flex items-center gap-3">
              <h4 class="text-sm font-semibold text-content">{{ key.name }}</h4>
              <div class="flex items-center gap-2">
                <span class="status-badge text-xs" :class="getStatusClass(key)">
                  <i class="fas" :class="key.is_active ? 'fa-check-circle' : 'fa-times-circle'" />
                  {{
                    key.status_text ||
                    (key.is_active ? $t('settings.apiKeys.status.active') : $t('settings.apiKeys.status.inactive'))
                  }}
                </span>
                <span v-if="key.is_expired" class="status-badge status-expired text-xs">
                  <i class="fas fa-clock" />
                  {{ $t('settings.apiKeys.status.expired') }}
                </span>
              </div>
            </div>
            <div class="flex items-center gap-4 text-xs text-brand-600">
              <span>
                <i class="fas fa-calendar-alt mr-1" />
                {{ formatDate(key.created_at) }}
              </span>
              <span v-if="key.folder_path">
                <i class="fas fa-folder-tree mr-1 text-brand-400" />
                {{ key.folder_path }}
              </span>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <button
              class="status-toggle-btn"
              :class="key.is_active ? 'status-active' : 'status-inactive'"
              :title="key.is_active ? $t('settings.apiKeys.tooltips.disable') : $t('settings.apiKeys.tooltips.enable')"
              @click="onToggleClick(key)"
            >
              <i class="fas fa-power-off text-xs" />
              <span class="ml-1 text-xs font-medium">
                {{ key.is_active ? $t('settings.apiKeys.actions.enable') : $t('settings.apiKeys.actions.disable') }}
              </span>
            </button>

            <div class="flex items-center gap-1">
              <CyberButton
                type="ghost"
                :title="$t('settings.apiKeys.actions.view')"
                no-border
                class="action-btn"
                @click="onDetailClick(key)"
              >
                <i class="fas fa-eye" />
              </CyberButton>
              <CyberButton
                type="ghost"
                :title="$t('settings.apiKeys.actions.edit')"
                no-border
                class="action-btn"
                @click="onEditClick(key)"
              >
                <i class="fas fa-edit" />
              </CyberButton>
              <CyberButton
                type="ghost"
                :title="$t('settings.apiKeys.actions.delete')"
                no-border
                class="action-btn action-btn-danger"
                @click="onDeleteClick(key)"
              >
                <i class="fas fa-trash" />
              </CyberButton>
            </div>
          </div>
        </div>

        <div v-if="key.key" class="api-key-display-section mb-3 mt-3  p-3">
          <div class="mb-2 flex items-center justify-between">
            <div class="flex items-center gap-2">
              <i class="fas fa-key text-sm text-content" />
              <span class="text-xs font-medium text-brand-800">
                {{ $t('settings.apiKeys.fields.apiKey') }}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <button
                class="visibility-btn"
                :title="isKeyVisible(key.id) ? $t('settings.apiKeys.tooltips.hideKey') : $t('settings.apiKeys.tooltips.showKey')"
                @click="toggleKeyVisibility(key.id)"
              >
                <i class="fas text-xs" :class="isKeyVisible(key.id) ? 'fa-eye-slash' : 'fa-eye'" />
              </button>
              <button class="copy-btn" :title="$t('settings.apiKeys.tooltips.copy')" @click="copyApiKey(key.key, key.name)">
                <i class="fas fa-copy text-xs" />
              </button>
            </div>
          </div>
          <div class="api-key-value">
            <code class="api-key-text">
              {{ isKeyVisible(key.id) ? key.key : maskApiKey(key.key) }}
            </code>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-2 md:grid-cols-3">
          <div class="usage-stat">
            <div class="stat-header">
              <div class="stat-icon">
                <i class="fas fa-hdd text-content" />
              </div>
              <div>
                <div class="stat-label">
                  {{ $t('settings.apiKeys.fields.storage') }}
                </div>
                <div class="stat-value">{{ getStorageText(key) }}</div>
              </div>
            </div>
            <div class="progress-container">
              <div class="progress-bar">
                <div class="progress-fill storage-progress" :style="{ width: `${getStorageProgress(key)}%` }" />
              </div>
              <div class="progress-text">{{ getStorageProgressText(key) }}</div>
            </div>
          </div>

          <div class="usage-stat">
            <div class="stat-header">
              <div class="stat-icon">
                <i class="fas fa-upload text-error-400" />
              </div>
              <div>
                <div class="stat-label">
                  {{ $t('settings.apiKeys.fields.uploadCount') }}
                </div>
                <div class="stat-value">{{ getUploadText(key) }}</div>
              </div>
            </div>
            <div class="progress-container">
              <div class="progress-bar">
                <div class="progress-fill upload-progress" :style="{ width: `${getUploadProgress(key)}%` }" />
              </div>
              <div class="progress-text">{{ getUploadProgressText(key) }}</div>
            </div>
          </div>

          <div class="usage-stat">
            <div class="stat-header">
              <div class="stat-icon">
                <i class="fas fa-file text-success-400" />
              </div>
              <div>
                <div class="stat-label">
                  {{ $t('settings.apiKeys.fields.singleFileLimit') }}
                </div>
                <div class="stat-value">{{ getSingleFileLimitText(key) }}</div>
              </div>
            </div>
            <div class="file-types-preview">
              <div class="text-xs text-brand-600">
                {{ $t('settings.apiKeys.fields.fileTypes') }}:
                {{ getFileTypesText(key) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="pagination.total_pages > 1" class="mt-6 flex justify-center">
        <CyberPagination
          :total="pagination.total"
          :size="pagination.size"
          :current-page="pagination.page"
          @update:current-page="changePage"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
  .empty-icon-container {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 4rem;
    height: 4rem;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.25) 0%, rgba(var(--color-brand-600-rgb), 0.35) 100%);
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow:
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.25),
      0 0 24px rgba(var(--color-brand-500-rgb), 0.15),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
    animation: iconPulse 2s ease-in-out infinite;
  }

  .empty-icon-container i {
    font-size: 1.5rem;
    color: var(--color-brand-400);
    filter: drop-shadow(0 2px 4px rgba(var(--color-brand-500-rgb), 0.5));
  }

  @keyframes iconPulse {
    0%,
    100% {
      transform: scale(1);
      box-shadow:
        0 4px 12px rgba(var(--color-brand-500-rgb), 0.25),
        0 0 24px rgba(var(--color-brand-500-rgb), 0.15),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);
    }
    50% {
      transform: scale(1.05);
      box-shadow:
        0 6px 16px rgba(var(--color-brand-500-rgb), 0.35),
        0 0 32px rgba(var(--color-brand-500-rgb), 0.25),
        inset 0 1px 0 rgba(255, 255, 255, 0.15);
    }
  }

  .api-key-card {
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.16);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.85) 0%,
      rgba(var(--color-background-900-rgb), 0.72) 100%
    );
  }

  .api-key-card:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.28);
    transform: translateY(-1px);
    box-shadow:
      0 4px 16px rgba(0, 0, 0, 0.18),
      0 0 24px rgba(var(--color-brand-500-rgb), 0.08);
  }

  .usage-stat {
    background: rgba(var(--color-background-800-rgb), 0.55);
    border-radius: var(--radius-sm);
    padding: 0.5rem;
    border: 1px solid var(--color-border-subtle);
  }

  .stat-header {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    margin-bottom: 0.375rem;
  }

  .stat-icon {
    width: 24px;
    height: 24px;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(var(--color-brand-500-rgb), 0.12);
    font-size: 0.75rem;
  }

  .stat-label {
    font-size: 0.65rem;
    color: var(--color-content-muted);
    opacity: 0.85;
    margin-bottom: 0.125rem;
    line-height: 1;
  }

  .stat-value {
    font-size: 0.75rem;
    font-weight: 600;
    font-family: ui-monospace, SFMono-Regular, 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace;
    line-height: 1.2;
  }

  .progress-container {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .progress-bar {
    flex: 1;
    height: 4px;
    background-color: rgba(var(--color-background-600-rgb), 0.65);
    border-radius: var(--radius-sm);
    overflow: hidden;
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
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.3) 50%, transparent 100%);
    animation: shimmer 2s infinite;
  }

  .storage-progress {
    background: linear-gradient(90deg, rgba(var(--color-brand-500-rgb), 1) 0%, rgba(var(--color-brand-600-rgb), 1) 100%);
  }

  .upload-progress {
    background: linear-gradient(90deg, rgba(var(--color-error-rgb), 1) 0%, rgba(var(--color-error-rgb), 1) 100%);
  }

  .progress-text {
    font-size: 0.65rem;
    color: var(--color-content-muted);
    opacity: 0.9;
    min-width: 40px;
    text-align: right;
    font-family: ui-monospace, SFMono-Regular, 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace;
    line-height: 1;
  }

  .file-types-preview {
    margin-top: 0.375rem;
    padding-top: 0.375rem;
    border-top: 1px solid var(--color-border-subtle);
  }

  .status-badge {
    display: inline-flex;
    align-items: center;
    padding: 0.125rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.65rem;
    font-weight: 500;
    gap: 0.25rem;
    line-height: 1;
    border: 1px solid transparent;
  }

  .status-badge.status-active {
    background-color: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-400);
    border-color: rgba(var(--color-success-rgb), 0.3);
  }

  .status-badge.status-inactive {
    background-color: rgba(var(--color-error-rgb), 0.15);
    color: var(--color-error-400);
    border-color: rgba(var(--color-error-rgb), 0.28);
  }

  .status-badge.status-expired {
    background-color: rgba(var(--color-warning-rgb), 0.18);
    color: var(--color-warning-400);
    border-color: rgba(var(--color-warning-rgb), 0.32);
  }

  .status-toggle-btn {
    display: flex;
    align-items: center;
    padding: 0.25rem 0.625rem;
    border-radius: var(--radius-sm);
    font-weight: 600;
    transition: all 0.2s ease;
    border: 1px solid transparent;
    min-width: 70px;
    justify-content: center;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    height: 28px;
  }

  .status-toggle-btn.status-active {
    background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.25) 0%, rgba(var(--color-success-600-rgb), 0.25) 100%);
    color: var(--color-success-400);
    border-color: rgba(var(--color-success-rgb), 0.4);
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .status-toggle-btn.status-active:hover {
    background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.35) 0%, rgba(var(--color-success-600-rgb), 0.35) 100%);
    border-color: rgba(var(--color-success-rgb), 0.6);
    transform: translateY(-1px);
    box-shadow: 0 2px 6px rgba(var(--color-success-rgb), 0.22);
  }

  .status-toggle-btn.status-inactive {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-600-rgb), 0.35) 0%,
      rgba(var(--color-background-700-rgb), 0.35) 100%
    );
    color: var(--color-content-muted);
    border-color: rgba(var(--color-background-700-rgb), 0.25);
  }

  .status-toggle-btn.status-inactive:hover {
    background: linear-gradient(135deg, rgba(var(--color-success-rgb), 0.25) 0%, rgba(var(--color-success-600-rgb), 0.25) 100%);
    color: var(--color-success-400);
    border-color: rgba(var(--color-success-rgb), 0.38);
    transform: translateY(-1px);
    box-shadow: 0 2px 6px rgba(var(--color-success-rgb), 0.18);
  }

  .action-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    color: var(--color-content-muted);
    transition: all 0.2s ease;
  }

  .action-btn:hover {
    background-color: rgba(var(--color-brand-500-rgb), 0.12);
    color: var(--color-brand-400);
    transform: translateY(-1px);
  }

  .action-btn-danger:hover {
    background-color: rgba(var(--color-error-rgb), 0.12);
    color: var(--color-error-400);
  }

  .api-key-display-section {
    background: rgba(var(--color-background-900-rgb), 0.35);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
  }

  .api-key-display-section:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
    box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .visibility-btn,
  .copy-btn {
    padding: 0.375rem;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    background: rgba(var(--color-brand-500-rgb), 0.05);
    color: rgba(var(--color-brand-500-rgb), 0.8);
    transition: all 0.2s ease;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
  }

  .visibility-btn:hover,
  .copy-btn:hover {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-400);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-1px);
    box-shadow: 0 2px 6px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .copy-btn:hover {
    background: rgba(var(--color-success-rgb), 0.15);
    color: var(--color-success-400);
    border-color: rgba(var(--color-success-rgb), 0.4);
    box-shadow: 0 2px 6px rgba(var(--color-success-rgb), 0.18);
  }

  .api-key-value {
    background: rgba(var(--color-background-900-rgb), 0.4);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    padding: 0.75rem;
    font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  }

  .api-key-text {
    color: rgba(var(--color-brand-500-rgb), 0.9);
    font-size: 0.875rem;
    letter-spacing: 0.5px;
    word-break: break-all;
    user-select: all;
    background: transparent;
    border: none;
    outline: none;
    display: block;
    width: 100%;
  }
</style>
