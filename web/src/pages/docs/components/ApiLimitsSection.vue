<script lang="ts">
  import { computed, defineComponent, type PropType } from 'vue'
  import ApiSection from './ApiSection/index.vue'
  import { useTexts } from '@/composables/useTexts'

  export default defineComponent({
    name: 'ApiLimitsSection',
    components: {
      ApiSection,
    },
    props: {
      supportedFormats: {
        type: Array as PropType<string[]>,
        required: true,
      },
    },
    setup(props) {
      const { $t } = useTexts()

      const fileFormatText = computed(() => {
        const formatsStr = Array.isArray(props.supportedFormats) ? props.supportedFormats.join(', ') : ''
        return $t('docs.limits.fileFormat', { formats: formatsStr })
      })

      return { $t, fileFormatText }
    },
  })
</script>

<template>
  <ApiSection id="api-limits" :title="$t('docs.limits.title')" icon-class="fas fa-exclamation-circle">
    <div class="limits-grid">
      <div class="limit-card">
        <h4><i class="fas fa-file-image" />{{ $t('docs.limits.fileLimit') }}</h4>
        <ul>
          <li>{{ fileFormatText }}</li>
          <li>{{ $t('docs.limits.singleFileSize') }}</li>
          <li>{{ $t('docs.limits.batchUpload') }}</li>
          <li>{{ $t('docs.limits.autoOptimize') }}</li>
        </ul>
      </div>

      <div class="limit-card">
        <h4><i class="fas fa-key" />{{ $t('docs.limits.apiKeyLimit') }}</h4>
        <ul>
          <li>{{ $t('docs.limits.storageCapacity') }}</li>
          <li>{{ $t('docs.limits.uploadCount') }}</li>
          <li>{{ $t('docs.limits.singleFileSizeLimit') }}</li>
          <li>{{ $t('docs.limits.defaultFolder') }}</li>
        </ul>
      </div>
    </div>

    <div class="error-codes">
      <h4><i class="fas fa-exclamation-triangle" />{{ $t('docs.limits.commonErrors') }}</h4>
      <div class="error-list">
        <div class="error-item">
          <span class="error-code">102</span>
          <div class="error-details">
            <span class="error-title">{{ $t('docs.limits.error102Title') }}</span>
            <span class="error-solution">{{ $t('docs.limits.error102Solution') }}</span>
          </div>
        </div>
        <div class="error-item">
          <span class="error-code">4000</span>
          <div class="error-details">
            <span class="error-title">{{ $t('docs.limits.error4000Title') }}</span>
            <span class="error-solution">{{ $t('docs.limits.error4000Solution') }}</span>
          </div>
        </div>
        <div class="error-item">
          <span class="error-code">4001</span>
          <div class="error-details">
            <span class="error-title">{{ $t('docs.limits.error4001Title') }}</span>
            <span class="error-solution">{{ $t('docs.limits.error4001Solution') }}</span>
          </div>
        </div>
        <div class="error-item">
          <span class="error-code">4008</span>
          <div class="error-details">
            <span class="error-title">{{ $t('docs.limits.error4008Title') }}</span>
            <span class="error-solution">{{ $t('docs.limits.error4008Solution') }}</span>
          </div>
        </div>
        <div class="error-item">
          <span class="error-code">4009</span>
          <div class="error-details">
            <span class="error-title">{{ $t('docs.limits.error4009Title') }}</span>
            <span class="error-solution">{{ $t('docs.limits.error4009Solution') }}</span>
          </div>
        </div>
        <div class="error-item">
          <span class="error-code">4010</span>
          <div class="error-details">
            <span class="error-title">{{ $t('docs.limits.error4010Title') }}</span>
            <span class="error-solution">{{ $t('docs.limits.error4010Solution') }}</span>
          </div>
        </div>
      </div>
    </div>
  </ApiSection>
</template>

<style scoped>
  .limits-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .limit-card {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    transition: all 0.2s ease;
  }

  .limit-card:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-2px);
  }

  .limit-card h4 {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 0 0 0.75rem 0;
    font-size: 1rem;
    color: var(--color-text-content);
  }

  .limit-card i {
    color: var(--color-brand-500);
  }

  .limit-card ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .limit-card li {
    padding: 0.25rem 0;
    color: var(--color-gray-200);
    font-size: 0.875rem;
    position: relative;
    padding-left: 1rem;
  }

  .error-codes {
    background: rgba(var(--color-badge-accent-text-rgb), 0.05);
    border: 1px solid rgba(var(--color-badge-accent-text-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
  }

  .error-codes h4 {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 0 0 0.75rem 0;
    font-size: 1rem;
    color: var(--color-error-500);
  }

  .error-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .error-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem;
    background: rgba(var(--color-background-900-rgb), 0.3);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
  }

  .error-item:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.25);
    background: rgba(var(--color-background-900-rgb), 0.5);
  }

  .error-code {
    color: var(--color-error-500);
    font-weight: 700;
    font-family: 'SF Mono', Monaco, monospace;
    font-size: 0.875rem;
    background: rgba(var(--color-badge-accent-text-rgb), 0.1);
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-badge-accent-text-rgb), 0.2);
    flex-shrink: 0;
    min-width: 50px;
    text-align: center;
  }

  .error-details {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    flex: 1;
  }

  .error-title {
    font-size: 0.875rem;
    color: var(--color-text-content);
    font-weight: 600;
  }

  .error-solution {
    font-size: 0.8125rem;
    color: var(--color-gray-200);
    opacity: 0.8;
  }

  @media (max-width: 768px) {
    .error-item {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.5rem;
    }

    .error-code {
      align-self: flex-start;
    }

    .error-details {
      width: 100%;
    }
  }
</style>
