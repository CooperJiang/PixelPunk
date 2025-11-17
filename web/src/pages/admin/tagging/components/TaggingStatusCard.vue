<script setup lang="ts">
  import { computed } from 'vue'
  import { useWebSocketStore } from '@/store/websocket'
  import type { StatusStatsItem } from '@/api/admin/tagging'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps<{
    stats: StatusStatsItem[]
    selectedStatus?: string
  }>()

  const emit = defineEmits<{
    (e: 'select', status: string): void
  }>()

  const webSocketStore = useWebSocketStore()

  const stats = computed<StatusStatsItem[]>(() => {
    if (webSocketStore.queueStats) {
      const queueStats = webSocketStore.queueStats.queue_stats
      return [
        { status: 'none', count: queueStats.none_count || 0 },
        { status: 'pending', count: queueStats.pending_count || 0 },
        { status: 'done', count: queueStats.done_count || 0 },
        { status: 'failed', count: queueStats.failed_count || 0 },
        { status: 'ignored', count: queueStats.ignored_count || 0 },
        { status: 'skipped', count: 0 },
      ]
    }

    return props.stats || []
  })

  const getStatusName = (status: string): string => {
    return $t(`admin.tagging.status.${status}`) || status
  }

  const getStatusDescription = (status: string): string => {
    return $t(`admin.tagging.statusDesc.${status}`) || ''
  }

  const getIconName = (status: string): string => {
    const iconMap: Record<string, string> = {
      none: 'fas fa-hourglass-start',
      pending: 'fas fa-spinner fa-spin',
      done: 'fas fa-check-circle',
      failed: 'fas fa-exclamation-triangle',
      skipped: 'fas fa-forward',
      ignored: 'fas fa-ban',
    }
    return iconMap[status] || 'fas fa-question-circle'
  }

  const getIconClass = (status: string): string => {
    const classMap: Record<string, string> = {
      none: 'text-content-content-muted',
      pending: 'text-error-500',
      done: 'text-green-400',
      failed: 'text-red-400',
      skipped: 'text-content-content-muted',
      ignored: 'text-content-content-disabled',
    }
    return classMap[status] || ''
  }

  const getStatusClass = (status: string): string => {
    const isSelected = props.selectedStatus === status
    const baseClass = 'border-brand-200 bg-background-600'
    const selectedClass = 'border-error-500 shadow-glow-pink'

    const statusClasses: Record<string, string> = {
      none: isSelected ? selectedClass : baseClass,
      pending: isSelected ? selectedClass : baseClass,
      done: isSelected ? selectedClass : baseClass,
      failed: isSelected ? selectedClass : baseClass,
      skipped: isSelected ? selectedClass : baseClass,
      ignored: isSelected ? selectedClass : baseClass,
    }

    return statusClasses[status] || baseClass
  }

  const selectStatus = (status: string) => {
    emit('select', status === props.selectedStatus ? '' : status)
  }
</script>

<template>
  <div class="mb-4 grid grid-cols-2 gap-3 md:grid-cols-5">
    <div
      v-for="stat in stats"
      :key="stat.status"
      class="tagging-status-card hover:shadow-glow cursor-pointer rounded-lg border p-3 transition-all duration-300"
      :class="getStatusClass(stat.status)"
      @click="selectStatus(stat.status)"
    >
      <div class="mb-1 flex items-center justify-between">
        <h3 class="text-sm font-medium">{{ getStatusName(stat.status) }}</h3>
        <div :class="getIconClass(stat.status)">
          <i :class="getIconName(stat.status)" />
        </div>
      </div>
      <div class="text-xl font-bold">{{ stat.count }}</div>
      <div class="text-xs opacity-70">{{ getStatusDescription(stat.status) }}</div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .tagging-status-card {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    backdrop-filter: blur(10px);
  }

  .tagging-status-card:hover {
    transform: translateY(-2px);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
  }

  .shadow-glow {
    box-shadow: 0 0 12px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .shadow-glow-pink {
    box-shadow: 0 0 12px rgba(255, 0, 110, 0.4);
  }
</style>
