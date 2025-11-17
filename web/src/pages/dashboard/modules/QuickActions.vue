<script setup lang="ts">
  import { computed, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { VueDraggable } from 'vue-draggable-plus'
  import { useTexts } from '@/composables/useTexts'
  import { useQuickActionsDrag } from '../composables/useQuickActionsDrag'

  defineOptions({
    name: 'QuickActions',
  })

  const router = useRouter()
  const { $t, currentLocale } = useTexts()
  const { createDraggableActions, STORAGE_KEY_COMMON, STORAGE_KEY_SETTINGS } = useQuickActionsDrag()

  /* 默认的常用操作列表 - 使用 computed 使其响应语言变化 */
  const defaultCommonActions = computed(() => [
    {
      id: 1,
      label: $t('dashboard.quickActions.actions.myFiles'),
      icon: 'fas fa-images',
      colorVar: 'brand',
      action: 'resource',
    },
    {
      id: 2,
      label: $t('dashboard.quickActions.actions.folders'),
      icon: 'fas fa-folder',
      colorVar: 'warning',
      action: 'folders',
    },
    {
      id: 3,
      label: $t('dashboard.quickActions.actions.tagManage'),
      icon: 'fas fa-tags',
      colorVar: 'info',
      action: 'tag-manage',
    },
    {
      id: 4,
      label: $t('dashboard.quickActions.actions.categoryManage'),
      icon: 'fas fa-layer-group',
      colorVar: 'success',
      action: 'category-manage',
    },
    {
      id: 5,
      label: $t('dashboard.quickActions.actions.openApi'),
      icon: 'fas fa-random',
      colorVar: 'error',
      action: 'open-api',
    },
    {
      id: 6,
      label: $t('dashboard.quickActions.actions.automation'),
      icon: 'fas fa-robot',
      colorVar: 'brand',
      action: 'automation',
    },
  ])

  /* 默认的设置快捷方式列表 - 使用 computed 使其响应语言变化 */
  const defaultSettingsShortcuts = computed(() => [
    {
      id: 1,
      title: $t('dashboard.quickActions.settings.api'),
      icon: 'fas fa-code',
      color: 'blue',
      route: '/settings#api',
    },
    {
      id: 2,
      title: $t('dashboard.quickActions.settings.profile'),
      icon: 'fas fa-user',
      color: 'green',
      route: '/settings#profile',
    },
    {
      id: 3,
      title: $t('dashboard.quickActions.settings.security'),
      icon: 'fas fa-shield-alt',
      color: 'orange',
      route: '/settings#security',
    },
    {
      id: 4,
      title: $t('dashboard.quickActions.settings.accessControl'),
      icon: 'fas fa-shield',
      color: 'purple',
      route: '/settings#access-control',
    },
    {
      id: 5,
      title: $t('dashboard.quickActions.settings.preferences'),
      icon: 'fas fa-cog',
      color: 'purple',
      route: '/settings#preferences',
    },
  ])

  const {
    items: commonActions,
    handleDragStart: handleCommonDragStart,
    handleDragEnd: handleCommonDragEnd,
  } = createDraggableActions(STORAGE_KEY_COMMON, defaultCommonActions.value)

  const {
    items: settingsShortcuts,
    handleDragStart: handleSettingsDragStart,
    handleDragEnd: handleSettingsDragEnd,
  } = createDraggableActions(STORAGE_KEY_SETTINGS, defaultSettingsShortcuts.value)

  // 监听语言变化，更新翻译
  watch(currentLocale, () => {
    // 更新 commonActions 的 label
    const currentOrder = commonActions.value.map((item) => item.id)
    const newActions = defaultCommonActions.value
    commonActions.value = currentOrder.map((id) => newActions.find((action) => action.id === id)!).filter(Boolean)

    // 更新 settingsShortcuts 的 title
    const currentSettingsOrder = settingsShortcuts.value.map((item) => item.id)
    const newSettings = defaultSettingsShortcuts.value
    settingsShortcuts.value = currentSettingsOrder.map((id) => newSettings.find((setting) => setting.id === id)!).filter(Boolean)
  })

  const handleAction = (action: string) => {
    switch (action) {
      case 'resource':
        router.push('/resource')
        break
      case 'folders':
        router.push('/folders')
        break
      case 'tag-manage':
        router.push('/tag-manage')
        break
      case 'category-manage':
        router.push('/category-manage')
        break
      case 'open-api':
        router.push('/open-api')
        break
      case 'automation':
        router.push('/automation')
        break
    }
  }

  const handleNavigation = (route: string) => {
    router.push(route)
  }
</script>

<template>
  <div class="quick-actions cyber-card">
    <div class="cyber-bg-pattern" />
    <div class="glow-orb" />

    <div class="actions-content">
      <div class="actions-header">
        <h3 class="section-title">
          <i class="fas fa-bolt" />
          {{ $t('dashboard.quickActions.title') }}
          <span class="drag-hint">
            <i class="fas fa-arrows-alt" />
            {{ $t('dashboard.quickActions.dragHint') }}
          </span>
        </h3>
      </div>

      <div class="main-actions">
        <VueDraggable
          v-model="commonActions"
          animation="150"
          ghost-class="action-ghost"
          chosen-class="action-chosen"
          drag-class="action-drag"
          class="actions-grid"
          @start="handleCommonDragStart"
          @end="handleCommonDragEnd"
        >
          <button
            v-for="action in commonActions"
            :key="action.id"
            class="action-btn primary-action"
            :data-color="action.colorVar"
            @click="handleAction(action.action)"
          >
            <div class="btn-bg-effect" />
            <div class="btn-content">
              <div class="btn-icon">
                <i :class="action.icon" />
              </div>
              <span class="btn-label">{{ action.label }}</span>
              <div class="btn-glow" />
            </div>
          </button>
        </VueDraggable>
      </div>

      <div class="settings-section">
        <div class="section-divider">
          <div class="divider-line" />
          <span class="divider-text">{{ $t('dashboard.quickActions.settings.divider') }}</span>
          <div class="divider-line" />
        </div>
        <VueDraggable
          v-model="settingsShortcuts"
          animation="150"
          ghost-class="setting-ghost"
          chosen-class="setting-chosen"
          drag-class="setting-drag"
          class="settings-grid"
          @start="handleSettingsDragStart"
          @end="handleSettingsDragEnd"
        >
          <button
            v-for="setting in settingsShortcuts"
            :key="setting.id"
            class="setting-btn"
            :class="setting.color"
            :title="setting.title"
            @click="handleNavigation(setting.route)"
          >
            <div class="setting-bg-effect" />
            <div class="setting-icon">
              <i :class="setting.icon" />
            </div>
            <span class="setting-label">{{ setting.title }}</span>
          </button>
        </VueDraggable>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .quick-actions {
    padding: 1.5rem;
    position: relative;
    overflow: hidden;
    background: linear-gradient(145deg, rgba(var(--color-background-800-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.15) inset,
      0 -1px 0 rgba(0, 0, 0, 0.15) inset,
      0 10px 20px rgba(0, 0, 0, 0.2),
      0 5px 10px rgba(var(--color-brand-500-rgb), 0.15),
      0 2px 5px rgba(0, 0, 0, 0.1);
    backdrop-filter: var(--backdrop-blur-md);

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      height: 4px;
      background: linear-gradient(90deg, transparent, var(--color-brand-500), transparent);
      opacity: 0.8;
      box-shadow: 0 2px 6px rgba(var(--color-brand-500-rgb), 0.6);
      z-index: 2;
    }
  }

  .cyber-bg-pattern {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image:
      radial-gradient(circle at 20% 30%, rgba(var(--color-brand-500-rgb), 0.03) 0%, transparent 50%),
      radial-gradient(circle at 80% 70%, rgba(var(--color-brand-500-rgb), 0.03) 0%, transparent 50%),
      linear-gradient(45deg, rgba(var(--color-brand-500-rgb), 0.02) 1px, transparent 1px),
      linear-gradient(-45deg, rgba(var(--color-brand-500-rgb), 0.02) 1px, transparent 1px);
    background-size:
      100% 100%,
      100% 100%,
      20px 20px,
      20px 20px;
    animation: patternShift 20s ease-in-out infinite;
    pointer-events: none;
  }

  .glow-orb {
    position: absolute;
    top: -30%;
    right: -30%;
    width: 60%;
    height: 60%;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.1) 0%, transparent 70%);
    border-radius: var(--radius-full);
    opacity: 0.6;
    animation: orbFloat 15s ease-in-out infinite;
    pointer-events: none;
  }

  @keyframes patternShift {
    0%,
    100% {
      transform: translate(0, 0);
    }
    50% {
      transform: translate(5px, 5px);
    }
  }

  @keyframes orbFloat {
    0%,
    100% {
      transform: translate(0, 0) scale(1);
      opacity: 0.6;
    }
    50% {
      transform: translate(-10px, -10px) scale(1.1);
      opacity: 0.4;
    }
  }

  .actions-content {
    position: relative;
    z-index: 1;
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .actions-header {
    margin-bottom: 8px;
  }

  .section-title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .section-title i {
    color: var(--color-brand-500);
    font-size: 14px;
  }

  .drag-hint {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    margin-left: 12px;
    padding: 3px 8px;
    background: rgba(var(--color-warning-rgb), 0.1);
    border: 1px solid rgba(var(--color-warning-rgb), 0.3);
    border-radius: var(--radius-sm);
    font-size: 10px;
    font-weight: 500;
    color: var(--color-warning-500);
    transition: all 0.2s ease;
    opacity: 0.7;
  }

  .drag-hint:hover {
    opacity: 1;
    background: rgba(var(--color-warning-rgb), 0.15);
  }

  .drag-hint i {
    font-size: 9px;
    color: var(--color-warning-500);
  }

  .main-actions {
    margin-bottom: 8px;
  }

  .actions-grid {
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 8px;
  }

  .primary-action {
    position: relative;
    border-radius: var(--radius-sm);
    cursor: pointer;
    overflow: hidden;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    background: linear-gradient(145deg, rgba(var(--color-background-700-rgb), 0.8), rgba(var(--color-background-700-rgb), 0.6));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    min-height: 70px;
    padding: 0;
    backdrop-filter: blur(10px);
    box-shadow:
      0 1px 0 rgba(255, 255, 255, 0.1) inset,
      0 -1px 0 rgba(0, 0, 0, 0.1) inset,
      0 3px 6px rgba(0, 0, 0, 0.15);
    transform: translateZ(0);

    &:hover {
      border-color: rgba(var(--color-brand-500-rgb), 0.3);
      box-shadow:
        0 1px 0 rgba(255, 255, 255, 0.15) inset,
        0 -1px 0 rgba(0, 0, 0, 0.15) inset,
        0 5px 10px rgba(0, 0, 0, 0.2);
    }
  }

  .btn-bg-effect {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.05) 0%, rgba(var(--color-brand-500-rgb), 0.05) 100%);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .primary-action:hover .btn-bg-effect {
    opacity: 1;
  }

  .primary-action:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow:
      0 10px 30px rgba(0, 0, 0, 0.3),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .btn-content {
    position: relative;
    z-index: 2;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6px;
    height: 100%;
    padding: 12px 8px;
  }

  .btn-icon {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    font-size: 14px;
    transition: all 0.3s ease;
    position: relative;
  }

  .btn-glow {
    position: absolute;
    top: -2px;
    left: -2px;
    right: -2px;
    bottom: -2px;
    border-radius: var(--radius-sm);
    opacity: 0;
    transition: opacity 0.3s ease;
    pointer-events: none;
  }

  .primary-action:hover .btn-glow {
    opacity: 1;
  }

  .btn-label {
    font-size: 11px;
    font-weight: 600;
    text-align: center;
    line-height: 1.2;
    transition: all 0.3s ease;
  }

  .btn-label {
    color: var(--color-content-default);
  }

  .primary-action[data-color='brand'] .btn-icon {
    background: linear-gradient(145deg, rgba(var(--color-brand-500-rgb), 0.2), rgba(var(--color-brand-500-rgb), 0.1));
    color: var(--color-brand-500);
    box-shadow:
      0 1px 0 rgba(var(--color-brand-500-rgb), 0.3) inset,
      0 2px 4px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .primary-action[data-color='warning'] .btn-icon {
    background: linear-gradient(145deg, rgba(var(--color-warning-rgb), 0.2), rgba(var(--color-warning-rgb), 0.1));
    color: var(--color-warning-500);
    box-shadow:
      0 1px 0 rgba(var(--color-warning-rgb), 0.3) inset,
      0 2px 4px rgba(var(--color-warning-rgb), 0.3);
  }

  .primary-action[data-color='info'] .btn-icon {
    background: linear-gradient(145deg, rgba(var(--color-info-rgb), 0.2), rgba(var(--color-info-rgb), 0.1));
    color: var(--color-info-500);
    box-shadow:
      0 1px 0 rgba(var(--color-info-rgb), 0.3) inset,
      0 2px 4px rgba(var(--color-info-rgb), 0.3);
  }

  .primary-action[data-color='success'] .btn-icon {
    background: linear-gradient(145deg, rgba(var(--color-success-rgb), 0.2), rgba(var(--color-success-rgb), 0.1));
    color: var(--color-success-500);
    box-shadow:
      0 1px 0 rgba(var(--color-success-rgb), 0.3) inset,
      0 2px 4px rgba(var(--color-success-rgb), 0.3);
  }

  .primary-action[data-color='error'] .btn-icon {
    background: linear-gradient(145deg, rgba(var(--color-error-rgb), 0.2), rgba(var(--color-error-rgb), 0.1));
    color: var(--color-error-500);
    box-shadow:
      0 1px 0 rgba(var(--color-error-rgb), 0.3) inset,
      0 2px 4px rgba(var(--color-error-rgb), 0.3);
  }

  .primary-action:hover[data-color='brand'] .btn-icon {
    background: rgba(var(--color-brand-500-rgb), 0.25);
    box-shadow: 0 0 20px rgba(var(--color-brand-500-rgb), 0.4);
  }

  .primary-action:hover[data-color='warning'] .btn-icon {
    background: rgba(var(--color-warning-rgb), 0.25);
    box-shadow: 0 0 20px rgba(var(--color-warning-rgb), 0.4);
  }

  .primary-action:hover[data-color='info'] .btn-icon {
    background: rgba(var(--color-info-rgb), 0.25);
    box-shadow: 0 0 20px rgba(var(--color-info-rgb), 0.4);
  }

  .primary-action:hover[data-color='success'] .btn-icon {
    background: rgba(var(--color-success-rgb), 0.25);
    box-shadow: 0 0 20px rgba(var(--color-success-rgb), 0.4);
  }

  .primary-action:hover[data-color='error'] .btn-icon {
    background: rgba(var(--color-error-rgb), 0.25);
    box-shadow: 0 0 20px rgba(var(--color-error-rgb), 0.4);
  }

  .primary-action:hover[data-color='brand'] .btn-label {
    color: var(--color-brand-500);
  }

  .primary-action:hover[data-color='warning'] .btn-label {
    color: var(--color-warning-500);
  }

  .primary-action:hover[data-color='info'] .btn-label {
    color: var(--color-info-500);
  }

  .primary-action:hover[data-color='success'] .btn-label {
    color: var(--color-success-500);
  }

  .primary-action:hover[data-color='error'] .btn-label {
    color: var(--color-error-500);
  }

  .primary-action[data-color='brand'] .btn-glow {
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.3) 0%, transparent 70%);
    box-shadow: 0 0 25px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .primary-action[data-color='warning'] .btn-glow {
    background: radial-gradient(circle, rgba(var(--color-warning-rgb), 0.3) 0%, transparent 70%);
    box-shadow: 0 0 25px rgba(var(--color-warning-rgb), 0.3);
  }

  .primary-action[data-color='info'] .btn-glow {
    background: radial-gradient(circle, rgba(var(--color-info-rgb), 0.3) 0%, transparent 70%);
    box-shadow: 0 0 25px rgba(var(--color-info-rgb), 0.3);
  }

  .primary-action[data-color='success'] .btn-glow {
    background: radial-gradient(circle, rgba(var(--color-success-rgb), 0.3) 0%, transparent 70%);
    box-shadow: 0 0 25px rgba(var(--color-success-rgb), 0.3);
  }

  .primary-action[data-color='error'] .btn-glow {
    background: radial-gradient(circle, rgba(var(--color-error-rgb), 0.3) 0%, transparent 70%);
    box-shadow: 0 0 25px rgba(var(--color-error-rgb), 0.3);
  }

  .section-divider {
    display: flex;
    align-items: center;
    gap: 12px;
    margin: 8px 0 16px 0;
  }

  .divider-line {
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent 0%, rgba(var(--color-brand-500-rgb), 0.3) 50%, transparent 100%);
  }

  .divider-text {
    font-size: 11px;
    font-weight: 600;
    color: rgba(var(--color-brand-500-rgb), 0.7);
    text-transform: uppercase;
    letter-spacing: 1px;
    padding: 0 8px;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1) 0%, rgba(var(--color-brand-500-rgb), 0.1) 100%);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    backdrop-filter: blur(10px);
  }

  .settings-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 6px;
    max-width: 100%;
  }

  .setting-btn {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    padding: 12px 8px;
    background: linear-gradient(145deg, rgba(var(--color-background-700-rgb), 0.9), rgba(var(--color-background-800-rgb), 0.8));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    min-height: 70px;
    overflow: hidden;
    backdrop-filter: blur(10px);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.1) inset,
      0 -1px 0 rgba(0, 0, 0, 0.2) inset,
      0 4px 8px rgba(0, 0, 0, 0.2),
      0 2px 4px rgba(var(--color-brand-500-rgb), 0.1);
  }

  .setting-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.5), transparent);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .setting-btn:hover::before {
    opacity: 1;
  }

  .setting-bg-effect {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .setting-btn:hover .setting-bg-effect {
    opacity: 1;
  }

  .setting-btn:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.2) inset,
      0 -1px 0 rgba(0, 0, 0, 0.3) inset,
      0 8px 20px rgba(0, 0, 0, 0.4),
      0 0 20px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .setting-icon {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    font-size: 12px;
    transition: all 0.3s ease;
    position: relative;
    z-index: 2;
  }

  .setting-label {
    font-size: 10px;
    font-weight: 600;
    text-align: center;
    line-height: 1.2;
    transition: all 0.3s ease;
    position: relative;
    z-index: 2;
  }

  .setting-bg-effect {
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-600-rgb), 0.08));
  }

  .setting-btn:hover .setting-bg-effect {
    opacity: 1;
  }

  .setting-icon {
    background: linear-gradient(145deg, rgba(var(--color-background-600-rgb), 0.8), rgba(var(--color-background-700-rgb), 0.6));
    color: var(--color-content-default);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-500-rgb), 0.05) inset,
      0 -1px 0 rgba(0, 0, 0, 0.2) inset,
      0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .setting-btn:hover .setting-icon {
    background: linear-gradient(145deg, rgba(var(--color-brand-500-rgb), 0.25), rgba(var(--color-brand-600-rgb), 0.15));
    color: var(--color-brand-300);
    box-shadow:
      0 2px 0 rgba(var(--color-brand-400-rgb), 0.4) inset,
      0 -1px 0 rgba(0, 0, 0, 0.3) inset,
      0 0 15px rgba(var(--color-brand-500-rgb), 0.5),
      0 2px 6px rgba(var(--color-brand-500-rgb), 0.4);
    transform: scale(1.1);
  }

  .setting-label {
    color: var(--color-content-default);
  }

  .setting-btn:hover .setting-label {
    color: var(--color-brand-400);
    text-shadow: 0 0 8px rgba(var(--color-brand-500-rgb), 0.5);
  }

  @media (max-width: 768px) {
    .actions-grid {
      grid-template-columns: repeat(2, 1fr);
      gap: 8px;
    }

    .primary-action {
      min-height: 60px;
    }

    .btn-content {
      padding: 10px 6px;
      gap: 5px;
    }

    .btn-icon {
      width: 28px;
      height: 28px;
      font-size: 12px;
    }

    .btn-label {
      font-size: 10px;
    }

    .settings-grid {
      grid-template-columns: repeat(3, 1fr);
      gap: 5px;
    }

    .setting-btn {
      padding: 10px 6px;
      min-height: 60px;
    }

    .setting-icon {
      width: 24px;
      height: 24px;
      font-size: 11px;
    }

    .setting-label {
      font-size: 9px;
    }
  }

  @media (max-width: 480px) {
    .quick-actions {
      padding: 16px;
    }

    .actions-content {
      gap: 16px;
    }

    .section-title {
      font-size: 14px;
    }

    .actions-grid {
      grid-template-columns: 1fr;
      gap: 6px;
    }

    .primary-action {
      min-height: 55px;
    }

    .settings-grid {
      grid-template-columns: repeat(5, 1fr);
      gap: 5px;
    }

    .setting-btn {
      padding: 8px 4px;
      min-height: 55px;
    }

    .setting-icon {
      width: 20px;
      height: 20px;
      font-size: 10px;
    }

    .setting-label {
      font-size: 8px;
    }
  }

  .action-ghost {
    opacity: 0.4 !important;
    background: rgba(var(--color-brand-500-rgb), 0.1) !important;
    border: 2px dashed rgba(var(--color-brand-500-rgb), 0.5) !important;
  }

  .action-chosen {
    cursor: move !important;
    transform: scale(1.05);
    box-shadow:
      0 8px 24px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.2) !important;
  }

  .action-drag {
    opacity: 0.8 !important;
    transform: rotate(3deg) !important;
    cursor: grabbing !important;
  }

  .setting-ghost {
    opacity: 0.4 !important;
    background: rgba(var(--color-brand-500-rgb), 0.1) !important;
    border: 2px dashed rgba(var(--color-brand-500-rgb), 0.5) !important;
  }

  .setting-chosen {
    cursor: move !important;
    transform: scale(1.05);
    box-shadow:
      0 8px 24px rgba(var(--color-brand-500-rgb), 0.3),
      0 0 30px rgba(var(--color-brand-500-rgb), 0.2) !important;
  }

  .setting-drag {
    opacity: 0.8 !important;
    transform: rotate(3deg) !important;
    cursor: grabbing !important;
  }

  .actions-grid,
  .settings-grid {
    cursor: grab;
  }

  .actions-grid:active,
  .settings-grid:active {
    cursor: grabbing;
  }
</style>
