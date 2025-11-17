/**
 * 水印配置处理工具
 * 负责准备和转换水印配置

 */
import { convertWatermarkConfigForActualFile, getFileDimensions } from '@/utils/file/watermarkUtils'
import { type CanvasWatermarkConfig, generateWatermark } from '@/utils/file/canvasWatermarkGenerator'
import { globalOptions } from '../state'
import type { UploadItem } from '../types'

/* 准备转换后的水印配置 */
export async function prepareWatermarkConfig(item: UploadItem) {
  const baseConfig = { ...(globalOptions.value.watermarkConfig || {}) }
  if (baseConfig.generatedImage && !baseConfig.fileBase64) {
    try {
      baseConfig.fileBase64 = String(baseConfig.generatedImage).replace(/^data:image\/\w+;base64,/, '')
    } catch {}
  }

  const sanitizeByType = (conf: any) => {
    const c = { ...conf }
    if (c.type === 'text') {
    } else {
      c.type = 'image'
      delete c.text
      delete c.fontSize
      delete c.fontFamily
      delete c.fontColor
      delete c.fontWeight
    }
    return c
  }

  if (!globalOptions.value.watermarkEnabled || !baseConfig || !item.dimensions) {
    const sanitized = sanitizeByType({ ...baseConfig, type: baseConfig?.type || 'image' })
    return {
      watermark: globalOptions.value.watermarkEnabled ? JSON.stringify(sanitized) : undefined,
    }
  }

  let configToConvert = sanitizeByType({ ...baseConfig })

  if (configToConvert.fileBase64 && configToConvert.generateMode === 'frontend') {
    configToConvert.type = 'image'

    const targetImageWidth = item.dimensions.width
    const watermarkOriginalWidth = configToConvert.generatedWidth || 100 // 水印原始宽度
    const relativeSize = configToConvert.relativeSize || 0.15 // 默认占图片宽度的15%
    const calculatedScale = (targetImageWidth * relativeSize) / watermarkOriginalWidth
    configToConvert.scale = Math.max(0.1, Math.min(5.0, calculatedScale)) // 限制在合理范围内

    configToConvert.rotation = 0 // 不旋转，水印图片已经旋转好了
    configToConvert.shadow = false // 不添加阴影，水印图片已经包含阴影
    configToConvert.shadowBlur = 0
    configToConvert.shadowOffsetX = 0
    configToConvert.shadowOffsetY = 0
  } else if (configToConvert.type === 'text' && configToConvert.text && configToConvert.text.trim()) {
    try {
      const canvasConfig: CanvasWatermarkConfig = {
        ...configToConvert,
        targetWidth: item.dimensions.width,
        targetHeight: item.dimensions.height,
      }

      const result = await generateWatermark(canvasConfig)

      if (result.success && result.base64Data) {
        const pureBase64 = result.base64Data.replace(/^data:image\/\w+;base64,/, '')

        const targetImageWidth = item.dimensions.width
        const watermarkOriginalWidth = result.width
        const relativeSize = configToConvert.relativeSize || 0.15
        const calculatedScale = (targetImageWidth * relativeSize) / watermarkOriginalWidth
        const finalScale = Math.max(0.1, Math.min(5.0, calculatedScale))

        configToConvert = {
          ...configToConvert,
          type: 'image', // 转换为文件类型
          generateMode: 'frontend',
          generatedImage: result.base64Data,
          generatedWidth: result.width,
          generatedHeight: result.height,
          fileBase64: pureBase64, // 设置为生成的文件数据
          imageUrl: undefined, // 清除默认URL
          imageFile: undefined, // 清除文件
          scale: finalScale, // 使用动态计算的缩放
          rotation: 0, // 水印已旋转
          shadow: false, // 水印已包含阴影
          shadowBlur: 0,
          shadowOffsetX: 0,
          shadowOffsetY: 0,
        }
      } else {
        console.error('[水印配置] 文本水印生成失败：生成结果无效')
        return {
          watermark: undefined,
        }
      }
    } catch (error) {
      console.error('[水印配置] 文本水印生成异常:', error)
      return {
        watermark: undefined,
      }
    }
  } else if (configToConvert.type !== 'text') {
    configToConvert.type = 'image'
    if (!configToConvert.fileBase64 && configToConvert.generatedFile) {
      try {
        configToConvert.fileBase64 = String(configToConvert.generatedFile).replace(/^data:image\/\w+;base64,/, '')
      } catch (error) {
        console.warn('[水印配置] generatedFile格式转换失败:', error)
      }
    }
  }

  const convertedConfig = convertWatermarkConfigForActualFile(configToConvert, item.dimensions.width, item.dimensions.height)

  const finalConfig = sanitizeByType(convertedConfig)
  if (finalConfig.type === 'image' && !finalConfig.fileBase64 && finalConfig.generatedFile) {
    try {
      finalConfig.fileBase64 = String(finalConfig.generatedFile).replace(/^data:image\/\w+;base64,/, '')
    } catch (error) {
      console.warn('[水印配置] 最终base64转换失败:', error)
    }
  }

  const backendConfig = {
    enabled: true,
    type: finalConfig.type, // 前后端统一使用 'image' 或 'text'
    fileBase64: finalConfig.fileBase64 || '',
    position: finalConfig.position || 'bottom-right',
    offsetX: finalConfig.offsetX ?? 20,
    offsetY: finalConfig.offsetY ?? 20,
    offsetUnit: 'px' as const,
    opacity: finalConfig.opacity ?? 0.7,
    scale: finalConfig.scale ?? 1.0, // 水印在目标图片上的缩放比例
  }

  return {
    watermark: JSON.stringify(backendConfig),
  }
}

export async function createUploadItemWithDimensions(
  file: File,
  itemId: string,
  type: 'regular' | 'chunked' = 'regular'
): Promise<any> {
  const baseItem: any = {
    id: itemId,
    file,
    status: 'pending',
    progress: 0,
    speed: 0,
    remainingTime: 0,
    type,
  }

  try {
    if (file.type.startsWith('image/')) {
      const dimensions = await getFileDimensions(file)
      baseItem.dimensions = dimensions
    }
  } catch (error) {
    console.warn(`⚠️ [UploadStore] 获取文件尺寸失败: ${file.name}`, error)
  }

  return baseItem
}
