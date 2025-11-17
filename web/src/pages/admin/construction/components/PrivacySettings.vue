<script setup lang="ts">
  import { onMounted, reactive, watch } from 'vue'
  import { MdEditor } from 'md-editor-v3'
  import 'md-editor-v3/lib/style.css'
  import type { Setting } from '@/api/admin/settings'
  import { useTexts } from '@/composables/useTexts'
  import { uploadAdminFile } from '@/api/admin/file'
  import { useToast } from '@/components/Toast/useToast'

  const { $t } = useTexts()
  const toast = useToast()

  const props = defineProps<{
    settings: Setting[]
  }>()

  const emit = defineEmits<{
    (e: 'update', settings: Setting[]): void
  }>()

  /* 本地设置对象 */
  const localSettings = reactive({
    privacy_policy_content: '',
  })

  async function handleUploadImg(files: File[], callback: (urls: string[]) => void) {
    try {
      const uploadPromises = files.map(async (file) => {
        const res = await uploadAdminFile(file)
        if (res.code === 200 && res.data) {
          return res.data.full_url
        } else {
          throw new Error(res.message || $t('admin.construction.privacy.uploadFailed'))
        }
      })

      const urls = await Promise.all(uploadPromises)
      callback(urls)
      toast.success($t('admin.construction.privacy.uploadSuccess').replace('{count}', urls.length.toString()))
    } catch (error: any) {
      toast.error($t('admin.construction.privacy.uploadError') + error.message)
      callback([])
    }
  }

  const getSettingsArray = (): Setting[] => {
    return [
      {
        key: 'privacy_policy_content',
        value: localSettings.privacy_policy_content,
        type: 'text',
        group: 'legal',
        description: $t('admin.construction.privacy.title'),
      },
    ]
  }

  const applySettings = (settings: Setting[]) => {
    settings.forEach((setting) => {
      const key = setting.key as keyof typeof localSettings
      if (key in localSettings) {
        const settingsRecord = localSettings as Record<string, unknown>
        settingsRecord[key] = setting.value
      }
    })
  }

  watch(
    () => props.settings,
    (newSettings) => {
      if (newSettings && newSettings.length > 0) {
        applySettings(newSettings)
      }
    },
    { deep: true, immediate: true }
  )

  watch(
    localSettings,
    () => {
      emit('update', getSettingsArray())
    },
    { deep: true }
  )

  onMounted(() => {
    if (!props.settings || props.settings.length === 0) {
      emit('update', getSettingsArray())
    }
  })
</script>

<template>
  <div class="privacy-settings">
    <div class="settings-section">
      <div class="section-header">
        <div class="section-title">
          <i class="fas fa-shield-alt" />
          <h3>{{ $t('admin.construction.privacy.title') }}</h3>
        </div>
        <p class="section-description">{{ $t('admin.construction.privacy.description') }}</p>
      </div>

      <div class="editor-container">
        <MdEditor
          v-model="localSettings.privacy_policy_content"
          language="zh-CN"
          :placeholder="$t('admin.construction.privacy.placeholder')"
          preview-theme="github"
          code-theme="github"
          :toolbars-exclude="['github', 'fullscreen', 'pageFullscreen']"
          :on-upload-img="handleUploadImg"
          :show-code-row-number="true"
          :auto-focus="false"
          :tab-width="2"
          style="height: 100%"
        />
      </div>

      <div class="editor-hint">
        <i class="fas fa-info-circle" />
        <span>{{ $t('admin.construction.privacy.hint') }}</span>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .privacy-settings {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
    height: 100%;
  }

  .settings-section {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    flex: 1;
    min-height: 0;
  }

  .section-header {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    color: var(--color-content-heading);

    i {
      font-size: var(--text-lg);
      color: var(--color-brand-500);
      filter: drop-shadow(0 0 4px rgba(var(--color-brand-500-rgb), 0.6));
    }

    h3 {
      font-size: var(--text-lg);
      font-weight: var(--font-semibold);
      margin: 0;
    }
  }

  .section-description {
    color: var(--color-content-muted);
    font-size: var(--text-sm);
    margin: 0;
    padding-left: calc(var(--text-lg) + var(--space-sm));
  }

  .editor-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .editor-hint {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    padding: var(--space-sm) var(--space-md);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-md);
    color: var(--color-content-muted);
    font-size: var(--text-sm);

    i {
      color: var(--color-brand-400);
      font-size: var(--text-base);
    }
  }
</style>
