<script setup lang="ts">
  import { reactive, ref } from 'vue'
  import { useAuthStore } from '@/store/auth'
  import { useToast } from '@/components/Toast/useToast'
  import { uploadAvatar } from '@/api/file'
  import { getUserProfile, updateUserProfile } from '@/api/user'
  import type { UpdateProfileDTO } from '@/api/types/index'
  import { useTexts } from '@/composables/useTexts'

  const authStore = useAuthStore()
  const toast = useToast()
  const { $t } = useTexts()

  const avatarInput = ref<HTMLInputElement | null>(null)
  const avatarUrl = ref<string>(authStore.userInfo?.avatarFullPath || authStore.userInfo?.avatar || '')
  const avatarPath = ref<string>(authStore.userInfo?.avatar || '')
  const isUploading = ref<boolean>(false)
  const isLoading = ref<boolean>(false)
  const isSaving = ref<boolean>(false)

  const profile = reactive({
    username: authStore.user?.username || '',
    email: '',
    nickname: '',
    website: '',
    bio: '',
    avatar: avatarPath.value || '',
  })

  const fetchUserProfile = async () => {
    try {
      isLoading.value = true
      const result = await getUserProfile()

      if (result.success) {
        const userData = result.data

        profile.username = userData.username || ''
        profile.email = userData.email || ''
        profile.bio = userData.bio || ''
        profile.website = userData.website || ''

        if (userData.avatar) {
          avatarPath.value = userData.avatar
          profile.avatar = userData.avatar
        }

        if (userData.avatarFullPath) {
          avatarUrl.value = userData.avatarFullPath
        } else {
          avatarUrl.value = userData.avatar || ''
        }

        authStore.setUserInfo(userData)
      }
    } catch (_error) {
    } finally {
      isLoading.value = false
    }
  }

  const triggerAvatarUpload = () => {
    avatarInput.value?.click()
  }

  const handleAvatarUpload = async (event: Event) => {
    const target = event.target as HTMLInputElement
    const { files } = target

    if (!files || files.length === 0) {
      return
    }

    const file = files[0]

    if (file.size > 2 * 1024 * 1024) {
      toast.error($t('settings.profile.avatar.sizeLimit'))
      return
    }

    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif']
    if (!allowedTypes.includes(file.type)) {
      toast.error($t('settings.profile.avatar.invalidType'))
      return
    }

    isUploading.value = true

    try {
      const result = await uploadAvatar(file)

      if (result.success) {
        const { data } = result

        avatarUrl.value = data.full_url || data.url

        avatarPath.value = data.avatar_url || data.url
        profile.avatar = data.avatar_url || data.url

        authStore.updateUserAvatar(data.avatar_url || data.url, data.full_url || data.url)

        toast.success($t('settings.profile.avatar.uploadSuccess'))
      }
    } catch (_error) {
    } finally {
      isUploading.value = false
      if (avatarInput.value) {
        avatarInput.value.value = ''
      }
    }
  }

  const saveProfile = async () => {
    try {
      isSaving.value = true

      const updateData: UpdateProfileDTO = {
        username: profile.username,
        bio: profile.bio,
        website: profile.website,
        avatar: avatarPath.value,
      }

      const result = await updateUserProfile(updateData)

      if (result.success) {
        if (result.data) {
          authStore.setUserInfo(result.data)
        }

        toast.success($t('settings.profile.messages.saveSuccess'))
      }
    } catch (_error) {
    } finally {
      isSaving.value = false
    }
  }

  defineExpose({
    fetchUserProfile,
  })
</script>

<template>
  <div>
    <CyberLoading v-if="isLoading" :visible="isLoading" />

    <div v-else class="mt-5 grid gap-6 md:grid-cols-[200px_1fr]">
      <div class="profile-avatar-card flex flex-col items-center border p-5 text-center">
        <CyberUserAvatar
          :avatar-url="avatarUrl"
          :username="profile.username"
          size="md"
          :custom-styles="{ width: '80px', height: '80px', fontSize: '1.75rem' }"
        />
        <input ref="avatarInput" type="file" class="hidden" accept=".jpg,.jpeg,.png,.webp,.gif" @change="handleAvatarUpload" />
        <CyberButton type="secondary" icon="upload" :loading="isUploading" class="mt-4 w-full" @click="triggerAvatarUpload">
          {{ isUploading ? $t('status.uploading') : $t('settings.profile.avatar.change') }}
        </CyberButton>
        <p class="mt-3 text-xs text-content-muted">{{ $t('settings.profile.avatar.hint') }}</p>
      </div>

      <div class="profile-form-card border p-5">
        <form class="space-y-4" @submit.prevent="saveProfile">
          <div class="space-y-2">
            <label class="block text-sm font-medium text-content">
              {{ $t('settings.profile.labels.email') }}
            </label>
            <CyberInput
              v-model="profile.email"
              type="email"
              disabled
              :placeholder="$t('settings.profile.placeholders.email')"
              prefix-icon="envelope"
            />
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-content">
              {{ $t('settings.profile.labels.username') }}
            </label>
            <CyberInput
              v-model="profile.username"
              type="text"
              :placeholder="$t('settings.profile.placeholders.username')"
              prefix-icon="user"
            />
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-content">
              {{ $t('settings.profile.labels.website') }}
            </label>
            <CyberInput
              v-model="profile.website"
              type="url"
              :placeholder="$t('settings.profile.placeholders.website')"
              prefix-icon="globe"
            />
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-content">
              {{ $t('settings.profile.labels.bio') }}
            </label>
            <CyberInput v-model="profile.bio" type="textarea" :placeholder="$t('settings.profile.placeholders.bio')" :rows="3" />
          </div>

          <div class="flex justify-end">
            <CyberButton type="primary" icon="save" :loading="isSaving">
              {{ isSaving ? $t('settings.profile.actions.saving') : $t('settings.profile.actions.save') }}
            </CyberButton>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .profile-avatar-card,
  .profile-form-card {
    border-radius: var(--radius-sm);
    background: linear-gradient(135deg, rgba(var(--color-background-700-rgb), 0.8), rgba(var(--color-background-800-rgb), 0.9));
    border-color: var(--color-border-subtle);
    box-shadow:
      0 4px 16px rgba(0, 0, 0, 0.08),
      inset 0 1px 0 rgba(255, 255, 255, 0.05);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .profile-avatar-card:hover,
  .profile-form-card:hover {
    border-color: var(--color-border-default);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.03);
    transform: translateY(-1px);
  }
</style>
