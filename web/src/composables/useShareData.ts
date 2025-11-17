import { ref, watch } from 'vue'
import { useRouter, type RouteLocationNormalizedLoaded } from 'vue-router'
import { shareApi } from '@/api'
import { StorageUtil } from '@/utils/storage/storage'
import { processFileUrls } from '@/utils/file/fileUtils'
import { STORAGE_KEYS } from '@/constants'
import { useTexts } from '@/composables/useTexts'

/**
 * 分享数据加载功能
 * 处理分享数据的加载、密码验证等
 */
export function useShareData(
  route: RouteLocationNormalizedLoaded,
  onError: (title: string, message: string, type?: string) => void
) {
  const router = useRouter()
  const { $t } = useTexts()

  const initialLoading = ref(true)
  const folderLoading = ref(false)
  const error = ref(false)
  const errorTitle = ref('')
  const errorMessage = ref('')
  const errorType = ref('error')

  const shareData = ref<any>(null)
  const shareKey = ref('')
  const accessToken = ref('')
  const shareName = ref('')

  const passwordDialogVisible = ref(false)
  const visitorDialogVisible = ref(false)

  const folders = ref<any[]>([])
  const images = ref<any[]>([])
  const rootImages = ref<any[]>([])

  const isVisitorDialogDismissed = () => {
    if (!shareKey.value) {
      return false
    }
    const storageKey = `${STORAGE_KEYS.VISITOR_DIALOG_DISMISSED}${shareKey.value}`
    return StorageUtil.has(storageKey)
  }

  const handleVisitorDismiss = () => {
    visitorDialogVisible.value = false

    if (shareKey.value) {
      const storageKey = `${STORAGE_KEYS.VISITOR_DIALOG_DISMISSED}${shareKey.value}`
      StorageUtil.set(storageKey, true, 24) // 保存24小时
    }
  }

  const loadShareData = async (currentFolderId?: string, isInitial = false) => {
    if (!shareKey.value) {
      onError($t('components.share.errors.invalidLink'), $t('components.share.errors.invalidLinkDesc'), 'not-found')
      return
    }

    if (isInitial) {
      initialLoading.value = true
    } else if (currentFolderId) {
      folderLoading.value = true
    }

    error.value = false

    try {
      const params: any = {}
      if (accessToken.value) {
        params.access_token = accessToken.value
      }

      if (currentFolderId) {
        params.folder_id = currentFolderId
      }

      try {
        const result = await shareApi.getPublicShare(shareKey.value, params)

        if (result.success) {
          const { data } = result

          if (data.require_password) {
            shareName.value = data.name || $t('components.share.defaults.shareName')
            passwordDialogVisible.value = true
            initialLoading.value = false
            folderLoading.value = false
            return
          }

          shareData.value = data

          folders.value = data.folders || []

          const processedFiles = (data.files || []).map((file: Event) => {
            const processedFile = processFileUrls(file, accessToken.value)

            if (processedFile.ai_info?.description) {
              processedFile.description = processedFile.ai_info.description
            }

            return processedFile
          })

          images.value = processedFiles

          if (isInitial || !currentFolderId) {
            rootImages.value = [...processedFiles]
          }

          if (data.share.collect_visitor_info && !isVisitorDialogDismissed()) {
            visitorDialogVisible.value = true
          }

          updatePageMetadata()
        } else {
          onError($t('components.share.errors.loadFailed'), result.message || $t('components.share.errors.loadFailedDesc'))
        }
      } catch (apiError: any) {
        if (apiError.code === 403) {
          if (apiError.message && apiError.message.includes($t('composables.useShareData.backendErrors.maxViewsKeyword'))) {
            onError($t('components.share.errors.maxViewsReached'), $t('components.share.errors.maxViewsReachedDesc'), 'max-views')
          } else if (apiError.message && apiError.message.includes($t('composables.useShareData.backendErrors.expiredKeyword'))) {
            onError($t('components.share.errors.shareExpired'), $t('components.share.errors.shareExpiredDesc'), 'expired')
          } else {
            onError(
              $t('components.share.errors.accessRestricted'),
              apiError.message || $t('components.share.errors.accessRestrictedDesc'),
              'unauthorized'
            )
          }
        } else {
          throw apiError // 重新抛出错误，由外部catch处理
        }
      }
    } catch (error: any) {
      const errorMessage = error?.message || error?.msg || String(error)

      if (error?.code === 404 || errorMessage.includes('404') || errorMessage.includes('Not Found')) {
        onError($t('components.share.errors.shareNotFound'), $t('components.share.errors.shareNotFoundDesc'), 'not-found')
      } else if (error?.code === 403 || errorMessage.includes('403')) {
        onError(
          $t('components.share.errors.accessRestricted'),
          $t('components.share.errors.accessRestrictedDesc'),
          'unauthorized'
        )
      } else if (error?.code === 109 || errorMessage.includes($t('composables.useShareData.backendErrors.invalidKeyword'))) {
        onError($t('components.share.errors.shareInvalid'), $t('components.share.errors.shareInvalidDesc'), 'expired')
      } else {
        onError($t('components.share.errors.loadFailed'), errorMessage || $t('components.share.errors.loadFailedDesc'))
      }
    } finally {
      initialLoading.value = false
      folderLoading.value = false
    }
  }

  const handlePasswordVerify = async (
    password: string,
    onSuccess?: (message: string) => void,
    onError?: (message: string) => void
  ) => {
    try {
      const result = await shareApi.verifySharePassword(shareKey.value, password)

      if (result.success) {
        const { data } = result
        passwordDialogVisible.value = false

        if (data && data.access_token) {
          accessToken.value = data.access_token

          router.replace({
            path: route.path,
            query: { access_token: accessToken.value },
          })

          if (onSuccess) {
            onSuccess($t('components.share.errors.passwordVerifySuccess'))
          }

          loadShareData(undefined, true)
        } else {
          onError?.($t('components.share.errors.passwordVerifyNoToken'))
        }
      } else {
        onError?.(result.message || $t('components.share.errors.passwordVerifyFailed'))
      }
    } catch (error: unknown) {
      onError?.(error.value?.message || $t('components.share.errors.passwordVerifyFailedRetry'))
    }
  }

  const handleVisitorSubmit = async (
    visitorInfo: { name: string; email?: string },
    onSuccess?: (message: string) => void,
    onError?: (message: string) => void
  ) => {
    try {
      const result = await shareApi.submitVisitorInfo(shareKey.value, visitorInfo)

      if (result.success) {
        visitorDialogVisible.value = false

        if (shareKey.value) {
          const storageKey = `${STORAGE_KEYS.VISITOR_DIALOG_DISMISSED}${shareKey.value}`
          StorageUtil.set(storageKey, true, 24) // 保存24小时
        }

        if (onSuccess) {
          onSuccess($t('components.share.errors.visitorSubmitSuccess'))
        }
      } else if (onError) {
        onError(result.message || $t('components.share.errors.visitorSubmitFailed'))
      }
    } catch (error: unknown) {
      if (onError) {
        onError(error.value?.message || $t('components.share.errors.visitorSubmitFailed'))
      }
    }
  }

  const updatePageMetadata = () => {
    if (!shareData.value) {
      return
    }

    const { share } = shareData.value
    const title = share.name || $t('components.share.defaults.shareTitle')
    const description = share.description || $t('components.share.defaults.shareDescription', { count: images.value.length })

    let previewImageUrl = ''
    if (images.value && images.value.length > 0) {
      const firstImage = images.value[0]
      previewImageUrl = firstImage.full_url || firstImage.url || ''
    }

    const currentUrl = window.location.href

    document.title = title

    const metaTags = [
      { name: 'description', content: description },
      { property: 'og:title', content: title },
      { property: 'og:description', content: description },
      { property: 'og:type', content: 'website' },
      { property: 'og:url', content: currentUrl },
      { property: 'og:image', content: previewImageUrl },
      { property: 'twitter:card', content: 'summary_large_image' },
      { property: 'twitter:title', content: title },
      { property: 'twitter:description', content: description },
      { property: 'twitter:image', content: previewImageUrl },
    ]

    metaTags.forEach((tag) => {
      const selector = tag.name ? `meta[name="${tag.name}"]` : `meta[property="${tag.property}"]`

      let element = document.querySelector(selector) as HTMLMetaElement

      if (element) {
        element.content = tag.content
      } else {
        element = document.createElement('meta')

        if (tag.name) {
          element.name = tag.name
        } else if (tag.property) {
          element.setAttribute('property', tag.property)
        }

        element.content = tag.content
        document.head.appendChild(element)
      }
    })
  }

  watch(
    () => route.params.key,
    (newKey) => {
      if (newKey && typeof newKey === 'string') {
        shareKey.value = newKey

        const queryAccessToken = route.query.access_token
        if (queryAccessToken && typeof queryAccessToken === 'string') {
          accessToken.value = queryAccessToken
        }

        loadShareData(undefined, true)
      }
    },
    { immediate: true }
  )

  return {
    initialLoading,
    folderLoading,
    error,
    errorTitle,
    errorMessage,
    errorType,
    shareData,
    shareKey,
    accessToken,
    shareName,
    passwordDialogVisible,
    visitorDialogVisible,
    folders,
    images,
    rootImages,
    loadShareData,
    handlePasswordVerify,
    handleVisitorSubmit,
    handleVisitorDismiss,
    updatePageMetadata,
  }
}
