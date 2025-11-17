<script setup lang="ts">
  import { ref, watch } from 'vue'
  import { MdEditor } from 'md-editor-v3'
  import 'md-editor-v3/lib/style.css'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import type { Announcement, CreateAnnouncementRequest } from '@/api/types/announcement'
  import { createAnnouncement, updateAnnouncement } from '@/api/admin/announcement'
  import { uploadAdminFile } from '@/api/admin/file'

  interface Props {
    modelValue: Announcement | null
    isCreating: boolean
    loading?: boolean
  }

  interface Emits {
    (e: 'submit'): void
    (e: 'update:modelValue', value: Announcement | null): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()
  const { $t } = useTexts()
  const toast = useToast()

  const form = ref<CreateAnnouncementRequest>({
    title: '',
    content: '',
    summary: '',
    is_pinned: false,
    status: 'draft',
  })

  function initForm() {
    if (props.modelValue) {
      form.value = {
        title: props.modelValue.title,
        content: props.modelValue.content,
        summary: props.modelValue.summary,
        is_pinned: props.modelValue.is_pinned,
        status: props.modelValue.status,
      }
    } else {
      form.value = {
        title: '',
        content: '',
        summary: '',
        is_pinned: false,
        status: 'draft',
      }
    }
  }

  /* 统一使用 Validator 进行校验 */
  import Validator, { getValidationRules } from '@/utils/validation/validator'
  const R = getValidationRules($t)
  function validateForm() {
    const titleCheck = Validator.validate(form.value.title || '', [R.required, R.maxLength(255)], $t)
    if (!titleCheck.valid) {
      toast.error(titleCheck.message || $t('utils.validation.errors.required'))
      return false
    }
    const contentCheck = Validator.validate(form.value.content || '', [R.required], $t)
    if (!contentCheck.valid) {
      toast.error(contentCheck.message || $t('utils.validation.errors.required'))
      return false
    }
    // 非必填：摘要长度限制
    if ((form.value.summary || '').length > 500) {
      const msg = R.maxLength(500).message
      toast.error(msg || $t('utils.validation.errors.maxLength', { n: '500' }))
      return false
    }
    return true
  }

  async function submitForm() {
    if (!validateForm()) {
      return
    }

    try {
      if (!props.isCreating && props.modelValue) {
        const res = await updateAnnouncement(props.modelValue.id, form.value)
        if (res.code === 200) {
          toast.success($t('admin.announcements.messages.updateSuccess'))
          emit('submit')
        } else {
          toast.error($t('admin.announcements.messages.updateError') + res.message)
        }
      } else {
        const res = await createAnnouncement(form.value)
        if (res.code === 200) {
          toast.success($t('admin.announcements.messages.createSuccess'))
          emit('submit')
        } else {
          toast.error($t('admin.announcements.messages.createError') + res.message)
        }
      }
    } catch (error: any) {
      toast.error($t('admin.announcements.messages.operationError') + error.message)
    }
  }

  async function handleUploadImg(files: File[], callback: (urls: string[]) => void) {
    try {
      const uploadPromises = files.map(async (file) => {
        const res = await uploadAdminFile(file)
        if (res.code === 200 && res.data) {
          return res.data.full_url
        } else {
          throw new Error(res.message || $t('admin.announcements.messages.uploadFailed'))
        }
      })

      const urls = await Promise.all(uploadPromises)
      callback(urls)
      toast.success($t('admin.announcements.messages.uploadSuccess').replace('{count}', urls.length.toString()))
    } catch (error: any) {
      toast.error($t('admin.announcements.messages.uploadError') + error.message)
      callback([])
    }
  }

  defineExpose({
    submitForm,
    validateForm,
  })

  watch(
    () => props.modelValue,
    () => {
      initForm()
    },
    { immediate: true, deep: true }
  )
</script>

<template>
  <div class="announcement-form-content">
    <div class="form-section">
      <div class="form-item">
        <label class="form-label required">{{ $t('admin.announcements.form.title.label') }}</label>
        <CyberInput v-model="form.title" :placeholder="$t('admin.announcements.form.title.placeholder')" :maxlength="255" />
      </div>

      <div class="form-item">
        <label class="form-label">{{ $t('admin.announcements.form.summary.label') }}</label>
        <CyberInput
          v-model="form.summary"
          type="textarea"
          :placeholder="$t('admin.announcements.form.summary.placeholder')"
          :rows="3"
          :maxlength="500"
        />
      </div>
    </div>

    <div class="form-section">
      <div class="form-item">
        <label class="form-label required">{{ $t('admin.announcements.form.content.label') }}</label>
        <div class="markdown-editor-wrapper">
          <MdEditor
            v-model="form.content"
            language="zh-CN"
            :placeholder="$t('admin.announcements.form.content.placeholder')"
            preview-theme="github"
            code-theme="github"
            :toolbars-exclude="['github', 'pageFullscreen']"
            :on-upload-img="handleUploadImg"
            :show-code-row-number="true"
            :auto-focus="false"
            :tab-width="2"
            style="height: 600px"
          />
        </div>
        <div class="form-hint">
          {{ $t('admin.announcements.form.content.hint') }}
        </div>
      </div>
    </div>

    <div class="form-section">
      <div class="form-row">
        <div class="form-item flex-1">
          <label class="form-label required">{{ $t('admin.announcements.form.status.label') }}</label>
          <CyberRadioGroup
            v-model="form.status"
            :options="[
              { label: $t('admin.announcements.status.draft'), value: 'draft' },
              { label: $t('admin.announcements.status.published'), value: 'published' },
              { label: $t('admin.announcements.status.archived'), value: 'archived' },
            ]"
            layout="horizontal"
          />
        </div>

        <div class="form-item flex-1">
          <label class="form-label">{{ $t('admin.announcements.form.options.label') }}</label>
          <div class="checkbox-inline-hint">
            <CyberCheckbox v-model="form.is_pinned">
              {{ $t('admin.announcements.form.options.pinned') }}
            </CyberCheckbox>
            <span class="inline-hint">
              <i class="fas fa-info-circle" />
              {{ $t('admin.announcements.form.options.pinnedHint') }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .announcement-form-content {
  }

  .form-section {
    margin-bottom: var(--space-lg);

    &:last-child {
      margin-bottom: 0;
    }
  }

  .form-row {
    display: flex;
    gap: var(--space-md);
    align-items: flex-start;
  }

  .flex-1 {
    flex: 1;
  }

  .form-item {
    margin-bottom: var(--space-md);

    &:last-child {
      margin-bottom: 0;
    }
  }

  .form-label {
    display: block;
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    color: var(--color-content-muted);
    margin-bottom: var(--space-xs);
  }

  .form-label.required::after {
    content: '*';
    color: var(--color-error-400);
    margin-left: var(--space-xs);
  }

  .form-hint {
    font-size: var(--text-xs);
    color: var(--color-content-muted);
    margin-top: var(--space-xs);
  }

  .checkbox-with-hint {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
  }

  .checkbox-inline-hint {
    display: flex;
    align-items: center;
    gap: var(--space-md);
  }

  .inline-hint {
    display: inline-flex;
    align-items: center;
    gap: var(--space-xs);
    font-size: var(--text-xs);
    color: var(--color-content-muted);

    i {
      color: var(--color-brand-500);
      font-size: var(--text-xs);
    }
  }

  .hint-with-icon {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    padding: var(--space-xs) var(--space-sm);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border-left: 2px solid var(--color-brand-500);
    border-radius: var(--radius-sm);
  }

  .hint-icon {
    color: var(--color-brand-500);
    font-size: var(--text-xs);
    flex-shrink: 0;
  }

  .markdown-editor-wrapper {
  }
</style>
