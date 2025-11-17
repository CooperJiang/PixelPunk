<script setup lang="ts">
  import type { Image } from '@/api/types'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineOptions({
    name: 'ImageGrid',
  })

  defineProps<Props>()

  defineEmits<{
    'select-image': [id: number]
    'edit-image': [image: Image]
    'delete-image': [image: Image]
    'toggle-recommendation': [image: Image]
    'view-details': [image: Image]
  }>()

  interface Props {
    images: Image[]
    loading: boolean
    selectMode: boolean
    selectedImages: number[]
    pageSize: number
  }
</script>

<template>
  <div class="images-grid-container">
    <div v-if="loading" class="images-grid">
      <CyberSkeleton type="card" :count="pageSize" :loading="loading" />
    </div>

    <div v-else class="images-grid">
      <FileCard
        v-for="image in images"
        :key="image.id"
        :image="image"
        :select-mode="selectMode"
        :selected="selectedImages.includes(image.id)"
        @select="$emit('select-image', image.id)"
        @edit="$emit('edit-image', image)"
        @delete="$emit('delete-image', image)"
        @toggle-recommendation="$emit('toggle-recommendation', image)"
        @view-details="$emit('view-details', image)"
      />
    </div>

    <div v-if="!loading && images.length === 0" class="empty-state">
      <div class="empty-icon">
        <i class="fas fa-image" />
      </div>
      <h3 class="empty-title">{{ $t('admin.files.list.empty') }}</h3>
      <p class="empty-desc">{{ $t('admin.files.list.emptyDesc') }}</p>
    </div>
  </div>
</template>

<style scoped>
  .images-grid-container {
    min-height: 400px;
  }

  .images-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: var(--space-lg);
    padding: var(--space-md) 0;
  }

  .skeleton-item {
    height: 320px;
    border-radius: var(--radius-lg);
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-2xl);
    text-align: center;
    min-height: 400px;
    background: rgba(var(--color-background-700-rgb), 0.3);
    border: 2px dashed var(--color-border-default);
    border-radius: var(--radius-lg);
    backdrop-filter: blur(4px);
  }

  .empty-icon {
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, var(--color-brand-500), var(--color-purple-500));
    border-radius: var(--radius-full);
    margin-bottom: var(--space-lg);
    box-shadow: 0 4px 20px rgba(var(--color-brand-500-rgb), 0.3);
  }

  .empty-icon i {
    font-size: 2rem;
    color: var(--color-content-heading);
  }

  .empty-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--color-content-default);
    margin: 0 0 var(--space-sm);
  }

  .empty-desc {
    font-size: 0.875rem;
    color: var(--color-content-default-secondary);
    line-height: 1.6;
    max-width: 400px;
    margin: 0;
  }

  @media (max-width: 1024px) {
    .images-grid {
      grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
      gap: var(--space-md);
    }
  }

  @media (max-width: 768px) {
    .images-grid {
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: var(--space-sm);
    }

    .empty-state {
      padding: var(--space-xl);
      min-height: 300px;
    }

    .empty-icon {
      width: 60px;
      height: 60px;
    }

    .empty-icon i {
      font-size: 1.5rem;
    }

    .empty-title {
      font-size: 1.125rem;
    }
  }

  @media (max-width: 480px) {
    .images-grid {
      grid-template-columns: 1fr;
      gap: var(--space-sm);
    }
  }
</style>
