<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { onMounted, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { shareApi } from '@/api'
  import { useToast } from '@/components/Toast/useToast'
  import { copyToClipboard } from '@/utils/file/clipboard'
  import type { ShareInfo } from '@/api/share/types'
  import VisitorInfoDialog from './components/VisitorInfoDialog.vue'
  import { useTexts } from '@/composables/useTexts'

  const router = useRouter()
  const toast = useToast()
  const { $t } = useTexts()

  const loading = ref(false)
  const shares = ref<ShareInfo[]>([])
  const pagination = ref<any>(null)
  const searchKeyword = ref('')
  const currentStatusFilter = ref<number | null>(null)
  const currentPage = ref(1)
  const pageSize = ref(10)

  const visitorDialogVisible = ref(false)
  const currentShareId = ref('')

  const deleteDialog = reactive({
    visible: false,
    shareId: '',
    shareName: '',
    force: false,
    loading: false,
  })

  const baseUrl = window.location.origin

  const filterTabs = [
    { label: $t('shares.filter.all'), value: null },
    { label: $t('shares.filter.normal'), value: 1 },
    { label: $t('shares.filter.expired'), value: 2 },
    { label: $t('shares.filter.deleted'), value: 3 },
    { label: $t('shares.filter.disabled'), value: 4 },
  ]

  const getStatusText = (status: number) => {
    switch (status) {
      case 1:
        return $t('shares.status.normal')
      case 2:
        return $t('shares.status.expired')
      case 3:
        return $t('shares.status.deleted')
      case 4:
        return $t('shares.status.disabled')
      default:
        return $t('shares.status.unknown')
    }
  }

  const getStatusClass = (status: number) => {
    switch (status) {
      case 1:
        return 'status-normal'
      case 2:
        return 'status-expired'
      case 3:
        return 'status-deleted'
      case 4:
        return 'status-disabled'
      default:
        return ''
    }
  }

  const formatDate = (dateString: string) => {
    if (!dateString) {
      return ''
    }

    const date = new Date(dateString)
    return date
      .toLocaleDateString(getCurrentLocale(), {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      })
      .replace(/\//g, '-')
  }

  const loadShares = async () => {
    loading.value = true

    const params: any = {
      page: currentPage.value,
      size: pageSize.value,
    }

    if (currentStatusFilter.value !== null) {
      params.status = currentStatusFilter.value
    }

    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    try {
      const result = await shareApi.getShareList(params)

      if (result.success) {
        const { data } = result

        shares.value = data.list || []
        pagination.value = {
          current_page: data.current_page || 1,
          last_page: data.last_page || 1,
          total: data.total || 0,
          from: (data.current_page - 1) * pageSize.value + 1 || 1,
          to: Math.min(data.current_page * pageSize.value, data.total) || 0,
        }
      }
    } catch (error) {
      console.error($t('shares.toast.loadFailed'), error)
    }

    loading.value = false
  }

  const handleFilterChange = (status: number | null) => {
    currentStatusFilter.value = status
    currentPage.value = 1
    loadShares()
  }

  const handleSearch = () => {
    currentPage.value = 1
    loadShares()
  }

  const handlePageSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    loadShares()
  }

  const handleCopyLink = async (share: ShareInfo) => {
    const url = share.url || `${baseUrl}/share/${share.share_key}`

    await copyToClipboard(url)
    toast.success($t('shares.toast.linkCopied'))
  }

  const handleOpenShare = (share: ShareInfo) => {
    const url = `/share/${share.share_key}`
    window.open(url, '_blank')
  }

  const confirmDelete = (share: ShareInfo) => {
    deleteDialog.shareId = share.id
    deleteDialog.shareName = share.name || $t('shares.unnamed')
    deleteDialog.force = false
    deleteDialog.visible = true
  }

  const confirmDeleteShare = async () => {
    try {
      deleteDialog.loading = true
      const result = await shareApi.deleteShare(deleteDialog.shareId, deleteDialog.force)

      if (result.success) {
        toast.success($t('shares.toast.deleteSuccess'))
        deleteDialog.visible = false
        loadShares()
      }
    } catch (error) {
      console.error($t('shares.toast.deleteFailed'), error)
    } finally {
      deleteDialog.loading = false
    }
  }

  const _handleDeleteShare = async (share: ShareInfo) => {
    confirmDelete(share)
  }

  const goToFolders = () => {
    router.push('/folders')
  }

  const handleViewVisitors = (share: ShareInfo) => {
    currentShareId.value = share.id
    visitorDialogVisible.value = true
  }

  onMounted(() => {
    loadShares()
  })
</script>

<template>
  <div class="shares-page">
    <div class="shares-container">
      <div class="header-section mb-4 rounded-lg p-4">
        <div class="flex items-center justify-between gap-4">
          <div class="flex items-center gap-3">
            <div class="page-header-icon">
              <i class="fas fa-share-alt text-brand-400" />
            </div>
            <div>
              <h1 class="text-lg font-bold text-content-heading">{{ $t('shares.page.title') }}</h1>
              <p class="text-xs text-content-muted">{{ $t('shares.page.subtitle') }}</p>
            </div>
          </div>
          <div class="flex gap-2">
            <CyberInput
              v-model="searchKeyword"
              :placeholder="$t('shares.search')"
              prefix-icon="search"
              clearable
              size="small"
              class="w-64"
              @input="handleSearch"
            />
          </div>
        </div>
      </div>

      <div class="filter-section mb-4 rounded-lg p-4">
        <div class="flex flex-wrap gap-2">
          <button
            v-for="tab in filterTabs"
            :key="tab.value"
            class="filter-tab rounded-md border px-4 py-2 text-sm font-medium transition-all duration-200"
            :class="
              currentStatusFilter === tab.value
                ? 'bg-brand-500/10 border-brand-400 text-brand-400'
                : 'border-default bg-background-700 text-content-muted hover:border-strong hover:text-content'
            "
            @click="handleFilterChange(tab.value)"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>

      <div class="content-area">
        <CyberSkeleton type="list" :count="4" :loading="loading && shares.length === 0" />

        <div v-if="!loading && shares.length > 0" class="flex flex-col gap-4">
          <div v-for="share in shares" :key="share.id" class="share-card rounded-lg border p-4 transition-all duration-200">
            <div class="mb-4 flex items-start justify-between">
              <div class="flex-1">
                <h3 class="mb-2 flex items-center gap-2 text-base font-semibold text-content-heading">
                  {{ share.name || $t('shares.unnamed') }}
                  <span class="share-status rounded px-2 py-0.5 text-xs font-normal" :class="getStatusClass(share.status)">
                    {{ getStatusText(share.status) }}
                  </span>
                </h3>
                <div class="mb-2 flex flex-wrap gap-3 text-xs text-content-muted">
                  <div class="flex items-center gap-1.5">
                    <i class="fas fa-calendar-alt" />
                    <span>{{ formatDate(share.created_at) }}</span>
                  </div>
                  <div class="flex items-center gap-1.5">
                    <i class="fas fa-eye" />
                    <span>{{ $t('shares.card.viewCount', { count: share.current_views }) }}</span>
                  </div>
                  <div v-if="share.max_views > 0" class="flex items-center gap-1.5">
                    <i class="fas fa-stop-circle" />
                    <span>{{ $t('shares.card.maxViews', { max: share.max_views }) }}</span>
                  </div>
                  <div v-if="share.expired_at" class="flex items-center gap-1.5">
                    <i class="fas fa-clock" />
                    <span>{{ formatDate(share.expired_at) }}</span>
                  </div>
                </div>
                <div v-if="share.description" class="text-sm text-content-muted">
                  {{ share.description }}
                </div>
              </div>
              <div class="flex flex-col items-end gap-2">
                <div class="flex flex-wrap justify-end gap-2">
                  <div
                    v-if="share.notification_on_access"
                    class="bg-brand-500/10 flex items-center gap-1.5 rounded border border-brand-400 px-2 py-1 text-xs font-medium text-brand-400"
                  >
                    <i class="fas fa-lock" />
                    <span>{{ $t('shares.card.passwordProtected') }}</span>
                  </div>
                  <div
                    v-if="share.collect_visitor_info"
                    class="bg-success-500/10 flex items-center gap-1.5 rounded border border-success-400 px-2 py-1 text-xs font-medium text-success-400"
                  >
                    <i class="fas fa-user-tag" />
                    <span>{{ $t('shares.card.collectVisitor') }}</span>
                  </div>
                  <div
                    v-if="share.notification_on_access"
                    class="bg-warning-500/10 flex items-center gap-1.5 rounded border border-warning-400 px-2 py-1 text-xs font-medium text-warning-400"
                  >
                    <i class="fas fa-bell" />
                    <span>{{ $t('shares.card.accessNotification') }}</span>
                  </div>
                </div>
                <CyberButton type="outlined" size="small" icon="copy" @click="handleCopyLink(share)">
                  {{ $t('shares.card.copyLink') }}
                </CyberButton>
              </div>
            </div>

            <div class="mb-3 flex gap-4 text-sm text-content-muted">
              <div class="flex items-center gap-1.5">
                <i class="fas fa-folder" />
                <span>{{ $t('shares.card.folderCount', { count: share.folder_count }) }}</span>
              </div>
              <div class="flex items-center gap-1.5">
                <i class="fas fa-image" />
                <span>{{ $t('shares.card.fileCount', { count: share.file_count }) }}</span>
              </div>
            </div>

            <div class="flex justify-end gap-2">
              <CyberButton type="outlined" size="small" icon="external-link-alt" @click="handleOpenShare(share)">
                {{ $t('shares.card.viewShare') }}
              </CyberButton>
              <CyberButton
                v-if="share.collect_visitor_info"
                type="secondary"
                size="small"
                icon="users"
                @click="handleViewVisitors(share)"
              >
                {{ $t('shares.card.visitorInfo') }}
              </CyberButton>
              <CyberButton type="danger" size="small" icon="trash" @click="confirmDelete(share)">
                {{ $t('shares.card.delete') }}
              </CyberButton>
            </div>
          </div>
        </div>

        <div v-if="pagination && pagination.total > 0" class="mt-4">
          <CyberPagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="pagination.total"
            :page-size-options="[5, 10, 20, 50]"
            :show-page-size-selector="true"
            :show-quick-jumper="true"
            @update:current-page="loadShares"
            @update:page-size="handlePageSizeChange"
          />
        </div>

        <div
          v-if="!loading && shares.length === 0"
          class="empty-state flex flex-col items-center justify-center rounded-lg p-16 text-center"
        >
          <div class="empty-icon-wrapper mb-4">
            <i class="fas fa-share-alt text-content-muted/30 text-4xl" />
          </div>
          <h3 class="mb-2 text-lg font-medium text-content-heading">{{ $t('shares.empty.title') }}</h3>
          <p class="mb-4 text-sm text-content-muted">{{ $t('shares.empty.subtitle') }}</p>
          <CyberButton type="primary" size="medium" icon="folder" @click="goToFolders">
            {{ $t('shares.empty.goToFolders') }}
          </CyberButton>
        </div>
      </div>
    </div>

    <VisitorInfoDialog v-model="visitorDialogVisible" :share-id="currentShareId" />

    <CyberDialog v-model="deleteDialog.visible" :title="$t('shares.dialog.deleteTitle')" width="400px">
      <div class="delete-dialog-content">
        <p>
          {{ $t('shares.delete.confirm', { name: deleteDialog.shareName }) }}
        </p>

        <div class="form-item force-delete p-4">
          <CyberCheckbox v-model="deleteDialog.force">{{ $t('shares.delete.forceDelete') }}</CyberCheckbox>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2 p-4">
          <CyberButton type="warning" @click="deleteDialog.visible = false">{{ $t('shares.delete.cancel') }}</CyberButton>
          <CyberButton type="danger" :loading="deleteDialog.loading" @click="confirmDeleteShare">{{
            $t('shares.delete.confirmButton')
          }}</CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>
</template>

<style scoped lang="scss">
  .shares-container {
    background: rgba(var(--color-background-800-rgb), 0.95);
    backdrop-filter: blur(10px);
    border: 1.5px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
    box-shadow:
      0 4px 12px var(--color-overlay-medium),
      0 2px 6px var(--color-overlay-light);
  }

  .header-section,
  .filter-section {
    background: rgba(var(--color-background-700-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .share-card {
    background: rgba(var(--color-background-700-rgb), 0.6);
    border-color: var(--color-border-default);
  }

  .share-card:hover {
    transform: translateY(-2px);
    border-color: var(--color-brand-400);
    box-shadow:
      0 4px 12px var(--color-overlay-medium),
      0 2px 6px var(--color-overlay-light);
  }

  .empty-state {
    background: rgba(var(--color-background-700-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .page-header-icon {
    width: 2.25rem;
    height: 2.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-brand-500-rgb), 0.05));
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
  }

  .share-status {
    border: 1px solid transparent;
  }

  .status-normal {
    background: rgba(var(--color-success-rgb), 0.1);
    color: var(--color-success-400);
    border-color: var(--color-success-400);
  }

  .status-expired {
    background: rgba(var(--color-warning-rgb), 0.1);
    color: var(--color-warning-400);
    border-color: var(--color-warning-400);
  }

  .status-deleted {
    background: rgba(var(--color-error-rgb), 0.1);
    color: var(--color-error-400);
    border-color: var(--color-error-400);
  }

  .status-disabled {
    background: rgba(var(--color-content-muted-rgb), 0.1);
    color: var(--color-content-muted);
    border-color: var(--color-content-muted);
  }

  .empty-icon-wrapper {
    width: 5rem;
    height: 5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.1), rgba(var(--color-brand-500-rgb), 0.05));
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-sm);
  }
</style>
