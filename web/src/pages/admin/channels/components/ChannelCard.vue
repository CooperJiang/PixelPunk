<script setup lang="ts">
  import type { StorageChannel } from '@/api/types/index'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineProps<{
    channel: StorageChannel
  }>()

  defineEmits<{
    (e: 'edit', channel: StorageChannel): void
    (e: 'export', channel: StorageChannel): void
    (e: 'toggle-status', channel: StorageChannel): void
    (e: 'set-default', channel: StorageChannel): void
    (e: 'delete', channel: StorageChannel): void
    (e: 'test', channel: StorageChannel): void
  }>()

  const toast = useToast()

  const getChannelIcon = (type: string) => {
    switch (type) {
      case 'oss':
        return 'fas fa-cloud'
      case 'cos':
        return 'fas fa-cloud'
      case 'rainyun':
        return 'fas fa-cloud-rain'
      case 'local':
        return 'fas fa-hdd'
      default:
        return 'fas fa-server'
    }
  }

  const getChannelIconColor = (type: string) => {
    switch (type) {
      case 'oss':
        return 'var(--color-brand-500)'
      case 'cos':
        return 'var(--color-success-400)'
      case 'rainyun':
        return 'var(--color-info-500)'
      case 'local':
        return 'var(--color-warning-500)'
      default:
        return 'var(--color-content-muted)'
    }
  }

  const getChannelTypeText = (type: string) => {
    const typeMap: Record<string, string> = {
      oss: $t('constants.storageOptions.channelTypes.oss'),
      cos: $t('constants.storageOptions.channelTypes.cos'),
      rainyun: $t('constants.storageOptions.channelTypes.rainyun'),
      local: $t('constants.storageOptions.channelTypes.local'),
    }
    return typeMap[type] || $t('constants.storageOptions.channelTypes.unknown')
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text).then(
      () => {
        toast.success($t('admin.channels.messages.copySuccess'))
      },
      () => {
        toast.error($t('admin.channels.messages.copyFailed'))
      }
    )
  }
</script>

<template>
  <div class="channel-card" :class="{ 'default-channel': channel.is_default, 'local-channel': channel.type === 'local' }">
    <div v-if="channel.is_default" class="default-glow" />

    <div class="card-section header-section">
      <div class="channel-header">
        <div class="channel-icon-wrapper">
          <div class="channel-icon">
            <i :class="getChannelIcon(channel.type)" :style="{ color: getChannelIconColor(channel.type) }" />
          </div>
          <div v-if="channel.is_default" class="default-crown">
            <i class="fas fa-crown" />
          </div>
        </div>

        <div class="channel-info">
          <h3 class="channel-name">
            {{ channel.name }}
          </h3>
          <div class="channel-meta">
            <span v-if="channel.is_default" class="default-badge">
              <i class="fas fa-star" />{{ $t('admin.channels.status.default') }}
            </span>
            <div class="channel-status" :class="{ active: channel.status === 1 }">
              <span class="status-dot" />
              {{ channel.status === 1 ? $t('admin.channels.status.enabled') : $t('admin.channels.status.disabled') }}
            </div>
          </div>
        </div>

        <div class="channel-actions">
          <button class="action-btn edit" :title="$t('admin.channels.actions.edit')" @click="$emit('edit', channel)">
            <i class="fas fa-edit" />
          </button>
          <button
            class="action-btn export"
            :disabled="channel.type === 'local'"
            :title="
              channel.type === 'local'
                ? $t('admin.channels.actions.exportNotSupported')
                : $t('admin.channels.actions.exportConfig')
            "
            @click="$emit('export', channel)"
          >
            <i class="fas fa-file-export" />
          </button>
          <CyberPopconfirm
            :title="channel.status === 1 ? $t('admin.channels.confirm.disable') : $t('admin.channels.confirm.enable')"
            @confirm="$emit('toggle-status', channel)"
          >
            <template #reference>
              <button
                class="action-btn"
                :class="channel.status === 1 ? 'deactivate' : 'activate'"
                :title="channel.status === 1 ? $t('admin.channels.actions.disable') : $t('admin.channels.actions.enable')"
              >
                <i :class="channel.status === 1 ? 'fas fa-ban' : 'fas fa-check'" />
              </button>
            </template>
          </CyberPopconfirm>
        </div>
      </div>
    </div>

    <!-- 详情区块 -->
    <div class="card-section details-section">
      <div class="channel-details">
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-icon">
              <i class="fas fa-hdd" />
            </div>
            <div class="stat-content">
              <span class="stat-label">{{ $t('admin.channels.info.type') }}</span>
              <span class="stat-value">{{ getChannelTypeText(channel.type) }}</span>
            </div>
          </div>

          <div class="stat-item">
            <div class="stat-icon">
              <i class="fas fa-images" />
            </div>
            <div class="stat-content">
              <span class="stat-label">{{ $t('admin.channels.info.file') }}</span>
              <span class="stat-value">{{ channel.file_count || 0 }}</span>
            </div>
          </div>

          <div class="stat-item">
            <div class="stat-icon">
              <i class="fas fa-database" />
            </div>
            <div class="stat-content">
              <span class="stat-label">{{ $t('admin.channels.info.bucket') }}</span>
              <span class="stat-value truncate">{{ channel.bucket || '-' }}</span>
            </div>
          </div>

          <div
            class="stat-item"
            :class="{ clickable: channel.custom_domain || channel.domain || channel.endpoint }"
            @click="
              channel.custom_domain || channel.domain || channel.endpoint
                ? copyToClipboard(channel.custom_domain || channel.domain || channel.endpoint)
                : null
            "
          >
            <div class="stat-icon">
              <i class="fas fa-globe" />
            </div>
            <div class="stat-content">
              <span class="stat-label">{{ $t('admin.channels.info.domain') }}</span>
              <span class="stat-value truncate">{{ channel.custom_domain || channel.domain || channel.endpoint || '-' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部区块 -->
    <div class="card-section footer-section">
      <div class="channel-footer">
        <div class="footer-left">
          <button v-if="!channel.is_default" class="footer-btn primary" @click="$emit('set-default', channel)">
            <i class="fas fa-star mr-1" />{{ $t('admin.channels.actions.setDefault') }}
          </button>

          <span v-if="channel.is_default" class="default-tip">
            <i class="fas fa-crown mr-1" />{{ $t('admin.channels.status.defaultChannel') }}
          </span>
        </div>

        <div class="footer-right">
          <button
            class="footer-btn test"
            :disabled="channel.status !== 1"
            :title="channel.status !== 1 ? $t('admin.channels.actions.testDisabledTip') : $t('admin.channels.actions.test')"
            @click="$emit('test', channel)"
          >
            <i class="fas fa-plug mr-1" />{{ $t('admin.channels.actions.test') }}
          </button>

          <CyberPopconfirm
            v-if="!channel.is_local"
            :title="$t('admin.channels.confirm.delete')"
            @confirm="$emit('delete', channel)"
          >
            <template #reference>
              <button class="footer-btn danger" :disabled="channel.is_default">
                <i class="fas fa-trash mr-1" />{{ $t('admin.channels.actions.delete') }}
              </button>
            </template>
          </CyberPopconfirm>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .channel-card {
    position: relative;
    overflow: hidden;
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal) var(--ease-out);
    background: rgba(var(--color-background-900-rgb), 0.3);
    border: 1px solid rgba(var(--color-border-default-rgb), 0.3);
    box-shadow:
      0 1px 3px rgba(0, 0, 0, 0.2),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.05);
    backdrop-filter: blur(8px);
  }

  .channel-card:hover {
    transform: translateY(-2px);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    box-shadow:
      0 4px 12px rgba(0, 0, 0, 0.3),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.15),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .channel-card.default-channel {
    border: 1px solid rgba(var(--color-warning-rgb), 0.35);
    box-shadow:
      0 1px 3px rgba(0, 0, 0, 0.2),
      0 0 0 1px rgba(var(--color-warning-rgb), 0.15),
      0 0 16px rgba(var(--color-warning-rgb), 0.15);
  }

  .channel-card.default-channel:hover {
    box-shadow:
      0 4px 12px rgba(0, 0, 0, 0.3),
      0 0 0 1px rgba(var(--color-warning-rgb), 0.25),
      0 0 24px rgba(var(--color-warning-rgb), 0.2);
  }

  .default-glow {
    pointer-events: none;
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at 50% 0%, rgba(var(--color-warning-rgb), 0.06) 0%, transparent 50%);
  }

  .channel-card.local-channel {
    border-color: rgba(var(--color-warning-rgb), 0.25);
  }

  .card-section {
    position: relative;
    z-index: 10;
  }

  .header-section {
    padding: 12px 14px;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-700-rgb), 0.35) 0%,
      rgba(var(--color-background-800-rgb), 0.25) 100%
    );
    border-bottom: 1px solid rgba(var(--color-border-default-rgb), 0.18);
  }

  .channel-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
  }

  .channel-icon-wrapper {
    position: relative;
    flex-shrink: 0;
  }

  .channel-icon {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    font-size: 16px;
    border-radius: var(--radius-md);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-600-rgb), 0.5) 0%,
      rgba(var(--color-background-700-rgb), 0.4) 100%
    );
    border: 1px solid rgba(var(--color-border-default-rgb), 0.25);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
    transition: all var(--transition-fast) var(--ease-out);
  }

  .channel-card:hover .channel-icon {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-600-rgb), 0.7) 0%,
      rgba(var(--color-background-700-rgb), 0.6) 100%
    );
    border-color: rgba(var(--color-border-default-rgb), 0.35);
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.2);
  }

  .default-crown {
    position: absolute;
    right: -4px;
    top: -4px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: var(--radius-full);
    font-size: 12px;
    background: linear-gradient(135deg, #ffd700 0%, #ffa500 100%);
    color: #ffffff;
    box-shadow:
      0 0 12px rgba(255, 215, 0, 0.6),
      0 2px 4px rgba(0, 0, 0, 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.3);
    border: 1px solid rgba(255, 215, 0, 0.4);
    animation: crown-glow 2s ease-in-out infinite;
  }

  @keyframes crown-glow {
    0%,
    100% {
      box-shadow:
        0 0 12px rgba(255, 215, 0, 0.6),
        0 2px 4px rgba(0, 0, 0, 0.3),
        inset 0 1px 0 rgba(255, 255, 255, 0.3);
    }
    50% {
      box-shadow:
        0 0 20px rgba(255, 215, 0, 0.8),
        0 2px 6px rgba(0, 0, 0, 0.4),
        inset 0 1px 0 rgba(255, 255, 255, 0.4);
    }
  }

  .channel-info {
    flex: 1;
    margin-left: 12px;
    margin-right: 8px;
    min-width: 0;
  }

  .channel-name {
    font-size: 15px;
    font-weight: 600;
    color: var(--color-content-default);
    line-height: 1.3;
    margin-bottom: 6px;
    letter-spacing: 0.01em;
  }

  .channel-meta {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .default-badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 2px 8px;
    font-size: 10px;
    font-weight: 600;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.18), rgba(var(--color-warning-rgb), 0.12));
    border: 1px solid rgba(var(--color-warning-rgb), 0.35);
    color: var(--color-warning-400);
    white-space: nowrap;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .channel-status {
    display: inline-flex;
    align-items: center;
    width: fit-content;
    padding: 2px 8px;
    font-size: 10px;
    font-weight: 600;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-error-rgb), 0.12);
    border: 1px solid rgba(var(--color-error-rgb), 0.35);
    color: var(--color-error-400);
  }

  .channel-status.active {
    background: rgba(var(--color-success-rgb), 0.12);
    border-color: rgba(var(--color-success-rgb), 0.4);
    color: var(--color-success-400);
  }

  .status-dot {
    width: 5px;
    height: 5px;
    margin-right: 5px;
    border-radius: var(--radius-full);
    background: currentColor;
    box-shadow: 0 0 6px currentColor;
    animation: pulse-dot 2s ease-in-out infinite;
  }

  @keyframes pulse-dot {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.6;
    }
  }

  .channel-actions {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
  }

  .action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 26px;
    height: 26px;
    font-size: 11px;
    border-radius: var(--radius-sm);
    border: 1px solid;
    background: rgba(var(--color-background-600-rgb), 0.2);
    backdrop-filter: var(--backdrop-blur-sm);
    transition: all var(--transition-fast) var(--ease-out);
    cursor: pointer;
  }

  .action-btn:active:not(:disabled) {
    transform: scale(0.9);
  }

  .action-btn.edit {
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
    color: rgba(var(--color-brand-rgb), 0.7);
  }

  .action-btn.edit:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    color: var(--color-brand-400);
    box-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .action-btn.export {
    border-color: rgba(var(--color-info-rgb), 0.25);
    color: rgba(var(--color-info-rgb), 0.7);
  }

  .action-btn.export:hover:not(:disabled) {
    background: rgba(var(--color-info-rgb), 0.2);
    border-color: rgba(var(--color-info-rgb), 0.5);
    color: var(--color-info-400);
    box-shadow: 0 0 8px rgba(var(--color-info-rgb), 0.3);
  }

  .action-btn.export:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .action-btn.activate {
    border-color: rgba(var(--color-success-rgb), 0.25);
    color: rgba(var(--color-success-rgb), 0.7);
  }

  .action-btn.activate:hover {
    background: rgba(var(--color-success-rgb), 0.2);
    border-color: rgba(var(--color-success-rgb), 0.5);
    color: var(--color-success-400);
    box-shadow: 0 0 8px rgba(var(--color-success-rgb), 0.3);
  }

  .action-btn.deactivate {
    border-color: rgba(var(--color-error-rgb), 0.25);
    color: rgba(var(--color-error-rgb), 0.7);
  }

  .action-btn.deactivate:hover {
    background: rgba(var(--color-error-rgb), 0.2);
    border-color: rgba(var(--color-error-rgb), 0.5);
    color: var(--color-error-400);
    box-shadow: 0 0 8px rgba(var(--color-error-rgb), 0.3);
  }

  .details-section {
    padding: 12px 14px;
    background: rgba(var(--color-background-800-rgb), 0.2);
    border-bottom: 1px solid rgba(var(--color-border-default-rgb), 0.15);
  }

  .channel-details {
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 7px;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 9px;
    padding: 9px;
    border-radius: var(--radius-md);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-600-rgb), 0.3) 0%,
      rgba(var(--color-background-700-rgb), 0.25) 100%
    );
    border: 1px solid rgba(var(--color-border-default-rgb), 0.18);
    transition: all var(--transition-fast) var(--ease-out);
  }

  .stat-item.clickable {
    cursor: pointer;
  }

  .stat-item.clickable:hover {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12) 0%, rgba(var(--color-brand-500-rgb), 0.08) 100%);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateY(-1px);
    box-shadow: 0 2px 6px rgba(var(--color-brand-500-rgb), 0.15);
  }

  .stat-icon {
    display: flex;
    flex-shrink: 0;
    align-items: center;
    justify-content: center;
    width: 30px;
    height: 30px;
    font-size: 12px;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.18) 0%, rgba(var(--color-brand-500-rgb), 0.12) 100%);
    color: var(--color-brand-400);
    box-shadow: 0 1px 3px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .stat-content {
    display: flex;
    flex: 1;
    min-width: 0;
    flex-direction: column;
    gap: 2px;
  }

  .stat-label {
    font-size: 11px;
    color: var(--color-content-muted);
    line-height: 1.2;
  }

  .stat-value {
    font-size: 13px;
    font-weight: 500;
    color: var(--color-content-default);
    line-height: 1.3;
  }

  .truncate {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .footer-section {
    padding: 12px 14px;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-900-rgb), 0.28) 0%,
      rgba(var(--color-background-800-rgb), 0.15) 100%
    );
  }

  .channel-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-sm);
  }

  .footer-left,
  .footer-right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .footer-btn {
    display: inline-flex;
    align-items: center;
    padding: 6px 12px;
    font-size: 12px;
    font-weight: 500;
    border-radius: var(--radius-md);
    border: 1px solid;
    transition: all var(--transition-fast) var(--ease-out);
    cursor: pointer;
    white-space: nowrap;
  }

  .footer-btn:active:not(:disabled) {
    transform: scale(0.98);
  }

  .footer-btn.primary {
    background: rgba(var(--color-brand-500-rgb), 0.12);
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-400);
  }

  .footer-btn.primary:hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-brand-300);
  }

  .footer-btn.test {
    background: rgba(var(--color-success-rgb), 0.12);
    border-color: rgba(var(--color-success-rgb), 0.3);
    color: var(--color-success-400);
  }

  .footer-btn.test:hover:not(:disabled) {
    background: rgba(var(--color-success-rgb), 0.2);
    border-color: rgba(var(--color-success-rgb), 0.4);
    color: var(--color-success-300);
  }

  .footer-btn.test:disabled {
    opacity: 0.4;
    cursor: not-allowed;
    color: rgba(var(--color-success-rgb), 0.5);
  }

  .footer-btn.danger {
    background: rgba(var(--color-error-rgb), 0.12);
    border-color: rgba(var(--color-error-rgb), 0.3);
    color: var(--color-error-400);
  }

  .footer-btn.danger:hover:not(:disabled) {
    background: rgba(var(--color-error-rgb), 0.2);
    border-color: rgba(var(--color-error-rgb), 0.4);
    color: var(--color-error-300);
  }

  .footer-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .default-tip {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 10px;
    font-size: 11px;
    font-weight: 600;
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.15), rgba(var(--color-warning-rgb), 0.1));
    border: 1px solid rgba(var(--color-warning-rgb), 0.3);
    color: var(--color-warning-400);
  }

  .mr-1 {
    margin-right: var(--space-xs);
  }
</style>
