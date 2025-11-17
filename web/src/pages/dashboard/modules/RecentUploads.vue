<script setup lang="ts">
  import { onMounted, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { type ImageInfo, getFileList } from '@/api/file'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import RecentFileItem from './RecentFileItem.vue'

  defineOptions({
    name: 'RecentUploads',
  })

  const props = defineProps<{
    refreshKey?: number
  }>()

  const router = useRouter()
  const toast = useToast()
  const { $t } = useTexts()
  const loading = ref(false)
  const recentImages = ref<ImageInfo[]>([])

  const copyAllUrls = async () => {
    if (recentImages.value.length === 0) {
      toast.warning($t('dashboard.recentUploads.toast.noImages'))
      return
    }

    const urls = recentImages.value
      .map((img) => img.full_url || img.url || '')
      .filter(Boolean)
      .join('\n')

    try {
      await navigator.clipboard.writeText(urls)
      toast.success($t('dashboard.recentUploads.toast.copySuccess').replace('{count}', recentImages.value.length.toString()))
    } catch {
      toast.error($t('dashboard.recentUploads.toast.copyFailed'))
    }
  }

  const fetchRecentImages = async () => {
    try {
      loading.value = true
      const result = await getFileList({
        page: 1,
        size: 10,
        sort: 'newest', // 按时间排序获取最新的
      })

      if (result.success) {
        recentImages.value = result.data?.items || []
      }
    } catch {
      toast.error($t('dashboard.recentUploads.toast.fetchFailed'))
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchRecentImages()
  })

  watch(
    () => props.refreshKey,
    () => {
      if (props.refreshKey !== undefined) {
        fetchRecentImages()
      }
    }
  )
</script>

<template>
  <div class="recent-uploads cyber-card">
    <div class="uploads-content">
      <div class="uploads-header">
        <h3 class="section-title">
          <i class="fas fa-history" />
          {{ $t('dashboard.recentUploads.title') }}
        </h3>
        <div class="header-actions">
          <button class="action-btn" :title="$t('dashboard.recentUploads.actions.copyAll')" @click="copyAllUrls">
            <i class="fas fa-copy" />
          </button>
          <button class="action-btn" :title="$t('dashboard.recentUploads.actions.viewAll')" @click="router.push('/folders')">
            <i class="fas fa-external-link-alt" />
          </button>
        </div>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-spinner">
          <i class="fas fa-spinner fa-spin" />
        </div>
        <div class="loading-text">{{ $t('dashboard.recentUploads.loading') }}</div>
      </div>

      <div v-else-if="recentImages.length > 0" class="recent-images-container">
        <RecentFileItem v-for="image in recentImages" :key="image.id" :image="image" />
      </div>

      <div v-else class="empty-state">
        <div class="empty-icon">
          <i class="fas fa-image" />
        </div>
        <div class="empty-text">{{ $t('dashboard.recentUploads.empty') }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .recent-uploads {
    padding: 20px;
  }

  .uploads-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
    height: 100%;
    max-height: 400px;
  }

  .uploads-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .section-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--color-content-heading);
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .section-title i {
    color: var(--color-badge-accent-text);
    font-size: 16px;
  }

  .header-actions {
    display: flex;
    gap: 8px;
  }

  .action-btn {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-hover-bg-neutral);
    border: 1px solid var(--color-border-default);
    border-radius: var(--radius-md);
    color: var(--color-content-muted);
    font-size: 12px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .action-btn:hover {
    background: var(--color-hover-bg);
    border-color: var(--color-hover-border);
    color: var(--color-badge-accent-text);
  }

  .recent-images-container {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    grid-template-rows: repeat(2, 100px);
    gap: 12px;
    justify-items: center;
    align-items: start;
    padding: 4px;
    max-height: 220px;
  }

  .empty-state {
    text-align: center;
    padding: 40px 20px;
    color: var(--color-content-muted);
  }

  .empty-icon {
    font-size: 36px;
    margin-bottom: 12px;
    color: var(--color-badge-accent-text);
    opacity: 0.3;
  }

  .empty-text {
    font-size: 14px;
    color: var(--color-content-disabled);
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
    color: var(--color-content-muted);
  }

  .loading-spinner {
    font-size: 24px;
    margin-bottom: 12px;
    color: var(--color-badge-accent-text);
  }

  .loading-spinner i {
    animation: spin 1s linear infinite;
  }

  .loading-text {
    font-size: 14px;
    color: var(--color-content-disabled);
  }

  @media (max-width: 1200px) {
    .recent-images-container {
      grid-template-columns: repeat(4, 1fr);
      grid-template-rows: repeat(2, 100px);
    }
  }

  @media (max-width: 768px) {
    .recent-uploads {
      padding: 16px;
    }

    .recent-images-container {
      grid-template-columns: repeat(3, 1fr);
      grid-template-rows: repeat(2, 90px);
      gap: 10px;
      max-height: 200px;
    }
  }

  @media (max-width: 480px) {
    .recent-images-container {
      grid-template-columns: repeat(2, 1fr);
      grid-template-rows: repeat(2, 80px);
      gap: 8px;
      max-height: 180px;
    }
  }
</style>
