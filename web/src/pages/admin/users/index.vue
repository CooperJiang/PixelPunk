<script setup lang="ts">
  import { computed, h, onMounted, reactive, ref } from 'vue'
  import { formatDate } from '@/utils/formatting/format'
  import { useToast } from '@/components/Toast/useToast'
  import { useTexts } from '@/composables/useTexts'
  import { getUserList, updateUser, toggleUserStatus as toggleUserStatusApi, type UserItem } from '@/api/admin/user'
  import UserForm from './components/UserForm.vue'
  import UserFilterPanel from './components/UserFilterPanel.vue'
  import CreateUserDialog from './components/CreateUserDialog.vue'
  import UserDetailDrawer from './components/UserDetailDrawer.vue'
  import BatchOperations from './components/BatchOperations.vue'
  import type { CyberTableColumn } from '@/components/Table/types'
  import { showConfirm } from '@/utils/dialog'

  defineOptions({
    name: 'AdminUsers',
  })

  /* 工具 */
  const { $t } = useTexts()
  const toast = useToast()

  /* 表格列配置 */
  const dataTableColumns = computed<CyberTableColumn<UserItem>[]>(() => [
    {
      key: 'username',
      title: $t('admin.users.columns.user'),
      width: 160,
      align: 'left',
      render: (value, record) =>
        h('div', { class: 'flex items-center' }, [
          h('div', { class: 'mr-3' }, [
            h('CyberUserAvatar', {
              avatarUrl: record.avatar_full_path,
              username: record.username,
              size: 'sm',
            }),
          ]),
          h('div', {}, [
            h('div', { class: 'font-medium', style: { color: 'var(--color-content-default)' } }, record.username),
            h('div', { class: 'text-xs', style: { color: 'var(--color-content-muted)' } }, `ID: ${record.id}`),
          ]),
        ]),
    },
    {
      key: 'email',
      title: $t('admin.users.columns.email'),
      dataIndex: 'email',
      width: 200,
      align: 'center',
      ellipsis: true,
    },
    {
      key: 'created_at',
      title: $t('admin.users.columns.registeredAt'),
      width: 180,
      align: 'center',
      render: (value) => formatDate(value),
    },
    {
      key: 'last_activity_at',
      title: $t('admin.users.columns.lastActivity'),
      width: 160,
      align: 'center',
      render: (value, record) =>
        h('div', { class: 'last-activity-cell' }, [
          record.last_activity_at
            ? h('div', { class: 'activity-time' }, formatDate(record.last_activity_at))
            : h('div', { class: 'no-activity' }, $t('admin.users.columns.neverActive')),
          record.last_activity_ip
            ? h('div', { class: 'activity-ip', title: record.last_activity_ip }, record.last_activity_ip)
            : null,
        ]),
    },
    {
      key: 'status',
      title: $t('admin.users.columns.status'),
      width: 80,
      align: 'center',
      render: (value, record) =>
        h(
          'span',
          {
            class: 'inline-block whitespace-nowrap rounded px-2 py-1 text-xs font-medium',
            style:
              record.status === 1
                ? {
                    background: 'rgba(var(--color-success-rgb), 0.2)',
                    color: 'var(--color-success-400)',
                  }
                : record.status === 2
                  ? {
                      background: 'rgba(var(--color-error-rgb), 0.2)',
                      color: 'var(--color-error-400)',
                    }
                  : {},
          },
          getStatusText(record.status)
        ),
    },
    {
      key: 'role',
      title: $t('admin.users.columns.role'),
      width: 100,
      align: 'center',
      render: (value, record) =>
        h(
          'span',
          {
            class: 'inline-block whitespace-nowrap text-xs font-medium',
            style: {
              color:
                record.role === 1
                  ? 'var(--color-error-500)'
                  : record.role === 2
                    ? 'var(--color-brand-500)'
                    : 'var(--color-content-muted)',
            },
          },
          getRoleText(record.role)
        ),
    },
    {
      key: 'storage',
      title: $t('admin.users.columns.storage'),
      width: 180,
      align: 'center',
      render: (value, record) =>
        h('div', { class: 'storage-cell' }, [
          h('div', { class: 'storage-item' }, [
            h('span', { class: 'storage-label' }, `${$t('admin.users.columns.storageLabel')}:`),
            h(
              'span',
              { class: 'storage-value' },
              `${formatStorageSize(record.used_storage || 0)}/${formatStorageSize(record.storage_limit || 0)}`
            ),
          ]),
          h('div', { class: 'storage-item' }, [
            h('span', { class: 'storage-label' }, `${$t('admin.users.columns.bandwidthLabel')}:`),
            h(
              'span',
              { class: 'storage-value' },
              `${formatStorageSize(record.used_bandwidth || 0)}/${formatStorageSize(record.bandwidth_limit || 0)}`
            ),
          ]),
        ]),
    },
    {
      key: 'actions',
      title: $t('admin.users.columns.actions'),
      width: 120,
      align: 'center',
      render: (value, record) =>
        h('div', { class: 'action-buttons flex gap-1 justify-center items-center' }, [
          h(
            'button',
            {
              class: 'inline-flex items-center justify-center w-7 h-7 text-xs border rounded transition-colors',
              style: {
                background: 'rgba(var(--color-brand-500-rgb), 0.2)',
                color: 'var(--color-brand-500)',
                borderColor: 'rgba(var(--color-brand-500-rgb), 0.3)',
              },
              title: $t('admin.users.actions.view'),
              onClick: () => viewUserDetail(record.id),
              onMouseenter: (e: MouseEvent) => {
                ;(e.target as HTMLElement).style.background = 'rgba(var(--color-brand-500-rgb), 0.3)'
              },
              onMouseleave: (e: MouseEvent) => {
                ;(e.target as HTMLElement).style.background = 'rgba(var(--color-brand-500-rgb), 0.2)'
              },
            },
            [h('i', { class: 'fas fa-eye' })]
          ),
          h(
            'button',
            {
              class:
                'inline-flex items-center justify-center w-7 h-7 text-xs border rounded transition-colors focus:outline-none',
              style: {
                background: 'rgba(var(--color-error-rgb), 0.2)',
                color: 'var(--color-error-500)',
                borderColor: 'rgba(var(--color-error-rgb), 0.4)',
              },
              title: $t('admin.users.actions.edit'),
              onClick: () => editUser(record),
              onMouseenter: (e: MouseEvent) => {
                ;(e.target as HTMLElement).style.background = 'rgba(var(--color-error-rgb), 0.3)'
              },
              onMouseleave: (e: MouseEvent) => {
                ;(e.target as HTMLElement).style.background = 'rgba(var(--color-error-rgb), 0.2)'
              },
            },
            [h('i', { class: 'fas fa-edit' })]
          ),
          h(
            'button',
            {
              class: 'inline-flex items-center justify-center w-7 h-7 text-xs border rounded transition-colors',
              style:
                record.status === 1
                  ? {
                      background: 'rgba(var(--color-error-rgb), 0.2)',
                      color: 'var(--color-error-400)',
                      borderColor: 'rgba(var(--color-error-rgb), 0.3)',
                    }
                  : {
                      background: 'rgba(var(--color-success-rgb), 0.2)',
                      color: 'var(--color-success-400)',
                      borderColor: 'rgba(var(--color-success-rgb), 0.3)',
                    },
              title: record.status === 1 ? $t('admin.users.actions.disable') : $t('admin.users.actions.enable'),
              onClick: () => {
                if (record.status === 1) {
                  userToDisable.value = record
                  showDisableConfirmDialog.value = true
                } else {
                  if (showConfirm($t('admin.users.confirm.enableUser'))) {
                    toggleUserStatus(record)
                  }
                }
              },
              onMouseenter: (e: MouseEvent) => {
                if (record.status === 1) {
                  ;(e.target as HTMLElement).style.background = 'rgba(var(--color-error-rgb), 0.3)'
                } else {
                  ;(e.target as HTMLElement).style.background = 'rgba(var(--color-success-rgb), 0.3)'
                }
              },
              onMouseleave: (e: MouseEvent) => {
                if (record.status === 1) {
                  ;(e.target as HTMLElement).style.background = 'rgba(var(--color-error-rgb), 0.2)'
                } else {
                  ;(e.target as HTMLElement).style.background = 'rgba(var(--color-success-rgb), 0.2)'
                }
              },
            },
            [h('i', { class: record.status === 1 ? 'fas fa-ban' : 'fas fa-check' })]
          ),
        ]),
    },
  ])

  const users = ref<UserItem[]>([])
  const loading = ref(false)
  const saveLoading = ref(false)

  const totalUsers = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  const showFilter = ref(false)
  const currentFilters = reactive({
    keyword: '',
    status: '0',
    role: '0',
    sort_by: '',
    start_date: '',
    end_date: '',
  })

  const showEditDialog = ref(false)
  const currentEditUser = ref<UserItem | null>(null)

  const showCreateDialog = ref(false)

  const showDetailDrawer = ref(false)
  const selectedDetailUserId = ref<number | null>(null)
  const selectedDetailUserData = ref<UserItem | null>(null)

  const showDisableConfirmDialog = ref(false)
  const userToDisable = ref<UserItem | null>(null)

  const selectedUserIds = ref<number[]>([])

  const _isAllSelected = computed(() => users.value.length > 0 && selectedUserIds.value.length === users.value.length)

  const fetchUserList = async () => {
    loading.value = true

    const params = {
      page: currentPage.value,
      size: pageSize.value,
      keyword: currentFilters.keyword || undefined,
      status: currentFilters.status !== '0' ? parseInt(currentFilters.status) : undefined,
      role: currentFilters.role !== '0' ? parseInt(currentFilters.role) : undefined,
      start_date: currentFilters.start_date || undefined,
      end_date: currentFilters.end_date || undefined,
      sort_by: currentFilters.sort_by || undefined,
    }

    try {
      const result = await getUserList(params)
      if (result.success) {
        const { data } = result
        users.value = data.list
        totalUsers.value = data.total

        selectedUserIds.value = selectedUserIds.value.filter((id) => users.value.some((user) => user.id === id))
      }
    } catch {}

    loading.value = false
  }

  const handleFilter = (filters) => {
    Object.assign(currentFilters, filters)
    currentPage.value = 1
    selectedUserIds.value = [] // 清空选择
    fetchUserList()
    showFilter.value = false
  }

  const handlePageChange = (page) => {
    currentPage.value = page
    fetchUserList()
  }

  const handlePageSizeChange = (size) => {
    pageSize.value = size
    currentPage.value = 1 // 重置到第一页
    fetchUserList()
  }

  const editUser = (user) => {
    currentEditUser.value = JSON.parse(JSON.stringify(user))
    showEditDialog.value = true
  }

  const saveUserChanges = async () => {
    if (!currentEditUser.value) {
      return
    }

    saveLoading.value = true
    const userData = {
      id: currentEditUser.value.id,
      username: currentEditUser.value.username,
      status: currentEditUser.value.status,
      role: currentEditUser.value.role,
    }

    try {
      const result = await updateUser(userData)
      if (result.success) {
        const index = users.value.findIndex((u) => u.id === currentEditUser.value?.id)
        if (index !== -1) {
          users.value[index] = { ...currentEditUser.value }
        }

        toast.success($t('admin.users.toast.userUpdated'))
        showEditDialog.value = false

        fetchUserList()
      }
    } catch {}
    saveLoading.value = false
  }

  const toggleUserStatus = async (user) => {
    const newStatus = user.status === 1 ? 2 : 1

    try {
      const result = await toggleUserStatusApi(user.id, newStatus)

      if (result.success) {
        const index = users.value.findIndex((u) => u.id === user.id)
        if (index !== -1) {
          users.value[index].status = newStatus
        }

        toast.success(newStatus === 1 ? $t('admin.users.toast.userEnabled') : $t('admin.users.toast.userDisabled'))
      }
    } catch {}
  }

  const confirmDisableUser = () => {
    if (userToDisable.value) {
      toggleUserStatus(userToDisable.value)
      showDisableConfirmDialog.value = false
      userToDisable.value = null
    }
  }

  const cancelDisableUser = () => {
    showDisableConfirmDialog.value = false
    userToDisable.value = null
  }

  const getStatusText = (status) => {
    switch (status) {
      case 1:
        return $t('admin.users.statusText.normal')
      case 2:
        return $t('admin.users.statusText.disabled')
      case 3:
        return $t('admin.users.statusText.deleted')
      default:
        return $t('admin.users.statusText.unknown')
    }
  }

  const getRoleText = (role) => {
    switch (role) {
      case 1:
        return $t('admin.users.role.superAdmin')
      case 2:
        return $t('admin.users.role.admin')
      case 3:
        return $t('admin.users.role.user')
      default:
        return $t('admin.users.role.unknown')
    }
  }

  const viewUserDetail = (userId: number) => {
    const userData = users.value.find((user) => user.id === userId)
    if (userData) {
      selectedDetailUserId.value = userId
      selectedDetailUserData.value = userData
      showDetailDrawer.value = true
    } else {
      toast.error($t('admin.users.toast.userDataNotFound'))
    }
  }

  const handleRefresh = () => {
    selectedUserIds.value = []
    fetchUserList()
  }

  const handleUserCreated = (_newUser: UserItem) => {
    fetchUserList()
  }

  const handleUserUpdated = (updatedUser: UserItem) => {
    const index = users.value.findIndex((u) => u.id === updatedUser.id)
    if (index !== -1) {
      users.value[index] = { ...updatedUser }
    }
  }

  const handleBatchCompleted = () => {
    fetchUserList() // 重新获取列表
  }

  const _toggleUserSelection = (userId: number, checked: boolean) => {
    if (checked) {
      if (!selectedUserIds.value.includes(userId)) {
        selectedUserIds.value.push(userId)
      }
    } else {
      const index = selectedUserIds.value.indexOf(userId)
      if (index > -1) {
        selectedUserIds.value.splice(index, 1)
      }
    }
  }

  const _toggleSelectAll = (checked: boolean) => {
    if (checked) {
      selectedUserIds.value = users.value.map((user) => user.id)
    } else {
      selectedUserIds.value = []
    }
  }

  const clearSelection = () => {
    selectedUserIds.value = []
  }

  const handleRowClick = (_row: UserItem, _index: number) => {}

  const handleSelectionChange = (selectedRows: (string | number)[]) => {
    selectedUserIds.value = selectedRows.map((id) => Number(id))
  }

  const formatStorageSize = (bytes: number) => {
    if (bytes === 0) {
      return '0 B'
    }
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  onMounted(() => {
    fetchUserList()
  })
</script>

<template>
  <div class="admin-users-page admin-page-container">
    <CyberAdminWrapper
      :title="$t('admin.users.title')"
      :subtitle="$t('admin.users.subtitle')"
      icon="fas fa-users"
      :show-topbar="showFilter"
    >
      <template #actions>
        <CyberButton type="primary" icon="plus" @click="showCreateDialog = true">
          {{ $t('admin.users.buttons.createUser') }}
        </CyberButton>

        <CyberButton type="secondary" :loading="loading" loading-mode="inline" icon="sync-alt" @click="handleRefresh">
          {{ $t('admin.users.actions.refresh') }}
        </CyberButton>

        <CyberButton
          type="outlined"
          icon="filter"
          :class="{ 'border-error-500 text-error-500': showFilter }"
          @click="showFilter = !showFilter"
        >
          {{ $t('admin.users.buttons.filter') }}
        </CyberButton>
      </template>

      <template #topbar>
        <UserFilterPanel :initial-filters="currentFilters" @filter="handleFilter" />
      </template>

      <template #content>
        <div class="admin-content-wrapper">
          <BatchOperations
            :selected-users="selectedUserIds"
            :all-users="users"
            @clear-selection="clearSelection"
            @batch-completed="handleBatchCompleted"
          />

          <div class="table-container">
            <div class="table-wrapper">
              <div v-if="loading" class="admin-loading">
                <CyberLoading :visible="true" :text="$t('admin.users.loading.users')" />
              </div>

              <div v-else-if="!loading && users.length === 0" class="empty-state">
                <div class="empty-content">
                  <div class="empty-icon-wrapper">
                    <i class="fas fa-users empty-icon" />
                  </div>
                  <h3 class="empty-title">{{ $t('admin.users.empty.title') }}</h3>
                  <p class="empty-desc">{{ $t('admin.users.empty.description') }}</p>
                  <CyberButton type="secondary" @click="handleRefresh">
                    <i class="fas fa-redo mr-2" />
                    {{ $t('admin.users.empty.resetFilter') }}
                  </CyberButton>
                </div>
              </div>

              <CyberDataTable
                v-else
                :data="users"
                :columns="dataTableColumns"
                :loading="loading"
                :selectable="true"
                :hoverable="true"
                :striped="true"
                :bordered="true"
                size="small"
                row-key="id"
                :selected-row-keys="selectedUserIds"
                :loading-text="$t('admin.users.loading.users')"
                @update:selected-row-keys="handleSelectionChange"
                @row-click="handleRowClick"
              />
            </div>

            <div v-if="totalUsers > 0" class="pagination-wrapper">
              <CyberPagination
                v-model:current-page="currentPage"
                :total="totalUsers"
                :page-size="pageSize"
                :show-page-size-selector="true"
                :show-quick-jumper="true"
                class="compact-pagination"
                @update:current-page="handlePageChange"
                @update:page-size="handlePageSizeChange"
              />
            </div>
          </div>
        </div>
      </template>
    </CyberAdminWrapper>

    <CyberDialog v-model="showEditDialog" :title="$t('admin.users.editDialog.title')" width="500px">
      <UserForm v-model="currentEditUser" />

      <template #footer>
        <div class="flex justify-end gap-2 p-4">
          <CyberButton type="outlined" @click="showEditDialog = false">{{ $t('admin.users.buttons.cancel') }}</CyberButton>
          <CyberButton type="primary" :loading="saveLoading" @click="saveUserChanges">{{
            $t('admin.users.editDialog.saveChanges')
          }}</CyberButton>
        </div>
      </template>
    </CyberDialog>

    <CyberDialog v-model="showDisableConfirmDialog" :title="$t('admin.users.disableWarning.title')" width="500px">
      <div class="disable-confirm-content">
        <div class="warning-icon">
          <i class="fas fa-exclamation-triangle" />
        </div>
        <div class="warning-message">
          <h3 class="warning-title">{{ $t('admin.users.disableWarning.warning', { username: userToDisable?.username }) }}</h3>
          <div class="warning-text">
            <p class="warning-item">
              <i class="fas fa-info-circle" />
              {{ $t('admin.users.disableWarning.consequence1') }}
            </p>
            <p class="warning-item">
              <i class="fas fa-sign-out-alt" />
              {{ $t('admin.users.disableWarning.consequence2') }}
            </p>
            <p class="warning-item">
              <i class="fas fa-lock" />
              {{ $t('admin.users.disableWarning.consequence3') }}
            </p>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2 p-4">
          <CyberButton type="outlined" @click="cancelDisableUser">{{ $t('admin.users.buttons.cancel') }}</CyberButton>
          <CyberButton type="danger" @click="confirmDisableUser">
            <i class="fas fa-ban mr-2" />
            {{ $t('admin.users.disableWarning.confirmButton') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>

    <CreateUserDialog v-model="showCreateDialog" @user-created="handleUserCreated" />

    <UserDetailDrawer
      v-model="showDetailDrawer"
      :user-id="selectedDetailUserId"
      :user-data="selectedDetailUserData"
      @user-updated="handleUserUpdated"
    />
  </div>
</template>

<style scoped lang="scss">
  .admin-users-page {
    color: var(--color-content);
  }

  .admin-content-wrapper {
    background: var(--color-background-900);
    padding: 0;
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;
  }

  .users-content {
    display: flex;
    flex-direction: column;
    gap: var(--space-xl);
  }

  .search-section {
    background: rgba(var(--color-background-800-rgb), 0.4);
    border-radius: var(--radius-sm);
    padding: var(--space-xl);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .table-container {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }

  .table-container:has(.cyber-table) {
    border-radius: var(--radius-sm);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
    overflow: hidden;
  }

  .table-wrapper {
    flex: 1;
    min-height: 0;
    overflow: hidden;
  }

  .table-wrapper:has(.cyber-table) {
    background: rgba(var(--color-background-800-rgb), 0.3);
  }

  .pagination-wrapper {
    flex-shrink: 0;
    display: flex;
    justify-content: center;
    padding: var(--space-sm) var(--space-md);
    background: rgba(var(--color-background-800-rgb), 0.25);
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.08);
    margin: 0;
  }

  .compact-pagination {
    :deep(.pagination-wrapper) {
      padding: var(--space-sm) var(--space-sm);
      min-height: auto;
    }

    :deep(.pagination-controls) {
      gap: var(--space-xs);
    }

    :deep(.pagination-btn) {
      width: var(--space-3xl);
      height: var(--space-3xl);
      font-size: var(--text-xs);
      min-width: var(--space-3xl);
    }

    :deep(.page-size-selector) {
      gap: var(--space-xs);
    }

    :deep(.select-trigger) {
      min-height: var(--space-3xl);
      padding: 0 var(--space-sm);
      font-size: var(--text-xs);
    }

    :deep(.jump-input) {
      width: 48px;
      height: var(--space-3xl);
      font-size: var(--text-xs);
    }

    :deep(.page-size-label),
    :deep(.jump-label),
    :deep(.total-info) {
      font-size: var(--text-xs);
    }
  }

  @media (max-width: 1024px) {
    .header-content {
      flex-direction: column;
      align-items: stretch;
      gap: var(--space-md);
    }

    .title-section {
      justify-content: space-between;
    }

    .action-buttons {
      justify-content: center;
    }
  }

  @media (max-width: 768px) {
    .header-content {
      padding: var(--space-md) var(--space-xl);
    }

    .users-content {
      padding: var(--space-md) var(--space-xl);
    }

    .pagination-wrapper {
      padding: var(--space-sm);
    }

    .title-section {
      flex-direction: column;
      align-items: flex-start;
      gap: var(--space-sm);
    }

    .page-title {
      font-size: var(--text-lg);
    }

    .action-buttons {
      flex-direction: column;
      width: 100%;
    }

    .search-section {
      padding: var(--space-md);
    }
  }

  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100%;
    height: 100%;
    padding: var(--space-4xl) var(--space-2xl);
  }

  .empty-content {
    text-align: center;
    max-width: 400px;
    animation: fadeInUp var(--transition-slow) var(--ease-out);
  }

  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .empty-icon-wrapper {
    margin-bottom: var(--space-xl);
    position: relative;
    display: inline-block;
  }

  .empty-icon {
    font-size: var(--text-4xl);
    color: var(--color-brand-500);
    filter: drop-shadow(0 0 20px rgba(var(--color-brand-500-rgb), 0.4));
    animation: floatIcon 3s var(--ease-in-out) infinite;
  }

  @keyframes floatIcon {
    0%,
    100% {
      transform: translateY(0px);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  .empty-title {
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    color: var(--color-content-default);
    margin-bottom: var(--space-sm);
    letter-spacing: 0.025em;
  }

  .empty-desc {
    font-size: var(--text-sm);
    color: var(--color-content-muted);
    line-height: 1.6;
    margin-bottom: var(--space-xl);
  }

  .last-activity-cell {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
    min-width: 0;
  }

  .activity-time {
    font-size: var(--text-xs);
    color: var(--color-content-default);
    white-space: nowrap;
  }

  .no-activity {
    font-size: var(--text-xs);
    color: var(--color-content-subtle);
    white-space: nowrap;
  }

  .activity-ip {
    font-size: 0.65rem;
    color: var(--color-content-muted);
    font-family: 'Courier New', monospace;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 120px;
  }

  .action-buttons {
    display: flex;
    gap: var(--space-xs);
    justify-content: center;
    align-items: center;
  }

  .storage-cell {
    display: flex;
    flex-direction: column;
    gap: var(--space-xs);
    font-size: var(--text-xs);
    min-width: 0;
  }

  .storage-item {
    display: flex;
    align-items: center;
    gap: var(--space-xs);
    white-space: nowrap;
  }

  .storage-label {
    color: var(--color-content-muted);
    font-weight: var(--font-medium);
    flex-shrink: 0;
  }

  .storage-value {
    color: var(--color-content-default);
    font-family: 'Courier New', monospace;
    font-size: 0.7rem;
  }

  :deep(.form-container) {
    background: rgba(var(--color-background-800-rgb), 0.4);
    border-radius: var(--radius-sm);
    padding: var(--space-xl);
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.1);
  }

  .disable-confirm-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: var(--space-2xl) var(--space-xl);
    gap: var(--space-xl);
  }

  .warning-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 64px;
    height: 64px;
    border-radius: var(--radius-full);
    background: linear-gradient(135deg, rgba(var(--color-error-rgb), 0.2) 0%, rgba(var(--color-error-rgb), 0.1) 100%);
    border: 2px solid rgba(var(--color-error-rgb), 0.3);

    i {
      font-size: var(--text-2xl);
      color: var(--color-error-400);
      animation: pulse-warning 2s var(--ease-in-out) infinite;
    }
  }

  @keyframes pulse-warning {
    0%,
    100% {
      opacity: 1;
      transform: scale(1);
    }
    50% {
      opacity: 0.8;
      transform: scale(1.05);
    }
  }

  .warning-message {
    width: 100%;
    text-align: center;
  }

  .warning-title {
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    color: var(--color-content-default);
    margin-bottom: var(--space-xl);
    line-height: 1.5;
  }

  .warning-text {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    text-align: left;
    background: rgba(var(--color-error-rgb), 0.05);
    border: 1px solid rgba(var(--color-error-rgb), 0.15);
    border-radius: var(--radius-sm);
    padding: var(--space-lg);
  }

  .warning-item {
    display: flex;
    align-items: flex-start;
    gap: var(--space-sm);
    color: var(--color-content-muted);
    font-size: var(--text-sm);
    line-height: 1.6;
    margin: 0;

    i {
      color: var(--color-error-400);
      margin-top: var(--space-xs);
      flex-shrink: 0;
    }

    strong {
      color: var(--color-content-default);
      font-weight: var(--font-semibold);
    }
  }
</style>
