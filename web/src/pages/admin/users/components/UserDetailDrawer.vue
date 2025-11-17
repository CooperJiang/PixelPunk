<script setup lang="ts">
  import { computed, reactive, ref, watch } from 'vue'
  import { useTexts } from '@/composables/useTexts'
  import { useToast } from '@/components/Toast/useToast'
  import { formatDate } from '@/utils/formatting/format'
  import {
    getUserDetail,
    resetUserPassword,
    sendUserEmail,
    toggleUserStatus,
    updateUserStorage,
    type UserItem,
  } from '@/api/admin/user'
  import SendEmailDialog from './SendEmailDialog.vue'

  const { $t } = useTexts()

  const props = defineProps<{
    modelValue: boolean
    userId: number | null
    userData?: UserItem | null
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    userUpdated: [user: UserItem]
  }>()

  const toast = useToast()

  const visible = ref(false)
  const userDetail = ref<UserItem | null>(null)
  const userStats = ref<Record<string, unknown>>({})
  const updateStorageLoading = ref(false)
  const resetPasswordLoading = ref(false)
  const sendEmailLoading = ref(false)
  const toggleStatusLoading = ref(false)
  const showEmailDialog = ref(false)

  /* 存储单位选项 */
  const storageUnitOptions = computed(() => [
    { label: $t('admin.users.detail.units.mb'), value: 'MB' },
    { label: $t('admin.users.detail.units.gb'), value: 'GB' },
    { label: $t('admin.users.detail.units.tb'), value: 'TB' },
  ])

  /* 带宽单位选项 */
  const bandwidthUnitOptions = computed(() => [
    { label: $t('admin.users.detail.units.mb'), value: 'MB' },
    { label: $t('admin.users.detail.units.gb'), value: 'GB' },
    { label: $t('admin.users.detail.units.tb'), value: 'TB' },
  ])

  /* 存储设置 */
  const storageSettings = reactive({
    storage_limit: 0,
    bandwidth_limit: 0,
    storage_unit: 'GB',
    bandwidth_unit: 'GB',
  })

  /* 监听props变化 */
  watch(
    () => props.modelValue,
    (newVal) => {
      visible.value = newVal
      if (newVal && (props.userData || props.userId)) {
        initUserDetail()
      }
    },
    { immediate: true }
  )

  watch(
    () => props.userData,
    (newVal) => {
      if (newVal && visible.value) {
        initUserDetail()
      }
    }
  )

  watch(visible, (newVal) => {
    emit('update:modelValue', newVal)
  })

  const convertBytesToBestUnit = (bytes: number) => {
    if (bytes === 0) {
      return { value: 0, unit: 'GB' }
    }

    const tb = bytes / (1024 * 1024 * 1024 * 1024)
    const gb = bytes / (1024 * 1024 * 1024)
    const mb = bytes / (1024 * 1024)

    if (tb >= 1) {
      return { value: Math.round(tb * 100) / 100, unit: 'TB' }
    } else if (gb >= 1) {
      return { value: Math.round(gb * 100) / 100, unit: 'GB' }
    }
    return { value: Math.round(mb * 100) / 100, unit: 'MB' }
  }

  const convertToBytes = (value: number, unit: string): number => {
    switch (unit) {
      case 'TB':
        return value * 1024 * 1024 * 1024 * 1024
      case 'GB':
        return value * 1024 * 1024 * 1024
      case 'MB':
        return value * 1024 * 1024
      default:
        return value
    }
  }

  const initUserDetail = async () => {
    if (props.userData) {
      userDetail.value = props.userData

      const storageConverted = convertBytesToBestUnit(props.userData.storage_limit || 5 * 1024 * 1024 * 1024)
      const bandwidthConverted = convertBytesToBestUnit(props.userData.bandwidth_limit || 100 * 1024 * 1024 * 1024)

      storageSettings.storage_limit = storageConverted.value
      storageSettings.storage_unit = storageConverted.unit
      storageSettings.bandwidth_limit = bandwidthConverted.value
      storageSettings.bandwidth_unit = bandwidthConverted.unit

      userStats.value = {
        total_images: props.userData.total_images || 0,
        used_storage: props.userData.used_storage || 0,
        monthly_bandwidth: props.userData.used_bandwidth || 0,
        total_views: props.userData.total_views || 0,
        total_shares: 0, // 这个字段后端还没有
        last_login: $t('admin.users.detail.noData'),
      }
    } else if (props.userId) {
      try {
        const result = await getUserDetail(props.userId)
        if (result.success) {
          const userDetailData = result.data
          userDetail.value = userDetailData

          const storageConverted = convertBytesToBestUnit(userDetailData.storage_limit || 5 * 1024 * 1024 * 1024)
          const bandwidthConverted = convertBytesToBestUnit(userDetailData.bandwidth_limit || 100 * 1024 * 1024 * 1024)

          storageSettings.storage_limit = storageConverted.value
          storageSettings.storage_unit = storageConverted.unit
          storageSettings.bandwidth_limit = bandwidthConverted.value
          storageSettings.bandwidth_unit = bandwidthConverted.unit

          userStats.value = {
            total_images: userDetailData.total_images || 0,
            used_storage: userDetailData.used_storage || 0,
            monthly_bandwidth: userDetailData.used_bandwidth || 0,
            total_views: userDetailData.total_views || 0,
            total_shares: 0, // 这个字段后端还没有
            last_login: $t('admin.users.detail.noData'),
          }
        }
      } catch (error: unknown) {
        const errorMessage =
          (error as any)?.response?.data?.message || (error as any)?.message || $t('admin.users.detail.errors.fetchFailed')
        toast.error(errorMessage)
      }
    }
  }

  const updateStorageSettings = async () => {
    if (!userDetail.value) {
      return
    }

    updateStorageLoading.value = true
    try {
      const storageBytes = convertToBytes(Number(storageSettings.storage_limit), storageSettings.storage_unit)
      const bandwidthBytes = convertToBytes(Number(storageSettings.bandwidth_limit), storageSettings.bandwidth_unit)

      await updateUserStorage({
        user_id: userDetail.value.id,
        storage_limit: storageBytes,
        bandwidth_limit: bandwidthBytes,
      })

      toast.success($t('admin.users.detail.messages.storageUpdateSuccess'))

      if (userDetail.value) {
        userDetail.value.storage_limit = storageBytes
        userDetail.value.bandwidth_limit = bandwidthBytes
      }

      emit('userUpdated', userDetail.value)
    } catch (error: unknown) {
      const errorMessage =
        (error as any)?.response?.data?.message || (error as any)?.message || $t('admin.users.detail.errors.storageUpdateFailed')
      toast.error(errorMessage)
    }
    updateStorageLoading.value = false
  }

  const resetPassword = async () => {
    if (!userDetail.value) {
      return
    }

    resetPasswordLoading.value = true
    try {
      const result = await resetUserPassword(userDetail.value.id)
      toast.success($t('admin.users.detail.messages.passwordResetSuccess', { password: result.new_password }))
    } catch (error: unknown) {
      const errorMessage =
        (error as any)?.response?.data?.message || (error as any)?.message || $t('admin.users.detail.errors.passwordResetFailed')
      toast.error(errorMessage)
    }
    resetPasswordLoading.value = false
  }

  const handleEmailSent = async (data: { userId: number; subject: string; content: string }) => {
    sendEmailLoading.value = true
    try {
      await sendUserEmail({
        user_id: data.userId,
        subject: data.subject,
        content: data.content,
      })

      toast.success($t('admin.users.detail.messages.emailSentSuccess'))
    } catch (error: unknown) {
      const errorMessage =
        (error as any)?.response?.data?.message || (error as any)?.message || $t('admin.users.detail.errors.emailSentFailed')
      toast.error(errorMessage)
    }
    sendEmailLoading.value = false
  }

  const toggleStatus = async () => {
    if (!userDetail.value) {
      return
    }

    toggleStatusLoading.value = true
    try {
      const newStatus = userDetail.value.status === 1 ? 2 : 1

      await toggleUserStatus(userDetail.value.id, newStatus)

      userDetail.value.status = newStatus
      toast.success($t(`admin.users.detail.messages.${newStatus === 1 ? 'userEnabled' : 'userDisabled'}`))

      emit('userUpdated', userDetail.value)
    } catch (error: unknown) {
      const errorMessage =
        (error as any)?.response?.data?.message || (error as any)?.message || $t('admin.users.detail.errors.operationFailed')
      toast.error(errorMessage)
    }
    toggleStatusLoading.value = false
  }

  const handleClose = () => {
    visible.value = false
    userDetail.value = null
    userStats.value = {}
  }

  const getStatusText = (status: number) => {
    switch (status) {
      case 1:
        return $t('admin.users.detail.status.active')
      case 2:
        return $t('admin.users.detail.status.disabled')
      case 3:
        return $t('admin.users.detail.status.deleted')
      default:
        return $t('admin.users.detail.status.unknown')
    }
  }

  const getRoleText = (role: number) => {
    switch (role) {
      case 1:
        return $t('admin.users.detail.roles.superAdmin')
      case 2:
        return $t('admin.users.detail.roles.admin')
      case 3:
        return $t('admin.users.detail.roles.user')
      default:
        return $t('admin.users.detail.roles.unknown')
    }
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
</script>

<template>
  <div class="user-detail-drawer">
    <CyberDrawer v-model="visible" :title="$t('admin.users.detail.title')" width="600px" @close="handleClose">
      <div v-if="userDetail" class="space-y-4 p-4">
        <div class="animate-slide-up rounded-lg border border-subtle bg-background-800 p-4" style="animation-delay: 0.05s">
          <h3 class="mb-4 flex items-center text-lg font-semibold text-content">
            <i class="fas fa-user mr-2" />
            {{ $t('admin.users.detail.basicInfo') }}
          </h3>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div class="flex items-center space-x-3">
              <div class="relative">
                <CyberUserAvatar :avatar-url="userDetail.avatar_full_path" :username="userDetail.username" size="lg" />
                <div class="absolute -bottom-1 -right-1 h-3 w-3 rounded-full border border-default bg-green-500" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate font-medium text-content">{{ userDetail.username }}</p>
                <p class="truncate text-sm text-content-muted">{{ userDetail.email }}</p>
              </div>
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span class="whitespace-nowrap text-content-muted">{{ $t('admin.users.detail.userId') }}:</span>
                <span class="font-mono text-content">#{{ userDetail.id }}</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="whitespace-nowrap text-content-muted">{{ $t('admin.users.detail.status.label') }}:</span>
                <span
                  class="rounded px-2 py-1 text-xs font-medium"
                  :class="[
                    userDetail.status === 1
                      ? 'bg-green-500/20 text-green-400'
                      : userDetail.status === 2
                        ? 'bg-red-500/20 text-red-400'
                        : 'bg-content/20 text-content-muted',
                  ]"
                >
                  {{ getStatusText(userDetail.status) }}
                </span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="whitespace-nowrap text-content-muted">{{ $t('admin.users.detail.role') }}:</span>
                <span
                  class="rounded px-2 py-1 text-xs font-medium"
                  :class="[
                    userDetail.role === 1
                      ? 'bg-purple-500/20 text-purple-400'
                      : userDetail.role === 2
                        ? 'bg-brand-500/20 text-brand-primary'
                        : 'bg-content/20 text-content-muted',
                  ]"
                >
                  {{ getRoleText(userDetail.role) }}
                </span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="whitespace-nowrap text-content-muted">{{ $t('admin.users.detail.registeredAt') }}:</span>
                <span class="truncate text-xs text-content">{{ formatDate(userDetail.created_at) }}</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="whitespace-nowrap text-content-muted">{{ $t('admin.users.detail.lastActivity') }}:</span>
                <span class="truncate text-xs text-content">{{
                  formatDate(userDetail.last_activity_at) || $t('admin.users.detail.never')
                }}</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="whitespace-nowrap text-content-muted">{{ $t('admin.users.detail.lastActivityIp') }}:</span>
                <span class="font-mono text-xs text-content">{{ userDetail.last_activity_ip || '--' }}</span>
              </div>
            </div>
          </div>
        </div>

        <div
          class="animate-slide-up rounded-xl border border-subtle bg-background-800 p-6 shadow-lg backdrop-blur-sm"
          style="animation-delay: 0.1s"
        >
          <h3 class="mb-6 flex items-center text-xl font-bold text-content">
            <div class="mr-3 flex h-8 w-8 items-center justify-center rounded-lg bg-green-500/20">
              <i class="fas fa-hdd text-green-400" />
            </div>
            {{ $t('admin.users.detail.storageManagement') }}
          </h3>
          <div class="space-y-6">
            <div class="grid grid-cols-1 gap-6">
              <div class="rounded-lg border border-subtle bg-background-700 p-4">
                <label class="mb-3 block flex items-center text-sm font-medium text-content">
                  <i class="fas fa-database mr-2 text-blue-400" />
                  {{ $t('admin.users.detail.storageLimit') }}
                </label>
                <div class="flex items-center gap-3">
                  <cyberInput
                    v-model="storageSettings.storage_limit"
                    type="number"
                    :min="0"
                    :max="99999"
                    :placeholder="$t('admin.users.detail.enterValue')"
                    class="flex-1"
                  />
                  <CyberDropdown
                    v-model="storageSettings.storage_unit"
                    :options="storageUnitOptions"
                    class="w-20"
                    :placeholder="$t('admin.users.detail.unit')"
                  />
                </div>
                <div class="mt-2 flex items-center text-xs text-content-muted">
                  <i class="fas fa-info-circle mr-1 text-blue-400" />
                  {{ $t('admin.users.detail.currentSetting') }}: {{ formatStorageSize(userDetail?.storage_limit || 0) }}
                </div>
              </div>

              <div class="rounded-lg border border-subtle bg-background-700 p-4">
                <label class="mb-3 block flex items-center text-sm font-medium text-content">
                  <i class="fas fa-wifi mr-2 text-purple-400" />
                  {{ $t('admin.users.detail.bandwidthLimit') }}
                </label>
                <div class="flex items-center gap-3">
                  <cyberInput
                    v-model="storageSettings.bandwidth_limit"
                    type="number"
                    :min="0"
                    :max="99999"
                    :placeholder="$t('admin.users.detail.enterValue')"
                    class="flex-1"
                  />
                  <CyberDropdown
                    v-model="storageSettings.bandwidth_unit"
                    :options="bandwidthUnitOptions"
                    class="w-20"
                    :placeholder="$t('admin.users.detail.unit')"
                  />
                </div>
                <div class="mt-2 flex items-center text-xs text-content-muted">
                  <i class="fas fa-info-circle mr-1 text-purple-400" />
                  {{ $t('admin.users.detail.currentSetting') }}: {{ formatStorageSize(userDetail?.bandwidth_limit || 0) }}
                </div>
              </div>
            </div>
            <div class="flex justify-end">
              <cyberButton type="primary" :loading="updateStorageLoading" @click="updateStorageSettings">
                {{ $t('admin.users.detail.updateStorageSettings') }}
              </cyberButton>
            </div>
          </div>
        </div>

        <div
          class="animate-slide-up rounded-xl border border-subtle bg-background-800 p-6 shadow-lg backdrop-blur-sm"
          style="animation-delay: 0.15s"
        >
          <h3 class="mb-6 flex items-center text-xl font-bold text-content">
            <div class="mr-3 flex h-8 w-8 items-center justify-center rounded-lg bg-orange-500/20">
              <i class="fas fa-chart-bar text-orange-400" />
            </div>
            {{ $t('admin.users.detail.usageStats') }}
          </h3>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-4">
              <div class="rounded-lg border border-subtle bg-background-700 p-3">
                <div class="flex items-center justify-between">
                  <span class="flex items-center text-sm text-content-muted">
                    <i class="fas fa-images mr-2 text-blue-400" />
                    {{ $t('admin.users.detail.uploadedFiles') }}
                  </span>
                  <span class="text-lg font-bold text-content">{{ userStats.total_images || 0 }}</span>
                </div>
              </div>

              <div class="rounded-lg border border-subtle bg-background-700 p-3">
                <div class="flex items-center justify-between">
                  <span class="flex items-center text-sm text-content-muted">
                    <i class="fas fa-hdd mr-2 text-green-400" />
                    {{ $t('admin.users.detail.usedStorage') }}
                  </span>
                  <span class="text-lg font-bold text-content">{{ formatStorageSize(userStats.used_storage || 0) }}</span>
                </div>
              </div>

              <div class="rounded-lg border border-subtle bg-background-700 p-3">
                <div class="flex items-center justify-between">
                  <span class="flex items-center text-sm text-content-muted">
                    <i class="fas fa-download mr-2 text-purple-400" />
                    {{ $t('admin.users.detail.monthlyBandwidth') }}
                  </span>
                  <span class="text-lg font-bold text-content">{{ formatStorageSize(userStats.monthly_bandwidth || 0) }}</span>
                </div>
              </div>
            </div>

            <div class="space-y-4">
              <div class="rounded-lg border border-subtle bg-background-700 p-3">
                <div class="flex items-center justify-between">
                  <span class="flex items-center text-sm text-content-muted">
                    <i class="fas fa-eye mr-2 text-yellow-400" />
                    {{ $t('admin.users.detail.totalViews') }}
                  </span>
                  <span class="text-lg font-bold text-content">{{ userStats.total_views || 0 }}</span>
                </div>
              </div>

              <div class="rounded-lg border border-subtle bg-background-700 p-3">
                <div class="flex items-center justify-between">
                  <span class="flex items-center text-sm text-content-muted">
                    <i class="fas fa-share mr-2 text-pink-400" />
                    {{ $t('admin.users.detail.shareCount') }}
                  </span>
                  <span class="text-lg font-bold text-content">{{ userStats.total_shares || 0 }}</span>
                </div>
              </div>

              <div class="rounded-lg border border-subtle bg-background-700 p-3">
                <div class="flex items-center justify-between">
                  <span class="flex items-center text-sm text-content-muted">
                    <i class="fas fa-clock mr-2 text-indigo-400" />
                    {{ $t('admin.users.detail.lastLogin') }}
                  </span>
                  <span class="text-sm font-medium text-content">{{
                    userStats.last_login || $t('admin.users.detail.never')
                  }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div
          class="animate-slide-up rounded-xl border border-subtle bg-background-800 p-6 shadow-lg backdrop-blur-sm"
          style="animation-delay: 0.2s"
        >
          <h3 class="mb-6 flex items-center text-xl font-bold text-content">
            <div class="mr-3 flex h-8 w-8 items-center justify-center rounded-lg bg-red-500/20">
              <i class="fas fa-tools text-red-400" />
            </div>
            {{ $t('admin.users.detail.quickActions') }}
          </h3>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div
              class="rounded-lg border border-subtle bg-background-700 p-4 transition-all hover:border-default hover:bg-background-600"
            >
              <cyberButton
                type="secondary"
                :loading="resetPasswordLoading"
                class="w-full justify-center border-0 shadow-sm"
                @click="resetPassword"
              >
                <i class="fas fa-key mr-2 text-yellow-400" />
                {{ $t('admin.users.detail.actions.resetPassword') }}
              </cyberButton>
              <p class="mt-2 text-center text-xs text-content-muted">{{ $t('admin.users.detail.actions.resetPasswordDesc') }}</p>
            </div>

            <div
              class="rounded-lg border border-subtle bg-background-700 p-4 transition-all hover:border-default hover:bg-background-600"
            >
              <cyberButton type="secondary" class="w-full justify-center border-0 shadow-sm" @click="showEmailDialog = true">
                <i class="fas fa-envelope mr-2 text-blue-400" />
                {{ $t('admin.users.detail.actions.sendEmail') }}
              </cyberButton>
              <p class="mt-2 text-center text-xs text-content-muted">{{ $t('admin.users.detail.actions.sendEmailDesc') }}</p>
            </div>

            <div
              class="rounded-lg border border-subtle bg-background-700 p-4 transition-all hover:border-default hover:bg-background-600"
            >
              <cyberButton
                :type="userDetail.status === 1 ? 'danger' : 'success'"
                :loading="toggleStatusLoading"
                class="w-full justify-center border-0 shadow-sm"
                @click="toggleStatus"
              >
                <i class="mr-2" :class="[userDetail.status === 1 ? 'fas fa-ban text-red-400' : 'fas fa-check text-green-400']" />
                {{
                  userDetail.status === 1
                    ? $t('admin.users.detail.actions.disableUser')
                    : $t('admin.users.detail.actions.enableUser')
                }}
              </cyberButton>
              <p class="mt-2 text-center text-xs text-content-muted">
                {{
                  userDetail.status === 1
                    ? $t('admin.users.detail.actions.disableUserDesc')
                    : $t('admin.users.detail.actions.enableUserDesc')
                }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="loading-container p-6 text-center">
        <cyberLoading :visible="true" :text="$t('admin.users.detail.loading')" />
      </div>
    </CyberDrawer>

    <SendEmailDialog v-model="showEmailDialog" :user="userDetail" @email-sent="handleEmailSent" />
  </div>
</template>

<style scoped lang="scss">
  .user-detail-drawer {
  }

  .user-detail-drawer :deep(.cyber-button) {
    transition: all 0.2s ease;
  }

  .user-detail-drawer :deep(.cyber-button:hover) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .user-detail-drawer .rounded-lg:hover {
    transition: all 0.2s ease;
  }

  .user-detail-drawer .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 200px;
    opacity: 0;
    animation: fade-in 0.3s ease 0.1s forwards;
  }
</style>
