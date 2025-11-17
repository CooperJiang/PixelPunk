import { ref, computed } from 'vue'
import { useToast } from '@/components/Toast/useToast'
import { useTexts } from '@/composables/useTexts'
import { API_BASE_URL } from '@/constants/env'
import type { ApiTestParams, ApiTestState, ApiResponse } from '../types'

export function useApiTest() {
  const toast = useToast()
  const { $t } = useTexts()

  const apiKey = ref('')
  const selectedFiles = ref<File[]>([])
  const params = ref<ApiTestParams>({
    folderId: '',
    filePath: '',
    accessLevel: 'private',
    optimize: false,
  })
  const response = ref<ApiResponse | null>(null)
  const isLoading = ref(false)
  const responseCopied = ref(false)
  const showValidation = ref(false)
  const isDragOver = ref(false)

  const currentDomain = computed(() => {
    const apiBaseUrl = API_BASE_URL
    if (apiBaseUrl.includes('/api')) {
      return apiBaseUrl.replace(/\/api.*$/, '')
    }
    if (apiBaseUrl.startsWith('/')) {
      return typeof window !== 'undefined' ? window.location.origin : 'http://localhost:9520'
    }
    return apiBaseUrl
  })

  const totalFileSize = computed(() => {
    return selectedFiles.value.reduce((total, file) => total + file.size, 0)
  })

  const responseStatusClass = computed(() => {
    if (!response.value) return ''
    const status = response.value.status
    if (status >= 200 && status < 300) return 'success'
    if (status >= 400 && status < 500) return 'error'
    if (status >= 500) return 'server-error'
    return ''
  })

  const uploadedImages = computed(() => {
    return response.value?.data?.data?.uploaded || []
  })

  const hasErrors = computed(() => {
    const data = response.value?.data?.data
    return !!(
      (data?.oversized_files && data.oversized_files.length > 0) ||
      (data?.unsupported_files && data.unsupported_files.length > 0) ||
      (data?.invalid_files && data.invalid_files.length > 0) ||
      (data?.upload_errors && data.upload_errors.length > 0)
    )
  })

  const oversizedFiles = computed(() => {
    return response.value?.data?.data?.oversized_files || []
  })

  const unsupportedFiles = computed(() => {
    return response.value?.data?.data?.unsupported_files || []
  })

  const invalidFiles = computed(() => {
    return response.value?.data?.data?.invalid_files || []
  })

  const uploadErrors = computed(() => {
    return response.value?.data?.data?.upload_errors || []
  })

  const sizeLimit = computed(() => {
    return response.value?.data?.data?.size_limit || ''
  })

  const sendRequest = async () => {
    if (!selectedFiles.value.length || !apiKey.value.trim()) {
      showValidation.value = true
      return
    }

    isLoading.value = true
    response.value = null
    responseCopied.value = false

    try {
      const formData = new FormData()

      selectedFiles.value.forEach((file) => {
        formData.append('files[]', file)
      })

      if (params.value.folderId) formData.append('folderId', params.value.folderId)
      if (params.value.filePath) formData.append('filePath', params.value.filePath)
      formData.append('access_level', params.value.accessLevel)
      formData.append('optimize', params.value.optimize ? 'true' : 'false')

      const startTime = Date.now()

      const res = await fetch(`${currentDomain.value}/api/v1/external/upload`, {
        method: 'POST',
        headers: {
          'x-pixelpunk-key': apiKey.value,
        },
        body: formData,
      })

      const endTime = Date.now()
      const data = await res.json()

      response.value = {
        status: res.status,
        statusText: res.statusText,
        data: data,
        duration: endTime - startTime,
      }
    } catch (error: any) {
      response.value = {
        status: 0,
        statusText: 'Network Error',
        data: { error: error.message },
        duration: 0,
      }
    } finally {
      isLoading.value = false
    }
  }

  const copyResponse = async () => {
    try {
      const jsonString = JSON.stringify(response.value!.data, null, 2)
      await navigator.clipboard.writeText(jsonString)
      responseCopied.value = true
      toast.success($t('docs.apiTester.toast.responseCopied'))
      setTimeout(() => {
        responseCopied.value = false
      }, 2000)
    } catch (error) {
      console.error('Copy failed:', error)
      try {
        const jsonString = JSON.stringify(response.value!.data, null, 2)
        const textArea = document.createElement('textarea')
        textArea.value = jsonString
        textArea.style.position = 'fixed'
        textArea.style.left = '-999999px'
        textArea.style.top = '-999999px'
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        const success = document.execCommand('copy')
        document.body.removeChild(textArea)

        if (success) {
          responseCopied.value = true
          toast.success($t('docs.apiTester.toast.responseCopied'))
          setTimeout(() => {
            responseCopied.value = false
          }, 2000)
        } else {
          throw new Error('execCommand failed')
        }
      } catch (fallbackError) {
        console.error('fallback复制也失败:', fallbackError)
        toast.error($t('docs.apiTester.toast.copyFailed'))
      }
    }
  }

  const copyImageUrl = async (url: string) => {
    try {
      await navigator.clipboard.writeText(url)
      toast.success($t('docs.apiTester.toast.imageLinkCopied'))
    } catch (error) {
      console.error('Copy failed:', error)
      try {
        const textArea = document.createElement('textarea')
        textArea.value = url
        textArea.style.position = 'fixed'
        textArea.style.left = '-999999px'
        textArea.style.top = '-999999px'
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        const success = document.execCommand('copy')
        document.body.removeChild(textArea)

        if (success) {
          toast.success($t('docs.apiTester.toast.imageLinkCopied'))
        } else {
          throw new Error('execCommand failed')
        }
      } catch (fallbackError) {
        console.error('fallback复制也失败:', fallbackError)
        toast.error($t('docs.apiTester.toast.copyLinkFailed'))
      }
    }
  }

  const openImageInNewTab = (url: string) => {
    try {
      const newWindow = window.open(url, '_blank')
      if (!newWindow) {
        toast.error($t('docs.apiTester.toast.unableToOpenWindow'))
      } else {
        toast.success($t('docs.apiTester.toast.openingNewWindow'))
      }
    } catch (error) {
      console.error('Failed to open new window:', error)
      toast.error($t('docs.apiTester.toast.failedToOpenWindow'))
    }
  }

  return {
    apiKey,
    selectedFiles,
    params,
    response,
    isLoading,
    responseCopied,
    showValidation,
    isDragOver,
    currentDomain,
    totalFileSize,
    responseStatusClass,
    uploadedImages,
    hasErrors,
    oversizedFiles,
    unsupportedFiles,
    invalidFiles,
    uploadErrors,
    sizeLimit,
    sendRequest,
    copyResponse,
    copyImageUrl,
    openImageInNewTab,
  }
}
