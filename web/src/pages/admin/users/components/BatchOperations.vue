<script setup lang="ts">
  import { nextTick, onMounted, onUnmounted, ref } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { useToast } from '@/components/Toast/useToast'
  import type { UserItem } from '@/api/admin/user'

  const props = defineProps<{
    selectedUsers: number[]
    allUsers: UserItem[]
  }>()

  const emit = defineEmits<{
    clearSelection: []
    batchCompleted: []
  }>()

  const toast = useToast()
  const { $t } = useTexts()

  /* 状态 - 每个按钮独立的loading状态 */
  const enableLoading = ref(false)
  const disableLoading = ref(false)
  const roleLoading = ref(false)
  const deleteLoading = ref(false)
  const exportLoading = ref(false)
  const batchRole = ref('')
  const showRoleSelector = ref(false)
  const selectedRoles = ref<string[]>([])
  const roleButtonRef = ref<HTMLElement | null>(null)
  const dropdownPosition = ref({ top: 0, left: 0, width: 0 })

  /* 批量禁用确认弹窗 */
  const showBatchDisableDialog = ref(false)

  /* 角色选项 */
  const roleOptions = [
    { label: $t('admin.users.roles.normal'), value: '3' },
    // { label: $t('admin.users.roles.admin'), value: '2' },
    { label: $t('admin.users.roles.superAdmin'), value: '1' },
  ]

  const updateDropdownPosition = () => {
    if (roleButtonRef.value) {
      const rect = roleButtonRef.value.getBoundingClientRect()
      dropdownPosition.value = {
        top: rect.bottom + 4,
        left: rect.left,
        width: rect.width,
      }
    }
  }

  const toggleRoleSelector = () => {
    showRoleSelector.value = !showRoleSelector.value
    if (showRoleSelector.value) {
      nextTick(() => {
        updateDropdownPosition()
      })
    }
  }

  const handlePositionUpdate = () => {
    if (showRoleSelector.value) {
      updateDropdownPosition()
    }
  }

  onMounted(() => {
    window.addEventListener('scroll', handlePositionUpdate, true)
    window.addEventListener('resize', handlePositionUpdate)
  })

  onUnmounted(() => {
    window.removeEventListener('scroll', handlePositionUpdate, true)
    window.removeEventListener('resize', handlePositionUpdate)
  })

  const clearSelection = () => {
    emit('clearSelection')
  }

  const handleBatchEnable = async () => {
    await performBatchOperation('enable', $t('admin.users.batch.enable'), undefined, enableLoading)
  }

  const handleBatchDisable = () => {
    if (props.selectedUsers.length === 0) {
      toast.warning($t('admin.users.batch.noSelection'))
      return
    }
    showBatchDisableDialog.value = true
  }

  const confirmBatchDisable = async () => {
    showBatchDisableDialog.value = false
    await performBatchOperation('disable', $t('admin.users.batch.disable'), undefined, disableLoading)
  }

  const cancelBatchDisable = () => {
    showBatchDisableDialog.value = false
  }

  const toggleRoleSelection = (value: string, checked: boolean) => {
    if (checked) {
      if (!selectedRoles.value.includes(value)) {
        selectedRoles.value.push(value)
      }
    } else {
      const index = selectedRoles.value.indexOf(value)
      if (index > -1) {
        selectedRoles.value.splice(index, 1)
      }
    }
  }

  const applyRoleSelection = async () => {
    if (selectedRoles.value.length === 0) {
      toast.warning($t('admin.users.batch.selectRole'))
      return
    }

    const roleToSet = selectedRoles.value[0]
    await performBatchOperation('set_role', $t('admin.users.batch.setRole'), parseInt(roleToSet), roleLoading)

    selectedRoles.value = []
    showRoleSelector.value = false
  }

  const cancelRoleSelection = () => {
    selectedRoles.value = []
    showRoleSelector.value = false
  }

  const _handleBatchSetRole = async () => {
    if (!batchRole.value) {
      toast.warning($t('admin.users.batch.selectRole'))
      return
    }

    await performBatchOperation('set_role', $t('admin.users.batch.setRole'), parseInt(batchRole.value), roleLoading)
    batchRole.value = ''
  }

  const handleBatchDelete = async () => {
    await performBatchOperation('delete', $t('admin.users.batch.delete'), undefined, deleteLoading)
  }

  const performBatchOperation = async (operation: string, operationName: string, role?: number, loadingRef?: any) => {
    if (props.selectedUsers.length === 0) {
      toast.warning($t('admin.users.batch.noSelection'))
      return
    }

    if (loadingRef) {
      loadingRef.value = true
    }

    try {
      const { batchOperateUsers } = await import('@/api/admin/user')

      const params: any = {
        user_ids: props.selectedUsers,
        operation: operation as any,
      }

      if (operation === 'set_role' && role) {
        params.role = role
      }

      const result = await batchOperateUsers(params)

      if (result.success) {
        toast.success($t('admin.users.toast.success', { operation: operationName }))
        emit('batchCompleted')
        emit('clearSelection')
      } else {
        toast.error(result.message || $t('admin.users.toast.error', { operation: operationName }))
      }
    } catch (error: any) {
      toast.error(error?.message || $t('admin.users.toast.error', { operation: operationName }))
    }

    if (loadingRef) {
      loadingRef.value = false
    }
  }

  const handleExport = async () => {
    if (props.selectedUsers.length === 0) {
      toast.warning($t('admin.users.batch.selectUsersToExport'))
      return
    }

    exportLoading.value = true
    try {
      const selectedUserData = props.allUsers.filter((user) => props.selectedUsers.includes(user.id))

      const exportData = selectedUserData.map((user) => ({
        [$t('admin.users.export.id')]: user.id,
        [$t('admin.users.export.username')]: user.username,
        [$t('admin.users.export.email')]: user.email,
        [$t('admin.users.export.status')]: getStatusText(user.status),
        [$t('admin.users.export.role')]: getRoleText(user.role),
        [$t('admin.users.export.createdAt')]: user.created_at,
        [$t('admin.users.export.updatedAt')]: user.updated_at,
      }))

      const csvContent = generateCSV(exportData)

      downloadCSV(csvContent, `${$t('admin.users.export.filename')}_${new Date().toISOString().split('T')[0]}.csv`)

      toast.success($t('admin.users.toast.exportSuccess'))
    } catch {
      toast.error($t('admin.users.toast.exportError'))
    }
    exportLoading.value = false
  }

  const generateCSV = (data: Record<string, unknown>[]) => {
    if (data.length === 0) {
      return ''
    }

    const headers = Object.keys(data[0])
    const csvRows = []

    csvRows.push(headers.join(','))

    for (const row of data) {
      const values = headers.map((header) => {
        const value = row[header]
        if (typeof value === 'string' && (value.includes(',') || value.includes('"'))) {
          return `"${value.replace(/"/g, '""')}"`
        }
        return value
      })
      csvRows.push(values.join(','))
    }

    return csvRows.join('\n')
  }

  const downloadCSV = (csvContent: string, filename: string) => {
    const blob = new Blob([`\uFEFF${csvContent}`], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')

    if (link.download !== undefined) {
      const url = URL.createObjectURL(blob)
      link.setAttribute('href', url)
      link.setAttribute('download', filename)
      link.style.visibility = 'hidden'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    }
  }

  const getStatusText = (status: number) => {
    switch (status) {
      case 1:
        return $t('admin.users.status.active')
      case 2:
        return $t('admin.users.status.disabled')
      case 3:
        return $t('admin.users.status.deleted')
      default:
        return $t('admin.users.status.unknown')
    }
  }

  const getRoleText = (role: number) => {
    switch (role) {
      case 1:
        return $t('admin.users.roles.superAdmin')
      case 2:
        return $t('admin.users.roles.admin')
      case 3:
        return $t('admin.users.roles.normal')
      default:
        return $t('admin.users.roles.unknown')
    }
  }
</script>

<template>
  <div v-if="selectedUsers.length > 0" class="batch-operations">
    <div class="batch-operations-container">
      <div class="flex items-center space-x-3">
        <span class="text-sm font-medium text-content">
          {{ $t('admin.users.batch.selected', { count: selectedUsers.length }) }}
        </span>
        <CyberButton type="ghost" class="text-xs" @click="clearSelection">
          <i class="fas fa-times mr-1" />
          {{ $t('admin.users.batch.clearSelection') }}
        </CyberButton>
      </div>

      <div class="flex flex-wrap items-center space-x-1.5">
        <CyberButton type="primary" :loading="enableLoading" @click="handleBatchEnable">
          <i class="fas fa-check mr-1" />
          {{ $t('admin.users.batch.enable') }}
        </CyberButton>

        <CyberButton type="danger" :loading="disableLoading" @click="handleBatchDisable">
          <i class="fas fa-ban mr-1" />
          {{ $t('admin.users.batch.disable') }}
        </CyberButton>

        <div ref="roleButtonRef" class="role-selector-dropdown">
          <CyberButton type="secondary" :loading="roleLoading" :class="{ active: showRoleSelector }" @click="toggleRoleSelector">
            <i class="fas fa-user-cog mr-1" />
            {{ $t('admin.users.batch.role') }}
            <i class="fas fa-chevron-down ml-1" :class="{ 'rotate-180': showRoleSelector }" />
          </CyberButton>
        </div>

        <CyberPopconfirm :title="$t('admin.users.confirm.delete')" @confirm="handleBatchDelete">
          <CyberButton type="danger" :loading="deleteLoading">
            <i class="fas fa-trash mr-1" />
            {{ $t('admin.users.batch.delete') }}
          </CyberButton>
        </CyberPopconfirm>

        <CyberButton type="secondary" :loading="exportLoading" @click="handleExport">
          <i class="fas fa-download mr-1" />
          {{ $t('admin.users.batch.export') }}
        </CyberButton>
      </div>
    </div>

    <Teleport to="body">
      <div
        v-if="showRoleSelector"
        class="role-selector-panel"
        :style="{
          top: `${dropdownPosition.top}px`,
          left: `${dropdownPosition.left}px`,
          minWidth: `${dropdownPosition.width}px`,
        }"
      >
        <div class="role-checkbox-group">
          <CyberCheckbox
            v-for="option in roleOptions"
            :key="option.value"
            :model-value="selectedRoles.includes(option.value)"
            class="role-checkbox-item"
            @update:model-value="(checked) => toggleRoleSelection(option.value, checked)"
          >
            {{ option.label }}
          </CyberCheckbox>

          <div class="role-actions">
            <CyberButton
              type="primary"
              size="small"
              :loading="roleLoading"
              :disabled="selectedRoles.length === 0"
              @click="applyRoleSelection"
            >
              {{ $t('admin.users.batch.apply') }}
            </CyberButton>
            <CyberButton type="secondary" size="small" @click="cancelRoleSelection">
              {{ $t('admin.users.batch.cancel') }}
            </CyberButton>
          </div>
        </div>
      </div>
    </Teleport>

    <CyberDialog v-model="showBatchDisableDialog" :title="$t('admin.users.confirm.batchDisable')" width="500px">
      <div class="disable-confirm-content">
        <div class="warning-icon">
          <i class="fas fa-exclamation-triangle" />
        </div>
        <div class="warning-message">
          <h3 class="warning-title">{{ $t('admin.users.confirm.disableTitle', { count: selectedUsers.length }) }}</h3>
          <div class="warning-text">
            <p class="warning-item">
              <i class="fas fa-info-circle" />
              {{ $t('admin.users.confirm.disableMessage1') }}
            </p>
            <p class="warning-item">
              <i class="fas fa-sign-out-alt" />
              {{ $t('admin.users.confirm.disableMessage2') }}
            </p>
            <p class="warning-item">
              <i class="fas fa-lock" />
              {{ $t('admin.users.confirm.disableMessage3') }}
            </p>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2 p-4">
          <CyberButton type="outlined" @click="cancelBatchDisable">{{ $t('admin.users.batch.cancel') }}</CyberButton>
          <CyberButton type="danger" :loading="disableLoading" @click="confirmBatchDisable">
            <i class="fas fa-ban mr-2" />
            {{ $t('admin.users.batch.confirmDisable') }}
          </CyberButton>
        </div>
      </template>
    </CyberDialog>
  </div>

  <template v-else>
    <div style="display: none"></div>
  </template>
</template>

<style scoped lang="scss">
  .batch-operations {
    flex-shrink: 0;
    margin-bottom: 1rem;
  }

  .batch-operations-container {
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.9) 0%,
      rgba(var(--color-background-900-rgb), 0.8) 50%,
      rgba(var(--color-brand-500-rgb), 0.05) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.3);
    border-radius: var(--radius-sm);
    padding: 0.75rem 1.25rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    box-shadow:
      0 4px 20px rgba(0, 0, 0, 0.3),
      0 0 40px rgba(var(--color-brand-500-rgb), 0.1);
    backdrop-filter: blur(10px);
  }

  @media (max-width: 1200px) {
    .batch-operations-container {
      padding: 0.875rem 1.25rem;
    }

    .batch-operations .flex.items-center {
      gap: 0.25rem;
    }

    .batch-operations :deep(.cyber-button) {
      padding: 0.375rem 0.5rem;
      font-size: 0.8rem;
    }
  }

  .role-selector-dropdown {
    position: relative;
  }

  .role-selector-dropdown .cyber-button.active {
    background-color: rgba(var(--color-brand-500-rgb), 0.2);
    border-color: rgba(var(--color-brand-500-rgb), 0.5);
  }

  .role-selector-dropdown .fa-chevron-down {
    transition: transform 0.2s ease;
  }

  .role-selector-dropdown .rotate-180 {
    transform: rotate(180deg);
  }

  .role-selector-panel {
    position: fixed;
    z-index: 9999;
    background: linear-gradient(
      135deg,
      rgba(var(--color-background-800-rgb), 0.95) 0%,
      rgba(var(--color-background-800-rgb), 0.9) 100%
    );
    border: 1px solid rgba(var(--color-brand-500-rgb), 0.4);
    border-radius: var(--radius-sm);
    padding: 12px;
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.4),
      0 0 40px rgba(var(--color-brand-500-rgb), 0.1);
    backdrop-filter: blur(12px);
  }

  .role-checkbox-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .role-checkbox-item {
    font-size: 0.85rem;
  }

  .role-actions {
    display: flex;
    gap: 6px;
    margin-top: 12px;
    padding-top: 8px;
    border-top: 1px solid rgba(var(--color-brand-500-rgb), 0.2);
  }

  @media (max-width: 768px) {
    .batch-operations-container {
      flex-direction: column;
      gap: 0.75rem;
      align-items: stretch;
    }

    .batch-operations .flex.items-center:last-child {
      justify-content: center;
      flex-wrap: wrap;
    }

    .role-selector-panel {
      position: fixed;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      min-width: 200px;
    }
  }

  .disable-confirm-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem 1.5rem;
    gap: 1.5rem;
  }

  .warning-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 64px;
    height: 64px;
    border-radius: var(--radius-full);
    background: linear-gradient(135deg, rgba(239, 68, 68, 0.2) 0%, rgba(220, 38, 38, 0.1) 100%);
    border: 2px solid rgba(239, 68, 68, 0.3);

    i {
      font-size: 2rem;
      color: rgb(239, 68, 68);
      animation: pulse-warning 2s ease-in-out infinite;
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
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--color-content);
    margin-bottom: 1.5rem;
    line-height: 1.5;
  }

  .warning-text {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    text-align: left;
    background: rgba(239, 68, 68, 0.05);
    border: 1px solid rgba(239, 68, 68, 0.15);
    border-radius: var(--radius-sm);
    padding: 1.25rem;
  }

  .warning-item {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    color: var(--color-content-muted);
    font-size: 0.875rem;
    line-height: 1.6;
    margin: 0;

    i {
      color: rgb(239, 68, 68);
      margin-top: 0.125rem;
      flex-shrink: 0;
    }

    strong {
      color: var(--color-content);
      font-weight: 600;
    }
  }
</style>
