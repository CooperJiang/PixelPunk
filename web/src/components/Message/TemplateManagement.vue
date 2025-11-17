<script setup lang="ts">
  import { onMounted, ref } from 'vue'
  import { deleteMessageTemplate, getMessageTemplates, toggleMessageTemplate, formatRelativeTime } from '@/api/message'
  import { getMessageTypeConfig, type MessageTemplate } from '@/api/message/types'
  import TemplateModal from './TemplateModal.vue'
  import { showConfirm } from '@/utils/dialog'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const messageTypeConfig = getMessageTypeConfig($t)

  const loading = ref(false)
  const templates = ref<MessageTemplate[]>([])
  const showModal = ref(false)
  const currentTemplate = ref<MessageTemplate | null>(null)
  const isEdit = ref(false)

  const loadTemplates = async () => {
    loading.value = true
    try {
      const response = await getMessageTemplates()
      templates.value = response.templates
    } catch {
    } finally {
      loading.value = false
    }
  }

  const openCreateModal = () => {
    currentTemplate.value = null
    isEdit.value = false
    showModal.value = true
  }

  const editTemplate = (template: MessageTemplate) => {
    currentTemplate.value = { ...template }
    isEdit.value = true
    showModal.value = true
  }

  const closeModal = () => {
    showModal.value = false
    currentTemplate.value = null
    isEdit.value = false
  }

  const handleTemplateSaved = () => {
    closeModal()
    loadTemplates()
  }

  const toggleTemplate = async (template: MessageTemplate) => {
    try {
      await toggleMessageTemplate(template.id.toString())
      template.is_enabled = !template.is_enabled
      template.updated_at = new Date().toISOString()
    } catch {}
  }

  const deleteTemplate = async (template: MessageTemplate) => {
    if (!showConfirm($t('components.templateManagement.confirmDelete', { name: template.name }))) {
      return
    }

    try {
      await deleteMessageTemplate(template.id.toString())
      const index = templates.value.findIndex((t) => t.id === template.id)
      if (index !== -1) {
        templates.value.splice(index, 1)
      }
    } catch {}
  }

  const getMessageTypeLabel = (type: string): string => messageTypeConfig[type as keyof typeof messageTypeConfig]?.label || type

  onMounted(() => {
    loadTemplates()
  })
</script>

<template>
  <div class="cyber-template-management">
    <div class="cyber-header">
      <div class="cyber-header-content">
        <div class="cyber-title-section">
          <h1 class="cyber-main-title">{{ $t('components.templateManagement.title') }}</h1>
          <p class="cyber-subtitle">{{ $t('components.templateManagement.subtitle') }}</p>
        </div>
        <div class="cyber-actions">
          <CyberButton type="primary" icon="plus" @click="openCreateModal">
            {{ $t('components.templateManagement.create') }}
          </CyberButton>
        </div>
      </div>
    </div>

    <div class="cyber-content">
      <CyberLoading v-if="loading" />

      <div v-else-if="templates.length === 0" class="cyber-empty">
        <i class="fas fa-file-alt mb-4 text-6xl text-content-disabled" />
        <h3 class="cyber-empty-title">{{ $t('components.templateManagement.empty.title') }}</h3>
        <p class="cyber-empty-desc">{{ $t('components.templateManagement.empty.description') }}</p>
      </div>

      <div v-else class="cyber-template-grid">
        <div
          v-for="template in templates"
          :key="template.id"
          class="cyber-template-card"
          :class="{ 'is-active': template.is_enabled, 'is-inactive': !template.is_enabled }"
        >
          <div class="cyber-card-header">
            <div class="cyber-template-info">
              <div class="cyber-template-meta">
                <h3 class="cyber-template-name">{{ template.name }}</h3>
                <div class="cyber-template-type">
                  <span class="cyber-type-badge">
                    {{ getMessageTypeLabel(template.type) }}
                  </span>
                  <span class="cyber-status-badge" :class="template.is_enabled ? 'enabled' : 'disabled'">
                    <i :class="template.is_enabled ? 'fas fa-check-circle' : 'fas fa-times-circle'" class="mr-1" />
                    {{
                      template.is_enabled
                        ? $t('components.templateManagement.enabled')
                        : $t('components.templateManagement.disabled')
                    }}
                  </span>
                </div>
              </div>
            </div>

            <div class="cyber-card-actions">
              <CyberButton type="ghost" icon="edit" size="small" @click="editTemplate(template)">
                {{ $t('components.templateManagement.edit') }}
              </CyberButton>

              <CyberButton
                :type="template.is_enabled ? 'outlined' : 'secondary'"
                :icon="template.is_enabled ? 'pause' : 'play'"
                size="small"
                @click="toggleTemplate(template)"
              >
                {{
                  template.is_enabled ? $t('components.templateManagement.disable') : $t('components.templateManagement.enable')
                }}
              </CyberButton>

              <CyberButton type="danger" icon="trash" size="small" @click="deleteTemplate(template)">
                {{ $t('components.templateManagement.delete') }}
              </CyberButton>
            </div>
          </div>

          <div class="cyber-card-body">
            <div class="cyber-template-content">
              <h4 class="cyber-content-title">{{ template.title }}</h4>
              <p class="cyber-content-text">{{ template.content }}</p>
            </div>

            <div class="cyber-template-config">
              <div class="cyber-config-row">
                <span class="cyber-config-item">
                  <i
                    :class="template.send_email ? 'fas fa-envelope text-cyan-400' : 'fas fa-ban text-content-disabled'"
                    class="mr-1"
                  />
                  {{
                    template.send_email
                      ? $t('components.templateManagement.emailEnabled')
                      : $t('components.templateManagement.emailDisabled')
                  }}
                </span>
                <span class="cyber-config-item">
                  <i
                    :class="template.show_toast ? 'fas fa-bell text-green-400' : 'fas fa-ban text-content-disabled'"
                    class="mr-1"
                  />
                  {{
                    template.show_toast
                      ? $t('components.templateManagement.toastEnabled')
                      : $t('components.templateManagement.toastDisabled')
                  }}
                </span>
              </div>

              <div v-if="template.action_text" class="cyber-action-config">
                <span class="cyber-action-label">
                  <i class="fas fa-mouse-pointer mr-1" style="color: var(--color-brand-500)" />
                  {{ $t('components.templateManagement.actionButton') }}
                </span>
                <span class="cyber-action-text">{{ template.action_text }}</span>
              </div>
            </div>
          </div>

          <div class="cyber-card-footer">
            <div class="cyber-template-stats">
              <span class="cyber-stat">
                <i class="fas fa-clock mr-1" />
                {{ $t('components.templateManagement.updated') }}: {{ formatRelativeTime(template.updated_at, $t) }}
              </span>
              <span class="cyber-stat">
                <i class="fas fa-user mr-1" />
                {{ $t('components.templateManagement.creator') }}: {{ template.created_by || 'System' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <TemplateModal
      :visible="showModal"
      :template="currentTemplate"
      :is-edit="isEdit"
      @close="closeModal"
      @saved="handleTemplateSaved"
    />
  </div>
</template>

<style scoped>
  .cyber-template-management {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1rem;
  }

  .cyber-header {
    padding: 2rem 0;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    margin-bottom: 2rem;
  }

  .cyber-header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 2rem;
  }

  .cyber-title-section {
    flex: 1;
  }

  .cyber-main-title {
    font-size: 1.875rem;
    font-weight: 700;
    color: var(--color-brand-500);
    text-shadow: 0 0 15px rgba(var(--color-brand-500-rgb), 0.3);
    margin: 0 0 0.5rem 0;
  }

  .cyber-subtitle {
    font-size: 0.875rem;
    color: var(--color-slate-400);
    margin: 0;
    line-height: 1.5;
  }

  .cyber-actions {
    flex-shrink: 0;
  }

  .cyber-content {
    min-height: 400px;
  }

  .cyber-empty {
    text-align: center;
    padding: 4rem 2rem;
  }

  .cyber-empty-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin: 0 0 0.5rem 0;
  }

  .cyber-empty-desc {
    color: var(--color-text-quaternary);
    margin: 0;
    font-size: 0.875rem;
  }

  .cyber-template-grid {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: repeat(auto-fill, minmax(450px, 1fr));
  }

  .cyber-template-card {
    background: rgba(24, 29, 38, 0.8);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: var(--radius-sm);
    padding: 1.5rem;
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
  }

  .cyber-template-card:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(var(--color-background-900-rgb), 0.3);
  }

  .cyber-template-card.is-active {
    border-color: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-template-card.is-inactive {
    border-color: rgba(107, 114, 128, 0.3);
  }

  .cyber-card-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 1rem;
    gap: 1rem;
  }

  .cyber-template-info {
    flex: 1;
    min-width: 0;
  }

  .cyber-template-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-content-default);
    margin: 0 0 0.5rem 0;
    line-height: 1.3;
  }

  .cyber-template-type {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    align-items: center;
  }

  .cyber-type-badge {
    display: inline-flex;
    align-items: center;
    padding: 3px 10px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    font-weight: 500;
    background: rgba(var(--color-error-rgb), 0.1);
    border: 1px solid rgba(var(--color-error-rgb), 0.3);
    color: var(--color-error-500);
  }

  .cyber-status-badge {
    display: inline-flex;
    align-items: center;
    padding: 3px 10px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    font-weight: 500;
  }

  .cyber-status-badge.enabled {
    background: rgba(16, 185, 129, 0.1);
    border: 1px solid rgba(16, 185, 129, 0.3);
    color: var(--color-success);
  }

  .cyber-status-badge.disabled {
    background: rgba(107, 114, 128, 0.1);
    border: 1px solid rgba(107, 114, 128, 0.3);
    color: var(--color-text-quaternary);
  }

  .cyber-card-actions {
    display: flex;
    gap: 0.5rem;
    flex-shrink: 0;
    flex-wrap: wrap;
  }

  .cyber-card-body {
    margin-bottom: 1rem;
  }

  .cyber-template-content {
    margin-bottom: 1rem;
  }

  .cyber-content-title {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin: 0 0 0.5rem 0;
    line-height: 1.3;
  }

  .cyber-content-text {
    font-size: 0.8rem;
    color: var(--color-gray-300);
    line-height: 1.4;
    margin: 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .cyber-template-config {
    padding: 0.75rem;
    background: rgba(var(--color-brand-500-rgb), 0.03);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
  }

  .cyber-config-row {
    display: flex;
    gap: 1rem;
    margin-bottom: 0.5rem;
  }

  .cyber-config-item {
    display: flex;
    align-items: center;
    font-size: 0.75rem;
    color: var(--color-slate-400);
    font-weight: 500;
  }

  .cyber-action-config {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding-top: 0.5rem;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .cyber-action-label {
    font-size: 0.75rem;
    color: var(--color-slate-400);
    font-weight: 500;
    display: flex;
    align-items: center;
  }

  .cyber-action-text {
    font-size: 0.75rem;
    color: var(--color-gray-300);
    font-weight: 500;
    padding: 2px 8px;
    background: rgba(147, 51, 234, 0.1);
    border: 1px solid rgba(147, 51, 234, 0.2);
    border-radius: var(--radius-sm);
  }

  .cyber-card-footer {
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    padding-top: 1rem;
  }

  .cyber-template-stats {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .cyber-stat {
    display: flex;
    align-items: center;
    font-size: 0.75rem;
    color: var(--color-text-quaternary);
  }

  @media (max-width: 1024px) {
    .cyber-template-grid {
      grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
    }
  }

  @media (max-width: 768px) {
    .cyber-header-content {
      flex-direction: column;
      align-items: stretch;
      gap: 1rem;
    }

    .cyber-template-grid {
      grid-template-columns: 1fr;
    }

    .cyber-card-header {
      flex-direction: column;
      align-items: stretch;
      gap: 1rem;
    }

    .cyber-card-actions {
      justify-content: space-between;
    }

    .cyber-config-row {
      flex-direction: column;
      gap: 0.5rem;
    }

    .cyber-template-stats {
      flex-direction: column;
      gap: 0.5rem;
    }
  }

  @media (max-width: 480px) {
    .cyber-template-management {
      padding: 0 0.5rem;
    }

    .cyber-template-card {
      padding: 1rem;
    }

    .cyber-card-actions {
      flex-direction: column;
    }
  }
</style>
