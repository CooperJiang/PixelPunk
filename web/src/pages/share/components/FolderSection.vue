<script setup lang="ts">
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  defineProps({
    folders: {
      type: Array,
      default: () => [],
    },
  })

  const emit = defineEmits(['folder-click'])

  const handleFolderClick = (folder: Record<string, unknown>) => {
    emit('folder-click', folder)
  }
</script>

<template>
  <div v-if="folders && folders.length > 0" class="share-section folders-section">
    <div class="section-header">
      <div class="section-title">
        <i class="fas fa-folder" />
        <span>{{ $t('share.folder.title') }}</span>
        <span class="count-badge">{{ folders.length }}</span>
      </div>
    </div>

    <div class="folders-grid">
      <CyberShareFolder v-for="folder in folders" :key="folder.id" :folder="folder" @click="handleFolderClick" />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .share-section {
    background: rgba(var(--color-background-800-rgb), 0.6);
    border-radius: var(--radius-lg);
    padding: var(--space-md);
    margin-bottom: var(--space-lg);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.15);
    box-shadow: var(--shadow-sm);
    backdrop-filter: var(--backdrop-blur-sm);
    transition: all var(--transition-normal) var(--ease-in-out);

    &:hover {
      border-color: rgba(var(--color-brand-500-rgb), 0.25);
      box-shadow: var(--shadow-md);
    }
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--space-md);
    position: relative;
    padding-bottom: var(--space-sm);
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .section-title {
    display: flex;
    align-items: center;
    color: var(--color-content-heading);
    font-size: var(--text-lg);
    font-weight: 600;
    gap: var(--space-sm);

    i {
      color: var(--color-brand-500);
      font-size: var(--text-lg);
    }
  }

  .count-badge {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
    border-radius: var(--radius-lg);
    padding: var(--space-xs) var(--space-sm);
    font-size: 0.625rem;
    font-weight: 500;
  }

  .folders-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: var(--space-md);
    width: 100%;
  }

  @media (min-width: 1600px) {
    .folders-grid {
      grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
      gap: var(--space-lg);
    }
  }

  @media (min-width: 1920px) {
    .folders-grid {
      grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    }
  }

  @media (max-width: 768px) {
    .folders-grid {
      grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
      gap: var(--space-sm);
    }
  }
</style>
