<script lang="ts">
  import { defineComponent } from 'vue'
  import ApiSection from './ApiSection/index.vue'
  import { useTexts } from '@/composables/useTexts'

  export default defineComponent({
    name: 'AuthenticationSection',
    components: {
      ApiSection,
    },
    setup() {
      const { $t } = useTexts()
      return { $t }
    },
  })
</script>

<template>
  <ApiSection id="api-authentication" :title="$t('docs.auth.title')" icon-class="fas fa-key">
    <div class="content-card">
      <div class="card-header">
        <i class="fas fa-key" />
        <h4>{{ $t('docs.auth.apiKeyAuth') }}</h4>
        <span class="badge">{{ $t('docs.auth.recommended') }}</span>
      </div>

      <p>{{ $t('docs.auth.twoMethods') }}</p>

      <div class="examples-grid">
        <div class="example-item">
          <h5>{{ $t('docs.auth.httpHeader') }}</h5>
          <div class="code-snippet">x-pixelpunk-key: YOUR_API_KEY</div>
        </div>

        <div class="example-item">
          <h5>{{ $t('docs.auth.queryParam') }}</h5>
          <div class="code-snippet">?key=YOUR_API_KEY</div>
        </div>
      </div>
    </div>

    <div class="tips-card">
      <h4><i class="fas fa-shield-alt" />{{ $t('docs.auth.securityTips') }}</h4>
      <ul>
        <li>{{ $t('docs.auth.tip1') }}</li>
        <li>{{ $t('docs.auth.tip2') }}</li>
        <li>{{ $t('docs.auth.tip3') }}</li>
        <li>{{ $t('docs.auth.tip4') }}</li>
      </ul>
    </div>
  </ApiSection>
</template>

<style scoped>
  .content-card,
  .tips-card {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .tips-card {
    background: rgba(var(--color-success-rgb), 0.05);
    border-color: rgba(var(--color-success-rgb), 0.2);
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .card-header i {
    color: var(--color-brand-500);
    font-size: 1.25rem;
  }

  .card-header h4 {
    margin: 0;
    font-size: 1.125rem;
    color: var(--color-text-content);
  }

  .badge {
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
    background: rgba(var(--color-success-rgb), 0.2);
    color: var(--color-success-500);
  }

  .examples-grid {
    display: grid;
    gap: 1rem;
    margin-top: 1rem;
  }

  .example-item h5 {
    font-size: 0.875rem;
    color: var(--color-brand-500);
    margin: 0 0 0.5rem 0;
    font-weight: 600;
  }

  .code-snippet {
    background: rgba(var(--color-background-900-rgb), 0.6);
    color: var(--color-brand-500);
    padding: 0.75rem;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    font-family: 'SF Mono', Monaco, monospace;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .tips-card h4 {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 0 0 0.75rem 0;
    font-size: 1rem;
    color: var(--color-success-500);
  }

  .tips-card ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .tips-card li {
    padding: 0.375rem 0;
    color: var(--color-gray-200);
    font-size: 0.875rem;
    position: relative;
    padding-left: 1.25rem;
  }
</style>
