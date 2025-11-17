<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import * as openApiService from '@/api/openapi'
  import type { RandomImageAPI } from '@/api/openapi'

  import ApiCard from './components/ApiCard.vue'
  import CreateApiModal from './components/CreateApiModal.vue'
  import EditApiModal from './components/EditApiModal.vue'

  const toast = useToast()
  const { $t } = useTexts()

  const isLoading = ref(false)
  const isSubmitting = ref(false)
  const isUpdating = ref(false)
  const isDeleting = ref(false)
  const totalApis = ref(0)
  const allApis = ref<RandomImageAPI[]>([])
  const searchKeyword = ref('')

  const showCreateModal = ref(false)
  const showEditModal = ref(false)
  const showDeleteConfirm = ref(false)
  const editingApi = ref<RandomImageAPI | null>(null)
  const deletingApi = ref<RandomImageAPI | null>(null)

  const filteredApis = computed(() => {
    if (!searchKeyword.value.trim()) {
      return allApis.value
    }
    const keyword = searchKeyword.value.toLowerCase().trim()
    return allApis.value.filter((api) => api.name.toLowerCase().includes(keyword))
  })

  const loadApis = async () => {
    try {
      isLoading.value = true
      const response = await openApiService.getRandomAPIList({
        page: 1,
        size: 1000,
      })
      allApis.value = response.data.items || []
      totalApis.value = response.data.total
    } catch (_error) {
      toast.error($t('openApi.toast.loadFailed'))
    } finally {
      isLoading.value = false
    }
  }

  const submitApiForm = async (formData: openApiService.CreateRandomAPIRequest) => {
    try {
      isSubmitting.value = true
      await openApiService.createRandomAPI(formData)
      toast.success($t('openApi.toast.createSuccess'))
      showCreateModal.value = false
      loadApis()
    } catch (_error) {
      toast.error($t('openApi.toast.createFailed'))
    } finally {
      isSubmitting.value = false
    }
  }

  const openEditModal = (api: RandomImageAPI) => {
    editingApi.value = api
    showEditModal.value = true
  }

  const submitEditForm = async (formData: openApiService.UpdateRandomAPIConfigRequest) => {
    if (!editingApi.value) {
      return
    }

    try {
      isUpdating.value = true
      await openApiService.updateRandomAPIConfig(editingApi.value.id, formData)
      toast.success($t('openApi.toast.updateSuccess'))
      showEditModal.value = false
      loadApis()
    } catch (_error) {
      toast.error($t('openApi.toast.updateFailed'))
    } finally {
      isUpdating.value = false
    }
  }

  const confirmDeleteApi = (api: RandomImageAPI) => {
    deletingApi.value = api
    showDeleteConfirm.value = true
  }

  const deleteApi = async () => {
    if (!deletingApi.value) {
      return
    }

    try {
      isDeleting.value = true
      await openApiService.deleteRandomAPI(deletingApi.value.id)
      toast.success($t('openApi.toast.deleteSuccess'))
      showDeleteConfirm.value = false
      loadApis()
    } catch (_error) {
      toast.error($t('openApi.toast.deleteFailed'))
    } finally {
      isDeleting.value = false
    }
  }

  const toggleApiStatus = async (api: RandomImageAPI) => {
    try {
      const newStatus = api.is_active ? 2 : 1
      await openApiService.updateRandomAPIStatus(api.id, { status: newStatus })
      toast.success($t('openApi.toast.toggleSuccess'))
      loadApis()
    } catch (_error) {
      toast.error($t('openApi.toast.toggleFailed'))
    }
  }

  const copyUrl = (url: string) => {
    navigator.clipboard.writeText(url).then(() => {
      toast.success($t('openApi.toast.copySuccess'))
    })
  }

  const getApiUrl = (apiKey: string) => {
    return `${window.location.origin}/api/v1/r/${apiKey}`
  }

  onMounted(() => {
    loadApis()
  })
</script>

<template>
  <div class="open-api-page">
    <div class="page-header mb-4 p-4">
      <div class="flex items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="page-header-icon flex h-9 w-9 items-center justify-center rounded-lg border">
            <i class="fas fa-random text-lg text-brand-400" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h1 class="text-lg font-bold text-content-heading">{{ $t('openApi.page.title') }}</h1>
              <span
                class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium text-brand-400"
                :style="{
                  backgroundColor: 'rgba(var(--color-brand-500-rgb), 0.1)',
                }"
              >
                {{ totalApis }}
              </span>
            </div>
            <p class="mt-0.5 text-xs text-content-muted">{{ $t('openApi.page.subtitle') }}</p>
          </div>
        </div>

        <div class="flex gap-2">
          <div class="search-input-wrapper">
            <CyberInput v-model="searchKeyword" :placeholder="$t('openApi.page.search')" size="small" clearable>
              <template #prefix>
                <i class="fas fa-search text-content-muted"></i>
              </template>
            </CyberInput>
          </div>
          <CyberButton type="primary" size="small" @click="showCreateModal = true">
            <i class="fas fa-plus mr-1.5" />
            {{ $t('openApi.page.create') }}
          </CyberButton>
        </div>
      </div>
    </div>

    <!-- API卡片内容区域 -->
    <div class="apis-container">
      <div class="apis-content">
        <div v-if="isLoading" class="loading-state">
          <i class="fas fa-spinner fa-spin text-2xl text-brand-500"></i>
          <p class="mt-2 text-content-muted">{{ $t('openApi.loading') }}</p>
        </div>

        <div v-else-if="filteredApis.length === 0" class="empty-state">
          <div class="empty-icon-wrapper">
            <i class="fas fa-random text-content-muted/30 text-5xl"></i>
          </div>
          <h3 class="mt-4 text-lg font-medium text-content">
            {{ $t('openApi.empty.title') }}
          </h3>
          <p class="mt-2 max-w-md text-center text-sm text-content-muted">
            {{ $t('openApi.empty.description') }}
          </p>
          <CyberButton v-if="!searchKeyword" type="primary" size="medium" class="mt-4" @click="showCreateModal = true">
            <i class="fas fa-plus mr-2" />
            {{ $t('openApi.page.create') }}
          </CyberButton>
        </div>

        <div v-else class="apis-grid">
          <ApiCard
            v-for="api in filteredApis"
            :key="api.id"
            :api="api"
            :url="getApiUrl(api.api_key)"
            @edit="openEditModal"
            @toggle-status="toggleApiStatus"
            @delete="confirmDeleteApi"
            @copy-url="copyUrl"
          />
        </div>
      </div>

      <div v-if="!isLoading && filteredApis.length > 0" class="apis-footer">
        <div class="total-count">{{ $t('openApi.footer.total', { count: filteredApis.length }) }}</div>
      </div>
    </div>

    <CreateApiModal v-model:visible="showCreateModal" :is-submitting="isSubmitting" @submit="submitApiForm" />

    <EditApiModal v-model:visible="showEditModal" :is-submitting="isUpdating" :api="editingApi" @submit="submitEditForm" />

    <CyberDialog
      :model-value="showDeleteConfirm"
      :title="$t('openApi.dialog.delete.title')"
      width="450px"
      :loading="isDeleting"
      :show-default-footer="true"
      @confirm="deleteApi"
      @update:model-value="(val) => (showDeleteConfirm = val)"
    >
      <p class="text-content-default">
        {{ $t('openApi.dialog.delete.message', { name: deletingApi?.name || '' }) }}
      </p>
      <p class="mt-2 text-sm text-content-muted">{{ $t('openApi.dialog.delete.warning') }}</p>
    </CyberDialog>
  </div>
</template>

<style scoped lang="scss">
  .open-api-page {
    .page-header {
      position: relative;
      overflow: hidden;
      background: rgba(var(--color-background-800-rgb), 0.6);
      border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
      border-radius: var(--radius-sm);
      box-shadow: var(--shadow-sm);
      backdrop-filter: var(--backdrop-blur-md);

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.05), transparent);
        pointer-events: none;
      }

      .page-header-icon {
        background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-500-rgb), 0.06));
        border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
        border-radius: var(--radius-sm);
      }
    }
  }

  .search-input-wrapper {
    width: 240px;
  }

  .apis-container {
    display: flex;
    flex-direction: column;
    position: relative;
    background: rgba(var(--color-background-800-rgb), 0.6);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-sm);
    backdrop-filter: var(--backdrop-blur-md);
    overflow: hidden;
  }

  .apis-content {
    flex: 1;
    min-height: 0;
    padding: 1.5rem;
  }

  .loading-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    text-align: center;
  }

  .empty-icon-wrapper {
    width: 5rem;
    height: 5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.12), rgba(var(--color-brand-500-rgb), 0.06));
    border: 2px solid rgba(var(--color-brand-500-rgb), 0.2);
    border-radius: var(--radius-2xl);
  }

  .apis-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 1.25rem;
  }

  @media (min-width: 1920px) {
    .apis-grid {
      grid-template-columns: repeat(4, 1fr);
    }
  }

  @media (min-width: 1440px) and (max-width: 1919px) {
    .apis-grid {
      grid-template-columns: repeat(3, 1fr);
    }
  }

  @media (min-width: 1024px) and (max-width: 1439px) {
    .apis-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  .apis-footer {
    flex-shrink: 0;
    display: flex;
    justify-content: center;
    background: rgba(var(--color-background-800-rgb), 0.5);
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    padding: 1rem;
    backdrop-filter: var(--backdrop-blur-sm);
  }

  .total-count {
    text-align: center;
    color: var(--color-content-muted);
    font-size: 0.875rem;
  }

  @media (max-width: 768px) {
    .apis-grid {
      grid-template-columns: 1fr;
    }

    .search-input-wrapper {
      width: 100%;
    }
  }
</style>
