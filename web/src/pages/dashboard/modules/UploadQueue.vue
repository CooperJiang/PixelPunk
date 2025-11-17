<script setup lang="ts">
  import { useRouter } from 'vue-router'
  import { useUploadStore } from '@/store/upload'
  import { useTexts } from '@/composables/useTexts'

  defineOptions({
    name: 'UploadStatistics',
  })

  const router = useRouter()
  const uploadStore = useUploadStore()
  const { $t } = useTexts()

  const goToUploadPage = () => {
    router.push('/upload')
  }
</script>

<template>
  <div class="upload-statistics cyber-card">
    <div class="stats-content">
      <div class="stats-header">
        <h3 class="section-title">
          <i class="fas fa-chart-bar" />
          {{ $t('dashboard.uploadQueue.statsTitle') }}
        </h3>
        <div class="header-actions">
          <span class="status-text">{{ $t('dashboard.uploadQueue.totalFiles', { count: uploadStore.statistics.total }) }}</span>
          <button class="upload-entry-btn" :title="$t('dashboard.uploadQueue.actions.goToUpload')" @click="goToUploadPage">
            <i class="fas fa-external-link-alt" />
          </button>
        </div>
      </div>

      <div class="global-progress">
        <div class="progress-info">
          <span class="progress-label">{{ $t('dashboard.uploadQueue.overallProgress') }}</span>
          <span class="progress-value">{{ uploadStore.globalProgress }}%</span>
        </div>
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: uploadStore.globalProgress + '%' }" />
        </div>
      </div>

      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-icon total">
            <i class="fas fa-list" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ uploadStore.statistics.total }}</div>
            <div class="stat-label">{{ $t('dashboard.uploadQueue.queueTotal') }}</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon uploading">
            <i class="fas fa-spinner" :class="{ 'fa-spin': uploadStore.statistics.uploading > 0 }" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ uploadStore.statistics.uploading }}</div>
            <div class="stat-label">{{ $t('dashboard.uploadQueue.uploading') }}</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon pending">
            <i class="fas fa-clock" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ uploadStore.statistics.pending }}</div>
            <div class="stat-label">{{ $t('dashboard.uploadQueue.pending') }}</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon completed">
            <i class="fas fa-check" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ uploadStore.statistics.completed }}</div>
            <div class="stat-label">{{ $t('dashboard.uploadQueue.completed') }}</div>
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-icon failed">
            <i class="fas fa-times" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ uploadStore.statistics.failed }}</div>
            <div class="stat-label">{{ $t('dashboard.uploadQueue.failed') }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .upload-statistics {
    padding: 20px;
    height: auto;
    min-height: 350px;
  }

  .stats-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
    height: 100%;
  }

  .stats-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .upload-entry-btn {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-background-700);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    color: var(--color-brand-500);
    font-size: 11px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .upload-entry-btn:hover {
    background: var(--color-hover-bg);
    border-color: var(--color-hover-border);
    color: var(--color-brand-400);
    transform: translateY(-1px);
  }

  .section-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 16px;
  }

  .status-text {
    font-size: 11px;
    color: var(--color-content-muted);
  }

  .global-progress {
    margin-bottom: 4px;
  }

  .progress-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 6px;
  }

  .progress-label {
    font-size: 12px;
    color: var(--color-content-default);
  }

  .progress-value {
    font-size: 11px;
    color: var(--color-brand-500);
    font-weight: 600;
    font-family: 'SF Mono', monospace;
  }

  .progress-bar {
    height: 4px;
    background: var(--color-background-900);
    border-radius: var(--radius-sm);
    overflow: hidden;
  }

  .progress-bar.mini {
    height: 2px;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--color-brand-500), var(--color-brand-400));
    border-radius: var(--radius-sm);
    transition: width 0.3s ease;
  }

  .stats-grid {
    display: flex;
    flex-direction: column;
    gap: 8px;
    flex: 1;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: var(--color-background-800);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-subtle);
    transition: all 0.2s ease;
    min-height: 50px;
  }

  .stat-item:hover {
    background: var(--color-hover-bg);
    border-color: var(--color-hover-border);
    transform: translateY(-1px);
  }

  .stat-icon {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    border-radius: var(--radius-sm);
    flex-shrink: 0;
  }

  .stat-icon.total {
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .stat-icon.uploading {
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .stat-icon.pending {
    color: var(--color-warning);
    background: rgba(var(--color-warning-rgb), 0.15);
    border: 1px solid rgba(var(--color-warning-rgb), 0.2);
  }

  .stat-icon.completed {
    color: var(--color-success);
    background: rgba(var(--color-success-rgb), 0.15);
    border: 1px solid rgba(var(--color-success-rgb), 0.2);
  }

  .stat-icon.failed {
    color: var(--color-error);
    background: rgba(var(--color-error-rgb), 0.15);
    border: 1px solid rgba(var(--color-error-rgb), 0.2);
  }

  .stat-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
    flex: 1;
  }

  .stat-value {
    font-size: 16px;
    font-weight: 700;
    color: var(--color-content-heading);
    font-family: 'SF Mono', monospace;
    line-height: 1;
  }

  .stat-label {
    font-size: 11px;
    color: var(--color-content-muted);
    font-weight: 500;
  }

  .empty-state {
    text-align: center;
    padding: 40px 20px;
    color: rgba(var(--color-content-rgb), 0.6);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 200px;
  }

  .empty-icon {
    font-size: 48px;
    margin-bottom: 16px;
    color: rgba(var(--color-badge-accent-text-rgb), 0.3);
    background: rgba(var(--color-badge-accent-text-rgb), 0.05);
    width: 80px;
    height: 80px;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(var(--color-badge-accent-text-rgb), 0.1);
  }

  .empty-text {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 8px;
    color: var(--color-content-heading);
  }

  .empty-subtitle {
    font-size: 12px;
    color: rgba(var(--color-content-rgb), 0.5);
  }

  @media (max-width: 768px) {
    .upload-statistics {
      padding: 16px;
      height: auto;
      min-height: 300px;
    }

    .stat-item {
      padding: 10px;
      min-height: 45px;
    }

    .stat-icon {
      width: 28px;
      height: 28px;
      font-size: 10px;
    }

    .stat-value {
      font-size: 14px;
    }

    .stat-label {
      font-size: 10px;
    }
  }
</style>
