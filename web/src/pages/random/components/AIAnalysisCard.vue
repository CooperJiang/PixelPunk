<script setup lang="ts">
  import type { ImageInfo } from '@/api/types'
  import { useTexts } from '@/composables/useTexts'

  interface Props {
    imageData: ImageInfo
  }

  defineProps<Props>()
  const { $t } = useTexts()
</script>

<template>
  <div v-if="imageData.ai_info" class="ai-analysis-card">
    <h3 class="section-title"><i class="fas fa-robot mr-2" />{{ $t('random.aiAnalysis.title') }}</h3>
    <div class="ai-content">
      <div v-if="imageData.ai_info.description" class="ai-description">
        <p>{{ imageData.ai_info.description }}</p>
      </div>

      <div v-if="imageData.ai_info.tags && imageData.ai_info.tags.length" class="ai-tags">
        <h4 class="tags-subtitle">{{ $t('random.aiAnalysis.smartTags') }}</h4>
        <div class="tags-container">
          <CyberTag v-for="tag in imageData.ai_info.tags.slice(0, 8)" :key="tag" variant="primary" size="small" :title="tag">
            {{ tag }}
          </CyberTag>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .ai-analysis-card {
    background: rgba(var(--color-background-700-rgb), 0.35);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    backdrop-filter: blur(12px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    transition: all 0.3s ease;
  }

  .ai-analysis-card:hover {
    background: rgba(var(--color-background-700-rgb), 0.45);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }

  .section-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin-bottom: 1rem;
    text-shadow: 0 0 10px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .ai-content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .ai-description p {
    color: var(--color-content);
    line-height: 1.6;
    font-size: 0.85rem;
    margin: 0;
  }

  .tags-subtitle {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--color-brand-500);
    margin-bottom: 0.5rem;
  }

  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    align-items: flex-start;
  }
</style>
