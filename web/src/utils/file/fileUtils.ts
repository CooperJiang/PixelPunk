/**
 * 文件处理相关工具函数
 */

/**
 * 处理文件URL，添加access_token，并处理AI信息
 * @param image 文件对象
 * @param accessToken 访问令牌
 * @returns 处理后的文件对象
 */
export const processFileUrls = (image: any, accessToken?: string): any => {
  const processedImage = { ...image }

  if (processedImage.ai_info && processedImage.ai_info.description && !processedImage.description) {
    processedImage.description = processedImage.ai_info.description
  }

  if (!accessToken) {
    return processedImage
  }

  const urlFields = ['url', 'full_url', 'thumb_url', 'full_thumb_url']

  urlFields.forEach((field) => {
    if (processedImage[field]) {
      const url = processedImage[field]
      const hasQuery = url.includes('?')
      processedImage[field] = `${url}${hasQuery ? '&' : '?'}access_token=${accessToken}`
    }
  })

  return processedImage
}

export const getFileDownloadUrl = (image: any, accessToken?: string): string => {
  const url = image.full_url || image.url || ''
  if (!url) {
    return ''
  }

  if (accessToken) {
    const hasQuery = url.includes('?')
    return `${url}${hasQuery ? '&' : '?'}access_token=${accessToken}`
  }

  return url
}

export const downloadSingleImage = (url: string, filename: string): void => {
  if (!url) {
    return
  }

  fetch(url)
    .then((response) => response.blob())
    .then((blob) => {
      const blobUrl = URL.createObjectURL(blob)

      const a = document.createElement('a')
      a.href = blobUrl
      a.download = filename || 'image'
      a.style.display = 'none'

      document.body.appendChild(a)
      a.click()

      setTimeout(() => {
        document.body.removeChild(a)
        URL.revokeObjectURL(blobUrl)
      }, 100)
    })
    .catch((error) => {
      console.error('下载文件失败:', error)

      try {
        const a = document.createElement('a')
        a.href = url
        a.download = filename || 'image'
        a.target = '_blank' // 确保在新窗口打开，避免导航离开当前页面
        a.rel = 'noopener noreferrer'
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
      } catch (e) {
        console.error('回退下载也失败:', e)
      }
    })
}

export const processImageUrls = processFileUrls

export const getImageDownloadUrl = getFileDownloadUrl

export const estimateImageHeight = (item: any): number => {
  if (item.width && item.height) {
    const aspectRatio = item.width / item.height
    const baseColumnWidth = 200
    const scaledHeight = baseColumnWidth / aspectRatio

    return Math.max(150, Math.min(scaledHeight, 400))
  }

  if (item.format) {
    const format = item.format.toLowerCase()
    if (['jpg', 'jpeg', 'png'].includes(format)) {
      return 220 // 普通文件
    } else if (['gif'].includes(format)) {
      return 250 // GIF 可能更高一些
    }
  }

  return 200
}
