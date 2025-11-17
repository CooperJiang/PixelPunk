<script lang="ts">
  import { defineComponent } from 'vue'
  import ApiSection from '../ApiSection/index.vue'
  import { useTexts } from '@/composables/useTexts'

  export default defineComponent({
    name: 'ApiOverview',
    components: {
      ApiSection,
    },
    props: {
      currentDomain: {
        type: String,
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
  <ApiSection id="api-overview" :title="$t('docs.overview.title')" icon-class="fas fa-book">
    <div class="section-intro">
      <p>{{ $t('docs.overview.description') }}</p>
    </div>

    <div class="info-cards">
      <div class="info-card">
        <div class="card-icon">
          <i class="fas fa-server" />
        </div>
        <div class="card-content">
          <h4>{{ $t('docs.overview.apiAddress') }}</h4>
          <code>{{ currentDomain }}/api/v1</code>
        </div>
      </div>

      <div class="info-card">
        <div class="card-icon">
          <i class="fas fa-shield-alt" />
        </div>
        <div class="card-content">
          <h4>{{ $t('docs.overview.authMethod') }}</h4>
          <p>{{ $t('docs.overview.apiKeyAuth') }}</p>
        </div>
      </div>

      <div class="info-card">
        <div class="card-icon">
          <i class="fas fa-zap" />
        </div>
        <div class="card-content">
          <h4>{{ $t('docs.overview.performance') }}</h4>
          <p>{{ $t('docs.overview.avgResponse') }}</p>
        </div>
      </div>
    </div>
  </ApiSection>
</template>

<style scoped>
  .section-intro {
    margin-bottom: 1.25rem;
  }

  .section-intro p {
    font-size: 0.95rem;
    color: var(--color-gray-200);
    line-height: 1.5;
    margin: 0;
  }

  .info-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 1rem;
    margin-bottom: 1.25rem;
  }

  .info-card {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 1rem;
    transition: all 0.2s ease;
  }

  .info-card:hover {
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    transform: translateY(-2px);
  }

  .card-icon {
    font-size: 1.25rem;
    color: var(--color-brand-500);
    margin-bottom: 0.5rem;
  }

  .card-content h4 {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--color-text-content);
    margin: 0 0 0.375rem 0;
  }

  .card-content p {
    font-size: 0.8125rem;
    color: var(--color-slate-500);
    margin: 0;
    line-height: 1.4;
  }

  .card-content code {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    color: var(--color-brand-500);
    padding: 0.25rem 0.375rem;
    border-radius: var(--radius-sm);
    font-family: 'SF Mono', Monaco, monospace;
    font-size: 0.8125rem;
    word-break: break-all;
  }
</style>
