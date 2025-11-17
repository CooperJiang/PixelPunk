<template>
  <div class="vector-action-panel-simplified">
    <div class="action-grid">
      <div class="action-card primary">
        <div class="card-header">
          <i class="fas fa-magic"></i>
          <span>{{ $t('admin.vectors.buttons.reset') }}</span>
        </div>
        <div class="card-content">
          <p class="description">{{ $t('admin.vectors.qdrant.systemInfo') }}</p>
          <div class="problem-stats" v-if="problemStats">
            <span v-if="problemStats.failed > 0" class="problem-item failed">
              <i class="fas fa-times-circle"></i>{{ problemStats.failed }} {{ $t('admin.vectors.status.error') }}
            </span>
            <span v-if="problemStats.missing > 0" class="problem-item missing">
              <i class="fas fa-exclamation-triangle"></i>{{ problemStats.missing }} {{ $t('admin.vectors.status.indexing') }}
            </span>
            <span v-if="problemStats.stuck > 0" class="problem-item stuck">
              <i class="fas fa-hourglass-half"></i>{{ problemStats.stuck }} {{ $t('admin.vectors.status.processing') }}
            </span>
            <span v-if="problemStats.unknown > 0" class="problem-item stuck">
              <i class="fas fa-hourglass-half"></i>{{ problemStats.unknown }} {{ $t('admin.vectors.filter.options.all') }}
            </span>
            <span v-if="problemStats.total === 0" class="no-problems">
              <i class="fas fa-check-circle"></i>{{ $t('admin.vectors.toast.noData') }}
            </span>
          </div>
        </div>
        <div class="card-actions">
          <CyberButton
            type="primary"
            :loading="processing"
            :disabled="!problemStats || problemStats.total === 0"
            @click="handleSmartFix"
            class="action-button"
          >
            <i class="fas fa-magic mr-1"></i>
            {{
              problemStats?.total > 0
                ? `${$t('admin.vectors.buttons.reset')} ${problemStats.total} ${$t('admin.vectors.settings.indexingSettings.indexTimeout')}`
                : $t('admin.vectors.empty.action')
            }}
          </CyberButton>
        </div>
      </div>

      <div class="action-card secondary">
        <div class="card-header">
          <i class="fas fa-search-plus"></i>
          <span>{{ $t('admin.vectors.settings.monitoringSettings.logLevel') }}</span>
        </div>
        <div class="card-content">
          <p class="description">{{ $t('admin.vectors.qdrant.description') }}</p>
          <div class="verification-info" v-if="verificationStats">
            <div class="sync-status" :class="getSyncStatusClass()">
              <i :class="getSyncIcon()"></i>
              {{ $t('admin.vectors.actions.syncRatio', { ratio: verificationStats.sync_ratio?.toFixed(1) || 0 }) }}
            </div>
          </div>
        </div>
        <div class="card-actions">
          <CyberButton type="secondary" :loading="verificationProcessing" @click="handleStartVerification" class="action-button">
            <i class="fas fa-search-plus mr-1"></i>{{ $t('admin.vectors.actions.startVerification') }}
          </CyberButton>
        </div>
      </div>

      <div class="action-card danger">
        <div class="card-header">
          <i class="fas fa-sync-alt"></i>
          <span>{{ $t('admin.vectors.actions.fullRebuild') }}</span>
        </div>
        <div class="card-content">
          <p class="description">{{ $t('admin.vectors.actions.fullRebuildDesc') }}</p>
          <div class="warning-text">
            <i class="fas fa-exclamation-triangle"></i>
            {{ $t('admin.vectors.actions.fullRebuildWarning') }}
          </div>
        </div>
        <div class="card-actions">
          <CyberButton type="danger" :loading="processing" @click="handleRegenerateAll" class="action-button">
            <i class="fas fa-sync-alt mr-1"></i>{{ $t('admin.vectors.actions.fullRebuild') }}
          </CyberButton>
        </div>
      </div>

      <div class="action-card secondary">
        <div class="card-header">
          <i class="fas fa-balance-scale"></i>
          <span>{{ $t('admin.vectors.actions.dataAlignment') }}</span>
        </div>
        <div class="card-content">
          <p class="description">{{ $t('admin.vectors.actions.dataAlignmentDesc') }}</p>
        </div>
        <div class="card-actions">
          <CyberButton type="secondary" :loading="processing" @click="handleReconcileMissing" class="action-button">
            <i class="fas fa-plus-circle mr-1"></i>{{ $t('admin.vectors.actions.reconcileMissing') }}
          </CyberButton>
          <CyberButton type="danger" :loading="processing" @click="handleCleanOrphans" class="action-button">
            <i class="fas fa-broom mr-1"></i>{{ $t('admin.vectors.actions.cleanOrphans') }}
          </CyberButton>
          <CyberButton type="primary" :loading="processing" @click="handleRebuildStale" class="action-button">
            <i class="fas fa-recycle mr-1"></i>{{ $t('admin.vectors.actions.rebuildStale') }}
          </CyberButton>
        </div>
      </div>
    </div>

    <div
      v-if="currentVerificationTask && currentVerificationTask.status === 'running'"
      class="mt-6 overflow-hidden rounded-xl border border-indigo-600/20 bg-gradient-to-r from-indigo-800/50 to-purple-700/30"
    >
      <div class="border-b border-indigo-600/20 bg-gradient-to-r from-indigo-600/20 to-purple-600/10 px-6 py-4">
        <h4 class="flex items-center text-lg font-semibold text-indigo-400">
          <div class="mr-3 flex h-6 w-6 items-center justify-center">
            <i class="fas fa-spinner fa-spin text-indigo-500"></i>
          </div>
          {{ $t('admin.vectors.actions.verificationInProgress') }}
        </h4>
      </div>

      <div class="p-6">
        <div class="mb-4">
          <div class="mb-2 flex items-center justify-between">
            <span class="text-sm font-medium text-slate-300">{{ $t('admin.vectors.actions.processingProgress') }}</span>
            <span class="text-sm font-bold text-indigo-400">
              {{
                (
                  ((currentVerificationTask.processed_count || 0) / Math.max(currentVerificationTask.total_count || 1, 1)) *
                  100
                ).toFixed(1)
              }}%
            </span>
          </div>

          <div class="h-3 w-full overflow-hidden rounded-full border border-slate-600/30 bg-slate-700/50">
            <div
              class="relative h-full overflow-hidden rounded-full bg-gradient-to-r from-indigo-500 to-purple-500 transition-all duration-300 ease-out"
              :style="{
                width: `${((currentVerificationTask.processed_count || 0) / Math.max(currentVerificationTask.total_count || 1, 1)) * 100}%`,
              }"
            >
              <div
                class="absolute inset-0 skew-x-12 animate-pulse bg-gradient-to-r from-transparent via-white/30 to-transparent"
              ></div>
            </div>
          </div>

          <div class="mt-2 flex items-center justify-between">
            <span class="text-xs text-slate-400">
              {{ $t('admin.vectors.actions.processed', { count: formatNumber(currentVerificationTask.processed_count || 0) }) }}
            </span>
            <span class="text-xs text-slate-400">
              {{ $t('admin.vectors.actions.total', { count: formatNumber(currentVerificationTask.total_count || 0) }) }}
            </span>
          </div>
        </div>

        <div class="rounded-lg border border-slate-600/30 bg-slate-800/40 p-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <div class="mr-3 flex h-8 w-8 items-center justify-center rounded-full bg-indigo-500/20">
                <i class="fas fa-tachometer-alt text-sm text-indigo-400"></i>
              </div>
              <span class="text-sm font-medium text-slate-300">{{ $t('admin.vectors.actions.realtimeStats') }}</span>
            </div>
            <div class="text-right">
              <div class="text-sm font-semibold text-indigo-400">
                {{ currentVerificationTask.processed_count || 0 }}/{{ currentVerificationTask.total_count || 0 }}
              </div>
              <div class="text-xs text-slate-400">
                {{
                  $t('admin.vectors.actions.remaining', {
                    count: (currentVerificationTask.total_count || 0) - (currentVerificationTask.processed_count || 0),
                  })
                }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <CyberDialog v-model="showConfirmDialog" :title="confirmData.title" width="400">
      <div class="confirm-content">
        <div class="confirm-text">
          <i class="confirm-icon fas fa-exclamation-triangle text-warning"></i>
          <div class="message-lines">
            <p v-for="line in confirmData.message" :key="line" class="message-line">{{ line }}</p>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex w-full justify-end gap-2 p-2">
          <CyberButton type="outlined" @click="showConfirmDialog = false">{{ $t('admin.vectors.actions.cancel') }}</CyberButton>
          <CyberButton :type="confirmData.type" :loading="processing" @click="executeConfirmedAction">
            {{ confirmData.confirmText }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, onUnmounted } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { useWebSocketStore } from '@/store/websocket'
  import {
    startVectorVerification,
    getVerificationProgress,
    getVerificationStatistics,
    batchRepairMissingVectors,
    type VectorVerificationTask,
    type VectorVerificationStats,
  } from '@/api/admin/vector-verification'
  import { recoverStuckTasks, reconcileMissing, cleanOrphans, rebuildStale } from '@/api/admin/vectors'

  const toast = useToast()
  const { $t } = useTexts()
  const wsStore = useWebSocketStore()
  const isPaused = computed(() => wsStore.vectorStats?.runtime?.paused === true)

  interface Props {
    processing?: boolean
    stats?: Record<string, unknown>
  }

  interface Emits {
    (e: 'regenerate-all'): void
    (e: 'refresh'): void
    (e: 'verification-completed'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    processing: false,
    stats: () => null,
  })

  const emit = defineEmits<Emits>()

  // const isRefreshing = ref(false)
  const verificationProcessing = ref(false)
  const currentVerificationTask = ref<VectorVerificationTask | null>(null)
  const verificationStats = ref<VectorVerificationStats | null>(null)
  const verificationInterval = ref<number | null>(null)

  const showConfirmDialog = ref(false)
  const confirmData = ref<{
    title: string
    message: string[]
    type: 'primary' | 'danger' | 'warning'
    confirmText: string
    action: string
  }>({
    title: '',
    message: [],
    type: 'primary',
    confirmText: $t('admin.vectors.actions.confirm'),
    action: '',
  })

  const problemStats = computed(() => {
    if (!props.stats && !verificationStats.value) return null

    const failed = props.stats?.failed_count || 0
    const stuck = props.stats?.running_count || 0 // 可能卡住的运行中任务
    const missing = verificationStats.value?.missing_count || 0
    const unknown = verificationStats.value?.unknown_count || 0

    const total = failed + stuck + missing + unknown

    return { failed, stuck, missing, unknown, total }
  })

  const handleSmartFix = () => {
    if (!problemStats.value || problemStats.value.total === 0) {
      toast.info($t('admin.vectors.toast.noData'))
      return
    }

    const problems = []
    if (problemStats.value.failed > 0) problems.push(`${problemStats.value.failed} ${$t('admin.vectors.status.error')}`)
    if (problemStats.value.stuck > 0) problems.push(`${problemStats.value.stuck} ${$t('admin.vectors.status.processing')}`)
    if (problemStats.value.missing > 0) problems.push(`${problemStats.value.missing} ${$t('admin.vectors.status.indexing')}`)
    if (problemStats.value.unknown > 0) problems.push(`${problemStats.value.unknown} ${$t('admin.vectors.filter.options.all')}`)

    confirmData.value = {
      title: $t('admin.vectors.confirm.reindex.title'),
      message: [
        `${$t('admin.vectors.actions.rebuild')}：`,
        '',
        ...problems,
        '',
        `${$t('admin.vectors.confirm.reindex.message')}?`,
      ],
      type: 'primary',
      confirmText: $t('admin.vectors.buttons.save'),
      action: 'smart-fix',
    }
    showConfirmDialog.value = true
  }

  const executeSmartFix = async () => {
    try {
      verificationProcessing.value = true
      let fixCount = 0

      if (problemStats.value?.failed > 0) {
        emit('retry-all-failed')
        fixCount += problemStats.value.failed
      }

      if (problemStats.value?.stuck > 0) {
        const result = await recoverStuckTasks()
        if (result.success) {
          fixCount += result.data?.affected_count || 0
        }
      }

      if (problemStats.value?.missing > 0 || problemStats.value?.unknown > 0) {
        const result = await batchRepairMissingVectors({})
        if (result.success) {
          fixCount += result.data?.affected_count || 0
        }
      }

      toast.success(
        `${$t('admin.vectors.buttons.save')} ${fixCount} ${$t('admin.vectors.settings.indexingSettings.indexTimeout')}`
      )
      if (isPaused.value) {
        toast.info($t('admin.vectors.qdrant.systemInfo'))
      }
      emit('verification-completed')
    } catch (error: unknown) {
      toast.error((error as Error).message || $t('admin.vectors.actions.smartFixFailed'))
    } finally {
      verificationProcessing.value = false
    }
  }

  const handleStartVerification = async () => {
    try {
      verificationProcessing.value = true

      const result = await startVectorVerification({
        task_type: 'manual',
        needs_verification: true,
      })

      if (result.success && result.data) {
        currentVerificationTask.value = result.data
        toast.success($t('admin.vectors.toast.validationStarted'))
        startProgressMonitoring()
      }
    } catch (error: unknown) {
      toast.error((error as Error).message || $t('admin.vectors.actions.startVerificationFailed'))
    } finally {
      verificationProcessing.value = false
    }
  }

  const handleRegenerateAll = () => {
    confirmData.value = {
      title: $t('admin.vectors.actions.fullRebuildWarningTitle'),
      message: [
        $t('admin.vectors.actions.dangerWarning'),
        '',
        $t('admin.vectors.actions.fullRebuildConfirm'),
        $t('admin.vectors.actions.thisWill'),
        $t('admin.vectors.actions.clearAllVectors'),
        $t('admin.vectors.actions.reprocessAllFiles'),
        $t('admin.vectors.actions.mayTakeLongTime'),
        '',
        $t('admin.vectors.actions.confirmContinue'),
      ],
      type: 'danger',
      confirmText: $t('admin.vectors.actions.confirmRebuild'),
      action: 'regenerate-all',
    }
    showConfirmDialog.value = true
  }

  /* const handleRefresh = () => {
    if (isRefreshing.value) return

    isRefreshing.value = true
    emit('refresh')
    loadVerificationStats()

    setTimeout(() => {
      isRefreshing.value = false
    }, 1000)
  } */

  const handleReconcileMissing = async () => {
    try {
      const res = await reconcileMissing({ limit: 1000, dry_run: false })
      if (res.success) {
        toast.success(
          $t('admin.vectors.toast.missingVectorsFound')
            .replace('{found}', (res.data?.found ?? 0).toString())
            .replace('{enqueued}', (res.data?.enqueued ?? 0).toString())
        )
        emit('refresh')
      }
    } catch (e: unknown) {
      toast.error((e as Error)?.message || $t('admin.vectors.actions.reconcileMissingFailed'))
    }
  }

  const handleCleanOrphans = async () => {
    try {
      const res = await cleanOrphans({ limit: 1000, dry_run: false })
      if (res.success) {
        toast.success(
          $t('admin.vectors.toast.orphansRemoved')
            .replace('{found}', (res.data?.found ?? 0).toString())
            .replace('{removed}', (res.data?.removed ?? 0).toString())
        )
        emit('refresh')
      }
    } catch (e: unknown) {
      toast.error((e as Error)?.message || $t('admin.vectors.actions.cleanOrphansFailed'))
    }
  }

  const handleRebuildStale = async () => {
    try {
      const res = await rebuildStale({ limit: 1000 })
      if (res.success) {
        toast.success($t('admin.vectors.toast.outdatedRebuilt').replace('{enqueued}', (res.data?.enqueued ?? 0).toString()))
        emit('refresh')
      }
    } catch (e: unknown) {
      toast.error((e as Error)?.message || $t('admin.vectors.actions.rebuildStaleFailed'))
    }
  }

  const executeConfirmedAction = () => {
    showConfirmDialog.value = false

    switch (confirmData.value.action) {
      case 'smart-fix':
        executeSmartFix()
        break
      case 'regenerate-all':
        emit('regenerate-all')
        break
    }
  }

  const startProgressMonitoring = () => {
    if (verificationInterval.value) return

    verificationInterval.value = window.setInterval(async () => {
      try {
        const result = await getVerificationProgress()

        if (result.data && result.data.is_running) {
          currentVerificationTask.value = result.data.task
        } else {
          currentVerificationTask.value = null
          stopProgressMonitoring()
          await loadVerificationStats()
          emit('verification-completed')
          toast.success($t('admin.vectors.toast.validationCompleted'))
        }
      } catch {}
    }, 2000)
  }

  const stopProgressMonitoring = () => {
    if (verificationInterval.value) {
      clearInterval(verificationInterval.value)
      verificationInterval.value = null
    }
  }

  const loadVerificationStats = async () => {
    try {
      const result = await getVerificationStatistics()
      if (result.data) {
        verificationStats.value = result.data
      }
    } catch {}
  }

  onMounted(() => {
    loadVerificationStats()
  })

  onUnmounted(() => {
    stopProgressMonitoring()
  })

  const formatNumber = (num: number): string => {
    return num.toLocaleString()
  }

  const getSyncStatusClass = (): string => {
    if (!verificationStats.value) return ''
    const ratio = verificationStats.value.sync_ratio || 0
    if (ratio >= 95) return 'sync-excellent'
    if (ratio >= 80) return 'sync-good'
    return 'sync-poor'
  }

  const getSyncIcon = (): string => {
    if (!verificationStats.value) return 'fas fa-question-circle'
    const ratio = verificationStats.value.sync_ratio || 0
    if (ratio >= 95) return 'fas fa-check-circle'
    if (ratio >= 80) return 'fas fa-exclamation-circle'
    return 'fas fa-times-circle'
  }
</script>

<style scoped lang="scss">
  .vector-action-panel-simplified {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.8) 0%,
      rgba(var(--color-background-700-rgb), 0.6) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 24px;
    backdrop-filter: blur(10px);
    box-shadow: var(--shadow-cyber-md);
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    .panel-title {
      font-size: 1.2rem;
      font-weight: 600;
      color: var(--color-brand-500);
      margin: 0;
      display: flex;
      align-items: center;
    }

    .panel-status {
      .status-processing {
        color: var(--color-warning);
        font-size: 0.9rem;
        display: flex;
        align-items: center;
      }
    }
  }

  .action-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 20px;
    align-items: stretch;
  }

  .action-card {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
    padding: 20px;
    transition: all 0.3s ease;
    display: flex;
    flex-direction: column;
    height: 100%;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
    }

    &.primary {
      border-color: rgba(var(--color-brand-500-rgb), 0.3);
      .card-header i,
      .card-header span {
        color: var(--color-brand-500);
      }
    }

    &.secondary {
      border-color: var(--color-border-default);
      .card-header i,
      .card-header span {
        color: var(--color-content-muted);
      }
    }

    &.danger {
      border-color: rgba(var(--color-error-rgb), 0.3);
      .card-header i,
      .card-header span {
        color: var(--color-danger);
      }
    }

    &.info {
      border-color: rgba(var(--color-info-rgb), 0.3);
      .card-header i,
      .card-header span {
        color: var(--color-cyan-500);
      }
    }
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
    font-size: 1.1rem;
    font-weight: 600;

    i {
      font-size: 1.2rem;
    }
  }

  .card-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    margin-bottom: 16px;

    .description {
      color: var(--color-slate-400);
      font-size: 0.9rem;
      margin-bottom: 12px;
      line-height: 1.4;
    }

    .problem-stats,
    .quick-stats,
    .system-preview {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-top: auto;
      padding-top: 8px;
    }

    .problem-item,
    .stat-item,
    .preview-item {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 0.8rem;
      padding: 2px 6px;

      i {
        font-size: 0.75rem;
        width: 12px;
      }
      border-radius: var(--radius-sm);

      &.failed {
        color: var(--color-danger);
        background: rgba(var(--color-error-rgb), 0.1);
      }

      &.missing {
        color: var(--color-warning);
        background: rgba(var(--color-warning-rgb), 0.1);
      }

      &.stuck {
        color: var(--color-brand-500);
        background: rgba(var(--color-brand-500-rgb), 0.1);
      }
    }

    .no-problems {
      color: var(--color-success);
      font-size: 0.8rem;
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .verification-info {
      margin-top: auto;
      padding-top: 8px;

      .sync-status {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 0.85rem;
        font-weight: 500;

        &.sync-excellent {
          color: var(--color-success);
        }

        &.sync-good {
          color: var(--color-warning);
        }

        &.sync-poor {
          color: var(--color-danger);
        }
      }
    }

    .warning-text {
      color: var(--color-warning);
      font-size: 0.8rem;
      display: flex;
      align-items: center;
      gap: 4px;
      margin-top: auto;
      padding-top: 8px;
    }
  }

  .card-actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    margin-top: auto;

    .action-button {
      flex: 1;
      min-width: 120px;
      min-height: 32px;
      font-size: 0.875rem;
    }

    .action-button-small {
      width: 32px;
      height: 32px;
      padding: 0;
    }
  }

  .system-info-panel {
    margin-top: 24px;
    padding: 20px;
    background: rgba(var(--color-background-800-rgb), 0.4);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-sm);
  }

  .info-sections {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 20px;
    margin-bottom: 20px;
  }

  .info-section {
    .section-title {
      font-size: 1rem;
      font-weight: 600;
      color: var(--color-brand-500);
      margin: 0 0 12px 0;
      display: flex;
      align-items: center;
    }
  }

  .model-list {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .model-item {
      padding: 8px 12px;
      background: var(--color-hover-bg-neutral);
      border: 1px solid var(--color-border-default);
      border-radius: var(--radius-sm);
      display: flex;
      flex-direction: column;
      gap: 2px;

      .model-name {
        font-weight: 500;
        color: var(--color-white);
        font-size: 0.85rem;
      }

      .model-desc {
        font-size: 0.75rem;
        color: var(--color-slate-400);
      }
    }
  }

  .qdrant-stats,
  .verification-stats {
    display: flex;
    flex-direction: column;
    gap: 6px;

    .stat-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 0.85rem;

      .label {
        color: var(--color-slate-400);
      }

      .value {
        font-weight: 600;
        color: var(--color-white);

        &.primary {
          color: var(--color-brand-500);
        }
        &.success {
          color: var(--color-success);
        }
        &.warning {
          color: var(--color-warning);
        }
        &.error {
          color: var(--color-danger);
        }
      }
    }
  }

  .verification-progress {
    border-top: 1px solid var(--color-border-default);
    padding-top: 20px;

    .progress-info {
      .progress-text {
        font-size: 0.85rem;
        color: var(--color-slate-400);
        margin-bottom: 8px;
      }

      .progress-bar {
        height: 6px;
        background: var(--color-hover-bg-neutral);
        border-radius: var(--radius-sm);
        overflow: hidden;

        .progress-fill {
          height: 100%;
          background: linear-gradient(90deg, var(--color-brand-500) 0%, var(--color-success) 100%);
          transition: width 0.3s ease;
        }
      }
    }
  }

  .confirm-content {
    padding: 20px 0;

    .confirm-text {
      display: flex;
      align-items: flex-start;
      gap: 12px;

      .confirm-icon {
        font-size: 1.2rem;
        margin-top: 2px;
      }

      .message-lines {
        flex: 1;

        .message-line {
          margin: 0 0 4px 0;
          color: var(--color-white);
          line-height: 1.4;

          &:last-child {
            margin-bottom: 0;
          }
        }
      }
    }
  }

  .no-data {
    text-align: center;
    color: var(--color-content-muted);
    font-size: 0.8rem;
    padding: 12px;
  }

  .info-panel-enter-active,
  .info-panel-leave-active {
    transition: all 0.3s ease;
  }

  .info-panel-enter-from,
  .info-panel-leave-to {
    opacity: 0;
    transform: translateY(-10px);
  }

  @media (max-width: 768px) {
    .action-grid {
      grid-template-columns: 1fr;
    }

    .info-sections {
      grid-template-columns: 1fr;
    }
  }
</style>
