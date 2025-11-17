<script setup lang="ts">
  import { Dialog } from '@/components/Dialog'
  import { useChannelManagement } from './composables/useChannelManagement'
  import { useTexts } from '@/composables/useTexts'

  import ChannelCard from './components/ChannelCard.vue'
  import ChannelForm from './components/ChannelForm.vue'

  defineOptions({
    name: 'ChannelsManagePage',
  })

  const { $t } = useTexts()

  const {
    channels,
    loading,
    sortedChannels,
    fetchChannels,

    currentChannel,
    channelForm,
    showAddChannelDialog,

    isDragOver,

    fileInputRef,

    channelOperations,
    formOperations,
    fileOperations,
  } = useChannelManagement()

  import { onMounted } from 'vue'

  onMounted(() => {
    fetchChannels()
  })
</script>

<template>
  <div class="admin-channels-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.channels.title')"
      :subtitle="$t('admin.channels.subtitle')"
      icon="fas fa-cloud-upload-alt"
    >
      <template #actions>
        <div class="action-buttons">
          <CyberButton
            type="primary"
            icon="plus"
            @click="
              () => {
                formOperations.resetForm()
                showAddChannelDialog = true
              }
            "
          >
            {{ $t('admin.channels.actions.addChannel') }}
          </CyberButton>

          <div class="button-divider" />

          <CyberButton type="ghost" icon="file-import" @click="fileOperations.importConfig">
            {{ $t('admin.channels.actions.importConfig') }}
          </CyberButton>

          <CyberButton type="ghost" icon="file-export" @click="fileOperations.exportAll">
            {{ $t('admin.channels.actions.exportAll') }}
          </CyberButton>

          <div class="button-divider" />

          <CyberButton
            type="ghost"
            :disabled="loading"
            :loading="loading"
            loading-mode="inline"
            icon="sync-alt"
            class="refresh-btn"
            :title="loading ? $t('admin.channels.actions.refreshing') : $t('admin.channels.actions.refresh')"
            @click="fetchChannels"
          >
            <span class="ml-2 hidden sm:inline">{{
              loading ? $t('admin.channels.actions.refreshing') : $t('admin.channels.actions.refresh')
            }}</span>
          </CyberButton>
        </div>
      </template>

      <template #toolbar>
        <div class="shortcut-tips">
          <div class="tip-item">
            <i class="fas fa-keyboard text-content" />
            <span
              >{{ $t('admin.channels.shortcuts.title') }} <kbd>Ctrl/Cmd + V</kbd> {{ $t('admin.channels.shortcuts.paste') }}</span
            >
          </div>
          <div class="tip-divider" />
          <div class="tip-item">
            <i class="fas fa-upload text-error-500" />
            <span>{{ $t('admin.channels.shortcuts.dragHint') }}</span>
          </div>
          <div class="tip-divider" />
          <div class="tip-item">
            <i class="fas fa-info-circle text-secondary" />
            <span>{{ $t('admin.channels.tips.multipleServices') }}</span>
          </div>
        </div>
      </template>

      <template #default>
        <div class="channels-content" :class="{ 'drag-over': isDragOver }">
          <div v-if="isDragOver" class="drag-overlay">
            <div class="drag-content">
              <i class="fas fa-cloud-upload-alt drag-icon" />
              <h3 class="drag-title">{{ $t('admin.channels.drag.title') }}</h3>
              <p class="drag-desc">{{ $t('admin.channels.drag.desc') }}</p>
            </div>
          </div>

          <div class="channels-grid">
            <div v-if="loading" class="loading-state">
              <CyberLoading :visible="true" :text="$t('admin.channels.loading')" />
            </div>

            <template v-else>
              <ChannelCard
                v-for="channel in sortedChannels"
                :key="channel.id"
                :channel="channel"
                @edit="formOperations.editChannel"
                @export="fileOperations.exportChannel"
                @toggle-status="channelOperations.toggleStatus"
                @set-default="channelOperations.setDefault"
                @delete="channelOperations.deleteChannel"
                @test="channelOperations.testConnection"
              />

              <div v-if="channels.length === 0" class="empty-content">
                <div class="empty-icon-wrapper">
                  <i class="fas fa-cloud-upload-alt" />
                  <div class="icon-glow" />
                </div>
                <h3 class="empty-title">{{ $t('admin.channels.empty.title') }}</h3>
                <p class="empty-desc">{{ $t('admin.channels.empty.desc') }}</p>
                <div class="empty-actions">
                  <CyberButton
                    type="primary"
                    size="medium"
                    icon="plus"
                    class="add-channel-btn"
                    @click="
                      () => {
                        formOperations.resetForm()
                        showAddChannelDialog = true
                      }
                    "
                  >
                    {{ $t('admin.channels.empty.action') }}
                  </CyberButton>
                  <p class="empty-tip">
                    <i class="fas fa-lightbulb mr-1" />
                    {{ $t('admin.channels.empty.tip') }}
                  </p>
                </div>
              </div>
            </template>
          </div>

          <Dialog
            v-model="showAddChannelDialog"
            :title="currentChannel ? $t('admin.channels.dialog.edit') : $t('admin.channels.dialog.add')"
            width="500px"
          >
            <ChannelForm
              :form="channelForm"
              :is-edit="!!currentChannel"
              @update:form="(newForm) => Object.assign(channelForm, newForm)"
            />

            <template #footer>
              <div class="dialog-footer flex justify-end space-x-2 p-4">
                <CyberButton type="secondary" class="text-xs" @click="showAddChannelDialog = false">
                  {{ $t('admin.channels.actions.cancel') }}
                </CyberButton>
                <CyberButton type="primary" class="text-xs" @click="formOperations.saveChannel">
                  {{ $t('admin.channels.actions.save') }}
                </CyberButton>
              </div>
            </template>
          </Dialog>

          <input ref="fileInputRef" type="file" accept=".json" style="display: none" @change="fileOperations.handleFileImport" />
        </div>
      </template>
    </CyberAdminWrapper>
  </div>
</template>

<style scoped lang="scss">
  .admin-channels-page {
    color: var(--color-content);
  }

  .channels-content {
    min-height: 500px;
    position: relative;
    color: var(--color-content);
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;
  }

  .channels-content.drag-over {
    @apply relative;
  }

  .drag-overlay {
    position: fixed;
    inset: 0;
    z-index: 50;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(var(--color-background-900-rgb), 0.95);
    backdrop-filter: var(--backdrop-blur-sm);
  }

  .drag-content {
    border-radius: var(--radius-lg);
    border: 1px dashed;
    padding: var(--space-3xl);
    text-align: center;
    border-color: rgba(var(--color-brand-500-rgb), 0.6);
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.98) 0%,
      rgba(var(--color-background-900-rgb), 0.95) 50%,
      rgba(var(--color-background-800-rgb), 0.98) 100%
    );
    animation: drag-pulse 2s var(--ease-in-out) infinite;
  }

  .drag-icon {
    margin-bottom: var(--space-md);
    font-size: var(--text-4xl);
    background: linear-gradient(45deg, var(--color-brand-500), var(--color-error-500));
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    animation: float 3s var(--ease-in-out) infinite;
  }

  .drag-title {
    margin-bottom: var(--space-sm);
    font-size: var(--text-xl);
    font-weight: var(--font-bold);
    background: linear-gradient(45deg, var(--color-brand-500), var(--color-error-500));
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
  }

  .drag-desc {
    font-size: var(--text-sm);
    opacity: var(--opacity-muted);
    color: var(--color-content-muted);
  }

  .shortcut-tips {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-md);
    padding: var(--space-md);
    background: rgba(var(--color-background-800-rgb), 0.3);
    border-radius: var(--radius-md);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    backdrop-filter: var(--backdrop-blur-sm);
  }

  .tip-item {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    font-size: var(--text-sm);
    color: var(--color-content-muted);
  }

  .tip-item i.text-content {
    color: var(--color-brand-500);
  }

  .tip-item i.text-accent {
    color: var(--color-error-500);
  }

  .tip-item i.text-secondary {
    color: var(--color-content-muted);
  }

  .tip-divider {
    height: var(--space-md);
    width: 1px;
    background: rgba(var(--color-brand-500-rgb), 0.2);
  }

  kbd {
    border-radius: var(--radius-sm);
    padding: var(--space-xs) var(--space-sm);
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    background: rgba(var(--color-background-800-rgb), 0.5);
    color: var(--color-brand-500);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    box-shadow: var(--shadow-cyber-sm);
  }

  .channels-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(min(100%, 320px), 1fr));
    gap: var(--space-lg);
  }

  .loading-state {
    grid-column: 1 / -1;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    padding: var(--space-3xl) var(--space-md);
  }

  @media (min-width: 1600px) {
    .channels-grid {
      grid-template-columns: repeat(4, 1fr);
    }
  }

  @media (min-width: 1200px) and (max-width: 1599px) {
    .channels-grid {
      grid-template-columns: repeat(3, 1fr);
    }
  }

  @media (min-width: 768px) and (max-width: 1199px) {
    .channels-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (max-width: 767px) {
    .channels-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 768px) {
    .shortcut-tips {
      flex-direction: column;
      gap: var(--space-sm);
    }
  }

  .empty-content {
    grid-column: 1 / -1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 384px;
    padding: var(--space-2xl);
    background: rgba(var(--color-background-800-rgb), 0.3);
    border-radius: var(--radius-lg);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    animation: fadeInUp var(--transition-slow) var(--ease-out);
  }

  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .empty-icon-wrapper {
    position: relative;
    margin-bottom: var(--space-xl);
  }

  .empty-icon-wrapper i {
    font-size: var(--text-4xl);
    background: linear-gradient(45deg, var(--color-brand-500), var(--color-error-500));
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    animation: float 3s var(--ease-in-out) infinite;
  }

  .icon-glow {
    position: absolute;
    inset: 0;
    z-index: -10;
    background: radial-gradient(circle, rgba(var(--color-brand-500-rgb), 0.3) 0%, transparent 70%);
    border-radius: var(--radius-full);
    animation: pulse 2s var(--ease-in-out) infinite;
  }

  .empty-title {
    margin-bottom: var(--space-sm);
    font-size: var(--text-xl);
    font-weight: var(--font-bold);
    color: var(--color-content-default);
  }

  .empty-desc {
    margin-bottom: var(--space-xl);
    font-size: var(--text-sm);
    color: var(--color-content-muted);
  }

  .empty-actions {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-md);
  }

  .empty-tip {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    font-size: var(--text-xs);
    color: var(--color-content-subtle);
  }

  .empty-tip i {
    color: var(--color-warning-400);
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0px);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 0.3;
    }
    50% {
      opacity: 0.6;
    }
  }

  .action-buttons {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }

  .button-divider {
    height: var(--space-xl);
    width: 1px;
    background: rgba(var(--color-brand-500-rgb), 0.2);
    opacity: var(--opacity-disabled);
  }

  @media (max-width: 768px) {
    .button-divider {
      display: none;
    }

    .action-buttons {
      flex-wrap: wrap;
      gap: var(--space-xs);
    }
  }

  :deep(.form-container) {
    background: rgba(var(--color-background-800-rgb), 0.4);
    border-radius: var(--radius-lg);
    padding: var(--space-lg);
    margin-bottom: var(--space-lg);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    box-shadow: var(--shadow-cyber-md);
    backdrop-filter: var(--backdrop-blur-md);
  }
</style>
