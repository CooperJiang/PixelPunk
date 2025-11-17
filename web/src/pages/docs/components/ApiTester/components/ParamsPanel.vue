<script setup lang="ts">
  import { computed } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import FileUpload from './FileUpload.vue'
  import type { ApiTestParams } from '../types'

  defineOptions({
    name: 'ParamsPanel',
  })

  interface Props {
    apiKey: string
    selectedFiles: File[]
    params: ApiTestParams
    totalFileSize: number
    showValidation: boolean
    isDragOver: boolean
  }

  const props = defineProps<Props>()

  const emit = defineEmits<{
    'update:apiKey': [value: string]
    'update:selectedFiles': [files: File[]]
    'update:params': [params: ApiTestParams]
    'update:showValidation': [value: boolean]
    'update:isDragOver': [value: boolean]
  }>()

  const { $t } = useTexts()

  /* 为 FileUpload 组件创建本地响应式引用 */
  const localSelectedFiles = computed({
    get: () => props.selectedFiles,
    set: (value) => emit('update:selectedFiles', value),
  })

  const localShowValidation = computed({
    get: () => props.showValidation,
    set: (value) => emit('update:showValidation', value),
  })

  const localIsDragOver = computed({
    get: () => props.isDragOver,
    set: (value) => emit('update:isDragOver', value),
  })

  const updateParam = (key: keyof ApiTestParams, value: any) => {
    const newParams = { ...props.params, [key]: value }
    emit('update:params', newParams)
  }
</script>

<template>
  <div class="params-cyber-panel">
    <div class="cyber-section">
      <div class="section-header">
        <i class="fas fa-shield-alt"></i>
        <h4>{{ $t('docs.params.authConfig') }}</h4>
        <div class="section-line"></div>
      </div>
      <div class="auth-container">
        <div class="input-group">
          <label class="cyber-label">
            <span>{{ $t('docs.params.apiKey') }}</span>
            <span class="required-star">*</span>
          </label>
          <CyberInput
            :model-value="apiKey"
            @update:model-value="$emit('update:apiKey', $event)"
            type="password"
            :placeholder="$t('docs.tester.apiKeyPlaceholder')"
            :error="!apiKey.trim() && showValidation"
          />
        </div>
      </div>
    </div>

    <FileUpload
      v-model:selectedFiles="localSelectedFiles"
      v-model:showValidation="localShowValidation"
      v-model:isDragOver="localIsDragOver"
      :totalFileSize="totalFileSize"
    />

    <div class="cyber-section">
      <div class="section-header">
        <i class="fas fa-cogs"></i>
        <h4>{{ $t('docs.params.paramsConfig') }}</h4>
        <div class="section-line"></div>
      </div>
      <div class="params-grid">
        <div class="input-group">
          <label class="cyber-label">{{ $t('docs.params.folderId') }}</label>
          <CyberInput
            :model-value="params.folderId"
            @update:model-value="updateParam('folderId', $event)"
            :placeholder="$t('docs.tester.folderIdPlaceholder')"
          />
        </div>

        <div class="input-group">
          <label class="cyber-label">{{ $t('docs.params.filePath') }}</label>
          <CyberInput
            :model-value="params.filePath"
            @update:model-value="updateParam('filePath', $event)"
            :placeholder="$t('docs.params.filePathPlaceholder')"
          />
        </div>

        <div class="input-group">
          <label class="cyber-label">{{ $t('docs.params.accessLevel') }}</label>
          <CyberRadioGroup
            :model-value="params.accessLevel"
            @update:model-value="updateParam('accessLevel', $event)"
            :options="[
              { label: $t('docs.params.public'), value: 'public' },
              { label: $t('docs.params.private'), value: 'private' },
              { label: $t('docs.params.protected'), value: 'protected' },
            ]"
          />
        </div>

        <div class="input-group full-width">
          <CyberCheckbox :model-value="params.optimize" @update:model-value="updateParam('optimize', $event)">
            {{ $t('docs.params.enableOptimize') }}
          </CyberCheckbox>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .params-cyber-panel {
    padding: 1rem;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.4) 0%,
      rgba(var(--color-background-700-rgb), 0.6) 100%
    );
    border-right: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .cyber-section {
    margin-bottom: 1.5rem;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.75rem;
    padding-bottom: 0.4rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  .section-header i {
    color: var(--color-brand-500);
    font-size: 1rem;
  }

  .section-header h4 {
    color: var(--color-text-content);
    font-size: 1rem;
    font-weight: 600;
    margin: 0;
  }

  .section-line {
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(var(--color-brand-500-rgb), 0.3), transparent);
  }

  .input-group {
    margin-bottom: 1rem;
  }

  .cyber-label {
    display: flex;
    align-items: center;
    color: var(--color-gray-200);
    font-size: 0.875rem;
    font-weight: 500;
    margin-bottom: 0.5rem;
  }

  .required-star {
    color: var(--color-error-500);
    margin-left: 0.25rem;
  }

  .params-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .full-width {
    grid-column: 1 / -1;
  }

  @media (max-width: 768px) {
    .params-cyber-panel {
      border-right: none;
      border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    }

    .params-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
