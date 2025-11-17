<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { useAutomationData } from './composables/useAutomationData'
  import TaggingStatusCard from './components/TaggingStatusCard.vue'
  import VectorStatusCard from './components/VectorStatusCard.vue'
  import TaggingHistoryTable from './components/TaggingHistoryTable.vue'

  defineOptions({
    name: 'AutomationPage',
  })

  const { $t } = useTexts()
  const { isLoading, overview, lastUpdated, refresh } = useAutomationData()

  /* 格式化最后更新时间 */
  const formattedLastUpdated = computed(() => {
    if (!lastUpdated.value) return ''
    return lastUpdated.value.toLocaleTimeString(getCurrentLocale())
  })

  /* 系统状态信息 */
  const systemStatus = computed(() => {
    if (!overview.value) return null
    return overview.value.system_status
  })
</script>

<template>
  <div class="automation-page">
    <div class="page-header">
      <div class="header-bg-pattern" />
      <div class="header-content">
        <div class="flex items-center gap-3">
          <div class="page-header-icon">
            <i class="fas fa-robot" />
            <div class="icon-glow" />
          </div>
          <div>
            <h1 class="header-title">{{ $t('automation.page.title') }}</h1>
            <p class="header-subtitle">
              {{ $t('automation.page.subtitle') }}
              <span v-if="formattedLastUpdated" class="update-time">
                <i class="fas fa-clock mr-1" />
                {{ formattedLastUpdated }}
              </span>
            </p>
          </div>
        </div>

        <div class="flex gap-2">
          <CyberButton type="outlined" :loading="isLoading" @click="refresh">
            <i class="fas fa-sync-alt mr-1.5" :class="{ 'fa-spin': isLoading }" />
            {{ $t('automation.page.refresh') }}
          </CyberButton>
        </div>
      </div>
    </div>

    <div v-if="isLoading && !overview" class="space-y-4">
      <CyberSkeleton type="card" :loading="true" />
      <CyberSkeleton type="card" :loading="true" />
    </div>

    <div v-else class="space-y-4">
      <TaggingStatusCard :stats="overview?.tagging || null" :system-enabled="systemStatus?.tagging_enabled || false" />

      <VectorStatusCard :stats="overview?.vector || null" :system-enabled="systemStatus?.vector_enabled || false" />

      <TaggingHistoryTable />

      <div class="info-tip">
        <div class="tip-icon">
          <i class="fas fa-lightbulb" />
        </div>
        <p>{{ $t('automation.page.tip') }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .automation-page {
    position: relative;
  }

  .page-header {
    position: relative;
    margin-bottom: 1.5rem;
    padding: 1.25rem 1.5rem;
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    backdrop-filter: blur(20px);
    overflow: hidden;
    box-shadow:
      0 1px 0 rgba(var(--color-brand-500-rgb), 0.1) inset,
      0 8px 16px rgba(0, 0, 0, 0.15),
      0 4px 8px rgba(var(--color-brand-500-rgb), 0.1),
      0 2px 4px rgba(0, 0, 0, 0.08);
    transform: translateZ(0);
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px) translateZ(0);
      box-shadow:
        0 1px 0 rgba(var(--color-brand-500-rgb), 0.15) inset,
        0 12px 24px rgba(0, 0, 0, 0.2),
        0 6px 12px rgba(var(--color-brand-500-rgb), 0.15),
        0 3px 6px rgba(0, 0, 0, 0.1);
    }
  }

  .header-bg-pattern {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image:
      radial-gradient(circle at 20% 50%, rgba(var(--color-brand-500-rgb), 0.03) 0%, transparent 50%),
      radial-gradient(circle at 80% 50%, rgba(var(--color-brand-500-rgb), 0.03) 0%, transparent 50%);
    pointer-events: none;
    animation: patternFloat 20s ease-in-out infinite;
  }

  @keyframes patternFloat {
    0%,
    100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-5px);
    }
  }

  .header-content {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    z-index: 1;
  }

  .page-header-icon {
    position: relative;
    width: 2.5rem;
    height: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(145deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.05));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    color: var(--color-brand-500);
    font-size: 1.125rem;
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.2) inset,
      0 -2px 0 rgba(0, 0, 0, 0.1) inset,
      0 4px 8px rgba(0, 0, 0, 0.15),
      0 2px 4px rgba(var(--color-brand-500-rgb), 0.2);
    transition: all 0.3s ease;
    transform: translateZ(0);

    i {
      position: relative;
      z-index: 2;
      filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
    }

    .icon-glow {
      position: absolute;
      inset: -2px;
      background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.4) 0%, transparent 70%);
      border-radius: var(--radius-sm);
      opacity: 0;
      transition: opacity 0.3s ease;
      animation: iconPulse 3s ease-in-out infinite;
    }

    &:hover {
      transform: translateY(-2px) scale(1.05) translateZ(0);
      box-shadow:
        0 3px 0 rgba(var(--color-brand-500-rgb), 0.3) inset,
        0 -2px 0 rgba(0, 0, 0, 0.15) inset,
        0 8px 16px rgba(0, 0, 0, 0.2),
        0 4px 8px rgba(var(--color-brand-500-rgb), 0.3);
    }
  }

  @keyframes iconPulse {
    0%,
    100% {
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
  }

  .header-title {
    font-size: 1.125rem;
    font-weight: 700;
    color: var(--color-content-heading);
    margin-bottom: 0.25rem;
    background: linear-gradient(135deg, var(--color-content-heading) 0%, var(--color-brand-500) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .header-subtitle {
    font-size: 0.8125rem;
    color: var(--color-content-muted);
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .update-time {
    display: inline-flex;
    align-items: center;
    padding: 0.125rem 0.5rem;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    font-size: 0.75rem;
    color: var(--color-brand-500);
  }

  .info-tip {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 1rem 1.25rem;
    background: linear-gradient(135deg, rgba(var(--color-warning-rgb), 0.08) 0%, rgba(var(--color-warning-rgb), 0.03) 100%);
    border: 1px solid rgba(var(--color-warning-rgb), 0.2);
    border-radius: var(--radius-sm);
    backdrop-filter: blur(10px);

    .tip-icon {
      flex-shrink: 0;
      width: 1.5rem;
      height: 1.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(var(--color-warning-rgb), 0.15);
      border-radius: var(--radius-full);
      color: var(--color-warning-500);
      font-size: 0.75rem;
    }

    p {
      flex: 1;
      margin: 0;
      font-size: 0.8125rem;
      line-height: 1.5;
      color: var(--color-content-default);
    }
  }
</style>
