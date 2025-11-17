/**
 * 分享功能相关工具函数
 */

/**
 * 复制分享链接到剪贴板
 * @param shareUrl 分享链接URL
 * @param onSuccess 成功回调
 * @param onError 失败回调
 */
export const copyShareLink = (shareUrl: string, onSuccess?: () => void, onError?: () => void): void => {
  if (navigator.clipboard) {
    navigator.clipboard
      .writeText(shareUrl)
      .then(() => {
        if (onSuccess) {
          onSuccess()
        }
      })
      .catch(() => {
        fallbackCopyTextToClipboard(shareUrl, onSuccess, onError)
      })
  } else {
    fallbackCopyTextToClipboard(shareUrl, onSuccess, onError)
  }
}

const fallbackCopyTextToClipboard = (text: string, onSuccess?: () => void, onError?: () => void): void => {
  const textArea = document.createElement('textarea')
  textArea.value = text

  textArea.style.position = 'fixed'
  textArea.style.left = '-999999px'
  textArea.style.top = '-999999px'
  document.body.appendChild(textArea)

  textArea.focus()
  textArea.select()

  try {
    const successful = document.execCommand('copy')
    if (successful) {
      if (onSuccess) {
        onSuccess()
      }
    } else if (onError) {
      onError()
    }
  } catch {
    if (onError) {
      onError()
    }
  }

  document.body.removeChild(textArea)
}

export const canUseShareAPI = (): boolean => typeof navigator !== 'undefined' && Boolean(navigator.share)

export const shareToSystem = async (
  shareData: { title: string; text: string; url: string },
  onSuccess?: () => void,
  onError?: () => void
): Promise<void> => {
  if (!canUseShareAPI()) {
    if (onError) {
      onError()
    }
    return
  }

  try {
    await navigator.share(shareData)
    if (onSuccess) {
      onSuccess()
    }
  } catch (error: unknown) {
    if (error.name !== 'AbortError') {
      if (onError) {
        onError()
      }
    }
  }
}

export const shareToWeChat = (shareUrl: string, _shareTitle: string): void => {
  const encodedUrl = encodeURIComponent(shareUrl)

  window.open(`https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=${encodedUrl}`, '_blank')
}

export const shareToWeibo = (shareUrl: string, shareTitle: string, shareDesc: string = ''): void => {
  const encodedUrl = encodeURIComponent(shareUrl)
  const encodedTitle = encodeURIComponent(shareTitle)
  const encodedDesc = encodeURIComponent(shareDesc)

  window.open(`http://service.weibo.com/share/share.php?url=${encodedUrl}&title=${encodedTitle}%20${encodedDesc}`, '_blank')
}

export const shareToQQ = (shareUrl: string, shareTitle: string, shareDesc: string = ''): void => {
  const encodedUrl = encodeURIComponent(shareUrl)
  const encodedTitle = encodeURIComponent(shareTitle)
  const encodedDesc = encodeURIComponent(shareDesc)

  window.open(
    `https://connect.qq.com/widget/shareqq/index.html?url=${encodedUrl}&title=${encodedTitle}&summary=${encodedDesc}`,
    '_blank'
  )
}
