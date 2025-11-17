<script lang="ts">
  import { defineComponent } from 'vue'
  import ApiSection from './ApiSection/index.vue'
  import CodeExample from './CodeExample/index.vue'
  import { useTexts } from '@/composables/useTexts'

  export default defineComponent({
    name: 'UploadApiSection',
    components: {
      ApiSection,
      CodeExample,
    },
    props: {
      supportedFormats: {
        type: Array,
        required: true,
      },
      jsonResponseExamples: {
        type: Object,
        required: true,
      },
    },
    setup() {
      const { $t } = useTexts()
      return { $t }
    },
  })
</script>

<template>
  <ApiSection id="api-upload" :title="$t('docs.upload.title')" icon-class="fas fa-upload">
    <div class="api-endpoint">
      <div class="endpoint-header">
        <span class="method post">POST</span>
        <code>/api/v1/external/upload</code>
      </div>

      <div class="endpoint-content">
        <h4>{{ $t('docs.upload.requestParams') }}</h4>
        <div class="params-table">
          <div class="param-row header">
            <span>{{ $t('docs.upload.paramName') }}</span>
            <span>{{ $t('docs.upload.paramType') }}</span>
            <span>{{ $t('docs.upload.paramRequired') }}</span>
            <span>{{ $t('docs.upload.paramDesc') }}</span>
          </div>
          <div class="param-row">
            <code>file</code>
            <span>File</span>
            <span class="required">{{ $t('docs.upload.yes') }}*</span>
            <span>{{ $t('docs.upload.fileDesc') }}</span>
          </div>
          <div class="param-row">
            <code>files[]</code>
            <span>File[]</span>
            <span class="required">{{ $t('docs.upload.yes') }}*</span>
            <span>{{ $t('docs.upload.filesDesc') }}</span>
          </div>
          <div class="param-row">
            <code>folderId</code>
            <span>String</span>
            <span>{{ $t('docs.upload.no') }}</span>
            <span>{{ $t('docs.upload.folderIdDesc') }}</span>
          </div>
          <div class="param-row">
            <code>filePath</code>
            <span>String</span>
            <span>{{ $t('docs.upload.no') }}</span>
            <span>{{ $t('docs.upload.filePathDesc') }}</span>
          </div>
          <div class="param-row">
            <code>access_level</code>
            <span>String</span>
            <span>{{ $t('docs.upload.no') }}</span>
            <span>{{ $t('docs.upload.accessLevelDesc') }}</span>
          </div>
          <div class="param-row">
            <code>optimize</code>
            <span>String</span>
            <span>{{ $t('docs.upload.no') }}</span>
            <span>{{ $t('docs.upload.optimizeDesc') }}</span>
          </div>
        </div>
        <p class="param-note">
          {{ $t('docs.upload.paramNote') }}
        </p>

        <h4>{{ $t('docs.upload.supportedFormats') }}</h4>
        <div class="format-tags">
          <span v-for="format in supportedFormats" :key="format" class="format-tag">{{ format }}</span>
        </div>

        <h4>{{ $t('docs.upload.successResponse') }}</h4>
        <CodeExample :title="$t('docs.upload.singleSuccess')" language="json">{{ jsonResponseExamples.single }}</CodeExample>
        <CodeExample :title="$t('docs.upload.batchSuccess')" language="json">{{ jsonResponseExamples.batch }}</CodeExample>
      </div>
    </div>
  </ApiSection>
</template>

<style scoped>
  .api-endpoint {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    overflow: hidden;
    margin-bottom: 1.5rem;
  }

  .endpoint-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.25rem;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
  }

  .method {
    padding: 0.375rem 0.75rem;
    border-radius: var(--radius-sm);
    font-weight: 700;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .method.post {
    background: linear-gradient(135deg, var(--color-warning-500) 0%, rgba(var(--color-warning-rgb), 0.9) 100%);
    color: var(--color-text-on-warning);
    box-shadow: 0 2px 8px rgba(var(--color-warning-rgb), 0.4);
    border: 1px solid rgba(var(--color-warning-rgb), 0.6);
  }

  .endpoint-header code {
    font-family: 'SF Mono', Monaco, monospace;
    font-size: 1rem;
    color: var(--color-brand-500);
    background: rgba(var(--color-brand-500-rgb), 0.2);
    padding: 0.375rem 0.5rem;
    border-radius: var(--radius-sm);
  }

  .endpoint-content {
    padding: 1.5rem;
  }

  .endpoint-content h4 {
    font-size: 1rem;
    color: var(--color-brand-500);
    margin: 0 0 0.75rem 0;
    font-weight: 600;
  }

  .params-table {
    background: rgba(var(--color-background-900-rgb), 0.5);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    overflow: hidden;
    margin-bottom: 0.75rem;
  }

  .param-row {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 2fr;
    gap: 0.75rem;
    padding: 0.75rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    font-size: 0.875rem;
  }

  .param-row:last-child {
    border-bottom: none;
  }

  .param-row.header {
    background: rgba(var(--color-brand-500-rgb), 0.08);
    font-weight: 600;
    color: var(--color-brand-500);
  }

  .param-row code {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    padding: 0.25rem 0.375rem;
    border-radius: var(--radius-sm);
  }

  .required {
    color: var(--color-error-500);
    font-weight: 600;
  }

  .param-note {
    font-size: 0.8125rem;
    color: var(--color-slate-500);
    font-style: italic;
  }

  .format-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
  }

  .format-tag {
    background: rgba(var(--color-brand-500-rgb), 0.3);
    color: var(--color-brand-500);
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
  }
</style>
