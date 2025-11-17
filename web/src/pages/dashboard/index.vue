<script setup lang="ts">
  import { onMounted, ref } from 'vue'
  import QuickActions from './modules/QuickActions.vue'
  import CompactStats from './modules/CompactStats.vue'
  import ActivityMonitor from './modules/ActivityMonitor.vue'
  import MessageCenter from './modules/MessageCenter.vue'
  import UploadQueue from './modules/UploadQueue.vue'
  import RecentUploads from './modules/RecentUploads.vue'
  import FolderDistribution from './modules/FolderDistribution.vue'

  defineOptions({
    name: 'WorkspaceDashboard',
  })

  const refreshKey = ref(0)

  onMounted(() => {
    setInterval(() => {
      refreshKey.value++
    }, 30000)
  })
</script>

<template>
  <div class="workspace-dashboard">
    <div class="cyber-decorations">
      <div class="corner-decoration top-left" />
      <div class="corner-decoration top-right" />
    </div>

    <div class="workspace-main">
      <div class="main-workspace">
        <div class="left-column">
          <div class="main-actions">
            <QuickActions />
          </div>
          <div class="stats-section">
            <CompactStats :key="refreshKey" />
          </div>
          <div class="recent-section">
            <RecentUploads />
          </div>
          <div class="folder-section">
            <FolderDistribution />
          </div>
        </div>
        <div class="side-panels">
          <div class="side-item">
            <MessageCenter :key="refreshKey" />
          </div>
          <div class="side-item">
            <ActivityMonitor />
          </div>
          <div class="side-item">
            <UploadQueue />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .workspace-dashboard {
    min-height: 100vh;
    position: relative;
    padding: 0;
    margin: 0;
    background:
      radial-gradient(circle at 20% 10%, rgba(var(--color-brand-500-rgb), 0.03) 0%, transparent 50%),
      radial-gradient(circle at 80% 90%, rgba(var(--color-brand-500-rgb), 0.02) 0%, transparent 50%);
  }

  .cyber-decorations {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    z-index: 1;
    opacity: 0.6;
  }

  .corner-decoration {
    position: absolute;
    width: 80px;
    height: 80px;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12) 0%, transparent 70%);
    clip-path: polygon(0 0, 100% 0, 0 100%);
    animation: pulseGlow 4s ease-in-out infinite;
  }

  .corner-decoration.top-left {
    top: 0;
    left: 0;
    transform: rotate(0deg);
  }

  .corner-decoration.top-right {
    top: 0;
    right: 0;
    transform: rotate(90deg);
    animation-delay: 2s;
  }

  @keyframes pulseGlow {
    0%,
    100% {
      opacity: 0.6;
      filter: blur(0px);
    }
    50% {
      opacity: 1;
      filter: blur(2px);
    }
  }

  .workspace-main {
    position: relative;
    z-index: 2;
    padding: 0;
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .main-workspace {
    display: grid;
    grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
    gap: 16px;
    width: 100%;
    align-items: start;
  }

  .left-column {
    display: flex;
    flex-direction: column;
    gap: 16px;
    width: 100%;
    min-width: 0;
  }

  .main-actions {
    width: 100%;
    animation: slideInUp 0.6s ease-out;
  }

  @keyframes slideInUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .side-panels {
    display: flex;
    flex-direction: column;
    gap: 16px;
    width: 100%;
    min-width: 0;
  }

  .side-item {
    width: 100%;
    min-height: 140px;
    animation: slideInRight 0.6s ease-out;
  }

  .side-item:nth-child(1) {
    animation-delay: 0.1s;
  }
  .side-item:nth-child(2) {
    animation-delay: 0.2s;
  }
  .side-item:nth-child(3) {
    animation-delay: 0.3s;
  }

  @keyframes slideInRight {
    from {
      opacity: 0;
      transform: translateX(30px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .stats-section {
    width: 100%;
    animation: slideInUp 0.6s ease-out 0.1s backwards;
  }

  .queue-section {
    width: 100%;
  }

  .recent-section {
    width: 100%;
    animation: slideInUp 0.6s ease-out 0.2s backwards;
  }

  .folder-section {
    width: 100%;
    min-height: 420px;
    animation: slideInUp 0.6s ease-out 0.3s backwards;
  }

  :deep(.cyber-card) {
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.95), rgba(var(--color-background-700-rgb), 0.8));
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-default);
    backdrop-filter: blur(12px);
    box-shadow:
      var(--shadow-cyber-md),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.1) inset;
    position: relative;
    overflow: hidden;
  }

  :deep(.cyber-card::before) {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.5), transparent);
    opacity: 0.8;
    box-shadow: 0 2px 8px rgba(var(--color-brand-500-rgb), 0.4);
  }

  :deep(.cyber-card:hover) {
  }

  :deep(.cyber-card):nth-child(even) {
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.9), rgba(var(--color-background-700-rgb), 0.85));
  }

  @media (max-width: 768px) {
    .workspace-main {
      padding: 16px;
      gap: 14px;
    }

    .main-workspace {
      grid-template-columns: 1fr;
      gap: 16px;
    }

    .left-column {
      gap: 16px;
    }

    .side-panels {
      flex-direction: row;
      gap: 16px;
    }

    .side-item {
      flex: 1;
      min-height: 160px;
    }

    .corner-decoration {
      width: 60px;
      height: 60px;
    }
  }

  @media (max-width: 640px) {
    .workspace-main {
      padding: 12px;
      gap: 12px;
    }

    .corner-decoration {
      width: 40px;
      height: 40px;
    }
  }
</style>
