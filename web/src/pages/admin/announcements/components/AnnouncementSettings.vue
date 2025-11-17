<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import type { AnnouncementSettings, AnnouncementSettingsUpdate } from '@/api/types/announcement'
  import { getAnnouncementSettings, updateAnnouncementSettings } from '@/api/admin/announcement'

  defineOptions({
    name: 'AnnouncementSettings',
  })

  const { $t } = useTexts()
  const toast = useToast()
  const loading = ref(false)
  const saving = ref(false)

  /* 配置项 */
  const settings = ref<AnnouncementSettings>({
    announcement_enabled: true,
    announcement_drawer_position: 'right',
    announcement_drawer_width: '450px',
    announcement_display_limit: 10,
    announcement_auto_show_delay: 3,
  })

  const fetchSettings = async () => {
    loading.value = true
    try {
      const result = await getAnnouncementSettings()
      if (result.code === 200 && result.data) {
        let drawerWidth = result.data.announcement_drawer_width ?? 600
        if (typeof drawerWidth === 'number') {
          drawerWidth = `${drawerWidth}px`
        }

        settings.value = {
          announcement_enabled: result.data.announcement_enabled ?? true,
          announcement_drawer_position: result.data.announcement_drawer_position ?? 'right',
          announcement_drawer_width: drawerWidth,
          announcement_display_limit: result.data.announcement_display_limit ?? 10,
          announcement_auto_show_delay: result.data.announcement_auto_show_delay ?? 3,
        }
      }
    } catch (error: any) {
      console.error($t('admin.announcements.settings.messages.fetchError'), error)
    } finally {
      loading.value = false
    }
  }

  const saveSettings = async () => {
    if (settings.value.announcement_display_limit < 1 || settings.value.announcement_display_limit > 20) {
      toast.error($t('admin.announcements.settings.messages.limitRange'))
      return
    }

    saving.value = true
    try {
      const updateData: AnnouncementSettingsUpdate = {
        announcement_enabled: settings.value.announcement_enabled,
        announcement_drawer_position: settings.value.announcement_drawer_position,
        announcement_drawer_width: '450px', // 固定宽度
        announcement_display_limit: settings.value.announcement_display_limit,
        announcement_auto_show_delay: settings.value.announcement_auto_show_delay,
      }

      const result = await updateAnnouncementSettings(updateData)
      if (result.code === 200) {
        toast.success($t('admin.announcements.settings.messages.saveSuccess'))
      } else {
        toast.error($t('admin.announcements.settings.messages.saveError') + result.message)
      }
    } catch (error: any) {
      toast.error($t('admin.announcements.settings.messages.saveError') + error.message)
    } finally {
      saving.value = false
    }
  }

  onMounted(() => {
    fetchSettings()
  })
</script>

<template>
  <div class="announcement-settings-card">
    <div v-if="loading" class="card-loading">
      <CyberSkeleton type="text" :count="1" />
    </div>

    <div v-else class="card-content">
      <div class="setting-item">
        <label class="setting-label">{{ $t('admin.announcements.settings.systemToggle') }}</label>
        <CyberSwitch v-model="settings.announcement_enabled" />
      </div>

      <div class="setting-item">
        <label class="setting-label">{{ $t('admin.announcements.settings.drawerPosition.label') }}</label>
        <CyberRadioGroup
          v-model="settings.announcement_drawer_position"
          :options="[
            { label: $t('admin.announcements.settings.drawerPosition.left'), value: 'left' },
            { label: $t('admin.announcements.settings.drawerPosition.right'), value: 'right' },
          ]"
          layout="horizontal"
        />
      </div>

      <div class="setting-item">
        <label class="setting-label">{{ $t('admin.announcements.settings.displayLimit') }}</label>
        <div class="slider-wrapper">
          <CyberSlider v-model="settings.announcement_display_limit" :min="1" :max="20" :step="1" width="150px" />
          <span class="slider-value">{{ settings.announcement_display_limit }}</span>
        </div>
      </div>

      <div class="setting-item">
        <label class="setting-label">{{ $t('admin.announcements.settings.autoShowDelay') }}</label>
        <div class="slider-wrapper">
          <CyberSlider v-model="settings.announcement_auto_show_delay" :min="0" :max="10" :step="1" width="150px" />
          <span class="slider-value"
            >{{ settings.announcement_auto_show_delay }}{{ $t('admin.announcements.settings.units.second') }}</span
          >
        </div>
      </div>

      <div class="setting-item save-button">
        <CyberButton type="primary" size="small" :loading="saving" @click="saveSettings">
          <i class="fas fa-save mr-1.5" />
          {{ $t('admin.announcements.settings.saveButton') }}
        </CyberButton>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .announcement-settings-card {
    background: rgba(var(--color-background-800-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-lg);
    margin-bottom: var(--space-md);
    padding: var(--space-md) var(--space-lg);
  }

  .card-loading {
    padding: var(--space-xs);
  }

  .card-content {
    display: flex;
    align-items: center;
    gap: var(--space-xl);
    flex-wrap: wrap;
  }

  .setting-item {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }

  .setting-label {
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--color-content-muted);
    white-space: nowrap;
  }

  .unit-hint {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
  }

  .save-button {
    margin-left: auto;
  }

  .slider-wrapper {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }

  .slider-value {
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--color-brand-500);
    min-width: 30px;
    text-align: center;
  }

  @media (max-width: 1200px) {
    .card-content {
      gap: var(--space-md);
    }

    .setting-item {
      flex: 1 1 auto;
      min-width: fit-content;
    }
  }

  @media (max-width: 768px) {
    .card-content {
      flex-direction: column;
      align-items: stretch;
    }

    .setting-item {
      width: 100%;
      justify-content: space-between;
    }

    .input-with-unit {
      width: auto;
      flex: 1;
    }
  }
</style>
