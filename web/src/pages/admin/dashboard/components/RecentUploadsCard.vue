<script setup lang="ts">
  import { computed, nextTick, onMounted, ref } from 'vue'
  import { getDashboardRecentUploads, type DashboardRecentUploadsResponse } from '@/api/admin/dashboard'
  import { useResponsivePageSize } from '@/hooks/useResponsivePageSize'
  import { useTexts } from '@/composables/useTexts'
  const { $t } = useTexts()

  const loading = ref(false)
  const showSkeleton = ref(false)
  const data = ref<DashboardRecentUploadsResponse>({
    recent_uploads: [],
  })
  const imageLoadingStates = ref<Record<string, boolean>>({})
  const showFilePreview = ref(false)
  const selectedFile = ref<any>(null)

  /* 动态页面大小计算 - 基于网格布局自动计算合理的文件数量 */
  const { pageSize: autoPageSize, recalc } = useResponsivePageSize({
    containerSelector: '.recent-uploads-card .card-content',
    gridSelector: '.recent-uploads-card .images-grid',
    childSelector: '.image-item',
    itemMinWidth: 200, // 增大最小宽度，减少列数
    columnGap: 16, // 增大间距到1rem
    rowMultiple: 2, // 减少到2行
    defaultSize: 6, // 默认6个资源
    debug: false,
    preferCssColumns: true,
    mode: 'once', // 只计算一次，避免频繁重算
  })

  const currentPreviewIndex = computed(() => {
    if (!selectedFile.value || !data.value.recent_uploads) return 0
    return data.value.recent_uploads.findIndex((img) => img.id === selectedFile.value?.id) || 0
  })

  const handlePrevImage = () => {
    if (!data.value.recent_uploads) return
    const currentIndex = currentPreviewIndex.value
    if (currentIndex > 0) {
      selectedFile.value = data.value.recent_uploads[currentIndex - 1]
    }
  }

  const handleNextImage = () => {
    if (!data.value.recent_uploads) return
    const currentIndex = currentPreviewIndex.value
    if (currentIndex < data.value.recent_uploads.length - 1) {
      selectedFile.value = data.value.recent_uploads[currentIndex + 1]
    }
  }

  function viewImage(image: any) {
    selectedFile.value = image
    showFilePreview.value = true
  }

  async function fetchData() {
    try {
      loading.value = true
      showSkeleton.value = true

      const limit = autoPageSize.value || 8
      const result = await getDashboardRecentUploads(limit)

      if (result.success && result.data && Array.isArray(result.data)) {
        const validUploads = result.data
          .filter((item) => {
            if (!item.id) {
              return false
            }
            return true
          })
          .map((item, index) => ({
            ...item,
            id: item.id || `temp-${Date.now()}-${index}`,
          }))

        data.value = { recent_uploads: validUploads }
      } else {
        data.value = { recent_uploads: [] }
      }
    } catch {
      data.value = { recent_uploads: [] }
    } finally {
      loading.value = false
      showSkeleton.value = false
    }
  }

  function refresh() {
    recalc(true)
    fetchData()
  }

  onMounted(async () => {
    loading.value = true
    showSkeleton.value = true

    await nextTick()
    setTimeout(() => {
      recalc(true)
      fetchData()
    }, 200)
  })
</script>

<template>
  <div class="recent-uploads-card">
    <div class="card-header">
      <div class="header-icon">
        <i class="fas fa-clock" />
      </div>
      <h3>{{ $t('admin.dashboard.recentUploads.title') }}</h3>
      <div class="refresh-btn" @click="refresh">
        <i class="fas fa-refresh" :class="{ 'animate-spin': loading }" />
      </div>
    </div>

    <div class="card-content">
      <div class="images-grid">
        <CyberSkeleton type="card" :count="autoPageSize || 6" :loading="showSkeleton" simple />
      </div>

      <div v-if="!loading && !showSkeleton && (!data.recent_uploads || data.recent_uploads.length === 0)" class="empty-container">
        <i class="fas fa-image" />
        <p>{{ $t('admin.dashboard.recentUploads.noRecentUploads') }}</p>
      </div>

      <div v-if="!showSkeleton && data.recent_uploads && data.recent_uploads.length > 0" class="images-grid">
        <div
          v-for="(image, index) in data.recent_uploads"
          :key="`image-${image.id || index}`"
          class="image-item"
          @click="viewImage(image)"
        >
          <div class="image-preview" v-loading="imageLoadingStates[image.id]">
            <CyberFile
              :src="image.full_thumb_url"
              :alt="image.display_name || image.file_name"
              :is-nsfw="image?.ai_info?.is_nsfw"
              width="100%"
              height="100%"
              fit-mode="cover"
              @loading="imageLoadingStates[image.id] = $event"
              :show-loading-text="false"
            />

            <cyber-file-expiry-tag
              v-if="image.is_time_limited"
              :expires-at="image.expires_at"
              :storage-duration="image.storage_duration"
              :is-time-limited="image.is_time_limited"
              position="top-left"
              mode="both"
              :show-icon="true"
            />

            <div class="image-badges">
              <div class="image-badge">{{ image.size_formatted }}</div>
              <div v-if="image.format" class="image-badge">{{ image.format.toUpperCase() }}</div>
            </div>
          </div>

          <div class="image-info">
            <div class="image-name" :title="image.display_name || image.file_name">
              {{ image.display_name || image.file_name }}
            </div>
            <div class="user-info">
              <i class="fas fa-user" />
              <span>{{ image.user?.username || $t('admin.dashboard.recentUploads.unknownUser') }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <cyberFileViewer
      v-model="showFilePreview"
      :file="selectedFile"
      :files="data.recent_uploads"
      :initial-index="currentPreviewIndex"
      :show-side-nav="true"
      :show-keyboard-tips="true"
      @close="showFilePreview = false"
      @prev="handlePrevImage"
      @next="handleNextImage"
    />
  </div>
</template>

<style scoped>
  .recent-uploads-card {
    @apply flex h-full flex-col;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.6) 0%,
      rgba(var(--color-background-900-rgb), 0.4) 100%
    );
    backdrop-filter: blur(10px);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-cyber-sm);
    overflow: hidden;
  }

  .card-header {
    @apply flex items-center gap-2;
    padding: 1rem 1rem 0.75rem 1rem;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .header-icon {
    @apply flex h-10 w-10 flex-shrink-0 items-center justify-center;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    color: var(--color-brand-500);
  }

  .card-header h3 {
    color: var(--color-content-heading);
    font-size: 0.95rem;
    font-weight: 600;
    margin-bottom: 2px;
    line-height: 1.2;
    flex: 1;
  }

  .refresh-btn {
    @apply flex h-8 w-8 cursor-pointer items-center justify-center transition-colors;
    color: var(--color-content-default);
  }

  .refresh-btn:hover {
    color: var(--color-brand-500);
  }

  .card-content {
    @apply flex-1;
    padding: 0.75rem 1rem 1rem 1rem;
  }

  .loading-container,
  .empty-container {
    @apply flex h-full flex-col items-center justify-center;
    color: var(--color-content-disabled);
  }

  .loading-spinner {
    width: 30px;
    height: 30px;
    border: 3px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-full);
    border-top-color: var(--color-brand-500);
    animation: spin 1s ease-in-out infinite;
    margin-bottom: 1rem;
  }

  .empty-container i {
    font-size: 2.5rem;
    margin-bottom: 0.5rem;
    opacity: 0.4;
  }

  .images-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
  }

  .image-item {
    background: var(--color-background-700);
    border-radius: var(--radius-sm);
    overflow: hidden;
    box-shadow: var(--shadow-cyber-sm);
    display: flex;
    flex-direction: column;
    height: 192px;
    border: 1px solid var(--color-border-default);
    cursor: pointer;
    transition: all 0.3s ease;
    position: relative;
  }

  .image-item::after {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: var(--radius-sm);
    opacity: 0;
    transition: opacity 0.3s ease;
    pointer-events: none;
    z-index: 1;
  }

  .image-item:hover {
    box-shadow:
      0 4px 12px rgba(var(--color-brand-500-rgb), 0.15),
      0 0 0 1px rgba(var(--color-brand-500-rgb), 0.2);
  }

  .image-item:hover::after {
    opacity: 1;
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
  }

  .image-item:hover .image-preview :deep(.cyber-file-container),
  .image-item:hover .image-preview :deep(img) {
    transform: scale(1.15) translateZ(30px);
  }

  .image-preview {
    position: relative;
    height: 140px;
    min-height: 140px;
    overflow: hidden;
    background-color: var(--color-background-800);
  }

  .image-preview :deep(.cyber-file-container),
  .image-preview :deep(img) {
    transition: transform 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    transform-style: preserve-3d;
  }

  .image-badges {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    display: flex;
    gap: 0.5rem;
    z-index: 10;
  }

  .image-badge {
    padding: 0.2rem 0.45rem;
    background: rgba(var(--color-background-900-rgb), 0.8);
    backdrop-filter: blur(4px);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.25);
    border-radius: var(--radius-sm);
    color: var(--color-content-heading);
    font-size: 0.65rem;
    font-weight: 600;
    white-space: nowrap;
    line-height: 1;
  }

  :deep(.cyber-file-expiry-tag) {
    z-index: 5;
  }

  .image-info {
    padding: 0.75rem;
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    z-index: 2;
    position: relative;
    min-height: 52px;
  }

  .image-name {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--color-content-heading);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.3;
    flex: 1;
    min-width: 0;
    max-width: 140px;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    flex-shrink: 0;
    font-size: 0.7rem;
    color: var(--color-content-default);
  }

  .user-info i {
    font-size: 0.7rem;
    color: rgba(var(--color-brand-500-rgb), 0.7);
    flex-shrink: 0;
    width: 12px;
    text-align: center;
  }

  .user-info span {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-weight: 500;
    max-width: 80px;
  }
</style>
