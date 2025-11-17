<script setup lang="ts">
  import { getCurrentLocale } from '@/utils/locale'
  import { computed, ref, watch } from 'vue'
  import { shareApi } from '@/api'
  import { useToast } from '@/components/Toast/useToast'
  import type { VisitorInfoDetail } from '@/api/share/types'
  import { useTexts } from '@/composables/useTexts'

  const { $t } = useTexts()

  const props = defineProps({
    modelValue: {
      type: Boolean,
      default: false,
    },
    shareId: {
      type: String,
      required: true,
    },
  })

  const emit = defineEmits(['update:modelValue'])

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  const loading = ref(false)
  const visitors = ref<VisitorInfoDetail[]>([])
  const searchKeyword = ref('')
  const currentPage = ref(1)
  const pagination = ref({
    current_page: 1,
    last_page: 1,
    total: 0,
    from: 1,
    to: 0,
  })

  const toast = useToast()

  watch(
    () => visible.value,
    (newValue) => {
      if (newValue) {
        loadVisitors()
      }
    }
  )

  const formatDate = (dateString: string) => {
    if (!dateString) {
      return '-'
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

  const loadVisitors = async () => {
    if (!props.shareId) {
      return
    }

    loading.value = true

    try {
      const params: any = {
        page: currentPage.value,
        size: 10,
      }

      if (searchKeyword.value) {
        params.keyword = searchKeyword.value
      }

      const result = await shareApi.getVisitorList(props.shareId, params)

      if (result.success && result.data) {
        visitors.value = result.data.list || []
        pagination.value = {
          current_page: result.data.current_page || 1,
          last_page: result.data.last_page || 1,
          total: result.data.total || 0,
          from: (result.data.current_page - 1) * 10 + 1 || 1,
          to: Math.min(result.data.current_page * 10, result.data.total) || 0,
        }
      }
    } catch (error) {
      console.error($t('shares.visitor.loadFailed'), error)
      toast.error($t('shares.visitor.loadFailed'))
    } finally {
      loading.value = false
    }
  }

  const handleSearch = () => {
    currentPage.value = 1
    loadVisitors()
  }

  const loadPage = (page: number) => {
    currentPage.value = page
    loadVisitors()
  }

  const handleDeleteVisitor = async (visitor: VisitorInfoDetail) => {
    try {
      await shareApi.deleteVisitor(props.shareId, visitor.id)

      toast.success($t('shares.visitor.deleteSuccess'))
      loadVisitors()
    } catch (error) {
      console.error($t('shares.visitor.deleteFailed'), error)
      toast.error($t('shares.visitor.deleteFailed'))
    }
  }

  const handleClose = () => {
    visible.value = false
    searchKeyword.value = ''
    currentPage.value = 1
  }
</script>

<template>
  <CyberDialog
    v-model="visible"
    :title="$t('shares.visitor.dialogTitle')"
    width="800px"
    :append-to-body="true"
    :show-default-footer="false"
    @close="handleClose"
  >
    <div class="visitor-info-dialog">
      <div class="search-bar">
        <CyberInput
          v-model="searchKeyword"
          :placeholder="$t('shares.visitor.search')"
          prefix-icon="search"
          clearable
          class="search-input"
          @input="handleSearch"
        />
      </div>

      <div v-if="visitors.length > 0" class="visitor-table">
        <div class="table-header">
          <div class="th name">{{ $t('shares.visitor.table.name') }}</div>
          <div class="th email">{{ $t('shares.visitor.table.email') }}</div>
          <div class="th ip">{{ $t('shares.visitor.table.ip') }}</div>
          <div class="th visits">{{ $t('shares.visitor.table.visits') }}</div>
          <div class="th date">{{ $t('shares.visitor.table.lastVisit') }}</div>
          <div class="th actions">{{ $t('shares.visitor.table.actions') }}</div>
        </div>
        <div class="table-body">
          <div v-for="visitor in visitors" :key="visitor.id" class="table-row">
            <div class="td name">{{ visitor.visitor_name }}</div>
            <div class="td email">{{ visitor.visitor_email || '-' }}</div>
            <div class="td ip">{{ visitor.ip_address || '-' }}</div>
            <div class="td visits">{{ visitor.visit_count }}</div>
            <div class="td date">{{ formatDate(visitor.last_visit_at) }}</div>
            <div class="td actions">
              <CyberPopconfirm
                :title="$t('shares.visitor.deleteConfirm')"
                :content="$t('shares.visitor.deleteHint')"
                @confirm="handleDeleteVisitor(visitor)"
              >
                <template #reference>
                  <button class="delete-btn">
                    <i class="fas fa-trash" />
                  </button>
                </template>
              </CyberPopconfirm>
            </div>
          </div>
        </div>
      </div>

      <div v-if="pagination.total > 0" class="pagination-container">
        <div class="pagination-info">
          {{ $t('shares.visitor.pagination.info', { total: pagination.total, from: pagination.from, to: pagination.to }) }}
        </div>
        <div class="pagination-controls">
          <button class="pagination-btn" :disabled="pagination.current_page === 1" @click="loadPage(pagination.current_page - 1)">
            <i class="fas fa-chevron-left" />
          </button>
          <span class="page-indicator">{{ pagination.current_page }} / {{ pagination.last_page }}</span>
          <button
            class="pagination-btn"
            :disabled="pagination.current_page === pagination.last_page"
            @click="loadPage(pagination.current_page + 1)"
          >
            <i class="fas fa-chevron-right" />
          </button>
        </div>
      </div>

      <div v-if="!loading && visitors.length === 0" class="empty-state">
        <div class="empty-icon">
          <i class="fas fa-users" />
        </div>
        <h3>{{ $t('shares.visitor.empty.title') }}</h3>
        <p>{{ $t('shares.visitor.empty.desc') }}</p>
      </div>

      <CyberLoading :visible="loading" :full-screen="false" :text="$t('status.loading')" />
    </div>
  </CyberDialog>
</template>

<style scoped>
  .visitor-info-dialog {
    max-height: 60vh;
    overflow-y: auto;
    padding-right: 10px;
    position: relative;
  }

  .visitor-info-dialog::-webkit-scrollbar {
    width: 4px;
  }

  .visitor-info-dialog::-webkit-scrollbar-track {
    background: rgba(var(--color-brand-500-rgb), 0.03);
    border-radius: var(--radius-sm);
  }

  .visitor-info-dialog::-webkit-scrollbar-thumb {
    background: rgba(var(--color-brand-500-rgb), 0.15);
    border-radius: var(--radius-sm);
  }

  .visitor-info-dialog::-webkit-scrollbar-thumb:hover {
    background: rgba(var(--color-brand-500-rgb), 0.3);
  }

  .search-bar {
    margin-bottom: 1rem;
  }

  .search-input {
    width: 300px;
  }

  .visitor-table {
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    border-radius: var(--radius-sm);
    overflow: hidden;
    background: rgba(var(--color-background-900-rgb), 0.5);
    margin-bottom: 1rem;
  }

  .table-header {
    display: grid;
    grid-template-columns: 2fr 2fr 1.5fr 1fr 2fr 1fr;
    padding: 0.6rem 1rem;
    background: rgba(var(--color-brand-500-rgb), 0.1);
    color: rgba(var(--color-brand-500-rgb), 0.9);
    font-weight: 600;
    font-size: 0.75rem;
    letter-spacing: 0.5px;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.15);
  }

  .table-row {
    display: grid;
    grid-template-columns: 2fr 2fr 1.5fr 1fr 2fr 1fr;
    padding: 0.6rem 1rem;
    align-items: center;
    border-bottom: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    transition: background 0.2s ease;
  }

  .table-row:last-child {
    border-bottom: none;
  }

  .table-row:hover {
    background: rgba(var(--color-brand-500-rgb), 0.05);
  }

  .th,
  .td {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .td {
    color: rgba(var(--color-white-rgb), 0.8);
    font-size: 0.85rem;
  }

  .td.actions {
    text-align: center;
    display: flex;
    justify-content: center;
  }

  .delete-btn {
    background: rgba(var(--color-error-rgb), 0.1);
    border: 1px solid rgba(var(--color-error-rgb), 0.2);
    color: rgba(var(--color-error-rgb), 0.8);
    width: 28px;
    height: 28px;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .delete-btn:hover {
    background: rgba(var(--color-error-rgb), 0.2);
    color: var(--color-error-500);
  }

  .pagination-container {
    margin-top: 1.25rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .pagination-info {
    color: rgba(var(--color-content-default-rgb), 0.6);
    font-size: 0.8rem;
  }

  .pagination-controls {
    display: flex;
    align-items: center;
    gap: 0.6rem;
  }

  .pagination-btn {
    width: 28px;
    height: 28px;
    border-radius: var(--radius-sm);
    background: rgba(var(--color-brand-500-rgb), 0.1);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
    color: rgba(var(--color-white-rgb), 0.8);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .pagination-btn:not(:disabled):hover {
    background: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.4);
    color: var(--color-brand-500);
  }

  .pagination-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .page-indicator {
    color: rgba(var(--color-content-default-rgb), 0.7);
    font-size: 0.8rem;
  }

  .empty-state {
    background: rgba(var(--color-background-900-rgb), 0.4);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    border-radius: var(--radius-sm);
    padding: 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    min-height: 200px;
  }

  .empty-icon {
    width: 50px;
    height: 50px;
    background: rgba(var(--color-brand-500-rgb), 0.05);
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 1rem;
    position: relative;
    overflow: hidden;
  }

  .empty-icon::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(var(--color-brand-500-rgb), 0.15), transparent);
    border-radius: var(--radius-full);
  }

  .empty-icon i {
    font-size: 1.5rem;
    color: rgba(var(--color-brand-500-rgb), 0.6);
    position: relative;
    z-index: 1;
  }

  .empty-state h3 {
    margin: 0 0 0.4rem;
    font-size: 1.1rem;
    color: rgba(var(--color-content-default-rgb), 0.85);
  }

  .empty-state p {
    margin: 0;
    color: rgba(var(--color-content-default-rgb), 0.5);
    font-size: 0.9rem;
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
  }
</style>
